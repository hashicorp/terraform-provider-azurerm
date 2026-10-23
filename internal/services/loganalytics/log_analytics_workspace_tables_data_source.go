// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogAnalyticsWorkspaceTablesDataSource struct{}

var _ sdk.DataSource = LogAnalyticsWorkspaceTablesDataSource{}

type LogAnalyticsWorkspaceTablesDataSourceModel struct {
	WorkspaceId string                            `tfschema:"workspace_id"`
	Names       []string                          `tfschema:"names"`
	Tables      []LogAnalyticsWorkspaceTableModel `tfschema:"tables"`
}

type LogAnalyticsWorkspaceTableModel struct {
	Name                 string `tfschema:"name"`
	Plan                 string `tfschema:"plan"`
	RetentionInDays      int64  `tfschema:"retention_in_days"`
	TotalRetentionInDays int64  `tfschema:"total_retention_in_days"`
}

func (LogAnalyticsWorkspaceTablesDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: tables.ValidateWorkspaceID,
		},
	}
}

func (LogAnalyticsWorkspaceTablesDataSource) Attributes() map[string]*pluginsdk.Schema {
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

func (LogAnalyticsWorkspaceTablesDataSource) ModelObject() any {
	return &LogAnalyticsWorkspaceTablesDataSourceModel{}
}

func (LogAnalyticsWorkspaceTablesDataSource) ResourceType() string {
	return "azurerm_log_analytics_workspace_tables"
}

func (LogAnalyticsWorkspaceTablesDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.TablesClient

			var state LogAnalyticsWorkspaceTablesDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := tables.ParseWorkspaceID(state.WorkspaceId)
			if err != nil {
				return err
			}

			resp, err := client.ListByWorkspace(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving tables for %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving tables for %s: `model` was nil", id)
			}

			tableList := pointer.From(resp.Model.Value)
			state.Tables = make([]LogAnalyticsWorkspaceTableModel, 0, len(tableList))
			state.Names = make([]string, 0, len(tableList))

			for _, table := range tableList {
				name := pointer.From(table.Name)
				state.Names = append(state.Names, name)
				state.Tables = append(state.Tables, flattenLogAnalyticsWorkspaceTable(name, table.Properties))
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func flattenLogAnalyticsWorkspaceTable(name string, properties *tables.TableProperties) LogAnalyticsWorkspaceTableModel {
	table := LogAnalyticsWorkspaceTableModel{
		Name: name,
	}

	if properties != nil {
		table.RetentionInDays = pointer.From(properties.RetentionInDays)
		table.TotalRetentionInDays = pointer.From(properties.TotalRetentionInDays)
		table.Plan = pointer.FromEnum(properties.Plan)
	}

	return table
}
