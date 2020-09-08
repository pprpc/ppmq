package model

import (
	pm "github.com/pprpc/ppmq/model"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQGetSublist"
)

// MPPMQGetSublist PPMQGetSublist  
func MPPMQGetSublist(req *PPMQGetSublist.Req, clientID string) (resp *PPMQGetSublist.Resp, code uint64, err error) {
	g := new(pm.Subscribe)
	g.ClientID = clientID
	t, rows, c, e := g.GetByClientID()
	if e != nil {
		code = c
		err = e
		return
	}
	resp = new(PPMQGetSublist.Resp)
	for _, row := range rows {
		_s := new(PPMQGetSublist.TopicInfo)
		_s.Topic = row.Topic
		_s.Qos = row.Qos
		_s.Cluster = row.Cluster
		_s.ClusterSubid = row.ClusterSubid

		resp.Topics = append(resp.Topics, _s)
	}
	resp.Total = t
	return
}
