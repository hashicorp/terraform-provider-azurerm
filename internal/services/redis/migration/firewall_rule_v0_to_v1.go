// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2023-04-01/firewallrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = FirewallRuleV0ToV1{}

type FirewallRuleV0ToV1 struct{}

func (FirewallRuleV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"redis_cache_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"start_ip": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"end_ip": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (FirewallRuleV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := firewallrules.ParseFirewallRuleIDInsensitively(oldIdRaw)
		if err != nil {
			return rawState, err
		}

		newId := oldId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
