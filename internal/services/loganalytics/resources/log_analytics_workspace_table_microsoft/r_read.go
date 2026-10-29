// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace_table_microsoft

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := tables.ParseTableID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			workspaceId := pointer.To(workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName))

			client := metadata.Client.LogAnalytics.TablesClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := WorkspaceTableMicrosoftResourceModel{
				Name:        id.TableName,
				WorkspaceId: workspaceId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if pointer.From(props.Plan) == tables.TablePlanEnumAnalytics {
						if asDefault := props.RetentionInDaysAsDefault; asDefault != nil && !*asDefault {
							state.RetentionInDays = pointer.From(props.RetentionInDays)
						}
						if asDefault := props.TotalRetentionInDaysAsDefault; asDefault != nil && !*asDefault {
							state.TotalRetentionInDays = pointer.From(props.TotalRetentionInDays)
						}
					}

					if schema := props.Schema; schema != nil {
						state.DisplayName = pointer.From(schema.DisplayName)
						state.Description = pointer.From(schema.Description)
						state.Labels = pointer.From(schema.Labels)
						state.Solutions = pointer.From(schema.Solutions)
						state.Columns = flattenWorkspaceTableMicrosoftColumns(schema.Columns)
						state.StandardColumns = flattenWorkspaceTableMicrosoftColumns(schema.StandardColumns)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenWorkspaceTableMicrosoftColumns(columns *[]tables.Column) []Column {
	if columns == nil {
		return []Column{}
	}
	result := make([]Column, 0, len(*columns))
	for _, column := range *columns {
		result = append(result, Column{
			Name:             pointer.From(column.Name),
			DisplayName:      pointer.From(column.DisplayName),
			Description:      pointer.From(column.Description),
			Type:             pointer.FromEnum(column.Type),
			Hidden:           pointer.From(column.IsHidden),
			DisplayByDefault: pointer.From(column.IsDefaultDisplay),
		})
	}
	return result
}

type Column struct {
	Name             string `tfschema:"name"`
	DisplayName      string `tfschema:"display_name"`
	Description      string `tfschema:"description"`
	Type             string `tfschema:"type"`
	Hidden           bool   `tfschema:"hidden"`
	DisplayByDefault bool   `tfschema:"display_by_default"`
}
