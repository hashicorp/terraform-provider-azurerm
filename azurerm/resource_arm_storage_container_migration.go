package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func resourceStorageContainerMigrateState(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Storage Container State v0; migrating to v1")
		return migrateStorageContainerStateV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateStorageContainerStateV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Storage Container Attributes before Migration: %#v", is.Attributes)

	environment := meta.(*ArmClient).Account.Environment

	containerName := is.Attributes["name"]
	storageAccountName := is.Attributes["storage_account_name"]
	newID := fmt.Sprintf("https://%s.blob.%s/%s", storageAccountName, environment.StorageEndpointSuffix, containerName)
	is.Attributes["id"] = newID
	is.ID = newID

	log.Printf("[DEBUG] ARM Storage Container Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
