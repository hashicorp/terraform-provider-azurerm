// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster

import (
	"maps"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenKubernetesClusterIdentityProfile(profile map[string]managedclusters.UserAssignedIdentity) ([]any, error) {
	if profile == nil {
		return []any{}, nil
	}

	kubeletIdentity := make([]any, 0)
	if kubeletidentity, ok := profile["kubeletidentity"]; ok {
		clientId := pointer.From(kubeletidentity.ClientId)

		objectId := pointer.From(kubeletidentity.ObjectId)

		userAssignedIdentityId := ""
		if resourceid := kubeletidentity.ResourceId; resourceid != nil {
			parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*resourceid)
			if err != nil {
				return nil, err
			}

			userAssignedIdentityId = parsedId.ID()
		}

		kubeletIdentity = append(kubeletIdentity, map[string]any{
			"client_id":                 clientId,
			"object_id":                 objectId,
			"user_assigned_identity_id": userAssignedIdentityId,
		})
	}

	return kubeletIdentity, nil
}

func flattenKubernetesClusterLinuxProfile(profile *managedclusters.ContainerServiceLinuxProfile) []any {
	if profile == nil {
		return []any{}
	}

	adminUsername := profile.AdminUsername

	sshKeys := make([]any, 0)
	ssh := profile.Ssh
	if keys := ssh.PublicKeys; keys != nil {
		for _, sshKey := range keys {
			keyData := ""
			if kd := sshKey.KeyData; kd != "" {
				keyData = kd
			}
			sshKeys = append(sshKeys, map[string]any{
				"key_data": keyData,
			})
		}
	}

	return []any{
		map[string]any{
			"admin_username": adminUsername,
			"ssh_key":        sshKeys,
		},
	}
}

func flattenKubernetesClusterAPIAccessProfile(profile *managedclusters.ManagedClusterAPIServerAccessProfile) []any {
	if profile == nil {
		return []any{}
	}

	// API access profile can be managed by other properties, only return it if one of the properties has been set
	if profile.AuthorizedIPRanges == nil && profile.EnableVnetIntegration == nil && profile.SubnetId == nil {
		return []any{}
	}

	return []any{map[string]any{
		"authorized_ip_ranges":                pluginsdk.FlattenSlice(profile.AuthorizedIPRanges),
		"virtual_network_integration_enabled": pointer.From(profile.EnableVnetIntegration),
		"subnet_id":                           pointer.From(profile.SubnetId),
	}}
}

func flattenBootstrapProfile(profile *managedclusters.ManagedClusterBootstrapProfile) ([]any, error) {
	if profile == nil || profile.ArtifactSource == nil {
		return []any{}, nil
	}

	var containerRegistryID string
	if profile.ContainerRegistryId != nil {
		id, err := registries.ParseRegistryID(*profile.ContainerRegistryId)
		if err != nil {
			return nil, err
		}
		containerRegistryID = id.ID()
	}
	return []any{
		map[string]any{
			"artifact_source":       profile.ArtifactSource,
			"container_registry_id": containerRegistryID,
		},
	}, nil
}

func flattenKubernetesClusterWindowsProfile(profile *managedclusters.ManagedClusterWindowsProfile, d *pluginsdk.ResourceData) []any {
	if profile == nil {
		return []any{}
	}

	adminUsername := profile.AdminUsername

	// admin password isn't returned, so let's look it up
	adminPassword := ""
	if v, ok := d.GetOk("windows_profile.0.admin_password"); ok {
		adminPassword = v.(string)
	}

	license := ""
	if profile.LicenseType != nil && *profile.LicenseType != managedclusters.LicenseTypeNone {
		license = string(*profile.LicenseType)
	}

	gmsaProfile := flattenGmsaProfile(profile.GmsaProfile)

	return []any{
		map[string]any{
			"admin_password": adminPassword,
			"admin_username": adminUsername,
			"license":        license,
			"gmsa":           gmsaProfile,
		},
	}
}

