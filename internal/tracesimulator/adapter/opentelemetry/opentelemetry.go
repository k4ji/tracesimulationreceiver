package opentelemetry

import (
	"fmt"
	simulator "github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/adapter"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/span"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.27.0"
)

const DefaultInstrumentationScopeName = "tracesimulator"

var _ simulator.Adapter[[]ptrace.Traces] = (*Adapter)(nil)

// Adapter is an OpenTelemetry adapter that transforms a tree of spans into OpenTelemetry format
type Adapter struct{}

func NewAdapter() *Adapter {
	return &Adapter{}
}

func (a *Adapter) Transform(rootSpans []*span.TreeNode) ([]ptrace.Traces, error) {
	var otelTraces []ptrace.Traces

	for _, rootSpan := range rootSpans {
		otelTrace := ptrace.NewTraces()
		if err := a.processNode(&otelTrace, rootSpan, nil); err != nil {
			return nil, fmt.Errorf("failed to transform root span '%s': %w", rootSpan.Name(), err)
		}
		otelTraces = append(otelTraces, otelTrace)
	}

	return otelTraces, nil
}

func (a *Adapter) processNode(otelTrace *ptrace.Traces, node *span.TreeNode, parentScopeSpans *ptrace.ScopeSpans) error {
	scopeSpans, err := a.createOrGetScopeSpans(otelTrace, node, parentScopeSpans)
	if err != nil {
		return fmt.Errorf("failed to process node '%s': %w", node.Name(), err)
	}

	a.addSpanToScope(*scopeSpans, node)

	for _, child := range node.Children() {
		if err := a.processNode(otelTrace, child, scopeSpans); err != nil {
			return err
		}
	}

	return nil
}

func (a *Adapter) createOrGetScopeSpans(otelTrace *ptrace.Traces, node *span.TreeNode, parentScopeSpans *ptrace.ScopeSpans) (*ptrace.ScopeSpans, error) {
	if node.IsResourceEntryPoint() {
		resourceSpans := otelTrace.ResourceSpans().AppendEmpty()
		resource := resourceSpans.Resource()
		resource.Attributes().PutStr(conventions.AttributeServiceName, node.Resource().Name())
		for k, v := range node.Resource().Attributes() {
			resource.Attributes().PutStr(k, v)
		}
		scopeSpans := resourceSpans.ScopeSpans().AppendEmpty()
		scopeSpans.Scope().SetName(DefaultInstrumentationScopeName)
		return &scopeSpans, nil
	}

	if parentScopeSpans == nil {
		return nil, fmt.Errorf("missing ScopeSpans for node '%s'", node.Name())
	}

	return parentScopeSpans, nil
}

func (a *Adapter) addSpanToScope(scopeSpans ptrace.ScopeSpans, node *span.TreeNode) {
	otelSpan := scopeSpans.Spans().AppendEmpty()
	otelSpan.SetTraceID(pcommon.TraceID(node.TraceID().Bytes()))
	otelSpan.SetSpanID(pcommon.SpanID(node.ID().Bytes()))
	otelSpan.SetName(node.Name())
	otelSpan.SetKind(toOtelKind(node.Kind()))
	otelSpan.SetStartTimestamp(pcommon.NewTimestampFromTime(node.StartTime()))
	otelSpan.SetEndTimestamp(pcommon.NewTimestampFromTime(node.EndTime()))
	setOtelStatusCode(&otelSpan, node.Status())

	for _, event := range node.Events() {
		otelEvent := otelSpan.Events().AppendEmpty()
		otelEvent.SetTimestamp(pcommon.NewTimestampFromTime(event.OccurredAt()))
		otelEvent.SetName(event.Name())
		for k, v := range event.Attributes() {
			otelEvent.Attributes().PutStr(k, v)
		}
	}

	for k, v := range node.Attributes() {
		otelSpan.Attributes().PutStr(k, v)
	}

	if node.ParentID() != nil {
		otelSpan.SetParentSpanID(pcommon.SpanID(node.ParentID().Bytes()))
	}

	for _, linked := range node.LinkedTo() {
		otelLink := otelSpan.Links().AppendEmpty()
		otelLink.SetTraceID(pcommon.TraceID(linked.TraceID().Bytes()))
		otelLink.SetSpanID(pcommon.SpanID(linked.ID().Bytes()))
	}
}

func toOtelKind(kind span.Kind) ptrace.SpanKind {
	switch kind {
	case span.KindServer:
		return ptrace.SpanKindServer
	case span.KindClient:
		return ptrace.SpanKindClient
	case span.KindProducer:
		return ptrace.SpanKindProducer
	case span.KindConsumer:
		return ptrace.SpanKindConsumer
	case span.KindInternal:
		return ptrace.SpanKindInternal
	default:
		return ptrace.SpanKindUnspecified
	}
}

func setOtelStatusCode(otelSpan *ptrace.Span, status span.Status) {
	switch status.Code() {
	case span.StatusCodeOK:
		otelSpan.Status().SetCode(ptrace.StatusCodeUnset)
	case span.StatusCodeError:
		otelSpan.Status().SetCode(ptrace.StatusCodeError)
		if status.Message() != nil {
			otelSpan.Status().SetMessage(*status.Message())
		}
	default:
		otelSpan.Status().SetCode(ptrace.StatusCodeUnset)
	}
}
