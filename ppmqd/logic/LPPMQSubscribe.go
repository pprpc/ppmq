package logic

import (
	"fmt"

	"xcthings.com/hjyz/logs"
	errc "xcthings.com/ppmq/common/errorcode"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQSubscribe"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"

	g "xcthings.com/ppmq/ppmqd/common/global"
	t "xcthings.com/ppmq/ppmqd/common/topic"
	m "xcthings.com/ppmq/ppmqd/model"
)

// LPPMQSubscribe PPMQSubscribe Business logic
func LPPMQSubscribe(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQSubscribe.Req, ci *g.ConnAttr) (resp *PPMQSubscribe.Resp, code uint64, err error) {
	if len(req.Topics) == 0 {
		code = errc.ParameterError
		err = fmt.Errorf("Incomplete parameters: Topics, Password")
		return
	}
	for _, v := range req.Topics {
		if err = t.CheckSubTopic(v.Topic); err != nil {
			code = errc.ParameterIllegal
			return
		}
	}
	//
	if g.PConf.Auth == 2 && g.PConf.CheckTopic == 1 {
		var topics []string
		for _, v := range req.Topics {
			topics = append(topics, v.Topic)
		}
		ans, lerr := m.CheckSub(ci.Account, topics)
		if lerr != nil {
			logs.Logger.Warnf("LPPMQSubscribe, CheckSub, %s.", lerr)
			code = errc.ParameterError
			err = fmt.Errorf("m.CheckSub(), %s", lerr)
			return
		}
		for _, v := range ans {
			if v.Code != 0 {
				code = errc.SUBREJECT
				err = fmt.Errorf("topic: %s, sub reject", v.Topic)
				return
			}
		}
	}

	// redis
	resp, code, err = m.MPPMQSubscribe(req, ci)
	if err != nil {
		return
	}
	if g.PConf.Redis.Addr != "" {
		err = g.TopicCache.Sub(ci.Account, ci.ClientID, req.Topics)
		if err != nil {
			pkg.Code = errc.DBERROR
			logs.Logger.Errorf("g.TopicCache.Sub(), %s.", err)
		}
	}
	return
}
