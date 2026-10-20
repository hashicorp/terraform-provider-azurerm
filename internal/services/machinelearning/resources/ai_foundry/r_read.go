// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ai_foundry

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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces

			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			hub := AIFoundryModel{
				Name:              id.WorkspaceName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				hub.Location = location.NormalizeNilable(model.Location)

				flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				hub.Identity = flattenedIdentity
				hub.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					if v := pointer.From(props.ApplicationInsights); v != "" {
						applicationInsightsId, err := componentsapis.ParseComponentIDInsensitively(v)
						if err != nil {
							return err
						}
						hub.ApplicationInsightsId = applicationInsightsId.ID()
					}

					if v := pointer.From(props.ContainerRegistry); v != "" {
						containerRegistryId, err := registries.ParseRegistryID(v)
						if err != nil {
							return err
						}
						hub.ContainerRegistryId = containerRegistryId.ID()
					}

					storageAccountId, err := commonids.ParseStorageAccountID(*props.StorageAccount)
					if err != nil {
						return err
					}
					hub.StorageAccountId = storageAccountId.ID()

					if props.KeyVault != nil {
						keyVaultId, err := commonids.ParseKeyVaultID(*props.KeyVault)
						if err != nil {
							return err
						}
						hub.KeyVaultId = keyVaultId.ID()
					}

					hub.Description = pointer.From(props.Description)
					hub.FriendlyName = pointer.From(props.FriendlyName)
					hub.HighBusinessImpactEnabled = pointer.From(props.HbiWorkspace)
					hub.PublicNetworkAccess = string(*props.PublicNetworkAccess)
					hub.DiscoveryUrl = pointer.From(props.DiscoveryURL)
					hub.WorkspaceId = pointer.From(props.WorkspaceId)
					hub.ManagedNetwork = flattenManagedNetwork(props.ManagedNetwork)

					if v := pointer.From(props.PrimaryUserAssignedIdentity); v != "" {
						userAssignedId, err := commonids.ParseUserAssignedIdentityID(v)
						if err != nil {
							return err
						}
						hub.PrimaryUserAssignedIdentity = userAssignedId.ID()
					}

					encryption, err := flattenEncryption(props.Encryption)
					if err != nil {
						return fmt.Errorf("flattening `encryption`: %+v", err)
					}
					hub.Encryption = encryption
				}
			}

			return metadata.Encode(&hub)
		},
	}
}

func flattenEncryption(input *workspaces.EncryptionProperty) ([]Encryption, error) {
	out := make([]Encryption, 0)

	if input == nil || input.Status != workspaces.EncryptionStatusEnabled {
		return out, nil
	}

	encryption := Encryption{}
	if v := input.KeyVaultProperties.KeyVaultArmId; v != "" {
		keyVaultId, err := commonids.ParseKeyVaultID(v)
		if err != nil {
			return nil, err
		}
		encryption.KeyVaultID = keyVaultId.ID()
	}
	if v := input.KeyVaultProperties.KeyIdentifier; v != "" {
		keyId, err := keyvault.ParseNestedItemID(v, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
		if err != nil {
			return nil, err
		}
		encryption.KeyID = keyId.ID()
	}

	if input.Identity != nil && input.Identity.UserAssignedIdentity != nil {
		userAssignedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*input.Identity.UserAssignedIdentity)
		if err != nil {
			return nil, err
		}
		encryption.IdentityClientID = userAssignedId.ID()
	}

	return append(out, encryption), nil
}

func flattenManagedNetwork(input *workspaces.ManagedNetworkSettings) []ManagedNetwork {
	out := make([]ManagedNetwork, 0)
	if input == nil {
		return out
	}

	return append(out, ManagedNetwork{
		IsolationMode: pointer.FromEnum(input.IsolationMode),
	})
}

type ManagedNetwork struct {
	IsolationMode string `tfschema:"isolation_mode"`
}

type Encryption struct {
	IdentityClientID string `tfschema:"user_assigned_identity_id"`
	KeyVaultID       string `tfschema:"key_vault_id"`
	KeyID            string `tfschema:"key_id"`
}
