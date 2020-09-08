## 项目说明

## 目录说明

目录名称| 说明
---|---
app| 存放start app的相关代码
common/global|存放APP需要的全局变量
commaon/app| 初始化start app所需要的资源
bin| 存放编译代码
conf| 存放相关配置文件
controller| 控制调度
logic|业务逻辑
model|数据结构，数据操作模型和方法

## 接口实现

CmdID|名称|可用|报文方向|说明|开发状态|测试状态
---|---|---|---|---|---|---
11| PPMQGetClientID | true | A-->S / D-->S | 获得唯一ClientID | yes | yes
12| PPMQConnect | true | A-->S / D-->S | 连接建立 | yes | yes
13| PPMQPublish | true | A-->S / D-->S | 发布消息 | yes | yes
14| PPMQPubRec | true | A-->S / D-->S | Qos 2 发布收到(保证交付第一步) | no | no
15| PPMQPubRel | true | A-->S / D-->S | Qos 2 发布释放(保证交付第二步) | no | no
16| PPMQPubComp | true | A-->S / D-->S | QoS 2 消息发布完成(保证交互第三步) | no | no
17| PPMQSubscribe | true | A-->S / D-->S | 客户端订阅请求 | yes | yes
18| PPMQUnSub | true | A-->S / D-->S | 客户端取消订阅请求 | yes | yes
19| PPMQPing | true | A-->S / D-->S | 心跳请求 | yes | yes
20| PPMQDisconnect | true | A-->S / D-->S | 客户端断开连接 | yes | yes
21| PPMQGetSublist | true | A-->S / D-->S | 客户端获取自己的订阅记录 | yes | yes
22| PPMQOAONotify| true | S-->S | 系统根据需要产生在离线消息通知 | pb | no
23| PPMQEXChangeMsg | true | S-->S | 系统用于进行消息交互使用的接口 | yes | yes


## 关于PPMQ,LocalMQ的说明

 * PPMQD: 用于作为对外的消息订阅服务器，可以看成是IOT Gateway在使用。IPC设备都是被分配到该服务器上，配置的端口需要外部环境能够访问
 * LocalMQD: 主要用于内部微服务之间的信息流转，用于微服务之间的接口解耦或者耗时调用的一些异步操作,端口不用对外部环境开放;例如:
   * user-ms 在接口`SendVcode`，会将信息投递到LocalMQ 的Topic: `/notify/vcode/`上.
   * sendvcode-ms 该微服务会连接LocalMQ,订阅topic:`/notify/vcode/`; 在接收到消息后触发发送验证码的服务.
 * 为了避免两个MQ的服务相互影响，分别使用独立的数据库；同时localMQD,的配置参数`auth_desc`建议设置为1,同时`check_topic`设置为0.
 