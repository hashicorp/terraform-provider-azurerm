// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2016-03-01/logprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = LogProfileUpgradeV0ToV1{}

type LogProfileUpgradeV0ToV1 struct{}

func (LogProfileUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return logProfileSchemaForV0AndV1()
}

func (LogProfileUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		// 	/subscriptions/{subscriptionId}/providers/microsoft.insights/logprofiles/profile1
		// new:
		// 	/subscriptions/{subscriptionId}/providers/Microsoft.Insights/logProfiles/profile1
		oldId := rawState["id"].(string)
		oldIdComponents := strings.Split(oldId, "/")

		if len(oldIdComponents) == 0 {
			return rawState, fmt.Errorf("old log profile id is empty or not formatted correctly: %s", oldId)
		}

		if len(oldIdComponents) != 7 {
			return rawState, fmt.Errorf("log profile id should have 6 segments, got %d: %s", len(oldIdComponents)-1, oldId)
		}

		newId := logprofiles.NewLogProfileID(oldIdComponents[2], oldIdComponents[6])

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId.ID())

		rawState["id"] = newId.ID()

		return rawState, nil
	}
}

func logProfileSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"storage_account_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"servicebus_rule_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"locations": {
			Type:     pluginsdk.TypeSet,
			MinItems: 1,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Set: pluginsdk.HashString,
		},
		"categories": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Set: pluginsdk.HashString,
		},
		"retention_policy": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"days": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},
	}
}
