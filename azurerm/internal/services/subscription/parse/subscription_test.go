package parse

import "testing"

func TestSubscriptionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SubscriptionId
	}{
		{
			// Empty
			Input: "",
			Error: true,
		},
		{
			// missing ID
			Input: "/subscriptions/",
			Error: true,
		},
		{
			// Invalid Resource ID UUID portion
			Input: "/subscriptions/abcdefgh-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			// Invalid UUID
			Input: "abcdefgh-0000-0000-0000-000000000000",
			Error: true,
		},
		{
			// Valid Resource ID
			Input: "/subscriptions/12345678-1234-5678-9012-123456789012",
			Expected: &SubscriptionId{
				SubscriptionID: "12345678-1234-5678-9012-123456789012",
			},
		},
		{
			// Valid UUID
			Input: "12345678-1234-5678-9012-123456789012",
			Expected: &SubscriptionId{
				SubscriptionID: "12345678-1234-5678-9012-123456789012",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SubscriptionID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("expected a value, but got an error: %+v", err)
		}

		if v.Error {
			t.Fatal("expected an error but did not get one")
		}

		if actual.SubscriptionID != v.Expected.SubscriptionID {
			t.Fatalf("expected %q but got %q for SubscriptionId", v.Expected.SubscriptionID, actual.SubscriptionID)
		}
	}
}
