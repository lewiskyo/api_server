package shutdown

import (
	"context"
	"sync"
)

type (
	Shutdown struct {
		ctx       context.Context
		cancel    context.CancelFunc
		waitGroup sync.WaitGroup
	}
)

func New() *Shutdown {
	ctx, cancel := context.WithCancel(context.TODO())
	return &Shutdown{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Shutdown) Run(run func(ctx context.Context)) {
	s.waitGroup.Add(1)
	defer s.waitGroup.Done()
	run(s.ctx)
}

func (s *Shutdown) Go(run func(ctx context.Context)) {
	s.waitGroup.Add(1)
	go func() {
		defer s.waitGroup.Done()
		run(s.ctx)
	}()
}

func (s *Shutdown) Shutdown() {
	s.cancel()
}

func (s *Shutdown) Wait() {
	s.waitGroup.Wait()
}

func (s *Shutdown) IsShutdown() bool {
	select {
	case <-s.ctx.Done():
		return true
	default:
		return false
	}
}
