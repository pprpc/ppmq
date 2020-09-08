package logic

import (
	"fmt"

	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQGetClientID"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	errc "github.com/pprpc/ppmq/common/errorcode"
	m "github.com/pprpc/ppmq/ppmqd/model"
)

// LPPMQGetClientID PPMQGetClientID Business logic
func LPPMQGetClientID(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQGetClientID.Req) (resp *PPMQGetClientID.Resp, code uint64, err error) {
	if req.Account == "" || req.HardwareInfo == "" {
		code = errc.ParameterError
		err = fmt.Errorf("Incomplete parameters: Account, HardwareInfo")
		return
	}
	resp, code, err = m.MPPMQGetClientID(req)
	return
}
