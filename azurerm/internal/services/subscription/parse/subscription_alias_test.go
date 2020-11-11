package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = SubscriptionAliasId{}

func TestAvailabilitySetIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	actual := NewSubscriptionAliasId("alias1").ID(subscriptionId)
	expected := "/providers/Microsoft.Subscription/aliases/alias1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSubscriptionAliasID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *SubscriptionAliasId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "Missing Alias Value",
			Input:    "/providers/Microsoft.Subscription/aliases",
			Expected: nil,
		},
		{
			Name:  "subscription Alias ID",
			Input: "/providers/Microsoft.Subscription/aliases/alias1",
			Expected: &SubscriptionAliasId{
				Name: "alias1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/providers/Microsoft.Subscription/Aliases/alias1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := SubscriptionAliasID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
