package model

import (
	"github.com/pprpc/util/common"
	pm "github.com/pprpc/ppmq/model"
	g "github.com/pprpc/ppmq/ppmqd/common/global"
	"github.com/pprpc/ppmq/ppmqd/common/topic"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQSubscribe"
)

// MPPMQSubscribe PPMQSubscribe
func MPPMQSubscribe(req *PPMQSubscribe.Req, ci *g.ConnAttr) (resp *PPMQSubscribe.Resp, code uint64, err error) {
	resp = new(PPMQSubscribe.Resp)

	q := new(pm.Subscribe)
	q.ClientID = ci.ClientID
	_, rows, c, e := q.GetByClientID()
	if e != nil {
		code = c
		err = e
		return
	}
	var topics []string
	for _, row := range rows {
		topics = append(topics, row.Topic)
	}
	curMS := common.GetTimeMs()
	var r []*pm.Subscribe
	var ans []*PPMQSubscribe.TopicAnswer
	for _, row := range req.Topics {
		if topic.Contains(row.Topic, topics) == false {
			_s := new(pm.Subscribe)
			_s.Topic = row.Topic
			_s.Account = ci.Account
			_s.ClientID = ci.ClientID
			_s.Qos = row.Qos
			_s.Cluster = row.Cluster
			_s.ClusterSubid = row.ClusterSubid
			_s.LastTime = curMS
			_s.GlobalSync = 0

			_ans := new(PPMQSubscribe.TopicAnswer)
			_ans.Topic = row.Topic
			_ans.Qos = 1

			r = append(r, _s)
			ans = append(ans, _ans)
		} else {
			_ans := new(PPMQSubscribe.TopicAnswer)
			_ans.Topic = row.Topic
			_ans.Qos = 1
			ans = append(ans, _ans)
		}
	}

	code, err = pm.SubscribeInsertRows(r)
	if err == nil {
		resp.Ans = ans
	}

	return
}
