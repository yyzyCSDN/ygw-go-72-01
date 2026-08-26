package service

import (
	"testing"

	"powergw/internal/channel"
	"powergw/internal/map"
	"powergw/internal/model"
	"powergw/internal/parse"
)

func newTestGateway() *Gateway {
	gw := NewGateway(Options{
		ChannelIDs: []string{"sta-a", "sta-b"},
		Clock:      func() int64 { return 1700000000 },
	})
	registry := mapper.NewTableRegistry()
	for _, table := range mapper.DefaultStationTables() {
		_ = registry.Register(table)
	}
	_ = registry.LoadAll(gw.Mapper)
	gw.Version.RegisterChannel("sta-a", model.ProtoIEC104)
	gw.Version.RegisterChannel("sta-b", model.ProtoIEC104)
	_ = gw.ApplyVersion(model.NewProtocolVersion("v1", model.ProtoIEC104, "table-a", 1001))
	return gw
}

func TestGatewayIngestAndStatus(t *testing.T) {
	gw := newTestGateway()
	raw := parse.BuildIEC104Frame(4001, 2205, 1)
	if err := gw.Ingest("sta-a", raw); err != nil {
		t.Fatalf("ingest failed: %v", err)
	}
	if gw.Upload.QueueLen() != 1 {
		t.Fatalf("queue len = %d", gw.Upload.QueueLen())
	}
	status := gw.Status()
	if len(status.Channels) != 2 {
		t.Fatalf("channels = %d", len(status.Channels))
	}
	if status.ActiveTable != "table-a" {
		t.Fatalf("active table = %q", status.ActiveTable)
	}
}

func TestGatewayDuplicateDetection(t *testing.T) {
	gw := newTestGateway()
	raw := parse.BuildIEC104Frame(4001, 2205, 7)
	if err := gw.Ingest("sta-a", raw); err != nil {
		t.Fatalf("first ingest failed: %v", err)
	}
	if err := gw.Ingest("sta-a", raw); err != ErrDuplicate {
		t.Fatalf("second ingest err = %v", err)
	}
}

func TestGatewayUnmappedPoint(t *testing.T) {
	gw := newTestGateway()
	raw := parse.BuildIEC104Frame(9999, 1, 3)
	if err := gw.Ingest("sta-a", raw); err != ErrUnmappedPoint {
		t.Fatalf("unmapped err = %v", err)
	}
}

func TestGatewayUnknownChannel(t *testing.T) {
	gw := newTestGateway()
	if err := gw.Ingest("missing", parse.BuildIEC104Frame(4001, 1, 1)); err != ErrUnknownChannel {
		t.Fatalf("unknown channel err = %v", err)
	}
}

func TestGatewayFlushAll(t *testing.T) {
	gw := newTestGateway()
	if err := gw.Ingest("sta-a", parse.BuildIEC104Frame(4001, 2205, 1)); err != nil {
		t.Fatal(err)
	}
	if err := gw.Ingest("sta-b", parse.BuildIEC104Frame(4002, 105, 1)); err != nil {
		t.Fatal(err)
	}
	count, err := gw.FlushAll()
	if err != nil {
		t.Fatalf("flush failed: %v", err)
	}
	if count != 2 {
		t.Fatalf("flush count = %d", count)
	}
	if gw.Sink.Count() != 2 {
		t.Fatalf("sink count = %d", gw.Sink.Count())
	}
	if gw.Upload.QueueLen() != 0 {
		t.Fatalf("queue len after flush = %d", gw.Upload.QueueLen())
	}
}

func TestGatewaySyncRotateRecover(t *testing.T) {
	gw := newTestGateway()
	count, err := gw.SyncAll()
	if err != nil || count != 2 {
		t.Fatalf("sync all = %d %v", count, err)
	}
	rotated, err := gw.RotateAll()
	if err != nil || rotated != 2 {
		t.Fatalf("rotate all = %d %v", rotated, err)
	}
	if gw.Sessions.OpenCount() != 2 {
		t.Fatalf("open sessions = %d", gw.Sessions.OpenCount())
	}
	if err := gw.Recover(); err != nil {
		t.Fatalf("recover failed: %v", err)
	}
	closed := gw.CloseSessions()
	if closed != 2 {
		t.Fatalf("closed sessions = %d", closed)
	}
}

