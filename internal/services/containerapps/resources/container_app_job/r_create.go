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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
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

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, job, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
