package msg

import (
	"context"
	"fmt"

	"xcthings.com/hjyz/common"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
	pb "xcthings.com/protoc/ftconnp2p/NotifyConn"
)

// NotifyConn .
func (msg *PPMQMsg) NotifyConn(r *pb.Req, did string) (msgid string, err error) {
	req := new(PPMQPublish.Req)
	req.MsgBody, err = pprpc.Marshal(r)
	if err != nil {
		err = fmt.Errorf("pprpc.Marshal(), %s", err)
		return
	}
	if did == "" {
		err = fmt.Errorf("Did is null")
		return
	}
	if r.ConnType == 1 {
		req.MsgId = getMsgid("p2p")
	} else {
		req.MsgId = getMsgid("relay")
	}
	req.Topic = "/notify/conn/" + did
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
