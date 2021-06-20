

create table user_info (
    uid int not null AUTO_INCREMENT,
    passwd varchar(16) not null DEFAULT '',
    authority TINYINT not NULL DEFAULT 0,
    department varchar(16) not NULL DEFAULT '',
    update_time int NOT NULL DEFAULT 0,
    create_time int not NULL DEFAULT 0,
    PRIMARY KEY(`uid`)
)  ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4


 CREATE TABLE `host_info` (
  `id` int NOT NULL AUTO_INCREMENT,
  `ip` varchar(32) NOT NULL DEFAULT '',
  `port` int NOT NULL DEFAULT '22',
  `ssh_user` varchar(32) NOT NULL DEFAULT '',
  `ssh_passwd` varchar(32) NOT NULL DEFAULT '',
  `update_time` int NOT NULL DEFAULT '0',
  `create_time` int NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ip` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4

alter table host_info add department varchar(32) not null default '' after ssh_passwd;
alter table host_info add system_os varchar(32) not null default '' after department;
alter table host_info add business varchar(32) not null DEFAULT '' after system_os;



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