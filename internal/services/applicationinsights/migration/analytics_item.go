// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

var _ pluginsdk.StateUpgrade = AnalyticsItemUpgradeV0ToV1{}

type AnalyticsItemUpgradeV0ToV1 struct{}

func (AnalyticsItemUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return analyticsItemSchemaForV0AndV1()
}

func (AnalyticsItemUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/components/component1/[myanalyticsItems|analyticsItems]/item1
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/component1/[myAnalyticsItems|analyticsItems]/item1
		oldId, err := azure.ParseAzureResourceID(rawState["id"].(string))
		if err != nil {
			return rawState, err
		}

		itemName := ""
		newId := ""
		for key, value := range oldId.Path {
			if strings.EqualFold(key, "analyticsitems") {
				itemName = value
				newId = parse.NewAnalyticsSharedItemID(oldId.SubscriptionID, oldId.ResourceGroup, oldId.Path["components"], itemName).ID()
				break
			} else if strings.EqualFold(key, "myanalyticsitems") {
				itemName = value
				newId = parse.NewAnalyticsUserItemID(oldId.SubscriptionID, oldId.ResourceGroup, oldId.Path["components"], itemName).ID()
				break
			}
		}

		if itemName == "" {
			return rawState, fmt.Errorf("couldn't find the `analyticsitems` or `myanalyticsitems` segment in the old resource id %q", oldId)
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId

		return rawState, nil
	}
}

func analyticsItemSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"application_insights_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"content": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"scope": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"function_alias": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"time_created": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_modified": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
