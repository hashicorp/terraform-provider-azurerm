// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/customdomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = CdnEndpointCustomDomainV0ToV1{}

type CdnEndpointCustomDomainV0ToV1 struct{}

func (CdnEndpointCustomDomainV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"cdn_endpoint_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"host_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"cdn_managed_https": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"certificate_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"protocol_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"tls_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"user_managed_https": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tls_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"key_vault_secret_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}

func (CdnEndpointCustomDomainV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// IDs imported while this resource parsed them with the legacy resourceids.ParseAzureResourceID
		// can contain non-canonically cased static segments (e.g. `resourcegroups`, `microsoft.cdn`),
		// which the case-sensitive SDK parser rejects - normalise them to the canonical casing
		oldId := rawState["id"].(string)
		id, err := customdomains.ParseEndpointCustomDomainIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from `%s` to `%s`", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
