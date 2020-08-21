package network

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func ResourceNetworkPacketCaptureMigrateState(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Network Packet Capture State v0; migrating to v1")
		return resourceNetworkPacketCaptureStateV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func resourceNetworkPacketCaptureStateV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Network Packet Capture before Migration: %#v", is.Attributes)

	newID := strings.Replace(is.Attributes["id"], "/NetworkPacketCaptures/", "/packetCaptures/", 1)
	is.Attributes["id"] = newID
	is.ID = newID

	log.Printf("[DEBUG] ARM Network Packet Capture Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
