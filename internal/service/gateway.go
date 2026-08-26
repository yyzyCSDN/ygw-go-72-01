package service

import (
	"time"

	"powergw/internal/channel"
	"powergw/internal/dedup"
	"powergw/internal/map"
	"powergw/internal/model"
	"powergw/internal/parse"
	"powergw/internal/sync"
	"powergw/internal/upload"
	"powergw/internal/version"
)

type MasterSink struct {
	items []model.Telemetry
	limit int
}

func NewMasterSink(limit int) *MasterSink {
	if limit <= 0 {
		limit = 512
	}
	return &MasterSink{limit: limit}
}

func (s *MasterSink) Send(tel model.Telemetry) error {
	if len(s.items) >= s.limit {
		s.items = s.items[1:]
	}
	s.items = append(s.items, tel)
	return nil
}

func (s *MasterSink) Items() []model.Telemetry {
	out := make([]model.Telemetry, 0, len(s.items))
	out = append(out, s.items...)
	return out
}

func (s *MasterSink) Count() int {
	return len(s.items)
}

func (s *MasterSink) Last() (model.Telemetry, bool) {
	if len(s.items) == 0 {
		return model.Telemetry{}, false
	}
	return s.items[len(s.items)-1], true
}

type Options struct {
	ChannelIDs []string
	Clock      func() int64
}

type Gateway struct {
	Parser     *parse.Parser
	Mapper     *mapper.Mapper
	Upload     *upload.Uploader
	Sync       *timesync.Sync
	Channels   *channel.Manager
	Version    *version.Manager
	Controller *version.Controller
	Dedup      *dedup.Deduper
	Sessions   *channel.SessionStore
	Recovery   *upload.Recovery
	Sink       *MasterSink
	Clock      func() int64
}

func NewGateway(opts Options) *Gateway {
	if opts.Clock == nil {
		opts.Clock = time.Now().Unix
	}
	mapperInstance := mapper.New()
	sink := NewMasterSink(512)
	source := timesync.NewClockSource(opts.Clock)
	versionManager := version.NewManager()
	return &Gateway{
		Parser:     parse.NewParser(parse.NewRegistry()),
		Mapper:     mapperInstance,
		Upload:     upload.NewUploader(mapperInstance, sink),
		Sync:       timesync.NewSync(source),
		Channels:   channel.NewManager(opts.ChannelIDs),
		Version:    versionManager,
		Controller: version.NewController(versionManager),
		Dedup:      dedup.NewDeduper(dedup.NewWindow(1024)),
		Sessions:   channel.NewSessionStore(),
		Recovery:   upload.NewRecovery(),
		Sink:       sink,
		Clock:      opts.Clock,
	}
}

func (g *Gateway) RegisterTable(table *model.PointTable) error {
	return g.Mapper.Load(table)
}

func (g *Gateway) ActivateTable(id string) error {
	return g.Mapper.ApplyTable(id)
}

func (g *Gateway) ActivateVersion(ver model.ProtocolVersion) error {
	return g.Controller.Activate(ver, g.Mapper)
}

func (g *Gateway) ApplyVersion(ver model.ProtocolVersion) error {
	return g.Controller.SwitchTo(ver, g.Mapper, g.Dedup)
}

func (g *Gateway) EstablishChannels() error {
	for _, id := range g.Channels.ChannelIDs() {
		if err := g.Channels.Connect(id); err != nil {
			return err
		}
		if err := g.Channels.StartSync(id); err != nil {
			return err
		}
		if err := g.Channels.Run(id); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gateway) FaultChannel(id string) error {
	return g.Channels.Fault(id, ErrChannelFaulted)
}

func (g *Gateway) Ingest(channelID string, raw []byte) error {
	if !g.Channels.Has(channelID) {
		return ErrUnknownChannel
	}
	message, err := g.Channels.Ingest(g.Parser, g.Version, channelID, raw)
	if err != nil {
		return err
	}
	key := dedup.Key(channelID, message.Addr, message.Seq)
	duplicate, err := g.Dedup.Check(channelID, key)
	if err != nil {
		return err
	}
	if duplicate {
		return ErrDuplicate
	}
	if err := g.Dedup.Mark(channelID, key); err != nil {
		return err
	}
	return SubmitMessage(g, message)
}

func (g *Gateway) ProcessQueued(channelID string) (int, error) {
	var messages []*model.Message
	for {
		message, ok := g.Channels.Dequeue(channelID)
		if !ok {
			break
		}
		messages = append(messages, message)
	}
	if len(messages) == 0 {
		return 0, nil
	}
	count, err := g.Channels.ProcessBatch(g.Parser, channelID, messages)
	if err != nil {
		return count, err
	}
	stale := make([]uint32, 0, len(messages))
	for _, message := range messages {
		stale = append(stale, message.Seq)
	}
	if _, err := g.Channels.WritebackBatch(g.Parser, channelID, stale); err != nil {
		return count, err
	}
	return count, nil
}

func (g *Gateway) FlushAll() (int, error) {
	for _, id := range g.Channels.ChannelIDs() {
		if _, err := g.ProcessQueued(id); err != nil {
			return 0, err
		}
	}
	snap := g.Mapper.Snapshot()
	total := 0
	for _, id := range g.Channels.ChannelIDs() {
		count, err := g.Channels.ProcessUpload(id, g.Upload, snap)
		if err != nil {
			return total, err
		}
		total += count
	}
	if err := g.Upload.Writeback(g.Mapper.Snapshot()); err != nil {
		return total, err
	}
	g.Recovery.Store(g.Mapper.Snapshot())
	return total, nil
}

func (g *Gateway) SyncAll() (int, error) {
	results := timesync.NewResultList()
	for _, id := range g.Channels.ChannelIDs() {
		results.Add(g.Sync.Generate(id))
	}
	return g.Sync.ApplyAll(results, g.Channels)
}

func (g *Gateway) RotateAll() (int, error) {
	count := 0
	for _, id := range g.Channels.ChannelIDs() {
		if _, err := g.Channels.RotateSession(g.Sessions, id); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (g *Gateway) Recover() error {
	if !g.Recovery.Has() {
		return nil
	}
	if err := g.Upload.Recover(g.Recovery.Latest()); err != nil {
		return err
	}
	g.Parser.ResetAll()
	g.Dedup.ResetAll()
	return nil
}

func (g *Gateway) CloseSessions() int {
	return g.Sessions.ReleaseAll()
}
