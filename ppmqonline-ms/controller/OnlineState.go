package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"github.com/pprpc/ppmq/protoc/ppmqonline/OnlineState"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	l "github.com/pprpc/ppmq/ppmqonline-ms/logic"
)

// OnlineStateer .
type OnlineStateer struct{}

// ReqHandle OnlineState request handle
func (t *OnlineStateer) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *OnlineState.Req) (err error) {
	var code uint64
	var resp *OnlineState.Resp
	resp, code, err = l.LOnlineState(c, pkg, req)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LOnlineState, code: %d, err: %s.", c, code, err)
	}

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}
	//if err != nil || pkg.Code != 0 {
	// 	return
	//}
	return
}

// RespHandle OnlineState response handle
func (t *OnlineStateer) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *OnlineState.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle OnlineState.
func (t *OnlineStateer) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, OnlineState, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
