package controller

import (
	"fmt"

	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	mqc "github.com/pprpc/ppmq/common"
	errc "github.com/pprpc/ppmq/common/errorcode"
	pm "github.com/pprpc/ppmq/model"
	g "github.com/pprpc/ppmq/ppmqd/common/global"
	l "github.com/pprpc/ppmq/ppmqd/logic"
	m "github.com/pprpc/ppmq/ppmqd/model"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQOAONotify"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQPublish"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
	"xcthings.com/xc-app/protoc/ipc/OfflineNotify"
)

func checkConn(c pprpc.RPCConn, pkg *packets.CmdPacket) (ci interface{}, err error) {
	ci, err = c.GetAttr()
	if err != nil || ci == nil {
		pkg.Code = errc.NOTRUNPPMQCONNECT
		err = fmt.Errorf("c.GetAttr(), %s", err)
		return
	}
	_t := ci.(*g.ConnAttr)
	if _t.IsAuth == false {
		pkg.Code = errc.NOTRUNPPMQCONNECT
		err = fmt.Errorf("IsAuth: %v, ", _t.IsAuth)
		return
	}

	return
}

// OnlineEvent .
func OnlineEvent(c pprpc.RPCConn, startMS int64) {
	_t, _ := c.GetAttr()
	ci := _t.(*g.ConnAttr)

	if ci.ClearSession == mqc.SESSIONCLEAR {
		l.ClearSub(ci.ClientID)
		l.UpdateStatus(ci.ClientID, mqc.MSGSTATUSClearSession)
	}

	if g.PConf.Ppmq.OnlineNotify {
		onlineNotify(ci)
	}
	if ci.HistorymsgType > mqc.ConsumeNew && ci.ClearSession == mqc.SESSIONKEEP {
		if g.PConf.Ppmq.MessageOrder == 1 {
			go l.LoopSendOfflineMsg(c, ci, writeOfflineMsg, startMS)
		} else {
			go l.LoopSendOfflineMsg(c, ci, writeMsg, startMS)
		}
	} else {
		ci.SendOfflinemsgEnd = true
	}
}

// OfflineEvent .
func OfflineEvent(c pprpc.RPCConn) {
	_t, _ := c.GetAttr()
	ci := _t.(*g.ConnAttr)
	err := m.SetOffline(ci.ClientID)
	if err != nil {
		logs.Logger.Errorf("%s, m.SetOffline(%s), error: %s.", c, ci.ClientID, err)
	}
	if g.PConf.Ppmq.OfflineNotify {
		offlineNotify(ci)
	}
}

// OfflineWillMsg .
// FIXME:
func OfflineWillMsg(clientID string) {

}

func onlineNotify(ci *g.ConnAttr) {
	pub, err := createNotify(ci, mqc.ONLINE)
	if err != nil {
		logs.Logger.Errorf("onlineNotify, createNotify(), error: %s.", err)
	}
	go Deliver2Customer(pub, "")
}

func offlineNotify(ci *g.ConnAttr) {
	pub, err := createNotify(ci, mqc.OFFLINE)
	if err != nil {
		logs.Logger.Errorf("offlineNotify, createNotify(), error: %s.", err)
	}
	go Deliver2Customer(pub, "")
}

func createNotify(ci *g.ConnAttr, s int32) (pub *PPMQPublish.Req, err error) {
	req := new(PPMQOAONotify.Req)
	req.Account = ci.Account
	req.ClientId = ci.ClientID
	req.IsSleep = ci.ISSleep
	req.UserId = ci.UserID
	req.Status = s
	req.ServerId = g.Conf.Public.ServerID
	req.Ipaddr = ci.RemoteIpaddr

	pub = new(PPMQPublish.Req)
	//pub.Topic = g.PConf.Ppmq.NotifyTopicPrefix + ci.Account
	pub.Topic = fmt.Sprintf("%s%s/%d", g.PConf.Ppmq.NotifyTopicPrefix, ci.Account, common.GetIID())
	if s == mqc.ONLINE {
		pub.MsgId = fmt.Sprintf("Online_%d", common.GetTimeNs())
	} else {
		pub.MsgId = fmt.Sprintf("Offline_%d", common.GetTimeNs())
	}
	pub.Format = int32(packets.TYPEPBBIN)
	pub.Cmdid = PPMQOAONotify.CmdID
	pub.CmdType = int32(packets.RPCREQ)
	pub.Timems = common.GetTimeMs()
	pub.MsgBody, _ = pprpc.Marshal(req)
	pub.Dup = 0
	pub.Retain = 0
	pub.Qos = 0

	// msg_info , msg_raw
	_, _, err = m.MPPMQPublish(pub, ci)

	return
}

func createOfflineNotify(account, cid string, ci *g.ConnAttr) (pub *PPMQPublish.Req, err error) {
	req := new(OfflineNotify.Req)
	req.OfflineCode = 100
	req.ClientId = cid

	pub = new(PPMQPublish.Req)
	pub.Topic = fmt.Sprintf("/iot/p/config/%s", account)
	pub.Format = int32(packets.TYPEPBBIN)
	pub.Cmdid = OfflineNotify.CmdID
	pub.CmdType = int32(packets.RPCREQ)
	pub.Timems = common.GetTimeMs()
	pub.MsgBody, _ = pprpc.Marshal(req)
	pub.Dup = 0
	pub.Retain = 0
	pub.Qos = 0

	// msg_info , msg_raw
	_, _, err = m.MPPMQPublish(pub, ci)

	return
}

// ClearSub .
func ClearSub(clientID string) {
	l.ClearSub(clientID)
}

// UpdateStatus .
func UpdateStatus(clientID string, status int32) {
	l.UpdateStatus(clientID, status)
}

//RunProcess .
func RunProcess(sec uint32) (err error) {
	err = onoffLineOnce(int(sec))
	return
}

// onoffLine .
func onoffLineOnce(sec int) error {
	if g.PConf.Ppmq.OfflineNotify {
		return l.CheckOnline(sec, offlineNotify)
	}
	return nil
}

// SetAllOffline .
func SetAllOffline(serverID string) error {
	return pm.SetOnlineStatusByServerID(serverID, mqc.OFFLINE)
}
