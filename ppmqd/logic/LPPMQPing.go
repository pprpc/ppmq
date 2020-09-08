package logic

import (
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQPing"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	m "github.com/pprpc/ppmq/ppmqd/model"
)

// LPPMQPing PPMQPing Business logic
func LPPMQPing(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQPing.Req) (resp *PPMQPing.Resp, code uint64, err error) {

	resp, code, err = m.MPPMQPing(req)
	return
}
