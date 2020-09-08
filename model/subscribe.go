package model

import (
	"fmt"

	errc "xcthings.com/ppmq/common/errorcode"
)

// Subscribe table name subscribe struct define
type Subscribe struct {
	ID           uint32 `json:"_" xorm:"id"`
	Account      string `json:"account"`
	ClientID     string `json:"client_id" xorm:"client_id"`
	Topic        string `json:"topic"`
	Qos          uint32 `json:"qos"`
	Cluster      uint32 `json:"cluster"`
	ClusterSubid string `json:"cluster_subid"`
	LastTime     int64  `json:"last_time"`
	GlobalSync   int32  `json:"global_sync"`
}

// Add add record
func (obj *Subscribe) Add() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Subscribe")
		return
	}
	_, err = Orm.Insert(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Get get record
func (obj *Subscribe) Get() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Subscribe")
		return
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("Incomplete parameters: ClientID Can not be empty")
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
		err = fmt.Errorf("The record of the search,ClientID: %vdoes not exist", obj.ClientID)
	}
	return
}

// Update update record
func (obj *Subscribe) Update() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Subscribe")
		return
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:ClientID Can not be empty")
		return
	}

	_, err = Orm.Where("client_id = ?", obj.ClientID).NoAutoCondition().Update(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Delete delete record
func (obj *Subscribe) Delete() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Subscribe")
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:ClientID Can not be empty")
		return
	}
	_, err = Orm.Where("client_id = ?", obj.ClientID).NoAutoCondition().Delete(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// DeleteTopic delete record
func (obj *Subscribe) DeleteTopic() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Subscribe")
	}
	if obj.ClientID == "" || obj.Topic == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:ClientID, Topic Can not be empty")
		return
	}
	_, err = Orm.Where("client_id = ?", obj.ClientID).And("topic = ?", obj.Topic).NoAutoCondition().Delete(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Reset reset struct
func (obj *Subscribe) Reset() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Subscribe")
		return
	}
	*obj = Subscribe{}
	return
}

// GetAdv get record
func (obj *Subscribe) GetAdv() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Subscribe")
		return
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:ClientID Can not be empty")
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
		err = fmt.Errorf("The record of the search,ClientID: %vdoes not exist", obj.ClientID)
	}
	return
}

// Set add or update record
func (obj *Subscribe) Set() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Subscribe")
		return
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:ClientID Can not be empty")
		return
	}

	_t := new(Subscribe)
	_t.ClientID = obj.ClientID
	code, err = _t.Get()
	if code == errc.NOTEXISTRECORD {
		code, err = obj.Add()
	} else if code == 0 {
		code, err = obj.Update()
	}
	return
}

// GetByClientID 获得多条记录
func (obj *Subscribe) GetByClientID() (total int64, lists []*Subscribe, code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Subscribe")
		return
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:ClientID Can not be empty")
		return
	}

	lists = make([]*Subscribe, 0)
	total, _ = Orm.Count(obj)
	err = Orm.Find(&lists, obj)

	if err != nil {
		code = errc.DBERROR
	}
	return
}

// SubscribeInsertRows .
func SubscribeInsertRows(rows []*Subscribe) (code uint64, err error) {
	_, err = Orm.Insert(&rows)
	if err != nil {
		code = errc.DBERROR
	}
	return
}
