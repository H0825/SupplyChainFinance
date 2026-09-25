/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 80040
 Source Host           : localhost:3306
 Source Schema         : finance

 Target Server Type    : MySQL
 Target Server Version : 80040
 File Encoding         : 65001

 Date: 11/03/2026 14:47:22
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for documents
-- ----------------------------
DROP TABLE IF EXISTS `documents`;
CREATE TABLE `documents`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `cid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `tag` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `operator` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `receivable_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `order_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `doc_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'other',
  `status` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'uploaded',
  `source` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'ipfs',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_cid`(`cid`) USING BTREE,
  INDEX `idx_doc_created_at`(`created_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of documents
-- ----------------------------
INSERT INTO `documents` VALUES (1, '屏幕截图 2025-03-26 224940.png', 'QmfZv9VZ4pHrfLKNDr4JdR2GhWHd5GhnJfpewLXQwwWxpp', 'R-1773154429824', 'user001', '2026-03-10 22:54:10', 'R-1773154429824', 'O-1773154429824', 'contract', 'uploaded', 'ipfs');
INSERT INTO `documents` VALUES (2, '屏幕截图 2026-02-27 001536.png', 'QmQvUd2bzkxihreCZVCdLNGqC7Zd68P4nvLvhbpRyvyMGC', 'R-1773154466625', 'user001', '2026-03-10 22:54:44', 'R-1773154466625', 'O-1773154466625', 'contract', 'uploaded', 'ipfs');
INSERT INTO `documents` VALUES (3, '屏幕截图 2026-01-19 205005.png', 'QmbS92eGQw1foF8JafTWUZ1wuBTqtApv4qJNWHzrqSprdo', 'R-1773154515255', 'user001', '2026-03-10 22:56:23', 'R-1773154515255', 'O-1773154515255', 'contract', 'uploaded', 'ipfs');
INSERT INTO `documents` VALUES (4, '屏幕截图 2026-02-27 001536.png', 'QmQvUd2bzkxihreCZVCdLNGqC7Zd68P4nvLvhbpRyvyMGC', 'R-1773155026927', 'user001', '2026-03-10 23:04:06', 'R-1773155026927', 'O-1773155026927', 'contract', 'uploaded', 'ipfs');
INSERT INTO `documents` VALUES (5, '屏幕截图 2025-06-15 152714.png', 'QmUgoec1BPQpYqPmEkVms61zsZGurNQQrMZyMLb93qVk5m', 'R-1773156986827', 'user001', '2026-03-10 23:36:55', 'R-1773156986827', 'O-1773156986827', 'contract', 'uploaded', 'ipfs');

