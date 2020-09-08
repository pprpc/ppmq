package app

import (
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQConnect"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQDisconnect"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQEXChangeMsg"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQGetClientID"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQGetSublist"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQPing"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQPublish"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQSubscribe"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQUnSub"

	g "github.com/pprpc/ppmq/ppmqd/common/global"
	ctrl "github.com/pprpc/ppmq/ppmqd/controller"
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
