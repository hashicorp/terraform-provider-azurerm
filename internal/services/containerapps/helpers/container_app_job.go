// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"reflect"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/jobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type JobTemplateModel struct {
	Containers     []Container       `tfschema:"container"`
	InitContainers []BaseContainer   `tfschema:"init_container"`
	Volumes        []ContainerVolume `tfschema:"volume"`
}

type EventTriggerConfiguration struct {
	Parallelism            int64        `tfschema:"parallelism"`
	ReplicaCompletionCount int64        `tfschema:"replica_completion_count"`
	Scale                  []ScaleModel `tfschema:"scale"`
}

type ScheduleTriggerConfiguration struct {
	CronExpression         string `tfschema:"cron_expression"`
	Parallelism            int64  `tfschema:"parallelism"`
	ReplicaCompletionCount int64  `tfschema:"replica_completion_count"`
}

type ManualTriggerConfiguration struct {
	Parallelism            int64 `tfschema:"parallelism"`
	ReplicaCompletionCount int64 `tfschema:"replica_completion_count"`
}

type ScaleModel struct {
	MaxExecutions   int64       `tfschema:"max_executions"`
	MinExecutions   int64       `tfschema:"min_executions"`
	PollingInterval int64       `tfschema:"polling_interval_in_seconds"`
	Rules           []ScaleRule `tfschema:"rules"`
}

type ScaleRule struct {
	Auth     []ScaleRuleAuth        `tfschema:"authentication"`
	Metadata map[string]interface{} `tfschema:"metadata"`
	Name     string                 `tfschema:"name"`
	Type     string                 `tfschema:"custom_rule_type"`
}

type ScaleRuleAuth struct {
	SecretReference  string `tfschema:"secret_name"`
	TriggerParameter string `tfschema:"trigger_parameter"`
}

func JobTemplateSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"container": ContainerAppContainerSchema(),

				"init_container": InitContainerAppContainerSchema(),

				"volume": ContainerVolumeSchema(),
			},
		},
	}
}

func ContainerAppsJobsScaleSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"max_executions": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      100,
					ValidateFunc: validation.IntAtLeast(1),
				},

				"min_executions": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      0,
					ValidateFunc: validation.IntAtLeast(0),
				},

				"polling_interval_in_seconds": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					Default:      30,
					ValidateFunc: validation.IntAtLeast(1),
				},

				"rules": CustomScaleRuleSchema(),
			},
		},
	}
}

func ExpandContainerAppJobSecrets(input []Secret) *[]jobs.Secret {
	if len(input) == 0 {
		return nil
	}

	result := make([]jobs.Secret, 0)

	for _, v := range input {
		result = append(result, jobs.Secret{
			Identity:    pointer.To(v.Identity),
			KeyVaultURL: pointer.To(v.KeyVaultSecretId),
			Name:        pointer.To(v.Name),
			Value:       pointer.To(v.Value),
		})
	}

	return &result
}

func ExpandContainerAppJobRegistries(input []Registry) (*[]jobs.RegistryCredentials, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]jobs.RegistryCredentials, 0)

	for _, v := range input {
		if err := ValidateContainerAppRegistry(v); err != nil {
			return nil, err
		}
		result = append(result, jobs.RegistryCredentials{
			Identity:          pointer.To(v.Identity),
			PasswordSecretRef: pointer.To(v.PasswordSecretRef),
			Server:            pointer.To(v.Server),
			Username:          pointer.To(v.UserName),
		})
	}

	return &result, nil
}

func ExpandContainerAppJobConfigurationManualTriggerConfig(input []ManualTriggerConfiguration) *jobs.JobConfigurationManualTriggerConfig {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	var manualTriggerConfig jobs.JobConfigurationManualTriggerConfig

	if v.Parallelism != 0 {
		manualTriggerConfig.Parallelism = pointer.To(v.Parallelism)
	}

	if v.ReplicaCompletionCount != 0 {
		manualTriggerConfig.ReplicaCompletionCount = pointer.To(v.ReplicaCompletionCount)
	}

	return pointer.To(manualTriggerConfig)
}

