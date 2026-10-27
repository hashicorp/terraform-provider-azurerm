// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_task

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/tasks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/base64"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.ContainerRegistryClient_v2019_06_01_preview.Tasks
			registryClient := metadata.Client.Containers.ContainerRegistryClient.Registries

			var model ContainerRegistryTaskModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			registryId, err := registries.ParseRegistryID(model.ContainerRegistryId)
			if err != nil {
				return err
			}

			registry, err := registryClient.Get(ctx, *registryId)
			if err != nil {
				return fmt.Errorf("getting registry %s: %+v", registryId, err)
			}

			id := tasks.NewTaskID(registryId.SubscriptionId, registryId.ResourceGroupName, registryId.RegistryName, model.Name)

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

			status := tasks.TaskStatusDisabled
			if model.Enabled {
				status = tasks.TaskStatusEnabled
			}

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]any))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			params := tasks.Task{
				Properties: &tasks.TaskProperties{
					Platform:           expandRegistryTaskPlatform(model.Platform),
					Step:               expandRegistryTaskStep(model),
					Trigger:            expandRegistryTaskTrigger(model),
					Status:             pointer.To(status),
					IsSystemTask:       &model.IsSystemTask,
					Timeout:            pointer.To(model.TimeoutInSec),
					Credentials:        expandRegistryTaskCredentials(model.RegistryCredential),
					AgentConfiguration: expandRegistryTaskAgentProperties(model.AgentConfig),
				},

				// The location of the task must be the same as the registry, otherwise the API will raise error complaining can't find the registry.
				Location: location.Normalize(registry.Model.Location),
				Identity: expandedIdentity,
				Tags:     &model.Tags,
			}

			if model.AgentPoolName != "" {
				params.Properties.AgentPoolName = &model.AgentPoolName
			}
			if model.LogTemplate != "" {
				params.Properties.LogTemplate = &model.LogTemplate
			}

			if err := client.CreateCallbackThenPoll(ctx, id, params, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func expandRegistryTaskTrigger(model ContainerRegistryTaskModel) *tasks.TriggerProperties {
	baseImageTrigger := expandRegistryTaskBaseImageTrigger(model.BaseImageTrigger)
	sourceTriggers := expandRegistryTaskSourceTriggers(model.SourceTrigger)
	timerTriggers := expandRegistryTaskTimerTriggers(model.TimerTrigger)
	if baseImageTrigger == nil && sourceTriggers == nil && timerTriggers == nil {
		return nil
	}
	return &tasks.TriggerProperties{
		BaseImageTrigger: baseImageTrigger,
		SourceTriggers:   sourceTriggers,
		TimerTriggers:    timerTriggers,
	}
}

func expandRegistryTaskBaseImageTrigger(triggers []BaseImageTrigger) *tasks.BaseImageTrigger {
	if len(triggers) == 0 {
		return nil
	}

	trigger := triggers[0]
	status := tasks.TriggerStatusDisabled
	if trigger.Enabled {
		status = tasks.TriggerStatusEnabled
	}
	out := &tasks.BaseImageTrigger{
		Name:                 trigger.Name,
		BaseImageTriggerType: tasks.BaseImageTriggerType(trigger.Type),
		Status:               &status,
	}
	if trigger.UpdateTriggerEndpoint != "" {
		out.UpdateTriggerEndpoint = &trigger.UpdateTriggerEndpoint
	}
	if trigger.UpdateTriggerPayloadType != "" {
		out.UpdateTriggerPayloadType = pointer.ToEnum[tasks.UpdateTriggerPayloadType](trigger.UpdateTriggerPayloadType)
	}
	return out
}

func expandRegistryTaskSourceTriggers(triggers []SourceTrigger) *[]tasks.SourceTrigger {
	if len(triggers) == 0 {
		return nil
	}
	out := make([]tasks.SourceTrigger, 0, len(triggers))
	for _, trigger := range triggers {
		status := tasks.TriggerStatusDisabled
		if trigger.Enabled {
			status = tasks.TriggerStatusEnabled
		}
		sourceTrigger := tasks.SourceTrigger{
			Name:   trigger.Name,
			Status: &status,
			SourceRepository: tasks.SourceProperties{
				SourceControlType: tasks.SourceControlType(trigger.SourceType),
				RepositoryURL:     trigger.RepositoryURL,
			},
		}
		if len(trigger.Events) != 0 {
			events := make([]tasks.SourceTriggerEvent, 0, len(trigger.Events))
			for _, event := range trigger.Events {
				events = append(events, tasks.SourceTriggerEvent(event))
			}
			sourceTrigger.SourceTriggerEvents = events
		}

		if trigger.Branch != "" {
			sourceTrigger.SourceRepository.Branch = &trigger.Branch
		}
		if len(trigger.Auth) != 0 {
			sourceTrigger.SourceRepository.SourceControlAuthProperties = expandRegistryTaskAuthInfo(trigger.Auth[0])
		}
		out = append(out, sourceTrigger)
	}
	return &out
}

func expandRegistryTaskAuthInfo(auth Auth) *tasks.AuthInfo {
	out := tasks.AuthInfo{
		TokenType: tasks.TokenType(auth.TokenType),
		Token:     auth.Token,
	}
	if auth.RefreshToken != "" {
		out.RefreshToken = &auth.RefreshToken
	}
	if auth.Scope != "" {
		out.Scope = &auth.Scope
	}
	if auth.ExpireInSec != 0 {
		out.ExpiresIn = pointer.To(auth.ExpireInSec)
	}
	return &out
}

func expandRegistryTaskTimerTriggers(triggers []TimerTrigger) *[]tasks.TimerTrigger {
	if len(triggers) == 0 {
		return nil
	}
	out := make([]tasks.TimerTrigger, 0, len(triggers))
	for _, trigger := range triggers {
		status := tasks.TriggerStatusDisabled
		if trigger.Enabled {
			status = tasks.TriggerStatusEnabled
		}
		timerTrigger := tasks.TimerTrigger{
			Name:     trigger.Name,
			Schedule: trigger.Schedule,
			Status:   &status,
		}
		out = append(out, timerTrigger)
	}
	return &out
}

func expandRegistryTaskStep(model ContainerRegistryTaskModel) tasks.TaskStepProperties {
	switch {
	case len(model.DockerStep) != 0:
		return expandRegistryTaskDockerStep(model.DockerStep[0])
	case len(model.FileTaskStep) != 0:
		return expandRegistryTaskFileTaskStep(model.FileTaskStep[0])
	case len(model.EncodedTaskStep) != 0:
		return expandRegistryTaskEncodedTaskStep(model.EncodedTaskStep[0])
	}
	return nil
}

func expandRegistryTaskDockerStep(step DockerStep) tasks.DockerBuildStep {
	out := tasks.DockerBuildStep{
		DockerFilePath: step.DockerfilePath,
		IsPushEnabled:  &step.IsPushEnabled,
		NoCache:        pointer.To(!step.IsCacheEnabled),
		Arguments:      expandRegistryTaskArguments(step.Arguments, step.SecretArguments),
	}
	if step.ContextPath != "" {
		out.ContextPath = &step.ContextPath
	}
	if step.ContextAccessToken != "" {
		out.ContextAccessToken = &step.ContextAccessToken
	}
	if len(step.ImageNames) != 0 {
		out.ImageNames = &step.ImageNames
	}
	if step.Target != "" {
		out.Target = &step.Target
	}

	return out
}

func expandRegistryTaskFileTaskStep(step FileTaskStep) tasks.FileTaskStep {
	out := tasks.FileTaskStep{
		TaskFilePath: step.TaskFilePath,
		Values:       expandRegistryTaskValues(step.Values, step.SecretValues),
	}
	if step.ValueFilePath != "" {
		out.ValuesFilePath = &step.ValueFilePath
	}
	if step.ContextPath != "" {
		out.ContextPath = &step.ContextPath
	}
	if step.ContextAccessToken != "" {
		out.ContextAccessToken = &step.ContextAccessToken
	}
	return out
}

func expandRegistryTaskEncodedTaskStep(step EncodedTaskStep) tasks.EncodedTaskStep {
	out := tasks.EncodedTaskStep{
		EncodedTaskContent: base64.EncodeIfNot(step.TaskContent),
		Values:             expandRegistryTaskValues(step.Values, step.SecretValues),
	}
	if step.ContextPath != "" {
		out.ContextPath = &step.ContextPath
	}
	if step.ContextAccessToken != "" {
		out.ContextAccessToken = &step.ContextAccessToken
	}
	if step.ValueContent != "" {
		out.EncodedValuesContent = pointer.To(base64.EncodeIfNot(step.ValueContent))
	}
	return out
}

func expandRegistryTaskArguments(arguments map[string]string, secretArguments map[string]string) *[]tasks.Argument {
	if len(arguments) == 0 && len(secretArguments) == 0 {
		return nil
	}
	out := make([]tasks.Argument, 0, len(arguments)+len(secretArguments))
	for k, v := range arguments {
		out = append(out, tasks.Argument{
			Name:     k,
			Value:    v,
			IsSecret: pointer.To(false),
		})
	}
	for k, v := range secretArguments {
		out = append(out, tasks.Argument{
			Name:     k,
			Value:    v,
			IsSecret: pointer.To(true),
		})
	}
	return &out
}

func expandRegistryTaskValues(values map[string]string, secretValues map[string]string) *[]tasks.SetValue {
	if len(values) == 0 && len(secretValues) == 0 {
		return nil
	}
	out := make([]tasks.SetValue, 0, len(values)+len(secretValues))
	for k, v := range values {
		out = append(out, tasks.SetValue{
			Name:     k,
			Value:    v,
			IsSecret: pointer.To(false),
		})
	}
	for k, v := range secretValues {
		out = append(out, tasks.SetValue{
			Name:     k,
			Value:    v,
			IsSecret: pointer.To(true),
		})
	}
	return &out
}

func expandRegistryTaskPlatform(input []Platform) *tasks.PlatformProperties {
	if len(input) == 0 {
		return nil
	}
	platform := input[0]
	out := &tasks.PlatformProperties{
		Os: tasks.OS(platform.OS),
	}
	if arch := platform.Architecture; arch != "" {
		out.Architecture = pointer.ToEnum[tasks.Architecture](arch)
	}
	if variant := platform.Variant; variant != "" {
		out.Variant = pointer.ToEnum[tasks.Variant](variant)
	}
	return out
}

func expandRegistryTaskCredentials(input []RegistryCredential) *tasks.Credentials {
	if len(input) == 0 {
		return nil
	}

	return &tasks.Credentials{
		SourceRegistry:   expandSourceRegistryCredential(input[0].Source),
		CustomRegistries: pointer.To(expandCustomRegistryCredential(input[0].Custom)),
	}
}

func expandSourceRegistryCredential(input []SourceRegistryCredential) *tasks.SourceRegistryCredentials {
	if len(input) == 0 {
		return nil
	}

	return &tasks.SourceRegistryCredentials{LoginMode: pointer.ToEnum[tasks.SourceRegistryLoginMode](input[0].LoginMode)}
}

func expandCustomRegistryCredential(input []CustomRegistryCredential) map[string]tasks.CustomRegistryCredentials {
	if len(input) == 0 {
		return nil
	}

	out := map[string]tasks.CustomRegistryCredentials{}
	for _, credential := range input {
		cred := tasks.CustomRegistryCredentials{}

		if credential.UserName != "" {
			usernameType := tasks.SecretObjectTypeOpaque
			if _, err := keyvault.ParseNestedItemID(credential.UserName, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeAny); err == nil {
				usernameType = tasks.SecretObjectTypeVaultsecret
			}
			cred.UserName = &tasks.SecretObject{
				Value: pointer.To(credential.UserName),
				Type:  &usernameType,
			}
		}
		if credential.Password != "" {
			passwordType := tasks.SecretObjectTypeOpaque
			if _, err := keyvault.ParseNestedItemID(credential.Password, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeAny); err == nil {
				passwordType = tasks.SecretObjectTypeVaultsecret
			}
			cred.Password = &tasks.SecretObject{
				Value: pointer.To(credential.Password),
				Type:  &passwordType,
			}
		}
		if credential.Identity != "" {
			cred.Identity = pointer.To(credential.Identity)
		}
		out[credential.LoginServer] = cred
	}
	return out
}

func expandRegistryTaskAgentProperties(input []AgentConfig) *tasks.AgentProperties {
	if len(input) == 0 {
		return nil
	}

	agentConfig := input[0]
	return &tasks.AgentProperties{Cpu: pointer.To(agentConfig.CPU)}
}

func patchRegistryTaskTriggerSourceTrigger(triggers []tasks.SourceTrigger, model ContainerRegistryTaskModel) *[]tasks.SourceTrigger {
	if len(triggers) != len(model.SourceTrigger) {
		return &triggers
	}

	result := make([]tasks.SourceTrigger, len(triggers))
	for i, trigger := range model.SourceTrigger {
		t := triggers[i]
		if len(trigger.Auth) == 0 {
			result[i] = t
			continue
		}

		t.SourceRepository.SourceControlAuthProperties = expandRegistryTaskAuthInfo(trigger.Auth[0])
		result[i] = t
	}
	return &result
}
