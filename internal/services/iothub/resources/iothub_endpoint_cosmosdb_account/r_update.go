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

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.IoTHub.ResourceClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id, err := parse.EndpointCosmosDBAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.IotHubName, IothubResourceName)
			defer locks.UnlockByName(id.IotHubName, IothubResourceName)

			var state IotHubEndpointCosmosDBAccountModel
			if err = metadata.Decode(&state); err != nil {
				return err
			}

			iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
			if err != nil {
				if response.WasNotFound(iothub.Response.Response) {
					return fmt.Errorf("%q was not found", id)
				}

				return fmt.Errorf("retrieving %q: %+v", id, err)
			}

			if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil || iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections == nil {
				return fmt.Errorf("%q was not found", id)
			}

			for i, endpoint := range pointer.From(iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections) {
				if strings.EqualFold(pointer.From(endpoint.Name), id.EndpointName) {
					authenticationType := devices.AuthenticationType(state.AuthenticationType)

					if authenticationType == devices.AuthenticationTypeKeyBased {
						if state.PrimaryKey == "" || state.SecondaryKey == "" {
							return fmt.Errorf("`primary_key` and `secondary_key` must be specified when `authentication_type` is `keyBased`")
						}

						endpoint.PrimaryKey = pointer.To(state.PrimaryKey)
						endpoint.SecondaryKey = pointer.To(state.SecondaryKey)
						endpoint.Identity = nil
					} else {
						if state.PrimaryKey != "" || state.SecondaryKey != "" {
							return fmt.Errorf("`primary_key` or `secondary_key` cannot be specified when `authentication_type` is `identityBased`")
						}

						endpoint.PrimaryKey = nil
						endpoint.SecondaryKey = nil
						if state.IdentityId != "" {
							endpoint.Identity = &devices.ManagedIdentity{
								UserAssignedIdentity: pointer.To(state.IdentityId),
							}
						} else {
							endpoint.Identity = nil
						}
					}

					if metadata.ResourceData.HasChange("authentication_type") {
						endpoint.AuthenticationType = authenticationType
					}

					if metadata.ResourceData.HasChange("partition_key_name") {
						if state.PartitionKeyName == "" {
							endpoint.PartitionKeyName = nil
						} else {
							endpoint.PartitionKeyName = pointer.To(state.PartitionKeyName)
						}
					}

					if metadata.ResourceData.HasChange("partition_key_template") {
						if state.PartitionKeyTemplate == "" {
							endpoint.PartitionKeyTemplate = nil
						} else {
							endpoint.PartitionKeyTemplate = pointer.To(state.PartitionKeyTemplate)
						}
					}

					// As `subscription_id` is `O+C`, `HasChange()` can't detect the change when it isn't specified. And `subscription_id` always needs to be set to the subscription ID used in the provider block when it isn't specified. So, `HasChange()` is not needed.
					// To align with the previous TF behaviour, `subscription_id` needs to be set with the provider's subscription Id when it isn't specified in the tf config, otherwise TF behaviour is different than before and it may block the existing users
					// From the business perspective, the raw config handling is only meant for the case that the user has an CosmosDB Account whose Endpoint's subscription is not the provider's one. Then the user wants to reset it to the provider's one by unset the subscription_id
					// From the TF code perspective, given `Computed: true` is enabled, TF would always get the value from the last apply when this property isn't set in the tf config. So `d.GetRawConfig()` is required to determine if it's set in the tf config
					if v := metadata.ResourceData.GetRawConfig().AsValueMap()["subscription_id"]; v.IsNull() {
						endpoint.SubscriptionID = pointer.To(subscriptionId)
					} else {
						endpoint.SubscriptionID = pointer.To(state.SubscriptionId)
					}

					(*iothub.Properties.Routing.Endpoints.CosmosDBSQLCollections)[i] = endpoint

					future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
					if err != nil {
						return fmt.Errorf("updating %s: %+v", id, err)
					}

					if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
						return fmt.Errorf("waiting for the completion of the update of %s: %+v", id, err)
					}

					return nil
				}
			}

			return fmt.Errorf("%q was not found", id)
		},
		Timeout: 30 * time.Minute,
	}
}
