package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = MonitorAADDiagnosticSettingId{}

func TestMonitorAADDiagnosticSettingIDFormatter(t *testing.T) {
	actual := NewMonitorAADDiagnosticSettingID("setting1").ID()
	expected := "/providers/Microsoft.AADIAM/diagnosticSettings/setting1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestMonitorAADDiagnosticSettingID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *MonitorAADDiagnosticSettingId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing prefix
			Input: "/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/providers/Microsoft.AADIAM/diagnosticSettings/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.AADIAM/diagnosticSettings/setting1",
			Expected: &MonitorAADDiagnosticSettingId{
				Name: "setting1",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.AADIAM/DIAGNOSTICSETTINGS/SETTING1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := MonitorAADDiagnosticSettingID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
