// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package discovery_bookshelf

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelves"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Discovery.Bookshelves

			id, err := bookshelves.ParseBookshelfID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DiscoveryBookshelfModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			payload := existing.Model
			props := payload.Properties
			props.ManagedOnBehalfOfConfiguration = nil

			if metadata.ResourceData.HasChange("customer_managed_key") && props.KeyVaultProperties != nil {
				keyVaultProperties, err := expandDiscoveryBookshelfCustomerManagedKey(model.CustomerManagedKey, props.KeyVaultProperties.IdentityClientId)
				if err != nil {
					return err
				}
				props.KeyVaultProperties = keyVaultProperties
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				publicNetworkAccess := bookshelves.PublicNetworkAccessDisabled
				if model.PublicNetworkAccessEnabled {
					publicNetworkAccess = bookshelves.PublicNetworkAccessEnabled
				}
				props.PublicNetworkAccess = pointer.To(publicNetworkAccess)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
