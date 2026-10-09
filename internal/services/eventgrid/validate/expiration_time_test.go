package validate_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/validate"
)

func TestExpirationTimeIfNotActivated(t *testing.T) {
	validTimes := []string{
		time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
	}
	for _, v := range validTimes {
		_, errors := validate.ExpirationTimeIfNotActivated()(v, "valid")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid `expiration_time_if_not_activated`: %q", v, errors)
		}
	}

	invalidTimes := []string{
		time.Now().Format(time.RFC3339),
		time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
		time.Now().Add(8 * 24 * time.Hour).Format(time.RFC3339),
		"invalid-time-format",
	}
	for _, v := range invalidTimes {
		_, errors := validate.ExpirationTimeIfNotActivated()(v, "invalid")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid `expiration_time_if_not_activated`", v)
		}
	}
}
