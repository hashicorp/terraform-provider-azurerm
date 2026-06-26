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
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-04-01/managedclusters"
	dnsValidate "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/applicationgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipprefixes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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
	DiskEncryptionSetID             string                              `tfschema:"disk_encryption_set_id"`
	DNSPrefix                       string                              `tfschema:"dns_prefix"`
	HostedSystemProfile             []HostedSystemProfile               `tfschema:"hosted_system"`
	HTTPProxyConfig                 []HTTPProxyConfigModel              `tfschema:"proxy"`
	Identity                        []identity.SystemOrUserAssignedList `tfschema:"identity"`
	ImageCleanerIntervalHours       int64                               `tfschema:"image_cleaner_interval_in_hours"`
	KeyManagementService            []KeyManagementServiceModel         `tfschema:"key_management_service"`
	KubeletIdentity                 []KubeletIdentityModel              `tfschema:"kubelet_identity"`
	KubernetesVersion               string                              `tfschema:"kubernetes_version"`
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
	AIToolchainOperatorEnabled      bool                                `tfschema:"ai_toolchain_operator_enabled"`

	// Addon fields
	ACIConnectorLinux         []ACIConnectorLinuxModel         `tfschema:"aci_connector_linux"`
	ConfidentialComputing     []ConfidentialComputingModel     `tfschema:"confidential_computing"`
	IngressApplicationGateway []IngressApplicationGatewayModel `tfschema:"ingress_application_gateway"`
	KeyVaultSecretsProvider   []KeyVaultSecretsProviderModel   `tfschema:"key_vault_secrets_provider"`
	OMSAgent                  []OMSAgentModel                  `tfschema:"oms_agent"`

	// Computed fields
	CurrentKubernetesVersion string            `tfschema:"current_kubernetes_version"`
	FQDN                     string            `tfschema:"fully_qualified_domain_name"`
	KubeAdminConfig          []KubeConfigModel `tfschema:"kube_admin_config"`
	KubeAdminConfigRaw       string            `tfschema:"kube_admin_config_raw"`
	KubeConfig               []KubeConfigModel `tfschema:"kube_config"`
	KubeConfigRaw            string            `tfschema:"kube_config_raw"`
	NodeResourceGroupID      string            `tfschema:"node_resource_group_id"`
	OIDCIssuerURL            string            `tfschema:"oidc_issuer_url"`
	PrivateFQDN              string            `tfschema:"private_fully_qualified_domain_name"`
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

