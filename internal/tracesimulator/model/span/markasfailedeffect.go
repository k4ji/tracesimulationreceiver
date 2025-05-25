package span

var _ Effect = (*MarkAsFailedEffect)(nil)

// MarkAsFailedEffect is a conditional definition effect that marks the span as failed.
type MarkAsFailedEffect struct {
	message *string
}

func (m MarkAsFailedEffect) Apply(node *TreeNode) error {
	node.status = StatusError(m.message)
	return nil
}

// NewMarkAsFailedEffect creates a new MarkAsFailedEffect with the given message.
func NewMarkAsFailedEffect(message *string) MarkAsFailedEffect {
	return MarkAsFailedEffect{
		message: message,
	}
}
