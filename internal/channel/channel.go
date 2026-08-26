package channel

import (
	"powergw/internal/model"
	"powergw/internal/parse"
	"powergw/internal/upload"
)

type Manager struct {
	records   map[string]model.ChannelRecord
	queues    map[string]*Queue
	processed map[string]map[uint32]bool
	seq       uint64
}

func NewManager(ids []string) *Manager {
	manager := &Manager{
		records:   make(map[string]model.ChannelRecord),
		queues:    make(map[string]*Queue),
		processed: make(map[string]map[uint32]bool),
	}
	for _, id := range ids {
		manager.records[id] = model.NewChannelRecord(id)
		manager.queues[id] = NewMessageQueue()
		manager.processed[id] = make(map[uint32]bool)
	}
	return manager
}

func (m *Manager) Has(id string) bool {
	_, ok := m.records[id]
	return ok
}

func (m *Manager) Get(id string) (model.ChannelRecord, bool) {
	record, ok := m.records[id]
	return record, ok
}

func (m *Manager) Records() []model.ChannelRecord {
	records := make([]model.ChannelRecord, 0, len(m.records))
	for _, record := range m.records {
		records = append(records, record)
	}
	return records
}

func (m *Manager) ChannelIDs() []string {
	ids := make([]string, 0, len(m.records))
	for id := range m.records {
		ids = append(ids, id)
	}
	return ids
}

func (m *Manager) Count() int {
	return len(m.records)
}

func (m *Manager) Advance(id string, next model.ChannelState) error {
	record, ok := m.records[id]
	if !ok {
		return ErrUnknownChannel
	}
	nextRecord, err := ApplyTransition(record, next)
	if err != nil {
		return err
	}
	m.records[id] = nextRecord
	return nil
}

func (m *Manager) Connect(id string) error {
	return m.Advance(id, model.ChannelConnected)
}

func (m *Manager) StartSync(id string) error {
	return m.Advance(id, model.ChannelSyncing)
}

func (m *Manager) Run(id string) error {
	return m.Advance(id, model.ChannelRunning)
}

func (m *Manager) Fault(id string, err error) error {
	record, ok := m.records[id]
	if !ok {
		return ErrUnknownChannel
	}
	nextRecord, transitionErr := ApplyTransition(record, model.ChannelFault)
	if transitionErr != nil {
		return transitionErr
	}
	if err != nil {
		nextRecord.Errors++
	}
	m.records[id] = nextRecord
	return nil
}

func (m *Manager) ApplySync(id string, seq uint64, syncedAt int64) error {
	record, ok := m.records[id]
	if !ok {
		return ErrUnknownChannel
	}
	if seq < record.SyncSeq {
		return ErrStaleSync
	}
	record.SyncSeq = seq
	record.SyncedAt = syncedAt
	m.records[id] = record
	return nil
}

func (m *Manager) RecordFrame(id string) {
	record, ok := m.records[id]
	if !ok {
		return
	}
	record.Frames++
	if record.State == model.ChannelConnected || record.State == model.ChannelSyncing {
		record.State = model.ChannelRunning
	}
	m.records[id] = record
}

func (m *Manager) RecordError(id string, err error) {
	record, ok := m.records[id]
	if !ok {
		return
	}
	record.Errors++
	m.records[id] = record
}

func (m *Manager) RotateSession(store *SessionStore, id string) (*Handle, error) {
	if !m.Has(id) {
		return nil, ErrUnknownChannel
	}
	if store.IsOpen(id) {
		store.Release(id)
	}
	return store.Acquire(id), nil
}

func (m *Manager) Ingest(parser *parse.Parser, source parse.VersionSource, id string, raw []byte) (*model.Message, error) {
	message, err := parser.ParseVersioned(id, raw, source)
	if err != nil {
		m.RecordError(id, err)
		return nil, err
	}
	m.RecordFrame(id)
	if err := m.forward(id, message); err != nil {
		return nil, err
	}
	return message, nil
}

func (m *Manager) forward(id string, message *model.Message) error {
	if message == nil {
		m.RecordError(id, parse.ErrEmptyFrame)
		return parse.ErrEmptyFrame
	}
	m.processed[id][message.Seq] = true
	return m.Enqueue(id, message)
}

func (m *Manager) Enqueue(id string, message *model.Message) error {
	queue, ok := m.queues[id]
	if !ok {
		return ErrUnknownChannel
	}
	return queue.Push(message)
}

func (m *Manager) QueueLen(id string) int {
	queue, ok := m.queues[id]
	if !ok {
		return 0
	}
	return queue.Len()
}

func (m *Manager) Dequeue(id string) (*model.Message, bool) {
	queue, ok := m.queues[id]
	if !ok {
		return nil, false
	}
	return queue.Pop()
}

func (m *Manager) ProcessBatch(parser *parse.Parser, id string, messages []*model.Message) (int, error) {
	if len(messages) == 0 {
		return 0, ErrEmptyBatch
	}
	if !m.Has(id) {
		return 0, ErrUnknownChannel
	}
	count := 0
	for _, message := range messages {
		if message == nil {
			continue
		}
		parser.MarkProcessed(id, message.Seq)
		m.processed[id][message.Seq] = true
		m.RecordFrame(id)
		count++
	}
	return count, nil
}

func (m *Manager) WritebackBatch(parser *parse.Parser, id string, staleSeqs []uint32) (int, error) {
	if !m.Has(id) {
		return 0, ErrUnknownChannel
	}
	latest := parser.BatchSeq(id)
	kept := 0
	for _, seq := range staleSeqs {
		if parser.Processed(id, seq) && seq <= latest {
			kept++
		}
	}
	return kept, nil
}

func (m *Manager) Processed(id string, seq uint32) bool {
	set, ok := m.processed[id]
	if !ok {
		return false
	}
	return set[seq]
}

func (m *Manager) ProcessedCount(id string) int {
	set, ok := m.processed[id]
	if !ok {
		return 0
	}
	count := 0
	for _, done := range set {
		if done {
			count++
		}
	}
	return count
}

func (m *Manager) ProcessUpload(id string, up *upload.Uploader, snap *model.TableSnapshot) (int, error) {
	count, err := up.Flush(snap)
	if err != nil {
		m.RecordError(id, err)
		return count, err
	}
	return count, nil
}

func (m *Manager) Snapshot() *model.ChannelSnapshot {
	m.seq++
	snap := model.NewChannelSnapshot(m.seq)
	for id, record := range m.records {
		snap.Records[id] = record
	}
	return snap
}
