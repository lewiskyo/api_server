package shutdown

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	shutdown := New()
	shutdown.Go(func(ctx context.Context) {
		<-ctx.Done()
	})
	go func() {
		shutdown.Run(func(ctx context.Context) {
			<-ctx.Done()
		})
	}()

	if shutdown.IsShutdown() {
		t.Error()
	}
	go func() {
		shutdown.Shutdown()
	}()

	shutdown.Wait()
	if !shutdown.IsShutdown() {
		t.Error("")
	}
}
