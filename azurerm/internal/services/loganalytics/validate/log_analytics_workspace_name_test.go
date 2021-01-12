package validate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccLogAnalyticsWorkspaceName_validation(t *testing.T) {
	str := acctest.RandString(63)
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "abc",
			ErrCount: 1,
		},
		{
			Value:    "Ab-c",
			ErrCount: 0,
		},
		{
			Value:    "-abc",
			ErrCount: 1,
		},
		{
			Value:    "abc-",
			ErrCount: 1,
		},
		{
			Value:    str,
			ErrCount: 0,
		},
		{
			Value:    str + "a",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := LogAnalyticsWorkspaceName(tc.Value, "azurerm_log_analytics")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Log Analytics Workspace Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}
