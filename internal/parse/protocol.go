package parse

import "powergw/internal/model"

type frameParser func(raw []byte) (*model.Message, error)

type Registry struct {
	parsers map[model.Protocol]frameParser
	order   []model.Protocol
}

func NewRegistry() *Registry {
	registry := &Registry{
		parsers: make(map[model.Protocol]frameParser),
	}
	registry.Register(model.ProtoIEC104, parseIEC104)
	registry.Register(model.ProtoModbus, parseModbus)
	return registry
}

func (r *Registry) Register(proto model.Protocol, fn frameParser) {
	if _, exists := r.parsers[proto]; !exists {
		r.order = append(r.order, proto)
	}
	r.parsers[proto] = fn
}

func (r *Registry) Parse(proto model.Protocol, raw []byte) (*model.Message, error) {
	fn, ok := r.parsers[proto]
	if !ok {
		return nil, ErrUnsupportedProtocol
	}
	return fn(raw)
}

func (r *Registry) Supports(proto model.Protocol) bool {
	_, ok := r.parsers[proto]
	return ok
}

func (r *Registry) Protocols() []model.Protocol {
	out := make([]model.Protocol, 0, len(r.order))
	out = append(out, r.order...)
	return out
}

func (r *Registry) Count() int {
	return len(r.order)
}
