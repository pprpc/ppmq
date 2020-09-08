package model

import (
	"fmt"

	"github.com/pprpc/util/logs"
	m "github.com/pprpc/ppmq/model"
	"github.com/pprpc/ppmq/protoc/ppmqonline/OnlineState"
)

// MOnlineState OnlineState
func MOnlineState(req *OnlineState.Req) (resp *OnlineState.Resp, code uint64, err error) {
	var rows []*m.Connection
	q := new(m.Connection)
	rows, code, err = q.QueryIn(req.Accounts)
	if err != nil {
		return
	}
	resp = new(OnlineState.Resp)
	for _, row := range rows {
		_s := new(OnlineState.OnlineInfo)
		_s.Account = row.Account
		_s.Online = row.ISOnline
		if _s.Online == 0 {
			_s.Online = 2
		}
		if row.ConnType == 1 {
			_s.ConnInfo = fmt.Sprintf("tcp://%s", row.ConnInfo)
		} else if row.ConnType == 2 {
			_s.ConnInfo = fmt.Sprintf("udp://%s", row.ConnInfo)
		} else {
			logs.Logger.Errorf("account: [%s], ConnType: %d.", row.Account, row.ConnType)
			continue
		}
		_s.LastTime = row.LastTime
		_s.Sleep = row.ISSleep
		_s.SrvId = row.ServerID
		logs.Logger.Debugf("account: [%s], ConnType: %d, ConnInfo: %s, Online: %d.",
			_s.Account, row.ConnType, _s.ConnInfo, _s.Online)
		resp.Stas = append(resp.Stas, _s)
	}
	return
}
