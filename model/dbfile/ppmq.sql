-- phpMyAdmin SQL Dump
-- version 4.7.3
-- https://www.phpmyadmin.net/
--
-- Host: localhost
-- Generation Time: 2018-09-24 06:25:04
-- 服务器版本： 10.2.17-MariaDB
-- PHP Version: 7.2.9

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";

--
-- Database: `ppmq`
--

-- --------------------------------------------------------

--
-- 表的结构 `account`
--

CREATE TABLE `account` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `account` varchar(64) NOT NULL COMMENT '账户信息',
  `password` varchar(256) NOT NULL COMMENT '账户对应密码，md5(text password)'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息系统： 账户信息';

-- --------------------------------------------------------

--
-- 表的结构 `clientid`
--

CREATE TABLE `clientid` (
  `id` int(11) NOT NULL,
  `client_id` varchar(64) NOT NULL COMMENT '全局唯一，确保相同的帐号+硬件得到的ClientId是一致',
  `account` varchar(64) NOT NULL COMMENT 'user_type = 1,这里是device_id; user_type = 2, 这里是用户电话号码',
  `hw_feature` varchar(256) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息: 记录分配的client_id';

-- --------------------------------------------------------

--
-- 表的结构 `connection`
--

CREATE TABLE `connection` (
  `id` int(11) NOT NULL,
  `server_id` varchar(255) NOT NULL COMMENT '每个ppmqd服务对应的一个编号，编号唯一: groupid_serverip',
  `clear_session` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0,保存会话信息;1: 会话仅仅维持当前连接',
  `will_flag` int(11) NOT NULL,
  `will_qos` int(11) NOT NULL,
  `will_retain` int(11) NOT NULL,
  `client_id` varchar(64) NOT NULL COMMENT '关联msg_clientid',
  `user_id` int(12) NOT NULL,
  `will_topic` varchar(256) NOT NULL,
  `will_body` longblob NOT NULL,
  `conn_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1,tcp ppmq;2,udp; 3, mqtt',
  `conn_info` varchar(256) NOT NULL COMMENT '连接信息',
  `historymsg_type` int(11) NOT NULL COMMENT '1 不获取，只用于接收新消息； 2 只用于获取离线消息； 3 同时获取新消息和离线消息.',
  `historymsg_count` int(11) NOT NULL COMMENT '可以设定获取离线消息的数量;0表示获取所有的历史消息',
  `is_online` tinyint(1) NOT NULL DEFAULT 0 COMMENT '1, 在线; 2, 离线',
  `last_time` bigint(20) NOT NULL COMMENT '最后更新状态的时间',
  `global_sync` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0,1,等待同步; 2,同步中; 3,同步完成 ',
  `account` varchar(80) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息: 连接';

-- --------------------------------------------------------

--
-- 表的结构 `msg_info`
--

CREATE TABLE `msg_info` (
  `id` int(11) NOT NULL,
  `msg_id` varchar(64) NOT NULL,
  `src_msgid` varchar(64) NOT NULL,
  `account` varchar(64) NOT NULL COMMENT '生产者ID; msg_account.account',
  `client_id` varchar(255) NOT NULL,
  `dup` tinyint(1) NOT NULL,
  `retain` tinyint(1) NOT NULL,
  `qos` tinyint(1) NOT NULL,
  `topic` varchar(256) NOT NULL,
  `format` tinyint(2) NOT NULL COMMENT '描述Body的编码格式: 4, PB-BIN; 5, PB-JSON;',
  `cmdid` int(11) NOT NULL COMMENT '描述Body的编码使用的命令ID： 0， 表示使用了非预定以命令，需要应用程序自己处理.',
  `cmd_type` tinyint(2) NOT NULL COMMENT '描述Body的内容的命令类型: 0 请求：  1 应答',
  `msg_timems` bigint(13) NOT NULL COMMENT '终端产生消息时间，单位UTC，ms',
  `create_time` bigint(14) NOT NULL COMMENT '创建消息的时间单位毫秒'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息：存放当前未完成流程的消息';

-- --------------------------------------------------------

--
-- 表的结构 `msg_log`
--

CREATE TABLE `msg_log` (
  `id` int(11) NOT NULL,
  `msg_id` varchar(64) NOT NULL,
  `account` varchar(64) NOT NULL COMMENT '消费者ID, msg_account.account',
  `client_id` varchar(64) NOT NULL,
  `server_id` varchar(255) NOT NULL,
  `status` tinyint(2) NOT NULL COMMENT 'QOS1: 1,send; 2, pubAns; QOS2: 1, pub; 2, pubrec；3, pubrel；4, pubcomp',
  `create_time` bigint(14) NOT NULL COMMENT '创建消息的时间单位毫秒'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息：投递日志，用于后期分析';

-- --------------------------------------------------------

--
-- 表的结构 `msg_raw`
--

CREATE TABLE `msg_raw` (
  `id` int(11) NOT NULL,
  `msg_id` varchar(64) NOT NULL,
  `msg_payload` longblob NOT NULL COMMENT '除开固定和可变报头的消息载荷'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息：存放当前未完成流程的消息';

-- --------------------------------------------------------

--
-- 表的结构 `msg_status`
--

CREATE TABLE `msg_status` (
  `id` int(11) NOT NULL,
  `msg_id` varchar(64) NOT NULL,
  `src_msgid` varchar(64) NOT NULL,
  `account` varchar(64) NOT NULL COMMENT '消费者ID, msg_account.account',
  `client_id` varchar(64) NOT NULL,
  `server_id` varchar(255) NOT NULL,
  `dup` tinyint(1) NOT NULL,
  `qos` tinyint(1) NOT NULL,
  `status` tinyint(2) NOT NULL COMMENT 'QOS1: 1,send; 2, pubAns; QOS2: 1, pub; 2, pubrec；3, pubrel；4, pubcomp',
  `create_time` bigint(14) NOT NULL COMMENT '创建消息的时间单位毫秒'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息：当前流程进行中的QOS消息(每条消息和clientid的组合只有一条记录)';

-- --------------------------------------------------------

--
-- 表的结构 `subscribe`
--

CREATE TABLE `subscribe` (
  `id` int(11) NOT NULL,
  `account` varchar(64) NOT NULL,
  `client_id` varchar(64) NOT NULL COMMENT '客户端标识，依照次字段确认是否重复订阅',
  `topic` varchar(256) NOT NULL,
  `qos` tinyint(1) NOT NULL,
  `cluster` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否集群消费',
  `cluster_subid` varchar(64) NOT NULL COMMENT '集群消费ID',
  `last_time` bigint(13) NOT NULL,
  `global_sync` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0,1,等待同步; 2,同步中; 3,同步完成 '
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息：订阅记录';

--
-- Indexes for dumped tables
--

--
-- Indexes for table `account`
--
ALTER TABLE `account`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `uid_account` (`user_id`,`account`),
  ADD KEY `account` (`account`),
  ADD KEY `user_id` (`user_id`);

--
-- Indexes for table `clientid`
--
ALTER TABLE `clientid`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `client_id` (`client_id`),
  ADD KEY `account` (`account`);

--
-- Indexes for table `connection`
--
ALTER TABLE `connection`
  ADD PRIMARY KEY (`id`),
  ADD KEY `client_id` (`client_id`),
  ADD KEY `server_id` (`server_id`);

--
-- Indexes for table `msg_info`
--
ALTER TABLE `msg_info`
  ADD PRIMARY KEY (`id`),
  ADD KEY `msg_id` (`msg_id`);

--
-- Indexes for table `msg_log`
--
ALTER TABLE `msg_log`
  ADD PRIMARY KEY (`id`),
  ADD KEY `msg_id` (`msg_id`),
  ADD KEY `client_id` (`client_id`);

--
-- Indexes for table `msg_raw`
--
ALTER TABLE `msg_raw`
  ADD PRIMARY KEY (`id`),
  ADD KEY `msg_id` (`msg_id`);

--
-- Indexes for table `msg_status`
--
ALTER TABLE `msg_status`
  ADD PRIMARY KEY (`id`),
  ADD KEY `msg_id` (`msg_id`),
  ADD KEY `client_id` (`client_id`);

--
-- Indexes for table `subscribe`
--
ALTER TABLE `subscribe`
  ADD PRIMARY KEY (`id`),
  ADD KEY `client_id` (`client_id`),
  ADD KEY `client_id_2` (`client_id`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `account`
--
ALTER TABLE `account`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=11;
--
-- 使用表AUTO_INCREMENT `clientid`
--
ALTER TABLE `clientid`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;
--
-- 使用表AUTO_INCREMENT `connection`
--
ALTER TABLE `connection`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;
--
-- 使用表AUTO_INCREMENT `msg_info`
--
ALTER TABLE `msg_info`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=255;
--
-- 使用表AUTO_INCREMENT `msg_log`
--
ALTER TABLE `msg_log`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=49;
--
-- 使用表AUTO_INCREMENT `msg_raw`
--
ALTER TABLE `msg_raw`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=255;
--
-- 使用表AUTO_INCREMENT `msg_status`
--
ALTER TABLE `msg_status`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=37;
--
-- 使用表AUTO_INCREMENT `subscribe`
--
ALTER TABLE `subscribe`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=16;COMMIT;
