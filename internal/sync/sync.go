package timesync

import (
	"powergw/internal/channel"
)

type Sync struct {
	source  TimeSource
	results map[string]uint64
	seq     uint64
}

func NewSync(source TimeSource) *Sync {
	return &Sync{
		source:  source,
		results: make(map[string]uint64),
	}
}

func (s *Sync) Generate(channelID string) Result {
	s.seq++
	return Result{
		ChannelID: channelID,
		Seq:       s.seq,
		SyncedAt:  s.source.Now(),
		Offset:    int64(s.seq % 100),
	}
}

func (s *Sync) Apply(result Result, channels *channel.Manager) error {
	if !result.Valid() {
		return ErrNoSource
	}
	last := s.results[result.ChannelID]
	if result.Seq < last {
		return ErrStaleResult
	}
	s.results[result.ChannelID] = result.Seq
	return channels.ApplySync(result.ChannelID, result.Seq, result.SyncedAt)
}

func (s *Sync) LastSeq(channelID string) uint64 {
	return s.results[channelID]
}

func (s *Sync) Count() int {
	return len(s.results)
}

func (s *Sync) Snapshot() map[string]uint64 {
	out := make(map[string]uint64, len(s.results))
	for channelID, seq := range s.results {
		out[channelID] = seq
	}
	return out
}

func (s *Sync) ApplyAll(results *ResultList, channels *channel.Manager) (int, error) {
	count := 0
	for _, result := range results.Items() {
		if err := s.Apply(result, channels); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}
