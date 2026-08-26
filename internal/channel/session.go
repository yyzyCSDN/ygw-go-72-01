package channel

type Handle struct {
	ChannelID string
	Open      bool
}

type SessionStore struct {
	handles map[string]*Handle
	order   []string
	open    int
	closed  int
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		handles: make(map[string]*Handle),
	}
}

func (s *SessionStore) Acquire(channelID string) *Handle {
	handle := &Handle{ChannelID: channelID, Open: true}
	s.handles[channelID] = handle
	s.order = append(s.order, channelID)
	s.open++
	return handle
}

func (s *SessionStore) Renew(channelID string) *Handle {
	handle := &Handle{ChannelID: channelID, Open: true}
	s.handles[channelID] = handle
	s.order = append(s.order, channelID)
	s.open++
	return handle
}

func (s *SessionStore) Release(channelID string) bool {
	handle, ok := s.handles[channelID]
	if !ok || !handle.Open {
		return false
	}
	handle.Open = false
	s.open--
	s.closed++
	return true
}

func (s *SessionStore) ReleaseAll() int {
	count := 0
	for _, channelID := range s.order {
		if s.Release(channelID) {
			count++
		}
	}
	return count
}

func (s *SessionStore) OpenCount() int {
	return s.open
}

func (s *SessionStore) ClosedCount() int {
	return s.closed
}

func (s *SessionStore) IsOpen(channelID string) bool {
	handle, ok := s.handles[channelID]
	return ok && handle.Open
}

func (s *SessionStore) Snapshot() map[string]bool {
	out := make(map[string]bool, len(s.handles))
	for channelID, handle := range s.handles {
		out[channelID] = handle.Open
	}
	return out
}
