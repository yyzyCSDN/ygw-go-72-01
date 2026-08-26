package parse

import "powergw/internal/model"

type VersionSource interface {
	ActiveFor(channelID string) (model.Protocol, string, error)
}

type Parser struct {
	registry *Registry
	states   map[string]*channelState
}

type channelState struct {
	versionID string
	proto     model.Protocol
	lastSeq   uint32
	batchSeq  uint32
	processed map[uint32]bool
	cached    bool
}

func NewParser(registry *Registry) *Parser {
	return &Parser{
		registry: registry,
		states:   make(map[string]*channelState),
	}
}

func (p *Parser) state(channelID string) *channelState {
	state, ok := p.states[channelID]
	if !ok {
		state = &channelState{processed: make(map[uint32]bool)}
		p.states[channelID] = state
	}
	return state
}

func (p *Parser) Parse(channelID string, raw []byte, proto model.Protocol, versionID string) (*model.Message, error) {
	if len(raw) == 0 {
		return nil, ErrEmptyFrame
	}
	if len(raw) > 512 {
		return nil, ErrFrameTooLong
	}
	state := p.state(channelID)
	message, err := p.registry.Parse(proto, raw)
	if err != nil {
		return nil, err
	}
	message.ChannelID = channelID
	if message.Seq == 0 {
		message.Seq = state.lastSeq + 1
	}
	message.VersionID = versionID
	if message.Seq > state.lastSeq {
		state.lastSeq = message.Seq
	}
	state.versionID = versionID
	state.proto = proto
	return message, nil
}

func (p *Parser) ParseVersioned(channelID string, raw []byte, source VersionSource) (*model.Message, error) {
	state := p.state(channelID)
	if !state.cached {
		if err := p.Refresh(channelID, source); err != nil {
			return nil, err
		}
	}
	return p.Parse(channelID, raw, state.proto, state.versionID)
}

func (p *Parser) Refresh(channelID string, source VersionSource) error {
	proto, versionID, err := source.ActiveFor(channelID)
	if err != nil {
		return err
	}
	state := p.state(channelID)
	state.proto = proto
	state.versionID = versionID
	state.cached = true
	return nil
}

func (p *Parser) MarkProcessed(channelID string, seq uint32) {
	state := p.state(channelID)
	state.processed[seq] = true
	if seq > state.batchSeq {
		state.batchSeq = seq
	}
}

func (p *Parser) Processed(channelID string, seq uint32) bool {
	state := p.state(channelID)
	return state.processed[seq]
}

func (p *Parser) LastSeq(channelID string) uint32 {
	return p.state(channelID).lastSeq
}

func (p *Parser) BatchSeq(channelID string) uint32 {
	return p.state(channelID).batchSeq
}

func (p *Parser) Reset(channelID string) {
	delete(p.states, channelID)
}

func (p *Parser) ResetAll() {
	p.states = make(map[string]*channelState)
}
