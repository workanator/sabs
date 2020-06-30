package sabs

import "context"

type Server interface {
	Serve(ctx context.Context) (err error)
}

type Reloader interface {
	Reload(ctx context.Context, data interface{}) (err error)
}

type Shutdowner interface {
	Shutdown(ctx context.Context) (err error)
}

type Service interface {
	Server
	Reloader
	Shutdowner
}
