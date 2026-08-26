package model

type Quality int

const (
	QualityGood Quality = iota
	QualityInvalid
	QualityStale
)

func (q Quality) String() string {
	switch q {
	case QualityGood:
		return "good"
	case QualityInvalid:
		return "invalid"
	case QualityStale:
		return "stale"
	default:
		return "unknown"
	}
}

type Telemetry struct {
	PointName string
	RawAddr   uint16
	Value     float64
	Quality   Quality
	Seq       uint32
	Timestamp int64
	Uploaded  bool
	VersionID string
}

func NewTelemetry(name string, addr uint16, value float64, seq uint32, timestamp int64) Telemetry {
	return Telemetry{
		PointName: name,
		RawAddr:   addr,
		Value:     value,
		Quality:   QualityGood,
		Seq:       seq,
		Timestamp: timestamp,
	}
}

type UploadState struct {
	SnapshotID string
	Seq        uint64
	Pending    map[uint16]bool
	Done       map[uint16]bool
}

func NewUploadState(snapshotID string, seq uint64) *UploadState {
	return &UploadState{
		SnapshotID: snapshotID,
		Seq:        seq,
		Pending:    make(map[uint16]bool),
		Done:       make(map[uint16]bool),
	}
}

func (s *UploadState) PendingCount() int {
	if s == nil {
		return 0
	}
	return len(s.Pending)
}

func (s *UploadState) DoneCount() int {
	if s == nil {
		return 0
	}
	return len(s.Done)
}

func (s *UploadState) Merge(fresh *UploadState) *UploadState {
	if s == nil {
		return fresh
	}
	if fresh == nil {
		return s
	}
	out := NewUploadState(fresh.SnapshotID, fresh.Seq)
	for addr := range s.Done {
		out.Done[addr] = true
	}
	for addr := range fresh.Done {
		out.Done[addr] = true
	}
	for addr := range fresh.Pending {
		if !out.Done[addr] {
			out.Pending[addr] = true
		}
	}
	for addr := range s.Pending {
		if !out.Done[addr] {
			out.Pending[addr] = true
		}
	}
	return out
}
