package service

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
	"yanglu/service/model"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type ExcelService struct {
	excelFile *excelize.File
}

func NewExcelService(file io.Reader) (*ExcelService, error) {
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	es := new(ExcelService)
	es.excelFile = xlsx
	return es, nil
}

func (es *ExcelService) GetHostInfos() ([]*model.HostInfo, error) {
	rows := es.excelFile.GetRows("Sheet1")
	if rows == nil {
		return nil, errors.New("缺少表")
	}
	// 读取第一行
	index := 0
	headers := []string{}
	list := []*model.HostInfo{}
	for _, row := range rows {

		hostInfo := new(model.HostInfo)

		for i, cell := range row {
			cell = strings.TrimSpace(cell)
			if cell == "" && headers[i] != "department" {
				return nil, errors.New("输入的值为空")
			}
			if index == 0 {
				headers = append(headers, cell)
			} else {
				switch i {
				case 0:
					hostInfo.Ip = cell
				case 1:
					hostInfo.Port, _ = strconv.Atoi(cell)
				case 2:
					hostInfo.SshUser = cell
				case 3:
					hostInfo.SshPasswd = cell
				case 4:
					hostInfo.Department = cell
				case 5:
					hostInfo.BusinessName = cell
				}
			}
		}
		if index == 0 {
			if !es.checkHeaders(headers) {
				return nil, errors.New("表头错误，格式 ip, port, name, passwd, department, business_name")
			}
			index++
		} else {
			hostInfo.CreateTime = time.Now().Unix()
			hostInfo.UpdateTime = hostInfo.CreateTime
			list = append(list, hostInfo)
		}
	}
	return list, nil
}

func (es *ExcelService) checkHeaders(headers []string) bool {
	columns := []string{"ip", "port", "name", "passwd", "department"}
	if len(headers) != len(columns) {
		return false
	}
	for i, v := range headers {
		if columns[i] != v {
			return false
		}
	}
	return true
}
