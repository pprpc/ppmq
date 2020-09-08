package logic

import (
	"fmt"
	"runtime"

	"github.com/go-xorm/xorm"
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	mqc "github.com/pprpc/ppmq/common"
	pm "github.com/pprpc/ppmq/model"
	g "github.com/pprpc/ppmq/ppmqd/common/global"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQPublish"
	"github.com/pprpc/core"
)

type writemsg func(c pprpc.RPCConn, req *PPMQPublish.Req, ci *g.ConnAttr)
type offlineNotify func(ci *g.ConnAttr)

// LoopSendOfflineMsg .
func LoopSendOfflineMsg(c pprpc.RPCConn, ci *g.ConnAttr, fn writemsg, endMS int64) {

	// recover
	defer func() {
		if err := recover(); err != nil {
			logs.Logger.Errorf("recover: %s.", err)
			var stack string
			for i := 1; ; i++ {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				stack = stack + fmt.Sprintln(file, line)
			}
			logs.Logger.Errorf("Stack:\n%s", stack)
		}
	}()

	if c.Type() == "M" &&
		g.PConf.Ppmq.OfflinemsgSendSleepms > 0 &&
		g.PConf.Ppmq.OfflinemsgSendSleepms < 30000 {

		common.SleepMs(g.PConf.Ppmq.OfflinemsgSendSleepms)
	}

	rows := make([]*pm.MsgStatus, 0)
	var err error
	var startMS int64
	curMS := common.GetTimeMs()
	if g.PConf.Ppmq.OfflinemsgTimeoutms > 0 {
		startMS = curMS - g.PConf.Ppmq.OfflinemsgTimeoutms
	} else {
		startMS = curMS - 86400000
	}

	if ci.HistorymsgCount > 0 {
		err = pm.Orm.Cols("msg_id").Where("status = ?", 1).
			And("client_id = ?", ci.ClientID).
			And("create_time < ?", endMS).
			And("create_time > ?", startMS).
			And("qos > 0").Desc("create_time").
			Limit(int(ci.HistorymsgCount), 0).Find(&rows, new(pm.MsgStatus))
	} else {
		err = pm.Orm.Cols("msg_id").Where("status = ?", 1).
			And("client_id = ?", ci.ClientID).
			And("create_time < ?", endMS).
			And("create_time > ?", startMS).
			And("qos > 0").Desc("create_time").Find(&rows, new(pm.MsgStatus))
	}
	if err != nil {
		logs.Logger.Errorf("LoopSendOfflineMsg, error: %s.", err)
		return
	}
	_l := len(rows) - 1

	logs.Logger.Debugf("account: [%s], offline message count: %d.", ci.Account, _l+1)
	//
	for i := _l; i >= 0; i-- {
		v := rows[i]
		//
		q := new(pm.MsgRaw)
		q.MsgID = v.MsgID
		_, err = q.Get()
		if err != nil {
			logs.Logger.Errorf("MsgRaw.Get(), error: %s.", err)
			continue
		}
		req := new(PPMQPublish.Req)
		err = pprpc.Unmarshal(q.MsgPayload, req)
		if err != nil {
			logs.Logger.Errorf("pprpc.Unmarshal(), MsgID: %s , error: %s.", v.MsgID, err)
			continue
		}
		//
		log := new(pm.MsgLog)
		log.ClientID = ci.ClientID
		log.MsgID = v.MsgID
		n, _, _ := log.GetCount()
		if n > 0 {
			req.Dup = 1
		} else {
			req.Dup = 0
		}

		fn(c, req, ci)
	}
	//
	if g.PConf.Ppmq.MessageOrder == 1 {
		if g.PConf.Ppmq.OfflinemsgSendIntervalms <= 0 {
			g.PConf.Ppmq.OfflinemsgSendIntervalms = 5000
		}
		for e := ci.HistoryCmdList.Front(); e != nil; e = e.Next() {
			pkg := e.Value.(g.HistoryList)
			logs.Logger.Debugf("account: [%s], msgid: [%s], type: %s, list loop", ci.Account, pkg.MsgID, c.Type())

			select {
			case <-c.HandleClose().Done():
				goto sendHistoryMsgEnd
			default:
				_, err = pkg.CmdPkg.Write(c)
				if err != nil {
					logs.Logger.Warnf("%s, msgid: [%s], account: [%s], pprpc write history msg, %s.",
						c, pkg.MsgID, ci.Account, err)
				}
			}
		}
		//
		RunMsgOrder(c, ci)
	}
sendHistoryMsgEnd:
}

// ClearSub  .
func ClearSub(clientID string) {
	d := new(pm.Subscribe)
	d.ClientID = clientID
	_, err := d.Delete()
	if err != nil {
		logs.Logger.Errorf("ClearSub(%s), error: %s.", clientID, err)
		return
	}

	// redis
	if g.PConf.Redis.Addr != "" {
		err = g.TopicCache.ClearSub(clientID)
		if err != nil {
			logs.Logger.Errorf("g.TopicCache.ClearSub(%s), %s.", clientID, err)
		}
	}
}

// UpdateStatus .
func UpdateStatus(cid string, status int32) {
	s := new(pm.MsgStatus)
	s.ClientID = cid
	_, err := s.UpdateStatusByClientID(status)
	if err != nil {
		logs.Logger.Errorf("MsgStatus.UpdateByClientID(), ClientID: %s, error: %s.", cid, err)
	}
}

// CheckOnline .
func CheckOnline(sec int, notify offlineNotify) (err error) {
	// waiting device reconnect.
	common.Sleep(sec * 2)

	s := new(pm.Connection)
	var rows *xorm.Rows
	rows, err = pm.Orm.Where("server_id = ?", g.Conf.Public.ServerID).And("is_online = ?", mqc.OFFLINE).Rows(s)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(s)
		if err != nil {
			logs.Logger.Errorf("CheckOnline,  rows.Scan(s), error: %s.", err)
			continue
		}
		if s.Account == "" && s.UserID == 0 {
			continue
		}
		ci := new(g.ConnAttr)
		ci.Account = s.Account
		ci.ClientID = s.ClientID
		ci.ISSleep = 0
		ci.UserID = s.UserID

		notify(ci)
	}
	return
}

// RunMsgOrder .
func RunMsgOrder(c pprpc.RPCConn, ci *g.ConnAttr) {
	var err error
	logs.Logger.Debugf("RunMsgOrder, account: [%s], CmdChan: %d,  type: %s.",
		ci.Account, len(ci.CmdChan), c.Type())

	if len(ci.CmdChan) == 0 {
		ci.SendOfflinemsgEnd = true
	}
	for {
		select {
		case <-c.HandleClose().Done():
			logs.Logger.Warn("RunMsgOrder, c.HandleClose().Done().")
			return
		case cmd := <-ci.CmdChan:
			if len(ci.CmdChan) == 0 {
				ci.SendOfflinemsgEnd = true
			}
			_, err = cmd.Write(c)
			if err != nil {
				logs.Logger.Warnf("RunMsgOrder, PPRPC, cmd.Write(), %s.", err)
			}

		}
	}
}
