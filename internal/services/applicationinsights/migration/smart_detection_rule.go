package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = SmartDetectionRuleUpgradeV0ToV1{}

type SmartDetectionRuleUpgradeV0ToV1 struct{}

func (SmartDetectionRuleUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return smartDetectionRuleSchemaForV0AndV1()
}

func (SmartDetectionRuleUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/components/component1/SmartDetectionRule/rule1
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/component1/smartDetectionRule/rule1
		oldId, err := azure.ParseAzureResourceID(rawState["id"].(string))
		if err != nil {
			return rawState, err
		}

		ruleName := ""
		for key, value := range oldId.Path {
			if strings.EqualFold(key, "SmartDetectionRule") {
				ruleName = value
				break
			}
		}

		if ruleName == "" {
			return rawState, fmt.Errorf("couldn't find the `SmartDetectionRule` segment in the old resource id %q", oldId)
		}

		newId := parse.NewSmartDetectionRuleID(oldId.SubscriptionID, oldId.ResourceGroup, oldId.Path["components"], ruleName)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId.ID())

		rawState["id"] = newId.ID()

		return rawState, nil
	}
}

func smartDetectionRuleSchemaForV0AndV1() map[string]*pluginsdk.Schema {
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
