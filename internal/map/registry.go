package mapper

import "powergw/internal/model"

type TableRegistry struct {
	tables map[string]*model.PointTable
	order  []string
}

func NewTableRegistry() *TableRegistry {
	return &TableRegistry{
		tables: make(map[string]*model.PointTable),
	}
}

func (r *TableRegistry) Register(table *model.PointTable) error {
	if table == nil || table.ID == "" {
		return ErrTableMissing
	}
	if _, exists := r.tables[table.ID]; exists {
		return ErrTableExists
	}
	r.tables[table.ID] = table
	r.order = append(r.order, table.ID)
	return nil
}

func (r *TableRegistry) Get(id string) (*model.PointTable, bool) {
	table, ok := r.tables[id]
	return table, ok
}

func (r *TableRegistry) All() []*model.PointTable {
	tables := make([]*model.PointTable, 0, len(r.order))
	for _, id := range r.order {
		tables = append(tables, r.tables[id])
	}
	return tables
}

func (r *TableRegistry) Count() int {
	return len(r.order)
}

func (r *TableRegistry) IDs() []string {
	out := make([]string, 0, len(r.order))
	out = append(out, r.order...)
	return out
}

func (r *TableRegistry) LoadAll(m *Mapper) error {
	for _, table := range r.All() {
		if err := m.Load(table); err != nil {
			return err
		}
	}
	return nil
}
