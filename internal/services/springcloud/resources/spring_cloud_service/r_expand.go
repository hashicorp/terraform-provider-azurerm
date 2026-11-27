// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_service

import (
	"errors"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func expandSpringCloudNetwork(input []any) *appplatform.NetworkProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)
	cidrRanges := pluginsdk.ExpandStringSlice(v["cidr_ranges"].([]any))
	network := &appplatform.NetworkProfile{
		ServiceRuntimeSubnetID: pointer.To(v["service_runtime_subnet_id"].(string)),
		AppSubnetID:            pointer.To(v["app_subnet_id"].(string)),
		ServiceCidr:            pointer.To(strings.Join(*cidrRanges, ",")),
		OutboundType:           pointer.To(v["outbound_type"].(string)),
	}
	if readTimeoutInSeconds := v["read_timeout_seconds"].(int); readTimeoutInSeconds != 0 {
		network.IngressConfig = &appplatform.IngressConfig{
			ReadTimeoutInSeconds: pointer.To(int32(readTimeoutInSeconds)),
		}
	}
	if serviceRuntimeNetworkResourceGroup := v["service_runtime_network_resource_group"].(string); serviceRuntimeNetworkResourceGroup != "" {
		network.ServiceRuntimeNetworkResourceGroup = pointer.To(serviceRuntimeNetworkResourceGroup)
	}
	if appNetworkResourceGroup := v["app_network_resource_group"].(string); appNetworkResourceGroup != "" {
		network.AppNetworkResourceGroup = pointer.To(appNetworkResourceGroup)
	}
	return network
}

func expandSpringCloudConfigServerGitProperty(input []any) (*appplatform.ConfigServerGitProperty, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	v := input[0].(map[string]any)
	if v == nil {
		return nil, nil
	}

	result := appplatform.ConfigServerGitProperty{
		URI: pointer.To(v["uri"].(string)),
	}

	if label := v["label"].(string); label != "" {
		result.Label = pointer.To(label)
	}
	if searchPaths := v["search_paths"].([]any); len(searchPaths) > 0 {
		result.SearchPaths = pluginsdk.ExpandStringSlice(searchPaths)
	}

	httpBasicAuth := v["http_basic_auth"].([]any)
	sshAuth := v["ssh_auth"].([]any)
	if len(httpBasicAuth) > 0 && len(sshAuth) > 0 {
		return nil, errors.New("can not set both `http_basic_auth` and `ssh_auth`")
	}
	if len(httpBasicAuth) > 0 {
		v := httpBasicAuth[0].(map[string]any)
		result.Username = pointer.To(v["username"].(string))
		result.Password = pointer.To(v["password"].(string))
	}
	if len(sshAuth) > 0 {
		v := sshAuth[0].(map[string]any)
		result.PrivateKey = pointer.To(v["private_key"].(string))
		result.StrictHostKeyChecking = pointer.To(v["strict_host_key_checking_enabled"].(bool))

		if hostKey := v["host_key"].(string); hostKey != "" {
			result.HostKey = pointer.To(hostKey)
		}
		if hostKeyAlgorithm := v["host_key_algorithm"].(string); hostKeyAlgorithm != "" {
			result.HostKeyAlgorithm = pointer.To(hostKeyAlgorithm)
		}
	}

	if v, ok := v["repository"]; ok {
		repositories, err := expandSpringCloudGitPatternRepository(v.([]any))
		if err != nil {
			return nil, err
		}
		result.Repositories = repositories
	}

	return &result, nil
}

func expandSpringCloudGitPatternRepository(input []any) (*[]appplatform.GitPatternRepository, error) {
	results := make([]appplatform.GitPatternRepository, 0)
	for _, item := range input {
		v := item.(map[string]any)

		result := appplatform.GitPatternRepository{
			Name: pointer.To(v["name"].(string)),
			URI:  pointer.To(v["uri"].(string)),
		}

		if label := v["label"].(string); len(label) > 0 {
			result.Label = pointer.To(label)
		}
		if pattern := v["pattern"].([]any); len(pattern) > 0 {
			result.Pattern = pluginsdk.ExpandStringSlice(pattern)
		}
		if searchPaths := v["search_paths"].([]any); len(searchPaths) > 0 {
			result.SearchPaths = pluginsdk.ExpandStringSlice(searchPaths)
		}

		httpBasicAuth := v["http_basic_auth"].([]any)
		sshAuth := v["ssh_auth"].([]any)
		if len(httpBasicAuth) > 0 && len(sshAuth) > 0 {
			return nil, errors.New("can not set both `http_basic_auth` and `ssh_auth` for the same repository")
		}
		if len(httpBasicAuth) > 0 {
			v := httpBasicAuth[0].(map[string]any)
			result.Username = pointer.To(v["username"].(string))
			result.Password = pointer.To(v["password"].(string))
		}
		if len(sshAuth) > 0 {
			v := sshAuth[0].(map[string]any)
			result.PrivateKey = pointer.To(v["private_key"].(string))
			result.StrictHostKeyChecking = pointer.To(v["strict_host_key_checking_enabled"].(bool))

			if hostKey := v["host_key"].(string); hostKey != "" {
				result.HostKey = pointer.To(hostKey)
			}
			if hostKeyAlgorithm := v["host_key_algorithm"].(string); hostKeyAlgorithm != "" {
				result.HostKeyAlgorithm = pointer.To(hostKeyAlgorithm)
			}
		}

		results = append(results, result)
	}
	return &results, nil
}

func expandSpringCloudTrace(input []any) *appplatform.MonitoringSettingProperties {
	if len(input) == 0 || input[0] == nil {
		return &appplatform.MonitoringSettingProperties{
			TraceEnabled: pointer.To(false),
		}
	}

	v := input[0].(map[string]any)
	return &appplatform.MonitoringSettingProperties{
		TraceEnabled:                  pointer.To(true),
		AppInsightsInstrumentationKey: pointer.To(v["connection_string"].(string)),
		AppInsightsSamplingRate:       pointer.To(v["sample_rate"].(float64)),
	}
}

func expandSpringCloudContainerRegistries(input []any) []appplatform.ContainerRegistryResource {
	if len(input) == 0 {
		return nil
	}
	out := make([]appplatform.ContainerRegistryResource, 0)
	for _, item := range input {
		v := item.(map[string]any)
		out = append(out, appplatform.ContainerRegistryResource{
			Name: pointer.To(v["name"].(string)),
			Properties: &appplatform.ContainerRegistryProperties{
				Credentials: &appplatform.ContainerRegistryBasicCredentials{
					Username: pointer.To(v["username"].(string)),
					Password: pointer.To(v["password"].(string)),
					Server:   pointer.To(v["server"].(string)),
				},
			},
		})
	}
	return out
}

func expandSpringCloudBuildService(input []any, springId commonids.SpringCloudServiceId) appplatform.BuildServiceProperties {
	if len(input) == 0 || input[0] == nil {
		return appplatform.BuildServiceProperties{}
	}
	v := input[0].(map[string]any)
	out := appplatform.BuildServiceProperties{}
	if value := v["container_registry_name"].(string); value != "" {
		out.ContainerRegistry = pointer.To(appplatform_rm.NewContainerRegistryID(springId.SubscriptionId, springId.ResourceGroupName, springId.ServiceName, value).ID())
	}
	return out
}

func expandSpringCloudMarketplaceResource(input []any) *appplatform.MarketplaceResource {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)
	return &appplatform.MarketplaceResource{
		Plan:      pointer.To(v["plan"].(string)),
		Publisher: pointer.To(v["publisher"].(string)),
		Product:   pointer.To(v["product"].(string)),
	}
}
