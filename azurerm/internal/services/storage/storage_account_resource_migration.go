package storage

import (
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func ResourceStorageAccountMigrateState(
	v int, is *terraform.InstanceState, _ interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Storage Account State v0; migrating to v1")
		return migrateStorageAccountStateV0toV1(is)
	case 1:
		log.Println("[INFO] Found AzureRM Storage Account State v1; migrating to v2")
		return migrateStorageAccountStateV1toV2(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateStorageAccountStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Storage Account Attributes before Migration: %#v", is.Attributes)

	accountType := is.Attributes["account_type"]
	split := strings.Split(accountType, "_")
	is.Attributes["account_tier"] = split[0]
	is.Attributes["account_replication_type"] = split[1]

	log.Printf("[DEBUG] ARM Storage Account Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}

func migrateStorageAccountStateV1toV2(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Storage Account Attributes before Migration: %#v", is.Attributes)

	is.Attributes["account_encryption_source"] = string(storage.MicrosoftStorage)

	log.Printf("[DEBUG] ARM Storage Account Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
