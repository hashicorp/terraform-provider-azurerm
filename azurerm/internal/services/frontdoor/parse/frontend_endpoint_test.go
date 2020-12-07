package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = FrontendEndpointId{}

func TestFrontendEndpointIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewFrontendEndpointID(subscriptionId, "group1", "frontdoor1", "endpoint1").ID("")
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/frontDoors/frontdoor1/frontendEndpoints/endpoint1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestFrontendEndpointIDParser(t *testing.T) {
	testData := []struct {
		input    string
		expected *FrontendEndpointId
	}{
		{
			// lower case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/frontdoors/frontDoor1/frontendendpoints/endpoint1",
			expected: nil,
		},
		{
			// camel case
			input: "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/frontDoors/frontDoor1/frontendEndpoints/endpoint1",
			expected: &FrontendEndpointId{
				ResourceGroup: "group1",
				FrontDoorName: "frontDoor1",
				Name:          "endpoint1",
			},
		},
		{
			// title case
			input: "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/Frontdoors/frontDoor1/FrontendEndpoints/endpoint1",
			expected: &FrontendEndpointId{
				ResourceGroup: "group1",
				FrontDoorName: "frontDoor1",
				Name:          "endpoint1",
			},
		},
		{
			// pascal case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/FrontDoors/frontDoor1/Frontendendpoints/endpoint1",
			expected: nil,
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q..", test.input)
		actual, err := FrontendEndpointIDInsensitively(test.input)
		if err != nil && test.expected == nil {
			continue
		} else {
			if err == nil && test.expected == nil {
				t.Fatalf("Expected an error but didn't get one")
			} else if err != nil && test.expected != nil {
				t.Fatalf("Expected no error but got: %+v", err)
			}
		}

		if actual.ResourceGroup != test.expected.ResourceGroup {
			t.Fatalf("Expected ResourceGroup to be %q but was %q", test.expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.FrontDoorName != test.expected.FrontDoorName {
			t.Fatalf("Expected FrontDoorName to be %q but was %q", test.expected.FrontDoorName, actual.FrontDoorName)
		}

		if actual.Name != test.expected.Name {
			t.Fatalf("Expected name to be %q but was %q", test.expected.Name, actual.Name)
		}
	}
}

func TestFrontendEndpointIDForImportParser(t *testing.T) {
	testData := []struct {
		input    string
		expected *FrontendEndpointId
	}{
		{
			// lower case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/frontdoors/frontDoor1/frontendendpoints/endpoint1",
			expected: nil,
		},
		{
			// camel case
			input: "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/frontDoors/frontDoor1/frontendEndpoints/endpoint1",
			expected: &FrontendEndpointId{
				ResourceGroup: "group1",
				FrontDoorName: "frontDoor1",
				Name:          "endpoint1",
			},
		},
		{
			// title case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/Frontdoors/frontDoor1/FrontendEndpoints/endpoint1",
			expected: nil,
		},
		{
			// pascal case
			input:    "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Network/FrontDoors/frontDoor1/Frontendendpoints/endpoint1",
			expected: nil,
		},
	}
	for _, test := range testData {
		t.Logf("Testing %q..", test.input)
		actual, err := FrontendEndpointID(test.input)
		if err != nil && test.expected == nil {
			continue
		} else {
			if err == nil && test.expected == nil {
				t.Fatalf("Expected an error but didn't get one")
			} else if err != nil && test.expected != nil {
				t.Fatalf("Expected no error but got: %+v", err)
			}
		}

		if actual.ResourceGroup != test.expected.ResourceGroup {
			t.Fatalf("Expected ResourceGroup to be %q but was %q", test.expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.FrontDoorName != test.expected.FrontDoorName {
			t.Fatalf("Expected FrontDoorName to be %q but was %q", test.expected.FrontDoorName, actual.FrontDoorName)
		}

		if actual.Name != test.expected.Name {
			t.Fatalf("Expected name to be %q but was %q", test.expected.Name, actual.Name)
		}
	}
}
