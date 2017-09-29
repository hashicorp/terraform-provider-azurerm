package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/terraform"
)

func resourceStorageAccountMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Storage Account State v0; migrating to v1")
		return migrateStorageAccountStateV0toV1(is)
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