func ExpandContainerAppJobConfigurationScheduleTriggerConfig(input []ScheduleTriggerConfiguration) *jobs.JobConfigurationScheduleTriggerConfig {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	var scheduleTriggerConfig jobs.JobConfigurationScheduleTriggerConfig

	scheduleTriggerConfig.CronExpression = v.CronExpression

	if v.Parallelism != 0 {
		scheduleTriggerConfig.Parallelism = pointer.To(v.Parallelism)
	}

	if v.ReplicaCompletionCount != 0 {
		scheduleTriggerConfig.ReplicaCompletionCount = pointer.To(v.ReplicaCompletionCount)
	}

	return pointer.To(scheduleTriggerConfig)
}

func ExpandContainerAppJobConfigurationEventTriggerConfig(input []EventTriggerConfiguration) *jobs.JobConfigurationEventTriggerConfig {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	var eventTriggerConfig jobs.JobConfigurationEventTriggerConfig

	if v.Parallelism != 0 {
		eventTriggerConfig.Parallelism = pointer.To(v.Parallelism)
	}

	if v.ReplicaCompletionCount != 0 {
		eventTriggerConfig.ReplicaCompletionCount = pointer.To(v.ReplicaCompletionCount)
	}

	eventTriggerConfig.Scale = ExpandContainerAppJobScale(v.Scale)

	return pointer.To(eventTriggerConfig)
}

func ExpandContainerAppJobScale(input []ScaleModel) *jobs.JobScale {
	if len(input) == 0 {
		return nil
	}
	v := input[0]
	var scale jobs.JobScale

	if v.MaxExecutions != 0 {
		scale.MaxExecutions = pointer.To(v.MaxExecutions)
	}

	if v.MinExecutions != 0 {
		scale.MinExecutions = pointer.To(v.MinExecutions)
	}

	if v.PollingInterval != 0 {
		scale.PollingInterval = pointer.To(v.PollingInterval)
	}

	scale.Rules = ExpandContainerAppJobScaleRules(v.Rules)

	return pointer.To(scale)
}

func ExpandContainerAppJobScaleRules(input []ScaleRule) *[]jobs.JobScaleRule {
	if len(input) == 0 {
		return nil
	}

	rules := make([]jobs.JobScaleRule, 0)
	for _, v := range input {
		var rule jobs.JobScaleRule

		rule.Auth = ExpandContainerAppJobScaleRulesAuth(v.Auth)

		if v.Metadata != nil {
			metadata := reflect.ValueOf(v.Metadata)
			rule.Metadata = pointer.To(metadata.Interface())
		}

		if v.Name != "" {
			rule.Name = pointer.To(v.Name)
		}

		if v.Type != "" {
			rule.Type = pointer.To(v.Type)
		}

		rules = append(rules, rule)
	}

	return pointer.To(rules)
}

func ExpandContainerAppJobScaleRulesAuth(input []ScaleRuleAuth) *[]jobs.ScaleRuleAuth {
	if len(input) == 0 {
		return nil
	}
	auth := make([]jobs.ScaleRuleAuth, 0)
	for _, v := range input {
		var ruleAuth jobs.ScaleRuleAuth

		if v.SecretReference != "" {
			ruleAuth.SecretRef = pointer.To(v.SecretReference)
		}

		if v.TriggerParameter != "" {
			ruleAuth.TriggerParameter = pointer.To(v.TriggerParameter)
		}

		auth = append(auth, ruleAuth)
	}

	return pointer.To(auth)
}

