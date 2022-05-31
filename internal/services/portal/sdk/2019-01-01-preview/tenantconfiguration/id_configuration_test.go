package tenantconfiguration

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ConfigurationId{}

func TestNewConfigurationID(t *testing.T) {
	id := NewConfigurationID("default")

	if id.ConfigurationName != "default" {
		t.Fatalf("Expected %q but got %q for Segment 'ConfigurationName'", id.ConfigurationName, "default")
	}
}

func TestFormatConfigurationID(t *testing.T) {
	actual := NewConfigurationID("default").ID()
	expected := "/providers/Microsoft.Portal/tenantConfigurations/default"
	if actual != expected {
		t.Fatalf("Expected the Formatted ID to be %q but got %q", expected, actual)
	}
}

func TestParseConfigurationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ConfigurationId
	}{
		{
			// Incomplete URI
			Input: "",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/providers",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/providers/Microsoft.Portal",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/providers/Microsoft.Portal/tenantConfigurations",
			Error: true,
		},
		{
			// Valid URI
			Input: "/providers/Microsoft.Portal/tenantConfigurations/default",
			Expected: &ConfigurationId{
				ConfigurationName: "default",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/providers/Microsoft.Portal/tenantConfigurations/default/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseConfigurationID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %+v", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.ConfigurationName != v.Expected.ConfigurationName {
			t.Fatalf("Expected %q but got %q for ConfigurationName", v.Expected.ConfigurationName, actual.ConfigurationName)
		}

	}
}

func TestParseConfigurationIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ConfigurationId
	}{
		{
			// Incomplete URI
			Input: "",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/providers",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/pRoViDeRs",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/providers/Microsoft.Portal",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/pRoViDeRs/mIcRoSoFt.pOrTaL",
			Error: true,
		},
		{
			// Incomplete URI
			Input: "/providers/Microsoft.Portal/tenantConfigurations",
			Error: true,
		},
		{
			// Incomplete URI (mIxEd CaSe since this is insensitive)
			Input: "/pRoViDeRs/mIcRoSoFt.pOrTaL/tEnAnTcOnFiGuRaTiOnS",
			Error: true,
		},
		{
			// Valid URI
			Input: "/providers/Microsoft.Portal/tenantConfigurations/default",
			Expected: &ConfigurationId{
				ConfigurationName: "default",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment)
			Input: "/providers/Microsoft.Portal/tenantConfigurations/default/extra",
			Error: true,
		},
		{
			// Valid URI (mIxEd CaSe since this is insensitive)
			Input: "/pRoViDeRs/mIcRoSoFt.pOrTaL/tEnAnTcOnFiGuRaTiOnS/dEfAuLt",
			Expected: &ConfigurationId{
				ConfigurationName: "default",
			},
		},
		{
			// Invalid (Valid Uri with Extra segment - mIxEd CaSe since this is insensitive)
			Input: "/pRoViDeRs/mIcRoSoFt.pOrTaL/tEnAnTcOnFiGuRaTiOnS/dEfAuLt/extra",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseConfigurationIDInsensitively(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %+v", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.ConfigurationName != v.Expected.ConfigurationName {
			t.Fatalf("Expected %q but got %q for ConfigurationName", v.Expected.ConfigurationName, actual.ConfigurationName)
		}

	}
}

func TestSegmentsForConfigurationId(t *testing.T) {
	segments := ConfigurationId{}.Segments()
	if len(segments) == 0 {
		t.Fatalf("ConfigurationId has no segments")
	}

	uniqueNames := make(map[string]struct{}, 0)
	for _, segment := range segments {
		uniqueNames[segment.Name] = struct{}{}
	}
	if len(uniqueNames) != len(segments) {
		t.Fatalf("Expected the Segments to be unique but got %q unique segments and %d total segments", len(uniqueNames), len(segments))
	}
}
