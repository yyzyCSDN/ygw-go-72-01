package model

type VersionState int

const (
	VersionDraft VersionState = iota
	VersionActive
	VersionSuperseded
)

func (s VersionState) String() string {
	switch s {
	case VersionDraft:
		return "draft"
	case VersionActive:
		return "active"
	case VersionSuperseded:
		return "superseded"
	default:
		return "unknown"
	}
}

type ProtocolVersion struct {
	ID       string
	State    VersionState
	Proto    Protocol
	TableID  string
	Checksum uint64
}

func NewProtocolVersion(id string, proto Protocol, tableID string, checksum uint64) ProtocolVersion {
	return ProtocolVersion{
		ID:       id,
		State:    VersionDraft,
		Proto:    proto,
		TableID:  tableID,
		Checksum: checksum,
	}
}

type VersionSnapshot struct {
	Current *ProtocolVersion
	History []ProtocolVersion
	TableID string
}

func (s *VersionSnapshot) Active() (ProtocolVersion, bool) {
	if s == nil || s.Current == nil {
		return ProtocolVersion{}, false
	}
	return *s.Current, true
}

func (s *VersionSnapshot) Count() int {
	if s == nil {
		return 0
	}
	return len(s.History)
}
