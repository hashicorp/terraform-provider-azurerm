// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_group

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2025-09-01/containerinstance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandContainerGroupInitContainers(d *pluginsdk.ResourceData, addedEmptyDirs map[string]bool) (*[]containerinstance.InitContainerDefinition, []containerinstance.Volume, error) {
	containersConfig := d.Get("init_container").([]any)
	containers := make([]containerinstance.InitContainerDefinition, 0)
	containerGroupVolumes := make([]containerinstance.Volume, 0)
	for _, containerConfig := range containersConfig {
		data := containerConfig.(map[string]any)

		name := data["name"].(string)
		image := data["image"].(string)

		container := containerinstance.InitContainerDefinition{
			Name: name,
			Properties: containerinstance.InitContainerPropertiesDefinition{
				Image:           pointer.To(image),
				SecurityContext: expandContainerSecurityContext(data["security"].([]any)),
			},
		}

		// Set both sensitive and non-secure environment variables
		var envVars *[]containerinstance.EnvironmentVariable
		var secEnvVars *[]containerinstance.EnvironmentVariable

		// Expand environment_variables into slice
		if v, ok := data["environment_variables"]; ok {
			envVars = expandContainerEnvironmentVariables(v, false)
		}

		// Expand secure_environment_variables into slice
		if v, ok := data["secure_environment_variables"]; ok {
			secEnvVars = expandContainerEnvironmentVariables(v, true)
		}

		// Combine environment variable slices
		*envVars = append(*envVars, *secEnvVars...)

		// Set both secure and non secure environment variables
		container.Properties.EnvironmentVariables = envVars

		if v, ok := data["commands"]; ok {
			c := v.([]any)
			command := make([]string, 0)
			for _, v := range c {
				command = append(command, v.(string))
			}

			container.Properties.Command = &command
		}

		if v, ok := data["volume"]; ok {
			volumeMounts, _, err := expandSingleContainerVolume(v)
			if err != nil {
				return nil, nil, err
			}
			container.Properties.VolumeMounts = volumeMounts

			expandedContainerGroupVolumes, err := expandContainerVolume(v, addedEmptyDirs, containerGroupVolumes)
			if err != nil {
				return nil, nil, err
			}
			containerGroupVolumes = expandedContainerGroupVolumes
		}

		containers = append(containers, container)
	}

	return &containers, containerGroupVolumes, nil
}

func expandContainerSecurityContext(input []any) *containerinstance.SecurityContextDefinition {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)

	return &containerinstance.SecurityContextDefinition{
		Privileged: pointer.To(raw["privilege_enabled"].(bool)),
	}
}

