// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_app_job

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/jobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/helpers"
)

func (r Resource) Update() sdk.ResourceFunc {
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
