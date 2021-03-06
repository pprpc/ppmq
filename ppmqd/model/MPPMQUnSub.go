package model

import (
	pm "github.com/pprpc/ppmq/model"
	g "github.com/pprpc/ppmq/ppmqd/common/global"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQUnSub"
)

// MPPMQUnSub PPMQUnSub  
func MPPMQUnSub(req *PPMQUnSub.Req, ci *g.ConnAttr) (resp *PPMQUnSub.Resp, code uint64, err error) {
	for _, topic := range req.Topics {
		d := new(pm.Subscribe)
		d.ClientID = ci.ClientID
		d.Topic = topic

		code, err = d.DeleteTopic()
		if err != nil {
			return
		}
	}
	resp = new(PPMQUnSub.Resp)
	return
}
