package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/terraform"
)

func resourceDataLakeStoreFileMigrateState(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Data Lake Store File State v0; migrating to v1")
		return resourceDataLakeStoreFileStateV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func resourceDataLakeStoreFileStateV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Data Lake Store File Attributes before Migration: %#v", is.Attributes)

	client := meta.(*ArmClient).datalake.StoreFilesClient

	storageAccountName := is.Attributes["account_name"]
	filePath := is.Attributes["remote_file_path"]
	newID := fmt.Sprintf("%s.%s%s", storageAccountName, client.AdlsFileSystemDNSSuffix, filePath)
	is.Attributes["id"] = newID
	is.ID = newID

	log.Printf("[DEBUG] ARM Data Lake Store File Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