func expandContainerGroupContainers(d *pluginsdk.ResourceData, addedEmptyDirs map[string]bool) ([]containerinstance.Container, []containerinstance.Port, []containerinstance.Volume, error) {
	containersConfig := d.Get("container").([]any)
	containers := make([]containerinstance.Container, 0)
	containerInstancePorts := make([]containerinstance.Port, 0)
	containerGroupPorts := make([]containerinstance.Port, 0)
	containerGroupVolumes := make([]containerinstance.Volume, 0)

	for _, containerConfig := range containersConfig {
		data := containerConfig.(map[string]any)

		name := data["name"].(string)
		cpu := data["cpu"].(float64)
		memory := data["memory"].(float64)

		container := containerinstance.Container{
			Name: name,
			Properties: containerinstance.ContainerProperties{
				Image: pointer.To(data["image"].(string)),
				Resources: &containerinstance.ResourceRequirements{
					Requests: containerinstance.ResourceRequests{
						MemoryInGB: memory,
						Cpu:        cpu,
					},
				},
				SecurityContext: expandContainerSecurityContext(data["security"].([]any)),
			},
		}

		cpuLimit := data["cpu_limit"].(float64)
		memLimit := data["memory_limit"].(float64)

		if !(cpuLimit == 0.0 && memLimit == 0.0) {
			limits := &containerinstance.ResourceLimits{}
			if cpuLimit != 0.0 {
				limits.Cpu = &cpuLimit
			}
			if memLimit != 0.0 {
				limits.MemoryInGB = &memLimit
			}

			container.Properties.Resources.Limits = limits
		}

		if v, ok := data["ports"].(*pluginsdk.Set); ok && len(v.List()) > 0 {
			var ports []containerinstance.ContainerPort
			for _, v := range v.List() {
				portObj := v.(map[string]any)

				port := int64(portObj["port"].(int))
				proto := portObj["protocol"].(string)

				ports = append(ports, containerinstance.ContainerPort{
					Port:     port,
					Protocol: pointer.ToEnum[containerinstance.ContainerNetworkProtocol](proto),
				})
				containerInstancePorts = append(containerInstancePorts, containerinstance.Port{
					Port:     port,
					Protocol: pointer.ToEnum[containerinstance.ContainerGroupNetworkProtocol](proto),
				})
			}
			container.Properties.Ports = &ports
		}

		// Set both sensitive and non-secure environment variables
		var envVars *[]containerinstance.EnvironmentVariable
		var secEnvVars *[]containerinstance.EnvironmentVariable

		// Expand environment_variables into slice
		if v, ok := data["environment_variables"]; ok {
			envVars = expandContainerEnvironmentVariables(v, false)
		}

		// Expand secure_environment_variables into slice
		if v, ok := data["secure_environment_variables"]; ok {
			secEnvVars = expandContainerEnvironmentVariables(v, true)
		}

		// Combine environment variable slices
		*envVars = append(*envVars, *secEnvVars...)

		// Set both secure and non secure environment variables
		container.Properties.EnvironmentVariables = envVars

		if v, ok := data["commands"]; ok {
			c := v.([]any)
			command := make([]string, 0)
			for _, v := range c {
				command = append(command, v.(string))
			}

			container.Properties.Command = &command
		}

		if v, ok := data["volume"]; ok {
			volumeMounts, _, err := expandSingleContainerVolume(v)
			if err != nil {
				return nil, nil, nil, err
			}
			container.Properties.VolumeMounts = volumeMounts

			expandedContainerGroupVolumes, err := expandContainerVolume(v, addedEmptyDirs, containerGroupVolumes)
			if err != nil {
				return nil, nil, nil, err
			}
			containerGroupVolumes = expandedContainerGroupVolumes
		}

		if v, ok := data["liveness_probe"]; ok {
			container.Properties.LivenessProbe = expandContainerProbe(v)
		}

		if v, ok := data["readiness_probe"]; ok {
			container.Properties.ReadinessProbe = expandContainerProbe(v)
		}

		containers = append(containers, container)
	}

	// Determine ports to be exposed on the group level, based on exposed_ports
	// and on what ports have been exposed on individual containers.
	if v, ok := d.Get("exposed_port").(*pluginsdk.Set); ok && len(v.List()) > 0 {
		cgpMap := make(map[int64]map[containerinstance.ContainerGroupNetworkProtocol]bool)
		for _, p := range containerInstancePorts {
			if p.Protocol == nil {
				continue
			}
			protocol := *p.Protocol

			if val, ok := cgpMap[p.Port]; ok {
				val[protocol] = true
				cgpMap[p.Port] = val
			} else {
				cgpMap[p.Port] = map[containerinstance.ContainerGroupNetworkProtocol]bool{protocol: true}
			}
		}

		for _, p := range v.List() {
			portConfig := p.(map[string]any)
			port := int64(portConfig["port"].(int))
			proto := portConfig["protocol"].(string)
			if !cgpMap[port][containerinstance.ContainerGroupNetworkProtocol(proto)] {
				return nil, nil, nil, fmt.Errorf(`port %d/%s is not exposed on any individual container in the container group.
					An exposed_ports block contains %d/%s, but no individual container has a ports block with the same port
					and protocol. Any ports exposed on the container group must also be exposed on an individual container`,
					port, proto, port, proto)
			}
			containerGroupPorts = append(containerGroupPorts, containerinstance.Port{
				Port:     port,
				Protocol: pointer.ToEnum[containerinstance.ContainerGroupNetworkProtocol](proto),
			})
		}
	} else {
		containerGroupPorts = containerInstancePorts // remove in 3.0 of the provider
	}

	return containers, containerGroupPorts, containerGroupVolumes, nil
}

