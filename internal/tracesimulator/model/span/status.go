package span

type StatusCode uint8

const (
	StatusCodeUnset StatusCode = iota
	StatusCodeOK
	StatusCodeError
)

const StatusCodeUnsetString = "unset"
const StatusCodeOKString = "ok"
const StatusCodeErrorString = "error"
const StatusCodeUnknownString = "unknown"

func (s StatusCode) String() string {
	switch s {
	case StatusCodeUnset:
		return StatusCodeUnsetString
	case StatusCodeOK:
		return StatusCodeOKString
	case StatusCodeError:
		return StatusCodeErrorString
	default:
		return StatusCodeUnknownString
	}
}

var (
	StatusOK    = NewStatus(StatusCodeOK, nil)
	StatusUnset = NewStatus(StatusCodeUnset, nil)
	StatusError = func(msg *string) Status { return NewStatus(StatusCodeError, msg) }
)

// Status represents the status of a span
type Status struct {
	code    StatusCode
	message *string
}

func NewStatus(code StatusCode, message *string) Status {
	return Status{
		code:    code,
		message: message,
	}
}

func (s *Status) Code() StatusCode {
	return s.code
}

func (s *Status) Message() *string {
	return s.message
}
