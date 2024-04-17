package nacosclient

import (
	"context"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosConfig struct {
	ServerHost  string
	ServerPort  uint64
	UserName    string
	Passwd      string
	NamespaceId string
	TimeOutMs   uint64
	Loglevel    string
}

type DataConfig struct {
	Group     string
	Namespace string
	DataId    string
	Data      string
}

type NacosClient config_client.IConfigClient

func (n *NacosConfig) NewClient() (NacosClient, error) {

	// create server config
	serverC := []constant.ServerConfig{
		*constant.NewServerConfig(n.ServerHost,
			n.ServerPort,
		),
	}

	// create client config
	clientC := constant.NewClientConfig(
		constant.WithUsername(n.UserName),
		constant.WithPassword(n.Passwd),
		constant.WithNamespaceId(n.NamespaceId),
		constant.WithTimeoutMs(n.TimeOutMs),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogLevel(n.Loglevel),
	)

	// set default

	if clientC.LogLevel == "" {
		clientC.LogLevel = "info"
	}
	if clientC.TimeoutMs == 0 {
		clientC.TimeoutMs = 5000
	}

	// create client
	return clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  clientC,
			ServerConfigs: serverC,
		},
	)

}

func GetConfig(nc NacosClient, dataId string, group string) (string, error) {

	// get config
	data, err := nc.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	return data, err
}

func ListenConfig(nc NacosClient, dataId string, group string, stopCh context.Context) (chan DataConfig, error) {

	// listen config
	dataCh := make(chan DataConfig, 1)

	err := nc.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			dataCh <- DataConfig{
				Namespace: namespace,
				Group:     group,
				DataId:    dataId,
				Data:      data,
			}
		},
	})

	return dataCh, err

}
