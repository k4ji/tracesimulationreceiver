# TraceSimulationReceiver Tutorials

These short tutorials will help you get started with the `TraceSimulationReceiver` by showing how to generate synthetic traces using a simple OpenTelemetry Collector config.

> 💡 A JSON Schema file is available to help you write and validate your YAML config for the OpenTelemetry Collector with `TraceSimulationReceiver`.  
> You can find it in the repository: `otel_collector_with_tracesimulationreceiver.schema.json`

---

## 1. Getting Started

Create a file called `config.yaml`:

```yaml
receivers:
  tracesimulationreceiver:
    global:
      interval: 5s
    blueprint:
      type: service
      service:
        services:
          - name: service-a
            spans:
              - name: root-span
                delay:
                  for: 0s
                  as: absolute
                duration:
                  for: 1s
                  as: absolute

exporters:
  debug:
    verbosity: detailed

service:
  pipelines:
    traces:
      receivers: [tracesimulationreceiver]
      processors: []
      exporters: [debug]
```

Run it with the OpenTelemetry Collector (built with this receiver) or your custom image:

```bash
docker run --rm -v "$(pwd)/config.yaml:/etc/otelcol/config.yaml"  ghcr.io/k4ji/otelcol-tracesimulationreceiver:v0.4.0 --config /etc/otelcol/config.yaml
```

You should see traces printed in the console every 5 seconds.

---

## 2. Trace Interval & Defaults

To change the generation frequency, adjust `global.interval`:

```yaml
  tracesimulationreceiver:
    global:
      interval: 2s
```

To avoid repeating values for each span:

```yaml
  tracesimulationreceiver:
    blueprint:
      service:
        default:
          delay:
            for: "0.1"
            as: relative
          duration:
            for: 1s
            as: absolute
```

All spans will now inherit these values unless they override them.

---

## 3. Span Timing & Relationships

To simulate parent-child or linked spans:

```yaml
  tracesimulationreceiver:
    blueprint:
      service:
        services:
          - name: service-a
            spans:
              - name: root-span
                ref: root
                delay:
                  for: 0s
                  as: absolute
                duration:
                  for: 2s
                  as: absolute
          - name: service-b
            spans:
              - name: child-span
                ref: child
                parent: root
```

You can also define links:

```yaml
  tracesimulationreceiver:
    blueprint:
      service:
        services:
          - name: service-c
            spans:
              - name: linked-span
                delay:
                  for: 3s
                  as: absolute
                duration:
                  for: 1s
                  as: absolute
                links:
                  - child
```

---

## 4. Attributes & Events

Add resource or span-level metadata:

```yaml
  tracesimulationreceiver:
    blueprint:
      service:
        services:
          - name: service-a
            resource:
              os: linux
            spans:
              - name: root-span
                attributes:
                  team: team-a
```

Add events to a span:

```yaml
  tracesimulationreceiver:
    blueprint:
      service:
        services:
          - name: service-a
            spans:
              - name: root-span
                events:
                  - name: db.query
                    delay:
                      for: "0.1"
                      as: relative
                    attributes:
                      db.system.name: postgresql
```

---

## 5. Conditional Effects

Apply behavior randomly using `conditional_effects`:

```yaml
  tracesimulationreceiver:
    blueprint:
      service:
        services:
          - name: service-c
            spans:
              - name: linked-span
                delay:
                  for: 3s
                  as: absolute
                duration:
                  for: 1s
                  as: absolute
                links:
                  - child
                conditional_effects:
                  - condition:
                      kind: probabilistic
                      probabilistic:
                        threshold: 0.2
                    effects:
                      - kind: mark_as_failed
                        mark_as_failed:
                          message: "Simulated failure"
                      - kind: annotate
                        annotate:
                          attributes:
                            error.type: "SimError"
                      - kind: record_event
                        record_event:
                          event:
                            name: exception
                            delay:
                              for: "0.5"
                              as: relative
                            attributes:
                              exception.message: "Synthetic error"
```

This will mark 20% of these spans as errors with extra metadata.

---

## 6. Shifting the Trace End Time with `end_time_offset`

The timing of each span is calculated **relative to the end time of the last span in the longest trace**. By default, this end time is set to the current time, and all other spans' start and end times are back-calculated based on their delays and durations.

You can control this end time directly using the `end_time_offset` field:

```yaml
receivers:
  tracesimulationreceiver:
    global:
      end_time_offset: -15m
```    

With this feature, you can simulate delayed or backdated traces to 
- Test how telemetry backends handle late-arriving or out-of-order data
- Validate time-based retention or ingestion behavior in OpenTelemetry pipelines
    
---

You're ready to start simulating real-world tracing scenarios. Happy tracing!
