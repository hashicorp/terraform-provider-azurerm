// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_service

import (
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func flattenSpringCloudConfigServerGitProperty(input *appplatform.ConfigServerProperties, d *pluginsdk.ResourceData) []any {
	if input == nil || input.ConfigServer == nil || input.ConfigServer.GitProperty == nil {
		return []any{}
	}

	gitProperty := input.ConfigServer.GitProperty

	// prepare old state to find sensitive props not returned by API.
	oldGitSetting := make(map[string]any)
	if oldGitSettings := d.Get("config_server_git_setting").([]any); len(oldGitSettings) > 0 {
		oldGitSetting = oldGitSettings[0].(map[string]any)
	}

	uri := pointer.From(gitProperty.URI)

	label := pointer.From(gitProperty.Label)

	searchPaths := pluginsdk.FlattenSlice(gitProperty.SearchPaths)

	httpBasicAuth := make([]any, 0)
	if gitProperty.Username != nil && gitProperty.Password != nil {
		// username and password returned by API are *
		// to avoid state diff, we get the props from old state
		username := ""
		password := ""
		if v, ok := oldGitSetting["http_basic_auth"]; ok {
			oldHTTPBasicAuth := v.([]any)
			if len(oldHTTPBasicAuth) > 0 {
				oldItem := oldHTTPBasicAuth[0].(map[string]any)
				username = oldItem["username"].(string)
				password = oldItem["password"].(string)
			}
		}

		httpBasicAuth = []any{
			map[string]any{
				"username": username,
				"password": password,
			},
		}
	}

	sshAuth := []any{}
	if gitProperty.PrivateKey != nil {
		// private_key, host_key and host_key_algorithm returned by API are *
		// to avoid state diff, we get the props from old state
		privateKey := ""
		hostKey := ""
		hostKeyAlgorithm := ""
		if v, ok := oldGitSetting["ssh_auth"]; ok {
			sshAuth := v.([]any)
			if len(sshAuth) > 0 {
				oldItem := sshAuth[0].(map[string]any)
				privateKey = oldItem["private_key"].(string)
				hostKey = oldItem["host_key"].(string)
				hostKeyAlgorithm = oldItem["host_key_algorithm"].(string)
			}
		}

		strictHostKeyChecking := pointer.From(gitProperty.StrictHostKeyChecking)

		sshAuth = []any{
			map[string]any{
				"private_key":                      privateKey,
				"host_key":                         hostKey,
				"host_key_algorithm":               hostKeyAlgorithm,
				"strict_host_key_checking_enabled": strictHostKeyChecking,
			},
		}
	}

	return []any{
		map[string]any{
			"uri":             uri,
			"label":           label,
			"search_paths":    searchPaths,
			"http_basic_auth": httpBasicAuth,
			"ssh_auth":        sshAuth,
			"repository":      flattenSpringCloudGitPatternRepository(gitProperty.Repositories, d),
		},
	}
}

func flattenSpringCloudGitPatternRepository(input *[]appplatform.GitPatternRepository, d *pluginsdk.ResourceData) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	// prepare old state to find sensitive props not returned by API.
	oldGitPatternRepositories := map[string]any{}
	if oldGitSettings := d.Get("config_server_git_setting").([]any); len(oldGitSettings) > 0 {
		oldGitSetting := oldGitSettings[0].(map[string]any)
		for _, r := range oldGitSetting["repository"].([]any) {
			repo := r.(map[string]any)
			if name, ok := repo["name"]; ok {
				oldGitPatternRepositories[name.(string)] = r
			}
		}
	}

	for _, item := range *input {
		name := pointer.From(item.Name)

		uri := pointer.From(item.URI)

		label := pointer.From(item.Label)

		// prepare old state to find sensitive props not returned by API.
		oldGitPatternRepository := make(map[string]any)
		if gpr, ok := oldGitPatternRepositories[name]; ok {
			oldGitPatternRepository = gpr.(map[string]any)
		}

		pattern := pluginsdk.FlattenSlice(item.Pattern)
		searchPaths := pluginsdk.FlattenSlice(item.SearchPaths)

		httpBasicAuth := []any{}
		if item.Username != nil && item.Password != nil {
			// username and password returned by API are *
			// to avoid state diff, we get the props from old state
			username := ""
			password := ""
			if v, ok := oldGitPatternRepository["http_basic_auth"]; ok {
				oldHTTPBasicAuth := v.([]any)
				if len(oldHTTPBasicAuth) > 0 {
					oldItem := oldHTTPBasicAuth[0].(map[string]any)
					username = oldItem["username"].(string)
					password = oldItem["password"].(string)
				}
			}

			httpBasicAuth = []any{
				map[string]any{
					"username": username,
					"password": password,
				},
			}
		}

		sshAuth := []any{}
		if item.PrivateKey != nil {
			// private_key, host_key and host_key_algorithm returned by API are *
			// to avoid state diff, we get the props from old state
			privateKey := ""
			hostKey := ""
			hostKeyAlgorithm := ""
			if v, ok := oldGitPatternRepository["ssh_auth"]; ok {
				sshAuth := v.([]any)
				if len(sshAuth) > 0 {
					oldItem := sshAuth[0].(map[string]any)
					privateKey = oldItem["private_key"].(string)
					hostKey = oldItem["host_key"].(string)
					hostKeyAlgorithm = oldItem["host_key_algorithm"].(string)
				}
			}

			strictHostKeyChecking := pointer.From(item.StrictHostKeyChecking)

			sshAuth = []any{
				map[string]any{
					"private_key":                      privateKey,
					"host_key":                         hostKey,
					"host_key_algorithm":               hostKeyAlgorithm,
					"strict_host_key_checking_enabled": strictHostKeyChecking,
				},
			}
		}

		results = append(results, map[string]any{
			"name":            name,
			"uri":             uri,
			"label":           label,
			"pattern":         pattern,
			"search_paths":    searchPaths,
			"http_basic_auth": httpBasicAuth,
			"ssh_auth":        sshAuth,
		})
	}

	return results
}

