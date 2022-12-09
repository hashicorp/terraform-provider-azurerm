package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = MedTechServiceFhirDestinationId{}

func TestMedTechServiceFhirDestinationIDFormatter(t *testing.T) {
	actual := NewMedTechServiceFhirDestinationID("12345678-1234-9876-4563-123456789012", "group1", "workspace1", "iotconnector1", "destination1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/iotconnector1/fhirDestinations/destination1"
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
			// missing IotConnectorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/",
			Error: true,
		},

		{
			// missing value for IotConnectorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/",
			Error: true,
		},

		{
			// missing FhirDestinationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/iotconnector1/",
			Error: true,
		},

		{
			// missing value for FhirDestinationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/iotconnector1/fhirDestinations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/iotconnector1/fhirDestinations/destination1",
			Expected: &MedTechServiceFhirDestinationId{
				SubscriptionId:      "12345678-1234-9876-4563-123456789012",
				ResourceGroup:       "group1",
				WorkspaceName:       "workspace1",
				IotConnectorName:    "iotconnector1",
				FhirDestinationName: "destination1",
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
		if actual.IotConnectorName != v.Expected.IotConnectorName {
			t.Fatalf("Expected %q but got %q for IotConnectorName", v.Expected.IotConnectorName, actual.IotConnectorName)
		}
		if actual.FhirDestinationName != v.Expected.FhirDestinationName {
			t.Fatalf("Expected %q but got %q for FhirDestinationName", v.Expected.FhirDestinationName, actual.FhirDestinationName)
		}
	}
}

func TestMedTechServiceFhirDestinationIDInsensitively(t *testing.T) {
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
			// missing IotConnectorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/",
			Error: true,
		},

		{
			// missing value for IotConnectorName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/",
			Error: true,
		},

		{
			// missing FhirDestinationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/iotconnector1/",
			Error: true,
		},

		{
			// missing value for FhirDestinationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/iotconnector1/fhirDestinations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotConnectors/iotconnector1/fhirDestinations/destination1",
			Expected: &MedTechServiceFhirDestinationId{
				SubscriptionId:      "12345678-1234-9876-4563-123456789012",
				ResourceGroup:       "group1",
				WorkspaceName:       "workspace1",
				IotConnectorName:    "iotconnector1",
				FhirDestinationName: "destination1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1/iotconnectors/iotconnector1/fhirdestinations/destination1",
			Expected: &MedTechServiceFhirDestinationId{
				SubscriptionId:      "12345678-1234-9876-4563-123456789012",
				ResourceGroup:       "group1",
				WorkspaceName:       "workspace1",
				IotConnectorName:    "iotconnector1",
				FhirDestinationName: "destination1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/WORKSPACES/workspace1/IOTCONNECTORS/iotconnector1/FHIRDESTINATIONS/destination1",
			Expected: &MedTechServiceFhirDestinationId{
				SubscriptionId:      "12345678-1234-9876-4563-123456789012",
				ResourceGroup:       "group1",
				WorkspaceName:       "workspace1",
				IotConnectorName:    "iotconnector1",
				FhirDestinationName: "destination1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/WoRkSpAcEs/workspace1/IoTcOnNeCtOrS/iotconnector1/FhIrDeStInAtIoNs/destination1",
			Expected: &MedTechServiceFhirDestinationId{
				SubscriptionId:      "12345678-1234-9876-4563-123456789012",
				ResourceGroup:       "group1",
				WorkspaceName:       "workspace1",
				IotConnectorName:    "iotconnector1",
				FhirDestinationName: "destination1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := MedTechServiceFhirDestinationIDInsensitively(v.Input)
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
		if actual.IotConnectorName != v.Expected.IotConnectorName {
			t.Fatalf("Expected %q but got %q for IotConnectorName", v.Expected.IotConnectorName, actual.IotConnectorName)
		}
		if actual.FhirDestinationName != v.Expected.FhirDestinationName {
			t.Fatalf("Expected %q but got %q for FhirDestinationName", v.Expected.FhirDestinationName, actual.FhirDestinationName)
		}
	}
}
