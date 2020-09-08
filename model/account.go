package model

import (
	"fmt"

	errc "xcthings.com/ppmq/common/errorcode"
)

// Account table name account struct define
type Account struct {
	ID       uint32 `json:"_" xorm:"id"`
	UserID   int64  `json:"user_id" xorm:"user_id"`
	Account  string `json:"account" xorm:"account"`
	Password string `json:"password"`
}

// : UserID,Account; SQL: alter table account add constraint uid_account unique(user_id,account);

// Add add record
func (obj *Account) Add() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Account")
		return
	}
	_, err = Orm.Insert(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Get get record
func (obj *Account) Get() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Account")
		return
	}
	if obj.Account == "" {
		code = errc.ParameterError
		err = fmt.Errorf("Incomplete parameters: Account Can not be empty")
		return
	}
	var has bool
	has, err = Orm.Where("account = ?", obj.Account).NoAutoCondition().Get(obj)
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

// Update update record
func (obj *Account) Update() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Account")
		return
	}
	if obj.Account == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:Account Can not be empty")
		return
	}

	_, err = Orm.Where("account = ?", obj.Account).NoAutoCondition().Update(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Delete delete record
func (obj *Account) Delete() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Account")
	}
	if obj.Account == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:Account Can not be empty")
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
	_, err = Orm.Where("account = ?", obj.Account).NoAutoCondition().Delete(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Reset reset struct
func (obj *Account) Reset() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Account")
		return
	}
	*obj = Account{}
	return
}

// GetAdv get record
func (obj *Account) GetAdv() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Account")
		return
	}
	if obj.Account == "" && obj.UserID == 0 {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:Account/UserID Cannot be empty at the same time")
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
func (obj *Account) Set() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Account")
		return
	}
	if obj.Account == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:Account Can not be empty")
		return
	}

	_t := new(Account)
	_t.Account = obj.Account
	code, err = _t.Get()
	if code == errc.NOTEXISTRECORD {
		code, err = obj.Add()
	} else if code == 0 {
		code, err = obj.Update()
	}
	return
}
