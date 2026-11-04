// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandKubernetesClusterLinuxProfile(input []any) *managedclusters.ContainerServiceLinuxProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]any)

	adminUsername := config["admin_username"].(string)
	linuxKeys := config["ssh_key"].([]any)

	keyData := ""
	if key, ok := linuxKeys[0].(map[string]any); ok {
		keyData = key["key_data"].(string)
	}

	return &managedclusters.ContainerServiceLinuxProfile{
		AdminUsername: adminUsername,
		Ssh: managedclusters.ContainerServiceSshConfiguration{
			PublicKeys: []managedclusters.ContainerServiceSshPublicKey{
				{
					KeyData: keyData,
				},
			},
		},
	}
}

func expandKubernetesClusterIdentityProfile(input []any) *map[string]managedclusters.UserAssignedIdentity {
	identityProfile := make(map[string]managedclusters.UserAssignedIdentity)
	if len(input) == 0 || input[0] == nil {
		return &identityProfile
	}

	values := input[0].(map[string]any)

	if identity.Type(values["user_assigned_identity_id"].(string)) != "" {
		identityProfile["kubeletidentity"] = managedclusters.UserAssignedIdentity{
			ResourceId: pointer.To(values["user_assigned_identity_id"].(string)),
			ClientId:   pointer.To(values["client_id"].(string)),
			ObjectId:   pointer.To(values["object_id"].(string)),
		}
	}

	return &identityProfile
}

func expandKubernetesClusterWindowsProfile(input []any) *managedclusters.ManagedClusterWindowsProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]any)

	license := managedclusters.LicenseTypeNone
	if v := config["license"].(string); v != "" {
		license = managedclusters.LicenseType(v)
	}

	gmsaProfile := expandGmsaProfile(config["gmsa"].([]any))

	return &managedclusters.ManagedClusterWindowsProfile{
		AdminUsername: config["admin_username"].(string),
		AdminPassword: pointer.To(config["admin_password"].(string)),
		LicenseType:   &license,
		GmsaProfile:   gmsaProfile,
	}
}

func expandKubernetesClusterAPIAccessProfile(d *pluginsdk.ResourceData) *managedclusters.ManagedClusterAPIServerAccessProfile {
	enablePrivateCluster := false
	if v, ok := d.GetOk("private_cluster_enabled"); ok {
		enablePrivateCluster = v.(bool)
	}

	apiAccessProfile := &managedclusters.ManagedClusterAPIServerAccessProfile{
		EnablePrivateCluster:           &enablePrivateCluster,
		EnablePrivateClusterPublicFQDN: pointer.To(d.Get("private_cluster_public_fqdn_enabled").(bool)),
		DisableRunCommand:              pointer.To(!d.Get("run_command_enabled").(bool)),
	}

	apiServerAccessProfileRaw := d.Get("api_server_access_profile").([]any)
	if len(apiServerAccessProfileRaw) == 0 || apiServerAccessProfileRaw[0] == nil {
		return apiAccessProfile
	}

	config := apiServerAccessProfileRaw[0].(map[string]any)
	if v := config["authorized_ip_ranges"]; v != nil {
		if apiServerAuthorizedIPRanges := pluginsdk.ExpandStringSlice(v.(*pluginsdk.Set).List()); len(*apiServerAuthorizedIPRanges) > 0 {
			apiAccessProfile.AuthorizedIPRanges = apiServerAuthorizedIPRanges
		}
	}

	apiAccessProfile.EnableVnetIntegration = pointer.To(config["virtual_network_integration_enabled"].(bool))

	if v, ok := config["subnet_id"]; ok {
		if s := v.(string); s != "" {
			apiAccessProfile.SubnetId = pointer.To(s)
		}
	}

	return apiAccessProfile
}

func expandBootstrapProfile(rawBootstrapProfile []any) *managedclusters.ManagedClusterBootstrapProfile {
	if len(rawBootstrapProfile) == 0 || rawBootstrapProfile[0] == nil {
		return nil
	}

	config := rawBootstrapProfile[0].(map[string]any)
	var containerRegistryID *string
	if v, exists := config["container_registry_id"]; exists && v != "" {
		containerRegistryID = pointer.To(v.(string))
	}

	return &managedclusters.ManagedClusterBootstrapProfile{
		ArtifactSource:      pointer.ToEnum[managedclusters.ArtifactSource](config["artifact_source"].(string)),
		ContainerRegistryId: containerRegistryID,
	}
}

