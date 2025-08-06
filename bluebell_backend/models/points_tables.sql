-- 用户积分表
DROP TABLE IF EXISTS `user_points`;
CREATE TABLE `user_points` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `user_id` bigint(20) NOT NULL COMMENT '用户ID',
    `points` int(11) NOT NULL DEFAULT '0' COMMENT '当前积分',
    `points_total` int(11) NOT NULL DEFAULT '0' COMMENT '总积分',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- 积分交易明细表
DROP TABLE IF EXISTS `user_points_transactions`;
CREATE TABLE `user_points_transactions` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `transaction_id` bigint(20) NOT NULL COMMENT '交易ID',
    `user_id` bigint(20) NOT NULL COMMENT '用户ID',
    `points_change` int(11) NOT NULL COMMENT '积分变化(正数为增加，负数为减少)',
    `current_balance` int(11) NOT NULL COMMENT '当前余额',
    `transaction_type` tinyint(4) NOT NULL COMMENT '交易类型(1:签到 2:发布帖子 3:评论 4:点赞 5:兑换奖励)',
    `description` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '描述',
    `ext_json` varchar(512) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '扩展信息JSON',
    `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_transaction_id` (`transaction_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_transaction_type` (`transaction_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;