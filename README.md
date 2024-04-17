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
# 文档这块做的不是很好，这还是我在google上查到的，官方文档的查询是百度，而不是内部查询
NACOS_AUTH_IDENTITY_KEY=username
NACOS_AUTH_IDENTITY_VALUE=shenchangqing
NACOS_AUTH_TOKEN=SmtjeENZN21nY1ZxdTRYWUN0SktrTEpBdXhuSHZDRUwK # 此值为原32位字符串base64加密后的值，不然启动报错
```
