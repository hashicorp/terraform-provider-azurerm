package azure

import (
	"testing"
)

func TestAzureRMApiManagementName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "a",
			ErrCount: 0,
		},
		{
			Value:    "abc",
			ErrCount: 0,
		},
		{
			Value:    "api1",
			ErrCount: 0,
		},
		{
			Value:    "company-api",
			ErrCount: 0,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "helloworld21!",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := ValidateApiManagementName(tc.Value, "azurerm_api_management")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the Api Management Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}