func expandContainerVolume(v any, addedEmptyDirs map[string]bool, containerGroupVolumes []containerinstance.Volume) ([]containerinstance.Volume, error) {
	_, containerVolumes, err := expandSingleContainerVolume(v)
	if err != nil {
		return nil, err
	}
	if containerVolumes != nil {
		for _, cgVol := range *containerVolumes {
			if cgVol.EmptyDir != nil {
				if addedEmptyDirs[cgVol.Name] {
					// empty_dir-volumes are allowed to overlap across containers, in fact that is their primary purpose,
					// but the containerGroup must not declare same name of such volumes twice.
					continue
				}
				addedEmptyDirs[cgVol.Name] = true
			}
			containerGroupVolumes = append(containerGroupVolumes, cgVol)
		}
	}
	return containerGroupVolumes, nil
}

func expandContainerEnvironmentVariables(input any, secure bool) *[]containerinstance.EnvironmentVariable {
	envVars := input.(map[string]any)
	output := make([]containerinstance.EnvironmentVariable, 0, len(envVars))

	if secure {
		for k, v := range envVars {
			ev := containerinstance.EnvironmentVariable{
				Name:        k,
				SecureValue: pointer.To(v.(string)),
			}

			output = append(output, ev)
		}
	} else {
		for k, v := range envVars {
			ev := containerinstance.EnvironmentVariable{
				Name:  k,
				Value: pointer.To(v.(string)),
			}

			output = append(output, ev)
		}
	}
	return &output
}

func expandContainerImageRegistryCredentials(d *pluginsdk.ResourceData) *[]containerinstance.ImageRegistryCredential {
	credsRaw := d.Get("image_registry_credential").([]any)
	if len(credsRaw) == 0 {
		return nil
	}

	output := make([]containerinstance.ImageRegistryCredential, 0, len(credsRaw))

	for _, c := range credsRaw {
		credConfig := c.(map[string]any)

		imageRegistryCredential := containerinstance.ImageRegistryCredential{}
		if v := credConfig["server"]; v != nil && v != "" {
			imageRegistryCredential.Server = v.(string)
		}
		if v := credConfig["username"]; v != nil && v != "" {
			imageRegistryCredential.Username = pointer.To(v.(string))
		}
		if v := credConfig["password"]; v != nil && v != "" {
			imageRegistryCredential.Password = pointer.To(v.(string))
		}
		if v := credConfig["user_assigned_identity_id"]; v != nil && v != "" {
			imageRegistryCredential.Identity = pointer.To(v.(string))
		}

		output = append(output, imageRegistryCredential)
	}

	return &output
}

