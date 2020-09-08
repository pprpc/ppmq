package main

import (
	"flag"
	"fmt"

	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"xcthings.com/ppmq/ppmqonline-ms/common/app"
	g "xcthings.com/ppmq/ppmqonline-ms/common/global"
)

var (
	bDate       string //
	ciHash      string //
	mainVersion string = "0.0.9"
)

var (
	etcdipaddr = flag.String("ipaddr", "127.0.0.1:2379", "etcd server ipaddr")
	region     = flag.String("region", "cn-shenzhen", "region name")
	ethname    = flag.String("i", "eth0", "set network device ")
	dbs        = flag.String("dbs", "ppmq", "set db name")

	ver      = flag.Bool("v", false, "show verion")
	confFile = flag.String("conf", "../conf/server.json", "set conf file path")
)

func main() {
	flag.Parse()
	// verion
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
	g.MSName = "ppmqonline"
	g.Dbs = *dbs

	err := app.Run(*confFile)
	if err != nil {
		logs.Logger.Errorf("app.Run(), error: %s.", err)
	}
	common.WaitCtrlC()
}
