package logic

import (
	"fmt"

	"xcthings.com/ppmq/protoc/ppmqd/PPMQGetClientID"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"

	errc "xcthings.com/ppmq/common/errorcode"
	m "xcthings.com/ppmq/ppmqd/model"
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
