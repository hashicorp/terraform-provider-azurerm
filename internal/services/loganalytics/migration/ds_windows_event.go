// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
)

var _ pluginsdk.StateUpgrade = WindowsEventV0ToV1{}

type WindowsEventV0ToV1 struct{}

func (WindowsEventV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId, err := datasources.ParseDataSourceIDInsensitively(rawState["id"].(string))
		if err != nil {
			return rawState, err
		}

		rawState["id"] = oldId.ID()
		return rawState, nil
	}
}

func (WindowsEventV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"workspace_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"event_log_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"event_types": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			Set:      set.HashStringIgnoreCase,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
