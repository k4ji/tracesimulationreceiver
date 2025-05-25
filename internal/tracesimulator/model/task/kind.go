package task

// Kind represents the kind of task
type Kind int

const (
	// KindUnknown represents an unknown task kind
	KindUnknown Kind = iota
	// KindClient represents a client task kind
	KindClient
	// KindServer represents a server task kind
	KindServer
	// KindProducer represents a producer task kind
	KindProducer
	// KindConsumer represents a consumer task kind
	KindConsumer
	// KindInternal represents an internal task kind
	KindInternal
)

func (k Kind) String() string {
	switch k {
	case KindClient:
		return "client"
	case KindServer:
		return "server"
	case KindProducer:
		return "producer"
	case KindConsumer:
		return "consumer"
	case KindInternal:
		return "internal"
	default:
		return "unknown"
	}
}

func FromString(kind string) Kind {
	switch kind {
	case "client":
		return KindClient
	case "server":
		return KindServer
	case "producer":
		return KindProducer
	case "consumer":
		return KindConsumer
	case "internal":
		return KindInternal
	default:
		return KindUnknown
	}
}
