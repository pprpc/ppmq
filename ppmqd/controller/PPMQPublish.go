package controller

import (
	"fmt"

	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	pm "xcthings.com/ppmq/model"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	g "xcthings.com/ppmq/ppmqd/common/global"
	l "xcthings.com/ppmq/ppmqd/logic"
)

// PPMQPublisher .
type PPMQPublisher struct{}

// ReqHandle PPMQPublish request handle
func (t *PPMQPublisher) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQPublish.Req) (err error) {
	var _t interface{}
	_t, err = checkConn(c, pkg)
	if err != nil {
		logs.Logger.Errorf("checkConn(), error: %s.", err)

		_, err = pprpc.WriteResp(c, pkg, nil)
		if err != nil {
			logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
			return
		}
		return
	}
	ci := _t.(*g.ConnAttr)

	var code uint64
	var resp *PPMQPublish.Resp
	resp, code, err = l.LPPMQPublish(c, pkg, req, ci)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LPPMQPublish, code: %d, err: %s.", c, code, err)
	}

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}
	if pkg.Code == 0 {
		logs.Logger.Debugf("%s, PPMQPublish, msgid: [%s], account: [%s] .",
			c, resp.GMsgid, ci.Account)
	}
	if err != nil || pkg.Code != 0 {
		return
	}
	if req.Dup == 0 {
		go Deliver2Customer(req, resp.MsgId)
	}

	return
}

// RespHandle PPMQPublish response handle
func (t *PPMQPublisher) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQPublish.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}

	// modify MsgLog, MsgStatus
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
	ci.StopTimer(fmt.Sprintf("%d", pkg.CmdSeq))

	// binwen: 暂时没有使用Redis存放ClientID对应的离线消息
	// if g.PConf.Redis.Addr != "" {
	// 	g.TopicCache.RemoveOfflineMsgID(ci.ClientID, resp.MsgId)
	// }

	_, err = pm.SetMsgStatus(resp.MsgId, ci.ClientID, ci.Account, g.Conf.Public.ServerID, 2)
	if err != nil {
		logs.Logger.Errorf("pm.SetMsgStatus(), error: %s.", err)
		return
	}
	logs.Logger.Debugf("account: %s, msgid: %s, puback ok.", ci.Account, resp.MsgId)

	return
}

// DestructHandle PPMQPublish.
func (t *PPMQPublisher) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQPublish, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
