package version

import "powergw/internal/model"

type History struct {
	versions []model.ProtocolVersion
	limit    int
}

func NewHistory(limit int) *History {
	if limit <= 0 {
		limit = 64
	}
	return &History{limit: limit}
}

func (h *History) Append(v model.ProtocolVersion) bool {
	if h == nil {
		return false
	}
	h.versions = append(h.versions, v)
	for len(h.versions) > h.limit {
		h.versions = h.versions[1:]
	}
	return true
}

func (h *History) List() []model.ProtocolVersion {
	if h == nil {
		return nil
	}
	out := make([]model.ProtocolVersion, 0, len(h.versions))
	out = append(out, h.versions...)
	return out
}

func (h *History) Len() int {
	if h == nil {
		return 0
	}
	return len(h.versions)
}

func (h *History) Latest() (model.ProtocolVersion, bool) {
	if h == nil || len(h.versions) == 0 {
		return model.ProtocolVersion{}, false
	}
	return h.versions[len(h.versions)-1], true
}

func (h *History) Limit() int {
	if h == nil {
		return 0
	}
	return h.limit
}
