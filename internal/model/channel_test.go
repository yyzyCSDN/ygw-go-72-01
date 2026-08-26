package model

import "testing"

func TestChannelStateString(t *testing.T) {
	wants := map[ChannelState]string{
		ChannelIdle:      "idle",
		ChannelConnected: "connected",
		ChannelSyncing:   "syncing",
		ChannelRunning:   "running",
		ChannelFault:     "fault",
		ChannelState(9):  "unknown",
	}
	for state, want := range wants {
		if got := state.String(); got != want {
			t.Fatalf("state %d string = %q", state, got)
		}
	}
}

func TestChannelStateCanAccept(t *testing.T) {
	if !ChannelRunning.CanAccept() || ChannelIdle.CanAccept() || ChannelFault.CanAccept() {
		t.Fatal("can accept check failed")
	}
	if !ChannelFault.Terminal() || ChannelRunning.Terminal() {
		t.Fatal("terminal check failed")
	}
}

func TestChannelRecordAdvanceState(t *testing.T) {
	record := NewChannelRecord("sta-a")
	if !record.AdvanceState(ChannelConnected) {
		t.Fatal("idle to connected failed")
	}
	if !record.AdvanceState(ChannelSyncing) {
		t.Fatal("connected to syncing failed")
	}
	if !record.AdvanceState(ChannelRunning) {
		t.Fatal("syncing to running failed")
	}
	if record.AdvanceState(ChannelConnected) {
		t.Fatal("invalid transition accepted")
	}
	if !record.AdvanceState(ChannelFault) {
		t.Fatal("running to fault failed")
	}
	if record.AdvanceState(ChannelRunning) {
		t.Fatal("fault transition accepted")
	}
	if record.AdvanceState(record.State) {
		t.Fatal("self transition accepted")
	}
}

func TestChannelSnapshotMerge(t *testing.T) {
	base := NewChannelSnapshot(1)
	base.Records["a"] = ChannelRecord{ID: "a", SyncSeq: 5}
	fresh := NewChannelSnapshot(2)
	fresh.Records["a"] = ChannelRecord{ID: "a", SyncSeq: 3}
	fresh.Records["b"] = ChannelRecord{ID: "b", SyncSeq: 1}
	merged := base.Merge(fresh)
	if record := merged.Records["a"]; record.SyncSeq != 5 {
		t.Fatalf("merge kept stale seq: %d", record.SyncSeq)
	}
	if _, ok := merged.Records["b"]; !ok {
		t.Fatal("merge dropped fresh record")
	}
}
