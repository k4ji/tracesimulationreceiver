package task

// MarkAsFailedEffect is a conditional definition effect that marks the task as failed.
type MarkAsFailedEffect struct {
	message *string
}

// NewMarkAsFailedEffect creates a new MarkAsFailedEffect with the given message.
func NewMarkAsFailedEffect(message *string) MarkAsFailedEffect {
	return MarkAsFailedEffect{
		message: message,
	}
}

func (m *MarkAsFailedEffect) Message() *string {
	return m.message
}
