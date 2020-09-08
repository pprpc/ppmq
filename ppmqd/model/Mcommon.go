package model

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/pprpc/util/common"
	"github.com/pprpc/util/crypto"
	"github.com/pprpc/util/logs"
	mqc "github.com/pprpc/ppmq/common"
	pm "github.com/pprpc/ppmq/model"
	g "github.com/pprpc/ppmq/ppmqd/common/global"
	"xcthings.com/protoc/authdevice/CheckDevice"
	"xcthings.com/protoc/authdevice/CheckDevicePPMQPub"
	"xcthings.com/protoc/authdevice/CheckDevicePPMQSub"
	"xcthings.com/protoc/authuser/CheckUser"
	"xcthings.com/protoc/authuser/CheckUserPPMQPub"
	"xcthings.com/protoc/authuser/CheckUserPPMQSub"
)

type TopicAns struct {
	Topic string
	Code  uint32
}

type OnlineInfo struct {
	ClientID string
	ServerID string
}

// CheckAccouont .
func CheckAccouont(account, password string) (bool, error) {
	if g.PConf.Auth == 1 {
		q := new(pm.Account)
		if g.PConf.DevicePrefix != "" {
			if account[:2] == g.PConf.DevicePrefix {
				q.Account = account
			} else {
				q.UserID = mqc.GetUserID(account)
				if q.UserID == 0 {
					return false, fmt.Errorf("Acccount: %s[%s], format error", account, account[:2])
				}
			}
		} else {
			if mqc.GetUserID(account) != 0 {
				q.UserID = mqc.GetUserID(account)
			} else {
				q.Account = account
			}
		}
		q.Password = password

		_, err := q.GetAdv()
		if err != nil {
			return false, err
		}
	} else {
		// auth micro service
		var flag int
		if g.PConf.DevicePrefix != "" {
			if account[:2] == g.PConf.DevicePrefix {
				flag = 1
			} else {
				flag = 2
			}
		} else {
			if mqc.GetUserID(account) != 0 {
				flag = 2
			} else {
				flag = 1
			}
		}
		logs.Logger.Debugf("account: %s, flag: %d, DevicePrefix: [%s].", account, flag, g.PConf.DevicePrefix)
		if flag == 1 {
			req := new(CheckDevice.Req)
			req.Did = account
			req.DidSignkey = password
			pkg, resp, err := g.MicrosConn.Invoke(context.Background(), CheckDevice.Module, CheckDevice.CmdID, req)
			if err != nil {
				err = fmt.Errorf("g.MicrosConn.Invoke(CheckDevice), %s", err)
				return false, err
			}
			if pkg.Code != 0 {
				return false, fmt.Errorf("rpc code: %d", pkg.Code)
			}
			if resp.(*CheckDevice.Resp).Code != 0 {
				return false, fmt.Errorf("Answer Code: %d", resp.(*CheckDevice.Resp).Code)
			}
		} else {
			uid := mqc.GetUserID(account)
			if uid == 0 {
				return false, fmt.Errorf("Acccount: %s[%s]，format error", account, account[:2])
			}
			req := new(CheckUser.Req)
			req.UserId = uid
			req.Password = password
			req.CountryCode = ""
			req.AccessKey = ""
			pkg, resp, err := g.MicrosConn.Invoke(context.Background(), CheckUser.Module, CheckUser.CmdID, req)
			if err != nil {
				err = fmt.Errorf("g.MicrosConn.Invoke(CheckUser), %s", err)
				return false, err
			}
			if pkg.Code != 0 {
				return false, fmt.Errorf("rpc code: %d", pkg.Code)
			}
			if resp.(*CheckUser.Resp).Code != 0 {
				return false, fmt.Errorf("Answer Code: %d", resp.(*CheckUser.Resp).Code)
			}
		}
	}
	return true, nil
}

// SetOffline .
func SetOffline(clientID string) error {
	q := new(pm.Connection)
	q.ClientID = clientID
	q.LastTime = common.GetTimeMs()
	q.ISOnline = mqc.OFFLINE
	q.ServerID = g.Conf.Public.ServerID

	_, err := q.SetOffline()
	return err
}

func getClientID(account, hwf string) string {
	return crypto.MD5([]byte(fmt.Sprintf("%s-%s", account, hwf)))
}

// GetMsgID .
func GetMsgID(clientID, topic string) string {
	return crypto.MD5([]byte(fmt.Sprintf("%s-%s-%d-%d",
		clientID, topic, common.GetTimeNs(), rand.New(rand.NewSource(time.Now().UnixNano())).Int63())))
}