func flattenKubernetesClusterWorkloadAutoscalerProfile(profile *managedclusters.ManagedClusterWorkloadAutoScalerProfile) []any {
	// The API always returns an empty WorkloadAutoScalerProfile object even if none of these values have ever been set
	if profile == nil || (profile.Keda == nil && profile.VerticalPodAutoscaler == nil) {
		return []any{}
	}

	kedaEnabled := false
	if v := profile.Keda; v != nil && v.Enabled {
		kedaEnabled = v.Enabled
	}

	vpaEnabled := false
	if v := profile.VerticalPodAutoscaler; v != nil && v.Enabled {
		vpaEnabled = v.Enabled
	}

	return []any{
		map[string]any{
			"keda_enabled":                    kedaEnabled,
			"vertical_pod_autoscaler_enabled": vpaEnabled,
		},
	}
}

func flattenKubernetesClusterNodeProvisioningProfile(profile *managedclusters.ManagedClusterNodeProvisioningProfile) []any {
	if profile == nil {
		return []any{}
	}
	return []any{map[string]any{
		"mode":               pointer.From(profile.Mode),
		"default_node_pools": pointer.From(profile.DefaultNodePools),
	}}
}

func flattenGmsaProfile(profile *managedclusters.WindowsGmsaProfile) []any {
	if profile == nil {
		return []any{}
	}

	dnsServer := pointer.From(profile.DnsServer)

	rootDomainName := pointer.From(profile.RootDomainName)

	return []any{
		map[string]any{
			"dns_server":  dnsServer,
			"root_domain": rootDomainName,
		},
	}
}

func flattenKubernetesClusterAdvancedNetworking(advancedNetworking *managedclusters.AdvancedNetworking) []any {
	if advancedNetworking == nil || !pointer.From(advancedNetworking.Enabled) {
		return []any{}
	}

	observabilityEnabled := false
	if advancedNetworking.Observability != nil {
		observabilityEnabled = pointer.From(advancedNetworking.Observability.Enabled)
	}

	securityEnabled := false
	if advancedNetworking.Security != nil {
		securityEnabled = pointer.From(advancedNetworking.Security.Enabled)
	}

	return []any{
		map[string]any{
			"observability_enabled": observabilityEnabled,
			"security_enabled":      securityEnabled,
		},
	}
}

