package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync/atomic"
	"time"
	"yanglu/config"
	"yanglu/service/model"

	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type HostInfoService struct {
	uid int
}

func NewHostInfoService() *HostInfoService {
	return &HostInfoService{}
}

func NewHostInfoServiceWithUid(uid int) *HostInfoService {
	return &HostInfoService{uid: uid}
}

func (hs *HostInfoService) GetHostNum() int {
	list, _ := model.NewHostInfo().ListAll()
	if list != nil {
		return len(list)
	}
	return 0
}

func (hs *HostInfoService) Add(ip string, port int, sshUser string, sshPasswd string, department string, businessName string) (*model.HostInfo, error) {

	hsDao, err := model.NewHostInfo().GetHostInfoByIp(ip)
	if err != nil {
		return nil, err
	}
	if hsDao.Id != 0 {
		return nil, errors.New("已经添加了该机器")
	}
	// 针对云端版，需要知道是谁添加的uid
	hsDao.Uid = hs.uid
	hsDao.Ip = ip
	hsDao.Port = port
	hsDao.SshUser = sshUser
	hsDao.SshPasswd = sshPasswd
	hsDao.Department = department
	hsDao.BusinessName = businessName
	// 检查有效性
	err = hs.CheckPass(hsDao)
	if err != nil {
		return nil, err
	}

	if config.LicenseInfoConf.NodeMax <= hs.GetHostNum() {
		return nil, fmt.Errorf("机器数量超过了最大限制%d台", config.LicenseInfoConf.NodeMax)
	}

	err = hsDao.Create()
	if err != nil {
		logrus.Error("Add err ", err)
		return nil, err
	}
	go func() {
		hs.Prepare([]*model.HostInfo{hsDao})
	}()
	return hsDao, nil
}

func (hs *HostInfoService) BatchAdd(list []*model.HostInfo) error {
	if len(list) == 0 {
		return errors.New("空数据")
	}
	err := hs.BatchCheck(list)
	if err != nil {
		return err
	}
	if config.LicenseInfoConf.NodeMax < hs.GetHostNum()+len(list) {
		return fmt.Errorf("机器数量超过了最大限制%d台", config.LicenseInfoConf.NodeMax)
	}
	for k, _ := range list {
		list[k].Uid = hs.uid
	}
	err = model.NewHostInfo().BatchCreate(list)
	if err != nil {
		logrus.Error("BatchAdd err ", err)
		return err
	}
	go func() {
		hs.Prepare(list)
	}()
	return nil
}

