package ppmqcli

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	mqc "xcthings.com/ppmq/common"
	errc "xcthings.com/ppmq/common/errorcode"
	clic "xcthings.com/ppmq/ppmqcli/common"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQConnect"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQGetClientID"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPing"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQSubscribe"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
)

// PpmqCli .
type PpmqCli struct {
	URI       string
	URL       *url.URL
	Account   string
	Password  string // md5(textpass)
	HWFeature string
	attr      *clic.PpmqConn
	Conn      pprpc.RPCCliConn

	ClientID        string
	historymsgType  int32
	historymsgCount int32

	outsideIpaddr string
	outsidePort   int32
	hbSec         uint32
}

// NewPpmqcli .
func NewPpmqcli(uri, account, password, hwFeature string) (cli *PpmqCli, err error) {
	cli = new(PpmqCli)

	if account == "" || password == "" || uri == "" || hwFeature == "" {
		err = fmt.Errorf("uri, account, password, hwFeature Can not be empty")
		return
	}
	cli.URL, err = url.ParseRequestURI(uri)
	if err != nil {
		err = fmt.Errorf("uri: %s, %s", uri, err)
		return
	}

	cli.attr = new(clic.PpmqConn)

	cli.URI = uri
	cli.Account = account
	cli.Password = password
	cli.HWFeature = hwFeature

	cli.historymsgType = mqc.ConsumeNew
	cli.historymsgCount = 0

	return
}

// SetRecivePublishCB .
func (cli *PpmqCli) SetRecivePublishCB(cb clic.RecivePublishCallBack) {
	cli.attr.PubCB = cb
}

// SetHistorymsg .
func (cli *PpmqCli) SetHistorymsg(t, c int32) (err error) {
	if t < mqc.ConsumeNew || t > mqc.ConsumeNewAndOld {
		err = fmt.Errorf("history msg type not support: %d", t)
		return
	}
	if c < 0 {
		err = fmt.Errorf("history msg count not support: %d", c)
		return
	}
	cli.historymsgType = t
	cli.historymsgCount = c

	return
}

// Dail .
func (cli *PpmqCli) Dail() (err error) {
	if cli.attr.PubCB == nil {
		err = fmt.Errorf("not set pub callback")
		return
	}
	switch cli.URL.Scheme {
	case "tcp":
		cli.Conn, err = pprpc.Dail(cli.URL, nil, service, 5*time.Second, cli.firstCmd)
	case "udp":
		cli.Conn, err = pprpc.DailUDP(cli.URL.Host, service, 30, cli.firstCmd)
	}
	if err != nil {
		logs.Logger.Warnf("pprpc.Dail(), %s.", err)
		common.Sleep(3)
	}
	cli.Conn.SetAutoHB(false)
	common.SleepMs(500)
	go cli.loopPing()
	return
}

// Close .
func (cli *PpmqCli) Close() (err error) {
	if cli.Conn == nil {
		err = fmt.Errorf("not set conn")
		return
	}
	cli.Conn.Close()
	return
}

// Subscribe .
func (cli *PpmqCli) Subscribe(req *PPMQSubscribe.Req) (resp *PPMQSubscribe.Resp, err error) {
	var pkg *packets.CmdPacket
	rpc := PPMQSubscribe.NewPPMQSubscribeClient(cli.Conn)
	resp, pkg, err = rpc.PPCall(context.Background(), req)
	if err != nil {
		logs.Logger.Errorf("PPMQSubscribe.PPCall(), error: %s.", err)
		return
	}
	if pkg.Code == errc.NOTRUNPPMQCONNECT {
		err = fmt.Errorf("PPMQSubscribe, Code: %d", pkg.Code)
		cli.handleError202()
		return
	} else if pkg.Code != 0 {
		logs.Logger.Errorf("PPMQSubscribe, Code: %d.", pkg.Code)
		err = fmt.Errorf("PPMQSubscribe, Code: %d", pkg.Code)
		return
	}
	return
}

// Publish .
func (cli *PpmqCli) Publish(req *PPMQPublish.Req) (resp *PPMQPublish.Resp, err error) {
	var pkg *packets.CmdPacket
	rpc := PPMQPublish.NewPPMQPublishClient(cli.Conn)
	resp, pkg, err = rpc.PPCall(context.Background(), req)
	if err != nil {
		logs.Logger.Errorf("PPMQPublish.PPCall(), error: %s.", err)
		return
	}
	if pkg.Code == errc.NOTRUNPPMQCONNECT {
		cli.handleError202()
		err = fmt.Errorf("PPMQPublish, Code: %d", pkg.Code)
		return
	} else if pkg.Code != 0 {
		logs.Logger.Errorf("PPMQPublish, Code: %d.", pkg.Code)
		err = fmt.Errorf("PPMQPublish, Code: %d", pkg.Code)
		return
	}
	return
}

