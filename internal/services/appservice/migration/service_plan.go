// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ServicePlanV0toV1 struct{}

var _ pluginsdk.StateUpgrade = ServicePlanV0toV1{}

func (s ServicePlanV0toV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"os_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"app_service_environment_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"per_site_scaling_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"worker_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"maximum_elastic_worker_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"zone_balancing_enabled": {
			Type:     pluginsdk.TypeBool,
			ForceNew: true,
			Optional: true,
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

func (s ServicePlanV0toV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		parsedId, err := commonids.ParseAppServicePlanIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}
		newId := parsedId.ID()
		rawState["id"] = newId
		return rawState, nil
	}
}
