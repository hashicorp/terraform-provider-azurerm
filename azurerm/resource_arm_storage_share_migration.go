package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/terraform"
)

func resourceStorageShareMigrateState(
	v int, is *terraform.InstanceState, _ interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Storage Share State v0; migrating to v1")
		return migrateStorageShareStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateStorageShareStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Storage Share Attributes before Migration: %#v", is.Attributes)

	name := is.Attributes["name"]
	resourceGroupName := is.Attributes["resource_group_name"]
	storageAccountName := is.Attributes["storage_account_name"]
	newID := fmt.Sprintf("%s/%s/%s", name, resourceGroupName, storageAccountName)
	is.Attributes["id"] = newID
	is.ID = newID

	log.Printf("[DEBUG] ARM Storage Share Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
