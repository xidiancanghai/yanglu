

create table user_info (
    uid int not null AUTO_INCREMENT,
    passwd varchar(16) not null DEFAULT '',
    authority TINYINT not NULL DEFAULT 0,
    department varchar(16) not NULL DEFAULT '',
    update_time int NOT NULL DEFAULT 0,
    create_time int not NULL DEFAULT 0,
    PRIMARY KEY(`uid`)
)  ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4