package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQUnSub"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	g "github.com/pprpc/ppmq/ppmqd/common/global"
	l "github.com/pprpc/ppmq/ppmqd/logic"
)

// PPMQUnSuber .
type PPMQUnSuber struct{}

// ReqHandle PPMQUnSub request handle
func (t *PPMQUnSuber) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQUnSub.Req) (err error) {
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
	var resp *PPMQUnSub.Resp
	resp, code, err = l.LPPMQUnSub(c, pkg, req, ci)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LPPMQUnSub, code: %d, err: %s.", c, code, err)
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

// RespHandle PPMQUnSub response handle
func (t *PPMQUnSuber) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQUnSub.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle PPMQUnSub.
func (t *PPMQUnSuber) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQUnSub, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