func ExpandContainerAppJobTemplate(input []JobTemplateModel) *jobs.JobTemplate {
	if len(input) != 1 {
		return nil
	}
	v := input[0]
	template := &jobs.JobTemplate{
		Containers:     expandContainerAppJobContainers(v.Containers),
		InitContainers: expandInitContainerAppJobContainers(v.InitContainers),
		Volumes:        expandContainerAppJobVolumes(v.Volumes),
	}

	return template
}

func FlattenContainerAppJobTemplate(input *jobs.JobTemplate) []JobTemplateModel {
	if input == nil {
		return []JobTemplateModel{}
	}

	template := JobTemplateModel{
		Containers:     flattenContainerAppJobContainers(input.Containers),
		InitContainers: flattenInitContainerAppJobContainers(input.InitContainers),
		Volumes:        flattenContainerAppJobVolumes(input.Volumes),
	}

	return []JobTemplateModel{template}
}

func expandContainerAppJobContainers(input []Container) *[]jobs.Container {
	if input == nil {
		return nil
	}

	result := make([]jobs.Container, 0)

	for _, v := range input {
		container := jobs.Container{
			Env:    expandContainerJobEnvVar(v),
			Image:  pointer.To(v.Image),
			Name:   pointer.To(v.Name),
			Probes: expandJobContainerProbes(v),
			Resources: &jobs.ContainerResources{
				Cpu:              pointer.To(v.CPU),
				EphemeralStorage: pointer.To(v.EphemeralStorage),
				Memory:           pointer.To(v.Memory),
			},
			VolumeMounts: expandContainerJobVolumeMounts(v.VolumeMounts),
		}
		if len(v.Args) != 0 {
			container.Args = pointer.To(v.Args)
		}
		if len(v.Command) != 0 {
			container.Command = pointer.To(v.Command)
		}

		result = append(result, container)
	}

	return &result
}

func expandInitContainerAppJobContainers(input []BaseContainer) *[]jobs.BaseContainer {
	if input == nil {
		return nil
	}

	result := make([]jobs.BaseContainer, 0)

	for _, v := range input {
		container := jobs.BaseContainer{
			Env:   expandInitContainerJobEnvVar(v),
			Image: pointer.To(v.Image),
			Name:  pointer.To(v.Name),
			Resources: &jobs.ContainerResources{
				Cpu:              pointer.To(v.CPU),
				EphemeralStorage: pointer.To(v.EphemeralStorage),
				Memory:           pointer.To(v.Memory),
			},
			VolumeMounts: expandContainerJobVolumeMounts(v.VolumeMounts),
		}
		if len(v.Args) != 0 {
			container.Args = pointer.To(v.Args)
		}
		if len(v.Command) != 0 {
			container.Command = pointer.To(v.Command)
		}

		result = append(result, container)
	}

	return pointer.To(result)
}

func flattenContainerAppJobContainers(input *[]jobs.Container) []Container {
	if input == nil || len(*input) == 0 {
		return []Container{}
	}
	result := make([]Container, 0)
	for _, v := range *input {
		container := Container{
			Name:         pointer.From(v.Name),
			Image:        pointer.From(v.Image),
			Args:         pointer.From(v.Args),
			Command:      pointer.From(v.Command),
			Env:          flattenContainerJobEnvVar(v.Env),
			VolumeMounts: flattenContainerJobVolumeMounts(v.VolumeMounts),
		}
		if v.Probes != nil {
			for _, p := range *v.Probes {
				switch *p.Type {
				case jobs.TypeLiveness:
					container.LivenessProbe = flattenContainerAppJobLivenessProbe(p)
				case jobs.TypeReadiness:
					container.ReadinessProbe = flattenContainerAppJobReadinessProbe(p)
				case jobs.TypeStartup:
					container.StartupProbe = flattenContainerAppJobStartupProbe(p)
				}
			}
		}

		if resources := v.Resources; resources != nil {
			container.CPU = pointer.From(resources.Cpu)
			container.EphemeralStorage = pointer.From(resources.EphemeralStorage)
			container.Memory = pointer.From(resources.Memory)
		}

		result = append(result, container)
	}
	return result
}

