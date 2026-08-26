package timesync

type Result struct {
	ChannelID string
	Seq       uint64
	SyncedAt  int64
	Offset    int64
}

func (r Result) Valid() bool {
	return r.ChannelID != "" && r.Seq > 0
}

type ResultList struct {
	items []Result
}

func NewResultList() *ResultList {
	return &ResultList{}
}

func (l *ResultList) Add(result Result) {
	l.items = append(l.items, result)
}

func (l *ResultList) Items() []Result {
	out := make([]Result, 0, len(l.items))
	out = append(out, l.items...)
	return out
}

func (l *ResultList) Len() int {
	return len(l.items)
}

func (l *ResultList) Clear() {
	l.items = nil
}
