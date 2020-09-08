package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pprpc/util/logs"
	"xcthings.com/micro/pprpcpool"
	"xcthings.com/micro/svc"
	g "github.com/pprpc/ppmq/ppmqd/common/global"
	"github.com/pprpc/ppmq/protoc/ppmqd/PPMQEXChangeMsg"
	"xcthings.com/protoc/authdevice/CheckDevice"
	"xcthings.com/protoc/authdevice/CheckDevicePPMQPub"
	"xcthings.com/protoc/authdevice/CheckDevicePPMQSub"
	"xcthings.com/protoc/authuser/CheckUser"
	"xcthings.com/protoc/authuser/CheckUserPPMQPub"
	"xcthings.com/protoc/authuser/CheckUserPPMQSub"

	ctrl "github.com/pprpc/ppmq/ppmqd/controller"
)

func initAuthConns() (err error) {
	authRegService()
	err = initMicroClientConn()
	if err != nil {
		return
	}
	err = etcdWatcher()
	if err != nil {
		return
	}
	// load
	loadRegister()
	return
}

func authRegService() {
	CheckDevice.RegisterService(g.AuthService, nil)
	CheckUser.RegisterService(g.AuthService, nil)
	CheckDevicePPMQPub.RegisterService(g.AuthService, nil)
	CheckDevicePPMQSub.RegisterService(g.AuthService, nil)
	CheckUserPPMQPub.RegisterService(g.AuthService, nil)
	CheckUserPPMQSub.RegisterService(g.AuthService, nil)
	PPMQEXChangeMsg.RegisterService(g.EXChangeMsgService, &ctrl.PPMQEXChangeMsger{})
}

func initMicroClientConn() (err error) {
	g.MicrosConn = pprpcpool.NewMicroClientConn(g.AuthService)
	for _, v := range g.PConf.Micros {
		err = g.MicrosConn.AddMicro(v.Name)
		if err != nil {
			err = fmt.Errorf("g.MicrosConn.AddMicro(%s), %s", v.Name, err)
			return
		}
	}
	if g.PConf.Ppmq.Mode == "cluster" {
		g.EXChangeConn = pprpcpool.NewMicroClientConn(g.EXChangeMsgService)
		err = g.EXChangeConn.AddMicro(g.MSName)
		if err != nil {
			err = fmt.Errorf("g.EXChangeConn.AddMicro(%s), %s", g.MSName, err)
			return
		}
	}
	return
}

func loadRegister() {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	key := fmt.Sprintf("/register/%s/", g.Region)
	kvs, err := g.SvcWatcher.GetValues(ctx, key)
	if err != nil {
		logs.Logger.Errorf("g.SvcAgent.GetValues(%s), %s.", key, err)
		return
	}

	for _, row := range kvs {
		_t := new(svc.ValueRegService)
		err := json.Unmarshal([]byte(row.Value), _t)
		if err != nil {
			logs.Logger.Errorf("json.Unmarshal(), %s.", err)
			continue
		}
		err = g.MicrosConn.AddHost(key, *_t)
		if err != nil {
			logs.Logger.Errorf("g.MicroConn.AddHost(%s), %s.", key, err)
		}
		if g.PConf.Ppmq.Mode == "cluster" {
			if _t.Name == g.MSName {

				if _t.LanIP != g.LanIP {
					err = g.EXChangeConn.AddHost(key, *_t)
					if err != nil {
						logs.Logger.Errorf("g.EXChangeConn.AddHost(%s), %s.", key, err)
					}
					return
				}
			}
		}
	}
}

func etcdWatcher() (err error) {
	var ep []string
	ep = append(ep, g.EtcdPoint)
	g.SvcWatcher, err = svc.NewWatcher("/register", ep, regcb)
	if err != nil {
		err = fmt.Errorf("svc.NewWatcher(/register), %s", err)
		return
	}
	go g.SvcWatcher.Start()

	return
}

func regcb(action, key, value string) {
	if action == "PUT" {
		_t := new(svc.ValueRegService)
		err := json.Unmarshal([]byte(value), _t)
		if err != nil {
			logs.Logger.Errorf("json.Unmarshal(), %s.", err)
			return
		}
		err = g.MicrosConn.AddHost(key, *_t)
		if err != nil {
			logs.Logger.Errorf("g.MicroConn.AddHost(%s), %s.", key, err)
		}
		if g.PConf.Ppmq.Mode == "cluster" {
			if _t.Name == g.MSName {

				if _t.LanIP != g.LanIP {
					err = g.EXChangeConn.AddHost(key, *_t)
					if err != nil {
						logs.Logger.Errorf("g.EXChangeConn.AddHost(%s), %s.", key, err)
					}
					return
				}
			}
		}
	} else {
		err := g.MicrosConn.DelHost(key)
		if err != nil {
			logs.Logger.Errorf("g.MicroConn.DelHost(%s), %s.", key, err)
		}
		if g.PConf.Ppmq.Mode == "cluster" {
			err = g.EXChangeConn.DelHost(key)
			if err != nil {
				logs.Logger.Errorf("g.EXChangeConn.AddHost(%s), %s.", key, err)
			}
		}
	}
}
