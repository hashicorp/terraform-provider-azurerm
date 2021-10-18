package signalr

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SubscriptionId{}

func TestSubscriptionIDFormatter(t *testing.T) {
	actual := NewSubscriptionID("{subscriptionId}").ID()
	expected := "/subscriptions/{subscriptionId}"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParseSubscriptionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SubscriptionId
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
			// valid
			Input: "/subscriptions/{subscriptionId}",
			Expected: &SubscriptionId{
				SubscriptionId: "{subscriptionId}",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/{SUBSCRIPTIONID}",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseSubscriptionID(v.Input)
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
	}
}

func TestParseSubscriptionIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SubscriptionId
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
			// valid
			Input: "/subscriptions/{subscriptionId}",
			Expected: &SubscriptionId{
				SubscriptionId: "{subscriptionId}",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/{subscriptionId}",
			Expected: &SubscriptionId{
				SubscriptionId: "{subscriptionId}",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/{subscriptionId}",
			Expected: &SubscriptionId{
				SubscriptionId: "{subscriptionId}",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/{subscriptionId}",
			Expected: &SubscriptionId{
				SubscriptionId: "{subscriptionId}",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ParseSubscriptionIDInsensitively(v.Input)
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
	}
}
