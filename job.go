package sabs

import "context"

type Starter interface {
	Start(ctx context.Context)
}

type Stopper interface {
	Stop(ctx context.Context) (err error)
}

type Job interface {
	Starter
	Stopper
}
