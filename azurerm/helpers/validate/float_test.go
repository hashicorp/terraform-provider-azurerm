package validate

import "testing"

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
			_, errors := FloatAtLeast(tc.MinValue)(tc.ActualValue, "floatValue")

			if len(errors) < tc.Errors {
				t.Fatalf("Expected FloatAtLeast to have %d not %d errors for %q", tc.Errors, len(errors), tc.Name)
			}
		})
	}
}