func flattenInitContainerAppJobContainers(input *[]jobs.BaseContainer) []BaseContainer {
	if input == nil || len(*input) == 0 {
		return []BaseContainer{}
	}
	result := make([]BaseContainer, 0)
	for _, v := range *input {
		container := BaseContainer{
			Name:         pointer.From(v.Name),
			Image:        pointer.From(v.Image),
			Args:         pointer.From(v.Args),
			Command:      pointer.From(v.Command),
			Env:          flattenContainerJobEnvVar(v.Env),
			VolumeMounts: flattenContainerJobVolumeMounts(v.VolumeMounts),
		}

		if resources := v.Resources; resources != nil {
			container.CPU = pointer.From(resources.Cpu)
			container.EphemeralStorage = pointer.From(resources.EphemeralStorage)
			container.Memory = pointer.From(resources.Memory)
		}

		result = append(result, container)
	}
	return result
}

func expandContainerJobEnvVar(input Container) *[]jobs.EnvironmentVar {
	envs := make([]jobs.EnvironmentVar, 0)
	if len(input.Env) == 0 {
		return &envs
	}

	for _, v := range input.Env {
		env := jobs.EnvironmentVar{
			Name: pointer.To(v.Name),
		}
		if v.SecretReference != "" {
			env.SecretRef = pointer.To(v.SecretReference)
		} else {
			env.Value = pointer.To(v.Value)
		}

		envs = append(envs, env)
	}

	return &envs
}

func expandInitContainerJobEnvVar(input BaseContainer) *[]jobs.EnvironmentVar {
	envs := make([]jobs.EnvironmentVar, 0)
	if len(input.Env) == 0 {
		return &envs
	}

	for _, v := range input.Env {
		env := jobs.EnvironmentVar{
			Name: pointer.To(v.Name),
		}
		if v.SecretReference != "" {
			env.SecretRef = pointer.To(v.SecretReference)
		} else {
			env.Value = pointer.To(v.Value)
		}

		envs = append(envs, env)
	}

	return &envs
}

func expandContainerAppJobVolumes(input []ContainerVolume) *[]jobs.Volume {
	if input == nil {
		return nil
	}

	volumes := make([]jobs.Volume, 0)

	for _, v := range input {
		volume := jobs.Volume{
			Name: pointer.To(v.Name),
		}
		if v.StorageName != "" {
			volume.StorageName = pointer.To(v.StorageName)
		}
		if v.StorageType != "" {
			storageType := jobs.StorageType(v.StorageType)
			volume.StorageType = &storageType
		}
		if v.MountOptions != "" {
			volume.MountOptions = pointer.To(v.MountOptions)
		}
		volumes = append(volumes, volume)
	}

	return &volumes
}

func expandJobContainerProbes(input Container) *[]jobs.ContainerAppProbe {
	probes := make([]jobs.ContainerAppProbe, 0)

	if len(input.LivenessProbe) == 1 {
		probes = append(probes, expandContainerAppJobLivenessProbe(input.LivenessProbe[0]))
	}

	if len(input.ReadinessProbe) == 1 {
		probes = append(probes, expandContainerAppJobReadinessProbe(input.ReadinessProbe[0]))
	}

	if len(input.StartupProbe) == 1 {
		probes = append(probes, expandContainerAppJobStartupProbe(input.StartupProbe[0]))
	}

	return &probes
}

func expandContainerJobVolumeMounts(input []ContainerVolumeMount) *[]jobs.VolumeMount {
	if input == nil {
		return nil
	}
	volumeMounts := make([]jobs.VolumeMount, 0)
	for _, v := range input {
		volumeMounts = append(volumeMounts, jobs.VolumeMount{
			MountPath:  pointer.To(v.Path),
			VolumeName: pointer.To(v.Name),
		})
	}

	return &volumeMounts
}

