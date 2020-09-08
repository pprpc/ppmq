package controller

import (
	"context"
	"time"

	"xcthings.com/hjyz/logs"
	mqc "xcthings.com/ppmq/common"
	g "xcthings.com/ppmq/ppmqd/common/global"
	lm "xcthings.com/ppmq/ppmqd/model"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQEXChangeMsg"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
)

// Deliver2Customer .
func Deliver2Customer(req *PPMQPublish.Req, srcmsgid string) {
	subs, err := lm.GetCustomerBYTopic(req.Topic, req.MsgId, srcmsgid)
	if err != nil {
		logs.Logger.Errorf("lm.GetCustomerBYTopic(%s), MsgID: %s ,error: %s.", req.Topic, req.MsgId, err)
		return
	}
	var lSubs, rSubs, sendSubs []lm.Sub
	var cids []string
	for _, v := range subs.Subs {
		cids = append(cids, v.ClientID)
	}

	logs.Logger.Debugf("topic: %s, msgid: %s, cids: %v.", req.Topic, req.MsgId, cids)
	//oncids, offcids := lm.GetOnlineByClientIDS(cids)
	oncids, _ := lm.GetOnlineByClientIDS(cids)

	// online clientid
	for _, v := range oncids {
		if v.ServerID == g.Conf.Public.ServerID {
			for _, v1 := range subs.Subs {
				if v1.ClientID == v.ClientID {
					v1.ServerID = v.ServerID
					lSubs = append(lSubs, v1)
				}
			}
		} else {
			// online to other serverid
			for _, v1 := range subs.Subs {
				if v1.ClientID == v.ClientID {
					v1.ServerID = v.ServerID
					rSubs = append(rSubs, v1)
				}
			}
		}
	}

	if len(rSubs) > 0 {
		go EXChangeMsg(req, rSubs)
	}

	fn := func(k, v interface{}) bool {
		c := v.(pprpc.RPCConn)
		_t, err := c.GetAttr()
		if err != nil {
			return true
		}
		ci := _t.(*g.ConnAttr)
		for _, row := range lSubs {
			if row.ClientID == ci.ClientID && (ci.HistorymsgType == mqc.ConsumeNew ||
				ci.HistorymsgType == mqc.ConsumeNewAndOld) {

				go writeMsg(c, req, ci)
				sendSubs = append(sendSubs, row)
			} else if row.ClientID == ci.ClientID {
				logs.Logger.Debugf("%s, ClientID: %s, HistorymsgType: %d.",
					c, ci.ClientID, ci.HistorymsgType)
			}
		}
		return true
	}
	g.Sess.Range(fn)
	err = lm.AddMsgLog(sendSubs, req.MsgId)
	if err != nil {
		logs.Logger.Errorf("lm.AddMsgLog(deliver,msgid), MsgID: %s , error: %s.", req.MsgId, err)
	}
}

// writeMsg
func writeMsg(c pprpc.RPCConn, req *PPMQPublish.Req, ci *g.ConnAttr) {

	logs.Logger.Debugf("msgid: [%s], account: [%s], MsgOrder: %d, SendOfflinemsgEnd: %v.",
		req.MsgId, ci.Account, g.PConf.Ppmq.MessageOrder, ci.SendOfflinemsgEnd)

	seq := pprpc.GetSeqID()
	cmd := packets.NewCmdPacket(ci.MessageType)
	cmd.CmdSeq = seq
	cmd.CmdID = PPMQPublish.CmdID
	cmd.EncType = ci.EncType
	cmd.RPCType = packets.RPCREQ

	var err error
	if ci.MessageType == packets.TYPEPBBIN {
		cmd.Payload, err = pprpc.Marshal(req)
	} else if ci.MessageType == packets.TYPEPBJSON {
		cmd.Payload, err = pprpc.MarshalJSON(req)
	}
	if err != nil {
		logs.Logger.Errorf("pprpc.Marshal/MarshalJSON(req), msgid: [%s] , error: %s.", req.MsgId, err)
		return
	}
	if c.Type() == "U" {
		cmd.FixHeader.SetProtocol(packets.PROTOUDP)
	}

	if g.PConf.Ppmq.MessageOrder == 1 && ci.SendOfflinemsgEnd == false {
		ci.CmdChan <- cmd
		logs.Logger.Debugf("%s, account: %s, msgid: [%s], write CmdChan ok.",
			c, ci.Account, req.MsgId)
		return
	}
	_, err = cmd.Write(c)
	if err != nil {
		logs.Logger.Errorf("cmd.Write(c), MsgID: %s ,error: %s.", req.MsgId, err)
	} else if c.Type() == "U" && g.PConf.Ppmq.UDPRespTimeoutms > 0 {
		fn := func() {
			cmd.Payload, _ = pprpc.Marshal(req)
			_, err := cmd.Write(c)
			if err != nil {
				logs.Logger.Errorf("UDP Timers, key: %d, cmd.Write(c), MsgID: %s ,error: %s.",
					cmd.CmdSeq, req.MsgId, err)
				return
			}
			logs.Logger.Warnf("UDP Timers, key: %d, cmd.Write(c), MsgID: %s, Timeout(ms): %d, OK.",
				cmd.CmdSeq, req.MsgId, g.PConf.Ppmq.UDPRespTimeoutms)
		}
		ci.StartTimer(time.Duration(g.PConf.Ppmq.UDPRespTimeoutms), cmd.CmdSeq, fn)
	}
	logs.Logger.Debugf("%s, account: %s, msgid: [%s], write publish ok.",
		c, ci.Account, req.MsgId)
}

