package PPMQSubscribe

import (
	"fmt"

	context "golang.org/x/net/context"
	"github.com/pprpc/util/common"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
)

const (
	CmdID   uint64 = 17
	CmdName string = "PPMQSubscribe"
	Module  string = "ppmqd"
)

// RegisterService 注册服务
func RegisterService(s *pprpc.Service, hook Rpcer) {
	s.RegisterService(&_PPRpc_Desc, hook)
}

type Rpcer interface {
	ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *Req) (err error)
	RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *Resp) (err error)
	DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64)
}

// Req
func _PPRpc_Req_Handler(h interface{}, c pprpc.RPCConn, pkg *packets.CmdPacket, callFlag bool, dec func(interface{}) error) (interface{}, error) {
	_s := common.GetTimeMs()
	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("code: %d, %s, RequestHook", pkg.Code, err)
		}
		if callFlag && h != nil {
			h.(Rpcer).DestructHandle(c, pkg, _s)
		}
	}()
	req := new(Req)
	if err = dec(req); err != nil {
		return nil, err
	}
	if callFlag && h != nil {
		err = h.(Rpcer).ReqHandle(c, pkg, req)
		if err != nil {
			return req, err
		}
	}
	return req, nil
}

// Resp
func _PPRpc_Resp_Handler(h interface{}, c pprpc.RPCConn, pkg *packets.CmdPacket, callFlag bool, dec func(interface{}) error) (interface{}, error) {
	_s := common.GetTimeMs()
	var err error
	defer func() {
		if err != nil {
			err = fmt.Errorf("code: %d, %s, RespHandle", pkg.Code, err)
		}
		if callFlag && h != nil {
			h.(Rpcer).DestructHandle(c, pkg, _s)
		}
	}()
	resp := new(Resp)
	if err = dec(resp); err != nil {
		return nil, err
	}
	if callFlag && h != nil {
		err = h.(Rpcer).RespHandle(c, pkg, resp)
		if err != nil {
			return resp, err
		}
	}
	return resp, nil
}

var _PPRpc_Desc = pprpc.ServiceDesc{
	CmdID:       CmdID,
	CmdName:     CmdName,
	ReqHandler:  _PPRpc_Req_Handler,
	RespHandler: _PPRpc_Resp_Handler,
}

// Client API for PPRpc service

type PPMQSubscribeClient interface {
	PPCall(ctx context.Context, in *Req) (*Resp, *packets.CmdPacket, error)
}

type pPMQSubscribeClient struct {
	cc pprpc.RPCCliConn
}

func NewPPMQSubscribeClient(cc pprpc.RPCCliConn) PPMQSubscribeClient {
	return &pPMQSubscribeClient{cc}
}

func (c *pPMQSubscribeClient) PPCall(ctx context.Context, in *Req) (*Resp, *packets.CmdPacket, error) {
	pkg, ans, err := c.cc.Invoke(ctx, CmdID, in)
	if err != nil {
		return nil, nil, err
	}
	return ans.(*Resp), pkg, nil
}