func expandKubernetesClusterWorkloadAutoscalerProfile(input []any, d *pluginsdk.ResourceData) *managedclusters.ManagedClusterWorkloadAutoScalerProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]any)

	var workloadAutoscalerProfile managedclusters.ManagedClusterWorkloadAutoScalerProfile

	if v := config["keda_enabled"].(bool); v || d.HasChange("workload_autoscaler_profile.0.keda_enabled") {
		workloadAutoscalerProfile.Keda = &managedclusters.ManagedClusterWorkloadAutoScalerProfileKeda{
			Enabled: v,
		}
	}

	if v := config["vertical_pod_autoscaler_enabled"].(bool); v || d.HasChange("workload_autoscaler_profile.0.vertical_pod_autoscaler_enabled") {
		workloadAutoscalerProfile.VerticalPodAutoscaler = &managedclusters.ManagedClusterWorkloadAutoScalerProfileVerticalPodAutoscaler{
			Enabled: v,
		}
	}

	return &workloadAutoscalerProfile
}

func expandKubernetesClusterNodeProvisioningProfile(input []any) *managedclusters.ManagedClusterNodeProvisioningProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	config := input[0].(map[string]any)
	profile := &managedclusters.ManagedClusterNodeProvisioningProfile{}
	if v := config["mode"].(string); v != "" {
		profile.Mode = pointer.ToEnum[managedclusters.NodeProvisioningMode](v)
	}
	if v := config["default_node_pools"].(string); v != "" {
		profile.DefaultNodePools = pointer.ToEnum[managedclusters.NodeProvisioningDefaultNodePools](v)
	}
	return profile
}

func expandGmsaProfile(input []any) *managedclusters.WindowsGmsaProfile {
	if len(input) == 0 {
		return nil
	}

	if input[0] == nil {
		return &managedclusters.WindowsGmsaProfile{
			Enabled:        pointer.To(true),
			DnsServer:      pointer.To(""),
			RootDomainName: pointer.To(""),
		}
	}

	config := input[0].(map[string]any)
	return &managedclusters.WindowsGmsaProfile{
		Enabled:        pointer.To(true),
		DnsServer:      pointer.To(config["dns_server"].(string)),
		RootDomainName: pointer.To(config["root_domain"].(string)),
	}
}

