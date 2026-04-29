// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containerapps

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironmentsstorages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ContainerAppJobDataSource struct{}

type ContainerAppJobDataSourceModel struct {
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

var _ sdk.DataSource = ContainerAppJobDataSource{}

func (r ContainerAppJobDataSource) ModelObject() interface{} {
	return &ContainerAppJobDataSourceModel{}
}

func (r ContainerAppJobDataSource) ResourceType() string {
	return "azurerm_container_app_job"
}

func (r ContainerAppJobDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r ContainerAppJobDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"container_app_environment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"workload_profile_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"template": helpers.JobTemplateSchemaComputed(),

		"replica_retry_limit": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"replica_timeout_in_seconds": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"secret": helpers.SecretsDataSourceSchema(),

		"registry": helpers.ContainerAppRegistrySchemaComputed(),

		"event_trigger_config": helpers.EventTriggerConfigurationSchemaComputed(),

		"manual_trigger_config": helpers.ManualTriggerConfigurationSchemaComputed(),

		"schedule_trigger_config": helpers.ScheduleTriggerConfigurationSchemaComputed(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

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

		"tags": commonschema.TagsDataSource(),
	}
}

func (r ContainerAppJobDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.JobClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ContainerAppJobDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			id := jobs.NewJobID(subscriptionId, state.ResourceGroup, state.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("reading %s: %+v", id, err)
			}

			state.Name = id.JobName
			state.ResourceGroup = id.ResourceGroupName

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)

				if model.Identity != nil {
					ident, err := identity.FlattenSystemAndUserAssignedMapToModel(pointer.To((identity.SystemAndUserAssignedMap)(*model.Identity)))
					if err != nil {
						return err
					}
					state.Identity = pointer.From(ident)
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
						state.ReplicaRetryLimit = pointer.From(config.ReplicaRetryLimit)

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
					state.OutboundIPAddresses = pointer.From(props.OutboundIPAddresses)
					state.EventStreamEndpoint = pointer.From(props.EventStreamEndpoint)
				}
			}

			secretResp, err := client.ListSecrets(ctx, id)
			if err != nil {
				return fmt.Errorf("listing secrets for %s: %+v", id, err)
			}
			state.Secrets = helpers.FlattenContainerAppJobSecrets(secretResp.Model)

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}
