package logic

import (
	"xcthings.com/hjyz/logs"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQUnSub"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"

	errc "xcthings.com/ppmq/common/errorcode"
	g "xcthings.com/ppmq/ppmqd/common/global"
	m "xcthings.com/ppmq/ppmqd/model"
)

// LPPMQUnSub PPMQUnSub Business logic
func LPPMQUnSub(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQUnSub.Req, ci *g.ConnAttr) (resp *PPMQUnSub.Resp, code uint64, err error) {

	resp, code, err = m.MPPMQUnSub(req, ci)
	if err != nil {
		return
	}

	// redis
	if g.PConf.Redis.Addr != "" {
		err = g.TopicCache.UnSub(ci.ClientID, req.Topics)
		if err != nil {
			pkg.Code = errc.DBERROR
			logs.Logger.Errorf("g.TopicCache.Sub(), %s.", err)
		}
	}
	return
}