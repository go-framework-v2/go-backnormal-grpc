CREATE TABLE `test_person` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '名称',
  `idcard` varchar(100) NOT NULL COMMENT '身份证号',
  `age` tinyint NOT NULL COMMENT '年龄',
  `gender` enum('1', '2') COMMENT '枚举，1男2女',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `valid` tinyint NOT NULL DEFAULT '1' COMMENT '逻辑删除标记，1有效0无效',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_testPerson_name_idcard_valid` (`name`,`idcard`,`valid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='人';