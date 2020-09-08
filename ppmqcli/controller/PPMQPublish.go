package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	clic "github.com/pprpc/ppmq/ppmqcli/common"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQPublish"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
)

// PPMQPublisher .
type PPMQPublisher struct{}

// ReqHandle PPMQPublish request handle
func (t *PPMQPublisher) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQPublish.Req) (err error) {
	resp := new(PPMQPublish.Resp)
	resp.MsgId = req.MsgId
	resp.GMsgid = ""
	logs.Logger.Debugf("MsgID: %s, Topic: %s , Cmdid: %d.", req.MsgId, req.Topic, req.Cmdid)

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}

	//logs.Logger.Debugf("Body: [%s].", req.MsgBody)
	_t, e := c.GetAttr()
	if e != nil {
		logs.Logger.Warnf("%s, not find set conn attr.", c)
		return
	}
	ci := _t.(*clic.PpmqConn)
	if ci.PubCB == nil {
		logs.Logger.Warnf("%s, not set PubCB.", c)
		return
	}
	go ci.PubCB(pkg, req)
	return
}

// RespHandle PPMQPublish response handle
func (t *PPMQPublisher) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQPublish.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	logs.Logger.Debugf("Msgid: %s, GMsgid: %s.", resp.MsgId, resp.GMsgid)

	return
}

// DestructHandle PPMQPublish.
func (t *PPMQPublisher) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQPublish, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