type HostedSystemProfile struct {
	NodeSubnetID       string `tfschema:"node_subnet_id"`
	SystemNodeSubnetID string `tfschema:"system_node_subnet_id"`
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

		client := metadata.Client.Containers_v2026_04_01.KubernetesClustersClient
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
			identityType := rd.Get("identity.0.type").(string)
			artifactSource := rd.Get("bootstrap_profile.0.artifact_source").(string)
			if outboundType == string(managedclusters.OutboundTypeNone) && artifactSource != string(managedclusters.ArtifactSourceCache) {
				return fmt.Errorf("when `network.outbound_type` is set to `none`, `bootstrap_profile.artifact_source` must be set to `Cache`")
			}

			hostedSystem := rd.Get("hosted_system").([]interface{})
			if len(hostedSystem) > 0 && hostedSystem[0] != nil {
				if !strings.EqualFold(identityType, string(identity.TypeUserAssigned)) {
					return fmt.Errorf("`hosted_system` requires `identity.type` to be `UserAssigned`")
				}
			}

			if outboundType == string(managedclusters.OutboundTypeManagedNATGateway) && len(hostedSystem) > 0 && hostedSystem[0] != nil {
				hostedSystemConfig := hostedSystem[0].(map[string]interface{})
				nodeSubnetID := hostedSystemConfig["node_subnet_id"].(string)
				systemNodeSubnetID := hostedSystemConfig["system_node_subnet_id"].(string)

				if nodeSubnetID != "" || systemNodeSubnetID != "" {
					return fmt.Errorf("network.outbound_type cannot be managedNATGateway when using hosted_system")
				}
			}

			if len(hostedSystem) == 0 {
				if !strings.EqualFold(identityType, string(identity.TypeSystemAssigned)) {
					return fmt.Errorf("when `hosted_system` is not configured, `identity.type` must be `SystemAssigned`")
				}

				if outboundType == string(managedclusters.OutboundTypeLoadBalancer) {
					return fmt.Errorf("when `hosted_system` is not configured, `network.outbound_type` cannot be `loadBalancer`")
				}
			}

			privateCluster := rd.Get("private_cluster").([]interface{})
			if len(privateCluster) > 0 && privateCluster[0] != nil {
				privateClusterConfig := privateCluster[0].(map[string]interface{})
				privateDNSZoneID := privateClusterConfig["private_dns_zone_id"].(string)
				dnsPrefix := privateClusterConfig["dns_prefix"].(string)

				if privateDNSZoneID != "" {
					if privateDNSZoneID != "System" && privateDNSZoneID != "None" && !strings.EqualFold(identityType, string(identity.TypeUserAssigned)) {
						return fmt.Errorf("a user assigned identity must be used when using a custom private dns zone")
					}
				}

				if dnsPrefix != "" && (privateDNSZoneID == "" || privateDNSZoneID == "System" || privateDNSZoneID == "None") {
					return fmt.Errorf("`private_cluster.0.dns_prefix` should only be set for private cluster with custom private dns zone")
				}
			}

			network := rd.Get("network").([]interface{})
			if len(network) > 0 && network[0] != nil {
				networkConfig := network[0].(map[string]interface{})
				loadBalancerSku := networkConfig["load_balancer_sku"].(string)

				loadBalancerProfile := networkConfig["load_balancer"].([]interface{})
				if len(loadBalancerProfile) > 0 && !strings.EqualFold(loadBalancerSku, "standard") {
					return fmt.Errorf("only load balancer SKU 'Standard' supports load balancer profiles. Provided load balancer type: %s", loadBalancerSku)
				}

				natGatewayProfile := networkConfig["nat_gateway"].([]interface{})
				if len(natGatewayProfile) > 0 && !strings.EqualFold(loadBalancerSku, "standard") {
					return fmt.Errorf("only load balancer SKU 'Standard' supports NAT Gateway profiles. Provided load balancer type: %s", loadBalancerSku)
				}
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
			ValidateFunc: commonids.ValidateDiskEncryptionSetID,
		},

		"dns_prefix": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"dns_prefix", "private_cluster.0.dns_prefix"},
			ValidateFunc: containerValidate.KubernetesDNSPrefix,
		},

		"hosted_system": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// O+C if no subnet ids are supplied, it will return the new, managed subnet ids
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"node_subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
					"system_node_subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"identity": commonschema.SystemOrUserAssignedIdentityRequired(),

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
						Default:  "System",
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
						// NOTE: O+C - Azure may populate default admin groups from tenant if not specified
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
			client := metadata.Client.Containers_v2026_04_01.KubernetesClustersClient
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

			securityProfile := &managedclusters.ManagedClusterSecurityProfile{
				Defender:                  expandKubernetesAutomaticClusterMicrosoftDefender(model.MicrosoftDefender, false),
				ImageCleaner:              expandKubernetesAutomaticClusterImageCleaner(model.ImageCleanerIntervalHours),
				CustomCATrustCertificates: expandKubernetesAutomaticClusterCustomCATrustCertificates(model.CustomCATrustCertificatesBase64),
			}

			kmsProfile, err := expandKubernetesAutomaticClusterKeyManagementService(model.KeyManagementService, ctx, keyVaultsClient, subscriptionId)
			if err != nil {
				return err
			}
			securityProfile.AzureKeyVaultKms = kmsProfile

			addonProfiles, err := expandKubernetesAutomaticClusterAddOns(&model, metadata.Client.Containers_v2026_04_01.Environment)
			if err != nil {
				return fmt.Errorf("expanding addons: %+v", err)
			}

			(*addonProfiles)["azurepolicy"] = managedclusters.ManagedClusterAddonProfile{
				Enabled: true,
			}

			parameters := managedclusters.ManagedCluster{
				Location: location.Normalize(model.Location),
				Sku: &managedclusters.ManagedClusterSKU{
					Name: pointer.To(managedclusters.ManagedClusterSKUName("automatic")),
					Tier: pointer.To(managedclusters.ManagedClusterSKUTier("standard")),
				},
				Properties: &managedclusters.ManagedClusterProperties{
					ApiServerAccessProfile: expandKubernetesAutomaticClusterAPIAccessProfile(model),
					AadProfile:             expandKubernetesAutomaticClusterCreateAADProfile(model.AzureActiveDirectoryRBAC),
					AddonProfiles:          addonProfiles,
					AutoScalerProfile:      expandKubernetesAutomaticClusterAutoScalerProfile(model.AutoScalerProfile),
					AzureMonitorProfile:    expandKubernetesAutomaticClusterAzureMonitorProfile(model.MonitorMetrics),
					KubernetesVersion:      pointer.To(model.KubernetesVersion),
					BootstrapProfile:       expandKubernetesAutomaticClusterBootstrapProfile(model.BootstrapProfile),
					HostedSystemProfile:    expandKubernetesAutomaticClusterHostedSystemProfile(model.HostedSystemProfile),
					MetricsProfile:         expandKubernetesAutomaticClusterMetricsProfile(model.CostAnalysisEnabled),
					NetworkProfile:         expandKubernetesAutomaticClusterNetworkProfile(model.NetworkProfile),
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

			parameters.Identity = expandIdentityModel(model.Identity)

			parameters.Properties.IdentityProfile = expandKubernetesAutomaticClusterIdentityProfile(model.KubeletIdentity)

			if len(model.PrivateCluster) > 0 && model.PrivateCluster[0].DNSPrefixPrivateCluster != "" {
				parameters.Properties.FqdnSubdomain = pointer.To(model.PrivateCluster[0].DNSPrefixPrivateCluster)
			} else {
				parameters.Properties.DnsPrefix = pointer.To(model.DNSPrefix)
			}

			if model.DiskEncryptionSetID != "" {
				parameters.Properties.DiskEncryptionSetID = pointer.To(model.DiskEncryptionSetID)
			}

			if model.SupportPlan != "" {
				parameters.Properties.SupportPlan = pointer.ToEnum[managedclusters.KubernetesSupportPlan](model.SupportPlan)
			}

			err = client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, managedclusters.DefaultCreateOrUpdateOperationOptions(), metadata.SetIDAndIdentityCallback(&id))
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
			client := metadata.Client.Containers_v2026_04_01.KubernetesClustersClient

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
	client := metadata.Client.Containers_v2026_04_01.KubernetesClustersClient

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
					state.IngressApplicationGateway,
					state.KeyVaultSecretsProvider,
					state.OMSAgent = flattenKubernetesAutomaticClusterAddOns(*props.AddonProfiles)
			}

			state.AutoScalerProfile, err = flattenKubernetesAutomaticClusterAutoScalerProfile(props.AutoScalerProfile)
			if err != nil {
				return fmt.Errorf("flattening `auto_scaler`: %+v", err)
			}

			state.MonitorMetrics = flattenKubernetesAutomaticClusterAzureMonitorProfile(props.AzureMonitorProfile)

			state.ServiceMeshProfile = flattenKubernetesAutomaticClusterServiceMeshProfile(props.ServiceMeshProfile)

			kubeletIdentity, err := flattenKubernetesAutomaticClusterIdentityProfile(pointer.From(props.IdentityProfile))
			if err != nil {
				return fmt.Errorf("flattening `kubelet_identity`: %+v", err)
			}
			state.KubeletIdentity = kubeletIdentity

			state.NetworkProfile = flattenKubernetesAutomaticClusterNetworkProfile(props.NetworkProfile)

			state.HTTPProxyConfig = flattenKubernetesAutomaticClusterHttpProxyConfig(props.HTTPProxyConfig)

			state.BootstrapProfile = flattenKubernetesAutomaticClusterBootstrapProfile(props.BootstrapProfile)

			state.HostedSystemProfile = flattenKubernetesAutomaticClusterHostedSystemProfile(props.HostedSystemProfile)

			state.UpgradeOverride = flattenKubernetesAutomaticClusterUpgradeOverride(props.UpgradeSettings)

			state.StorageProfile = flattenKubernetesAutomaticClusterStorageProfile(props.StorageProfile)

			state.WebAppRouting = flattenKubernetesAutomaticClusterWebAppRouting(props.IngressProfile)

			state.MicrosoftDefender = flattenKubernetesAutomaticClusterMicrosoftDefender(props.SecurityProfile)

			if securityProfile := props.SecurityProfile; securityProfile != nil {
				if securityProfile.AzureKeyVaultKms != nil {
					state.KeyManagementService = flattenKubernetesAutomaticClusterKeyManagementService(securityProfile.AzureKeyVaultKms)
				}
				if securityProfile.ImageCleaner != nil {
					state.ImageCleanerIntervalHours = pointer.From(securityProfile.ImageCleaner.IntervalHours)
				}
			}

			state.CostAnalysisEnabled = flattenKubernetesAutomaticClusterMetricsProfile(props.MetricsProfile)

			state.AzureActiveDirectoryRBAC = flattenKubernetesAutomaticClusterAzureActiveDirectoryRBAC(props.AadProfile)

			if props.AiToolchainOperatorProfile != nil {
				state.AIToolchainOperatorEnabled = pointer.From(props.AiToolchainOperatorProfile.Enabled)
			}

			state.SupportPlan = string(pointer.From(props.SupportPlan))
		}

		state.Identity = flattenIdentityModel(model.Identity)

		kubeConfigRaw, kubeConfig := flattenKubernetesClusterCredentialsTyped(credentials.Model, "clusterUser")
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
			client := metadata.Client.Containers_v2026_04_01
			clusterClient := client.KubernetesClustersClient
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
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
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

			if metadata.ResourceData.HasChanges("aci_connector_linux",
				"confidential_computing",
				"oms_agent",
				"ingress_application_gateway",
				"key_vault_secrets_provider") {
				addonProfiles, err := expandKubernetesAutomaticClusterAddOns(&model, metadata.Client.Containers_v2026_04_01.Environment)
				if err != nil {
					return fmt.Errorf("expanding addons: %w", err)
				}
				props.AddonProfiles = addonProfiles
				updateCluster = true
			}

			if metadata.ResourceData.HasChanges("api_server_access", "private_cluster", "run_command_enabled") {
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
				props.NetworkProfile = expandKubernetesAutomaticClusterNetworkProfile(model.NetworkProfile)
				updateCluster = true
			}

			if metadata.ResourceData.HasChange("proxy") {
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

			if metadata.ResourceData.HasChanges("microsoft_defender",
				"image_cleaner_interval_in_hours",
				"key_management_service",
				"custom_ca_trust_certificates_base64") {
				if props.SecurityProfile == nil {
					props.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
				}

				if metadata.ResourceData.HasChange("microsoft_defender") {
					props.SecurityProfile.Defender = expandKubernetesAutomaticClusterMicrosoftDefender(model.MicrosoftDefender, metadata.ResourceData.HasChange("microsoft_defender"))
				}

				if metadata.ResourceData.HasChange("image_cleaner_interval_in_hours") {
					props.SecurityProfile.ImageCleaner = expandKubernetesAutomaticClusterImageCleaner(model.ImageCleanerIntervalHours)
				}

				if metadata.ResourceData.HasChange("key_management_service") {
					props.SecurityProfile.AzureKeyVaultKms, err = expandKubernetesAutomaticClusterKeyManagementService(model.KeyManagementService, ctx, keyVaultsClient, subscriptionId)
					if err != nil {
						return err
					}
				}

				if metadata.ResourceData.HasChange("custom_ca_trust_certificates_base64") {
					props.SecurityProfile.CustomCATrustCertificates = expandKubernetesAutomaticClusterCustomCATrustCertificates(model.CustomCATrustCertificatesBase64)
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
				props.SupportPlan = pointer.ToEnum[managedclusters.KubernetesSupportPlan](model.SupportPlan)
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
			client := metadata.Client.Containers_v2026_04_01.KubernetesClustersClient

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
	apiServerAccessProfile := make([]APIServerAccessProfileModel, 0, 1)

	if profile == nil {
		return apiServerAccessProfile, []PrivateClusterModel{}, true
	}

	// Extract private cluster settings
	enablePrivateCluster := pointer.From(profile.EnablePrivateCluster)

	enablePrivateClusterPublicFQDN := pointer.From(profile.EnablePrivateClusterPublicFQDN)

	runCommandEnabled := !pointer.From(profile.DisableRunCommand)

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

func expandKubernetesAutomaticClusterHostedSystemProfile(input []HostedSystemProfile) *managedclusters.ManagedClusterHostedSystemProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]

	profile := &managedclusters.ManagedClusterHostedSystemProfile{
		Enabled: pointer.To(true),
	}

	if config.NodeSubnetID != "" {
		profile.NodeSubnetID = pointer.To(config.NodeSubnetID)
	}

	if config.SystemNodeSubnetID != "" {
		profile.SystemNodeSubnetID = pointer.To(config.SystemNodeSubnetID)
	}

	return profile
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

		userAssignedIdentityId := ""
		if resourceid := kubeletidentity.ResourceId; resourceid != nil {
			parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*resourceid)
			if err != nil {
				return nil, err
			}

			userAssignedIdentityId = parsedId.ID()
		}

		kubeletIdentity = append(kubeletIdentity, KubeletIdentityModel{
			ClientID:               pointer.From(kubeletidentity.ClientId),
			ObjectID:               pointer.From(kubeletidentity.ObjectId),
			UserAssignedIdentityID: userAssignedIdentityId,
		})
	}

	return kubeletIdentity, nil
}