func flattenSpringCloudTrace(input *appplatform.MonitoringSettingProperties) []any {
	if input == nil {
		return []any{}
	}

	enabled := false
	connectionString := ""
	samplingRate := 0.0
	if input.TraceEnabled != nil {
		enabled = *input.TraceEnabled
	}
	if input.AppInsightsInstrumentationKey != nil {
		connectionString = *input.AppInsightsInstrumentationKey
	}
	if input.AppInsightsSamplingRate != nil {
		samplingRate = *input.AppInsightsSamplingRate
	}

	if !enabled {
		return []any{}
	}

	return []any{
		map[string]any{
			"connection_string": connectionString,
			"sample_rate":       samplingRate,
		},
	}
}

func flattenSpringCloudNetwork(input *appplatform.NetworkProfile) []any {
	if input == nil {
		return []any{}
	}

	var serviceRuntimeSubnetID, appSubnetID, serviceRuntimeNetworkResourceGroup, appNetworkResourceGroup string
	var readTimeoutInSeconds int32
	var cidrRanges []any
	if input.ServiceRuntimeSubnetID != nil {
		serviceRuntimeSubnetID = *input.ServiceRuntimeSubnetID
	}
	if input.AppSubnetID != nil {
		appSubnetID = *input.AppSubnetID
	}
	if input.ServiceCidr != nil {
		cidrs := strings.Split(*input.ServiceCidr, ",")
		cidrRanges = pluginsdk.FlattenSlice(&cidrs)
	}
	if input.ServiceRuntimeNetworkResourceGroup != nil {
		serviceRuntimeNetworkResourceGroup = *input.ServiceRuntimeNetworkResourceGroup
	}
	if input.AppNetworkResourceGroup != nil {
		appNetworkResourceGroup = *input.AppNetworkResourceGroup
	}

	if ingressConfig := input.IngressConfig; ingressConfig != nil {
		if ingressConfig.ReadTimeoutInSeconds != nil {
			readTimeoutInSeconds = *ingressConfig.ReadTimeoutInSeconds
		}
	}

	outboundType := "loadBalancer"
	if input.OutboundType != nil {
		outboundType = *input.OutboundType
	}

	if serviceRuntimeSubnetID == "" && appSubnetID == "" && serviceRuntimeNetworkResourceGroup == "" && appNetworkResourceGroup == "" && len(cidrRanges) == 0 {
		return []any{}
	}

	return []any{
		map[string]any{
			"app_subnet_id":                          appSubnetID,
			"service_runtime_subnet_id":              serviceRuntimeSubnetID,
			"cidr_ranges":                            cidrRanges,
			"app_network_resource_group":             appNetworkResourceGroup,
			"read_timeout_seconds":                   readTimeoutInSeconds,
			"service_runtime_network_resource_group": serviceRuntimeNetworkResourceGroup,
			"outbound_type":                          outboundType,
		},
	}
}

