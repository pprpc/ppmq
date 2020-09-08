package app

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap/zapcore"
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"xcthings.com/micro/svc"
	mmodel "github.com/pprpc/ppmq/model"
	g "github.com/pprpc/ppmq/ppmqd/common/global"
	ctrl "github.com/pprpc/ppmq/ppmqd/controller"
	"github.com/pprpc/ppmq/ppmqd/model"
	"github.com/pprpc/core/sess"

	rsub "github.com/pprpc/ppmq/model/redis"
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
	//
	pprofInit()
	if g.EtcdPoint != "" {
		err = etcdReg()
		if err != nil {
			return
		}
	}
	_t, _ := json.Marshal(&g.Conf)
	logs.Logger.Debugf("Config: %s.", _t)
	//
	err = model.InitEngine(g.Conf)
	if err != nil {
		return
	}
	g.Sess = sess.NewSessions(g.PConf.Ppmq.MaxSessions)

	//
	if g.PConf.Redis.Addr != "" {
		g.TopicCache = rsub.NewRSub(g.PConf.Redis)
	}
	err = initAuthConns()
	if err != nil {
		return
	}

	//
	regService()

	// set all offline
	if err = ctrl.SetAllOffline(g.Conf.Public.ServerID); err != nil {
		return
	}
	//
	serverInit()
	//
	ctrl.RunProcess(g.PConf.Ppmq.UDPHbsec)
	//
	go clearTempDB()

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
	reg.Listen = g.Conf.Listen //svc.GetListenURI(g.Conf.Listen)
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

func clearTempDB() {
	var code uint64
	var err error
	for {
		//curMs = g.Ppmq.ClearMsgExpiredTimems
		tms := common.GetTimeMs() - g.PConf.Ppmq.TempdbExpiredms
		code, err = mmodel.ClearTempDB(tms)
		if err != nil {
			logs.Logger.Warnf("mmodel.ClearTempDB(%d), code: %d, %s.", tms, code, err)
		}
		common.Sleep(int(g.PConf.Ppmq.TempdbClearIntervalsec))
	}
}
