// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/servicebus/2022-10-01-preview/namespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = NamespaceNetworkRuleSetV0ToV1{}

type NamespaceNetworkRuleSetV0ToV1 struct{}

func (NamespaceNetworkRuleSetV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"namespace_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"default_action": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"ip_rules": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"trusted_services_allowed": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"network_rules": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"ignore_missing_vnet_service_endpoint": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},
	}
}

func (NamespaceNetworkRuleSetV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// removing the constant URI suffix from the id since it isn't needed
		oldId := strings.TrimSuffix(rawState["id"].(string), "/networkrulesets/default")

		id, err := namespaces.ParseNamespaceID(oldId)
		if err != nil {
			return nil, err
		}

		rawState["id"] = id.ID()

		return rawState, nil
	}
}
