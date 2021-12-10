package validate

import "testing"

func TestAppServicePlanName_validation(t *testing.T) {
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
		_, errors := AppServicePlanName(tc.Value, "azurerm_app_service_plan")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected the App Service Plan Name to trigger a validation error for '%s'", tc.Value)
		}
	}
}
