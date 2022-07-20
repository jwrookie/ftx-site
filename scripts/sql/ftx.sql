CREATE DATABASE IF NOT EXISTS `ftx` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `lucky` (
    `lucky_id` bigint(20) unsigned NOT NULL COMMENT '主键',
    `email` varchar(255) NOT NULL DEFAULT '' COMMENT '邮箱',
    `kyc_level` varchar(10) NOT NULL DEFAULT '' COMMENT 'kyc 账号等级',
    `personality` varchar(10) NOT NULL DEFAULT '' COMMENT '人格',
    `prize` varchar(255) NOT NULL DEFAULT '' COMMENT '奖品',
    `clothes_size` varchar(10) NOT NULL DEFAULT '' COMMENT '衣服尺码，仅当奖品为衣服时，该字段才有意义',
    `user_name` varchar(255)  NOT NULL DEFAULT '' COMMENT '收件人姓名',
    `user_phone` varchar(255)  NOT NULL DEFAULT '' COMMENT '收件人联系方式',
    `country` varchar(255)  NOT NULL DEFAULT '' COMMENT '国家',
    `region` varchar(255)  NOT NULL DEFAULT '' COMMENT '行政区域',
    `postal_code` varchar(255)  NOT NULL DEFAULT '' COMMENT '邮编',
    `address` varchar(255) NOT NULL DEFAULT '' COMMENT '收件地址',
    `inviter_email` varchar(255) NOT NULL DEFAULT '' COMMENT '邀请人邮箱',
    `created_at` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
    `updated_at` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
    `deleted_at` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '删除时间',
    PRIMARY KEY (`lucky_id`) USING BTREE,
    UNIQUE KEY `email` (`email`),
    KEY `inviter_email` (`inviter_email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='中奖信息表';