func expandSingleContainerVolume(input any) (*[]containerinstance.VolumeMount, *[]containerinstance.Volume, error) {
	volumesRaw := input.([]any)

	if len(volumesRaw) == 0 {
		return nil, nil, nil
	}

	volumeMounts := make([]containerinstance.VolumeMount, 0)
	containerGroupVolumes := make([]containerinstance.Volume, 0)

	for _, volumeRaw := range volumesRaw {
		volumeConfig := volumeRaw.(map[string]any)

		name := volumeConfig["name"].(string)
		mountPath := volumeConfig["mount_path"].(string)
		readOnly := volumeConfig["read_only"].(bool)
		emptyDir := volumeConfig["empty_dir"].(bool)
		shareName := volumeConfig["share_name"].(string)
		storageAccountName := volumeConfig["storage_account_name"].(string)
		storageAccountKey := volumeConfig["storage_account_key"].(string)

		vm := containerinstance.VolumeMount{
			Name:      name,
			MountPath: mountPath,
			ReadOnly:  pointer.To(readOnly),
		}

		volumeMounts = append(volumeMounts, vm)

		cv := containerinstance.Volume{
			Name: name,
		}

		secret := expandSecrets(volumeConfig["secret"].(map[string]any))

		gitRepoVolume := expandGitRepoVolume(volumeConfig["git_repo"].([]any))

		switch {
		case emptyDir:
			if shareName != "" || storageAccountName != "" || storageAccountKey != "" || secret != nil || gitRepoVolume != nil {
				return nil, nil, fmt.Errorf("only one of `empty_dir` volume, `git_repo` volume, `secret` volume or storage account volume (`share_name`, `storage_account_name`, and `storage_account_key`) can be specified")
			}
			var m any = map[string]string{}
			cv.EmptyDir = &m
		case gitRepoVolume != nil:
			if shareName != "" || storageAccountName != "" || storageAccountKey != "" || secret != nil {
				return nil, nil, fmt.Errorf("only one of `empty_dir` volume, `git_repo` volume, `secret` volume or storage account volume (`share_name`, `storage_account_name`, and `storage_account_key`) can be specified")
			}
			cv.GitRepo = gitRepoVolume
		case secret != nil:
			if shareName != "" || storageAccountName != "" || storageAccountKey != "" {
				return nil, nil, fmt.Errorf("only one of `empty_dir` volume, `git_repo` volume, `secret` volume or storage account volume (`share_name`, `storage_account_name`, and `storage_account_key`) can be specified")
			}
			cv.Secret = &secret
		default:
			if shareName == "" && storageAccountName == "" && storageAccountKey == "" {
				return nil, nil, fmt.Errorf("exactly one of `empty_dir` volume, `git_repo` volume, `secret` volume or storage account volume (`share_name`, `storage_account_name`, and `storage_account_key`) must be specified")
			} else if shareName == "" || storageAccountName == "" || storageAccountKey == "" {
				return nil, nil, fmt.Errorf("when using a storage account volume, all of `share_name`, `storage_account_name`, `storage_account_key` must be specified")
			}
			cv.AzureFile = &containerinstance.AzureFileVolume{
				ShareName:          shareName,
				ReadOnly:           pointer.To(readOnly),
				StorageAccountName: storageAccountName,
				StorageAccountKey:  pointer.To(storageAccountKey),
			}
		}

		containerGroupVolumes = append(containerGroupVolumes, cv)
	}

	return &volumeMounts, &containerGroupVolumes, nil
}

func expandGitRepoVolume(input []any) *containerinstance.GitRepoVolume {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)
	gitRepoVolume := &containerinstance.GitRepoVolume{
		Repository: v["url"].(string),
	}
	if directory := v["directory"].(string); directory != "" {
		gitRepoVolume.Directory = pointer.To(directory)
	}
	if revision := v["revision"].(string); revision != "" {
		gitRepoVolume.Revision = pointer.To(revision)
	}
	return gitRepoVolume
}

func expandSecrets(secretsMap map[string]any) map[string]string {
	if len(secretsMap) == 0 {
		return nil
	}
	output := make(map[string]string, len(secretsMap))

	for name, value := range secretsMap {
		output[name] = value.(string)
	}

	return output
}

func expandContainerProbe(input any) *containerinstance.ContainerProbe {
	probe := containerinstance.ContainerProbe{}
	probeRaw := input.([]any)

	if len(probeRaw) == 0 {
		return nil
	}

	for _, p := range probeRaw {
		if p == nil {
			continue
		}
		probeConfig := p.(map[string]any)

		if v := probeConfig["initial_delay_seconds"].(int); v > 0 {
			probe.InitialDelaySeconds = pointer.To(int64(v))
		}

		if v := probeConfig["period_seconds"].(int); v > 0 {
			probe.PeriodSeconds = pointer.To(int64(v))
		}

		if v := probeConfig["failure_threshold"].(int); v > 0 {
			probe.FailureThreshold = pointer.To(int64(v))
		}

		if v := probeConfig["success_threshold"].(int); v > 0 {
			probe.SuccessThreshold = pointer.To(int64(v))
		}

		if v := probeConfig["timeout_seconds"].(int); v > 0 {
			probe.TimeoutSeconds = pointer.To(int64(v))
		}

		commands := probeConfig["exec"].([]any)
		if len(commands) > 0 {
			exec := containerinstance.ContainerExec{
				Command: pluginsdk.ExpandStringSlice(commands),
			}
			probe.Exec = &exec
		}

		httpRaw := probeConfig["http_get"].([]any)
		if len(httpRaw) > 0 {
			for _, httpget := range httpRaw {
				if httpget == nil {
					continue
				}
				x := httpget.(map[string]any)

				path := x["path"].(string)
				port := x["port"].(int)
				scheme := x["scheme"].(string)

				probe.HTTPGet = &containerinstance.ContainerHTTPGet{
					Path:        pointer.To(path),
					Port:        int64(port),
					Scheme:      pointer.ToEnum[containerinstance.Scheme](scheme),
					HTTPHeaders: expandContainerProbeHttpHeaders(x["http_headers"].(map[string]any)),
				}
			}
		}
	}
	return &probe
}

