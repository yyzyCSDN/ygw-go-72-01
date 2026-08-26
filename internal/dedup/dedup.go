package dedup

type Deduper struct {
	seen  map[string]map[uint64]bool
	order map[string][]uint64
	limit int
}

func NewDeduper(window *Window) *Deduper {
	return &Deduper{
		seen:  make(map[string]map[uint64]bool),
		order: make(map[string][]uint64),
		limit: window.Size(),
	}
}

func (d *Deduper) Check(channelID string, key uint64) (bool, error) {
	set := d.seen[channelID]
	if set == nil {
		return false, nil
	}
	return set[key], nil
}

func (d *Deduper) Mark(channelID string, key uint64) error {
	set := d.seen[channelID]
	if set == nil {
		set = make(map[uint64]bool)
		d.seen[channelID] = set
	}
	if set[key] {
		return nil
	}
	set[key] = true
	d.order[channelID] = append(d.order[channelID], key)
	for len(d.order[channelID]) > d.limit {
		oldest := d.order[channelID][0]
		d.order[channelID] = d.order[channelID][1:]
		delete(set, oldest)
	}
	return nil
}

func (d *Deduper) CheckOrMark(channelID string, key uint64) (bool, error) {
	duplicate, err := d.Check(channelID, key)
	if err != nil {
		return false, err
	}
	if duplicate {
		return true, nil
	}
	if err := d.Mark(channelID, key); err != nil {
		return false, err
	}
	return false, nil
}

func (d *Deduper) ResetAll() error {
	d.seen = make(map[string]map[uint64]bool)
	d.order = make(map[string][]uint64)
	return nil
}

func (d *Deduper) ResetChannel(channelID string) error {
	delete(d.seen, channelID)
	delete(d.order, channelID)
	return nil
}

func (d *Deduper) SeenCount() int {
	count := 0
	for _, set := range d.seen {
		count += len(set)
	}
	return count
}

func (d *Deduper) ChannelCount() int {
	return len(d.seen)
}

func (d *Deduper) Limit() int {
	return d.limit
}
