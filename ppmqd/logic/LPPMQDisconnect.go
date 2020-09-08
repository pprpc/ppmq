package logic

import (
	"xcthings.com/ppmq/protoc/ppmqd/PPMQDisconnect"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"

	m "xcthings.com/ppmq/ppmqd/model"
)

// LPPMQDisconnect PPMQDisconnect Business logic
func LPPMQDisconnect(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQDisconnect.Req) (resp *PPMQDisconnect.Resp, code uint64, err error) {

	resp, code, err = m.MPPMQDisconnect(req)
	return
}
