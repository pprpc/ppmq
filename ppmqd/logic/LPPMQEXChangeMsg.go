package logic

import (
	"xcthings.com/ppmq/protoc/ppmqd/PPMQEXChangeMsg"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"

	m "xcthings.com/ppmq/ppmqd/model"
)

// LPPMQEXChangeMsg PPMQEXChangeMsg Business logic
func LPPMQEXChangeMsg(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQEXChangeMsg.Req) (resp *PPMQEXChangeMsg.Resp, code uint64, err error) {

	resp, code, err = m.MPPMQEXChangeMsg(req)
	return
}