func flattenKubernetesClusterNetworkProfile(profile *managedclusters.ContainerServiceNetworkProfile) []any {
	if profile == nil {
		return []any{}
	}

	dnsServiceIP := pointer.From(profile.DnsServiceIP)

	serviceCidr := pointer.From(profile.ServiceCidr)

	podCidr := pointer.From(profile.PodCidr)

	networkPlugin := ""
	if profile.NetworkPlugin != nil {
		networkPlugin = string(*profile.NetworkPlugin)
	}

	networkMode := ""
	if profile.NetworkMode != nil {
		networkMode = string(*profile.NetworkMode)
	}

	networkPolicy := ""
	if profile.NetworkPolicy != nil {
		networkPolicy = string(*profile.NetworkPolicy)

		// convert "none" network policy to "", because "none" is only supported in preview api-version
		if networkPolicy == "none" {
			networkPolicy = ""
		}
	}

	outboundType := ""
	if profile.OutboundType != nil {
		outboundType = string(*profile.OutboundType)
	}

	lbProfiles := make([]any, 0)
	if lbp := profile.LoadBalancerProfile; lbp != nil {
		lb := make(map[string]any)

		if v := lbp.AllocatedOutboundPorts; v != nil {
			lb["outbound_ports_allocated"] = v
		}

		if v := lbp.IdleTimeoutInMinutes; v != nil {
			lb["idle_timeout_in_minutes"] = v
		}

		if ips := lbp.ManagedOutboundIPs; ips != nil {
			if count := ips.Count; count != nil {
				lb["managed_outbound_ip_count"] = count
			}

			if countIPv6 := ips.CountIPv6; countIPv6 != nil {
				lb["managed_outbound_ipv6_count"] = countIPv6
			}
		}

		if oip := lbp.OutboundIPs; oip != nil {
			if poip := oip.PublicIPs; poip != nil {
				lb["outbound_ip_address_ids"] = resourceReferencesToIds(poip)
			}
		}

		if oip := lbp.OutboundIPPrefixes; oip != nil {
			if pip := oip.PublicIPPrefixes; pip != nil {
				lb["outbound_ip_prefix_ids"] = resourceReferencesToIds(pip)
			}
		}

		if v := lbp.BackendPoolType; v != nil {
			lb["backend_pool_type"] = v
		}

		lb["effective_outbound_ips"] = resourceReferencesToIds(profile.LoadBalancerProfile.EffectiveOutboundIPs)
		lbProfiles = append(lbProfiles, lb)
	}

	ngwProfiles := make([]any, 0)
	if ngwp := profile.NatGatewayProfile; ngwp != nil {
		ng := make(map[string]any)

		if v := ngwp.IdleTimeoutInMinutes; v != nil {
			ng["idle_timeout_in_minutes"] = v
		}

		if ips := ngwp.ManagedOutboundIPProfile; ips != nil {
			if count := ips.Count; count != nil {
				ng["managed_outbound_ip_count"] = count
			}
		}

		ng["effective_outbound_ips"] = resourceReferencesToIds(profile.NatGatewayProfile.EffectiveOutboundIPs)
		ngwProfiles = append(ngwProfiles, ng)
	}

	ipVersions := make([]any, 0)
	if ipfs := profile.IPFamilies; ipfs != nil {
		for _, item := range *ipfs {
			ipVersions = append(ipVersions, item)
		}
	}

	// TODO - Remove the workaround below once issue https://github.com/Azure/azure-rest-api-specs/issues/18056 is resolved
	sku := profile.LoadBalancerSku
	for _, v := range managedclusters.PossibleValuesForLoadBalancerSku() {
		if strings.EqualFold(v, string(*sku)) {
			sku = pointer.ToEnum[managedclusters.LoadBalancerSku](v)
		}
	}

	networkPluginMode := ""
	if profile.NetworkPluginMode != nil {
		// The returned value has inconsistent casing
		// TODO: Remove the normalization codes once the following issue is fixed.
		// Issue: https://github.com/Azure/azure-rest-api-specs/issues/21810
		if strings.EqualFold(string(*profile.NetworkPluginMode), string(managedclusters.NetworkPluginModeOverlay)) {
			networkPluginMode = string(managedclusters.NetworkPluginModeOverlay)
		}
	}

	networkDataPlane := string(managedclusters.NetworkDataplaneAzure)
	if v := profile.NetworkDataplane; v != nil {
		networkDataPlane = pointer.FromEnum(v)
	}

	advancedNetworking := flattenKubernetesClusterAdvancedNetworking(profile.AdvancedNetworking)

	result := map[string]any{
		"dns_service_ip":        dnsServiceIP,
		"network_data_plane":    networkDataPlane,
		"load_balancer_sku":     string(*sku),
		"load_balancer_profile": lbProfiles,
		"nat_gateway_profile":   ngwProfiles,
		"ip_versions":           ipVersions,
		"network_plugin":        networkPlugin,
		"network_plugin_mode":   networkPluginMode,
		"network_mode":          networkMode,
		"network_policy":        networkPolicy,
		"pod_cidr":              podCidr,
		"pod_cidrs":             pluginsdk.FlattenSlice(profile.PodCidrs),
		"service_cidr":          serviceCidr,
		"service_cidrs":         pluginsdk.FlattenSlice(profile.ServiceCidrs),
		"outbound_type":         outboundType,
		"advanced_networking":   advancedNetworking,
	}

	return []any{result}
}

