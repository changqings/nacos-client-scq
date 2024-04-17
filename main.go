package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var (
	nacos_server_host  = "192.168.1.15"
	nacos_server_port  = 8848
	nacos_client_user  = "dev"
	nacos_client_pwd   = "devdevdev"
	nacos_client_ns_id = "1b87cfbc-25ee-4981-9494-34db9606b32f" // id of namespace dev

)

type DataConfig struct {
	Group     string
	Namespace string
	DataId    string
	Data      string
}

func main() {
	// create stopCh

	stopCh, cannel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cannel()

	// create server config
	ns := []constant.ServerConfig{
		*constant.NewServerConfig(nacos_server_host,
			uint64(nacos_server_port)),
	}

	// create client config

	nc := *constant.NewClientConfig(
		constant.WithUsername(nacos_client_user),
		constant.WithPassword(nacos_client_pwd),
		constant.WithNamespaceId(nacos_client_ns_id),
		constant.WithTimeoutMs(12000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogLevel("info"),
	)

	// create client

	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &nc,
			ServerConfigs: ns,
		},
	)

	if err != nil {
		panic(err)
	}

	// get config
	c1, err := client.GetConfig(vo.ConfigParam{
		DataId: "test01",
		Group:  "dev",
	})

	if err != nil {
		slog.Error("get config", "msg", err)
		os.Exit(1)
	}
	fmt.Printf("get config test01: %s\n", c1)

	// listen config
	dataCh := make(chan DataConfig, 1)

	lis_data_id := "test02"
	lis_group := "dev"

	errListen := client.ListenConfig(vo.ConfigParam{
		DataId: lis_data_id,
		Group:  lis_group,
		OnChange: func(namespace, group, dataId, data string) {
			dataCh <- DataConfig{
				Namespace: namespace,
				Group:     group,
				DataId:    dataId,
				Data:      data,
			}
		},
	})

	if errListen != nil {
		slog.Error("listen config", "msg", errListen)
		os.Exit(1)
	}

	for {
		select {
		case d := <-dataCh:
			fmt.Printf("config changed, ns: %s, group: %s, dataId: %s, new data: %s\n",
				d.Namespace, d.Group, d.DataId, d.Data)
		case <-stopCh.Done():
			os.Exit(0)
		}
	}
}