func expandContainerAppJobLivenessProbe(input ContainerAppLivenessProbe) jobs.ContainerAppProbe {
	probeType := jobs.TypeLiveness
	result := jobs.ContainerAppProbe{
		Type:                &probeType,
		InitialDelaySeconds: pointer.To(input.InitialDelay),
		PeriodSeconds:       pointer.To(input.Interval),
		TimeoutSeconds:      pointer.To(input.Timeout),
		FailureThreshold:    pointer.To(input.FailureThreshold),
	}

	switch p := strings.ToUpper(input.Transport); p {
	case "HTTP", "HTTPS":
		scheme := jobs.Scheme(p)
		result.HTTPGet = &jobs.ContainerAppProbeHTTPGet{
			Host:   pointer.To(input.Host),
			Path:   pointer.To(input.Path),
			Port:   input.Port,
			Scheme: &scheme,
		}
		if input.Headers != nil {
			headers := make([]jobs.ContainerAppProbeHTTPGetHTTPHeadersInlined, 0)

			for _, h := range input.Headers {
				headers = append(headers, jobs.ContainerAppProbeHTTPGetHTTPHeadersInlined{
					Name:  h.Name,
					Value: h.Value,
				})
			}
			result.HTTPGet.HTTPHeaders = &headers
		}

	default:
		result.TcpSocket = &jobs.ContainerAppProbeTcpSocket{
			Host: pointer.To(input.Host),
			Port: input.Port,
		}
	}

	return result
}

func expandContainerAppJobReadinessProbe(input ContainerAppReadinessProbe) jobs.ContainerAppProbe {
	probeType := jobs.TypeReadiness
	result := jobs.ContainerAppProbe{
		Type:             &probeType,
		PeriodSeconds:    pointer.To(input.Interval),
		TimeoutSeconds:   pointer.To(input.Timeout),
		FailureThreshold: pointer.To(input.FailureThreshold),
		SuccessThreshold: pointer.To(input.SuccessThreshold),
	}

	switch p := strings.ToUpper(input.Transport); p {
	case "HTTP", "HTTPS":
		scheme := jobs.Scheme(p)
		result.HTTPGet = &jobs.ContainerAppProbeHTTPGet{
			Host:   pointer.To(input.Host),
			Path:   pointer.To(input.Path),
			Port:   input.Port,
			Scheme: &scheme,
		}
		if input.Headers != nil {
			headers := make([]jobs.ContainerAppProbeHTTPGetHTTPHeadersInlined, 0)

			for _, h := range input.Headers {
				headers = append(headers, jobs.ContainerAppProbeHTTPGetHTTPHeadersInlined{
					Name:  h.Name,
					Value: h.Value,
				})
			}
			result.HTTPGet.HTTPHeaders = &headers
		}

	default:
		result.TcpSocket = &jobs.ContainerAppProbeTcpSocket{
			Host: pointer.To(input.Host),
			Port: input.Port,
		}
	}

	return result
}

func expandContainerAppJobStartupProbe(input ContainerAppStartupProbe) jobs.ContainerAppProbe {
	probeType := jobs.TypeStartup
	result := jobs.ContainerAppProbe{
		Type:             &probeType,
		PeriodSeconds:    pointer.To(input.Interval),
		TimeoutSeconds:   pointer.To(input.Timeout),
		FailureThreshold: pointer.To(input.FailureThreshold),
	}

	switch p := strings.ToUpper(input.Transport); p {
	case "HTTP", "HTTPS":
		scheme := jobs.Scheme(p)
		result.HTTPGet = &jobs.ContainerAppProbeHTTPGet{
			Host:   pointer.To(input.Host),
			Path:   pointer.To(input.Path),
			Port:   input.Port,
			Scheme: &scheme,
		}
		if input.Headers != nil {
			headers := make([]jobs.ContainerAppProbeHTTPGetHTTPHeadersInlined, 0)

			for _, h := range input.Headers {
				headers = append(headers, jobs.ContainerAppProbeHTTPGetHTTPHeadersInlined{
					Name:  h.Name,
					Value: h.Value,
				})
			}
			result.HTTPGet.HTTPHeaders = &headers
		}

	default:
		result.TcpSocket = &jobs.ContainerAppProbeTcpSocket{
			Host: pointer.To(input.Host),
			Port: input.Port,
		}
	}

	return result
}

