package model

import (
	"fmt"

	"xcthings.com/micro/svc"
	ppmq "github.com/pprpc/ppmq/model"
)

//InitEngine  db engine
func InitEngine(cfg svc.MSConfig) (err error) {
	for _, row := range cfg.Dbs {
		if row.Type != "mysql" && row.Type != "sqlite3" {
			err = fmt.Errorf("InitEngine, not support type: %s", row.Type)
			return
		}
		switch row.ConfName {
		case "ppmq":
			if row.Type == "mysql" {
				err = ppmq.MySQL(row)
				if err != nil {
					return
				}
			} else {
				err = ppmq.SQLite3(row)
				if err != nil {
					return
				}
			}
		default:
			err = fmt.Errorf("InitEngine, not support conf_name: %s", row.ConfName)
			return
		}
	}
	return
}
