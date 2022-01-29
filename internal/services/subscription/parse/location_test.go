package parse

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = LocationId{}

func TestLocationIDFormatter(t *testing.T) {
	actual := NewLocationID("12345678-1234-9876-4563-123456789012").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/locations"

	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestLocationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *LocationId
	}{
		{
			Input: "",
			Error: true,
		},
		{
			Input: "/",
			Error: true,
		},
		{
			Input: "/subscriptions/",
			Error: true,
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},
		{
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/locations",
			Expected: &LocationId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
			},
		},
		{
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/LOCATIONS",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := LocationID(v.Input)
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
