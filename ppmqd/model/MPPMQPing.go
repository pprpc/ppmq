package model

import (
	pm "xcthings.com/ppmq/model"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPing"
)

// MPPMQPing PPMQPing  
func MPPMQPing(req *PPMQPing.Req) (resp *PPMQPing.Resp, code uint64, err error) {
	resp = new(PPMQPing.Resp)
	return
}

func UpdateSleep(cid string, sleep int32) (code uint64, err error) {
	row := new(pm.Connection)
	row.ClientID = cid
	row.ISSleep = sleep

	code, err = row.UpdateSleep()
	return
}
