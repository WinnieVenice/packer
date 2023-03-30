-- Active: 1664641776568@@127.0.0.1@3306@cq

CREATE DATABASE cq
    DEFAULT CHARACTER SET = 'utf8mb4';

CREATE TABLE `user_info` (
    `id` INT UNSIGNED AUTO_INCREMENT,
    `qq` INT UNSIGNED NOT NULL,
    `cf` CHAR(20),
    `at` CHAR(20),
    `lg` CHAR(20),
    `vj` CHAR(20),
    `lc` CHAR(20),
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

SELECT CONCAT ('DROP TABLE IF EXISTS', table_name, ';') FROM information_schema.tables WHERE TABLE_SCHEMA = 'cq';

select * FROM information_schema.tables WHERE TABLE_SCHEMA = 'cq';

select * from user_infos;

DROP Table user_infos;

select * from user_infos WHERE cf_id = 'WinnieVenice';

select * from platform_user_infos;

insert into user_infos (id, nick_name, codeforces_id) VALUES (2602306472, "wwj", "baddog");

insert into user_infos (id, nick_name, codeforces_id) VALUES (2633092012, "lxl", "kejunyu")