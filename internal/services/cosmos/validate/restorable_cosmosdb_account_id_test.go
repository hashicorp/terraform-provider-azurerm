package validate

import "testing"

func TestRestorableCosmosdbAccountID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// empty
			Input: "",
			Valid: false,
		},
		{
			// missing SubscriptionId
			Input: "/",
			Valid: false,
		},
		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Valid: false,
		},
		{
			// missing LocationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/",
			Valid: false,
		},
		{
			// missing value for LocationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/",
			Valid: false,
		},
		{
			// missing RestorableDatabaseAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/location1/",
			Valid: false,
		},
		{
			// missing value for RestorableDatabaseAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/location1/restorableDatabaseAccounts/",
			Valid: false,
		},
		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/location1/restorableDatabaseAccounts/restorableDatabaseAccount1",
			Valid: true,
		},
		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.DOCUMENTDB/LOCATIONS/LOCATION1/RESTORABLEDATABASEACCOUNTS/RESTORABLEDATABASEACCOUNT1",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := RestorableCosmosdbAccountID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
