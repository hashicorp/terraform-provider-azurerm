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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/agentpools"
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
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type KubernetesAutomaticClusterModel struct {
	Name                            string                              `tfschema:"name"`
	Location                        string                              `tfschema:"location"`
	ResourceGroupName               string                              `tfschema:"resource_group_name"`
	APIServerAccessProfile          []APIServerAccessProfileModel       `tfschema:"api_server_access"`
	AutoScalerProfile               []AutoScalerProfileModel            `tfschema:"auto_scaler"`
	AzureActiveDirectoryRBAC        []AzureActiveDirectoryRBACModel     `tfschema:"azure_active_directory_role_based_access_control"`
	BootstrapProfile                []BootstrapProfileModel             `tfschema:"bootstrap"`
	CostAnalysisEnabled             bool                                `tfschema:"cost_analysis_enabled"`
	CustomCATrustCertificatesBase64 []string                            `tfschema:"custom_ca_trust_certificates_base64"`
	DefaultNodePool                 []DefaultNodePoolModel              `tfschema:"default_node_pool"`
	DiskEncryptionSetID             string                              `tfschema:"disk_encryption_set_id"`
	DNSPrefix                       string                              `tfschema:"dns_prefix"`
	DNSPrefixPrivateCluster         string                              `tfschema:"dns_prefix_private_cluster"`
	HTTPProxyConfig                 []HTTPProxyConfigModel              `tfschema:"http_proxy_config"`
	Identity                        []identity.SystemOrUserAssignedList `tfschema:"identity"`
	ImageCleanerIntervalHours       int64                               `tfschema:"image_cleaner_interval_in_hours"`
	KeyManagementService            []KeyManagementServiceModel         `tfschema:"key_management_service"`
	KubeletIdentity                 []KubeletIdentityModel              `tfschema:"kubelet_identity"`
	KubernetesVersion               string                              `tfschema:"kubernetes_version"`
	LinuxProfile                    []LinuxProfileModel                 `tfschema:"linux_profile"`
	// MaintenanceWindow               []MaintenanceWindowModel            `tfschema:"maintenance_window"`
	// MaintenanceWindowAutoUpgrade    []MaintenanceWindowAutoUpgradeModel `tfschema:"maintenance_window_auto_upgrade"`
	// MaintenanceWindowNodeOS         []MaintenanceWindowNodeOSModel      `tfschema:"maintenance_window_node_os"`
	MicrosoftDefender          []MicrosoftDefenderModel  `tfschema:"microsoft_defender"`
	MonitorMetrics             []MonitorMetricsModel     `tfschema:"monitor_metrics"`
	NetworkProfile             []NetworkProfileModel     `tfschema:"network"`
	NodeResourceGroup          string                    `tfschema:"node_resource_group_name"`
	PrivateCluster             []PrivateClusterModel     `tfschema:"private_cluster"`
	RunCommandEnabled          bool                      `tfschema:"run_command_enabled"`
	ServiceMeshProfile         []ServiceMeshProfileModel `tfschema:"service_mesh"`
	StorageProfile             []StorageProfileModel     `tfschema:"storage"`
	SupportPlan                string                    `tfschema:"support_plan"`
	Tags                       map[string]interface{}    `tfschema:"tags"`
	UpgradeOverride            []UpgradeOverrideModel    `tfschema:"upgrade_override"`
	WebAppRouting              []WebAppRoutingModel      `tfschema:"web_app_routing"`
	WindowsProfile             []WindowsProfileModel     `tfschema:"windows_profile"`
	AIToolchainOperatorEnabled bool                      `tfschema:"ai_toolchain_operator_enabled"`

	// Addon fields
	ACIConnectorLinux             []ACIConnectorLinuxModel         `tfschema:"aci_connector_linux"`
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

type APIServerAccessProfileModel struct {
	AuthorizedIPRanges []string `tfschema:"authorized_ip_ranges"`
	SubnetID           string   `tfschema:"subnet_id"`
}

type PrivateClusterModel struct {
	PrivateClusterPublicFQDNEnabled bool   `tfschema:"public_fully_qualified_domain_name_enabled"`
	PrivateDNSZoneID                string `tfschema:"private_dns_zone_id"`
}

type AutoScalerProfileModel struct {
	BalanceSimilarNodeGroups                 bool    `tfschema:"balance_similar_node_groups_enabled"`
	DaemonsetEvictionForEmptyNodesEnabled    bool    `tfschema:"daemonset_eviction_for_empty_nodes_enabled"`
	DaemonsetEvictionForOccupiedNodesEnabled bool    `tfschema:"daemonset_eviction_for_occupied_nodes_enabled"`
	Expander                                 string  `tfschema:"expander"`
	IgnoreDaemonsetsUtilizationEnabled       bool    `tfschema:"daemonset_ignore_utilization_enabled"`
	MaxGracefulTerminationSec                string  `tfschema:"maximum_graceful_termination_in_seconds"`
	MaxNodeProvisioningTime                  string  `tfschema:"maximum_node_provisioning_time"`
	MaxUnreadyNodes                          int64   `tfschema:"maximum_unready_nodes"`
	MaxUnreadyPercentage                     float64 `tfschema:"maximum_unready_percentage"`
	NewPodScaleUpDelay                       string  `tfschema:"new_pod_scale_up_delay"`
	ScanInterval                             string  `tfschema:"scan_interval"`
	ScaleDownDelayAfterAdd                   string  `tfschema:"scale_down_delay_after_add"`
	ScaleDownDelayAfterDelete                string  `tfschema:"scale_down_delay_after_delete"`
	ScaleDownDelayAfterFailure               string  `tfschema:"scale_down_delay_after_failure"`
	ScaleDownUnneeded                        string  `tfschema:"scale_down_unneeded"`
	ScaleDownUnready                         string  `tfschema:"scale_down_unready"`
	ScaleDownUtilizationThreshold            string  `tfschema:"scale_down_utilization_threshold"`
	EmptyBulkDeleteMax                       string  `tfschema:"maximum_empty_bulk_delete"`
	SkipNodesWithLocalStorage                bool    `tfschema:"skip_nodes_with_local_storage_enabled"`
	SkipNodesWithSystemPods                  bool    `tfschema:"skip_nodes_with_system_pods_enabled"`
}

type AzureActiveDirectoryRBACModel struct {
	TenantID            string   `tfschema:"tenant_id"`
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
	DNSServiceIP string   `tfschema:"dns_service_ip"`
	PodCIDR      string   `tfschema:"pod_cidr"`
	PodCIDRs     []string `tfschema:"pod_cidrs"`
	ServiceCIDR  string   `tfschema:"service_cidr"`
	ServiceCIDRs []string `tfschema:"service_cidrs"`
	// IPVersions          []string                   `tfschema:"ip_versions"`
	OutboundType        string                     `tfschema:"outbound_type"`
	LoadBalancerSKU     string                     `tfschema:"load_balancer_sku"`
	LoadBalancerProfile []LoadBalancerProfileModel `tfschema:"load_balancer"`
	NATGatewayProfile   []NATGatewayProfileModel   `tfschema:"nat_gateway"`
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

type ServiceMeshProfileModel struct {
	Mode                          string                      `tfschema:"mode"`
	Revisions                     []string                    `tfschema:"revisions"`
	InternalIngressGatewayEnabled bool                        `tfschema:"internal_ingress_gateway_enabled"`
	ExternalIngressGatewayEnabled bool                        `tfschema:"external_ingress_gateway_enabled"`
	CertificateAuthority          []CertificateAuthorityModel `tfschema:"certificate_authority"`
}

type CertificateAuthorityModel struct {
	KeyVaultID          string `tfschema:"key_vault_id"`
	RootCertObjectName  string `tfschema:"root_certificate_object_name"`
	CertObjectName      string `tfschema:"certificate_object_name"`
	CertChainObjectName string `tfschema:"certificate_chain_object_name"`
	KeyObjectName       string `tfschema:"key_object_name"`
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
	_ sdk.ResourceWithUpdate         = KubernetesAutomaticClusterResource{}
	_ sdk.ResourceWithIdentity       = KubernetesAutomaticClusterResource{}
	_ sdk.ResourceWithCustomImporter = KubernetesAutomaticClusterResource{}
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

func (r KubernetesAutomaticClusterResource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := commonids.ParseKubernetesClusterID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		client := metadata.Client.Containers.KubernetesClustersClient
		resp, err := client.Get(ctx, *id)
		if err != nil || resp.Model == nil || resp.Model.Sku == nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if pointer.From(resp.Model.Sku.Name) != managedclusters.ManagedClusterSKUNameAutomatic {
			return fmt.Errorf("importing %s: specified Kubernetes Cluster is not using the SKU `Automatic`, got `%s`", id, *resp.Model.Sku.Name)
		}

		return nil
	}
}

func (r KubernetesAutomaticClusterResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if rd.Id() != "" {
				// The behaviour of the API requires this, but this could be removed when https://github.com/Azure/azure-rest-api-specs/issues/27373 has been addressed
				// Check default_node_pool upgrade_settings drain_timeout_in_minutes
				if rd.HasChange("default_node_pool.0.upgrade_settings.0.drain_timeout_in_minutes") {
					old, new := rd.GetChange("default_node_pool.0.upgrade_settings.0.drain_timeout_in_minutes")
					if old.(int) != 0 && new.(int) == 0 {
						if err := metadata.ResourceDiff.ForceNew("default_node_pool.0.upgrade_settings.0.drain_timeout_in_minutes"); err != nil {
							return err
						}
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
							if err := metadata.ResourceDiff.ForceNew("default_node_pool.0.name"); err != nil {
								return err
							}
						}
					} else {
						if err := metadata.ResourceDiff.ForceNew("default_node_pool.0.name"); err != nil {
							return err
						}
					}
				}

				// Check windows_profile gmsa changes
				if rd.HasChange("windows_profile.0.gmsa") {
					old, new := rd.GetChange("windows_profile.0.gmsa")
					oldList := old.([]interface{})
					newList := new.([]interface{})
					if len(oldList) > 0 && len(newList) == 0 {
						if err := metadata.ResourceDiff.ForceNew("windows_profile.gmsa"); err != nil {
							return err
						}
					}
				}

				if rd.HasChange("windows_profile.0.gmsa.0.dns_server") {
					old, new := rd.GetChange("windows_profile.0.gmsa.0.dns_server")
					if old.(string) != "" && new.(string) == "" {
						if err := metadata.ResourceDiff.ForceNew("windows_profile.gmsa.dns_server"); err != nil {
							return err
						}
					}
				}

				if rd.HasChange("windows_profile.0.gmsa.0.root_domain") {
					old, new := rd.GetChange("windows_profile.0.gmsa.0.root_domain")
					if old.(string) != "" && new.(string) == "" {
						if err := metadata.ResourceDiff.ForceNew("windows_profile.gmsa.root_domain"); err != nil {
							return err
						}
					}
				}

				// Check api_server_access subnet_id changes
				if rd.HasChange("api_server_access.0.subnet_id") {
					old, new := rd.GetChange("api_server_access.0.subnet_id")
					if old.(string) != "" && new.(string) == "" {
						if err := metadata.ResourceDiff.ForceNew("api_server_access.0.subnet_id"); err != nil {
							return err
						}
					}
				}
			}

			// Validate outbound_type and bootstrap artifact_source
			outboundType := rd.Get("network.0.outbound_type").(string)
			artifactSource := rd.Get("bootstrap.0.artifact_source").(string)
			if outboundType == string(managedclusters.OutboundTypeNone) && artifactSource != string(managedclusters.ArtifactSourceCache) {
				return fmt.Errorf("when `network.outbound_type` is set to `none`, `bootstrap.artifact_source` must be set to `Cache`")
			}

			return nil
		},
	}
}

