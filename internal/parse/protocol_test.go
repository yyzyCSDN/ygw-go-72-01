package parse

import (
	"testing"

	"powergw/internal/model"
)

func TestRegistryParsesIEC104(t *testing.T) {
	registry := NewRegistry()
	raw := BuildIEC104Frame(4001, 2205, 12)
	message, err := registry.Parse(model.ProtoIEC104, raw)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if message.Addr != 4001 {
		t.Fatalf("addr = %d", message.Addr)
	}
	if message.Seq != 12 {
		t.Fatalf("seq = %d", message.Seq)
	}
	if message.Kind != model.KindTelemetry {
		t.Fatalf("kind = %v", message.Kind)
	}
	value, ok := message.Field("value")
	if !ok || value.Value < 22.0 || value.Value > 22.1 {
		t.Fatalf("value = %v", value)
	}
}

func TestRegistryParsesModbus(t *testing.T) {
	registry := NewRegistry()
	raw := BuildModbusFrame(0x03, 42, 87)
	message, err := registry.Parse(model.ProtoModbus, raw)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(message.Fields) != 2 {
		t.Fatalf("fields = %d", len(message.Fields))
	}
	if message.Fields[0].Value != 42 || message.Fields[1].Value != 87 {
		t.Fatalf("register values wrong: %v", message.Fields)
	}
	if message.Addr != 1 {
		t.Fatalf("modbus addr = %d", message.Addr)
	}
}

func TestRegistryRejectsBadFrames(t *testing.T) {
	registry := NewRegistry()
	if _, err := registry.Parse(model.ProtoIEC104, []byte{0x68, 0x05}); err == nil {
		t.Fatal("bad iec104 frame accepted")
	}
	if _, err := registry.Parse(model.ProtoModbus, []byte{0x01, 0x03, 0x02, 0x00}); err == nil {
		t.Fatal("bad modbus frame accepted")
	}
	if _, err := registry.Parse(model.ProtoModbus, []byte{0x01, 0x10, 0x02, 0x00, 0x01}); err == nil {
		t.Fatal("unsupported function accepted")
	}
}

func TestRegistryProtocols(t *testing.T) {
	registry := NewRegistry()
	if !registry.Supports(model.ProtoIEC104) || !registry.Supports(model.ProtoModbus) {
		t.Fatal("supported protocols missing")
	}
	if registry.Supports(model.Protocol(8)) {
		t.Fatal("unsupported protocol reported supported")
	}
	if registry.Count() != 2 {
		t.Fatalf("protocol count = %d", registry.Count())
	}
	if len(registry.Protocols()) != 2 {
		t.Fatalf("protocols = %v", registry.Protocols())
	}
}

func TestParserAssignsSequence(t *testing.T) {
	parser := NewParser(NewRegistry())
	message, err := parser.Parse("sta-a", BuildModbusFrame(0x03, 10), model.ProtoModbus, "v1")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if message.Seq != 1 {
		t.Fatalf("first seq = %d", message.Seq)
	}
	second, err := parser.Parse("sta-a", BuildModbusFrame(0x03, 20), model.ProtoModbus, "v1")
	if err != nil {
		t.Fatalf("second parse failed: %v", err)
	}
	if second.Seq != 2 {
		t.Fatalf("second seq = %d", second.Seq)
	}
	if parser.LastSeq("sta-a") != 2 {
		t.Fatalf("last seq = %d", parser.LastSeq("sta-a"))
	}
}

func TestParserStateReset(t *testing.T) {
	parser := NewParser(NewRegistry())
	if _, err := parser.Parse("sta-a", BuildModbusFrame(0x03, 10), model.ProtoModbus, "v1"); err != nil {
		t.Fatal(err)
	}
	parser.Reset("sta-a")
	message, err := parser.Parse("sta-a", BuildModbusFrame(0x03, 11), model.ProtoModbus, "v1")
	if err != nil {
		t.Fatal(err)
	}
	if message.Seq != 1 {
		t.Fatalf("seq after reset = %d", message.Seq)
	}
	parser.ResetAll()
	if parser.LastSeq("sta-a") != 0 {
		t.Fatal("reset all failed")
	}
}
