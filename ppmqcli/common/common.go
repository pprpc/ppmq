package common

import (
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"xcthings.com/pprpc/packets"
)

type PpmqConn struct {
	PubCB RecivePublishCallBack
}

type RecivePublishCallBack func(pkg *packets.CmdPacket, req *PPMQPublish.Req)
