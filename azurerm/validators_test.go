package azurerm

import (
	"testing"
)

func TestValidateRFC3339Date(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "",
			ErrCount: 1,
		},
		{
			Value:    "Random",
			ErrCount: 1,
		},
		{
			Value:    "2017-01-01",
			ErrCount: 1,
		},
		{
			Value:    "2017-01-01T01:23:45",
			ErrCount: 1,
		},
		{
			Value:    "2017-01-01T01:23:45+00:00",
			ErrCount: 0,
		},
		{
			Value:    "2017-01-01T01:23:45Z",
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validateRFC3339Date(tc.Value, "example")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected validateRFC3339Date to trigger '%d' errors for '%s' - got '%d'", tc.ErrCount, tc.Value, len(errors))
		}
	}
}

func TestValidateIso8601Duration(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			// Date components only
			Value:  "P1Y2M3D",
			Errors: 0,
		},
		{
			// Time components only
			Value:  "PT7H42M3S",
			Errors: 0,
		},
		{
			// Date and time components
			Value:  "P1Y2M3DT7H42M3S",
			Errors: 0,
		},
		{
			// Invalid prefix
			Value:  "1Y2M3DT7H42M3S",
			Errors: 1,
		},
		{
			// Wrong order of components, i.e. invalid format
			Value:  "PT7H42M3S1Y2M3D",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateIso8601Duration()(tc.Value, "example")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected validateIso8601Duration to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}

func TestValidateCollation(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "en-US",
			Errors: 1,
		},
		{
			Value:  "en_US",
			Errors: 0,
		},
		{
			Value:  "en US",
			Errors: 0,
		},
		{
			Value:  "English_United States.1252",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validateCollation()(tc.Value, "collation")
		if len(errors) != tc.Errors {
			t.Fatalf("Expected validateCollation to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}

func TestValidateAzureVirtualMachineTimeZone(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "",
			Errors: 0,
		},
		{
			Value:  "UTC",
			Errors: 0,
		},
		{
			Value:  "China Standard Time",
			Errors: 0,
		},
		{
			// Valid UTC time zone
			Value:  "utc-11",
			Errors: 0,
		},
		{
			// Invalid UTC time zone
			Value:  "UTC-30",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		_, errors := validateAzureVirtualMachineTimeZone()(tc.Value, "unittest")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected validateAzureVMTimeZone to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}

func TestValidateAzureDataLakeStoreRemoteFilePath(t *testing.T) {
	cases := []struct {
		Value  string
		Errors int
	}{
		{
			Value:  "bad",
			Errors: 1,
		},
		{
			Value:  "/good/file/path",
			Errors: 0,
		},
	}

	for _, tc := range cases {
		_, errors := validateFilePath()(tc.Value, "unittest")

		if len(errors) != tc.Errors {
			t.Fatalf("Expected validateFilePath to trigger '%d' errors for '%s' - got '%d'", tc.Errors, tc.Value, len(errors))
		}
	}
}
