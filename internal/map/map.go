package mapper

import (
	"errors"

	"powergw/internal/model"
)

var (
	ErrNoActiveTable = errors.New("no active table")
	ErrTableMissing  = errors.New("table missing")
	ErrTableExists   = errors.New("table already registered")
	ErrPointMissing  = errors.New("point missing")
	ErrEmptyTable    = errors.New("empty table")
	ErrDuplicateAddr = errors.New("duplicate point address")
	ErrBrokenIndex   = errors.New("broken point index")
)

type Mapper struct {
	tables     map[string]*model.PointTable
	snapshots  map[string]*model.TableSnapshot
	active     *model.PointTable
	activeSnap *model.TableSnapshot
	seq        uint64
}

func New() *Mapper {
	return &Mapper{
		tables:    make(map[string]*model.PointTable),
		snapshots: make(map[string]*model.TableSnapshot),
	}
}

func (m *Mapper) Load(table *model.PointTable) error {
	if table == nil || table.ID == "" {
		return ErrTableMissing
	}
	if _, exists := m.tables[table.ID]; exists {
		return ErrTableExists
	}
	m.tables[table.ID] = table.Copy()
	return nil
}

func (m *Mapper) ApplyTable(id string) error {
	table, ok := m.tables[id]
	if !ok {
		return ErrTableMissing
	}
	m.seq++
	m.active = table
	m.activeSnap = model.NewTableSnapshot(table, m.seq)
	m.snapshots[id] = m.activeSnap
	return nil
}

func (m *Mapper) Resolve(addr uint16) (model.Point, bool) {
	if m.active == nil {
		return model.Point{}, false
	}
	return m.active.Get(addr)
}

func (m *Mapper) MarkUploaded(addr uint16, version string) error {
	if m.active == nil {
		return ErrNoActiveTable
	}
	point, ok := m.active.Points[addr]
	if !ok {
		return ErrPointMissing
	}
	point.Uploaded = true
	point.Version = version
	m.active.Points[addr] = point
	m.refresh()
	return nil
}

func (m *Mapper) Snapshot() *model.TableSnapshot {
	if m.activeSnap == nil {
		return model.NewTableSnapshot(m.active, m.seq)
	}
	return m.activeSnap
}

func (m *Mapper) ActiveID() string {
	if m.active == nil {
		return ""
	}
	return m.active.ID
}

func (m *Mapper) ActiveVersion() string {
	if m.active == nil {
		return ""
	}
	return m.active.Version
}

func (m *Mapper) Table(id string) (*model.PointTable, bool) {
	table, ok := m.tables[id]
	return table, ok
}

func (m *Mapper) TableIDs() []string {
	ids := make([]string, 0, len(m.tables))
	for id := range m.tables {
		ids = append(ids, id)
	}
	return ids
}

func (m *Mapper) PointCount() int {
	if m.active == nil {
		return 0
	}
	return m.active.Len()
}

func (m *Mapper) UploadedCount() int {
	if m.active == nil {
		return 0
	}
	return m.active.UploadedCount()
}

func (m *Mapper) refresh() {
	if m.active == nil {
		return
	}
	m.activeSnap = model.NewTableSnapshot(m.active, m.seq)
	m.snapshots[m.active.ID] = m.activeSnap
}