func (hs *HostInfoService) CheckPass(host *model.HostInfo) error {

	hostStr := fmt.Sprintf("%s:%d", host.Ip, host.Port)
	client, err := ssh.Dial("tcp", hostStr, &ssh.ClientConfig{
		User:            host.SshUser,
		Auth:            []ssh.AuthMethod{ssh.Password(host.SshPasswd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		logrus.Error("CheckPassB err ", err)
		return errors.New("请检查账号密码")
	}
	defer client.Close()
	return nil
}

func (hs *HostInfoService) BatchCheck(list []*model.HostInfo) error {
	res := make([]func() error, len(list))
	num := new(int32)
	*num = int32(len(list))
	for i := 0; i < len(list); i++ {
		res[i] = hs.Future(hs.CheckPass, list[i], num)
	}
	t := time.Now().Unix()
	for atomic.LoadInt32(num) != 0 {
		time.Sleep(time.Microsecond * 200)
		now := time.Now().Unix()
		if now-t > 4 {
			return errors.New("添加失败，请稍后重试")
		}
	}
	msg := ""
	for i := 0; i < len(list); i++ {
		err := res[i]()
		if err == nil {
			continue
		}
		msg = msg + fmt.Sprintf("%s机器添加错误，请核对;", list[i].Ip)
	}
	if msg != "" {
		return errors.New(msg)
	}
	return nil
}

func (hs *HostInfoService) Future(f func(host *model.HostInfo) error, host *model.HostInfo, num *int32) func() error {
	var err error
	c := make(chan struct{}, 1)
	go func(count *int32) {
		defer func() {
			close(c)
			atomic.AddInt32(num, -1)
		}()
		err = f(host)
	}(num)
	return func() error {
		<-c
		return err
	}
}

func (hs *HostInfoService) UpdateDepartment(ip string, department string) error {
	hsDao, err := model.NewHostInfo().GetHostInfoByIp(ip)
	if err != nil {
		logrus.Error("Department err ", err)
		return err
	}
	if config.IsCloud() && hsDao.Uid != hs.uid {
		return errors.New("你不是该台主机管理员")
	}
	hsDao.Department = department
	hsDao.UpdateTime = time.Now().Unix()
	m := map[string]interface{}{
		"department":  department,
		"update_time": hsDao.UpdateTime,
	}
	err = hsDao.Updates(m)
	if err != nil {
		logrus.Error("Department err ", err)
		return err
	}
	return nil
}

func (h *HostInfoService) GetSystemInfo(host *model.HostInfo) (string, error) {
	hostStr := fmt.Sprintf("%s:%d", host.Ip, host.Port)
	client, err := ssh.Dial("tcp", hostStr, &ssh.ClientConfig{
		User:            host.SshUser,
		Auth:            []ssh.AuthMethod{ssh.Password(host.SshPasswd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		logrus.Error("GetSystemInfo err ", err)
		return "", errors.New("请检查账号密码")
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		logrus.Error("GetSystemInfo err ", err)
		return "", err
	}
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	if err = session.Run("cat /etc/issue"); err != nil {
		logrus.Error("GetSystemInfo err ", err)
		return "", err
	}
	ss := strings.Split(strings.TrimSpace(b.String()), " ")
	if len(ss) > 0 {
		return ss[0], nil
	}
	return "", errors.New("找不到系统信息")
}

func (h *HostInfoService) Prepare(list []*model.HostInfo) error {
	if len(list) == 0 {
		return nil
	}
	var err error = nil
	for k, v := range list {
		if v.Ip == "" {
			continue
		}
		list[k].SystemOs, err = h.GetSystemInfo(list[k])
		if err != nil {
			logrus.Error("UpdateSystemInfo err ", err)
			continue
		}
		err = list[k].Updates(map[string]interface{}{
			"system_os": list[k].SystemOs,
		})
		if err != nil {
			logrus.Error("UpdateSystemInfo err ", err)
			continue
		}
		// 拷贝文件
		// h.CpFile(list[k], "trivy_dir.sh", "/var")
		// h.Cmd(list[k], "bash /var/trivy_dir.sh")
		cmd := `
		cd /var;path="/var/trivy";if [ ! -d "$path" ];then     mkdir "$path";     echo "ok"; else     echo "file_exists"; fi
		`
		h.Cmd(list[k], cmd)
		h.CpFileBySftp(list[k], "trivy_0.16.0_Linux-64bit.tar.gz", "/var/trivy")
		h.Cmd(list[k], "cd /var/trivy;tar -xzvf trivy_0.16.0_Linux-64bit.tar.gz")

		cmd = `
		cd /root/.cache;path="/root/.cache/trivy/db";if [ ! -d "$path" ];then     mkdir -p "$path";     echo "ok"; else     echo "file_exists"; fi
		`

		h.Cmd(list[k], cmd)
		h.CpFileBySftp(list[k], "trivy-offline.db.tgz", "/root/.cache/trivy/db")
		h.Cmd(list[k], "cd /root/.cache/trivy/db;tar zxvf trivy-offline.db.tgz")
	}
	return nil
}

func (h *HostInfoService) GetClient(host *model.HostInfo) (*ssh.Client, error) {
	hostStr := fmt.Sprintf("%s:%d", host.Ip, host.Port)
	client, err := ssh.Dial("tcp", hostStr, &ssh.ClientConfig{
		User:            host.SshUser,
		Auth:            []ssh.AuthMethod{ssh.Password(host.SshPasswd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		logrus.Error("GetClient err ", err)
		return nil, errors.New("请检查账号密码")
	}
	return client, nil
}

func (h *HostInfoService) Cmd(host *model.HostInfo, cmd string) (string, error) {
	client, err := h.GetClient(host)
	if err != nil {
		return "", err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		logrus.Error("cmd err ", err)
		return "", err
	}
	defer session.Close()
	result, err := session.Output(cmd)
	if err != nil {
		logrus.Error("Cmd = ", cmd)
	}
	return string(result), err
}

func (h *HostInfoService) CpFile(host *model.HostInfo, fileName string, filePath string) error {

	client, err := h.GetClient(host)
	if err != nil {
		return err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		logrus.Error("CpTrivyFile err ", err)
		return err
	}
	defer session.Close()

	go func() {
		buf := make([]byte, 1024)
		w, _ := session.StdinPipe()
		defer w.Close()
		path := fmt.Sprintf("../software_package/%s", fileName)
		file, err := os.Open(path)
		if err != nil {
			logrus.Error("CpTrivyFile err ", err)
			return
		}
		info, _ := file.Stat()
		fmt.Fprintln(w, "C0644", info.Size(), fileName)
		for {
			n, err := file.Read(buf)
			fmt.Fprint(w, string(buf[:n]))
			if err != nil {
				if err == io.EOF {
					return
				} else {
					logrus.Panic(err)
				}
			}
		}
	}()
	if err := session.Run(fmt.Sprintf("/usr/bin/scp -qrt %s", filePath)); err != nil {
		if err != nil {
			if err.Error() != "Process exited with: 1. Reason was:  ()" {
				logrus.Error("CpTrivyFile err ", err, " file = ", fileName)
			}
			return err
		}
		// 解压文案
	}
	return nil
}

func (h *HostInfoService) GetResult(host *model.HostInfo) (string, error) {

	client, err := h.GetClient(host)
	if err != nil {
		logrus.Error("GetResult err = ", err)
		return "", err
	}
	defer client.Close()
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		logrus.Error("GetResult err ", err)
		return "", err
	}
	path := "/var/trivy/results.json"
	srcFile, err := sftpClient.Open(path)
	if err != nil {
		logrus.Error("GetResult err ", err)
		return "", err
	}
	defer srcFile.Close()
	bytes, err := ioutil.ReadAll(srcFile)
	if err != nil {
		logrus.Error("GetResult err ", err)
		return "", err
	}
	return string(bytes), nil
}

func (h *HostInfoService) CpFileBySftp(host *model.HostInfo, sourceFile string, descPath string) error {
	client, err := h.GetClient(host)
	if err != nil {
		logrus.Error("CpFileBySftp err = ", err)
		return err
	}
	defer client.Close()
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		logrus.Error("CpFileBySftp err ", err)
		return err
	}

	localFile := "../software_package/" + sourceFile
	srcFile, err := os.Open(localFile)
	if err != nil {
		logrus.Error("CpFileBySftp ", err)
		return err
	}
	defer srcFile.Close()
	dstFile, err := sftpClient.Create(path.Join(descPath, sourceFile))
	if err != nil {
		logrus.Error("CpFileBySftp ", err)
		return err
	}
	defer dstFile.Close()

	// buf := make([]byte, 1024*1024)

	// for {
	// 	n, _ := srcFile.Read(buf)
	// 	logrus.Debug("n = ", n)
	// 	if n == 0 {
	// 		break
	// 	}
	// 	dstFile.Write(buf)
	// }
	_, err = io.Copy(dstFile, srcFile)

	if err != nil {
		logrus.Error(" err = ", err)
	}

	return nil
}

func (hc *HostInfoService) ListAll() (interface{}, error) {
	hosts, err := model.NewHostInfo().ListAll()
	if err != nil {
		logrus.Error("ListAll err ", err)
		return nil, err
	}
	if len(hosts) == 0 {
		return nil, errors.New("暂未有任何机器")
	}
	type Item struct {
		*model.HostInfo
		CheckStatus int `json:"check_status"`
	}
	list := make([]Item, len(hosts))
	ips := make([]string, len(hosts))
	for k, v := range hosts {
		ips[k] = v.Ip
	}
	checkStatus, err := model.NewVulnerabilityLog().HostCheckStatus(ips)
	if err != nil {
		logrus.Error("ListAll err ", err)
		return nil, err
	}
	for k, v := range hosts {
		list[k].HostInfo = hosts[k]
		if t, ok := checkStatus[v.Ip]; ok {
			list[k].CheckStatus = t
		}
	}
	res := make([]Item, 0)
	if config.IsCloud() {
		for k, v := range list {
			if v.Uid == hc.uid {
				res = append(res, list[k])
			}
		}
	} else {
		res = list
	}
	return res, nil
}

func (hs *HostInfoService) SystemOsDistribute(uid int) (interface{}, error) {
	return model.NewHostInfo().SystemOsDistribute(uid)
}

func (hs *HostInfoService) Delete(ips []string) error {
	return model.NewHostInfo().BatchDelete(ips)
}
