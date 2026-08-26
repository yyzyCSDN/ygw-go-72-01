package mapper

import (
	"testing"

	"powergw/internal/model"
)

func TestTableRegistry(t *testing.T) {
	registry := NewTableRegistry()
	table := BuildTable("t1", "v1", []model.Point{{RawAddr: 1, Name: "电压"}})
	if err := registry.Register(table); err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if err := registry.Register(table); err != ErrTableExists {
		t.Fatalf("duplicate register err = %v", err)
	}
	if registry.Count() != 1 {
		t.Fatalf("count = %d", registry.Count())
	}
	if got, ok := registry.Get("t1"); !ok || got.ID != "t1" {
		t.Fatalf("get failed: %v %v", got, ok)
	}
	if _, ok := registry.Get("missing"); ok {
		t.Fatal("missing table found")
	}
	if len(registry.All()) != 1 || len(registry.IDs()) != 1 {
		t.Fatalf("all/ids wrong: %d %d", len(registry.All()), len(registry.IDs()))
	}
}

func TestTableRegistryLoadAll(t *testing.T) {
	registry := NewTableRegistry()
	_ = registry.Register(BuildTable("t1", "v1", []model.Point{{RawAddr: 1, Name: "电压"}}))
	m := New()
	if err := registry.LoadAll(m); err != nil {
		t.Fatalf("load all failed: %v", err)
	}
	if _, ok := m.Table("t1"); !ok {
		t.Fatal("table not loaded")
	}
	if len(m.TableIDs()) != 1 {
		t.Fatalf("table ids = %v", m.TableIDs())
	}
}