// firstCmd .
func (cli *PpmqCli) firstCmd(c pprpc.RPCCliConn) {
	if cli.ClientID == "" {
		err := cli.getClientID(c)
		if err != nil {
			logs.Logger.Errorf("GetClientID(), error: %s.", err)
			return
		}
	}
	cli.connect(c)
	if cli.attr.PubCB != nil {
		c.SetAttr(cli.attr)
	}
}

func (cli *PpmqCli) loopPing() {
	for {
		if cli.hbSec == 0 {
			common.SleepMs(500)
			continue
		}

		var tCli *pprpc.TCPCliConn
		var uCli *pprpc.UDPCliConn
		switch cli.Conn.(type) {
		case *pprpc.TCPCliConn:
			tCli = cli.Conn.(*pprpc.TCPCliConn)
		case *pprpc.UDPCliConn:
			uCli = cli.Conn.(*pprpc.UDPCliConn)
		}

		if tCli != nil {
			select {
			case <-tCli.Ctx.Done():
				logs.Logger.Warn("cli.Ctx.Done(), waiting reconnect.")
				cli.hbSec = 0
				goto ForEnd
			case <-time.After(time.Second * time.Duration(cli.hbSec)):
				cli.ping()
			}
		} else if uCli != nil {
			select {
			case <-uCli.Ctx.Done():
				logs.Logger.Warn("cli.Ctx.Done(), waiting reconnect.")
				cli.hbSec = 0
				goto ForEnd
			case <-time.After(time.Second * time.Duration(cli.hbSec)):
				cli.ping()
			}
		}
	ForEnd:
	}
	//Stop:
}

func (cli *PpmqCli) ping() {
	req := new(PPMQPing.Req)
	req.IsSleep = false

	rpc := PPMQPing.NewPPMQPingClient(cli.Conn)
	resp, pkg, err := rpc.PPCall(context.Background(), req)
	if err != nil {
		logs.Logger.Errorf("PPMQPing.PPCall(), error: %s.", err)
		return
	}
	if pkg.Code == errc.NOTRUNPPMQCONNECT {
		cli.handleError202()
		return
	} else if pkg.Code != 0 {
		logs.Logger.Errorf("PPMQPing, Code: %d.", pkg.Code)
		return
	}
	cli.outsideIpaddr = resp.OutsideIpaddr
	logs.Logger.Debugf("PPMQPing OK.")
}

func (cli *PpmqCli) getClientID(c pprpc.RPCCliConn) (err error) {
	req := new(PPMQGetClientID.Req)
	req.Account = cli.Account
	req.HardwareInfo = cli.HWFeature

	dc := PPMQGetClientID.NewPPMQGetClientIDClient(c)
	resp, pkg, e := dc.PPCall(context.Background(), req)
	if e != nil {
		logs.Logger.Errorf("PPMQGetClientID.PPCall(), error: %s.", e)
		err = e
		logs.Logger.Debugf("Encrypt: %d .", pkg.EncType)
		return
	}
	cli.ClientID = resp.ClientId
	logs.Logger.Debugf("PPMQGetClientID, ClientID: %s.", resp.ClientId)
	return
}

func (cli *PpmqCli) connect(c pprpc.RPCCliConn) {
	if cli.ClientID == "" {
		logs.Logger.Warnf("ClientID not empty.")
		return
	}
	connect := new(PPMQConnect.Req)
	connect.ClientId = cli.ClientID
	connect.Account = cli.Account
	connect.Password = cli.Password
	connect.HistorymsgType = cli.historymsgType
	connect.HistorymsgCount = cli.historymsgCount
	connect.ClearSession = mqc.SESSIONKEEP
	connect.WillFlag = 0
	connect.WillQos = 1
	connect.WillRetain = 0
	connect.WillTopic = ""
	connect.WillBody = ""

	creq := PPMQConnect.NewPPMQConnectClient(c)
	resp, _, err := creq.PPCall(context.Background(), connect)
	if err != nil {
		logs.Logger.Errorf("PPMQConnect.PPCall(), error: %s.", err)
		cli.handleError202()
		return
	}
	cli.hbSec = resp.HbInterval
	logs.Logger.Debugf("PPMQConnect, hbSec: %d.", resp.HbInterval)
	return
}

func (cli *PpmqCli) handleError202() {
	common.SleepMs(800)
	cli.firstCmd(cli.Conn)
}
