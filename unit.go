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
	unit.ctx, unit.cancel = context.WithCancel(ctx)
	for _, res := range unit.resources {
		unit.ctx = context.WithValue(unit.ctx, res.ID, res.Value)
	}

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

	if err = unit.service.Shutdown(ctx); err != nil {
		return
	}

	for _, res := range unit.resources {
		if v, ok := res.Value.(Shutdowner); ok && v != nil {
			if err = v.Shutdown(ctx); err != nil {
				return ErrResourceShutdownFailure{
					ID:     res.ID,
					Reason: err,
				}
			}
		} else if v, ok := res.Value.(io.Closer); ok && v != nil {
			if err = v.Close(); err != nil {
				return ErrResourceShutdownFailure{
					ID:     res.ID,
					Reason: err,
				}
			}
		}
	}

	return nil
}
