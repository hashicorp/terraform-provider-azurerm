// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_group

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2025-09-01/containerinstance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenPorts(ports []any) *pluginsdk.Set {
	if len(ports) > 0 {
		flatPorts := make([]any, 0)
		for _, p := range ports {
			port := make(map[string]any)
			switch t := p.(type) {
			case containerinstance.Port:
				port["port"] = int(t.Port)
				proto := ""
				if t.Protocol != nil {
					proto = string(*t.Protocol)
				}
				port["protocol"] = proto
			case containerinstance.ContainerPort:
				port["port"] = int(t.Port)
				proto := ""
				if t.Protocol != nil {
					proto = string(*t.Protocol)
				}
				port["protocol"] = proto
			}
			flatPorts = append(flatPorts, port)
		}
		return pluginsdk.NewSet(resourceContainerGroupPortsHash, flatPorts)
	}
	return pluginsdk.NewSet(resourceContainerGroupPortsHash, make([]any, 0))
}

func flattenContainerSecurityContext(input *containerinstance.SecurityContextDefinition) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"privilege_enabled": pointer.From(input.Privileged),
		},
	}
}

func flattenContainerProbeHttpHeaders(input *[]containerinstance.HTTPHeader) map[string]any {
	if input == nil {
		return map[string]any{}
	}

	output := map[string]any{}
	for _, header := range *input {
		name := pointer.From(header.Name)
		output[name] = pointer.From(header.Value)
	}
	return output
}

func flattenContainerImageRegistryCredentials(d *pluginsdk.ResourceData, input *[]containerinstance.ImageRegistryCredential) []any {
	if input == nil {
		return []any{}
	}
	configsOld := d.Get("image_registry_credential").([]any)

	output := make([]any, 0)
	for i, cred := range *input {
		credConfig := make(map[string]any)
		credConfig["server"] = cred.Server
		credConfig["username"] = cred.Username
		credConfig["user_assigned_identity_id"] = cred.Identity

		if len(configsOld) > i {
			data := configsOld[i].(map[string]any)
			oldServer := data["server"].(string)
			if cred.Server == oldServer {
				if v, ok := d.GetOk(fmt.Sprintf("image_registry_credential.%d.password", i)); ok {
					credConfig["password"] = v.(string)
				}
			}
		}

		output = append(output, credConfig)
	}
	return output
}

func flattenContainerGroupInitContainers(d *pluginsdk.ResourceData, initContainers *[]containerinstance.InitContainerDefinition, containerGroupVolumes *[]containerinstance.Volume) []any {
	if initContainers == nil {
		return []any{}
	}
	// map old container names to index so we can look up things up
	nameIndexMap := map[string]int{}
	for i, c := range d.Get("init_container").([]any) {
		cfg := c.(map[string]any)
		nameIndexMap[cfg["name"].(string)] = i
	}

	containerCfg := make([]any, 0, len(*initContainers))
	for _, container := range *initContainers {
		name := container.Name

		// get index from name
		index := nameIndexMap[name]

		containerConfig := make(map[string]any)
		containerConfig["name"] = name

		if v := container.Properties.Image; v != nil {
			containerConfig["image"] = *v
		}

		if container.Properties.EnvironmentVariables != nil {
			if len(*container.Properties.EnvironmentVariables) > 0 {
				containerConfig["environment_variables"] = flattenContainerEnvironmentVariables(container.Properties.EnvironmentVariables)
				containerConfig["secure_environment_variables"] = flattenContainerSecureEnvironmentVariables(container.Properties.EnvironmentVariables, d, index, "init_container")
			}
		}

		commands := make([]string, 0)
		if command := container.Properties.Command; command != nil {
			commands = *command
		}
		containerConfig["commands"] = commands

		if containerGroupVolumes != nil && container.Properties.VolumeMounts != nil {
			containersConfigRaw := d.Get("container").([]any)
			flattenContainerVolume(containerConfig, containersConfigRaw, container.Name, container.Properties.VolumeMounts, containerGroupVolumes)
		}

		containerConfig["security"] = flattenContainerSecurityContext(container.Properties.SecurityContext)

		containerCfg = append(containerCfg, containerConfig)
	}

	return containerCfg
}

