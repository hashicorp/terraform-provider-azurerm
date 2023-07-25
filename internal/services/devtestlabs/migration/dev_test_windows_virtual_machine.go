// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = DevTestWindowsVirtualMachineUpgradeV0ToV1{}

type DevTestWindowsVirtualMachineUpgradeV0ToV1 struct{}

func (DevTestWindowsVirtualMachineUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return devTestWindowsVirtualMachineSchemaForV0AndV1()
}

func (DevTestWindowsVirtualMachineUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.devtestlab/labs/{labName}/virtualmachines/{virtualMachineName}
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/virtualMachines/{virtualMachineName}
		oldId := rawState["id"].(string)
		id, err := virtualmachines.ParseVirtualMachineIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}

func devTestWindowsVirtualMachineSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"lab_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"location": commonschema.Location(),

		"size": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"username": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"password": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"storage_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"lab_subnet_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"lab_virtual_network_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"allow_claim": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"disallow_public_ip_address": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"gallery_image_reference": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"offer": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"publisher": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"sku": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"version": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"inbound_nat_rule": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"protocol": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"backend_port": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"frontend_port": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"notes": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": tags.Schema(),

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"unique_identifier": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