func UnpackContainerJobSecretsCollection(input *jobs.JobSecretsCollection) *[]jobs.Secret {
	if input == nil || len(input.Value) == 0 {
		return nil
	}

	result := make([]jobs.Secret, 0)
	result = append(result, input.Value...)

	return &result
}

func flattenContainerJobEnvVar(input *[]jobs.EnvironmentVar) []ContainerEnvVar {
	if input == nil || len(*input) == 0 {
		return []ContainerEnvVar{}
	}

	result := make([]ContainerEnvVar, 0)

	for _, v := range *input {
		result = append(result, ContainerEnvVar{
			Name:            pointer.From(v.Name),
			SecretReference: pointer.From(v.SecretRef),
			Value:           pointer.From(v.Value),
		})
	}

	return result
}

func flattenContainerJobVolumeMounts(input *[]jobs.VolumeMount) []ContainerVolumeMount {
	if input == nil || len(*input) == 0 {
		return []ContainerVolumeMount{}
	}

	result := make([]ContainerVolumeMount, 0)
	for _, v := range *input {
		result = append(result, ContainerVolumeMount{
			Name: pointer.From(v.VolumeName),
			Path: pointer.From(v.MountPath),
		})
	}

	return result
}

func flattenContainerAppJobLivenessProbe(input jobs.ContainerAppProbe) []ContainerAppLivenessProbe {
	result := make([]ContainerAppLivenessProbe, 0)
	probe := ContainerAppLivenessProbe{
		InitialDelay:           pointer.From(input.InitialDelaySeconds),
		Interval:               pointer.From(input.PeriodSeconds),
		Timeout:                pointer.From(input.TimeoutSeconds),
		FailureThreshold:       pointer.From(input.FailureThreshold),
		TerminationGracePeriod: pointer.From(input.TerminationGracePeriodSeconds),
	}
	if httpGet := input.HTTPGet; httpGet != nil {
		if httpGet.Scheme != nil {
			probe.Transport = string(*httpGet.Scheme)
		}
		probe.Host = pointer.From(httpGet.Host)
		probe.Port = httpGet.Port
		probe.Path = pointer.From(httpGet.Path)

		if httpGet.HTTPHeaders != nil {
			headers := make([]HttpHeader, 0)
			for _, h := range *httpGet.HTTPHeaders {
				headers = append(headers, HttpHeader{
					Name:  h.Name,
					Value: h.Value,
				})
			}
			probe.Headers = headers
		}
	}

	if tcpSocket := input.TcpSocket; tcpSocket != nil {
		probe.Transport = "TCP"
		probe.Host = pointer.From(tcpSocket.Host)
		probe.Port = tcpSocket.Port
	}

	result = append(result, probe)

	return result
}

