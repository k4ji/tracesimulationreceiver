package task

import "testing"

func TestNewExternalID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		expectErr bool
	}{
		{"alphanumeric are valid", "validID123", false},
		{"underscores are valid", "valid_id_123", false},
		{"hyphens are valid", "valid-id-123", false},
		{"spaces are invalid", "invalid id", true},
		{"special characters are invalid", "invalid@id", true},
		{"empty string is invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			externalID, err := NewExternalID(tt.id)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for invalid ID %s, got nil", tt.id)
				}
				if externalID != nil {
					t.Errorf("expected nil ExternalID for invalid ID %s, got %v", tt.id, externalID)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error for valid ID %s, got %v", tt.id, err)
				}
				if externalID == nil || externalID.Value() != tt.id {
					t.Errorf("expected ExternalID value to be %s, got %v", tt.id, externalID)
				}
			}
		})
	}
}
