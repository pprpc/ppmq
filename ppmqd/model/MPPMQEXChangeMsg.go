package model

import "github.com/pprpc/ppmq/protoc/ppmqd/PPMQEXChangeMsg"

// MPPMQEXChangeMsg PPMQEXChangeMsg  
func MPPMQEXChangeMsg(req *PPMQEXChangeMsg.Req) (resp *PPMQEXChangeMsg.Resp, code uint64, err error) {

	resp = new(PPMQEXChangeMsg.Resp)
	return
}
