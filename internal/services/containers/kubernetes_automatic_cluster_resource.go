// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/managedclusters"
	dnsValidate "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	keyVaultClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type KubernetesAutomaticClusterModel struct {
	Name                   string                        `tfschema:"name"`
	Location               string                        `tfschema:"location"`
	ResourceGroupName      string                        `tfschema:"resource_group_name"`
	APIServerAccessProfile []APIServerAccessProfileModel `tfschema:"api_server_access_profile"`
	// AutomaticUpgradeChannel         string                              `tfschema:"automatic_upgrade_channel"`
	AutoScalerProfile               []AutoScalerProfileModel        `tfschema:"auto_scaler_profile"`
	AzureActiveDirectoryRBAC        []AzureActiveDirectoryRBACModel `tfschema:"azure_active_directory_role_based_access_control"`
	BootstrapProfile                []BootstrapProfileModel         `tfschema:"bootstrap_profile"`
	CostAnalysisEnabled             bool                            `tfschema:"cost_analysis_enabled"`
	CustomCATrustCertificatesBase64 []string                        `tfschema:"custom_ca_trust_certificates_base64"`
	DefaultNodePool                 []DefaultNodePoolModel          `tfschema:"default_node_pool"`
	DiskEncryptionSetID             string                          `tfschema:"disk_encryption_set_id"`
	DNSPrefix                       string                          `tfschema:"dns_prefix"`
	DNSPrefixPrivateCluster         string                          `tfschema:"dns_prefix_private_cluster"`
	EdgeZone                        string                          `tfschema:"edge_zone"`
	HTTPProxyConfig                 []HTTPProxyConfigModel          `tfschema:"http_proxy_config"`
	Identity                        []IdentityModel                 `tfschema:"identity"`
	// ImageCleanerEnabled             bool                                `tfschema:"image_cleaner_enabled"`
	ImageCleanerIntervalHours    int64                               `tfschema:"image_cleaner_interval_hours"`
	KeyManagementService         []KeyManagementServiceModel         `tfschema:"key_management_service"`
	KubeletIdentity              []KubeletIdentityModel              `tfschema:"kubelet_identity"`
	KubernetesVersion            string                              `tfschema:"kubernetes_version"`
	LinuxProfile                 []LinuxProfileModel                 `tfschema:"linux_profile"`
	MaintenanceWindow            []MaintenanceWindowModel            `tfschema:"maintenance_window"`
	MaintenanceWindowAutoUpgrade []MaintenanceWindowAutoUpgradeModel `tfschema:"maintenance_window_auto_upgrade"`
	MaintenanceWindowNodeOS      []MaintenanceWindowNodeOSModel      `tfschema:"maintenance_window_node_os"`
	MicrosoftDefender            []MicrosoftDefenderModel            `tfschema:"microsoft_defender"`
	MonitorMetrics               []MonitorMetricsModel               `tfschema:"monitor_metrics"`
	NetworkProfile               []NetworkProfileModel               `tfschema:"network_profile"`
	NodeOSUpgradeChannel         string                              `tfschema:"node_os_upgrade_channel"`
	NodeProvisioningProfile      []NodeProvisioningProfileModel      `tfschema:"node_provisioning_profile"`
	NodeResourceGroup            string                              `tfschema:"node_resource_group"`
	// OIDCIssuerEnabled               bool                                `tfschema:"oidc_issuer_enabled"`
	// PrivateClusterEnabled           bool                      `tfschema:"private_cluster_enabled"`
	PrivateClusterPublicFQDNEnabled bool   `tfschema:"private_cluster_public_fqdn_enabled"`
	PrivateDNSZoneID                string `tfschema:"private_dns_zone_id"`
	// RoleBasedAccessControlEnabled   bool                      `tfschema:"role_based_access_control_enabled"`
	RunCommandEnabled  bool                      `tfschema:"run_command_enabled"`
	ServiceMeshProfile []ServiceMeshProfileModel `tfschema:"service_mesh_profile"`
	ServicePrincipal   []ServicePrincipalModel   `tfschema:"service_principal"`
	// SKUTier                         string                              `tfschema:"sku_tier"`
	// SKUName                         string                              `tfschema:"sku_name"`
	StorageProfile  []StorageProfileModel  `tfschema:"storage_profile"`
	SupportPlan     string                 `tfschema:"support_plan"`
	Tags            map[string]interface{} `tfschema:"tags"`
	UpgradeOverride []UpgradeOverrideModel `tfschema:"upgrade_override"`
	WebAppRouting   []WebAppRoutingModel   `tfschema:"web_app_routing"`
	WindowsProfile  []WindowsProfileModel  `tfschema:"windows_profile"`
	// WorkloadAutoscalerProfile  []WorkloadAutoscalerProfileModel `tfschema:"workload_autoscaler_profile"`
	AIToolchainOperatorEnabled bool `tfschema:"ai_toolchain_operator_enabled"`
	// WorkloadIdentityEnabled    bool                             `tfschema:"workload_identity_enabled"`

	// Addon fields
	ACIConnectorLinux []ACIConnectorLinuxModel `tfschema:"aci_connector_linux"`
	// AzurePolicyEnabled            bool                             `tfschema:"azure_policy_enabled"`
	ConfidentialComputing         []ConfidentialComputingModel     `tfschema:"confidential_computing"`
	HTTPApplicationRoutingEnabled bool                             `tfschema:"http_application_routing_enabled"`
	IngressApplicationGateway     []IngressApplicationGatewayModel `tfschema:"ingress_application_gateway"`
	KeyVaultSecretsProvider       []KeyVaultSecretsProviderModel   `tfschema:"key_vault_secrets_provider"`
	OpenServiceMeshEnabled        bool                             `tfschema:"open_service_mesh_enabled"`
	OMSAgent                      []OMSAgentModel                  `tfschema:"oms_agent"`

	// Computed fields
	CurrentKubernetesVersion       string            `tfschema:"current_kubernetes_version"`
	FQDN                           string            `tfschema:"fqdn"`
	HTTPApplicationRoutingZoneName string            `tfschema:"http_application_routing_zone_name"`
	KubeAdminConfig                []KubeConfigModel `tfschema:"kube_admin_config"`
	KubeAdminConfigRaw             string            `tfschema:"kube_admin_config_raw"`
	KubeConfig                     []KubeConfigModel `tfschema:"kube_config"`
	KubeConfigRaw                  string            `tfschema:"kube_config_raw"`
	NodeResourceGroupID            string            `tfschema:"node_resource_group_id"`
	OIDCIssuerURL                  string            `tfschema:"oidc_issuer_url"`
	PrivateFQDN                    string            `tfschema:"private_fqdn"`
}

type IdentityModel struct {
	Type        string   `tfschema:"type"`
	IdentityIds []string `tfschema:"identity_ids"`
	PrincipalId string   `tfschema:"principal_id"`
	TenantId    string   `tfschema:"tenant_id"`
}

type APIServerAccessProfileModel struct {
	AuthorizedIPRanges []string `tfschema:"authorized_ip_ranges"`
	// VirtualNetworkIntegrationEnabled bool     `tfschema:"virtual_network_integration_enabled"`
	SubnetID string `tfschema:"subnet_id"`
}

type AutoScalerProfileModel struct {
	BalanceSimilarNodeGroups                 bool    `tfschema:"balance_similar_node_groups"`
	DaemonsetEvictionForEmptyNodesEnabled    bool    `tfschema:"daemonset_eviction_for_empty_nodes_enabled"`
	DaemonsetEvictionForOccupiedNodesEnabled bool    `tfschema:"daemonset_eviction_for_occupied_nodes_enabled"`
	Expander                                 string  `tfschema:"expander"`
	IgnoreDaemonsetsUtilizationEnabled       bool    `tfschema:"ignore_daemonsets_utilization_enabled"`
	MaxGracefulTerminationSec                string  `tfschema:"max_graceful_termination_sec"`
	MaxNodeProvisioningTime                  string  `tfschema:"max_node_provisioning_time"`
	MaxUnreadyNodes                          int64   `tfschema:"max_unready_nodes"`
	MaxUnreadyPercentage                     float64 `tfschema:"max_unready_percentage"`
	NewPodScaleUpDelay                       string  `tfschema:"new_pod_scale_up_delay"`
	ScanInterval                             string  `tfschema:"scan_interval"`
	ScaleDownDelayAfterAdd                   string  `tfschema:"scale_down_delay_after_add"`
	ScaleDownDelayAfterDelete                string  `tfschema:"scale_down_delay_after_delete"`
	ScaleDownDelayAfterFailure               string  `tfschema:"scale_down_delay_after_failure"`
	ScaleDownUnneeded                        string  `tfschema:"scale_down_unneeded"`
	ScaleDownUnready                         string  `tfschema:"scale_down_unready"`
	ScaleDownUtilizationThreshold            string  `tfschema:"scale_down_utilization_threshold"`
	EmptyBulkDeleteMax                       string  `tfschema:"empty_bulk_delete_max"`
	SkipNodesWithLocalStorage                bool    `tfschema:"skip_nodes_with_local_storage"`
	SkipNodesWithSystemPods                  bool    `tfschema:"skip_nodes_with_system_pods"`
}

type AzureActiveDirectoryRBACModel struct {
	TenantID string `tfschema:"tenant_id"`
	// AzureRBACEnabled    bool     `tfschema:"azure_rbac_enabled"`
	AdminGroupObjectIDs []string `tfschema:"admin_group_object_ids"`
}

type BootstrapProfileModel struct {
	ArtifactSource      string `tfschema:"artifact_source"`
	ContainerRegistryID string `tfschema:"container_registry_id"`
}

type HTTPProxyConfigModel struct {
	HTTPProxy  string   `tfschema:"http_proxy"`
	HTTPSProxy string   `tfschema:"https_proxy"`
	NoProxy    []string `tfschema:"no_proxy"`
	TrustedCA  string   `tfschema:"trusted_ca"`
}

type KeyManagementServiceModel struct {
	KeyVaultKeyID         string `tfschema:"key_vault_key_id"`
	KeyVaultNetworkAccess string `tfschema:"key_vault_network_access"`
}

type KubeletIdentityModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type LinuxProfileModel struct {
	AdminUsername string        `tfschema:"admin_username"`
	SSHKey        []SSHKeyModel `tfschema:"ssh_key"`
}

type SSHKeyModel struct {
	KeyData string `tfschema:"key_data"`
}

type MaintenanceWindowModel struct {
	Allowed    []MaintenanceWindowAllowedModel    `tfschema:"allowed"`
	NotAllowed []MaintenanceWindowNotAllowedModel `tfschema:"not_allowed"`
}

type MaintenanceWindowAllowedModel struct {
	Day   string  `tfschema:"day"`
	Hours []int64 `tfschema:"hours"`
}

type MaintenanceWindowNotAllowedModel struct {
	End   string `tfschema:"end"`
	Start string `tfschema:"start"`
}

type MaintenanceWindowAutoUpgradeModel struct {
	Frequency  string                             `tfschema:"frequency"`
	Interval   int64                              `tfschema:"interval"`
	DayOfWeek  string                             `tfschema:"day_of_week"`
	Duration   int64                              `tfschema:"duration"`
	WeekIndex  string                             `tfschema:"week_index"`
	DayOfMonth int64                              `tfschema:"day_of_month"`
	StartDate  string                             `tfschema:"start_date"`
	StartTime  string                             `tfschema:"start_time"`
	UTCOffset  string                             `tfschema:"utc_offset"`
	NotAllowed []MaintenanceWindowNotAllowedModel `tfschema:"not_allowed"`
}

type MaintenanceWindowNodeOSModel struct {
	Frequency  string                             `tfschema:"frequency"`
	Interval   int64                              `tfschema:"interval"`
	DayOfWeek  string                             `tfschema:"day_of_week"`
	Duration   int64                              `tfschema:"duration"`
	WeekIndex  string                             `tfschema:"week_index"`
	DayOfMonth int64                              `tfschema:"day_of_month"`
	StartDate  string                             `tfschema:"start_date"`
	StartTime  string                             `tfschema:"start_time"`
	UTCOffset  string                             `tfschema:"utc_offset"`
	NotAllowed []MaintenanceWindowNotAllowedModel `tfschema:"not_allowed"`
}

type MicrosoftDefenderModel struct {
	LogAnalyticsWorkspaceID string `tfschema:"log_analytics_workspace_id"`
}

type MonitorMetricsModel struct {
	MonitorMetricsEnabled bool   `tfschema:"monitor_metrics_enabled"`
	AnnotationsAllowed    string `tfschema:"annotations_allowed"`
	LabelsAllowed         string `tfschema:"labels_allowed"`
}

type NetworkProfileModel struct {
	// NetworkPlugin       string                     `tfschema:"network_plugin"`
	NetworkPluginMode string `tfschema:"network_plugin_mode"`
	NetworkPolicy     string `tfschema:"network_policy"`
	// NetworkDataPlane    string                     `tfschema:"network_data_plane"`
	NetworkMode         string                     `tfschema:"network_mode"`
	DNSServiceIP        string                     `tfschema:"dns_service_ip"`
	PodCIDR             string                     `tfschema:"pod_cidr"`
	PodCIDRs            []string                   `tfschema:"pod_cidrs"`
	ServiceCIDR         string                     `tfschema:"service_cidr"`
	ServiceCIDRs        []string                   `tfschema:"service_cidrs"`
	IPVersions          []string                   `tfschema:"ip_versions"`
	OutboundType        string                     `tfschema:"outbound_type"`
	LoadBalancerSKU     string                     `tfschema:"load_balancer_sku"`
	LoadBalancerProfile []LoadBalancerProfileModel `tfschema:"load_balancer_profile"`
	NATGatewayProfile   []NATGatewayProfileModel   `tfschema:"nat_gateway_profile"`
	AdvancedNetworking  []AdvancedNetworkingModel  `tfschema:"advanced_networking"`
}

type LoadBalancerProfileModel struct {
	ManagedOutboundIPCount int64    `tfschema:"managed_outbound_ip_count"`
	OutboundIPAddressIDs   []string `tfschema:"outbound_ip_address_ids"`
	OutboundIPPrefixIDs    []string `tfschema:"outbound_ip_prefix_ids"`
	OutboundPortsAllocated int64    `tfschema:"outbound_ports_allocated"`
	IdleTimeoutInMinutes   int64    `tfschema:"idle_timeout_in_minutes"`
	BackendPoolType        string   `tfschema:"backend_pool_type"`
	EffectiveOutboundIPs   []string `tfschema:"effective_outbound_ips"`
}

type NATGatewayProfileModel struct {
	ManagedOutboundIPCount int64    `tfschema:"managed_outbound_ip_count"`
	IdleTimeoutInMinutes   int64    `tfschema:"idle_timeout_in_minutes"`
	EffectiveOutboundIPs   []string `tfschema:"effective_outbound_ips"`
}

type AdvancedNetworkingModel struct {
	ObservabilityEnabled bool `tfschema:"observability_enabled"`
	SecurityEnabled      bool `tfschema:"security_enabled"`
}

type NodeProvisioningProfileModel struct {
	// Mode             string `tfschema:"mode"`
	DefaultNodePools string `tfschema:"default_node_pools"`
}

type ServiceMeshProfileModel struct {
	Mode                          string                      `tfschema:"mode"`
	Revisions                     []string                    `tfschema:"revisions"`
	InternalIngressGatewayEnabled bool                        `tfschema:"internal_ingress_gateway_enabled"`
	ExternalIngressGatewayEnabled bool                        `tfschema:"external_ingress_gateway_enabled"`
	CertificateAuthority          []CertificateAuthorityModel `tfschema:"certificate_authority"`
}

type CertificateAuthorityModel struct {
	KeyVaultID          string `tfschema:"key_vault_id"`
	RootCertObjectName  string `tfschema:"root_cert_object_name"`
	CertObjectName      string `tfschema:"cert_object_name"`
	CertChainObjectName string `tfschema:"cert_chain_object_name"`
	KeyObjectName       string `tfschema:"key_object_name"`
}

type ServicePrincipalModel struct {
	ClientID     string `tfschema:"client_id"`
	ClientSecret string `tfschema:"client_secret"`
}

type StorageProfileModel struct {
	BlobDriverEnabled         bool `tfschema:"blob_driver_enabled"`
	DiskDriverEnabled         bool `tfschema:"disk_driver_enabled"`
	FileDriverEnabled         bool `tfschema:"file_driver_enabled"`
	SnapshotControllerEnabled bool `tfschema:"snapshot_controller_enabled"`
}

type UpgradeOverrideModel struct {
	ForceUpgradeEnabled bool   `tfschema:"force_upgrade_enabled"`
	EffectiveUntil      string `tfschema:"effective_until"`
}

type WebAppRoutingModel struct {
	DNSZoneIDs             []string                     `tfschema:"dns_zone_ids"`
	DefaultNginxController string                       `tfschema:"default_nginx_controller"`
	WebAppRoutingIdentity  []WebAppRoutingIdentityModel `tfschema:"web_app_routing_identity"`
}

type WebAppRoutingIdentityModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type WindowsProfileModel struct {
	AdminUsername string      `tfschema:"admin_username"`
	AdminPassword string      `tfschema:"admin_password"`
	License       string      `tfschema:"license"`
	GMSA          []GMSAModel `tfschema:"gmsa"`
}

type GMSAModel struct {
	DNSServer          string `tfschema:"dns_server"`
	RootDomain         string `tfschema:"root_domain"`
	GMSAProfileEnabled bool   `tfschema:"gmsa_profile_enabled"`
}

// type WorkloadAutoscalerProfileModel struct {
// 	KEDAEnabled                  bool `tfschema:"keda_enabled"`
// 	VerticalPodAutoscalerEnabled bool `tfschema:"vertical_pod_autoscaler_enabled"`
// }

type KubeConfigModel struct {
	Host                 string `tfschema:"host"`
	Username             string `tfschema:"username"`
	Password             string `tfschema:"password"`
	ClientCertificate    string `tfschema:"client_certificate"`
	ClientKey            string `tfschema:"client_key"`
	ClusterCACertificate string `tfschema:"cluster_ca_certificate"`
}

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name azurerm_kubernetes_automatic_cluster -properties "name,resource_group_name"
type KubernetesAutomaticClusterResource struct{}

var (
	_ sdk.ResourceWithUpdate   = KubernetesAutomaticClusterResource{}
	_ sdk.ResourceWithIdentity = KubernetesAutomaticClusterResource{}
)

func (r KubernetesAutomaticClusterResource) ResourceType() string {
	return "azurerm_kubernetes_automatic_cluster"
}

func (r KubernetesAutomaticClusterResource) ModelObject() interface{} {
	return &KubernetesAutomaticClusterModel{}
}

func (r KubernetesAutomaticClusterResource) Identity() resourceids.ResourceId {
	return &commonids.KubernetesClusterId{}
}

func (r KubernetesAutomaticClusterResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			// Check if this is a new resource (no old state)
			isNewResource := rd.Id() == ""

			if !isNewResource {
				// The behaviour of the API requires this, but this could be removed when https://github.com/Azure/azure-rest-api-specs/issues/27373 has been addressed
				// Check default_node_pool upgrade_settings drain_timeout_in_minutes
				if rd.HasChange("default_node_pool.0.upgrade_settings.0.drain_timeout_in_minutes") {
					old, new := rd.GetChange("default_node_pool.0.upgrade_settings.0.drain_timeout_in_minutes")
					if old.(int) != 0 && new.(int) == 0 {
						return fmt.Errorf("changing `default_node_pool.upgrade_settings.drain_timeout_in_minutes` from a non-zero value to zero requires recreation")
					}
				}

				// Check default_node_pool name changes with temporary_name_for_rotation
				if rd.HasChange("default_node_pool.0.name") {
					oldName, newName := rd.GetChange("default_node_pool.0.name")
					defaultName := rd.Get("default_node_pool.0.name")
					tempName := rd.Get("default_node_pool.0.temporary_name_for_rotation")

					// if the default node pool name has been set to temporary_name_for_rotation it means resizing failed
					// we should not try to recreate the cluster, another apply will attempt the resize again
					if oldName != "" && oldName == tempName {
						if newName != defaultName {
							return fmt.Errorf("changing `default_node_pool.name` requires recreation")
						}
					} else {
						return fmt.Errorf("changing `default_node_pool.name` requires recreation")
					}
				}

				// Migration of `identity` to `service_principal` is not allowed
				if rd.HasChange("service_principal.0.client_id") {
					old, _ := rd.GetChange("service_principal.0.client_id")
					oldStr := old.(string)
					if oldStr == "msi" || oldStr == "" {
						return fmt.Errorf("changing `service_principal.client_id` from MSI requires recreation")
					}
				}

				// Check windows_profile gmsa changes
				if rd.HasChange("windows_profile.0.gmsa") {
					old, new := rd.GetChange("windows_profile.0.gmsa")
					oldList := old.([]interface{})
					newList := new.([]interface{})
					if len(oldList) > 0 && len(newList) == 0 {
						return fmt.Errorf("removing `windows_profile.gmsa` requires recreation")
					}
				}

				if rd.HasChange("windows_profile.0.gmsa.0.dns_server") {
					old, new := rd.GetChange("windows_profile.0.gmsa.0.dns_server")
					if old.(string) != "" && new.(string) == "" {
						return fmt.Errorf("changing `windows_profile.gmsa.dns_server` from a non-empty value to empty requires recreation")
					}
				}

				if rd.HasChange("windows_profile.0.gmsa.0.root_domain") {
					old, new := rd.GetChange("windows_profile.0.gmsa.0.root_domain")
					if old.(string) != "" && new.(string) == "" {
						return fmt.Errorf("changing `windows_profile.gmsa.root_domain` from a non-empty value to empty requires recreation")
					}
				}

				// Check api_server_access_profile subnet_id changes
				if rd.HasChange("api_server_access_profile.0.subnet_id") {
					old, new := rd.GetChange("api_server_access_profile.0.subnet_id")
					if old.(string) != "" && new.(string) == "" {
						return fmt.Errorf("changing `api_server_access_profile.subnet_id` from a non-empty value to empty requires recreation")
					}
				}

				// Check network_plugin_mode changes
				if rd.HasChange("network_profile.0.network_plugin_mode") {
					_, new := rd.GetChange("network_profile.0.network_plugin_mode")
					if !strings.EqualFold(new.(string), string(managedclusters.NetworkPluginModeOverlay)) {
						return fmt.Errorf("changing `network_profile.network_plugin_mode` to a value other than overlay requires recreation")
					}
				}

				// Check network_policy changes
				if rd.HasChange("network_profile.0.network_policy") {
					old, new := rd.GetChange("network_profile.0.network_policy")
					oldStr := old.(string)
					newStr := new.(string)

					if oldStr != "" {
						// Azure supports in-place upgrade from Azure Network Policy Manager to Cilium and Calico to Cilium
						allowed := (oldStr == "azure" && newStr == "cilium") ||
							(oldStr == "calico" && newStr == "cilium")
						if !allowed {
							return fmt.Errorf("changing `network_profile.network_policy` requires recreation")
						}
					}
				}

				// Check pod_cidr changes
				if rd.HasChange("network_profile.0.pod_cidr") {
					old, _ := rd.GetChange("network_profile.0.pod_cidr")
					if old.(string) != "" {
						return fmt.Errorf("changing `network_profile.pod_cidr` requires recreation")
					}
				}

				// Check pod_cidrs changes
				if rd.HasChange("network_profile.0.pod_cidrs") {
					old, _ := rd.GetChange("network_profile.0.pod_cidrs")
					oldList := old.([]interface{})
					if len(oldList) > 0 {
						return fmt.Errorf("changing `network_profile.pod_cidrs` requires recreation")
					}
				}
			}

			// Validate outbound_type and bootstrap_profile artifact_source
			outboundType := rd.Get("network_profile.0.outbound_type").(string)
			artifactSource := rd.Get("bootstrap_profile.0.artifact_source").(string)

			if outboundType == string(managedclusters.OutboundTypeNone) && artifactSource != string(managedclusters.ArtifactSourceCache) {
				return fmt.Errorf("when `network_profile.outbound_type` is set to `none`, `bootstrap_profile.artifact_source` must be set to `Cache`")
			}

			return nil
		},
	}
}

func (r KubernetesAutomaticClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(val interface{}, key string) (warns []string, errs []error) {
		idRaw, ok := val.(string)
		if !ok {
			errs = append(errs, fmt.Errorf("expected `id` to be a string but got %+v", val))
			return
		}

		if _, err := commonids.ParseKubernetesClusterID(idRaw); err != nil {
			errs = append(errs, fmt.Errorf("parsing %q: %+v", idRaw, err))
		}
		return
	}
}

func (r KubernetesAutomaticClusterResource) Arguments() map[string]*pluginsdk.Schema {
	arguments := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"node_provisioning_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// "mode": {
					// 	Type:         pluginsdk.TypeString,
					// 	Optional:     true,
					// 	Default:      managedclusters.NodeProvisioningModeAuto,
					//	ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForNodeProvisioningMode(), false),
					//	AtLeastOneOf: []string{"node_provisioning_profile.0.mode", "node_provisioning_profile.0.default_node_pools"},
					// },
					"default_node_pools": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      managedclusters.NodeProvisioningDefaultNodePoolsAuto,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForNodeProvisioningDefaultNodePools(), false),
					},
				},
			},
		},

		"default_node_pool": SchemaDefaultAutomaticClusterNodePoolTyped(),

		"dns_prefix": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"dns_prefix", "dns_prefix_private_cluster"},
			ValidateFunc: containerValidate.KubernetesDNSPrefix,
		},

		"dns_prefix_private_cluster": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"dns_prefix", "dns_prefix_private_cluster"},
		},

		"kubernetes_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"node_resource_group": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		// "sku_tier": {
		// 	Type:     pluginsdk.TypeString,
		// 	Optional: true,
		// 	Default:  string(managedclusters.ManagedClusterSKUTierStandard),
		//	ValidateFunc: validation.StringInSlice([]string{
		//		string(managedclusters.ManagedClusterSKUTierStandard),
		//		string(managedclusters.ManagedClusterSKUTierPremium),
		//	}, false),
		// },

		// "sku_name": {
		// 	Type:     pluginsdk.TypeString,
		// 	Optional: true,
		// 	Default:  string(managedclusters.ManagedClusterSKUNameAutomatic),
		//	ValidateFunc: validation.StringInSlice([]string{
		//		string(managedclusters.ManagedClusterSKUNameAutomatic),
		//	}, false),
		// },

		"disk_encryption_set_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: computeValidate.DiskEncryptionSetID,
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"tags": commonschema.Tags(),

		// "automatic_upgrade_channel": {
		// 	Type:     pluginsdk.TypeString,
		// 	Optional: true,
		// 	ValidateFunc: validation.StringInSlice([]string{
		//		string(managedclusters.UpgradeChannelStable),
		//	}, false),
		//	Default: string(managedclusters.UpgradeChannelStable),
		// },

		"node_os_upgrade_channel": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(managedclusters.NodeOSUpgradeChannelNodeImage),
			ValidateFunc: validation.StringInSlice([]string{
				string(managedclusters.NodeOSUpgradeChannelNodeImage),
				string(managedclusters.NodeOSUpgradeChannelNone),
				string(managedclusters.NodeOSUpgradeChannelSecurityPatch),
				string(managedclusters.NodeOSUpgradeChannelUnmanaged),
			}, false),
		},

		"cost_analysis_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"custom_ca_trust_certificates_base64": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 10,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsBase64,
			},
		},

		// "image_cleaner_enabled": {
		// 	Type:     pluginsdk.TypeBool,
		// 	Optional: true,
		// 	Default:  true,
		// },

		"image_cleaner_interval_hours": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(168, 2160),
			Default:      168,
		},

		// "oidc_issuer_enabled": {
		// 	Type:     pluginsdk.TypeBool,
		// 	Optional: true,
		// 	Default:  true,
		// },

		// "private_cluster_enabled": {
		// 	Type:     pluginsdk.TypeBool,
		// 	Optional: true,
		// 	ForceNew: true,
		//	Default:  true,
		// },

		"private_cluster_public_fqdn_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"private_dns_zone_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				privatezones.ValidatePrivateDnsZoneID,
				validation.StringInSlice([]string{
					"System",
					"None",
				}, false),
			),
		},

		// "role_based_access_control_enabled": {
		// 	Type:     pluginsdk.TypeBool,
		// 	Optional: true,
		// 	Default:  true,
		//	ForceNew: true,
		// },

		"run_command_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"support_plan": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(managedclusters.KubernetesSupportPlanKubernetesOfficial),
			ValidateFunc: validation.StringInSlice([]string{
				string(managedclusters.KubernetesSupportPlanKubernetesOfficial),
				string(managedclusters.KubernetesSupportPlanAKSLongTermSupport),
			}, false),
		},

		"ai_toolchain_operator_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// "workload_identity_enabled": {
		// 	Type:     pluginsdk.TypeBool,
		// 	Optional: true,
		// 	Default:  true,
		// },

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"service_principal": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ExactlyOneOf: []string{"identity", "service_principal"},
			MaxItems:     1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: containerValidate.ClientID,
					},
					"client_secret": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"kubelet_identity": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
						RequiredWith: []string{
							"kubelet_identity.0.object_id",
							"kubelet_identity.0.user_assigned_identity_id",
							"identity.0.identity_ids",
						},
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"object_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
						RequiredWith: []string{
							"kubelet_identity.0.client_id",
							"kubelet_identity.0.user_assigned_identity_id",
							"identity.0.identity_ids",
						},
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"user_assigned_identity_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
						RequiredWith: []string{
							"kubelet_identity.0.client_id",
							"kubelet_identity.0.object_id",
							"identity.0.identity_ids",
						},
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},

		"azure_active_directory_role_based_access_control": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.Any(validation.IsUUID, validation.StringIsEmpty),
						AtLeastOneOf: []string{
							"azure_active_directory_role_based_access_control.0.tenant_id",
							"azure_active_directory_role_based_access_control.0.admin_group_object_ids",
						},
					},
					// "azure_rbac_enabled": {
					// 	Type:     pluginsdk.TypeBool,
					// 	Optional: true,
					// 	Default:  true,
					// },
					"admin_group_object_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsUUID,
						},
						AtLeastOneOf: []string{
							"azure_active_directory_role_based_access_control.0.tenant_id",
							"azure_active_directory_role_based_access_control.0.admin_group_object_ids",
						},
					},
				},
			},
		},

		"api_server_access_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authorized_ip_ranges": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.CIDR,
						},
					},
					// "virtual_network_integration_enabled": {
					// 	Type:     pluginsdk.TypeBool,
					// 	Optional: true,
					// },
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"http_proxy_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"http_proxy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"https_proxy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"no_proxy": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"trusted_ca": {
						Type:      pluginsdk.TypeString,
						Optional:  true,
						Sensitive: true,
					},
				},
			},
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// "network_plugin": {
					// 	Type:     pluginsdk.TypeString,
					// 	Optional: true,
					// 	ForceNew: true,
					//	Default:  managedclusters.NetworkPluginAzure,
					//	ValidateFunc: validation.StringInSlice([]string{
					//		string(managedclusters.NetworkPluginAzure),
					//		string(managedclusters.NetworkPluginKubenet),
					//		string(managedclusters.NetworkPluginNone),
					//	}, false),
					// },
					"network_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.NetworkModeBridge),
							string(managedclusters.NetworkModeTransparent),
						}, false),
					},
					"network_policy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.NetworkPolicyAzure),
							string(managedclusters.NetworkPolicyCilium),
						}, false),
					},
					"dns_service_ip": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validate.IPv4Address,
					},
					// "network_data_plane": {
					// 	Type:     pluginsdk.TypeString,
					// 	Optional: true,
					// 	Default:  string(managedclusters.NetworkDataplaneCilium),
					//	ValidateFunc: validation.StringInSlice(
					//		managedclusters.PossibleValuesForNetworkDataplane(),
					//		false),
					// },
					"network_plugin_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.NetworkPluginModeOverlay),
						}, false),
					},
					"pod_cidr": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validate.CIDR,
					},
					"pod_cidrs": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"service_cidr": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validate.CIDR,
					},
					"service_cidrs": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"load_balancer_sku": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(managedclusters.LoadBalancerSkuStandard),
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.LoadBalancerSkuBasic),
							string(managedclusters.LoadBalancerSkuStandard),
						}, false),
					},
					"outbound_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(managedclusters.OutboundTypeManagedNATGateway),
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.OutboundTypeLoadBalancer),
							string(managedclusters.OutboundTypeUserDefinedRouting),
							string(managedclusters.OutboundTypeManagedNATGateway),
							string(managedclusters.OutboundTypeUserAssignedNATGateway),
							string(managedclusters.OutboundTypeNone),
						}, false),
					},
					"load_balancer_profile": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						ForceNew: true,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"outbound_ports_allocated": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntBetween(0, 64000),
								},
								"idle_timeout_in_minutes": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      30,
									ValidateFunc: validation.IntBetween(4, 100),
								},
								"managed_outbound_ip_count": {
									Type:          pluginsdk.TypeInt,
									Optional:      true,
									Computed:      true,
									ValidateFunc:  validation.IntBetween(1, 100),
									ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids", "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"},
								},
								"managed_outbound_ipv6_count": {
									Type:          pluginsdk.TypeInt,
									Optional:      true,
									Computed:      true,
									ValidateFunc:  validation.IntBetween(1, 100),
									ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids", "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"},
								},
								"outbound_ip_prefix_ids": {
									Type:          pluginsdk.TypeSet,
									Optional:      true,
									ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.managed_outbound_ip_count", "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"},
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: azure.ValidateResourceID,
									},
								},
								"outbound_ip_address_ids": {
									Type:          pluginsdk.TypeSet,
									Optional:      true,
									ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.managed_outbound_ip_count", "network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids"},
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: azure.ValidateResourceID,
									},
								},
								"effective_outbound_ips": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"backend_pool_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  string(managedclusters.BackendPoolTypeNodeIPConfiguration),
									ValidateFunc: validation.StringInSlice([]string{
										string(managedclusters.BackendPoolTypeNodeIPConfiguration),
										string(managedclusters.BackendPoolTypeNodeIP),
									}, false),
								},
							},
						},
					},
					"nat_gateway_profile": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						ForceNew: true,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"idle_timeout_in_minutes": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      4,
									ValidateFunc: validation.IntBetween(4, 120),
								},
								"managed_outbound_ip_count": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Computed:     true,
									ValidateFunc: validation.IntBetween(1, 100),
								},
								"effective_outbound_ips": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"ip_versions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(managedclusters.IPFamilyIPvFour),
								string(managedclusters.IPFamilyIPvSix),
							}, false),
						},
					},
					"advanced_networking": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"observability_enabled": {
									Type:         pluginsdk.TypeBool,
									Optional:     true,
									Default:      false,
									AtLeastOneOf: []string{"network_profile.0.advanced_networking.0.observability_enabled", "network_profile.0.advanced_networking.0.security_enabled"},
								},
								"security_enabled": {
									Type:         pluginsdk.TypeBool,
									Optional:     true,
									Default:      false,
									AtLeastOneOf: []string{"network_profile.0.advanced_networking.0.observability_enabled", "network_profile.0.advanced_networking.0.security_enabled"},
								},
							},
						},
					},
				},
			},
		},

		"linux_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: containerValidate.KubernetesAdminUserName,
					},
					"ssh_key": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_data": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},

		"windows_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
					"admin_password": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringLenBetween(8, 123),
					},
					"license": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.LicenseTypeWindowsServer),
						}, false),
					},
					"gmsa": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"gmsa_profile_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  true,
								},
								"dns_server": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"root_domain": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"bootstrap_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"artifact_source": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForArtifactSource(), false),
						Default:      managedclusters.ArtifactSourceDirect,
					},
					"container_registry_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: registries.ValidateRegistryID,
					},
				},
			},
		},

		"auto_scaler_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"balance_similar_node_groups": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"daemonset_eviction_for_empty_nodes_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"daemonset_eviction_for_occupied_nodes_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"expander": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(managedclusters.ExpanderRandom),
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.ExpanderLeastNegativewaste),
							string(managedclusters.ExpanderMostNegativepods),
							string(managedclusters.ExpanderPriority),
							string(managedclusters.ExpanderRandom),
						}, false),
					},
					"ignore_daemonsets_utilization_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"max_graceful_termination_sec": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"max_node_provisioning_time": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "15m",
						ValidateFunc: containerValidate.Duration,
					},
					"max_unready_nodes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      3,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"max_unready_percentage": {
						Type:         pluginsdk.TypeFloat,
						Optional:     true,
						Default:      45,
						ValidateFunc: validation.FloatBetween(0, 100),
					},
					"new_pod_scale_up_delay": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: containerValidate.Duration,
					},
					"scan_interval": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: containerValidate.Duration,
					},
					"scale_down_delay_after_add": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: containerValidate.Duration,
					},
					"scale_down_delay_after_delete": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: containerValidate.Duration,
					},
					"scale_down_delay_after_failure": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: containerValidate.Duration,
					},
					"scale_down_unneeded": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: containerValidate.Duration,
					},
					"scale_down_unready": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: containerValidate.Duration,
					},
					"scale_down_utilization_threshold": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"empty_bulk_delete_max": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"skip_nodes_with_local_storage": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"skip_nodes_with_system_pods": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		"monitor_metrics": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"monitor_metrics_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"annotations_allowed": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"labels_allowed": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"key_management_service": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: false,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey),
					},
					"key_vault_network_access": {
						Type:         pluginsdk.TypeString,
						Default:      string(managedclusters.KeyVaultNetworkAccessTypesPublic),
						Optional:     true,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForKeyVaultNetworkAccessTypes(), false),
					},
				},
			},
		},

		"microsoft_defender": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"log_analytics_workspace_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: workspaces.ValidateWorkspaceID,
					},
				},
			},
		},

		"web_app_routing": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dns_zone_ids": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								dnsValidate.ValidateDnsZoneID,
								privatezones.ValidatePrivateDnsZoneID,
								validation.StringIsEmpty,
							),
						},
					},
					"default_nginx_controller": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  managedclusters.NginxIngressControllerTypeAnnotationControlled,
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.NginxIngressControllerTypeAnnotationControlled),
							string(managedclusters.NginxIngressControllerTypeInternal),
							string(managedclusters.NginxIngressControllerTypeExternal),
							string(managedclusters.NginxIngressControllerTypeNone),
						}, false),
					},
					"web_app_routing_identity": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"object_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"user_assigned_identity_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"maintenance_window": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed": {
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"maintenance_window.0.allowed", "maintenance_window.0.not_allowed"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"day": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(maintenanceconfigurations.WeekDaySunday),
										string(maintenanceconfigurations.WeekDayMonday),
										string(maintenanceconfigurations.WeekDayTuesday),
										string(maintenanceconfigurations.WeekDayWednesday),
										string(maintenanceconfigurations.WeekDayThursday),
										string(maintenanceconfigurations.WeekDayFriday),
										string(maintenanceconfigurations.WeekDaySaturday),
									}, false),
								},
								"hours": {
									Type:     pluginsdk.TypeSet,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeInt,
										ValidateFunc: validation.IntBetween(0, 23),
									},
								},
							},
						},
					},
					"not_allowed": {
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"maintenance_window.0.allowed", "maintenance_window.0.not_allowed"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"end": {
									Type:             pluginsdk.TypeString,
									Required:         true,
									DiffSuppressFunc: suppress.RFC3339Time,
									ValidateFunc:     validation.IsRFC3339Time,
								},
								"start": {
									Type:             pluginsdk.TypeString,
									Required:         true,
									DiffSuppressFunc: suppress.RFC3339Time,
									ValidateFunc:     validation.IsRFC3339Time,
								},
							},
						},
					},
				},
			},
		},

		"maintenance_window_auto_upgrade": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"frequency": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Daily",
							"Weekly",
							"RelativeMonthly",
							"AbsoluteMonthly",
						}, false),
					},
					"interval": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"day_of_week": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(maintenanceconfigurations.PossibleValuesForWeekDay(), false),
					},
					"duration": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(4, 24),
					},
					"week_index": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(maintenanceconfigurations.PossibleValuesForType(), false),
					},
					"day_of_month": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 31),
					},
					"start_date": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"start_time": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"utc_offset": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"not_allowed": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"end": {
									Type:             pluginsdk.TypeString,
									Required:         true,
									DiffSuppressFunc: suppress.RFC3339Time,
									ValidateFunc:     validation.IsRFC3339Time,
								},
								"start": {
									Type:             pluginsdk.TypeString,
									Required:         true,
									DiffSuppressFunc: suppress.RFC3339Time,
									ValidateFunc:     validation.IsRFC3339Time,
								},
							},
						},
					},
				},
			},
		},

		"maintenance_window_node_os": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"frequency": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Weekly",
							"RelativeMonthly",
							"AbsoluteMonthly",
							"Daily",
						}, false),
					},
					"interval": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"day_of_week": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(maintenanceconfigurations.PossibleValuesForWeekDay(), false),
					},
					"duration": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(4, 24),
					},
					"week_index": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(maintenanceconfigurations.PossibleValuesForType(), false),
					},
					"day_of_month": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 31),
					},
					"start_date": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Computed:         true,
						DiffSuppressFunc: suppress.RFC3339Time,
						ValidateFunc:     validation.IsRFC3339Time,
					},
					"start_time": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"utc_offset": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"not_allowed": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"end": {
									Type:             pluginsdk.TypeString,
									Required:         true,
									DiffSuppressFunc: suppress.RFC3339Time,
									ValidateFunc:     validation.IsRFC3339Time,
								},
								"start": {
									Type:             pluginsdk.TypeString,
									Required:         true,
									DiffSuppressFunc: suppress.RFC3339Time,
									ValidateFunc:     validation.IsRFC3339Time,
								},
							},
						},
					},
				},
			},
		},

		"service_mesh_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.ServiceMeshModeIstio),
						}, false),
					},
					"internal_ingress_gateway_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"external_ingress_gateway_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"certificate_authority": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}),
								"root_cert_object_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"cert_chain_object_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"cert_object_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"key_object_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"revisions": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						MaxItems: 2,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},

		"storage_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"blob_driver_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"disk_driver_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"file_driver_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"snapshot_controller_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		// "workload_autoscaler_profile": {
		// 	Type:     pluginsdk.TypeList,
		// 	Optional: true,
		// 	Computed: true,
		//	MaxItems: 1,
		//	Elem: &pluginsdk.Resource{
		//		Schema: map[string]*pluginsdk.Schema{
		//			"keda_enabled": {
		//				Type:     pluginsdk.TypeBool,
		//				Optional: true,
		//				Default:  true,
		//			},
		//			"vertical_pod_autoscaler_enabled": {
		//				Type:     pluginsdk.TypeBool,
		//				Optional: true,
		//				Computed: true,
		//			},
		//		},
		//	},
		// },

		"upgrade_override": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"force_upgrade_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"effective_until": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
				},
			},
		},
	}

	for k, v := range schemaKubernetesAutomaticClusterAddOnsTyped() {
		arguments[k] = v
	}

	return arguments
}

