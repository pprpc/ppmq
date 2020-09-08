package logic

import (
	"fmt"

	"github.com/pprpc/ppmq/protoc/ppmqonline/OnlineState"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
	errc "xcthings.com/xc-app/common/errorcode"

	m "github.com/pprpc/ppmq/ppmqonline-ms/model"
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
