package app

import (
	"github.com/BarnabyCharles/frame/config"
	"github.com/BarnabyCharles/frame/databases/mysql"
	"github.com/BarnabyCharles/frame/es"
)

func Init(ServerName, group, NamespaceId, host string, port int, str ...string) error {
	var err error

	err = config.ClientNacos(NamespaceId, host, port)
	if err != nil {
		return err
	}

	for _, val := range str {
		switch val {
		case "mysql":
			err = mysql.InitMysql(ServerName, group)
		case "es":
			err = es.ConnectionEs(ServerName, group)
		}
	}
	return err
}
