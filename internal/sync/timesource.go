package timesync

type TimeSource interface {
	Now() int64
}

type ClockSource struct {
	clock func() int64
}

func NewClockSource(clock func() int64) *ClockSource {
	return &ClockSource{clock: clock}
}

func (c *ClockSource) Now() int64 {
	if c == nil || c.clock == nil {
		return 0
	}
	return c.clock()
}
