package servicebus

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func ResourceAzureRMServiceBusNamespaceMigrateState(
	v int, is *terraform.InstanceState, _ interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM ServiceBus Namespace State v0; migrating to v1")
		return migrateAzureRMServiceBusNamespaceStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateAzureRMServiceBusNamespaceStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM ServiceBus Namespace Attributes before Migration: %#v", is.Attributes)

	skuName := strings.ToLower(is.Attributes["sku"])
	premiumSku := strings.ToLower(string(servicebus.Premium))

	if skuName != premiumSku {
		delete(is.Attributes, "capacity")
	}

	log.Printf("[DEBUG] ARM ServiceBus Namespace Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
