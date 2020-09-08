package model

import (
	"xcthings.com/hjyz/common"
	mqc "xcthings.com/ppmq/common"
	pm "xcthings.com/ppmq/model"
	g "xcthings.com/ppmq/ppmqd/common/global"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQConnect"
	"xcthings.com/pprpc"
)

// MPPMQConnect PPMQConnect  
func MPPMQConnect(c pprpc.RPCConn, req *PPMQConnect.Req) (resp *PPMQConnect.Resp, code uint64, err error) {

	resp = new(PPMQConnect.Resp)

	q := new(pm.Connection)
	q.ClientID = req.ClientId
	uid := mqc.GetUserID(req.Account)
	if uid == 0 {
		q.Account = req.Account
	} else {
		q.UserID = uid
	}
	q.ServerID = g.Conf.Public.ServerID
	q.ClearSession = req.ClearSession
	q.WillBody = req.WillBody
	q.WillFlag = req.WillFlag
	q.WillQOS = req.WillQos
	q.WillRetain = req.WillRetain
	q.WillTopic = req.WillTopic

	q.ISOnline = mqc.ONLINE
	q.ConnType = 1
	if c.Type() == "U" {
		q.ConnType = 2
	} else if c.Type() == "M" {
		q.ConnType = 3
	}
	q.ConnInfo = c.RemoteAddr().String()
	q.HistorymsgType = req.HistorymsgType
	q.HistorymsgCount = req.HistorymsgCount
	q.LastTime = common.GetTimeMs()
	q.GlobalSync = 0

	code, err = q.Set()
	if err != nil {
		return
	}
	resp.Code = 0
	resp.Sp = 0

	return
}
