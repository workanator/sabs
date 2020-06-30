package sabs

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	DefaultTerminationSignals = []os.Signal{
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGINT,
	}
	DefaultReloadSignals = []os.Signal{
		syscall.SIGHUP,
	}
)

func WaitTerminationSignal(
	ctx context.Context,
	srv Service,
	shutdownTimeout time.Duration,
	signals ...os.Signal,
) error {
	// Set signal channel.
	if len(signals) == 0 {
		signals = DefaultTerminationSignals
	}

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)

	// Wait for the one of termination signals and initiate service shutdown.
	select {
	case <-sigChan:
		var shutCtx context.Context
		if shutdownTimeout > 0 {
			var shutCancel context.CancelFunc
			shutCtx, shutCancel = context.WithTimeout(ctx, shutdownTimeout)
			defer shutCancel()
		} else {
			shutCtx = ctx
		}
		return srv.Shutdown(shutCtx)
	case <-ctx.Done():
		return ctx.Err()
	}
}

func WaitReloadSignal(
	ctx context.Context,
	srv Service,
	reloadTimeout time.Duration,
	data interface{},
	signals ...os.Signal,
) error {
	// Set signal channel.
	if len(signals) == 0 {
		signals = DefaultReloadSignals
	}

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)

	// Wait for the one of reload signals and initiate service reload.
	select {
	case <-sigChan:
		var shutCtx context.Context
		if reloadTimeout > 0 {
			var shutCancel context.CancelFunc
			shutCtx, shutCancel = context.WithTimeout(ctx, reloadTimeout)
			defer shutCancel()
		} else {
			shutCtx = ctx
		}
		return srv.Reload(shutCtx, data)
	case <-ctx.Done():
		return ctx.Err()
	}
}
