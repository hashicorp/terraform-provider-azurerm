// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package fluid_relay_server

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelayservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) (err error) {
			client := metadata.Client.FluidRelay.FluidRelayServers

			var model ServerModel
			if err = metadata.Decode(&model); err != nil {
				return err
			}

			id := fluidrelayservers.NewFluidRelayServerID(metadata.Client.Account.SubscriptionId, model.ResourceGroup, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if !response.WasNotFound(existing.HttpResponse) {
					if err != nil {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}

					return metadata.ResourceRequiresImport(s.ResourceType(), id)
				}
			}

			payload := fluidrelayservers.FluidRelayServer{
				Location:   location.Normalize(model.Location),
				Properties: &fluidrelayservers.FluidRelayServerProperties{},
				Tags:       &model.Tags,
			}
			payload.Identity, err = identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding user identities: %+v", err)
			}

			if model.StorageSKU != "" {
				payload.Properties.Storagesku = pointer.ToEnum[fluidrelayservers.StorageSKU](model.StorageSKU)
			}

			if customerManagedKey := expandFluidRelayServerCustomerManagedKey(model.CustomerManagedKey); customerManagedKey != nil {
				payload.Properties.Encryption = customerManagedKey
			}

			if _, err = client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %v err: %+v", id, err)
			}
			metadata.SetID(id)

			return nil
		},
	}
}

func expandFluidRelayServerCustomerManagedKey(input []CustomerManagedKey) *fluidrelayservers.EncryptionProperties {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	return &fluidrelayservers.EncryptionProperties{
		CustomerManagedKeyEncryption: &fluidrelayservers.CustomerManagedKeyEncryptionProperties{
			KeyEncryptionKeyURL: pointer.To(v.KeyVaultKeyID),
			KeyEncryptionKeyIdentity: &fluidrelayservers.CustomerManagedKeyEncryptionPropertiesKeyEncryptionKeyIdentity{
				IdentityType:                   pointer.To(fluidrelayservers.CmkIdentityTypeUserAssigned),
				UserAssignedIdentityResourceId: pointer.To(v.UserAssignedIdentityId),
			},
		},
	}
}