func expandKubernetesClusterNetworkProfile(input []any, d *pluginsdk.ResourceData) (*managedclusters.ContainerServiceNetworkProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	config := input[0].(map[string]any)

	networkPlugin := config["network_plugin"].(string)
	networkMode := config["network_mode"].(string)
	if networkPlugin != "azure" && networkMode != "" {
		return nil, fmt.Errorf("`network_mode` cannot be set if `network_plugin` is not `azure`")
	}
	networkPolicy := config["network_policy"].(string)
	loadBalancerProfileRaw := config["load_balancer_profile"].([]any)
	loadBalancerSku := config["load_balancer_sku"].(string)
	natGatewayProfileRaw := config["nat_gateway_profile"].([]any)
	outboundType := config["outbound_type"].(string)
	ipVersions, err := expandIPVersions(config["ip_versions"].([]any))
	if err != nil {
		return nil, err
	}

	networkProfile := managedclusters.ContainerServiceNetworkProfile{
		NetworkPlugin:   pointer.ToEnum[managedclusters.NetworkPlugin](networkPlugin),
		NetworkMode:     pointer.ToEnum[managedclusters.NetworkMode](networkMode),
		NetworkPolicy:   pointer.ToEnum[managedclusters.NetworkPolicy](networkPolicy),
		LoadBalancerSku: pointer.ToEnum[managedclusters.LoadBalancerSku](loadBalancerSku),
		OutboundType:    pointer.ToEnum[managedclusters.OutboundType](outboundType),
		IPFamilies:      ipVersions,
	}

	if networkDataPlane := config["network_data_plane"].(string); networkDataPlane != "" {
		networkProfile.NetworkDataplane = pointer.ToEnum[managedclusters.NetworkDataplane](networkDataPlane)
	}

	if networkPluginMode := config["network_plugin_mode"].(string); networkPluginMode != "" {
		networkProfile.NetworkPluginMode = pointer.ToEnum[managedclusters.NetworkPluginMode](networkPluginMode)
	}

	if len(loadBalancerProfileRaw) > 0 {
		if !strings.EqualFold(loadBalancerSku, "standard") {
			return nil, fmt.Errorf("only load balancer SKU 'Standard' supports load balancer profiles. Provided load balancer type: %s", loadBalancerSku)
		}

		networkProfile.LoadBalancerProfile = expandLoadBalancerProfile(loadBalancerProfileRaw)
	}

	if len(natGatewayProfileRaw) > 0 {
		if !strings.EqualFold(loadBalancerSku, "standard") {
			return nil, fmt.Errorf("only load balancer SKU 'Standard' supports NAT Gateway profiles. Provided load balancer type: %s", loadBalancerSku)
		}

		networkProfile.NatGatewayProfile = expandNatGatewayProfile(natGatewayProfileRaw)
	}

	if v, ok := config["dns_service_ip"]; ok && v.(string) != "" {
		dnsServiceIP := v.(string)
		networkProfile.DnsServiceIP = pointer.To(dnsServiceIP)
	}

	if v, ok := config["pod_cidr"]; ok && v.(string) != "" {
		podCidr := v.(string)
		networkProfile.PodCidr = pointer.To(podCidr)
	}

	if v, ok := config["pod_cidrs"]; ok {
		networkProfile.PodCidrs = pluginsdk.ExpandStringSlice(v.([]any))
	}

	if v, ok := config["service_cidr"]; ok && v.(string) != "" {
		serviceCidr := v.(string)
		networkProfile.ServiceCidr = pointer.To(serviceCidr)
	}

	if v, ok := config["service_cidrs"]; ok {
		networkProfile.ServiceCidrs = pluginsdk.ExpandStringSlice(v.([]any))
	}

	if v, ok := config["advanced_networking"]; ok {
		networkProfile.AdvancedNetworking = expandKubernetesClusterAdvancedNetworking(v.([]any), d)
	}

	return &networkProfile, nil
}

func expandKubernetesClusterAdvancedNetworking(input []any, d *pluginsdk.ResourceData) *managedclusters.AdvancedNetworking {
	if len(input) == 0 || input[0] == nil {
		o, n := d.GetChange("network_profile.0.advanced_networking")
		if o != nil && len(o.([]any)) == 1 && n == nil || len(n.([]any)) == 0 {
			return &managedclusters.AdvancedNetworking{
				Enabled: pointer.To(false),
				Observability: &managedclusters.AdvancedNetworkingObservability{
					Enabled: pointer.To(false),
				},
				Security: &managedclusters.AdvancedNetworkingSecurity{
					Enabled: pointer.To(false),
				},
			}
		}
		return nil
	}

	config := input[0].(map[string]any)
	observabilityEnabled := config["observability_enabled"].(bool)
	securityEnabled := config["security_enabled"].(bool)

	return &managedclusters.AdvancedNetworking{
		Enabled: pointer.To(true),
		Observability: &managedclusters.AdvancedNetworkingObservability{
			Enabled: pointer.To(observabilityEnabled),
		},
		Security: &managedclusters.AdvancedNetworkingSecurity{
			Enabled: pointer.To(securityEnabled),
		},
	}
}

