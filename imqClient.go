package imq

import (
	"context"
	"github.com/eskpil/imq/pkg/common"
	"github.com/eskpil/imq/pkg/protocol"
)

type client interface {
	Publish(ctx context.Context, message protocol.Message) error
	Listen(queue string, handler common.Handler) error
	Call(ctx context.Context, message protocol.Message) (*protocol.Message, error)
}
