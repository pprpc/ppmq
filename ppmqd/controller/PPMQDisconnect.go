package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQDisconnect"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
)

// PPMQDisconnecter .
type PPMQDisconnecter struct{}

// ReqHandle PPMQDisconnect request handle
func (t *PPMQDisconnecter) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQDisconnect.Req) (err error) {
	var resp *PPMQDisconnect.Resp
	resp = new(PPMQDisconnect.Resp)

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}
	logs.Logger.Warnf("%s, PPMQDisconnect.", c)
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
