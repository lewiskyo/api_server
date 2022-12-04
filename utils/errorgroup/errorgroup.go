package errgroup

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

// Execute represents a task.
type Execute func(ctx context.Context) error

// Group is a collection of tasks.
type Group struct {
	ctx     context.Context
	cancel  func()
	err     error
	errOnce sync.Once
	wg      sync.WaitGroup
	limit   int
	taskQ   chan Execute
	tasks   []Execute
}

// Option configures a group
type Option func(g *Group)

// WithCancel configures a group.
func WithCancel(ctx context.Context) Option {
	return func(g *Group) {
		g.ctx, g.cancel = context.WithCancel(ctx)
	}
}

// WithContext configures a group.
func WithContext(ctx context.Context) Option {
	return func(g *Group) {
		g.ctx = ctx
	}
}

// WithLimit limits the number of goroutines.
func WithLimit(l int) Option {
	return func(g *Group) {
		g.limit = l
	}
}

// New a group, and a group can't be reused if u have already called wait method.
func New(oo ...Option) *Group {
	g := &Group{}
	g.ctx = context.Background()
	g.limit = runtime.NumCPU()

	for _, o := range oo {
		o(g)
	}
	g.initCheck()

	g.schedule()

	return g
}

// Go execute a task or cache a task.
func (g *Group) Go(exec Execute) {
	g.wg.Add(1)
	select {
	case g.taskQ <- exec:
	default:
		g.tasks = append(g.tasks, exec)
	}
}

// Wait all tasks to be executed and returns the first err if any.
func (g *Group) Wait() error {
	for _, exec := range g.tasks {
		g.taskQ <- exec
	}
	g.wg.Wait()

	close(g.taskQ)

	if g.cancel != nil {
		g.cancel()
	}

	return g.err
}

func (g *Group) initCheck() {
	if g.limit <= 0 {
		panic("errgroup: limit should be greater than 0")
	}
}

func (g *Group) schedule() {
	g.taskQ = make(chan Execute, g.limit)
	for i := 0; i < g.limit; i++ {
		go func() {
			for f := range g.taskQ {
				g.do(f)
			}
		}()
	}
}

func (g *Group) do(exec Execute) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 65536) //nolint:gomnd
			buf = buf[:runtime.Stack(buf, false)]
			err = fmt.Errorf("errgroup recover: %s\n%s", r, buf)
		}
		if err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
		g.wg.Done()
	}()
	err = exec(g.ctx)
}
