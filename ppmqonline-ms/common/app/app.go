package app

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap/zapcore"
	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	"xcthings.com/micro/svc"
	g "xcthings.com/ppmq/ppmqonline-ms/common/global"
	"xcthings.com/ppmq/ppmqonline-ms/model"
)

// Run start app
func Run(cp string) (err error) {

	g.Conf, err = g.LoadConf(cp)
	if err != nil {
		logs.Logger.Errorf("g.LoadConf(%s), error: %s.", cp, err)
		return
	}
	//
	logc := g.Conf.Log
	if logc.File != "" {
		logs.Logger.SetLogFile(logc.File, logc.MaxSize, logc.MaxBackups, logc.MaxAge, logc.Caller)
	}
	var lev zapcore.Level
	lev, err = getLevel(logc.Level)
	if err != nil {
		return
	}
	logs.Logger.SetLevel(lev)

	// register service
	err = etcdReg()
	if err != nil {
		return
	}
	_t, _ := json.Marshal(&g.Conf)
	logs.Logger.Debugf("Config: %s.", _t)

	//
	pprofInit()
	//
	err = model.InitEngine(g.Conf)
	if err != nil {
		return
	}
	//

	//
	regService()

	//
	err = serverInit()

	return
}

func getLevel(l int8) (lev zapcore.Level, err error) {
	switch l {
	case -1:
		lev = zapcore.DebugLevel
	case 0:
		lev = zapcore.InfoLevel
	case 1:
		lev = zapcore.WarnLevel
	case 2:
		lev = zapcore.ErrorLevel
	case 3:
		lev = zapcore.DPanicLevel
	case 4:
		lev = zapcore.PanicLevel
	case 5:
		lev = zapcore.FatalLevel
	default:
		err = fmt.Errorf("not support: %d", l)
	}
	return
}

func etcdReg() (err error) {
	if g.EtcdPoint == "" {
		logs.Logger.Warn("not set etcdpoint.")
		return
	}

	var ep, ips []string
	ep = append(ep, g.EtcdPoint)
	var reg svc.ValueRegService
	reg.Region = g.Region
	reg.Listen = g.Conf.Listen
	reg.Name = g.MSName
	reg.ResSrv = svc.GetListenResID(g.Conf.Listen)
	ips, err = common.GetIPAddrByName(g.Ethname)
	if err != nil {
		err = fmt.Errorf("common.GetIPAddrByName(%s), %s", g.Ethname, err)
		return
	}
	reg.LanIP = ips[0]

	g.SvcAgent, err = svc.NewAgent(reg, 5, ep)
	if err != nil {
		err = fmt.Errorf("svc.NewAgent(), %s", err)
		return
	}
	go g.SvcAgent.Start()

	return
}
