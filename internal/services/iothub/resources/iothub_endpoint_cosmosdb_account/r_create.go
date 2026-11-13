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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state IotHubEndpointCosmosDBAccountModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.IoTHub.ResourceClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			iotHubId, err := parse.IotHubID(state.IothubId)
			if err != nil {
				return err
			}

			id := parse.NewEndpointCosmosDBAccountID(subscriptionId, iotHubId.ResourceGroup, iotHubId.Name, state.Name)

			locks.ByName(iotHubId.Name, IothubResourceName)
			defer locks.UnlockByName(iotHubId.Name, IothubResourceName)

			iothub, err := client.Get(ctx, iotHubId.ResourceGroup, iotHubId.Name)
			if err != nil {
				if response.WasNotFound(iothub.Response.Response) {
					return fmt.Errorf("%q was not found", iotHubId)
				}

				return fmt.Errorf("retrieving %q: %+v", iotHubId, err)
			}

			authenticationType := devices.AuthenticationType(state.AuthenticationType)
			cosmosDBAccountEndpoint := devices.RoutingCosmosDBSQLAPIProperties{
				Name:               pointer.To(id.EndpointName),
				ResourceGroup:      pointer.To(state.ResourceGroupName),
				AuthenticationType: authenticationType,
				CollectionName:     pointer.To(state.ContainerName),
				DatabaseName:       pointer.To(state.DatabaseName),
				EndpointURI:        pointer.To(state.EndpointUri),
			}

			// To align with the previous TF behaviour, `subscription_id` needs to be set with the provider's subscription Id when it isn't specified in the tf config, otherwise TF behaviour is different than before and it may block the existing users
			// From the business perspective, the raw config handling is only meant for the case that the user has an CosmosDB Account whose Endpoint's subscription is not the provider's one. Then the user wants to reset it to the provider's one by unset the subscription_id
			// From the TF code perspective, given `Computed: true` is enabled, TF would always get the value from the last apply when this property isn't set in the tf config. So `d.GetRawConfig()` is required to determine if it's set in the tf config
			if v := metadata.ResourceData.GetRawConfig().AsValueMap()["subscription_id"]; v.IsNull() {
				cosmosDBAccountEndpoint.SubscriptionID = pointer.To(subscriptionId)
			} else {
				cosmosDBAccountEndpoint.SubscriptionID = pointer.To(state.SubscriptionId)
			}

			if state.PartitionKeyName != "" {
				cosmosDBAccountEndpoint.PartitionKeyName = pointer.To(state.PartitionKeyName)
			}

			if state.PartitionKeyTemplate != "" {
				cosmosDBAccountEndpoint.PartitionKeyTemplate = pointer.To(state.PartitionKeyTemplate)
			}

			if authenticationType == devices.AuthenticationTypeKeyBased {
				if state.PrimaryKey == "" || state.SecondaryKey == "" {
					return fmt.Errorf("`primary_key` and `secondary_key` must be specified when `authentication_type` is `keyBased`")
				}
				cosmosDBAccountEndpoint.PrimaryKey = pointer.To(state.PrimaryKey)
				cosmosDBAccountEndpoint.SecondaryKey = pointer.To(state.SecondaryKey)
			} else {
				if state.PrimaryKey != "" || state.SecondaryKey != "" {
					return fmt.Errorf("`primary_key` or `secondary_key` cannot be specified when `authentication_type` is `identityBased`")
				}

				if state.IdentityId != "" {
					cosmosDBAccountEndpoint.Identity = &devices.ManagedIdentity{
						UserAssignedIdentity: pointer.To(state.IdentityId),
					}
				}
			}

			routing := iothub.Properties.Routing
			if routing == nil {
				routing = &devices.RoutingProperties{}
			}

			if routing.Endpoints == nil {
				routing.Endpoints = &devices.RoutingEndpoints{}
			}

			if routing.Endpoints.CosmosDBSQLCollections == nil {
				cosmosDBAccounts := make([]devices.RoutingCosmosDBSQLAPIProperties, 0)
				routing.Endpoints.CosmosDBSQLCollections = &cosmosDBAccounts
			}

			endpoints := make([]devices.RoutingCosmosDBSQLAPIProperties, 0)

			for _, existingEndpoint := range pointer.From(routing.Endpoints.CosmosDBSQLCollections) {
				if strings.EqualFold(pointer.From(existingEndpoint.Name), id.EndpointName) {
					if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
						return tf.ImportAsExistsError(r.ResourceType(), id.ID())
					}
				}
				endpoints = append(endpoints, existingEndpoint)
			}

			endpoints = append(endpoints, cosmosDBAccountEndpoint)
			routing.Endpoints.CosmosDBSQLCollections = &endpoints

			future, err := client.CreateOrUpdate(ctx, iotHubId.ResourceGroup, iotHubId.Name, iothub, "")
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the completion of the creation of %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