func (r KubernetesAutomaticClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"current_kubernetes_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"node_resource_group_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"oidc_issuer_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"portal_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"http_application_routing_zone_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kube_admin_config": {
			Type:      pluginsdk.TypeList,
			Computed:  true,
			Sensitive: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"username": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"password": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"client_certificate": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"client_key": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"cluster_ca_certificate": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},

		"kube_admin_config_raw": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"kube_config": {
			Type:      pluginsdk.TypeList,
			Computed:  true,
			Sensitive: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"username": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"password": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"client_certificate": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"client_key": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
					"cluster_ca_certificate": {
						Type:      pluginsdk.TypeString,
						Computed:  true,
						Sensitive: true,
					},
				},
			},
		},

		"kube_config_raw": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}

func (r KubernetesAutomaticClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesClustersClient
			keyVaultsClient := metadata.Client.KeyVault
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model KubernetesAutomaticClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := commonids.NewKubernetesClusterID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			if err := validateKubernetesAutomaticClusterTyped(&model, nil); err != nil {
				return fmt.Errorf("validating configuration: %+v", err)
			}

			location := location.Normalize(model.Location)
			kubernetesVersion := model.KubernetesVersion

			linuxProfile := expandKubernetesAutomaticClusterLinuxProfile(model.LinuxProfile)

			windowsProfile := expandKubernetesAutomaticClusterWindowsProfile(model.WindowsProfile)

			bootstrapProfile := expandKubernetesAutomaticClusterBootstrapProfile(model.BootstrapProfile)

			autoScalerProfile := expandKubernetesAutomaticClusterAutoScalerProfile(model.AutoScalerProfile)

			networkProfile, err := expandKubernetesAutomaticClusterNetworkProfile(model.NetworkProfile)
			if err != nil {
				return fmt.Errorf("expanding network profile: %+v", err)
			}

			// workloadAutoscalerProfile := expandKubernetesAutomaticClusterWorkloadAutoscalerProfile(model.WorkloadAutoscalerProfile)

			nodeProvisioningProfile := expandKubernetesAutomaticClusterNodeProvisioningProfile(model.NodeProvisioningProfile)

			azureMonitorProfile := expandKubernetesAutomaticClusterAzureMonitorProfile(model.MonitorMetrics)

			httpProxyConfig := expandKubernetesAutomaticClusterHttpProxyConfig(model.HTTPProxyConfig)

			apiAccessProfile := expandKubernetesAutomaticClusterAPIAccessProfile(model)
			// if !(*apiAccessProfile.EnablePrivateCluster) && model.DNSPrefix == "" {
			// 	return fmt.Errorf("`dns_prefix` should be set if it is not a private cluster")
			// }

			// enableOidcIssuer := model.OIDCIssuerEnabled
			// var oidcIssuerProfile *managedclusters.ManagedClusterOIDCIssuerProfile
			// if enableOidcIssuer {
			// 	oidcIssuerProfile = expandKubernetesAutomaticClusterOidcIssuerProfile(enableOidcIssuer)
			//}

			storageProfile := expandKubernetesAutomaticClusterStorageProfile(model.StorageProfile)

			upgradeOverrideSetting := expandKubernetesAutomaticClusterUpgradeOverride(model.UpgradeOverride)

			securityProfile := &managedclusters.ManagedClusterSecurityProfile{}

			securityProfile.Defender = expandKubernetesAutomaticClusterMicrosoftDefender(model.MicrosoftDefender, false)

			// workloadIdentity := model.WorkloadIdentityEnabled
			// if workloadIdentity {
			// 	if !enableOidcIssuer {
			// 		return fmt.Errorf("`oidc_issuer_enabled` must be set to `true` to enable Azure AD Workload Identity")
			//	}
			//	securityProfile.WorkloadIdentity = &managedclusters.ManagedClusterSecurityProfileWorkloadIdentity{
			//		Enabled: &workloadIdentity,
			//	}
			//}

			// if model.ImageCleanerEnabled {
			securityProfile.ImageCleaner = &managedclusters.ManagedClusterSecurityProfileImageCleaner{
				Enabled:       pointer.To(true),
				IntervalHours: pointer.To(model.ImageCleanerIntervalHours),
			}
			//}

			securityProfile.AzureKeyVaultKms, err = expandKubernetesAutomaticClusterKeyManagementService(model.KeyManagementService, ctx, keyVaultsClient, subscriptionId)
			if err != nil {
				return err
			}

			if len(model.CustomCATrustCertificatesBase64) > 0 {
				securityProfile.CustomCATrustCertificates = &model.CustomCATrustCertificatesBase64
			}

			autoUpgradeProfile := &managedclusters.ManagedClusterAutoUpgradeProfile{}
			// autoChannelUpgrade := model.AutomaticUpgradeChannel
			nodeOsChannelUpgrade := model.NodeOSUpgradeChannel
			//
			// if nodeOsChannelUpgrade != "" && autoChannelUpgrade != "" {
			// 	if autoChannelUpgrade == string(managedclusters.UpgradeChannelNodeNegativeimage) && nodeOsChannelUpgrade != string(managedclusters.NodeOSUpgradeChannelNodeImage) {
			// 		return fmt.Errorf("`node_os_upgrade_channel` cannot be set to a value other than `NodeImage` if `automatic_upgrade_channel` is set to `node-image`")
			// 	}
			// }

			// if autoChannelUpgrade != "" {
			// 	autoUpgradeProfile.UpgradeChannel = pointer.To(managedclusters.UpgradeChannel(autoChannelUpgrade))
			// } else {
			// 	autoUpgradeProfile.UpgradeChannel = pointer.To(managedclusters.UpgradeChannelNone)
			//}

			if nodeOsChannelUpgrade != "" {
				autoUpgradeProfile.NodeOSUpgradeChannel = pointer.To(managedclusters.NodeOSUpgradeChannel(nodeOsChannelUpgrade))
			}

			metricsProfile := expandKubernetesAutomaticClusterMetricsProfile(model.CostAnalysisEnabled)

			var ingressProfile *managedclusters.ManagedClusterIngressProfile
			if len(model.WebAppRouting) > 0 {
				ingressProfile = expandKubernetesAutomaticClusterWebAppRouting(model.WebAppRouting, false)
			}

			var serviceMeshProfile *managedclusters.ServiceMeshProfile
			if len(model.ServiceMeshProfile) > 0 {
				serviceMeshProfile = expandKubernetesAutomaticClusterServiceMeshProfile(model.ServiceMeshProfile, nil)
			}

			agentProfiles, err := ExpandDefaultNodePoolTyped(model.DefaultNodePool)
			if err != nil {
				return fmt.Errorf("expanding `default_node_pool`: %+v", err)
			}

			agentProfile := ConvertDefaultNodePoolToAgentPool(agentProfiles)
			if prop := agentProfile.Properties; prop != nil {
				if nodePoolVersion := prop.CurrentOrchestratorVersion; nodePoolVersion != nil {
					if model.KubernetesVersion != "" && model.KubernetesVersion != *nodePoolVersion {
						return fmt.Errorf("version mismatch between the control plane running %s and default node pool running %s, they must use the same kubernetes versions", model.KubernetesVersion, *nodePoolVersion)
					}
				}
			}

			addonProfiles, err := expandKubernetesAddOnsTyped(&model, metadata.Client.Containers.Environment)
			if err != nil {
				return fmt.Errorf("expanding addons: %+v", err)
			}

			if addonProfiles == nil {
				addonProfiles = &map[string]managedclusters.ManagedClusterAddonProfile{}
			}

			(*addonProfiles)["azurepolicy"] = managedclusters.ManagedClusterAddonProfile{
				Enabled: true,
			}

			var azureADProfile *managedclusters.ManagedClusterAADProfile
			if len(model.AzureActiveDirectoryRBAC) > 0 {
				azureADProfile = &managedclusters.ManagedClusterAADProfile{
					Managed: pointer.To(true),
					// EnableAzureRBAC:     pointer.To(model.RoleBasedAccessControlEnabled),
					AdminGroupObjectIDs: &model.AzureActiveDirectoryRBAC[0].AdminGroupObjectIDs,
				}
			}

			parameters := managedclusters.ManagedCluster{
				ExtendedLocation: expandKubernetesAutomaticClusterEdgeZone(model.EdgeZone),
				Location:         location,
				Sku: &managedclusters.ManagedClusterSKU{
					Name: pointer.To(managedclusters.ManagedClusterSKUName("automatic")),
					Tier: pointer.To(managedclusters.ManagedClusterSKUTier("standard")),
				},
				Properties: &managedclusters.ManagedClusterProperties{
					ApiServerAccessProfile:  apiAccessProfile,
					AadProfile:              azureADProfile,
					AddonProfiles:           addonProfiles,
					AgentPoolProfiles:       agentProfiles,
					NodeProvisioningProfile: nodeProvisioningProfile,
					AutoScalerProfile:       autoScalerProfile,
					AutoUpgradeProfile:      autoUpgradeProfile,
					AzureMonitorProfile:     azureMonitorProfile,
					DnsPrefix:               pointer.To(model.DNSPrefix),
					// EnableRBAC:              pointer.To(model.RoleBasedAccessControlEnabled),
					KubernetesVersion: pointer.To(kubernetesVersion),
					BootstrapProfile:  bootstrapProfile,
					LinuxProfile:      linuxProfile,
					WindowsProfile:    windowsProfile,
					MetricsProfile:    metricsProfile,
					NetworkProfile:    networkProfile,
					NodeResourceGroup: pointer.To(model.NodeResourceGroup),
					// DisableLocalAccounts:      pointer.To(requiredValues.DisableLocalAccounts),
					HTTPProxyConfig: httpProxyConfig,
					// OidcIssuerProfile:         oidcIssuerProfile,
					SecurityProfile: securityProfile,
					StorageProfile:  storageProfile,
					UpgradeSettings: upgradeOverrideSetting,
					// WorkloadAutoScalerProfile: workloadAutoscalerProfile,
					IngressProfile:     ingressProfile,
					ServiceMeshProfile: serviceMeshProfile,
				},
				Tags: tags.Expand(model.Tags),
			}

			if model.AIToolchainOperatorEnabled {
				parameters.Properties.AiToolchainOperatorProfile = &managedclusters.ManagedClusterAIToolchainOperatorProfile{
					Enabled: pointer.To(true),
				}
			}

			if len(model.Identity) > 0 {
				parameters.Identity = expandIdentityModel(model.Identity)
				parameters.Properties.ServicePrincipalProfile = &managedclusters.ManagedClusterServicePrincipalProfile{
					ClientId: "msi",
				}
			}

			if len(model.KubeletIdentity) > 0 {
				parameters.Properties.IdentityProfile = expandKubernetesAutomaticClusterIdentityProfile(model.KubeletIdentity)
			}

			servicePrincipalSet := false
			if len(model.ServicePrincipal) > 0 {
				sp := model.ServicePrincipal[0]
				parameters.Properties.ServicePrincipalProfile = &managedclusters.ManagedClusterServicePrincipalProfile{
					ClientId: sp.ClientID,
					Secret:   pointer.To(sp.ClientSecret),
				}
				servicePrincipalSet = true
			}

			if model.PrivateDNSZoneID != "" {
				if (parameters.Identity == nil && !servicePrincipalSet) || (model.PrivateDNSZoneID != "System" && model.PrivateDNSZoneID != "None" && (!servicePrincipalSet && parameters.Identity.Type != identity.TypeUserAssigned)) {
					return fmt.Errorf("a user assigned identity or a service principal must be used when using a custom private dns zone")
				}
				apiAccessProfile.PrivateDNSZone = pointer.To(model.PrivateDNSZoneID)
			}

			if model.DNSPrefixPrivateCluster != "" {
				if apiAccessProfile.PrivateDNSZone == nil || *apiAccessProfile.PrivateDNSZone == "System" || *apiAccessProfile.PrivateDNSZone == "None" {
					return fmt.Errorf("`dns_prefix_private_cluster` should only be set for private cluster with custom private dns zone")
				}
				parameters.Properties.FqdnSubdomain = pointer.To(model.DNSPrefixPrivateCluster)
			}

			if model.DiskEncryptionSetID != "" {
				parameters.Properties.DiskEncryptionSetID = pointer.To(model.DiskEncryptionSetID)
			}

			if model.SupportPlan != "" {
				parameters.Properties.SupportPlan = pointer.To(managedclusters.KubernetesSupportPlan(model.SupportPlan))
			}

			err = client.CreateOrUpdateThenPoll(ctx, id, parameters, managedclusters.DefaultCreateOrUpdateOperationOptions())
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if len(model.MaintenanceWindow) > 0 {
				maintenanceClient := metadata.Client.Containers.MaintenanceConfigurationsClient
				maintenanceParams := maintenanceconfigurations.MaintenanceConfiguration{
					Properties: expandKubernetesAutomaticClusterMaintenanceConfigurationDefault(model.MaintenanceWindow),
				}
				maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
				if _, err := maintenanceClient.CreateOrUpdate(ctx, maintenanceId, maintenanceParams); err != nil {
					return fmt.Errorf("creating/updating default maintenance config for %s: %+v", id, err)
				}
			}

			if len(model.MaintenanceWindowAutoUpgrade) > 0 {
				maintenanceClient := metadata.Client.Containers.MaintenanceConfigurationsClient
				maintenanceParams := maintenanceconfigurations.MaintenanceConfiguration{
					Properties: expandKubernetesAutomaticClusterMaintenanceConfigurationAutoUpgrade(model.MaintenanceWindowAutoUpgrade),
				}
				maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
				if _, err := maintenanceClient.CreateOrUpdate(ctx, maintenanceId, maintenanceParams); err != nil {
					return fmt.Errorf("creating/updating auto upgrade schedule maintenance config for %s: %+v", id, err)
				}
			}

			if len(model.MaintenanceWindowNodeOS) > 0 {
				maintenanceClient := metadata.Client.Containers.MaintenanceConfigurationsClient
				maintenanceParams := maintenanceconfigurations.MaintenanceConfiguration{
					Properties: expandKubernetesAutomaticClusterMaintenanceConfigurationNodeOS(model.MaintenanceWindowNodeOS),
				}
				maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
				if _, err := maintenanceClient.CreateOrUpdate(ctx, maintenanceId, maintenanceParams); err != nil {
					return fmt.Errorf("creating/updating node os upgrade schedule maintenance config for %s: %+v", id, err)
				}
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return metadata.Encode(&model)
		},
	}
}

func (r KubernetesAutomaticClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesClustersClient

			id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			return r.flatten(ctx, metadata, id, resp.Model)
		},
	}
}

