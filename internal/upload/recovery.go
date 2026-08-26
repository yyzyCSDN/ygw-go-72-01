package upload

import "powergw/internal/model"

type Recovery struct {
	last *model.TableSnapshot
}

func NewRecovery() *Recovery {
	return &Recovery{}
}

func (r *Recovery) Store(snap *model.TableSnapshot) {
	r.last = snap
}

func (r *Recovery) Latest() *model.TableSnapshot {
	return r.last
}

func (r *Recovery) Has() bool {
	return r.last != nil
}

func (r *Recovery) Clear() {
	r.last = nil
}
