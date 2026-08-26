package mapper

import (
	"testing"

	"powergw/internal/model"
)

func TestBuildTable(t *testing.T) {
	table := BuildTable("t1", "v1", []model.Point{
		{RawAddr: 1, Name: "电压"},
		{RawAddr: 2, Name: "电流"},
	})
	if table.Len() != 2 {
		t.Fatalf("len = %d", table.Len())
	}
	if err := ValidateTable(table); err != nil {
		t.Fatalf("validate failed: %v", err)
	}
}

func TestValidateTableErrors(t *testing.T) {
	if err := ValidateTable(nil); err != ErrTableMissing {
		t.Fatalf("nil table err = %v", err)
	}
	empty := model.NewPointTable("t1", "v1")
	if err := ValidateTable(empty); err != ErrEmptyTable {
		t.Fatalf("empty table err = %v", err)
	}
	dup := BuildTable("t1", "v1", []model.Point{{RawAddr: 1, Name: "电压"}})
	dup.Ordered = append(dup.Ordered, 1)
	if err := ValidateTable(dup); err != ErrDuplicateAddr {
		t.Fatalf("dup table err = %v", err)
	}
	broken := BuildTable("t1", "v1", []model.Point{{RawAddr: 1, Name: "电压"}})
	broken.Ordered = append(broken.Ordered, 9)
	if err := ValidateTable(broken); err != ErrBrokenIndex {
		t.Fatalf("broken table err = %v", err)
	}
}

func TestDefaultStationTables(t *testing.T) {
	tables := DefaultStationTables()
	if len(tables) < 2 {
		t.Fatalf("tables = %d", len(tables))
	}
	ids := map[string]bool{}
	for _, table := range tables {
		ids[table.ID] = true
		if err := ValidateTable(table); err != nil {
			t.Fatalf("table %s invalid: %v", table.ID, err)
		}
	}
	if !ids["table-a"] || !ids["table-modbus"] {
		t.Fatalf("missing default tables: %v", ids)
	}
}
