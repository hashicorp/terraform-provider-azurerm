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
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/capacityreservationgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/snapshots"
	dnsValidate "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/applicationgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipprefixes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/applicationsecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	BootstrapProfile                []BootstrapProfileModel             `tfschema:"bootstrap_profile"`
	CostAnalysisEnabled             bool                                `tfschema:"cost_analysis_enabled"`
	CustomCATrustCertificatesBase64 []string                            `tfschema:"custom_ca_trust_certificates_base64"`
	DefaultNodePool                 []DefaultNodePoolModel              `tfschema:"default_node_pool"`
	DiskEncryptionSetID             string                              `tfschema:"disk_encryption_set_id"`
	DNSPrefix                       string                              `tfschema:"dns_prefix"`
	HTTPProxyConfig                 []HTTPProxyConfigModel              `tfschema:"proxy"`
	Identity                        []identity.SystemOrUserAssignedList `tfschema:"identity"`
	ImageCleanerIntervalHours       int64                               `tfschema:"image_cleaner_interval_in_hours"`
	KeyManagementService            []KeyManagementServiceModel         `tfschema:"key_management_service"`
	KubeletIdentity                 []KubeletIdentityModel              `tfschema:"kubelet_identity"`
	KubernetesVersion               string                              `tfschema:"kubernetes_version"`
	LinuxProfile                    []LinuxProfileModel                 `tfschema:"linux_profile"`
	MicrosoftDefender               []MicrosoftDefenderModel            `tfschema:"microsoft_defender"`
	MonitorMetrics                  []MonitorMetricsModel               `tfschema:"monitor_metrics"`
	NetworkProfile                  []NetworkProfileModel               `tfschema:"network"`
	NodeResourceGroup               string                              `tfschema:"node_resource_group_name"`
	PrivateCluster                  []PrivateClusterModel               `tfschema:"private_cluster"`
	RunCommandEnabled               bool                                `tfschema:"run_command_enabled"`
	ServiceMeshProfile              []ServiceMeshProfileModel           `tfschema:"service_mesh"`
	StorageProfile                  []StorageProfileModel               `tfschema:"storage_profile"`
	SupportPlan                     string                              `tfschema:"support_plan"`
	Tags                            map[string]interface{}              `tfschema:"tags"`
	UpgradeOverride                 []UpgradeOverrideModel              `tfschema:"upgrade_override"`
	WebAppRouting                   []WebAppRoutingModel                `tfschema:"web_app_routing"`
	WindowsProfile                  []WindowsProfileModel               `tfschema:"windows_profile"`
	AIToolchainOperatorEnabled      bool                                `tfschema:"ai_toolchain_operator_enabled"`

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
	FQDN                           string            `tfschema:"fully_qualified_domain_name"`
	HTTPApplicationRoutingZoneName string            `tfschema:"http_application_routing_zone_name"`
	KubeAdminConfig                []KubeConfigModel `tfschema:"kube_admin_config"`
	KubeAdminConfigRaw             string            `tfschema:"kube_admin_config_raw"`
	KubeConfig                     []KubeConfigModel `tfschema:"kube_config"`
	KubeConfigRaw                  string            `tfschema:"kube_config_raw"`
	NodeResourceGroupID            string            `tfschema:"node_resource_group_id"`
	OIDCIssuerURL                  string            `tfschema:"oidc_issuer_url"`
	PrivateFQDN                    string            `tfschema:"private_fully_qualified_domain_name"`
}

type APIServerAccessProfileModel struct {
	AuthorizedIPRanges []string `tfschema:"authorized_ip_ranges"`
	SubnetID           string   `tfschema:"subnet_id"`
}

type PrivateClusterModel struct {
	PrivateClusterPublicFQDNEnabled bool   `tfschema:"public_fully_qualified_domain_name_enabled"`
	PrivateDNSZoneID                string `tfschema:"private_dns_zone_id"`
	DNSPrefixPrivateCluster         string `tfschema:"dns_prefix"`
}

type AutoScalerProfileModel struct {
	BalanceSimilarNodeGroups                 bool    `tfschema:"balance_similar_node_groups_enabled"`
	DaemonsetEvictionForEmptyNodesEnabled    bool    `tfschema:"daemonset_eviction_for_empty_nodes_enabled"`
	DaemonsetEvictionForOccupiedNodesEnabled bool    `tfschema:"daemonset_eviction_for_occupied_nodes_enabled"`
	Expander                                 string  `tfschema:"expander"`
	IgnoreDaemonsetsUtilizationEnabled       bool    `tfschema:"daemonset_ignore_utilization_enabled"`
	MaxGracefulTerminationSec                int64   `tfschema:"maximum_graceful_termination_in_seconds"`
	MaxNodeProvisioningTime                  int64   `tfschema:"maximum_node_provisioning_in_minutes"`
	MaxUnreadyNodes                          int64   `tfschema:"maximum_unready_nodes"`
	MaxUnreadyPercentage                     float64 `tfschema:"maximum_unready_percentage"`
	NewPodScaleUpDelay                       string  `tfschema:"new_pod_scale_up_delay_duration"`
	ScanInterval                             int64   `tfschema:"scan_interval_in_seconds"`
	ScaleDownDelayAfterAdd                   int64   `tfschema:"scale_down_delay_after_add_in_minutes"`
	ScaleDownDelayAfterDelete                int64   `tfschema:"scale_down_delay_after_delete_in_seconds"`
	ScaleDownDelayAfterFailure               int64   `tfschema:"scale_down_delay_after_failure_in_minutes"`
	ScaleDownUnneeded                        int64   `tfschema:"scale_down_unneeded_in_minutes"`
	ScaleDownUnready                         int64   `tfschema:"scale_down_unready_in_minutes"`
	ScaleDownUtilizationThreshold            float64 `tfschema:"scale_down_utilization_threshold"`
	EmptyBulkDeleteMax                       int64   `tfschema:"maximum_empty_bulk_delete"`
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
	TrustedCA  string   `tfschema:"trusted_certificate_authority"`
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
	AdminUsername string `tfschema:"admin_username"`
	SSHKeyData    string `tfschema:"ssh_key_data"`
}

type MicrosoftDefenderModel struct {
	LogAnalyticsWorkspaceID string `tfschema:"log_analytics_workspace_id"`
}

type MonitorMetricsModel struct {
	MonitorMetricsEnabled bool   `tfschema:"enabled"`
	AnnotationsAllowed    string `tfschema:"annotations_allowed"`
	LabelsAllowed         string `tfschema:"labels_allowed"`
}

type NetworkProfileModel struct {
	DNSServiceIP        string                     `tfschema:"dns_service_ip"`
	PodCIDR             string                     `tfschema:"pod_cidr"`
	ServiceCIDR         string                     `tfschema:"service_cidr"`
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
	EffectiveOutboundIPs   []string `tfschema:"effective_outbound_ip_ids"`
}

type NATGatewayProfileModel struct {
	ManagedOutboundIPCount int64    `tfschema:"managed_outbound_ip_count"`
	IdleTimeoutInMinutes   int64    `tfschema:"idle_timeout_in_minutes"`
	EffectiveOutboundIPs   []string `tfschema:"effective_outbound_ip_ids"`
}

type AdvancedNetworkingModel struct {
	ObservabilityEnabled bool `tfschema:"observability_enabled"`
	SecurityEnabled      bool `tfschema:"security_enabled"`
}

type ServiceMeshProfileModel struct {
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
	GMSA          []GMSAModel `tfschema:"group_managed_service_accounts"`
}

type GMSAModel struct {
	DNSServer          string `tfschema:"dns_server"`
	RootDomain         string `tfschema:"root_domain"`
	GMSAProfileEnabled bool   `tfschema:"enabled"`
}

type KubeConfigModel struct {
	Host                 string `tfschema:"host"`
	Username             string `tfschema:"username"`
	Password             string `tfschema:"password"`
	ClientCertificate    string `tfschema:"client_certificate"`
	ClientKey            string `tfschema:"client_key"`
	ClusterCACertificate string `tfschema:"cluster_ca_certificate"`
}

type DefaultNodePoolModel struct {
	Name                       string                    `tfschema:"name"`
	TemporaryNameForRotation   string                    `tfschema:"temporary_name_for_rotation"`
	VMSize                     string                    `tfschema:"virtual_machine_size"`
	CapacityReservationGroupID string                    `tfschema:"capacity_reservation_group_id"`
	KubeletConfig              []KubeletConfigModel      `tfschema:"kubelet_config"`
	LinuxOSConfig              []LinuxOSConfigModel      `tfschema:"linux_os_config"`
	FipsEnabled                bool                      `tfschema:"fips_enabled"`
	GPUInstance                string                    `tfschema:"gpu_instance"`
	GPUDriver                  bool                      `tfschema:"gpu_driver_enabled"`
	KubeletDiskType            string                    `tfschema:"kubelet_disk_type"`
	MaxPods                    int64                     `tfschema:"maximum_pods"`
	NodeNetworkProfile         []NodeNetworkProfileModel `tfschema:"node_network_profile"`
	NodeCount                  int64                     `tfschema:"node_count"`
	NodeLabels                 map[string]string         `tfschema:"node_labels"`
	NodePublicIPPrefixID       string                    `tfschema:"node_public_ip_prefix_id"`
	Tags                       map[string]interface{}    `tfschema:"tags"`
	OSDiskSizeGB               int64                     `tfschema:"os_disk_size_in_gb"`
	OSSKU                      string                    `tfschema:"os_sku"`
	UltraSSDEnabled            bool                      `tfschema:"ultra_ssd_enabled"`
	VnetSubnetID               string                    `tfschema:"subnet_id"`
	OrchestratorVersion        string                    `tfschema:"orchestrator_version"`
	ProximityPlacementGroupID  string                    `tfschema:"proximity_placement_group_id"`
	SnapshotID                 string                    `tfschema:"snapshot_id"`
	HostGroupID                string                    `tfschema:"host_group_id"`
	UpgradeSettings            []UpgradeSettingsModel    `tfschema:"upgrade_settings"`
	NodePublicIPEnabled        bool                      `tfschema:"node_public_ip_enabled"`
	HostEncryptionEnabled      bool                      `tfschema:"host_encryption_enabled"`
}

type KubeletConfigModel struct {
	CPUManagerPolicy      bool     `tfschema:"cpu_manager_policy_static_enabled"`
	CPUCfsQuotaEnabled    bool     `tfschema:"cpu_cfs_quota_enabled"`
	CPUCfsQuotaPeriod     string   `tfschema:"cpu_cfs_quota_period"`
	ImageGcHighThreshold  int64    `tfschema:"image_gc_high_threshold"`
	ImageGcLowThreshold   int64    `tfschema:"image_gc_low_threshold"`
	TopologyManagerPolicy string   `tfschema:"topology_manager_policy"`
	AllowedUnsafeSysctls  []string `tfschema:"allowed_unsafe_sysctls"`
	ContainerLogMaxSizeMB int64    `tfschema:"container_log_maximum_size_in_mb"`
	ContainerLogMaxFiles  int64    `tfschema:"container_log_maximum_files"`
	PodMaxPid             int64    `tfschema:"pod_maximum_pid"`
}

type LinuxOSConfigModel struct {
	SysctlConfig              []SysctlConfigModel `tfschema:"sysctl_config"`
	TransparentHugePage       string              `tfschema:"transparent_huge_page"`
	TransparentHugePageDefrag string              `tfschema:"transparent_huge_page_defragmentation"`
	SwapFileSizeMB            int64               `tfschema:"swap_file_size_in_mb"`
}