func flattenContainerAppJobReadinessProbe(input jobs.ContainerAppProbe) []ContainerAppReadinessProbe {
	result := make([]ContainerAppReadinessProbe, 0)
	probe := ContainerAppReadinessProbe{
		Interval:         pointer.From(input.PeriodSeconds),
		Timeout:          pointer.From(input.TimeoutSeconds),
		FailureThreshold: pointer.From(input.FailureThreshold),
		SuccessThreshold: pointer.From(input.SuccessThreshold),
	}

	if httpGet := input.HTTPGet; httpGet != nil {
		if httpGet.Scheme != nil {
			probe.Transport = string(*httpGet.Scheme)
		}
		probe.Host = pointer.From(httpGet.Host)
		probe.Port = httpGet.Port
		probe.Path = pointer.From(httpGet.Path)

		if httpGet.HTTPHeaders != nil {
			headers := make([]HttpHeader, 0)
			for _, h := range *httpGet.HTTPHeaders {
				headers = append(headers, HttpHeader{
					Name:  h.Name,
					Value: h.Value,
				})
			}
			probe.Headers = headers
		}
	}

	if tcpSocket := input.TcpSocket; tcpSocket != nil {
		probe.Transport = "TCP"
		probe.Host = pointer.From(tcpSocket.Host)
		probe.Port = tcpSocket.Port
	}

	result = append(result, probe)

	return result
}

func flattenContainerAppJobStartupProbe(input jobs.ContainerAppProbe) []ContainerAppStartupProbe {
	result := make([]ContainerAppStartupProbe, 0)
	probe := ContainerAppStartupProbe{
		Interval:               pointer.From(input.PeriodSeconds),
		Timeout:                pointer.From(input.TimeoutSeconds),
		FailureThreshold:       pointer.From(input.FailureThreshold),
		TerminationGracePeriod: pointer.From(input.TerminationGracePeriodSeconds),
	}

	if httpGet := input.HTTPGet; httpGet != nil {
		if httpGet.Scheme != nil {
			probe.Transport = string(*httpGet.Scheme)
		}
		probe.Host = pointer.From(httpGet.Host)
		probe.Port = httpGet.Port
		probe.Path = pointer.From(httpGet.Path)

		if httpGet.HTTPHeaders != nil {
			headers := make([]HttpHeader, 0)
			for _, h := range *httpGet.HTTPHeaders {
				headers = append(headers, HttpHeader{
					Name:  h.Name,
					Value: h.Value,
				})
			}
			probe.Headers = headers
		}
	}

	if tcpSocket := input.TcpSocket; tcpSocket != nil {
		probe.Transport = "TCP"
		probe.Host = pointer.From(tcpSocket.Host)
		probe.Port = tcpSocket.Port
	}

	result = append(result, probe)

	return result
}

func flattenContainerAppJobVolumes(input *[]jobs.Volume) []ContainerVolume {
	if input == nil || len(*input) == 0 {
		return []ContainerVolume{}
	}

	result := make([]ContainerVolume, 0)
	for _, v := range *input {
		containerVolume := ContainerVolume{
			Name:        pointer.From(v.Name),
			StorageName: pointer.From(v.StorageName),
		}
		if v.StorageType != nil {
			containerVolume.StorageType = string(*v.StorageType)
		}
		if v.MountOptions != nil {
			containerVolume.MountOptions = pointer.From(v.MountOptions)
		}

		result = append(result, containerVolume)
	}

	return result
}

func FlattenContainerAppJobRegistries(input *[]jobs.RegistryCredentials) []Registry {
	if input == nil || len(*input) == 0 {
		return []Registry{}
	}

	result := make([]Registry, 0)

	for _, v := range *input {
		result = append(result, Registry{
			Identity:          pointer.From(v.Identity),
			PasswordSecretRef: pointer.From(v.PasswordSecretRef),
			Server:            pointer.From(v.Server),
			UserName:          pointer.From(v.Username),
		})
	}

	return result
}

func FlattenContainerAppJobConfigurationEventTriggerConfig(input *jobs.JobConfigurationEventTriggerConfig) []EventTriggerConfiguration {
	if input == nil {
		return []EventTriggerConfiguration{}
	}

	result := make([]EventTriggerConfiguration, 0)

	eventTriggerConfig := EventTriggerConfiguration{
		Parallelism:            pointer.From(input.Parallelism),
		ReplicaCompletionCount: pointer.From(input.ReplicaCompletionCount),
	}

	if input.Scale != nil {
		scale := flattenContainerAppJobScale(input.Scale)
		eventTriggerConfig.Scale = scale
	}

	result = append(result, eventTriggerConfig)

	return result
}

