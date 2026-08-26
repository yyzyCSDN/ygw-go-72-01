package service

import "time"

type Runner struct {
	gateway  *Gateway
	interval time.Duration
	stop     chan struct{}
	done     chan struct{}
}

func NewRunner(gateway *Gateway, interval time.Duration) *Runner {
	if interval <= 0 {
		interval = time.Second
	}
	return &Runner{
		gateway:  gateway,
		interval: interval,
		stop:     make(chan struct{}),
		done:     make(chan struct{}),
	}
}

func (r *Runner) Start() {
	go func() {
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()
		defer close(r.done)
		for {
			select {
			case <-ticker.C:
				_ = r.Cycle()
			case <-r.stop:
				return
			}
		}
	}()
}

func (r *Runner) Stop() {
	close(r.stop)
	<-r.done
}

func (r *Runner) Cycle() error {
	if _, err := r.gateway.SyncAll(); err != nil {
		return err
	}
	if _, err := r.gateway.FlushAll(); err != nil {
		return err
	}
	if _, err := r.gateway.RotateAll(); err != nil {
		return err
	}
	return nil
}
