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

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient.Tokens

			id, err := parse.ContainerRegistryTokenPasswordID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			tokenId := tokens.NewTokenID(id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TokenName)

			locks.ByID(tokenId.ID())
			defer locks.UnlockByID(tokenId.ID())

			param := tokens.TokenUpdateParameters{
				Properties: &tokens.TokenUpdateProperties{
					Credentials: &tokens.TokenCredentialsProperties{
						Passwords: &[]tokens.TokenPassword{},
					},
				},
			}

			if err := client.UpdateThenPoll(ctx, tokenId, param); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
