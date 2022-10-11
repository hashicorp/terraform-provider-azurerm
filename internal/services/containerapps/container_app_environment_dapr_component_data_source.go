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
)

type ContainerAppEnvironmentDaprComponentDataSource struct{}

type ContainerAppEnvironmentDaprComponentDataSourceModel struct {
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

var _ sdk.DataSource = ContainerAppEnvironmentDaprComponentDataSource{}

func (r ContainerAppEnvironmentDaprComponentDataSource) ModelObject() interface{} {
	return &ContainerAppEnvironmentDaprComponentDataSourceModel{}
}

func (r ContainerAppEnvironmentDaprComponentDataSource) ResourceType() string {
	return "azurerm_container_app_environment_dapr_component"
}

func (r ContainerAppEnvironmentDaprComponentDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: helpers.ValidateDaprComponentName,
			Description:  "The name for this Dapr component.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: daprcomponents.ValidateManagedEnvironmentID,
			Description:  "The Container App Managed Environment ID to configure this Dapr component on.",
		},
	}
}

func (r ContainerAppEnvironmentDaprComponentDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"type": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Dapr Component Type.",
		},

		"version": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The version of the component.",
		},

		"init_timeout": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The component initialisation timeout in ISO8601 format. e.g. `5s`, `2h`, `1m`. Defaults to `5s`",
		},

		"ignore_errors": {
			Type:        pluginsdk.TypeBool,
			Computed:    true,
			Description: "Should the Dapr sidecar to continue initialisation if the component fails to load. Defaults to `false`",
		},

		"secret": helpers.SecretsDataSourceSchema(),

		"metadata": helpers.ContainerAppEnvironmentDaprMetadataDataSourceSchema(),

		"scopes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Description: "A list of scopes to which this component applies. e.g. a Container App's `dapr.app_id` value.",
		},
	}
}

func (r ContainerAppEnvironmentDaprComponentDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.DaprComponentsClient

			var daprComponent ContainerAppEnvironmentDaprComponentDataSourceModel
			if err := metadata.Decode(&daprComponent); err != nil {
				return err
			}

			envId, err := daprcomponents.ParseManagedEnvironmentID(daprComponent.ManagedEnvironmentId)
			if err != nil {
				return err
			}

			id := daprcomponents.NewDaprComponentID(envId.SubscriptionId, envId.ResourceGroupName, envId.EnvironmentName, daprComponent.Name)

			daprComponentResp, err := client.Get(ctx, id)
			if err != nil || daprComponentResp.Model == nil {
				if response.WasNotFound(daprComponentResp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := daprComponentResp.Model

			daprComponent.Name = id.ComponentName
			daprComponent.ManagedEnvironmentId = envId.ID()
			if props := model.Properties; props != nil {
				daprComponent.Version = pointer.From(props.Version)
				daprComponent.Type = pointer.From(props.ComponentType)
				daprComponent.Scopes = scopesPtr(props.Scopes)
				daprComponent.InitTimeout = pointer.From(props.InitTimeout)
				daprComponent.IgnoreErrors = pointer.From(props.IgnoreErrors)
				daprComponent.Metadata = flattenDaprComponentPropertiesMetadata(props.Metadata)
			}

			secretsResp, err := client.ListSecrets(ctx, id)
			if err != nil || secretsResp.Model == nil {
				if secretsResp.HttpResponse == nil || secretsResp.HttpResponse.StatusCode != http.StatusNoContent {
					return fmt.Errorf("retrieving secrets for %s: %+v", id, err)
				}
			}
			daprComponent.Secrets = helpers.FlattenContainerAppDaprSecrets(secretsResp.Model)

			metadata.SetID(id)

			return metadata.Encode(&daprComponent)
		},
	}
}
