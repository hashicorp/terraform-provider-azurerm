package migration

import (
	"context"
	"testing"
)

func TestSecurityCenterSubscriptionPricingMigrateState(t *testing.T) {
	inputAttributes := map[string]interface{}{
		"id": "/subscriptions/00000000-0000-0000-0000-000000000000/pricings/default",
	}
	expectedId := "/subscriptions/00000000-0000-0000-0000-000000000000/pricings/VirtualMachines"

	rawState, _ := SubscriptionPricingV0ToV1{}.UpgradeFunc()(context.TODO(), inputAttributes, nil)
	if rawState["id"].(string) != expectedId {
		t.Fatalf("ResourceType migration failed, expected %q, got: %q", expectedId, rawState["id"].(string))
	}
}
