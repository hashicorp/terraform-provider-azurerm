// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_endpoint_cosmosdb_account

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.ResourceClient
			id, err := parse.EndpointCosmosDBAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.IotHubName, IothubResourceName)
			defer locks.UnlockByName(id.IotHubName, IothubResourceName)

			iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
			if err != nil {
				if response.WasNotFound(iothub.Response.Response) {
					return fmt.Errorf("%q was not found", id)
				}
				return fmt.Errorf("retrieving %q: %+v", id, err)
			}

			if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil || iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections == nil {
				return nil
			}

			updatedEndpoints := make([]devices.RoutingCosmosDBSQLAPIProperties, 0)
			for _, endpoint := range pointer.From(iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections) {
				if !strings.EqualFold(pointer.From(endpoint.Name), id.EndpointName) {
					updatedEndpoints = append(updatedEndpoints, endpoint)
				}
			}
			iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections = &updatedEndpoints

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the deletion of %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
