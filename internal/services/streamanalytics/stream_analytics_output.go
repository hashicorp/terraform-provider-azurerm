// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importStreamAnalyticsOutput(expectTypes []outputs.OutputDataSource) pluginsdk.ImporterFunc {
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

		if props := resp.Model.Properties; props != nil {
			var actualType outputs.OutputDataSource

			if datasource, ok := props.Datasource.(outputs.BlobOutputDataSource); ok {
				actualType = datasource
			} else if datasource, ok := props.Datasource.(outputs.AzureTableOutputDataSource); ok {
				actualType = datasource
			} else if datasource, ok := props.Datasource.(outputs.EventHubOutputDataSource); ok {
				actualType = datasource
			} else if datasource, ok := props.Datasource.(outputs.EventHubV2OutputDataSource); ok {
				actualType = datasource
			} else if datasource, ok := props.Datasource.(outputs.AzureSqlDatabaseOutputDataSource); ok {
				actualType = datasource
			} else if datasource, ok := props.Datasource.(outputs.AzureSynapseOutputDataSource); ok {
				actualType = datasource
			} else if datasource, ok := props.Datasource.(outputs.DocumentDbOutputDataSource); ok {
				actualType = datasource
			} else if datasource, ok := props.Datasource.(outputs.AzureFunctionOutputDataSource); ok {
				actualType = datasource
			} else if datasource, ok := props.Datasource.(outputs.ServiceBusQueueOutputDataSource); ok {
				actualType = datasource
			} else if datasource, ok := props.Datasource.(outputs.ServiceBusTopicOutputDataSource); ok {
				actualType = datasource
			} else if datasource, ok := props.Datasource.(outputs.PowerBIOutputDataSource); ok {
				actualType = datasource
			} else if datasource, ok := props.Datasource.(outputs.AzureDataLakeStoreOutputDataSource); ok {
				actualType = datasource
			} else {
				return nil, fmt.Errorf("unable to convert output data source: %+v", props.Datasource)
			}

			var typeAsExpected bool
			var expectedTypeNames []string
			for _, expectType := range expectTypes {
				expectedTypeNames = append(expectedTypeNames, fmt.Sprintf("%T", expectType.OutputDataSource().Type))
				if reflect.TypeOf(actualType) != reflect.TypeOf(expectType) {
					continue
				}
				typeAsExpected = true
			}

			if !typeAsExpected {
				return nil, fmt.Errorf("stream analytics output has mismatched type, expected: %q, got %q", expectedTypeNames, actualType.OutputDataSource().Type)
			}
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}
