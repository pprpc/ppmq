package controller

import (
	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQDisconnect"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
)

// PPMQDisconnecter .
type PPMQDisconnecter struct{}

// ReqHandle PPMQDisconnect request handle
func (t *PPMQDisconnecter) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQDisconnect.Req) (err error) {
	return
}

// RespHandle PPMQDisconnect response handle
func (t *PPMQDisconnecter) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQDisconnect.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle PPMQDisconnect.
func (t *PPMQDisconnecter) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQDisconnect, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
	c.Close()
}
