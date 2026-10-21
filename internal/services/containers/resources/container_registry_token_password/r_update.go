// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_token_password

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/tokens"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.ContainerRegistryTokenPasswordID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			tokenId := tokens.NewTokenID(id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TokenName)

			var plan ContainerRegistryTokenPasswordModel
			if err := metadata.Decode(&plan); err != nil {
				return err
			}

			passwords, err := r.expandContainerRegistryTokenPassword(plan)
			if err != nil {
				return fmt.Errorf("expanding `password`: %v", err)
			}

			locks.ByID(tokenId.ID())
			defer locks.UnlockByID(tokenId.ID())

			genPasswords, err := r.generatePassword(ctx, *metadata.Client.Containers, tokenId, *passwords)
			if err != nil {
				return err
			}

			// The password is only known right after it is generated, therefore setting it to the resource data here.
			password1, password2 := r.flattenContainerRegistryTokenPassword(&genPasswords)
			plan.Password1 = password1
			plan.Password2 = password2
			if err := metadata.Encode(&plan); err != nil {
				return fmt.Errorf("encoding model and store into state: %+v", err)
			}

			return nil
		},
	}
}
