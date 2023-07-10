package common

import "github.com/eskpil/imq/pkg/protocol"

type Handler func(ctx Context, message protocol.Message) error
