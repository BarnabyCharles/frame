package config

import (
	"errors"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var Client config_client.IConfigClient

func ClientNacos(NamespaceId string, host string, port int) error {
	//create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(host, uint64(port), constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(NamespaceId),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)
	// create config client
	var err error
	Client, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return errors.New("创建nacos连接客户端失败！" + err.Error())
	}
	return nil
}

func GetNacosConfig(SeverName string, group string) (string, error) {

	content, err := Client.GetConfig(vo.ConfigParam{
		DataId: SeverName,
		Group:  group,
	})
	if err != nil {
		return "", nil
	}
	err = Client.ListenConfig(vo.ConfigParam{
		DataId: SeverName,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			content = data
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
		},
	})

	return content, nil
}
