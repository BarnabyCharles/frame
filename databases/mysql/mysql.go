package mysql

import (
	"errors"
	"fmt"
	"github.com/BarnabyCharles/frame/config"
	"github.com/ghodss/yaml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMysql(serverName, group string) error {
	nacos, err := config.GetNacosConfig(serverName, group)
	if err != nil {
		return err
	}
	var mysqConfig config.AppConfig
	err = yaml.Unmarshal([]byte(nacos), &mysqConfig)

	if err != nil {
		return errors.New("将yaml文件转换为结构体格式失败！" + err.Error())
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		mysqConfig.Mysql.Username,
		mysqConfig.Mysql.Password,
		mysqConfig.Mysql.Host,
		mysqConfig.Mysql.Port,
		mysqConfig.Mysql.Database,
	)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		return errors.New("连接数据库失败！" + err.Error())
	}

	sql, err := db.DB()
	if err != nil {
		return err
	}
	defer sql.Close()
	DB = db

	return err
}