func expandContainerProbeHttpHeaders(input map[string]any) *[]containerinstance.HTTPHeader {
	if len(input) == 0 {
		return nil
	}

	headers := []containerinstance.HTTPHeader{}
	for k, v := range input {
		header := containerinstance.HTTPHeader{
			Name:  pointer.To(k),
			Value: pointer.To(v.(string)),
		}
		headers = append(headers, header)
	}
	return &headers
}

func expandContainerGroupDiagnostics(input []any) *containerinstance.ContainerGroupDiagnostics {
	if len(input) == 0 {
		return nil
	}

	vs := input[0].(map[string]any)

	analyticsVs := vs["log_analytics"].([]any)
	analyticsV := analyticsVs[0].(map[string]any)

	workspaceId := analyticsV["workspace_id"].(string)
	workspaceKey := analyticsV["workspace_key"].(string)

	logAnalytics := containerinstance.LogAnalytics{
		WorkspaceId:  workspaceId,
		WorkspaceKey: workspaceKey,
	}

	if logType := analyticsV["log_type"].(string); logType != "" {
		logAnalytics.LogType = pointer.ToEnum[containerinstance.LogAnalyticsLogType](logType)

		metadataMap := analyticsV["metadata"].(map[string]any)
		metadata := make(map[string]string)
		for k, v := range metadataMap {
			metadata[k] = v.(string)
		}

		logAnalytics.Metadata = &metadata
	}

	return &containerinstance.ContainerGroupDiagnostics{LogAnalytics: &logAnalytics}
}

func expandContainerGroupDnsConfig(input any) *containerinstance.DnsConfiguration {
	dnsConfigRaw := input.([]any)
	if len(dnsConfigRaw) > 0 && dnsConfigRaw[0] != nil {
		config := dnsConfigRaw[0].(map[string]any)

		nameservers := []string{}
		for _, v := range config["nameservers"].([]any) {
			nameservers = append(nameservers, v.(string))
		}
		options := []string{}
		for _, v := range config["options"].(*pluginsdk.Set).List() {
			options = append(options, v.(string))
		}
		searchDomains := []string{}
		for _, v := range config["search_domains"].(*pluginsdk.Set).List() {
			searchDomains = append(searchDomains, v.(string))
		}

		return &containerinstance.DnsConfiguration{
			Options:       pointer.To(strings.Join(options, " ")),
			SearchDomains: pointer.To(strings.Join(searchDomains, " ")),
			NameServers:   nameservers,
		}
	}

	return nil
}

func expandContainerGroupSubnets(input []any) (*[]containerinstance.ContainerGroupSubnetId, error) {
	if len(input) == 0 {
		return nil, nil
	}

	results := make([]containerinstance.ContainerGroupSubnetId, 0)
	for _, item := range input {
		id, err := commonids.ParseSubnetID(item.(string))
		if err != nil {
			return nil, fmt.Errorf(`parsing subnet id %q: %v`, item, err)
		}

		results = append(results, containerinstance.ContainerGroupSubnetId{
			Id: id.ID(),
		})
	}
	return &results, nil
}
