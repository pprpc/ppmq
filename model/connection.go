package model

import (
	"fmt"

	errc "github.com/pprpc/ppmq/common/errorcode"
)

// Connection table name connection struct define
type Connection struct {
	ID           uint32 `json:"_" xorm:"id"`
	UserID       int64  `json:"user_id" xorm:"user_id"`
	Account      string `json:"account"`
	ServerID     string `json:"server_id" xorm:"server_id"`
	ClearSession int32  `json:"clear_session" xorm:"clear_session"` //
	ClientID     string `json:"client_id" xorm:"client_id"`
	// will message
	WillFlag   uint32 `json:"will_flag"`
	WillRetain uint32 `json:"will_retain"`
	WillQOS    uint32 `json:"will_qos" xorm:"will_qos"`
	WillTopic  string `json:"will_topic"`
	WillBody   string `json:"will_body"`
	// connection info
	ConnType int32  `json:"conn_type"` //  	1,tcp ppmq;2,udp; 3, mqtt
	ConnInfo string `json:"conn_info"`
	ISOnline int32  `json:"is_online" xorm:"is_online"` // 1, online 2, offline
	ISSleep  int32  `json:"is_sleep" xorm:"is_sleep"`   // 0: normal; 1 sleep
	// history message
	HistorymsgType  int32 `json:"historymsg_type"`
	HistorymsgCount int32 `json:"historymsg_count"`
	LastTime        int64 `json:"last_time"`
	GlobalSync      int32 `json:"global_sync"`
}

// Add add record
func (obj *Connection) Add() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Connection")
		return
	}
	_, err = Orm.Insert(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Get get record
func (obj *Connection) Get() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Connection")
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
func (obj *Connection) Update() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Connection")
		return
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:ClientID Can not be empty")
		return
	}

	_, err = Orm.Where("client_id = ?", obj.ClientID).NoAutoCondition().MustCols("clear_session").Update(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Delete delete record
func (obj *Connection) Delete() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Connection")
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:ClientID Can not be empty")
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
func (obj *Connection) Reset() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Connection")
		return
	}
	*obj = Connection{}
	return
}

// GetAdv get record
func (obj *Connection) GetAdv() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Connection")
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
func (obj *Connection) Set() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Connection")
		return
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:ClientID Can not be empty")
		return
	}

	_t := new(Connection)
	_t.ClientID = obj.ClientID
	code, err = _t.Get()
	if code == errc.NOTEXISTRECORD {
		code, err = obj.Add()
	} else if code == 0 {
		code, err = obj.Update()
	}
	return
}

// SetOffline 设置设备离线
func (obj *Connection) SetOffline() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Connection")
		return
	}
	if obj.ClientID == "" || obj.LastTime == 0 || obj.ISOnline == 0 {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:ClientID, LastTime, ISOnline")
		return
	}
	_, err = Orm.Where("client_id = ?", obj.ClientID).And("server_id = ?", obj.ServerID).NoAutoCondition().Cols("is_online", "last_time").Update(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// SetOfflineByServerID .
func (obj *Connection) SetOfflineByServerID() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Connection")
		return
	}
	if obj.ServerID == "" || obj.LastTime == 0 || obj.ISOnline == 0 {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:ServerID, LastTime, ISOnline")
		return
	}
	_, err = Orm.Where("server_id = ?", obj.ServerID).NoAutoCondition().Cols("is_online", "last_time").Update(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// QueryIn .
func (obj *Connection) QueryIn(accounts []string) (rows []*Connection, code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Connection")
		return
	}
	rows = make([]*Connection, 0)
	// Orm.Desc("id").Limit(int(pageSize), int((page-1)*pageSize)).Find(&dms)
	// Orm.Where("status = ?", tab.Status).Desc("id").Limit(int(pageSize), int((page-1)*pageSize)).Find(&dms)
	// total, err = Orm.Count(tab)
	err = Orm.In("account", accounts).Find(&rows)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// UpdateSleep .
func (obj *Connection) UpdateSleep() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Connection")
		return
	}
	if obj.ClientID == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:ClientID Can not be empty")
		return
	}

	_, err = Orm.Where("client_id = ?", obj.ClientID).Cols("is_sleep").NoAutoCondition().Update(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// GetConnectionByClientID .
func GetConnectionByClientID(cids []string) (rows []*Connection, err error) {
	rows = make([]*Connection, 0)
	err = Orm.In("client_id", cids).Find(&rows)
	return
}

// SetOnlineStatusByServerID set device online stataus by serverid
func SetOnlineStatusByServerID(serverID string, isOnline int32) (err error) {
	obj := new(Connection)
	obj.ISOnline = isOnline

	_, err = Orm.Where("server_id = ?", serverID).Cols("is_online").NoAutoCondition().Update(obj)
	return
}

// GetRowsByAccount .
func GetRowsByAccount(account string) (rows []*Connection, code uint64, err error) {
	if account == "" {
		code = errc.ParameterError
		err = fmt.Errorf("Incomplete parameters: account Can not be empty")
		return
	}
	rows = make([]*Connection, 0)
	err = Orm.Where("account = ?", account).NoAutoCondition().Find(&rows)

	if err != nil {
		code = errc.DBERROR
	}
	return
}

// DeleteConnection .
func DeleteConnection(cid string) (code uint64, err error) {
	if cid == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:cid Can not be empty")
		return
	}
	_, err = Orm.Where("client_id = ?", cid).NoAutoCondition().Delete(new(Connection))
	if err != nil {
		code = errc.DBERROR
	}
	return
}