type SysctlConfigModel struct {
	FsAioMaxNr                     int64 `tfschema:"fs_aio_max_nr"`
	FsFileMax                      int64 `tfschema:"fs_file_max"`
	FsInotifyMaxUserWatches        int64 `tfschema:"fs_inotify_max_user_watches"`
	FsNrOpen                       int64 `tfschema:"fs_nr_open"`
	KernelThreadsMax               int64 `tfschema:"kernel_threads_max"`
	NetCoreNetdevMaxBacklog        int64 `tfschema:"net_core_netdev_max_backlog"`
	NetCoreOptmemMax               int64 `tfschema:"net_core_optmem_max"`
	NetCoreRmemDefault             int64 `tfschema:"net_core_rmem_default"`
	NetCoreRmemMax                 int64 `tfschema:"net_core_rmem_max"`
	NetCoreSomaxconn               int64 `tfschema:"net_core_somaxconn"`
	NetCoreWmemDefault             int64 `tfschema:"net_core_wmem_default"`
	NetCoreWmemMax                 int64 `tfschema:"net_core_wmem_max"`
	NetIPv4IPLocalPortRangeMin     int64 `tfschema:"net_ipv4_ip_local_port_range_min"`
	NetIPv4IPLocalPortRangeMax     int64 `tfschema:"net_ipv4_ip_local_port_range_maximum"`
	NetIPv4NeighDefaultGcThresh1   int64 `tfschema:"net_ipv4_neigh_default_gc_thresh1"`
	NetIPv4NeighDefaultGcThresh2   int64 `tfschema:"net_ipv4_neigh_default_gc_thresh2"`
	NetIPv4NeighDefaultGcThresh3   int64 `tfschema:"net_ipv4_neigh_default_gc_thresh3"`
	NetIPv4TCPFinTimeout           int64 `tfschema:"net_ipv4_tcp_fin_timeout"`
	NetIPv4TCPKeepaliveIntvl       int64 `tfschema:"net_ipv4_tcp_keepalive_intvl"`
	NetIPv4TCPKeepaliveProbes      int64 `tfschema:"net_ipv4_tcp_keepalive_probes"`
	NetIPv4TCPKeepaliveTime        int64 `tfschema:"net_ipv4_tcp_keepalive_time"`
	NetIPv4TCPMaxSynBacklog        int64 `tfschema:"net_ipv4_tcp_max_syn_backlog"`
	NetIPv4TCPMaxTwBuckets         int64 `tfschema:"net_ipv4_tcp_max_tw_buckets"`
	NetIPv4TCPTwReuse              bool  `tfschema:"net_ipv4_tcp_tw_reuse_enabled"`
	NetNetfilterNfConntrackBuckets int64 `tfschema:"net_netfilter_nf_conntrack_buckets"`
	NetNetfilterNfConntrackMax     int64 `tfschema:"net_netfilter_nf_conntrack_max"`
	VMMaxMapCount                  int64 `tfschema:"vm_max_map_count"`
	VMSwappiness                   int64 `tfschema:"vm_swappiness"`
	VMVfsCachePressure             int64 `tfschema:"vm_vfs_cache_pressure"`
}

type NodeNetworkProfileModel struct {
	AllowedHostPorts            []AllowedHostPortsModel `tfschema:"allowed_host_ports"`
	ApplicationSecurityGroupIDs []string                `tfschema:"application_security_group_ids"`
	NodePublicIPTags            map[string]string       `tfschema:"node_public_ip_tags"`
}

type AllowedHostPortsModel struct {
	PortStart int64  `tfschema:"port_start"`
	PortEnd   int64  `tfschema:"port_end"`
	Protocol  string `tfschema:"protocol"`
}

type UpgradeSettingsModel struct {
	MaxSurge                  string `tfschema:"maximum_surge"`
	DrainTimeoutInMinutes     int64  `tfschema:"drain_timeout_in_minutes"`
	NodeSoakDurationInMinutes int64  `tfschema:"node_soak_duration_in_minutes"`
	UndrainableNodeBehavior   string `tfschema:"undrainable_node_behavior"`
}

type ACIConnectorLinuxModel struct {
	SubnetName        string                   `tfschema:"subnet_name"`
	ConnectorIdentity []ConnectorIdentityModel `tfschema:"connector_identity"`
}

type ConnectorIdentityModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type ConfidentialComputingModel struct {
	SGXQuoteHelperEnabled bool `tfschema:"sgx_quote_helper_enabled"`
}

type IngressApplicationGatewayModel struct {
	GatewayID                         string                                   `tfschema:"gateway_id"`
	GatewayName                       string                                   `tfschema:"gateway_name"`
	SubnetCIDR                        string                                   `tfschema:"subnet_cidr"`
	SubnetID                          string                                   `tfschema:"subnet_id"`
	EffectiveGatewayID                string                                   `tfschema:"effective_gateway_id"`
	IngressApplicationGatewayIdentity []IngressApplicationGatewayIdentityModel `tfschema:"identity"`
}

type IngressApplicationGatewayIdentityModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type KeyVaultSecretsProviderModel struct {
	SecretRotationInterval string                `tfschema:"secret_rotation_interval"`
	SecretIdentity         []SecretIdentityModel `tfschema:"secret_identity"`
}

type SecretIdentityModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type OMSAgentModel struct {
	LogAnalyticsWorkspaceID     string                  `tfschema:"log_analytics_workspace_id"`
	MSIAuthForMonitoringEnabled *bool                   `tfschema:"msi_auth_for_monitoring_enabled"`
	OMSAgentIdentity            []OMSAgentIdentityModel `tfschema:"identity"`
}

type OMSAgentIdentityModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
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
		if err != nil || resp.Model == nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}

		if resp.Model.Sku == nil || resp.Model.Sku.Name == nil {
			return fmt.Errorf("importing %s: SKU information is missing", id)
		}

		if pointer.From(resp.Model.Sku.Name) != managedclusters.ManagedClusterSKUNameAutomatic {
			return fmt.Errorf("importing %s: specified Kubernetes Cluster is not using the SKU `Automatic`, got `%s`", id, pointer.From(resp.Model.Sku.Name))
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
				if rd.HasChange("windows_profile.0.group_managed_service_accounts") {
					old, new := rd.GetChange("windows_profile.0.group_managed_service_accounts")
					oldList := old.([]interface{})
					newList := new.([]interface{})
					if len(oldList) > 0 && len(newList) == 0 {
						if err := metadata.ResourceDiff.ForceNew("windows_profile.group_managed_service_accounts"); err != nil {
							return err
						}
					}
				}

				if rd.HasChange("windows_profile.0.group_managed_service_accounts.0.dns_server") {
					old, new := rd.GetChange("windows_profile.0.group_managed_service_accounts.0.dns_server")
					if old.(string) != "" && new.(string) == "" {
						if err := metadata.ResourceDiff.ForceNew("windows_profile.group_managed_service_accounts.dns_server"); err != nil {
							return err
						}
					}
				}

				if rd.HasChange("windows_profile.0.group_managed_service_accounts.0.root_domain") {
					old, new := rd.GetChange("windows_profile.0.group_managed_service_accounts.0.root_domain")
					if old.(string) != "" && new.(string) == "" {
						if err := metadata.ResourceDiff.ForceNew("windows_profile.group_managed_service_accounts.root_domain"); err != nil {
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
			artifactSource := rd.Get("bootstrap_profile.0.artifact_source").(string)
			if outboundType == string(managedclusters.OutboundTypeNone) && artifactSource != string(managedclusters.ArtifactSourceCache) {
				return fmt.Errorf("when `network.outbound_type` is set to `none`, `bootstrap_profile.artifact_source` must be set to `Cache`")
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
			ValidateFunc: containerValidate.KubernetesClusterName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"default_node_pool": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: containerValidate.KubernetesAgentPoolName,
					},

					"capacity_reservation_group_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: capacityreservationgroups.ValidateCapacityReservationGroupID,
					},

					"fips_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"gpu_instance": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.GPUInstanceProfileMIGOneg),
							string(managedclusters.GPUInstanceProfileMIGTwog),
							string(managedclusters.GPUInstanceProfileMIGThreeg),
							string(managedclusters.GPUInstanceProfileMIGFourg),
							string(managedclusters.GPUInstanceProfileMIGSeveng),
						}, false),
					},

					"gpu_driver_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      false,
						RequiredWith: []string{"default_node_pool.0.gpu_instance"},
					},

					"host_encryption_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"host_group_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: computeValidate.HostGroupID,
					},

					"kubelet_config": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"cpu_manager_policy_static_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"cpu_cfs_quota_enabled": {
									Type:     pluginsdk.TypeBool,
									Default:  true,
									Optional: true,
								},

								"cpu_cfs_quota_period": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"image_gc_high_threshold": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntBetween(0, 100),
								},

								"image_gc_low_threshold": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntBetween(1, 100),
								},

								"topology_manager_policy": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										"best-effort",
										"restricted",
										"single-numa-node",
									}, false),
								},

								"allowed_unsafe_sysctls": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"container_log_maximum_size_in_mb": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(0),
								},

								"container_log_maximum_files": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(2),
								},

								"pod_maximum_pid": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(0),
								},
							},
						},
					},

					"kubelet_disk_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  managedclusters.KubeletDiskTypeOS,
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.KubeletDiskTypeOS),
							string(managedclusters.KubeletDiskTypeTemporary),
						}, false),
					},
					"linux_os_config": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"sysctl_config": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											// Abbreviations acceptable because refers to system named attributes
											"fs_aio_max_nr": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(65536, 6553500),
											},

											"fs_file_max": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(8192, 12000500),
											},

											"fs_inotify_max_user_watches": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(781250, 2097152),
											},

											"fs_nr_open": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(8192, 20000500),
											},

											"kernel_threads_max": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(20, 513785),
											},

											"net_core_netdev_max_backlog": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1000, 3240000),
											},

											"net_core_optmem_max": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(20480, 4194304),
											},

											"net_core_rmem_default": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(212992, 134217728),
											},

											"net_core_rmem_max": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(212992, 134217728),
											},

											"net_core_somaxconn": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(4096, 3240000),
											},

											"net_core_wmem_default": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(212992, 134217728),
											},

											"net_core_wmem_max": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(212992, 134217728),
											},

											"net_ipv4_ip_local_port_range_min": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1024, 60999),
											},

											"net_ipv4_ip_local_port_range_maximum": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(32768, 65535),
											},

											"net_ipv4_neigh_default_gc_thresh1": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(128, 80000),
											},

											"net_ipv4_neigh_default_gc_thresh2": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(512, 90000),
											},

											"net_ipv4_neigh_default_gc_thresh3": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1024, 100000),
											},

											"net_ipv4_tcp_fin_timeout": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(5, 120),
											},

											"net_ipv4_tcp_keepalive_intvl": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(10, 90),
											},

											"net_ipv4_tcp_keepalive_probes": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 15),
											},

											"net_ipv4_tcp_keepalive_time": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(30, 432000),
											},

											"net_ipv4_tcp_max_syn_backlog": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(128, 3240000),
											},

											"net_ipv4_tcp_max_tw_buckets": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(8000, 1440000),
											},

											"net_ipv4_tcp_tw_reuse_enabled": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
												Default:  false,
											},

											"net_netfilter_nf_conntrack_buckets": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(65536, 524288),
											},

											"net_netfilter_nf_conntrack_max": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(131072, 2097152),
											},

											"vm_max_map_count": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(65530, 262144),
											},

											"vm_swappiness": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(0, 100),
											},

											"vm_vfs_cache_pressure": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(0, 100),
											},
										},
									},
								},

								"transparent_huge_page": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									// omitted never as option, set as default
									ValidateFunc: validation.StringInSlice([]string{
										"always",
										"madvise",
									}, false),
								},

								"transparent_huge_page_defragmentation": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									// omitted never as option, set as default
									ValidateFunc: validation.StringInSlice([]string{
										"always",
										"defer",
										"defer+madvise",
										"madvise",
									}, false),
								},

								"swap_file_size_in_mb": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(0),
								},
							},
						},
					},

					"maximum_pods": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						// O+C - Azure assigns a default maximum pods value based on network plugin if not specified
						Computed:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"node_count": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						// O+C - Azure may adjust node count based on autoscaling configuration if not explicitly set
						Computed:     true,
						ValidateFunc: validation.IntBetween(1, 1000),
					},

					"node_labels": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						// O+C - Azure may apply default labels to nodes if not specified
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"node_network_profile": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"allowed_host_ports": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"port_start": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 65535),
											},

											"port_end": {
												Type:         pluginsdk.TypeInt,
												Optional:     true,
												ValidateFunc: validation.IntBetween(1, 65535),
											},

											"protocol": {
												Type:     pluginsdk.TypeString,
												Optional: true,
												ValidateFunc: validation.StringInSlice([]string{
													string(agentpools.ProtocolTCP),
													string(agentpools.ProtocolUDP),
												}, false),
											},
										},
									},
								},

								"application_security_group_ids": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: applicationsecuritygroups.ValidateApplicationSecurityGroupID,
									},
								},

								"node_public_ip_tags": {
									Type:     pluginsdk.TypeMap,
									Optional: true,
									ForceNew: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},

					"node_public_ip_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"node_public_ip_prefix_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
						RequiredWith: []string{"default_node_pool.0.node_public_ip_enabled"},
					},

					"orchestrator_version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// O+C - Azure uses the cluster's Kubernetes version if not specified
						Computed:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"os_disk_size_in_gb": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						// O+C - Azure calculates default OS disk size based on VM size if not specified
						Computed:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"os_sku": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  agentpools.OSSKUAzureLinux,
						ValidateFunc: validation.StringInSlice([]string{
							string(agentpools.OSSKUAzureLinux),
							string(agentpools.OSSKUAzureLinuxThree),
						}, false),
					},

					"proximity_placement_group_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
					},

					"snapshot_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: snapshots.ValidateSnapshotID,
					},

					"tags": commonschema.Tags(),

					"temporary_name_for_rotation": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: containerValidate.KubernetesAgentPoolName,
					},

					"ultra_ssd_enabled": {
						Type:     pluginsdk.TypeBool,
						Default:  false,
						Optional: true,
					},

					"upgrade_settings": {
						Type: pluginsdk.TypeList,
						// O+C - Azure provides default upgrade settings if not specified
						Optional: true,
						Computed: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"maximum_surge": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "10%",
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"drain_timeout_in_minutes": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntAtLeast(0),
								},
								"node_soak_duration_in_minutes": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntBetween(0, 30),
								},
								"undrainable_node_behavior": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(agentpools.PossibleValuesForUndrainableNodeBehavior(), true),
								},
							},
						},
					},

					"virtual_machine_size": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// O+C can be left blank and azure will select an available vm size
						Computed:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},
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
			ExactlyOneOf: []string{"dns_prefix", "private_cluster.0.dns_prefix"},
			ValidateFunc: containerValidate.KubernetesDNSPrefix,
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
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C - Azure generates a default name (MC_<rg>_<cluster>_<location>) if not specified
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
					"public_fully_qualified_domain_name_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"private_dns_zone_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						ValidateFunc: validation.Any(
							privatezones.ValidatePrivateDnsZoneID,
							validation.StringInSlice([]string{
								"System",
								"None",
							}, false),
						),
					},
					"dns_prefix": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ExactlyOneOf: []string{"dns_prefix", "private_cluster.0.dns_prefix"},
						ValidateFunc: containerValidate.KubernetesDNSPrefix,
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
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      600,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"maximum_node_provisioning_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      15,
						ValidateFunc: validation.IntAtLeast(0),
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
					"new_pod_scale_up_delay_duration": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "0s",
						ValidateFunc: containerValidate.Duration,
					},
					"scan_interval_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      10,
						ValidateFunc: validation.IntAtLeast(1),
					},
					"scale_down_delay_after_add_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      10,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"scale_down_delay_after_delete_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      3,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"scale_down_delay_after_failure_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      3,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"scale_down_unneeded_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      10,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"scale_down_unready_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      20,
						ValidateFunc: validation.IntAtLeast(0),
					},
					"scale_down_utilization_threshold": {
						Type:         pluginsdk.TypeFloat,
						Optional:     true,
						Default:      0.5,
						ValidateFunc: validation.FloatAtLeast(0),
					},
					"maximum_empty_bulk_delete": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      10,
						ValidateFunc: validation.IntAtLeast(0),
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
			// NOTE: O+C - Azure populates default AAD settings from subscription context if not fully specified
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// NOTE: O+C - Azure uses subscription's tenant ID if not specified
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
						// // NOTE: O+C - Azure may populate default admin groups from tenant if not specified
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

		"bootstrap_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// NOTE: O+C - Azure provides default bootstrap configuration if not specified
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

		"proxy": {
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
							Type:         schema.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"trusted_certificate_authority": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
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
			// NOTE: O+C - Azure creates a managed identity for kubelet if not specified
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// NOTE: O+C - Azure generates client_id when creating managed identity
						Computed: true,
						ForceNew: true,
						RequiredWith: []string{
							"kubelet_identity.0.object_id",
							"kubelet_identity.0.user_assigned_identity_id",
							"identity.0.identity_ids",
						},
						ValidateFunc: validation.IsUUID,
					},
					"object_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// NOTE: O+C - Azure generates object_id when creating managed identity
						Computed: true,
						ForceNew: true,
						RequiredWith: []string{
							"kubelet_identity.0.client_id",
							"kubelet_identity.0.user_assigned_identity_id",
							"identity.0.identity_ids",
						},
						ValidateFunc: validation.IsUUID,
					},
					"user_assigned_identity_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// NOTE: O+C - Azure generates identity resource ID when creating managed identity
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
					"ssh_key_data": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
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
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
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
			// NOTE: O+C - Azure provides default network configuration (Azure CNI) if not specified
			Computed: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dns_service_ip": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// NOTE: O+C - Azure calculates a default DNS service IP from service CIDR if not specified
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validate.IPv4Address,
					},
					"pod_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// NOTE: O+C - Azure calculates a default pod CIDR if not specified
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validate.CIDR,
					},
					"service_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// NOTE: O+C - Azure calculates a default service CIDR if not specified
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validate.CIDR,
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
							// expose none as an option because as default it would be breaking
							string(managedclusters.OutboundTypeNone),
						}, false),
					},
					"load_balancer": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						ForceNew: true,
						Optional: true,
						// NOTE: O+C - Azure provides default load balancer configuration if not specified
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
									Type:     pluginsdk.TypeInt,
									Optional: true,
									// NOTE: O+C - Azure assigns default outbound IP count if not specified
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
										ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
									},
								},
								"outbound_ip_address_ids": {
									Type:          pluginsdk.TypeSet,
									Optional:      true,
									ConflictsWith: []string{"network.0.load_balancer.0.managed_outbound_ip_count", "network.0.load_balancer.0.outbound_ip_prefix_ids"},
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: commonids.ValidatePublicIPAddressID,
									},
								},
								"effective_outbound_ip_ids": {
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
						// NOTE: O+C - Azure provides default NAT gateway configuration if not specified
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
									Type:     pluginsdk.TypeInt,
									Optional: true,
									// NOTE: O+C - Azure assigns default outbound IP count for NAT gateway if not specified
									Computed:     true,
									ValidateFunc: validation.IntBetween(1, 100),
								},
								"effective_outbound_ip_ids": {
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
					"internal_ingress_gateway_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"external_ingress_gateway_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
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
							ValidateFunc: validation.StringStartsWithOneOf("asm-"),
						},
					},
				},
			},
		},

		"storage_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// NOTE: O+C - Azure provides default storage profile configuration if not specified
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"blob_driver_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      false,
						AtLeastOneOf: []string{"storage_profile.0.blob_driver_enabled", "storage_profile.0.disk_driver_enabled", "storage_profile.0.file_driver_enabled", "storage_profile.0.snapshot_controller_enabled"},
					},
					"disk_driver_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      true,
						AtLeastOneOf: []string{"storage_profile.0.blob_driver_enabled", "storage_profile.0.disk_driver_enabled", "storage_profile.0.file_driver_enabled", "storage_profile.0.snapshot_controller_enabled"},
					},
					"file_driver_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      true,
						AtLeastOneOf: []string{"storage_profile.0.blob_driver_enabled", "storage_profile.0.disk_driver_enabled", "storage_profile.0.file_driver_enabled", "storage_profile.0.snapshot_controller_enabled"},
					},
					"snapshot_controller_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      true,
						AtLeastOneOf: []string{"storage_profile.0.blob_driver_enabled", "storage_profile.0.disk_driver_enabled", "storage_profile.0.file_driver_enabled", "storage_profile.0.snapshot_controller_enabled"},
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
			// NOTE: O+C - Azure provides default web app routing configuration if not specified
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
							// Allow none as option because none is breaking in most cases
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
					"group_managed_service_accounts": {
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
								"dns_server": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"root_domain": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},
		"aci_connector_linux": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"connector_identity": {
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
		"confidential_computing": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"sgx_quote_helper_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},
		"http_application_routing_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
		"ingress_application_gateway": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"gateway_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ConflictsWith: []string{
							"ingress_application_gateway.0.subnet_cidr",
							"ingress_application_gateway.0.subnet_id",
							"ingress_application_gateway.0.gateway_name",
						},
						AtLeastOneOf: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_cidr",
							"ingress_application_gateway.0.subnet_id",
						},
						ValidateFunc: applicationgateways.ValidateApplicationGatewayID,
					},
					"gateway_name": {
						Type:          pluginsdk.TypeString,
						Optional:      true,
						ValidateFunc:  validation.StringIsNotEmpty,
						ConflictsWith: []string{"ingress_application_gateway.0.gateway_id"},
					},
					"subnet_cidr": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ConflictsWith: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_id",
							"ingress_application_gateway.0.gateway_name",
						},
						AtLeastOneOf: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_cidr",
							"ingress_application_gateway.0.subnet_id",
						},
						ValidateFunc: validate.CIDR,
					},
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ConflictsWith: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_cidr",
							"ingress_application_gateway.0.gateway_name",
						},
						AtLeastOneOf: []string{
							"ingress_application_gateway.0.gateway_id",
							"ingress_application_gateway.0.subnet_cidr",
							"ingress_application_gateway.0.subnet_id",
						},
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"effective_gateway_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"identity": {
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
		"key_vault_secrets_provider": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			// NOTE: O+C - Azure provides default Key Vault secrets provider configuration if not specified
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"secret_rotation_interval": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "2m",
						ValidateFunc: containerValidate.Duration,
					},
					"secret_identity": {
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
		"oms_agent": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"log_analytics_workspace_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: workspaces.ValidateWorkspaceID,
					},
					"msi_auth_for_monitoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"identity": {
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
		"open_service_mesh_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}

	return arguments
}

