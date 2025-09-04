// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/labs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = DevTestLabUpgradeV0ToV1{}

type DevTestLabUpgradeV0ToV1 struct{}

func (DevTestLabUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return devTestLabSchemaForV0AndV1()
}

func (DevTestLabUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.devtestlab/labs/{labName}
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}
		oldId := rawState["id"].(string)
		id, err := labs.ParseLabIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}

func devTestLabSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": commonschema.Location(),

		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"storage_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": tags.Schema(),

		"artifacts_storage_account_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"default_storage_account_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"default_premium_storage_account_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"key_vault_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"premium_data_disk_storage_account_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"unique_identifier": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
