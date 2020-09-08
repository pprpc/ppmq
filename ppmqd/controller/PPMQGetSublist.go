package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQGetSublist"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	g "github.com/pprpc/ppmq/ppmqd/common/global"
	l "github.com/pprpc/ppmq/ppmqd/logic"
)

// PPMQGetSublister .
type PPMQGetSublister struct{}

// ReqHandle PPMQGetSublist request handle
func (t *PPMQGetSublister) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQGetSublist.Req) (err error) {
	var _t interface{}
	_t, err = checkConn(c, pkg)
	if err != nil {
		_, err = pprpc.WriteResp(c, pkg, nil)
		if err != nil {
			logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
		}
		logs.Logger.Errorf("checkConn(), error: %s.", err)
		return
	}
	ci := _t.(*g.ConnAttr)

	var code uint64
	var resp *PPMQGetSublist.Resp
	resp, code, err = l.LPPMQGetSublist(c, pkg, req, ci.ClientID)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LPPMQGetSublist, code: %d, err: %s.", c, code, err)
	}

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}

	return
}

// RespHandle PPMQGetSublist response handle
func (t *PPMQGetSublister) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQGetSublist.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle PPMQGetSublist.
func (t *PPMQGetSublister) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQGetSublist, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
