// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKubernetesClusterUpdate(d *pluginsdk.ResourceData, meta any) error {
	containersClient := meta.(*clients.Client).Containers
	nodePoolsClient := containersClient.AgentPoolsClient
	clusterClient := containersClient.KubernetesClustersClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	env := containersClient.Environment
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKubernetesClusterID(d.Id())
	if err != nil {
		return err
	}

	// we need to conditionally update the cluster
	existing, err := clusterClient.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}
	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving existing %s: `properties` was nil", *id)
	}
	props := existing.Model.Properties

	if err := helpers.ValidateKubernetesCluster(d, existing.Model, id.ResourceGroupName, id.ManagedClusterName); err != nil {
		return err
	}

	// when update, we should set the value of `Identity.UserAssignedIdentities` empty
	// otherwise the rest api will report error - this is tracked here: https://github.com/Azure/azure-rest-api-specs/issues/13631
	if existing.Model.Identity != nil && existing.Model.Identity.IdentityIds != nil {
		for k := range existing.Model.Identity.IdentityIds {
			existing.Model.Identity.IdentityIds[k] = identity.UserAssignedIdentityDetails{}
		}
	}

	if d.HasChange("service_principal") && !d.HasChange("identity") {
		servicePrincipals := d.Get("service_principal").([]any)
		// we'll be rotating the Service Principal - removing the SP block is handled by the validate function
		servicePrincipalRaw := servicePrincipals[0].(map[string]any)

		clientId := servicePrincipalRaw["client_id"].(string)
		clientSecret := servicePrincipalRaw["client_secret"].(string)
		params := managedclusters.ManagedClusterServicePrincipalProfile{
			ClientId: clientId,
			Secret:   pointer.To(clientSecret),
		}

		if err = clusterClient.ResetServicePrincipalProfileThenPoll(ctx, *id, params); err != nil {
			return fmt.Errorf("updating Service Principal for %s: %+v", *id, err)
		}

		// since we're patching it, re-retrieve the latest version of the cluster
		existing, err = clusterClient.Get(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving updated %s: %+v", *id, err)
		}
		if existing.Model == nil || existing.Model.Properties == nil {
			return fmt.Errorf("retrieving updated %s: `properties` was nil", *id)
		}
	}

	// since there's multiple reasons why we could be called into Update, we use this to only update if something's changed that's not SP/Version
	updateCluster := false

	// RBAC profile updates need to be handled atomically before any call to createUpdate as a diff there will create a PropertyChangeNotAllowed error
	if d.HasChange("role_based_access_control_enabled") {
		// check if we can determine current EnableRBAC state - don't do anything destructive if we can't be sure
		if props.EnableRBAC == nil {
			return fmt.Errorf("updating %s: RBAC Enabled was nil", *id)
		}
		rbacEnabled := d.Get("role_based_access_control_enabled").(bool)

		// changing rbacEnabled must still force cluster recreation
		if *props.EnableRBAC == rbacEnabled {
			props.EnableRBAC = pointer.To(rbacEnabled)
		} else {
			updateCluster = true
		}
	}

	if d.HasChange("azure_active_directory_role_based_access_control") {
		azureADRaw := d.Get("azure_active_directory_role_based_access_control").([]any)
		azureADProfile := expandKubernetesClusterAzureActiveDirectoryRoleBasedAccessControl(azureADRaw)

		props.AadProfile = azureADProfile
		if props.AadProfile == nil || (props.AadProfile.Managed == nil || !*props.AadProfile.Managed) {
			props.AadProfile = &managedclusters.ManagedClusterAADProfile{}
			if err = clusterClient.ResetAADProfileThenPoll(ctx, *id, *props.AadProfile); err != nil {
				return fmt.Errorf("updating Managed Kubernetes Cluster AAD Profile for %s: %+v", *id, err)
			}
		}

		if props.AadProfile != nil && props.AadProfile.Managed != nil && *props.AadProfile.Managed {
			existing.Model.Properties.AadProfile = azureADProfile
			updateCluster = true
		}
	}

	// lintignore:R019 // deliberate subset: only the addon-related fields that require rebuilding the addonProfiles payload
	if d.HasChanges(
		"aci_connector_linux",
		"azure_policy_enabled",
		"confidential_computing",
		"http_application_routing_enabled",
		"oms_agent",
		"ingress_application_gateway",
		"open_service_mesh_enabled",
		"key_vault_secrets_provider",
	) {
		updateCluster = true
		addOns := collectKubernetesAddons(d)
		addonProfiles, err := expandKubernetesAddOns(d, addOns, env)
		if err != nil {
			return err
		}
		existing.Model.Properties.AddonProfiles = addonProfiles
	}

	if d.HasChanges("run_command_enabled", "private_cluster_public_fqdn_enabled", "api_server_access_profile") {
		updateCluster = true

		existing.Model.Properties.ApiServerAccessProfile = expandKubernetesClusterAPIAccessProfile(d)
	}

	if d.HasChange("auto_scaler_profile") {
		updateCluster = true
		autoScalerProfileRaw := d.Get("auto_scaler_profile").([]any)

		existing.Model.Properties.AutoScalerProfile = expandKubernetesClusterAutoScalerProfile(autoScalerProfileRaw)
	}

	if d.HasChange("monitor_metrics") {
		updateCluster = true
		azureMonitorKubernetesMetricsRaw := d.Get("monitor_metrics").([]any)

		existing.Model.Properties.AzureMonitorProfile = expandKubernetesClusterAzureMonitorProfile(azureMonitorKubernetesMetricsRaw)
	}

	if d.HasChange("linux_profile") {
		updateCluster = true
		linuxProfileRaw := d.Get("linux_profile").([]any)
		existing.Model.Properties.LinuxProfile = expandKubernetesClusterLinuxProfile(linuxProfileRaw)
	}

	if d.HasChange("local_account_disabled") {
		updateCluster = true
		existing.Model.Properties.DisableLocalAccounts = pointer.To(d.Get("local_account_disabled").(bool))
	}

	if d.HasChange("cost_analysis_enabled") {
		updateCluster = true
		metricsProfile, err := expandKubernetesClusterMetricsProfile(d.Get("cost_analysis_enabled").(bool), d.Get("sku_tier").(string))
		if err != nil {
			return err
		}
		existing.Model.Properties.MetricsProfile = metricsProfile
	}

	if d.HasChange("network_profile") {
		updateCluster = true

		if existing.Model.Properties.NetworkProfile == nil {
			return fmt.Errorf("updating %s: `network_profile` was nil", *id)
		}

		networkProfile := *existing.Model.Properties.NetworkProfile

		if networkProfile.LoadBalancerProfile != nil {
			loadBalancerProfile := *networkProfile.LoadBalancerProfile

			if key := "network_profile.0.network_policy"; d.HasChange(key) {
				networkPolicy := d.Get(key).(string)
				existing.Model.Properties.NetworkProfile.NetworkPolicy = pointer.ToEnum[managedclusters.NetworkPolicy](networkPolicy)
			}

			if key := "network_profile.0.load_balancer_profile.0.effective_outbound_ips"; d.HasChange(key) {
				loadBalancerProfile.EffectiveOutboundIPs = idsToResourceReferences(d.Get(key))
			}

			if key := "network_profile.0.load_balancer_profile.0.idle_timeout_in_minutes"; d.HasChange(key) {
				idleTimeoutInMinutes := d.Get(key).(int)
				loadBalancerProfile.IdleTimeoutInMinutes = pointer.To(int64(idleTimeoutInMinutes))
			}

			if key := "network_profile.0.load_balancer_profile.0.managed_outbound_ip_count"; d.HasChange(key) {
				managedOutboundIPCount := d.Get(key).(int)
				loadBalancerProfile.ManagedOutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileManagedOutboundIPs{
					Count: pointer.To(int64(managedOutboundIPCount)),
				}

				// fixes: Load balancer profile must specify one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs.
				loadBalancerProfile.OutboundIPs = nil
				loadBalancerProfile.OutboundIPPrefixes = nil
			}

			if key := "network_profile.0.load_balancer_profile.0.managed_outbound_ipv6_count"; d.HasChange(key) {
				managedOutboundIPV6Count := d.Get(key).(int)
				if loadBalancerProfile.ManagedOutboundIPs == nil {
					loadBalancerProfile.ManagedOutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileManagedOutboundIPs{}
				}
				loadBalancerProfile.ManagedOutboundIPs.CountIPv6 = pointer.To(int64(managedOutboundIPV6Count))

				// fixes: Load balancer profile must specify one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs.
				loadBalancerProfile.OutboundIPs = nil
				loadBalancerProfile.OutboundIPPrefixes = nil
			}

			if key := "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"; d.HasChange(key) {
				outboundIPAddress := d.Get(key)
				if v := outboundIPAddress.(*pluginsdk.Set).List(); len(v) == 0 {
					// sending [] to unset `outbound_ip_address_ids` results in 400 / Bad Request
					// instead we default back to AKS managed outbound which is the default of the AKS API when nothing is provided
					loadBalancerProfile.ManagedOutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileManagedOutboundIPs{
						Count: pointer.To(int64(1)),
					}
					loadBalancerProfile.OutboundIPs = nil
					loadBalancerProfile.OutboundIPPrefixes = nil
				} else {
					publicIPAddressIDs := idsToResourceReferences(d.Get(key))
					loadBalancerProfile.OutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileOutboundIPs{
						PublicIPs: publicIPAddressIDs,
					}

					// fixes: Load balancer profile must specify one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs.
					loadBalancerProfile.ManagedOutboundIPs = nil
					loadBalancerProfile.OutboundIPPrefixes = nil
				}
			}

			if key := "network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids"; d.HasChange(key) {
				outboundIPPrefixes := d.Get(key)
				if v := outboundIPPrefixes.(*pluginsdk.Set).List(); len(v) == 0 {
					// sending [] to unset `outbound_ip_address_ids` results in 400 / Bad Request
					// instead we default back to AKS managed outbound which is the default of the AKS API when nothing is specified
					loadBalancerProfile.ManagedOutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileManagedOutboundIPs{
						Count: pointer.To(int64(1)),
					}
					loadBalancerProfile.OutboundIPs = nil
					loadBalancerProfile.OutboundIPPrefixes = nil
				} else {
					outboundIPPrefixIDs := idsToResourceReferences(d.Get(key))
					loadBalancerProfile.OutboundIPPrefixes = &managedclusters.ManagedClusterLoadBalancerProfileOutboundIPPrefixes{
						PublicIPPrefixes: outboundIPPrefixIDs,
					}

					// fixes: Load balancer profile must specify one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs.
					loadBalancerProfile.ManagedOutboundIPs = nil
					loadBalancerProfile.OutboundIPs = nil
				}
			}

			if key := "network_profile.0.load_balancer_profile.0.outbound_ports_allocated"; d.HasChange(key) {
				allocatedOutboundPorts := d.Get(key).(int)
				loadBalancerProfile.AllocatedOutboundPorts = pointer.To(int64(allocatedOutboundPorts))
			}

			if key := "network_profile.0.load_balancer_profile.0.backend_pool_type"; d.HasChange(key) {
				backendPoolType := d.Get(key).(string)
				loadBalancerProfile.BackendPoolType = pointer.ToEnum[managedclusters.BackendPoolType](backendPoolType)
			}

			existing.Model.Properties.NetworkProfile.LoadBalancerProfile = &loadBalancerProfile
		}

		if networkProfile.NatGatewayProfile != nil {
			natGatewayProfile := *networkProfile.NatGatewayProfile

			if key := "network_profile.0.nat_gateway_profile.0.idle_timeout_in_minutes"; d.HasChange(key) {
				idleTimeoutInMinutes := d.Get(key).(int)
				natGatewayProfile.IdleTimeoutInMinutes = pointer.To(int64(idleTimeoutInMinutes))
			}

			if key := "network_profile.0.nat_gateway_profile.0.managed_outbound_ip_count"; d.HasChange(key) {
				managedOutboundIPCount := d.Get(key).(int)
				natGatewayProfile.ManagedOutboundIPProfile = &managedclusters.ManagedClusterManagedOutboundIPProfile{
					Count: pointer.To(int64(managedOutboundIPCount)),
				}
				natGatewayProfile.EffectiveOutboundIPs = nil
			}

			existing.Model.Properties.NetworkProfile.NatGatewayProfile = &natGatewayProfile
		}

		if d.HasChange("network_profile.0.network_data_plane") {
			existing.Model.Properties.NetworkProfile.NetworkDataplane = pointer.ToEnum[managedclusters.NetworkDataplane](d.Get("network_profile.0.network_data_plane").(string))
		}

		if key := "network_profile.0.network_plugin"; d.HasChange(key) {
			networkPlugin := d.Get(key).(string)
			existing.Model.Properties.NetworkProfile.NetworkPlugin = pointer.ToEnum[managedclusters.NetworkPlugin](networkPlugin)
		}

		if key := "network_profile.0.network_plugin_mode"; d.HasChange(key) {
			networkPluginMode := d.Get(key).(string)
			existing.Model.Properties.NetworkProfile.NetworkPluginMode = pointer.ToEnum[managedclusters.NetworkPluginMode](networkPluginMode)
		}

		if key := "network_profile.0.pod_cidr"; d.HasChange(key) {
			podCidr := d.Get(key).(string)
			existing.Model.Properties.NetworkProfile.PodCidr = pointer.To(podCidr)
		}

		if key := "network_profile.0.pod_cidrs"; d.HasChange(key) {
			podCidrs := d.Get(key).([]any)
			existing.Model.Properties.NetworkProfile.PodCidrs = pluginsdk.ExpandStringSlice(podCidrs)
		}

		if key := "network_profile.0.outbound_type"; d.HasChange(key) {
			outboundType := managedclusters.OutboundType(d.Get(key).(string))
			existing.Model.Properties.NetworkProfile.OutboundType = pointer.To(outboundType)
			if outboundType != managedclusters.OutboundTypeLoadBalancer {
				existing.Model.Properties.NetworkProfile.LoadBalancerProfile = nil
			}
			if outboundType != managedclusters.OutboundTypeManagedNATGateway && outboundType != managedclusters.OutboundTypeUserAssignedNATGateway {
				existing.Model.Properties.NetworkProfile.NatGatewayProfile = nil
			}
		}

		if d.HasChange("network_profile.0.advanced_networking") {
			existing.Model.Properties.NetworkProfile.AdvancedNetworking = expandKubernetesClusterAdvancedNetworking(d.Get("network_profile.0.advanced_networking").([]any), d)
		}
	}
	if d.HasChange("service_mesh_profile") {
		updateCluster = true
		if serviceMeshProfile := expandKubernetesClusterServiceMeshProfile(d.Get("service_mesh_profile").([]any), existing.Model.Properties.ServiceMeshProfile); serviceMeshProfile != nil {
			existing.Model.Properties.ServiceMeshProfile = serviceMeshProfile
		}
	}

	if d.HasChange("tags") {
		updateCluster = true
		t := d.Get("tags").(map[string]any)
		existing.Model.Tags = tags.Expand(t)
	}

	if d.HasChange("windows_profile") {
		updateCluster = true
		windowsProfileRaw := d.Get("windows_profile").([]any)
		existing.Model.Properties.WindowsProfile = expandKubernetesClusterWindowsProfile(windowsProfileRaw)
	}

	if d.HasChange("identity") {
		updateCluster = true
		managedClusterIdentityRaw := d.Get("identity").([]any)

		expandedIdentity, err := expandKubernetesClusterManagedClusterIdentity(managedClusterIdentityRaw)
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		existing.Model.Identity = expandedIdentity
	}

	if d.HasChange("sku_tier") {
		updateCluster = true
		if existing.Model.Sku == nil {
			existing.Model.Sku = &managedclusters.ManagedClusterSKU{
				Name: pointer.To(managedclusters.ManagedClusterSKUNameBase),
			}
		}

		skuTier := managedclusters.ManagedClusterSKUTierFree
		if v := d.Get("sku_tier").(string); v != "" {
			skuTier = managedclusters.ManagedClusterSKUTier(v)
		}
		existing.Model.Sku.Tier = &skuTier
	}

	autoUpgradeChannel := "automatic_upgrade_channel"
	nodeOsUpgradeChannel := "node_os_upgrade_channel"
	if d.HasChange(autoUpgradeChannel) {
		updateCluster = true
		if existing.Model.Properties.AutoUpgradeProfile == nil {
			existing.Model.Properties.AutoUpgradeProfile = &managedclusters.ManagedClusterAutoUpgradeProfile{}
		}
		channel := d.Get(autoUpgradeChannel).(string)
		if channel == "" {
			channel = string(managedclusters.UpgradeChannelNone)
		}
		existing.Model.Properties.AutoUpgradeProfile.UpgradeChannel = pointer.ToEnum[managedclusters.UpgradeChannel](channel)
	}

	if d.HasChange(nodeOsUpgradeChannel) {
		updateCluster = true
		if d.Get(autoUpgradeChannel).(string) == string(managedclusters.UpgradeChannelNodeNegativeimage) && d.Get(nodeOsUpgradeChannel).(string) != string(managedclusters.NodeOSUpgradeChannelNodeImage) {
			return fmt.Errorf("`node_os_upgrade_channel` must be set to `NodeImage` if `automatic_upgrade_channel` is set to `node-image`")
		}
		if existing.Model.Properties.AutoUpgradeProfile == nil {
			existing.Model.Properties.AutoUpgradeProfile = &managedclusters.ManagedClusterAutoUpgradeProfile{}
		}
		existing.Model.Properties.AutoUpgradeProfile.NodeOSUpgradeChannel = pointer.ToEnum[managedclusters.NodeOSUpgradeChannel](d.Get(nodeOsUpgradeChannel).(string))
	}

	if d.HasChange("http_proxy_config") {
		updateCluster = true
		httpProxyConfigRaw := d.Get("http_proxy_config").([]any)
		existing.Model.Properties.HTTPProxyConfig = expandKubernetesClusterHttpProxyConfig(httpProxyConfigRaw)
	}

	if d.HasChange("oidc_issuer_enabled") {
		updateCluster = true
		oidcIssuerEnabled := d.Get("oidc_issuer_enabled").(bool)
		existing.Model.Properties.OidcIssuerProfile = expandKubernetesClusterOidcIssuerProfile(oidcIssuerEnabled)
	}

	if d.HasChanges("key_management_service") {
		updateCluster = true
		azureKeyVaultKmsRaw := d.Get("key_management_service").([]any)
		azureKeyVaultKms, _ := expandKubernetesClusterAzureKeyVaultKms(ctx, keyVaultsClient, id.SubscriptionId, d, azureKeyVaultKmsRaw)
		if existing.Model.Properties.SecurityProfile == nil {
			existing.Model.Properties.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
		}
		existing.Model.Properties.SecurityProfile.AzureKeyVaultKms = azureKeyVaultKms
	}

	if d.HasChanges("custom_ca_trust_certificates_base64") {
		updateCluster = true
		customCaTrustCertListRaw := d.Get("custom_ca_trust_certificates_base64").([]any)
		if existing.Model.Properties.SecurityProfile == nil {
			existing.Model.Properties.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
		}
		existing.Model.Properties.SecurityProfile.CustomCATrustCertificates = pluginsdk.ExpandStringSlice(customCaTrustCertListRaw)
	}

	if d.HasChanges("microsoft_defender") {
		updateCluster = true
		microsoftDefenderRaw := d.Get("microsoft_defender").([]any)
		if existing.Model.Properties.SecurityProfile == nil {
			existing.Model.Properties.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
		}
		existing.Model.Properties.SecurityProfile.Defender = expandKubernetesClusterMicrosoftDefender(d, microsoftDefenderRaw)
	}

	if d.HasChanges("storage_profile") {
		updateCluster = true
		storageProfileRaw := d.Get("storage_profile").([]any)
		existing.Model.Properties.StorageProfile = expandStorageProfile(storageProfileRaw)
	}

	if d.HasChange("workload_autoscaler_profile") {
		updateCluster = true
		workloadAutoscalerProfileRaw := d.Get("workload_autoscaler_profile").([]any)
		workloadAutoscalerProfile := expandKubernetesClusterWorkloadAutoscalerProfile(workloadAutoscalerProfileRaw, d)
		if workloadAutoscalerProfile == nil {
			existing.Model.Properties.WorkloadAutoScalerProfile = &managedclusters.ManagedClusterWorkloadAutoScalerProfile{
				Keda: &managedclusters.ManagedClusterWorkloadAutoScalerProfileKeda{},
			}
		} else {
			existing.Model.Properties.WorkloadAutoScalerProfile = workloadAutoscalerProfile
		}
	}

	if d.HasChange("ai_toolchain_operator_enabled") {
		updateCluster = true
		existing.Model.Properties.AiToolchainOperatorProfile = &managedclusters.ManagedClusterAIToolchainOperatorProfile{
			Enabled: pointer.To(d.Get("ai_toolchain_operator_enabled").(bool)),
		}
	}

	if d.HasChange("node_provisioning_profile") {
		updateCluster = true
		nodeProvisioningProfileRaw := d.Get("node_provisioning_profile").([]any)
		existing.Model.Properties.NodeProvisioningProfile = expandKubernetesClusterNodeProvisioningProfile(nodeProvisioningProfileRaw)
	}

	if d.HasChanges("workload_identity_enabled") {
		updateCluster = true
		workloadIdentity := d.Get("workload_identity_enabled").(bool)
		if existing.Model.Properties.SecurityProfile == nil {
			existing.Model.Properties.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
		}
		existing.Model.Properties.SecurityProfile.WorkloadIdentity = &managedclusters.ManagedClusterSecurityProfileWorkloadIdentity{
			Enabled: &workloadIdentity,
		}
	}

	if d.HasChanges("image_cleaner_enabled", "image_cleaner_interval_hours") {
		updateCluster = true
		if existing.Model.Properties.SecurityProfile == nil {
			existing.Model.Properties.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
		}
		existing.Model.Properties.SecurityProfile.ImageCleaner = &managedclusters.ManagedClusterSecurityProfileImageCleaner{
			Enabled:       pointer.To(d.Get("image_cleaner_enabled").(bool)),
			IntervalHours: pointer.To(int64(d.Get("image_cleaner_interval_hours").(int))),
		}
	}

	if d.HasChange("bootstrap_profile") {
		bootstrapProfileRaw := d.Get("bootstrap_profile").([]any)
		profile := expandBootstrapProfile(bootstrapProfileRaw)

		// If profile is removed in the config, we should set ArtifactSource to Direct as it's the default value in the service side.
		if profile == nil {
			profile = &managedclusters.ManagedClusterBootstrapProfile{
				ArtifactSource: pointer.To(managedclusters.ArtifactSourceDirect),
			}
		}

		updateCluster = true
		existing.Model.Properties.BootstrapProfile = profile
	}

	if d.HasChange("upgrade_override") {
		upgradeOverrideSettingRaw := d.Get("upgrade_override").([]any)

		if len(upgradeOverrideSettingRaw) == 0 {
			return fmt.Errorf("`upgrade_override` cannot be unset")
		}

		updateCluster = true
		existing.Model.Properties.UpgradeSettings = expandKubernetesClusterUpgradeOverrideSetting(upgradeOverrideSettingRaw)
	}

	if d.HasChange("web_app_routing") {
		updateCluster = true
		existing.Model.Properties.IngressProfile = expandKubernetesClusterIngressProfile(d, d.Get("web_app_routing").([]any))
	}

	if d.HasChange("support_plan") {
		updateCluster = true
		existing.Model.Properties.SupportPlan = pointer.ToEnum[managedclusters.KubernetesSupportPlan](d.Get("support_plan").(string))
	}

	if updateCluster {
		// If Defender was explicitly disabled in a prior update then we should strip SecurityProfile.AzureDefender from the request
		// body to prevent errors in cases where Defender is disabled for the entire subscription
		if !d.HasChanges("microsoft_defender") && len(d.Get("microsoft_defender").([]any)) == 0 {
			if existing.Model.Properties.SecurityProfile == nil {
				existing.Model.Properties.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
			}
			existing.Model.Properties.SecurityProfile.Defender = nil
		}

		if err = clusterClient.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, managedclusters.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}
	}

	// then roll the version of Kubernetes if necessary
	if d.HasChange("kubernetes_version") {
		existing, err = clusterClient.Get(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving existing %s: %+v", *id, err)
		}
		if existing.Model == nil || existing.Model.Properties == nil {
			return fmt.Errorf("retrieving existing %s: `properties` was nil", *id)
		}

		kubernetesVersion := d.Get("kubernetes_version").(string)
		log.Printf("[DEBUG] Upgrading the version of Kubernetes to %q..", kubernetesVersion)
		existing.Model.Properties.KubernetesVersion = pointer.To(kubernetesVersion)

		if err = clusterClient.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, managedclusters.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating Kubernetes Version for %s: %+v", *id, err)
		}

		log.Printf("[DEBUG] Upgraded the version of Kubernetes to %q..", kubernetesVersion)
	}

	// update the node pool using the separate API
	if d.HasChange("default_node_pool") {
		agentProfiles, err := helpers.ExpandDefaultNodePool(d)
		if err != nil {
			return fmt.Errorf("expanding `default_node_pool`: %+v", err)
		}
		agentProfile := helpers.ConvertDefaultNodePoolToAgentPool(agentProfiles)
		defaultNodePoolId := agentpools.NewAgentPoolID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, *agentProfile.Name)

		// if a users specified a version - confirm that version is supported on the cluster
		if nodePoolVersion := agentProfile.Properties.CurrentOrchestratorVersion; nodePoolVersion != nil {
			existingNodePool, err := nodePoolsClient.Get(ctx, defaultNodePoolId)
			if err != nil {
				return fmt.Errorf("retrieving Default Node Pool %s: %+v", defaultNodePoolId, err)
			}
			currentNodePoolVersion := pointer.From(existingNodePool.Model.Properties.OrchestratorVersion)

			if err := helpers.ValidateNodePoolSupportsVersion(ctx, containersClient, currentNodePoolVersion, defaultNodePoolId, *nodePoolVersion); err != nil {
				return err
			}
		}

		hostEncryptionEnabled := "default_node_pool.0.host_encryption_enabled"
		nodePublicIpEnabled := "default_node_pool.0.node_public_ip_enabled"

		cycleNodePoolProperties := []string{
			"default_node_pool.0.name",
			hostEncryptionEnabled,
			nodePublicIpEnabled,
			"default_node_pool.0.fips_enabled",
			"default_node_pool.0.kubelet_config",
			"default_node_pool.0.kubelet_disk_type",
			"default_node_pool.0.linux_os_config",
			"default_node_pool.0.max_pods",
			"default_node_pool.0.only_critical_addons_enabled",
			"default_node_pool.0.os_disk_size_gb",
			"default_node_pool.0.os_disk_type",
			"default_node_pool.0.pod_subnet_id",
			"default_node_pool.0.snapshot_id",
			"default_node_pool.0.ultra_ssd_enabled",
			"default_node_pool.0.vnet_subnet_id",
			"default_node_pool.0.vm_size",
			"default_node_pool.0.zones",
		}

		// if the default node pool name has changed, it means the initial attempt at resizing failed
		cycleNodePool := d.HasChanges(cycleNodePoolProperties...)
		// os_sku could only be updated if the current and new os_sku are either Ubuntu or AzureLinux
		if d.HasChange("default_node_pool.0.os_sku") {
			oldOsSkuRaw, newOsSkuRaw := d.GetChange("default_node_pool.0.os_sku")
			oldOsSku := oldOsSkuRaw.(string)
			newOsSku := newOsSkuRaw.(string)
			if !strings.HasPrefix(oldOsSku, string(managedclusters.OSSKUUbuntu)) && !strings.HasPrefix(oldOsSku, string(managedclusters.OSSKUAzureLinux)) {
				cycleNodePool = true
			}
			if !strings.HasPrefix(newOsSku, string(managedclusters.OSSKUUbuntu)) && !strings.HasPrefix(newOsSku, string(managedclusters.OSSKUAzureLinux)) {
				cycleNodePool = true
			}
		}
		if cycleNodePool {
			log.Printf("[DEBUG] Cycling Default Node Pool..")
			// to provide a seamless updating experience for the vm size of the default node pool we need to cycle the default
			// node pool by provisioning a temporary system node pool, tearing down the former default node pool and then
			// bringing up the new one.

			if v := d.Get("default_node_pool.0.temporary_name_for_rotation").(string); v == "" {
				return fmt.Errorf("`temporary_name_for_rotation` must be specified when updating any of the following properties %q", cycleNodePoolProperties)
			}

			temporaryNodePoolName := d.Get("default_node_pool.0.temporary_name_for_rotation").(string)
			tempNodePoolId := agentpools.NewAgentPoolID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, temporaryNodePoolName)

			tempExisting, err := nodePoolsClient.Get(ctx, tempNodePoolId)
			if !response.WasNotFound(tempExisting.HttpResponse) && err != nil {
				return fmt.Errorf("checking for existing temporary %s: %+v", tempNodePoolId, err)
			}

			defaultExisting, err := nodePoolsClient.Get(ctx, defaultNodePoolId)
			if !response.WasNotFound(defaultExisting.HttpResponse) && err != nil {
				return fmt.Errorf("checking for existing default %s: %+v", defaultNodePoolId, err)
			}

			tempAgentProfile := agentProfile
			tempAgentProfile.Name = &temporaryNodePoolName

			if tempAgentProfile.Properties != nil {
				tempAgentProfile.Properties.NodeImageVersion = nil
			}
			if agentProfile.Properties != nil {
				agentProfile.Properties.NodeImageVersion = nil
			}

			// if the temp node pool already exists due to a previous failure, don't bother spinning it up
			if tempExisting.Model == nil {
				if err := helpers.RetryNodePoolCreation(ctx, nodePoolsClient, tempNodePoolId, tempAgentProfile); err != nil {
					return fmt.Errorf("creating temporary %s: %+v", tempNodePoolId, err)
				}
			}

			// delete the old default node pool if it exists
			if defaultExisting.Model != nil {
				if err := nodePoolsClient.DeleteThenPoll(ctx, defaultNodePoolId, agentpools.DefaultDeleteOperationOptions()); err != nil {
					return fmt.Errorf("deleting default %s: %+v", defaultNodePoolId, err)
				}
			}

			// create the default node pool with the new vm size
			if err := helpers.RetryNodePoolCreation(ctx, nodePoolsClient, defaultNodePoolId, agentProfile); err != nil {
				// if creation of the default node pool fails we automatically fall back to the temporary node pool
				// in func findDefaultNodePool
				log.Printf("[DEBUG] Creation of resized default node pool failed")
				return fmt.Errorf("creating default %s: %+v", defaultNodePoolId, err)
			}

			if err := nodePoolsClient.DeleteThenPoll(ctx, tempNodePoolId, agentpools.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting temporary %s: %+v", tempNodePoolId, err)
			}

			log.Printf("[DEBUG] Cycled Default Node Pool..")
		} else {
			if err := nodePoolsClient.CreateOrUpdateThenPoll(ctx, defaultNodePoolId, agentProfile, agentpools.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating Default Node Pool %s %+v", defaultNodePoolId, err)
			}
		}
	}

	if d.HasChange("maintenance_window") {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		maintenanceWindowProperties := expandKubernetesClusterMaintenanceConfigurationDefault(d.Get("maintenance_window").([]any))
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
		if maintenanceWindowProperties != nil {
			parameters := maintenanceconfigurations.MaintenanceConfiguration{
				Properties: maintenanceWindowProperties,
			}
			if _, err := client.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
				return fmt.Errorf("creating/updating Maintenance Configuration for Managed Kubernetes Cluster (%q): %+v", id, err)
			}
		} else {
			if _, err := client.Delete(ctx, maintenanceId); err != nil {
				return fmt.Errorf("deleting Maintenance Configuration for %s: %+v", id, err)
			}
		}
	}

	if d.HasChange("maintenance_window_auto_upgrade") {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
		existing, err := client.Get(ctx, maintenanceId)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("retrieving Auto Upgrade Schedule Maintenance Configuration for %s: %+v", id, err)
		}
		var existingProperties *maintenanceconfigurations.MaintenanceConfigurationProperties
		if existing.Model != nil {
			existingProperties = existing.Model.Properties
		}
		maintenanceWindowProperties := expandKubernetesClusterMaintenanceConfigurationForUpdate(d.Get("maintenance_window_auto_upgrade").([]any), existingProperties)
		if maintenanceWindowProperties != nil {
			parameters := maintenanceconfigurations.MaintenanceConfiguration{
				Properties: maintenanceWindowProperties,
			}
			if _, err := client.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
				return fmt.Errorf("creating/updating Auto Upgrade Schedule Maintenance Configuration for %s: %+v", id, err)
			}
		} else {
			if _, err := client.Delete(ctx, maintenanceId); err != nil {
				return fmt.Errorf("deleting Auto Upgrade Schedule Maintenance Configuration for %s: %+v", id, err)
			}
		}
	}

	if d.HasChange("maintenance_window_node_os") {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
		existing, err := client.Get(ctx, maintenanceId)
		if err != nil && !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("retrieving Node OS Upgrade Schedule Maintenance Configuration for %s: %+v", id, err)
		}
		var existingProperties *maintenanceconfigurations.MaintenanceConfigurationProperties
		if existing.Model != nil {
			existingProperties = existing.Model.Properties
		}
		maintenanceWindowProperties := expandKubernetesClusterMaintenanceConfigurationForUpdate(d.Get("maintenance_window_node_os").([]any), existingProperties)
		if maintenanceWindowProperties != nil {
			parameters := maintenanceconfigurations.MaintenanceConfiguration{
				Properties: maintenanceWindowProperties,
			}
			if _, err := client.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
				return fmt.Errorf("creating/updating Node OS Upgrade Schedule Maintenance Configuration for %s: %+v", id, err)
			}
		} else {
			if _, err := client.Delete(ctx, maintenanceId); err != nil {
				return fmt.Errorf("deleting Node OS Upgrade Schedule Maintenance Configuration for %s: %+v", id, err)
			}
		}
	}

	return resourceKubernetesClusterRead(d, meta)
}

