package logic

import (
	"fmt"

	"github.com/pprpc/util/logs"
	errc "xcthings.com/ppmq/common/errorcode"
	pm "xcthings.com/ppmq/model"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	mqc "xcthings.com/ppmq/common"
	g "xcthings.com/ppmq/ppmqd/common/global"
	t "xcthings.com/ppmq/ppmqd/common/topic"
	m "xcthings.com/ppmq/ppmqd/model"
)

// LPPMQPublish PPMQPublish Business logic
func LPPMQPublish(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQPublish.Req, ci *g.ConnAttr) (resp *PPMQPublish.Resp, code uint64, err error) {

	if req.MsgId == "" || req.Topic == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:MsgId")
		return
	}
	if req.Format != int32(packets.TYPEPBBIN) &&
		req.Format != int32(packets.TYPEPBJSON) {
		code = errc.ParameterIllegal
		err = fmt.Errorf("The parameter is invalid: Format(%d)", req.Format)
		return
	}
	if req.CmdType != int32(packets.RPCREQ) && req.CmdType != int32(packets.RPCRESP) {
		code = errc.ParameterIllegal
		err = fmt.Errorf("The parameter is invalid: CmdType(%d)", req.CmdType)
		return
	}
	if err = t.CheckPubTopic(req.Topic); err != nil {
		code = errc.ParameterIllegal
		return
	}

	if g.PConf.Auth == 2 && g.PConf.CheckTopic == 1 {
		// check topic
		pc, lerr := m.CheckPub(ci.Account, req.Topic)
		if lerr != nil {
			logs.Logger.Warnf("LPPMQPublish, CheckPub, %s.", lerr)
			code = errc.ParameterError
			err = fmt.Errorf("m.CheckPub(), %s", lerr)
			return
		}
		if pc != 0 {
			code = errc.PUBREJECT
			err = fmt.Errorf("topic: %s, pub reject", req.Topic)
			return
		}
	}
	// dup
	if req.Dup == mqc.REPEATSEND {
		q := new(pm.MsgInfo)
		q.SrcMsgid = req.MsgId
		q.ClientID = ci.ClientID
		code, _ = q.GetAdv()
		if code == 0 {
			resp = new(PPMQPublish.Resp)
			resp.MsgId = req.MsgId
			resp.GMsgid = q.MsgID
			return
		}
	} else {

		req.Dup = 0
	}

	resp, code, err = m.MPPMQPublish(req, ci)
	return
}
