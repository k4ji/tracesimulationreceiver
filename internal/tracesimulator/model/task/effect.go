package task

// EffectKind defines the type of effect to apply.
type EffectKind string

const (
	EffectKindMarkAsFailed EffectKind = "markAsFailed"
	EffectKindRecordEvent  EffectKind = "recordEvent"
	EffectKindAnnotate     EffectKind = "annotate"
)

type Effect struct {
	kind         EffectKind
	markAsFailed *MarkAsFailedEffect
	recordEvent  *RecordEventEffect
	annotate     *AnnotateEffect
}

func FromMarkAsFailedEffect(markAsFailed MarkAsFailedEffect) Effect {
	return Effect{
		kind:         EffectKindMarkAsFailed,
		markAsFailed: &markAsFailed,
	}
}

func FromRecordEventEffect(recordEvent RecordEventEffect) Effect {
	return Effect{
		kind:        EffectKindRecordEvent,
		recordEvent: &recordEvent,
	}
}

func FromAnnotateEffect(annotate AnnotateEffect) Effect {
	return Effect{
		kind:     EffectKindAnnotate,
		annotate: &annotate,
	}
}

func (e *Effect) Kind() EffectKind {
	return e.kind
}

func (e *Effect) MarkAsFailedEffect() *MarkAsFailedEffect {
	return e.markAsFailed
}

func (e *Effect) RecordEventEffect() *RecordEventEffect {
	return e.recordEvent
}

func (e *Effect) AnnotateEffect() *AnnotateEffect {
	return e.annotate
}