func flattenOutboundPublicIPAddresses(input *appplatform.NetworkProfile) []any {
	if input == nil || input.OutboundIPs == nil {
		return []any{}
	}

	return pluginsdk.FlattenSlice(input.OutboundIPs.PublicIPs)
}

func flattenRequiredTraffic(input *appplatform.NetworkProfile) []any {
	if input == nil || input.RequiredTraffics == nil {
		return []any{}
	}

	result := make([]any, 0)
	for _, v := range *input.RequiredTraffics {
		protocol := pointer.From(v.Protocol)

		port := 0
		if v.Port != nil {
			port = int(*v.Port)
		}

		result = append(result, map[string]any{
			"protocol":     protocol,
			"port":         port,
			"ip_addresses": pluginsdk.FlattenSlice(v.Ips),
			"fqdns":        pluginsdk.FlattenSlice(v.Fqdns),
			"direction":    string(v.Direction),
		})
	}
	return result
}

func flattenSpringCloudContainerRegistries(state []appplatform.ContainerRegistryResource, input []appplatform.ContainerRegistryResource) []any {
	if input == nil {
		return []any{}
	}

	statePasswordMap := make(map[string]string)
	for _, v := range state {
		if v.Name == nil || v.Properties == nil || v.Properties.Credentials == nil {
			continue
		}
		if basicCredential, ok := v.Properties.Credentials.AsContainerRegistryBasicCredentials(); ok && basicCredential.Password != nil {
			statePasswordMap[*v.Name] = *basicCredential.Password
		}
	}
	result := make([]any, 0)
	for _, v := range input {
		name := ""
		username := ""
		password := ""
		server := ""
		if v.Name != nil {
			name = *v.Name
			password = statePasswordMap[name]
		}
		if v.Properties != nil && v.Properties.Credentials != nil {
			if basicCredential, ok := v.Properties.Credentials.AsContainerRegistryBasicCredentials(); ok {
				if basicCredential.Username != nil {
					username = *basicCredential.Username
				}
				if basicCredential.Server != nil {
					server = *basicCredential.Server
				}
			}
		}

		result = append(result, map[string]any{
			"name":     name,
			"username": username,
			"password": password,
			"server":   server,
		})
	}
	return result
}

func flattenSpringCloudBuildService(input *appplatform.BuildServiceProperties) []any {
	if input == nil || input.ContainerRegistry == nil {
		return []any{}
	}
	id, err := appplatform_rm.ParseContainerRegistryIDInsensitively(*input.ContainerRegistry)
	if err == nil {
		return []any{
			map[string]any{
				"container_registry_name": id.ContainerRegistryName,
			},
		}
	}
	return []any{}
}

func flattenSpringCloudMarketplaceResource(input *appplatform.MarketplaceResource) []any {
	if input == nil {
		return []any{}
	}
	plan := ""
	publisher := ""
	product := ""
	if input.Plan != nil {
		plan = *input.Plan
	}
	if input.Publisher != nil {
		publisher = *input.Publisher
	}
	if input.Product != nil {
		product = *input.Product
	}
	return []any{
		map[string]any{
			"plan":      plan,
			"publisher": publisher,
			"product":   product,
		},
	}
}
