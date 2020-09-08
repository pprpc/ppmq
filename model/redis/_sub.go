package rsub

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-redis/redis"
	"github.com/pprpc/util/logs"
	"xcthings.com/micro/svc"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQSubscribe"
)

type RSub struct {
	rcli *redis.Client
}
type Sub struct {
	Account      string
	ClientID     string
	Cluster      uint32
	ClusterSubid string
	Qos          uint32
}

/*
// Sub 订阅者信息
type Sub struct {
	Account      string
	ClientID     string
	Cluster      uint32
    ClusterSubid string
    Qos         uint
}
*/

// NewRSub connect redis server
func NewRSub(c svc.RedisConf) (s *RSub) {
	s = new(RSub)
	s.rcli = redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		PoolSize:     c.PoolSize,
		MinIdleConns: c.IdleConn,
	})
	return
}

// CloseRSub close redis sub conn
func CloseRSub(s *RSub) error {
	return s.rcli.Close()
}

// Sub 进行订阅
func (s *RSub) Sub(account, clientid string, topics []*PPMQSubscribe.TopicInfo) (err error) {
	for _, row := range topics {
		key := row.Topic
		value := Sub{
			Account:      account,
			ClientID:     clientid,
			Cluster:      row.Cluster,
			ClusterSubid: row.ClusterSubid,
			Qos:          row.Qos,
		}
		_s, lerr := json.Marshal(value)
		if lerr != nil {
			err = fmt.Errorf("json.Marshal(Sub), %s", lerr)
			return
		}
		s.rcli.Set(fmt.Sprintf("%s-%s", clientid, row.Topic), _s, 0)
		s.rcli.SAdd(key, clientid)
		s.rcli.SAdd(clientid, row.Topic)
		logs.Logger.Debugf("Sub set1, key: %s, v: %s.", fmt.Sprintf("%s-%s", clientid, row.Topic), _s)
		logs.Logger.Debugf("Sub set2, key: %s, v: %v.", key, clientid)
		logs.Logger.Debugf("Sub set3, key: %s, v: %v.", clientid, row.Topic)
	}
	return
}

// UnSub 取消订阅
func (s *RSub) UnSub(clientid string, topics []string) (err error) {
	for _, row := range topics {
		err = s.rcli.SRem(row, clientid).Err()
		if err != nil {
			return
		}
		err = s.rcli.Del(fmt.Sprintf("%s-%s", clientid, row)).Err()
		if err != nil {
			return
		}
		err = s.rcli.SRem(clientid, row).Err()
		if err != nil {
			return
		}
	}
	return
}

// ClearSub 取消所有clientid的订阅
func (s *RSub) ClearSub(clientid string) (err error) {
	topics, lerr := s.rcli.SMembers(clientid).Result()
	if lerr != nil {
		err = fmt.Errorf("s.rcli.Smembers(%s), %s", clientid, lerr)
		return
	}
	for _, row := range topics {
		// 删除 topic 下的 clientid
		err = s.rcli.SRem(row, clientid).Err()
		if err != nil {
			return
		}
		// 删除该topic的订阅详细信息
		err = s.rcli.Del(fmt.Sprintf("%s-%s", clientid, row)).Err()
		if err != nil {
			return
		}
	}
	// 删除 clientid 的所有订阅 topic
	err = s.rcli.Del(clientid).Err()

	return
}

// GetSubByTopic 通过Topic，找到对应的订阅者
func (s *RSub) GetSubByTopic(topic string) (subs []Sub, err error) {
	// 1. 切分 2. 将订阅合并
	// 3. 需要获取订阅者的详细订阅信息
	logs.Logger.Debugf("GetSubByTopic, topic: [%s].", topic)
	var keys []string
	keys, err = getKey(topic)
	if err != nil {
		err = fmt.Errorf("topic: %s, %s", topic, err)
		return
	}
	for _, key := range keys {
		logs.Logger.Debugf("GetSubByTopic, key: %s.", key)
		cids, lerr := s.rcli.SMembers(key).Result()
		if lerr != nil {
			err = fmt.Errorf("rcli.SMembers(%s).Result(), %s", key, lerr)
			return
		}

		for _, cid := range cids {
			k := fmt.Sprintf("%s-%s", cid, key)
			v, e := s.rcli.Get(k).Result()
			if e != nil {
				err = fmt.Errorf("s.rcli.Get(%s), %s", k, e)
				return
			}
			var s Sub
			err = json.Unmarshal([]byte(v), &s)
			if err != nil {
				err = fmt.Errorf("json.Unmarshal(%s), %s", v, err)
				return
			}
			subs = append(subs, s)
		}
	}
	return
	// cids, lerr := s.rcli.SUnion(keys...).Result()
	// if lerr != nil {
	// 	err = fmt.Errorf("SUnion(%v).Result(), %s", keys, lerr)
	// 	return
	// }
	// for _, cid := range cids {
	// 	s.rcli.Get(smt.Sprintf("%s-%s", cid, topic))
	// }
}

// AddOfflineMsgID 增加离线消息
func (s *RSub) AddOfflineMsgID(clientid, msgid string) (err error) {
	// 保存的消息数量,保留的时间
	err = s.rcli.SAdd(fmt.Sprintf("%s-offline", clientid), msgid).Err()
	return
}

// RemoveOfflineMsgID 移除离线消息
func (s *RSub) RemoveOfflineMsgID(clientid, msgid string) (err error) {
	err = s.rcli.SRem(fmt.Sprintf("%s-offline", clientid), msgid).Err()
	return
}

// GetOfflineMsgID 获取离线信息
func (s *RSub) GetOfflineMsgID(clientid string) (msgids []string, err error) {
	key := fmt.Sprintf("%s-offline", clientid)

	msgids, err = s.rcli.SMembers(key).Result()
	if err != nil {
		err = fmt.Errorf("rcli.SMembers(%s).Result(), %s", key, err)
	}
	return
}
func getKey(topic string) (mtop []string, err error) {
	sepArr := strings.Split(topic, "/")
	var _mtop []string
	l := len(sepArr)
	for i, row := range sepArr {
		if i == 0 {
			if row != "" {
				err = fmt.Errorf("idx: %2d, v: %s, format error", i, row)
				return
			}
			continue
		} else if i == 1 {
			if row == "" {
				err = fmt.Errorf("idx: %2d, v: %s, format error", i, row)
				return
			}
			mtop = append(mtop, fmt.Sprintf("/%s/", row))
			_mtop = append(_mtop, fmt.Sprintf("/%s/#", row))
		} else if i+1 < l {
			if row == "" {
				err = fmt.Errorf("idx: %2d, v: %s, format error", i, row)
				return
			}
			mtop = append(mtop, fmt.Sprintf("%s%s/", mtop[i-2], row))
			_mtop = append(_mtop, fmt.Sprintf("%s%s/#", mtop[i-2], row))
		} else {
			if row == "" {
				err = fmt.Errorf("idx: %2d, v: %s, format error", i, row)
				return
			}
			mtop = append(mtop, fmt.Sprintf("%s%s", mtop[i-2], row))
			_mtop = append(_mtop, fmt.Sprintf("%s%s/#", mtop[i-2], row))
		}
	}
	mtop = append(mtop, _mtop...)
	return
}
