package model

type Point struct {
	RawAddr  uint16
	Name     string
	Unit     string
	Scale    float64
	Offset   float64
	Uploaded bool
	Version  string
}

type PointTable struct {
	ID      string
	Version string
	Points  map[uint16]Point
	Ordered []uint16
}

func NewPointTable(id, version string) *PointTable {
	return &PointTable{
		ID:      id,
		Version: version,
		Points:  make(map[uint16]Point),
	}
}

func (t *PointTable) Add(point Point) *PointTable {
	if t == nil {
		return t
	}
	if _, exists := t.Points[point.RawAddr]; !exists {
		t.Ordered = append(t.Ordered, point.RawAddr)
	}
	point.Version = t.Version
	t.Points[point.RawAddr] = point
	return t
}

func (t *PointTable) Get(addr uint16) (Point, bool) {
	if t == nil {
		return Point{}, false
	}
	point, ok := t.Points[addr]
	return point, ok
}

func (t *PointTable) Len() int {
	if t == nil {
		return 0
	}
	return len(t.Ordered)
}

func (t *PointTable) Contains(addr uint16) bool {
	_, ok := t.Get(addr)
	return ok
}

func (t *PointTable) Copy() *PointTable {
	if t == nil {
		return nil
	}
	out := NewPointTable(t.ID, t.Version)
	for _, addr := range t.Ordered {
		out.Add(t.Points[addr])
	}
	return out
}

func (t *PointTable) Merge(other *PointTable) *PointTable {
	if t == nil {
		return other
	}
	if other == nil {
		return t
	}
	out := t.Copy()
	for _, addr := range other.Ordered {
		point := other.Points[addr]
		if existing, ok := out.Points[addr]; ok && existing.Uploaded {
			point.Uploaded = true
		}
		out.Add(point)
	}
	return out
}

func (t *PointTable) UploadedCount() int {
	if t == nil {
		return 0
	}
	count := 0
	for _, point := range t.Points {
		if point.Uploaded {
			count++
		}
	}
	return count
}
