package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/changqings/nacos-client-scq/nacosclient"
)

var (
	clientConfig = nacosclient.NacosConfig{
		ServerHost:  "192.168.1.15",
		ServerPort:  8848,
		UserName:    "dev",
		Passwd:      "devdevdev",
		NamespaceId: "1b87cfbc-25ee-4981-9494-34db9606b32f",
	}
)

func main() {

	// example usage:
	stopCh, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	client, err := clientConfig.NewClient()
	if err != nil {
		slog.Error("get client", "msg", err)
		return
	}

	// get config
	dataId, group := "test01", "dev"
	data, err := nacosclient.GetConfig(client, dataId, group)
	if err != nil {
		slog.Error("get config", "dataId", dataId, "group", group, "msg", err)
		return
	}
	slog.Info("get config", "dataId", dataId, "group", group, "data", data)

	// listen config, shoud block blow
	dataId, group = "test02", "dev"
	lisData, err := nacosclient.ListenConfig(client, dataId, group, stopCh)
	if err != nil {
		slog.Error("listen config", "dataId", dataId, "group", group, "msg", err)
	}

	//
	for {
		select {
		case d := <-lisData:
			fmt.Printf("config changed, ns: %s, group: %s, dataId: %s, new data: %s\n",
				d.Namespace, d.Group, d.DataId, d.Data)
		case <-stopCh.Done():
			os.Exit(0)
		}
	}

}