func expandLoadBalancerProfile(d []any) *managedclusters.ManagedClusterLoadBalancerProfile {
	if d[0] == nil {
		return nil
	}

	config := d[0].(map[string]any)

	profile := &managedclusters.ManagedClusterLoadBalancerProfile{}

	if mins, ok := config["idle_timeout_in_minutes"]; ok && mins.(int) != 0 {
		profile.IdleTimeoutInMinutes = pointer.To(int64(mins.(int)))
	}

	if port, ok := config["outbound_ports_allocated"].(int); ok {
		profile.AllocatedOutboundPorts = pointer.To(int64(port))
	}

	if ipCount := config["managed_outbound_ip_count"]; ipCount != nil {
		if c := int64(ipCount.(int)); c > 0 {
			profile.ManagedOutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileManagedOutboundIPs{Count: &c}
		}
	}

	if ipv6Count := config["managed_outbound_ipv6_count"]; ipv6Count != nil {
		if c := int64(ipv6Count.(int)); c > 0 {
			if profile.ManagedOutboundIPs == nil {
				profile.ManagedOutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileManagedOutboundIPs{}
			}
			profile.ManagedOutboundIPs.CountIPv6 = &c
		}
	}

	if ipPrefixes := idsToResourceReferences(config["outbound_ip_prefix_ids"]); ipPrefixes != nil {
		profile.OutboundIPPrefixes = &managedclusters.ManagedClusterLoadBalancerProfileOutboundIPPrefixes{PublicIPPrefixes: ipPrefixes}
	}

	if outIps := idsToResourceReferences(config["outbound_ip_address_ids"]); outIps != nil {
		profile.OutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileOutboundIPs{PublicIPs: outIps}
	}

	if backendPoolType, ok := config["backend_pool_type"].(string); ok {
		profile.BackendPoolType = pointer.ToEnum[managedclusters.BackendPoolType](backendPoolType)
	}

	return profile
}

func expandIPVersions(input []any) (*[]managedclusters.IPFamily, error) {
	if len(input) == 0 {
		return nil, nil
	}

	ipv := make([]managedclusters.IPFamily, 0)
	for _, data := range input {
		ipv = append(ipv, managedclusters.IPFamily(data.(string)))
	}

	if len(ipv) == 1 && ipv[0] == managedclusters.IPFamilyIPvSix {
		return nil, fmt.Errorf("`ip_versions` must be `IPv4` or `IPv4` and `IPv6`. `IPv6` alone is not supported")
	}

	return &ipv, nil
}

func expandNatGatewayProfile(d []any) *managedclusters.ManagedClusterNATGatewayProfile {
	if d[0] == nil {
		return nil
	}

	config := d[0].(map[string]any)

	profile := &managedclusters.ManagedClusterNATGatewayProfile{}

	if mins, ok := config["idle_timeout_in_minutes"]; ok && mins.(int) != 0 {
		profile.IdleTimeoutInMinutes = pointer.To(int64(mins.(int)))
	}

	if ipCount := config["managed_outbound_ip_count"]; ipCount != nil {
		if c := int64(ipCount.(int)); c > 0 {
			profile.ManagedOutboundIPProfile = &managedclusters.ManagedClusterManagedOutboundIPProfile{Count: &c}
		}
	}

	return profile
}

func expandKubernetesClusterAzureActiveDirectoryRoleBasedAccessControl(input []any) *managedclusters.ManagedClusterAADProfile {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	azureAdRaw := input[0].(map[string]any)

	adminGroupObjectIdsRaw := azureAdRaw["admin_group_object_ids"].([]any)
	adminGroupObjectIds := pluginsdk.ExpandStringSlice(adminGroupObjectIdsRaw)
	return &managedclusters.ManagedClusterAADProfile{
		TenantID:            pointer.To(azureAdRaw["tenant_id"].(string)),
		Managed:             pointer.To(true),
		AdminGroupObjectIDs: adminGroupObjectIds,
		EnableAzureRBAC:     pointer.To(azureAdRaw["azure_rbac_enabled"].(bool)),
	}
}

func expandKubernetesClusterManagedClusterIdentity(input []any) (*identity.SystemOrUserAssignedMap, error) {
	expanded, err := identity.ExpandSystemOrUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := identity.SystemOrUserAssignedMap{
		Type: identity.Type(string(expanded.Type)),
	}
	if expanded.Type == identity.TypeUserAssigned {
		out.IdentityIds = make(map[string]identity.UserAssignedIdentityDetails)
		for k := range expanded.IdentityIds {
			out.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				// intentionally empty
			}
		}
	}
	return &out, nil
}

