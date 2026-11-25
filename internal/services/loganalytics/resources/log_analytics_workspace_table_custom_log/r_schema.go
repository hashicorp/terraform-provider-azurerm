// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace_table_custom_log

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var table WorkspaceTableCustomLogResourceModel
			if err := metadata.DecodeDiff(&table); err != nil {
				return err
			}

			if table.Plan == string(tables.TablePlanEnumBasic) {
				if _, ok := metadata.ResourceDiff.GetOk("retention_in_days"); ok {
					return errors.New("`retention_in_days` cannot be set when `plan` is set to `Basic`")
				}
			}

			return nil
		},
	}
}

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Log Analytics Workspace Table Custom Log.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`_CL$`), "must end with '_CL'."),
		},

		"column": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: workspaceTableColumnSchema(),
			},
		},

		"workspace_id": {
			Description:  "The object ID of the Log Analytics Workspace that contains the table.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"description": {
			Description:  "The description of the table.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Description:  "The display name of the table.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"plan": {
			Description:  "The system how to handle and charge the logs ingested to the table.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(tables.TablePlanEnumAnalytics),
			ValidateFunc: validation.StringInSlice(tables.PossibleValuesForTablePlanEnum(), false),
		},

		"retention_in_days": {
			Description:  "The table's retention in days.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(4, 730),
		},

		"total_retention_in_days": {
			Description:  "The table's total retention in days.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.Any(validation.IntBetween(4, 730), validation.IntInSlice([]int{1095, 1460, 1826, 2191, 2556, 2922, 3288, 3653, 4018, 4383})),
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"solutions": {
			Description: "A list of solutions associated with the table.",
			Type:        pluginsdk.TypeSet,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"standard_column": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: workspaceTableColumnSchemaComputed(),
			},
		},
	}
}
