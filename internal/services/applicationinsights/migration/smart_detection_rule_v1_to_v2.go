// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	smartdetection "github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/componentproactivedetectionapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = SmartDetectionRuleUpgradeV1ToV2{}

type SmartDetectionRuleUpgradeV1ToV2 struct{}

func (SmartDetectionRuleUpgradeV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return smartDetectionRuleSchemaForV1AndV2()
}

func (SmartDetectionRuleUpgradeV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/components/component1/SmartDetectionRule/rule1
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/component1/proactiveDetectionConfigs/rule1
		oldIdRaw := rawState["id"].(string)
		oldId, err := parse.SmartDetectionRuleIDInsensitively(oldIdRaw)
		if err != nil {
			return rawState, err
		}

		id := smartdetection.NewProactiveDetectionConfigID(oldId.SubscriptionId, oldId.ResourceGroup, oldId.ComponentName, oldId.SmartDetectionRuleName)

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}

func smartDetectionRuleSchemaForV1AndV2() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"application_insights_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"send_emails_to_subscription_owners": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"additional_email_recipients": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
	}
}
