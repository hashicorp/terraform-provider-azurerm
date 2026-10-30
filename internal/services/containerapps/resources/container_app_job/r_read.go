// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_job

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironmentsstorages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
)

func (r Resource) Read() sdk.ResourceFunc {
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
						ident, err := identity.FlattenSystemAndUserAssignedMapToModel(pointer.To(identity.SystemAndUserAssignedMap(*model.Identity)))
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
					state.EventStreamEndpoint = pointer.From(props.EventStreamEndpoint)
					state.OutboundIPAddresses = pointer.From(props.OutboundIPAddresses)
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
