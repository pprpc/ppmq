package global

import (
	"container/list"
	"fmt"
	"time"

	"github.com/pprpc/util/cache"
	"github.com/pprpc/util/logs"
	"xcthings.com/micro/pprpcpool"
	"xcthings.com/micro/svc"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
	"github.com/pprpc/core/sess"

	rsub "xcthings.com/ppmq/model/redis"
)

// ConnAttr .
type ConnAttr struct {
	ClientID     string
	Account      string
	IsAuth       bool
	UserID       int64
	ConnectFlag  bool
	OfflineEvent bool

	EncType     uint8
	MessageType uint8

	ClearSession int32
	WillFlag     uint32
	WillQos      uint32
	WillRetain   uint32
	WillTopic    string
	WillBody     string
	HbInterval   uint32

	HistorymsgType  int32
	HistorymsgCount int32

	ConnStatus int32
	ConnType   int32
	LastTime   int64
	ISSleep    int32

	RemoteIpaddr string
	// time.Timer
	Timers            *cache.Cache
	CmdChan           chan *packets.CmdPacket
	HistoryCmdList    *list.List
	HistoryPubackSync chan uint16
	RunSub            bool
	SendOfflinemsgEnd bool
}

type HistoryList struct {
	MsgID  string
	CmdPkg *packets.CmdPacket
}

// Conf .
var Conf svc.MSConfig
var PConf PrivateConf
var EtcdPoint, Region, Ethname, MSName, Dbs, LanIP string
var SvcAgent *svc.Agent
var SvcWatcher *svc.Watcher

// Service，AuthService global service
var Service, AuthService, EXChangeMsgService *pprpc.Service

// MicrosConn micro service connections.
var MicrosConn, EXChangeConn *pprpcpool.MicroClientConn

// Sess  存放所有连接
var Sess *sess.Sessions

var TopicCache *rsub.RSub

func init() {
	Service = pprpc.NewService()
	AuthService = pprpc.NewService()
	EXChangeMsgService = pprpc.NewService()
}

//ClearTimer clear timer
func (attr *ConnAttr) ClearTimer() {
	if attr.Timers == nil {
		return
	}
	fn := func(k, v interface{}) bool {
		v.(*time.Timer).Stop()
		return true
	}
	attr.Timers.Range(fn)
	logs.Logger.Debugf("ClientID: %s, ClearTimer OK", attr.ClientID)
}

// StartTimer start timer
func (attr *ConnAttr) StartTimer(d time.Duration, seq uint64, fn func()) {
	if attr.Timers == nil {
		return
	}

	t := time.AfterFunc(d*time.Millisecond, fn)
	attr.Timers.Set(fmt.Sprintf("%d", seq), t)
}

// StopTimer stop timer by key.
func (attr *ConnAttr) StopTimer(key string) {
	if attr.Timers == nil {
		return
	}

	v, e := attr.Timers.Get(key)
	if e != nil {
		return
	}
	v.(*time.Timer).Stop()
	attr.Timers.Del(key)
	logs.Logger.Debugf("ClientID: %s, Key: %s, StopTimer OK", attr.ClientID, key)
}
