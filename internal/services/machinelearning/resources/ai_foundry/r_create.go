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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AIFoundryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := workspaces.NewWorkspaceID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError("azurerm_ai_foundry", id.ID())
				}
			}

			storageAccountId, err := commonids.ParseStorageAccountID(model.StorageAccountId)
			if err != nil {
				return err
			}

			keyVaultId, err := commonids.ParseKeyVaultID(model.KeyVaultId)
			if err != nil {
				return err
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]any))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			payload := workspaces.Workspace{
				Name:     pointer.To(id.WorkspaceName),
				Location: pointer.To(location.Normalize(model.Location)),
				Identity: expandedIdentity,
				Tags:     tags.Expand(model.Tags),
				Kind:     pointer.To("Hub"),
				Properties: &workspaces.WorkspaceProperties{
					KeyVault:            pointer.To(keyVaultId.ID()),
					PublicNetworkAccess: pointer.ToEnum[workspaces.PublicNetworkAccess](model.PublicNetworkAccess),
					StorageAccount:      pointer.To(storageAccountId.ID()),
				},
			}

			if model.ApplicationInsightsId != "" {
				applicationInsightsId, err := componentsapis.ParseComponentID(model.ApplicationInsightsId)
				if err != nil {
					return err
				}
				payload.Properties.ApplicationInsights = pointer.To(applicationInsightsId.ID())
			}

			if model.ContainerRegistryId != "" {
				containerRegistryId, err := registries.ParseRegistryID(model.ContainerRegistryId)
				if err != nil {
					return err
				}
				payload.Properties.ContainerRegistry = pointer.To(containerRegistryId.ID())
			}

			if model.Description != "" {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if model.FriendlyName != "" {
				payload.Properties.FriendlyName = pointer.To(model.FriendlyName)
			}

			if model.HighBusinessImpactEnabled {
				payload.Properties.HbiWorkspace = pointer.To(model.HighBusinessImpactEnabled)
			}

			if model.PrimaryUserAssignedIdentity != "" {
				userAssignedId, err := commonids.ParseUserAssignedIdentityID(model.PrimaryUserAssignedIdentity)
				if err != nil {
					return err
				}
				payload.Properties.PrimaryUserAssignedIdentity = pointer.To(userAssignedId.ID())
			}

			if len(model.Encryption) > 0 {
				payload.Properties.Encryption = expandEncryption(model.Encryption)
			}

			if len(model.ManagedNetwork) > 0 {
				payload.Properties.ManagedNetwork = expandManagedNetwork(model.ManagedNetwork)
			}

			if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandEncryption(input []Encryption) *workspaces.EncryptionProperty {
	encryption := input[0]
	out := workspaces.EncryptionProperty{
		Identity: &workspaces.IdentityForCmk{},
		KeyVaultProperties: workspaces.EncryptionKeyVaultProperties{
			KeyVaultArmId: encryption.KeyVaultID,
			KeyIdentifier: encryption.KeyID,
		},
		Status: workspaces.EncryptionStatusEnabled,
	}

	if encryption.IdentityClientID != "" {
		out.Identity.UserAssignedIdentity = pointer.To(encryption.IdentityClientID)
	}

	return &out
}

func expandManagedNetwork(input []ManagedNetwork) *workspaces.ManagedNetworkSettings {
	network := input[0]

	return &workspaces.ManagedNetworkSettings{
		IsolationMode: pointer.ToEnum[workspaces.IsolationMode](network.IsolationMode),
	}
}
