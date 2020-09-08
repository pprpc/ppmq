package model

import (
	"fmt"

	errc "xcthings.com/ppmq/common/errorcode"
)

// MsgLog table name msg_log struct define
type MsgLog struct {
	ID         uint32 `json:"_" xorm:"id"`
	MsgID      string `json:"msg_id" xorm:"msg_id"`
	Account    string `json:"account"`
	ClientID   string `json:"client_id" xorm:"client_id"`
	ServerID   string `json:"server_id" xorm:"server_id"`
	Status     int32  `json:"status"` // QOS1: 1,send; 2, pubAns; QOS2: 1, pub; 2, pubrec；3, pubrel；4, pubcomp
	CreateTime int64  `json:"create_time"`
}

// Add add record
func (obj *MsgLog) Add() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgLog")
		return
	}
	_, err = Orm.Insert(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Get get record
func (obj *MsgLog) Get() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgLog")
		return
	}
	if obj.MsgID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("Incomplete parameters: MsgID Can not be empty")
		return
	}
	var has bool
	has, err = Orm.Where("msg_id = ?", obj.MsgID).NoAutoCondition().Get(obj)
	if err != nil {
		code = errc.DBERROR
		return
	}
	if has != true {
		code = errc.NOTEXISTRECORD
		err = fmt.Errorf("The record of the search,MsgID: %vdoes not exist", obj.MsgID)
	}
	return
}

// Update update record
func (obj *MsgLog) Update() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgLog")
		return
	}
	if obj.MsgID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:MsgID Can not be empty")
		return
	}

	_, err = Orm.Where("msg_id = ?", obj.MsgID).NoAutoCondition().Update(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Delete delete record
func (obj *MsgLog) Delete() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgLog")
	}
	if obj.MsgID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:MsgID Can not be empty")
		return
	}
	/*
		up := new(MsgConnect)
		up.IsDel = 2
		_, err = Orm.Where("code = ?", obj.Code).NoAutoCondition().Update(up)
		if err != nil {
			code = errc.CodeDB
		}
	*/
	_, err = Orm.Where("msg_id = ?", obj.MsgID).NoAutoCondition().Delete(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Reset reset struct
func (obj *MsgLog) Reset() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgLog")
		return
	}
	*obj = MsgLog{}
	return
}

// GetAdv get record
func (obj *MsgLog) GetAdv() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgLog")
		return
	}
	if obj.MsgID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:MsgID Can not be empty")
		return
	}
	var has bool
	has, err = Orm.Get(obj)
	if err != nil {
		code = errc.DBERROR
		return
	}
	if has != true {
		code = errc.NOTEXISTRECORD
		err = fmt.Errorf("The record of the search,MsgID: %vdoes not exist", obj.MsgID)
	}
	return
}

// Set add or update record
func (obj *MsgLog) Set() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgLog")
		return
	}
	if obj.MsgID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:MsgID Can not be empty")
		return
	}

	_t := new(MsgLog)
	_t.MsgID = obj.MsgID
	code, err = _t.Get()
	if code == errc.NOTEXISTRECORD {
		code, err = obj.Add()
	} else if code == 0 {
		code, err = obj.Update()
	}
	return
}

// GetCount .
func (obj *MsgLog) GetCount() (count int64, code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgLog")
		return
	}
	if obj.MsgID == "" || obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:MsgID / ClientID Can not be empty")
		return
	}
	count, _ = Orm.Count(obj)
	return
}

// MsgLogInsertRows .
func MsgLogInsertRows(rows []*MsgLog) (code uint64, err error) {
	_, err = Orm.Insert(&rows)
	if err != nil {
		code = errc.DBERROR
	}
	return
}