func writeOfflineMsg(c pprpc.RPCConn, req *PPMQPublish.Req, ci *g.ConnAttr) {

	logs.Logger.Debugf("write offline msg, type: %s, msgid: [%s], account: [%s].",
		c.Type(), req.MsgId, ci.Account)

	seq := pprpc.GetSeqID()
	cmd := packets.NewCmdPacket(ci.MessageType)
	cmd.CmdSeq = seq
	cmd.CmdID = PPMQPublish.CmdID
	cmd.EncType = ci.EncType
	cmd.RPCType = packets.RPCREQ

	var err error
	if ci.MessageType == packets.TYPEPBBIN {
		cmd.Payload, err = pprpc.Marshal(req)
	} else if ci.MessageType == packets.TYPEPBJSON {
		cmd.Payload, err = pprpc.MarshalJSON(req)
	}
	if err != nil {
		logs.Logger.Errorf("pprpc.Marshal/MarshalJSON(req), MsgID: %s , error: %s.", req.MsgId, err)
		return
	}
	if c.Type() == "U" {
		cmd.FixHeader.SetProtocol(packets.PROTOUDP)
	}
	ci.HistoryCmdList.PushBack(g.HistoryList{MsgID: req.MsgId, CmdPkg: cmd})
}

// EXChangeMsg .
func EXChangeMsg(req *PPMQPublish.Req, subs []lm.Sub) {
	if g.PConf.Ppmq.Mode == "single" {
		return
	}
	// FIXME:
	var err error
	for _, v := range subs {
		rreq := new(PPMQEXChangeMsg.Req)
		rreq.SignKey = g.PConf.Ppmq.ClusterSignkey
		rreq.ServerId = v.ServerID
		rreq.ClientId = append(rreq.ClientId, v.ClientID)
		rreq.Payload, err = pprpc.Marshal(req)
		if err != nil {
			logs.Logger.Warnf("pprpc.Marshal(PPMQPublish.Req), %s.", err)
			continue
		}
		pkg, _, err := g.EXChangeConn.InvokeServerID(context.Background(), g.MSName, v.ServerID, PPMQEXChangeMsg.CmdID, rreq)
		if err != nil {
			logs.Logger.Warnf("g.EXChangeConn.InvokeServerID(ctx, %s, %s, %s,req), msgid: [%s], %s.",
				g.MSName, v.ServerID, PPMQEXChangeMsg.CmdID, req.MsgId, err)

		} else {
			if pkg.Code != 0 {
				logs.Logger.Warnf("g.EXChangeConn.InvokeServerID(ctx, %s, %d, %d,req), msgid: [%s], code: %d.",
					g.MSName, v.ServerID, req.MsgId, pkg.Code)
			} else {
				logs.Logger.Debugf("PPMQEXChangeMsg: %s, msgid: [%s], OK.", v.ServerID, req.MsgId)
			}
		}
	}
}
