package controller

import (
	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	mqc "xcthings.com/ppmq/common"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPing"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"

	g "xcthings.com/ppmq/ppmqd/common/global"

	m "xcthings.com/ppmq/ppmqd/model"
)

// PPMQPinger .
type PPMQPinger struct{}

// ReqHandle PPMQPing request handle
func (t *PPMQPinger) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQPing.Req) (err error) {
	var _t interface{}
	_t, err = checkConn(c, pkg)
	if err != nil {
		logs.Logger.Errorf("checkConn(), Code: %d, error: %s.", pkg.Code, err)
		_, err = pprpc.WriteResp(c, pkg, nil)
		if err != nil {
			logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
		}
		return
	}
	ci := _t.(*g.ConnAttr)
	oldSleep := ci.ISSleep
	if req.IsSleep {
		ci.ISSleep = mqc.PINGSLEEP
	} else {
		ci.ISSleep = mqc.PINGNORMAL
	}
	ci.LastTime = common.GetTimeMs()

	resp := new(PPMQPing.Resp)
	resp.OutsideIpaddr = common.GetIPAddr(c.RemoteAddr())

	if ci.RemoteIpaddr != resp.OutsideIpaddr {
		ci.RemoteIpaddr = resp.OutsideIpaddr
	} else {
		resp.OutsideIpaddr = ""
	}

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}
	if oldSleep != ci.ISSleep {
		_, err = m.UpdateSleep(ci.ClientID, ci.ISSleep)
		if err != nil {
			logs.Logger.Warnf("client_id: %s, sleep: %d, UpdateSleep(), %s.",
				ci.ClientID, ci.ISSleep, err)
		}
		go onlineNotify(ci)
	}
	return
}

// RespHandle PPMQPing response handle
func (t *PPMQPinger) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQPing.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle PPMQPing.
func (t *PPMQPinger) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQPing, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
