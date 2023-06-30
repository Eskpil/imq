package imq

import (
	"context"
	"github.com/eskpil/imq/pkg/protocol"
)

type Handler func(ctx Context, message protocol.Message)

type client interface {
	Publish(ctx context.Context, message protocol.Message) error
	Listen(queue string, handler Handler) error
	Call(ctx context.Context, message protocol.Message) (*protocol.Message, error)
}