func expandKubernetesAutomaticClusterImageCleaner(intervalHours int64) *managedclusters.ManagedClusterSecurityProfileImageCleaner {
	return &managedclusters.ManagedClusterSecurityProfileImageCleaner{
		Enabled:       pointer.To(true),
		IntervalHours: pointer.To(intervalHours),
	}
}

func expandKubernetesAutomaticClusterCustomCATrustCertificates(input []string) *[]string {
	if len(input) == 0 {
		return nil
	}

	return pointer.To(input)
}

func idsToAutomaticResourceReferences(ids []string) *[]managedclusters.ResourceReference {
	if len(ids) == 0 {
		return nil
	}

	results := make([]managedclusters.ResourceReference, 0, len(ids))
	for _, id := range ids {
		results = append(results, managedclusters.ResourceReference{Id: pointer.To(id)})
	}

	return &results
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

func expandKubernetesAutomaticClusterNetworkProfile(input []NetworkProfileModel) *managedclusters.ContainerServiceNetworkProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0]

	loadBalancerSku := config.LoadBalancerSKU

	networkProfile := managedclusters.ContainerServiceNetworkProfile{
		LoadBalancerSku: pointer.To(managedclusters.LoadBalancerSku(loadBalancerSku)),
		OutboundType:    pointer.ToEnum[managedclusters.OutboundType](config.OutboundType),
		IPFamilies:      &[]managedclusters.IPFamily{"IPv4"},
	}

	networkProfile.LoadBalancerProfile = expandAutomaticLoadBalancerProfile(config.LoadBalancerProfile)
	networkProfile.NatGatewayProfile = expandAutomaticNatGatewayProfile(config.NATGatewayProfile)

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

	return &networkProfile
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

	profile.AllocatedOutboundPorts = pointer.To(config.OutboundPortsAllocated)

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

	istioIngressGatewaysList := make([]managedclusters.IstioIngressGateway, 0, 2)

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

