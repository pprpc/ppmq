## PPMQ设计说明

<!-- TOC -->

- [PPMQ设计说明](#ppmq设计说明)
- [1. 数据库设计](#1-数据库设计)
  - [1.1. account](#11-account)
  - [1.2. clientid](#12-clientid)
  - [1.3. connection](#13-connection)
  - [1.4. msg_info](#14-msg_info)
  - [1.5. msg_log](#15-msg_log)
  - [1.6. msg_raw](#16-msg_raw)
  - [1.7. msg_status](#17-msg_status)
  - [1.7. subscribe](#17-subscribe)
- [2. 配置参数说明](#2-配置参数说明)
- [3. 访问PPMQD的微服务定义](#3-访问ppmqd的微服务定义)
  - [3.1. p2p-ms](#31-p2p-ms)
  - [3.1. relay-ms](#31-relay-ms)
- [4. 访问localmqd的微服务说明](#4-访问localmqd的微服务说明)
  - [4.1. sendvcode-ms](#41-sendvcode-ms)
- [5. 关于PPMQD,Localmqd使用说明](#5-关于ppmqdlocalmqd使用说明)
- [6. 资源ID标号](#6-资源id标号)

<!-- /TOC -->

## 1. 数据库设计

### 1.1. account

> 说明
 * 该表主要用于存放连接PPMQ的账户信息.
 * 如果配置参数`auth=1`的时候，微服务使用该表内的参数认证.
   * 目前`localmq.account`没有使用`user_id`进行认证，都是使用`account`的字符串账号进行的认证.
 * PPMQ支持连接外部的认证服务: `authdevice`,`authuser`; 也支持微服务内验证，具体使用那种验证方式在配置文件中设定
   * `authdevice-ms`, 通过 `device.device_id` 进行认证
   * `authuser-ms`, 通过`user.user_baseinfo` 进行认证

> 数据表定义
~~~
CREATE TABLE `account` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `account` varchar(64) NOT NULL COMMENT '账户信息',
  `password` varchar(256) NOT NULL COMMENT '账户对应密码，md5(text password)'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息系统： 账户信息';

~~~

### 1.2. clientid

>说明
 * 每一个用户(设备)连接到PPMQ需要确定一个唯一的ClientID,该ClientID全局唯一，后续通过ClientID进行消息投递，连接区分等.
 * 该表主要记录已经分配的ClientID，通过接口:`PPMQGetClientID`调用.
 * 实际使用中可以不用调用，因为我们的设备ID已经全局唯一，可以直接使用用户ID或者设备ID作为连接PPMQ的ClientID来使用.
 * 更多的ClientID的约定与`MQTT`协议的ClientID一致.

> 数据表定义
~~~
CREATE TABLE `clientid` (
  `id` int(11) NOT NULL,
  `client_id` varchar(64) NOT NULL COMMENT '全局唯一，确保相同的帐号+硬件得到的ClientId是一致',
  `account` varchar(64) NOT NULL COMMENT 'user_type = 1,这里是device_id; user_type = 2, 这里是用户电话号码',
  `hw_feature` varchar(256) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息: 记录分配的client_id';

~~~

**注意:**
 * `hw_feature` 设备的硬件特性描述，用于平台生成ClientID的要素之一.
 * 平台如何生成ClientID参见代码: `getClientID` 方法

### 1.3. connection

> 说明
 * 每个设备或者用户连接到PPMQD后，都会在该数据库中存在一条记录.
 * 每个ClientID存在一条记录.
 * 设备/用户是否在线可以通过查阅这张表获得相关的连接信息.
 * 所有`will_xxxx`字段表示遗嘱消息，同MQTT的遗嘱消息定义一致(我们没有实现遗嘱消息的相关功能，因为在我们应用场景中不需要该特性).


>数据表定义
~~~


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
  `is_online` tinyint(1) NOT NULL DEFAULT 0 COMMENT '1, 离线；2, 在线; ',
  `is_sleep` tinyint(2) NOT NULL COMMENT '0, 正常心跳： 1, 睡眠心跳',
  `last_time` bigint(20) NOT NULL COMMENT '最后更新状态的时间',
  `global_sync` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0,1,等待同步; 2,同步中; 3,同步完成 ',
  `account` varchar(80) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息: 连接';

~~~

### 1.4. msg_info

>说明
 * 存放具体消息的摘要
 * 该数据表PPMQD会定期清理

> 数据表定义
~~~
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
~~~


### 1.5. msg_log

>说明
 * 消息日志表
 * 消息在投递给消费者的时候进行记录，用于后期进行消息投递的耗时分析
 * 该数据表由PPMQD维护，会定期清空历史的无效数据

> 数据表定义
~~~

CREATE TABLE `msg_log` (
  `id` int(11) NOT NULL,
  `msg_id` varchar(64) NOT NULL,
  `account` varchar(64) NOT NULL COMMENT '消费者ID, msg_account.account',
  `client_id` varchar(64) NOT NULL,
  `server_id` varchar(255) NOT NULL,
  `status` tinyint(2) NOT NULL COMMENT 'QOS1: 1,send; 2, pubAns; QOS2: 1, pub; 2, pubrec；3, pubrel；4, pubcomp',
  `create_time` bigint(14) NOT NULL COMMENT '创建消息的时间单位毫秒'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息：投递日志，用于后期分析';

~~~

### 1.6. msg_raw

>说明
 * 存放原始消息
 * 在进行消息的离线投放的时候需要使用该数据表
 * 该数据表由PPMQD维护，会定期清空历史的无效数据

> 数据表定义
~~~
CREATE TABLE `msg_raw` (
  `id` int(11) NOT NULL,
  `msg_id` varchar(64) NOT NULL,
  `msg_payload` longblob NOT NULL COMMENT '除开固定和可变报头的消息载荷'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息：存放当前未完成流程的消息';
~~~

### 1.7. msg_status

>说明
 * 记录消息状态
 * 记录消费者消费消息的时候的状态，对于没有投递成功的消息在消费者上线后检查该表，并根据设置的参数进行信息的再次投递
 * 该数据表由PPMQD维护，会定期清空历史的无效数据

> 数据表定义
~~~
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

~~~


### 1.7. subscribe

>说明
 * 订阅数据表
 * 记录ClientID的订阅记录，订阅参数

> 数据表定义
~~~
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

~~~
## 2. 配置参数说明

 * 每个微服务都有这几个独立的配置文件.

> listen.json
~~~
[{
    "uri": "tcp://0.0.0.0:53"
    ,"read_timeout": 270
    ,"res_id": 1
},{
    "uri": "udp://0.0.0.0:53"
    ,"read_timeout": 35
    ,"res_id": 2
},{
    "uri": "tcp://0.0.0.0:1053"
    ,"read_timeout": 270
    ,"res_id": 1
},{
    "uri": "udp://0.0.0.0:1053"
    ,"read_timeout": 35
    ,"res_id": 2
},{
    "uri": "tls://0.0.0.0:8053"
    ,"read_timeout": 270
    ,"tls_crt": "../conf/tls/server.crt"
    ,"tls_key": "../conf/tls/server.key"
    ,"res_id": 0
}]
~~~

字段|说明
---|---
uri | 设置微服务监听的URL: tcp://0.0.0.0:53; 表示监听TCP的53端口，在机器的所有网口上
read_timeout | 连接读超时时间，单位秒; 如果连接读取数据超时，则服务端会断开连接,释放相关资源
res_id | 资源ID; 注册到ETCD的服务资源号，该资源ID定义在文件:`pprpc_service.md`
tls_crt | 如果监听的TLS，这里设置证书
tls_key | 如果监听的TLS，这里设置证书

> log.json
~~~
{
    "file": "/var/log/ppmq/ppmqd.log"
    ,"max_size": 2
    ,"max_backups": 5
    ,"max_age": 0
    ,"caller": true
    ,"level": -1
}
~~~

字段|说明
---|---
file | 日志存放路径
max_size | 每个日志文件的大小，单位MB
max_backups | 最多几个日志文件
max_age | 最大日志存放几天
caller | 日志文件是否显示源码调用行数和调用文件名称
level | 日志等级: -1, debug; 0, Info; 1, Warn; 2, Error; 5, Fatal

> public.json
~~~
{
    "report_interval": 60
    ,"max_go": 100000
    ,"admin_prof" : true
    ,"admin_port": 20053
    ,"run_go": true
}
~~~

字段|说明
---|---
report_interval | 自动记录当前微服务的运行状态信息的间隔时间，单位秒
max_go | 最大允许的GO程，目前没有使用
admin_prof | 是否启用Pprof性能分析; true: 启用; false: 不启用
admin_port | 性能分析使用的端口
run_go | RPC响应的时候是否启用GO程; true: 启用： false: 不启用

> private.json
~~~
{
    "max_session": 1000000
    ,"ppmq":{
        "online_notify": true
        ,"offline_notify": true
        ,"_notify_topic_prefix": "在离线信息的前缀: online and offline = oao; /oao/clientid"
        ,"notify_topic_prefix": "/oao/"
        ,"udp_hbsec": 28
        ,"tcp_hbsec": 150
        ,"max_sessions": 12345
        ,"_mode": "single: 单机版; cluster: 集群版本"
        ,"mode": "single"
        ,"cluster_signkey": "xxxxxxxxxxxxxxxx"
        ,"qos": true
        ,"udp_resp_timeoutms": 600
        ,"offlinemsg_timeoutms_desc": "当前时间先前推该时间作为开始时间"
        ,"offlinemsg_timeoutms": 3600000
        ,"offlinemsg_endms_desc": "当前时间往前推该时间作为结束时间"
        ,"offlinemsg_endms": 10000
        ,"offlinemsg_send_sleepms_desc": "Connect成功后多久开始发送离线消息"
        ,"offlinemsg_send_sleepms": 1500
        ,"offlinemsg_send_intervalms_desc": "发送离线消息时候的间隔时间"
        ,"offlinemsg_send_intervalms": 100
        ,"_sso_enable": "true: 开启账号单点登录; false: 允许多次登录"
        ,"sso_enable": true
        ,"tempdb_expiredms": 604800000
        ,"tempdb_clear_intervalsec": 3600
    }
    ,"redis":{
        "addr": "127.0.0.1:6379"
        ,"password":""
        ,"db": 0
        ,"pool_size": 5
        ,"idle_conn": 5
    }
    ,"_clear_sameid":"去掉了sso_enable功能，使用该功能解决相同ID上线问题"
    ,"clear_sameid": {
        "clear": true
        ,"white_account": ["PPdevice2"]
    }
    ,"device_prefix": ""
    ,"auth_desc": "1, local auth; 2, auth micro service"
    ,"auth": 1
    ,"check_topic_desc": "auth=2的时候: 0 不检查topic; 1 检查topic"
    ,"check_topic": 1
    ,"micros":[{
        "name": "authuser"
        ,"uris": ["tcp://127.0.0.1:6004"]
    },{
        "name": "authdevice"
        ,"uris": ["tcp://192.168.6.217:6000"]
    },{
        "name": "authdevice"
        ,"uris": ["tcp://192.168.6.217:6000"]
    }]
}
~~~

字段|说明
---|---
max_session | 当前PPMQD最大允许的连接数
online_notify | 是否开启设备上线通知; 如果开启，设备/用户上线会自动产生一条消息
offline_notify | 是否开启设备的离线通知; 如果开启，设备/用户离线会自动产生一条消息
notify_topic_prefix | 在离线通知产生的topic前缀; 需要设备在离线的业务订阅这个topic就可以知道设备的在离线信息
udp_hbsec | 如果设备通过UDP连接，则需要维持UDP心跳的间隔时间
tcp_hbsec | 如果设备通过TCP连接，则需要维持TCP心跳的间隔时间
max_sessions | 忽略
mode | PPMQD的运行模式: single(单机版); cluster(集群版本)
cluster_signkey | PPMQD之间进行数据交互时使用的签名
qos | 是否启用QOS; 程序只实现了MQTT的Qos=1的清空，目前该字段忽略
udp_resp_timeoutms | 这里设置的应答超时时间，表示向消费者投递消息的时候如果超时没有回应，则会再次重新发送一次消息
offlinemsg_timeoutms | 如果消费者消费离线消息，这里设置可以消防多场时间范围内的离线消息
offlinemsg_endms | 如果消费者消费离线消息，这里设置可以消防多场时间范围内的离线消息
offlinemsg_send_sleepms | Connect成功后多久开始发送离线消息
offlinemsg_send_intervalms | 发送离线消息时候的间隔时间，如果下发离线信息过快可能导致消费者处理异常
sso_enable | 是否开启单点登录; true: 开启账号单点登录; false: 允许多次登录; 通过ClientID确定; 暂时启用，使用`clear_sameid`配置
tempdb_expiredms | 清理数据表: `msg_info`,`msg_log`,`msg_raw`,`msg_status`的超时时间
tempdb_clear_intervalsec | 清理数据表的检测周期
redis | 参见Redis设置;主要用于PPMQD存放订阅信息
clear_sameid | 去掉了sso_enable功能，使用该功能解决相同ID上线问题
device_prefix | 设备ID的前缀，目前参数弃用
auth | 账号认证方式: 1, 内部认证(检索数据表:`account`); 2, 通过微服务认证(authdevice,authuser);
check_topic | 是否进行topic权限检查(检查是否由权限产生或者订阅相关的topic); 0 不检查topic; 1 检查topic
micros | 需要使用的微服务;目前主要使用`authdevice`,`authuser`；用于进行账号验证和Topic验证.

>Redis

字段|说明
---|---
addr | Redis服务器
password | 密码
db | 数据库
pool_size | 连接池大小
idle_conn | 空闲连接大小

>clear_sameid

字段|说明
---|---
clear | 是否开启: true 开启; false 关闭
white_account | 允许的白名单，在白名单中的账号可以存在多个连接


## 3. 访问PPMQD的微服务定义

 * 相关微服务访问ppmqd的配置参数都在各个微服务的配置文件`ppmqcli.json`文件中定义的

### 3.1. p2p-ms

 * 由于建立连接需要和设备进行沟通，该微服务需要使用到ppmqd

>ppmqcli.json
~~~
[{
    "class": "ppmqd"
    ,"url": "tcp://localhost:8053"
    ,"account": "PPftconnp2p"
    ,"password": "2399ee1cd21036b74ec4c9961d1d1cww"
    ,"hw_feature": "52:54:00:7f:ca:df1"
}]
~~~

字段|说明
---|---
class| 连接的PPMQD的服务类型: ppmqd, localmqd
url| 微服务其的地址，目前通过ETCD进行服务发现的时候会自动找到连接地址，该值会被忽略
account| 连接ppmqd的认证账号，该账号需要在`device.device_id`中存在，否则该微服务连接PPMQD会失败
password| 账号对应的密码
hw_feature| p2p-ms连接ppmqd的获取`clientid`的时候的要素之一，需要确保不同微服务使用相同的`account`连接ppmqd的时候`hw_feature`不能相同 

### 3.1. relay-ms
 
  * 由于建立连接需要和设备进行沟通，该微服务需要使用到ppmqd
  * 其余配置信息说明同`p2p-ms`一致
 

## 4. 访问localmqd的微服务说明

 * 相关微服务访问localmqd的配置参数都在各个微服务的配置文件`ppmqcli.json`文件中定义的

### 4.1. sendvcode-ms

 * 当需要发送验证码的时候需要解耦为异步的方式

>ppmqcli.json
~~~
[{
    "class": "localmqd"
    ,"url": "tcp://localhost:8053"
    ,"account": "FDlocalmqd-sendvcode"
    ,"password": "mli3uBPWlZi4mli3uBPWlZi4"
    ,"hw_feature": "52:54:00:7f:ca:dfsendvcode"
    ,"topic_prefix": "/notify/vcode/"
}]
~~~

字段|说明
---|---
class| 连接的PPMQD的服务类型: ppmqd, localmqd
url| 微服务其的地址，目前通过ETCD进行服务发现的时候会自动找到连接地址，该值会被忽略
account| 连接localmqd的认证账号，该账号需要在`localmq.account`中存在，否则该微服务连接localmq会失败
password| 账号对应的密码
hw_feature| sendvcode-ms连接ppmqd的获取`clientid`的时候的要素之一，需要确保不同微服务使用相同的`account`连接ppmqd的时候`hw_feature`不能相同 
topic_prefix | 订阅的`topic`

**注意：**
 * 其余微服务需要访问`localmqd`的时候配置参数只是在`account`,`password`,`hw_feature` 三个参数进行调整即可.


## 5. 关于PPMQD,Localmqd使用说明

 * PPMQD: 用于作为对外的消息订阅服务器，可以看成是IOT Gateway在使用。IPC设备都是被分配到该服务器上，配置的端口需要外部环境能够访问
 * LocalMQD: 主要用于内部微服务之间的信息流转，用于微服务之间的接口解耦或者耗时调用的一些异步操作,端口不用对外部环境开放;例如:
   * user-ms 在接口`SendVcode`，会将信息投递到LocalMQ 的Topic: `/notify/vcode/`上.
   * sendvcode-ms 该微服务会连接LocalMQ,订阅topic:`/notify/vcode/`; 在接收到消息后触发发送验证码的服务.
 * 为了避免两个MQ的服务相互影响，分别使用独立的数据库；同时localMQD,的配置参数`auth_desc`建议设置为1,同时`check_topic`设置为0.

## 6. 资源ID标号

资源编号(res_srv)|微服务名称
---|---
1|ppmqd(TCP)
2|ppmqd(UDP)
3|ppmqd(mqtt)
4|apigw-mob(grpc)
5|apigw-mob(http)
6|apigw-mob(pprpc)
8|ftconnnat
9|ftconnrelay
10|ftconnp2p
11|ftlives
13|glbs(tcp)
14|glbs(udp)
104|apigw-mob(grpctls)
105|apigw-mob(https)
106|apigw-mob(pprpc_tls)