func expandKubernetesClusterAutoScalerProfile(input []any) *managedclusters.ManagedClusterPropertiesAutoScalerProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]any)

	balanceSimilarNodeGroups := config["balance_similar_node_groups"].(bool)
	daemonsetEvictionForEmptyNodes := config["daemonset_eviction_for_empty_nodes_enabled"].(bool)
	daemonsetEvictionForOccupiedNodes := config["daemonset_eviction_for_occupied_nodes_enabled"].(bool)
	expander := config["expander"].(string)
	ignoreDaemonsetsUtilization := config["ignore_daemonsets_utilization_enabled"].(bool)
	maxGracefulTerminationSec := config["max_graceful_termination_sec"].(string)
	maxNodeProvisionTime := config["max_node_provisioning_time"].(string)
	maxUnreadyNodes := fmt.Sprint(config["max_unready_nodes"].(int))
	maxUnreadyPercentage := fmt.Sprint(config["max_unready_percentage"].(float64))
	newPodScaleUpDelay := config["new_pod_scale_up_delay"].(string)
	scaleDownDelayAfterAdd := config["scale_down_delay_after_add"].(string)
	scaleDownDelayAfterDelete := config["scale_down_delay_after_delete"].(string)
	scaleDownDelayAfterFailure := config["scale_down_delay_after_failure"].(string)
	scaleDownUnneededTime := config["scale_down_unneeded"].(string)
	scaleDownUnreadyTime := config["scale_down_unready"].(string)
	scaleDownUtilizationThreshold := config["scale_down_utilization_threshold"].(string)
	emptyBulkDeleteMax := config["empty_bulk_delete_max"].(string)
	scanInterval := config["scan_interval"].(string)
	skipNodesWithLocalStorage := config["skip_nodes_with_local_storage"].(bool)
	skipNodesWithSystemPods := config["skip_nodes_with_system_pods"].(bool)

	return &managedclusters.ManagedClusterPropertiesAutoScalerProfile{
		BalanceSimilarNodeGroups:          pointer.To(strconv.FormatBool(balanceSimilarNodeGroups)),
		DaemonsetEvictionForEmptyNodes:    pointer.To(daemonsetEvictionForEmptyNodes),
		DaemonsetEvictionForOccupiedNodes: pointer.To(daemonsetEvictionForOccupiedNodes),
		Expander:                          pointer.ToEnum[managedclusters.Expander](expander),
		IgnoreDaemonsetsUtilization:       pointer.To(ignoreDaemonsetsUtilization),
		MaxGracefulTerminationSec:         pointer.To(maxGracefulTerminationSec),
		MaxNodeProvisionTime:              pointer.To(maxNodeProvisionTime),
		MaxTotalUnreadyPercentage:         pointer.To(maxUnreadyPercentage),
		NewPodScaleUpDelay:                pointer.To(newPodScaleUpDelay),
		OkTotalUnreadyCount:               pointer.To(maxUnreadyNodes),
		ScaleDownDelayAfterAdd:            pointer.To(scaleDownDelayAfterAdd),
		ScaleDownDelayAfterDelete:         pointer.To(scaleDownDelayAfterDelete),
		ScaleDownDelayAfterFailure:        pointer.To(scaleDownDelayAfterFailure),
		ScaleDownUnneededTime:             pointer.To(scaleDownUnneededTime),
		ScaleDownUnreadyTime:              pointer.To(scaleDownUnreadyTime),
		ScaleDownUtilizationThreshold:     pointer.To(scaleDownUtilizationThreshold),
		MaxEmptyBulkDelete:                pointer.To(emptyBulkDeleteMax),
		ScanInterval:                      pointer.To(scanInterval),
		SkipNodesWithLocalStorage:         pointer.To(strconv.FormatBool(skipNodesWithLocalStorage)),
		SkipNodesWithSystemPods:           pointer.To(strconv.FormatBool(skipNodesWithSystemPods)),
	}
}

func expandKubernetesClusterAzureKeyVaultKms(ctx context.Context, keyVaultsClient *client.Client, subscriptionId string, d *pluginsdk.ResourceData, input []any) (*managedclusters.AzureKeyVaultKms, error) {
	if len(input) == 0 && d.HasChanges("key_management_service") {
		return &managedclusters.AzureKeyVaultKms{
			Enabled: pointer.To(false),
		}, nil
	} else if len(input) == 0 {
		return nil, nil
	}

	raw := input[0].(map[string]any)
	kvAccess := managedclusters.KeyVaultNetworkAccessTypes(raw["key_vault_network_access"].(string))

	azureKeyVaultKms := &managedclusters.AzureKeyVaultKms{
		Enabled:               pointer.To(true),
		KeyId:                 pointer.To(raw["key_vault_key_id"].(string)),
		KeyVaultNetworkAccess: &kvAccess,
	}

	// Set Key vault Resource ID in case public access is disabled
	if kvAccess == managedclusters.KeyVaultNetworkAccessTypesPrivate {
		subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)

		keyVaultKeyId, err := keyvault.ParseNestedItemID(*azureKeyVaultKms.KeyId, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
		if err != nil {
			return nil, err
		}
		keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, keyVaultKeyId.KeyVaultBaseURL)
		if err != nil {
			return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", keyVaultKeyId.KeyVaultBaseURL, err)
		}

		azureKeyVaultKms.KeyVaultResourceId = keyVaultID
	}

	return azureKeyVaultKms, nil
}

