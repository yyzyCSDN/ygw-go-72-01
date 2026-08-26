package channel

import (
	"testing"

	"powergw/internal/model"
)

func TestSessionStoreAcquireRelease(t *testing.T) {
	store := NewSessionStore()
	handle := store.Acquire("sta-a")
	if !handle.Open || !store.IsOpen("sta-a") || store.OpenCount() != 1 {
		t.Fatalf("acquire failed: %v %d", handle, store.OpenCount())
	}
	if !store.Release("sta-a") {
		t.Fatal("release failed")
	}
	if store.IsOpen("sta-a") || store.OpenCount() != 0 {
		t.Fatalf("release not effective: %d", store.OpenCount())
	}
	if store.Release("sta-a") {
		t.Fatal("double release accepted")
	}
	if store.ClosedCount() != 1 {
		t.Fatalf("closed count = %d", store.ClosedCount())
	}
}

func TestSessionStoreReleaseAll(t *testing.T) {
	store := NewSessionStore()
	store.Acquire("sta-a")
	store.Acquire("sta-b")
	if count := store.ReleaseAll(); count != 2 {
		t.Fatalf("release all = %d", count)
	}
	if store.OpenCount() != 0 {
		t.Fatalf("open count = %d", store.OpenCount())
	}
	if count := store.ReleaseAll(); count != 0 {
		t.Fatalf("second release all = %d", count)
	}
}

func TestSessionStoreSnapshot(t *testing.T) {
	store := NewSessionStore()
	store.Acquire("sta-a")
	store.Acquire("sta-b")
	store.Release("sta-b")
	snap := store.Snapshot()
	if !snap["sta-a"] || snap["sta-b"] {
		t.Fatalf("snapshot wrong: %v", snap)
	}
}

func TestMessageQueue(t *testing.T) {
	queue := NewMessageQueue()
	message := &model.Message{ChannelID: "sta-a", Seq: 1}
	if err := queue.Push(message); err != nil {
		t.Fatal(err)
	}
	got, ok := queue.Pop()
	if !ok || got.Seq != 1 {
		t.Fatalf("pop failed: %v %v", got, ok)
	}
	if queue.Len() != 0 {
		t.Fatalf("len = %d", queue.Len())
	}
	queue.Push(message)
	queue.Clear()
	if queue.Len() != 0 {
		t.Fatal("clear failed")
	}
}
