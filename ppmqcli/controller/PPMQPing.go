package controller

import (
	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPing"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
)

// PPMQPinger .
type PPMQPinger struct{}

// ReqHandle PPMQPing request handle
func (t *PPMQPinger) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQPing.Req) (err error) {

	return
}

// RespHandle PPMQPing response handle
func (t *PPMQPinger) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQPing.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	if resp.OutsideIpaddr != "" {
		logs.Logger.Debugf("WanIpaddr: %s.", resp.OutsideIpaddr)
	}
	return
}

// DestructHandle PPMQPing.
func (t *PPMQPinger) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQPing, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
