package servicebus_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	`github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus`
)

func TestAccAzureRMServiceBusNamespaceMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
		ID           string
		Attributes   map[string]string
		Expected     map[string]string
	}{
		"v0_1_without_value": {
			StateVersion: 0,
			ID:           "some_id",
			Attributes:   map[string]string{},
			Expected:     map[string]string{},
		},
		"v0_1_basic_sku_without_value": {
			StateVersion: 0,
			ID:           "some_id",
			Attributes: map[string]string{
				"sku": "basic",
			},
			Expected: map[string]string{
				"sku": "basic",
			},
		},
		"v0_1_basic_sku_with_value": {
			StateVersion: 0,
			ID:           "some_id",
			Attributes: map[string]string{
				"sku":      "basic",
				"capacity": "3",
			},
			Expected: map[string]string{
				"sku": "basic",
			},
		},
		"v0_1_standard_sku_without_value": {
			StateVersion: 0,
			ID:           "some_id",
			Attributes: map[string]string{
				"sku": "standard",
			},
			Expected: map[string]string{
				"sku": "standard",
			},
		},
		"v0_1_standard_sku_with_value": {
			StateVersion: 0,
			ID:           "some_id",
			Attributes: map[string]string{
				"sku":      "standard",
				"capacity": "3",
			},
			Expected: map[string]string{
				"sku": "standard",
			},
		},
		"v0_1_premium_sku_without_value": {
			StateVersion: 0,
			ID:           "some_id",
			Attributes: map[string]string{
				"sku": "premium",
			},
			Expected: map[string]string{
				"sku": "premium",
			},
		},
		"v0_1_premium_sku_with_value": {
			StateVersion: 0,
			ID:           "some_id",
			Attributes: map[string]string{
				"sku":      "premium",
				"capacity": "3",
			},
			Expected: map[string]string{
				"sku":      "premium",
				"capacity": "3",
			},
		},
		"v0_1_premium_sku_with_value_casing": {
			StateVersion: 0,
			ID:           "some_id",
			Attributes: map[string]string{
				"sku":      "Premium",
				"capacity": "3",
			},
			Expected: map[string]string{
				"sku":      "Premium",
				"capacity": "3",
			},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.Attributes,
		}
		is, err := servicebus.ResourceAzureRMServiceBusNamespaceMigrateState(tc.StateVersion, is, nil)
		if err != nil {
			t.Fatalf("bad: %q, err: %#v", tn, err)
		}

		if !reflect.DeepEqual(tc.Expected, is.Attributes) {
			t.Fatalf("Bad Service Bus Namespace Migrate\n\n. Got: %+v\n\n expected: %+v", is.Attributes, tc.Expected)
		}
	}
}