func FlattenContainerAppJobConfigurationManualTriggerConfig(input *jobs.JobConfigurationManualTriggerConfig) []ManualTriggerConfiguration {
	if input == nil {
		return []ManualTriggerConfiguration{}
	}

	result := make([]ManualTriggerConfiguration, 0)

	manualTriggerConfig := ManualTriggerConfiguration{
		Parallelism:            pointer.From(input.Parallelism),
		ReplicaCompletionCount: pointer.From(input.ReplicaCompletionCount),
	}

	result = append(result, manualTriggerConfig)

	return result
}

func FlattenContainerAppJobConfigurationScheduleTriggerConfig(input *jobs.JobConfigurationScheduleTriggerConfig) []ScheduleTriggerConfiguration {
	if input == nil {
		return []ScheduleTriggerConfiguration{}
	}

	result := make([]ScheduleTriggerConfiguration, 0)

	scheduleTriggerConfig := ScheduleTriggerConfiguration{
		Parallelism:            pointer.From(input.Parallelism),
		ReplicaCompletionCount: pointer.From(input.ReplicaCompletionCount),
	}

	if input.CronExpression != "" {
		scheduleTriggerConfig.CronExpression = input.CronExpression
	}

	result = append(result, scheduleTriggerConfig)

	return result
}

func flattenContainerAppJobScale(input *jobs.JobScale) []ScaleModel {
	if input == nil {
		return []ScaleModel{}
	}

	result := make([]ScaleModel, 0)

	scale := ScaleModel{
		MaxExecutions:   pointer.From(input.MaxExecutions),
		MinExecutions:   pointer.From(input.MinExecutions),
		PollingInterval: pointer.From(input.PollingInterval),
	}

	if input.Rules != nil {
		rules := flattenContainerAppJobScaleRules(input.Rules)
		scale.Rules = rules
	}

	result = append(result, scale)

	return result
}

func flattenContainerAppJobScaleRules(input *[]jobs.JobScaleRule) []ScaleRule {
	if input == nil || len(*input) == 0 {
		return []ScaleRule{}
	}

	result := make([]ScaleRule, 0)

	for _, v := range *input {
		rule := ScaleRule{
			Name: pointer.From(v.Name),
			Type: pointer.From(v.Type),
		}

		if v.Metadata != nil {
			metadata := pointer.From(v.Metadata)
			if reflect.TypeOf(metadata).Kind() == reflect.Map {
				rule.Metadata = metadata.(map[string]interface{})
			}
		}

		if v.Auth != nil {
			auth := flattenContainerAppJobScaleRulesAuth(v.Auth)
			rule.Auth = auth
		}

		result = append(result, rule)
	}

	return result
}

func flattenContainerAppJobScaleRulesAuth(input *[]jobs.ScaleRuleAuth) []ScaleRuleAuth {
	if input == nil || len(*input) == 0 {
		return []ScaleRuleAuth{}
	}

	result := make([]ScaleRuleAuth, 0)

	for _, v := range *input {
		auth := ScaleRuleAuth{
			SecretReference:  pointer.From(v.SecretRef),
			TriggerParameter: pointer.From(v.TriggerParameter),
		}

		result = append(result, auth)
	}

	return result
}

func FlattenContainerAppJobSecrets(input *jobs.JobSecretsCollection) []Secret {
	if input == nil || input.Value == nil {
		return []Secret{}
	}

	result := make([]Secret, 0)

	for _, v := range input.Value {
		secret := Secret{
			Identity:         pointer.From(v.Identity),
			KeyVaultSecretId: pointer.From(v.KeyVaultURL),
			Name:             pointer.From(v.Name),
		}
		if v.KeyVaultURL == nil {
			secret.Value = pointer.From(v.Value)
		}
		result = append(result, secret)
	}

	return result
}
