package version

import (
	"testing"

	"powergw/internal/model"
)

func TestHistoryAppend(t *testing.T) {
	history := NewHistory(3)
	for index := 0; index < 5; index++ {
		history.Append(model.NewProtocolVersion("v", model.ProtoIEC104, "t", uint64(index)))
	}
	if history.Len() != 3 {
		t.Fatalf("len = %d", history.Len())
	}
	latest, ok := history.Latest()
	if !ok || latest.Checksum != 4 {
		t.Fatalf("latest = %v %v", latest, ok)
	}
	if history.Limit() != 3 {
		t.Fatalf("limit = %d", history.Limit())
	}
	if len(history.List()) != 3 {
		t.Fatalf("list len = %d", len(history.List()))
	}
}

func TestHistoryEmpty(t *testing.T) {
	history := NewHistory(0)
	if history.Len() != 0 {
		t.Fatalf("len = %d", history.Len())
	}
	if _, ok := history.Latest(); ok {
		t.Fatal("empty history latest")
	}
	var nilHistory *History
	if nilHistory.Len() != 0 || nilHistory.Limit() != 0 {
		t.Fatal("nil history failed")
	}
	if nilHistory.Append(model.ProtocolVersion{}) {
		t.Fatal("nil history append accepted")
	}
}

func TestVersionStateString(t *testing.T) {
	if model.VersionDraft.String() != "draft" || model.VersionActive.String() != "active" || model.VersionSuperseded.String() != "superseded" {
		t.Fatal("version state strings wrong")
	}
}

func TestProtocolVersionFactory(t *testing.T) {
	ver := model.NewProtocolVersion("v1", model.ProtoIEC104, "t1", 99)
	if ver.State != model.VersionDraft {
		t.Fatalf("state = %v", ver.State)
	}
	if ver.Proto != model.ProtoIEC104 || ver.TableID != "t1" || ver.Checksum != 99 {
		t.Fatalf("fields wrong: %+v", ver)
	}
}
