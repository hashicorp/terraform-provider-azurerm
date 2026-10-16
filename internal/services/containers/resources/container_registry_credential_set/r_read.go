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

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.CredentialSetsClient
			id, err := credentialsets.ParseCredentialSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			registryId := registries.NewRegistryID(id.SubscriptionId, id.ResourceGroupName, id.RegistryName)

			var config ContainerRegistryCredentialSetModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			config.Name = id.CredentialSetName
			config.ContainerRegistryId = registryId.ID()

			if model := resp.Model; model != nil {
				config.Identity = identity.FlattenSystemAssignedToModel(model.Identity)
				if props := model.Properties; props != nil {
					config.LoginServer = pointer.From(props.LoginServer)
					config.AuthenticationCredential = flattenAuthCredentials(props.AuthCredentials)
				}
			}
			return metadata.Encode(&config)
		},
	}
}

func flattenAuthCredentials(input *[]credentialsets.AuthCredential) []AuthenticationCredential {
	output := make([]AuthenticationCredential, 0)
	if input == nil {
		return output
	}
	for _, v := range *input {
		output = append(output, AuthenticationCredential{
			UsernameSecretId: pointer.From(v.UsernameSecretIdentifier),
			PasswordSecretId: pointer.From(v.PasswordSecretIdentifier),
		})
	}
	return output
}

type AuthenticationCredential struct {
	UsernameSecretId string `tfschema:"username_secret_id"`
	PasswordSecretId string `tfschema:"password_secret_id"`
}
