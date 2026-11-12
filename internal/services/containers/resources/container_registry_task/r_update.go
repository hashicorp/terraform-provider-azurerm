// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_task

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/tasks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := tasks.ParseTaskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Containers.ContainerRegistryClient_v2019_06_01_preview.Tasks
			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: +%v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("model is nil for %s", *id)
			}

			var model ContainerRegistryTaskModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("platform") {
				existing.Model.Properties.Platform = expandRegistryTaskPlatform(model.Platform)
			}
			if metadata.ResourceData.HasChanges("docker_step", "file_step", "encoded_step") {
				existing.Model.Properties.Step = expandRegistryTaskStep(model)
			}

			if metadata.ResourceData.HasChanges("base_image_trigger", "source_trigger", "timer_trigger") {
				existing.Model.Properties.Trigger = expandRegistryTaskTrigger(model)
			}

			if existing.Model.Properties.Trigger != nil {
				if !metadata.ResourceData.HasChange("source_trigger") && existing.Model.Properties.Trigger.SourceTriggers != nil {
					// For update that is not affecting source_trigger, we need to patch the source_trigger to include the properties missing in the response of GET.
					existing.Model.Properties.Trigger.SourceTriggers = patchRegistryTaskTriggerSourceTrigger(*existing.Model.Properties.Trigger.SourceTriggers, model)
				}
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]any))
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				existing.Model.Identity = expandedIdentity
			}

			// Deliberately always set "registry_credential" as the custom registry's credentials are not returned by API, but are required for a PUT request.
			existing.Model.Properties.Credentials = expandRegistryTaskCredentials(model.RegistryCredential)

			if metadata.ResourceData.HasChange("agent_setting") {
				existing.Model.Properties.AgentConfiguration = expandRegistryTaskAgentProperties(model.AgentConfig)
			}
			if metadata.ResourceData.HasChange("agent_pool_name") && model.AgentPoolName != "" {
				existing.Model.Properties.AgentPoolName = &model.AgentPoolName
			}
			if metadata.ResourceData.HasChange("enabled") {
				status := tasks.TaskStatusDisabled
				if model.Enabled {
					status = tasks.TaskStatusEnabled
				}
				existing.Model.Properties.Status = pointer.To(status)
			}
			if metadata.ResourceData.HasChange("log_template") && model.LogTemplate != "" {
				existing.Model.Properties.LogTemplate = &model.LogTemplate
			}
			if metadata.ResourceData.HasChange("timeout_in_seconds") {
				existing.Model.Properties.Timeout = pointer.To(model.TimeoutInSec)
			}
			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = &model.Tags
			}

			// Due to the fact that the service doesn't honour explicitly set to null fields in the PATCH request,
			// we can not use PATCH (i.e. the Update) here.
			if err := client.CreateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
