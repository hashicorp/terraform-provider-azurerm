// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_registry_task

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview/tasks"
)

func flattenRegistryTaskBaseImageTrigger(trigger *tasks.BaseImageTrigger, model ContainerRegistryTaskModel) []BaseImageTrigger {
	if trigger == nil {
		return []BaseImageTrigger{}
	}

	payloadType := ""
	if v := trigger.UpdateTriggerPayloadType; v != nil {
		payloadType = string(*v)
	}
	obj := BaseImageTrigger{
		Type:                     string(trigger.BaseImageTriggerType),
		Enabled:                  *trigger.Status == tasks.TriggerStatusEnabled,
		UpdateTriggerPayloadType: payloadType,
		Name:                     trigger.Name,
	}

	// UpdateTriggerEndpoint is not returned from API, setting it from config.
	if len(model.BaseImageTrigger) == 1 {
		obj.UpdateTriggerEndpoint = model.BaseImageTrigger[0].UpdateTriggerEndpoint
	}

	return []BaseImageTrigger{obj}
}

func flattenRegistryTaskSourceTriggers(triggers *[]tasks.SourceTrigger, model ContainerRegistryTaskModel) []SourceTrigger {
	if triggers == nil {
		return []SourceTrigger{}
	}
	out := make([]SourceTrigger, 0, len(*triggers))
	for i, trigger := range *triggers {
		obj := SourceTrigger{
			Enabled: *trigger.Status == tasks.TriggerStatusEnabled,
		}
		obj.Name = trigger.Name

		if trigger.SourceTriggerEvents != nil {
			events := make([]string, 0, len(trigger.SourceTriggerEvents))
			for _, event := range trigger.SourceTriggerEvents {
				events = append(events, string(event))
			}
			obj.Events = events
		}

		obj.SourceType = string(trigger.SourceRepository.SourceControlType)
		obj.RepositoryURL = trigger.SourceRepository.RepositoryURL
		if trigger.SourceRepository.Branch != nil {
			obj.Branch = *trigger.SourceRepository.Branch
		}

		// Auth is not returned from API, setting it from config.
		if len(model.SourceTrigger) > i {
			obj.Auth = model.SourceTrigger[i].Auth
		}

		out = append(out, obj)
	}
	return out
}

func flattenRegistryTaskTimerTriggers(triggers *[]tasks.TimerTrigger) []TimerTrigger {
	if triggers == nil {
		return []TimerTrigger{}
	}
	out := make([]TimerTrigger, 0, len(*triggers))
	for _, trigger := range *triggers {
		obj := TimerTrigger{
			Enabled: *trigger.Status == tasks.TriggerStatusEnabled,
		}
		obj.Name = trigger.Name
		obj.Schedule = trigger.Schedule
		out = append(out, obj)
	}
	return out
}

func flattenRegistryTaskDockerStep(step tasks.TaskStepProperties, model ContainerRegistryTaskModel) []DockerStep {
	if step == nil {
		return []DockerStep{}
	}

	dockerStep, ok := step.(tasks.DockerBuildStep)
	if !ok {
		return []DockerStep{}
	}

	obj := DockerStep{
		Arguments: flattenRegistryTaskArguments(dockerStep.Arguments),
	}

	if dockerStep.ContextPath != nil {
		obj.ContextPath = *dockerStep.ContextPath
	}
	obj.DockerfilePath = dockerStep.DockerFilePath

	if dockerStep.ImageNames != nil {
		obj.ImageNames = *dockerStep.ImageNames
	}
	if dockerStep.IsPushEnabled != nil {
		obj.IsPushEnabled = *dockerStep.IsPushEnabled
	}
	if dockerStep.NoCache != nil {
		obj.IsCacheEnabled = !(*dockerStep.NoCache)
	}
	if dockerStep.Target != nil {
		obj.Target = *dockerStep.Target
	}

	if len(model.DockerStep) == 1 {
		// The ContextAccessToken is sensitive and won't return from API, setting it from the config.
		obj.ContextAccessToken = model.DockerStep[0].ContextAccessToken

		// The SecretArguments is sensitive and won't return from API, setting it from the config.
		obj.SecretArguments = model.DockerStep[0].SecretArguments
	}

	return []DockerStep{obj}
}