func expandKubernetesClusterMaintenanceConfigurationDefault(input []any) *maintenanceconfigurations.MaintenanceConfigurationProperties {
	if len(input) == 0 {
		return nil
	}
	value := input[0].(map[string]any)
	return &maintenanceconfigurations.MaintenanceConfigurationProperties{
		NotAllowedTime: expandKubernetesClusterMaintenanceConfigurationTimeSpans(value["not_allowed"].(*pluginsdk.Set).List()),
		TimeInWeek:     expandKubernetesClusterMaintenanceConfigurationTimeInWeeks(value["allowed"].(*pluginsdk.Set).List()),
	}
}

func expandKubernetesClusterMaintenanceConfigurationTimeSpans(input []any) *[]maintenanceconfigurations.TimeSpan {
	results := make([]maintenanceconfigurations.TimeSpan, 0)
	for _, item := range input {
		v := item.(map[string]any)
		start, _ := time.Parse(time.RFC3339, v["start"].(string))
		end, _ := time.Parse(time.RFC3339, v["end"].(string))
		results = append(results, maintenanceconfigurations.TimeSpan{
			Start: pointer.To(start.Format("2006-01-02T15:04:05Z07:00")),
			End:   pointer.To(end.Format("2006-01-02T15:04:05Z07:00")),
		})
	}
	return &results
}

func expandKubernetesClusterMaintenanceConfigurationDateSpans(input []any) *[]maintenanceconfigurations.DateSpan {
	results := make([]maintenanceconfigurations.DateSpan, 0)
	for _, item := range input {
		v := item.(map[string]any)
		start, _ := time.Parse(time.RFC3339, v["start"].(string))
		end, _ := time.Parse(time.RFC3339, v["end"].(string))
		results = append(results, maintenanceconfigurations.DateSpan{
			Start: start.Format("2006-01-02"),
			End:   end.Format("2006-01-02"),
		})
	}
	return &results
}

func expandKubernetesClusterMaintenanceConfigurationTimeInWeeks(input []any) *[]maintenanceconfigurations.TimeInWeek {
	results := make([]maintenanceconfigurations.TimeInWeek, 0)
	for _, item := range input {
		v := item.(map[string]any)
		results = append(results, maintenanceconfigurations.TimeInWeek{
			Day:       pointer.ToEnum[maintenanceconfigurations.WeekDay](v["day"].(string)),
			HourSlots: pluginsdk.ExpandInt64Slice(v["hours"].(*pluginsdk.Set).List()),
		})
	}
	return &results
}

func expandKubernetesClusterHttpProxyConfig(input []any) *managedclusters.ManagedClusterHTTPProxyConfig {
	httpProxyConfig := managedclusters.ManagedClusterHTTPProxyConfig{}
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]any)

	httpProxyConfig.HTTPProxy = pointer.To(config["http_proxy"].(string))
	httpProxyConfig.HTTPSProxy = pointer.To(config["https_proxy"].(string))
	if value := config["trusted_ca"].(string); len(value) != 0 {
		httpProxyConfig.TrustedCa = pointer.To(value)
	}

	noProxyRaw := config["no_proxy"].(*pluginsdk.Set).List()
	httpProxyConfig.NoProxy = pluginsdk.ExpandStringSlice(noProxyRaw)

	return &httpProxyConfig
}

func expandKubernetesClusterOidcIssuerProfile(input bool) *managedclusters.ManagedClusterOIDCIssuerProfile {
	oidcIssuerProfile := managedclusters.ManagedClusterOIDCIssuerProfile{}
	oidcIssuerProfile.Enabled = &input

	return &oidcIssuerProfile
}

