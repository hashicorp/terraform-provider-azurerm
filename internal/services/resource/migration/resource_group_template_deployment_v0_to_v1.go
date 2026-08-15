// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ResourceGroupTemplateDeploymentV0ToV1{}

type ResourceGroupTemplateDeploymentV0ToV1 struct{}

func (ResourceGroupTemplateDeploymentV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"deployment_mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"template_content": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"template_spec_version_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"debug_level": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"parameters_content": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"output_content": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (ResourceGroupTemplateDeploymentV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// IDs imported while this resource parsed them with the legacy resourceids.ParseAzureResourceID
		// can contain non-canonically cased static segments (e.g. `resourcegroups`, `microsoft.resources`),
		// which the case-sensitive SDK parser rejects - normalise them to the canonical casing
		oldId := rawState["id"].(string)
		id, err := deployments.ParseResourceGroupProviderDeploymentIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from `%s` to `%s`", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
