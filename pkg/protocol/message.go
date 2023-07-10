package protocol

type Command uint64

const (
	CommandUnknown Command = iota

	// CommandAck used by both sides
	CommandAck

	// CommandListen is only used by an external client. It must be acked
	CommandListen

	// CommandPublish is used by both sides. It must be acked
	CommandPublish
)

type Publication struct {
	Queue string `json:"queue"`
}

type Message struct {
	Transaction string `json:"transaction"`

	Command     `json:"command"`
	Publication `json:"publication"`

	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}
