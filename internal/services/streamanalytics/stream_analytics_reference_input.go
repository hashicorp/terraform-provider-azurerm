package streamanalytics

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/streamanalytics/mgmt/2020-03-01-preview/streamanalytics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importStreamAnalyticsReferenceInput(expectType streamanalytics.TypeBasicReferenceInputDataSource) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.StreamInputID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).StreamAnalytics.InputsClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.StreamingjobName, id.InputName)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if props := resp.Properties; props != nil {
			v, ok := props.AsReferenceInputProperties()
			if !ok {
				return nil, fmt.Errorf("converting properties to a Reference Input: %+v", err)
			}

			var actualType streamanalytics.TypeBasicReferenceInputDataSource

			if inputMsSql, ok := v.Datasource.AsAzureSQLReferenceInputDataSource(); ok {
				actualType = inputMsSql.Type
			} else if inputBlob, ok := v.Datasource.AsBlobReferenceInputDataSource(); ok {
				actualType = inputBlob.Type
			} else {
				return nil, fmt.Errorf("unable to convert input data source: %+v", v)
			}

			if actualType != expectType {
				return nil, fmt.Errorf("stream analytics reference input has mismatched type, expected: %q, got %q", expectType, actualType)
			}
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}
