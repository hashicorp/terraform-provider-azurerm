// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
				if _, ok := rd.GetOk("retention_in_days"); ok {
					return fmt.Errorf("cannot set retention_in_days because the retention is fixed at eight days on Basic plan")
				}
			}

			return nil
		},
	}
}

func (r LogAnalyticsWorkspaceTableResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"plan": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(tables.TablePlanEnumAnalytics),
			ValidateFunc: validation.StringInSlice([]string{
				string(tables.TablePlanEnumAnalytics),
				string(tables.TablePlanEnumBasic),
			}, false),
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

func (r LogAnalyticsWorkspaceTableResource) ModelObject() interface{} {
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
			subscriptionId := metadata.Client.Account.SubscriptionId

			tableName := model.Name
			log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Workspace Table %s update.", tableName)

			workspaceId, err := workspaces.ParseWorkspaceID(model.WorkspaceId)
			if err != nil {
				return fmt.Errorf("invalid workspace object ID for table %s: %s", tableName, err)
			}

			id := tables.NewTableID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, tableName)

			updateInput := tables.Table{
				Properties: &tables.TableProperties{
					Plan: pointer.To(tables.TablePlanEnum(model.Plan)),
				},
			}

			if model.Plan == string(tables.TablePlanEnumAnalytics) {
				updateInput.Properties.RetentionInDays = pointer.To(model.RetentionInDays)
			}

			if model.TotalRetentionInDays != 0 {
				updateInput.Properties.TotalRetentionInDays = pointer.To(model.TotalRetentionInDays)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, updateInput); err != nil {
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
				return fmt.Errorf("reading Log Analytics Workspace Table %s: %v", id, err)
			}

			if model := existing.Model; model != nil {
				if props := model.Properties; props != nil {
					updateInput := tables.Table{
						Properties: &tables.TableProperties{
							Plan: props.Plan,
						},
					}

					if metadata.ResourceData.HasChange("plan") {
						updateInput.Properties.Plan = pointer.To(tables.TablePlanEnum(state.Plan))
					}

					if state.Plan == string(tables.TablePlanEnumAnalytics) {
						updateInput.Properties.RetentionInDays = existing.Model.Properties.RetentionInDays

						if metadata.ResourceData.HasChange("retention_in_days") {
							updateInput.Properties.RetentionInDays = pointer.To(state.RetentionInDays)
						}
					}

					if metadata.ResourceData.HasChange("total_retention_in_days") {
						updateInput.Properties.TotalRetentionInDays = pointer.To(state.TotalRetentionInDays)
					}

					if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput); err != nil {
						return fmt.Errorf("failed to update table: %s: %+v", id.TableName, err)
					}
				}
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
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			workspaceId, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Get("workspace_id").(string))
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

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
					if pointer.From(props.Plan) == tables.TablePlanEnumAnalytics {
						state.RetentionInDays = pointer.From(props.RetentionInDays)
					}
					state.TotalRetentionInDays = pointer.From(props.TotalRetentionInDays)
					state.Plan = string(pointer.From(props.Plan))
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
			retentionInDays := utils.Int64(-1)
			totalRetentionInDays := utils.Int64(-1)

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
