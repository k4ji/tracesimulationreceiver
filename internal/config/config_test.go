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
									Delay: &service.Delay{
										Value: ptrString("100ms"),
										Mode:  ptrString("absolute"),
									},
									Duration: &service.Duration{
										Value: ptrString("200ms"),
										Mode:  ptrString("absolute"),
									},
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
						Delay: &service.Delay{
							Value: ptrString("0"),
							Mode:  ptrString("absolute"),
						},
						Duration: &service.Duration{
							Value: ptrString("1ns"),
							Mode:  ptrString("absolute"),
						},
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
									Delay: &service.Delay{
										Value: ptrString("0"),
										Mode:  ptrString("absolute"),
									},
									Duration: &service.Duration{
										Value: ptrString("1ns"),
										Mode:  ptrString("absolute"),
									},
								},
							},
						},
						{
							Name: "service2",
							Tasks: []service.Task{
								{
									Name:       "task2",
									ExternalID: &duplicateID,
									Delay: &service.Delay{
										Value: ptrString("0"),
										Mode:  ptrString("absolute"),
									},
									Duration: &service.Duration{
										Value: ptrString("1ns"),
										Mode:  ptrString("absolute"),
									},
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
						Delay:    nil,
						Duration: nil,
						FailWith: service.FailureCondition{
							Probability: nil,
						},
					},
					Services: []service.Service{
						{
							Name: "service1",
							Tasks: []service.Task{
								{
									Name: "task1",
									Delay: &service.Delay{
										Value: ptrString("0"),
										Mode:  ptrString("absolute"),
									}, Duration: &service.Duration{
										Value: ptrString("1ns"),
										Mode:  ptrString("absolute"),
									},
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

		t.Run("invalid Delay", func(t *testing.T) {
			cfg.Blueprint.ServiceBlueprint.Services[0].Tasks[0].Delay = &service.Delay{
				Value: ptrString("-1ns"),
				Mode:  ptrString("absolute"),
			}
			err := cfg.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid delay: absolute delay must be non-negative")
		})

		t.Run("invalid Duration", func(t *testing.T) {
			cfg.Blueprint.ServiceBlueprint.Services[0].Tasks[0].Delay = &service.Delay{
				Value: ptrString("100ms"),
				Mode:  ptrString("absolute"),
			}
			cfg.Blueprint.ServiceBlueprint.Services[0].Tasks[0].Duration = &service.Duration{
				Value: ptrString("0"),
				Mode:  ptrString("absolute"),
			}
			err := cfg.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "duration must be greater than 0")
		})

		t.Run("invalid FailWithProbability", func(t *testing.T) {
			cfg.Blueprint.ServiceBlueprint.Services[0].Tasks[0].Duration = &service.Duration{
				Value: ptrString("200ms"),
				Mode:  ptrString("absolute"),
			}
			cfg.Blueprint.ServiceBlueprint.Services[0].Tasks[0].FailWith.Probability = ptrFloat(1.5)
			err := cfg.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "must have a FailWithProbability value between 0.0 and 1.0")
		})
	})
}

func ptrFloat(f float64) *float64 { return &f }

func ptrString(s string) *string {
	return &s
}
