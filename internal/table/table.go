package table

import (
	"fmt"
	"github.com/eskpil/imq"
)

type EntryKind uint64

const (
	EntryKindLocale EntryKind = iota
	EntryKindRemote
)

type LocalEntry struct {
	Handler imq.Handler
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

func (t *Table) AddInternalEntry(queue string, handler imq.Handler) error {
	if _, has := t.Entries[queue]; !has {
		return fmt.Errorf("entry already exists")
	}

	entry := new(Entry)

	entry.Kind = EntryKindLocale
	entry.LocalEntry.Handler = handler
	entry.Queue = queue

	t.Entries[queue] = entry

	return nil
}

func (t *Table) FindEntry(queue string) (*Entry, error) {
	return nil, nil
}
