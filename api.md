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