func (r KubernetesAutomaticClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"current_kubernetes_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"fully_qualified_domain_name": {
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

		"portal_fully_qualified_domain_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_fully_qualified_domain_name": {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			if err := validateKubernetesAutomaticClusterTyped(&model, nil); err != nil {
				return fmt.Errorf("validating configuration: %+v", err)
			}

			networkProfile, err := expandKubernetesAutomaticClusterNetworkProfile(model.NetworkProfile)
			if err != nil {
				return fmt.Errorf("expanding network profile: %+v", err)
			}

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

			apiAccessProfile := expandKubernetesAutomaticClusterAPIAccessProfile(model)

			var azureADProfile *managedclusters.ManagedClusterAADProfile
			if len(model.AzureActiveDirectoryRBAC) > 0 {
				azureADProfile = &managedclusters.ManagedClusterAADProfile{
					Managed:             pointer.To(true),
					AdminGroupObjectIDs: &model.AzureActiveDirectoryRBAC[0].AdminGroupObjectIDs,
				}
			}

			parameters := managedclusters.ManagedCluster{
				Location: location.Normalize(model.Location),
				Sku: &managedclusters.ManagedClusterSKU{
					Name: pointer.To(managedclusters.ManagedClusterSKUName("automatic")),
					Tier: pointer.To(managedclusters.ManagedClusterSKUTier("standard")),
				},
				Properties: &managedclusters.ManagedClusterProperties{
					ApiServerAccessProfile: apiAccessProfile,
					AadProfile:             azureADProfile,
					AddonProfiles:          addonProfiles,
					AgentPoolProfiles:      agentProfiles,
					AutoScalerProfile:      expandKubernetesAutomaticClusterAutoScalerProfile(model.AutoScalerProfile),
					AutoUpgradeProfile:     pointer.To(managedclusters.ManagedClusterAutoUpgradeProfile{}),
					AzureMonitorProfile:    expandKubernetesAutomaticClusterAzureMonitorProfile(model.MonitorMetrics),
					KubernetesVersion:      pointer.To(model.KubernetesVersion),
					BootstrapProfile:       expandKubernetesAutomaticClusterBootstrapProfile(model.BootstrapProfile),
					LinuxProfile:           expandKubernetesAutomaticClusterLinuxProfile(model.LinuxProfile),
					WindowsProfile:         expandKubernetesAutomaticClusterWindowsProfile(model.WindowsProfile),
					MetricsProfile:         expandKubernetesAutomaticClusterMetricsProfile(model.CostAnalysisEnabled),
					NetworkProfile:         networkProfile,
					NodeResourceGroup:      pointer.To(model.NodeResourceGroup),
					HTTPProxyConfig:        expandKubernetesAutomaticClusterHttpProxyConfig(model.HTTPProxyConfig),
					SecurityProfile:        securityProfile,
					StorageProfile:         expandKubernetesAutomaticClusterStorageProfile(model.StorageProfile),
					UpgradeSettings:        expandKubernetesAutomaticClusterUpgradeOverride(model.UpgradeOverride),
					IngressProfile:         expandKubernetesAutomaticClusterWebAppRouting(model.WebAppRouting, false),
					ServiceMeshProfile:     expandKubernetesAutomaticClusterServiceMeshProfile(model.ServiceMeshProfile, nil),
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

			if len(model.PrivateCluster) > 0 && model.PrivateCluster[0].DNSPrefixPrivateCluster != "" {
				if apiAccessProfile.PrivateDNSZone == nil || *apiAccessProfile.PrivateDNSZone == "System" || *apiAccessProfile.PrivateDNSZone == "None" {
					return fmt.Errorf("`private_cluster.0.dns_prefix` should only be set for private cluster with custom private dns zone")
				}
				parameters.Properties.FqdnSubdomain = pointer.To(model.PrivateCluster[0].DNSPrefixPrivateCluster)
			} else {
				parameters.Properties.DnsPrefix = pointer.To(model.DNSPrefix)
			}

			if model.DiskEncryptionSetID != "" {
				parameters.Properties.DiskEncryptionSetID = pointer.To(model.DiskEncryptionSetID)
			}

			if model.SupportPlan != "" {
				parameters.Properties.SupportPlan = pointer.To(managedclusters.KubernetesSupportPlan(model.SupportPlan))
			}

			err = client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, managedclusters.DefaultCreateOrUpdateOperationOptions(), metadata.SetIDCallback(&id))
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
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
	if err := metadata.Decode(pointer.To(config)); err != nil {
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

			if len(state.PrivateCluster) > 0 {
				state.PrivateCluster[0].DNSPrefixPrivateCluster = pointer.From(props.FqdnSubdomain)
			}

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
					"default_node_pool.0.os_disk_size_in_gb",
					"default_node_pool.0.pod_subnet_id",
					"default_node_pool.0.snapshot_id",
					"default_node_pool.0.ultra_ssd_enabled",
					"default_node_pool.0.subnet_id",
					"default_node_pool.0.virtual_machine_size",
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
				metadata.ResourceData.HasChange("confidential_computing_sgx_quote_helper_enabled") ||
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
	privateDNSZoneID := pointer.From(profile.PrivateDNSZone)
	switch {
	case profile.PrivateDNSZone != nil && strings.EqualFold("System", *profile.PrivateDNSZone):
		privateDNSZoneID = "System"
	case profile.PrivateDNSZone != nil && strings.EqualFold("None", *profile.PrivateDNSZone):
		privateDNSZoneID = "None"
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
	if len(config.SSHKeyData) > 0 {
		keyData = config.SSHKeyData
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

	sshKeyData := ""
	ssh := profile.Ssh
	if keys := ssh.PublicKeys; keys != nil {
		for _, sshKey := range keys {
			keyData := ""
			if kd := sshKey.KeyData; kd != "" {
				keyData = kd
			}
			sshKeyData = keyData
		}
	}

	return []LinuxProfileModel{
		{
			AdminUsername: adminUsername,
			SSHKeyData:    sshKeyData,
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

	if config.ServiceCIDR != "" {
		networkProfile.ServiceCidr = pointer.To(config.ServiceCIDR)
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

	// TODO - Remove the workaround below once issue https://github.com/Azure/azure-rest-api-specs/issues/18056 is resolved
	sku := profile.LoadBalancerSku
	for _, v := range managedclusters.PossibleValuesForLoadBalancerSku() {
		if strings.EqualFold(v, string(*sku)) {
			lsSku := managedclusters.LoadBalancerSku(v)
			sku = pointer.To(lsSku)
		}
	}

	advancedNetworking := flattenKubernetesAutomaticClusterAdvancedNetworking(profile.AdvancedNetworking)

	return []NetworkProfileModel{
		{
			DNSServiceIP:        dnsServiceIP,
			LoadBalancerSKU:     string(pointer.From(sku)),
			LoadBalancerProfile: lbProfiles,
			NATGatewayProfile:   ngwProfiles,
			PodCIDR:             podCidr,
			ServiceCIDR:         serviceCidr,
			OutboundType:        outboundType,
			AdvancedNetworking:  advancedNetworking,
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

	if config.MaxGracefulTerminationSec != 0 {
		profile.MaxGracefulTerminationSec = pointer.To(fmt.Sprintf("%d", config.MaxGracefulTerminationSec))
	}

	if config.MaxNodeProvisioningTime > 0 {
		profile.MaxNodeProvisionTime = pointer.To(fmt.Sprintf("%dm", config.MaxNodeProvisioningTime))
	}

	if config.NewPodScaleUpDelay != "" {
		profile.NewPodScaleUpDelay = pointer.To(config.NewPodScaleUpDelay)
	}

	if config.ScaleDownDelayAfterAdd > 0 {
		profile.ScaleDownDelayAfterAdd = pointer.To(fmt.Sprintf("%dm", config.ScaleDownDelayAfterAdd))
	}

	if config.ScaleDownDelayAfterDelete > 0 {
		profile.ScaleDownDelayAfterDelete = pointer.To(fmt.Sprintf("%ds", config.ScaleDownDelayAfterDelete))
	}

	if config.ScaleDownDelayAfterFailure > 0 {
		profile.ScaleDownDelayAfterFailure = pointer.To(fmt.Sprintf("%dm", config.ScaleDownDelayAfterFailure))
	}

	if config.ScaleDownUnneeded > 0 {
		profile.ScaleDownUnneededTime = pointer.To(fmt.Sprintf("%dm", config.ScaleDownUnneeded))
	}

	if config.ScaleDownUnready > 0 {
		profile.ScaleDownUnreadyTime = pointer.To(fmt.Sprintf("%dm", config.ScaleDownUnready))
	}

	if config.ScaleDownUtilizationThreshold != 0.0 {
		profile.ScaleDownUtilizationThreshold = pointer.To(strconv.FormatFloat(config.ScaleDownUtilizationThreshold, 'f', 1, 64))
	}

	if config.EmptyBulkDeleteMax != 0 {
		profile.MaxEmptyBulkDelete = pointer.To(strconv.FormatInt(config.EmptyBulkDeleteMax, 10))
	}

	if config.SkipNodesWithLocalStorage {
		profile.SkipNodesWithLocalStorage = pointer.To(strconv.FormatBool(config.SkipNodesWithLocalStorage))
	}

	if config.SkipNodesWithSystemPods {
		profile.SkipNodesWithSystemPods = pointer.To(strconv.FormatBool(config.SkipNodesWithSystemPods))
	}

	if config.ScanInterval > 0 {
		profile.ScanInterval = pointer.To(fmt.Sprintf("%ds", config.ScanInterval))
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

	maxGracefulTerminationSec := int64(0)
	if profile.MaxGracefulTerminationSec != nil {
		if val, err := strconv.ParseInt(pointer.From(profile.MaxGracefulTerminationSec), 10, 64); err == nil {
			maxGracefulTerminationSec = val
		}
	}

	MaxNodeProvisioningTime := int64(0)
	if profile.MaxNodeProvisionTime != nil {
		// Parse "15m" to 15
		timeStr := pointer.From(profile.MaxNodeProvisionTime)
		if strings.HasSuffix(timeStr, "m") {
			if val, err := strconv.ParseInt(strings.TrimSuffix(timeStr, "m"), 10, 64); err == nil {
				MaxNodeProvisioningTime = val
			}
		}
	}

	newPodScaleUpDelay := ""
	if profile.NewPodScaleUpDelay != nil {
		newPodScaleUpDelay = pointer.From(profile.NewPodScaleUpDelay)
	}

	scaleDownDelayAfterAdd := int64(0)
	if profile.ScaleDownDelayAfterAdd != nil {
		timeStr := pointer.From(profile.ScaleDownDelayAfterAdd)
		if strings.HasSuffix(timeStr, "m") {
			if val, err := strconv.ParseInt(strings.TrimSuffix(timeStr, "m"), 10, 64); err == nil {
				scaleDownDelayAfterAdd = val
			}
		}
	}

	scaleDownDelayAfterDelete := int64(0)
	if profile.ScaleDownDelayAfterDelete != nil {
		timeStr := pointer.From(profile.ScaleDownDelayAfterDelete)
		if strings.HasSuffix(timeStr, "s") {
			if val, err := strconv.ParseInt(strings.TrimSuffix(timeStr, "s"), 10, 64); err == nil {
				scaleDownDelayAfterDelete = val
			}
		}
	}

	scaleDownDelayAfterFailure := int64(0)
	if profile.ScaleDownDelayAfterFailure != nil {
		timeStr := pointer.From(profile.ScaleDownDelayAfterFailure)
		if strings.HasSuffix(timeStr, "m") {
			if val, err := strconv.ParseInt(strings.TrimSuffix(timeStr, "m"), 10, 64); err == nil {
				scaleDownDelayAfterFailure = val
			}
		}
	}

	scaleDownUnneededTime := int64(0)
	if profile.ScaleDownUnneededTime != nil {
		timeStr := pointer.From(profile.ScaleDownUnneededTime)
		if strings.HasSuffix(timeStr, "m") {
			if val, err := strconv.ParseInt(strings.TrimSuffix(timeStr, "m"), 10, 64); err == nil {
				scaleDownUnneededTime = val
			}
		}
	}

	scaleDownUnreadyTime := int64(0)
	if profile.ScaleDownUnreadyTime != nil {
		timeStr := pointer.From(profile.ScaleDownUnreadyTime)
		if strings.HasSuffix(timeStr, "m") {
			if val, err := strconv.ParseInt(strings.TrimSuffix(timeStr, "m"), 10, 64); err == nil {
				scaleDownUnreadyTime = val
			}
		}
	}

	scaleDownUtilizationThreshold := float64(0)
	if profile.ScaleDownUtilizationThreshold != nil {
		if val, err := strconv.ParseFloat(pointer.From(profile.ScaleDownUtilizationThreshold), 64); err == nil {
			scaleDownUtilizationThreshold = val
		}
	}

	emptyBulkDeleteMax := int64(0)
	if profile.MaxEmptyBulkDelete != nil {
		if val, err := strconv.ParseInt(pointer.From(profile.MaxEmptyBulkDelete), 10, 64); err == nil {
			emptyBulkDeleteMax = val
		}
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

	scanInterval := int64(0)
	if profile.ScanInterval != nil {
		timeStr := pointer.From(profile.ScanInterval)
		if strings.HasSuffix(timeStr, "s") {
			if val, err := strconv.ParseInt(strings.TrimSuffix(timeStr, "s"), 10, 64); err == nil {
				scanInterval = val
			}
		}
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
		annotationsAllowed = pointer.From(input.Metrics.KubeStateMetrics.MetricAnnotationsAllowList)
	}

	labelsAllowed := ""
	if input.Metrics.KubeStateMetrics.MetricLabelsAllowlist != nil {
		labelsAllowed = pointer.From(input.Metrics.KubeStateMetrics.MetricLabelsAllowlist)
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
		httpProxy = pointer.From(httpProxyConfig.HTTPProxy)
	}

	httpsProxy := ""
	if httpProxyConfig.HTTPSProxy != nil {
		httpsProxy = pointer.From(httpProxyConfig.HTTPSProxy)
	}

	noProxyList := make([]string, 0)
	if httpProxyConfig.NoProxy != nil {
		noProxyList = append(noProxyList, pointer.From(httpProxyConfig.NoProxy)...)
	}

	trustedCa := ""
	if httpProxyConfig.TrustedCa != nil {
		trustedCa = pointer.From(httpProxyConfig.TrustedCa)
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

		keyVaultKeyId, err := keyvault.ParseNestedItemID(pointer.From(kms.KeyId), keyvault.VersionTypeVersioned, nestedItemType)
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
		keyVaultKeyID = pointer.From(kms.KeyId)
	}

	keyVaultNetworkAccess := ""
	if kms.KeyVaultNetworkAccess != nil {
		keyVaultNetworkAccess = string(pointer.From(kms.KeyVaultNetworkAccess))
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
	if v := pointer.From(input.Defender.LogAnalyticsWorkspaceResourceId); v != "" {
		logAnalyticsWorkspace = v
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
		dnsZoneIDs = pointer.From(input.WebAppRouting.DnsZoneResourceIds)
	}

	defaultNginxController := managedclusters.NginxIngressControllerTypeAnnotationControlled
	if input.WebAppRouting.Nginx != nil {
		defaultNginxController = pointer.From(input.WebAppRouting.Nginx.DefaultIngressControllerType)
	}

	webAppRoutingIdentity := make([]WebAppRoutingIdentityModel, 0)
	if input.WebAppRouting.Identity != nil {
		clientID := ""
		if input.WebAppRouting.Identity.ClientId != nil {
			clientID = pointer.From(input.WebAppRouting.Identity.ClientId)
		}
		objectID := ""
		if input.WebAppRouting.Identity.ObjectId != nil {
			objectID = pointer.From(input.WebAppRouting.Identity.ObjectId)
		}
		userAssignedIdentityID := ""
		if input.WebAppRouting.Identity.ResourceId != nil {
			userAssignedIdentityID = pointer.From(input.WebAppRouting.Identity.ResourceId)
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

	profile := managedclusters.ServiceMeshProfile{}

	profile.Mode = managedclusters.ServiceMeshModeIstio
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

	return &profile
}

func flattenKubernetesAutomaticClusterServiceMeshProfile(profile *managedclusters.ServiceMeshProfile) []ServiceMeshProfileModel {
	if profile == nil || profile.Mode != managedclusters.ServiceMeshModeIstio || profile.Istio == nil {
		return []ServiceMeshProfileModel{}
	}

	revisions := make([]string, 0)
	if profile.Istio.Revisions != nil {
		revisions = pointer.From(profile.Istio.Revisions)
	}

	internalIngressGatewayEnabled := false
	externalIngressGatewayEnabled := false

	if profile.Istio.Components != nil && profile.Istio.Components.IngressGateways != nil {
		for _, gateway := range pointer.From(profile.Istio.Components.IngressGateways) {
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

	return []CertificateAuthorityModel{{
		KeyVaultID:          pointer.From(plugin.KeyVaultId),
		RootCertObjectName:  pointer.From(plugin.RootCertObjectName),
		CertChainObjectName: pointer.From(plugin.CertChainObjectName),
		CertObjectName:      pointer.From(plugin.CertObjectName),
		KeyObjectName:       pointer.From(plugin.KeyObjectName),
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

	return []UpgradeOverrideModel{{
		ForceUpgradeEnabled: forceUpgradeEnabled,
		EffectiveUntil:      pointer.From(input.OverrideSettings.Until),
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
		result.AdminGroupObjectIDs = pointer.From(profile.AdminGroupObjectIDs)
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

	return []identity.SystemOrUserAssignedList{{
		Type:        input.Type,
		IdentityIds: identityIds,
		PrincipalId: input.PrincipalId,
		TenantId:    input.TenantId,
	}}
}

func expandClusterNodePoolKubeletConfigTyped(input []KubeletConfigModel) *managedclusters.KubeletConfig {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	result := &managedclusters.KubeletConfig{
		CpuCfsQuota:          pointer.To(config.CPUCfsQuotaEnabled),
		FailSwapOn:           pointer.To(false), // must be false to enable swap file on nodes
		AllowedUnsafeSysctls: pointer.To(config.AllowedUnsafeSysctls),
	}

	CPUManagerPolicy := "none"
	if config.CPUManagerPolicy {
		CPUManagerPolicy = "static"
	}
	result.CpuManagerPolicy = pointer.To(CPUManagerPolicy)

	if config.CPUCfsQuotaPeriod != "" {
		result.CpuCfsQuotaPeriod = pointer.To(config.CPUCfsQuotaPeriod)
	}
	if config.ImageGcHighThreshold != 0 {
		result.ImageGcHighThreshold = pointer.To(config.ImageGcHighThreshold)
	}
	if config.ImageGcLowThreshold != 0 {
		result.ImageGcLowThreshold = pointer.To(config.ImageGcLowThreshold)
	}

	TopologyManagerPolicy := "none"
	if v := config.TopologyManagerPolicy; v != "" {
		TopologyManagerPolicy = v
	}
	result.TopologyManagerPolicy = pointer.To(TopologyManagerPolicy)

	if config.ContainerLogMaxSizeMB != 0 {
		result.ContainerLogMaxSizeMB = pointer.To(config.ContainerLogMaxSizeMB)
	}
	if config.ContainerLogMaxFiles != 0 {
		result.ContainerLogMaxFiles = pointer.To(config.ContainerLogMaxFiles)
	}
	if config.PodMaxPid != 0 {
		result.PodMaxPids = pointer.To(config.PodMaxPid)
	}

	return result
}

func expandAutomaticClusterNodePoolLinuxOSConfig(input []LinuxOSConfigModel) (*managedclusters.LinuxOSConfig, error) {
	if len(input) == 0 {
		return nil, nil
	}

	config := input[0]
	sysctlConfig, err := expandClusterNodePoolSysctlConfigTyped(config.SysctlConfig)
	if err != nil {
		return nil, err
	}

	TransparentHugePage := "never"
	if v := config.TransparentHugePage; v != "" {
		TransparentHugePage = v
	}

	TransparentHugePageDefrag := "never"
	if v := config.TransparentHugePageDefrag; v != "" {
		TransparentHugePageDefrag = v
	}

	result := &managedclusters.LinuxOSConfig{
		Sysctls:                    sysctlConfig,
		TransparentHugePageDefrag:  pointer.To(TransparentHugePageDefrag),
		TransparentHugePageEnabled: pointer.To(TransparentHugePage),
	}

	if config.SwapFileSizeMB != 0 {
		result.SwapFileSizeMB = pointer.To(config.SwapFileSizeMB)
	}

	return result, nil
}

func expandClusterNodePoolSysctlConfigTyped(input []SysctlConfigModel) (*managedclusters.SysctlConfig, error) {
	if len(input) == 0 {
		return nil, nil
	}

	config := input[0]
	result := &managedclusters.SysctlConfig{
		NetIPv4TcpTwReuse: pointer.To(config.NetIPv4TCPTwReuse),
	}

	if config.NetCoreSomaxconn != 0 {
		result.NetCoreSomaxconn = pointer.To(config.NetCoreSomaxconn)
	}
	if config.NetCoreNetdevMaxBacklog != 0 {
		result.NetCoreNetdevMaxBacklog = pointer.To(config.NetCoreNetdevMaxBacklog)
	}
	if config.NetCoreRmemDefault != 0 {
		result.NetCoreRmemDefault = pointer.To(config.NetCoreRmemDefault)
	}
	if config.NetCoreRmemMax != 0 {
		result.NetCoreRmemMax = pointer.To(config.NetCoreRmemMax)
	}
	if config.NetCoreWmemDefault != 0 {
		result.NetCoreWmemDefault = pointer.To(config.NetCoreWmemDefault)
	}
	if config.NetCoreWmemMax != 0 {
		result.NetCoreWmemMax = pointer.To(config.NetCoreWmemMax)
	}
	if config.NetCoreOptmemMax != 0 {
		result.NetCoreOptmemMax = pointer.To(config.NetCoreOptmemMax)
	}
	if config.NetIPv4TCPMaxSynBacklog != 0 {
		result.NetIPv4TcpMaxSynBacklog = pointer.To(config.NetIPv4TCPMaxSynBacklog)
	}
	if config.NetIPv4TCPMaxTwBuckets != 0 {
		result.NetIPv4TcpMaxTwBuckets = pointer.To(config.NetIPv4TCPMaxTwBuckets)
	}
	if config.NetIPv4TCPFinTimeout != 0 {
		result.NetIPv4TcpFinTimeout = pointer.To(config.NetIPv4TCPFinTimeout)
	}
	if config.NetIPv4TCPKeepaliveTime != 0 {
		result.NetIPv4TcpKeepaliveTime = pointer.To(config.NetIPv4TCPKeepaliveTime)
	}
	if config.NetIPv4TCPKeepaliveProbes != 0 {
		result.NetIPv4TcpKeepaliveProbes = pointer.To(config.NetIPv4TCPKeepaliveProbes)
	}
	if config.NetIPv4TCPKeepaliveIntvl != 0 {
		result.NetIPv4TcpkeepaliveIntvl = pointer.To(config.NetIPv4TCPKeepaliveIntvl)
	}

	if (config.NetIPv4IPLocalPortRangeMin != 0 && config.NetIPv4IPLocalPortRangeMax == 0) ||
		(config.NetIPv4IPLocalPortRangeMin == 0 && config.NetIPv4IPLocalPortRangeMax != 0) {
		return nil, fmt.Errorf("`net_ipv4_ip_local_port_range_min` and `net_ipv4_ip_local_port_range_maximum` should both be set or unset")
	}
	if config.NetIPv4IPLocalPortRangeMin > config.NetIPv4IPLocalPortRangeMax {
		return nil, fmt.Errorf("`net_ipv4_ip_local_port_range_min` should be no larger than `net_ipv4_ip_local_port_range_maximum`")
	}
	if config.NetIPv4IPLocalPortRangeMin != 0 && config.NetIPv4IPLocalPortRangeMax != 0 {
		result.NetIPv4IPLocalPortRange = pointer.To(fmt.Sprintf("%d %d", config.NetIPv4IPLocalPortRangeMin, config.NetIPv4IPLocalPortRangeMax))
	}

	if config.NetIPv4NeighDefaultGcThresh1 != 0 {
		result.NetIPv4NeighDefaultGcThresh1 = pointer.To(config.NetIPv4NeighDefaultGcThresh1)
	}
	if config.NetIPv4NeighDefaultGcThresh2 != 0 {
		result.NetIPv4NeighDefaultGcThresh2 = pointer.To(config.NetIPv4NeighDefaultGcThresh2)
	}
	if config.NetIPv4NeighDefaultGcThresh3 != 0 {
		result.NetIPv4NeighDefaultGcThresh3 = pointer.To(config.NetIPv4NeighDefaultGcThresh3)
	}
	if config.NetNetfilterNfConntrackMax != 0 {
		result.NetNetfilterNfConntrackMax = pointer.To(config.NetNetfilterNfConntrackMax)
	}
	if config.NetNetfilterNfConntrackBuckets != 0 {
		result.NetNetfilterNfConntrackBuckets = pointer.To(config.NetNetfilterNfConntrackBuckets)
	}
	if config.FsAioMaxNr != 0 {
		result.FsAioMaxNr = pointer.To(config.FsAioMaxNr)
	}
	if config.FsInotifyMaxUserWatches != 0 {
		result.FsInotifyMaxUserWatches = pointer.To(config.FsInotifyMaxUserWatches)
	}
	if config.FsFileMax != 0 {
		result.FsFileMax = pointer.To(config.FsFileMax)
	}
	if config.FsNrOpen != 0 {
		result.FsNrOpen = pointer.To(config.FsNrOpen)
	}
	if config.KernelThreadsMax != 0 {
		result.KernelThreadsMax = pointer.To(config.KernelThreadsMax)
	}
	if config.VMMaxMapCount != 0 {
		result.VMMaxMapCount = pointer.To(config.VMMaxMapCount)
	}
	if config.VMSwappiness != 0 {
		result.VMSwappiness = pointer.To(config.VMSwappiness)
	}
	if config.VMVfsCachePressure != 0 {
		result.VMVfsCachePressure = pointer.To(config.VMVfsCachePressure)
	}

	return result, nil
}

func expandClusterNodePoolUpgradeSettingsTyped(input []UpgradeSettingsModel) *managedclusters.AgentPoolUpgradeSettings {
	if len(input) == 0 {
		return nil
	}

	config := input[0]
	result := &managedclusters.AgentPoolUpgradeSettings{}

	if config.MaxSurge != "" {
		result.MaxSurge = pointer.To(config.MaxSurge)
	}
	if config.DrainTimeoutInMinutes != 0 {
		result.DrainTimeoutInMinutes = pointer.To(config.DrainTimeoutInMinutes)
	}
	if config.NodeSoakDurationInMinutes != 0 {
		result.NodeSoakDurationInMinutes = pointer.To(config.NodeSoakDurationInMinutes)
	}

	if config.UndrainableNodeBehavior != "" {
		result.UndrainableNodeBehavior = pointer.To(managedclusters.UndrainableNodeBehavior(config.UndrainableNodeBehavior))
	}

	return result
}

func flattenClusterNodePoolKubeletConfigTyped(input *managedclusters.KubeletConfig) []KubeletConfigModel {
	if input == nil {
		return []KubeletConfigModel{}
	}

	cpuCfsQuotaEnabled := false
	if input.CpuCfsQuota != nil {
		cpuCfsQuotaEnabled = pointer.From(input.CpuCfsQuota)
	}

	allowedUnsafeSysctls := []string{}
	if input.AllowedUnsafeSysctls != nil {
		allowedUnsafeSysctls = pointer.From(input.AllowedUnsafeSysctls)
	}

	result := KubeletConfigModel{
		CPUCfsQuotaEnabled:   cpuCfsQuotaEnabled,
		AllowedUnsafeSysctls: allowedUnsafeSysctls,
	}

	CpuManagerPolicy := false
	if v := input.CpuManagerPolicy; v != nil && v != pointer.To("none") {
		CpuManagerPolicy = true
	}
	result.CPUManagerPolicy = CpuManagerPolicy

	if input.CpuCfsQuotaPeriod != nil {
		result.CPUCfsQuotaPeriod = pointer.From(input.CpuCfsQuotaPeriod)
	}
	if input.ImageGcHighThreshold != nil {
		result.ImageGcHighThreshold = pointer.From(input.ImageGcHighThreshold)
	}
	if input.ImageGcLowThreshold != nil {
		result.ImageGcLowThreshold = pointer.From(input.ImageGcLowThreshold)
	}

	TopologyManagerPolicy := pointer.To("")
	if v := input.TopologyManagerPolicy; v != nil && v != pointer.To("none") {
		TopologyManagerPolicy = v
	}
	result.TopologyManagerPolicy = pointer.From(TopologyManagerPolicy)

	if input.ContainerLogMaxSizeMB != nil {
		result.ContainerLogMaxSizeMB = pointer.From(input.ContainerLogMaxSizeMB)
	}
	if input.ContainerLogMaxFiles != nil {
		result.ContainerLogMaxFiles = pointer.From(input.ContainerLogMaxFiles)
	}
	if input.PodMaxPids != nil {
		result.PodMaxPid = *input.PodMaxPids
	}

	return []KubeletConfigModel{result}
}

func flattenAutomaticClusterNodePoolLinuxOSConfig(input *managedclusters.LinuxOSConfig) []LinuxOSConfigModel {
	if input == nil {
		return []LinuxOSConfigModel{}
	}

	sysctlConfig := flattenClusterNodePoolSysctlConfigTyped(input.Sysctls)

	result := LinuxOSConfigModel{
		SysctlConfig: sysctlConfig,
	}

	TransparentHugePageEnabled := pointer.To("")
	if v := input.TransparentHugePageEnabled; v != nil && v != pointer.To("never") {
		TransparentHugePageEnabled = v
	}
	result.TransparentHugePage = pointer.From(TransparentHugePageEnabled)

	TransparentHugePageDefrag := pointer.To("")
	if v := input.TransparentHugePageDefrag; v != nil && v != pointer.To("never") {
		TransparentHugePageDefrag = v
	}
	result.TransparentHugePageDefrag = pointer.From(TransparentHugePageDefrag)

	if input.SwapFileSizeMB != nil {
		result.SwapFileSizeMB = pointer.From(input.SwapFileSizeMB)
	}

	return []LinuxOSConfigModel{result}
}

func flattenClusterNodePoolSysctlConfigTyped(input *managedclusters.SysctlConfig) []SysctlConfigModel {
	if input == nil {
		return []SysctlConfigModel{}
	}

	netIPv4TcpTwReuse := false
	if input.NetIPv4TcpTwReuse != nil {
		netIPv4TcpTwReuse = pointer.From(input.NetIPv4TcpTwReuse)
	}

	result := SysctlConfigModel{
		NetIPv4TCPTwReuse: netIPv4TcpTwReuse,
	}

	if input.NetCoreSomaxconn != nil {
		result.NetCoreSomaxconn = pointer.From(input.NetCoreSomaxconn)
	}
	if input.NetCoreNetdevMaxBacklog != nil {
		result.NetCoreNetdevMaxBacklog = pointer.From(input.NetCoreNetdevMaxBacklog)
	}
	if input.NetCoreRmemDefault != nil {
		result.NetCoreRmemDefault = pointer.From(input.NetCoreRmemDefault)
	}
	if input.NetCoreRmemMax != nil {
		result.NetCoreRmemMax = pointer.From(input.NetCoreRmemMax)
	}
	if input.NetCoreWmemDefault != nil {
		result.NetCoreWmemDefault = pointer.From(input.NetCoreWmemDefault)
	}
	if input.NetCoreWmemMax != nil {
		result.NetCoreWmemMax = pointer.From(input.NetCoreWmemMax)
	}
	if input.NetCoreOptmemMax != nil {
		result.NetCoreOptmemMax = pointer.From(input.NetCoreOptmemMax)
	}
	if input.NetIPv4TcpMaxSynBacklog != nil {
		result.NetIPv4TCPMaxSynBacklog = pointer.From(input.NetIPv4TcpMaxSynBacklog)
	}
	if input.NetIPv4TcpMaxTwBuckets != nil {
		result.NetIPv4TCPMaxTwBuckets = pointer.From(input.NetIPv4TcpMaxTwBuckets)
	}
	if input.NetIPv4TcpFinTimeout != nil {
		result.NetIPv4TCPFinTimeout = pointer.From(input.NetIPv4TcpFinTimeout)
	}
	if input.NetIPv4TcpKeepaliveTime != nil {
		result.NetIPv4TCPKeepaliveTime = pointer.From(input.NetIPv4TcpKeepaliveTime)
	}
	if input.NetIPv4TcpKeepaliveProbes != nil {
		result.NetIPv4TCPKeepaliveProbes = pointer.From(input.NetIPv4TcpKeepaliveProbes)
	}
	if input.NetIPv4TcpkeepaliveIntvl != nil {
		result.NetIPv4TCPKeepaliveIntvl = pointer.From(input.NetIPv4TcpkeepaliveIntvl)
	}

	if input.NetIPv4IPLocalPortRange != nil {
		portRange := *input.NetIPv4IPLocalPortRange
		var min, max int64
		if _, err := fmt.Sscanf(portRange, "%d %d", &min, &max); err == nil {
			result.NetIPv4IPLocalPortRangeMin = min
			result.NetIPv4IPLocalPortRangeMax = max
		}
	}

	if input.NetIPv4NeighDefaultGcThresh1 != nil {
		result.NetIPv4NeighDefaultGcThresh1 = pointer.From(input.NetIPv4NeighDefaultGcThresh1)
	}
	if input.NetIPv4NeighDefaultGcThresh2 != nil {
		result.NetIPv4NeighDefaultGcThresh2 = pointer.From(input.NetIPv4NeighDefaultGcThresh2)
	}
	if input.NetIPv4NeighDefaultGcThresh3 != nil {
		result.NetIPv4NeighDefaultGcThresh3 = pointer.From(input.NetIPv4NeighDefaultGcThresh3)
	}
	if input.NetNetfilterNfConntrackMax != nil {
		result.NetNetfilterNfConntrackMax = pointer.From(input.NetNetfilterNfConntrackMax)
	}
	if input.NetNetfilterNfConntrackBuckets != nil {
		result.NetNetfilterNfConntrackBuckets = pointer.From(input.NetNetfilterNfConntrackBuckets)
	}
	if input.FsAioMaxNr != nil {
		result.FsAioMaxNr = pointer.From(input.FsAioMaxNr)
	}
	if input.FsInotifyMaxUserWatches != nil {
		result.FsInotifyMaxUserWatches = pointer.From(input.FsInotifyMaxUserWatches)
	}
	if input.FsFileMax != nil {
		result.FsFileMax = pointer.From(input.FsFileMax)
	}
	if input.FsNrOpen != nil {
		result.FsNrOpen = pointer.From(input.FsNrOpen)
	}
	if input.KernelThreadsMax != nil {
		result.KernelThreadsMax = pointer.From(input.KernelThreadsMax)
	}
	if input.VMMaxMapCount != nil {
		result.VMMaxMapCount = pointer.From(input.VMMaxMapCount)
	}
	if input.VMSwappiness != nil {
		result.VMSwappiness = pointer.From(input.VMSwappiness)
	}
	if input.VMVfsCachePressure != nil {
		result.VMVfsCachePressure = pointer.From(input.VMVfsCachePressure)
	}

	return []SysctlConfigModel{result}
}

func flattenClusterNodePoolUpgradeSettingsTyped(input *managedclusters.AgentPoolUpgradeSettings) []UpgradeSettingsModel {
	if input == nil || (input.MaxSurge == nil && input.DrainTimeoutInMinutes == nil && input.NodeSoakDurationInMinutes == nil) {
		return []UpgradeSettingsModel{}
	}

	result := UpgradeSettingsModel{}

	if input.MaxSurge != nil {
		result.MaxSurge = *input.MaxSurge
	}
	if input.DrainTimeoutInMinutes != nil {
		result.DrainTimeoutInMinutes = pointer.From(input.DrainTimeoutInMinutes)
	}
	if input.NodeSoakDurationInMinutes != nil {
		result.NodeSoakDurationInMinutes = pointer.From(input.NodeSoakDurationInMinutes)
	}

	if input.UndrainableNodeBehavior != nil {
		result.UndrainableNodeBehavior = string(pointer.From(input.UndrainableNodeBehavior))
	}

	return []UpgradeSettingsModel{result}
}

func findDefaultNodePoolTyped(input *[]managedclusters.ManagedClusterAgentPoolProfile) (*managedclusters.ManagedClusterAgentPoolProfile, error) {
	if input == nil {
		return nil, fmt.Errorf("agent pool profiles is nil")
	}

	var agentPool *managedclusters.ManagedClusterAgentPoolProfile
	for _, v := range *input {
		if v.Name == "" {
			continue
		}
		if v.Mode == nil || *v.Mode != managedclusters.AgentPoolModeSystem {
			continue
		}

		agentPool = &v
		break
	}

	if agentPool == nil {
		return nil, fmt.Errorf("unable to determine default agent pool - no System mode pool found")
	}

	return agentPool, nil
}

func flattenClusterPoolNetworkProfileTyped(input *managedclusters.AgentPoolNetworkProfile) []NodeNetworkProfileModel {
	if input == nil || (input.NodePublicIPTags == nil && input.AllowedHostPorts == nil && input.ApplicationSecurityGroups == nil) {
		return []NodeNetworkProfileModel{}
	}
	results := make([]NodeNetworkProfileModel, 0)
	result := NodeNetworkProfileModel{
		AllowedHostPorts:            flattenClusterPoolNetworkProfileAllowedHostPortsTyped(input.AllowedHostPorts),
		ApplicationSecurityGroupIDs: []string{},
		NodePublicIPTags:            flattenClusterPoolNetworkProfileNodePublicIPTagsTyped(input.NodePublicIPTags),
	}

	if input.ApplicationSecurityGroups != nil {
		result.ApplicationSecurityGroupIDs = pointer.From(input.ApplicationSecurityGroups)
	}
	results = append(results, result)
	return results
}

func flattenClusterPoolNetworkProfileAllowedHostPortsTyped(input *[]managedclusters.PortRange) []AllowedHostPortsModel {
	if input == nil {
		return []AllowedHostPortsModel{}
	}

	result := make([]AllowedHostPortsModel, 0)
	for _, portRange := range *input {
		model := AllowedHostPortsModel{}
		if portRange.PortEnd != nil {
			model.PortEnd = pointer.From(portRange.PortEnd)
		}
		if portRange.PortStart != nil {
			model.PortStart = pointer.From(portRange.PortStart)
		}
		if portRange.Protocol != nil {
			model.Protocol = string(*portRange.Protocol)
		}
		result = append(result, model)
	}
	return result
}

func flattenClusterPoolNetworkProfileNodePublicIPTagsTyped(input *[]managedclusters.IPTag) map[string]string {
	if input == nil {
		return map[string]string{}
	}

	result := make(map[string]string)
	for _, tag := range *input {
		if tag.IPTagType != nil && tag.Tag != nil {
			result[*tag.IPTagType] = *tag.Tag
		}
	}

	return result
}

func expandClusterPoolNetworkProfileTyped(input []NodeNetworkProfileModel) *managedclusters.AgentPoolNetworkProfile {
	if len(input) == 0 {
		return nil
	}

	profile := input[0]
	result := &managedclusters.AgentPoolNetworkProfile{
		AllowedHostPorts:          expandClusterPoolNetworkProfileAllowedHostPortsTyped(profile.AllowedHostPorts),
		ApplicationSecurityGroups: pointer.To(profile.ApplicationSecurityGroupIDs),
		NodePublicIPTags:          expandClusterPoolNetworkProfileNodePublicIPTagsTyped(profile.NodePublicIPTags),
	}

	return result
}

func expandClusterPoolNetworkProfileAllowedHostPortsTyped(input []AllowedHostPortsModel) *[]managedclusters.PortRange {
	if len(input) == 0 {
		return nil
	}

	out := make([]managedclusters.PortRange, 0, len(input))
	for _, v := range input {
		out = append(out, managedclusters.PortRange{
			PortEnd:   pointer.To(v.PortEnd),
			PortStart: pointer.To(v.PortStart),
			Protocol:  pointer.To(managedclusters.Protocol(v.Protocol)),
		})
	}
	return &out
}

func expandClusterPoolNetworkProfileNodePublicIPTagsTyped(input map[string]string) *[]managedclusters.IPTag {
	if len(input) == 0 {
		return nil
	}

	out := make([]managedclusters.IPTag, 0, len(input))
	for key, val := range input {
		ipTag := managedclusters.IPTag{
			IPTagType: pointer.To(key),
			Tag:       pointer.To(val),
		}
		out = append(out, ipTag)
	}
	return &out
}

// ExpandDefaultNodePoolTyped converts a DefaultNodePoolModel to ManagedClusterAgentPoolProfile
func ExpandDefaultNodePoolTyped(input []DefaultNodePoolModel) (*[]managedclusters.ManagedClusterAgentPoolProfile, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("default_node_pool must be specified")
	}

	raw := input[0]

	nodeLabels := pointer.To(raw.NodeLabels)
	var nodeTaints *[]string

	profile := managedclusters.ManagedClusterAgentPoolProfile{
		EnableFIPS:             pointer.To(raw.FipsEnabled),
		EnableNodePublicIP:     pointer.To(raw.NodePublicIPEnabled),
		EnableEncryptionAtHost: pointer.To(raw.HostEncryptionEnabled),
		KubeletDiskType:        pointer.To(managedclusters.KubeletDiskType(raw.KubeletDiskType)),
		Name:                   raw.Name,
		NodeLabels:             nodeLabels,
		NodeTaints:             nodeTaints,
		Tags:                   tags.Expand(raw.Tags),
		Type:                   pointer.To(managedclusters.AgentPoolTypeVirtualMachineScaleSets),
		VMSize:                 pointer.To(raw.VMSize),

		// at this time the default node pool has to be Linux or the AKS cluster fails to provision with:
		// Pods not in Running status: coredns-7fc597cc45-v5z7x,coredns-autoscaler-7ccc76bfbd-djl7j,metrics-server-cbd95f966-5rl97,tunnelfront-7d9884977b-wpbvn
		// Windows agents can be configured via the separate node pool resource
		OsType: pointer.To(managedclusters.OSTypeLinux),

		// without this set the API returns:
		// Code="MustDefineAtLeastOneSystemPool" Message="Must define at least one system pool."
		// since this is the "default" node pool we can assume this is a system node pool
		Mode: pointer.To(managedclusters.AgentPoolModeSystem),

		UpgradeSettings: expandClusterNodePoolUpgradeSettingsTyped(raw.UpgradeSettings),
	}

	if raw.MaxPods > 0 {
		profile.MaxPods = pointer.To(raw.MaxPods)
	}

	if raw.NodePublicIPPrefixID != "" {
		profile.NodePublicIPPrefixID = pointer.To(raw.NodePublicIPPrefixID)
	}

	if raw.OSDiskSizeGB > 0 {
		profile.OsDiskSizeGB = pointer.To(raw.OSDiskSizeGB)
	}

	if raw.OSSKU != "" {
		profile.OsSKU = pointer.To(managedclusters.OSSKU(raw.OSSKU))
	}

	profile.ScaleDownMode = pointer.To(managedclusters.ScaleDownModeDelete)

	if raw.SnapshotID != "" {
		profile.CreationData = &managedclusters.CreationData{
			SourceResourceId: pointer.To(raw.SnapshotID),
		}
	}

	if raw.UltraSSDEnabled {
		profile.EnableUltraSSD = pointer.To(raw.UltraSSDEnabled)
	}

	if raw.VnetSubnetID != "" {
		profile.VnetSubnetID = pointer.To(raw.VnetSubnetID)
	}

	if raw.HostGroupID != "" {
		profile.HostGroupID = pointer.To(raw.HostGroupID)
	}

	if raw.OrchestratorVersion != "" {
		profile.OrchestratorVersion = pointer.To(raw.OrchestratorVersion)
	}

	if raw.ProximityPlacementGroupID != "" {
		profile.ProximityPlacementGroupID = pointer.To(raw.ProximityPlacementGroupID)
	}

	profile.WorkloadRuntime = pointer.To(managedclusters.WorkloadRuntimeOCIContainer)

	if raw.CapacityReservationGroupID != "" {
		profile.CapacityReservationGroupID = pointer.To(raw.CapacityReservationGroupID)
	}

	if raw.GPUInstance != "" {
		profile.GpuInstanceProfile = pointer.To(managedclusters.GPUInstanceProfile(raw.GPUInstance))
	}

	if raw.GPUDriver {
		profile.GpuProfile = &managedclusters.GPUProfile{
			Driver: pointer.To(managedclusters.GPUDriverInstall),
		}
	}

	profile.Count = pointer.To(raw.NodeCount)

	if len(raw.KubeletConfig) > 0 {
		profile.KubeletConfig = expandClusterNodePoolKubeletConfigTyped(raw.KubeletConfig)
	}

	if len(raw.LinuxOSConfig) > 0 {
		linuxOSConfig, err := expandAutomaticClusterNodePoolLinuxOSConfig(raw.LinuxOSConfig)
		if err != nil {
			return nil, err
		}
		profile.LinuxOSConfig = linuxOSConfig
	}

	if len(raw.NodeNetworkProfile) > 0 {
		profile.NetworkProfile = expandClusterPoolNetworkProfileTyped(raw.NodeNetworkProfile)
	}

	return &[]managedclusters.ManagedClusterAgentPoolProfile{
		profile,
	}, nil
}

func FlattenDefaultNodePoolTyped(input *[]managedclusters.ManagedClusterAgentPoolProfile, metadata *sdk.ResourceMetaData) ([]DefaultNodePoolModel, error) {
	if input == nil {
		return []DefaultNodePoolModel{}, nil
	}

	agentPool, err := findDefaultNodePoolTyped(input)
	if err != nil {
		return nil, err
	}

	result := DefaultNodePoolModel{
		Name: agentPool.Name,
	}

	// Preserve temporary_name_for_rotation from existing state since it's not returned by the API
	if metadata != nil {
		var existingModel KubernetesAutomaticClusterModel
		if err := metadata.Decode(&existingModel); err == nil {
			if len(existingModel.DefaultNodePool) > 0 {
				result.TemporaryNameForRotation = existingModel.DefaultNodePool[0].TemporaryNameForRotation
			}
		}
	}

	if agentPool.Count != nil {
		result.NodeCount = pointer.From(agentPool.Count)
	}

	if agentPool.EnableUltraSSD != nil {
		result.UltraSSDEnabled = pointer.From(agentPool.EnableUltraSSD)
	}

	if agentPool.EnableFIPS != nil {
		result.FipsEnabled = pointer.From(agentPool.EnableFIPS)
	}

	if agentPool.EnableNodePublicIP != nil {
		result.NodePublicIPEnabled = pointer.From(agentPool.EnableNodePublicIP)
	}

	if agentPool.EnableEncryptionAtHost != nil {
		result.HostEncryptionEnabled = pointer.From(agentPool.EnableEncryptionAtHost)
	}

	if agentPool.GpuInstanceProfile != nil {
		result.GPUInstance = string(pointer.From(agentPool.GpuInstanceProfile))
	}

	if agentPool.GpuProfile != nil && agentPool.GpuProfile.Driver != nil {
		if pointer.From(agentPool.GpuProfile.Driver) == managedclusters.GPUDriverInstall {
			result.GPUDriver = true
		}
	}

	if agentPool.MaxPods != nil {
		result.MaxPods = pointer.From(agentPool.MaxPods)
	}

	if agentPool.NodeLabels != nil {
		result.NodeLabels = make(map[string]string)
		for k, v := range pointer.From(agentPool.NodeLabels) {
			result.NodeLabels[k] = v
		}
	}

	if agentPool.NodePublicIPPrefixID != nil {
		result.NodePublicIPPrefixID = pointer.From(agentPool.NodePublicIPPrefixID)
	}

	if agentPool.OsDiskSizeGB != nil {
		result.OSDiskSizeGB = pointer.From(agentPool.OsDiskSizeGB)
	}

	if agentPool.VnetSubnetID != nil {
		result.VnetSubnetID = pointer.From(agentPool.VnetSubnetID)
	}

	if agentPool.HostGroupID != nil {
		result.HostGroupID = pointer.From(agentPool.HostGroupID)
	}

	// NOTE: workaround for migration from 2022-01-02-preview (<3.12.0) to 2022-03-02-preview (>=3.12.0)
	// Before terraform apply is run against the new API, Azure will respond only with currentOrchestratorVersion
	if agentPool.OrchestratorVersion != nil {
		result.OrchestratorVersion = pointer.From(agentPool.OrchestratorVersion)
	} else if agentPool.CurrentOrchestratorVersion != nil {
		result.OrchestratorVersion = pointer.From(agentPool.CurrentOrchestratorVersion)
	}

	if agentPool.ProximityPlacementGroupID != nil {
		result.ProximityPlacementGroupID = pointer.From(agentPool.ProximityPlacementGroupID)
	}

	if agentPool.CreationData != nil && agentPool.CreationData.SourceResourceId != nil {
		id, err := snapshots.ParseSnapshotIDInsensitively(pointer.From(agentPool.CreationData.SourceResourceId))
		if err != nil {
			return nil, err
		}
		result.SnapshotID = id.ID()
	}

	if agentPool.VMSize != nil {
		result.VMSize = pointer.From(agentPool.VMSize)
	}

	if agentPool.CapacityReservationGroupID != nil {
		result.CapacityReservationGroupID = pointer.From(agentPool.CapacityReservationGroupID)
	}

	if agentPool.KubeletDiskType != nil {
		result.KubeletDiskType = string(pointer.From(agentPool.KubeletDiskType))
	}

	if agentPool.OsSKU != nil {
		result.OSSKU = string(pointer.From(agentPool.OsSKU))
	}

	result.UpgradeSettings = flattenClusterNodePoolUpgradeSettingsTyped(agentPool.UpgradeSettings)

	result.LinuxOSConfig = flattenAutomaticClusterNodePoolLinuxOSConfig(agentPool.LinuxOSConfig)

	result.KubeletConfig = flattenClusterNodePoolKubeletConfigTyped(agentPool.KubeletConfig)

	result.NodeNetworkProfile = flattenClusterPoolNetworkProfileTyped(agentPool.NetworkProfile)

	if agentPool.Tags != nil {
		result.Tags = tags.Flatten(agentPool.Tags)
	}

	return []DefaultNodePoolModel{result}, nil
}

func expandKubernetesAutomaticClusterAddOns(input *KubernetesAutomaticClusterModel, env environments.Environment) (*map[string]managedclusters.ManagedClusterAddonProfile, error) {
	addonProfiles := map[string]managedclusters.ManagedClusterAddonProfile{}

	addonProfiles[openServiceMeshKey] = managedclusters.ManagedClusterAddonProfile{
		Enabled: input.OpenServiceMeshEnabled,
	}

	addonProfiles[confidentialComputingKey] = managedclusters.ManagedClusterAddonProfile{
		Enabled: false,
	}
	if len(input.ConfidentialComputing) > 0 {
		cc := input.ConfidentialComputing[0]
		config := make(map[string]string)
		quoteHelperEnabled := "false"
		if cc.SGXQuoteHelperEnabled {
			quoteHelperEnabled = "true"
		}
		config["ACCSGXQuoteHelperEnabled"] = quoteHelperEnabled
		addonProfiles[confidentialComputingKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  pointer.To(config),
		}
	}

	if input.HTTPApplicationRoutingEnabled {
		addonProfiles[httpApplicationRoutingKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: input.HTTPApplicationRoutingEnabled,
		}
	}

	addonProfiles[omsAgentKey] = managedclusters.ManagedClusterAddonProfile{
		Enabled: false,
	}
	if len(input.OMSAgent) > 0 {
		oms := input.OMSAgent[0]
		config := make(map[string]string)

		if oms.LogAnalyticsWorkspaceID != "" {
			lawid, err := workspaces.ParseWorkspaceIDInsensitively(oms.LogAnalyticsWorkspaceID)
			if err != nil {
				return nil, fmt.Errorf("parsing Log Analytics Workspace ID: %w", err)
			}
			config["logAnalyticsWorkspaceResourceID"] = lawid.ID()
		}

		if oms.MSIAuthForMonitoringEnabled != nil {
			config["useAADAuth"] = fmt.Sprintf("%t", *oms.MSIAuthForMonitoringEnabled)
		}

		addonProfiles[omsAgentKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  pointer.To(config),
		}
	}

	addonProfiles[aciConnectorKey] = managedclusters.ManagedClusterAddonProfile{
		Enabled: false,
	}
	if len(input.ACIConnectorLinux) > 0 {
		aci := input.ACIConnectorLinux[0]
		config := make(map[string]string)

		if aci.SubnetName != "" {
			config["SubnetName"] = aci.SubnetName
		}

		addonProfiles[aciConnectorKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  pointer.To(config),
		}
	}

	addonProfiles[ingressApplicationGatewayKey] = managedclusters.ManagedClusterAddonProfile{
		Enabled: false,
	}
	if len(input.IngressApplicationGateway) > 0 {
		iag := input.IngressApplicationGateway[0]
		config := make(map[string]string)

		if iag.GatewayID != "" {
			config["applicationGatewayId"] = iag.GatewayID
		}

		if iag.GatewayName != "" {
			config["applicationGatewayName"] = iag.GatewayName
		}

		if iag.SubnetCIDR != "" {
			config["subnetCIDR"] = iag.SubnetCIDR
		}

		if iag.SubnetID != "" {
			config["subnetId"] = iag.SubnetID
		}

		addonProfiles[ingressApplicationGatewayKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  pointer.To(config),
		}
	}

	if len(input.KeyVaultSecretsProvider) > 0 {
		kvsp := input.KeyVaultSecretsProvider[0]
		config := make(map[string]string)

		config["enableSecretRotation"] = fmt.Sprintf("%t", true)
		config["rotationPollInterval"] = kvsp.SecretRotationInterval

		addonProfiles[azureKeyvaultSecretsProviderKey] = managedclusters.ManagedClusterAddonProfile{
			Enabled: true,
			Config:  pointer.To(config),
		}
	}

	return filterUnsupportedKubernetesAddOns(addonProfiles, env)
}

func flattenKubernetesAutomaticClusterAddOns(profile map[string]managedclusters.ManagedClusterAddonProfile) (
	aciConnectorLinux []ACIConnectorLinuxModel,
	confidentialComputing []ConfidentialComputingModel,
	httpApplicationRoutingEnabled bool,
	httpApplicationRoutingZoneName string,
	ingressApplicationGateway []IngressApplicationGatewayModel,
	keyVaultSecretsProvider []KeyVaultSecretsProviderModel,
	omsAgent []OMSAgentModel,
	openServiceMeshEnabled bool,
) {
	aciConnector := kubernetesAddonProfileLocate(profile, aciConnectorKey)
	if aciConnector.Enabled {
		subnetName := ""
		if v := aciConnector.Config; v != nil && (*v)["SubnetName"] != "" {
			subnetName = (*v)["SubnetName"]
		}

		identity := flattenKubernetesClusterAddOnIdentityProfileTyped(aciConnector.Identity)

		aciConnectorLinux = []ACIConnectorLinuxModel{{
			SubnetName:        subnetName,
			ConnectorIdentity: identity,
		}}
	}

	confidentialComputingProfile := kubernetesAddonProfileLocate(profile, confidentialComputingKey)
	if confidentialComputingProfile.Enabled {
		quoteHelperEnabled := false
		if v := kubernetesAddonProfilelocateInConfig(confidentialComputingProfile.Config, "ACCSGXQuoteHelperEnabled"); v != "" && v != "false" {
			quoteHelperEnabled = true
		}
		confidentialComputing = []ConfidentialComputingModel{{
			SGXQuoteHelperEnabled: quoteHelperEnabled,
		}}
	}

	httpApplicationRouting := kubernetesAddonProfileLocate(profile, httpApplicationRoutingKey)
	httpApplicationRoutingEnabled = httpApplicationRouting.Enabled

	if v := kubernetesAddonProfilelocateInConfig(httpApplicationRouting.Config, "HTTPApplicationRoutingZoneName"); v != "" {
		httpApplicationRoutingZoneName = v
	}

	omsAgentProfile := kubernetesAddonProfileLocate(profile, omsAgentKey)
	if omsAgentProfile.Enabled {
		workspaceID := ""
		useAADAuth := false

		if v := kubernetesAddonProfilelocateInConfig(omsAgentProfile.Config, "logAnalyticsWorkspaceResourceID"); v != "" {
			if lawid, err := workspaces.ParseWorkspaceIDInsensitively(v); err == nil {
				workspaceID = lawid.ID()
			}
		}

		if v := kubernetesAddonProfilelocateInConfig(omsAgentProfile.Config, "useAADAuth"); v != "false" && v != "" {
			useAADAuth = true
		}

		omsAgentIdentity := flattenKubernetesClusterAddOnIdentityProfileTyped(omsAgentProfile.Identity)

		omsAgent = []OMSAgentModel{{
			LogAnalyticsWorkspaceID:     workspaceID,
			MSIAuthForMonitoringEnabled: pointer.To(useAADAuth),
			OMSAgentIdentity:            flattenOMSAgentIdentityTyped(omsAgentIdentity),
		}}
	}

	ingressApplicationGatewayProfile := kubernetesAddonProfileLocate(profile, ingressApplicationGatewayKey)
	if ingressApplicationGatewayProfile.Enabled {
		gatewayId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGatewayProfile.Config, "applicationGatewayId"); v != "" {
			gatewayId = v
		}

		gatewayName := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGatewayProfile.Config, "applicationGatewayName"); v != "" {
			gatewayName = v
		}

		effectiveGatewayId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGatewayProfile.Config, "effectiveApplicationGatewayId"); v != "" {
			effectiveGatewayId = v
		}

		subnetCIDR := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGatewayProfile.Config, "subnetCIDR"); v != "" {
			subnetCIDR = v
		}

		subnetId := ""
		if v := kubernetesAddonProfilelocateInConfig(ingressApplicationGatewayProfile.Config, "subnetId"); v != "" {
			subnetId = v
		}

		ingressApplicationGatewayIdentity := flattenKubernetesClusterAddOnIdentityProfileTyped(ingressApplicationGatewayProfile.Identity)

		ingressApplicationGateway = []IngressApplicationGatewayModel{{
			GatewayID:                         gatewayId,
			GatewayName:                       gatewayName,
			EffectiveGatewayID:                effectiveGatewayId,
			SubnetCIDR:                        subnetCIDR,
			SubnetID:                          subnetId,
			IngressApplicationGatewayIdentity: flattenIngressApplicationGatewayIdentityTyped(ingressApplicationGatewayIdentity),
		}}
	}

	openServiceMesh := kubernetesAddonProfileLocate(profile, openServiceMeshKey)
	openServiceMeshEnabled = openServiceMesh.Enabled

	azureKeyVaultSecretsProviderProfile := kubernetesAddonProfileLocate(profile, azureKeyvaultSecretsProviderKey)
	if azureKeyVaultSecretsProviderProfile.Enabled {
		rotationPollInterval := ""

		if v := kubernetesAddonProfilelocateInConfig(azureKeyVaultSecretsProviderProfile.Config, "rotationPollInterval"); v != "" {
			rotationPollInterval = v
		}

		azureKeyvaultSecretsProviderIdentity := flattenKubernetesClusterAddOnIdentityProfileTyped(azureKeyVaultSecretsProviderProfile.Identity)

		keyVaultSecretsProvider = []KeyVaultSecretsProviderModel{{
			SecretRotationInterval: rotationPollInterval,
			SecretIdentity:         flattenSecretIdentityTyped(azureKeyvaultSecretsProviderIdentity),
		}}
	}

	return
}

func flattenKubernetesClusterAddOnIdentityProfileTyped(profile *managedclusters.UserAssignedIdentity) []ConnectorIdentityModel {
	if profile == nil {
		return []ConnectorIdentityModel{}
	}

	clientID := ""
	if profile.ClientId != nil {
		clientID = *profile.ClientId
	}

	objectID := ""
	if profile.ObjectId != nil {
		objectID = *profile.ObjectId
	}

	userAssignedIdentityID := ""
	if profile.ResourceId != nil {
		userAssignedIdentityID = *profile.ResourceId
	}

	return []ConnectorIdentityModel{{
		ClientID:               clientID,
		ObjectID:               objectID,
		UserAssignedIdentityID: userAssignedIdentityID,
	}}
}

func flattenOMSAgentIdentityTyped(input []ConnectorIdentityModel) []OMSAgentIdentityModel {
	if len(input) == 0 {
		return []OMSAgentIdentityModel{}
	}
	return []OMSAgentIdentityModel{{
		ClientID:               input[0].ClientID,
		ObjectID:               input[0].ObjectID,
		UserAssignedIdentityID: input[0].UserAssignedIdentityID,
	}}
}

func flattenIngressApplicationGatewayIdentityTyped(input []ConnectorIdentityModel) []IngressApplicationGatewayIdentityModel {
	if len(input) == 0 {
		return []IngressApplicationGatewayIdentityModel{}
	}
	return []IngressApplicationGatewayIdentityModel{{
		ClientID:               input[0].ClientID,
		ObjectID:               input[0].ObjectID,
		UserAssignedIdentityID: input[0].UserAssignedIdentityID,
	}}
}

func flattenSecretIdentityTyped(input []ConnectorIdentityModel) []SecretIdentityModel {
	if len(input) == 0 {
		return []SecretIdentityModel{}
	}
	return []SecretIdentityModel{{
		ClientID:               input[0].ClientID,
		ObjectID:               input[0].ObjectID,
		UserAssignedIdentityID: input[0].UserAssignedIdentityID,
	}}
}
