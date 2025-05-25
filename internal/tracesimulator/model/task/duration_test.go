package task_test

import (
	"testing"
	"time"

	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task"
	"github.com/k4ji/tracesimulationreceiver/internal/tracesimulator/model/task/taskduration"
	"github.com/stretchr/testify/assert"
)

func TestNewDuration(t *testing.T) {
	t.Run("valid expression", func(t *testing.T) {
		expr, _ := taskduration.NewAbsoluteDuration(5 * time.Second)
		duration, err := task.NewDuration(expr)

		assert.NoError(t, err)
		assert.NotNil(t, duration)
	})

	t.Run("nil expression", func(t *testing.T) {
		assert.Panics(t, func() {
			_, _ = task.NewDuration(nil)
		}, "Expected panic when expression is nil")
	})
}

func TestDuration_Resolve(t *testing.T) {
	t.Run("resolve absolute duration", func(t *testing.T) {
		expr, _ := taskduration.NewAbsoluteDuration(5 * time.Second)
		duration, _ := task.NewDuration(expr)

		result, err := duration.Resolve(nil)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 5*time.Second, *result)
	})

	t.Run("resolve absolute duration of zero", func(t *testing.T) {
		expr, _ := taskduration.NewAbsoluteDuration(0)
		duration, _ := task.NewDuration(expr)

		result, err := duration.Resolve(nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "duration must be greater than 0")
	})

	t.Run("resolve relative duration with valid context", func(t *testing.T) {
		baseDuration := 10 * time.Second
		expr, _ := taskduration.NewRelativeDuration(0.5)
		duration, _ := task.NewDuration(expr)

		result, err := duration.Resolve(&baseDuration)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 5*time.Second, *result)
	})

	t.Run("resolve relative duration with invalid context type", func(t *testing.T) {
		expr, _ := taskduration.NewRelativeDuration(0.5)
		duration, _ := task.NewDuration(expr)

		result, err := duration.Resolve("invalid context")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid context type")
	})

	t.Run("resolve relative duration with nil context", func(t *testing.T) {
		expr, _ := taskduration.NewRelativeDuration(0.5)
		duration, _ := task.NewDuration(expr)

		result, err := duration.Resolve(nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid context type")
	})

	t.Run("resolve relative duration of zero", func(t *testing.T) {
		baseDuration := 10 * time.Second
		expr, _ := taskduration.NewRelativeDuration(0.0)
		duration, _ := task.NewDuration(expr)

		result, err := duration.Resolve(&baseDuration)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "duration must be greater than 0")
	})
}
