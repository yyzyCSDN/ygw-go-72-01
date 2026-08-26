package mapper

import "powergw/internal/model"

func BuildTable(id, version string, points []model.Point) *model.PointTable {
	table := model.NewPointTable(id, version)
	for _, point := range points {
		table.Add(point)
	}
	return table
}

func ValidateTable(table *model.PointTable) error {
	if table == nil || table.ID == "" {
		return ErrTableMissing
	}
	if len(table.Ordered) == 0 {
		return ErrEmptyTable
	}
	seen := make(map[uint16]bool, len(table.Ordered))
	for _, addr := range table.Ordered {
		if seen[addr] {
			return ErrDuplicateAddr
		}
		seen[addr] = true
		if _, ok := table.Points[addr]; !ok {
			return ErrBrokenIndex
		}
	}
	return nil
}

func DefaultStationTables() []*model.PointTable {
	return []*model.PointTable{
		BuildTable("table-a", "v1", []model.Point{
			{RawAddr: 4001, Name: "线路电压", Unit: "kV", Scale: 0.01},
			{RawAddr: 4002, Name: "线路电流", Unit: "A", Scale: 0.01},
			{RawAddr: 4003, Name: "有功功率", Unit: "MW", Scale: 0.1},
			{RawAddr: 4004, Name: "无功功率", Unit: "MVar", Scale: 0.1},
			{RawAddr: 4005, Name: "频率", Unit: "Hz", Scale: 0.001},
			{RawAddr: 4006, Name: "温度", Unit: "C", Scale: 0.1},
		}),
		BuildTable("table-b", "v2", []model.Point{
			{RawAddr: 5001, Name: "母线电压", Unit: "kV", Scale: 0.01},
			{RawAddr: 5002, Name: "母线电流", Unit: "A", Scale: 0.01},
			{RawAddr: 5003, Name: "功率因数", Unit: "", Scale: 0.001},
			{RawAddr: 5004, Name: "视在功率", Unit: "MVA", Scale: 0.1},
			{RawAddr: 5005, Name: "油温", Unit: "C", Scale: 0.1},
		}),
		BuildTable("table-modbus", "m1", []model.Point{
			{RawAddr: 0, Name: "母线电压", Unit: "kV", Scale: 0.01},
			{RawAddr: 1, Name: "母线电流", Unit: "A", Scale: 0.01},
			{RawAddr: 2, Name: "有功功率", Unit: "MW", Scale: 0.1},
			{RawAddr: 3, Name: "频率", Unit: "Hz", Scale: 0.001},
		}),
	}
}
