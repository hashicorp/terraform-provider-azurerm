// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_app_job

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

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

func (r DataSource) ModelObject() any {
	return &ContainerAppJobModel{}
}

func (r DataSource) ResourceType() string {
	return "azurerm_container_app_job"
}

func (r DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"container_app_environment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"event_stream_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"event_trigger_config": helpers.EventTriggerConfigurationSchemaComputed(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"manual_trigger_config": helpers.ManualTriggerConfigurationSchemaComputed(),

		"outbound_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"registry": helpers.ContainerAppRegistrySchemaComputed(),

		"replica_retry_limit": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"replica_timeout_in_seconds": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"schedule_trigger_config": helpers.ScheduleTriggerConfigurationSchemaComputed(),

		"secret": helpers.SecretsDataSourceSchema(),

		"template": helpers.JobTemplateSchemaComputed(),

		"workload_profile_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerApps.JobClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ContainerAppJobModel
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
					ident, err := identity.FlattenSystemAndUserAssignedMapToModel(pointer.To(identity.SystemAndUserAssignedMap(*model.Identity)))
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
