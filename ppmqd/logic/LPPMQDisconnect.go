package logic

import (
	"xcthings.com/ppmq/protoc/ppmqd/PPMQDisconnect"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	m "xcthings.com/ppmq/ppmqd/model"
)

// LPPMQDisconnect PPMQDisconnect Business logic
func LPPMQDisconnect(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQDisconnect.Req) (resp *PPMQDisconnect.Resp, code uint64, err error) {

	resp, code, err = m.MPPMQDisconnect(req)
	return
}
