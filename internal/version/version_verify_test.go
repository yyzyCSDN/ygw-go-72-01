package version

import (
	"testing"

	"powergw/internal/model"
)

type noopTables struct{}

func (n *noopTables) ApplyTable(id string) error {
	return nil
}

func (n *noopTables) StageTable(id string) error {
	return nil
}

func TestVersionStateUsesLatestProtocol(t *testing.T) {
	manager := NewManager()
	manager.RegisterChannel("sta-a", model.ProtoIEC104)
	tables := &noopTables{}
	if err := manager.Apply(model.NewProtocolVersion("v1", model.ProtoIEC104, "t1", 1), tables); err != nil {
		t.Fatal(err)
	}
	if err := manager.Apply(model.NewProtocolVersion("v2", model.ProtoModbus, "t2", 2), tables); err != nil {
		t.Fatal(err)
	}
	proto, _, err := manager.ActiveFor("sta-a")
	if err != nil {
		t.Fatal(err)
	}
	if proto != model.ProtoModbus {
		t.Fatalf("proto = %v, want modbus", proto)
	}
}
