// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2023-07-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = SubscriptionTemplateDeploymentV0ToV1{}

type SubscriptionTemplateDeploymentV0ToV1 struct{}

func (SubscriptionTemplateDeploymentV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
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

func (SubscriptionTemplateDeploymentV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// IDs imported while this resource parsed them with the legacy resourceids.ParseAzureResourceID
		// can contain a non-canonically cased provider segment (e.g. `microsoft.resources`),
		// which the case-sensitive SDK parser rejects - normalise it to the canonical casing
		oldId := rawState["id"].(string)
		id, err := deployments.ParseProviderDeploymentIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from `%s` to `%s`", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
