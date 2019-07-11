package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/terraform"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/file/shares"
)

func resourceStorageShareMigrateState(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Storage Share State v0; migrating to v1")
		return migrateStorageShareStateV0toV1(is)
	case 1:
		log.Println("[INFO] Found AzureRM Storage Share State v1; migrating to v2")
		return migrateStorageShareStateV1toV2(is, meta)
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

func migrateStorageShareStateV1toV2(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Storage Share Attributes before Migration: %#v", is.Attributes)

	environment := meta.(*ArmClient).environment
	client := shares.NewWithEnvironment(environment)

	shareName := is.Attributes["name"]
	storageAccountName := is.Attributes["storage_account_name"]
	newID := client.GetResourceID(storageAccountName, shareName)
	is.Attributes["id"] = newID
	is.ID = newID

	log.Printf("[DEBUG] ARM Storage Share Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
