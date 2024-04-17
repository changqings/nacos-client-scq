## nacos 使用记录

## 测试使用，使用单机模式，standalone

搭建使用 docker-compose, `https://github.com/nacos-group/nacos-docker`
修改需要认证登录，修改`env/nacos-standalone-mysql.env`
构建失败，拉不下来 sql 脚本，手动处理，注意构建时的使用相对路径

```
NACOS_AUTH_ENABLE=true
# 认证开关，开启后，ui默认nacos/nacos,权限为admin,可以添加新用户后
# 修改数据库表`roles`,`users`这两个表中的nacos为新用户来禁用nacos用户

## 以下identity_key_value
# 官方说明，https://nacos.io/zh-cn/blog/announcement-token-secret-key.html
# 文档这块做的不是很好，这还是我在google上查到的，官方文档的查询是百度，而不是内部查询，LJ
NACOS_AUTH_IDENTITY_KEY=username
NACOS_AUTH_IDENTITY_VALUE=shenchangqing
NACOS_AUTH_TOKEN=SmtjeENZN21nY1ZxdTRYWUN0SktrTEpBdXhuSHZDRUwK # 此值为原32位字符串base64加密后的值，不然启动报错
```

## 引入

```
go get github.com/changqings/nacos-client-scq/nacosclient
```

### example

```go
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
```
