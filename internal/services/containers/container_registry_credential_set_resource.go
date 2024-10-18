package containers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/credentialsets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
		"auth_credentials": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"username_secret_identifier": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"password_secret_identifier": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}

func (ContainerRegistryCredentialSetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"identity": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

type Identity struct {
	Type        string `tfschema:"type"`
	TenantId    string `tfschema:"tenant_id"`
	PrincipalId string `tfschema:"principal_id"`
}

type AuthCredential struct {
	UsernameSecretIdentifier string `tfschema:"username_secret_identifier"`
	PasswordSecretIdentifier string `tfschema:"password_secret_identifier"`
}

type ContainerRegistryCredentialSetModel struct {
	Name                string           `tfschema:"name"`
	ContainerRegistryId string           `tfschema:"container_registry_id"`
	LoginServer         string           `tfschema:"login_server"`
	AuthCredentials     []AuthCredential `tfschema:"auth_credentials"`
	Identity            []Identity       `tfschema:"identity"`
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

			registryId, err := registries.ParseRegistryID(metadata.ResourceData.Get("container_registry_id").(string))
			if err != nil {
				return err
			}

			id := credentialsets.NewCredentialSetID(subscriptionId,
				registryId.ResourceGroupName,
				registryId.RegistryName,
				metadata.ResourceData.Get("name").(string),
			)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := credentialsets.CredentialSet{
				Name: &id.CredentialSetName,
				Properties: &credentialsets.CredentialSetProperties{
					LoginServer:     pointer.To(config.LoginServer),
					AuthCredentials: expandAuthCredentials(config.AuthCredentials),
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

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			var model ContainerRegistryCredentialSetModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("login_server") {
				existing.Model.Properties.LoginServer = pointer.To(model.LoginServer)
			}
			if metadata.ResourceData.HasChange("auth_credentials") {
				existing.Model.Properties.AuthCredentials = expandAuthCredentials(model.AuthCredentials)
			}

			param := credentialsets.CredentialSetUpdateParameters{
				Identity: expandIdentity(model.Identity),
				Properties: &credentialsets.CredentialSetUpdateProperties{
					AuthCredentials: existing.Model.Properties.AuthCredentials,
				},
			}
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
				config.Identity = flattenIdentity(model.Identity)
				if props := model.Properties; props != nil {
					config.LoginServer = pointer.From(props.LoginServer)
					config.AuthCredentials = flattenAuthCredentials(props.AuthCredentials)
				}
			}
			if err := metadata.Encode(&config); err != nil {
				return fmt.Errorf("encoding: %+v", err)
			}
			return nil
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

func expandAuthCredentials(input []AuthCredential) *[]credentialsets.AuthCredential {
	if len(input) == 0 {
		return nil
	}
	output := make([]credentialsets.AuthCredential, 0, len(input))
	for _, v := range input {
		output = append(output, credentialsets.AuthCredential{
			Name:                     pointer.To(credentialsets.CredentialNameCredentialOne),
			UsernameSecretIdentifier: pointer.To(v.UsernameSecretIdentifier),
			PasswordSecretIdentifier: pointer.To(v.PasswordSecretIdentifier),
		})
	}
	return &output
}

func flattenAuthCredentials(input *[]credentialsets.AuthCredential) []AuthCredential {
	if input == nil {
		return nil
	}
	output := make([]AuthCredential, len(*input))
	for i, v := range *input {
		output[i] = AuthCredential{
			UsernameSecretIdentifier: pointer.From(v.UsernameSecretIdentifier),
			PasswordSecretIdentifier: pointer.From(v.PasswordSecretIdentifier),
		}
	}
	return output
}

func flattenIdentity(input *identity.SystemAndUserAssignedMap) []Identity {
	if input == nil {
		return nil
	}
	output := make([]Identity, 1)
	output[0] = Identity{
		Type:        string(input.Type),
		TenantId:    input.TenantId,
		PrincipalId: input.PrincipalId,
	}
	return output
}

func expandIdentity(input []Identity) *identity.SystemAndUserAssignedMap {
	if len(input) == 0 {
		return nil
	}
	output := identity.SystemAndUserAssignedMap{
		Type: identity.TypeSystemAssigned,
	}
	return &output
}
