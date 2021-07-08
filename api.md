## 接口文档说明
 <br/>

### 1, 获取验证码id
    curl -X GET http://127.0.0.1:8090/util/get_captcha_id

    {
        "code":0,
        "message":"ok",
        "data":{
            "id":"bxEqIF4TbIynJeSgccFM"
        }
    }

    说明：id是验证码的id，通过该id向后台请求一张验证码图片

### 2，获取验证码
    curl -X GET http://127.0.0.1:8090/util/get_captcha?id=bxEqIF4TbIynJeSgccFM

    说明：将会返回一个验证码图片


### 3, 登陆说明

    curl -X POST http://127.0.0.1:8090/user/login -d 'name=secadmin&passwd=secadmin@123$&captcha_id=ZMbCKtXkjxuKzDp2bMm5&captcha_value=961066'

    
    {
        "code":0,
        "message":"ok",
        "data":{
            "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVaWQiOjMwLCJleHAiOjE2MjUyNzY0NjB9.6esJoJwcIDulSJNUTGl3ScqmwbiAReP8oWcFM_8i2B4"
        }
    }
    
### 4, 添加主机

     curl -X POST http://127.0.0.1:8090/host/add -d 'ip=47.104.213.134&port=22&ssh_user=root&ssh_passwd=ylysec.coM515!@#' -H "Token:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVaWQiOjMwLCJleHAiOjE2MjUyNzczOTZ9.PpEH2SPxv8wG6wmDbAfH-jiajkTj8qxuY2K06lptyzg"

    {
        "code":0,
        "message":"ok",
        "data":{
        }
    }
    

### 5, 权限配置

    curl -X GET http://127.0.0.1:8090/config/get_const_config
    
    {
        "code":0,
        "message":"ok",
        "data":{
            "add_host":2,
            "add_user":5,
            "check_docker":11,
            "check_log":6,
            "check_soft":3,
            "create_security_task":4,
            "create_smart_task":8,
            "create_user_group":10,
            "delete_user":7,
            "update_soft":9
        }
    }   
    

### 6, 用户列表
    curl -X GET http://127.0.0.1:8090/user/list_users

        {
            "code":0,
            "message":"ok",
            "data":{
            "list":[
                {
                    "uid":30,
                    "name":"secadmin",
                    "authority":[
                        1
                    ],
                    "department":""
                },
                {
                    "uid":31,
                    "name":"fudake",
                    "authority":[
                        2,
                        3
                    ],
                    "department":"安全部门"
                }
            ]
        }
    }

### 7  搜索接口
    按照ip搜索   curl -X POST http://127.0.0.1:8090/host/search_host -d 'type=0&condition=47.104'
    按照部门搜索  curl -X POST http://127.0.0.1:8090/host/search_host -d 'type=1&condition=测试部门'
    
    {
    "code":0,
    "message":"ok",
    "data":{
        "list":[
            {
                "ip":"112.125.25.235",
                "port":22,
                "ssh_user":"root",
                "department":"测试部门",
                "system_os":"Ubuntu"
            }
        ]
    }
}

### 8 当前任务情况
    curl -X GET http://47.104.213.134:8080/task/curl_task_info
    
    {
        "code":0,
        "message":"ok",
        "data":{
            "all":1,
            "checking":0,
            "queue_task":0,
            "planing_task":1
        }
    }

### 9 系统日志
    curl -X GET http://127.0.0.1:8090/log/list?last_id=-1
    说明：last_id是上一页日志的最后一条id，第一次给-1
    {
        "code":0,
        "message":"ok",
        "data":{
            "list":[
                {
                    "id":18,
                    "uid":30,
                    "detail":"secadmin用户添加了主机47.104.213.134",
                    "create_time":1624672685
                },
                {
                    "id":17,
                    "uid":30,
                    "detail":"secadmin用户登陆",
                    "create_time":1624672596
                }
            ]
        }
    }

### 10 修改密码
     curl -X POST http://127.0.0.1:8090/user/reset_passwd -d 'uid=31&pass_wd=1233'
     {
         "code":0,
         "message":"ok",
         "data":{}
    }

### 11 查看软件包信息
    curl -X GET http://127.0.0.1:8090/host/get_vulnerability_info?ip=47.104.213.134
    {   
        "code":0,
        "message":"ok",
        "data":{
            "list":[
                {
                    "installed_version":"0.6.55-0ubuntu12~20.04.4",
                    "pkg_name":"accountsservice",
                    "severity":"LOW",
                    "vulnerability_id":"CVE-2012-6655"
                }
            ]
        }
    }


###  12 查看系统信息

    curl -X GET http://127.0.0.1:8090/util/get_system_info

    {   
        "code":0,
        "message":"ok",
        "data":{
            "edition":0, // 系统版本 0， 免费版， 1 企业版
            "max_node":2
        }
    }

### 13 用户登陆

    curl -X POST http://127.0.0.1:8090/user/login -d 'name=fudake&passwd=fudake&captcha_id=123345&captcha_value=1234'

    {
        "code": 0,
        "message": "ok",
        "data": {
            "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVaWQiOjMyLCJleHAiOjE2MjU3MjQ5Mjl9.cGRvc4t3QermU9QGG83OnnWB5GySpT_8aZK0A_o_Gq0"
        }
    }

### 14 找回密码

    curl -X POST http://127.0.0.1:8090/user/find_passwd -d 'account=13152015823'

    {
        "code":0,
        "message":"ok",
        "data":{
            "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVaWQiOjEwMDAwMTgsImV4cCI6MTYyNTg4OTk5N30.fNEz_hW3qpZQsY-6DICdVPYFH5YACW7WfqJUuy-7EbM"
        }
    }
    account 是手机号或邮箱


### 15 用户注册
    curl -X POST http://127.0.0.1:8090/user/register -d 'company=123&phone=13155555555&emal=123@qq.com&passwd=1234&captcha_id=1234&captcha_value=123456'

    {
        "code":0,
        "message":"ok",
        "data":{
            "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVaWQiOjMwLCJleHAiOjE2MjUyNzY0NjB9.6esJoJwcIDulSJNUTGl3ScqmwbiAReP8oWcFM_8i2B4"
        }
### 16 删除主机

    curl -X POST http://127.0.0.1:8090/host/delete -d 'ip=112.125.25.235'

    {   
        "code":0,
        "message":"ok",
        "data":{}
    }