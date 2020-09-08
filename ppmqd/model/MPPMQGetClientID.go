package model

import (
	pm "github.com/pprpc/ppmq/model"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQGetClientID"
)

// MPPMQGetClientID PPMQGetClientID  
func MPPMQGetClientID(req *PPMQGetClientID.Req) (resp *PPMQGetClientID.Resp, code uint64, err error) {
	resp = new(PPMQGetClientID.Resp)

	q := new(pm.Clientid)
	q.Account = req.Account
	q.HwFeature = req.HardwareInfo
	code, err = q.GetAdv()
	if code == 0 {
		resp.ClientId = q.ClientID
	} else {
		err = nil
		ins := new(pm.Clientid)
		ins.Account = req.Account
		ins.HwFeature = req.HardwareInfo
		ins.ClientID = getClientID(req.Account, req.HardwareInfo)

		code, err = ins.Add()
		if code == 0 {
			resp.ClientId = ins.ClientID
		}
	}
	return
}
