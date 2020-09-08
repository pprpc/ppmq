package controller

import (
	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQGetClientID"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
)

// PPMQGetClientIDer .
type PPMQGetClientIDer struct{}

// ReqHandle PPMQGetClientID request handle
func (t *PPMQGetClientIDer) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQGetClientID.Req) (err error) {
	return
}

// RespHandle PPMQGetClientID response handle
func (t *PPMQGetClientIDer) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQGetClientID.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	logs.Logger.Debugf("%s, ClientID: [%s].", c, resp.ClientId)
	return
}

// DestructHandle PPMQGetClientID.
func (t *PPMQGetClientIDer) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQGetClientID, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
