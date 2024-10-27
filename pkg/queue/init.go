package queue

import (
	"context"
	"errors"
)

var (
	ErrAlreadyStarted  = errors.New("already started")
	ErrNilEventHandler = errors.New("event handler is nil")
)

type OnEventHandler func(ctx context.Context, key, value []byte) error
