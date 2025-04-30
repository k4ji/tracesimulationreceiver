package config

import (
	"github.com/k4ji/tracesimulationreceiver/internal/config/blueprint"
	"github.com/k4ji/tracesimulationreceiver/internal/config/blueprint/service"
	"github.com/k4ji/tracesimulationreceiver/internal/config/global"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	t.Run("valid configuration", func(t *testing.T) {
		cfg := Config{
			Global: global.Default(),
			Blueprint: blueprint.Blueprint{
				Type: "service",
				ServiceBlueprint: &service.Blueprint{
					Services: []service.Service{
						{
							Name: "service1",
							Tasks: []service.Task{
								{
									Name:       "task1",
									ExternalID: ptrString("task1-id"),
									StartAfter: ptrDuration(100 * time.Millisecond),
									Duration:   ptrDuration(200 * time.Millisecond),
									FailWith: service.FailureCondition{
										Probability: ptrFloat(0.5),
									},
								},
							},
						},
					},
				},
			},
		}
		err := cfg.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid global interval", func(t *testing.T) {
		cfg := Config{
			Global: global.Global{
				Interval: time.Duration(0),
			},
		}
		err := cfg.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "global interval must be greater than 0")
	})

	t.Run("duplicate task ids", func(t *testing.T) {
		duplicateID := "task-id"
		cfg := Config{
			Global: global.Default(),
			Blueprint: blueprint.Blueprint{
				Type: "service",
				ServiceBlueprint: &service.Blueprint{
					Default: service.DefaultValues{
						StartAfter: ptrDuration(0),
						Duration:   ptrDuration(1),
						FailWith: service.FailureCondition{
							Probability: ptrFloat(0.0),
						},
					},
					Services: []service.Service{
						{
							Name: "service1",
							Tasks: []service.Task{
								{
									Name:       "task1",
									ExternalID: &duplicateID,
								},
							},
						},
						{
							Name: "service2",
							Tasks: []service.Task{
								{
									Name:       "task2",
									ExternalID: &duplicateID,
								},
							},
						},
					},
				},
			},
		}
		err := cfg.Validate()
		assert.EqualError(t, err, "blueprint validation failed: service blueprint validation failed: duplicate task ID task-id found")
	})

	t.Run("invalid task properties", func(t *testing.T) {
		cfg := Config{
			Global: global.Default(),
			Blueprint: blueprint.Blueprint{
				Type: "service",
				ServiceBlueprint: &service.Blueprint{
					Default: service.DefaultValues{
						StartAfter: nil,
						Duration:   nil,
						FailWith: service.FailureCondition{
							Probability: nil,
						},
					},
					Services: []service.Service{
						{
							Name: "service1",
							Tasks: []service.Task{
								{
									Name:       "task1",
									StartAfter: ptrDuration(0),
									Duration:   ptrDuration(1),
									FailWith: service.FailureCondition{
										Probability: ptrFloat(1.0),
									},
								},
							},
						},
					},
				},
			},
		}

		t.Run("invalid StartAfter", func(t *testing.T) {
			cfg.Blueprint.ServiceBlueprint.Services[0].Tasks[0].StartAfter = ptrDuration(-1)
			err := cfg.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "must have a StartAfter value greater than or equal to 0")
		})

		t.Run("invalid Duration", func(t *testing.T) {
			cfg.Blueprint.ServiceBlueprint.Services[0].Tasks[0].StartAfter = ptrDuration(100 * time.Millisecond)
			cfg.Blueprint.ServiceBlueprint.Services[0].Tasks[0].Duration = ptrDuration(0)
			err := cfg.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "must have a Duration value greater than 0")
		})

		t.Run("invalid FailWithProbability", func(t *testing.T) {
			cfg.Blueprint.ServiceBlueprint.Services[0].Tasks[0].Duration = ptrDuration(200 * time.Millisecond)
			cfg.Blueprint.ServiceBlueprint.Services[0].Tasks[0].FailWith.Probability = ptrFloat(1.5)
			err := cfg.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "must have a FailWithProbability value between 0.0 and 1.0")
		})
	})
}

func ptrDuration(d time.Duration) *time.Duration {
	return &d
}

func ptrFloat(f float64) *float64 { return &f }

func ptrString(s string) *string {
	return &s
}
