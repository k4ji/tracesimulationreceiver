package span

import (
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
)

// Kind to represent the kind of task, which corresponds to the kind of span defined in OpenTelemetry
type Kind int

const (
	// KindUnknown represents an unknown span kind
	KindUnknown Kind = iota
	// KindClient represents a client span kind
	KindClient
	// KindServer represents a server span kind
	KindServer
	// KindProducer represents a producer span kind
	KindProducer
	// KindConsumer represents a consumer span kind
	KindConsumer
	// KindInternal represents an internal span kind
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

func FromTaskKind(kind task.Kind) Kind {
	switch kind {
	case task.KindClient:
		return KindClient
	case task.KindServer:
		return KindServer
	case task.KindProducer:
		return KindProducer
	case task.KindConsumer:
		return KindConsumer
	case task.KindInternal:
		return KindInternal
	default:
		return KindUnknown
	}
}
