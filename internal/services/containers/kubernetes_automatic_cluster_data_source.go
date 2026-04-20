// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-10-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/kubernetes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = KubernetesAutomaticClusterDataSource{}

type KubernetesAutomaticClusterDataSource struct{}

type KubernetesAutomaticClusterDataSourceModel struct {
	Name                                       string                                     `tfschema:"name"`
	ResourceGroupName                          string                                     `tfschema:"resource_group_name"`
	Location                                   string                                     `tfschema:"location"`
	ACIConnectorLinux                          []ACIConnectorLinuxDataSourceModel         `tfschema:"aci_connector_linux"`
	AgentPoolProfile                           []AgentPoolProfileDataSourceModel          `tfschema:"agent_pool_profile"`
	AzureActiveDirectoryRoleBasedAccessControl []AzureActiveDirectoryRBACDataSourceModel  `tfschema:"azure_active_directory_role_based_access_control"`
	AzurePolicyEnabled                         bool                                       `tfschema:"azure_policy_enabled"`
	CurrentKubernetesVersion                   string                                     `tfschema:"current_kubernetes_version"`
	DNSPrefix                                  string                                     `tfschema:"dns_prefix"`
	FQDN                                       string                                     `tfschema:"fqdn"`
	HTTPApplicationRoutingEnabled              bool                                       `tfschema:"http_application_routing_enabled"`
	HTTPApplicationRoutingZoneName             string                                     `tfschema:"http_application_routing_zone_name"`
	IngressApplicationGateway                  []IngressApplicationGatewayDataSourceModel `tfschema:"ingress_application_gateway"`
	KeyVaultSecretsProvider                    []KeyVaultSecretsProviderDataSourceModel   `tfschema:"key_vault_secrets_provider"`
	APIServerAuthorizedIPRanges                []string                                   `tfschema:"api_server_authorized_ip_ranges"`
	DiskEncryptionSetID                        string                                     `tfschema:"disk_encryption_set_id"`
	MicrosoftDefender                          []MicrosoftDefenderDataSourceModel         `tfschema:"microsoft_defender"`
	OMSAgent                                   []OMSAgentDataSourceModel                  `tfschema:"oms_agent"`
	OpenServiceMeshEnabled                     bool                                       `tfschema:"open_service_mesh_enabled"`
	PrivateClusterEnabled                      bool                                       `tfschema:"private_cluster_enabled"`
	PrivateFQDN                                string                                     `tfschema:"private_fqdn"`
	Identity                                   []IdentityDataSourceModel                  `tfschema:"identity"`
	KeyManagementService                       []KeyManagementServiceDataSourceModel      `tfschema:"key_management_service"`
	KubernetesVersion                          string                                     `tfschema:"kubernetes_version"`
	KubeAdminConfig                            []KubeConfigModel                          `tfschema:"kube_admin_config"`
	KubeAdminConfigRaw                         string                                     `tfschema:"kube_admin_config_raw"`
	KubeConfig                                 []KubeConfigModel                          `tfschema:"kube_config"`
	KubeConfigRaw                              string                                     `tfschema:"kube_config_raw"`
	KubeletIdentity                            []KubeletIdentityDataSourceModel           `tfschema:"kubelet_identity"`
	LinuxProfile                               []LinuxProfileDataSourceModel              `tfschema:"linux_profile"`
	WindowsProfile                             []WindowsProfileDataSourceModel            `tfschema:"windows_profile"`
	NetworkProfile                             []NetworkProfileDataSourceModel            `tfschema:"network_profile"`
	NodeResourceGroup                          string                                     `tfschema:"node_resource_group"`
	NodeResourceGroupID                        string                                     `tfschema:"node_resource_group_id"`
	OIDCIssuerEnabled                          bool                                       `tfschema:"oidc_issuer_enabled"`
	OIDCIssuerURL                              string                                     `tfschema:"oidc_issuer_url"`
	RoleBasedAccessControlEnabled              bool                                       `tfschema:"role_based_access_control_enabled"`
	ServicePrincipal                           []ServicePrincipalDataSourceModel          `tfschema:"service_principal"`
	StorageProfile                             []StorageProfileDataSourceModel            `tfschema:"storage_profile"`
	ServiceMeshProfile                         []ServiceMeshProfileDataSourceModel        `tfschema:"service_mesh_profile"`
	Tags                                       map[string]interface{}                     `tfschema:"tags"`
}

type ACIConnectorLinuxDataSourceModel struct {
	SubnetName string `tfschema:"subnet_name"`
}

type AgentPoolProfileDataSourceModel struct {
	Name                 string                           `tfschema:"name"`
	Type                 string                           `tfschema:"type"`
	Count                int64                            `tfschema:"count"`
	MaxCount             int64                            `tfschema:"max_count"`
	MinCount             int64                            `tfschema:"min_count"`
	AutoScalingEnabled   bool                             `tfschema:"auto_scaling_enabled"`
	VMSize               string                           `tfschema:"vm_size"`
	Tags                 map[string]interface{}           `tfschema:"tags"`
	OSDiskSizeGB         int64                            `tfschema:"os_disk_size_gb"`
	VnetSubnetID         string                           `tfschema:"vnet_subnet_id"`
	OSType               string                           `tfschema:"os_type"`
	OrchestratorVersion  string                           `tfschema:"orchestrator_version"`
	MaxPods              int64                            `tfschema:"max_pods"`
	NodeLabels           map[string]string                `tfschema:"node_labels"`
	NodeTaints           []string                         `tfschema:"node_taints"`
	NodePublicIPEnabled  bool                             `tfschema:"node_public_ip_enabled"`
	NodePublicIPPrefixID string                           `tfschema:"node_public_ip_prefix_id"`
	UpgradeSettings      []UpgradeSettingsDataSourceModel `tfschema:"upgrade_settings"`
	Zones                []string                         `tfschema:"zones"`
}

