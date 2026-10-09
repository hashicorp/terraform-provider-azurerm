// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LogAnalyticsWorkspaceTableResource struct{}

var (
	_ sdk.ResourceWithUpdate        = LogAnalyticsWorkspaceTableResource{}
	_ sdk.ResourceWithCustomizeDiff = LogAnalyticsWorkspaceTableResource{}
)

type LogAnalyticsWorkspaceTableResourceModel struct {
	Name                 string `tfschema:"name"`
	WorkspaceId          string `tfschema:"workspace_id"`
	Plan                 string `tfschema:"plan"`
	RetentionInDays      int64  `tfschema:"retention_in_days"`
	TotalRetentionInDays int64  `tfschema:"total_retention_in_days"`
}

func (r LogAnalyticsWorkspaceTableResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if string(tables.TablePlanEnumBasic) == rd.Get("plan").(string) {
				if v, ok := rd.GetOk("retention_in_days"); ok && v.(int) != 30 {
					return fmt.Errorf("cannot set retention_in_days because the retention is fixed at 30 days on Basic plan")
				}
			}

			return nil
		},
	}
}

func (r LogAnalyticsWorkspaceTableResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"plan": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(tables.TablePlanEnumAnalytics),
			ValidateFunc: validation.StringInSlice(tables.PossibleValuesForTablePlanEnum(), false),
		},

		"retention_in_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(4, 730),
		},

		"total_retention_in_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.Any(validation.IntBetween(4, 730), validation.IntInSlice([]int{1095, 1460, 1826, 2191, 2556, 2922, 3288, 3653, 4018, 4383})),
		},
	}
}

func (r LogAnalyticsWorkspaceTableResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LogAnalyticsWorkspaceTableResource) ModelObject() any {
	return &LogAnalyticsWorkspaceTableResourceModel{}
}

func (r LogAnalyticsWorkspaceTableResource) ResourceType() string {
	return "azurerm_log_analytics_workspace_table"
}

func (r LogAnalyticsWorkspaceTableResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return tables.ValidateTableID
}

func (r LogAnalyticsWorkspaceTableResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LogAnalyticsWorkspaceTableResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.LogAnalytics.TablesClient

			tableName := model.Name

			workspaceId, err := workspaces.ParseWorkspaceID(model.WorkspaceId)
			if err != nil {
				return fmt.Errorf("invalid workspace object ID for table %s: %s", tableName, err)
			}

			// TODO: Import check

			id := tables.NewTableID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, tableName)

			updateInput := tables.Table{
				Properties: &tables.TableProperties{
					Plan: pointer.ToEnum[tables.TablePlanEnum](model.Plan),
				},
			}

			if model.Plan == string(tables.TablePlanEnumAnalytics) {
				// The service will return HTTP 400 if it's specified `0` in payload, to keep it as default, we need to pass `-1`
				updateInput.Properties.RetentionInDays = pointer.To(int64(-1))
				// `0` is not a valid value for `retention_in_days`, so we can use it to validate if it's specified.
				if model.RetentionInDays != 0 {
					updateInput.Properties.RetentionInDays = pointer.To(model.RetentionInDays)
				}
			}

			updateInput.Properties.TotalRetentionInDays = pointer.To(int64(-1))
			if model.TotalRetentionInDays != 0 {
				updateInput.Properties.TotalRetentionInDays = pointer.To(model.TotalRetentionInDays)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, updateInput, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("failed to update table %s in workspace %s in resource group %s: %s", tableName, workspaceId.WorkspaceName, workspaceId.ResourceGroupName, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LogAnalyticsWorkspaceTableResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.TablesClient
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LogAnalyticsWorkspaceTableResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			props := existing.Model.Properties
			if props.Schema != nil {
				// Azure rejects StandardColumns in CreateOrUpdate requests.
				props.Schema.StandardColumns = nil
			}

			if pointer.From(props.RetentionInDaysAsDefault) {
				props.RetentionInDays = pointer.To(int64(-1))
			}
			if pointer.From(props.TotalRetentionInDaysAsDefault) {
				props.TotalRetentionInDays = pointer.To(int64(-1))
			}

			if metadata.ResourceData.HasChange("plan") {
				props.Plan = pointer.ToEnum[tables.TablePlanEnum](state.Plan)
			}

			if state.Plan == string(tables.TablePlanEnumAnalytics) {
				if metadata.ResourceData.HasChange("retention_in_days") {
					props.RetentionInDays = pointer.To(state.RetentionInDays)
					// Azure uses -1 to reset retention to the workspace default.
					if state.RetentionInDays == 0 {
						props.RetentionInDays = pointer.To(int64(-1))
					}
				}
			} else {
				// Retention is read-only for Basic tables.
				props.RetentionInDays = nil
			}

			if metadata.ResourceData.HasChange("total_retention_in_days") {
				props.TotalRetentionInDays = pointer.To(state.TotalRetentionInDays)
				// Azure uses -1 to reset total retention to the table retention.
				if state.TotalRetentionInDays == 0 {
					props.TotalRetentionInDays = pointer.To(int64(-1))
				}
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LogAnalyticsWorkspaceTableResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

			client := metadata.Client.LogAnalytics.TablesClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving Log Analytics Workspace Table %s: %+v", *id, err)
			}

			state := LogAnalyticsWorkspaceTableResourceModel{
				Name:        id.TableName,
				WorkspaceId: workspaceId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if !pointer.From(props.RetentionInDaysAsDefault) && pointer.From(props.Plan) == tables.TablePlanEnumAnalytics {
						state.RetentionInDays = pointer.From(props.RetentionInDays)
					}

					if !pointer.From(props.TotalRetentionInDaysAsDefault) {
						state.TotalRetentionInDays = pointer.From(props.TotalRetentionInDays)
					}

					state.Plan = pointer.FromEnum(props.Plan)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LogAnalyticsWorkspaceTableResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LogAnalyticsWorkspaceTableResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.LogAnalytics.TablesClient
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			// We do not delete the resource here, just set the retention to workspace default value, which is
			// achieved by setting the value to `-1`
			retentionInDays := pointer.To(int64(-1))
			totalRetentionInDays := pointer.To(int64(-1))

			updateInput := tables.Table{
				Properties: &tables.TableProperties{
					RetentionInDays:      retentionInDays,
					TotalRetentionInDays: totalRetentionInDays,
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput); err != nil {
				return fmt.Errorf("failed to update table %s in workspace %s in resource group %s: %s", id.TableName, id.WorkspaceName, id.ResourceGroupName, err)
			}

			return nil
		},
	}
}
