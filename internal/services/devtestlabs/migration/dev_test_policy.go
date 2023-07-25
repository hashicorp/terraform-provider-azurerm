// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/policies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = DevTestLabPolicyUpgradeV0ToV1{}

type DevTestLabPolicyUpgradeV0ToV1 struct{}

func (DevTestLabPolicyUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return devTestLabPolicySchemaForV0AndV1()
}

func (DevTestLabPolicyUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.devtestlab/labs/{labName}/policysets/{policySetName}/policies/{policyName}
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}/policySets/{policySetName}/policies/{policyName}
		oldId := rawState["id"].(string)
		id, err := policies.ParsePolicyIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}

func devTestLabPolicySchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"policy_set_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"lab_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"threshold": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"evaluator_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"fact_data": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": tags.Schema(),
	}
}
