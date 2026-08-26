package mapper

import (
	"testing"

	"powergw/internal/model"
)

func TestPendingAndLatestUploaded(t *testing.T) {
	table := BuildTable("t1", "v1", []model.Point{
		{RawAddr: 1, Name: "电压", Uploaded: true},
		{RawAddr: 2, Name: "电流"},
	})
	snap := model.NewTableSnapshot(table, 1)
	pending := PendingAddrs(snap)
	if len(pending) != 1 || pending[0] != 2 {
		t.Fatalf("pending = %v", pending)
	}
	uploaded := LatestUploaded(snap)
	if len(uploaded) != 1 || uploaded[0] != 1 {
		t.Fatalf("uploaded = %v", uploaded)
	}
	if PendingAddrs(nil) != nil || LatestUploaded(nil) != nil {
		t.Fatal("nil snapshot handling failed")
	}
}

func TestApplyUploaded(t *testing.T) {
	m := New()
	table := BuildTable("t1", "v1", []model.Point{
		{RawAddr: 1, Name: "电压"},
		{RawAddr: 2, Name: "电流"},
	})
	if err := m.Load(table); err != nil {
		t.Fatal(err)
	}
	if err := m.ApplyTable("t1"); err != nil {
		t.Fatal(err)
	}
	if err := ApplyUploaded(m, []uint16{1}, "v1"); err != nil {
		t.Fatal(err)
	}
	point, ok := m.Resolve(1)
	if !ok || !point.Uploaded {
		t.Fatalf("point not uploaded: %v %v", point, ok)
	}
	if err := ApplyUploaded(m, []uint16{9}, "v1"); err != ErrPointMissing {
		t.Fatalf("missing point err = %v", err)
	}
}

func TestMapperTableOperations(t *testing.T) {
	m := New()
	if m.ActiveID() != "" || m.ActiveVersion() != "" {
		t.Fatal("empty mapper state wrong")
	}
	table := BuildTable("t1", "v1", []model.Point{{RawAddr: 1, Name: "电压"}})
	if err := m.Load(table); err != nil {
		t.Fatal(err)
	}
	if err := m.Load(table); err != ErrTableExists {
		t.Fatalf("duplicate load err = %v", err)
	}
	if err := m.ApplyTable("missing"); err != ErrTableMissing {
		t.Fatalf("missing table err = %v", err)
	}
	if err := m.ApplyTable("t1"); err != nil {
		t.Fatal(err)
	}
	if m.PointCount() != 1 {
		t.Fatalf("point count = %d", m.PointCount())
	}
	if err := m.MarkUploaded(9, "v1"); err != ErrPointMissing {
		t.Fatalf("mark missing err = %v", err)
	}
	snap := SnapshotTable(m)
	if snap.TableID != "t1" {
		t.Fatalf("snapshot table = %q", snap.TableID)
	}
}