func flattenContainerGroupContainers(d *pluginsdk.ResourceData, containers *[]containerinstance.Container, containerGroupVolumes *[]containerinstance.Volume) []any {
	// map old container names to index so we can look up things up
	nameIndexMap := map[string]int{}
	for i, c := range d.Get("container").([]any) {
		cfg := c.(map[string]any)
		nameIndexMap[cfg["name"].(string)] = i
	}

	containerCfg := make([]any, 0, len(*containers))
	for _, container := range *containers {
		name := container.Name

		// get index from name
		index := nameIndexMap[name]

		containerConfig := make(map[string]any)
		containerConfig["name"] = name

		containerConfig["image"] = container.Properties.Image

		resources := container.Properties.Resources
		resourceRequests := resources.Requests
		containerConfig["cpu"] = resourceRequests.Cpu
		containerConfig["memory"] = resourceRequests.MemoryInGB

		if resourceLimits := resources.Limits; resourceLimits != nil {
			if v := resourceLimits.Cpu; v != nil {
				containerConfig["cpu_limit"] = *v
			}
			if v := resourceLimits.MemoryInGB; v != nil {
				containerConfig["memory_limit"] = *v
			}
		}

		containerPorts := make([]any, len(*container.Properties.Ports))
		if container.Properties.Ports != nil {
			for i := range *container.Properties.Ports {
				containerPorts[i] = (*container.Properties.Ports)[i]
			}
		}
		containerConfig["ports"] = flattenPorts(containerPorts)

		if container.Properties.EnvironmentVariables != nil {
			if len(*container.Properties.EnvironmentVariables) > 0 {
				containerConfig["environment_variables"] = flattenContainerEnvironmentVariables(container.Properties.EnvironmentVariables)
				containerConfig["secure_environment_variables"] = flattenContainerSecureEnvironmentVariables(container.Properties.EnvironmentVariables, d, index, "container")
			}
		}

		commands := make([]string, 0)
		if command := container.Properties.Command; command != nil {
			commands = *command
		}
		containerConfig["commands"] = commands

		if containerGroupVolumes != nil && container.Properties.VolumeMounts != nil {
			containersConfigRaw := d.Get("container").([]any)
			flattenContainerVolume(containerConfig, containersConfigRaw, container.Name, container.Properties.VolumeMounts, containerGroupVolumes)
		}

		containerConfig["liveness_probe"] = flattenContainerProbes(container.Properties.LivenessProbe)
		containerConfig["readiness_probe"] = flattenContainerProbes(container.Properties.ReadinessProbe)
		containerConfig["security"] = flattenContainerSecurityContext(container.Properties.SecurityContext)

		containerCfg = append(containerCfg, containerConfig)
	}

	return containerCfg
}

func flattenContainerVolume(containerConfig map[string]any, containersConfigRaw []any, containerName string, volumeMounts *[]containerinstance.VolumeMount, containerGroupVolumes *[]containerinstance.Volume) {
	// Also pass in the container volume config from schema
	var containerVolumesConfig *[]any
	for _, containerConfigRaw := range containersConfigRaw {
		data := containerConfigRaw.(map[string]any)
		nameRaw := data["name"].(string)
		if nameRaw == containerName {
			// found container config for current container
			// extract volume mounts from config
			if v, ok := data["volume"]; ok {
				containerVolumesConfig = pointer.To(v.([]any))
			}
		}
	}
	volumeConfigs := make([]any, 0)

	if volumeMounts == nil {
		containerConfig["volume"] = nil
		return
	}

	for _, vm := range *volumeMounts {
		volumeConfig := make(map[string]any)
		volumeConfig["name"] = vm.Name
		volumeConfig["mount_path"] = vm.MountPath
		if vm.ReadOnly != nil {
			volumeConfig["read_only"] = *vm.ReadOnly
		}

		// find corresponding volume in container group volumes
		// and use the data
		if containerGroupVolumes != nil {
			for _, cgv := range *containerGroupVolumes {
				if cgv.Name == vm.Name {
					if file := cgv.AzureFile; file != nil {
						volumeConfig["share_name"] = file.ShareName
						volumeConfig["storage_account_name"] = file.StorageAccountName
						// skip storage_account_key, is always nil
					}

					if cgv.EmptyDir != nil {
						volumeConfig["empty_dir"] = true
					}

					volumeConfig["git_repo"] = flattenGitRepoVolume(cgv.GitRepo)
				}
			}
		}

		// find corresponding volume in config
		// and use the data
		if containerVolumesConfig != nil {
			for _, cvr := range *containerVolumesConfig {
				cv := cvr.(map[string]any)
				rawName := cv["name"].(string)
				if vm.Name == rawName {
					volumeConfig["storage_account_key"] = cv["storage_account_key"].(string)
					volumeConfig["secret"] = cv["secret"]
				}
			}
		}

		volumeConfigs = append(volumeConfigs, volumeConfig)
	}

	containerConfig["volume"] = volumeConfigs
}

func flattenContainerSecureEnvironmentVariables(input *[]containerinstance.EnvironmentVariable, d *pluginsdk.ResourceData, oldContainerIndex int, rootPropName string) map[string]any {
	output := make(map[string]any)

	if input == nil {
		return output
	}

	for _, envVar := range *input {
		if envVar.Value == nil {
			output[envVar.Name] = d.Get(fmt.Sprintf("%s.%d.secure_environment_variables.%s", rootPropName, oldContainerIndex, envVar.Name))
		}
	}

	return output
}

