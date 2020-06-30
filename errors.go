package sabs

import "fmt"

type ErrResourceReloadFailure struct {
	ID     interface{}
	Reason error
}

func (e ErrResourceReloadFailure) Error() string {
	return "resource reload failure: " + fmt.Sprint(e.ID) + ", " + e.Reason.Error()
}

type ErrResourceShutdownFailure struct {
	ID     interface{}
	Reason error
}

func (e ErrResourceShutdownFailure) Error() string {
	return "resource shutdown failure: " + fmt.Sprint(e.ID) + ", " + e.Reason.Error()
}
