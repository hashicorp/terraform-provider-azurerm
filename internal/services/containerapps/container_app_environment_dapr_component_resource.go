package containerapps

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2022-03-01/daprcomponents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppEnvironmentDaprComponentResource struct{}

type ContainerAppEnvironmentDaprComponentModel struct {
	Name                 string                 `tfschema:"name"`
	ManagedEnvironmentId string                 `tfschema:"container_app_environment_id"`
	Type                 string                 `tfschema:"type"`
	Version              string                 `tfschema:"version"`
	IgnoreErrors         bool                   `tfschema:"ignore_errors"`
	InitTimeout          string                 `tfschema:"init_timeout"`
	Secrets              []helpers.Secret       `tfschema:"secret"`
	Scopes               []string               `tfschema:"scopes"`
	Metadata             []helpers.DaprMetadata `tfschema:"metadata"`
}

var _ sdk.ResourceWithUpdate = ContainerAppEnvironmentDaprComponentResource{}

func (r ContainerAppEnvironmentDaprComponentResource) ModelObject() interface{} {
	return &ContainerAppEnvironmentDaprComponentModel{}
}

func (r ContainerAppEnvironmentDaprComponentResource) ResourceType() string {
	return "azurerm_container_app_environment_dapr_component"
}

func (r ContainerAppEnvironmentDaprComponentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return daprcomponents.ValidateDaprComponentID
}

func (r ContainerAppEnvironmentDaprComponentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: helpers.ValidateDaprComponentName,
			Description:  "The name for this Dapr Component.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: daprcomponents.ValidateManagedEnvironmentID,
			Description:  "The Container App Managed Environment ID to configure this Dapr component on.",
		},

		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The Dapr Component Type. For example `state.azure.blobstorage`.",
		},

		"version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The version of the component.",
		},

		"init_timeout": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "5s",
			ValidateFunc: helpers.ValidateInitTimeout,
			Description:  "The component initialisation timeout in ISO8601 format. e.g. `5s`, `2h`, `1m`. Defaults to `5s`.",
		},

		"ignore_errors": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Should the Dapr sidecar to continue initialisation if the component fails to load. Defaults to `false`",
		},

		"secret": helpers.SecretsSchema(),

		"metadata": helpers.ContainerAppEnvironmentDaprMetadataSchema(),

		"scopes": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A list of scopes to which this component applies. e.g. a Container App's `dapr.app_id` value.",
		},
	}
}

func (r ContainerAppEnvironmentDaprComponentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ContainerAppEnvironmentDaprComponentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.DaprComponentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var daprComponent ContainerAppEnvironmentDaprComponentModel

			if err := metadata.Decode(&daprComponent); err != nil {
				return err
			}

			managedEnvironmentId, err := daprcomponents.ParseManagedEnvironmentID(daprComponent.ManagedEnvironmentId)
			if err != nil {
				return err
			}

			id := daprcomponents.NewDaprComponentID(subscriptionId, managedEnvironmentId.ResourceGroupName, managedEnvironmentId.EnvironmentName, daprComponent.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			daprComponentRequest := daprcomponents.DaprComponent{
				Properties: &daprcomponents.DaprComponentProperties{
					ComponentType: pointer.To(daprComponent.Type),
					IgnoreErrors:  pointer.To(daprComponent.IgnoreErrors),
					InitTimeout:   pointer.To(daprComponent.InitTimeout),
					Metadata:      expandDaprComponentPropertiesMetadata(daprComponent.Metadata),
					Secrets:       helpers.ExpandDaprSecrets(daprComponent.Secrets),
					Scopes:        expandDaprComponentPropertiesScopes(daprComponent.Scopes),
					Version:       pointer.To(daprComponent.Version),
				},
			}

			if len(daprComponent.Scopes) > 0 {
				daprComponentRequest.Properties.Scopes = &daprComponent.Scopes
			}

			if _, err := client.CreateOrUpdate(ctx, id, daprComponentRequest); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ContainerAppEnvironmentDaprComponentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.DaprComponentsClient
			id, err := daprcomponents.ParseDaprComponentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			daprComponentResp, err := client.Get(ctx, *id)
			if err != nil || daprComponentResp.Model == nil {
				if response.WasNotFound(daprComponentResp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := daprComponentResp.Model

			var state ContainerAppEnvironmentDaprComponentModel

			state.Name = id.ComponentName
			state.ManagedEnvironmentId = daprcomponents.NewManagedEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.EnvironmentName).ID()
			if props := model.Properties; props != nil {
				state.Version = pointer.From(props.Version)
				state.Type = pointer.From(props.ComponentType)
				state.Scopes = pointer.From(props.Scopes)
				state.InitTimeout = pointer.From(props.InitTimeout)
				state.IgnoreErrors = pointer.From(props.IgnoreErrors)
				state.Metadata = flattenDaprComponentPropertiesMetadata(props.Metadata)
			}

			secretsResp, err := client.ListSecrets(ctx, *id)
			if err != nil || secretsResp.Model == nil {
				if secretsResp.HttpResponse == nil || secretsResp.HttpResponse.StatusCode != http.StatusNoContent {
					return fmt.Errorf("retrieving secrets for %s: %+v", *id, err)
				}
			}
			state.Secrets = helpers.FlattenContainerAppDaprSecrets(secretsResp.Model)

			return metadata.Encode(&state)
		},
	}
}

func (r ContainerAppEnvironmentDaprComponentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.DaprComponentsClient
			id, err := daprcomponents.ParseDaprComponentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if resp, err := client.Delete(ctx, *id); err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("deleting %s: %+v", *id, err)
				}
			}
			return nil
		},
	}
}

func (r ContainerAppEnvironmentDaprComponentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.DaprComponentsClient
			id, err := daprcomponents.ParseDaprComponentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ContainerAppEnvironmentDaprComponentModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s for update: %+v", *id, err)
			}

			// Populate the secrets from the List API to prevent accidental removal.
			secretsResp, err := client.ListSecrets(ctx, *id)
			if err != nil || secretsResp.Model == nil {
				if secretsResp.HttpResponse == nil || secretsResp.HttpResponse.StatusCode != http.StatusNoContent {
					return fmt.Errorf("retrieving secrets for %s: %+v", *id, err)
				}
			}

			existing.Model.Properties.Secrets = helpers.UnpackContainerDaprSecretsCollection(secretsResp.Model)

			if metadata.ResourceData.HasChange("version") {
				existing.Model.Properties.Version = pointer.To(state.Version)
			}

			if metadata.ResourceData.HasChange("init_timeout") {
				existing.Model.Properties.InitTimeout = pointer.To(state.InitTimeout)
			}

			if metadata.ResourceData.HasChange("ignore_errors") {
				existing.Model.Properties.IgnoreErrors = pointer.To(state.IgnoreErrors)
			}

			if metadata.ResourceData.HasChange("secret") {
				existing.Model.Properties.Secrets = helpers.ExpandDaprSecrets(state.Secrets)
			}

			if metadata.ResourceData.HasChange("metadata") {
				existing.Model.Properties.Metadata = expandDaprComponentPropertiesMetadata(state.Metadata)
			}

			if metadata.ResourceData.HasChange("scopes") {
				existing.Model.Properties.Scopes = expandDaprComponentPropertiesScopes(state.Scopes)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandDaprComponentPropertiesMetadata(input []helpers.DaprMetadata) *[]daprcomponents.DaprMetadata {
	if len(input) == 0 {
		return nil
	}

	result := make([]daprcomponents.DaprMetadata, 0)

	for _, v := range input {
		d := daprcomponents.DaprMetadata{
			Name: pointer.To(v.Name),
		}
		if v.Value != "" {
			d.Value = pointer.To(v.Value)
		}
		if v.SecretName != "" {
			d.SecretRef = pointer.To(v.SecretName)
		}
		result = append(result, d)
	}

	return &result
}

func flattenDaprComponentPropertiesMetadata(input *[]daprcomponents.DaprMetadata) []helpers.DaprMetadata {
	if input == nil {
		return nil
	}

	result := make([]helpers.DaprMetadata, 0)
	for _, v := range *input {
		result = append(result, helpers.DaprMetadata{
			Name:       pointer.From(v.Name),
			SecretName: pointer.From(v.SecretRef),
			Value:      pointer.From(v.Value),
		})
	}

	return result
}

func expandDaprComponentPropertiesScopes(input []string) *[]string {
	if input == nil {
		return nil
	}

	return &input
}
