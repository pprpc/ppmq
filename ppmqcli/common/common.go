package common

import (
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQPublish"
	"github.com/pprpc/core/packets"
)

type PpmqConn struct {
	PubCB RecivePublishCallBack
}

type RecivePublishCallBack func(pkg *packets.CmdPacket, req *PPMQPublish.Req)
