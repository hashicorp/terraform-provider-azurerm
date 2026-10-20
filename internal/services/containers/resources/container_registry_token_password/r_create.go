// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_registry_token_password

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/tokens"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient
			var plan ContainerRegistryTokenPasswordModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			tokenId, err := tokens.ParseTokenID(plan.TokenId)
			if err != nil {
				return err
			}

			id := parse.NewContainerRegistryTokenPasswordID(tokenId.SubscriptionId, tokenId.ResourceGroupName, tokenId.RegistryName, tokenId.TokenName, "password")

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				pwds, _, err := r.readPassword(ctx, client, *tokenId)
				if err != nil {
					return err
				}
				// ACR token with no password returns an empty array for ".password"
				if len(pwds) != 0 {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			passwords, err := r.expandContainerRegistryTokenPassword(plan)
			if err != nil {
				return fmt.Errorf("expanding `password`: %v", err)
			}

			locks.ByID(tokenId.ID())
			defer locks.UnlockByID(tokenId.ID())

			// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
			genPasswords, err := r.generatePassword(ctx, *metadata.Client.Containers, *tokenId, *passwords)
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

			metadata.SetID(id)
			return nil
		},
	}
}

func (r Resource) expandContainerRegistryTokenPassword(plan ContainerRegistryTokenPasswordModel) (*[]tokens.TokenPassword, error) {
	result := make([]tokens.TokenPassword, 0)

	expandPassword := func(name string, password []ContainerRegistryTokenPassword) (*tokens.TokenPassword, error) {
		if len(password) == 1 {
			password := password[0]
			ret := &tokens.TokenPassword{
				Name:  pointer.ToEnum[tokens.TokenPasswordName](name),
				Value: pointer.To(password.Value),
			}
			if v := password.Expiry; v != "" {
				t, err := time.Parse(time.RFC3339, v)
				if err != nil {
					return nil, err
				}
				ret.Expiry = pointer.To(date.Time{Time: t}.String())
			}
			return ret, nil
		}
		return nil, nil
	}
	// TODO: Use below SDK enum once the following issue is resolved: https://github.com/Azure/azure-rest-api-specs/issues/18339
	// tokens.PasswordNamePassword
	v, err := expandPassword("password1", plan.Password1)
	if err != nil {
		return nil, err
	}
	if v != nil {
		result = append(result, *v)
	}

	v, err = expandPassword(string(tokens.TokenPasswordNamePasswordTwo), plan.Password2)
	if err != nil {
		return nil, err
	}
	if v != nil {
		result = append(result, *v)
	}
	return &result, nil
}
