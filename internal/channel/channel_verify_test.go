package channel

import "testing"

func TestChannelHandleClosed(t *testing.T) {
	manager := NewManager([]string{"sta-a"})
	store := NewSessionStore()
	first, err := manager.RotateSession(store, "sta-a")
	if err != nil {
		t.Fatal(err)
	}
	if !first.Open {
		t.Fatal("first handle not open")
	}
	second, err := manager.RotateSession(store, "sta-a")
	if err != nil {
		t.Fatal(err)
	}
	if !second.Open {
		t.Fatal("second handle not open")
	}
	if first.Open {
		t.Fatal("old handle still open after rotation")
	}
	if store.OpenCount() != 1 {
		t.Fatalf("open count = %d, want 1", store.OpenCount())
	}
}
