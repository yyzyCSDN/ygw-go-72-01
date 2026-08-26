package dedup

import (
	"testing"

	"github.com/cespare/xxhash/v2"
)

func TestKeyDeterministic(t *testing.T) {
	first := Key("sta-a", 4001, 12)
	second := Key("sta-a", 4001, 12)
	if first != second {
		t.Fatalf("key not deterministic: %d %d", first, second)
	}
	if Key("sta-a", 4001, 12) == Key("sta-a", 4001, 13) {
		t.Fatal("seq not part of key")
	}
	if Key("sta-a", 4001, 12) == Key("sta-b", 4001, 12) {
		t.Fatal("channel not part of key")
	}
}

func TestWindowSize(t *testing.T) {
	window := NewWindow(128)
	if window.Size() != 128 || window.Half() != 64 {
		t.Fatalf("window wrong: %d %d", window.Size(), window.Half())
	}
	defaultWindow := NewWindow(0)
	if defaultWindow.Size() != DefaultWindow {
		t.Fatalf("default window = %d", defaultWindow.Size())
	}
	var nilWindow *Window
	if nilWindow.Size() != DefaultWindow {
		t.Fatal("nil window size wrong")
	}
}

func TestDeduperCheckOrMark(t *testing.T) {
	deduper := NewDeduper(NewWindow(4))
	key := Key("sta-a", 4001, 1)
	duplicate, err := deduper.CheckOrMark("sta-a", key)
	if err != nil || duplicate {
		t.Fatalf("first mark failed: %v %v", duplicate, err)
	}
	duplicate, err = deduper.CheckOrMark("sta-a", key)
	if err != nil || !duplicate {
		t.Fatalf("second check failed: %v %v", duplicate, err)
	}
	if deduper.SeenCount() != 1 || deduper.ChannelCount() != 1 {
		t.Fatalf("counts wrong: %d %d", deduper.SeenCount(), deduper.ChannelCount())
	}
	if deduper.Limit() != 4 {
		t.Fatalf("limit = %d", deduper.Limit())
	}
}

func TestDeduperEviction(t *testing.T) {
	deduper := NewDeduper(NewWindow(2))
	for index := 0; index < 5; index++ {
		_, _ = deduper.CheckOrMark("sta-a", uint64(index))
	}
	duplicate, _ := deduper.Check("sta-a", 0)
	if duplicate {
		t.Fatal("oldest key not evicted")
	}
	duplicate, _ = deduper.Check("sta-a", 4)
	if !duplicate {
		t.Fatal("newest key missing")
	}
	if deduper.SeenCount() != 2 {
		t.Fatalf("seen count = %d", deduper.SeenCount())
	}
}

func TestDeduperReset(t *testing.T) {
	deduper := NewDeduper(NewWindow(8))
	key := Key("sta-a", 4001, 1)
	_, _ = deduper.CheckOrMark("sta-a", key)
	if err := deduper.ResetChannel("sta-a"); err != nil {
		t.Fatal(err)
	}
	if deduper.SeenCount() != 0 {
		t.Fatalf("seen after reset channel = %d", deduper.SeenCount())
	}
	_, _ = deduper.CheckOrMark("sta-a", key)
	if err := deduper.ResetAll(); err != nil {
		t.Fatal(err)
	}
	if deduper.SeenCount() != 0 {
		t.Fatalf("seen after reset all = %d", deduper.SeenCount())
	}
}

func TestKeyStringStable(t *testing.T) {
	first := KeyString("sta-a", 4001, 12)
	second := KeyString("sta-a", 4001, 12)
	if first != second {
		t.Fatalf("key string not stable: %q %q", first, second)
	}
	digest := xxhash.New()
	addr := uint16(4001)
	seq := uint32(12)
	digest.Write([]byte("sta-a"))
	digest.Write([]byte{byte(addr), byte(addr >> 8)})
	digest.Write([]byte{byte(seq), byte(seq >> 8), byte(seq >> 16), byte(seq >> 24)})
	if Key("sta-a", addr, seq) != digest.Sum64() {
		t.Fatal("key does not match xxhash of parts")
	}
}
