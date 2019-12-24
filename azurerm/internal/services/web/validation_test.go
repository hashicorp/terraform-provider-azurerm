package web

import "testing"

func TestAzureRMAppServicePlanName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "ab",
			ErrCount: 0,
		},
		{
			Value:    "abc",
			ErrCount: 0,
		},
		{
			Value:    "webapp1",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
			ErrCount: 0,
		},
		{
			Value:    "hello_world",
			ErrCount: 0,
		},
		{
			Value:    "helloworld21!",
			ErrCount: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateAppServicePlanName(tc.Value, "azurerm_app_service_plan")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the App Service Plan Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}

func TestAzureRMAppServiceName_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "ab",
			ErrCount: 0,
		},
		{
			Value:    "abc",
			ErrCount: 0,
		},
		{
			Value:    "webapp1",
			ErrCount: 0,
		},
		{
			Value:    "hello-world",
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
		_, errors := validateAppServiceName(tc.Value, "azurerm_app_service")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the App Service Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}