func (r KubernetesAutomaticClusterResource) flatten(ctx context.Context, metadata sdk.ResourceMetaData, id *commonids.KubernetesClusterId, model *managedclusters.ManagedCluster) error {
	client := metadata.Client.Containers.KubernetesClustersClient

	credentials, err := client.ListClusterUserCredentials(ctx, *id, managedclusters.ListClusterUserCredentialsOperationOptions{})
	if err != nil {
		return fmt.Errorf("retrieving User Credentials for %s: %+v", id, err)
	}
	if credentials.Model == nil {
		return fmt.Errorf("retrieving User Credentials for %s: payload is empty", id)
	}

	var config KubernetesAutomaticClusterModel
	if err := metadata.Decode(&config); err != nil {
		return fmt.Errorf("decoding %+v", err)
	}

	state := KubernetesAutomaticClusterModel{
		Name:              id.ManagedClusterName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.EdgeZone = flattenKubernetesAutomaticClusterEdgeZone(model.ExtendedLocation)
		// Only set tags if non-empty to avoid empty map in state
		// if model.Tags != nil {
		state.Tags = tags.Flatten(model.Tags)
		// }

		// skuTier := string(managedclusters.ManagedClusterSKUTierStandard)
		// skuName := string(managedclusters.ManagedClusterSKUNameAutomatic)
		//// #TODO dont need this
		// if model.Sku != nil {
		//	if model.Sku.Tier != nil && *model.Sku.Tier != "" {
		//		skuTier = string(*model.Sku.Tier)
		//	}
		//	if model.Sku.Name != nil && *model.Sku.Name != "" {
		//		skuName = string(*model.Sku.Name)
		//	}
		// }
		// state.SKUTier = string(*model.Sku.Tier)
		// state.SKUName = string(*model.Sku.Name)

		if props := model.Properties; props != nil {
			state.DNSPrefix = pointer.From(props.DnsPrefix)
			state.DNSPrefixPrivateCluster = pointer.From(props.FqdnSubdomain)
			state.FQDN = pointer.From(props.Fqdn)
			state.PrivateFQDN = pointer.From(props.PrivateFQDN)
			state.DiskEncryptionSetID = pointer.From(props.DiskEncryptionSetID)
			state.KubernetesVersion = pointer.From(props.KubernetesVersion)
			state.CurrentKubernetesVersion = pointer.From(props.CurrentKubernetesVersion)

			state.NodeResourceGroup = pointer.From(props.NodeResourceGroup)
			if state.NodeResourceGroup != "" {
				state.NodeResourceGroupID = commonids.NewResourceGroupID(id.SubscriptionId, state.NodeResourceGroup).ID()
			}

			// upgradeChannel := ""
			// nodeOSUpgradeChannel := ""
			// if profile := props.AutoUpgradeProfile; profile != nil {
			//	if profile.UpgradeChannel != nil && *profile.UpgradeChannel != managedclusters.UpgradeChannelNone {
			//		upgradeChannel = string(*profile.UpgradeChannel)
			//	}
			//	if profile.NodeOSUpgradeChannel != nil {
			//		nodeOSUpgradeChannel = string(*profile.NodeOSUpgradeChannel)
			//	}
			//}
			// state.AutomaticUpgradeChannel = string(*props.AutoUpgradeProfile.UpgradeChannel)
			state.NodeOSUpgradeChannel = string(*props.AutoUpgradeProfile.NodeOSUpgradeChannel)

			if props.SecurityProfile != nil && props.SecurityProfile.CustomCATrustCertificates != nil {
				state.CustomCATrustCertificatesBase64 = *props.SecurityProfile.CustomCATrustCertificates
			}

			// enablePrivateCluster := false
			enablePrivateClusterPublicFQDN := false
			runCommandEnabled := true
			privateDnsZoneId := ""

			apiServerAccessProfile := flattenKubernetesAutomaticClusterAPIAccessProfile(props.ApiServerAccessProfile)
			state.APIServerAccessProfile = apiServerAccessProfile

			if accessProfile := props.ApiServerAccessProfile; accessProfile != nil {
				// if accessProfile.EnablePrivateCluster != nil {
				// 	enablePrivateCluster = *accessProfile.EnablePrivateCluster
				// }
				if accessProfile.EnablePrivateClusterPublicFQDN != nil {
					enablePrivateClusterPublicFQDN = pointer.From(accessProfile.EnablePrivateClusterPublicFQDN)
				}
				if accessProfile.DisableRunCommand != nil {
					runCommandEnabled = !pointer.From(accessProfile.DisableRunCommand)
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
			state.PrivateDNSZoneID = privateDnsZoneId
			// state.PrivateClusterEnabled = enablePrivateCluster
			state.PrivateClusterPublicFQDNEnabled = enablePrivateClusterPublicFQDN
			state.RunCommandEnabled = runCommandEnabled

			if props.AddonProfiles != nil {
				state.ACIConnectorLinux,
					state.ConfidentialComputing,
					state.HTTPApplicationRoutingEnabled,
					state.HTTPApplicationRoutingZoneName,
					state.IngressApplicationGateway,
					state.KeyVaultSecretsProvider,
					state.OMSAgent,
					state.OpenServiceMeshEnabled = flattenKubernetesAddOnsTyped(*props.AddonProfiles)
			}

			autoScalerProfile, err := flattenKubernetesAutomaticClusterAutoScalerProfile(props.AutoScalerProfile)
			if err != nil {
				return fmt.Errorf("flattening `auto_scaler_profile`: %+v", err)
			}
			state.AutoScalerProfile = autoScalerProfile
			state.MonitorMetrics = flattenKubernetesAutomaticClusterAzureMonitorProfile(props.AzureMonitorProfile)
			state.ServiceMeshProfile = flattenKubernetesAutomaticClusterServiceMeshProfile(props.ServiceMeshProfile)

			if props.AgentPoolProfiles != nil {
				defaultNodePool, err := FlattenDefaultNodePoolTyped(props.AgentPoolProfiles, &metadata)
				if err != nil {
					return fmt.Errorf("flattening `default_node_pool`: %+v", err)
				}
				state.DefaultNodePool = defaultNodePool
			}

			kubeletIdentity, err := flattenKubernetesAutomaticClusterIdentityProfile(pointer.From(props.IdentityProfile))
			if err != nil {
				return fmt.Errorf("flattening `kubelet_identity`: %+v", err)
			}
			state.KubeletIdentity = kubeletIdentity

			state.LinuxProfile = flattenKubernetesAutomaticClusterLinuxProfile(props.LinuxProfile)
			state.NetworkProfile = flattenKubernetesAutomaticClusterNetworkProfile(props.NetworkProfile)

			state.WindowsProfile = flattenKubernetesAutomaticClusterWindowsProfile(props.WindowsProfile, config)

			// state.WorkloadAutoscalerProfile = flattenKubernetesAutomaticClusterWorkloadAutoscalerProfile(props.WorkloadAutoScalerProfile)
			state.NodeProvisioningProfile = flattenKubernetesAutomaticClusterNodeProvisioningProfile(props.NodeProvisioningProfile)
			state.HTTPProxyConfig = flattenKubernetesAutomaticClusterHttpProxyConfig(props.HTTPProxyConfig)

			state.BootstrapProfile = flattenKubernetesAutomaticClusterBootstrapProfile(props.BootstrapProfile)
			state.UpgradeOverride = flattenKubernetesAutomaticClusterUpgradeOverride(props.UpgradeSettings)

			if props.StorageProfile != nil {
				state.StorageProfile = flattenKubernetesAutomaticClusterStorageProfile(props.StorageProfile)
			}

			state.WebAppRouting = flattenKubernetesAutomaticClusterWebAppRouting(props.IngressProfile)
			state.MicrosoftDefender = flattenKubernetesAutomaticClusterMicrosoftDefender(props.SecurityProfile)

			// if props.SecurityProfile != nil && props.SecurityProfile.AzureKeyVaultKms != nil {
			state.KeyManagementService = flattenKubernetesAutomaticClusterKeyManagementService(props.SecurityProfile.AzureKeyVaultKms)
			//}

			state.CostAnalysisEnabled = flattenKubernetesAutomaticClusterMetricsProfile(props.MetricsProfile)

			// rbacEnabled := true
			// if props.EnableRBAC != nil {
			// 	rbacEnabled = *props.EnableRBAC
			// }
			// state.RoleBasedAccessControlEnabled = rbacEnabled

			state.AzureActiveDirectoryRBAC = flattenKubernetesAutomaticClusterAzureActiveDirectoryRBAC(props.AadProfile)

			if props.ServicePrincipalProfile != nil &&
				props.ServicePrincipalProfile.ClientId != "" &&
				props.ServicePrincipalProfile.ClientId != "msi" {
				state.ServicePrincipal = []ServicePrincipalModel{{
					ClientID: props.ServicePrincipalProfile.ClientId,
					// ClientSecret is not returned by the API
				}}
			}

			// if props.SecurityProfile != nil && props.SecurityProfile.ImageCleaner != nil {
			// if props.SecurityProfile.ImageCleaner.Enabled != nil {
			// 	state.ImageCleanerEnabled = *props.SecurityProfile.ImageCleaner.Enabled
			// }
			// if props.SecurityProfile.ImageCleaner.IntervalHours != nil {
			state.ImageCleanerIntervalHours = pointer.From(props.SecurityProfile.ImageCleaner.IntervalHours)
			//}
			//}

			// state.OIDCIssuerEnabled, state.OIDCIssuerURL = flattenKubernetesAutomaticClusterOidcIssuerProfile(props.OidcIssuerProfile)

			// workloadIdentity := false
			// if props.SecurityProfile != nil && props.SecurityProfile.WorkloadIdentity != nil {
			// 	workloadIdentity = pointer.From(props.SecurityProfile.WorkloadIdentity.Enabled)
			// }
			// state.WorkloadIdentityEnabled = workloadIdentity

			aiToolchainOperatorEnabled := false
			if props.AiToolchainOperatorProfile != nil {
				aiToolchainOperatorEnabled = pointer.From(props.AiToolchainOperatorProfile.Enabled)
			}
			state.AIToolchainOperatorEnabled = aiToolchainOperatorEnabled

			state.SupportPlan = string(pointer.From(props.SupportPlan))

			// if props.AadProfile != nil && (props.DisableLocalAccounts == nil || !*props.DisableLocalAccounts) {
			// 	adminCredentials, err := client.ListClusterAdminCredentials(ctx, *id, managedclusters.ListClusterAdminCredentialsOperationOptions{})
			// 	if err != nil {
			// 		return fmt.Errorf("retrieving Admin Credentials for %s: %+v", id, err)
			//	}
			//	adminKubeConfigRaw, adminKubeConfig := flattenKubernetesClusterCredentials(adminCredentials.Model, "clusterAdmin")
			//	state.KubeAdminConfigRaw = pointer.From(adminKubeConfigRaw)
			//	for _, item := range adminKubeConfig {
			//		if config, ok := item.(map[string]interface{}); ok {
			//			state.KubeAdminConfig = append(state.KubeAdminConfig, KubeConfigModel{
			//				Host:                 config["host"].(string),
			//				Username:             config["username"].(string),
			//				Password:             config["password"].(string),
			//				ClientCertificate:    config["client_certificate"].(string),
			//				ClientKey:            config["client_key"].(string),
			//				ClusterCACertificate: config["cluster_ca_certificate"].(string),
			//			})
			//		}
			//	}
			//}
		}

		state.Identity = flattenIdentityModel(model.Identity)

		kubeConfigRaw, kubeConfig := flattenKubernetesClusterCredentials(credentials.Model, "clusterUser")
		state.KubeConfigRaw = pointer.From(kubeConfigRaw)
		for _, item := range kubeConfig {
			if config, ok := item.(map[string]interface{}); ok {
				state.KubeConfig = append(state.KubeConfig, KubeConfigModel{
					Host:                 config["host"].(string),
					Username:             config["username"].(string),
					Password:             config["password"].(string),
					ClientCertificate:    config["client_certificate"].(string),
					ClientKey:            config["client_key"].(string),
					ClusterCACertificate: config["cluster_ca_certificate"].(string),
				})
			}
		}

		maintenanceClient := metadata.Client.Containers.MaintenanceConfigurationsClient

		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
		configResp, _ := maintenanceClient.Get(ctx, maintenanceId)
		if configurationBody := configResp.Model; configurationBody != nil && configurationBody.Properties != nil {
			state.MaintenanceWindow = flattenKubernetesAutomaticClusterMaintenanceConfigurationDefault(configurationBody.Properties)
		}

		maintenanceId = maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
		configResp, _ = maintenanceClient.Get(ctx, maintenanceId)
		if configurationBody := configResp.Model; configurationBody != nil && configurationBody.Properties != nil && configurationBody.Properties.MaintenanceWindow != nil {
			state.MaintenanceWindowAutoUpgrade = flattenKubernetesAutomaticClusterMaintenanceConfiguration(configurationBody.Properties.MaintenanceWindow)
		}

		maintenanceId = maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
		configResp, _ = maintenanceClient.Get(ctx, maintenanceId)
		if configurationBody := configResp.Model; configurationBody != nil && configurationBody.Properties != nil && configurationBody.Properties.MaintenanceWindow != nil {
			autoUpgradeConfig := flattenKubernetesAutomaticClusterMaintenanceConfiguration(configurationBody.Properties.MaintenanceWindow)
			if len(autoUpgradeConfig) > 0 {
				au := autoUpgradeConfig[0]
				state.MaintenanceWindowNodeOS = []MaintenanceWindowNodeOSModel{{
					Frequency:  au.Frequency,
					Interval:   au.Interval,
					DayOfWeek:  au.DayOfWeek,
					Duration:   au.Duration,
					WeekIndex:  au.WeekIndex,
					DayOfMonth: au.DayOfMonth,
					StartDate:  au.StartDate,
					StartTime:  au.StartTime,
					UTCOffset:  au.UTCOffset,
					NotAllowed: au.NotAllowed,
				}}
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r KubernetesAutomaticClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers
			clusterClient := client.KubernetesClustersClient
			nodePoolsClient := client.AgentPoolsClient
			keyVaultsClient := metadata.Client.KeyVault
			subscriptionId := metadata.Client.Account.SubscriptionId

			id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model KubernetesAutomaticClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := clusterClient.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}

			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving existing %s: properties was nil", *id)
			}

			if err := validateKubernetesAutomaticClusterTyped(&model, existing.Model); err != nil {
				return err
			}

			if existing.Model.Identity != nil && existing.Model.Identity.IdentityIds != nil {
				for k := range existing.Model.Identity.IdentityIds {
					existing.Model.Identity.IdentityIds[k] = identity.UserAssignedIdentityDetails{}
				}
			}

			if metadata.ResourceData.HasChange("service_principal") && !metadata.ResourceData.HasChange("identity") {
				if len(model.ServicePrincipal) > 0 {
					sp := model.ServicePrincipal[0]
					params := managedclusters.ManagedClusterServicePrincipalProfile{
						ClientId: sp.ClientID,
						Secret:   pointer.To(sp.ClientSecret),
					}
					if err := clusterClient.ResetServicePrincipalProfileThenPoll(ctx, *id, params); err != nil {
						return fmt.Errorf("updating Service Principal for %s: %+v", *id, err)
					}

					existing, err = clusterClient.Get(ctx, *id)
					if err != nil {
						return fmt.Errorf("retrieving updated %s: %+v", *id, err)
					}
					if existing.Model == nil || existing.Model.Properties == nil {
						return fmt.Errorf("retrieving updated %s: properties was nil", *id)
					}
				}
			}

			updateCluster := false
			props := existing.Model.Properties

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = tags.Expand(model.Tags)
				updateCluster = true
			}

			// if metadata.ResourceData.HasChange("sku_tier") {
			// 	if existing.Model.Sku == nil {
			// 		existing.Model.Sku = &managedclusters.ManagedClusterSKU{}
			// 	}
			//	tier := managedclusters.ManagedClusterSKUTier(model.SKUTier)
			//	existing.Model.Sku.Tier = &tier
			//	updateCluster = true
			//}

			if metadata.ResourceData.HasChange("kubernetes_version") {
				props.KubernetesVersion = pointer.To(model.KubernetesVersion)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("auto_scaler_profile") {
				props.AutoScalerProfile = expandKubernetesAutomaticClusterAutoScalerProfile(model.AutoScalerProfile)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("monitor_metrics") {
				props.AzureMonitorProfile = expandKubernetesAutomaticClusterAzureMonitorProfile(model.MonitorMetrics)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("service_mesh_profile") {
				props.ServiceMeshProfile = expandKubernetesAutomaticClusterServiceMeshProfile(model.ServiceMeshProfile, props.ServiceMeshProfile)
				updateCluster = true
			}

			// if metadata.ResourceData.HasChange("workload_autoscaler_profile") {
			// 	props.WorkloadAutoScalerProfile = expandKubernetesAutomaticClusterWorkloadAutoscalerProfile(model.WorkloadAutoscalerProfile)
			// 	updateCluster = true
			// }

			if metadata.ResourceData.HasChange("node_provisioning_profile") {
				props.NodeProvisioningProfile = expandKubernetesAutomaticClusterNodeProvisioningProfile(model.NodeProvisioningProfile)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("default_node_pool") {
				props.AgentPoolProfiles, err = ExpandDefaultNodePoolTyped(model.DefaultNodePool)
				if err != nil {
					return err
				}
				agentProfile := ConvertDefaultNodePoolToAgentPool(props.AgentPoolProfiles)
				defaultNodePoolId := agentpools.NewAgentPoolID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, *agentProfile.Name)

				// if a users specified a version - confirm that version is supported on the cluster
				if nodePoolVersion := agentProfile.Properties.CurrentOrchestratorVersion; nodePoolVersion != nil {
					existingNodePool, err := nodePoolsClient.Get(ctx, defaultNodePoolId)
					if err != nil {
						return fmt.Errorf("retrieving Default Node Pool %s: %+v", defaultNodePoolId, err)
					}
					currentNodePoolVersion := ""
					if v := existingNodePool.Model.Properties.OrchestratorVersion; v != nil {
						currentNodePoolVersion = *v
					}

					if err := validateNodePoolAutomaticSupportsVersion(ctx, client, currentNodePoolVersion, defaultNodePoolId, *nodePoolVersion); err != nil {
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
					// "default_node_pool.0.os_disk_type",
					"default_node_pool.0.pod_subnet_id",
					"default_node_pool.0.snapshot_id",
					"default_node_pool.0.ultra_ssd_enabled",
					"default_node_pool.0.vnet_subnet_id",
					"default_node_pool.0.vm_size",
					// "default_node_pool.0.zones",
				}

				// if the default node pool name has changed, it means the initial attempt at resizing failed
				cycleNodePool := metadata.ResourceData.HasChanges(cycleNodePoolProperties...)
				// os_sku could only be updated if the current and new os_sku are either Ubuntu or AzureLinux
				if metadata.ResourceData.HasChange("default_node_pool.0.os_sku") {
					oldOsSkuRaw, newOsSkuRaw := metadata.ResourceData.GetChange("default_node_pool.0.os_sku")
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
					// to provide a seamless updating experience for the vm size of the default node pool we need to cycle the default
					// node pool by provisioning a temporary system node pool, tearing down the former default node pool and then
					// bringing up the new one.

					if v := model.DefaultNodePool[0].TemporaryNameForRotation; v == "" {
						return fmt.Errorf("`temporary_name_for_rotation` must be specified when updating any of the following properties %q", cycleNodePoolProperties)
					}

					temporaryNodePoolName := model.DefaultNodePool[0].TemporaryNameForRotation
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
						if err := retryNodePoolCreation(ctx, nodePoolsClient, tempNodePoolId, tempAgentProfile); err != nil {
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
					if err := retryNodePoolCreation(ctx, nodePoolsClient, defaultNodePoolId, agentProfile); err != nil {
						// if creation of the default node pool fails we automatically fall back to the temporary node pool
						// in func findDefaultNodePool
						return fmt.Errorf("creating default %s: %+v", defaultNodePoolId, err)
					}

					if err := nodePoolsClient.DeleteThenPoll(ctx, tempNodePoolId, agentpools.DefaultDeleteOperationOptions()); err != nil {
						return fmt.Errorf("deleting temporary %s: %+v", tempNodePoolId, err)
					}
				} else {
					if err := nodePoolsClient.CreateOrUpdateThenPoll(ctx, defaultNodePoolId, agentProfile, agentpools.DefaultCreateOrUpdateOperationOptions()); err != nil {
						return fmt.Errorf("updating Default Node Pool %s %+v", defaultNodePoolId, err)
					}
				}
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("node_os_upgrade_channel") {
				if props.AutoUpgradeProfile == nil {
					props.AutoUpgradeProfile = &managedclusters.ManagedClusterAutoUpgradeProfile{}
				}
				if model.NodeOSUpgradeChannel != "" {
					channel := managedclusters.NodeOSUpgradeChannel(model.NodeOSUpgradeChannel)
					props.AutoUpgradeProfile.NodeOSUpgradeChannel = &channel
				}
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("aci_connector_linux") ||
				// metadata.ResourceData.HasChange("azure_policy_enabled") ||
				metadata.ResourceData.HasChange("confidential_computing") ||
				metadata.ResourceData.HasChange("http_application_routing_enabled") ||
				metadata.ResourceData.HasChange("oms_agent") ||
				metadata.ResourceData.HasChange("ingress_application_gateway") ||
				metadata.ResourceData.HasChange("key_vault_secrets_provider") ||
				metadata.ResourceData.HasChange("open_service_mesh_enabled") {
				addonProfiles, err := expandKubernetesAddOnsTyped(&model, metadata.Client.Containers.Environment)
				if err != nil {
					return fmt.Errorf("expanding addons: %+v", err)
				}
				props.AddonProfiles = addonProfiles
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("api_server_access_profile") ||
				metadata.ResourceData.HasChange("run_command_enabled") ||
				metadata.ResourceData.HasChange("private_cluster_public_fqdn_enabled") {
				props.ApiServerAccessProfile = expandKubernetesAutomaticClusterAPIAccessProfile(model)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("azure_active_directory_role_based_access_control") {
				var aadProfile *managedclusters.ManagedClusterAADProfile
				if len(model.AzureActiveDirectoryRBAC) > 0 {
					aad := model.AzureActiveDirectoryRBAC[0]
					aadProfile = &managedclusters.ManagedClusterAADProfile{
						Managed:             pointer.To(true),
						TenantID:            pointer.To(aad.TenantID),
						EnableAzureRBAC:     pointer.To(true),
						AdminGroupObjectIDs: &aad.AdminGroupObjectIDs,
					}
				}
				props.AadProfile = aadProfile

				if props.AadProfile == nil || (props.AadProfile.Managed == nil || !*props.AadProfile.Managed) {
					props.AadProfile = &managedclusters.ManagedClusterAADProfile{}
					if err := clusterClient.ResetAADProfileThenPoll(ctx, *id, *props.AadProfile); err != nil {
						return fmt.Errorf("updating AAD Profile for %s: %+v", *id, err)
					}
				}

				if props.AadProfile != nil && props.AadProfile.Managed != nil && *props.AadProfile.Managed {
					updateCluster = true
				}
			}

			if metadata.ResourceData.HasChange("network_profile") {
				networkProfile, err := expandKubernetesAutomaticClusterNetworkProfile(model.NetworkProfile)
				if err != nil {
					return fmt.Errorf("expanding network profile: %+v", err)
				}
				props.NetworkProfile = networkProfile
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("http_proxy_config") {
				props.HTTPProxyConfig = expandKubernetesAutomaticClusterHttpProxyConfig(model.HTTPProxyConfig)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("storage_profile") {
				props.StorageProfile = expandKubernetesAutomaticClusterStorageProfile(model.StorageProfile)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("web_app_routing") {
				props.IngressProfile = expandKubernetesAutomaticClusterWebAppRouting(model.WebAppRouting, metadata.ResourceData.HasChange("web_app_routing"))
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("upgrade_override") {
				props.UpgradeSettings = expandKubernetesAutomaticClusterUpgradeOverride(model.UpgradeOverride)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("bootstrap_profile") {
				props.BootstrapProfile = expandKubernetesAutomaticClusterBootstrapProfile(model.BootstrapProfile)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("microsoft_defender") ||
				// metadata.ResourceData.HasChange("workload_identity_enabled") ||
				// metadata.ResourceData.HasChange("image_cleaner_enabled") ||
				metadata.ResourceData.HasChange("image_cleaner_interval_hours") ||
				metadata.ResourceData.HasChange("key_management_service") ||
				metadata.ResourceData.HasChange("custom_ca_trust_certificates_base64") {
				if props.SecurityProfile == nil {
					props.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
				}

				if metadata.ResourceData.HasChange("microsoft_defender") {
					props.SecurityProfile.Defender = expandKubernetesAutomaticClusterMicrosoftDefender(model.MicrosoftDefender, metadata.ResourceData.HasChange("microsoft_defender"))
				}

				// if metadata.ResourceData.HasChange("workload_identity_enabled") {
				// 	if props.SecurityProfile.WorkloadIdentity == nil {
				// 		props.SecurityProfile.WorkloadIdentity = &managedclusters.ManagedClusterSecurityProfileWorkloadIdentity{}
				// 	}
				// 	props.SecurityProfile.WorkloadIdentity.Enabled = pointer.To(model.WorkloadIdentityEnabled)
				// }

				if metadata.ResourceData.HasChange("image_cleaner_interval_hours") {
					// if props.SecurityProfile.ImageCleaner == nil {
					// 	props.SecurityProfile.ImageCleaner = &managedclusters.ManagedClusterSecurityProfileImageCleaner{}
					// }
					props.SecurityProfile.ImageCleaner.Enabled = pointer.To(true)
					if model.ImageCleanerIntervalHours > 0 {
						props.SecurityProfile.ImageCleaner.IntervalHours = pointer.To(model.ImageCleanerIntervalHours)
					}
				}

				if metadata.ResourceData.HasChange("key_management_service") {
					props.SecurityProfile.AzureKeyVaultKms, err = expandKubernetesAutomaticClusterKeyManagementService(model.KeyManagementService, ctx, keyVaultsClient, subscriptionId)
					if err != nil {
						return err
					}
				}

				if metadata.ResourceData.HasChange("custom_ca_trust_certificates_base64") {
					if len(model.CustomCATrustCertificatesBase64) > 0 {
						props.SecurityProfile.CustomCATrustCertificates = &model.CustomCATrustCertificatesBase64
					} else {
						props.SecurityProfile.CustomCATrustCertificates = nil
					}
				}

				updateCluster = true
			}

			// if metadata.ResourceData.HasChange("oidc_issuer_enabled") {
			// 	props.OidcIssuerProfile = expandKubernetesAutomaticClusterOidcIssuerProfile(model.OIDCIssuerEnabled)
			// 	updateCluster = true
			// }

			if metadata.ResourceData.HasChange("ai_toolchain_operator_enabled") {
				if props.AiToolchainOperatorProfile == nil {
					props.AiToolchainOperatorProfile = &managedclusters.ManagedClusterAIToolchainOperatorProfile{}
				}
				props.AiToolchainOperatorProfile.Enabled = pointer.To(model.AIToolchainOperatorEnabled)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("support_plan") {
				plan := managedclusters.KubernetesSupportPlan(model.SupportPlan)
				props.SupportPlan = &plan
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("cost_analysis_enabled") {
				if props.MetricsProfile == nil {
					props.MetricsProfile = &managedclusters.ManagedClusterMetricsProfile{}
				}
				if props.MetricsProfile.CostAnalysis == nil {
					props.MetricsProfile.CostAnalysis = &managedclusters.ManagedClusterCostAnalysis{}
				}
				props.MetricsProfile.CostAnalysis.Enabled = pointer.To(model.CostAnalysisEnabled)
				updateCluster = true
			}

			if updateCluster {
				if err := clusterClient.CreateOrUpdateThenPoll(ctx, *id, *existing.Model, managedclusters.DefaultCreateOrUpdateOperationOptions()); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			maintenanceClient := metadata.Client.Containers.MaintenanceConfigurationsClient

			if metadata.ResourceData.HasChange("maintenance_window") {
				maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
				if len(model.MaintenanceWindow) > 0 {
					parameters := maintenanceconfigurations.MaintenanceConfiguration{
						Properties: expandKubernetesAutomaticClusterMaintenanceConfigurationDefault(model.MaintenanceWindow),
					}
					if _, err := maintenanceClient.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
						return fmt.Errorf("updating default maintenance config for %s: %+v", id, err)
					}
				} else {
					if _, err := maintenanceClient.Delete(ctx, maintenanceId); err != nil {
						return fmt.Errorf("deleting default maintenance config for %s: %+v", id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("maintenance_window_auto_upgrade") {
				maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
				if len(model.MaintenanceWindowAutoUpgrade) > 0 {
					parameters := maintenanceconfigurations.MaintenanceConfiguration{
						Properties: expandKubernetesAutomaticClusterMaintenanceConfigurationAutoUpgrade(model.MaintenanceWindowAutoUpgrade),
					}
					if _, err := maintenanceClient.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
						return fmt.Errorf("updating auto upgrade maintenance config for %s: %+v", id, err)
					}
				} else {
					if _, err := maintenanceClient.Delete(ctx, maintenanceId); err != nil {
						return fmt.Errorf("deleting auto upgrade maintenance config for %s: %+v", id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("maintenance_window_node_os") {
				maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
				if len(model.MaintenanceWindowNodeOS) > 0 {
					parameters := maintenanceconfigurations.MaintenanceConfiguration{
						Properties: expandKubernetesAutomaticClusterMaintenanceConfigurationNodeOS(model.MaintenanceWindowNodeOS),
					}
					if _, err := maintenanceClient.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
						return fmt.Errorf("updating node os maintenance config for %s: %+v", id, err)
					}
				} else {
					if _, err := maintenanceClient.Delete(ctx, maintenanceId); err != nil {
						return fmt.Errorf("deleting node os maintenance config for %s: %+v", id, err)
					}
				}
			}

			return nil
		},
	}
}

func (r KubernetesAutomaticClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesClustersClient

			id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model KubernetesAutomaticClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			maintenanceClient := metadata.Client.Containers.MaintenanceConfigurationsClient

			if len(model.MaintenanceWindow) > 0 {
				maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
				if _, err := maintenanceClient.Delete(ctx, maintenanceId); err != nil {
					return fmt.Errorf("deleting default maintenance configuration for %s: %+v", *id, err)
				}
			}

			if len(model.MaintenanceWindowAutoUpgrade) > 0 {
				maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
				if _, err := maintenanceClient.Delete(ctx, maintenanceId); err != nil {
					return fmt.Errorf("deleting auto-upgrade maintenance configuration for %s: %+v", *id, err)
				}
			}

			if len(model.MaintenanceWindowNodeOS) > 0 {
				maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
				if _, err := maintenanceClient.Delete(ctx, maintenanceId); err != nil {
					return fmt.Errorf("deleting node OS maintenance configuration for %s: %+v", *id, err)
				}
			}

			if err := client.DeleteThenPoll(ctx, *id, managedclusters.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandKubernetesAutomaticClusterAPIAccessProfile(model KubernetesAutomaticClusterModel) *managedclusters.ManagedClusterAPIServerAccessProfile {
	apiAccessProfile := &managedclusters.ManagedClusterAPIServerAccessProfile{
		// EnablePrivateCluster:           pointer.To(model.PrivateClusterEnabled),
		EnablePrivateClusterPublicFQDN: pointer.To(model.PrivateClusterPublicFQDNEnabled),
		DisableRunCommand:              pointer.To(!model.RunCommandEnabled),
	}

	if len(model.APIServerAccessProfile) == 0 {
		return apiAccessProfile
	}

	config := model.APIServerAccessProfile[0]
	apiAccessProfile.AuthorizedIPRanges = pointer.To(config.AuthorizedIPRanges)
	// apiAccessProfile.EnableVnetIntegration = pointer.To(true)

	if config.SubnetID != "" {
		apiAccessProfile.SubnetId = pointer.To(config.SubnetID)
	}

	return apiAccessProfile
}

func flattenKubernetesAutomaticClusterAPIAccessProfile(profile *managedclusters.ManagedClusterAPIServerAccessProfile) []APIServerAccessProfileModel {
	apiServerAccessProfile := make([]APIServerAccessProfileModel, 0)

	if profile == nil {
		return apiServerAccessProfile
	}

	// API access profile can be managed by other properties, only return it if one of the properties has been set
	if profile.AuthorizedIPRanges == nil && profile.EnableVnetIntegration == nil && profile.SubnetId == nil {
		return apiServerAccessProfile
	}

	apiServerAccessProfile = append(apiServerAccessProfile, APIServerAccessProfileModel{
		AuthorizedIPRanges: pointer.From(profile.AuthorizedIPRanges),
		// VirtualNetworkIntegrationEnabled: pointer.From(profile.EnableVnetIntegration),
		SubnetID: pointer.From(profile.SubnetId),
	})

	return apiServerAccessProfile
}

func expandKubernetesAutomaticClusterBootstrapProfile(input []BootstrapProfileModel) *managedclusters.ManagedClusterBootstrapProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	var containerRegistryID *string
	if config.ContainerRegistryID != "" {
		containerRegistryID = pointer.To(config.ContainerRegistryID)
	}

	return &managedclusters.ManagedClusterBootstrapProfile{
		ArtifactSource:      pointer.ToEnum[managedclusters.ArtifactSource](config.ArtifactSource),
		ContainerRegistryId: containerRegistryID,
	}
}

func expandKubernetesAutomaticClusterLinuxProfile(input []LinuxProfileModel) *managedclusters.ContainerServiceLinuxProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]

	keyData := ""
	if len(config.SSHKey) > 0 {
		keyData = config.SSHKey[0].KeyData
	}

	return &managedclusters.ContainerServiceLinuxProfile{
		AdminUsername: config.AdminUsername,
		Ssh: managedclusters.ContainerServiceSshConfiguration{
			PublicKeys: []managedclusters.ContainerServiceSshPublicKey{
				{
					KeyData: keyData,
				},
			},
		},
	}
}

func flattenKubernetesAutomaticClusterLinuxProfile(profile *managedclusters.ContainerServiceLinuxProfile) []LinuxProfileModel {
	if profile == nil {
		return []LinuxProfileModel{}
	}

	adminUsername := profile.AdminUsername

	sshKeys := make([]SSHKeyModel, 0)
	ssh := profile.Ssh
	if keys := ssh.PublicKeys; keys != nil {
		for _, sshKey := range keys {
			keyData := ""
			if kd := sshKey.KeyData; kd != "" {
				keyData = kd
			}
			sshKeys = append(sshKeys, SSHKeyModel{
				KeyData: keyData,
			})
		}
	}

	return []LinuxProfileModel{
		{
			AdminUsername: adminUsername,
			SSHKey:        sshKeys,
		},
	}
}

func expandKubernetesAutomaticClusterIdentityProfile(input []KubeletIdentityModel) *map[string]managedclusters.UserAssignedIdentity {
	identityProfile := make(map[string]managedclusters.UserAssignedIdentity)
	if len(input) == 0 {
		return &identityProfile
	}

	values := input[0]

	if values.UserAssignedIdentityID != "" {
		identityProfile["kubeletidentity"] = managedclusters.UserAssignedIdentity{
			ResourceId: pointer.To(values.UserAssignedIdentityID),
			ClientId:   pointer.To(values.ClientID),
			ObjectId:   pointer.To(values.ObjectID),
		}
	}

	return &identityProfile
}

func flattenKubernetesAutomaticClusterIdentityProfile(profile map[string]managedclusters.UserAssignedIdentity) ([]KubeletIdentityModel, error) {
	if profile == nil {
		return []KubeletIdentityModel{}, nil
	}

	kubeletIdentity := make([]KubeletIdentityModel, 0)
	if kubeletidentity, ok := profile["kubeletidentity"]; ok {
		clientId := ""
		if clientid := kubeletidentity.ClientId; clientid != nil {
			clientId = *clientid
		}

		objectId := ""
		if objectid := kubeletidentity.ObjectId; objectid != nil {
			objectId = *objectid
		}

		userAssignedIdentityId := ""
		if resourceid := kubeletidentity.ResourceId; resourceid != nil {
			parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*resourceid)
			if err != nil {
				return nil, err
			}

			userAssignedIdentityId = parsedId.ID()
		}

		kubeletIdentity = append(kubeletIdentity, KubeletIdentityModel{
			ClientID:               clientId,
			ObjectID:               objectId,
			UserAssignedIdentityID: userAssignedIdentityId,
		})
	}

	return kubeletIdentity, nil
}

func expandKubernetesAutomaticClusterWindowsProfile(input []WindowsProfileModel) *managedclusters.ManagedClusterWindowsProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]

	license := managedclusters.LicenseTypeNone
	if config.License != "" {
		license = managedclusters.LicenseType(config.License)
	}

	gmsaProfile := expandAutomaticGMSAProfile(config.GMSA)

	return &managedclusters.ManagedClusterWindowsProfile{
		AdminUsername: config.AdminUsername,
		AdminPassword: pointer.To(config.AdminPassword),
		LicenseType:   pointer.To(license),
		GmsaProfile:   gmsaProfile,
	}
}

func flattenKubernetesAutomaticClusterWindowsProfile(profile *managedclusters.ManagedClusterWindowsProfile, config KubernetesAutomaticClusterModel) []WindowsProfileModel {
	if profile == nil {
		return []WindowsProfileModel{}
	}

	adminUsername := profile.AdminUsername

	adminPassword := ""
	if len(config.WindowsProfile) != 0 {
		adminPassword = config.WindowsProfile[0].AdminPassword
	}

	license := ""
	if profile.LicenseType != nil && pointer.From(profile.LicenseType) != managedclusters.LicenseTypeNone {
		license = string(pointer.From(profile.LicenseType))
	}

	gmsaProfile := flattenAutomaticGMSAProfile(profile.GmsaProfile)

	return []WindowsProfileModel{
		{
			AdminUsername: adminUsername,
			AdminPassword: adminPassword,
			License:       license,
			GMSA:          gmsaProfile,
		},
	}
}

func expandAutomaticGMSAProfile(input []GMSAModel) *managedclusters.WindowsGmsaProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]

	if !config.GMSAProfileEnabled {
		return nil
	}

	return &managedclusters.WindowsGmsaProfile{
		Enabled:        pointer.To(true),
		DnsServer:      pointer.To(config.DNSServer),
		RootDomainName: pointer.To(config.RootDomain),
	}
}

func flattenAutomaticGMSAProfile(profile *managedclusters.WindowsGmsaProfile) []GMSAModel {
	if profile == nil {
		return []GMSAModel{}
	}

	return []GMSAModel{
		{
			GMSAProfileEnabled: pointer.From(profile.Enabled),
			DNSServer:          pointer.From(profile.DnsServer),
			RootDomain:         pointer.From(profile.RootDomainName),
		},
	}
}

// func expandAutomaticBootstrapProfile(input []BootstrapProfileModel) *managedclusters.ManagedClusterBootstrapProfile {
// 	if len(input) == 0 {
// 		return nil
// 	}

// 	config := input[0]
// 	var containerRegistryID *string
// 	if config.ContainerRegistryID != "" {
// 		containerRegistryID = pointer.To(config.ContainerRegistryID)
// 	}

// 	return &managedclusters.ManagedClusterBootstrapProfile{
// 		ArtifactSource:      pointer.ToEnum[managedclusters.ArtifactSource](config.ArtifactSource),
// 		ContainerRegistryId: containerRegistryID,
// 	}
// }

// func flattenAutomaticBootstrapProfile(profile *managedclusters.ManagedClusterBootstrapProfile) ([]BootstrapProfileModel, error) {
// 	if profile == nil || profile.ArtifactSource == nil {
// 		return []BootstrapProfileModel{}, nil
// 	}

// 	var containerRegistryID string
// 	if profile.ContainerRegistryId != nil {
// 		id, err := registries.ParseRegistryID(*profile.ContainerRegistryId)
// 		if err != nil {
// 			return nil, err
// 		}
// 		containerRegistryID = id.ID()
// 	}

// 	return []BootstrapProfileModel{
// 		{
// 			ArtifactSource:      string(*profile.ArtifactSource),
// 			ContainerRegistryID: containerRegistryID,
// 		},
// 	}, nil
// }

func idsToAutomaticResourceReferences(ids []string) *[]managedclusters.ResourceReference {
	if len(ids) == 0 {
		return nil
	}

	results := make([]managedclusters.ResourceReference, 0)
	for _, id := range ids {
		results = append(results, managedclusters.ResourceReference{Id: pointer.To(id)})
	}

	if len(results) > 0 {
		return &results
	}

	return nil
}

func automaticResourceReferencesToIds(refs *[]managedclusters.ResourceReference) []string {
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

func expandKubernetesAutomaticClusterNetworkProfile(input []NetworkProfileModel) (*managedclusters.ContainerServiceNetworkProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	config := input[0]

	// networkPlugin := config.NetworkPlugin
	networkMode := config.NetworkMode
	// if networkPlugin != "azure" && networkMode != "" {
	// 	return nil, fmt.Errorf("`network_mode` cannot be set if `network_plugin` is not `azure`")
	// }

	networkPolicy := config.NetworkPolicy
	loadBalancerSku := config.LoadBalancerSKU
	outboundType := config.OutboundType

	ipVersions, err := expandAutomaticIPVersions(config.IPVersions)
	if err != nil {
		return nil, err
	}

	networkProfile := managedclusters.ContainerServiceNetworkProfile{
		// NetworkPlugin:   pointer.To(managedclusters.NetworkPlugin(networkPlugin)),
		NetworkMode:     pointer.To(managedclusters.NetworkMode(networkMode)),
		NetworkPolicy:   pointer.To(managedclusters.NetworkPolicy(networkPolicy)),
		LoadBalancerSku: pointer.To(managedclusters.LoadBalancerSku(loadBalancerSku)),
		OutboundType:    pointer.To(managedclusters.OutboundType(outboundType)),
		IPFamilies:      ipVersions,
	}

	// if config.NetworkDataPlane != "" {
	// 	networkProfile.NetworkDataplane = pointer.To(managedclusters.NetworkDataplane(config.NetworkDataPlane))
	// }

	if config.NetworkPluginMode != "" {
		networkProfile.NetworkPluginMode = pointer.To(managedclusters.NetworkPluginMode(config.NetworkPluginMode))
	}

	if len(config.LoadBalancerProfile) > 0 {
		if !strings.EqualFold(loadBalancerSku, "standard") {
			return nil, fmt.Errorf("only load balancer SKU 'Standard' supports load balancer profiles. Provided load balancer type: %s", loadBalancerSku)
		}

		networkProfile.LoadBalancerProfile = expandAutomaticLoadBalancerProfile(config.LoadBalancerProfile)
	}

	if len(config.NATGatewayProfile) > 0 {
		if !strings.EqualFold(loadBalancerSku, "standard") {
			return nil, fmt.Errorf("only load balancer SKU 'Standard' supports NAT Gateway profiles. Provided load balancer type: %s", loadBalancerSku)
		}

		networkProfile.NatGatewayProfile = expandAutomaticNatGatewayProfile(config.NATGatewayProfile)
	}

	if config.DNSServiceIP != "" {
		networkProfile.DnsServiceIP = pointer.To(config.DNSServiceIP)
	}

	if config.PodCIDR != "" {
		networkProfile.PodCidr = pointer.To(config.PodCIDR)
	}

	if len(config.PodCIDRs) > 0 {
		networkProfile.PodCidrs = &config.PodCIDRs
	}

	if config.ServiceCIDR != "" {
		networkProfile.ServiceCidr = pointer.To(config.ServiceCIDR)
	}

	if len(config.ServiceCIDRs) > 0 {
		networkProfile.ServiceCidrs = &config.ServiceCIDRs
	}

	networkProfile.AdvancedNetworking = expandKubernetesAutomaticClusterAdvancedNetworking(config.AdvancedNetworking)

	return &networkProfile, nil
}

func expandKubernetesAutomaticClusterAdvancedNetworking(input []AdvancedNetworkingModel) *managedclusters.AdvancedNetworking {
	advancedNetworkingConfig := managedclusters.AdvancedNetworking{
		Enabled: pointer.To(false),
		Observability: &managedclusters.AdvancedNetworkingObservability{
			Enabled: pointer.To(false),
		},
		Security: &managedclusters.AdvancedNetworkingSecurity{
			Enabled: pointer.To(false),
		},
	}
	if len(input) != 0 {
		config := input[0]
		observabilityEnabled := config.ObservabilityEnabled
		securityEnabled := config.SecurityEnabled

		advancedNetworkingConfig = managedclusters.AdvancedNetworking{
			Enabled: pointer.To(true),
			Observability: &managedclusters.AdvancedNetworkingObservability{
				Enabled: pointer.To(observabilityEnabled),
			},
			Security: &managedclusters.AdvancedNetworkingSecurity{
				Enabled: pointer.To(securityEnabled),
			},
		}
	}

	return &advancedNetworkingConfig
}

func flattenKubernetesAutomaticClusterAdvancedNetworking(advancedNetworking *managedclusters.AdvancedNetworking) []AdvancedNetworkingModel {
	if advancedNetworking == nil || !pointer.From(advancedNetworking.Enabled) {
		return []AdvancedNetworkingModel{}
	}

	observabilityEnabled := false
	if advancedNetworking.Observability != nil {
		observabilityEnabled = pointer.From(advancedNetworking.Observability.Enabled)
	}

	securityEnabled := false
	if advancedNetworking.Security != nil {
		securityEnabled = pointer.From(advancedNetworking.Security.Enabled)
	}

	return []AdvancedNetworkingModel{
		{
			ObservabilityEnabled: observabilityEnabled,
			SecurityEnabled:      securityEnabled,
		},
	}
}

func expandAutomaticLoadBalancerProfile(input []LoadBalancerProfileModel) *managedclusters.ManagedClusterLoadBalancerProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	profile := &managedclusters.ManagedClusterLoadBalancerProfile{}

	if config.IdleTimeoutInMinutes != 0 {
		profile.IdleTimeoutInMinutes = pointer.To(config.IdleTimeoutInMinutes)
	}

	if config.OutboundPortsAllocated != 0 {
		profile.AllocatedOutboundPorts = pointer.To(config.OutboundPortsAllocated)
	}

	if config.ManagedOutboundIPCount > 0 {
		profile.ManagedOutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileManagedOutboundIPs{
			Count: pointer.To(config.ManagedOutboundIPCount),
		}
	}

	if len(config.OutboundIPPrefixIDs) > 0 {
		profile.OutboundIPPrefixes = &managedclusters.ManagedClusterLoadBalancerProfileOutboundIPPrefixes{
			PublicIPPrefixes: idsToAutomaticResourceReferences(config.OutboundIPPrefixIDs),
		}
	}

	if len(config.OutboundIPAddressIDs) > 0 {
		profile.OutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileOutboundIPs{
			PublicIPs: idsToAutomaticResourceReferences(config.OutboundIPAddressIDs),
		}
	}

	if config.BackendPoolType != "" {
		profile.BackendPoolType = pointer.To(managedclusters.BackendPoolType(config.BackendPoolType))
	}

	return profile
}

func expandAutomaticIPVersions(input []string) (*[]managedclusters.IPFamily, error) {
	if len(input) == 0 {
		return nil, nil
	}

	ipv := make([]managedclusters.IPFamily, 0)
	for _, data := range input {
		ipv = append(ipv, managedclusters.IPFamily(data))
	}

	if len(ipv) == 1 && ipv[0] == managedclusters.IPFamilyIPvSix {
		return nil, fmt.Errorf("`ip_versions` must be `IPv4` or `IPv4` and `IPv6`. `IPv6` alone is not supported")
	}

	return &ipv, nil
}

func expandAutomaticNatGatewayProfile(input []NATGatewayProfileModel) *managedclusters.ManagedClusterNATGatewayProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	profile := &managedclusters.ManagedClusterNATGatewayProfile{}

	if config.IdleTimeoutInMinutes != 0 {
		profile.IdleTimeoutInMinutes = pointer.To(config.IdleTimeoutInMinutes)
	}

	if config.ManagedOutboundIPCount > 0 {
		profile.ManagedOutboundIPProfile = &managedclusters.ManagedClusterManagedOutboundIPProfile{
			Count: pointer.To(config.ManagedOutboundIPCount),
		}
	}

	return profile
}

func flattenKubernetesAutomaticClusterNetworkProfile(profile *managedclusters.ContainerServiceNetworkProfile) []NetworkProfileModel {
	if profile == nil {
		return []NetworkProfileModel{}
	}

	dnsServiceIP := ""
	if profile.DnsServiceIP != nil {
		dnsServiceIP = *profile.DnsServiceIP
	}

	serviceCidr := ""
	if profile.ServiceCidr != nil {
		serviceCidr = *profile.ServiceCidr
	}

	podCidr := ""
	if profile.PodCidr != nil {
		podCidr = *profile.PodCidr
	}

	// networkPlugin := ""
	// if profile.NetworkPlugin != nil {
	// 	networkPlugin = string(*profile.NetworkPlugin)
	// }

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

	lbProfiles := make([]LoadBalancerProfileModel, 0)
	if lbp := profile.LoadBalancerProfile; lbp != nil {
		lb := LoadBalancerProfileModel{}

		if v := lbp.AllocatedOutboundPorts; v != nil {
			lb.OutboundPortsAllocated = *v
		}

		if v := lbp.IdleTimeoutInMinutes; v != nil {
			lb.IdleTimeoutInMinutes = *v
		}

		if ips := lbp.ManagedOutboundIPs; ips != nil {
			if count := ips.Count; count != nil {
				lb.ManagedOutboundIPCount = *count
			}
		}

		if oip := lbp.OutboundIPs; oip != nil {
			if poip := oip.PublicIPs; poip != nil {
				lb.OutboundIPAddressIDs = automaticResourceReferencesToIds(poip)
			}
		}

		if oip := lbp.OutboundIPPrefixes; oip != nil {
			if pip := oip.PublicIPPrefixes; pip != nil {
				lb.OutboundIPPrefixIDs = automaticResourceReferencesToIds(pip)
			}
		}

		if v := lbp.BackendPoolType; v != nil {
			lb.BackendPoolType = string(*v)
		}

		lb.EffectiveOutboundIPs = automaticResourceReferencesToIds(profile.LoadBalancerProfile.EffectiveOutboundIPs)
		lbProfiles = append(lbProfiles, lb)
	}

	ngwProfiles := make([]NATGatewayProfileModel, 0)
	if ngwp := profile.NatGatewayProfile; ngwp != nil {
		ng := NATGatewayProfileModel{}

		if v := ngwp.IdleTimeoutInMinutes; v != nil {
			ng.IdleTimeoutInMinutes = *v
		}

		if ips := ngwp.ManagedOutboundIPProfile; ips != nil {
			if count := ips.Count; count != nil {
				ng.ManagedOutboundIPCount = *count
			}
		}

		ng.EffectiveOutboundIPs = automaticResourceReferencesToIds(profile.NatGatewayProfile.EffectiveOutboundIPs)
		ngwProfiles = append(ngwProfiles, ng)
	}

	ipVersions := make([]string, 0)
	if ipfs := profile.IPFamilies; ipfs != nil {
		for _, item := range *ipfs {
			ipVersions = append(ipVersions, string(item))
		}
	}

	// TODO - Remove the workaround below once issue https://github.com/Azure/azure-rest-api-specs/issues/18056 is resolved
	sku := profile.LoadBalancerSku
	for _, v := range managedclusters.PossibleValuesForLoadBalancerSku() {
		if strings.EqualFold(v, string(*sku)) {
			lsSku := managedclusters.LoadBalancerSku(v)
			sku = &lsSku
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

	// networkDataPlane := string(managedclusters.NetworkDataplaneAzure)
	// if v := profile.NetworkDataplane; v != nil {
	// 	networkDataPlane = string(pointer.From(v))
	// }

	advancedNetworking := flattenKubernetesAutomaticClusterAdvancedNetworking(profile.AdvancedNetworking)

	podCidrs := []string{}
	if profile.PodCidrs != nil {
		podCidrs = *profile.PodCidrs
	}

	serviceCidrs := []string{}
	if profile.ServiceCidrs != nil {
		serviceCidrs = *profile.ServiceCidrs
	}

	return []NetworkProfileModel{
		{
			DNSServiceIP: dnsServiceIP,
			// NetworkDataPlane:    networkDataPlane,
			LoadBalancerSKU:     string(*sku),
			LoadBalancerProfile: lbProfiles,
			NATGatewayProfile:   ngwProfiles,
			IPVersions:          ipVersions,
			// NetworkPlugin:       networkPlugin,
			NetworkPluginMode:  networkPluginMode,
			NetworkMode:        networkMode,
			NetworkPolicy:      networkPolicy,
			PodCIDR:            podCidr,
			PodCIDRs:           podCidrs,
			ServiceCIDR:        serviceCidr,
			ServiceCIDRs:       serviceCidrs,
			OutboundType:       outboundType,
			AdvancedNetworking: advancedNetworking,
		},
	}
}

func expandKubernetesAutomaticClusterAutoScalerProfile(input []AutoScalerProfileModel) *managedclusters.ManagedClusterPropertiesAutoScalerProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	profile := &managedclusters.ManagedClusterPropertiesAutoScalerProfile{}

	if config.BalanceSimilarNodeGroups {
		profile.BalanceSimilarNodeGroups = pointer.To(strconv.FormatBool(config.BalanceSimilarNodeGroups))
	}

	if config.Expander != "" {
		profile.Expander = pointer.To(managedclusters.Expander(config.Expander))
	}

	if config.MaxGracefulTerminationSec != "" {
		profile.MaxGracefulTerminationSec = pointer.To(config.MaxGracefulTerminationSec)
	}

	if config.MaxNodeProvisioningTime != "" {
		profile.MaxNodeProvisionTime = pointer.To(config.MaxNodeProvisioningTime)
	}

	if config.NewPodScaleUpDelay != "" {
		profile.NewPodScaleUpDelay = pointer.To(config.NewPodScaleUpDelay)
	}

	if config.ScaleDownDelayAfterAdd != "" {
		profile.ScaleDownDelayAfterAdd = pointer.To(config.ScaleDownDelayAfterAdd)
	}

	if config.ScaleDownDelayAfterDelete != "" {
		profile.ScaleDownDelayAfterDelete = pointer.To(config.ScaleDownDelayAfterDelete)
	}

	if config.ScaleDownDelayAfterFailure != "" {
		profile.ScaleDownDelayAfterFailure = pointer.To(config.ScaleDownDelayAfterFailure)
	}

	if config.ScaleDownUnneeded != "" {
		profile.ScaleDownUnneededTime = pointer.To(config.ScaleDownUnneeded)
	}

	if config.ScaleDownUnready != "" {
		profile.ScaleDownUnreadyTime = pointer.To(config.ScaleDownUnready)
	}

	if config.ScaleDownUtilizationThreshold != "" {
		profile.ScaleDownUtilizationThreshold = pointer.To(config.ScaleDownUtilizationThreshold)
	}

	if config.EmptyBulkDeleteMax != "" {
		profile.MaxEmptyBulkDelete = pointer.To(config.EmptyBulkDeleteMax)
	}

	if config.SkipNodesWithLocalStorage {
		profile.SkipNodesWithLocalStorage = pointer.To(strconv.FormatBool(config.SkipNodesWithLocalStorage))
	}

	if config.SkipNodesWithSystemPods {
		profile.SkipNodesWithSystemPods = pointer.To(strconv.FormatBool(config.SkipNodesWithSystemPods))
	}

	return profile
}

func flattenKubernetesAutomaticClusterAutoScalerProfile(profile *managedclusters.ManagedClusterPropertiesAutoScalerProfile) ([]AutoScalerProfileModel, error) {
	if profile == nil {
		return []AutoScalerProfileModel{}, nil
	}

	var balanceSimilarNodeGroups bool
	if profile.BalanceSimilarNodeGroups != nil {
		b, err := strconv.ParseBool(*profile.BalanceSimilarNodeGroups)
		if err != nil {
			return nil, fmt.Errorf("parsing BalanceSimilarNodeGroups: %+v", err)
		}
		balanceSimilarNodeGroups = b
	}

	expander := ""
	if profile.Expander != nil {
		expander = string(*profile.Expander)
	}

	maxGracefulTerminationSec := ""
	if profile.MaxGracefulTerminationSec != nil {
		maxGracefulTerminationSec = *profile.MaxGracefulTerminationSec
	}

	MaxNodeProvisioningTime := ""
	if profile.MaxNodeProvisionTime != nil {
		MaxNodeProvisioningTime = *profile.MaxNodeProvisionTime
	}

	newPodScaleUpDelay := ""
	if profile.NewPodScaleUpDelay != nil {
		newPodScaleUpDelay = *profile.NewPodScaleUpDelay
	}

	scaleDownDelayAfterAdd := ""
	if profile.ScaleDownDelayAfterAdd != nil {
		scaleDownDelayAfterAdd = *profile.ScaleDownDelayAfterAdd
	}

	scaleDownDelayAfterDelete := ""
	if profile.ScaleDownDelayAfterDelete != nil {
		scaleDownDelayAfterDelete = *profile.ScaleDownDelayAfterDelete
	}

	scaleDownDelayAfterFailure := ""
	if profile.ScaleDownDelayAfterFailure != nil {
		scaleDownDelayAfterFailure = *profile.ScaleDownDelayAfterFailure
	}

	scaleDownUnneededTime := ""
	if profile.ScaleDownUnneededTime != nil {
		scaleDownUnneededTime = *profile.ScaleDownUnneededTime
	}

	scaleDownUnreadyTime := ""
	if profile.ScaleDownUnreadyTime != nil {
		scaleDownUnreadyTime = *profile.ScaleDownUnreadyTime
	}

	scaleDownUtilizationThreshold := ""
	if profile.ScaleDownUtilizationThreshold != nil {
		scaleDownUtilizationThreshold = *profile.ScaleDownUtilizationThreshold
	}

	emptyBulkDeleteMax := ""
	if profile.MaxEmptyBulkDelete != nil {
		emptyBulkDeleteMax = *profile.MaxEmptyBulkDelete
	}

	var skipNodesWithLocalStorage bool
	if profile.SkipNodesWithLocalStorage != nil {
		b, err := strconv.ParseBool(*profile.SkipNodesWithLocalStorage)
		if err != nil {
			return nil, fmt.Errorf("parsing SkipNodesWithLocalStorage: %+v", err)
		}
		skipNodesWithLocalStorage = b
	}

	var skipNodesWithSystemPods bool
	if profile.SkipNodesWithSystemPods != nil {
		b, err := strconv.ParseBool(*profile.SkipNodesWithSystemPods)
		if err != nil {
			return nil, fmt.Errorf("parsing SkipNodesWithSystemPods: %+v", err)
		}
		skipNodesWithSystemPods = b
	}

	scanInterval := ""
	if profile.ScanInterval != nil {
		scanInterval = *profile.ScanInterval
	}

	var daemonsetEvictionForEmptyNodesEnabled bool
	if profile.DaemonsetEvictionForEmptyNodes != nil {
		daemonsetEvictionForEmptyNodesEnabled = *profile.DaemonsetEvictionForEmptyNodes
	}

	var daemonsetEvictionForOccupiedNodesEnabled bool
	if profile.DaemonsetEvictionForOccupiedNodes != nil {
		daemonsetEvictionForOccupiedNodesEnabled = *profile.DaemonsetEvictionForOccupiedNodes
	}

	var ignoreDaemonsetsUtilizationEnabled bool
	if profile.IgnoreDaemonsetsUtilization != nil {
		ignoreDaemonsetsUtilizationEnabled = *profile.IgnoreDaemonsetsUtilization
	}

	return []AutoScalerProfileModel{{
		BalanceSimilarNodeGroups:                 balanceSimilarNodeGroups,
		DaemonsetEvictionForEmptyNodesEnabled:    daemonsetEvictionForEmptyNodesEnabled,
		DaemonsetEvictionForOccupiedNodesEnabled: daemonsetEvictionForOccupiedNodesEnabled,
		Expander:                                 expander,
		IgnoreDaemonsetsUtilizationEnabled:       ignoreDaemonsetsUtilizationEnabled,
		MaxGracefulTerminationSec:                maxGracefulTerminationSec,
		MaxNodeProvisioningTime:                  MaxNodeProvisioningTime,
		MaxUnreadyNodes:                          0,
		MaxUnreadyPercentage:                     0.0,
		NewPodScaleUpDelay:                       newPodScaleUpDelay,
		ScanInterval:                             scanInterval,
		ScaleDownDelayAfterAdd:                   scaleDownDelayAfterAdd,
		ScaleDownDelayAfterDelete:                scaleDownDelayAfterDelete,
		ScaleDownDelayAfterFailure:               scaleDownDelayAfterFailure,
		ScaleDownUnneeded:                        scaleDownUnneededTime,
		ScaleDownUnready:                         scaleDownUnreadyTime,
		ScaleDownUtilizationThreshold:            scaleDownUtilizationThreshold,
		EmptyBulkDeleteMax:                       emptyBulkDeleteMax,
		SkipNodesWithLocalStorage:                skipNodesWithLocalStorage,
		SkipNodesWithSystemPods:                  skipNodesWithSystemPods,
	}}, nil
}

// func expandKubernetesAutomaticClusterWorkloadAutoscalerProfile(input []WorkloadAutoscalerProfileModel) *managedclusters.ManagedClusterWorkloadAutoScalerProfile {
// 	if len(input) == 0 {
// 		return nil
// 	}
//
//	config := input[0]
//	profile := &managedclusters.ManagedClusterWorkloadAutoScalerProfile{}
//
//	if config.KEDAEnabled {
//		profile.Keda = &managedclusters.ManagedClusterWorkloadAutoScalerProfileKeda{
//			Enabled: config.KEDAEnabled,
//		}
//	}
//
//	if config.VerticalPodAutoscalerEnabled {
//		profile.VerticalPodAutoscaler = &managedclusters.ManagedClusterWorkloadAutoScalerProfileVerticalPodAutoscaler{
//			Enabled: config.VerticalPodAutoscalerEnabled,
//		}
//	}
//
//	return profile
//}

// func flattenKubernetesAutomaticClusterWorkloadAutoscalerProfile(profile *managedclusters.ManagedClusterWorkloadAutoScalerProfile) []WorkloadAutoscalerProfileModel {
// 	// The API always returns an empty WorkloadAutoScalerProfile object even if none of these values have ever been set
// 	if profile == nil || (profile.Keda == nil && profile.VerticalPodAutoscaler == nil) {
// 		return []WorkloadAutoscalerProfileModel{}
//	}
//
//	kedaEnabled := false
//	if profile.Keda != nil && profile.Keda.Enabled {
//		kedaEnabled = profile.Keda.Enabled
//	}
//
//	vpaEnabled := false
//	if profile.VerticalPodAutoscaler != nil && profile.VerticalPodAutoscaler.Enabled {
//		vpaEnabled = profile.VerticalPodAutoscaler.Enabled
//	}
//
//	return []WorkloadAutoscalerProfileModel{{
//		KEDAEnabled:                  kedaEnabled,
//		VerticalPodAutoscalerEnabled: vpaEnabled,
//	}}
//}

func expandKubernetesAutomaticClusterNodeProvisioningProfile(input []NodeProvisioningProfileModel) *managedclusters.ManagedClusterNodeProvisioningProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	profile := &managedclusters.ManagedClusterNodeProvisioningProfile{}

	// if config.Mode != "" {
	// 	profile.Mode = pointer.To(managedclusters.NodeProvisioningMode(config.Mode))
	// }

	if config.DefaultNodePools != "" {
		profile.DefaultNodePools = pointer.To(managedclusters.NodeProvisioningDefaultNodePools(config.DefaultNodePools))
	}

	return profile
}

func flattenKubernetesAutomaticClusterNodeProvisioningProfile(profile *managedclusters.ManagedClusterNodeProvisioningProfile) []NodeProvisioningProfileModel {
	if profile == nil {
		return []NodeProvisioningProfileModel{}
	}

	// mode := ""
	// if profile.Mode != nil {
	// 	mode = string(*profile.Mode)
	// }

	defaultNodePools := ""
	if profile.DefaultNodePools != nil {
		defaultNodePools = string(*profile.DefaultNodePools)
	}

	return []NodeProvisioningProfileModel{{
		// Mode:             mode,
		DefaultNodePools: defaultNodePools,
	}}
}

func expandKubernetesAutomaticClusterAzureMonitorProfile(input []MonitorMetricsModel) *managedclusters.ManagedClusterAzureMonitorProfile {
	if len(input) == 0 {
		return &managedclusters.ManagedClusterAzureMonitorProfile{
			Metrics: &managedclusters.ManagedClusterAzureMonitorProfileMetrics{
				Enabled: false,
			},
		}
	}

	config := input[0]

	profile := &managedclusters.ManagedClusterAzureMonitorProfile{
		Metrics: &managedclusters.ManagedClusterAzureMonitorProfileMetrics{
			Enabled: config.MonitorMetricsEnabled,
			KubeStateMetrics: &managedclusters.ManagedClusterAzureMonitorProfileKubeStateMetrics{
				MetricAnnotationsAllowList: pointer.To(config.AnnotationsAllowed),
				MetricLabelsAllowlist:      pointer.To(config.LabelsAllowed),
			},
		},
	}

	return profile
}

func flattenKubernetesAutomaticClusterAzureMonitorProfile(input *managedclusters.ManagedClusterAzureMonitorProfile) []MonitorMetricsModel {
	if input == nil || input.Metrics == nil || !input.Metrics.Enabled {
		return []MonitorMetricsModel{}
	}

	if input.Metrics.KubeStateMetrics == nil {
		return []MonitorMetricsModel{{
			MonitorMetricsEnabled: true,
			AnnotationsAllowed:    "",
			LabelsAllowed:         "",
		}}
	}

	annotationsAllowed := ""
	if input.Metrics.KubeStateMetrics.MetricAnnotationsAllowList != nil {
		annotationsAllowed = *input.Metrics.KubeStateMetrics.MetricAnnotationsAllowList
	}

	labelsAllowed := ""
	if input.Metrics.KubeStateMetrics.MetricLabelsAllowlist != nil {
		labelsAllowed = *input.Metrics.KubeStateMetrics.MetricLabelsAllowlist
	}

	return []MonitorMetricsModel{{
		MonitorMetricsEnabled: true,
		AnnotationsAllowed:    annotationsAllowed,
		LabelsAllowed:         labelsAllowed,
	}}
}

func expandKubernetesAutomaticClusterMetricsProfile(input bool) *managedclusters.ManagedClusterMetricsProfile {
	return &managedclusters.ManagedClusterMetricsProfile{
		CostAnalysis: &managedclusters.ManagedClusterCostAnalysis{
			Enabled: pointer.To(input),
		},
	}
}

func flattenKubernetesAutomaticClusterMetricsProfile(input *managedclusters.ManagedClusterMetricsProfile) bool {
	if input == nil || input.CostAnalysis == nil {
		return false
	}
	return pointer.From(input.CostAnalysis.Enabled)
}

func expandKubernetesAutomaticClusterMaintenanceConfigurationDefault(input []MaintenanceWindowModel) *maintenanceconfigurations.MaintenanceConfigurationProperties {
	if len(input) == 0 {
		return nil
	}
	value := input[0]
	return &maintenanceconfigurations.MaintenanceConfigurationProperties{
		NotAllowedTime: expandKubernetesAutomaticClusterMaintenanceConfigurationTimeSpans(value.NotAllowed),
		TimeInWeek:     expandKubernetesAutomaticClusterMaintenanceConfigurationTimeInWeeks(value.Allowed),
	}
}

func expandKubernetesAutomaticClusterMaintenanceConfigurationForCreate(input []MaintenanceWindowAutoUpgradeModel) *maintenanceconfigurations.MaintenanceConfigurationProperties {
	if len(input) == 0 {
		return nil
	}
	value := input[0]

	var schedule maintenanceconfigurations.Schedule

	if value.Frequency == "Daily" {
		schedule = maintenanceconfigurations.Schedule{
			Daily: &maintenanceconfigurations.DailySchedule{
				IntervalDays: value.Interval,
			},
		}
	}
	if value.Frequency == "Weekly" {
		schedule = maintenanceconfigurations.Schedule{
			Weekly: &maintenanceconfigurations.WeeklySchedule{
				IntervalWeeks: value.Interval,
				DayOfWeek:     maintenanceconfigurations.WeekDay(value.DayOfWeek),
			},
		}
	}
	if value.Frequency == "AbsoluteMonthly" {
		schedule = maintenanceconfigurations.Schedule{
			AbsoluteMonthly: &maintenanceconfigurations.AbsoluteMonthlySchedule{
				DayOfMonth:     value.DayOfMonth,
				IntervalMonths: value.Interval,
			},
		}
	}
	if value.Frequency == "RelativeMonthly" {
		schedule = maintenanceconfigurations.Schedule{
			RelativeMonthly: &maintenanceconfigurations.RelativeMonthlySchedule{
				DayOfWeek:      maintenanceconfigurations.WeekDay(value.DayOfWeek),
				WeekIndex:      maintenanceconfigurations.Type(value.WeekIndex),
				IntervalMonths: value.Interval,
			},
		}
	}

	output := &maintenanceconfigurations.MaintenanceConfigurationProperties{
		MaintenanceWindow: &maintenanceconfigurations.MaintenanceWindow{
			StartTime:       value.StartTime,
			UtcOffset:       pointer.To(value.UTCOffset),
			NotAllowedDates: expandKubernetesAutomaticClusterMaintenanceConfigurationDateSpans(value.NotAllowed),
			Schedule:        schedule,
		},
	}

	if value.StartDate != "" {
		startDate, _ := time.Parse(time.RFC3339, value.StartDate)
		output.MaintenanceWindow.StartDate = pointer.To(startDate.Format("2006-01-02"))
	}

	if value.Duration != 0 {
		output.MaintenanceWindow.DurationHours = value.Duration
	}

	return output
}

func expandKubernetesAutomaticClusterMaintenanceConfigurationAutoUpgrade(input []MaintenanceWindowAutoUpgradeModel) *maintenanceconfigurations.MaintenanceConfigurationProperties {
	return expandKubernetesAutomaticClusterMaintenanceConfigurationForCreate(input)
}

func expandKubernetesAutomaticClusterMaintenanceConfigurationNodeOS(input []MaintenanceWindowNodeOSModel) *maintenanceconfigurations.MaintenanceConfigurationProperties {
	if len(input) == 0 {
		return nil
	}
	// Convert MaintenanceWindowNodeOSModel to MaintenanceWindowAutoUpgradeModel since they have the same structure
	converted := []MaintenanceWindowAutoUpgradeModel{{
		Frequency:  input[0].Frequency,
		Interval:   input[0].Interval,
		DayOfWeek:  input[0].DayOfWeek,
		Duration:   input[0].Duration,
		WeekIndex:  input[0].WeekIndex,
		DayOfMonth: input[0].DayOfMonth,
		StartDate:  input[0].StartDate,
		StartTime:  input[0].StartTime,
		UTCOffset:  input[0].UTCOffset,
		NotAllowed: input[0].NotAllowed,
	}}
	return expandKubernetesAutomaticClusterMaintenanceConfigurationForCreate(converted)
}

// func expandKubernetesAutomaticClusterMaintenanceConfigurationForUpdate(input []MaintenanceWindowAutoUpgradeModel, existing *maintenanceconfigurations.MaintenanceConfigurationProperties) *maintenanceconfigurations.MaintenanceConfigurationProperties {
// 	if len(input) == 0 {
// 		return nil
// 	}
// 	value := input[0]

// 	var schedule maintenanceconfigurations.Schedule

// 	if value.Frequency == "Daily" {
// 		schedule = maintenanceconfigurations.Schedule{
// 			Daily: &maintenanceconfigurations.DailySchedule{
// 				IntervalDays: value.Interval,
// 			},
// 		}
// 	}
// 	if value.Frequency == "Weekly" {
// 		schedule = maintenanceconfigurations.Schedule{
// 			Weekly: &maintenanceconfigurations.WeeklySchedule{
// 				IntervalWeeks: value.Interval,
// 				DayOfWeek:     maintenanceconfigurations.WeekDay(value.DayOfWeek),
// 			},
// 		}
// 	}
// 	if value.Frequency == "AbsoluteMonthly" {
// 		schedule = maintenanceconfigurations.Schedule{
// 			AbsoluteMonthly: &maintenanceconfigurations.AbsoluteMonthlySchedule{
// 				DayOfMonth:     value.DayOfMonth,
// 				IntervalMonths: value.Interval,
// 			},
// 		}
// 	}
// 	if value.Frequency == "RelativeMonthly" {
// 		schedule = maintenanceconfigurations.Schedule{
// 			RelativeMonthly: &maintenanceconfigurations.RelativeMonthlySchedule{
// 				DayOfWeek:      maintenanceconfigurations.WeekDay(value.DayOfWeek),
// 				WeekIndex:      maintenanceconfigurations.Type(value.WeekIndex),
// 				IntervalMonths: value.Interval,
// 			},
// 		}
// 	}

// 	output := &maintenanceconfigurations.MaintenanceConfigurationProperties{
// 		MaintenanceWindow: &maintenanceconfigurations.MaintenanceWindow{
// 			StartTime:       value.StartTime,
// 			UtcOffset:       pointer.To(value.UTCOffset),
// 			NotAllowedDates: expandKubernetesAutomaticClusterMaintenanceConfigurationDateSpans(value.NotAllowed),
// 			Schedule:        schedule,
// 		},
// 	}

// 	if value.StartDate != "" {
// 		startDate, _ := time.Parse(time.RFC3339, value.StartDate)
// 		startDateStr := startDate.Format("2006-01-02")
// 		// start_date is an Optional+Computed property, the default value returned by the API could be invalid during update, so we only set it if it's different from the existing value
// 		if existing == nil || existing.MaintenanceWindow.StartDate == nil || *existing.MaintenanceWindow.StartDate != startDateStr {
// 			output.MaintenanceWindow.StartDate = pointer.To(startDateStr)
// 		}
// 	}

// 	if value.Duration != 0 {
// 		output.MaintenanceWindow.DurationHours = value.Duration
// 	}

// 	return output
// }

func expandKubernetesAutomaticClusterMaintenanceConfigurationTimeSpans(input []MaintenanceWindowNotAllowedModel) *[]maintenanceconfigurations.TimeSpan {
	results := make([]maintenanceconfigurations.TimeSpan, 0)
	for _, item := range input {
		start, _ := time.Parse(time.RFC3339, item.Start)
		end, _ := time.Parse(time.RFC3339, item.End)
		results = append(results, maintenanceconfigurations.TimeSpan{
			Start: pointer.To(start.Format("2006-01-02T15:04:05Z07:00")),
			End:   pointer.To(end.Format("2006-01-02T15:04:05Z07:00")),
		})
	}
	return &results
}

func expandKubernetesAutomaticClusterMaintenanceConfigurationDateSpans(input []MaintenanceWindowNotAllowedModel) *[]maintenanceconfigurations.DateSpan {
	results := make([]maintenanceconfigurations.DateSpan, 0)
	for _, item := range input {
		start, _ := time.Parse(time.RFC3339, item.Start)
		end, _ := time.Parse(time.RFC3339, item.End)
		results = append(results, maintenanceconfigurations.DateSpan{
			Start: start.Format("2006-01-02"),
			End:   end.Format("2006-01-02"),
		})
	}
	return &results
}

func expandKubernetesAutomaticClusterMaintenanceConfigurationTimeInWeeks(input []MaintenanceWindowAllowedModel) *[]maintenanceconfigurations.TimeInWeek {
	results := make([]maintenanceconfigurations.TimeInWeek, 0)
	for _, item := range input {
		results = append(results, maintenanceconfigurations.TimeInWeek{
			Day:       pointer.To(maintenanceconfigurations.WeekDay(item.Day)),
			HourSlots: pointer.To(item.Hours),
		})
	}
	return &results
}

func flattenKubernetesAutomaticClusterMaintenanceConfiguration(input *maintenanceconfigurations.MaintenanceWindow) []MaintenanceWindowAutoUpgradeModel {
	results := make([]MaintenanceWindowAutoUpgradeModel, 0)
	if input == nil {
		return results
	}

	startDate := ""
	if input.StartDate != nil {
		startDate = *input.StartDate + "T00:00:00Z"
	}
	utcOffset := ""
	if input.UtcOffset != nil {
		utcOffset = *input.UtcOffset
	}

	windowModel := MaintenanceWindowAutoUpgradeModel{
		NotAllowed: flattenKubernetesAutomaticClusterMaintenanceConfigurationDateSpans(input.NotAllowedDates),
		Duration:   input.DurationHours,
		StartDate:  startDate,
		StartTime:  input.StartTime,
		UTCOffset:  utcOffset,
	}

	scheduleProps := flattenKubernetesAutomaticClusterMaintenanceConfigurationSchedule(input.Schedule)
	windowModel.Frequency = scheduleProps["frequency"].(string)
	windowModel.Interval = scheduleProps["interval"].(int64)
	windowModel.DayOfWeek = scheduleProps["day_of_week"].(string)
	windowModel.WeekIndex = scheduleProps["week_index"].(string)
	windowModel.DayOfMonth = scheduleProps["day_of_month"].(int64)

	return append(results, windowModel)
}

func flattenKubernetesAutomaticClusterMaintenanceConfigurationSchedule(input maintenanceconfigurations.Schedule) map[string]interface{} {
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

	dayOfMonth := int64(0)
	if input.AbsoluteMonthly != nil {
		frequency = "AbsoluteMonthly"
		interval = input.AbsoluteMonthly.IntervalMonths
		dayOfMonth = input.AbsoluteMonthly.DayOfMonth
	}

	weekIndex := ""
	if input.RelativeMonthly != nil {
		frequency = "RelativeMonthly"
		interval = input.RelativeMonthly.IntervalMonths
		dayOfWeek = string(input.RelativeMonthly.DayOfWeek)
		weekIndex = string(input.RelativeMonthly.WeekIndex)
	}

	return map[string]interface{}{
		"frequency":    frequency,
		"interval":     interval,
		"day_of_week":  dayOfWeek,
		"week_index":   weekIndex,
		"day_of_month": dayOfMonth,
	}
}

func flattenKubernetesAutomaticClusterMaintenanceConfigurationDefault(input *maintenanceconfigurations.MaintenanceConfigurationProperties) []MaintenanceWindowModel {
	results := make([]MaintenanceWindowModel, 0)
	if input == nil {
		return results
	}
	return append(results, MaintenanceWindowModel{
		NotAllowed: flattenKubernetesAutomaticClusterMaintenanceConfigurationTimeSpans(input.NotAllowedTime),
		Allowed:    flattenKubernetesAutomaticClusterMaintenanceConfigurationTimeInWeeks(input.TimeInWeek),
	})
}

func flattenKubernetesAutomaticClusterMaintenanceConfigurationTimeSpans(input *[]maintenanceconfigurations.TimeSpan) []MaintenanceWindowNotAllowedModel {
	results := make([]MaintenanceWindowNotAllowedModel, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var end string
		if item.End != nil {
			end = *item.End
		}
		var start string
		if item.Start != nil {
			start = *item.Start
		}
		results = append(results, MaintenanceWindowNotAllowedModel{
			End:   end,
			Start: start,
		})
	}
	return results
}

func flattenKubernetesAutomaticClusterMaintenanceConfigurationDateSpans(input *[]maintenanceconfigurations.DateSpan) []MaintenanceWindowNotAllowedModel {
	results := make([]MaintenanceWindowNotAllowedModel, 0)
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
		results = append(results, MaintenanceWindowNotAllowedModel{
			End:   end + "T00:00:00Z",
			Start: start + "T00:00:00Z",
		})
	}
	return results
}

func flattenKubernetesAutomaticClusterMaintenanceConfigurationTimeInWeeks(input *[]maintenanceconfigurations.TimeInWeek) []MaintenanceWindowAllowedModel {
	results := make([]MaintenanceWindowAllowedModel, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		day := ""
		if item.Day != nil {
			day = string(*item.Day)
		}
		hours := make([]int64, 0)
		if item.HourSlots != nil {
			hours = *item.HourSlots
		}
		results = append(results, MaintenanceWindowAllowedModel{
			Day:   day,
			Hours: hours,
		})
	}
	return results
}

func expandKubernetesAutomaticClusterHttpProxyConfig(input []HTTPProxyConfigModel) *managedclusters.ManagedClusterHTTPProxyConfig {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	httpProxyConfig := managedclusters.ManagedClusterHTTPProxyConfig{
		HTTPProxy:  pointer.To(config.HTTPProxy),
		HTTPSProxy: pointer.To(config.HTTPSProxy),
		NoProxy:    pointer.To(config.NoProxy),
	}

	if config.TrustedCA != "" {
		httpProxyConfig.TrustedCa = pointer.To(config.TrustedCA)
	}

	return &httpProxyConfig
}

func flattenKubernetesAutomaticClusterHttpProxyConfig(httpProxyConfig *managedclusters.ManagedClusterHTTPProxyConfig) []HTTPProxyConfigModel {
	if httpProxyConfig == nil {
		return []HTTPProxyConfigModel{}
	}

	httpProxy := ""
	if httpProxyConfig.HTTPProxy != nil {
		httpProxy = *httpProxyConfig.HTTPProxy
	}

	httpsProxy := ""
	if httpProxyConfig.HTTPSProxy != nil {
		httpsProxy = *httpProxyConfig.HTTPSProxy
	}

	noProxyList := make([]string, 0)
	if httpProxyConfig.NoProxy != nil {
		noProxyList = append(noProxyList, *httpProxyConfig.NoProxy...)
	}

	trustedCa := ""
	if httpProxyConfig.TrustedCa != nil {
		trustedCa = *httpProxyConfig.TrustedCa
	}

	return []HTTPProxyConfigModel{{
		HTTPProxy:  httpProxy,
		HTTPSProxy: httpsProxy,
		NoProxy:    noProxyList,
		TrustedCA:  trustedCa,
	}}
}

func expandKubernetesAutomaticClusterKeyManagementService(input []KeyManagementServiceModel, ctx context.Context, keyVaultsClient *keyVaultClient.Client, subscriptionId string) (*managedclusters.AzureKeyVaultKms, error) {
	if len(input) == 0 {
		return &managedclusters.AzureKeyVaultKms{
			Enabled: pointer.To(false),
		}, nil
	}

	config := input[0]
	kvAccess := managedclusters.KeyVaultNetworkAccessTypes(config.KeyVaultNetworkAccess)

	kms := &managedclusters.AzureKeyVaultKms{
		Enabled:               pointer.To(true),
		KeyId:                 pointer.To(config.KeyVaultKeyID),
		KeyVaultNetworkAccess: pointer.To(kvAccess),
	}

	// Set Key vault Resource ID in case public access is disabled
	if kvAccess == managedclusters.KeyVaultNetworkAccessTypesPrivate {
		subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)

		nestedItemType := keyvault.NestedItemTypeKey

		keyVaultKeyId, err := keyvault.ParseNestedItemID(*kms.KeyId, keyvault.VersionTypeVersioned, nestedItemType)
		if err != nil {
			return nil, err
		}
		keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, keyVaultKeyId.KeyVaultBaseURL)
		if err != nil {
			return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", keyVaultKeyId.KeyVaultBaseURL, err)
		}

		kms.KeyVaultResourceId = keyVaultID
	}

	return kms, nil
}

func flattenKubernetesAutomaticClusterKeyManagementService(kms *managedclusters.AzureKeyVaultKms) []KeyManagementServiceModel {
	if kms == nil || kms.Enabled == nil || !*kms.Enabled {
		return []KeyManagementServiceModel{}
	}

	keyVaultKeyID := ""
	if kms.KeyId != nil {
		keyVaultKeyID = *kms.KeyId
	}

	keyVaultNetworkAccess := ""
	if kms.KeyVaultNetworkAccess != nil {
		keyVaultNetworkAccess = string(*kms.KeyVaultNetworkAccess)
	}

	return []KeyManagementServiceModel{{
		KeyVaultKeyID:         keyVaultKeyID,
		KeyVaultNetworkAccess: keyVaultNetworkAccess,
	}}
}

func expandKubernetesAutomaticClusterMicrosoftDefender(input []MicrosoftDefenderModel, hasChange bool) *managedclusters.ManagedClusterSecurityProfileDefender {
	if (len(input) == 0) && hasChange {
		return &managedclusters.ManagedClusterSecurityProfileDefender{
			SecurityMonitoring: &managedclusters.ManagedClusterSecurityProfileDefenderSecurityMonitoring{
				Enabled: pointer.To(false),
			},
		}
	} else if len(input) == 0 {
		return nil
	}

	config := input[0]
	return &managedclusters.ManagedClusterSecurityProfileDefender{
		SecurityMonitoring: &managedclusters.ManagedClusterSecurityProfileDefenderSecurityMonitoring{
			Enabled: pointer.To(true),
		},
		LogAnalyticsWorkspaceResourceId: pointer.To(config.LogAnalyticsWorkspaceID),
	}
}

func flattenKubernetesAutomaticClusterMicrosoftDefender(input *managedclusters.ManagedClusterSecurityProfile) []MicrosoftDefenderModel {
	if input == nil || input.Defender == nil || input.Defender.SecurityMonitoring == nil || input.Defender.SecurityMonitoring.Enabled == nil || !*input.Defender.SecurityMonitoring.Enabled {
		return []MicrosoftDefenderModel{}
	}

	logAnalyticsWorkspace := ""
	if v := input.Defender.LogAnalyticsWorkspaceResourceId; v != nil {
		logAnalyticsWorkspace = *v
	}

	return []MicrosoftDefenderModel{{
		LogAnalyticsWorkspaceID: logAnalyticsWorkspace,
	}}
}

func expandKubernetesAutomaticClusterWebAppRouting(input []WebAppRoutingModel, hasChange bool) *managedclusters.ManagedClusterIngressProfile {
	if len(input) == 0 && hasChange {
		return &managedclusters.ManagedClusterIngressProfile{
			WebAppRouting: &managedclusters.ManagedClusterIngressProfileWebAppRouting{
				Enabled: pointer.To(false),
			},
		}
	} else if len(input) == 0 {
		return nil
	}

	config := input[0]
	out := managedclusters.ManagedClusterIngressProfile{
		WebAppRouting: &managedclusters.ManagedClusterIngressProfileWebAppRouting{
			Enabled: pointer.To(true),
			Nginx: &managedclusters.ManagedClusterIngressProfileNginx{
				DefaultIngressControllerType: (*managedclusters.NginxIngressControllerType)(pointer.To(config.DefaultNginxController)),
			},
		},
	}

	if len(config.DNSZoneIDs) > 0 {
		out.WebAppRouting.DnsZoneResourceIds = pointer.To(config.DNSZoneIDs)
	}

	return &out
}

func flattenKubernetesAutomaticClusterWebAppRouting(input *managedclusters.ManagedClusterIngressProfile) []WebAppRoutingModel {
	if input == nil || input.WebAppRouting == nil || input.WebAppRouting.Enabled == nil || !*input.WebAppRouting.Enabled {
		return []WebAppRoutingModel{}
	}

	dnsZoneIDs := make([]string, 0)
	if input.WebAppRouting.DnsZoneResourceIds != nil {
		dnsZoneIDs = *input.WebAppRouting.DnsZoneResourceIds
	}

	defaultNginxController := managedclusters.NginxIngressControllerTypeAnnotationControlled
	if input.WebAppRouting.Nginx != nil {
		defaultNginxController = *input.WebAppRouting.Nginx.DefaultIngressControllerType
	}

	webAppRoutingIdentity := make([]WebAppRoutingIdentityModel, 0)
	if input.WebAppRouting.Identity != nil {
		clientID := ""
		if input.WebAppRouting.Identity.ClientId != nil {
			clientID = *input.WebAppRouting.Identity.ClientId
		}
		objectID := ""
		if input.WebAppRouting.Identity.ObjectId != nil {
			objectID = *input.WebAppRouting.Identity.ObjectId
		}
		userAssignedIdentityID := ""
		if input.WebAppRouting.Identity.ResourceId != nil {
			userAssignedIdentityID = *input.WebAppRouting.Identity.ResourceId
		}

		webAppRoutingIdentity = append(webAppRoutingIdentity, WebAppRoutingIdentityModel{
			ClientID:               clientID,
			ObjectID:               objectID,
			UserAssignedIdentityID: userAssignedIdentityID,
		})
	}

	return []WebAppRoutingModel{{
		DNSZoneIDs:             dnsZoneIDs,
		DefaultNginxController: string(defaultNginxController),
		WebAppRoutingIdentity:  webAppRoutingIdentity,
	}}
}

func expandKubernetesAutomaticClusterStorageProfile(input []StorageProfileModel) *managedclusters.ManagedClusterStorageProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	profile := managedclusters.ManagedClusterStorageProfile{
		BlobCSIDriver: &managedclusters.ManagedClusterStorageProfileBlobCSIDriver{
			Enabled: pointer.To(config.BlobDriverEnabled),
		},
		DiskCSIDriver: &managedclusters.ManagedClusterStorageProfileDiskCSIDriver{
			Enabled: pointer.To(config.DiskDriverEnabled),
		},
		FileCSIDriver: &managedclusters.ManagedClusterStorageProfileFileCSIDriver{
			Enabled: pointer.To(config.FileDriverEnabled),
		},
		SnapshotController: &managedclusters.ManagedClusterStorageProfileSnapshotController{
			Enabled: pointer.To(config.SnapshotControllerEnabled),
		},
	}

	return &profile
}

func flattenKubernetesAutomaticClusterStorageProfile(profile *managedclusters.ManagedClusterStorageProfile) []StorageProfileModel {
	if profile == nil {
		return []StorageProfileModel{}
	}

	blobDriverEnabled := false
	if profile.BlobCSIDriver != nil && profile.BlobCSIDriver.Enabled != nil {
		blobDriverEnabled = *profile.BlobCSIDriver.Enabled
	}

	diskDriverEnabled := false
	if profile.DiskCSIDriver != nil && profile.DiskCSIDriver.Enabled != nil {
		diskDriverEnabled = *profile.DiskCSIDriver.Enabled
	}

	fileDriverEnabled := false
	if profile.FileCSIDriver != nil && profile.FileCSIDriver.Enabled != nil {
		fileDriverEnabled = *profile.FileCSIDriver.Enabled
	}

	snapshotControllerEnabled := false
	if profile.SnapshotController != nil && profile.SnapshotController.Enabled != nil {
		snapshotControllerEnabled = *profile.SnapshotController.Enabled
	}

	return []StorageProfileModel{{
		BlobDriverEnabled:         blobDriverEnabled,
		DiskDriverEnabled:         diskDriverEnabled,
		FileDriverEnabled:         fileDriverEnabled,
		SnapshotControllerEnabled: snapshotControllerEnabled,
	}}
}

func expandKubernetesAutomaticClusterServiceMeshProfile(input []ServiceMeshProfileModel, existing *managedclusters.ServiceMeshProfile) *managedclusters.ServiceMeshProfile {
	if len(input) == 0 {
		// explicitly disable istio if it was enabled before
		if existing != nil && existing.Mode == managedclusters.ServiceMeshModeIstio {
			return &managedclusters.ServiceMeshProfile{
				Mode: managedclusters.ServiceMeshModeDisabled,
			}
		}
		return nil
	}

	config := input[0]
	mode := config.Mode

	profile := managedclusters.ServiceMeshProfile{}

	if managedclusters.ServiceMeshMode(mode) == managedclusters.ServiceMeshModeIstio {
		profile.Mode = managedclusters.ServiceMeshMode(mode)
		profile.Istio = &managedclusters.IstioServiceMesh{}
		profile.Istio.Components = &managedclusters.IstioComponents{}

		istioIngressGatewaysList := make([]managedclusters.IstioIngressGateway, 0)

		ingressGatewayElementInternal := managedclusters.IstioIngressGateway{
			Enabled: config.InternalIngressGatewayEnabled,
			Mode:    managedclusters.IstioIngressGatewayModeInternal,
		}
		istioIngressGatewaysList = append(istioIngressGatewaysList, ingressGatewayElementInternal)

		ingressGatewayElementExternal := managedclusters.IstioIngressGateway{
			Enabled: config.ExternalIngressGatewayEnabled,
			Mode:    managedclusters.IstioIngressGatewayModeExternal,
		}
		istioIngressGatewaysList = append(istioIngressGatewaysList, ingressGatewayElementExternal)

		profile.Istio.Components.IngressGateways = &istioIngressGatewaysList

		if len(config.CertificateAuthority) > 0 {
			certificateAuthority := expandKubernetesAutomaticClusterServiceMeshProfileCertificateAuthority(config.CertificateAuthority)
			profile.Istio.CertificateAuthority = certificateAuthority
		}

		if len(config.Revisions) > 0 {
			profile.Istio.Revisions = pointer.To(config.Revisions)
		}
	}

	return &profile
}

func flattenKubernetesAutomaticClusterServiceMeshProfile(profile *managedclusters.ServiceMeshProfile) []ServiceMeshProfileModel {
	if profile == nil || profile.Mode != managedclusters.ServiceMeshModeIstio || profile.Istio == nil {
		return []ServiceMeshProfileModel{}
	}

	mode := string(profile.Mode)
	revisions := make([]string, 0)
	if profile.Istio.Revisions != nil {
		revisions = *profile.Istio.Revisions
	}

	internalIngressGatewayEnabled := false
	externalIngressGatewayEnabled := false

	if profile.Istio.Components != nil && profile.Istio.Components.IngressGateways != nil {
		for _, gateway := range *profile.Istio.Components.IngressGateways {
			if gateway.Mode == managedclusters.IstioIngressGatewayModeInternal {
				internalIngressGatewayEnabled = gateway.Enabled
			}
			if gateway.Mode == managedclusters.IstioIngressGatewayModeExternal {
				externalIngressGatewayEnabled = gateway.Enabled
			}
		}
	}

	certificateAuthority := flattenKubernetesAutomaticClusterServiceMeshProfileCertificateAuthority(profile.Istio.CertificateAuthority)

	return []ServiceMeshProfileModel{{
		Mode:                          mode,
		Revisions:                     revisions,
		InternalIngressGatewayEnabled: internalIngressGatewayEnabled,
		ExternalIngressGatewayEnabled: externalIngressGatewayEnabled,
		CertificateAuthority:          certificateAuthority,
	}}
}

func expandKubernetesAutomaticClusterServiceMeshProfileCertificateAuthority(input []CertificateAuthorityModel) *managedclusters.IstioCertificateAuthority {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	return &managedclusters.IstioCertificateAuthority{
		Plugin: &managedclusters.IstioPluginCertificateAuthority{
			KeyVaultId:          pointer.To(config.KeyVaultID),
			RootCertObjectName:  pointer.To(config.RootCertObjectName),
			CertChainObjectName: pointer.To(config.CertChainObjectName),
			CertObjectName:      pointer.To(config.CertObjectName),
			KeyObjectName:       pointer.To(config.KeyObjectName),
		},
	}
}

func flattenKubernetesAutomaticClusterServiceMeshProfileCertificateAuthority(certificateAuthority *managedclusters.IstioCertificateAuthority) []CertificateAuthorityModel {
	if certificateAuthority == nil || certificateAuthority.Plugin == nil {
		return []CertificateAuthorityModel{}
	}

	plugin := certificateAuthority.Plugin

	keyVaultID := ""
	if plugin.KeyVaultId != nil {
		keyVaultID = *plugin.KeyVaultId
	}

	rootCertObjectName := ""
	if plugin.RootCertObjectName != nil {
		rootCertObjectName = *plugin.RootCertObjectName
	}

	certChainObjectName := ""
	if plugin.CertChainObjectName != nil {
		certChainObjectName = *plugin.CertChainObjectName
	}

	certObjectName := ""
	if plugin.CertObjectName != nil {
		certObjectName = *plugin.CertObjectName
	}

	keyObjectName := ""
	if plugin.KeyObjectName != nil {
		keyObjectName = *plugin.KeyObjectName
	}

	return []CertificateAuthorityModel{{
		KeyVaultID:          keyVaultID,
		RootCertObjectName:  rootCertObjectName,
		CertChainObjectName: certChainObjectName,
		CertObjectName:      certObjectName,
		KeyObjectName:       keyObjectName,
	}}
}

func expandKubernetesAutomaticClusterUpgradeOverride(input []UpgradeOverrideModel) *managedclusters.ClusterUpgradeSettings {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	return &managedclusters.ClusterUpgradeSettings{
		OverrideSettings: &managedclusters.UpgradeOverrideSettings{
			ForceUpgrade: pointer.To(config.ForceUpgradeEnabled),
			Until:        pointer.To(config.EffectiveUntil),
		},
	}
}

func flattenKubernetesAutomaticClusterUpgradeOverride(input *managedclusters.ClusterUpgradeSettings) []UpgradeOverrideModel {
	if input == nil || input.OverrideSettings == nil {
		return []UpgradeOverrideModel{}
	}

	forceUpgradeEnabled := false
	if input.OverrideSettings.ForceUpgrade != nil {
		forceUpgradeEnabled = pointer.From(input.OverrideSettings.ForceUpgrade)
	}

	effectiveUntil := ""
	if input.OverrideSettings.Until != nil {
		effectiveUntil = pointer.From(input.OverrideSettings.Until)
	}

	return []UpgradeOverrideModel{{
		ForceUpgradeEnabled: forceUpgradeEnabled,
		EffectiveUntil:      effectiveUntil,
	}}
}

func flattenKubernetesAutomaticClusterBootstrapProfile(profile *managedclusters.ManagedClusterBootstrapProfile) []BootstrapProfileModel {
	if profile == nil {
		return []BootstrapProfileModel{}
	}

	return []BootstrapProfileModel{{
		ArtifactSource:      string(pointer.From(profile.ArtifactSource)),
		ContainerRegistryID: pointer.From(profile.ContainerRegistryId),
	}}
}

func flattenKubernetesAutomaticClusterAzureActiveDirectoryRBAC(profile *managedclusters.ManagedClusterAADProfile) []AzureActiveDirectoryRBACModel {
	if profile == nil || profile.Managed == nil || !*profile.Managed {
		return []AzureActiveDirectoryRBACModel{}
	}

	result := AzureActiveDirectoryRBACModel{
		TenantID: pointer.From(profile.TenantID),
	}

	if profile.AdminGroupObjectIDs != nil {
		result.AdminGroupObjectIDs = *profile.AdminGroupObjectIDs
	}

	return []AzureActiveDirectoryRBACModel{result}
}

// func flattenKubernetesAutomaticClusterKubeConfig(input []interface{}) []KubeConfigModel {
// 	results := make([]KubeConfigModel, 0)
// 	if len(input) == 0 {
// 		return results
// 	}

// 	for _, item := range input {
// 		config := item.(map[string]interface{})
// 		results = append(results, KubeConfigModel{
// 			Host:                 config["host"].(string),
// 			Username:             config["username"].(string),
// 			Password:             config["password"].(string),
// 			ClientCertificate:    config["client_certificate"].(string),
// 			ClientKey:            config["client_key"].(string),
// 			ClusterCACertificate: config["cluster_ca_certificate"].(string),
// 		})
// 	}

// 	return results
// }

// func flattenKubernetesAutomaticClusterMaintenanceConfigurationNodeOS(input *maintenanceconfigurations.MaintenanceWindow) []MaintenanceWindowNodeOSModel {
// 	results := make([]MaintenanceWindowNodeOSModel, 0)
// 	if input == nil {
// 		return results
// 	}

// 	scheduleData := flattenKubernetesAutomaticClusterMaintenanceConfigurationSchedule(input.Schedule)

// 	result := MaintenanceWindowNodeOSModel{
// 		Frequency:  scheduleData["frequency"].(string),
// 		Interval:   int64(scheduleData["interval"].(int)),
// 		StartTime:  input.StartTime,
// 		UTCOffset:  pointer.From(input.UtcOffset),
// 		Duration:   input.DurationHours,
// 		NotAllowed: flattenKubernetesAutomaticClusterMaintenanceConfigurationDateSpans(input.NotAllowedDates),
// 	}

// 	if v, ok := scheduleData["day_of_week"].(string); ok {
// 		result.DayOfWeek = v
// 	}
// 	if v, ok := scheduleData["week_index"].(string); ok {
// 		result.WeekIndex = v
// 	}
// 	if v, ok := scheduleData["day_of_month"].(int); ok {
// 		result.DayOfMonth = int64(v)
// 	}
// 	if input.StartDate != nil {
// 		result.StartDate = pointer.From(input.StartDate) + "T00:00:00Z"
// 	}

// 	results = append(results, result)
// 	return results
// }

// func expandKubernetesAutomaticClusterOidcIssuerProfile(input bool) *managedclusters.ManagedClusterOIDCIssuerProfile {
// 	return &managedclusters.ManagedClusterOIDCIssuerProfile{
// 		Enabled: &input,
// 	}
// }

// func flattenKubernetesAutomaticClusterOidcIssuerProfile(profile *managedclusters.ManagedClusterOIDCIssuerProfile) (bool, string) {
// 	if profile == nil {
// 		return false, ""
// 	}
//
//	enabled := false
//	if profile.Enabled != nil {
//		enabled = pointer.From(profile.Enabled)
//	}
//
//	issuerURL := ""
//	if profile.IssuerURL != nil {
//		issuerURL = pointer.From(profile.IssuerURL)
//	}
//
//	return enabled, issuerURL
//}

func expandKubernetesAutomaticClusterEdgeZone(input string) *edgezones.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &edgezones.Model{
		Name: normalized,
	}
}

func flattenKubernetesAutomaticClusterEdgeZone(input *edgezones.Model) string {
	if input == nil || input.Name == "" {
		return ""
	}
	return edgezones.NormalizeNilable(&input.Name)
}

func expandIdentityModel(input []IdentityModel) *identity.SystemOrUserAssignedMap {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	identityType := identity.Type(config.Type)

	identityIds := make(map[string]identity.UserAssignedIdentityDetails)
	for _, id := range config.IdentityIds {
		identityIds[id] = identity.UserAssignedIdentityDetails{}
	}

	return &identity.SystemOrUserAssignedMap{
		Type:        identityType,
		IdentityIds: identityIds,
	}
}

func flattenIdentityModel(input *identity.SystemOrUserAssignedMap) []IdentityModel {
	if input == nil {
		return []IdentityModel{}
	}

	// Only set IdentityIds for UserAssigned type to avoid empty array in plan
	var identityIds []string
	if input.Type == identity.TypeUserAssigned || input.Type == identity.TypeSystemAssignedUserAssigned {
		if len(input.IdentityIds) > 0 {
			identityIds = make([]string, 0, len(input.IdentityIds))
			for id := range input.IdentityIds {
				identityIds = append(identityIds, id)
			}
		}
	}

	principalId := ""
	if input.PrincipalId != "" {
		principalId = input.PrincipalId
	}

	tenantId := ""
	if input.TenantId != "" {
		tenantId = input.TenantId
	}

	return []IdentityModel{{
		Type:        string(input.Type),
		IdentityIds: identityIds,
		PrincipalId: principalId,
		TenantId:    tenantId,
	}}
}
