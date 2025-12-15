// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policydefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PolicyDefinitionV0ToV1 struct{}

var _ pluginsdk.StateUpgrade = PolicyDefinitionV0ToV1{}

func (p PolicyDefinitionV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"policy_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"management_group_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"policy_rule": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"parameters": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"role_definition_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"metadata": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},
	}
}

func (p PolicyDefinitionV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		parsedId, err := policydefinitions.ParseProviderPolicyDefinitionIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		newId := parsedId.ID()
		log.Printf("[DEBUG] Updating ID from `%s` to `%s`", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
