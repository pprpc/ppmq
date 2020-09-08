package global

import (
	"encoding/json"
	"fmt"
	"strings"

	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	"xcthings.com/micro/svc"
)

// MicroClient micrl service client
type MicroClient struct {
	Name string   `json:"name,omitempty"`
	URIS []string `json:"uris,omitempty"`
}

// ClearSameConf clear same id
type ClearSameConf struct {
	Clear        bool     `json:"clear,omitempty"`
	WhiteAccount []string `json:"white_account,omitempty"`
}

// PpmqConf .
type PpmqConf struct {
	OnlineNotify      bool   `json:"online_notify"`
	OfflineNotify     bool   `json:"offline_notify"`
	UDPHbsec          uint32 `json:"udp_hbsec"`
	TCPHbsec          uint32 `json:"tcp_hbsec"`
	NotifyTopicPrefix string `json:"notify_topic_prefix"` // Online Offline 通知事件的前缀 Prefix/clientid
	//ServerID          string `json:"server_id"`
	MaxSessions              int32  `json:"max_sessions"`
	Mode                     string `json:"mode"`
	ClusterSignkey           string `json:"cluster_signkey,omitempty"`
	Qos                      bool   `json:"qos,omitempty"`
	SSOEnable                bool   `json:"sso_enable,omitempty"`
	UDPRespTimeoutms         int64  `json:"udp_resp_timeoutms,omitempty"`
	OfflinemsgTimeoutms      int64  `json:"offlinemsg_timeoutms,omitempty"`
	OfflinemsgSendSleepms    int64  `json:"offlinemsg_send_sleepms,omitempty"`
	OfflinemsgSendIntervalms int64  `json:"offlinemsg_send_intervalms,omitempty"`
	MessageOrder             int64  `json:"message_order,omitempty"`
	MessageOrderLength       int    `json:"message_order_length,omitempty"`
	TempdbExpiredms          int64  `json:"tempdb_expiredms,omitempty"`
	TempdbClearIntervalsec   int64  `json:"tempdb_clear_intervalsec,omitempty"`
}

// PrivateConf private conf
type PrivateConf struct {
	Ppmq         PpmqConf      `json:"ppmq"`
	MaxSession   int32         `json:"max_session"`
	DevicePrefix string        `json:"device_prefix,omitempty"`
	Auth         int32         `json:"auth"` // 1 local db; 2 auth micro service.
	CheckTopic   int32         `json:"check_topic"`
	Micros       []MicroClient `json:"micros,omitempty"` // authuser, authdevice
	Redis        svc.RedisConf `json:"redis"`
	ClearDIDConf ClearSameConf `json:"clear_sameid"`
}

// LoadConf load conf
func LoadConf(filePath string) (conf svc.MSConfig, err error) {
	defer func() {
		setConfDefault(&conf, &PConf)
	}()

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
	var ep, ips []string
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
	logs.Logger.Debugf("ips: %v.", ips)

	dbs := strings.Split(Dbs, ",")
	LanIP = ips[0]

	cfg, err = svc.NewConfig(ag, Region, LanIP, MSName, dbs, true)
	// if MSName == "ppmqd" {
	// 	cfg, err = svc.NewConfig(ag, Region, ips[0], MSName, []string{"ppmq"}, true)
	// } else {
	// 	cfg, err = svc.NewConfig(ag, Region, ips[0], MSName, []string{"localmq"}, true)
	// }
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

func setConfDefault(conf *svc.MSConfig, p *PrivateConf) {
	if p.Ppmq.UDPRespTimeoutms < 500 && p.Ppmq.UDPRespTimeoutms != 0 {
		p.Ppmq.UDPRespTimeoutms = 500
	}
	// if p.DevicePrefix == "" {
	// 	p.DevicePrefix = "PP"
	// }

	return
}