func flattenKubernetesClusterAzureActiveDirectoryRoleBasedAccessControl(input *managedclusters.ManagedClusterProperties) []any {
	results := make([]any, 0)

	if profile := input.AadProfile; profile != nil {
		adminGroupObjectIds := pluginsdk.FlattenSlice(profile.AdminGroupObjectIDs)
		results = append(results, map[string]any{
			"admin_group_object_ids": adminGroupObjectIds,
			"tenant_id":              pointer.From(profile.TenantID),
			"azure_rbac_enabled":     pointer.From(profile.EnableAzureRBAC),
		})
	}

	return results
}

func flattenAzureRmKubernetesClusterServicePrincipalProfile(profile *managedclusters.ManagedClusterServicePrincipalProfile, d *pluginsdk.ResourceData) []any {
	if profile == nil {
		return []any{}
	}

	clientId := profile.ClientId

	if strings.EqualFold(clientId, "msi") {
		return []any{}
	}

	// client secret isn't returned by the API so pass the existing value along
	clientSecret := ""
	if sp, ok := d.GetOk("service_principal"); ok {
		var val []any

		// prior to 1.34 this was a *pluginsdk.Set, now it's a List - try both
		if v, ok := sp.([]any); ok {
			val = v
		} else if v, ok := sp.(*pluginsdk.Set); ok {
			val = v.List()
		}

		if len(val) > 0 && val[0] != nil {
			raw := val[0].(map[string]any)
			clientSecret = raw["client_secret"].(string)
		}
	}

	return []any{
		map[string]any{
			"client_id":     clientId,
			"client_secret": clientSecret,
		},
	}
}

func flattenKubernetesClusterAutoScalerProfile(profile *managedclusters.ManagedClusterPropertiesAutoScalerProfile) ([]any, error) {
	if profile == nil {
		return []any{}, nil
	}

	balanceSimilarNodeGroups := false
	if profile.BalanceSimilarNodeGroups != nil {
		// @tombuildsstuff: presumably this'll get converted to a Boolean at some point
		//					at any rate we should use the proper type users expect here
		balanceSimilarNodeGroups = strings.EqualFold(*profile.BalanceSimilarNodeGroups, "true")
	}

	daemonsetEvictionForEmptyNodes := pointer.From(profile.DaemonsetEvictionForEmptyNodes)

	daemonsetEvictionForOccupiedNodes := true
	if profile.DaemonsetEvictionForOccupiedNodes != nil {
		daemonsetEvictionForOccupiedNodes = pointer.From(profile.DaemonsetEvictionForOccupiedNodes)
	}

	expander := ""
	if profile.Expander != nil {
		expander = string(*profile.Expander)
	}

	ignoreDaemonsetsUtilization := pointer.From(profile.IgnoreDaemonsetsUtilization)

	maxGracefulTerminationSec := pointer.From(profile.MaxGracefulTerminationSec)

	maxNodeProvisionTime := pointer.From(profile.MaxNodeProvisionTime)

	maxUnreadyNodes := 0
	if profile.OkTotalUnreadyCount != nil {
		var err error
		maxUnreadyNodes, err = strconv.Atoi(*profile.OkTotalUnreadyCount)
		if err != nil {
			return nil, err
		}
	}

	maxUnreadyPercentage := 0.0
	if profile.MaxTotalUnreadyPercentage != nil {
		var err error
		maxUnreadyPercentage, err = strconv.ParseFloat(*profile.MaxTotalUnreadyPercentage, 64)
		if err != nil {
			return nil, err
		}
	}

	newPodScaleUpDelay := pointer.From(profile.NewPodScaleUpDelay)

	scaleDownDelayAfterAdd := pointer.From(profile.ScaleDownDelayAfterAdd)

	scaleDownDelayAfterDelete := pointer.From(profile.ScaleDownDelayAfterDelete)

	scaleDownDelayAfterFailure := pointer.From(profile.ScaleDownDelayAfterFailure)

	scaleDownUnneededTime := pointer.From(profile.ScaleDownUnneededTime)

	scaleDownUnreadyTime := pointer.From(profile.ScaleDownUnreadyTime)

	scaleDownUtilizationThreshold := pointer.From(profile.ScaleDownUtilizationThreshold)

	emptyBulkDeleteMax := pointer.From(profile.MaxEmptyBulkDelete)

	scanInterval := pointer.From(profile.ScanInterval)

	skipNodesWithLocalStorage := true
	if profile.SkipNodesWithLocalStorage != nil {
		skipNodesWithLocalStorage = strings.EqualFold(*profile.SkipNodesWithLocalStorage, "true")
	}

	skipNodesWithSystemPods := true
	if profile.SkipNodesWithSystemPods != nil {
		skipNodesWithSystemPods = strings.EqualFold(*profile.SkipNodesWithSystemPods, "true")
	}

	return []any{
		map[string]any{
			"balance_similar_node_groups":                   balanceSimilarNodeGroups,
			"daemonset_eviction_for_empty_nodes_enabled":    daemonsetEvictionForEmptyNodes,
			"daemonset_eviction_for_occupied_nodes_enabled": daemonsetEvictionForOccupiedNodes,
			"expander":                              expander,
			"ignore_daemonsets_utilization_enabled": ignoreDaemonsetsUtilization,
			"max_graceful_termination_sec":          maxGracefulTerminationSec,
			"max_node_provisioning_time":            maxNodeProvisionTime,
			"max_unready_nodes":                     maxUnreadyNodes,
			"max_unready_percentage":                maxUnreadyPercentage,
			"new_pod_scale_up_delay":                newPodScaleUpDelay,
			"scale_down_delay_after_add":            scaleDownDelayAfterAdd,
			"scale_down_delay_after_delete":         scaleDownDelayAfterDelete,
			"scale_down_delay_after_failure":        scaleDownDelayAfterFailure,
			"scale_down_unneeded":                   scaleDownUnneededTime,
			"scale_down_unready":                    scaleDownUnreadyTime,
			"scale_down_utilization_threshold":      scaleDownUtilizationThreshold,
			"empty_bulk_delete_max":                 emptyBulkDeleteMax,
			"scan_interval":                         scanInterval,
			"skip_nodes_with_local_storage":         skipNodesWithLocalStorage,
			"skip_nodes_with_system_pods":           skipNodesWithSystemPods,
		},
	}, nil
}

