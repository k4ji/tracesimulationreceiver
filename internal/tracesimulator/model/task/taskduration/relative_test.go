package taskduration

import (
	"testing"
	"time"
)

func TestRelativeDuration_Resolve(t *testing.T) {
	tests := []struct {
		name        string
		value       float64
		context     interface{}
		expected    *time.Duration
		expectError bool
	}{
		{
			name:        "Valid relative duration",
			value:       0.5,
			context:     2 * time.Second,
			expected:    func() *time.Duration { d := 1 * time.Second; return &d }(),
			expectError: false,
		},
		{
			name:        "Invalid context type",
			value:       0.5,
			context:     "invalid",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Nil base duration",
			value:       0.5,
			context:     nil,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Negative base duration",
			value:       0.5,
			context:     -1 * time.Second,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "Zero base duration",
			value:       0.5,
			context:     time.Duration(0),
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd, err := NewRelativeDuration(tt.value)
			if err != nil {
				t.Fatalf("failed to create RelativeDuration: %v", err)
			}
			result, err := rd.Resolve(tt.context)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %v", err)
				}
				if result == nil || *result != *tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}