-- ----------------------------
-- Table structure for enterprises
-- ----------------------------
DROP TABLE IF EXISTS `enterprises`;
CREATE TABLE `enterprises`  (
  `enterprise_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `wallet` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `is_core` tinyint(1) NOT NULL DEFAULT 0,
  `is_financial` tinyint(1) NOT NULL DEFAULT 0,
  `register_time` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`enterprise_id`) USING BTREE,
  UNIQUE INDEX `uk_enterprise_wallet`(`wallet`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of enterprises
-- ----------------------------
INSERT INTO `enterprises` VALUES ('ENT-20260310233626-522', '天天海信', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', 1, 0, 1773156998);

-- ----------------------------
-- Table structure for financing_records
-- ----------------------------
DROP TABLE IF EXISTS `financing_records`;
CREATE TABLE `financing_records`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `receivable_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `financial_inst` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `amount` decimal(30, 0) NOT NULL DEFAULT 0,
  `interest_rate` decimal(30, 0) NOT NULL DEFAULT 0,
  `financing_time` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_financing_receivable`(`receivable_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of financing_records
-- ----------------------------
INSERT INTO `financing_records` VALUES (1, 'R-1773156986827', '0x7469ea0121c75C46fF3CC7cC4E90F9e929C3f15d', 120000, 420, 1773157073);

-- ----------------------------
-- Table structure for receivables
-- ----------------------------
DROP TABLE IF EXISTS `receivables`;
CREATE TABLE `receivables`  (
  `receivable_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `order_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `issuer` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `payer` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `amount` decimal(30, 0) NOT NULL DEFAULT 0,
  `issue_time` bigint NOT NULL DEFAULT 0,
  `due_time` bigint NOT NULL DEFAULT 0,
  `status` tinyint UNSIGNED NOT NULL DEFAULT 0,
  PRIMARY KEY (`receivable_id`) USING BTREE,
  INDEX `idx_order_id`(`order_id`) USING BTREE,
  INDEX `idx_status`(`status`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of receivables
-- ----------------------------
INSERT INTO `receivables` VALUES ('R-1773155026927', 'O-1773155026927', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '0xe1C6eA9f8bF3a3e6e7F6aC3DB68017Eb24CB598a', 120000, 1773155048, 1778339026, 3);
INSERT INTO `receivables` VALUES ('R-1773156986827', 'O-1773156986827', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '0xe1C6eA9f8bF3a3e6e7F6aC3DB68017Eb24CB598a', 120000, 1773157017, 1778340986, 2);

-- ----------------------------
-- Table structure for tx_logs
-- ----------------------------
DROP TABLE IF EXISTS `tx_logs`;
CREATE TABLE `tx_logs`  (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `action` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `business_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `tx_hash` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` int NOT NULL,
  `message` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `operator_addr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_business_id`(`business_id`) USING BTREE,
  INDEX `idx_created_at`(`created_at`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of tx_logs
-- ----------------------------
INSERT INTO `tx_logs` VALUES (1, 'register_enterprise', 'ENT-20260310222733-763', '0xf4d816b416f859f10fddae8f513518ab8f04ba409a4b660601f7e97cdcf47333', 0, 'register enterprise success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:27:43');
INSERT INTO `tx_logs` VALUES (2, 'issue_receivable', 'R-1773152853387', '0x614074c44f58b9f047ab0ce256af9e449ccc67e26e23b22c0ee48e96dda35d44', 0, 'issue receivable success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:28:23');
INSERT INTO `tx_logs` VALUES (3, 'confirm_receivable', 'R-1773152853387', '0xe41093903065edb88c927d76ea7baaed081ccb5424c88b97252529c9384a9442', 0, 'confirm receivable success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:28:30');
INSERT INTO `tx_logs` VALUES (4, 'confirm_receivable', 'R-1773152853387', '0xd7f3c69fdf1c5675a71bef40784de2dfc7895accb90071527aae7b92c4d627e6', 22, 'confirm receivable reverted', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:28:38');
INSERT INTO `tx_logs` VALUES (5, 'register_enterprise', 'ENT-20260310222856-820', '0xf71450c89374610a37a1df40a163fec6264322d7a0ae8da01478492deb074a80', 22, 'register enterprise reverted', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:29:05');
INSERT INTO `tx_logs` VALUES (6, 'issue_receivable', 'R-1773153000276', '0x9769bfb4777a2727d2a3c41ce3f570400c1bc08e63a7e864d5957a726485198b', 0, 'issue receivable success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:30:15');
INSERT INTO `tx_logs` VALUES (7, 'register_enterprise', 'ENT-20260310223102-808', '0x8406b42f3f3c72ab2b79e0c162a1500b1dd83628173aa403e39a69bb5e7a398e', 22, 'register enterprise reverted', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:31:05');
INSERT INTO `tx_logs` VALUES (8, 'issue_receivable', 'R-1773153054667', '0x180ceea90f4d4ec13eba652246afb3e2d9f40f64db36f015f08eef64bd682909', 0, 'issue receivable success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:31:44');
INSERT INTO `tx_logs` VALUES (9, 'confirm_receivable', 'R-1773153054667', '0x089c22bf63c7f2b234f6cca8d2e1417538119c6d904d2f1d4f8465c37fc0860f', 0, 'confirm receivable success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:31:49');
INSERT INTO `tx_logs` VALUES (10, 'confirm_receivable', 'R-1773153054667', '0x374bf7cfbad8e5abda00f880d28dac95219f857aad68711f3844bc42ad118a49', 22, 'confirm receivable reverted', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:31:55');
INSERT INTO `tx_logs` VALUES (11, 'record_settlement', 'R-1773153054667', '0xc1f2f8a3a89b10a169aee5a3be2ed4dcdadd71f53ad62718100375e7b0bf9b40', 0, 'record settlement success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:32:00');
INSERT INTO `tx_logs` VALUES (12, 'register_enterprise', 'ENT-20260310225349-813', '0x0d2b7a049032d688766bb1a240b81ac1ca436f7973ee9556a487718f4e85ac28', 0, 'register enterprise success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:53:58');
INSERT INTO `tx_logs` VALUES (13, 'issue_receivable', 'R-1773154429824', '0xb1cbf6914e5db9796d603fe0aa2e02a7b3ba7ef4a032ca181f1297f84ed23d01', 0, 'issue receivable success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:54:15');
INSERT INTO `tx_logs` VALUES (14, 'issue_receivable', 'R-1773154466625', '0x89198c521d60bfe48758afe2a315b9c64675ff7971514522a99c4f7f57696a84', 0, 'issue receivable success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:54:45');
INSERT INTO `tx_logs` VALUES (15, 'confirm_receivable', 'R-1773154466625', '0xd65b32b13b3b3b17b2d52de03520f8051aac7182e0e7ee2948eacb5d854a50c3', 0, 'confirm receivable success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:54:48');
INSERT INTO `tx_logs` VALUES (16, 'record_settlement', 'R-1773154466625', '0xb8a95463d60f1007ac4e7cf01368a068ba83e0f7716805ab6fba83e17dbfbadb', 0, 'record settlement success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 22:55:01');
INSERT INTO `tx_logs` VALUES (17, 'issue_receivable', 'R-1773154515255', '0x94df88e435ce72cc1d7339172c6db8dcf1da128e15d8800931930142df560bd5', 0, 'issue receivable success', '0xc13FC67f57046a56252Ef63deD46eb9540EF0292', '2026-03-10 22:56:25');
INSERT INTO `tx_logs` VALUES (18, 'confirm_receivable', 'R-1773154515255', '0x9d0b1a259c6534092018f4850e382eff629fc5f04e3ec86f5dc55ba671814904', 0, 'confirm receivable success', '0xc13FC67f57046a56252Ef63deD46eb9540EF0292', '2026-03-10 22:56:29');
INSERT INTO `tx_logs` VALUES (19, 'record_settlement', 'R-1773154515255', '0x954ad8d498bbb645184b02a2c2a82361138d997c6f3f42d03db83a8f978ec49c', 0, 'record settlement success', '0xc13FC67f57046a56252Ef63deD46eb9540EF0292', '2026-03-10 22:56:31');
INSERT INTO `tx_logs` VALUES (20, 'register_enterprise', 'ENT-20260310230346-385', '0x85ea8d8041f930734ee955d66c319af0597e1ad9b970ca993fef147bbfe9a650', 22, 'register enterprise reverted', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 23:03:57');
INSERT INTO `tx_logs` VALUES (21, 'issue_receivable', 'R-1773155026927', '0x2871d52a1e5441b4ecac8bc71563c480f1339b4ad51d2b098639b8f7c1678d52', 0, 'issue receivable success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 23:04:08');
INSERT INTO `tx_logs` VALUES (22, 'confirm_receivable', 'R-1773155026927', '0xa42be89a6d4479b6c3f640ea25b1cc7aad9b6b33573750f3a6b0c4cfdefc40e5', 0, 'confirm receivable success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 23:04:20');
INSERT INTO `tx_logs` VALUES (23, 'record_settlement', 'R-1773155026927', '0xcb1c65e367df331602fbe3c629ebcd8f155e1efb36152f5a063874105fa43fae', 0, 'record settlement success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 23:10:20');
INSERT INTO `tx_logs` VALUES (24, 'register_enterprise', 'ENT-20260310233626-522', '0x37f5ecec6fe8665a4b0ed13b7579df0d4fa43f8e09fd80d6b37137ecd0528277', 0, 'register enterprise success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 23:36:38');
INSERT INTO `tx_logs` VALUES (25, 'issue_receivable', 'R-1773156986827', '0x9bf7639a2d664ea8ae64334c7100ef76951a97fd2a323fa760e427039228c865', 0, 'issue receivable success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 23:36:57');
INSERT INTO `tx_logs` VALUES (26, 'confirm_receivable', 'R-1773156986827', '0xb54669a198220a03f92f46fc2537a8d80ba4dde6d84d62d747b7fa73eb05490f', 0, 'confirm receivable success', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', '2026-03-10 23:37:01');
INSERT INTO `tx_logs` VALUES (27, 'record_financing', 'R-1773156986827', '0x31a6653e2831cbef9dd2bf0d4a064730d6e4d1a2a3b5b39ef67b4dde2bc994d9', 0, 'record financing success', '0x7469ea0121c75C46fF3CC7cC4E90F9e929C3f15d', '2026-03-10 23:37:53');

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `role` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'business',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uk_username`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (1, 'user001', '$2a$10$gQw6kvvoWq79/c23IAHs9uOVWM4xRmvegZ2bmiVqvbO88EVTjosyG', '0xE25eA7CbDaeA4E606435BBC1a96e8caE259c988b', 'business');
INSERT INTO `users` VALUES (2, 'jr001', '$2a$10$nhB0DfFDGgDcRta2/Rw8Sev05LY4I4n8OoUdzuuylP8TZxiHMaxjq', '0x7469ea0121c75C46fF3CC7cC4E90F9e929C3f15d', 'finance');
INSERT INTO `users` VALUES (3, 'admin001', '$2a$10$iqo7gWNZ.N/RHItPsy.sWOuayzESuBsAsBiYb.wnhGa64NCojTKpq', '0xe1C6eA9f8bF3a3e6e7F6aC3DB68017Eb24CB598a', 'admin');

SET FOREIGN_KEY_CHECKS = 1;
