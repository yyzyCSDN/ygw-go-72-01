package channel

import (
	"testing"

	"powergw/internal/model"
)

func TestStateGuards(t *testing.T) {
	idle := model.NewChannelRecord("sta-a")
	if !CanConnect(idle) || CanSync(idle) || CanRun(idle) || CanFault(idle) {
		t.Fatal("idle guards wrong")
	}
	connected := model.ChannelRecord{ID: "sta-a", State: model.ChannelConnected}
	if CanConnect(connected) || !CanSync(connected) || CanRun(connected) || !CanFault(connected) {
		t.Fatal("connected guards wrong")
	}
	syncing := model.ChannelRecord{ID: "sta-a", State: model.ChannelSyncing}
	if !CanRun(syncing) || !CanFault(syncing) {
		t.Fatal("syncing guards wrong")
	}
	running := model.ChannelRecord{ID: "sta-a", State: model.ChannelRunning}
	if !CanFault(running) || CanRun(running) {
		t.Fatal("running guards wrong")
	}
	fault := model.ChannelRecord{ID: "sta-a", State: model.ChannelFault}
	if CanFault(fault) || CanRun(fault) {
		t.Fatal("fault guards wrong")
	}
}

func TestApplyTransition(t *testing.T) {
	record := model.NewChannelRecord("sta-a")
	next, err := ApplyTransition(record, model.ChannelConnected)
	if err != nil || next.State != model.ChannelConnected {
		t.Fatalf("transition failed: %v %v", next, err)
	}
	if _, err := ApplyTransition(record, model.ChannelRunning); err != ErrInvalidState {
		t.Fatalf("invalid transition err = %v", err)
	}
	if _, err := ApplyTransition(next, model.ChannelIdle); err != ErrInvalidState {
		t.Fatalf("backward transition err = %v", err)
	}
}
