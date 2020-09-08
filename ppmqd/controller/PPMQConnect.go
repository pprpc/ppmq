package controller

import (
	"container/list"
	"fmt"

	"github.com/pprpc/util/cache"
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQConnect"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	mqc "xcthings.com/ppmq/common"
	pm "xcthings.com/ppmq/model"
	g "xcthings.com/ppmq/ppmqd/common/global"
	l "xcthings.com/ppmq/ppmqd/logic"
)

// PPMQConnecter .
type PPMQConnecter struct{}

// ReqHandle PPMQConnect request handle
func (t *PPMQConnecter) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQConnect.Req) (err error) {
	startMS := common.GetTimeMs()
	var code uint64
	var resp *PPMQConnect.Resp
	resp, code, err = l.LPPMQConnect(c, pkg, req)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LPPMQConnect(), code: %d, err: %s.", c, code, err)
	}

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}

	if err != nil || pkg.Code != 0 {
		return
	}

	ci := new(g.ConnAttr)
	ci.ClientID = req.ClientId
	ci.Account = req.Account
	ci.UserID = mqc.GetUserID(req.Account)

	ci.IsAuth = true
	ci.ConnectFlag = false

	ci.ClearSession = req.ClearSession
	ci.MessageType = pkg.MessageType
	ci.EncType = pkg.EncType
	ci.HbInterval = resp.HbInterval

	ci.WillFlag = req.WillFlag
	ci.WillQos = req.WillQos
	ci.WillRetain = req.WillRetain
	ci.WillTopic = req.WillTopic
	ci.WillBody = req.WillBody

	ci.HistorymsgType = req.HistorymsgType
	ci.HistorymsgCount = req.HistorymsgCount

	ci.ConnStatus = mqc.SESSIONKEEP
	ci.ConnType = 1
	if c.Type() == "U" {
		ci.ConnType = 2
	} else if c.Type() == "M" {
		ci.ConnType = 3
	}
	ci.LastTime = common.GetTimeMs()
	ci.ISSleep = mqc.PINGNORMAL
	ci.RemoteIpaddr = common.GetIPAddr(c.RemoteAddr())
	ci.OfflineEvent = true
	ci.Timers = cache.NewCache(50)
	ci.RunSub = false
	if req.ClearSession == mqc.SESSIONCLEAR {
		ci.SendOfflinemsgEnd = true
	} else {
		ci.SendOfflinemsgEnd = false
	}
	if g.PConf.Ppmq.MessageOrder == 1 {
		ci.CmdChan = make(chan *packets.CmdPacket, g.PConf.Ppmq.MessageOrderLength)
		ci.HistoryCmdList = list.New()
		ci.HistoryPubackSync = make(chan uint16, 100)
	}

	// check:The same device ID is on-line
	checkSameDeviceIDOnLine(req, ci)

	err = c.SetAttr(ci)
	if err != nil {
		logs.Logger.Errorf("c.SetAttr(ci), error: %s.", err)
		return
	}

	v, err := g.Sess.Get(req.ClientId)
	if err == nil {
		_t := v.(pprpc.RPCConn)
		_ci, e := _t.GetAttr()
		if e == nil {
			if _ci != nil {
				_ci.(*g.ConnAttr).OfflineEvent = false
				common.SleepMs(200)
			}
			_t.Close()
		}
	}
	err = nil
	logs.Logger.Debugf("%s, Sess.Push(%s, c), Account: [%s].", c, req.ClientId, req.Account)
	g.Sess.Push(req.ClientId, c)
	ci.ConnectFlag = true
	//
	go OnlineEvent(c, startMS)

	return
}

// RespHandle PPMQConnect response handle
func (t *PPMQConnecter) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *PPMQConnect.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle PPMQConnect.
func (t *PPMQConnecter) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, PPMQConnect, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}

func checkSameDeviceIDOnLine(req *PPMQConnect.Req, ci *g.ConnAttr) {
	if g.PConf.ClearDIDConf.Clear == false {
		return
	}
	// white account
	for _, row := range g.PConf.ClearDIDConf.WhiteAccount {
		if req.Account == row {
			return
		}
	}

	// check db notify other deviceid offline
	rows, code, err := pm.GetRowsByAccount(req.Account)
	if err != nil {
		logs.Logger.Warnf("pm.GetRowsByAccount(%s), code: %d, error: %s.", req.Account, code, err)
		return
	}
	logs.Logger.Debugf("account: %s, length: %d.", req.Account, len(rows))
	//
	for _, row := range rows {
		logs.Logger.Debugf("account: %s, row.ClientID: %s, req.ClientId: %s, IsOnline: %d.", req.Account, row.ClientID, req.ClientId, row.ISOnline)
		if row.ClientID == req.ClientId || row.ISOnline == mqc.OFFLINE {
			continue
		}

		// 2019-0615
		// delete same did ,if did is offline
		if row.ISOnline == mqc.OFFLINE && row.ClientID != req.ClientId {
			_, err = pm.DeleteConnection(row.ClientID)
			if err != nil {
				logs.Logger.Errorf("pm.DeleteConnection(%s), %s", row.ClientID, err)
			} else {
				logs.Logger.Debugf("pm.DeleteConnection(%s), account: %s(offline), clear.", row.ClientID, req.Account)
			}
			continue
		}

		pub, err := createOfflineNotify(req.Account, row.ClientID, ci)
		if err != nil {
			logs.Logger.Errorf("checkSameDeviceIDOnLine, createOfflineNotify(%s), error: %s.", req.Account, err)
			return
		}
		go Deliver2Customer(pub, fmt.Sprintf("checkSameDeviceIDOnLine-%d", common.GetTimeSec()))
		logs.Logger.Debugf("OfflineNotify: [%s].", pub.Topic)
	}
}
