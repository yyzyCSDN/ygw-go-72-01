package model

import "testing"

func TestPointTableAddGet(t *testing.T) {
	table := NewPointTable("t1", "v1")
	table.Add(Point{RawAddr: 1, Name: "电压", Unit: "kV"})
	table.Add(Point{RawAddr: 2, Name: "电流", Unit: "A"})
	if table.Len() != 2 {
		t.Fatalf("len = %d", table.Len())
	}
	point, ok := table.Get(2)
	if !ok || point.Name != "电流" || point.Version != "v1" {
		t.Fatalf("get failed: %v %v", point, ok)
	}
	table.Add(Point{RawAddr: 1, Name: "电压修正", Unit: "kV"})
	if table.Len() != 2 {
		t.Fatalf("duplicate add changed len: %d", table.Len())
	}
	if point, _ := table.Get(1); point.Name != "电压修正" {
		t.Fatalf("duplicate add did not overwrite: %v", point)
	}
}

func TestPointTableCopyMerge(t *testing.T) {
	source := NewPointTable("t1", "v1")
	source.Add(Point{RawAddr: 1, Name: "电压", Uploaded: true})
	source.Add(Point{RawAddr: 2, Name: "电流"})
	copyTable := source.Copy()
	if copyTable.ID != source.ID || copyTable.Len() != source.Len() {
		t.Fatalf("copy mismatch: %v %v", copyTable, source)
	}
	other := NewPointTable("t2", "v2")
	other.Add(Point{RawAddr: 2, Name: "电流", Uploaded: true})
	other.Add(Point{RawAddr: 3, Name: "功率"})
	merged := source.Merge(other)
	if merged.Len() != 3 {
		t.Fatalf("merged len = %d", merged.Len())
	}
	if point, _ := merged.Get(2); !point.Uploaded {
		t.Fatalf("merged point 2 not uploaded: %v", point)
	}
	if !source.Contains(1) || source.Contains(9) {
		t.Fatalf("contains check failed")
	}
}

func TestPointTableUploadedCount(t *testing.T) {
	table := NewPointTable("t1", "v1")
	table.Add(Point{RawAddr: 1, Name: "电压", Uploaded: true})
	table.Add(Point{RawAddr: 2, Name: "电流"})
	if table.UploadedCount() != 1 {
		t.Fatalf("uploaded count = %d", table.UploadedCount())
	}
	var nilTable *PointTable
	if nilTable.UploadedCount() != 0 || nilTable.Len() != 0 {
		t.Fatal("nil table count failed")
	}
}
