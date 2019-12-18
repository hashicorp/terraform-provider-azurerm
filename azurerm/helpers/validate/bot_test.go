package validate

import (
	"testing"
)

func TestValidateBotMSTeamsCallingWebHook(t *testing.T) {
	tests := []struct {
		webhook string
		valid   bool
	}{
		{
			webhook: "bad webhook",
			valid:   false,
		},
		{
			webhook: "http://badwebhook.com",
			valid:   false,
		},
		{
			webhook: "http://badwebhook.com/",
			valid:   false,
		},
		{
			webhook: "https://badwebhook.com",
			valid:   false,
		},
		{
			webhook: "https://goodwebhook.com/",
			valid:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.webhook, func(t *testing.T) {
			_, err := ValidateBotMSTeamsCallingWebHook()(tt.webhook, "")
			valid := err == nil
			if valid != tt.valid {
				t.Errorf("Expected valid status %t but got %t for input %s", tt.valid, valid, tt.webhook)
			}
		})
	}
}
