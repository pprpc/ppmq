package logic

import (
	"xcthings.com/ppmq/protoc/ppmqd/PPMQGetSublist"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"

	m "xcthings.com/ppmq/ppmqd/model"
)

// LPPMQGetSublist PPMQGetSublist Business logic
func LPPMQGetSublist(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQGetSublist.Req, clientID string) (resp *PPMQGetSublist.Resp, code uint64, err error) {

	resp, code, err = m.MPPMQGetSublist(req, clientID)
	return
}
