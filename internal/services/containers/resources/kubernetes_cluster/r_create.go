// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/edgezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKubernetesClusterCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Containers.KubernetesClustersClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	env := meta.(*clients.Client).Containers.Environment
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewKubernetesClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kubernetes_cluster", id.ID())
		}
	}

	if err := helpers.ValidateKubernetesCluster(d, nil, id.ResourceGroupName, id.ManagedClusterName); err != nil {
		return err
	}

	location := location.Normalize(d.Get("location").(string))
	dnsPrefix := d.Get("dns_prefix").(string)
	kubernetesVersion := d.Get("kubernetes_version").(string)

	linuxProfileRaw := d.Get("linux_profile").([]any)
	linuxProfile := expandKubernetesClusterLinuxProfile(linuxProfileRaw)

	// NOTE: we /could/ validate the default node pool version here - but since the entire cluster deployment
	// will fail here this should be fine to omit for the Create
	agentProfiles, err := helpers.ExpandDefaultNodePool(d)
	if err != nil {
		return fmt.Errorf("expanding `default_node_pool`: %+v", err)
	}

	// the AKS API will create the default node pool with the same version as the control plane regardless of what is
	// supplied by the user which will result in a diff in some cases, so if versions have been supplied check that they
	// are identical
	agentProfile := helpers.ConvertDefaultNodePoolToAgentPool(agentProfiles)
	if prop := agentProfile.Properties; prop != nil {
		if nodePoolVersion := prop.CurrentOrchestratorVersion; nodePoolVersion != nil {
			if kubernetesVersion != "" && kubernetesVersion != *nodePoolVersion {
				return fmt.Errorf("version mismatch between the control plane running %s and default node pool running %s, they must use the same kubernetes versions", kubernetesVersion, *nodePoolVersion)
			}
		}
	}

	var addonProfiles *map[string]managedclusters.ManagedClusterAddonProfile
	addOns := collectKubernetesAddons(d)
	addonProfiles, err = expandKubernetesAddOns(d, addOns, env)
	if err != nil {
		return err
	}

	networkProfileRaw := d.Get("network_profile").([]any)
	networkProfile, err := expandKubernetesClusterNetworkProfile(networkProfileRaw, d)
	if err != nil {
		return err
	}

	metricsProfile, err := expandKubernetesClusterMetricsProfile(d.Get("cost_analysis_enabled").(bool), d.Get("sku_tier").(string))
	if err != nil {
		return err
	}

	var azureADProfile *managedclusters.ManagedClusterAADProfile
	if v, ok := d.GetOk("azure_active_directory_role_based_access_control"); ok {
		azureADProfile = expandKubernetesClusterAzureActiveDirectoryRoleBasedAccessControl(v.([]any))
	}

	t := d.Get("tags").(map[string]any)

	windowsProfileRaw := d.Get("windows_profile").([]any)
	windowsProfile := expandKubernetesClusterWindowsProfile(windowsProfileRaw)

	workloadAutoscalerProfileRaw := d.Get("workload_autoscaler_profile").([]any)
	workloadAutoscalerProfile := expandKubernetesClusterWorkloadAutoscalerProfile(workloadAutoscalerProfileRaw, d)

	nodeProvisioningProfileRaw := d.Get("node_provisioning_profile").([]any)
	nodeProvisioningProfile := expandKubernetesClusterNodeProvisioningProfile(nodeProvisioningProfileRaw)

	apiAccessProfile := expandKubernetesClusterAPIAccessProfile(d)
	if !(*apiAccessProfile.EnablePrivateCluster) && dnsPrefix == "" {
		return fmt.Errorf("`dns_prefix` should be set if it is not a private cluster")
	}

	nodeResourceGroup := d.Get("node_resource_group").(string)

	autoScalerProfileRaw := d.Get("auto_scaler_profile").([]any)
	autoScalerProfile := expandKubernetesClusterAutoScalerProfile(autoScalerProfileRaw)

	azureMonitorKubernetesMetricsRaw := d.Get("monitor_metrics").([]any)
	azureMonitorProfile := expandKubernetesClusterAzureMonitorProfile(azureMonitorKubernetesMetricsRaw)

	bootstrapProfileRaw := d.Get("bootstrap_profile").([]any)
	bootstrapProfile := expandBootstrapProfile(bootstrapProfileRaw)

	httpProxyConfigRaw := d.Get("http_proxy_config").([]any)
	httpProxyConfig := expandKubernetesClusterHttpProxyConfig(httpProxyConfigRaw)

	enableOidcIssuer := d.Get("oidc_issuer_enabled").(bool)
	oidcIssuerProfile := expandKubernetesClusterOidcIssuerProfile(enableOidcIssuer)

	storageProfileRaw := d.Get("storage_profile").([]any)
	storageProfile := expandStorageProfile(storageProfileRaw)

	upgradeOverrideSettingRaw := d.Get("upgrade_override").([]any)
	upgradeOverrideSetting := expandKubernetesClusterUpgradeOverrideSetting(upgradeOverrideSettingRaw)

	// assemble securityProfile (Defender, WorkloadIdentity, ImageCleaner, AzureKeyVaultKms)
	securityProfile := &managedclusters.ManagedClusterSecurityProfile{}

	microsoftDefenderRaw := d.Get("microsoft_defender").([]any)
	securityProfile.Defender = expandKubernetesClusterMicrosoftDefender(d, microsoftDefenderRaw)

	workloadIdentity := false
	if v, ok := d.GetOk("workload_identity_enabled"); ok {
		workloadIdentity = v.(bool)

		if workloadIdentity && !enableOidcIssuer {
			return fmt.Errorf("`oidc_issuer_enabled` must be set to `true` to enable Azure AD Workload Identity")
		}

		securityProfile.WorkloadIdentity = &managedclusters.ManagedClusterSecurityProfileWorkloadIdentity{
			Enabled: &workloadIdentity,
		}
	}

	if d.Get("image_cleaner_enabled").(bool) {
		securityProfile.ImageCleaner = &managedclusters.ManagedClusterSecurityProfileImageCleaner{
			Enabled:       pointer.To(d.Get("image_cleaner_enabled").(bool)),
			IntervalHours: pointer.To(int64(d.Get("image_cleaner_interval_hours").(int))),
		}
	}

	azureKeyVaultKmsRaw := d.Get("key_management_service").([]any)
	securityProfile.AzureKeyVaultKms, err = expandKubernetesClusterAzureKeyVaultKms(ctx, keyVaultsClient, id.SubscriptionId, d, azureKeyVaultKmsRaw)
	if err != nil {
		return err
	}

	autoUpgradeProfile := &managedclusters.ManagedClusterAutoUpgradeProfile{}

	autoChannelUpgrade := d.Get("automatic_upgrade_channel").(string)
	nodeOsChannelUpgrade := d.Get("node_os_upgrade_channel").(string)

	// this check needs to be separate and gated since node_os_channel_upgrade is a preview feature
	if nodeOsChannelUpgrade != "" && autoChannelUpgrade != "" {
		if autoChannelUpgrade == string(managedclusters.UpgradeChannelNodeNegativeimage) && nodeOsChannelUpgrade != string(managedclusters.NodeOSUpgradeChannelNodeImage) {
			return fmt.Errorf("`node_os_upgrade_channel` cannot be set to a value other than `NodeImage` if `automatic_upgrade_channel` is set to `node-image`")
		}
	}

	if autoChannelUpgrade != "" {
		autoUpgradeProfile.UpgradeChannel = pointer.ToEnum[managedclusters.UpgradeChannel](autoChannelUpgrade)
	} else {
		autoUpgradeProfile.UpgradeChannel = pointer.To(managedclusters.UpgradeChannelNone)
	}

	if nodeOsChannelUpgrade != "" {
		autoUpgradeProfile.NodeOSUpgradeChannel = pointer.ToEnum[managedclusters.NodeOSUpgradeChannel](nodeOsChannelUpgrade)
	}

	if customCaTrustCertListRaw := d.Get("custom_ca_trust_certificates_base64").([]any); len(customCaTrustCertListRaw) > 0 {
		securityProfile.CustomCATrustCertificates = pluginsdk.ExpandStringSlice(customCaTrustCertListRaw)
	}

	parameters := managedclusters.ManagedCluster{
		ExtendedLocation: edgezones.Expand(d.Get("edge_zone").(string)),
		Location:         location,
		Sku: &managedclusters.ManagedClusterSKU{
			Name: pointer.To(managedclusters.ManagedClusterSKUNameBase), // the only possible value at this point
			Tier: pointer.ToEnum[managedclusters.ManagedClusterSKUTier](d.Get("sku_tier").(string)),
		},
		Properties: &managedclusters.ManagedClusterProperties{
			ApiServerAccessProfile:    apiAccessProfile,
			AadProfile:                azureADProfile,
			AddonProfiles:             addonProfiles,
			AgentPoolProfiles:         agentProfiles,
			NodeProvisioningProfile:   nodeProvisioningProfile,
			AutoScalerProfile:         autoScalerProfile,
			AutoUpgradeProfile:        autoUpgradeProfile,
			AzureMonitorProfile:       azureMonitorProfile,
			DnsPrefix:                 pointer.To(dnsPrefix),
			EnableRBAC:                pointer.To(d.Get("role_based_access_control_enabled").(bool)),
			KubernetesVersion:         pointer.To(kubernetesVersion),
			BootstrapProfile:          bootstrapProfile,
			LinuxProfile:              linuxProfile,
			WindowsProfile:            windowsProfile,
			MetricsProfile:            metricsProfile,
			NetworkProfile:            networkProfile,
			NodeResourceGroup:         pointer.To(nodeResourceGroup),
			DisableLocalAccounts:      pointer.To(d.Get("local_account_disabled").(bool)),
			HTTPProxyConfig:           httpProxyConfig,
			OidcIssuerProfile:         oidcIssuerProfile,
			SecurityProfile:           securityProfile,
			StorageProfile:            storageProfile,
			UpgradeSettings:           upgradeOverrideSetting,
			WorkloadAutoScalerProfile: workloadAutoscalerProfile,
		},
		Tags: tags.Expand(t),
	}

	if d.Get("ai_toolchain_operator_enabled").(bool) {
		parameters.Properties.AiToolchainOperatorProfile = &managedclusters.ManagedClusterAIToolchainOperatorProfile{
			Enabled: pointer.To(true),
		}
	}
	managedClusterIdentityRaw := d.Get("identity").([]any)
	kubernetesClusterIdentityRaw := d.Get("kubelet_identity").([]any)
	servicePrincipalProfileRaw := d.Get("service_principal").([]any)

	if len(managedClusterIdentityRaw) == 0 && len(servicePrincipalProfileRaw) == 0 {
		return fmt.Errorf("either an `identity` or `service_principal` block must be specified for cluster authentication")
	}

	if len(managedClusterIdentityRaw) > 0 {
		expandedIdentity, err := expandKubernetesClusterManagedClusterIdentity(managedClusterIdentityRaw)
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		parameters.Identity = expandedIdentity
		parameters.Properties.ServicePrincipalProfile = &managedclusters.ManagedClusterServicePrincipalProfile{
			ClientId: "msi",
		}
	}
	if len(kubernetesClusterIdentityRaw) > 0 {
		parameters.Properties.IdentityProfile = expandKubernetesClusterIdentityProfile(kubernetesClusterIdentityRaw)
	}

	servicePrincipalSet := false
	if len(servicePrincipalProfileRaw) > 0 {
		servicePrincipalProfileVal := servicePrincipalProfileRaw[0].(map[string]any)
		parameters.Properties.ServicePrincipalProfile = &managedclusters.ManagedClusterServicePrincipalProfile{
			ClientId: servicePrincipalProfileVal["client_id"].(string),
			Secret:   pointer.To(servicePrincipalProfileVal["client_secret"].(string)),
		}
		servicePrincipalSet = true
	}

	if v, ok := d.GetOk("private_dns_zone_id"); ok {
		if (parameters.Identity == nil && !servicePrincipalSet) || (v.(string) != "System" && v.(string) != "None" && (!servicePrincipalSet && parameters.Identity.Type != identity.TypeUserAssigned)) {
			return fmt.Errorf("a user assigned identity or a service principal must be used when using a custom private dns zone")
		}
		apiAccessProfile.PrivateDNSZone = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("dns_prefix_private_cluster"); ok {
		if !(*apiAccessProfile.EnablePrivateCluster) || apiAccessProfile.PrivateDNSZone == nil || *apiAccessProfile.PrivateDNSZone == "System" || *apiAccessProfile.PrivateDNSZone == "None" {
			return fmt.Errorf("`dns_prefix_private_cluster` should only be set for private cluster with custom private dns zone")
		}
		parameters.Properties.FqdnSubdomain = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("disk_encryption_set_id"); ok && v.(string) != "" {
		parameters.Properties.DiskEncryptionSetID = pointer.To(v.(string))
	}

	if v := d.Get("support_plan").(string); v != "" {
		parameters.Properties.SupportPlan = pointer.ToEnum[managedclusters.KubernetesSupportPlan](v)
	}

	if ingressProfile := expandKubernetesClusterIngressProfile(d, d.Get("web_app_routing").([]any)); ingressProfile != nil {
		parameters.Properties.IngressProfile = ingressProfile
	}

	if serviceMeshProfile := expandKubernetesClusterServiceMeshProfile(d.Get("service_mesh_profile").([]any), &managedclusters.ServiceMeshProfile{}); serviceMeshProfile != nil {
		parameters.Properties.ServiceMeshProfile = serviceMeshProfile
	}

	if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, managedclusters.DefaultCreateOrUpdateOperationOptions(), sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	if maintenanceConfigRaw, ok := d.GetOk("maintenance_window"); ok {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		parameters := maintenanceconfigurations.MaintenanceConfiguration{
			Properties: expandKubernetesClusterMaintenanceConfigurationDefault(maintenanceConfigRaw.([]any)),
		}
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
		if _, err := client.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
			return fmt.Errorf("creating/updating default maintenance config for %s: %+v", id, err)
		}
	}

	if maintenanceConfigRaw, ok := d.GetOk("maintenance_window_auto_upgrade"); ok {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		parameters := maintenanceconfigurations.MaintenanceConfiguration{
			Properties: expandKubernetesClusterMaintenanceConfigurationForCreate(maintenanceConfigRaw.([]any)),
		}
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
		if _, err := client.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
			return fmt.Errorf("creating/updating auto upgrade schedule maintenance config for %s: %+v", id, err)
		}
	}

	if maintenanceConfigRaw, ok := d.GetOk("maintenance_window_node_os"); ok {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		parameters := maintenanceconfigurations.MaintenanceConfiguration{
			Properties: expandKubernetesClusterMaintenanceConfigurationForCreate(maintenanceConfigRaw.([]any)),
		}
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
		if _, err := client.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
			return fmt.Errorf("creating/updating node os upgrade schedule maintenance config for %s: %+v", id, err)
		}
	}

	return resourceKubernetesClusterRead(d, meta)
}

func expandKubernetesClusterMaintenanceConfigurationForCreate(input []any) *maintenanceconfigurations.MaintenanceConfigurationProperties {
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
		output.MaintenanceWindow.StartDate = pointer.To(startDate.Format("2006-01-02"))
	}

	if duration := value["duration"]; duration != nil && duration.(int) != 0 {
		output.MaintenanceWindow.DurationHours = int64(duration.(int))
	}

	return output
}
