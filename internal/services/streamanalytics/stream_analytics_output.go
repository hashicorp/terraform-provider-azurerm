package streamanalytics

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importStreamAnalyticsOutput(expectType outputs.OutputDataSource) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := outputs.ParseOutputID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).StreamAnalytics.OutputsClient
		resp, err := client.Get(ctx, *id)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		var actualType outputs.OutputDataSource
		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				if ds := props.Datasource; ds != nil {
					var datasource outputs.OutputDataSource
					var ok bool

					if datasource, ok = (*ds).(outputs.BlobOutputDataSource); ok {
						actualType = datasource
					} else if datasource, ok = (*ds).(outputs.AzureTableOutputDataSource); ok {
						actualType = datasource
					} else if datasource, ok = (*ds).(outputs.EventHubOutputDataSource); ok {
						actualType = datasource
					} else if datasource, ok = (*ds).(outputs.EventHubV2OutputDataSource); ok {
						actualType = datasource
					} else if datasource, ok = (*ds).(outputs.AzureSqlDatabaseOutputDataSource); ok {
						actualType = datasource
					} else if datasource, ok = (*ds).(outputs.AzureSynapseOutputDataSource); ok {
						actualType = datasource
					} else if datasource, ok = (*ds).(outputs.DocumentDbOutputDataSource); ok {
						actualType = datasource
					} else if datasource, ok = (*ds).(outputs.AzureFunctionOutputDataSource); ok {
						actualType = datasource
					} else if datasource, ok = (*ds).(outputs.ServiceBusQueueOutputDataSource); ok {
						actualType = datasource
					} else if datasource, ok = (*ds).(outputs.ServiceBusTopicOutputDataSource); ok {
						actualType = datasource
					} else if datasource, ok = (*ds).(outputs.PowerBIOutputDataSource); ok {
						actualType = datasource
					} else if datasource, ok = (*ds).(outputs.AzureDataLakeStoreOutputDataSource); ok {
						actualType = datasource
					}

					if !ok {
						return nil, fmt.Errorf("unable to convert output data source: %+v", props.Datasource)
					}
				}
			}
		}

		if reflect.TypeOf(actualType) != reflect.TypeOf(expectType) {
			return nil, fmt.Errorf("stream analytics output has mismatched type, expected: %q, got %q", expectType, actualType)
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}