func flattenKubernetesClusterMaintenanceConfiguration(input *maintenanceconfigurations.MaintenanceWindow) any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	startDate := ""
	if input.StartDate != nil {
		startDate = *input.StartDate + "T00:00:00Z"
	}
	utcOffset := pointer.From(input.UtcOffset)

	windowProperties := map[string]any{
		"not_allowed": flattenKubernetesClusterMaintenanceConfigurationDateSpans(input.NotAllowedDates),
		"duration":    int(input.DurationHours),
		"start_date":  startDate,
		"start_time":  input.StartTime,
		"utc_offset":  utcOffset,
	}
	// Add flattened schedule properties
	maps.Copy(windowProperties, flattenKubernetesClusterMaintenanceConfigurationSchedule(input.Schedule))

	return append(results, windowProperties)
}

func flattenKubernetesClusterMaintenanceConfigurationSchedule(input maintenanceconfigurations.Schedule) map[string]any {
	frequency := ""
	interval := int64(0)
	if input.Daily != nil {
		frequency = "Daily"
		interval = input.Daily.IntervalDays
	}

	dayOfWeek := ""
	if input.Weekly != nil {
		frequency = "Weekly"
		interval = input.Weekly.IntervalWeeks
		dayOfWeek = string(input.Weekly.DayOfWeek)
	}

	dayOfMonth := 0
	if input.AbsoluteMonthly != nil {
		frequency = "AbsoluteMonthly"
		interval = input.AbsoluteMonthly.IntervalMonths
		dayOfMonth = int(input.AbsoluteMonthly.DayOfMonth)
	}

	weekIndex := ""
	if input.RelativeMonthly != nil {
		frequency = "RelativeMonthly"
		interval = input.RelativeMonthly.IntervalMonths
		dayOfWeek = string(input.RelativeMonthly.DayOfWeek)
		weekIndex = string(input.RelativeMonthly.WeekIndex)
	}

	return map[string]any{
		"frequency":    frequency,
		"interval":     interval,
		"day_of_week":  dayOfWeek,
		"week_index":   weekIndex,
		"day_of_month": dayOfMonth,
	}
}

