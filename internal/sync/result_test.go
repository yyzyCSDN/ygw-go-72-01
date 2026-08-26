package timesync

import "testing"

func TestResultValid(t *testing.T) {
	valid := Result{ChannelID: "sta-a", Seq: 3, SyncedAt: 100}
	if !valid.Valid() {
		t.Fatal("valid result rejected")
	}
	invalid := Result{ChannelID: "sta-a", Seq: 0}
	if invalid.Valid() {
		t.Fatal("invalid result accepted")
	}
}

func TestResultList(t *testing.T) {
	list := NewResultList()
	list.Add(Result{ChannelID: "a", Seq: 1})
	list.Add(Result{ChannelID: "b", Seq: 2})
	if list.Len() != 2 {
		t.Fatalf("len = %d", list.Len())
	}
	items := list.Items()
	if len(items) != 2 || items[1].ChannelID != "b" {
		t.Fatalf("items = %v", items)
	}
	list.Clear()
	if list.Len() != 0 {
		t.Fatalf("len after clear = %d", list.Len())
	}
}

func TestClockSource(t *testing.T) {
	source := NewClockSource(func() int64 { return 42 })
	if source.Now() != 42 {
		t.Fatalf("now = %d", source.Now())
	}
	var nilSource *ClockSource
	if nilSource.Now() != 0 {
		t.Fatal("nil source now != 0")
	}
}
