// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package fluid_relay_server

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26/fluidrelayservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (s Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, meta sdk.ResourceMetaData) error {
			client := meta.Client.FluidRelay.FluidRelayServers

			id, err := fluidrelayservers.ParseFluidRelayServerID(meta.ResourceData.Id())
			if err != nil {
				return err
			}

			server, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(server.HttpResponse) {
					return meta.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			keys, err := client.ListKeys(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving keys for %s: %+v", *id, err)
			}

			output := &ServerModel{
				Name:          id.FluidRelayServerName,
				ResourceGroup: id.ResourceGroup,
			}
			if model := server.Model; model != nil {
				output.Location = location.Normalize(model.Location)

				if err = output.flattenIdentity(model.Identity); err != nil {
					return fmt.Errorf("flattening `identity`: %v", err)
				}
				if server.Model.Tags != nil {
					output.Tags = *server.Model.Tags
				}
				if prop := model.Properties; prop != nil {
					output.CustomerManagedKey, err = flattenFluidRelayServerCustomerManagedKey(prop.Encryption)
					if err != nil {
						return fmt.Errorf("flattening `customer_managed_key`: %v", err)
					}

					if prop.FrsTenantId != nil {
						output.FrsTenantId = *prop.FrsTenantId
					}
					if points := prop.FluidRelayEndpoints; points != nil {
						if points.OrdererEndpoints != nil {
							output.OrdererEndpoints = *points.OrdererEndpoints
						}
						if points.StorageEndpoints != nil {
							output.StorageEndpoints = *points.StorageEndpoints
						}

						if points.ServiceEndpoints != nil {
							output.ServiceEndpoints = *points.ServiceEndpoints
						}
					}
				}
				if val, ok := meta.ResourceData.GetOk("storage_sku"); ok {
					output.StorageSKU = val.(string)
				}
			}

			if model := keys.Model; model != nil {
				output.PrimaryKey = pointer.From(model.Key1)
				output.SecondaryKey = pointer.From(model.Key2)
			}

			return meta.Encode(output)
		},
	}
}

func (s *ServerModel) flattenIdentity(input *identity.SystemAndUserAssignedMap) error {
	if input == nil {
		return nil
	}
	config := identity.SystemAndUserAssignedMap{
		Type:        input.Type,
		PrincipalId: input.PrincipalId,
		TenantId:    input.TenantId,
		IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
	}
	for k, v := range input.IdentityIds {
		config.IdentityIds[k] = identity.UserAssignedIdentityDetails{
			ClientId:    v.ClientId,
			PrincipalId: v.PrincipalId,
		}
	}
	model, err := identity.FlattenSystemAndUserAssignedMapToModel(&config)
	if err != nil {
		return err
	}
	s.Identity = *model
	return nil
}

func flattenFluidRelayServerCustomerManagedKey(input *fluidrelayservers.EncryptionProperties) ([]CustomerManagedKey, error) {
	if input == nil || input.CustomerManagedKeyEncryption == nil {
		return []CustomerManagedKey{}, nil
	}

	customerManagedKey := CustomerManagedKey{}

	if input.CustomerManagedKeyEncryption.KeyEncryptionKeyURL != nil {
		if v := pointer.From(input.CustomerManagedKeyEncryption.KeyEncryptionKeyURL); v != "" {
			id, err := keyvault.ParseNestedItemID(v, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
			if err != nil {
				return nil, err
			}
			customerManagedKey.KeyVaultKeyID = id.ID()
		}
	}
	if input.CustomerManagedKeyEncryption.KeyEncryptionKeyIdentity != nil && input.CustomerManagedKeyEncryption.KeyEncryptionKeyIdentity.UserAssignedIdentityResourceId != nil {
		if v := pointer.From(input.CustomerManagedKeyEncryption.KeyEncryptionKeyIdentity.UserAssignedIdentityResourceId); v != "" {
			id, err := commonids.ParseUserAssignedIdentityID(v)
			if err != nil {
				return nil, err
			}
			customerManagedKey.UserAssignedIdentityId = id.ID()
		}
	}

	return []CustomerManagedKey{customerManagedKey}, nil
}

type CustomerManagedKey struct {
	UserAssignedIdentityId string `tfschema:"user_assigned_identity_id"`
	KeyVaultKeyID          string `tfschema:"key_vault_key_id"`
}
