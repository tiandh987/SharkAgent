
CREATE DATABASE IF NOT EXISTS `iam`
USE `iam`;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
    `instanceID` varchar(32) DEFAULT NULL,
    `name` varchar(45) NOT NULL,
    `status` int(1) DEFAULT 1 COMMENT '1:可用, 0:不可用',
    `nickname` varchar(30) NOT NULL,
    `password` varchar(255) NOT NULL,
    `email` varchar(256) NOT NULL,
    `phone` varchar(20) DEFAULT NULL,
    `isAdmin` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '1: administrator, 0: non-administrator',
    `extendShadow` longtext DEFAULT NULL,
    `loginedAt` timestamp NULL DEFAULT NULL COMMENT 'last login time',
    `createdAt` timestamp NOT NULL DEFAULT current_timestamp(),
    `updatedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_name` (`name`),
    UNIQUE KEY `instanceID_UNIQUE` (`instanceID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

INSERT INTO `user` VALUES (1,'user-admin','admin',1,'admin','$2a$10$WnQD2DCfWVhlGmkQ8pdLkesIGPf9KJB7N1mhSOqulbgN7ZMo44Mv2','admin@foxmail.com','1812884xxxx',1,'{}',now(),'2022-04-23 17:27:40','2022-04-23 17:27:40');
