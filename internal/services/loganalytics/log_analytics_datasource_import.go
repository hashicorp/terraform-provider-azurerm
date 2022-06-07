package loganalytics

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importLogAnalyticsDataSource(kind operationalinsights.DataSourceKind) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.DataSourceID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Kind != kind {
			return nil, fmt.Errorf(`log analytics Data Source "kind" mismatch, expected "%s", got "%s"`, kind, resp.Kind)
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}
