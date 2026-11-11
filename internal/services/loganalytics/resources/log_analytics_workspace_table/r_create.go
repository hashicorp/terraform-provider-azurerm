// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace_table

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
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
