package network

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func ResourceNetworkConnectionMonitorMigrateState(v int, is *terraform.InstanceState, _ interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Network Connection Monitor State v0; migrating to v1")
		return resourceNetworkConnectionMonitorStateV0toV1(is)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func resourceNetworkConnectionMonitorStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Network Connection Monitor before Migration: %#v", is.Attributes)

	newID := strings.Replace(is.Attributes["id"], "/NetworkConnectionMonitors/", "/connectionMonitors/", 1)
	is.Attributes["id"] = newID
	is.ID = newID

	log.Printf("[DEBUG] ARM Network Connection Monitor Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