func flattenRegistryTaskFileTaskStep(step tasks.TaskStepProperties, model ContainerRegistryTaskModel) []FileTaskStep {
	if step == nil {
		return []FileTaskStep{}
	}

	fileTaskStep, ok := step.(tasks.FileTaskStep)
	if !ok {
		return []FileTaskStep{}
	}

	obj := FileTaskStep{
		Values: flattenRegistryTaskValues(fileTaskStep.Values),
	}

	if fileTaskStep.ContextPath != nil {
		obj.ContextPath = *fileTaskStep.ContextPath
	}
	obj.TaskFilePath = fileTaskStep.TaskFilePath

	if fileTaskStep.ValuesFilePath != nil {
		obj.ValueFilePath = *fileTaskStep.ValuesFilePath
	}

	if len(model.FileTaskStep) == 1 {
		// The ContextAccessToken is sensitive and won't return from API, setting it from the config.
		obj.ContextAccessToken = model.FileTaskStep[0].ContextAccessToken

		// The SecretValues is sensitive and won't return from API, setting it from the config.
		obj.SecretValues = model.FileTaskStep[0].SecretValues
	}

	return []FileTaskStep{obj}
}

func flattenRegistryTaskEncodedTaskStep(step tasks.TaskStepProperties, model ContainerRegistryTaskModel) []EncodedTaskStep {
	if step == nil {
		return []EncodedTaskStep{}
	}

	encodedTaskStep, ok := step.(tasks.EncodedTaskStep)
	if !ok {
		return []EncodedTaskStep{}
	}

	obj := EncodedTaskStep{
		Values: flattenRegistryTaskValues(encodedTaskStep.Values),
	}

	if encodedTaskStep.ContextPath != nil {
		obj.ContextPath = *encodedTaskStep.ContextPath
	}
	if encodedTaskStep.ContextAccessToken != nil {
		obj.ContextAccessToken = *encodedTaskStep.ContextAccessToken
	}
	obj.TaskContent = encodedTaskStep.EncodedTaskContent

	if encodedTaskStep.EncodedValuesContent != nil {
		obj.ValueContent = *encodedTaskStep.EncodedValuesContent
	}

	if len(model.EncodedTaskStep) == 1 {
		// The ContextAccessToken is sensitive and won't return from API, setting it from the config.
		obj.ContextAccessToken = model.EncodedTaskStep[0].ContextAccessToken

		// The SecretValues is sensitive and won't return from API, setting it from the config.
		obj.SecretValues = model.EncodedTaskStep[0].SecretValues
	}

	return []EncodedTaskStep{obj}
}

func flattenRegistryTaskArguments(arguments *[]tasks.Argument) map[string]string {
	if arguments == nil {
		return map[string]string{}
	}

	args := map[string]string{}

	for _, argument := range *arguments {
		var (
			k        string
			v        string
			isSecret bool
		)
		k = argument.Name
		v = argument.Value

		if argument.IsSecret != nil {
			isSecret = *argument.IsSecret
		}

		// GET response won't return the value of secret arguments
		if isSecret {
			continue
		}

		args[k] = v
	}
	return args
}

func flattenRegistryTaskValues(values *[]tasks.SetValue) map[string]string {
	if values == nil {
		return map[string]string{}
	}

	vals := map[string]string{}

	for _, value := range *values {
		var (
			k        string
			v        string
			isSecret bool
		)
		k = value.Name
		v = value.Value

		if value.IsSecret != nil {
			isSecret = *value.IsSecret
		}

		if isSecret {
			// GET response won't return the value of secret values
			continue
		}
		vals[k] = v
	}
	return vals
}

func flattenRegistryTaskPlatform(platform *tasks.PlatformProperties) []Platform {
	if platform == nil {
		return []Platform{}
	}

	architecture := ""
	if v := platform.Architecture; v != nil {
		architecture = string(*v)
	}

	variant := ""
	if v := platform.Variant; v != nil {
		variant = string(*v)
	}

	return []Platform{{
		OS:           string(platform.Os),
		Architecture: architecture,
		Variant:      variant,
	}}
}

func flattenRegistryTaskCredentials(input *tasks.Credentials, model ContainerRegistryTaskModel) []RegistryCredential {
	if input == nil {
		return []RegistryCredential{}
	}

	// The customRegistryCredentials is sensitive and won't return from API, setting it from the config.
	var custom []CustomRegistryCredential
	if len(model.RegistryCredential) == 1 {
		custom = model.RegistryCredential[0].Custom
	}

	return []RegistryCredential{
		{
			Source: flattenSourceRegistryCredential(input.SourceRegistry),
			Custom: custom,
		},
	}
}

func flattenSourceRegistryCredential(input *tasks.SourceRegistryCredentials) []SourceRegistryCredential {
	if input == nil || input.LoginMode == nil {
		return []SourceRegistryCredential{}
	}

	return []SourceRegistryCredential{{LoginMode: string(*input.LoginMode)}}
}

func flattenRegistryTaskAgentProperties(input *tasks.AgentProperties) []AgentConfig {
	if input == nil {
		return []AgentConfig{}
	}

	return []AgentConfig{{CPU: pointer.From(input.Cpu)}}
}