// CheckSub .
func CheckSub(account string, topics []string) (ans []TopicAns, err error) {
	if mqc.GetUserID(account) == 0 {
		req := new(CheckDevicePPMQSub.Req)
		req.Did = account
		req.Topics = topics
		req.Ipaddr = ""
		pkg, resp, err := g.MicrosConn.Invoke(context.Background(),
			CheckDevicePPMQSub.Module, CheckDevicePPMQSub.CmdID, req)
		if err != nil {
			err = fmt.Errorf("g.MicrosConn.Invoke(CheckDevicePPMQSub), %s", err)
			return nil, err
		}
		if pkg.Code != 0 {
			return nil, fmt.Errorf("rpc code: %d", pkg.Code)
		}
		for _, v := range resp.(*CheckDevicePPMQSub.Resp).Topics {
			s := TopicAns{
				Topic: v.Topic,
				Code:  v.Code,
			}
			ans = append(ans, s)
		}
	} else {
		uid := mqc.GetUserID(account)
		if uid == 0 {
			return nil, fmt.Errorf("Acccount: %s，format error", account)
		}
		req := new(CheckUserPPMQSub.Req)
		req.UserId = uid
		req.Topics = topics
		req.Ipaddr = ""
		pkg, resp, err := g.MicrosConn.Invoke(context.Background(),
			CheckUserPPMQSub.Module, CheckUserPPMQSub.CmdID, req)
		if err != nil {
			err = fmt.Errorf("g.MicrosConn.Invoke(CheckUserPPMQSub), %s", err)
			return nil, err
		}
		if pkg.Code != 0 {
			return nil, fmt.Errorf("rpc code: %d", pkg.Code)
		}
		for _, v := range resp.(*CheckUserPPMQSub.Resp).Topics {
			s := TopicAns{
				Topic: v.Topic,
				Code:  v.Code,
			}
			ans = append(ans, s)
		}
	}
	return
}

// CheckPub .
func CheckPub(account string, topic string) (code int32, err error) {
	//if account[:2] == g.PConf.DevicePrefix {
	if mqc.GetUserID(account) == 0 {
		req := new(CheckDevicePPMQPub.Req)
		req.Did = account
		req.Topic = topic
		req.Ipaddr = ""
		pkg, resp, err := g.MicrosConn.Invoke(context.Background(),
			CheckDevicePPMQPub.Module, CheckDevicePPMQPub.CmdID, req)
		if err != nil {
			err = fmt.Errorf("g.MicrosConn.Invoke(CheckDevicePPMQPub), %s", err)
			return 0, err
		}
		if pkg.Code != 0 {
			return 0, fmt.Errorf("rpc code: %d", pkg.Code)
		}
		code = resp.(*CheckDevicePPMQPub.Resp).Code

	} else {
		uid := mqc.GetUserID(account)
		if uid == 0 {
			return 0, fmt.Errorf("Acccount: %s，format error", account)
		}
		req := new(CheckUserPPMQPub.Req)
		req.UserId = uid
		req.Topic = topic
		req.Ipaddr = ""
		pkg, resp, err := g.MicrosConn.Invoke(context.Background(),
			CheckUserPPMQPub.Module, CheckUserPPMQPub.CmdID, req)
		if err != nil {
			err = fmt.Errorf("g.MicrosConn.Invoke(CheckUserPPMQPub), %s", err)
			return 0, err
		}
		if pkg.Code != 0 {
			return 0, fmt.Errorf("rpc code: %d", pkg.Code)
		}
		code = resp.(*CheckUserPPMQPub.Resp).Code
	}
	return
}

//GetOnlineByClientIDS .
func GetOnlineByClientIDS(cids []string) (oncids, offcids []OnlineInfo) {
	if len(cids) < 1 {
		logs.Logger.Debug("pm.GetConnectionByClientID(), cids is null.")
		return
	}
	rows, err := pm.GetConnectionByClientID(cids)
	if err != nil {
		logs.Logger.Debugf("pm.GetConnectionByClientID(), %d, %s", len(cids), err)
		return
	}
	for _, v := range rows {
		if v.ISOnline == mqc.ONLINE {
			oncids = append(oncids, OnlineInfo{ClientID: v.ClientID, ServerID: v.ServerID})
		} else {
			offcids = append(offcids, OnlineInfo{ClientID: v.ClientID, ServerID: v.ServerID})
		}
	}

	return
}
