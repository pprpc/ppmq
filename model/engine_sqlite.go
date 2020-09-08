package model

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"xcthings.com/hjyz/common"
	"xcthings.com/micro/svc"
	//_ "github.com/mattn/go-sqlite3"
)

// SQLite3  .
func SQLite3(cfg svc.ValueDbconf) (err error) {
	if cfg.Type != "sqlite3" {
		err = fmt.Errorf("Engine SQLite3, not support: %s", cfg.Type)
		return
	}
	if cfg.Debug == true {
		Orm.ShowSQL(cfg.Debug)
	}
	Orm, err = xorm.NewEngine(cfg.Type, cfg.Name)
	if err != nil {
		err = fmt.Errorf("xorm.NewEngine(%s, %s), %s", cfg.Type, cfg.Name, err)
		return
	}
	if common.PathIsExist(cfg.Name) == false {
		err = Orm.Sync2(new(Account), new(Clientid), new(Connection),
			new(MsgInfo), new(MsgLog), new(MsgRaw), new(MsgStatus), new(Subscribe))
		//err = Orm.Sync2()
		if err != nil {
			err = fmt.Errorf("Orm.Sync2(), %s", err)
			return
		}
	}
	return
}
