
CREATE TABLE `user_info` (
  `uid` int NOT NULL AUTO_INCREMENT,
  `name` varchar(16) NOT NULL DEFAULT '',
  `passwd` varchar(16) NOT NULL DEFAULT '',
  `authority` varchar(32) NOT NULL DEFAULT '[]',
  `department` varchar(16) NOT NULL DEFAULT '',
  `is_delete` tinyint NOT NULL DEFAULT '0',
  `update_time` int NOT NULL DEFAULT '0',
  `create_time` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`uid`),
  KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4

 CREATE TABLE `host_info` (
  `id` int NOT NULL AUTO_INCREMENT,
  `ip` varchar(32) NOT NULL DEFAULT '',
  `port` int NOT NULL DEFAULT '22',
  `ssh_user` varchar(32) NOT NULL DEFAULT '',
  `ssh_passwd` varchar(32) NOT NULL DEFAULT '',
  `department` varchar(32) NOT NULL DEFAULT '',
  `business_name` varchar(32) NOT NULL DEFAULT '',
  `system_os` varchar(32) NOT NULL DEFAULT '',
  `uid` int not null DEFAULT 0,
  `update_time` int NOT NULL DEFAULT '0',
  `create_time` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4

alter table host_info add department varchar(32) not null default '' after ssh_passwd;
alter table host_info add system_os varchar(32) not null default '' after department;
alter table host_info add business_name varchar(32) not null DEFAULT '' after system_os;
ALTER table host_info add uid int not null after system_os;



 CREATE TABLE `task_id` (
  `id` int NOT NULL AUTO_INCREMENT,
  `ip` varchar(32) NOT NULL DEFAULT '',
  `is_repeate` tinyint(1) NOT NULL DEFAULT 0,
  `execu_time` varchar(32) NOT NULL DEFAULT '0',
  `create_time` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `ip` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4

create table check_log (
  `id` int NOT NULL AUTO_INCREMENT,
  `ip` varchar(32) NOT NULL DEFAULT '',
  `task_id` int NOT NULL DEFAULT 0,
  `result` text not NULL,
  `create_time` int not null DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `ip` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4


create table vulnerability_log (
  `id` int NOT NULL AUTO_INCREMENT,
  `ip` varchar(32) NOT NULL DEFAULT '',
  `task_id` int NOT NULL DEFAULT 0,
  `result` text not NULL,
  `create_time` int not null DEFAULT 0,
)

CREATE TABLE `vulnerability_log` (
  `id` int NOT NULL AUTO_INCREMENT,
  `ip` varchar(32) NOT NULL DEFAULT '',
  `task_id` int NOT NULL DEFAULT '0',
  `batch_id` int NOT NULL DEFAULT '0',
  `vulnerability_id` varchar(64) NOT NULL DEFAULT '',
  `pkg_name` varchar(256) NOT NULL DEFAULT '',
  `installed_version` varchar(256) NOT NULL DEFAULT '',
  `severity` varchar(32) NOT NULL DEFAULT '',
  `create_time` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `ip` (`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

alter table  vulnerability_log add fixed_version varchar(256) not null default '' after installed_version;


select *, case pkg_name  when @pre_pkg then @curl_rank := @curl_rank + 1 else  @curl_rank := 1 end as c_rank, @pre_pkg := pkg_name   from (select * from vulnerability_log where ip = "112.125.25.235" and batch_id = 1  order by pkg_name, case severity when 'HIGH' then 3 when 'MEDIUM' then 2 else 0 end desc)
 as p, (select @pre_pkg := "", @curl_rank := '') as t 


select id, ip, task_id, batch_id, vulnerability_id, pkg_name, installed_version, severity, create_time, c_rank   from (
 select *, case when @pre_pkg = pkg_name then @curl_rank := @curl_rank + 1 else  @curl_rank := 1 end as c_rank, @pre_pkg := pkg_name   from (select * from vulnerability_log where ip = "112.125.25.235" and batch_id = 1  order by pkg_name, case severity when 'HIGH' then 3 when 'MEDIUM' then 2 else 0 end desc)
 as p, (select @pre_pkg := "", @curl_rank := '') as t ) as t1



 create table task_item_info (
    `id` int NOT NULL AUTO_INCREMENT,
    `ip` varchar(32) NOT NULL DEFAULT '',
    `task_id` int NOT NULL DEFAULT '0',
    `status` int NOT NULL DEFAULT '0',
    `update_time` int NOT NULL DEFAULT '0',
    `create_time` int NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `ip` (`ip`),
    KEY `task_id` (`task_id`)
 ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4

 alter table vulnerability_log change batch_id task_item_id int not null default 0;


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


create table order_info  (
    `id` int not null AUTO_INCREMENT,
    `uid` int not null DEFAULT '',
    `money` int not null 
    `passwd` varchar(16) not null DEFAULT '',
    `authority` varchar(32) NOT NULL DEFAULT '[]',
    `create_time` int NOT NULL DEFAULT '0',
    PRIMARY KEY(`uid`)
 ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4
 
alter table host_info add is_delete tinyint not null default 0 after uid;


create table article_info (
    `id` int not null AUTO_INCREMENT,
    `uid` int not null DEFAULT 0,
    `content` text not null,
    `is_delete` tinyint not null DEFAULT 0,
    `update_time` int NOT NULL DEFAULT '0',
    `create_time` int NOT NULL DEFAULT '0',
    PRIMARY KEY(`id`),
    KEY(`uid`)
)  ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4

CREATE TABLE `action_log` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uid` int NOT NULL DEFAULT '0',
  `type` tinyint NOT NULL DEFAULT '0',
  `ip` char(32) NOT NULL DEFAULT '',
  `detail` varchar(128) NOT NULL DEFAULT '0',
  `create_time` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `uid` (`uid`),
  KEY `type` (`type`),
  KEY `ip` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=93 DEFAULT CHARSET=utf8mb4