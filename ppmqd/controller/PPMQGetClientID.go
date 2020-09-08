package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQGetClientID"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	l "github.com/pprpc/ppmq/ppmqd/logic"
)

// PPMQGetClientIDer .
type PPMQGetClientIDer struct{}

// ReqHandle PPMQGetClientID request handle
func (t *PPMQGetClientIDer) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQGetClientID.Req) (err error) {
	var code uint64
	var resp *PPMQGetClientID.Resp
	resp, code, err = l.LPPMQGetClientID(c, pkg, req)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LPPMQGetClientID(), code: %d, err: %s.", c, code, err)
	}

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}
	logs.Logger.Debugf("%s, %s, c.WriteResp.seqID: %d.", c, pkg.CmdName, pkg.CmdSeq)
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
