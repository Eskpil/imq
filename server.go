package imq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eskpil/imq/internal/table"
	"github.com/eskpil/imq/pkg/common"
	"github.com/eskpil/imq/pkg/protocol"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"sync"
)

type Server struct {
	table       *table.Table
	connections map[string]net.Conn

	mutex sync.Mutex
}

func NewServer() (*Server, error) {
	server := new(Server)

	server.table = table.New()
	server.connections = make(map[string]net.Conn, 32)

	return server, nil
}

func (s *Server) Listen(queue string, handler common.Handler) error {
	return s.table.AddInternalEntry(queue, handler)
}

func (s *Server) Publish(ctx context.Context, message protocol.Message) error {
	entry, err := s.table.FindEntry(message.Publication.Queue)
	if err != nil {
		return err
	}

	switch entry.Kind {
	case table.EntryKindLocal:
		{
			c := common.UpgradeContext(ctx)

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

func ackConnection(transaction string, conn net.Conn) error {
	message := protocol.Message{
		Transaction: transaction,
		Command:     protocol.CommandAck,
	}

	bytes, err := json.Marshal(&message)
	if err != nil {
		return err
	}

	nwritten, err := conn.Write(bytes)
	if err != nil {
		return err
	}

	if nwritten != len(bytes) {
		return fmt.Errorf("could not write bytes")
	}

	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	id := uuid.New().String()

	s.mutex.Lock()
	s.connections[id] = conn
	s.mutex.Unlock()

	// Handle the incoming connection
	for {
		// Read data from the connection
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			fmt.Println("Error reading from connection:", err.Error())
			break
		}

		message := new(protocol.Message)
		if err := json.Unmarshal(buffer[:n], message); err != nil {
			log.Fatalf("could not unmarshal data: %v", err)
			break
		}

		switch message.Command {
		case protocol.CommandPublish:
			if err := s.Publish(context.Background(), *message); err != nil {
				log.Fatalf("could not publish message: %v", err)
				break
			}
		case protocol.CommandListen:
			if err := s.table.AddRemoteEntry(message.Publication.Queue, id); err != nil {
				break
			}
		}

		if err := ackConnection(message.Transaction, conn); err != nil {
			log.Fatalf("could not ack connection: %v", err)
			break
		}

	}

	// Remove the connection from the server's connections map
	s.mutex.Lock()
	delete(s.connections, id)
	s.mutex.Unlock()

	// Close the connection when done
	conn.Close()
}

func (s *Server) Start() error {
	// Start listening for incoming connections
	listener, err := net.Listen("tcp", ":2000")
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Println("Listening on :2000")

	// Accept and handle incoming connections in separate goroutines
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		// Start a new goroutine to handle the connection concurrently
		go s.handleConnection(conn)
	}
}
