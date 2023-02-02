# nps 服务端

> 适配 `npc` `0.26.10` 后的客户端

# 使用方法

> 选择 `npc` `0.26.10` 后的版本 根据文档配置进行使用

- https://github.com/ehang-io/nps/releases
- https://github.com/yisier/nps/releases
- https://hub.docker.com/r/yisier1/npc


# protoc

```shell
protoc --proto_path=../ --go_out=./api --go-grpc_out=./api ../protos/account/*.proto
```

# 基于

- https://github.com/ehang-io/nps
- https://github.com/yisier/nps
