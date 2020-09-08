package msg

import (
	"context"
	"fmt"

	"github.com/pprpc/util/common"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
	pb "xcthings.com/protoc/user/ProcessVcode"
)

// ProcessVcode .
func (msg *PPMQMsg) ProcessVcode(r *pb.Req) (msgid string, err error) {
	req := new(PPMQPublish.Req)
	req.MsgBody, err = pprpc.Marshal(r)
	if err != nil {
		err = fmt.Errorf("pprpc.Marshal(), %s", err)
		return
	}

	req.MsgId = getMsgid("vcode")
	req.Topic = fmt.Sprintf("/notify/vcode/%d", common.GetIID())
	req.Format = int32(packets.TYPEPBBIN)
	req.Cmdid = uint64(pb.CmdID)
	req.CmdType = int32(packets.RPCREQ)
	req.Timems = common.GetTimeMs()
	req.Dup = 0
	req.Retain = 0
	req.Qos = 0

	rpc := PPMQPublish.NewPPMQPublishClient(msg.Cli.Conn)
	resp, _, err := rpc.PPCall(context.Background(), req)
	if err != nil {
		err = fmt.Errorf("PPMQPublish.PPCall(), %s", err)
		return
	}
	msgid = resp.GMsgid
	return
}
