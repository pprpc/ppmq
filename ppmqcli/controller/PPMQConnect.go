package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQConnect"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
)

// PPMQConnecter .
type PPMQConnecter struct{}

// ReqHandle PPMQConnect request handle
func (t *PPMQConnecter) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQConnect.Req) (err error) {
	return
}

// RespHandle PPMQConnect response handle
func (t *PPMQConnecter) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQConnect.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle PPMQConnect.
func (t *PPMQConnecter) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQConnect, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
