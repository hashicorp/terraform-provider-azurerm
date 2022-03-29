package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = MedTechServiceFhirDestinationId{}

func TestMedTechServiceFhirDestinationIDFormatter(t *testing.T) {
	actual := NewMedTechServiceFhirDestinationID("12345678-1234-9876-4563-123456789012", "group1", "workspace1", "iotconnector1", "destination1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotconnectors/iotconnector1/fhirdestinations/destination1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestMedTechServiceFhirDestinationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *MedTechServiceFhirDestinationId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing WorkspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/",
			Error: true,
		},

		{
			// missing value for WorkspaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/",
			Error: true,
		},

		{
			// missing IotconnectorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/",
			Error: true,
		},

		{
			// missing value for IotconnectorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotconnectors/",
			Error: true,
		},

		{
			// missing FhirdestinationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotconnectors/iotconnector1/",
			Error: true,
		},

		{
			// missing value for FhirdestinationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotconnectors/iotconnector1/fhirdestinations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotconnectors/iotconnector1/fhirdestinations/destination1",
			Expected: &MedTechServiceFhirDestinationId{
				SubscriptionId:      "12345678-1234-9876-4563-123456789012",
				ResourceGroup:       "group1",
				WorkspaceName:       "workspace1",
				IotconnectorName:    "iotconnector1",
				FhirdestinationName: "destination1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.HEALTHCAREAPIS/WORKSPACES/WORKSPACE1/IOTCONNECTORS/IOTCONNECTOR1/FHIRDESTINATIONS/DESTINATION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := MedTechServiceFhirDestinationID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.WorkspaceName != v.Expected.WorkspaceName {
			t.Fatalf("Expected %q but got %q for WorkspaceName", v.Expected.WorkspaceName, actual.WorkspaceName)
		}
		if actual.IotconnectorName != v.Expected.IotconnectorName {
			t.Fatalf("Expected %q but got %q for IotconnectorName", v.Expected.IotconnectorName, actual.IotconnectorName)
		}
		if actual.FhirdestinationName != v.Expected.FhirdestinationName {
			t.Fatalf("Expected %q but got %q for FhirdestinationName", v.Expected.FhirdestinationName, actual.FhirdestinationName)
		}
	}
}