func TestGatewayProcessQueued(t *testing.T) {
	gw := newTestGateway()
	if err := gw.Ingest("sta-a", parse.BuildIEC104Frame(4001, 2200, 1)); err != nil {
		t.Fatal(err)
	}
	if err := gw.Ingest("sta-a", parse.BuildIEC104Frame(4002, 100, 2)); err != nil {
		t.Fatal(err)
	}
	if gw.Channels.QueueLen("sta-a") != 2 {
		t.Fatalf("channel queue len = %d", gw.Channels.QueueLen("sta-a"))
	}
	count, err := gw.ProcessQueued("sta-a")
	if err != nil {
		t.Fatalf("process queued failed: %v", err)
	}
	if count != 2 {
		t.Fatalf("processed count = %d", count)
	}
	if gw.Channels.QueueLen("sta-a") != 0 {
		t.Fatalf("queue len after processing = %d", gw.Channels.QueueLen("sta-a"))
	}
}

func TestGatewayEstablishAndFault(t *testing.T) {
	gw := newTestGateway()
	if err := gw.EstablishChannels(); err != nil {
		t.Fatalf("establish failed: %v", err)
	}
	record, ok := gw.Channels.Get("sta-a")
	if !ok {
		t.Fatal("channel missing")
	}
	if record.State != model.ChannelRunning {
		t.Fatalf("state = %v", record.State)
	}
	if err := gw.FaultChannel("sta-a"); err != nil {
		t.Fatalf("fault failed: %v", err)
	}
	record, _ = gw.Channels.Get("sta-a")
	if record.State != model.ChannelFault {
		t.Fatalf("state after fault = %v", record.State)
	}
	if err := gw.FaultChannel("missing"); err != channel.ErrUnknownChannel {
		t.Fatalf("missing channel fault err = %v", err)
	}
}

func TestGatewayActivateVersion(t *testing.T) {
	gw := newTestGateway()
	if err := gw.ActivateVersion(model.NewProtocolVersion("v2", model.ProtoModbus, "table-modbus", 2002)); err != nil {
		t.Fatalf("activate version failed: %v", err)
	}
	current, ok := gw.Version.Active()
	if !ok {
		t.Fatal("no active version")
	}
	if current.ID != "v2" {
		t.Fatalf("active version = %q", current.ID)
	}
}

func TestGatewayProcessChannel(t *testing.T) {
	gw := newTestGateway()
	frames := [][]byte{
		parse.BuildIEC104Frame(4001, 2200, 1),
		parse.BuildIEC104Frame(4002, 100, 2),
		parse.BuildIEC104Frame(4003, 500, 3),
	}
	count, err := ProcessChannel(gw, "sta-a", frames)
	if err != nil || count != 3 {
		t.Fatalf("process channel = %d %v", count, err)
	}
}

func TestGatewayBuildTelemetry(t *testing.T) {
	gw := newTestGateway()
	message, err := gw.Parser.Parse("sta-a", parse.BuildIEC104Frame(4001, 2205, 1), model.ProtoIEC104, "v1")
	if err != nil {
		t.Fatal(err)
	}
	point, ok := MapMessage(gw, message)
	if !ok {
		t.Fatal("map failed")
	}
	tel := BuildTelemetry(gw, message, point)
	if tel.PointName != "线路电压" {
		t.Fatalf("point name = %q", tel.PointName)
	}
	if tel.VersionID != "v1" {
		t.Fatalf("version id = %q", tel.VersionID)
	}
	if tel.Value < 22.04 || tel.Value > 22.06 {
		t.Fatalf("value = %v", tel.Value)
	}
}

func TestGatewayFeedDemo(t *testing.T) {
	gw := newTestGateway()
	count, err := FeedDemo(gw, 3)
	if err != nil {
		t.Fatalf("feed demo failed: %v", err)
	}
	if count != 6 {
		t.Fatalf("feed count = %d", count)
	}
	if _, err := gw.FlushAll(); err != nil {
		t.Fatalf("flush after demo failed: %v", err)
	}
	if gw.Sink.Count() != 6 {
		t.Fatalf("sink count = %d", gw.Sink.Count())
	}
}

func TestMasterSink(t *testing.T) {
	sink := NewMasterSink(3)
	for index := 0; index < 5; index++ {
		if err := sink.Send(model.NewTelemetry("点", uint16(index), float64(index), uint32(index), int64(index))); err != nil {
			t.Fatal(err)
		}
	}
	if sink.Count() != 3 {
		t.Fatalf("count = %d", sink.Count())
	}
	last, ok := sink.Last()
	if !ok || last.RawAddr != 4 {
		t.Fatalf("last = %v %v", last, ok)
	}
	if len(sink.Items()) != 3 {
		t.Fatalf("items = %d", len(sink.Items()))
	}
}

func TestGatewayRunnerCycle(t *testing.T) {
	gw := newTestGateway()
	runner := NewRunner(gw, 0)
	if err := runner.Cycle(); err != nil {
		t.Fatalf("cycle failed: %v", err)
	}
	if gw.Sessions.OpenCount() != 2 {
		t.Fatalf("open sessions after cycle = %d", gw.Sessions.OpenCount())
	}
}
