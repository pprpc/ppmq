package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQUnSub"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
)

// PPMQUnSuber .
type PPMQUnSuber struct{}

// ReqHandle PPMQUnSub request handle
func (t *PPMQUnSuber) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQUnSub.Req) (err error) {

	return
}

// RespHandle PPMQUnSub response handle
func (t *PPMQUnSuber) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQUnSub.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle PPMQUnSub.
func (t *PPMQUnSuber) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQUnSub, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
