package imq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eskpil/imq/internal/table"
	"github.com/eskpil/imq/pkg/protocol"
	"net"
)

type Server struct {
	table       *table.Table
	connections map[string]*net.TCPConn
}

func (s *Server) Listen(queue string, handler Handler) error {
	return s.table.AddInternalEntry(queue, handler)
}

func (s *Server) Publish(ctx context.Context, message protocol.Message) error {
	entry, err := s.table.FindEntry(message.Publication.Queue)
	if err != nil {
		return err
	}

	switch entry.Kind {
	case table.EntryKindLocale:
		{
			c := upgradeContext(ctx)

			entry.LocalEntry.Handler(c, message)
		}
	case table.EntryKindRemote:
		{
			conn, ok := s.connections[entry.RemoteEntry.SessionId]
			if !ok {
				return fmt.Errorf("could not find handler for queue")
			}

			bytes, err := json.Marshal(message)
			if err != nil {
				return err
			}

			nwritten, err := conn.Write(bytes)
			if err != nil || nwritten != len(bytes) {
				return err
			}
		}
	}

	return nil
}
