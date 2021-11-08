package streamanalytics

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/streamanalytics/mgmt/2020-03-01-preview/streamanalytics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importStreamAnalyticsOutput(expectType streamanalytics.TypeBasicOutputDataSource) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.OutputID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).StreamAnalytics.OutputsClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.Name)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if props := resp.OutputProperties; props != nil {
			var actualType streamanalytics.TypeBasicOutputDataSource

			if datasource, ok := props.Datasource.AsBlobOutputDataSource(); ok {
				actualType = datasource.Type
			} else if datasource, ok := props.Datasource.AsAzureTableOutputDataSource(); ok {
				actualType = datasource.Type
			} else if datasource, ok := props.Datasource.AsEventHubOutputDataSource(); ok {
				actualType = datasource.Type
			} else if datasource, ok := props.Datasource.AsEventHubV2OutputDataSource(); ok {
				actualType = datasource.Type
			} else if datasource, ok := props.Datasource.AsAzureSQLDatabaseOutputDataSource(); ok {
				actualType = datasource.Type
			} else if datasource, ok := props.Datasource.AsAzureSynapseOutputDataSource(); ok {
				actualType = datasource.Type
			} else if datasource, ok := props.Datasource.AsDocumentDbOutputDataSource(); ok {
				actualType = datasource.Type
			} else if datasource, ok := props.Datasource.AsAzureFunctionOutputDataSource(); ok {
				actualType = datasource.Type
			} else if datasource, ok := props.Datasource.AsServiceBusQueueOutputDataSource(); ok {
				actualType = datasource.Type
			} else if datasource, ok := props.Datasource.AsServiceBusTopicOutputDataSource(); ok {
				actualType = datasource.Type
			} else if datasource, ok := props.Datasource.AsPowerBIOutputDataSource(); ok {
				actualType = datasource.Type
			} else if datasource, ok := props.Datasource.AsAzureDataLakeStoreOutputDataSource(); ok {
				actualType = datasource.Type
			} else if datasource, ok := props.Datasource.AsOutputDataSource(); ok {
				actualType = datasource.Type
			} else {
				return nil, fmt.Errorf("unable to convert output data source: %+v", props.Datasource)
			}

			if actualType != expectType {
				return nil, fmt.Errorf("stream analytics output has mismatched type, expected: %q, got %q", expectType, actualType)
			}
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}
