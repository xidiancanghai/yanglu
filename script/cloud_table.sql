

 create table cloud_user_info (
    `uid` int not null AUTO_INCREMENT,
    `company` varchar(32) not null DEFAULT '',
    `phone` varchar(16) not null DEFAULT '',
    `email` varchar(32) not null DEFAULT '',
    `passwd` varchar(16) not null DEFAULT '',
    `authority` varchar(32) NOT NULL DEFAULT '[]',
    `create_time` int NOT NULL DEFAULT '0',
    PRIMARY KEY(`uid`)
 ) ENGINE=InnoDB AUTO_INCREMENT=1000001 DEFAULT CHARSET=utf8mb4

create table user_tmp_passwd (
    `uid` int not null AUTO_INCREMENT,
    `pass_wd` varchar(16) not null DEFAULT '',
    `is_delete` tinyint (1) not null DEFAULT 0,
    `update_time` int NOT NULL DEFAULT '0',
    `create_time` int NOT NULL DEFAULT '0',
    PRIMARY KEY(`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4