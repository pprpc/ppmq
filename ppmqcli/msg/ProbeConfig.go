package msg

import (
	"context"
	"fmt"

	"xcthings.com/hjyz/common"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
	pb "xcthings.com/protoc/ftconnnat/ProbeConfig"
)

// ProbeConfig .
func (msg *PPMQMsg) ProbeConfig(r *pb.Req, serverid string) (msgid string, err error) {
	req := new(PPMQPublish.Req)
	req.MsgBody, err = pprpc.Marshal(r)
	if err != nil {
		err = fmt.Errorf("pprpc.Marshal(), %s", err)
		return
	}

	req.MsgId = getMsgid("nat")
	req.Topic = "/notify/nat/" + serverid
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
