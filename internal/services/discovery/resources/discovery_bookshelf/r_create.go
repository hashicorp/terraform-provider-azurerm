// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package discovery_bookshelf

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelves"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Discovery.Bookshelves
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model DiscoveryBookshelfModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := bookshelves.NewBookshelfID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			customerManagedKeys := bookshelves.CustomerManagedKeysDisabled
			var keyVaultProperties *bookshelves.BookshelfKeyVaultProperties
			if len(model.CustomerManagedKey) > 0 {
				customerManagedKeys = bookshelves.CustomerManagedKeysEnabled

				identityId, err := commonids.ParseUserAssignedIdentityID(model.CustomerManagedKey[0].UserAssignedIdentityId)
				if err != nil {
					return err
				}

				identity, err := metadata.Client.ManagedIdentity.V20241130.Identities.UserAssignedIdentitiesGet(ctx, *identityId)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", *identityId, err)
				}

				if identity.Model == nil || identity.Model.Properties == nil || identity.Model.Properties.ClientId == nil {
					return fmt.Errorf("retrieving %s: `clientId` was nil", *identityId)
				}

				keyVaultProperties, err = expandDiscoveryBookshelfCustomerManagedKey(model.CustomerManagedKey, *identity.Model.Properties.ClientId)
				if err != nil {
					return err
				}
			}

			publicNetworkAccess := bookshelves.PublicNetworkAccessDisabled
			if model.PublicNetworkAccessEnabled {
				publicNetworkAccess = bookshelves.PublicNetworkAccessEnabled
			}

			payload := bookshelves.Bookshelf{
				Location: location.Normalize(model.Location),
				Properties: &bookshelves.BookshelfProperties{
					CustomerManagedKeys:     pointer.To(customerManagedKeys),
					KeyVaultProperties:      keyVaultProperties,
					PrivateEndpointSubnetId: pointer.To(model.PrivateEndpointSubnetId),
					PublicNetworkAccess:     pointer.To(publicNetworkAccess),
					SearchSubnetId:          pointer.To(model.SearchSubnetId),
					WorkloadIdentities:      expandDiscoveryBookshelfWorkloadIdentities(model.WorkloadIdentityIds),
				},
				Tags: pointer.To(model.Tags),
			}

			if model.LogAnalyticsClusterId != "" {
				payload.Properties.LogAnalyticsClusterId = pointer.To(model.LogAnalyticsClusterId)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
