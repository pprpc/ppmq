package model

import (
	"fmt"

	"xcthings.com/hjyz/common"
	pm "xcthings.com/ppmq/model"
	g "xcthings.com/ppmq/ppmqd/common/global"
	"xcthings.com/ppmq/ppmqd/common/topic"
)

// Sub .
type Sub struct {
	Account      string
	ClientID     string
	Cluster      uint32
	ClusterSubid string
	ServerID     string
}

// Subers .
type Subers struct {
	Subs []Sub
}

// NewSubers .
func NewSubers() (c *Subers) {
	c = new(Subers)
	return
}

// Add .
func (c *Subers) Add(sub Sub) {
	for _, v := range c.Subs {
		if v.ClientID == sub.ClientID {
			return
		}
	}
	c.Subs = append(c.Subs, sub)
}

// Del .
func (c *Subers) Del(k string) {
	for i, v := range c.Subs {
		if v.ClientID == k {
			if i == 0 {
				c.Subs = c.Subs[i+1:]
			} else if i+1 < len(c.Subs) {
				c.Subs = append(c.Subs[:i], c.Subs[i+1:]...)
			} else {
				c.Subs = c.Subs[:i]
			}
		}
	}
}

// GetCustomerBYTopic .
func GetCustomerBYTopic(t, msgid, srcmsgid string) (subs *Subers, err error) {
	subs = NewSubers()

	if g.PConf.Redis.Addr != "" {
		s, lerr := g.TopicCache.GetSubByTopic(t)
		if lerr != nil {
			err = fmt.Errorf("g.TopicCache.GetSubByTopic(%s), %s", t, lerr)
			return
		}
		for _, v := range s {
			var sub Sub
			sub.ClientID = v.ClientID
			sub.Cluster = v.Cluster
			sub.ClusterSubid = v.ClusterSubid
			sub.Account = v.Account
			subs.Add(sub)
		}
	} else {
		ms := new(pm.Subscribe)
		rows, e := pm.Orm.Rows(ms)
		if e != nil {
			err = e
			return
		}
		defer rows.Close()
		var mcode uint32

		for rows.Next() {
			err = rows.Scan(ms)
			if ms.Topic == "" {
				continue
			}
			mcode = topic.MatchV2(t, ms.Topic)
			//logs.Logger.Debugf("msgid: %s, t: %s, topic: %s, mcode: %d, cid: %s.", msgid, t, ms.Topic, mcode, ms.ClientID)
			if mcode == 2 || mcode == 3 {
				var sub Sub
				sub.ClientID = ms.ClientID
				sub.Cluster = ms.Cluster
				sub.ClusterSubid = ms.ClusterSubid
				sub.Account = ms.Account
				subs.Add(sub)
			}
		}
	}

	si := pm.Orm.NewSession()
	defer si.Close()
	err = si.Begin()
	cusMS := common.GetTimeMs()
	for _, v := range subs.Subs {

		_t := new(pm.MsgStatus)
		_t.MsgID = msgid
		_t.SrcMsgid = srcmsgid
		_t.Account = v.Account
		_t.ClientID = v.ClientID
		_t.ServerID = v.ServerID //g.Conf.Public.ServerID
		_t.Dup = 0
		_t.Qos = 1
		_t.Status = 1
		_t.CreateTime = cusMS

		_, err = si.Insert(_t)
		if err != nil {
			si.Rollback()
			return
		}
	}
	err = si.Commit()
	if err != nil {
		si.Rollback()
	}
	return
}

// AddMsgLog .
func AddMsgLog(d []Sub, msgID string) (err error) {
	var rows []*pm.MsgLog
	cusMS := common.GetTimeMs()
	for _, row := range d {
		_t := new(pm.MsgLog)
		_t.MsgID = msgID
		_t.Account = row.Account
		_t.ClientID = row.ClientID
		_t.ServerID = g.Conf.Public.ServerID
		_t.Status = 1
		_t.CreateTime = cusMS

		rows = append(rows, _t)
	}
	_, err = pm.MsgLogInsertRows(rows)
	return
}
