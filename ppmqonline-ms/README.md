## 项目说明

 * `ftconnrelay` 用于连接的转发

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
10000| SLCtrl | true | A-->S | 开锁| yes | no
10001| SLPwdAdd | true | A-->S | 添加密码 | yes | no
10002| SLPwdEdit | true | A-->S | 编辑密码 | yes | no
10003| SLPwdDel | true | A-->S | 删除密码 | yes | no
10004| SLConfigGet | true | A-->S | 获取锁的配置信息 | yes | no
10005| SLUnlockList | true | A-->S | 获取开锁记录 | no | no



## TODO
