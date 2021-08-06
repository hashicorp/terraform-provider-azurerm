package validate

import "testing"

func TestMonitorAADDiagnosticSettingID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// missing prefix
			Input: "/",
			Valid: false,
		},

		{
			// missing value for Name
			Input: "/providers/Microsoft.AADIAM/diagnosticSettings/",
			Valid: false,
		},

		{
			// valid
			Input: "/providers/Microsoft.AADIAM/diagnosticSettings/setting1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.AADIAM/DIAGNOSTICSETTINGS/SETTING1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := MonitorAADDiagnosticSettingID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
