package model

import (
	"fmt"

	errc "xcthings.com/ppmq/common/errorcode"
)

// Clientid table name client_id struct define
type Clientid struct {
	ID        uint32 `json:"_" xorm:"id"`
	Account   string `json:"account" xorm:"account"`
	ClientID  string `json:"client_id" xorm:"client_id"`
	HwFeature string `json:"hw_feature"`
}

// Add add record
func (obj *Clientid) Add() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Clientid")
		return
	}
	_, err = Orm.Insert(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Get get record
func (obj *Clientid) Get() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Clientid")
		return
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("Incomplete parameters: Clientid Can not be empty")
		return
	}
	var has bool
	has, err = Orm.Where("client_id = ?", obj.ClientID).NoAutoCondition().Get(obj)
	if err != nil {
		code = errc.DBERROR
		return
	}
	if has != true {
		code = errc.NOTEXISTRECORD
		err = fmt.Errorf("The record of the search,Clientid: %vdoes not exist", obj.ClientID)
	}
	return
}

// Update update record
func (obj *Clientid) Update() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Clientid")
		return
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:Clientid Can not be empty")
		return
	}

	_, err = Orm.Where("client_id = ?", obj.ClientID).NoAutoCondition().Update(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Delete delete record
func (obj *Clientid) Delete() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Clientid")
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:Clientid Can not be empty")
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
	_, err = Orm.Where("client_id = ?", obj.ClientID).NoAutoCondition().Delete(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Reset reset struct
func (obj *Clientid) Reset() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Clientid")
		return
	}
	*obj = Clientid{}
	return
}

// GetAdv get record
func (obj *Clientid) GetAdv() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Clientid")
		return
	}
	if obj.Account == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:Account Can not be empty")
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
		err = fmt.Errorf("The record of the search,Account: %vdoes not exist", obj.Account)
	}
	return
}

// Set add or update record
func (obj *Clientid) Set() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Clientid")
		return
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:Clientid Can not be empty")
		return
	}

	_t := new(Clientid)
	_t.ClientID = obj.ClientID
	code, err = _t.Get()
	if code == errc.NOTEXISTRECORD {
		code, err = obj.Add()
	} else if code == 0 {
		code, err = obj.Update()
	}
	return
}
