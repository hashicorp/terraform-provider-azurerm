package signalr

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = LocationId{}

func TestLocationIDFormatter(t *testing.T) {
	actual := NewLocationID("{subscriptionId}", "{location}").ID()
	expected := "/subscriptions/{subscriptionId}/providers/Microsoft.SignalRService/locations/{location}"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParseLocationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *LocationId
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
			// missing Name
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.SignalRService/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.SignalRService/locations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.SignalRService/locations/{location}",
			Expected: &LocationId{
				SubscriptionId: "{subscriptionId}",
				Name:           "{location}",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/{SUBSCRIPTIONID}/PROVIDERS/MICROSOFT.SIGNALRSERVICE/LOCATIONS/{LOCATION}",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseLocationID(v.Input)
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
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}

func TestParseLocationIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *LocationId
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
			// missing Name
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.SignalRService/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.SignalRService/locations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.SignalRService/locations/{location}",
			Expected: &LocationId{
				SubscriptionId: "{subscriptionId}",
				Name:           "{location}",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.SignalRService/locations/{location}",
			Expected: &LocationId{
				SubscriptionId: "{subscriptionId}",
				Name:           "{location}",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.SignalRService/LOCATIONS/{location}",
			Expected: &LocationId{
				SubscriptionId: "{subscriptionId}",
				Name:           "{location}",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/{subscriptionId}/providers/Microsoft.SignalRService/LoCaTiOnS/{location}",
			Expected: &LocationId{
				SubscriptionId: "{subscriptionId}",
				Name:           "{location}",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseLocationIDInsensitively(v.Input)
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
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
