package app

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"strconv"

	"xcthings.com/hjyz/logs"
	"xcthings.com/micro/svc"
	mqc "xcthings.com/ppmq/common"
	g "xcthings.com/ppmq/ppmqd/common/global"
	ctrl "xcthings.com/ppmq/ppmqd/controller"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
	"xcthings.com/pprpc/pptcp"
	"xcthings.com/pprpc/ppudp"
)

func serverInit() {
	var err error
	for _, lis := range g.Conf.Listen {
		err = runServer(lis)
		if err != nil {
			logs.Logger.Errorf("runServer(lis), error: %s.", err)
		}
	}
}

func runServer(lis svc.LisConf) error {
	u, e := url.ParseRequestURI(lis.URI)
	if e != nil {
		return e
	}
	switch u.Scheme {
	case "udp":
		p := u.Port()
		_t, e := strconv.Atoi(p)
		if e != nil {
			return e
		}
		usrv, err := pprpc.NewRPCUDPServer(u.Hostname(), int(_t), int(g.PConf.MaxSession))
		if err != nil {
			return fmt.Errorf("pprpc.NewRPCUDPServer(), error: %s", err)
		}
		usrv.HBCB = cb
		usrv.DisconnectCB = closeUDP
		// DisconnectCB func(RPCConn)
		usrv.Service = g.Service
		usrv.SetReadTimeout(lis.ReadTimeout)
		logs.Logger.Infof("Listen UDPServer: %s.", lis.URI)
		go usrv.Serve()
		udplis = append(udplis, usrv)
	default:
		var tlsc *tls.Config
		if lis.TLSCrt != "" && lis.TLSKey != "" {
			tlsc, e = pprpc.GetTLSConfig(lis.TLSCrt, lis.TLSKey)
			if e != nil {
				return fmt.Errorf("pprpc.GetTLSConfig(), %s", e)
			}
		} else {
			tlsc = nil
		}

		srv, err := pprpc.NewRPCTCPServer(u, tlsc)
		if err != nil {
			return fmt.Errorf("pprpc.NewRPCTCPServer(), error: %s", err)
		}
		srv.HBCB = cb
		srv.DisconnectCB = closeTCP
		srv.Service = g.Service
		srv.SetReadTimeout(lis.ReadTimeout)
		logs.Logger.Infof("Listen TCPServer: %s.", lis.URI)
		go srv.Serve()
		tcplis = append(tcplis, srv)
	}
	return nil
}

//
func cb(pkg *packets.HBPacket, c pprpc.RPCConn) error {
	logs.Logger.Debugf("%s, HBPacket, MessageType: %d.", c, pkg.MessageType)
	_, err := pkg.Write(c)
	return err
}

func closeUDP(conn *ppudp.Connection) {
	disconnectCB(conn)
}

func closeTCP(conn *pptcp.Connection) {
	disconnectCB(conn)
}

func disconnectCB(c pprpc.RPCConn) {

	_t, e := c.GetAttr()
	if e != nil {
		logs.Logger.Debugf("c.GetAttr(), error: %s.", e)
		return
	}
	if _t == nil {
		logs.Logger.Debugf("conn.attr is nil.")
		return
	}
	ci := _t.(*g.ConnAttr)

	if ci.ClearSession == mqc.SESSIONCLEAR {
		ctrl.ClearSub(ci.ClientID)
		ctrl.UpdateStatus(ci.ClientID, mqc.MSGSTATUSClearSession)
	}
	if ci.WillFlag == 1 {
		ctrl.OfflineWillMsg(ci.ClientID)
	}

	if ci.OfflineEvent == true {
		ctrl.OfflineEvent(c)
	}
	// clear timer
	ci.ClearTimer()

	v, e := g.Sess.Get(ci.ClientID)
	if e == nil {
		if v.(pprpc.RPCConn).RemoteAddr() == c.RemoteAddr() && v.(pprpc.RPCConn).Type() == c.Type() {
			logs.Logger.Infof("%s, Sess.Remove(%s).", c, ci.ClientID)
			g.Sess.Remove(ci.ClientID)

			// FIXME: panic: close of closed channel
			// if ci.IsAuth == true && ci.ClearSession == mqc.SESSIONKEEP && g.PConf.Ppmq.MessageOrder == 1 {
			// 	// M MQTT
			// 	// S TLS
			// 	// T TCP
			// 	// Q QUIC
			// 	// U UDP
			// 	if c.Type() == "M" {
			// 		close(ci.MqttCmdChan)
			// 		// } else if c.Type() == "S" || c.Type() == "T" || c.Type() == "Q" {
			// 	} else {
			// 		close(ci.CmdChan)
			// 	}
			// }
		} else {
			logs.Logger.Infof("%s, ClientID: %s , conn not match, not remove.", c, ci.ClientID)
		}
	}

}
