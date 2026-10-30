// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/managedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceKubernetesClusterRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKubernetesClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	credentials, err := client.ListClusterUserCredentials(ctx, *id, managedclusters.ListClusterUserCredentialsOperationOptions{})
	if err != nil {
		return fmt.Errorf("retrieving User Credentials for %s: %+v", id, err)
	}
	if credentials.Model == nil {
		return fmt.Errorf("retrieving User Credentials for %s: payload is empty", id)
	}

	d.Set("name", id.ManagedClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("edge_zone", flattenEdgeZone(model.ExtendedLocation))
		d.Set("location", location.Normalize(model.Location))

		skuTier := string(managedclusters.ManagedClusterSKUTierFree)
		if model.Sku != nil && model.Sku.Tier != nil && *model.Sku.Tier != "" {
			skuTier = string(*model.Sku.Tier)
		}
		d.Set("sku_tier", skuTier)

		if props := model.Properties; props != nil {
			d.Set("dns_prefix", props.DnsPrefix)
			d.Set("dns_prefix_private_cluster", props.FqdnSubdomain)
			d.Set("fqdn", props.Fqdn)
			d.Set("private_fqdn", props.PrivateFQDN)
			d.Set("portal_fqdn", props.AzurePortalFQDN)
			d.Set("disk_encryption_set_id", props.DiskEncryptionSetID)
			d.Set("kubernetes_version", props.KubernetesVersion)
			d.Set("current_kubernetes_version", props.CurrentKubernetesVersion)
			d.Set("local_account_disabled", props.DisableLocalAccounts)

			nodeResourceGroup := ""
			if v := props.NodeResourceGroup; v != nil {
				nodeResourceGroup = *props.NodeResourceGroup
			}
			d.Set("node_resource_group", nodeResourceGroup)

			nodeResourceGroupId := commonids.NewResourceGroupID(id.SubscriptionId, nodeResourceGroup)
			d.Set("node_resource_group_id", nodeResourceGroupId.ID())

			upgradeChannel := ""
			nodeOSUpgradeChannel := ""
			if profile := props.AutoUpgradeProfile; profile != nil {
				if profile.UpgradeChannel != nil && *profile.UpgradeChannel != managedclusters.UpgradeChannelNone {
					upgradeChannel = string(*profile.UpgradeChannel)
				}
				if profile.NodeOSUpgradeChannel != nil {
					nodeOSUpgradeChannel = string(*profile.NodeOSUpgradeChannel)
				}
			}

			d.Set("automatic_upgrade_channel", upgradeChannel)
			d.Set("node_os_upgrade_channel", nodeOSUpgradeChannel)

			customCaTrustCertList := make([]any, 0)
			if props.SecurityProfile != nil {
				customCaTrustCertList = pluginsdk.FlattenSlice(props.SecurityProfile.CustomCATrustCertificates)
			}
			if err := d.Set("custom_ca_trust_certificates_base64", customCaTrustCertList); err != nil {
				return fmt.Errorf("setting `custom_ca_trust_certificates_base64`: %+v", err)
			}

			enablePrivateCluster := false
			enablePrivateClusterPublicFQDN := false
			runCommandEnabled := true
			privateDnsZoneId := ""

			if err := d.Set("api_server_access_profile", flattenKubernetesClusterAPIAccessProfile(props.ApiServerAccessProfile)); err != nil {
				return fmt.Errorf("setting `api_server_access_profile`: %+v", err)
			}
			if accessProfile := props.ApiServerAccessProfile; accessProfile != nil {
				if accessProfile.EnablePrivateCluster != nil {
					enablePrivateCluster = *accessProfile.EnablePrivateCluster
				}
				if accessProfile.EnablePrivateClusterPublicFQDN != nil {
					enablePrivateClusterPublicFQDN = *accessProfile.EnablePrivateClusterPublicFQDN
				}
				if accessProfile.DisableRunCommand != nil {
					runCommandEnabled = !*accessProfile.DisableRunCommand
				}
				switch {
				case accessProfile.PrivateDNSZone != nil && strings.EqualFold("System", *accessProfile.PrivateDNSZone):
					privateDnsZoneId = "System"
				case accessProfile.PrivateDNSZone != nil && strings.EqualFold("None", *accessProfile.PrivateDNSZone):
					privateDnsZoneId = "None"
				default:
					privateDnsZoneId = pointer.From(accessProfile.PrivateDNSZone)
				}
			}
			d.Set("private_dns_zone_id", privateDnsZoneId)
			d.Set("private_cluster_enabled", enablePrivateCluster)
			d.Set("private_cluster_public_fqdn_enabled", enablePrivateClusterPublicFQDN)
			d.Set("run_command_enabled", runCommandEnabled)

			if props.AddonProfiles != nil {
				addOns := flattenKubernetesAddOns(*props.AddonProfiles)
				d.Set("aci_connector_linux", addOns["aci_connector_linux"])
				d.Set("azure_policy_enabled", addOns["azure_policy_enabled"].(bool))
				d.Set("confidential_computing", addOns["confidential_computing"])
				d.Set("http_application_routing_enabled", addOns["http_application_routing_enabled"].(bool))
				d.Set("http_application_routing_zone_name", addOns["http_application_routing_zone_name"])
				d.Set("oms_agent", addOns["oms_agent"])
				d.Set("ingress_application_gateway", addOns["ingress_application_gateway"])
				d.Set("open_service_mesh_enabled", addOns["open_service_mesh_enabled"].(bool))
				d.Set("key_vault_secrets_provider", addOns["key_vault_secrets_provider"])
			}
			autoScalerProfile, err := flattenKubernetesClusterAutoScalerProfile(props.AutoScalerProfile)
			if err != nil {
				return err
			}
			if err := d.Set("auto_scaler_profile", autoScalerProfile); err != nil {
				return fmt.Errorf("setting `auto_scaler_profile`: %+v", err)
			}

			if err := d.Set("monitor_metrics", flattenKubernetesClusterAzureMonitorProfile(props.AzureMonitorProfile)); err != nil {
				return fmt.Errorf("setting `monitor_metrics`: %+v", err)
			}

			if err := d.Set("service_mesh_profile", flattenKubernetesClusterAzureServiceMeshProfile(props.ServiceMeshProfile)); err != nil {
				return fmt.Errorf("setting `service_mesh_profile`: %+v", err)
			}

			flattenedDefaultNodePool, err := helpers.FlattenDefaultNodePool(props.AgentPoolProfiles, d)
			if err != nil {
				return fmt.Errorf("flattening `default_node_pool`: %+v", err)
			}
			if err := d.Set("default_node_pool", flattenedDefaultNodePool); err != nil {
				return fmt.Errorf("setting `default_node_pool`: %+v", err)
			}

			kubeletIdentity := []any{}
			if identityProfile := props.IdentityProfile; identityProfile != nil {
				kubeletIdentity, err = flattenKubernetesClusterIdentityProfile(*props.IdentityProfile)
				if err != nil {
					return err
				}
			}

			if err := d.Set("kubelet_identity", kubeletIdentity); err != nil {
				return fmt.Errorf("setting `kubelet_identity`: %+v", err)
			}

			if err := d.Set("linux_profile", flattenKubernetesClusterLinuxProfile(props.LinuxProfile)); err != nil {
				return fmt.Errorf("setting `linux_profile`: %+v", err)
			}

			if err := d.Set("network_profile", flattenKubernetesClusterNetworkProfile(props.NetworkProfile)); err != nil {
				return fmt.Errorf("setting `network_profile`: %+v", err)
			}

			if err := d.Set("cost_analysis_enabled", flattenKubernetesClusterMetricsProfile(props.MetricsProfile)); err != nil {
				return fmt.Errorf("setting `cost_analysis_enabled`: %+v", err)
			}

			rbacEnabled := true
			if props.EnableRBAC != nil {
				rbacEnabled = *props.EnableRBAC
			}
			d.Set("role_based_access_control_enabled", rbacEnabled)

			if err := d.Set("azure_active_directory_role_based_access_control", flattenKubernetesClusterAzureActiveDirectoryRoleBasedAccessControl(props)); err != nil {
				return fmt.Errorf("setting `azure_active_directory_role_based_access_control`: %+v", err)
			}

			if err := d.Set("service_principal", flattenAzureRmKubernetesClusterServicePrincipalProfile(props.ServicePrincipalProfile, d)); err != nil {
				return fmt.Errorf("setting `service_principal`: %+v", err)
			}

			if err := d.Set("windows_profile", flattenKubernetesClusterWindowsProfile(props.WindowsProfile, d)); err != nil {
				return fmt.Errorf("setting `windows_profile`: %+v", err)
			}

			if err := d.Set("upgrade_override", flattenKubernetesClusterUpgradeOverrideSetting(props.UpgradeSettings)); err != nil {
				return fmt.Errorf("setting `upgrade_override`: %+v", err)
			}

			if err := d.Set("workload_autoscaler_profile", flattenKubernetesClusterWorkloadAutoscalerProfile(props.WorkloadAutoScalerProfile)); err != nil {
				return fmt.Errorf("setting `workload_autoscaler_profile`: %+v", err)
			}

			if err := d.Set("node_provisioning_profile", flattenKubernetesClusterNodeProvisioningProfile(props.NodeProvisioningProfile)); err != nil {
				return fmt.Errorf("setting `node_provisioning_profile`: %+v", err)
			}

			aiToolchainOperatorEnabled := false
			if props.AiToolchainOperatorProfile != nil {
				aiToolchainOperatorEnabled = pointer.From(props.AiToolchainOperatorProfile.Enabled)
			}
			d.Set("ai_toolchain_operator_enabled", aiToolchainOperatorEnabled)

			if props.SecurityProfile != nil && props.SecurityProfile.ImageCleaner != nil {
				if props.SecurityProfile.ImageCleaner.Enabled != nil {
					d.Set("image_cleaner_enabled", props.SecurityProfile.ImageCleaner.Enabled)
				}
				if props.SecurityProfile.ImageCleaner.IntervalHours != nil {
					d.Set("image_cleaner_interval_hours", props.SecurityProfile.ImageCleaner.IntervalHours)
				}
			}

			if err := d.Set("http_proxy_config", flattenKubernetesClusterHttpProxyConfig(props)); err != nil {
				return fmt.Errorf("setting `http_proxy_config`: %+v", err)
			}

			oidcIssuerEnabled := false
			oidcIssuerUrl := ""
			if props.OidcIssuerProfile != nil {
				if props.OidcIssuerProfile.Enabled != nil {
					oidcIssuerEnabled = *props.OidcIssuerProfile.Enabled
				}
				if props.OidcIssuerProfile.IssuerURL != nil {
					oidcIssuerUrl = *props.OidcIssuerProfile.IssuerURL
				}
			}

			d.Set("oidc_issuer_enabled", oidcIssuerEnabled)
			d.Set("oidc_issuer_url", oidcIssuerUrl)

			if err := d.Set("microsoft_defender", flattenKubernetesClusterMicrosoftDefender(props.SecurityProfile)); err != nil {
				return fmt.Errorf("setting `microsoft_defender`: %+v", err)
			}

			if err := d.Set("web_app_routing", flattenKubernetesClusterIngressProfile(props.IngressProfile)); err != nil {
				return fmt.Errorf("setting `web_app_routing`: %+v", err)
			}

			workloadIdentity := false
			if props.SecurityProfile != nil && props.SecurityProfile.WorkloadIdentity != nil {
				workloadIdentity = *props.SecurityProfile.WorkloadIdentity.Enabled
			}
			d.Set("workload_identity_enabled", workloadIdentity)

			azureKeyVaultKms := flattenKubernetesClusterDataSourceKeyVaultKms(props.SecurityProfile)
			if len(azureKeyVaultKms) == 1 {
				// the Key Vault Resource ID is not returned by the API in all cases, so the configured value is kept
				azureKeyVaultKms[0].(map[string]any)["key_vault_id"] = d.Get("key_management_service.0.key_vault_id").(string)
			}
			if err := d.Set("key_management_service", azureKeyVaultKms); err != nil {
				return fmt.Errorf("setting `key_management_service`: %+v", err)
			}

			// adminProfile is only available for RBAC enabled clusters with AAD and local account is not disabled
			var adminKubeConfigRaw *string
			adminKubeConfig := make([]any, 0)
			if props.AadProfile != nil && (props.DisableLocalAccounts == nil || !*props.DisableLocalAccounts) {
				adminCredentials, err := client.ListClusterAdminCredentials(ctx, *id, managedclusters.ListClusterAdminCredentialsOperationOptions{})
				if err != nil {
					return fmt.Errorf("retrieving Admin Credentials for %s: %+v", id, err)
				}
				adminKubeConfigRaw, adminKubeConfig = flattenKubernetesClusterCredentials(adminCredentials.Model, "clusterAdmin")
			}

			d.Set("kube_admin_config_raw", adminKubeConfigRaw)
			if err := d.Set("kube_admin_config", adminKubeConfig); err != nil {
				return fmt.Errorf("setting `kube_admin_config`: %+v", err)
			}

			d.Set("support_plan", pointer.From(props.SupportPlan))

			bootstrapProfile, err := flattenBootstrapProfile(props.BootstrapProfile)
			if err != nil {
				return fmt.Errorf("flattening `bootstrap_profile`: %+v", err)
			}
			if err := d.Set("bootstrap_profile", bootstrapProfile); err != nil {
				return fmt.Errorf("setting `bootstrap_profile`: %+v", err)
			}
		}

		identity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}
		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		kubeConfigRaw, kubeConfig := flattenKubernetesClusterCredentials(credentials.Model, "clusterUser")
		d.Set("kube_config_raw", kubeConfigRaw)
		if err := d.Set("kube_config", kubeConfig); err != nil {
			return fmt.Errorf("setting `kube_config`: %+v", err)
		}

		var maintenanceWindow any
		maintenanceConfigurationsClient := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
		configResp, _ := maintenanceConfigurationsClient.Get(ctx, maintenanceId)
		if configurationBody := configResp.Model; configurationBody != nil && configurationBody.Properties != nil {
			maintenanceWindow = flattenKubernetesClusterMaintenanceConfigurationDefault(configurationBody.Properties)
		}
		d.Set("maintenance_window", maintenanceWindow)

		var maintenanceWindowAutoUpgrade any
		maintenanceId = maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
		configResp, _ = maintenanceConfigurationsClient.Get(ctx, maintenanceId)
		if configurationBody := configResp.Model; configurationBody != nil && configurationBody.Properties != nil && configurationBody.Properties.MaintenanceWindow != nil {
			maintenanceWindowAutoUpgrade = flattenKubernetesClusterMaintenanceConfiguration(configurationBody.Properties.MaintenanceWindow)
		}
		d.Set("maintenance_window_auto_upgrade", maintenanceWindowAutoUpgrade)

		var maintenanceWindowNodeOS any
		maintenanceId = maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
		configResp, _ = maintenanceConfigurationsClient.Get(ctx, maintenanceId)
		if configurationBody := configResp.Model; configurationBody != nil && configurationBody.Properties != nil && configurationBody.Properties.MaintenanceWindow != nil {
			maintenanceWindowNodeOS = flattenKubernetesClusterMaintenanceConfiguration(configurationBody.Properties.MaintenanceWindow)
		}
		d.Set("maintenance_window_node_os", maintenanceWindowNodeOS)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func idsToResourceReferences(set any) *[]managedclusters.ResourceReference {
	if set == nil {
		return nil
	}

	s := set.(*pluginsdk.Set)
	results := make([]managedclusters.ResourceReference, 0)

	for _, element := range s.List() {
		results = append(results, managedclusters.ResourceReference{Id: pointer.To(element.(string))})
	}

	if len(results) > 0 {
		return &results
	}

	return nil
}

func resourceReferencesToIds(refs *[]managedclusters.ResourceReference) []string {
	if refs == nil {
		return nil
	}

	ids := make([]string, 0)

	for _, ref := range *refs {
		if ref.Id != nil {
			ids = append(ids, *ref.Id)
		}
	}

	if len(ids) > 0 {
		return ids
	}

	return nil
}
