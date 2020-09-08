package logic

import (
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPing"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"

	m "xcthings.com/ppmq/ppmqd/model"
)

// LPPMQPing PPMQPing Business logic
func LPPMQPing(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQPing.Req) (resp *PPMQPing.Resp, code uint64, err error) {

	resp, code, err = m.MPPMQPing(req)
	return
}
