package model

type TableSnapshot struct {
	TableID string
	Version string
	Seq     uint64
	Points  []Point
}

func NewTableSnapshot(table *PointTable, seq uint64) *TableSnapshot {
	snap := &TableSnapshot{
		TableID: "",
		Version: "",
		Seq:     seq,
	}
	if table == nil {
		return snap
	}
	snap.TableID = table.ID
	snap.Version = table.Version
	snap.Points = make([]Point, 0, table.Len())
	for _, addr := range table.Ordered {
		if point, ok := table.Points[addr]; ok {
			snap.Points = append(snap.Points, point)
		}
	}
	return snap
}

func (s *TableSnapshot) Find(addr uint16) (Point, bool) {
	if s == nil {
		return Point{}, false
	}
	for _, point := range s.Points {
		if point.RawAddr == addr {
			return point, true
		}
	}
	return Point{}, false
}

func (s *TableSnapshot) Len() int {
	if s == nil {
		return 0
	}
	return len(s.Points)
}

func (s *TableSnapshot) Uploaded(addr uint16) bool {
	point, ok := s.Find(addr)
	if !ok {
		return false
	}
	return point.Uploaded
}

func (s *TableSnapshot) Merge(fresh *TableSnapshot) *TableSnapshot {
	if s == nil {
		return fresh
	}
	if fresh == nil {
		return s
	}
	byAddr := make(map[uint16]Point, len(s.Points)+len(fresh.Points))
	for _, point := range s.Points {
		byAddr[point.RawAddr] = point
	}
	for _, point := range fresh.Points {
		prev, ok := byAddr[point.RawAddr]
		if ok && prev.Uploaded {
			point.Uploaded = true
		}
		byAddr[point.RawAddr] = point
	}
	out := &TableSnapshot{
		TableID: fresh.TableID,
		Version: fresh.Version,
		Seq:     fresh.Seq,
	}
	seen := make(map[uint16]bool, len(byAddr))
	for _, point := range s.Points {
		if current, ok := byAddr[point.RawAddr]; ok && !seen[point.RawAddr] {
			out.Points = append(out.Points, current)
			seen[point.RawAddr] = true
		}
	}
	for _, point := range fresh.Points {
		if !seen[point.RawAddr] {
			out.Points = append(out.Points, point)
			seen[point.RawAddr] = true
		}
	}
	return out
}

func (s *TableSnapshot) WithUploaded(addr uint16) *TableSnapshot {
	if s == nil {
		return nil
	}
	out := &TableSnapshot{
		TableID: s.TableID,
		Version: s.Version,
		Seq:     s.Seq,
		Points:  make([]Point, 0, len(s.Points)),
	}
	for _, point := range s.Points {
		if point.RawAddr == addr {
			point.Uploaded = true
		}
		out.Points = append(out.Points, point)
	}
	return out
}
