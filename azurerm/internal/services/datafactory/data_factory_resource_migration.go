package datafactory

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func ResourceDataFactoryMigrateState(
	v int, is *terraform.InstanceState, _ interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM DataFactory State v0; migrating to v1")
		return migrateDataFactoryStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateDataFactoryStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Data Factory Attributes before Migration: %#v", is.Attributes)

	is.Attributes["public_network_enabled"] = string(datafactory.PublicNetworkAccessEnabled)

	log.Printf("[DEBUG] ARM Data Factory Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
