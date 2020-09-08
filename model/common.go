package model

import (
	"fmt"

	"github.com/pprpc/util/crypto"
)

func getClientID(account, hwf string) string {
	return crypto.MD5([]byte(fmt.Sprintf("%s-%s", account, hwf)))
}

func getMsgStatus(cid, msgID string) (row *MsgStatus, err error) {
	row = new(MsgStatus)
	row.ClientID = cid
	row.MsgID = msgID
	_, err = row.GetAdv()
	return
}
