package loganalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/tables"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LogAnalyticsWorkspaceTableDataSource struct{}

type LogAnalyticsWorkspaceTableDataSourceModel struct {
	ResourceGroupName    string `tfschema:"resource_group_name"`
	Name                 string `tfschema:"name"`
	WorkspaceId          string `tfschema:"workspace_id"`
	Plan                 string `tfschema:"plan"`
	RetentionInDays      int64  `tfschema:"retention_in_days"`
	TotalRetentionInDays int64  `tfschema:"total_retention_in_days"`
}

func (LogAnalyticsWorkspaceTableDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (LogAnalyticsWorkspaceTableDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"retention_in_days": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
		"total_retention_in_days": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
		"plan": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (LogAnalyticsWorkspaceTableDataSource) ModelObject() interface{} {
	return &LogAnalyticsWorkspaceTableDataSourceModel{}
}

func (LogAnalyticsWorkspaceTableDataSource) ResourceType() string {
	return "azurerm_log_analytics_workspace_table"
}

func (LogAnalyticsWorkspaceTableDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.TablesClient

			var state LogAnalyticsWorkspaceTableDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := workspaces.ParseWorkspaceID(state.WorkspaceId)
			if err != nil {
				return fmt.Errorf("invalid workspace object ID for table %s: %s", state.Name, err)
			}

			id := tables.NewTableID(metadata.Client.Account.SubscriptionId, state.ResourceGroupName, workspaceId.WorkspaceName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.RetentionInDays = pointer.From(props.RetentionInDays)
					state.TotalRetentionInDays = pointer.From(props.TotalRetentionInDays)
					state.Plan = string(pointer.From(props.Plan))
				}
			}
			return metadata.Encode(&state)
		},
	}
}
