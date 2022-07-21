package typescript

import (
	"context"

	"github.com/dop251/goja"
)

func startInterruptable(ctx context.Context, vm *goja.Runtime) chan struct{} {
	done := make(chan struct{})
	started := make(chan struct{})
	go func() {
		// Inform the parent go-routine that we've started, this prevents a race condition where the
		// runtime would beat the context cancellation in unit tests even though the context started
		// out in a 'cancelled' state.
		close(started)
		select {
		case <-ctx.Done():
			vm.Interrupt("context halt")
		case <-done:
			return
		}
	}()
	<-started
	return done
}
