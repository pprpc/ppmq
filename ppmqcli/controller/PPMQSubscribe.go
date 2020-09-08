package controller

import (
	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQSubscribe"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
)

// PPMQSubscribeer .
type PPMQSubscribeer struct{}

// ReqHandle PPMQSubscribe request handle
func (t *PPMQSubscribeer) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQSubscribe.Req) (err error) {

	return
}

// RespHandle PPMQSubscribe response handle
func (t *PPMQSubscribeer) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQSubscribe.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle PPMQSubscribe.
func (t *PPMQSubscribeer) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQSubscribe, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}