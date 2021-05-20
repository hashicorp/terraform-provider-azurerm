package loganalytics

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func importLogAnalyticsDataSource(kind operationalinsights.DataSourceKind) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.LogAnalyticsDataSourceID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve Log Analytics Data Source %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
		}

		if resp.Kind != kind {
			return nil, fmt.Errorf(`Log Analytics Data Source "kind" mismatch, expected "%s", got "%s"`, kind, resp.Kind)
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}
