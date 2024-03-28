package config

import (
	"errors"

	"github.com/spf13/viper"
)

func InitViper(fileName, filePath, fileType string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(fileName)
	v.SetConfigFile(filePath)

	v.SetConfigType(fileType)
	err := v.ReadInConfig()
	if err != nil {
		return nil, errors.New("获取配置信息失败！" + err.Error())
	}

	return v, nil
}
