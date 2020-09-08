package global

import (
	"encoding/json"
	"fmt"
	"strings"

	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	"xcthings.com/micro/svc"
)

// PrivateConf private conf
type PrivateConf struct {
	MaxSession int32 `json:"max_session"`
}

// LoadConf load conf
func LoadConf(filePath string) (conf svc.MSConfig, err error) {
	if EtcdPoint == "" {
		conf, err = loadConfFromFile(filePath)
		if err != nil {
			return
		}
		err = json.Unmarshal(conf.PrivateConfig, &PConf)
		if err != nil {
			err = fmt.Errorf("LoadConf, loadConfFromFile.PrivateConfig json.Unmarshal, %s", err)
			return
		}
		return
	}

	conf, err = loadConfFromETCD()
	if err != nil {
		logs.Logger.Warnf("LoadConf, loadConfFromETCD(), %s.", err)
		//fmt.Fprintf(os.Stderr, "LoadConf, loadConfFromETCD(), %s.", err)
		conf, err = loadConfFromFile(filePath)
		if err != nil {
			return
		}
		err = json.Unmarshal(conf.PrivateConfig, &PConf)
		if err != nil {
			err = fmt.Errorf("LoadConf, loadConfFromFile.PrivateConfig json.Unmarshal, %s", err)
			return
		}
	} else {
		logs.Logger.Debug("load configure from etcd ok.")
	}

	return
}

func loadConfFromETCD() (conf svc.MSConfig, err error) {
	var ep []string
	ep = append(ep, EtcdPoint)
	var ag *svc.Agent
	var cfg *svc.Config

	ag, err = svc.NewAgent(svc.ValueRegService{}, 5, ep)
	if err != nil {
		err = fmt.Errorf("svc.NewAgent(), %s", err)
		return
	}
	defer ag.Close()

	ips, e := common.GetIPAddrByName(Ethname)
	if e != nil {
		err = fmt.Errorf("common.GetIPAddrByName(%s), %s", Ethname, e)
		return
	}

	dbs := strings.Split(Dbs, ",")

	cfg, err = svc.NewConfig(ag, Region, ips[0], MSName, dbs, true)
	if err != nil {
		err = fmt.Errorf("svc.NewConfig(), %s", err)
		return
	}

	err = cfg.GetAll()
	if err != nil {
		err = fmt.Errorf("cfg.GetAll(), %s", err)
		return
	}
	conf = *cfg.Conf

	err = json.Unmarshal(conf.PrivateConfig, &PConf)

	return
}

func loadConfFromFile(filePath string) (conf svc.MSConfig, err error) {
	if common.PathIsExist(filePath) != true {
		err = fmt.Errorf("conf file not exist")
		return
	}
	var buf []byte
	buf, err = common.LoadFileToByte(filePath)
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, &conf)
	err = json.Unmarshal(conf.PrivateConfig, &PConf)
	return
}
