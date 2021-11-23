package migration

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/parse"
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
		oldId, err := azure.ParseAzureResourceID(strings.Replace(rawState["id"].(string), "/policysets/", "/policySets/", 1))
		if err != nil {
			return rawState, err
		}

		labName, err := oldId.PopSegment("labs")
		if err != nil {
			return rawState, err
		}

		policySet, err := oldId.PopSegment("policySets")
		if err != nil {
			return rawState, err
		}

		policyName, err := oldId.PopSegment("policies")
		if err != nil {
			return rawState, err
		}

		newId := parse.NewDevTestLabPolicyID(oldId.SubscriptionID, oldId.ResourceGroup, labName, policySet, policyName)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId.ID())

		rawState["id"] = newId.ID()

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
