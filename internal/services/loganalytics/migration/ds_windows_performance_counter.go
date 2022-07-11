package migration

import (
	"context"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
)

var _ pluginsdk.StateUpgrade = WindowsPerformanceCounterV0ToV1{}

type WindowsPerformanceCounterV0ToV1 struct{}

func (WindowsPerformanceCounterV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId, err := parse.DataSourceIDInsensitively(rawState["id"].(string))
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

		"resource_group_name": azure.SchemaResourceGroupName(),

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
