// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironmentsstorages"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
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
	Template                  []helpers.JobTemplateModel                 `tfschema:"template"`
	ReplicaRetryLimit         int64                                      `tfschema:"replica_retry_limit"`
	ReplicaTimeoutInSeconds   int64                                      `tfschema:"replica_timeout_in_seconds"`
	Secrets                   []helpers.Secret                           `tfschema:"secret"`
	Registries                []helpers.Registry                         `tfschema:"registry"`
	EventTriggerConfig        []helpers.EventTriggerConfiguration        `tfschema:"event_trigger_config"`
	ManualTriggerConfig       []helpers.ManualTriggerConfiguration       `tfschema:"manual_trigger_config"`
	ScheduleTriggerConfig     []helpers.ScheduleTriggerConfiguration     `tfschema:"schedule_trigger_config"`
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
			ValidateFunc: validate.ContainerAppJobName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"container_app_environment_id": commonschema.ResourceIDReferenceRequiredForceNew(&certificates.ManagedEnvironmentId{}),

		"replica_timeout_in_seconds": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"workload_profile_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"template": helpers.JobTemplateSchema(),

		"secret": helpers.SecretsSchema(),

		"replica_retry_limit": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},

		"registry": helpers.ContainerAppRegistrySchema(),

		"event_trigger_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			ExactlyOneOf: []string{
				"event_trigger_config",
				"manual_trigger_config",
				"schedule_trigger_config",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"parallelism": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"replica_completion_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"scale": helpers.ContainerAppsJobsScaleSchema(),
				},
			},
		},

		"schedule_trigger_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			ExactlyOneOf: []string{
				"event_trigger_config",
				"manual_trigger_config",
				"schedule_trigger_config",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cron_expression": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"parallelism": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"replica_completion_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

		"manual_trigger_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			ExactlyOneOf: []string{
				"event_trigger_config",
				"manual_trigger_config",
				"schedule_trigger_config",
			},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"parallelism": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"replica_completion_count": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

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

			registries, err := helpers.ExpandContainerAppJobRegistries(model.Registries)
			if err != nil {
				return fmt.Errorf("expanding registry config for %s: %v", id, err)
			}

			job := jobs.Job{
				Location: location.Normalize(model.Location),
				Properties: &jobs.JobProperties{
					Configuration: &jobs.JobConfiguration{
						ReplicaRetryLimit: pointer.To(model.ReplicaRetryLimit),
						ReplicaTimeout:    model.ReplicaTimeoutInSeconds,
						Secrets:           helpers.ExpandContainerAppJobSecrets(model.Secrets),
						Registries:        registries,
					},
					EnvironmentId: pointer.To(model.ContainerAppEnvironmentId),
					Template:      helpers.ExpandContainerAppJobTemplate(model.Template),
				},
				Tags: tags.Expand(model.Tags),
			}

			var triggerType jobs.TriggerType
			if len(model.ManualTriggerConfig) > 0 {
				triggerType = jobs.TriggerTypeManual
				job.Properties.Configuration.ManualTriggerConfig = helpers.ExpandContainerAppJobConfigurationManualTriggerConfig(model.ManualTriggerConfig)
			}
			if len(model.EventTriggerConfig) > 0 {
				triggerType = jobs.TriggerTypeEvent
				job.Properties.Configuration.EventTriggerConfig = helpers.ExpandContainerAppJobConfigurationEventTriggerConfig(model.EventTriggerConfig)
			}
			if len(model.ScheduleTriggerConfig) > 0 {
				triggerType = jobs.TriggerTypeSchedule
				job.Properties.Configuration.ScheduleTriggerConfig = helpers.ExpandContainerAppJobConfigurationScheduleTriggerConfig(model.ScheduleTriggerConfig)
			}
			job.Properties.Configuration.TriggerType = triggerType

			ident, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			job.Identity = ident

			if model.WorkloadProfileName != "" {
				job.Properties.WorkloadProfileName = pointer.To(model.WorkloadProfileName)
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
				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)
				if model.Identity != nil {
					if model.Identity != nil {
						ident, err := identity.FlattenSystemAndUserAssignedMapToModel(pointer.To((identity.SystemAndUserAssignedMap)(*model.Identity)))
						if err != nil {
							return err
						}
						state.Identity = pointer.From(ident)
					}
				}

				if props := model.Properties; props != nil {
					envId, err := managedenvironmentsstorages.ParseManagedEnvironmentIDInsensitively(pointer.From(props.EnvironmentId))
					if err != nil {
						return err
					}
					state.ContainerAppEnvironmentId = envId.ID()
					state.Template = helpers.FlattenContainerAppJobTemplate(props.Template)
					if config := props.Configuration; config != nil {
						state.Registries = helpers.FlattenContainerAppJobRegistries(config.Registries)
						state.ReplicaTimeoutInSeconds = config.ReplicaTimeout
						if config.ReplicaRetryLimit != nil {
							state.ReplicaRetryLimit = pointer.From(config.ReplicaRetryLimit)
						}

						switch config.TriggerType {
						case jobs.TriggerTypeEvent:
							state.EventTriggerConfig = helpers.FlattenContainerAppJobConfigurationEventTriggerConfig(config.EventTriggerConfig)
						case jobs.TriggerTypeManual:
							state.ManualTriggerConfig = helpers.FlattenContainerAppJobConfigurationManualTriggerConfig(config.ManualTriggerConfig)
						case jobs.TriggerTypeSchedule:
							state.ScheduleTriggerConfig = helpers.FlattenContainerAppJobConfigurationScheduleTriggerConfig(config.ScheduleTriggerConfig)
						}
					}
					state.WorkloadProfileName = pointer.From(props.WorkloadProfileName)
				}
			}

			secretResp, err := client.ListSecrets(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing secrets for %s: %+v", *id, err)
			}
			state.Secrets = helpers.FlattenContainerAppJobSecrets(secretResp.Model)

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

			if model.Properties.Configuration == nil {
				model.Properties.Configuration = &jobs.JobConfiguration{}
			}

			secretsResp, err := client.ListSecrets(ctx, *id)
			if err != nil || secretsResp.Model == nil {
				if !response.WasStatusCode(secretsResp.HttpResponse, http.StatusNoContent) {
					return fmt.Errorf("retrieving secrets for update for %s: %+v", *id, err)
				}
			}
			model.Properties.Configuration.Secrets = helpers.UnpackContainerJobSecretsCollection(secretsResp.Model)

			d := metadata.ResourceData

			if d.HasChange("secret") {
				model.Properties.Configuration.Secrets = helpers.ExpandContainerAppJobSecrets(state.Secrets)
			}

			if d.HasChange("registry") {
				model.Properties.Configuration.Registries, err = helpers.ExpandContainerAppJobRegistries(state.Registries)
				if err != nil {
					return fmt.Errorf("invalid registry config for %s: %v", id, err)
				}
			}

			if d.HasChange("replica_retry_limit") {
				model.Properties.Configuration.ReplicaRetryLimit = pointer.To(state.ReplicaRetryLimit)
			}

			if d.HasChange("replica_timeout_in_seconds") {
				model.Properties.Configuration.ReplicaTimeout = state.ReplicaTimeoutInSeconds
			}

			if d.HasChange("event_trigger_config") {
				model.Properties.Configuration.EventTriggerConfig = helpers.ExpandContainerAppJobConfigurationEventTriggerConfig(state.EventTriggerConfig)
			}

			if d.HasChange("manual_trigger_config") {
				model.Properties.Configuration.ManualTriggerConfig = helpers.ExpandContainerAppJobConfigurationManualTriggerConfig(state.ManualTriggerConfig)
			}

			if d.HasChange("schedule_trigger_config") {
				model.Properties.Configuration.ScheduleTriggerConfig = helpers.ExpandContainerAppJobConfigurationScheduleTriggerConfig(state.ScheduleTriggerConfig)
			}

			if d.HasChange("identity") {
				ident, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				model.Identity = ident
			}

			if d.HasChange("workload_profile_name") {
				model.Properties.WorkloadProfileName = pointer.To(state.WorkloadProfileName)
			}

			if d.HasChange("tags") {
				model.Tags = tags.Expand(state.Tags)
			}

			model.Properties.Template = helpers.ExpandContainerAppJobTemplate(state.Template)

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
