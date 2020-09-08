package ppmqcli

import (
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQConnect"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQDisconnect"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQGetClientID"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQGetSublist"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQPing"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQPublish"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQSubscribe"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQUnSub"
	"github.com/pprpc/core"

	ctrl "github.com/pprpc/ppmq/ppmqcli/controller"
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