func flattenKubernetesClusterMaintenanceConfigurationDefault(input *maintenanceconfigurations.MaintenanceConfigurationProperties) any {
	results := make([]any, 0)
	if input == nil {
		return results
	}
	results = append(results, map[string]any{
		"not_allowed": flattenKubernetesClusterMaintenanceConfigurationTimeSpans(input.NotAllowedTime),
		"allowed":     flattenKubernetesClusterMaintenanceConfigurationTimeInWeeks(input.TimeInWeek),
	})
	return results
}

func flattenKubernetesClusterMaintenanceConfigurationTimeSpans(input *[]maintenanceconfigurations.TimeSpan) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		end := pointer.From(item.End)
		start := pointer.From(item.Start)
		results = append(results, map[string]any{
			"end":   end,
			"start": start,
		})
	}
	return results
}

func flattenKubernetesClusterMaintenanceConfigurationDateSpans(input *[]maintenanceconfigurations.DateSpan) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var end string
		if item.End != "" {
			end = item.End
		}
		var start string
		if item.Start != "" {
			start = item.Start
		}
		results = append(results, map[string]any{
			"end":   end + "T00:00:00Z",
			"start": start + "T00:00:00Z",
		})
	}
	return results
}

func flattenKubernetesClusterMaintenanceConfigurationTimeInWeeks(input *[]maintenanceconfigurations.TimeInWeek) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		day := ""
		if item.Day != nil {
			day = string(*item.Day)
		}
		results = append(results, map[string]any{
			"day":   day,
			"hours": pluginsdk.FlattenSlice(item.HourSlots),
		})
	}
	return results
}

func flattenKubernetesClusterHttpProxyConfig(props *managedclusters.ManagedClusterProperties) []any {
	if props == nil || props.HTTPProxyConfig == nil {
		return []any{}
	}

	httpProxyConfig := props.HTTPProxyConfig

	httpProxy := pointer.From(httpProxyConfig.HTTPProxy)

	httpsProxy := pointer.From(httpProxyConfig.HTTPSProxy)

	noProxyList := make([]string, 0)
	if httpProxyConfig.NoProxy != nil {
		noProxyList = append(noProxyList, *httpProxyConfig.NoProxy...)
	}

	trustedCa := pointer.From(httpProxyConfig.TrustedCa)

	results := []any{}
	return append(results, map[string]any{
		"http_proxy":  httpProxy,
		"https_proxy": httpsProxy,
		"no_proxy":    noProxyList,
		"trusted_ca":  trustedCa,
	})
}

func flattenKubernetesClusterMicrosoftDefender(input *managedclusters.ManagedClusterSecurityProfile) []any {
	if input == nil || input.Defender == nil || input.Defender.SecurityMonitoring == nil || input.Defender.SecurityMonitoring.Enabled == nil || !*input.Defender.SecurityMonitoring.Enabled {
		return []any{}
	}

	return []any{
		map[string]any{
			"log_analytics_workspace_id": pointer.From(input.Defender.LogAnalyticsWorkspaceResourceId),
		},
	}
}