func flattenKubernetesAutomaticClusterHostedSystemProfile(profile *managedclusters.ManagedClusterHostedSystemProfile) []HostedSystemProfile {
	if profile == nil {
		return []HostedSystemProfile{}
	}

	return []HostedSystemProfile{{
		NodeSubnetID:       pointer.From(profile.NodeSubnetID),
		SystemNodeSubnetID: pointer.From(profile.SystemNodeSubnetID),
	}}
}

func expandKubernetesAutomaticClusterCreateAADProfile(input []AzureActiveDirectoryRBACModel) *managedclusters.ManagedClusterAADProfile {
	if len(input) == 0 {
		return nil
	}

	return &managedclusters.ManagedClusterAADProfile{
		Managed:             pointer.To(true),
		AdminGroupObjectIDs: &input[0].AdminGroupObjectIDs,
	}
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

func expandKubernetesAutomaticClusterAddOns(input *KubernetesAutomaticClusterModel, env environments.Environment) (*map[string]managedclusters.ManagedClusterAddonProfile, error) {
	addonProfiles := map[string]managedclusters.ManagedClusterAddonProfile{}

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

	return filterUnsupportedKubernetesAddOnsTyped(addonProfiles, env)
}

func flattenKubernetesAutomaticClusterAddOns(profile map[string]managedclusters.ManagedClusterAddonProfile) (
	aciConnectorLinux []ACIConnectorLinuxModel,
	confidentialComputing []ConfidentialComputingModel,
	ingressApplicationGateway []IngressApplicationGatewayModel,
	keyVaultSecretsProvider []KeyVaultSecretsProviderModel,
	omsAgent []OMSAgentModel,
) {
	aciConnector := kubernetesAddonProfileLocateTyped(profile, aciConnectorKey)
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

	confidentialComputingProfile := kubernetesAddonProfileLocateTyped(profile, confidentialComputingKey)
	if confidentialComputingProfile.Enabled {
		quoteHelperEnabled := false
		if v := kubernetesAddonProfileLocateTypedInConfig(confidentialComputingProfile.Config, "ACCSGXQuoteHelperEnabled"); v != "" && v != "false" {
			quoteHelperEnabled = true
		}
		confidentialComputing = []ConfidentialComputingModel{{
			SGXQuoteHelperEnabled: quoteHelperEnabled,
		}}
	}

	omsAgentProfile := kubernetesAddonProfileLocateTyped(profile, omsAgentKey)
	if omsAgentProfile.Enabled {
		workspaceID := ""
		useAADAuth := false

		if v := kubernetesAddonProfileLocateTypedInConfig(omsAgentProfile.Config, "logAnalyticsWorkspaceResourceID"); v != "" {
			if lawid, err := workspaces.ParseWorkspaceIDInsensitively(v); err == nil {
				workspaceID = lawid.ID()
			}
		}

		if v := kubernetesAddonProfileLocateTypedInConfig(omsAgentProfile.Config, "useAADAuth"); v != "false" && v != "" {
			useAADAuth = true
		}

		omsAgentIdentity := flattenKubernetesClusterAddOnIdentityProfileTyped(omsAgentProfile.Identity)

		omsAgent = []OMSAgentModel{{
			LogAnalyticsWorkspaceID:     workspaceID,
			MSIAuthForMonitoringEnabled: pointer.To(useAADAuth),
			OMSAgentIdentity:            flattenOMSAgentIdentityTyped(omsAgentIdentity),
		}}
	}

	ingressApplicationGatewayProfile := kubernetesAddonProfileLocateTyped(profile, ingressApplicationGatewayKey)
	if ingressApplicationGatewayProfile.Enabled {
		gatewayId := ""
		if v := kubernetesAddonProfileLocateTypedInConfig(ingressApplicationGatewayProfile.Config, "applicationGatewayId"); v != "" {
			gatewayId = v
		}

		gatewayName := ""
		if v := kubernetesAddonProfileLocateTypedInConfig(ingressApplicationGatewayProfile.Config, "applicationGatewayName"); v != "" {
			gatewayName = v
		}

		effectiveGatewayId := ""
		if v := kubernetesAddonProfileLocateTypedInConfig(ingressApplicationGatewayProfile.Config, "effectiveApplicationGatewayId"); v != "" {
			effectiveGatewayId = v
		}

		subnetCIDR := ""
		if v := kubernetesAddonProfileLocateTypedInConfig(ingressApplicationGatewayProfile.Config, "subnetCIDR"); v != "" {
			subnetCIDR = v
		}

		subnetId := ""
		if v := kubernetesAddonProfileLocateTypedInConfig(ingressApplicationGatewayProfile.Config, "subnetId"); v != "" {
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

	azureKeyVaultSecretsProviderProfile := kubernetesAddonProfileLocateTyped(profile, azureKeyvaultSecretsProviderKey)
	if azureKeyVaultSecretsProviderProfile.Enabled {
		rotationPollInterval := ""

		if v := kubernetesAddonProfileLocateTypedInConfig(azureKeyVaultSecretsProviderProfile.Config, "rotationPollInterval"); v != "" {
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

// when the Kubernetes Cluster is updated in the Portal - Azure updates the casing on the keys
// meaning what's submitted could be different to what's returned..
func kubernetesAddonProfileLocateTyped(profile map[string]managedclusters.ManagedClusterAddonProfile, key string) managedclusters.ManagedClusterAddonProfile {
	for k, v := range profile {
		if strings.EqualFold(k, key) {
			return v
		}
	}

	return managedclusters.ManagedClusterAddonProfile{}
}

// when the Kubernetes Cluster is updated in the Portal - Azure updates the casing on the keys
// meaning what's submitted could be different to what's returned..
// Related issue: https://github.com/Azure/azure-rest-api-specs/issues/10716
func kubernetesAddonProfileLocateTypedInConfig(config *map[string]string, key string) string {
	if config == nil {
		return ""
	}

	for k, v := range *config {
		if strings.EqualFold(k, key) {
			return v
		}
	}

	return ""
}

func filterUnsupportedKubernetesAddOnsTyped(input map[string]managedclusters.ManagedClusterAddonProfile, env environments.Environment) (*map[string]managedclusters.ManagedClusterAddonProfile, error) {
	filter := func(input map[string]managedclusters.ManagedClusterAddonProfile, key string) (map[string]managedclusters.ManagedClusterAddonProfile, error) {
		output := input
		if v, ok := output[key]; ok {
			if v.Enabled {
				return nil, fmt.Errorf("the addon %q is not supported for a Kubernetes Cluster located in %q", key, env.Name)
			}

			// otherwise it's disabled by default, so just remove it
			delete(output, key)
		}

		return output, nil
	}

	output := input
	if unsupportedAddons, ok := unsupportedAddonsForEnvironment[env.Name]; ok {
		for _, key := range unsupportedAddons {
			out, err := filter(output, key)
			if err != nil {
				return nil, err
			}

			output = out
		}
	}
	return &output, nil
}