func flattenContainerEnvironmentVariables(input *[]containerinstance.EnvironmentVariable) map[string]any {
	output := make(map[string]any)

	if input == nil {
		return output
	}

	for _, envVar := range *input {
		if envVar.Value != nil {
			log.Printf("[DEBUG] NOT SECURE: Name: %s - Value: %s", envVar.Name, *envVar.Value)
			output[envVar.Name] = *envVar.Value
		}
	}

	return output
}

func flattenGitRepoVolume(input *containerinstance.GitRepoVolume) []any {
	if input == nil {
		return []any{}
	}
	var revision, directory, repository string
	if input.Directory != nil {
		directory = *input.Directory
	}
	if input.Revision != nil {
		revision = *input.Revision
	}
	repository = input.Repository
	return []any{
		map[string]any{
			"url":       repository,
			"directory": directory,
			"revision":  revision,
		},
	}
}

func flattenContainerProbes(input *containerinstance.ContainerProbe) []any {
	outputs := make([]any, 0)
	if input == nil {
		return outputs
	}

	output := make(map[string]any)

	if v := input.Exec; v != nil {
		output["exec"] = *v.Command
	}

	httpGets := make([]any, 0)
	if get := input.HTTPGet; get != nil {
		httpGet := make(map[string]any)
		if v := get.Path; v != nil {
			httpGet["path"] = *v
		}
		httpGet["port"] = get.Port
		httpGet["scheme"] = get.Scheme
		httpGet["http_headers"] = flattenContainerProbeHttpHeaders(get.HTTPHeaders)
		httpGets = append(httpGets, httpGet)
	}
	output["http_get"] = httpGets

	if v := input.FailureThreshold; v != nil {
		output["failure_threshold"] = *v
	}

	if v := input.InitialDelaySeconds; v != nil {
		output["initial_delay_seconds"] = *v
	}

	if v := input.PeriodSeconds; v != nil {
		output["period_seconds"] = *v
	}

	if v := input.SuccessThreshold; v != nil {
		output["success_threshold"] = *v
	}

	if v := input.TimeoutSeconds; v != nil {
		output["timeout_seconds"] = *v
	}

	outputs = append(outputs, output)
	return outputs
}

func flattenContainerGroupDiagnostics(d *pluginsdk.ResourceData, input *containerinstance.ContainerGroupDiagnostics) []any {
	if input == nil {
		return []any{}
	}

	logAnalytics := make([]any, 0)

	if la := input.LogAnalytics; la != nil {
		output := make(map[string]any)

		logType := ""
		if la.LogType != nil {
			logType = string(*la.LogType)
		}
		output["log_type"] = logType

		metadata := make(map[string]any)
		if la.Metadata != nil {
			for k, v := range *la.Metadata {
				metadata[k] = v
			}
		}
		output["metadata"] = metadata
		output["workspace_id"] = la.WorkspaceId

		// the existing config may not exist at Import time, protect against it.
		workspaceKey := ""
		if existingDiags := d.Get("diagnostics").([]any); len(existingDiags) > 0 {
			existingDiag := existingDiags[0].(map[string]any)
			if existingLA := existingDiag["log_analytics"].([]any); len(existingLA) > 0 {
				vs := existingLA[0].(map[string]any)
				if key := vs["workspace_key"]; key != nil && key.(string) != "" {
					workspaceKey = key.(string)
				}
			}
		}
		output["workspace_key"] = workspaceKey

		logAnalytics = append(logAnalytics, output)
	}

	return []any{
		map[string]any{
			"log_analytics": logAnalytics,
		},
	}
}

func flattenContainerGroupDnsConfig(input *containerinstance.DnsConfiguration) []any {
	output := make(map[string]any)

	if input == nil {
		return make([]any, 0)
	}

	// We're converting to TypeSet here from an API response that looks like "a b c" (assumes space delimited)
	var searchDomains []string
	if input.SearchDomains != nil {
		searchDomains = strings.Fields(*input.SearchDomains)
	}
	output["search_domains"] = searchDomains

	// We're converting to TypeSet here from an API response that looks like "a b c" (assumes space delimited)
	var options []string
	if input.Options != nil {
		options = strings.Fields(*input.Options)
	}
	output["options"] = options
	output["nameservers"] = input.NameServers

	return []any{output}
}

func flattenContainerGroupSubnets(input *[]containerinstance.ContainerGroupSubnetId) ([]any, error) {
	subnetIDs := make([]any, 0)
	if input == nil {
		return subnetIDs, nil
	}

	for _, resourceRef := range *input {
		if resourceRef.Id == "" {
			continue
		}

		id, err := commonids.ParseSubnetIDInsensitively(resourceRef.Id)
		if err != nil {
			return nil, fmt.Errorf(`parsing subnet id %q: %v`, resourceRef.Id, err)
		}

		subnetIDs = append(subnetIDs, id.ID())
	}

	return subnetIDs, nil
}
