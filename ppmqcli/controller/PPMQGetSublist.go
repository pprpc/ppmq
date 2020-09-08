package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQGetSublist"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
)

// PPMQGetSublister .
type PPMQGetSublister struct{}

// ReqHandle PPMQGetSublist request handle
func (t *PPMQGetSublister) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQGetSublist.Req) (err error) {

	return
}

// RespHandle PPMQGetSublist response handle
func (t *PPMQGetSublister) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQGetSublist.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	for _, v := range resp.Topics {
		logs.Logger.Debugf("Topic: %s.", v.Topic)
	}
	return
}

// DestructHandle PPMQGetSublist.
func (t *PPMQGetSublister) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQGetSublist, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
