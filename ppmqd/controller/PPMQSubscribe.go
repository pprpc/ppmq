package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQSubscribe"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	g "github.com/pprpc/ppmq/ppmqd/common/global"
	l "github.com/pprpc/ppmq/ppmqd/logic"
)

// PPMQSubscribeer .
type PPMQSubscribeer struct{}

// ReqHandle PPMQSubscribe request handle
func (t *PPMQSubscribeer) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQSubscribe.Req) (err error) {
	var _t interface{}
	_t, err = checkConn(c, pkg)
	if err != nil {
		logs.Logger.Errorf("checkConn(), error: %s.", err)
		_, err = pprpc.WriteResp(c, pkg, nil)
		if err != nil {
			logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
		}
		return
	}
	ci := _t.(*g.ConnAttr)

	var code uint64
	var resp *PPMQSubscribe.Resp
	resp, code, err = l.LPPMQSubscribe(c, pkg, req, ci)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LPPMQSubscribe, code: %d, err: %s.", c, code, err)
	}
	ci.RunSub = true

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}
	//if err != nil || pkg.Code != 0 {
	// 	return
	//}
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
