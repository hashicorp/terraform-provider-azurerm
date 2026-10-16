// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_credential_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/credentialsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.CredentialSetsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ContainerRegistryCredentialSetModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			registryId, err := registries.ParseRegistryID(config.ContainerRegistryId)
			if err != nil {
				return err
			}

			id := credentialsets.NewCredentialSetID(
				subscriptionId,
				registryId.ResourceGroupName,
				registryId.RegistryName,
				config.Name,
			)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			identityExpanded, err := identity.ExpandSystemAssignedFromModel(config.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			param := credentialsets.CredentialSet{
				Name: pointer.To(id.CredentialSetName),
				Properties: &credentialsets.CredentialSetProperties{
					LoginServer:     pointer.To(config.LoginServer),
					AuthCredentials: expandAuthCredentials(config.AuthenticationCredential),
				},
				Identity: identityExpanded,
			}

			if err := client.CreateCallbackThenPoll(ctx, id, param, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandAuthCredentials(input []AuthenticationCredential) *[]credentialsets.AuthCredential {
	output := make([]credentialsets.AuthCredential, 0)
	if len(input) == 0 {
		return &output
	}
	for _, v := range input {
		output = append(output, credentialsets.AuthCredential{
			Name:                     pointer.To(credentialsets.CredentialNameCredentialOne),
			UsernameSecretIdentifier: pointer.To(v.UsernameSecretId),
			PasswordSecretIdentifier: pointer.To(v.PasswordSecretId),
		})
	}
	return &output
}
