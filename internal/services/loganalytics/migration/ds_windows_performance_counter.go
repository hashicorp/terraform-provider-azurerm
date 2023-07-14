// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

var _ pluginsdk.StateUpgrade = WindowsPerformanceCounterV0ToV1{}

type WindowsPerformanceCounterV0ToV1 struct{}

func (WindowsPerformanceCounterV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId, err := datasources.ParseDataSourceIDInsensitively(rawState["id"].(string))
		if err != nil {
			return rawState, err
		}

		rawState["id"] = oldId.ID()
		return rawState, nil
	}
}

func (WindowsPerformanceCounterV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"workspace_name": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"counter_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"instance_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"interval_seconds": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"object_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}
