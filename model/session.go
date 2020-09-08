package model

import (
	"fmt"

	"xcthings.com/hjyz/common"
	errc "xcthings.com/ppmq/common/errorcode"
)

// SaveMsg .
func SaveMsg(raw *MsgRaw, info *MsgInfo) (code uint64, err error) {
	si := Orm.NewSession()
	defer si.Close()
	err = si.Begin()
	if err != nil {
		code = errc.DBERROR
		err = fmt.Errorf("si.Begin(), %s", err)
		return
	}

	_, err = si.Insert(raw)
	if err != nil {
		code = errc.DBERROR
		err = fmt.Errorf("si.Insert(MsgRaw), %s", err)
		si.Rollback()
		return
	}

	_, err = si.Insert(info)
	if err != nil {
		code = errc.DBERROR
		err = fmt.Errorf("si.Insert(MsgInfo), %s", err)
		si.Rollback()
		return
	}

	err = si.Commit()
	if err != nil {
		code = errc.DBERROR
		si.Rollback()
		err = fmt.Errorf("si.Commit(), %s", err)
	}
	return
}

// SetMsgStatus .
func SetMsgStatus(msgID, cid, account, sid string, status int32) (code uint64, err error) {
	row, le := getMsgStatus(cid, msgID)
	if le != nil {
		code = errc.DBERROR
		err = fmt.Errorf("getMsgStatus(%s,%s), %s", cid, msgID, le)
		return
	}

	si := Orm.NewSession()
	defer si.Close()
	err = si.Begin()
	if err != nil {
		code = errc.DBERROR
		err = fmt.Errorf("si.Begin(), %s", err)
		return
	}
	cusMS := common.GetTimeMs()

	log := new(MsgLog)
	log.MsgID = msgID
	log.Account = account
	log.ClientID = cid
	log.ServerID = sid
	log.Status = status
	log.CreateTime = cusMS

	_, err = si.Insert(log)
	if err != nil {
		code = errc.DBERROR
		err = fmt.Errorf("si.Insert(MsgLog), %s", err)
		si.Rollback()
		return
	}

	//_, err = si.Where("client_id = ?", cid).And("msg_id = ?", msgID).NoAutoCondition().Update(u)
	if (row.Qos == 1 && status == 2) || (row.Qos == 2 && status == 4) {
		_, err = si.Where("client_id = ?", cid).And("msg_id = ?", msgID).
			NoAutoCondition().Delete(new(MsgStatus))
		if err != nil {
			code = errc.DBERROR
			err = fmt.Errorf("si.Delete(MsgStatus), %s", err)
			si.Rollback()
			return
		}
	} else {
		u := new(MsgStatus)
		u.Status = status
		_, err = si.Where("client_id = ?", cid).And("msg_id = ?", msgID).
			NoAutoCondition().Update(u)
		if err != nil {
			code = errc.DBERROR
			err = fmt.Errorf("si.Update(MsgStatus), %s", err)
			si.Rollback()
			return
		}
	}

	err = si.Commit()
	if err != nil {
		code = errc.DBERROR
		si.Rollback()
		err = fmt.Errorf("si.Commit(), %s", err)
	}
	return

}

// ClearTempDB 清空临时数据
func ClearTempDB(timems int64) (code uint64, err error) {
	si := Orm.NewSession()
	defer si.Close()
	err = si.Begin()
	if err != nil {
		code = errc.DBERROR
		err = fmt.Errorf("si.Begin(), %s", err)
		return
	}

	// msg_info
	si.Where("create_time < ?", timems).NoAutoCondition().Delete(new(MsgInfo))
	// msg_log
	si.Where("create_time < ?", timems).NoAutoCondition().Delete(new(MsgLog))
	// msg_raw
	si.Where("create_time < ?", timems).NoAutoCondition().Delete(new(MsgRaw))
	// msg_status
	si.Where("create_time < ?", timems).NoAutoCondition().Delete(new(MsgStatus))

	err = si.Commit()
	if err != nil {
		code = errc.DBERROR
		si.Rollback()
		err = fmt.Errorf("si.Commit(), %s", err)
	}
	return
}
