package containerapps

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2023-05-01/managedenvironmentsstorages"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppJobResource struct{}

type ContainerAppJobModel struct {
	Name                      string                                     `tfschema:"name"`
	ResourceGroup             string                                     `tfschema:"resource_group_name"`
	Location                  string                                     `tfschema:"location"`
	ContainerAppEnvironmentId string                                     `tfschema:"container_app_environment_id"`
	WorkloadProfileName       string                                     `tfschema:"workload_profile_name"`
	Template                  []TemplateModel                            `tfschema:"template"`
	Configuration             []ConfigurationModel                       `tfschema:"configuration"`
	Identity                  []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Tags                      map[string]interface{}                     `tfschema:"tags"`

	OutboundIPAddresses []string `tfschema:"outbound_ip_addresses"`
	EventStreamEndpoint string   `tfschema:"event_stream_endpoint"`
}

var _ sdk.ResourceWithUpdate = ContainerAppJobResource{}

func (r ContainerAppJobResource) ModelObject() interface{} {
	return &ContainerAppJobModel{}
}

func (r ContainerAppJobResource) ResourceType() string {
	return "azurerm_container_app_job"
}

func (r ContainerAppJobResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return jobs.ValidateJobID
}

func (r ContainerAppJobResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: certificates.ValidateManagedEnvironmentID,
		},

		"workload_profile_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"template": templateSchema(),

		"configuration": configurationSchema(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}

func (r ContainerAppJobResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"event_stream_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ContainerAppJobResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.JobClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ContainerAppJobModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := jobs.NewJobID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			job := jobs.Job{
				Location: model.Location,
				Name:     pointer.To(model.Name),
				Properties: &jobs.JobProperties{
					EnvironmentId: pointer.To(model.ContainerAppEnvironmentId),
				},
				Tags: tags.Expand(model.Tags),
				Type: nil,
			}

			ident, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			job.Identity = pointer.To(identity.LegacySystemAndUserAssignedMap(*ident))

			if model.WorkloadProfileName != "" {
				job.Properties.WorkloadProfileName = pointer.To(model.WorkloadProfileName)
			}

			if model.Configuration != nil {
				config, err := expandContainerAppJobConfiguration(model.Configuration)
				if err != nil {
					return fmt.Errorf("expanding `configuration`: %+v", err)
				}
				job.Properties.Configuration = config
			}

			if model.Template != nil {
				template, err := expandContainerAppJobTemplate(model.Template)
				if err != nil {
					return fmt.Errorf("expanding `template`: %+v", err)
				}
				job.Properties.Template = template
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, job); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ContainerAppJobResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.JobClient

			id, err := jobs.ParseJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			var state ContainerAppJobModel

			state.Name = id.JobName
			state.ResourceGroup = id.ResourceGroupName

			if model := existing.Model; model != nil {
				state.ContainerAppEnvironmentId = *model.Properties.EnvironmentId
				state.Location = model.Location
				state.Tags = tags.Flatten(model.Tags)
				ident, err := identity.FlattenSystemAndUserAssignedMapToModel(pointer.To(identity.SystemAndUserAssignedMap(*model.Identity)))
				if err != nil {
					return err
				}
				state.Identity = pointer.From(ident)

				if props := model.Properties; props != nil {
					envId, err := managedenvironmentsstorages.ParseManagedEnvironmentIDInsensitively(pointer.From(props.EnvironmentId))
					if err != nil {
						return err
					}
					state.ContainerAppEnvironmentId = envId.ID()
					state.Template = flattenContainerAppJobTemplate(props.Template)
					state.Configuration = flattenContainerAppJobConfiguration(props.Configuration)
					state.WorkloadProfileName = pointer.From(props.WorkloadProfileName)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ContainerAppJobResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.JobClient

			id, err := jobs.ParseJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ContainerAppJobModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := existing.Model

			if model.Properties == nil {
				return fmt.Errorf("retrieving properties for %s for update: %+v", *id, err)
			}

			d := metadata.ResourceData

			if d.HasChange("template") {
				template, err := expandContainerAppJobTemplate(state.Template)
				if err != nil {
					return fmt.Errorf("expanding `template`: %+v", err)
				}
				model.Properties.Template = template
			}

			if d.HasChange("configuration") {
				config, err := expandContainerAppJobConfiguration(state.Configuration)
				if err != nil {
					return fmt.Errorf("expanding `configuration`: %+v", err)
				}
				model.Properties.Configuration = config
			}

			if d.HasChange("identity") {
				ident, err := identity.ExpandSystemAndUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				model.Identity = pointer.To(identity.LegacySystemAndUserAssignedMap(*ident))
			}

			if d.HasChange("workload_profile_name") {
				model.Properties.WorkloadProfileName = pointer.To(state.WorkloadProfileName)
			}

			if d.HasChange("tags") {
				model.Tags = tags.Expand(state.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ContainerAppJobResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.JobClient

			id, err := jobs.ParseJobID(metadata.ResourceData.Id())
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