func expandKubernetesClusterMicrosoftDefender(d *pluginsdk.ResourceData, input []any) *managedclusters.ManagedClusterSecurityProfileDefender {
	if (len(input) == 0 || input[0] == nil) && d.HasChange("microsoft_defender") {
		return &managedclusters.ManagedClusterSecurityProfileDefender{
			SecurityMonitoring: &managedclusters.ManagedClusterSecurityProfileDefenderSecurityMonitoring{
				Enabled: pointer.To(false),
			},
		}
	} else if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]any)
	return &managedclusters.ManagedClusterSecurityProfileDefender{
		SecurityMonitoring: &managedclusters.ManagedClusterSecurityProfileDefenderSecurityMonitoring{
			Enabled: pointer.To(true),
		},
		LogAnalyticsWorkspaceResourceId: pointer.To(config["log_analytics_workspace_id"].(string)),
	}
}

func expandStorageProfile(input []any) *managedclusters.ManagedClusterStorageProfile {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]any)

	profile := managedclusters.ManagedClusterStorageProfile{
		BlobCSIDriver: &managedclusters.ManagedClusterStorageProfileBlobCSIDriver{
			Enabled: pointer.To(raw["blob_driver_enabled"].(bool)),
		},
		DiskCSIDriver: &managedclusters.ManagedClusterStorageProfileDiskCSIDriver{
			Enabled: pointer.To(raw["disk_driver_enabled"].(bool)),
		},
		FileCSIDriver: &managedclusters.ManagedClusterStorageProfileFileCSIDriver{
			Enabled: pointer.To(raw["file_driver_enabled"].(bool)),
		},
		SnapshotController: &managedclusters.ManagedClusterStorageProfileSnapshotController{
			Enabled: pointer.To(raw["snapshot_controller_enabled"].(bool)),
		},
	}

	return &profile
}

func expandEdgeZone(input string) *edgezones.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &edgezones.Model{
		Name: normalized,
	}
}

func expandKubernetesClusterServiceMeshProfile(input []any, existing *managedclusters.ServiceMeshProfile) *managedclusters.ServiceMeshProfile {
	if len(input) == 0 {
		// explicitly disable istio if it was enabled before
		if existing != nil && existing.Mode == managedclusters.ServiceMeshModeIstio {
			return &managedclusters.ServiceMeshProfile{
				Mode: managedclusters.ServiceMeshModeDisabled,
			}
		}

		return nil
	}

	raw := input[0].(map[string]any)

	mode := raw["mode"].(string)

	profile := managedclusters.ServiceMeshProfile{}

	if managedclusters.ServiceMeshMode(mode) == managedclusters.ServiceMeshModeIstio {
		profile.Mode = managedclusters.ServiceMeshMode(mode)
		profile.Istio = &managedclusters.IstioServiceMesh{}

		profile.Istio.Components = &managedclusters.IstioComponents{}

		istioIngressGatewaysList := make([]managedclusters.IstioIngressGateway, 0)

		if raw["internal_ingress_gateway_enabled"] != nil {
			ingressGatewayElementInternal := managedclusters.IstioIngressGateway{
				Enabled: raw["internal_ingress_gateway_enabled"].(bool),
				Mode:    managedclusters.IstioIngressGatewayModeInternal,
			}

			istioIngressGatewaysList = append(istioIngressGatewaysList, ingressGatewayElementInternal)
		}

		if raw["external_ingress_gateway_enabled"] != nil {
			ingressGatewayElementExternal := managedclusters.IstioIngressGateway{
				Enabled: raw["external_ingress_gateway_enabled"].(bool),
				Mode:    managedclusters.IstioIngressGatewayModeExternal,
			}

			istioIngressGatewaysList = append(istioIngressGatewaysList, ingressGatewayElementExternal)
		}

		profile.Istio.Components.IngressGateways = &istioIngressGatewaysList

		if raw["certificate_authority"] != nil {
			profile.Istio.CertificateAuthority = expandKubernetesClusterServiceMeshProfileCertificateAuthority(raw["certificate_authority"].([]any))
		}

		if raw["revisions"] != nil {
			profile.Istio.Revisions = pluginsdk.ExpandStringSlice(raw["revisions"].([]any))
		}
	}

	return &profile
}

