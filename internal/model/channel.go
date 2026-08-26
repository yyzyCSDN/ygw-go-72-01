package model

type ChannelState int

const (
	ChannelIdle ChannelState = iota
	ChannelConnected
	ChannelSyncing
	ChannelRunning
	ChannelFault
)

func (s ChannelState) String() string {
	switch s {
	case ChannelIdle:
		return "idle"
	case ChannelConnected:
		return "connected"
	case ChannelSyncing:
		return "syncing"
	case ChannelRunning:
		return "running"
	case ChannelFault:
		return "fault"
	default:
		return "unknown"
	}
}

func (s ChannelState) CanAccept() bool {
	return s == ChannelConnected || s == ChannelSyncing || s == ChannelRunning
}

func (s ChannelState) Terminal() bool {
	return s == ChannelFault
}

type ChannelRecord struct {
	ID       string
	State    ChannelState
	Session  string
	SyncedAt int64
	SyncSeq  uint64
	Frames   uint64
	Errors   uint64
}

func NewChannelRecord(id string) ChannelRecord {
	return ChannelRecord{
		ID:      id,
		State:   ChannelIdle,
		SyncSeq: 0,
	}
}

func (r *ChannelRecord) AdvanceState(next ChannelState) bool {
	if r == nil {
		return false
	}
	if r.State == next {
		return false
	}
	allowed := map[ChannelState]bool{
		ChannelConnected: r.State == ChannelIdle,
		ChannelSyncing:   r.State == ChannelConnected,
		ChannelRunning:   r.State == ChannelSyncing,
		ChannelFault:     r.State == ChannelConnected || r.State == ChannelSyncing || r.State == ChannelRunning,
	}
	if !allowed[next] {
		return false
	}
	r.State = next
	return true
}

type ChannelSnapshot struct {
	Records map[string]ChannelRecord
	Seq     uint64
}

func NewChannelSnapshot(seq uint64) *ChannelSnapshot {
	return &ChannelSnapshot{
		Records: make(map[string]ChannelRecord),
		Seq:     seq,
	}
}

func (s *ChannelSnapshot) Get(id string) (ChannelRecord, bool) {
	if s == nil {
		return ChannelRecord{}, false
	}
	record, ok := s.Records[id]
	return record, ok
}

func (s *ChannelSnapshot) Merge(fresh *ChannelSnapshot) *ChannelSnapshot {
	if s == nil {
		return fresh
	}
	if fresh == nil {
		return s
	}
	out := NewChannelSnapshot(fresh.Seq)
	for id, record := range s.Records {
		out.Records[id] = record
	}
	for id, record := range fresh.Records {
		if current, ok := out.Records[id]; ok && current.SyncSeq > record.SyncSeq {
			continue
		}
		out.Records[id] = record
	}
	return out
}
