package containers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/credentialsets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = ContainerRegistryCredentialSetResource{}

type ContainerRegistryCredentialSetResource struct{}

func (ContainerRegistryCredentialSetResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The name of the credential set.",
		},
		"container_registry_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: registries.ValidateRegistryID,
		},
		"login_server": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"authentication_credentials": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"username_secret_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyVaultValidate.VersionlessNestedItemId,
					},
					"password_secret_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyVaultValidate.VersionlessNestedItemId,
					},
				},
			},
		},
	}
}

func (ContainerRegistryCredentialSetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),
	}
}

type AuthenticationCredential struct {
	UsernameSecretId string `tfschema:"username_secret_id"`
	PasswordSecretId string `tfschema:"password_secret_id"`
}

type ContainerRegistryCredentialSetModel struct {
	Name                     string                                     `tfschema:"name"`
	ContainerRegistryId      string                                     `tfschema:"container_registry_id"`
	LoginServer              string                                     `tfschema:"login_server"`
	AuthenticationCredential []AuthenticationCredential                 `tfschema:"authentication_credentials"`
	Identity                 []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
}

func (ContainerRegistryCredentialSetResource) ModelObject() interface{} {
	return &ContainerRegistryCredentialSetModel{}
}

func (ContainerRegistryCredentialSetResource) ResourceType() string {
	return "azurerm_container_registry_credential_set"
}

func (r ContainerRegistryCredentialSetResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.CredentialSetsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ContainerRegistryCredentialSetModel
			if err := metadata.Decode(&config); err != nil {
				return err
			}

			log.Printf("[INFO] preparing arguments for Container Registry Credential Set creation.")

			registryId, err := registries.ParseRegistryID(config.ContainerRegistryId)
			if err != nil {
				return err
			}

			id := credentialsets.NewCredentialSetID(subscriptionId,
				registryId.ResourceGroupName,
				registryId.RegistryName,
				config.Name,
			)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := credentialsets.CredentialSet{
				Name: pointer.To(id.CredentialSetName),
				Properties: &credentialsets.CredentialSetProperties{
					LoginServer:     pointer.To(config.LoginServer),
					AuthCredentials: expandAuthCredentials(config.AuthenticationCredential),
				},
				Identity: &identity.SystemAndUserAssignedMap{
					Type: identity.TypeSystemAssigned,
				},
			}

			if err := client.CreateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ContainerRegistryCredentialSetResource) Update() sdk.ResourceFunc {
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

			if err := client.UpdateThenPoll(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (ContainerRegistryCredentialSetResource) Read() sdk.ResourceFunc {
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
				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				config.Identity = pointer.From(flattenedIdentity)
				if props := model.Properties; props != nil {
					config.LoginServer = pointer.From(props.LoginServer)
					config.AuthenticationCredential = flattenAuthCredentials(props.AuthCredentials)
				}
			}
			return metadata.Encode(&config)
		},
	}
}

func (ContainerRegistryCredentialSetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.CredentialSetsClient
			id, err := credentialsets.ParseCredentialSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (ContainerRegistryCredentialSetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return credentialsets.ValidateCredentialSetID
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