func flattenKubernetesClusterIngressProfile(input *managedclusters.ManagedClusterIngressProfile) []any {
	if input == nil || input.WebAppRouting == nil || (input.WebAppRouting.Enabled != nil && !*input.WebAppRouting.Enabled) {
		return []any{}
	}

	out := map[string]any{}
	out["dns_zone_ids"] = pluginsdk.FlattenSlice(input.WebAppRouting.DnsZoneResourceIds)

	webAppRoutingIdentity := []any{}

	if v := input.WebAppRouting.Identity; v != nil {
		webAppRoutingIdentity = flattenKubernetesClusterAddOnIdentityProfile(v)
	}

	out["web_app_routing_identity"] = webAppRoutingIdentity

	if input.WebAppRouting.Nginx != nil {
		out["default_nginx_controller"] = pointer.From(input.WebAppRouting.Nginx.DefaultIngressControllerType)
	}

	return []any{out}
}

func flattenKubernetesClusterAzureServiceMeshProfile(input *managedclusters.ServiceMeshProfile) []any {
	if input == nil || input.Mode != managedclusters.ServiceMeshModeIstio {
		return []any{}
	}

	returnMap := map[string]any{
		"mode": string(managedclusters.ServiceMeshModeIstio),
	}

	if (input.Istio.Components.IngressGateways != nil) && len(*input.Istio.Components.IngressGateways) > 0 {
		for _, value := range *input.Istio.Components.IngressGateways {
			mode := value.Mode

			var currentIngressKey string

			if mode == managedclusters.IstioIngressGatewayModeExternal {
				currentIngressKey = "external_ingress_gateway_enabled"
			} else {
				currentIngressKey = "internal_ingress_gateway_enabled"
			}

			returnMap[currentIngressKey] = value.Enabled
		}
	}

	if input.Istio.CertificateAuthority != nil {
		returnMap["certificate_authority"] = flattenKubernetesClusterServiceMeshProfileCertificateAuthority(input.Istio.CertificateAuthority)
	}

	if input.Istio.Revisions != nil {
		returnMap["revisions"] = pluginsdk.FlattenSlice(input.Istio.Revisions)
	}

	return []any{returnMap}
}

func flattenKubernetesClusterServiceMeshProfileCertificateAuthority(certificateAuthority *managedclusters.IstioCertificateAuthority) any {
	if certificateAuthority == nil || certificateAuthority.Plugin == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"key_vault_id":           pointer.From(certificateAuthority.Plugin.KeyVaultId),
			"root_cert_object_name":  pointer.From(certificateAuthority.Plugin.RootCertObjectName),
			"cert_chain_object_name": pointer.From(certificateAuthority.Plugin.CertChainObjectName),
			"cert_object_name":       pointer.From(certificateAuthority.Plugin.CertObjectName),
			"key_object_name":        pointer.From(certificateAuthority.Plugin.KeyObjectName),
		},
	}
}

func flattenKubernetesClusterAzureMonitorProfile(input *managedclusters.ManagedClusterAzureMonitorProfile) []any {
	if input == nil || input.Metrics == nil || !input.Metrics.Enabled {
		return []any{}
	}
	if input.Metrics.KubeStateMetrics == nil {
		return []any{
			map[string]any{},
		}
	}
	annotationAllowList := pointer.From(input.Metrics.KubeStateMetrics.MetricAnnotationsAllowList)
	labelAllowList := pointer.From(input.Metrics.KubeStateMetrics.MetricLabelsAllowlist)
	return []any{
		map[string]any{
			"annotations_allowed": annotationAllowList,
			"labels_allowed":      labelAllowList,
		},
	}
}

func flattenKubernetesClusterMetricsProfile(input *managedclusters.ManagedClusterMetricsProfile) bool {
	if input == nil || input.CostAnalysis == nil {
		return false
	}
	return pointer.From(input.CostAnalysis.Enabled)
}

func flattenKubernetesClusterUpgradeOverrideSetting(input *managedclusters.ClusterUpgradeSettings) []any {
	if input == nil || input.OverrideSettings == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"effective_until":       pointer.From(input.OverrideSettings.Until),
			"force_upgrade_enabled": pointer.From(input.OverrideSettings.ForceUpgrade),
		},
	}
}
