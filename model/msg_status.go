package model

import (
	"fmt"

	errc "xcthings.com/ppmq/common/errorcode"
)

// MsgStatus table name msg_status struct define
type MsgStatus struct {
	ID         uint32 `json:"_" xorm:"id"`
	MsgID      string `json:"msg_id" xorm:"msg_id"`
	SrcMsgid   string `json:"src_msgid" xorm:"src_msgid"`
	Account    string `json:"account"`
	ClientID   string `json:"client_id" xorm:"client_id"`
	ServerID   string `json:"server_id" xorm:"server_id"`
	Dup        int32  `json:"dup"`
	Qos        int32  `json:"qos"`
	Status     int32  `json:"status"`      // QOS1: 1,send; 2, pubAns; QOS2: 1, pub; 2, pubrec；3, pubrel；4, pubcomp
	CreateTime int64  `json:"create_time"` //
}

// Add add record
func (obj *MsgStatus) Add() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgStatus")
		return
	}
	_, err = Orm.Insert(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Get get record
func (obj *MsgStatus) Get() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgStatus")
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
func (obj *MsgStatus) Update() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgStatus")
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
func (obj *MsgStatus) Delete() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgStatus")
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
func (obj *MsgStatus) Reset() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgStatus")
		return
	}
	*obj = MsgStatus{}
	return
}

// GetAdv get record
func (obj *MsgStatus) GetAdv() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgStatus")
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
func (obj *MsgStatus) Set() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgStatus")
		return
	}
	if obj.MsgID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:MsgID Can not be empty")
		return
	}

	_t := new(MsgStatus)
	_t.MsgID = obj.MsgID
	code, err = _t.Get()
	if code == errc.NOTEXISTRECORD {
		code, err = obj.Add()
	} else if code == 0 {
		code, err = obj.Update()
	}
	return
}

// UpdateStatusByClientID update record
func (obj *MsgStatus) UpdateStatusByClientID(s int32) (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgStatus")
		return
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:ClientID Can not be empty")
		return
	}
	_t := new(MsgStatus)
	_t.ClientID = obj.ClientID
	_t.Status = s
	_, err = Orm.Where("client_id = ?", _t.ClientID).And("status = ?", 1).NoAutoCondition().Update(_t)
	if err != nil {
		code = errc.DBERROR
	}
	return
}
