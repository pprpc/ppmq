package app

import (
	g "github.com/pprpc/ppmq/ppmqonline-ms/common/global"
	"github.com/pprpc/ppmq/protoc/ppmqonline/OnlineState"

	ctrl "github.com/pprpc/ppmq/ppmqonline-ms/controller"
)

func regService() {
	OnlineState.RegisterService(g.Service, &ctrl.OnlineStateer{})
}
