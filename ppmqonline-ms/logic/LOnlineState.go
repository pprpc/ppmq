package logic

import (
	"fmt"

	"xcthings.com/ppmq/protoc/ppmqonline/OnlineState"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
	errc "xcthings.com/xc-app/common/errorcode"

	m "xcthings.com/ppmq/ppmqonline-ms/model"
)

// LOnlineState OnlineState Business logic
func LOnlineState(c pprpc.RPCConn, pkg *packets.CmdPacket, req *OnlineState.Req) (resp *OnlineState.Resp, code uint64, err error) {
	if len(req.Accounts) == 0 {
		code = errc.ParameterError
		err = fmt.Errorf("invalid: accounts must be set")
		return
	}
	resp, code, err = m.MOnlineState(req)
	return
}