func expandKubernetesClusterMaintenanceConfigurationForUpdate(input []any, existing *maintenanceconfigurations.MaintenanceConfigurationProperties) *maintenanceconfigurations.MaintenanceConfigurationProperties {
	if len(input) == 0 {
		return nil
	}
	value := input[0].(map[string]any)

	var schedule maintenanceconfigurations.Schedule

	if value["frequency"] == "Daily" {
		schedule = maintenanceconfigurations.Schedule{
			Daily: &maintenanceconfigurations.DailySchedule{
				IntervalDays: int64(value["interval"].(int)),
			},
		}
	}
	if value["frequency"] == "Weekly" {
		schedule = maintenanceconfigurations.Schedule{
			Weekly: &maintenanceconfigurations.WeeklySchedule{
				IntervalWeeks: int64(value["interval"].(int)),
				DayOfWeek:     maintenanceconfigurations.WeekDay(value["day_of_week"].(string)),
			},
		}
	}
	if value["frequency"] == "AbsoluteMonthly" {
		schedule = maintenanceconfigurations.Schedule{
			AbsoluteMonthly: &maintenanceconfigurations.AbsoluteMonthlySchedule{
				DayOfMonth:     int64(value["day_of_month"].(int)),
				IntervalMonths: int64(value["interval"].(int)),
			},
		}
	}
	if value["frequency"] == "RelativeMonthly" {
		schedule = maintenanceconfigurations.Schedule{
			RelativeMonthly: &maintenanceconfigurations.RelativeMonthlySchedule{
				DayOfWeek:      maintenanceconfigurations.WeekDay(value["day_of_week"].(string)),
				WeekIndex:      maintenanceconfigurations.Type(value["week_index"].(string)),
				IntervalMonths: int64(value["interval"].(int)),
			},
		}
	}

	output := &maintenanceconfigurations.MaintenanceConfigurationProperties{
		MaintenanceWindow: &maintenanceconfigurations.MaintenanceWindow{
			StartTime:       value["start_time"].(string),
			UtcOffset:       pointer.To(value["utc_offset"].(string)),
			NotAllowedDates: expandKubernetesClusterMaintenanceConfigurationDateSpans(value["not_allowed"].(*pluginsdk.Set).List()),
			Schedule:        schedule,
		},
	}

	if startDateRaw := value["start_date"]; startDateRaw != nil && startDateRaw.(string) != "" {
		startDate, _ := time.Parse(time.RFC3339, startDateRaw.(string))
		startDateStr := startDate.Format("2006-01-02")
		// start_date is an Optional+Computed property, the default value returned by the API could be invalid during update, so we only set it if it's different from the existing value
		if existing == nil || existing.MaintenanceWindow.StartDate == nil || *existing.MaintenanceWindow.StartDate != startDateStr {
			output.MaintenanceWindow.StartDate = pointer.To(startDateStr)
		}
	}

	if duration := value["duration"]; duration != nil && duration.(int) != 0 {
		output.MaintenanceWindow.DurationHours = int64(duration.(int))
	}

	return output
}
