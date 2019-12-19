package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func resourceStorageBlobMigrateState(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Storage Blob State v0; migrating to v1")
		return migrateStorageBlobStateV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateStorageBlobStateV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Storage Blob Attributes before Migration: %#v", is.Attributes)

	environment := meta.(*clients.Client).Account.Environment

	blobName := is.Attributes["name"]
	containerName := is.Attributes["storage_container_name"]
	storageAccountName := is.Attributes["storage_account_name"]
	newID := fmt.Sprintf("https://%s.blob.%s/%s/%s", storageAccountName, environment.StorageEndpointSuffix, containerName, blobName)
	is.Attributes["id"] = newID
	is.ID = newID

	log.Printf("[DEBUG] ARM Storage Blob Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
