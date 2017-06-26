package azurerm

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/terraform"
)

func resourceAzureRMRedisCacheMigrateState(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Redis Cache State v0; migrating to v1")
		return migrateAzureRMRedisCacheStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateAzureRMRedisCacheStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Redis Cache Attributes before Migration: %#v", is.Attributes)

	is.Attributes["redis_configuration.0.rdb_backup_enabled"] = "false"

	log.Printf("[DEBUG] ARM Redis Cache Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
