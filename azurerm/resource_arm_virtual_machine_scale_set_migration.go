package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func resourceVirtualMachineScaleSetMigrateState(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Scale Set State v0; migrating to v1")
		return resourceVirtualMachineScaleSetStateV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func resourceVirtualMachineScaleSetStateV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Virtual Machine Scale Set Attributes before Migration: %#v", is.Attributes)

	client := meta.(*ArmClient).Compute.VMScaleSetClient
	ctx, cancel := context.WithTimeout(meta.(*ArmClient).StopContext, 5*time.Minute)
	defer cancel()

	resGroup := is.Attributes["resource_group_name"]
	name := is.Attributes["name"]

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return is, err
	}

	is.ID = *read.ID

	log.Printf("[DEBUG] ARM Virtual Machine Scale Set Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}
