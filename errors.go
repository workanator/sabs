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

type ErrJobStartFailure struct {
	ID     interface{}
	Reason error
}

func (e ErrJobStartFailure) Error() string {
	return "job start failure: " + fmt.Sprint(e.ID) + ", " + e.Reason.Error()
}

type ErrJobStopFailure struct {
	ID     interface{}
	Reason error
}

func (e ErrJobStopFailure) Error() string {
	return "job stop failure: " + fmt.Sprint(e.ID) + ", " + e.Reason.Error()
}
