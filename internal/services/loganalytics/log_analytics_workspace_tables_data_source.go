// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogAnalyticsWorkspaceTablesDataSource struct{}

var _ sdk.DataSource = LogAnalyticsWorkspaceTablesDataSource{}

type LogAnalyticsWorkspaceTablesDataSourceModel struct {
	WorkspaceId string                       `tfschema:"workspace_id"`
	Names       []string                     `tfschema:"names"`
	Tables      []TableEntityDataSourceModel `tfschema:"tables"`
}

type TableEntityDataSourceModel struct {
	Name                 string `tfschema:"name"`
	RetentionInDays      int64  `tfschema:"retention_in_days"`
	TotalRetentionInDays int64  `tfschema:"total_retention_in_days"`
	Plan                 string `tfschema:"plan"`
}

func (k LogAnalyticsWorkspaceTablesDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},
	}
}

func (k LogAnalyticsWorkspaceTablesDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"names": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		
		"tables": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"plan": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"retention_in_days": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"total_retention_in_days": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},
	}
}

func (k LogAnalyticsWorkspaceTablesDataSource) ModelObject() interface{} {
	return &LogAnalyticsWorkspaceTablesDataSourceModel{}
}

func (k LogAnalyticsWorkspaceTablesDataSource) ResourceType() string {
	return "azurerm_log_analytics_workspace_tables"
}

func (k LogAnalyticsWorkspaceTablesDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state LogAnalyticsWorkspaceTablesDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding model: %+v", err)
			}

			client := metadata.Client.LogAnalytics.TablesClient

			id, err := workspaces.ParseWorkspaceID(state.WorkspaceId)
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			workspaceId := tables.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
			resp, err := client.ListByWorkspace(ctx, workspaceId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", workspaceId)
				}

				return fmt.Errorf("retrieving tables by %s: %+v", workspaceId, err)
			}

			metadata.SetID(workspaceId)

			if model := resp.Model; model != nil {
				tables := make([]TableEntityDataSourceModel, 0)
				names := make([]string, 0)

				if model.Value != nil && len(*model.Value) > 0 {
					for _, props := range *model.Value {
						name := pointer.From(props.Name)
						names = append(names, name)
						tables = append(tables, flattenLogAnalyticsWorkspaceTable(name, props.Properties))
					}
				}

				state.Tables = tables
				state.Names = names
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenLogAnalyticsWorkspaceTable(name string, properties *tables.TableProperties) TableEntityDataSourceModel {
	table := TableEntityDataSourceModel{
		Name: name,
	}

	if properties != nil {
		table.RetentionInDays = pointer.From(properties.RetentionInDays)
		table.TotalRetentionInDays = pointer.From(properties.TotalRetentionInDays)
		table.Plan = string(pointer.From(properties.Plan))
	}

	return table
}
