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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.ResourceClient
			id, err := parse.EndpointCosmosDBAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var oldState IotHubEndpointCosmosDBAccountModel
			if err = metadata.Decode(&oldState); err != nil {
				return err
			}

			iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
			if err != nil {
				if response.WasNotFound(iothub.Response.Response) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %q: %+v", id, err)
			}

			if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil || iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections == nil {
				return metadata.MarkAsGone(id)
			}

			for _, endpoint := range pointer.From(iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections) {
				if strings.EqualFold(pointer.From(endpoint.Name), id.EndpointName) {
					state := &IotHubEndpointCosmosDBAccountModel{
						Name:                 id.EndpointName,
						ResourceGroupName:    pointer.From(endpoint.ResourceGroup),
						IothubId:             parse.NewIotHubID(id.SubscriptionId, id.ResourceGroup, id.IotHubName).ID(),
						ContainerName:        pointer.From(endpoint.CollectionName),
						DatabaseName:         pointer.From(endpoint.DatabaseName),
						EndpointUri:          pointer.From(endpoint.EndpointURI),
						PartitionKeyName:     pointer.From(endpoint.PartitionKeyName),
						PartitionKeyTemplate: pointer.From(endpoint.PartitionKeyTemplate),
						PrimaryKey:           oldState.PrimaryKey,
						SecondaryKey:         oldState.SecondaryKey,
						SubscriptionId:       pointer.From(endpoint.SubscriptionID),
					}

					authenticationType := string(devices.AuthenticationTypeKeyBased)
					if string(endpoint.AuthenticationType) != "" {
						authenticationType = string(endpoint.AuthenticationType)
					}
					state.AuthenticationType = authenticationType

					identityId := ""
					if endpoint.Identity != nil && endpoint.Identity.UserAssignedIdentity != nil {
						identityId = pointer.From(endpoint.Identity.UserAssignedIdentity)
					}
					state.IdentityId = identityId

					return metadata.Encode(state)
				}
			}

			return metadata.MarkAsGone(id)
		},
		Timeout: 5 * time.Minute,
	}
}
