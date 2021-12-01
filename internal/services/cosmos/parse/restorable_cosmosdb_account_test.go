package parse

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
)

var _ resourceid.Formatter = RestorableCosmosdbAccountId{}

func TestRestorableCosmosdbAccountIDFormatter(t *testing.T) {
	actual := NewRestorableCosmosdbAccountID("12345678-1234-9876-4563-123456789012", "location1", "restorableDatabaseAccount1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/location1/restorableDatabaseAccounts/restorableDatabaseAccount1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestRestorableCosmosdbAccountID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *RestorableCosmosdbAccountId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/",
			Error: true,
		},
		{
			// missing value for LocationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/",
			Error: true,
		},
		{
			// missing RestorableDatabaseAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/location1/",
			Error: true,
		},
		{
			// missing value for RestorableDatabaseAccountName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/location1/restorableDatabaseAccounts/",
			Error: true,
		},
		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.DocumentDB/locations/location1/restorableDatabaseAccounts/restorableDatabaseAccount1",
			Expected: &RestorableCosmosdbAccountId{
				SubscriptionId:                "12345678-1234-9876-4563-123456789012",
				LocationName:                  "location1",
				RestorableDatabaseAccountName: "restorableDatabaseAccount1",
			},
		},
		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.DOCUMENTDB/LOCATIONS/LOCATION1/RESTORABLEDATABASEACCOUNTS/RESTORABLEDATABASEACCOUNT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := RestorableCosmosdbAccountID(v.Input)
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
		if actual.RestorableDatabaseAccountName != v.Expected.RestorableDatabaseAccountName {
			t.Fatalf("Expected %q but got %q for RestorableDatabaseAccountName", v.Expected.RestorableDatabaseAccountName, actual.RestorableDatabaseAccountName)
		}
	}
}
