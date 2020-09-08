package global

import (
	"xcthings.com/micro/svc"
	"github.com/pprpc/core"
)

// Conf .
var Conf svc.MSConfig
var PConf PrivateConf
var EtcdPoint, Region, Ethname, MSName, Dbs string

var SvcAgent *svc.Agent

// Service global service
var Service *pprpc.Service

func init() {
	Service = pprpc.NewService()
}