func expandKubernetesClusterServiceMeshProfileCertificateAuthority(input []any) *managedclusters.IstioCertificateAuthority {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]any)

	return &managedclusters.IstioCertificateAuthority{
		Plugin: &managedclusters.IstioPluginCertificateAuthority{
			KeyVaultId:          pointer.To(config["key_vault_id"].(string)),
			RootCertObjectName:  pointer.To(config["root_cert_object_name"].(string)),
			CertChainObjectName: pointer.To(config["cert_chain_object_name"].(string)),
			CertObjectName:      pointer.To(config["cert_object_name"].(string)),
			KeyObjectName:       pointer.To(config["key_object_name"].(string)),
		},
	}
}

func expandKubernetesClusterIngressProfile(d *pluginsdk.ResourceData, input []any) *managedclusters.ManagedClusterIngressProfile {
	if len(input) == 0 && d.HasChange("web_app_routing") {
		return &managedclusters.ManagedClusterIngressProfile{
			WebAppRouting: &managedclusters.ManagedClusterIngressProfileWebAppRouting{
				Enabled: pointer.To(false),
			},
		}
	} else if len(input) == 0 {
		return nil
	}

	out := managedclusters.ManagedClusterIngressProfile{
		WebAppRouting: &managedclusters.ManagedClusterIngressProfileWebAppRouting{
			Enabled: pointer.To(true),
		},
	}
	if input[0] != nil {
		config := input[0].(map[string]any)
		if v := config["dns_zone_ids"]; v != nil {
			if dnsZoneResourceIds, ok := v.([]any); ok && len(dnsZoneResourceIds) > 0 {
				out.WebAppRouting.DnsZoneResourceIds = pluginsdk.ExpandStringSlice(dnsZoneResourceIds)
			}
		}
		if v := config["default_nginx_controller"]; v != nil {
			out.WebAppRouting.Nginx = &managedclusters.ManagedClusterIngressProfileNginx{
				DefaultIngressControllerType: pointer.ToEnum[managedclusters.NginxIngressControllerType](v.(string)),
			}
		}
	}
	return &out
}

func expandKubernetesClusterAzureMonitorProfile(input []any) *managedclusters.ManagedClusterAzureMonitorProfile {
	if len(input) == 0 {
		return &managedclusters.ManagedClusterAzureMonitorProfile{
			Metrics: &managedclusters.ManagedClusterAzureMonitorProfileMetrics{},
		}
	}
	if input[0] == nil {
		return &managedclusters.ManagedClusterAzureMonitorProfile{
			Metrics: &managedclusters.ManagedClusterAzureMonitorProfileMetrics{
				Enabled: true,
			},
		}
	}
	config := input[0].(map[string]any)
	return &managedclusters.ManagedClusterAzureMonitorProfile{
		Metrics: &managedclusters.ManagedClusterAzureMonitorProfileMetrics{
			Enabled: true,
			KubeStateMetrics: &managedclusters.ManagedClusterAzureMonitorProfileKubeStateMetrics{
				MetricAnnotationsAllowList: pointer.To(config["annotations_allowed"].(string)),
				MetricLabelsAllowlist:      pointer.To(config["labels_allowed"].(string)),
			},
		},
	}
}

func expandKubernetesClusterMetricsProfile(input bool, skuTier string) (*managedclusters.ManagedClusterMetricsProfile, error) {
	if input && skuTier != "Standard" && skuTier != "Premium" {
		return nil, fmt.Errorf("`sku_tier` must be either `Standard` or `Premium` when cost analysis is enabled")
	}
	return &managedclusters.ManagedClusterMetricsProfile{
		CostAnalysis: &managedclusters.ManagedClusterCostAnalysis{
			Enabled: pointer.To(input),
		},
	}, nil
}

func expandKubernetesClusterUpgradeOverrideSetting(input []any) *managedclusters.ClusterUpgradeSettings {
	if len(input) == 0 || input[0] == nil {
		// Return nil only when upgrade_override block is not set
		return nil
	}

	raw := input[0].(map[string]any)
	return &managedclusters.ClusterUpgradeSettings{
		OverrideSettings: &managedclusters.UpgradeOverrideSettings{
			ForceUpgrade: pointer.To(raw["force_upgrade_enabled"].(bool)),
			Until:        pointer.To(raw["effective_until"].(string)),
		},
	}
}
