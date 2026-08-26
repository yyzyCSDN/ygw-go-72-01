基于 Go 实现的电力规约转换与遥测前置网关系统项目，一款电力数据网关服务，完成规约解析转换、点表映射、遥测上送、通道管理、对时与版本切换。

PowerProtocolGateway 部署在变电站前置机侧，接入 IEC 60870-5-104 与 Modbus 通道，
把从站设备上报的原始遥测/遥信报文解析为统一消息，按点表把原始地址映射到主站点名，
再上送主站；通道按周期对时，规约版本升级后点表同步更新，重复报文被去重拦截。

## 构建

依赖已 vendor 离线，构建时无需联网：

    go build -mod=vendor ./...

## 运行

    go run ./cmd/gateway -addr 127.0.0.1:8090 -dir ./data

启动后访问 http://127.0.0.1:8090/console 打开前置网关控制台。

## 验证

    go test -mod=vendor ./...
    go vet -mod=vendor ./...
