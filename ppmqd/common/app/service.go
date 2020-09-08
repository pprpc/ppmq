package app

import (
	"xcthings.com/ppmq/protoc/ppmqd/PPMQConnect"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQDisconnect"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQEXChangeMsg"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQGetClientID"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQGetSublist"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPing"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQSubscribe"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQUnSub"

	g "xcthings.com/ppmq/ppmqd/common/global"
	ctrl "xcthings.com/ppmq/ppmqd/controller"
)

func regService() {
	PPMQGetClientID.RegisterService(g.Service, &ctrl.PPMQGetClientIDer{})
	PPMQConnect.RegisterService(g.Service, &ctrl.PPMQConnecter{})
	PPMQPublish.RegisterService(g.Service, &ctrl.PPMQPublisher{})
	PPMQSubscribe.RegisterService(g.Service, &ctrl.PPMQSubscribeer{})
	PPMQUnSub.RegisterService(g.Service, &ctrl.PPMQUnSuber{})
	PPMQPing.RegisterService(g.Service, &ctrl.PPMQPinger{})
	PPMQDisconnect.RegisterService(g.Service, &ctrl.PPMQDisconnecter{})
	PPMQGetSublist.RegisterService(g.Service, &ctrl.PPMQGetSublister{})
	PPMQEXChangeMsg.RegisterService(g.Service, &ctrl.PPMQEXChangeMsger{})
}
