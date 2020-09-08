package model

import (
	"fmt"

	errc "github.com/pprpc/ppmq/common/errorcode"
)

// MsgInfo table name msg_info struct define
type MsgInfo struct {
	ID         uint32 `json:"_" xorm:"id"`
	MsgID      string `json:"msg_id" xorm:"msg_id"`
	SrcMsgid   string `json:"src_msgid" xorm:"src_msgid"`
	Account    string `json:"account"`
	ClientID   string `json:"client_id" xorm:"client_id"`
	Dup        int32  `json:"dup"`
	Retain     int32  `json:"retain"`
	Qos        int32  `json:"qos"`
	Topic      string `json:"topic"`
	Format     int32  `json:"format"`      //   4, PB-BIN; 5, PB-JSON;
	Cmdid      uint64 `json:"cmdid"`       //
	CmdType    int32  `json:"cmd_type"`    //   0 request：  1 response
	MsgTimems  int64  `json:"msg_timems"`  //
	CreateTime int64  `json:"create_time"` //
}

// Add add record
func (obj *MsgInfo) Add() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgInfo")
		return
	}
	_, err = Orm.Insert(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Get get record
func (obj *MsgInfo) Get() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgInfo")
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
func (obj *MsgInfo) Update() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgInfo")
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
func (obj *MsgInfo) Delete() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgInfo")
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
func (obj *MsgInfo) Reset() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgInfo")
		return
	}
	*obj = MsgInfo{}
	return
}

// GetAdv get record
func (obj *MsgInfo) GetAdv() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgInfo")
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
func (obj *MsgInfo) Set() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgInfo")
		return
	}
	if obj.MsgID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:MsgID Can not be empty")
		return
	}

	_t := new(MsgInfo)
	_t.MsgID = obj.MsgID
	code, err = _t.Get()
	if code == errc.NOTEXISTRECORD {
		code, err = obj.Add()
	} else if code == 0 {
		code, err = obj.Update()
	}
	return
}
