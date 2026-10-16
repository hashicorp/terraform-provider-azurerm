// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_credential_set

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/credentialsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.CredentialSetsClient
			id, err := credentialsets.ParseCredentialSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			param := credentialsets.CredentialSetUpdateParameters{}

			var model ContainerRegistryCredentialSetModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			properties := credentialsets.CredentialSetUpdateProperties{}

			if metadata.ResourceData.HasChange("authentication_credentials") {
				properties.AuthCredentials = expandAuthCredentials(model.AuthenticationCredential)
			}

			param.Properties = &properties

			if metadata.ResourceData.HasChange("identity") {
				identityExpanded, err := identity.ExpandSystemAssignedFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				param.Identity = identityExpanded
			}

			if err := client.UpdateThenPoll(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}
