// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = DataFactoryV0ToV1{}

type DataFactoryV0ToV1 struct{}

func (DataFactoryV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"github_configuration": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"vsts_configuration"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"account_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"branch_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"git_url": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"repository_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"root_folder": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"vsts_configuration": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"github_configuration"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"account_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"branch_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"project_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"repository_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"root_folder": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (DataFactoryV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[DEBUG] Updating `public_network_enabled` to %q", factories.PublicNetworkAccessEnabled)

		rawState["public_network_enabled"] = true

		return rawState, nil
	}
}
