package model

import "testing"

func TestNewTableSnapshot(t *testing.T) {
	table := NewPointTable("t1", "v1")
	table.Add(Point{RawAddr: 1, Name: "电压"})
	table.Add(Point{RawAddr: 2, Name: "电流"})
	snap := NewTableSnapshot(table, 7)
	if snap.TableID != "t1" || snap.Version != "v1" || snap.Seq != 7 || snap.Len() != 2 {
		t.Fatalf("snapshot fields wrong: %+v", snap)
	}
	point, ok := snap.Find(1)
	if !ok || point.Name != "电压" {
		t.Fatalf("find failed: %v %v", point, ok)
	}
	if _, ok := snap.Find(99); ok {
		t.Fatal("missing point found")
	}
	if snap.Uploaded(1) {
		t.Fatal("point reported uploaded")
	}
}

func TestTableSnapshotMerge(t *testing.T) {
	base := NewPointTable("t1", "v1")
	base.Add(Point{RawAddr: 1, Name: "电压", Uploaded: true})
	base.Add(Point{RawAddr: 2, Name: "电流"})
	freshTable := NewPointTable("t1", "v1")
	freshTable.Add(Point{RawAddr: 2, Name: "电流", Uploaded: true})
	freshTable.Add(Point{RawAddr: 3, Name: "功率"})
	merged := NewTableSnapshot(base, 1).Merge(NewTableSnapshot(freshTable, 2))
	if merged.Len() != 3 {
		t.Fatalf("merged len = %d", merged.Len())
	}
	if !merged.Uploaded(1) || !merged.Uploaded(2) {
		t.Fatal("merged uploaded flags wrong")
	}
	withOne := merged.WithUploaded(3)
	if !withOne.Uploaded(3) {
		t.Fatal("with uploaded failed")
	}
	if merged.Uploaded(3) {
		t.Fatal("original snapshot mutated")
	}
}

func TestUploadStateMerge(t *testing.T) {
	base := NewUploadState("s1", 1)
	base.Done[1] = true
	base.Pending[2] = true
	fresh := NewUploadState("s1", 2)
	fresh.Done[2] = true
	fresh.Pending[3] = true
	merged := base.Merge(fresh)
	if merged.PendingCount() != 1 || merged.DoneCount() != 2 {
		t.Fatalf("merge counts wrong: pending=%d done=%d", merged.PendingCount(), merged.DoneCount())
	}
	if merged.Pending[2] {
		t.Fatal("done point still pending")
	}
	if !merged.Done[1] || !merged.Done[2] {
		t.Fatal("done flags missing")
	}
}
