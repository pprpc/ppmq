package model

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	errc "github.com/pprpc/ppmq/common/errorcode"
	pm "github.com/pprpc/ppmq/model"
	g "github.com/pprpc/ppmq/ppmqd/common/global"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQPublish"
	"github.com/pprpc/core"
)

// MPPMQPublish PPMQPublish
func MPPMQPublish(req *PPMQPublish.Req, ci *g.ConnAttr) (resp *PPMQPublish.Resp, code uint64, err error) {
	msgid := GetMsgID(ci.ClientID, req.Topic)
	resp = new(PPMQPublish.Resp)
	resp.GMsgid = msgid
	resp.MsgId = req.MsgId
	req.MsgId = msgid

	if req.Qos == 0 && g.PConf.Ppmq.Qos == true {
		logs.Logger.Debugf("MsgId: %s(%s), Qos: 0, Topic: %s, Body Length: %d.",
			req.MsgId, msgid, req.Topic, len(req.MsgBody))
		return
	}
	srcMsgid := req.MsgId
	raw := new(pm.MsgRaw)
	raw.MsgID = msgid
	raw.MsgPayload, err = pprpc.Marshal(req)
	if err != nil {
		code = errc.MarshalError
		return
	}

	info := new(pm.MsgInfo)
	info.MsgID = req.MsgId
	info.SrcMsgid = srcMsgid
	info.Account = ci.Account
	info.ClientID = ci.ClientID
	info.Dup = req.Dup
	info.Retain = req.Retain
	info.Qos = req.Qos
	info.Topic = req.Topic
	info.Format = req.Format
	info.Cmdid = req.Cmdid
	info.CmdType = req.CmdType
	info.MsgTimems = req.Timems
	info.CreateTime = common.GetTimeMs()

	code, err = pm.SaveMsg(raw, info)
	if err != nil {
		return
	}

	return
}
