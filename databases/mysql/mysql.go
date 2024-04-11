package mysql

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ghodss/yaml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/BarnabyCharles/frame/config"
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

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return errors.New("连接数据库失败！" + err.Error())
	}

	return err
}
