// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/onlineendpoint"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MachineLearningOnlineEndpoint struct{}

type MachineLearningOnlineEndpointModel struct {
	Name                       string                                     `tfschema:"name"`
	MachineLearningWorkspaceId string                                     `tfschema:"machine_learning_workspace_id"`
	Location                   string                                     `tfschema:"location"`
	AuthenticationMode         string                                     `tfschema:"authentication_mode"`
	Description                string                                     `tfschema:"description"`
	Identity                   []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PublicNetworkAccessEnabled bool                                       `tfschema:"public_network_access_enabled"`
	Properties                 map[string]interface{}                     `tfschema:"properties"`
	RestEndpoint               string                                     `tfschema:"rest_endpoint"`
	SwaggerUri                 string                                     `tfschema:"swagger_uri"`
	Tags                       map[string]interface{}                     `tfschema:"tags"`
}

func (r MachineLearningOnlineEndpoint) ModelObject() interface{} {
	return &MachineLearningOnlineEndpointModel{}
}

func (r MachineLearningOnlineEndpoint) ResourceType() string {
	return "azurerm_machine_learning_online_endpoint"
}

func (r MachineLearningOnlineEndpoint) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return onlineendpoint.ValidateOnlineEndpointID
}

var (
	_ sdk.ResourceWithUpdate        = MachineLearningOnlineEndpoint{}
	_ sdk.ResourceWithCustomizeDiff = MachineLearningOnlineEndpoint{}
)

func (r MachineLearningOnlineEndpoint) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[A-Za-z][A-Za-z0-9-]{1,30}[A-Za-z0-9]$`),
				"`name` must start with a letter, contain only letters, digits, or dashes, end with a letter or digit, and be between 3 and 32 characters"),
		},

		"machine_learning_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"location": commonschema.Location(),

		"authentication_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(onlineendpoint.EndpointAuthModeKey),
			ValidateFunc: validation.StringInSlice(onlineendpoint.PossibleValuesForEndpointAuthMode(), false),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityRequiredForceNew(),

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r MachineLearningOnlineEndpoint) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"rest_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"swagger_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r MachineLearningOnlineEndpoint) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff.Id() != "" && metadata.ResourceDiff.HasChange("description") {
				return fmt.Errorf("the `description` of a Machine Learning Online Endpoint cannot be changed after creation")
			}

			return nil
		},
	}
}

func (r MachineLearningOnlineEndpoint) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.OnlineEndpoints
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MachineLearningOnlineEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := workspaces.ParseWorkspaceID(model.MachineLearningWorkspaceId)
			if err != nil {
				return fmt.Errorf("parsing `machine_learning_workspace_id`: %+v", err)
			}

			id := onlineendpoint.NewOnlineEndpointID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_machine_learning_online_endpoint", id.ID())
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			publicNetworkAccess := onlineendpoint.PublicNetworkAccessTypeEnabled
			if !model.PublicNetworkAccessEnabled {
				publicNetworkAccess = onlineendpoint.PublicNetworkAccessTypeDisabled
			}

			properties := onlineendpoint.OnlineEndpoint{
				AuthMode:            onlineendpoint.EndpointAuthMode(model.AuthenticationMode),
				Properties:          expandMachineLearningOnlineEndpointProperties(model.Properties),
				PublicNetworkAccess: pointer.To(publicNetworkAccess),
				Description:         pointer.To(model.Description),
			}

			payload := onlineendpoint.OnlineEndpointTrackedResource{
				Identity:   expandedIdentity,
				Location:   location.Normalize(model.Location),
				Properties: properties,
				Tags:       tags.Expand(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r MachineLearningOnlineEndpoint) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.OnlineEndpoints

			id, err := onlineendpoint.ParseOnlineEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing online endpoint ID `%s`: %+v", metadata.ResourceData.Id(), err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := MachineLearningOnlineEndpointModel{
				Name:                       id.OnlineEndpointName,
				MachineLearningWorkspaceId: workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID(),
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			state.Location = location.Normalize(model.Location)

			flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			state.Identity = flattenedIdentity

			props := model.Properties
			state.AuthenticationMode = string(props.AuthMode)
			state.Description = pointer.From(props.Description)
			state.Properties = flattenMachineLearningOnlineEndpointProperties(props.Properties)
			state.PublicNetworkAccessEnabled = pointer.From(props.PublicNetworkAccess) == onlineendpoint.PublicNetworkAccessTypeEnabled
			state.RestEndpoint = pointer.From(props.ScoringUri)
			state.SwaggerUri = pointer.From(props.SwaggerUri)
			state.Tags = tags.Flatten(model.Tags)

			return metadata.Encode(&state)
		},
	}
}

func (r MachineLearningOnlineEndpoint) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.OnlineEndpoints

			id, err := onlineendpoint.ParseOnlineEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing online endpoint ID `%s`: %+v", metadata.ResourceData.Id(), err)
			}

			var model MachineLearningOnlineEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			publicNetworkAccess := onlineendpoint.PublicNetworkAccessTypeEnabled
			if !model.PublicNetworkAccessEnabled {
				publicNetworkAccess = onlineendpoint.PublicNetworkAccessTypeDisabled
			}

			properties := onlineendpoint.OnlineEndpoint{
				AuthMode:            onlineendpoint.EndpointAuthMode(model.AuthenticationMode),
				Properties:          expandMachineLearningOnlineEndpointProperties(model.Properties),
				PublicNetworkAccess: pointer.To(publicNetworkAccess),
				Description:         pointer.To(model.Description),
			}

			payload := onlineendpoint.OnlineEndpointTrackedResource{
				Identity:   expandedIdentity,
				Location:   location.Normalize(model.Location),
				Properties: properties,
				Tags:       tags.Expand(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r MachineLearningOnlineEndpoint) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.OnlineEndpoints

			id, err := onlineendpoint.ParseOnlineEndpointID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing online endpoint ID `%s`: %+v", metadata.ResourceData.Id(), err)
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

var machineLearningOnlineEndpointSystemProperties = map[string]struct{}{
	"azureasyncoperationuri":   {},
	"azureml.onlineendpointid": {},
}

func isMachineLearningOnlineEndpointSystemProperty(key string) bool {
	_, ok := machineLearningOnlineEndpointSystemProperties[strings.ToLower(key)]
	return ok
}

func expandMachineLearningOnlineEndpointProperties(input map[string]interface{}) *map[string]string {
	if len(input) == 0 {
		return nil
	}

	result := make(map[string]string, len(input))
	for k, v := range input {
		if isMachineLearningOnlineEndpointSystemProperty(k) {
			continue
		}
		result[k] = v.(string)
	}

	if len(result) == 0 {
		return nil
	}

	return &result
}

func flattenMachineLearningOnlineEndpointProperties(input *map[string]string) map[string]interface{} {
	if input == nil {
		return map[string]interface{}{}
	}

	result := make(map[string]interface{}, len(*input))
	for k, v := range *input {
		if isMachineLearningOnlineEndpointSystemProperty(k) {
			continue
		}
		result[k] = v
	}

	return result
}
