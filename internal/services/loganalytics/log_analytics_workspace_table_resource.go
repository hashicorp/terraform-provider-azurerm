// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogAnalyticsWorkspaceTableResource struct {
}

var _ sdk.ResourceWithUpdate = LogAnalyticsWorkspaceTableResource{}

type LogAnalyticsWorkspaceTableResourceModel struct {
	Name            string `tfschema:"name"`
	WorkspaceId     string `tfschema:"workspace_id"`
	RetentionInDays int64  `tfschema:"retention_in_days"`
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

		"retention_in_days": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.Any(validation.IntBetween(30, 730), validation.IntInSlice([]int{7})),
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

			retentionInDays := model.RetentionInDays
			updateInput := tables.Table{
				Properties: &tables.TableProperties{
					RetentionInDays: &retentionInDays,
				},
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

			updateInput := tables.Table{
				Properties: &tables.TableProperties{
					RetentionInDays: existing.Model.Properties.RetentionInDays,
				},
			}

			if metadata.ResourceData.HasChange("retention_in_days") {
				updateInput.Properties.RetentionInDays = &state.RetentionInDays
			}
			if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput); err != nil {
				return fmt.Errorf("failed to update table: %s: %+v", id.TableName, err)
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
				if model.Properties.RetentionInDays != nil {
					state.RetentionInDays = *model.Properties.RetentionInDays
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

			updateInput := tables.Table{
				Properties: &tables.TableProperties{
					RetentionInDays: retentionInDays,
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, updateInput); err != nil {
				return fmt.Errorf("failed to update table %s in workspace %s in resource group %s: %s", id.TableName, id.WorkspaceName, id.ResourceGroupName, err)
			}

			return nil
		},
	}
}
