package streamanalytics

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importStreamAnalyticsReferenceInput(expectType string) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := inputs.ParseInputID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).StreamAnalytics.InputsClient
		resp, err := client.Get(ctx, *id)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		var actualType string
		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				if reference, ok := (*props).(inputs.ReferenceInputProperties); ok {
					if ds := reference.Datasource; ds != nil {
						if _, ok := (*ds).(inputs.BlobReferenceInputDataSource); ok {
							actualType = "Microsoft.Storage/Blob"
						}
						if _, ok := (*ds).(inputs.AzureSqlReferenceInputDataSource); ok {
							actualType = "Microsoft.Sql/Server/Database"
						}
					}
				}
			}
		}

		if actualType != expectType {
			return nil, fmt.Errorf("expected the Reference Input to be a %q but got got %q", expectType, actualType)
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}
