-- 添加已读/未读字段到消息表
ALTER TABLE `lc_chat_message`
ADD COLUMN `is_read` TINYINT(1) DEFAULT 0 COMMENT '是否已读' AFTER `faq_id`,
ADD COLUMN `read_at` DATETIME(0) NULL COMMENT '阅读时间' AFTER `is_read`,
ADD INDEX `idx_is_read` (`is_read`);
