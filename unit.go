package sabs

import (
	"context"
	"io"
)

type Resource struct {
	ID, Value interface{}
}

func Attach(id, value interface{}) Resource {
	return Resource{
		ID:    id,
		Value: value,
	}
}

type Unit struct {
	service   Service
	resources []Resource
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewUnit(srv Service, resources ...Resource) *Unit {
	return &Unit{
		service:   srv,
		resources: resources,
	}
}

func (unit *Unit) Serve(ctx context.Context) (err error) {
	// Propagate resources to the context.
	unit.ctx, unit.cancel = context.WithCancel(ctx)
	for _, res := range unit.resources {
		unit.ctx = context.WithValue(unit.ctx, res.ID, res.Value)
	}

	// Start all resources which conforms Starter interface.
	for _, res := range unit.resources {
		if starter, ok := res.Value.(Starter); ok {
			go starter.Start(unit.ctx)
		}
	}

	// Run the service.
	return unit.service.Serve(unit.ctx)
}

func (unit *Unit) Reload(ctx context.Context, data interface{}) (err error) {
	if err = unit.service.Reload(ctx, data); err != nil {
		return
	}

	for _, res := range unit.resources {
		if v, ok := res.Value.(Reloader); ok && v != nil {
			if err = v.Reload(ctx, data); err != nil {
				return ErrResourceReloadFailure{
					ID:     res.ID,
					Reason: err,
				}
			}
		}
	}

	return nil
}

func (unit *Unit) Shutdown(ctx context.Context) (err error) {
	defer unit.cancel()

	// Shutdown the server first.
	if err = unit.service.Shutdown(ctx); err != nil {
		return
	}

	// Shutdown resources.
	for _, res := range unit.resources {
		// Stop the job if the resource confirms Stopper interface.
		if stopper, ok := res.Value.(Stopper); ok && stopper != nil {
			if err = stopper.Stop(ctx); err != nil {
				return ErrJobStopFailure{
					ID:     res.ID,
					Reason: err,
				}
			}
		}

		// Shutdown the resource, or close it if possible.
		if shutdowner, ok := res.Value.(Shutdowner); ok && shutdowner != nil {
			if err = shutdowner.Shutdown(ctx); err != nil {
				return ErrResourceShutdownFailure{
					ID:     res.ID,
					Reason: err,
				}
			}
		} else if closer, ok := res.Value.(io.Closer); ok && closer != nil {
			if err = closer.Close(); err != nil {
				return ErrResourceShutdownFailure{
					ID:     res.ID,
					Reason: err,
				}
			}
		}
	}

	return nil
}
