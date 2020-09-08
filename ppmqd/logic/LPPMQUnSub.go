package logic

import (
	"github.com/pprpc/util/logs"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQUnSub"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	errc "github.com/pprpc/ppmq/common/errorcode"
	g "github.com/pprpc/ppmq/ppmqd/common/global"
	m "github.com/pprpc/ppmq/ppmqd/model"
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
