// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package log_analytics_workspace_table_microsoft

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkspaceTableMicrosoftResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}
			client := metadata.Client.LogAnalytics.TablesClient

			tableName := model.Name

			workspaceId, err := workspaces.ParseWorkspaceID(model.WorkspaceId)
			if err != nil {
				return err
			}

			id := tables.NewTableID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, tableName)
			if _, err = client.Get(ctx, id); err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", tableName, err)
			}
			// Microsoft tables are always automatically provisioned whenever log analytics workspaces are provisioned, so there's no point in returning a 'resource already exists' error

			param := tables.Table{
				Properties: &tables.TableProperties{
					Plan:                 pointer.To(tables.TablePlanEnumAnalytics),
					RetentionInDays:      helpers.DefaultRetentionInDays,
					TotalRetentionInDays: helpers.DefaultRetentionInDays,
					Schema: &tables.Schema{
						Columns:     expandWorkspaceTableMicrosoftColumns(model.Columns),
						DisplayName: pointer.To(model.DisplayName),
						Description: pointer.To(model.Description),
						Labels:      pointer.To(model.Labels),
						Name:        pointer.To(tableName),
						TableType:   pointer.To(tables.TableTypeEnumMicrosoft),
					},
				},
			}

			if model.RetentionInDays > 0 {
				param.Properties.RetentionInDays = pointer.To(model.RetentionInDays)
			}

			if model.TotalRetentionInDays > 0 {
				param.Properties.TotalRetentionInDays = pointer.To(model.TotalRetentionInDays)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandWorkspaceTableMicrosoftColumns(columns []Column) *[]tables.Column {
	result := make([]tables.Column, 0, len(columns))
	for _, column := range columns {
		columnToAdd := tables.Column{
			Name:             pointer.To(column.Name),
			Type:             pointer.ToEnum[tables.ColumnTypeEnum](column.Type),
			IsHidden:         pointer.To(column.Hidden),
			IsDefaultDisplay: pointer.To(column.DisplayByDefault),
		}
		if column.DisplayName != "" {
			columnToAdd.DisplayName = pointer.To(column.DisplayName)
		}
		if column.Description != "" {
			columnToAdd.Description = pointer.To(column.Description)
		}
		result = append(result, columnToAdd)
	}
	return pointer.To(result)
}