type UpgradeSettingsDataSourceModel struct {
	MaxSurge                  string `tfschema:"max_surge"`
	MaxUnavailable            string `tfschema:"max_unavailable"`
	DrainTimeoutInMinutes     int64  `tfschema:"drain_timeout_in_minutes"`
	NodeSoakDurationInMinutes int64  `tfschema:"node_soak_duration_in_minutes"`
	UndrainableNodeBehavior   string `tfschema:"undrainable_node_behavior"`
}

type AzureActiveDirectoryRBACDataSourceModel struct {
	TenantID            string   `tfschema:"tenant_id"`
	AzureRBACEnabled    bool     `tfschema:"azure_rbac_enabled"`
	AdminGroupObjectIDs []string `tfschema:"admin_group_object_ids"`
}

type IngressApplicationGatewayDataSourceModel struct {
	GatewayID                         string                                             `tfschema:"gateway_id"`
	GatewayName                       string                                             `tfschema:"gateway_name"`
	SubnetCIDR                        string                                             `tfschema:"subnet_cidr"`
	SubnetID                          string                                             `tfschema:"subnet_id"`
	EffectiveGatewayID                string                                             `tfschema:"effective_gateway_id"`
	IngressApplicationGatewayIdentity []IngressApplicationGatewayIdentityDataSourceModel `tfschema:"ingress_application_gateway_identity"`
}

type IngressApplicationGatewayIdentityDataSourceModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type KeyVaultSecretsProviderDataSourceModel struct {
	SecretRotationEnabled  bool                            `tfschema:"secret_rotation_enabled"`
	SecretRotationInterval string                          `tfschema:"secret_rotation_interval"`
	SecretIdentity         []SecretIdentityDataSourceModel `tfschema:"secret_identity"`
}

type SecretIdentityDataSourceModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type MicrosoftDefenderDataSourceModel struct {
	LogAnalyticsWorkspaceID string `tfschema:"log_analytics_workspace_id"`
}

type OMSAgentDataSourceModel struct {
	LogAnalyticsWorkspaceID     string                            `tfschema:"log_analytics_workspace_id"`
	MSIAuthForMonitoringEnabled bool                              `tfschema:"msi_auth_for_monitoring_enabled"`
	OMSAgentIdentity            []OMSAgentIdentityDataSourceModel `tfschema:"oms_agent_identity"`
}

type OMSAgentIdentityDataSourceModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type IdentityDataSourceModel struct {
	Type        string   `tfschema:"type"`
	IdentityIDs []string `tfschema:"identity_ids"`
	PrincipalID string   `tfschema:"principal_id"`
	TenantID    string   `tfschema:"tenant_id"`
}

type KeyManagementServiceDataSourceModel struct {
	KeyVaultKeyID         string `tfschema:"key_vault_key_id"`
	KeyVaultNetworkAccess string `tfschema:"key_vault_network_access"`
}

