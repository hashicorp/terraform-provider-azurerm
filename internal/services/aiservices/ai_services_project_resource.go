package aiservices

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AIServicesProject struct{}

type AIServicesProjectModel struct {
	Name                      string                                     `tfschema:"name"`
	Location                  string                                     `tfschema:"location"`
	AIServicesHubId           string                                     `tfschema:"ai_services_hub_id"`
	Identity                  []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	HighBusinessImpactEnabled bool                                       `tfschema:"high_business_impact_enabled"`
	ImageBuildComputeName     string                                     `tfschema:"image_build_compute_name"`
	Description               string                                     `tfschema:"description"`
	FriendlyName              string                                     `tfschema:"friendly_name"`
	ProjectId                 string                                     `tfschema:"project_id"`
	Tags                      map[string]interface{}                     `tfschema:"tags"`
}

func (r AIServicesProject) ModelObject() interface{} {
	return &AIServicesProjectModel{}
}

func (r AIServicesProject) ResourceType() string {
	return "azurerm_ai_services_project"
}

func (r AIServicesProject) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return workspaces.ValidateWorkspaceID
}

func (r AIServicesProject) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.MachineLearning.Workspaces
		resp, err := client.Get(ctx, *id)
		if err != nil || resp.Model == nil || resp.Model.Kind == nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if !strings.EqualFold(*resp.Model.Kind, "Project") {
			return fmt.Errorf("importing %s: specified workspace is not of kind `Project`, got `%s`", id, *resp.Model.Kind)
		}

		return nil
	}
}

var _ sdk.ResourceWithUpdate = AIServicesProject{}

var _ sdk.ResourceWithCustomImporter = AIServicesProject{}

func (r AIServicesProject) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_-]{2,32}$"),
				"AI Services Project name must be 2 - 32 characters long, contain only letters, numbers and hyphens.",
			),
		},

		"location": commonschema.Location(),

		"ai_services_hub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"high_business_impact_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			// NOTE: O+C creating a project that has encryption enabled with system assigned identity will set this property to true
			Computed: true,
			ForceNew: true,
		},

		"image_build_compute_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"friendly_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r AIServicesProject) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"project_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r AIServicesProject) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model AIServicesProjectModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			hubId, err := workspaces.ParseWorkspaceID(model.AIServicesHubId)
			if err != nil {
				return err
			}

			id := workspaces.NewWorkspaceID(subscriptionId, hubId.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_ai_services_project", id.ID())
			}

			payload := workspaces.Workspace{
				Name:     pointer.To(id.WorkspaceName),
				Location: pointer.To(location.Normalize(model.Location)),
				Tags:     tags.Expand(model.Tags),
				Kind:     pointer.To("Project"),
				Properties: &workspaces.WorkspaceProperties{
					HubResourceId: pointer.To(hubId.ID()),
				},
			}

			if len(model.Identity) > 0 {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if model.Description != "" {
				payload.Properties.Description = pointer.To(model.Description)
			}

			if model.FriendlyName != "" {
				payload.Properties.FriendlyName = pointer.To(model.FriendlyName)
			}

			if model.HighBusinessImpactEnabled {
				payload.Properties.HbiWorkspace = pointer.To(model.HighBusinessImpactEnabled)
			}

			if model.ImageBuildComputeName != "" {
				payload.Properties.ImageBuildCompute = pointer.To(model.ImageBuildComputeName)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AIServicesProject) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces

			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state AIServicesProjectModel
			if err := metadata.Decode(&state); err != nil {
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

			payload := existing.Model

			// Hubs and Projects share the same API where Projects inherit the KV/Storage/AppInsights/CR/Network/Encryption settings from the Hub
			// When updating a Project the API will error when trying to send the inherited settings that get returned when we retrieve the resource for patching in changes
			// This is a hack to work around this behaviour and design, so that we can continue to support the use of `ignore_changes` on the resource
			payload.Properties.ManagedNetwork = nil
			payload.Properties.KeyVault = nil
			payload.Properties.StorageAccount = nil
			payload.Properties.ContainerRegistry = nil
			payload.Properties.ApplicationInsights = nil
			payload.Properties.Encryption = nil

			if metadata.ResourceData.HasChange("image_build_compute_name") {
				payload.Properties.ImageBuildCompute = pointer.To(state.ImageBuildComputeName)
			}

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = pointer.To(state.Description)
			}

			if metadata.ResourceData.HasChange("friendly_name") {
				payload.Properties.FriendlyName = pointer.To(state.FriendlyName)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]interface{}))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = tags.Expand(state.Tags)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r AIServicesProject) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces

			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			hub := AIServicesProjectModel{
				Name: id.WorkspaceName,
			}

			if model := resp.Model; model != nil {
				hub.Location = location.NormalizeNilable(model.Location)
				flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				hub.Identity = flattenedIdentity

				hub.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					if v := pointer.From(props.HubResourceId); v != "" {
						hubId, err := workspaces.ParseWorkspaceID(v)
						if err != nil {
							return err
						}
						hub.AIServicesHubId = hubId.ID()
					}

					hub.Description = pointer.From(props.Description)
					hub.FriendlyName = pointer.From(props.FriendlyName)
					hub.HighBusinessImpactEnabled = pointer.From(props.HbiWorkspace)
					hub.ImageBuildComputeName = pointer.From(props.ImageBuildCompute)
					hub.ProjectId = pointer.From(props.WorkspaceId)
				}
			}

			return metadata.Encode(&hub)
		},
	}
}

func (r AIServicesProject) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.Workspaces

			id, err := workspaces.ParseWorkspaceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, workspaces.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
