package ppmqcli

import (
	"xcthings.com/ppmq/protoc/ppmqd/PPMQConnect"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQDisconnect"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQGetClientID"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQGetSublist"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPing"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQSubscribe"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQUnSub"
	"github.com/pprpc/core"

	ctrl "xcthings.com/ppmq/ppmqcli/controller"
)

var service *pprpc.Service

func init() {
	service = pprpc.NewService()
	PPMQGetClientID.RegisterService(service, &ctrl.PPMQGetClientIDer{})
	PPMQConnect.RegisterService(service, &ctrl.PPMQConnecter{})
	PPMQDisconnect.RegisterService(service, &ctrl.PPMQDisconnecter{})
	PPMQPing.RegisterService(service, &ctrl.PPMQPinger{})
	PPMQGetSublist.RegisterService(service, &ctrl.PPMQGetSublister{})
	PPMQUnSub.RegisterService(service, &ctrl.PPMQUnSuber{})
	PPMQSubscribe.RegisterService(service, &ctrl.PPMQSubscribeer{})
	PPMQPublish.RegisterService(service, &ctrl.PPMQPublisher{})

}