func (r KubernetesAutomaticClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateKubernetesClusterID
}

func (r KubernetesAutomaticClusterResource) Arguments() map[string]*pluginsdk.Schema {
	arguments := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 63),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"default_node_pool": SchemaDefaultAutomaticClusterNodePoolTyped(),

		"ai_toolchain_operator_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"cost_analysis_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"custom_ca_trust_certificates_base64": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 10,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsBase64,
			},
		},

		"disk_encryption_set_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: computeValidate.DiskEncryptionSetID,
		},

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

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"image_cleaner_interval_in_hours": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      168,
			ValidateFunc: validation.IntBetween(168, 2160),
		},

		"kubernetes_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C - Azure may update the version automatically based on upgrade policies
			Computed:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"node_resource_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: resourcegroups.ValidateName,
		},

		"private_cluster": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"public_fully_qualified_domain_name_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"private_dns_zone_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "System",
						ForceNew: true,
						ValidateFunc: validation.Any(
							privatezones.ValidatePrivateDnsZoneID,
							validation.StringInSlice([]string{
								"System",
								"None",
							}, false),
						),
					},
				},
			},
		},

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

		"tags": commonschema.Tags(),

		"api_server_access": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authorized_ip_ranges": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						AtLeastOneOf: []string{
							"api_server_access.0.authorized_ip_ranges",
							"api_server_access.0.subnet_id",
						},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.CIDR,
						},
					},
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"api_server_access.0.authorized_ip_ranges",
							"api_server_access.0.subnet_id",
						},
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"auto_scaler": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"balance_similar_node_groups_enabled": {
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
					"daemonset_ignore_utilization_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"maximum_graceful_termination_in_seconds": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "600",
					},
					"maximum_node_provisioning_time": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "15m",
						ValidateFunc: containerValidate.Duration,
					},
					"maximum_unready_nodes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      3,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"maximum_unready_percentage": {
						Type:         pluginsdk.TypeFloat,
						Optional:     true,
						Default:      45,
						ValidateFunc: validation.FloatBetween(0, 100),
					},
					"new_pod_scale_up_delay": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "0s",
						ValidateFunc: containerValidate.Duration,
					},
					"scan_interval_in_seconds": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "10",
						ValidateFunc: containerValidate.Duration,
					},
					"scale_down_delay_after_add": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "10m",
						ValidateFunc: containerValidate.Duration,
					},
					"scale_down_delay_after_delete": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "3m",
						ValidateFunc: containerValidate.Duration,
					},
					"scale_down_delay_after_failure": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "3m",
						ValidateFunc: containerValidate.Duration,
					},
					"scale_down_unneeded": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "10m",
						ValidateFunc: containerValidate.Duration,
					},
					"scale_down_unready": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "20m",
						ValidateFunc: containerValidate.Duration,
					},
					"scale_down_utilization_threshold": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "0.5",
					},
					"maximum_empty_bulk_delete": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "10",
					},
					"skip_nodes_with_local_storage_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"skip_nodes_with_system_pods_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
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

		"bootstrap": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"artifact_source": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      managedclusters.ArtifactSourceDirect,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForArtifactSource(), false),
					},
					"container_registry_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: registries.ValidateRegistryID,
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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
					"https_proxy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
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

		"key_management_service": {
			Type:     pluginsdk.TypeList,
			Optional: true,
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

		"network": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dns_service_ip": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validate.IPv4Address,
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
					"load_balancer": {
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
									ConflictsWith: []string{"network.0.load_balancer.0.outbound_ip_prefix_ids", "network.0.load_balancer.0.outbound_ip_address_ids"},
								},
								"outbound_ip_prefix_ids": {
									Type:          pluginsdk.TypeSet,
									Optional:      true,
									ConflictsWith: []string{"network.0.load_balancer.0.managed_outbound_ip_count", "network.0.load_balancer.0.outbound_ip_address_ids"},
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: azure.ValidateResourceID,
									},
								},
								"outbound_ip_address_ids": {
									Type:          pluginsdk.TypeSet,
									Optional:      true,
									ConflictsWith: []string{"network.0.load_balancer.0.managed_outbound_ip_count", "network.0.load_balancer.0.outbound_ip_prefix_ids"},
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
					"nat_gateway": {
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
									AtLeastOneOf: []string{"network.0.advanced_networking.0.observability_enabled", "network.0.advanced_networking.0.security_enabled"},
								},
								"security_enabled": {
									Type:         pluginsdk.TypeBool,
									Optional:     true,
									Default:      false,
									AtLeastOneOf: []string{"network.0.advanced_networking.0.observability_enabled", "network.0.advanced_networking.0.security_enabled"},
								},
							},
						},
					},
				},
			},
		},

		"service_mesh": {
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
								"root_certificate_object_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"certificate_chain_object_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"certificate_object_name": {
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

		"storage": {
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

		"web_app_routing": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dns_zone_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
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

		"windows_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
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
						ValidateFunc: validation.StringLenBetween(14, 123),
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
	}

	for k, v := range schemaKubernetesAutomaticClusterAddOnsTyped() {
		arguments[k] = v
	}

	return arguments
}

func (r KubernetesAutomaticClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"current_kubernetes_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"fqdn": {
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

		"node_resource_group_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"oidc_issuer_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"portal_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
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

			azureMonitorProfile := expandKubernetesAutomaticClusterAzureMonitorProfile(model.MonitorMetrics)

			httpProxyConfig := expandKubernetesAutomaticClusterHttpProxyConfig(model.HTTPProxyConfig)

			apiAccessProfile := expandKubernetesAutomaticClusterAPIAccessProfile(model)

			storageProfile := expandKubernetesAutomaticClusterStorageProfile(model.StorageProfile)

			upgradeOverrideSetting := expandKubernetesAutomaticClusterUpgradeOverride(model.UpgradeOverride)

			securityProfile := &managedclusters.ManagedClusterSecurityProfile{}

			securityProfile.Defender = expandKubernetesAutomaticClusterMicrosoftDefender(model.MicrosoftDefender, false)

			securityProfile.ImageCleaner = &managedclusters.ManagedClusterSecurityProfileImageCleaner{
				Enabled:       pointer.To(true),
				IntervalHours: pointer.To(model.ImageCleanerIntervalHours),
			}

			securityProfile.AzureKeyVaultKms, err = expandKubernetesAutomaticClusterKeyManagementService(model.KeyManagementService, ctx, keyVaultsClient, subscriptionId)
			if err != nil {
				return err
			}

			if len(model.CustomCATrustCertificatesBase64) > 0 {
				securityProfile.CustomCATrustCertificates = pointer.To(model.CustomCATrustCertificatesBase64)
			}

			metricsProfile := expandKubernetesAutomaticClusterMetricsProfile(model.CostAnalysisEnabled)

			ingressProfile := expandKubernetesAutomaticClusterWebAppRouting(model.WebAppRouting, false)

			serviceMeshProfile := expandKubernetesAutomaticClusterServiceMeshProfile(model.ServiceMeshProfile, nil)

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

			addonProfiles, err := expandKubernetesAutomaticClusterAddOns(&model, metadata.Client.Containers.Environment)
			if err != nil {
				return fmt.Errorf("expanding addons: %+v", err)
			}

			(*addonProfiles)["azurepolicy"] = managedclusters.ManagedClusterAddonProfile{
				Enabled: true,
			}

			var azureADProfile *managedclusters.ManagedClusterAADProfile
			if len(model.AzureActiveDirectoryRBAC) > 0 {
				azureADProfile = &managedclusters.ManagedClusterAADProfile{
					Managed:             pointer.To(true),
					AdminGroupObjectIDs: &model.AzureActiveDirectoryRBAC[0].AdminGroupObjectIDs,
				}
			}

			parameters := managedclusters.ManagedCluster{
				Location: location,
				Sku: &managedclusters.ManagedClusterSKU{
					Name: pointer.To(managedclusters.ManagedClusterSKUName("automatic")),
					Tier: pointer.To(managedclusters.ManagedClusterSKUTier("standard")),
				},
				Properties: &managedclusters.ManagedClusterProperties{
					ApiServerAccessProfile: apiAccessProfile,
					AadProfile:             azureADProfile,
					AddonProfiles:          addonProfiles,
					AgentPoolProfiles:      agentProfiles,
					AutoScalerProfile:      autoScalerProfile,
					AutoUpgradeProfile:     pointer.To(managedclusters.ManagedClusterAutoUpgradeProfile{}),
					AzureMonitorProfile:    azureMonitorProfile,
					KubernetesVersion:      pointer.To(kubernetesVersion),
					BootstrapProfile:       bootstrapProfile,
					LinuxProfile:           linuxProfile,
					WindowsProfile:         windowsProfile,
					MetricsProfile:         metricsProfile,
					NetworkProfile:         networkProfile,
					NodeResourceGroup:      pointer.To(model.NodeResourceGroup),
					HTTPProxyConfig:        httpProxyConfig,
					SecurityProfile:        securityProfile,
					StorageProfile:         storageProfile,
					UpgradeSettings:        upgradeOverrideSetting,
					IngressProfile:         ingressProfile,
					ServiceMeshProfile:     serviceMeshProfile,
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
			}

			if len(model.KubeletIdentity) > 0 {
				parameters.Properties.IdentityProfile = expandKubernetesAutomaticClusterIdentityProfile(model.KubeletIdentity)
			}

			if len(model.PrivateCluster) > 0 && model.PrivateCluster[0].PrivateDNSZoneID != "" {
				privateDNSZoneID := model.PrivateCluster[0].PrivateDNSZoneID
				if (parameters.Identity == nil) || (privateDNSZoneID != "System" && privateDNSZoneID != "None" && (parameters.Identity.Type != identity.TypeUserAssigned)) {
					return fmt.Errorf("a user assigned identity must be used when using a custom private dns zone")
				}
			}

			if model.DNSPrefixPrivateCluster != "" {
				if apiAccessProfile.PrivateDNSZone == nil || *apiAccessProfile.PrivateDNSZone == "System" || *apiAccessProfile.PrivateDNSZone == "None" {
					return fmt.Errorf("`dns_prefix_private_cluster` should only be set for private cluster with custom private dns zone")
				}
				parameters.Properties.FqdnSubdomain = pointer.To(model.DNSPrefixPrivateCluster)
			} else {
				parameters.Properties.DnsPrefix = pointer.To(model.DNSPrefix)
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

			//if len(model.MaintenanceWindow) > 0 {
			//	maintenanceClient := metadata.Client.Containers.MaintenanceConfigurationsClient
			//	maintenanceParams := maintenanceconfigurations.MaintenanceConfiguration{
			//		Properties: expandKubernetesAutomaticClusterMaintenanceConfigurationDefault(model.MaintenanceWindow),
			//	}
			//	maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
			//	if _, err := maintenanceClient.CreateOrUpdate(ctx, maintenanceId, maintenanceParams); err != nil {
			//		return fmt.Errorf("creating/updating default maintenance config for %s: %+v", id, err)
			//	}
			//}

			//if len(model.MaintenanceWindowAutoUpgrade) > 0 {
			//	maintenanceClient := metadata.Client.Containers.MaintenanceConfigurationsClient
			//	maintenanceParams := maintenanceconfigurations.MaintenanceConfiguration{
			//		Properties: expandKubernetesAutomaticClusterMaintenanceConfigurationAutoUpgrade(model.MaintenanceWindowAutoUpgrade),
			//	}
			//	maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
			//	if _, err := maintenanceClient.CreateOrUpdate(ctx, maintenanceId, maintenanceParams); err != nil {
			//		return fmt.Errorf("creating/updating auto upgrade schedule maintenance config for %s: %+v", id, err)
			//	}
			//}

			//if len(model.MaintenanceWindowNodeOS) > 0 {
			//	maintenanceClient := metadata.Client.Containers.MaintenanceConfigurationsClient
			//	maintenanceParams := maintenanceconfigurations.MaintenanceConfiguration{
			//		Properties: expandKubernetesAutomaticClusterMaintenanceConfigurationNodeOS(model.MaintenanceWindowNodeOS),
			//	}
			//	maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
			//	if _, err := maintenanceClient.CreateOrUpdate(ctx, maintenanceId, maintenanceParams); err != nil {
			//		return fmt.Errorf("creating/updating node os upgrade schedule maintenance config for %s: %+v", id, err)
			//	}
			//}

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

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
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
		return fmt.Errorf("decoding: %+v", err)
	}

	state := KubernetesAutomaticClusterModel{
		Name:              id.ManagedClusterName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)

		state.Tags = tags.Flatten(model.Tags)

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

			if props.SecurityProfile != nil && props.SecurityProfile.CustomCATrustCertificates != nil {
				state.CustomCATrustCertificatesBase64 = *props.SecurityProfile.CustomCATrustCertificates
			}

			apiServerAccessProfile, privateCluster, runCommandEnabled := flattenKubernetesAutomaticClusterAPIAccessProfile(props.ApiServerAccessProfile)
			state.APIServerAccessProfile = apiServerAccessProfile
			state.PrivateCluster = privateCluster
			state.RunCommandEnabled = runCommandEnabled

			if props.AddonProfiles != nil {
				state.ACIConnectorLinux,
					state.ConfidentialComputing,
					state.HTTPApplicationRoutingEnabled,
					state.HTTPApplicationRoutingZoneName,
					state.IngressApplicationGateway,
					state.KeyVaultSecretsProvider,
					state.OMSAgent,
					state.OpenServiceMeshEnabled = flattenKubernetesAutomaticClusterAddOns(*props.AddonProfiles)
			}

			state.AutoScalerProfile, err = flattenKubernetesAutomaticClusterAutoScalerProfile(props.AutoScalerProfile)
			if err != nil {
				return fmt.Errorf("flattening `auto_scaler`: %+v", err)
			}

			state.MonitorMetrics = flattenKubernetesAutomaticClusterAzureMonitorProfile(props.AzureMonitorProfile)

			state.ServiceMeshProfile = flattenKubernetesAutomaticClusterServiceMeshProfile(props.ServiceMeshProfile)

			if props.AgentPoolProfiles != nil {
				state.DefaultNodePool, err = FlattenDefaultNodePoolTyped(props.AgentPoolProfiles, &metadata)
				if err != nil {
					return fmt.Errorf("flattening `default_node_pool`: %+v", err)
				}
			}

			kubeletIdentity, err := flattenKubernetesAutomaticClusterIdentityProfile(pointer.From(props.IdentityProfile))
			if err != nil {
				return fmt.Errorf("flattening `kubelet_identity`: %+v", err)
			}
			state.KubeletIdentity = kubeletIdentity

			state.LinuxProfile = flattenKubernetesAutomaticClusterLinuxProfile(props.LinuxProfile)

			state.NetworkProfile = flattenKubernetesAutomaticClusterNetworkProfile(props.NetworkProfile)

			state.WindowsProfile = flattenKubernetesAutomaticClusterWindowsProfile(props.WindowsProfile, config)

			state.HTTPProxyConfig = flattenKubernetesAutomaticClusterHttpProxyConfig(props.HTTPProxyConfig)

			state.BootstrapProfile = flattenKubernetesAutomaticClusterBootstrapProfile(props.BootstrapProfile)

			state.UpgradeOverride = flattenKubernetesAutomaticClusterUpgradeOverride(props.UpgradeSettings)

			state.StorageProfile = flattenKubernetesAutomaticClusterStorageProfile(props.StorageProfile)

			state.WebAppRouting = flattenKubernetesAutomaticClusterWebAppRouting(props.IngressProfile)

			state.MicrosoftDefender = flattenKubernetesAutomaticClusterMicrosoftDefender(props.SecurityProfile)

			if props.SecurityProfile != nil && props.SecurityProfile.AzureKeyVaultKms != nil {
				state.KeyManagementService = flattenKubernetesAutomaticClusterKeyManagementService(props.SecurityProfile.AzureKeyVaultKms)
			}

			state.CostAnalysisEnabled = flattenKubernetesAutomaticClusterMetricsProfile(props.MetricsProfile)

			state.AzureActiveDirectoryRBAC = flattenKubernetesAutomaticClusterAzureActiveDirectoryRBAC(props.AadProfile)

			if props.SecurityProfile != nil && props.SecurityProfile.ImageCleaner != nil {
				state.ImageCleanerIntervalHours = pointer.From(props.SecurityProfile.ImageCleaner.IntervalHours)
			}

			aiToolchainOperatorEnabled := false
			if props.AiToolchainOperatorProfile != nil {
				aiToolchainOperatorEnabled = pointer.From(props.AiToolchainOperatorProfile.Enabled)
			}
			state.AIToolchainOperatorEnabled = aiToolchainOperatorEnabled

			state.SupportPlan = string(pointer.From(props.SupportPlan))
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

		//maintenanceClient := metadata.Client.Containers.MaintenanceConfigurationsClient
		//
		//maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
		//configResp, _ := maintenanceClient.Get(ctx, maintenanceId)
		//if configurationBody := configResp.Model; configurationBody != nil && configurationBody.Properties != nil {
		//	state.MaintenanceWindow = flattenKubernetesAutomaticClusterMaintenanceConfigurationDefault(configurationBody.Properties)
		//}
		//
		//maintenanceId = maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
		//configResp, _ = maintenanceClient.Get(ctx, maintenanceId)
		//if configurationBody := configResp.Model; configurationBody != nil && configurationBody.Properties != nil && configurationBody.Properties.MaintenanceWindow != nil {
		//	state.MaintenanceWindowAutoUpgrade = flattenKubernetesAutomaticClusterMaintenanceConfiguration(configurationBody.Properties.MaintenanceWindow)
		//}
		//
		//maintenanceId = maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
		//configResp, _ = maintenanceClient.Get(ctx, maintenanceId)
		//if configurationBody := configResp.Model; configurationBody != nil && configurationBody.Properties != nil && configurationBody.Properties.MaintenanceWindow != nil {
		//	autoUpgradeConfig := flattenKubernetesAutomaticClusterMaintenanceConfiguration(configurationBody.Properties.MaintenanceWindow)
		//	if len(autoUpgradeConfig) > 0 {
		//		au := autoUpgradeConfig[0]
		//		state.MaintenanceWindowNodeOS = []MaintenanceWindowNodeOSModel{{
		//			Frequency:  au.Frequency,
		//			Interval:   au.Interval,
		//			DayOfWeek:  au.DayOfWeek,
		//			Duration:   au.Duration,
		//			WeekIndex:  au.WeekIndex,
		//			DayOfMonth: au.DayOfMonth,
		//			StartDate:  au.StartDate,
		//			StartTime:  au.StartTime,
		//			UTCOffset:  au.UTCOffset,
		//			NotAllowed: au.NotAllowed,
		//		}}
		//	}
		//}
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

			updateCluster := false
			props := existing.Model.Properties

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = tags.Expand(model.Tags)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("kubernetes_version") {
				props.KubernetesVersion = pointer.To(model.KubernetesVersion)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("auto_scaler") {
				props.AutoScalerProfile = expandKubernetesAutomaticClusterAutoScalerProfile(model.AutoScalerProfile)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("monitor_metrics") {
				props.AzureMonitorProfile = expandKubernetesAutomaticClusterAzureMonitorProfile(model.MonitorMetrics)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("service_mesh") {
				props.ServiceMeshProfile = expandKubernetesAutomaticClusterServiceMeshProfile(model.ServiceMeshProfile, props.ServiceMeshProfile)
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
						return fmt.Errorf("retrieving Default Node Pool %s: %w", defaultNodePoolId, err)
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
					"default_node_pool.0.maximum_pods",
					"default_node_pool.0.os_disk_size_gb",
					"default_node_pool.0.pod_subnet_id",
					"default_node_pool.0.snapshot_id",
					"default_node_pool.0.ultra_ssd_enabled",
					"default_node_pool.0.vnet_subnet_id",
					"default_node_pool.0.vm_size",
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
						return fmt.Errorf("checking for existing temporary %s: %w", tempNodePoolId, err)
					}

					defaultExisting, err := nodePoolsClient.Get(ctx, defaultNodePoolId)
					if !response.WasNotFound(defaultExisting.HttpResponse) && err != nil {
						return fmt.Errorf("checking for existing default %s: %w", defaultNodePoolId, err)
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
							return fmt.Errorf("creating temporary %s: %w", tempNodePoolId, err)
						}
					}

					// delete the old default node pool if it exists
					if defaultExisting.Model != nil {
						if err := nodePoolsClient.DeleteThenPoll(ctx, defaultNodePoolId, agentpools.DefaultDeleteOperationOptions()); err != nil {
							return fmt.Errorf("deleting default %s: %w", defaultNodePoolId, err)
						}
					}

					// create the default node pool with the new vm size
					if err := retryNodePoolCreation(ctx, nodePoolsClient, defaultNodePoolId, agentProfile); err != nil {
						// if creation of the default node pool fails we automatically fall back to the temporary node pool
						// in func findDefaultNodePool
						return fmt.Errorf("creating default %s: %w", defaultNodePoolId, err)
					}

					if err := nodePoolsClient.DeleteThenPoll(ctx, tempNodePoolId, agentpools.DefaultDeleteOperationOptions()); err != nil {
						return fmt.Errorf("deleting temporary %s: %w", tempNodePoolId, err)
					}
				} else {
					if err := nodePoolsClient.CreateOrUpdateThenPoll(ctx, defaultNodePoolId, agentProfile, agentpools.DefaultCreateOrUpdateOperationOptions()); err != nil {
						return fmt.Errorf("updating Default Node Pool %s: %w", defaultNodePoolId, err)
					}
				}
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("aci_connector_linux") ||
				metadata.ResourceData.HasChange("confidential_computing") ||
				metadata.ResourceData.HasChange("http_application_routing_enabled") ||
				metadata.ResourceData.HasChange("oms_agent") ||
				metadata.ResourceData.HasChange("ingress_application_gateway") ||
				metadata.ResourceData.HasChange("key_vault_secrets_provider") ||
				metadata.ResourceData.HasChange("open_service_mesh_enabled") {
				addonProfiles, err := expandKubernetesAutomaticClusterAddOns(&model, metadata.Client.Containers.Environment)
				if err != nil {
					return fmt.Errorf("expanding addons: %w", err)
				}
				props.AddonProfiles = addonProfiles
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("api_server_access") ||
				metadata.ResourceData.HasChange("private_cluster") ||
				metadata.ResourceData.HasChange("run_command_enabled") {
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
						return fmt.Errorf("updating AAD Profile for %s: %w", *id, err)
					}
				}

				if props.AadProfile != nil && props.AadProfile.Managed != nil && *props.AadProfile.Managed {
					updateCluster = true
				}
			}

			if metadata.ResourceData.HasChange("network") {
				networkProfile, err := expandKubernetesAutomaticClusterNetworkProfile(model.NetworkProfile)
				if err != nil {
					return fmt.Errorf("expanding network profile: %w", err)
				}
				props.NetworkProfile = networkProfile
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("http_proxy_config") {
				props.HTTPProxyConfig = expandKubernetesAutomaticClusterHttpProxyConfig(model.HTTPProxyConfig)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("storage") {
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

			if metadata.ResourceData.HasChange("bootstrap") {
				props.BootstrapProfile = expandKubernetesAutomaticClusterBootstrapProfile(model.BootstrapProfile)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("microsoft_defender") ||
				metadata.ResourceData.HasChange("image_cleaner_interval_in_hours") ||
				metadata.ResourceData.HasChange("key_management_service") ||
				metadata.ResourceData.HasChange("custom_ca_trust_certificates_base64") {
				if props.SecurityProfile == nil {
					props.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
				}

				if metadata.ResourceData.HasChange("microsoft_defender") {
					props.SecurityProfile.Defender = expandKubernetesAutomaticClusterMicrosoftDefender(model.MicrosoftDefender, metadata.ResourceData.HasChange("microsoft_defender"))
				}

				if metadata.ResourceData.HasChange("image_cleaner_interval_in_hours") {
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
					return fmt.Errorf("updating %s: %w", *id, err)
				}
			}

			// maintenanceClient := metadata.Client.Containers.MaintenanceConfigurationsClient

			//if metadata.ResourceData.HasChange("maintenance_window") {
			//	maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
			//	if len(model.MaintenanceWindow) > 0 {
			//		parameters := maintenanceconfigurations.MaintenanceConfiguration{
			//			Properties: expandKubernetesAutomaticClusterMaintenanceConfigurationDefault(model.MaintenanceWindow),
			//		}
			//		if _, err := maintenanceClient.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
			//			return fmt.Errorf("updating default maintenance config for %s: %w", id, err)
			//		}
			//	} else {
			//		if _, err := maintenanceClient.Delete(ctx, maintenanceId); err != nil {
			//			return fmt.Errorf("deleting default maintenance config for %s: %w", id, err)
			//		}
			//	}
			//}
			//
			//if metadata.ResourceData.HasChange("maintenance_window_auto_upgrade") {
			//	maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
			//	if len(model.MaintenanceWindowAutoUpgrade) > 0 {
			//		parameters := maintenanceconfigurations.MaintenanceConfiguration{
			//			Properties: expandKubernetesAutomaticClusterMaintenanceConfigurationAutoUpgrade(model.MaintenanceWindowAutoUpgrade),
			//		}
			//		if _, err := maintenanceClient.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
			//			return fmt.Errorf("updating auto upgrade maintenance config for %s: %w", id, err)
			//		}
			//	} else {
			//		if _, err := maintenanceClient.Delete(ctx, maintenanceId); err != nil {
			//			return fmt.Errorf("deleting auto upgrade maintenance config for %s: %w", id, err)
			//		}
			//	}
			//}
			//
			//if metadata.ResourceData.HasChange("maintenance_window_node_os") {
			//	maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
			//	if len(model.MaintenanceWindowNodeOS) > 0 {
			//		parameters := maintenanceconfigurations.MaintenanceConfiguration{
			//			Properties: expandKubernetesAutomaticClusterMaintenanceConfigurationNodeOS(model.MaintenanceWindowNodeOS),
			//		}
			//		if _, err := maintenanceClient.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
			//			return fmt.Errorf("updating node os maintenance config for %s: %w", id, err)
			//		}
			//	} else {
			//		if _, err := maintenanceClient.Delete(ctx, maintenanceId); err != nil {
			//			return fmt.Errorf("deleting node os maintenance config for %s: %w", id, err)
			//		}
			//	}
			//}

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
				return fmt.Errorf("decoding: %w", err)
			}

			//maintenanceClient := metadata.Client.Containers.MaintenanceConfigurationsClient
			//
			//if len(model.MaintenanceWindow) > 0 {
			//	maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
			//	if _, err := maintenanceClient.Delete(ctx, maintenanceId); err != nil {
			//		return fmt.Errorf("deleting default maintenance configuration for %s: %w", *id, err)
			//	}
			//}
			//
			//if len(model.MaintenanceWindowAutoUpgrade) > 0 {
			//	maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
			//	if _, err := maintenanceClient.Delete(ctx, maintenanceId); err != nil {
			//		return fmt.Errorf("deleting auto-upgrade maintenance configuration for %s: %w", *id, err)
			//	}
			//}
			//
			//if len(model.MaintenanceWindowNodeOS) > 0 {
			//	maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
			//	if _, err := maintenanceClient.Delete(ctx, maintenanceId); err != nil {
			//		return fmt.Errorf("deleting node OS maintenance configuration for %s: %w", *id, err)
			//	}
			//}

			if err := client.DeleteThenPoll(ctx, *id, managedclusters.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %w", *id, err)
			}

			return nil
		},
	}
}

func expandKubernetesAutomaticClusterAPIAccessProfile(model KubernetesAutomaticClusterModel) *managedclusters.ManagedClusterAPIServerAccessProfile {
	enablePrivateCluster, enablePrivateClusterPublicFQDN, privateDNSZoneID := expandKubernetesAutomaticClusterPrivateCluster(model.PrivateCluster)

	apiAccessProfile := &managedclusters.ManagedClusterAPIServerAccessProfile{
		EnablePrivateCluster:           pointer.To(enablePrivateCluster),
		EnablePrivateClusterPublicFQDN: pointer.To(enablePrivateClusterPublicFQDN),
		DisableRunCommand:              pointer.To(!model.RunCommandEnabled),
	}

	if privateDNSZoneID != "" {
		apiAccessProfile.PrivateDNSZone = pointer.To(privateDNSZoneID)
	}

	if len(model.APIServerAccessProfile) == 0 {
		return apiAccessProfile
	}

	config := model.APIServerAccessProfile[0]
	apiAccessProfile.AuthorizedIPRanges = pointer.To(config.AuthorizedIPRanges)

	if config.SubnetID != "" {
		apiAccessProfile.SubnetId = pointer.To(config.SubnetID)
	}

	return apiAccessProfile
}

func flattenKubernetesAutomaticClusterAPIAccessProfile(profile *managedclusters.ManagedClusterAPIServerAccessProfile) ([]APIServerAccessProfileModel, []PrivateClusterModel, bool) {
	apiServerAccessProfile := make([]APIServerAccessProfileModel, 0)
	runCommandEnabled := true

	if profile == nil {
		return apiServerAccessProfile, []PrivateClusterModel{}, runCommandEnabled
	}

	// Extract private cluster settings
	enablePrivateCluster := false
	enablePrivateClusterPublicFQDN := false
	privateDNSZoneID := ""

	if profile.EnablePrivateCluster != nil {
		enablePrivateCluster = *profile.EnablePrivateCluster
	}
	if profile.EnablePrivateClusterPublicFQDN != nil {
		enablePrivateClusterPublicFQDN = *profile.EnablePrivateClusterPublicFQDN
	}
	if profile.DisableRunCommand != nil {
		runCommandEnabled = !*profile.DisableRunCommand
	}

	// Handle PrivateDNSZone normalization
	switch {
	case profile.PrivateDNSZone != nil && strings.EqualFold("System", *profile.PrivateDNSZone):
		privateDNSZoneID = "System"
	case profile.PrivateDNSZone != nil && strings.EqualFold("None", *profile.PrivateDNSZone):
		privateDNSZoneID = "None"
	default:
		privateDNSZoneID = pointer.From(profile.PrivateDNSZone)
	}

	privateCluster := flattenKubernetesAutomaticClusterPrivateCluster(enablePrivateCluster, enablePrivateClusterPublicFQDN, privateDNSZoneID)

	// API access profile can be managed by other properties, only return it if one of the properties has been set
	hasAuthorizedIPRanges := profile.AuthorizedIPRanges != nil && len(*profile.AuthorizedIPRanges) > 0
	hasSubnetId := profile.SubnetId != nil && *profile.SubnetId != ""

	if !hasAuthorizedIPRanges && !hasSubnetId {
		return apiServerAccessProfile, privateCluster, runCommandEnabled
	}

	apiServerAccessProfile = append(apiServerAccessProfile, APIServerAccessProfileModel{
		AuthorizedIPRanges: pointer.From(profile.AuthorizedIPRanges),
		SubnetID:           pointer.From(profile.SubnetId),
	})

	return apiServerAccessProfile, privateCluster, runCommandEnabled
}

func expandKubernetesAutomaticClusterPrivateCluster(model []PrivateClusterModel) (bool, bool, string) {
	if len(model) == 0 {
		return false, false, ""
	}

	config := model[0]
	return true, config.PrivateClusterPublicFQDNEnabled, config.PrivateDNSZoneID
}

func flattenKubernetesAutomaticClusterPrivateCluster(enablePrivateCluster bool, enablePrivateClusterPublicFQDN bool, privateDNSZoneID string) []PrivateClusterModel {
	if !enablePrivateCluster {
		return []PrivateClusterModel{}
	}

	return []PrivateClusterModel{
		{
			PrivateClusterPublicFQDNEnabled: enablePrivateClusterPublicFQDN,
			PrivateDNSZoneID:                privateDNSZoneID,
		},
	}
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

	loadBalancerSku := config.LoadBalancerSKU
	outboundType := config.OutboundType

	//ipVersions, err := expandAutomaticIPVersions(config.IPVersions)
	//if err != nil {
	//	return nil, err
	//}

	networkProfile := managedclusters.ContainerServiceNetworkProfile{
		LoadBalancerSku: pointer.To(managedclusters.LoadBalancerSku(loadBalancerSku)),
		OutboundType:    pointer.To(managedclusters.OutboundType(outboundType)),
		IPFamilies:      &[]managedclusters.IPFamily{"IPv4"},
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

//func expandAutomaticIPVersions(input []string) (*[]managedclusters.IPFamily, error) {
//	if len(input) == 0 {
//		return nil, nil
//	}
//
//	ipv := make([]managedclusters.IPFamily, 0)
//	for _, data := range input {
//		ipv = append(ipv, managedclusters.IPFamily(data))
//	}
//
//	if len(ipv) == 1 && ipv[0] == managedclusters.IPFamilyIPvSix {
//		return nil, fmt.Errorf("`ip_versions` must be `IPv4` or `IPv4` and `IPv6`. `IPv6` alone is not supported")
//	}
//
//	return &ipv, nil
//}

func flattenAutomaticLoadBalancerProfile(profile *managedclusters.ManagedClusterLoadBalancerProfile) []LoadBalancerProfileModel {
	if profile == nil {
		return []LoadBalancerProfileModel{}
	}

	result := LoadBalancerProfileModel{}

	if profile.AllocatedOutboundPorts != nil {
		result.OutboundPortsAllocated = pointer.From(profile.AllocatedOutboundPorts)
	}

	if profile.IdleTimeoutInMinutes != nil {
		result.IdleTimeoutInMinutes = pointer.From(profile.IdleTimeoutInMinutes)
	}

	if profile.ManagedOutboundIPs != nil && profile.ManagedOutboundIPs.Count != nil {
		result.ManagedOutboundIPCount = pointer.From(profile.ManagedOutboundIPs.Count)
	}

	if profile.OutboundIPs != nil && profile.OutboundIPs.PublicIPs != nil {
		result.OutboundIPAddressIDs = automaticResourceReferencesToIds(profile.OutboundIPs.PublicIPs)
	}

	if profile.OutboundIPPrefixes != nil && profile.OutboundIPPrefixes.PublicIPPrefixes != nil {
		result.OutboundIPPrefixIDs = automaticResourceReferencesToIds(profile.OutboundIPPrefixes.PublicIPPrefixes)
	}

	if profile.BackendPoolType != nil {
		result.BackendPoolType = string(pointer.From(profile.BackendPoolType))
	}

	result.EffectiveOutboundIPs = automaticResourceReferencesToIds(profile.EffectiveOutboundIPs)

	return []LoadBalancerProfileModel{result}
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

func flattenAutomaticNatGatewayProfile(profile *managedclusters.ManagedClusterNATGatewayProfile) []NATGatewayProfileModel {
	if profile == nil {
		return []NATGatewayProfileModel{}
	}

	result := NATGatewayProfileModel{}

	if profile.IdleTimeoutInMinutes != nil {
		result.IdleTimeoutInMinutes = pointer.From(profile.IdleTimeoutInMinutes)
	}

	if profile.ManagedOutboundIPProfile != nil && profile.ManagedOutboundIPProfile.Count != nil {
		result.ManagedOutboundIPCount = pointer.From(profile.ManagedOutboundIPProfile.Count)
	}

	result.EffectiveOutboundIPs = automaticResourceReferencesToIds(profile.EffectiveOutboundIPs)

	return []NATGatewayProfileModel{result}
}

func flattenKubernetesAutomaticClusterNetworkProfile(profile *managedclusters.ContainerServiceNetworkProfile) []NetworkProfileModel {
	if profile == nil {
		return []NetworkProfileModel{}
	}

	dnsServiceIP := ""
	if profile.DnsServiceIP != nil {
		dnsServiceIP = pointer.From(profile.DnsServiceIP)
	}

	serviceCidr := ""
	if profile.ServiceCidr != nil {
		serviceCidr = pointer.From(profile.ServiceCidr)
	}

	podCidr := ""
	if profile.PodCidr != nil {
		podCidr = pointer.From(profile.PodCidr)
	}

	outboundType := ""
	if profile.OutboundType != nil {
		outboundType = string(pointer.From(profile.OutboundType))
	}

	lbProfiles := flattenAutomaticLoadBalancerProfile(profile.LoadBalancerProfile)

	ngwProfiles := flattenAutomaticNatGatewayProfile(profile.NatGatewayProfile)

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

	advancedNetworking := flattenKubernetesAutomaticClusterAdvancedNetworking(profile.AdvancedNetworking)

	podCidrs := []string{}
	if profile.PodCidrs != nil {
		podCidrs = pointer.From(profile.PodCidrs)
	}

	serviceCidrs := []string{}
	if profile.ServiceCidrs != nil {
		serviceCidrs = pointer.From(profile.ServiceCidrs)
	}

	return []NetworkProfileModel{
		{
			DNSServiceIP:        dnsServiceIP,
			LoadBalancerSKU:     string(pointer.From(sku)),
			LoadBalancerProfile: lbProfiles,
			NATGatewayProfile:   ngwProfiles,
			// IPVersions:          ipVersions,
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
			return nil, fmt.Errorf("parsing BalanceSimilarNodeGroups: %w", err)
		}
		balanceSimilarNodeGroups = b
	}

	expander := ""
	if profile.Expander != nil {
		expander = string(pointer.From(profile.Expander))
	}

	maxGracefulTerminationSec := ""
	if profile.MaxGracefulTerminationSec != nil {
		maxGracefulTerminationSec = pointer.From(profile.MaxGracefulTerminationSec)
	}

	MaxNodeProvisioningTime := ""
	if profile.MaxNodeProvisionTime != nil {
		MaxNodeProvisioningTime = pointer.From(profile.MaxNodeProvisionTime)
	}

	newPodScaleUpDelay := ""
	if profile.NewPodScaleUpDelay != nil {
		newPodScaleUpDelay = pointer.From(profile.NewPodScaleUpDelay)
	}

	scaleDownDelayAfterAdd := ""
	if profile.ScaleDownDelayAfterAdd != nil {
		scaleDownDelayAfterAdd = pointer.From(profile.ScaleDownDelayAfterAdd)
	}

	scaleDownDelayAfterDelete := ""
	if profile.ScaleDownDelayAfterDelete != nil {
		scaleDownDelayAfterDelete = pointer.From(profile.ScaleDownDelayAfterDelete)
	}

	scaleDownDelayAfterFailure := ""
	if profile.ScaleDownDelayAfterFailure != nil {
		scaleDownDelayAfterFailure = pointer.From(profile.ScaleDownDelayAfterFailure)
	}

	scaleDownUnneededTime := ""
	if profile.ScaleDownUnneededTime != nil {
		scaleDownUnneededTime = pointer.From(profile.ScaleDownUnneededTime)
	}

	scaleDownUnreadyTime := ""
	if profile.ScaleDownUnreadyTime != nil {
		scaleDownUnreadyTime = pointer.From(profile.ScaleDownUnreadyTime)
	}

	scaleDownUtilizationThreshold := ""
	if profile.ScaleDownUtilizationThreshold != nil {
		scaleDownUtilizationThreshold = pointer.From(profile.ScaleDownUtilizationThreshold)
	}

	emptyBulkDeleteMax := ""
	if profile.MaxEmptyBulkDelete != nil {
		emptyBulkDeleteMax = pointer.From(profile.MaxEmptyBulkDelete)
	}

	var skipNodesWithLocalStorage bool
	if profile.SkipNodesWithLocalStorage != nil {
		b, err := strconv.ParseBool(pointer.From(profile.SkipNodesWithLocalStorage))
		if err != nil {
			return nil, fmt.Errorf("parsing SkipNodesWithLocalStorage: %w", err)
		}
		skipNodesWithLocalStorage = b
	}

	var skipNodesWithSystemPods bool
	if profile.SkipNodesWithSystemPods != nil {
		b, err := strconv.ParseBool(pointer.From(profile.SkipNodesWithSystemPods))
		if err != nil {
			return nil, fmt.Errorf("parsing SkipNodesWithSystemPods: %w", err)
		}
		skipNodesWithSystemPods = b
	}

	scanInterval := ""
	if profile.ScanInterval != nil {
		scanInterval = pointer.From(profile.ScanInterval)
	}

	var daemonsetEvictionForEmptyNodesEnabled bool
	if profile.DaemonsetEvictionForEmptyNodes != nil {
		daemonsetEvictionForEmptyNodesEnabled = pointer.From(profile.DaemonsetEvictionForEmptyNodes)
	}

	var daemonsetEvictionForOccupiedNodesEnabled bool
	if profile.DaemonsetEvictionForOccupiedNodes != nil {
		daemonsetEvictionForOccupiedNodesEnabled = pointer.From(profile.DaemonsetEvictionForOccupiedNodes)
	}

	var ignoreDaemonsetsUtilizationEnabled bool
	if profile.IgnoreDaemonsetsUtilization != nil {
		ignoreDaemonsetsUtilizationEnabled = pointer.From(profile.IgnoreDaemonsetsUtilization)
	}

	maxUnreadyNodes := 0
	if profile.OkTotalUnreadyCount != nil {
		var err error
		maxUnreadyNodes, err = strconv.Atoi(pointer.From(profile.OkTotalUnreadyCount))
		if err != nil {
			return nil, err
		}
	}

	maxUnreadyPercentage := 0.0
	if profile.MaxTotalUnreadyPercentage != nil {
		var err error
		maxUnreadyPercentage, err = strconv.ParseFloat(pointer.From(profile.MaxTotalUnreadyPercentage), 64)
		if err != nil {
			return nil, err
		}
	}

	return []AutoScalerProfileModel{{
		BalanceSimilarNodeGroups:                 balanceSimilarNodeGroups,
		DaemonsetEvictionForEmptyNodesEnabled:    daemonsetEvictionForEmptyNodesEnabled,
		DaemonsetEvictionForOccupiedNodesEnabled: daemonsetEvictionForOccupiedNodesEnabled,
		Expander:                                 expander,
		IgnoreDaemonsetsUtilizationEnabled:       ignoreDaemonsetsUtilizationEnabled,
		MaxGracefulTerminationSec:                maxGracefulTerminationSec,
		MaxNodeProvisioningTime:                  MaxNodeProvisioningTime,
		MaxUnreadyNodes:                          int64(maxUnreadyNodes),
		MaxUnreadyPercentage:                     maxUnreadyPercentage,
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

//func expandKubernetesAutomaticClusterMaintenanceConfigurationDefault(input []MaintenanceWindowModel) *maintenanceconfigurations.MaintenanceConfigurationProperties {
//	if len(input) == 0 {
//		return nil
//	}
//	value := input[0]
//	return &maintenanceconfigurations.MaintenanceConfigurationProperties{
//		NotAllowedTime: expandKubernetesAutomaticClusterMaintenanceConfigurationTimeSpans(value.NotAllowed),
//		TimeInWeek:     expandKubernetesAutomaticClusterMaintenanceConfigurationTimeInWeeks(value.Allowed),
//	}
//}

//
//func expandKubernetesAutomaticClusterMaintenanceConfigurationForCreate(input []MaintenanceWindowAutoUpgradeModel) *maintenanceconfigurations.MaintenanceConfigurationProperties {
//	if len(input) == 0 {
//		return nil
//	}
//	value := input[0]
//
//	var schedule maintenanceconfigurations.Schedule
//
//	if value.Frequency == "Daily" {
//		schedule = maintenanceconfigurations.Schedule{
//			Daily: &maintenanceconfigurations.DailySchedule{
//				IntervalDays: value.Interval,
//			},
//		}
//	}
//	if value.Frequency == "Weekly" {
//		schedule = maintenanceconfigurations.Schedule{
//			Weekly: &maintenanceconfigurations.WeeklySchedule{
//				IntervalWeeks: value.Interval,
//				DayOfWeek:     maintenanceconfigurations.WeekDay(value.DayOfWeek),
//			},
//		}
//	}
//	if value.Frequency == "AbsoluteMonthly" {
//		schedule = maintenanceconfigurations.Schedule{
//			AbsoluteMonthly: &maintenanceconfigurations.AbsoluteMonthlySchedule{
//				DayOfMonth:     value.DayOfMonth,
//				IntervalMonths: value.Interval,
//			},
//		}
//	}
//	if value.Frequency == "RelativeMonthly" {
//		schedule = maintenanceconfigurations.Schedule{
//			RelativeMonthly: &maintenanceconfigurations.RelativeMonthlySchedule{
//				DayOfWeek:      maintenanceconfigurations.WeekDay(value.DayOfWeek),
//				WeekIndex:      maintenanceconfigurations.Type(value.WeekIndex),
//				IntervalMonths: value.Interval,
//			},
//		}
//	}
//
//	output := &maintenanceconfigurations.MaintenanceConfigurationProperties{
//		MaintenanceWindow: &maintenanceconfigurations.MaintenanceWindow{
//			StartTime:       value.StartTime,
//			UtcOffset:       pointer.To(value.UTCOffset),
//			NotAllowedDates: expandKubernetesAutomaticClusterMaintenanceConfigurationDateSpans(value.NotAllowed),
//			Schedule:        schedule,
//		},
//	}
//
//	if value.StartDate != "" {
//		startDate, _ := time.Parse(time.RFC3339, value.StartDate)
//		output.MaintenanceWindow.StartDate = pointer.To(startDate.Format("2006-01-02"))
//	}
//
//	if value.Duration != 0 {
//		output.MaintenanceWindow.DurationHours = value.Duration
//	}
//
//	return output
//}
//
//func expandKubernetesAutomaticClusterMaintenanceConfigurationAutoUpgrade(input []MaintenanceWindowAutoUpgradeModel) *maintenanceconfigurations.MaintenanceConfigurationProperties {
//	return expandKubernetesAutomaticClusterMaintenanceConfigurationForCreate(input)
//}
//
//func expandKubernetesAutomaticClusterMaintenanceConfigurationNodeOS(input []MaintenanceWindowNodeOSModel) *maintenanceconfigurations.MaintenanceConfigurationProperties {
//	if len(input) == 0 {
//		return nil
//	}
//	// Convert MaintenanceWindowNodeOSModel to MaintenanceWindowAutoUpgradeModel since they have the same structure
//	converted := []MaintenanceWindowAutoUpgradeModel{{
//		Frequency:  input[0].Frequency,
//		Interval:   input[0].Interval,
//		DayOfWeek:  input[0].DayOfWeek,
//		Duration:   input[0].Duration,
//		WeekIndex:  input[0].WeekIndex,
//		DayOfMonth: input[0].DayOfMonth,
//		StartDate:  input[0].StartDate,
//		StartTime:  input[0].StartTime,
//		UTCOffset:  input[0].UTCOffset,
//		NotAllowed: input[0].NotAllowed,
//	}}
//	return expandKubernetesAutomaticClusterMaintenanceConfigurationForCreate(converted)
//}
//
//func expandKubernetesAutomaticClusterMaintenanceConfigurationTimeSpans(input []MaintenanceWindowNotAllowedModel) *[]maintenanceconfigurations.TimeSpan {
//	results := make([]maintenanceconfigurations.TimeSpan, 0)
//	for _, item := range input {
//		start, _ := time.Parse(time.RFC3339, item.Start)
//		end, _ := time.Parse(time.RFC3339, item.End)
//		results = append(results, maintenanceconfigurations.TimeSpan{
//			Start: pointer.To(start.Format("2006-01-02T15:04:05Z07:00")),
//			End:   pointer.To(end.Format("2006-01-02T15:04:05Z07:00")),
//		})
//	}
//	return &results
//}
//
//func expandKubernetesAutomaticClusterMaintenanceConfigurationDateSpans(input []MaintenanceWindowNotAllowedModel) *[]maintenanceconfigurations.DateSpan {
//	results := make([]maintenanceconfigurations.DateSpan, 0)
//	for _, item := range input {
//		start, _ := time.Parse(time.RFC3339, item.Start)
//		end, _ := time.Parse(time.RFC3339, item.End)
//		results = append(results, maintenanceconfigurations.DateSpan{
//			Start: start.Format("2006-01-02"),
//			End:   end.Format("2006-01-02"),
//		})
//	}
//	return &results
//}
//
//func expandKubernetesAutomaticClusterMaintenanceConfigurationTimeInWeeks(input []MaintenanceWindowAllowedModel) *[]maintenanceconfigurations.TimeInWeek {
//	results := make([]maintenanceconfigurations.TimeInWeek, 0)
//	for _, item := range input {
//		results = append(results, maintenanceconfigurations.TimeInWeek{
//			Day:       pointer.To(maintenanceconfigurations.WeekDay(item.Day)),
//			HourSlots: pointer.To(item.Hours),
//		})
//	}
//	return &results
//}
//
//func flattenKubernetesAutomaticClusterMaintenanceConfiguration(input *maintenanceconfigurations.MaintenanceWindow) []MaintenanceWindowAutoUpgradeModel {
//	results := make([]MaintenanceWindowAutoUpgradeModel, 0)
//	if input == nil {
//		return results
//	}
//
//	startDate := ""
//	if input.StartDate != nil {
//		startDate = *input.StartDate + "T00:00:00Z"
//	}
//	utcOffset := ""
//	if input.UtcOffset != nil {
//		utcOffset = *input.UtcOffset
//	}
//
//	windowModel := MaintenanceWindowAutoUpgradeModel{
//		NotAllowed: flattenKubernetesAutomaticClusterMaintenanceConfigurationDateSpans(input.NotAllowedDates),
//		Duration:   input.DurationHours,
//		StartDate:  startDate,
//		StartTime:  input.StartTime,
//		UTCOffset:  utcOffset,
//	}
//
//	scheduleProps := flattenKubernetesAutomaticClusterMaintenanceConfigurationSchedule(input.Schedule)
//	windowModel.Frequency = scheduleProps["frequency"].(string)
//	windowModel.Interval = scheduleProps["interval"].(int64)
//	windowModel.DayOfWeek = scheduleProps["day_of_week"].(string)
//	windowModel.WeekIndex = scheduleProps["week_index"].(string)
//	windowModel.DayOfMonth = scheduleProps["day_of_month"].(int64)
//
//	return append(results, windowModel)
//}
//
//func flattenKubernetesAutomaticClusterMaintenanceConfigurationSchedule(input maintenanceconfigurations.Schedule) map[string]interface{} {
//	frequency := ""
//	interval := int64(0)
//	if input.Daily != nil {
//		frequency = "Daily"
//		interval = input.Daily.IntervalDays
//	}
//
//	dayOfWeek := ""
//	if input.Weekly != nil {
//		frequency = "Weekly"
//		interval = input.Weekly.IntervalWeeks
//		dayOfWeek = string(input.Weekly.DayOfWeek)
//	}
//
//	dayOfMonth := int64(0)
//	if input.AbsoluteMonthly != nil {
//		frequency = "AbsoluteMonthly"
//		interval = input.AbsoluteMonthly.IntervalMonths
//		dayOfMonth = input.AbsoluteMonthly.DayOfMonth
//	}
//
//	weekIndex := ""
//	if input.RelativeMonthly != nil {
//		frequency = "RelativeMonthly"
//		interval = input.RelativeMonthly.IntervalMonths
//		dayOfWeek = string(input.RelativeMonthly.DayOfWeek)
//		weekIndex = string(input.RelativeMonthly.WeekIndex)
//	}
//
//	return map[string]interface{}{
//		"frequency":    frequency,
//		"interval":     interval,
//		"day_of_week":  dayOfWeek,
//		"week_index":   weekIndex,
//		"day_of_month": dayOfMonth,
//	}
//}
//
//func flattenKubernetesAutomaticClusterMaintenanceConfigurationDefault(input *maintenanceconfigurations.MaintenanceConfigurationProperties) []MaintenanceWindowModel {
//	results := make([]MaintenanceWindowModel, 0)
//	if input == nil {
//		return results
//	}
//	return append(results, MaintenanceWindowModel{
//		NotAllowed: flattenKubernetesAutomaticClusterMaintenanceConfigurationTimeSpans(input.NotAllowedTime),
//		Allowed:    flattenKubernetesAutomaticClusterMaintenanceConfigurationTimeInWeeks(input.TimeInWeek),
//	})
//}
//
//func flattenKubernetesAutomaticClusterMaintenanceConfigurationTimeSpans(input *[]maintenanceconfigurations.TimeSpan) []MaintenanceWindowNotAllowedModel {
//	results := make([]MaintenanceWindowNotAllowedModel, 0)
//	if input == nil {
//		return results
//	}
//
//	for _, item := range *input {
//		var end string
//		if item.End != nil {
//			end = *item.End
//		}
//		var start string
//		if item.Start != nil {
//			start = *item.Start
//		}
//		results = append(results, MaintenanceWindowNotAllowedModel{
//			End:   end,
//			Start: start,
//		})
//	}
//	return results
//}
//
//func flattenKubernetesAutomaticClusterMaintenanceConfigurationDateSpans(input *[]maintenanceconfigurations.DateSpan) []MaintenanceWindowNotAllowedModel {
//	results := make([]MaintenanceWindowNotAllowedModel, 0)
//	if input == nil {
//		return results
//	}
//
//	for _, item := range *input {
//		var end string
//		if item.End != "" {
//			end = item.End + "T23:59:59Z"
//		}
//		var start string
//		if item.Start != "" {
//			start = item.Start + "T00:00:00Z"
//		}
//		results = append(results, MaintenanceWindowNotAllowedModel{
//			End:   end,
//			Start: start,
//		})
//	}
//	return results
//}
//
//func flattenKubernetesAutomaticClusterMaintenanceConfigurationTimeInWeeks(input *[]maintenanceconfigurations.TimeInWeek) []MaintenanceWindowAllowedModel {
//	results := make([]MaintenanceWindowAllowedModel, 0)
//	if input == nil {
//		return results
//	}
//
//	for _, item := range *input {
//		day := ""
//		if item.Day != nil {
//			day = string(*item.Day)
//		}
//		hours := make([]int64, 0)
//		if item.HourSlots != nil {
//			hours = *item.HourSlots
//		}
//		results = append(results, MaintenanceWindowAllowedModel{
//			Day:   day,
//			Hours: hours,
//		})
//	}
//	return results
//}

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
		blobDriverEnabled = pointer.From(profile.BlobCSIDriver.Enabled)
	}

	diskDriverEnabled := false
	if profile.DiskCSIDriver != nil && profile.DiskCSIDriver.Enabled != nil {
		diskDriverEnabled = pointer.From(profile.DiskCSIDriver.Enabled)
	}

	fileDriverEnabled := false
	if profile.FileCSIDriver != nil && profile.FileCSIDriver.Enabled != nil {
		fileDriverEnabled = pointer.From(profile.FileCSIDriver.Enabled)
	}

	snapshotControllerEnabled := false
	if profile.SnapshotController != nil && profile.SnapshotController.Enabled != nil {
		snapshotControllerEnabled = pointer.From(profile.SnapshotController.Enabled)
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

func expandIdentityModel(input []identity.SystemOrUserAssignedList) *identity.SystemOrUserAssignedMap {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	identityType := config.Type

	identityIds := make(map[string]identity.UserAssignedIdentityDetails)
	for _, id := range config.IdentityIds {
		identityIds[id] = identity.UserAssignedIdentityDetails{}
	}

	return &identity.SystemOrUserAssignedMap{
		Type:        identityType,
		IdentityIds: identityIds,
	}
}

func flattenIdentityModel(input *identity.SystemOrUserAssignedMap) []identity.SystemOrUserAssignedList {
	if input == nil {
		return []identity.SystemOrUserAssignedList{}
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

	return []identity.SystemOrUserAssignedList{{
		Type:        input.Type,
		IdentityIds: identityIds,
		PrincipalId: principalId,
		TenantId:    tenantId,
	}}
}
