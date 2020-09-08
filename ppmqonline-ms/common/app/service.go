package app

import (
	g "xcthings.com/ppmq/ppmqonline-ms/common/global"
	"xcthings.com/ppmq/protoc/ppmqonline/OnlineState"

	ctrl "xcthings.com/ppmq/ppmqonline-ms/controller"
)

func regService() {
	OnlineState.RegisterService(g.Service, &ctrl.OnlineStateer{})
}
