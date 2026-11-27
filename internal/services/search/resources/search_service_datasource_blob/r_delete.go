// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package search_service_datasource_blob

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/data-plane/search/2025-09-01/datasources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	searchClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/client"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			resourceId, err := datasources.ParseDatasourceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := searchClient.From(metadata.Client).SearchDataPlaneClient.DataSources.Clone(resourceId.BaseURI)

			if _, err := client.Delete(ctx, *resourceId, datasources.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *resourceId, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
