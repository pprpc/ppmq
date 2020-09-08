package logic

import (
	"fmt"

	"xcthings.com/hjyz/logs"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQConnect"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"

	mqc "xcthings.com/ppmq/common"
	errc "xcthings.com/ppmq/common/errorcode"
	pm "xcthings.com/ppmq/model"
	g "xcthings.com/ppmq/ppmqd/common/global"
	m "xcthings.com/ppmq/ppmqd/model"
)

// LPPMQConnect PPMQConnect Business logic
func LPPMQConnect(c pprpc.RPCConn, pkg *packets.CmdPacket, req *PPMQConnect.Req) (resp *PPMQConnect.Resp, code uint64, err error) {
	if req.Account == "" || req.Password == "" {
		code = errc.ParameterError
		err = fmt.Errorf("Incomplete parameters: Account, Password")
		return
	}
	_, err = m.CheckAccouont(req.Account, req.Password)
	if err != nil {
		code = errc.ACCOUNTNOTMATCH
		return
	}

	if req.ClientId == "" {
		code = errc.ParameterError
		err = fmt.Errorf("Incomplete parameters: ClientId")
		return
	}

	//
	cid := new(pm.Clientid)
	cid.ClientID = req.ClientId
	code, _ = cid.Get()
	if code == errc.NOTEXISTRECORD {
		cid.Account = req.Account
		cid.HwFeature = "defult"
		code, err = cid.Add()
		if err != nil {
			return
		}
	}

	if req.ClearSession != mqc.SESSIONKEEP && req.ClearSession != mqc.SESSIONCLEAR {
		code = errc.ParameterIllegal
		err = fmt.Errorf("The parameter is invalid, ClearSession: %d", req.ClearSession)
		return
	}
	if req.HistorymsgType < mqc.ConsumeNew || req.HistorymsgType > mqc.ConsumeNewAndOld {
		code = errc.ParameterIllegal
		err = fmt.Errorf("The parameter is invalid, HistorymsgType: %d", req.HistorymsgType)
		return
	}
	if req.HistorymsgCount < 0 {
		code = errc.ParameterIllegal
		err = fmt.Errorf("The parameter is invalid: HistorymsgType")
		return
	}
	// 190324: sso enable func
	// 190508： disable
	// if g.PConf.Ppmq.SSOEnable == true {
	// 	err = ssoEnable(req.Account)
	// 	if err != nil {
	// 		logs.Logger.Warnf("ssoEnable(%s), %s.", req.Account, err)
	// 	}
	// }
	// FIXME: 忽略遗嘱的处理，后续再检查输入合法性

	resp, code, err = m.MPPMQConnect(c, req)
	if err == nil {
		if c.Type() == "U" {
			resp.HbInterval = g.PConf.Ppmq.UDPHbsec
		} else {
			resp.HbInterval = g.PConf.Ppmq.TCPHbsec
		}
	}
	return
}

func ssoEnable(account string) (err error) {
	var rows []*pm.Connection
	rows, _, err = pm.GetRowsByAccount(account)
	if err != nil {
		return
	}
	for _, v := range rows {
		// 1, delete
		if v.ISOnline == mqc.OFFLINE {
			_, err = pm.DeleteConnection(v.ClientID)
			if err != nil {
				err = fmt.Errorf("pm.DeleteConnection(%s), %s", v.ClientID, err)
				return
			}
			logs.Logger.Debugf("pm.DeleteConnection(%s), account: %s(offline), clear.", v.ClientID, account)
			continue
		}

		logs.Logger.Warnf("SSOEnable, cid: %s, account: %s(online), clear and closeconn.", v.ClientID, account)
		// 2, get connection and close
		if v.ServerID == g.Conf.Public.ServerID {
			_v, lerr := g.Sess.Get(v.ClientID)
			if lerr != nil {
				err = fmt.Errorf("g.Sess.Get(%s,), %s", v.ClientID, lerr)
				return
			}

			// Close
			_v.(pprpc.RPCConn).Close()
		} else {
			// FIXME： notice other serverid, close conn;
		}
	}
	return
}
