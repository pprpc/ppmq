package model

import "github.com/pprpc/ppmq/protoc/ppmqd/PPMQDisconnect"

// MPPMQDisconnect PPMQDisconnect  
func MPPMQDisconnect(req *PPMQDisconnect.Req) (resp *PPMQDisconnect.Resp, code uint64, err error) {

	resp = new(PPMQDisconnect.Resp)
	return
}
