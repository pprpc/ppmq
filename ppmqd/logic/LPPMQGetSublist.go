package logic

import (
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQGetSublist"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	m "github.com/pprpc/ppmq/ppmqd/model"
)

// LPPMQGetSublist PPMQGetSublist Business logic
func LPPMQGetSublist(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQGetSublist.Req, clientID string) (resp *PPMQGetSublist.Resp, code uint64, err error) {

	resp, code, err = m.MPPMQGetSublist(req, clientID)
	return
}
