package validate

import "testing"

func TestIsAutomationAccountID(t *testing.T) {
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
			// invalid resource id
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Automation/notAnAutomationAccount/Automation1",
			Valid: false,
		},

		{
			// text
			Input: "justSomeText",
			Valid: false,
		},

		{
			// valid automation resource id
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Automation/automationAccounts/Automation1",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		valid := IsAutomationAccountID(tc.Input)

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
