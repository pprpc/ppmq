package logic

import (
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQEXChangeMsg"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	m "github.com/pprpc/ppmq/ppmqd/model"
)

// LPPMQEXChangeMsg PPMQEXChangeMsg Business logic
func LPPMQEXChangeMsg(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQEXChangeMsg.Req) (resp *PPMQEXChangeMsg.Resp, code uint64, err error) {

	resp, code, err = m.MPPMQEXChangeMsg(req)
	return
}
