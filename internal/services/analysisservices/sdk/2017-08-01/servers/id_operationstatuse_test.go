package servers

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = OperationstatuseId{}

func TestOperationstatuseIDFormatter(t *testing.T) {
	actual := NewOperationstatuseID("{subscriptionId}", "{location}", "{operationId}").ID()
	expected := "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/locations/{location}/operationstatuses/{operationId}"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParseOperationstatuseID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *OperationstatuseId
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
			// missing LocationName
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/",
			Error: true,
		},

		{
			// missing value for LocationName
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/locations/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/locations/{location}/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/locations/{location}/operationstatuses/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/locations/{location}/operationstatuses/{operationId}",
			Expected: &OperationstatuseId{
				SubscriptionId: "{subscriptionId}",
				LocationName:   "{location}",
				Name:           "{operationId}",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/{SUBSCRIPTIONID}/PROVIDERS/MICROSOFT.ANALYSISSERVICES/LOCATIONS/{LOCATION}/OPERATIONSTATUSES/{OPERATIONID}",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseOperationstatuseID(v.Input)
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
		if actual.LocationName != v.Expected.LocationName {
			t.Fatalf("Expected %q but got %q for LocationName", v.Expected.LocationName, actual.LocationName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}

func TestParseOperationstatuseIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *OperationstatuseId
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
			// missing LocationName
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/",
			Error: true,
		},

		{
			// missing value for LocationName
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/locations/",
			Error: true,
		},

		{
			// missing Name
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/locations/{location}/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/locations/{location}/operationstatuses/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/locations/{location}/operationstatuses/{operationId}",
			Expected: &OperationstatuseId{
				SubscriptionId: "{subscriptionId}",
				LocationName:   "{location}",
				Name:           "{operationId}",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/locations/{location}/operationstatuses/{operationId}",
			Expected: &OperationstatuseId{
				SubscriptionId: "{subscriptionId}",
				LocationName:   "{location}",
				Name:           "{operationId}",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/LOCATIONS/{location}/OPERATIONSTATUSES/{operationId}",
			Expected: &OperationstatuseId{
				SubscriptionId: "{subscriptionId}",
				LocationName:   "{location}",
				Name:           "{operationId}",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.AnalysisServices/LoCaTiOnS/{location}/OpErAtIoNsTaTuSeS/{operationId}",
			Expected: &OperationstatuseId{
				SubscriptionId: "{subscriptionId}",
				LocationName:   "{location}",
				Name:           "{operationId}",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseOperationstatuseIDInsensitively(v.Input)
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
		if actual.LocationName != v.Expected.LocationName {
			t.Fatalf("Expected %q but got %q for LocationName", v.Expected.LocationName, actual.LocationName)
		}
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
