package securitycenter_test

import (
	"testing"

	securityCenter "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter"
)

func TestSecurityCenterSubscriptionPricingMigrateState(t *testing.T) {
	inputAttributes := map[string]interface{}{
		"id": "/subscriptions/00000000-0000-0000-0000-000000000000/pricings/default",
	}
	expectedId := "/subscriptions/00000000-0000-0000-0000-000000000000/pricings/VirtualMachines"

	rawState, _ := securityCenter.ResourceSecurityCenterSubscriptionPricingUpgradeV0ToV1(inputAttributes, nil)
	if rawState["id"].(string) != expectedId {
		t.Fatalf("ResourceType migration failed, expected %q, got: %q", expectedId, rawState["id"].(string))
	}
}
