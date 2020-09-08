package model

import (
	"fmt"

	errc "github.com/pprpc/ppmq/common/errorcode"
)

// MsgRaw table name msg_raw struct define
type MsgRaw struct {
	ID         uint32 `json:"_" xorm:"id"`
	MsgID      string `json:"msg_id" xorm:"msg_id"`
	MsgPayload []byte `json:"msg_payload" xorm:"msg_payload"`
}

// Add add record
func (obj *MsgRaw) Add() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgRaw")
		return
	}
	_, err = Orm.Insert(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Get get record
func (obj *MsgRaw) Get() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgRaw")
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
func (obj *MsgRaw) Update() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgRaw")
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
func (obj *MsgRaw) Delete() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgRaw")
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
func (obj *MsgRaw) Reset() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgRaw")
		return
	}
	*obj = MsgRaw{}
	return
}

// GetAdv get record
func (obj *MsgRaw) GetAdv() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgRaw")
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
func (obj *MsgRaw) Set() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： MsgRaw")
		return
	}
	if obj.MsgID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:MsgID Can not be empty")
		return
	}

	_t := new(MsgRaw)
	_t.MsgID = obj.MsgID
	code, err = _t.Get()
	if code == errc.NOTEXISTRECORD {
		code, err = obj.Add()
	} else if code == 0 {
		code, err = obj.Update()
	}
	return
}
