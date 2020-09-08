package main

import (
	"flag"
	"fmt"

	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	"xcthings.com/ppmq/ppmqd/common/app"
	g "xcthings.com/ppmq/ppmqd/common/global"
	ctrl "xcthings.com/ppmq/ppmqd/controller"
)

var (
	bDate       string             //
	ciHash      string             //
	mainVersion string = "1.02.01" //
)

var (
	etcdipaddr = flag.String("ipaddr", "127.0.0.1:2379", "etcd server ipaddr")
	region     = flag.String("region", "cn-shenzhen", "region name")
	ethname    = flag.String("i", "eth0", "network device name")
	msname     = flag.String("msname", "ppmqd", "micro name: ppmqd, localmqd")
	dbs        = flag.String("dbs", "ppmq", "dbs")

	ver      = flag.Bool("v", false, "show version")
	confFile = flag.String("conf", "../conf/server.json", "Specify configuration files")
)

func main() {
	flag.Parse()
	// show version
	if *ver {
		version := mainVersion
		if len(bDate) > 0 {
			version += ("+" + bDate)
		}
		fmt.Println("version:", version)

		if len(ciHash) > 0 {
			fmt.Println("git commit hash:", ciHash)
		}
		return
	}
	defer logs.Logger.Flush()

	g.EtcdPoint = *etcdipaddr
	g.Region = *region
	g.Ethname = *ethname
	g.MSName = *msname
	g.Dbs = *dbs

	err := app.Run(*confFile)
	if err != nil {
		logs.Logger.Errorf("app.Run(), error: %s.", err)
	} else {
		defer ctrl.SetAllOffline(g.Conf.Public.ServerID)
	}
	common.WaitCtrlC()
}