type KubeletIdentityDataSourceModel struct {
	ClientID               string `tfschema:"client_id"`
	ObjectID               string `tfschema:"object_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

type LinuxProfileDataSourceModel struct {
	AdminUsername string                  `tfschema:"admin_username"`
	SSHKey        []SSHKeyDataSourceModel `tfschema:"ssh_key"`
}

type SSHKeyDataSourceModel struct {
	KeyData string `tfschema:"key_data"`
}

type WindowsProfileDataSourceModel struct {
	AdminUsername string `tfschema:"admin_username"`
}

type NetworkProfileDataSourceModel struct {
	NetworkPlugin    string `tfschema:"network_plugin"`
	NetworkPolicy    string `tfschema:"network_policy"`
	ServiceCIDR      string `tfschema:"service_cidr"`
	DNSServiceIP     string `tfschema:"dns_service_ip"`
	DockerBridgeCIDR string `tfschema:"docker_bridge_cidr"`
	PodCIDR          string `tfschema:"pod_cidr"`
	LoadBalancerSKU  string `tfschema:"load_balancer_sku"`
}

type ServicePrincipalDataSourceModel struct {
	ClientID string `tfschema:"client_id"`
}

type StorageProfileDataSourceModel struct {
	BlobDriverEnabled         bool `tfschema:"blob_driver_enabled"`
	DiskDriverEnabled         bool `tfschema:"disk_driver_enabled"`
	FileDriverEnabled         bool `tfschema:"file_driver_enabled"`
	SnapshotControllerEnabled bool `tfschema:"snapshot_controller_enabled"`
}

type ServiceMeshProfileDataSourceModel struct {
	Mode                          string                                `tfschema:"mode"`
	InternalIngressGatewayEnabled bool                                  `tfschema:"internal_ingress_gateway_enabled"`
	ExternalIngressGatewayEnabled bool                                  `tfschema:"external_ingress_gateway_enabled"`
	CertificateAuthority          []CertificateAuthorityDataSourceModel `tfschema:"certificate_authority"`
	Revisions                     []string                              `tfschema:"revisions"`
}

type CertificateAuthorityDataSourceModel struct {
	KeyVaultID          string `tfschema:"key_vault_id"`
	RootCertObjectName  string `tfschema:"root_cert_object_name"`
	CertChainObjectName string `tfschema:"cert_chain_object_name"`
	CertObjectName      string `tfschema:"cert_object_name"`
	KeyObjectName       string `tfschema:"key_object_name"`
}

func (KubernetesAutomaticClusterDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (KubernetesAutomaticClusterDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),
		"aci_connector_linux": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"agent_pool_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"max_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"min_count": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"auto_scaling_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"vm_size": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"tags": commonschema.TagsDataSource(),

					"os_disk_size_gb": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"vnet_subnet_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"os_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"orchestrator_version": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"max_pods": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},

					"node_labels": {
						Type:     pluginsdk.TypeMap,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"node_taints": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},

					"node_public_ip_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"node_public_ip_prefix_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"upgrade_settings": upgradeSettingsForDataSourceSchema(),

					"zones": commonschema.ZonesMultipleComputed(),
				},
			},
		},

		"azure_active_directory_role_based_access_control": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"azure_rbac_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"admin_group_object_ids": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"azure_policy_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"current_kubernetes_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"dns_prefix": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"http_application_routing_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"http_application_routing_zone_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ingress_application_gateway": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"gateway_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"gateway_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"subnet_cidr": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"effective_gateway_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"ingress_application_gateway_identity": {
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
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"secret_rotation_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
					"secret_rotation_interval": {
						Type:     pluginsdk.TypeString,
						Computed: true,
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

		"api_server_authorized_ip_ranges": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"disk_encryption_set_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"microsoft_defender": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"log_analytics_workspace_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"oms_agent": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"log_analytics_workspace_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"msi_auth_for_monitoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
					"oms_agent_identity": {
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
			Computed: true,
		},

		"private_cluster_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"private_fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

		"key_management_service": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"key_vault_network_access": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"kubernetes_version": {
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

		"kubelet_identity": {
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

		"linux_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"ssh_key": {
						Type:     pluginsdk.TypeList,
						Computed: true,

						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_data": {
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
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_plugin": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"network_policy": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"service_cidr": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"dns_service_ip": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"docker_bridge_cidr": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"pod_cidr": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"load_balancer_sku": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"node_resource_group": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"node_resource_group_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"oidc_issuer_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"oidc_issuer_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"role_based_access_control_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"service_principal": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"storage_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"blob_driver_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
					"disk_driver_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
					"file_driver_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
					"snapshot_controller_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
				},
			},
		},

		"service_mesh_profile": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"internal_ingress_gateway_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
					"external_ingress_gateway_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
					"certificate_authority": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_vault_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"root_cert_object_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"cert_chain_object_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"cert_object_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"key_object_name": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
					"revisions": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (KubernetesAutomaticClusterDataSource) ModelObject() interface{} {
	return &KubernetesAutomaticClusterDataSourceModel{}
}

func (KubernetesAutomaticClusterDataSource) ResourceType() string {
	return "azurerm_kubernetes_automatic_cluster"
}

func (KubernetesAutomaticClusterDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.KubernetesClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state KubernetesAutomaticClusterDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := commonids.NewKubernetesClusterID(subscriptionId, state.ResourceGroupName, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			userCredentialsResp, err := client.ListClusterUserCredentials(ctx, id, managedclusters.ListClusterUserCredentialsOperationOptions{})
			if err != nil && !response.WasStatusCode(userCredentialsResp.HttpResponse, http.StatusForbidden) {
				return fmt.Errorf("retrieving User Credentials for %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					state.DNSPrefix = pointer.From(props.DnsPrefix)
					state.FQDN = pointer.From(props.Fqdn)
					state.PrivateFQDN = pointer.From(props.PrivateFQDN)
					state.KubernetesVersion = pointer.From(props.KubernetesVersion)
					state.CurrentKubernetesVersion = pointer.From(props.CurrentKubernetesVersion)
					state.NodeResourceGroup = pointer.From(props.NodeResourceGroup)
					state.NodeResourceGroupID = pointer.From(props.NodeResourceGroup)
					state.RoleBasedAccessControlEnabled = pointer.From(props.EnableRBAC)
					state.PrivateClusterEnabled = pointer.From(props.ApiServerAccessProfile.EnablePrivateCluster)

					// Flatten complex nested structures
					state.AgentPoolProfile = flattenKubernetesAutomaticClusterDataSourceAgentPoolProfiles(props.AgentPoolProfiles)
					state.NetworkProfile = flattenKubernetesAutomaticClusterDataSourceNetworkProfile(props.NetworkProfile)
					state.LinuxProfile = flattenKubernetesAutomaticClusterDataSourceLinuxProfile(props.LinuxProfile)
					state.WindowsProfile = flattenKubernetesAutomaticClusterDataSourceWindowsProfile(props.WindowsProfile)
					state.ServicePrincipal = flattenKubernetesAutomaticClusterDataSourceServicePrincipalProfile(props.ServicePrincipalProfile)
					state.AzureActiveDirectoryRoleBasedAccessControl = flattenKubernetesAutomaticClusterDataSourceAzureActiveDirectoryRoleBasedAccessControl(props)

					// Flatten add-ons
					if props.AddonProfiles != nil {
						state.ACIConnectorLinux = flattenKubernetesAutomaticClusterDataSourceACIConnectorLinux(*props.AddonProfiles)
						state.AzurePolicyEnabled = flattenKubernetesAutomaticClusterDataSourceAzurePolicy(*props.AddonProfiles)
						state.HTTPApplicationRoutingEnabled, state.HTTPApplicationRoutingZoneName = flattenKubernetesAutomaticClusterDataSourceHTTPApplicationRouting(*props.AddonProfiles)
						state.IngressApplicationGateway = flattenKubernetesAutomaticClusterDataSourceIngressApplicationGateway(*props.AddonProfiles)
						state.KeyVaultSecretsProvider = flattenKubernetesAutomaticClusterDataSourceKeyVaultSecretsProvider(*props.AddonProfiles)
						state.OMSAgent = flattenKubernetesAutomaticClusterDataSourceOMSAgent(*props.AddonProfiles)
						state.OpenServiceMeshEnabled = flattenKubernetesAutomaticClusterDataSourceOpenServiceMesh(*props.AddonProfiles)
					}

					// Flatten security and other profiles
					if props.SecurityProfile != nil {
						state.MicrosoftDefender = flattenKubernetesAutomaticClusterDataSourceMicrosoftDefender(props.SecurityProfile)
						state.KeyManagementService = flattenKubernetesAutomaticClusterDataSourceKeyVaultKms(props.SecurityProfile.AzureKeyVaultKms)
					}

					state.StorageProfile = flattenKubernetesAutomaticClusterDataSourceStorageProfile(props.StorageProfile)
					state.ServiceMeshProfile = flattenKubernetesAutomaticClusterDataSourceServiceMeshProfile(props.ServiceMeshProfile)

					// OIDC Issuer
					if props.OidcIssuerProfile != nil {
						state.OIDCIssuerEnabled = pointer.From(props.OidcIssuerProfile.Enabled)
						state.OIDCIssuerURL = pointer.From(props.OidcIssuerProfile.IssuerURL)
					}

					// API Server settings
					if props.ApiServerAccessProfile != nil {
						state.APIServerAuthorizedIPRanges = pointer.From(props.ApiServerAccessProfile.AuthorizedIPRanges)
					}

					// Disk encryption
					if props.DiskEncryptionSetID != nil {
						state.DiskEncryptionSetID = pointer.From(props.DiskEncryptionSetID)
					}
				}

				// Flatten identity
				state.Identity = flattenKubernetesAutomaticClusterDataSourceIdentity(model.Identity)

				// Flatten kubelet identity
				if model.Properties != nil && model.Properties.IdentityProfile != nil {
					kubeletIdentity, err := flattenKubernetesAutomaticClusterDataSourceIdentityProfile(model.Properties.IdentityProfile)
					if err != nil {
						return fmt.Errorf("flattening `kubelet_identity`: %+v", err)
					}
					state.KubeletIdentity = kubeletIdentity
				}

				// Flatten kube configs
				if userCredentialsResp.Model != nil {
					kubeConfigRaw, kubeConfig := flattenKubernetesAutomaticClusterCredentials(userCredentialsResp.Model, "clusterUser")
					state.KubeConfigRaw = pointer.From(kubeConfigRaw)
					state.KubeConfig = kubeConfig
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenKubernetesAutomaticClusterDataSourceStorageProfile(input *managedclusters.ManagedClusterStorageProfile) []StorageProfileDataSourceModel {
	if input == nil {
		return []StorageProfileDataSourceModel{}
	}

	blobEnabled := false
	if input.BlobCSIDriver != nil && input.BlobCSIDriver.Enabled != nil {
		blobEnabled = pointer.From(input.BlobCSIDriver.Enabled)
	}

	diskEnabled := true
	if input.DiskCSIDriver != nil && input.DiskCSIDriver.Enabled != nil {
		diskEnabled = pointer.From(input.DiskCSIDriver.Enabled)
	}

	fileEnabled := true
	if input.FileCSIDriver != nil && input.FileCSIDriver.Enabled != nil {
		fileEnabled = pointer.From(input.FileCSIDriver.Enabled)
	}

	snapshotController := true
	if input.SnapshotController != nil && input.SnapshotController.Enabled != nil {
		snapshotController = pointer.From(input.SnapshotController.Enabled)
	}

	return []StorageProfileDataSourceModel{{
		BlobDriverEnabled:         blobEnabled,
		DiskDriverEnabled:         diskEnabled,
		FileDriverEnabled:         fileEnabled,
		SnapshotControllerEnabled: snapshotController,
	}}
}

func flattenKubernetesAutomaticClusterCredentials(model *managedclusters.CredentialResults, configName string) (*string, []KubeConfigModel) {
	if model == nil || model.Kubeconfigs == nil || len(*model.Kubeconfigs) < 1 {
		return nil, []KubeConfigModel{}
	}

	for _, c := range *model.Kubeconfigs {
		if c.Name == nil || *c.Name != configName {
			continue
		}
		if kubeConfigRaw := c.Value; kubeConfigRaw != nil {
			rawConfig := *kubeConfigRaw
			if base64IsEncoded(*kubeConfigRaw) {
				rawConfig = base64Decode(*kubeConfigRaw)
			}

			var flattenedKubeConfig []KubeConfigModel

			if strings.Contains(rawConfig, "apiserver-id:") || strings.Contains(rawConfig, "exec") {
				kubeConfigAAD, err := kubernetes.ParseKubeConfigAAD(rawConfig)
				if err != nil {
					return pointer.To(rawConfig), []KubeConfigModel{}
				}

				flattenedKubeConfig = flattenKubernetesAutomaticClusterDataSourceKubeConfigAAD(*kubeConfigAAD)
			} else {
				kubeConfig, err := kubernetes.ParseKubeConfig(rawConfig)
				if err != nil {
					return pointer.To(rawConfig), []KubeConfigModel{}
				}

				flattenedKubeConfig = flattenKubernetesAutomaticClusterDataSourceKubeConfig(*kubeConfig)
			}

			return pointer.To(rawConfig), flattenedKubeConfig
		}
	}

	return nil, []KubeConfigModel{}
}

func flattenKubernetesAutomaticClusterDataSourceACIConnectorLinux(profile map[string]managedclusters.ManagedClusterAddonProfile) []ACIConnectorLinuxDataSourceModel {
	aciConnector := kubernetesAddonProfileLocate(profile, aciConnectorKey)
	if !aciConnector.Enabled {
		return []ACIConnectorLinuxDataSourceModel{}
	}

	subnetName := ""
	if v := aciConnector.Config; v != nil && (*v)["SubnetName"] != "" {
		subnetName = (*v)["SubnetName"]
	}

	return []ACIConnectorLinuxDataSourceModel{{
		SubnetName: subnetName,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceAzurePolicy(profile map[string]managedclusters.ManagedClusterAddonProfile) bool {
	azurePolicy := kubernetesAddonProfileLocate(profile, azurePolicyKey)
	return azurePolicy.Enabled
}

func flattenKubernetesAutomaticClusterDataSourceHTTPApplicationRouting(profile map[string]managedclusters.ManagedClusterAddonProfile) (bool, string) {
	httpApplicationRouting := kubernetesAddonProfileLocate(profile, httpApplicationRoutingKey)
	enabled := httpApplicationRouting.Enabled
	zoneName := kubernetesAddonProfilelocateInConfig(httpApplicationRouting.Config, "HTTPApplicationRoutingZoneName")
	return enabled, zoneName
}

func flattenKubernetesAutomaticClusterDataSourceOMSAgent(profile map[string]managedclusters.ManagedClusterAddonProfile) []OMSAgentDataSourceModel {
	omsAgent := kubernetesAddonProfileLocate(profile, omsAgentKey)
	if !omsAgent.Enabled {
		return []OMSAgentDataSourceModel{}
	}

	workspaceID := ""
	if v := kubernetesAddonProfilelocateInConfig(omsAgent.Config, "logAnalyticsWorkspaceResourceID"); v != "" {
		if lawid, err := workspaces.ParseWorkspaceID(v); err == nil {
			workspaceID = lawid.ID()
		}
	}

	useAADAuth := false
	if v := kubernetesAddonProfilelocateInConfig(omsAgent.Config, "useAADAuth"); v != "false" && v != "" {
		useAADAuth = true
	}

	omsAgentIdentity := flattenKubernetesAutomaticClusterDataSourceAddOnIdentity(omsAgent.Identity)

	return []OMSAgentDataSourceModel{{
		LogAnalyticsWorkspaceID:     workspaceID,
		MSIAuthForMonitoringEnabled: useAADAuth,
		OMSAgentIdentity:            omsAgentIdentity,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceIngressApplicationGateway(profile map[string]managedclusters.ManagedClusterAddonProfile) []IngressApplicationGatewayDataSourceModel {
	ingressApplicationGateway := kubernetesAddonProfileLocate(profile, ingressApplicationGatewayKey)
	if !ingressApplicationGateway.Enabled {
		return []IngressApplicationGatewayDataSourceModel{}
	}

	gatewayId := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "applicationGatewayId")
	gatewayName := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "applicationGatewayName")
	effectiveGatewayId := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "effectiveApplicationGatewayId")
	subnetCIDR := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "subnetCIDR")
	subnetId := kubernetesAddonProfilelocateInConfig(ingressApplicationGateway.Config, "subnetId")

	ingressApplicationGatewayIdentity := flattenKubernetesAutomaticClusterDataSourceIngressAppGatewayIdentity(ingressApplicationGateway.Identity)

	return []IngressApplicationGatewayDataSourceModel{{
		GatewayID:                         gatewayId,
		GatewayName:                       gatewayName,
		EffectiveGatewayID:                effectiveGatewayId,
		SubnetCIDR:                        subnetCIDR,
		SubnetID:                          subnetId,
		IngressApplicationGatewayIdentity: ingressApplicationGatewayIdentity,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceOpenServiceMesh(profile map[string]managedclusters.ManagedClusterAddonProfile) bool {
	openServiceMesh := kubernetesAddonProfileLocate(profile, openServiceMeshKey)
	return openServiceMesh.Enabled
}

func flattenKubernetesAutomaticClusterDataSourceKeyVaultSecretsProvider(profile map[string]managedclusters.ManagedClusterAddonProfile) []KeyVaultSecretsProviderDataSourceModel {
	azureKeyVaultSecretsProvider := kubernetesAddonProfileLocate(profile, azureKeyvaultSecretsProviderKey)
	if !azureKeyVaultSecretsProvider.Enabled {
		return []KeyVaultSecretsProviderDataSourceModel{}
	}

	enableSecretRotation := false
	if v := kubernetesAddonProfilelocateInConfig(azureKeyVaultSecretsProvider.Config, "enableSecretRotation"); v != "false" && v != "" {
		enableSecretRotation = true
	}

	rotationPollInterval := kubernetesAddonProfilelocateInConfig(azureKeyVaultSecretsProvider.Config, "rotationPollInterval")

	azureKeyvaultSecretsProviderIdentity := flattenKubernetesAutomaticClusterDataSourceSecretIdentity(azureKeyVaultSecretsProvider.Identity)

	return []KeyVaultSecretsProviderDataSourceModel{{
		SecretRotationEnabled:  enableSecretRotation,
		SecretRotationInterval: rotationPollInterval,
		SecretIdentity:         azureKeyvaultSecretsProviderIdentity,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceAddOnIdentity(input *managedclusters.UserAssignedIdentity) []OMSAgentIdentityDataSourceModel {
	if input == nil {
		return []OMSAgentIdentityDataSourceModel{}
	}

	clientId := pointer.From(input.ClientId)
	objectId := pointer.From(input.ObjectId)
	userAssignedIdentityId := ""
	if resourceId := input.ResourceId; resourceId != nil {
		if parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*resourceId); err == nil {
			userAssignedIdentityId = parsedId.ID()
		}
	}

	return []OMSAgentIdentityDataSourceModel{{
		ClientID:               clientId,
		ObjectID:               objectId,
		UserAssignedIdentityID: userAssignedIdentityId,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceIngressAppGatewayIdentity(input *managedclusters.UserAssignedIdentity) []IngressApplicationGatewayIdentityDataSourceModel {
	if input == nil {
		return []IngressApplicationGatewayIdentityDataSourceModel{}
	}

	clientId := pointer.From(input.ClientId)
	objectId := pointer.From(input.ObjectId)
	userAssignedIdentityId := ""
	if resourceId := input.ResourceId; resourceId != nil {
		if parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*resourceId); err == nil {
			userAssignedIdentityId = parsedId.ID()
		}
	}

	return []IngressApplicationGatewayIdentityDataSourceModel{{
		ClientID:               clientId,
		ObjectID:               objectId,
		UserAssignedIdentityID: userAssignedIdentityId,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceSecretIdentity(input *managedclusters.UserAssignedIdentity) []SecretIdentityDataSourceModel {
	if input == nil {
		return []SecretIdentityDataSourceModel{}
	}

	clientId := pointer.From(input.ClientId)
	objectId := pointer.From(input.ObjectId)
	userAssignedIdentityId := ""
	if resourceId := input.ResourceId; resourceId != nil {
		if parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(*resourceId); err == nil {
			userAssignedIdentityId = parsedId.ID()
		}
	}

	return []SecretIdentityDataSourceModel{{
		ClientID:               clientId,
		ObjectID:               objectId,
		UserAssignedIdentityID: userAssignedIdentityId,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceAgentPoolProfiles(input *[]managedclusters.ManagedClusterAgentPoolProfile) []AgentPoolProfileDataSourceModel {
	agentPoolProfiles := make([]AgentPoolProfileDataSourceModel, 0)

	if input == nil {
		return agentPoolProfiles
	}

	for _, profile := range *input {
		count := int64(0)
		if profile.Count != nil {
			count = pointer.From(profile.Count)
		}

		enableNodePublicIP := false
		if profile.EnableNodePublicIP != nil {
			enableNodePublicIP = *profile.EnableNodePublicIP
		}

		minCount := int64(0)
		if profile.MinCount != nil {
			minCount = pointer.From(profile.MinCount)
		}

		maxCount := int64(0)
		if profile.MaxCount != nil {
			maxCount = pointer.From(profile.MaxCount)
		}

		enableAutoScaling := false
		if profile.EnableAutoScaling != nil {
			enableAutoScaling = *profile.EnableAutoScaling
		}

		name := profile.Name

		nodePublicIPPrefixID := pointer.From(profile.NodePublicIPPrefixID)

		osDiskSizeGb := int64(0)
		if profile.OsDiskSizeGB != nil {
			osDiskSizeGb = pointer.From(profile.OsDiskSizeGB)
		}

		vnetSubnetId := pointer.From(profile.VnetSubnetID)
		orchestratorVersion := pointer.From(profile.OrchestratorVersion)

		maxPods := int64(0)
		if profile.MaxPods != nil {
			maxPods = pointer.From(profile.MaxPods)
		}

		nodeLabels := make(map[string]string)
		if profile.NodeLabels != nil {
			for k, v := range *profile.NodeLabels {
				if v == "" {
					continue
				}
				nodeLabels[k] = v
			}
		}

		nodeTaints := pointer.From(profile.NodeTaints)
		vmSize := pointer.From(profile.VMSize)

		osType := ""
		if profile.OsType != nil {
			osType = string(*profile.OsType)
		}
		profileType := ""
		if profile.Type != nil {
			profileType = string(*profile.Type)
		}

		// Convert zones from []interface{} to []string
		zonesList := make([]string, 0)
		if profile.AvailabilityZones != nil {
			zonesList = append(zonesList, *profile.AvailabilityZones...)
		}

		out := AgentPoolProfileDataSourceModel{
			Name:                 name,
			Type:                 profileType,
			Count:                count,
			MaxCount:             maxCount,
			MinCount:             minCount,
			AutoScalingEnabled:   enableAutoScaling,
			VMSize:               vmSize,
			Tags:                 tags.Flatten(profile.Tags),
			OSDiskSizeGB:         osDiskSizeGb,
			VnetSubnetID:         vnetSubnetId,
			OSType:               osType,
			OrchestratorVersion:  orchestratorVersion,
			MaxPods:              maxPods,
			NodeLabels:           nodeLabels,
			NodeTaints:           nodeTaints,
			NodePublicIPEnabled:  enableNodePublicIP,
			NodePublicIPPrefixID: nodePublicIPPrefixID,
			UpgradeSettings:      flattenKubernetesAutomaticClusterDataSourceUpgradeSettings(profile.UpgradeSettings),
			Zones:                zonesList,
		}

		agentPoolProfiles = append(agentPoolProfiles, out)
	}

	return agentPoolProfiles
}

func flattenKubernetesAutomaticClusterDataSourceAzureActiveDirectoryRoleBasedAccessControl(input *managedclusters.ManagedClusterProperties) []AzureActiveDirectoryRBACDataSourceModel {
	if input == nil || input.AadProfile == nil {
		return []AzureActiveDirectoryRBACDataSourceModel{}
	}

	profile := input.AadProfile
	adminGroupObjectIds := pointer.From(profile.AdminGroupObjectIDs)

	azureRbacEnabled := pointer.From(profile.EnableAzureRBAC)
	tenantId := pointer.From(profile.TenantID)

	return []AzureActiveDirectoryRBACDataSourceModel{{
		TenantID:            tenantId,
		AzureRBACEnabled:    azureRbacEnabled,
		AdminGroupObjectIDs: adminGroupObjectIds,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceIdentityProfile(profile *map[string]managedclusters.UserAssignedIdentity) ([]KubeletIdentityDataSourceModel, error) {
	if profile == nil || *profile == nil {
		return []KubeletIdentityDataSourceModel{}, nil
	}

	kubeletidentity := (*profile)["kubeletidentity"]

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

	return []KubeletIdentityDataSourceModel{{
		ClientID:               clientId,
		ObjectID:               objectId,
		UserAssignedIdentityID: userAssignedIdentityId,
	}}, nil
}

func flattenKubernetesAutomaticClusterDataSourceLinuxProfile(input *managedclusters.ContainerServiceLinuxProfile) []LinuxProfileDataSourceModel {
	if input == nil {
		return []LinuxProfileDataSourceModel{}
	}

	sshKeys := make([]SSHKeyDataSourceModel, 0)
	if input.Ssh.PublicKeys != nil {
		for _, sshKey := range input.Ssh.PublicKeys {
			if sshKey.KeyData != "" {
				sshKeys = append(sshKeys, SSHKeyDataSourceModel{
					KeyData: sshKey.KeyData,
				})
			}
		}
	}

	return []LinuxProfileDataSourceModel{{
		AdminUsername: input.AdminUsername,
		SSHKey:        sshKeys,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceWindowsProfile(input *managedclusters.ManagedClusterWindowsProfile) []WindowsProfileDataSourceModel {
	if input == nil {
		return []WindowsProfileDataSourceModel{}
	}

	return []WindowsProfileDataSourceModel{{
		AdminUsername: input.AdminUsername,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceNetworkProfile(profile *managedclusters.ContainerServiceNetworkProfile) []NetworkProfileDataSourceModel {
	if profile == nil {
		return []NetworkProfileDataSourceModel{}
	}

	networkPlugin := ""
	if profile.NetworkPlugin != nil {
		networkPlugin = string(*profile.NetworkPlugin)
	}

	networkPolicy := ""
	if profile.NetworkPolicy != nil {
		networkPolicy = string(*profile.NetworkPolicy)
	}

	loadBalancerSku := ""
	if profile.LoadBalancerSku != nil {
		loadBalancerSku = string(*profile.LoadBalancerSku)
	}

	return []NetworkProfileDataSourceModel{{
		NetworkPlugin:    networkPlugin,
		NetworkPolicy:    networkPolicy,
		ServiceCIDR:      pointer.From(profile.ServiceCidr),
		DNSServiceIP:     pointer.From(profile.DnsServiceIP),
		DockerBridgeCIDR: "", // DockerBridgeCidr field doesn't exist in this API version
		PodCIDR:          pointer.From(profile.PodCidr),
		LoadBalancerSKU:  loadBalancerSku,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceServicePrincipalProfile(profile *managedclusters.ManagedClusterServicePrincipalProfile) []ServicePrincipalDataSourceModel {
	if profile == nil {
		return []ServicePrincipalDataSourceModel{}
	}

	clientID := profile.ClientId

	return []ServicePrincipalDataSourceModel{{
		ClientID: clientID,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceKubeConfig(config kubernetes.KubeConfig) []KubeConfigModel {
	cluster := config.Clusters[0].Cluster
	user := config.Users[0].User
	name := config.Users[0].Name

	return []KubeConfigModel{{
		Host:                 cluster.Server,
		Username:             name,
		Password:             user.Token,
		ClientCertificate:    user.ClientCertificteData,
		ClientKey:            user.ClientKeyData,
		ClusterCACertificate: cluster.ClusterAuthorityData,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceKubeConfigAAD(config kubernetes.KubeConfigAAD) []KubeConfigModel {
	cluster := config.Clusters[0].Cluster
	name := config.Users[0].Name

	return []KubeConfigModel{{
		Host:                 cluster.Server,
		Username:             name,
		Password:             "",
		ClientCertificate:    "",
		ClientKey:            "",
		ClusterCACertificate: cluster.ClusterAuthorityData,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceMicrosoftDefender(input *managedclusters.ManagedClusterSecurityProfile) []MicrosoftDefenderDataSourceModel {
	if input == nil || input.Defender == nil || input.Defender.SecurityMonitoring == nil || (input.Defender.SecurityMonitoring.Enabled != nil && !*input.Defender.SecurityMonitoring.Enabled) {
		return []MicrosoftDefenderDataSourceModel{}
	}

	logAnalyticsWorkspace := pointer.From(input.Defender.LogAnalyticsWorkspaceResourceId)

	return []MicrosoftDefenderDataSourceModel{{
		LogAnalyticsWorkspaceID: logAnalyticsWorkspace,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceUpgradeSettings(input *managedclusters.AgentPoolUpgradeSettings) []UpgradeSettingsDataSourceModel {
	if input == nil {
		return []UpgradeSettingsDataSourceModel{}
	}

	maxSurge := pointer.From(input.MaxSurge)
	maxUnavailable := pointer.From(input.MaxUnavailable)
	drainTimeout := int64(0)
	if input.DrainTimeoutInMinutes != nil {
		drainTimeout = pointer.From(input.DrainTimeoutInMinutes)
	}
	nodeSoakDuration := int64(0)
	if input.NodeSoakDurationInMinutes != nil {
		nodeSoakDuration = pointer.From(input.NodeSoakDurationInMinutes)
	}
	undrainableBehavior := ""
	if input.UndrainableNodeBehavior != nil {
		undrainableBehavior = string(*input.UndrainableNodeBehavior)
	}

	return []UpgradeSettingsDataSourceModel{{
		MaxSurge:                  maxSurge,
		MaxUnavailable:            maxUnavailable,
		DrainTimeoutInMinutes:     drainTimeout,
		NodeSoakDurationInMinutes: nodeSoakDuration,
		UndrainableNodeBehavior:   undrainableBehavior,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceServiceMeshProfile(profile *managedclusters.ServiceMeshProfile) []ServiceMeshProfileDataSourceModel {
	if profile == nil || profile.Mode != managedclusters.ServiceMeshModeIstio || profile.Istio == nil {
		return []ServiceMeshProfileDataSourceModel{}
	}

	mode := string(profile.Mode)
	revisions := pointer.From(profile.Istio.Revisions)

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

	certificateAuthority := flattenKubernetesAutomaticClusterDataSourceServiceMeshProfileCertificateAuthority(profile.Istio.CertificateAuthority)

	return []ServiceMeshProfileDataSourceModel{{
		Mode:                          mode,
		InternalIngressGatewayEnabled: internalIngressGatewayEnabled,
		ExternalIngressGatewayEnabled: externalIngressGatewayEnabled,
		CertificateAuthority:          certificateAuthority,
		Revisions:                     revisions,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceServiceMeshProfileCertificateAuthority(certificateAuthority *managedclusters.IstioCertificateAuthority) []CertificateAuthorityDataSourceModel {
	if certificateAuthority == nil || certificateAuthority.Plugin == nil {
		return []CertificateAuthorityDataSourceModel{}
	}

	plugin := certificateAuthority.Plugin

	return []CertificateAuthorityDataSourceModel{{
		KeyVaultID:          pointer.From(plugin.KeyVaultId),
		RootCertObjectName:  pointer.From(plugin.RootCertObjectName),
		CertChainObjectName: pointer.From(plugin.CertChainObjectName),
		CertObjectName:      pointer.From(plugin.CertObjectName),
		KeyObjectName:       pointer.From(plugin.KeyObjectName),
	}}
}

func flattenKubernetesAutomaticClusterDataSourceIdentity(input *identity.SystemOrUserAssignedMap) []IdentityDataSourceModel {
	if input == nil {
		return []IdentityDataSourceModel{}
	}

	identityType := ""
	if input.Type != "" {
		identityType = string(input.Type)
	}

	identityIDs := make([]string, 0)
	if input.IdentityIds != nil {
		for id := range input.IdentityIds {
			identityIDs = append(identityIDs, id)
		}
	}

	return []IdentityDataSourceModel{{
		Type:        identityType,
		IdentityIDs: identityIDs,
		PrincipalID: input.PrincipalId,
		TenantID:    input.TenantId,
	}}
}

func flattenKubernetesAutomaticClusterDataSourceKeyVaultKms(input *managedclusters.AzureKeyVaultKms) []KeyManagementServiceDataSourceModel {
	if input == nil || !pointer.From(input.Enabled) {
		return []KeyManagementServiceDataSourceModel{}
	}

	keyVaultKeyID := pointer.From(input.KeyId)
	keyVaultNetworkAccess := ""
	if input.KeyVaultNetworkAccess != nil {
		keyVaultNetworkAccess = string(*input.KeyVaultNetworkAccess)
	}

	return []KeyManagementServiceDataSourceModel{{
		KeyVaultKeyID:         keyVaultKeyID,
		KeyVaultNetworkAccess: keyVaultNetworkAccess,
	}}
}
