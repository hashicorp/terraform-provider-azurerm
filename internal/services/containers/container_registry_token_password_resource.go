package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2021-08-01-preview/containerregistry"
	"github.com/Azure/go-autorest/autorest/date"
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
			ValidateFunc: validate.ContainerRegistryTokenID,
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
			client := metadata.Client.Containers.TokensClient
			var plan ContainerRegistryTokenPasswordModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			tokenId, err := parse.ContainerRegistryTokenID(plan.TokenId)
			if err != nil {
				return fmt.Errorf("parsing container registry token id %q: %+v", plan.TokenId, err)
			}

			id := parse.NewContainerRegistryTokenPasswordID(tokenId.SubscriptionId, tokenId.ResourceGroup, tokenId.RegistryName, tokenId.TokenName, "password")

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
			client := metadata.Client.Containers.TokensClient
			id, err := parse.ContainerRegistryTokenPasswordID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			tokenId := parse.NewContainerRegistryTokenID(id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TokenName)

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
				existingPasswords[string(containerregistry.TokenPasswordNamePassword2)] = state.Password2[0]
			}

			model := ContainerRegistryTokenPasswordModel{
				TokenId: parse.NewContainerRegistryTokenID(id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TokenName).ID(),
			}
			for _, pwd := range pwds {
				name := string(pwd.Name)
				v := ContainerRegistryTokenPassword{}
				if pwd.Expiry != nil {
					v.Expiry = pwd.Expiry.String()
				}
				// The value is not returned from the API, hence setting it based on state.
				if e, ok := existingPasswords[name]; ok {
					v.Value = e.Value
				}
				switch name {
				case "password1":
					model.Password1 = []ContainerRegistryTokenPassword{v}
				case string(containerregistry.TokenPasswordNamePassword2):
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
			client := metadata.Client.Containers.TokensClient

			id, err := parse.ContainerRegistryTokenPasswordID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			tokenId := parse.NewContainerRegistryTokenID(id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TokenName)

			locks.ByID(tokenId.ID())
			defer locks.UnlockByID(tokenId.ID())

			param := containerregistry.TokenUpdateParameters{
				TokenUpdateProperties: &containerregistry.TokenUpdateProperties{
					Credentials: &containerregistry.TokenCredentialsProperties{
						Passwords: &[]containerregistry.TokenPassword{},
					},
				},
			}
			future, err := client.Update(ctx, id.ResourceGroup, id.RegistryName, id.TokenName, param)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for removal of %s: %+v", id, err)
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

			tokenId := parse.NewContainerRegistryTokenID(id.SubscriptionId, id.ResourceGroup, id.RegistryName, id.TokenName)

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

func (r ContainerRegistryTokenPasswordResource) expandContainerRegistryTokenPassword(plan ContainerRegistryTokenPasswordModel) (*[]containerregistry.TokenPassword, error) {
	result := make([]containerregistry.TokenPassword, 0)

	expandPassword := func(name string, password []ContainerRegistryTokenPassword) (*containerregistry.TokenPassword, error) {
		if len(password) == 1 {
			password := password[0]
			ret := &containerregistry.TokenPassword{
				Name:  containerregistry.TokenPasswordName(name),
				Value: utils.String(password.Value),
			}
			if v := password.Expiry; v != "" {
				t, err := time.Parse(time.RFC3339, v)
				if err != nil {
					return nil, err
				}
				ret.Expiry = &date.Time{Time: t}
			}
			return ret, nil
		}
		return nil, nil
	}
	// TODO: Use below SDK enum once the following issue is resolved: https://github.com/Azure/azure-rest-api-specs/issues/18339
	// containerregistry.PasswordNamePassword
	v, err := expandPassword("password1", plan.Password1)
	if err != nil {
		return nil, err
	}
	if v != nil {
		result = append(result, *v)
	}

	v, err = expandPassword(string(containerregistry.PasswordNamePassword2), plan.Password2)
	if err != nil {
		return nil, err
	}
	if v != nil {
		result = append(result, *v)
	}
	return &result, nil
}

func (r ContainerRegistryTokenPasswordResource) flattenContainerRegistryTokenPassword(input *[]containerregistry.TokenPassword) (password1, password2 []ContainerRegistryTokenPassword) {
	if input == nil {
		return nil, nil
	}

	for _, e := range *input {
		password := ContainerRegistryTokenPassword{}
		if e.Expiry != nil {
			password.Expiry = e.Expiry.String()
		}
		if e.Value != nil {
			password.Value = *e.Value
		}
		switch e.Name {
		case "password1":
			password1 = []ContainerRegistryTokenPassword{password}
		case containerregistry.TokenPasswordNamePassword2:
			password2 = []ContainerRegistryTokenPassword{password}
		}
	}
	return
}

func (r ContainerRegistryTokenPasswordResource) readPassword(ctx context.Context, client *containerregistry.TokensClient, id parse.ContainerRegistryTokenId) ([]containerregistry.TokenPassword, error) {
	existing, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.TokenName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	props := existing.TokenProperties
	if props == nil {
		return nil, fmt.Errorf("checking for presence of existing %s: unexpected nil tokenProperties", id)
	}
	cred := props.Credentials
	if cred == nil {
		return nil, fmt.Errorf("checking for presence of existing %s: unexpected nil tokenProperties.credentials", id)
	}
	pwds := cred.Passwords
	if pwds == nil {
		return nil, fmt.Errorf("checking for presence of existing %s: unexpected nil tokenProperties.credentials.passwords", id)
	}
	return *pwds, nil
}

func (r ContainerRegistryTokenPasswordResource) generatePassword(ctx context.Context, clients client.Client, id parse.ContainerRegistryTokenId, passwords []containerregistry.TokenPassword) ([]containerregistry.TokenPassword, error) {
	var genPasswords []containerregistry.TokenPassword

	existingPasswords, err := r.readPassword(ctx, clients.TokensClient, id)
	if err != nil {
		return nil, fmt.Errorf("reading existing passwords: %+v", err)
	}

	// The token password API has the following behavior:
	// - To remove password, one uses the PATCH of the ACR token endpoint
	// - To add password, one uses the POST of the ACR's generate credential endpoint
	// Hence we'd have to check whether there is any password to clean up before we try to update/create passwords.
	if len(existingPasswords) > len(passwords) {
		param := containerregistry.TokenUpdateParameters{
			TokenUpdateProperties: &containerregistry.TokenUpdateProperties{
				Credentials: &containerregistry.TokenCredentialsProperties{
					Passwords: &passwords,
				},
			},
		}
		future, err := clients.TokensClient.Update(ctx, id.ResourceGroup, id.RegistryName, id.TokenName, param)
		if err != nil {
			return nil, fmt.Errorf("deleting %s: %+v", id, err)
		}
		if err := future.WaitForCompletionRef(ctx, clients.TokensClient.Client); err != nil {
			return nil, fmt.Errorf("waiting for removal of %s: %+v", id, err)
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
				if pwd.Expiry == nil || pwd.Expiry.Time.Equal(password.Expiry.Time) {
					genPasswords = append(genPasswords, password)
					continue PasswordGenLoop
				}
				break
			}
		}

		param := containerregistry.GenerateCredentialsParameters{
			TokenID: utils.String(id.ID()),
			Expiry:  password.Expiry,
			Name:    password.Name,
		}
		future, err := clients.RegistriesClient.GenerateCredentials(ctx, id.ResourceGroup, id.RegistryName, param)
		if err != nil {
			return nil, fmt.Errorf("generating password credential %s: %v", password.Name, err)
		}
		if err := future.WaitForCompletionRef(ctx, clients.RegistriesClient.Client); err != nil {
			return nil, fmt.Errorf("waiting for password credential generation for %s: %v", password.Name, err)
		}

		result, err := future.Result(*clients.RegistriesClient)
		if err != nil {
			return nil, fmt.Errorf("getting password credential after creation for %s: %v", password.Name, err)
		}

		genPasswords = append(genPasswords, containerregistry.TokenPassword{
			Expiry: password.Expiry,
			Name:   password.Name,
			Value:  (*result.Passwords)[idx].Value,
		})
	}
	return genPasswords, nil
}
