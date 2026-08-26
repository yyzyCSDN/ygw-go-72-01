package version

import (
	"errors"

	"powergw/internal/model"
)

type TableApplier interface {
	ApplyTable(id string) error
}

type DedupResetter interface {
	ResetAll() error
}

type Manager struct {
	current  *model.ProtocolVersion
	history  *History
	channels map[string]model.Protocol
	seq      uint64
}

func NewManager() *Manager {
	return &Manager{
		history:  NewHistory(64),
		channels: make(map[string]model.Protocol),
	}
}

func (m *Manager) RegisterChannel(channelID string, proto model.Protocol) {
	m.channels[channelID] = proto
}

func (m *Manager) Apply(ver model.ProtocolVersion, tables TableApplier) error {
	if m.current != nil && m.current.State == model.VersionSuperseded {
		return ErrSuperseded
	}
	if m.hasID(ver.ID) {
		return ErrVersionExists
	}
	ver.State = model.VersionActive
	m.seq++
	m.current = &ver
	m.history.Append(ver)
	if tables != nil {
		if err := tables.ApplyTable(ver.TableID); err != nil {
			return err
		}
	}
	for channelID := range m.channels {
		m.channels[channelID] = ver.Proto
	}
	return nil
}

func (m *Manager) Switch(ver model.ProtocolVersion, tables TableApplier, resetter DedupResetter) error {
	if err := m.Apply(ver, tables); err != nil {
		return err
	}
	if resetter != nil {
		if err := resetter.ResetAll(); err != nil {
			return errors.Join(ErrWriteback, err)
		}
	}
	return nil
}

func (m *Manager) ActiveFor(channelID string) (model.Protocol, string, error) {
	if m.current == nil {
		return 0, "", ErrNoActiveVersion
	}
	if proto, ok := m.channels[channelID]; ok {
		return proto, m.current.ID, nil
	}
	return m.current.Proto, m.current.ID, nil
}

func (m *Manager) Active() (model.ProtocolVersion, bool) {
	if m.current == nil {
		return model.ProtocolVersion{}, false
	}
	return *m.current, true
}

func (m *Manager) Snapshot() *model.VersionSnapshot {
	snap := &model.VersionSnapshot{
		History: m.history.List(),
	}
	if m.current != nil {
		current := *m.current
		snap.Current = &current
		snap.TableID = current.TableID
	}
	return snap
}

func (m *Manager) hasID(id string) bool {
	for _, item := range m.history.List() {
		if item.ID == id {
			return true
		}
	}
	return false
}

func (m *Manager) History() []model.ProtocolVersion {
	return m.history.List()
}

func (m *Manager) Count() int {
	return m.history.Len()
}

func (m *Manager) Supersede() {
	if m.current == nil {
		return
	}
	m.current.State = model.VersionSuperseded
}

func (m *Manager) ResetState() {
	m.current = nil
	m.history = NewHistory(64)
	m.seq = 0
}
