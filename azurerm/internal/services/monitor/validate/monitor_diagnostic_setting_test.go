package validate

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestMonitorDiagnosticSettingName(t *testing.T) {
	cases := []struct {
		Name   string
		Errors int
	}{
		{
			Name:   "somename",
			Errors: 0,
		},
		{
			Name:   "",
			Errors: 1,
		},
		{
			Name:   acceptance.RandString(261),
			Errors: 1,
		},
		{
			Name:   "some<name",
			Errors: 1,
		},
		{
			Name:   "some>name",
			Errors: 1,
		},
		{
			Name:   "some*name",
			Errors: 1,
		},
		{
			Name:   "some%name",
			Errors: 1,
		},
		{
			Name:   "some&name",
			Errors: 1,
		},
		{
			Name:   "some:name",
			Errors: 1,
		},
		{
			Name:   "some\\name",
			Errors: 1,
		},
		{
			Name:   "some?name",
			Errors: 1,
		},
		{
			Name:   "some+name",
			Errors: 1,
		},
		{
			Name:   "some/name",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			_, errors := MonitorDiagnosticSettingName(tc.Name, "test")

			if len(errors) != tc.Errors {
				t.Fatalf("Expected Name to return %d error(s) not %d", tc.Errors, len(errors))
			}
		})
	}
}
