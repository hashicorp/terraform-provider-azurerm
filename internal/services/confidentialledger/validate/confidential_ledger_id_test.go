package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestConfidentialLedgerID(t *testing.T) {
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
			// missing ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Valid: false,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Valid: false,
		},

		{
			// missing LedgerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/ledgerRG/providers/Microsoft.ConfidentialLedger/",
			Valid: false,
		},

		{
			// missing value for LedgerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/ledgerRG/providers/Microsoft.ConfidentialLedger/Ledgers/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/ledgerRG/providers/Microsoft.ConfidentialLedger/Ledgers/testLedger",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/LEDGERRG/PROVIDERS/MICROSOFT.CONFIDENTIALLEDGER/LEDGERS/TESTLEDGER",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ConfidentialLedgerID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
