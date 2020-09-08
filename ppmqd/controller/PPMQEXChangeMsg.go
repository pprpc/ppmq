package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	mqc "xcthings.com/ppmq/common"
	g "xcthings.com/ppmq/ppmqd/common/global"
	lm "xcthings.com/ppmq/ppmqd/model"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQEXChangeMsg"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
)

// PPMQEXChangeMsger .
type PPMQEXChangeMsger struct{}

// ReqHandle PPMQEXChangeMsg request handle
func (t *PPMQEXChangeMsger) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQEXChangeMsg.Req) (err error) {

	resp := new(PPMQEXChangeMsg.Resp)
	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}
	if req.SignKey != g.PConf.Ppmq.ClusterSignkey {
		logs.Logger.Errorf("signkey: %s, conf signkey: %s, not match.", req.SignKey, g.PConf.Ppmq.ClusterSignkey)
	}
	go DeliverExChangeMsg(req)
	return
}

// RespHandle PPMQEXChangeMsg response handle
func (t *PPMQEXChangeMsger) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQEXChangeMsg.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle PPMQEXChangeMsg.
func (t *PPMQEXChangeMsger) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQEXChangeMsg, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}

//DeliverExChangeMsg
func DeliverExChangeMsg(req *PPMQEXChangeMsg.Req) {
	pub := new(PPMQPublish.Req)
	err := pprpc.Unmarshal(req.Payload, pub)
	if err != nil {
		logs.Logger.Errorf("pprpc.Unmarshal(), error: %s.", err)
		return
	}
	var deliver []lm.Sub
	fn := func(k, v interface{}) bool {
		c := v.(pprpc.RPCConn)
		_t, err := c.GetAttr()
		if err != nil {
			return false
		}
		for _, cid := range req.ClientId {
			ci := _t.(*g.ConnAttr)
			if cid == ci.ClientID && (ci.HistorymsgType == mqc.ConsumeNew || ci.HistorymsgType == mqc.ConsumeNewAndOld) {
				go writeMsg(c, pub, ci)
				deliver = append(deliver, lm.Sub{Account: ci.Account, ClientID: ci.ClientID})
			} else if cid == ci.ClientID {
				logs.Logger.Debugf("%s, ClientID: %s, HistorymsgType: %d, not Publish msg.",
					c, ci.ClientID, ci.HistorymsgType)
			}
		}
		return true
	}
	g.Sess.Range(fn)
	err = lm.AddMsgLog(deliver, pub.MsgId)
	if err != nil {
		logs.Logger.Errorf("lm.AddMsgLog(deliver,msgid), MsgID: %s , error: %s.", pub.MsgId, err)
	}

}
