// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_token_password

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	containerregistry_v2025_11_01 "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/tokens"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient
			id, err := parse.ContainerRegistryTokenPasswordID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			tokenId := tokens.NewTokenID(id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TokenName)

			pwds, notFound, err := r.readPassword(ctx, client, tokenId)
			if notFound || (pwds != nil && len(pwds) == 0) {
				return metadata.MarkAsGone(id)
			}
			if err != nil {
				return err
			}
			// ACR token with no password returns a empty array for ".password"
			if len(pwds) == 0 {
				return metadata.MarkAsGone(id)
			}

			var state ContainerRegistryTokenPasswordModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding from state %+v", err)
			}
			existingPasswords := map[string]ContainerRegistryTokenPassword{}
			if len(state.Password1) == 1 {
				existingPasswords["password1"] = state.Password1[0]
			}
			if len(state.Password2) == 1 {
				existingPasswords[string(tokens.TokenPasswordNamePasswordTwo)] = state.Password2[0]
			}

			model := ContainerRegistryTokenPasswordModel{
				TokenId: tokens.NewTokenID(id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TokenName).ID(),
			}
			for _, pwd := range pwds {
				name := string(*pwd.Name)
				v := ContainerRegistryTokenPassword{}
				if pwd.Expiry != nil {
					v.Expiry = *pwd.Expiry
				}
				// The value is not returned from the API, hence setting it based on state.
				if e, ok := existingPasswords[name]; ok {
					v.Value = e.Value
				}
				switch name {
				case "password1":
					model.Password1 = []ContainerRegistryTokenPassword{v}
				case string(tokens.TokenPasswordNamePasswordTwo):
					model.Password2 = []ContainerRegistryTokenPassword{v}
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r Resource) flattenContainerRegistryTokenPassword(input *[]tokens.TokenPassword) (password1, password2 []ContainerRegistryTokenPassword) {
	if input == nil {
		return []ContainerRegistryTokenPassword{}, []ContainerRegistryTokenPassword{}
	}

	for _, e := range *input {
		password := ContainerRegistryTokenPassword{}
		if e.Expiry != nil {
			password.Expiry = *e.Expiry
		}
		if e.Value != nil {
			password.Value = *e.Value
		}
		if e.Name != nil {
			switch *e.Name {
			case "password1":
				password1 = []ContainerRegistryTokenPassword{password}
			case tokens.TokenPasswordNamePasswordTwo:
				password2 = []ContainerRegistryTokenPassword{password}
			}
		}
	}
	return
}

type ContainerRegistryTokenPassword struct {
	Expiry string `tfschema:"expiry"`
	Value  string `tfschema:"value"`
}

func (r Resource) readPassword(ctx context.Context, client *containerregistry_v2025_11_01.Client, id tokens.TokenId) ([]tokens.TokenPassword, bool, error) {
	existing, err := client.Tokens.Get(ctx, id)
	if err != nil {
		return nil, response.WasNotFound(existing.HttpResponse), fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if existing.Model == nil {
		return nil, false, fmt.Errorf("checking for presence of existing %s: model is nil", id)
	}

	if existing.Model.Properties == nil {
		return nil, false, fmt.Errorf("checking for presence of existing %s: properties is nil", id)
	}

	if existing.Model.Properties.Credentials == nil {
		return nil, false, fmt.Errorf("checking for presence of existing %s: credentials is nil", id)
	}

	passwords := existing.Model.Properties.Credentials.Passwords
	if passwords == nil {
		return nil, false, fmt.Errorf("checking for presence of existing %s: passwords is nil", id)
	}

	return *passwords, false, nil
}

func (r Resource) generatePassword(ctx context.Context, clients client.Client, id tokens.TokenId, passwords []tokens.TokenPassword) ([]tokens.TokenPassword, error) {
	var genPasswords []tokens.TokenPassword

	existingPasswords, _, err := r.readPassword(ctx, clients.ContainerRegistryClient, id)
	if err != nil {
		return nil, fmt.Errorf("reading existing passwords: %+v", err)
	}

	// The token password API has the following behaviour:
	// - To remove password, one uses the PATCH of the ACR token endpoint
	// - To add password, one uses the POST of the ACR's generate credential endpoint
	// Hence we'd have to check whether there is any password to clean up before we try to update/create passwords.
	if len(existingPasswords) > len(passwords) {
		param := tokens.TokenUpdateParameters{
			Properties: &tokens.TokenUpdateProperties{
				Credentials: &tokens.TokenCredentialsProperties{
					Passwords: &passwords,
				},
			},
		}
		if err := clients.ContainerRegistryClient.Tokens.UpdateThenPoll(ctx, id, param); err != nil {
			return nil, fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	// Iterate and generate password planned to be created.
PasswordGenLoop:
	for idx, password := range passwords {
		// For each password specified in the config, check whether it is changed compared to its remote state (i.e. is the `expiry` changed).
		// If nothing is changed, we should skip it - rather than generating a new password.
		for _, pwd := range existingPasswords {
			if pwd.Name == password.Name {
				if (pwd.Expiry == nil) != (password.Expiry == nil) {
					break
				}

				pwdExpiryTime, err := pwd.GetExpiryAsTime()
				if err != nil {
					return nil, fmt.Errorf("unable to get expiry time for %s: %+v", string(*pwd.Name), err)
				}
				passwordExpiryTime, err := password.GetExpiryAsTime()
				if err != nil {
					return nil, fmt.Errorf("unable to get expiry time for %s: %+v", string(*password.Name), err)
				}
				if pwd.Expiry == nil || pwdExpiryTime.Equal(*passwordExpiryTime) {
					genPasswords = append(genPasswords, password)
					continue PasswordGenLoop
				}
				break
			}
		}

		param := registries.GenerateCredentialsParameters{
			TokenId: pointer.To(id.ID()),
			Expiry:  password.Expiry,
			Name:    (*registries.TokenPasswordName)(password.Name),
		}

		registryId := registries.NewRegistryID(id.SubscriptionId, id.ResourceGroupName, id.RegistryName)

		result, err := clients.ContainerRegistryClient.Registries.GenerateCredentials(ctx, registryId, param)
		if err != nil {
			return nil, fmt.Errorf("generating password credential %s: %v", string(*password.Name), err)
		}

		if err := result.Poller.PollUntilDone(ctx); err != nil {
			return nil, fmt.Errorf("polling generation of password credential %s: %v", string(*password.Name), err)
		}

		var res registries.GenerateCredentialsResult
		if err := json.NewDecoder(result.HttpResponse.Body).Decode(&res); err != nil {
			return nil, fmt.Errorf("decoding generated password credentials: %v", err)
		}

		value := ""
		if res.Passwords != nil && len(*res.Passwords) > idx && (*res.Passwords)[idx].Value != nil {
			value = *(*res.Passwords)[idx].Value
		}

		genPasswords = append(genPasswords, tokens.TokenPassword{
			Expiry: password.Expiry,
			Name:   password.Name,
			Value:  &value,
		})
	}
	return genPasswords, nil
}
