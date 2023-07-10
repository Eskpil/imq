package table

import (
	"fmt"
	"github.com/eskpil/imq/pkg/common"
)

type EntryKind uint64

const (
	EntryKindLocal EntryKind = iota
	EntryKindRemote
)

type LocalEntry struct {
	Handler common.Handler
}

type RemoteEntry struct {
	// Our server should know this. And can convert it into a socket connection.
	SessionId string
}

type Entry struct {
	Kind  EntryKind
	Queue string

	LocalEntry
	RemoteEntry
}

type Table struct {
	Entries map[string]*Entry
}

func New() *Table {
	table := new(Table)

	table.Entries = make(map[string]*Entry, 1000)

	return table
}

func (t *Table) AddInternalEntry(queue string, handler common.Handler) error {
	if _, has := t.Entries[queue]; has {
		return fmt.Errorf("entry already exists")
	}

	entry := new(Entry)

	entry.Kind = EntryKindLocal
	entry.LocalEntry.Handler = handler
	entry.Queue = queue

	t.Entries[queue] = entry

	return nil
}

func (t *Table) AddRemoteEntry(queue string, sessionId string) error {
	if _, has := t.Entries[queue]; has {
		return fmt.Errorf("entry already exists")
	}

	entry := new(Entry)

	entry.Kind = EntryKindRemote
	entry.RemoteEntry.SessionId = sessionId
	entry.Queue = queue

	t.Entries[queue] = entry

	return nil
}

func (t *Table) FindEntry(queue string) (*Entry, error) {
	entry, has := t.Entries[queue]
	if !has {
		return nil, fmt.Errorf("can not find queue")
	}

	return entry, nil
}
