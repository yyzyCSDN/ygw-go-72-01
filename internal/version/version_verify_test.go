package version

import (
	"testing"

	"powergw/internal/model"
)

type failingResetter struct {
	calls int
}

func (r *failingResetter) ResetAll() error {
	r.calls++
	if r.calls == 2 {
		return ErrWriteback
	}
	return nil
}

type applyTables struct{}

func (a *applyTables) ApplyTable(id string) error {
	return nil
}

func TestVersionSwitchWritebackErrorNotSwallowed(t *testing.T) {
	manager := NewManager()
	resetter := &failingResetter{}
	if err := manager.Switch(model.NewProtocolVersion("v1", model.ProtoIEC104, "t1", 1), &applyTables{}, resetter); err != nil {
		t.Fatal(err)
	}
	if err := manager.Switch(model.NewProtocolVersion("v2", model.ProtoModbus, "t2", 2), &applyTables{}, resetter); err == nil {
		t.Fatal("version switch writeback error was swallowed")
	}
}
