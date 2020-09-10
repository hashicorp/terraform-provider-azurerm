package loganalytics

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2020-03-01-preview/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func importLogAnalyticsDataSource(kind operationalinsights.DataSourceKind) func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
	return func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
		id, err := parse.LogAnalyticsDataSourceID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).LogAnalytics.DataSourcesClient
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		resp, err := client.Get(ctx, id.ResourceGroup, id.Workspace, id.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve Log Analytics Data Source %q (Resource Group %q / Workspace: %q): %+v", id.Name, id.ResourceGroup, id.Workspace, err)
		}

		if resp.Kind != kind {
			return nil, fmt.Errorf(`Log Analytics Data Source "kind" mismatch, expected "%s", got "%s"`, kind, resp.Kind)
		}
		return []*schema.ResourceData{d}, nil
	}
}
