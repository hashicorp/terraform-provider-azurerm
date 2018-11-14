package azure

import "testing"

func TestHelper_AzureResourceID(t *testing.T) {
	cases := []struct {
		ID     string
		Errors int
	}{
		{
			ID:     "",
			Errors: 1,
		},
		{
			ID:     "nonsense",
			Errors: 1,
		},
		{
			ID:     "/slash",
			Errors: 1,
		},
		{
			ID:     "/path/to/nothing",
			Errors: 1,
		},
		{
			ID:     "/subscriptions",
			Errors: 1,
		},
		{
			ID:     "/providers",
			Errors: 1,
		},
		{
			ID:     "/subscriptions/not-a-guid",
			Errors: 0,
		},
		{
			ID:     "/providers/test",
			Errors: 0,
		},
		{
			ID:     "/subscriptions/00000000-0000-0000-0000-00000000000/",
			Errors: 0,
		},
		{
			ID:     "/providers/provider.name/",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.ID, func(t *testing.T) {
			_, errors := ValidateResourceID(tc.ID, "test")

			if len(errors) < tc.Errors {
				t.Fatalf("Expected ValidateResourceID to have %d not %d errors for %q", tc.Errors, len(errors), tc.ID)
			}
		})
	}
}

func TestAzureResourceIDOrEmpty(t *testing.T) {
	cases := []struct {
		ID     string
		Errors int
	}{
		{
			ID:     "",
			Errors: 0,
		},
		{
			ID:     "nonsense",
			Errors: 1,
		},
		//as this function just calls TestAzureResourceId lets not be as comprehensive
		{
			ID:     "/providers/provider.name/",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.ID, func(t *testing.T) {
			_, errors := ValidateResourceIDOrEmpty(tc.ID, "test")

			if len(errors) < tc.Errors {
				t.Fatalf("Expected TestAzureResourceIdOrEmpty to have %d not %d errors for %q", tc.Errors, len(errors), tc.ID)
			}
		})
	}
}

func TestAzureValidateMsSqlServiceName(t *testing.T) {
	cases := []struct {
		ServiceName string
		Errors      int
	}{
		{
			ServiceName: "as",
			Errors:      3,
		},
		{
			ServiceName: "Asd",
			Errors:      3,
		},
		{
			ServiceName: "asd",
			Errors:      0,
		},
		{
			ServiceName: "-asd",
			Errors:      3,
		},
		{
			ServiceName: "asd-",
			Errors:      3,
		},
		{
			ServiceName: "asd-1",
			Errors:      0,
		},
		{
			ServiceName: "asd--1",
			Errors:      1,
		},
		{
			ServiceName: "asd--1-",
			Errors:      4,
		},
		{
			ServiceName: "asdfghjklzasdfghjklzasdfghjklzasdfghjklzasdfghjklz",
			Errors:      0,
		},
		{
			ServiceName: "asdfghjklzasdfghjklzasdfghjklzasdfghjklzasdfghjklz1",
			Errors:      3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.ServiceName, func(t *testing.T) {
			_, errors := ValidateMsSqlServiceName(tc.ServiceName, "name")

			if len(errors) < tc.Errors {
				t.Fatalf("Expected TestAzureValidateMsSqlServiceName to have %d not %d errors for %q", tc.Errors, len(errors), tc.ServiceName)
			}
		})
	}
}

func TestAzureFloatAtLeast(t *testing.T) {
	cases := []struct {
		Name        string
		MinValue    float64
		ActualValue float64
		Errors      int
	}{
		{
			Name:        "Min_Full_Stop_Zero_Greater",
			MinValue:    0.0,
			ActualValue: 1.0,
			Errors:      0,
		},
		{
			Name:        "Min_One_Full_Stop_Zero_Lesser",
			MinValue:    1.0,
			ActualValue: 0.0,
			Errors:      1,
		},
		{
			Name:        "Min_Full_Stop_Two_Five_Greater",
			MinValue:    0.25,
			ActualValue: 0.26,
			Errors:      0,
		},
		{
			Name:        "Min_Full_Stop_Two_Five_Equal",
			MinValue:    0.25,
			ActualValue: 0.25,
			Errors:      0,
		},
		{
			Name:        "Min_Full_Stop_Two_Five_Lesser",
			MinValue:    0.25,
			ActualValue: 0.24,
			Errors:      1,
		},
		{
			Name:        "Min_Full_Stop_Long_Zero_Lesser",
			MinValue:    0.0000000000000000000000000000000000000001,
			ActualValue: 0,
			Errors:      1,
		},
		{
			Name:        "Min_Full_Stop_Long_Greater",
			MinValue:    0.0000000000000000000000000000000000000001,
			ActualValue: -0,
			Errors:      1,
		},
		{
			Name:        "Min_Negative_Full_Stop_Two_Five_Greater",
			MinValue:    -0.25,
			ActualValue: 1,
			Errors:      0,
		},
		{
			Name:        "Min_Zero_No_Full_Stop_Equal",
			MinValue:    0,
			ActualValue: -0,
			Errors:      0,
		},
		{
			Name:        "Min_Negative_Full_Stop_Two_Five_Lesser",
			MinValue:    -0.25,
			ActualValue: -0.26,
			Errors:      1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			validateFunc := FloatAtLeast(tc.MinValue)
			_, errors := validateFunc(tc.ActualValue, "floatValue")

			if len(errors) < tc.Errors {
				t.Fatalf("Expected FloatAtLeast to have %d not %d errors for %q", tc.Errors, len(errors), tc.Name)
			}
		})
	}
}
