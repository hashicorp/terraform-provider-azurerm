// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	containterregistry_v2021_08_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/tokens"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerRegistryTokenPasswordResource struct{}

var _ sdk.ResourceWithUpdate = ContainerRegistryTokenPasswordResource{}

type ContainerRegistryTokenPasswordModel struct {
	TokenId   string                           `tfschema:"container_registry_token_id"`
	Password1 []ContainerRegistryTokenPassword `tfschema:"password1"`
	Password2 []ContainerRegistryTokenPassword `tfschema:"password2"`
}

type ContainerRegistryTokenPassword struct {
	Expiry string `tfschema:"expiry"`
	Value  string `tfschema:"value"`
}

func (r ContainerRegistryTokenPasswordResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"container_registry_token_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: tokens.ValidateTokenID,
		},

		"password1": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expiry": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// TODO: Remove the force new and add customize diff to SetNewComputed on the `value` once https://github.com/hashicorp/terraform-plugin-sdk/issues/459 is addressed.
						ForceNew:         true,
						ValidateFunc:     validation.IsRFC3339Time,
						DiffSuppressFunc: suppress.RFC3339Time,
					},

					"value": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},

		"password2": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expiry": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// TODO: Remove the force new and add customize diff to SetNewComputed on the `value` once https://github.com/hashicorp/terraform-plugin-sdk/issues/459 is addressed.
						ForceNew:         true,
						ValidateFunc:     validation.IsRFC3339Time,
						DiffSuppressFunc: suppress.RFC3339Time,
					},

					"value": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},
	}
}

func (r ContainerRegistryTokenPasswordResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ContainerRegistryTokenPasswordResource) ResourceType() string {
	return "azurerm_container_registry_token_password"
}

func (r ContainerRegistryTokenPasswordResource) ModelObject() interface{} {
	return &ContainerRegistryTokenPasswordModel{}
}

func (r ContainerRegistryTokenPasswordResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ContainerRegistryTokenPasswordID
}

func (r ContainerRegistryTokenPasswordResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient_v2023_06_01_preview
			var plan ContainerRegistryTokenPasswordModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			tokenId, err := tokens.ParseTokenID(plan.TokenId)
			if err != nil {
				return err
			}

			id := parse.NewContainerRegistryTokenPasswordID(tokenId.SubscriptionId, tokenId.ResourceGroupName, tokenId.RegistryName, tokenId.TokenName, "password")

			pwds, err := r.readPassword(ctx, client, *tokenId)
			if err != nil {
				return err
			}
			// ACR token with no password returns a empty array for ".password"
			if len(pwds) != 0 {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			passwords, err := r.expandContainerRegistryTokenPassword(plan)
			if err != nil {
				return fmt.Errorf("expanding `password`: %v", err)
			}

			locks.ByID(tokenId.ID())
			defer locks.UnlockByID(tokenId.ID())

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

func (r ContainerRegistryTokenPasswordResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient_v2023_06_01_preview
			id, err := parse.ContainerRegistryTokenPasswordID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			tokenId := tokens.NewTokenID(id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TokenName)

			pwds, err := r.readPassword(ctx, client, tokenId)
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

func (r ContainerRegistryTokenPasswordResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient_v2023_06_01_preview.Tokens

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

func (r ContainerRegistryTokenPasswordResource) Update() sdk.ResourceFunc {
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

func (r ContainerRegistryTokenPasswordResource) expandContainerRegistryTokenPassword(plan ContainerRegistryTokenPasswordModel) (*[]tokens.TokenPassword, error) {
	result := make([]tokens.TokenPassword, 0)

	expandPassword := func(name string, password []ContainerRegistryTokenPassword) (*tokens.TokenPassword, error) {
		if len(password) == 1 {
			password := password[0]
			ret := &tokens.TokenPassword{
				Name:  pointer.To(tokens.TokenPasswordName(name)),
				Value: utils.String(password.Value),
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

func (r ContainerRegistryTokenPasswordResource) flattenContainerRegistryTokenPassword(input *[]tokens.TokenPassword) (password1, password2 []ContainerRegistryTokenPassword) {
	if input == nil {
		return nil, nil
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

func (r ContainerRegistryTokenPasswordResource) readPassword(ctx context.Context, client *containterregistry_v2021_08_01_preview.Client, id tokens.TokenId) ([]tokens.TokenPassword, error) {
	existing, err := client.Tokens.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if existing.Model == nil {
		return nil, fmt.Errorf("checking for presence of existing %s: model is nil", id)
	}

	if existing.Model.Properties == nil {
		return nil, fmt.Errorf("checking for presence of existing %s: properties is nil", id)
	}

	if existing.Model.Properties.Credentials == nil {
		return nil, fmt.Errorf("checking for presence of existing %s: credentials is nil", id)
	}

	passwords := existing.Model.Properties.Credentials.Passwords
	if passwords == nil {
		return nil, fmt.Errorf("checking for presence of existing %s: passwords is nil", id)
	}

	return *passwords, nil
}

func (r ContainerRegistryTokenPasswordResource) generatePassword(ctx context.Context, clients client.Client, id tokens.TokenId, passwords []tokens.TokenPassword) ([]tokens.TokenPassword, error) {
	var genPasswords []tokens.TokenPassword

	existingPasswords, err := r.readPassword(ctx, clients.ContainerRegistryClient_v2023_06_01_preview, id)
	if err != nil {
		return nil, fmt.Errorf("reading existing passwords: %+v", err)
	}

	// The token password API has the following behavior:
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
		if err := clients.ContainerRegistryClient_v2023_06_01_preview.Tokens.UpdateThenPoll(ctx, id, param); err != nil {
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
			TokenId: utils.String(id.ID()),
			Expiry:  password.Expiry,
			Name:    (*registries.TokenPasswordName)(password.Name),
		}

		registryId := registries.NewRegistryID(id.SubscriptionId, id.ResourceGroupName, id.RegistryName)

		result, err := clients.ContainerRegistryClient_v2023_06_01_preview.Registries.GenerateCredentials(ctx, registryId, param)
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
