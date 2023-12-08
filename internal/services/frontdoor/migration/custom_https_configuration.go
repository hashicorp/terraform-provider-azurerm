// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = CustomHttpsConfigurationV0ToV1{}

type CustomHttpsConfigurationV0ToV1 struct{}

func (CustomHttpsConfigurationV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"frontend_endpoint_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"custom_https_provisioning_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		//lintignore:XS003
		"custom_https_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"certificate_source": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"azure_key_vault_certificate_secret_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"azure_key_vault_certificate_secret_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"azure_key_vault_certificate_vault_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"minimum_tls_version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"provisioning_state": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"provisioning_substate": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (CustomHttpsConfigurationV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// this was: fmt.Sprintf("%s/customHttpsConfiguration/%s", frontDoorId, frontendEndpointName
		oldId := rawState["id"].(string)
		id, err := parse.CustomHttpsConfigurationIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
