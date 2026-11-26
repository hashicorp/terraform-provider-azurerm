// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace_table_custom_log

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2025-07-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.TablesClient

			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := WorkspaceTableCustomLogResourceModel{
				WorkspaceId: pointer.To(workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)).ID(),
				Name:        id.TableName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Plan = pointer.FromEnum(props.Plan)

					if asDefault := props.RetentionInDaysAsDefault; asDefault != nil && !*asDefault {
						state.RetentionInDays = pointer.From(props.RetentionInDays)
					}
					if asDefault := props.TotalRetentionInDaysAsDefault; asDefault != nil && !*asDefault {
						state.TotalRetentionInDays = pointer.From(props.TotalRetentionInDays)
					}

					if schema := props.Schema; schema != nil {
						state.DisplayName = pointer.From(schema.DisplayName)
						state.Description = pointer.From(schema.Description)
						state.Solutions = pointer.From(schema.Solutions)
						state.Columns = flattenWorkspaceTableColumns(schema.Columns)
						state.StandardColumns = flattenWorkspaceTableColumns(schema.StandardColumns)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
