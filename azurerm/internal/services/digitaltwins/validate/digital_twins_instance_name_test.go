package validate

import (
	"testing"
)

func TestDigitalTwinsInstanceName(t *testing.T) {
	tests := []struct {
		Name  string
		Input string
		Valid bool
	}{
		{
			Name:  "Empty",
			Input: "",
			Valid: false,
		},
		{
			Name:  "Too short",
			Input: "a",
			Valid: false,
		},
		{
			Name:  "Invalid character",
			Input: "digital_twins",
			Valid: false,
		},
		{
			Name:  "Valid Name",
			Input: "Digital-12-Twins",
			Valid: true,
		},
		{
			Name:  "End with `-`",
			Input: "Digital-12-",
			Valid: false,
		},
		{
			Name:  "Start with `-`",
			Input: "-Digital-12",
			Valid: false,
		},
		{
			Name:  "Invalid character",
			Input: "digital.twins",
			Valid: false,
		},
		{
			Name:  "Too long",
			Input: "digitalTwinsdigitalTwinsdigitalTwinsdigitalTwinsdigitalTwinsdigi",
			Valid: false,
		},
		{
			Name:  "Max characters",
			Input: "digitalTwinsdigitalTwinsdigitalTwinsdigitalTwinsdigitalTwins123",
			Valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := DigitalTwinsInstanceName(tt.Input, "Name")
			valid := err == nil
			if valid != tt.Valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.Valid, valid, tt.Input)
			}
		})
	}
}
