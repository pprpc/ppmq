## 设计

 * clientid+topic 作为一个Key,value 为订阅信息存放在Redis中
 * topic 作为一个key, value 为set 类型,其中的值为 clientid；这样可以通过topic 找到所有相关订阅者
 * 查找的时候通过 topic 进行切分，分别查找 `topic` 对应的 clientid, 再找到 clientid对应该topic的订阅信息.
 * key: clientid, value 是一个set集合，用于保存该clientid所有的订阅topic 

## 数据结构

 * 两个set集合，一个key是topic(value 是订阅该topic的所有clientid),一个key 是 clientid(value 是记录该clientid订阅的所有topic)
 * 一个普通key-value数据,key: clientid+topic; value 是具体的订阅信息.

## 需要发送的离线消息记录

 * key: clientid+offline, set集合; value: msgid