package main

import (
	"time"

	"powergw/internal/map"
	"powergw/internal/model"
	"powergw/internal/service"
)

func BuildGateway() *service.Gateway {
	gw := service.NewGateway(service.Options{
		ChannelIDs: []string{"sta-a", "sta-b"},
		Clock:      time.Now().Unix,
	})
	registry := mapper.NewTableRegistry()
	for _, table := range mapper.DefaultStationTables() {
		_ = registry.Register(table)
	}
	for _, table := range registry.All() {
		_ = gw.RegisterTable(table)
	}
	_ = gw.ActivateTable("table-a")
	gw.Version.RegisterChannel("sta-a", model.ProtoIEC104)
	gw.Version.RegisterChannel("sta-b", model.ProtoIEC104)
	_ = gw.ActivateVersion(model.NewProtocolVersion("v1", model.ProtoIEC104, "table-a", 1001))
	_ = gw.EstablishChannels()
	_ = service.Restore(gw, service.SnapshotForRecovery(gw))
	return gw
}
