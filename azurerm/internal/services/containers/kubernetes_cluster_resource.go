package containers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-03-01/containerservice"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	computeValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/kubernetes"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/parse"
	containerValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/containers/validate"
	msiparse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
	msivalidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/validate"
	privateDnsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/privatedns/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceKubernetesCluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKubernetesClusterCreate,
		Read:   resourceKubernetesClusterRead,
		Update: resourceKubernetesClusterUpdate,
		Delete: resourceKubernetesClusterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ClusterID(id)
			return err
		}),

		CustomizeDiff: pluginsdk.CustomDiffInSequence(
			// Downgrade from Paid to Free is not supported and requires rebuild to apply
			pluginsdk.ForceNewIfChange("sku_tier", func(ctx context.Context, old, new, meta interface{}) bool {
				return new == "Free"
			}),
			// Migration of `identity` to `service_principal` is not allowed, the other way around is
			pluginsdk.ForceNewIfChange("service_principal.0.client_id", func(ctx context.Context, old, new, meta interface{}) bool {
				return old == "msi" || old == ""
			}),
		),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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

			"default_node_pool": SchemaDefaultNodePool(),

			// Optional
			"addon_profile": schemaKubernetesAddOnProfiles(),

			"api_server_authorized_ip_ranges": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.CIDR,
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
						"expander": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.ExpanderLeastWaste),
								string(containerservice.ExpanderMostPods),
								string(containerservice.ExpanderPriority),
								string(containerservice.ExpanderRandom),
							}, false),
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
							Default:  true,
						},
						"skip_nodes_with_system_pods": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"disk_encryption_set_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.DiskEncryptionSetID,
			},

			"enable_pod_security_policy": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"identity": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				ExactlyOneOf: []string{"identity", "service_principal"},
				MaxItems:     1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.ResourceIdentityTypeSystemAssigned),
								string(containerservice.ResourceIdentityTypeUserAssigned),
							}, false),
						},
						"user_assigned_identity_id": {
							Type:         pluginsdk.TypeString,
							ValidateFunc: msivalidate.UserAssignedIdentityID,
							Optional:     true,
						},
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							RequiredWith: []string{"kubelet_identity.0.object_id", "kubelet_identity.0.user_assigned_identity_id", "identity.0.user_assigned_identity_id"},
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"object_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							RequiredWith: []string{"kubelet_identity.0.client_id", "kubelet_identity.0.user_assigned_identity_id", "identity.0.user_assigned_identity_id"},
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"user_assigned_identity_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							RequiredWith: []string{"kubelet_identity.0.client_id", "kubelet_identity.0.object_id", "identity.0.user_assigned_identity_id"},
							ValidateFunc: msivalidate.UserAssignedIdentityID,
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
							ForceNew: true,
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

			"network_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"network_plugin": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.NetworkPluginAzure),
								string(containerservice.NetworkPluginKubenet),
							}, false),
						},

						"network_mode": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								// https://github.com/Azure/AKS/issues/1954#issuecomment-759306712
								// Transparent is already the default and only option for CNI
								// Bridge is only kept for backward compatibility
								string(containerservice.NetworkModeBridge),
								string(containerservice.NetworkModeTransparent),
							}, false),
						},

						"network_policy": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.NetworkPolicyCalico),
								string(containerservice.NetworkPolicyAzure),
							}, false),
						},

						"dns_service_ip": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validate.IPv4Address,
						},

						"docker_bridge_cidr": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validate.CIDR,
						},

						"pod_cidr": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validate.CIDR,
						},

						"service_cidr": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validate.CIDR,
						},

						"load_balancer_sku": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(containerservice.LoadBalancerSkuStandard),
							ForceNew: true,
							// TODO: fix the casing in the Swagger
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.LoadBalancerSkuBasic),
								string(containerservice.LoadBalancerSkuStandard),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"outbound_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  string(containerservice.OutboundTypeLoadBalancer),
							ValidateFunc: validation.StringInSlice([]string{
								string(containerservice.OutboundTypeLoadBalancer),
								string(containerservice.OutboundTypeUserDefinedRouting),
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
										ValidateFunc: validation.IntBetween(4, 120),
									},
									"managed_outbound_ip_count": {
										Type:          pluginsdk.TypeInt,
										Optional:      true,
										Computed:      true,
										ValidateFunc:  validation.IntBetween(1, 100),
										ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids", "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"},
									},
									"outbound_ip_prefix_ids": {
										Type:          pluginsdk.TypeSet,
										Optional:      true,
										Computed:      true,
										ConfigMode:    pluginsdk.SchemaConfigModeAttr,
										ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.managed_outbound_ip_count", "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"},
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: azure.ValidateResourceID,
										},
									},
									"outbound_ip_address_ids": {
										Type:          pluginsdk.TypeSet,
										Optional:      true,
										Computed:      true,
										ConfigMode:    pluginsdk.SchemaConfigModeAttr,
										ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.managed_outbound_ip_count", "network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids"},
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: azure.ValidateResourceID,
										},
									},
									"effective_outbound_ips": {
										Type:       pluginsdk.TypeSet,
										Computed:   true,
										ConfigMode: pluginsdk.SchemaConfigModeAttr,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},

			"node_resource_group": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"private_fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"private_link_enabled": {
				Type:          pluginsdk.TypeBool,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"private_cluster_enabled"},
				Deprecated:    "Deprecated in favour of `private_cluster_enabled`", // TODO -- remove this in next major version
			},

			"private_cluster_enabled": {
				Type:          pluginsdk.TypeBool,
				Optional:      true,
				ForceNew:      true,
				Computed:      true, // TODO -- remove this when deprecation resolves
				ConflictsWith: []string{"private_link_enabled"},
			},

			"private_dns_zone_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true, // a Private Cluster is `System` by default even if unspecified
				ForceNew: true,
				ValidateFunc: validation.Any(
					privateDnsValidate.PrivateDnsZoneID,
					validation.StringInSlice([]string{
						"System",
						"None",
					}, false),
				),
			},

			"role_based_access_control": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
							ForceNew: true,
						},
						"azure_active_directory": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"client_app_id": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsUUID,
										AtLeastOneOf: []string{"role_based_access_control.0.azure_active_directory.0.client_app_id", "role_based_access_control.0.azure_active_directory.0.server_app_id",
											"role_based_access_control.0.azure_active_directory.0.server_app_secret", "role_based_access_control.0.azure_active_directory.0.tenant_id",
											"role_based_access_control.0.azure_active_directory.0.managed", "role_based_access_control.0.azure_active_directory.0.admin_group_object_ids",
										},
									},

									"server_app_id": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsUUID,
										AtLeastOneOf: []string{"role_based_access_control.0.azure_active_directory.0.client_app_id", "role_based_access_control.0.azure_active_directory.0.server_app_id",
											"role_based_access_control.0.azure_active_directory.0.server_app_secret", "role_based_access_control.0.azure_active_directory.0.tenant_id",
											"role_based_access_control.0.azure_active_directory.0.managed", "role_based_access_control.0.azure_active_directory.0.admin_group_object_ids",
										},
									},

									"server_app_secret": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
										AtLeastOneOf: []string{"role_based_access_control.0.azure_active_directory.0.client_app_id", "role_based_access_control.0.azure_active_directory.0.server_app_id",
											"role_based_access_control.0.azure_active_directory.0.server_app_secret", "role_based_access_control.0.azure_active_directory.0.tenant_id",
											"role_based_access_control.0.azure_active_directory.0.managed", "role_based_access_control.0.azure_active_directory.0.admin_group_object_ids",
										},
									},

									"tenant_id": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Computed: true,
										// OrEmpty since this can be sourced from the client config if it's not specified
										ValidateFunc: validation.Any(validation.IsUUID, validation.StringIsEmpty),
										AtLeastOneOf: []string{"role_based_access_control.0.azure_active_directory.0.client_app_id", "role_based_access_control.0.azure_active_directory.0.server_app_id",
											"role_based_access_control.0.azure_active_directory.0.server_app_secret", "role_based_access_control.0.azure_active_directory.0.tenant_id",
											"role_based_access_control.0.azure_active_directory.0.managed", "role_based_access_control.0.azure_active_directory.0.admin_group_object_ids",
										},
									},

									"managed": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										AtLeastOneOf: []string{"role_based_access_control.0.azure_active_directory.0.client_app_id", "role_based_access_control.0.azure_active_directory.0.server_app_id",
											"role_based_access_control.0.azure_active_directory.0.server_app_secret", "role_based_access_control.0.azure_active_directory.0.tenant_id",
											"role_based_access_control.0.azure_active_directory.0.managed", "role_based_access_control.0.azure_active_directory.0.admin_group_object_ids",
										},
									},

									"azure_rbac_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
									},

									"admin_group_object_ids": {
										Type:       pluginsdk.TypeSet,
										Optional:   true,
										ConfigMode: pluginsdk.SchemaConfigModeAttr,
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.IsUUID,
										},
										AtLeastOneOf: []string{"role_based_access_control.0.azure_active_directory.0.client_app_id", "role_based_access_control.0.azure_active_directory.0.server_app_id",
											"role_based_access_control.0.azure_active_directory.0.server_app_secret", "role_based_access_control.0.azure_active_directory.0.tenant_id",
											"role_based_access_control.0.azure_active_directory.0.managed", "role_based_access_control.0.azure_active_directory.0.admin_group_object_ids",
										},
									},
								},
							},
						},
					},
				},
			},

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

			"sku_tier": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// @tombuildsstuff (2020-05-29) - Preview limitations:
				//  * Currently, there is no way to remove Uptime SLA from an AKS cluster after creation with it enabled.
				//  * Private clusters aren't currently supported.
				// @jackofallops (2020-07-21) - Update:
				//  * sku_tier can now be upgraded in place, downgrade requires rebuild
				Default: string(containerservice.ManagedClusterSKUTierFree),
				ValidateFunc: validation.StringInSlice([]string{
					string(containerservice.ManagedClusterSKUTierFree),
					string(containerservice.ManagedClusterSKUTierPaid),
				}, false),
			},

			"tags": tags.Schema(),

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
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringLenBetween(8, 123),
						},
					},
				},
			},

			"automatic_channel_upgrade": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(containerservice.UpgradeChannelPatch),
					string(containerservice.UpgradeChannelRapid),
					string(containerservice.UpgradeChannelStable),
				}, false),
			},

			// Computed
			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"kube_admin_config": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"host": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"username": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"host": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"username": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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
		},
	}
}

func resourceKubernetesClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.KubernetesClustersClient
	env := meta.(*clients.Client).Containers.Environment
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	tenantId := meta.(*clients.Client).Account.TenantId

	log.Printf("[INFO] preparing arguments for Managed Kubernetes Cluster create.")

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	existing, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Kubernetes Cluster %q (Resource Group %q): %s", name, resGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_kubernetes_cluster", *existing.ID)
	}

	if err := validateKubernetesCluster(d, nil, resGroup, name); err != nil {
		return err
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	dnsPrefix := d.Get("dns_prefix").(string)
	kubernetesVersion := d.Get("kubernetes_version").(string)

	linuxProfileRaw := d.Get("linux_profile").([]interface{})
	linuxProfile := expandKubernetesClusterLinuxProfile(linuxProfileRaw)

	// NOTE: we /could/ validate the default node pool version here - but since the entire cluster deployment
	// will fail here this should be fine to omit for the Create
	agentProfiles, err := ExpandDefaultNodePool(d)
	if err != nil {
		return fmt.Errorf("expanding `default_node_pool`: %+v", err)
	}

	addOnProfilesRaw := d.Get("addon_profile").([]interface{})
	addonProfiles, err := expandKubernetesAddOnProfiles(addOnProfilesRaw, env)
	if err != nil {
		return err
	}

	networkProfileRaw := d.Get("network_profile").([]interface{})
	networkProfile, err := expandKubernetesClusterNetworkProfile(networkProfileRaw)
	if err != nil {
		return err
	}

	rbacRaw := d.Get("role_based_access_control").([]interface{})
	rbacEnabled, azureADProfile, err := expandKubernetesClusterRoleBasedAccessControl(rbacRaw, tenantId)
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]interface{})

	windowsProfileRaw := d.Get("windows_profile").([]interface{})
	windowsProfile := expandKubernetesClusterWindowsProfile(windowsProfileRaw)

	apiServerAuthorizedIPRangesRaw := d.Get("api_server_authorized_ip_ranges").(*pluginsdk.Set).List()
	apiServerAuthorizedIPRanges := utils.ExpandStringSlice(apiServerAuthorizedIPRangesRaw)

	enablePrivateCluster := false
	if v, ok := d.GetOk("private_link_enabled"); ok {
		enablePrivateCluster = v.(bool)
	}
	if v, ok := d.GetOk("private_cluster_enabled"); ok {
		enablePrivateCluster = v.(bool)
	}

	if !enablePrivateCluster && dnsPrefix == "" {
		return fmt.Errorf("`dns_prefix` should be set if it is not a private cluster")
	}

	apiAccessProfile := containerservice.ManagedClusterAPIServerAccessProfile{
		EnablePrivateCluster: &enablePrivateCluster,
		AuthorizedIPRanges:   apiServerAuthorizedIPRanges,
	}

	nodeResourceGroup := d.Get("node_resource_group").(string)

	if d.Get("enable_pod_security_policy").(bool) {
		return fmt.Errorf("The AKS API has removed support for this field on 2020-10-15 and is no longer possible to configure this the Pod Security Policy - as such you'll need to set `enable_pod_security_policy` to `false`")
	}

	autoScalerProfileRaw := d.Get("auto_scaler_profile").([]interface{})
	autoScalerProfile := expandKubernetesClusterAutoScalerProfile(autoScalerProfileRaw)

	parameters := containerservice.ManagedCluster{
		Name:     &name,
		Location: &location,
		Sku: &containerservice.ManagedClusterSKU{
			Name: containerservice.ManagedClusterSKUNameBasic, // the only possible value at this point
			Tier: containerservice.ManagedClusterSKUTier(d.Get("sku_tier").(string)),
		},
		ManagedClusterProperties: &containerservice.ManagedClusterProperties{
			APIServerAccessProfile: &apiAccessProfile,
			AadProfile:             azureADProfile,
			AddonProfiles:          *addonProfiles,
			AgentPoolProfiles:      agentProfiles,
			AutoScalerProfile:      autoScalerProfile,
			DNSPrefix:              utils.String(dnsPrefix),
			EnableRBAC:             utils.Bool(rbacEnabled),
			KubernetesVersion:      utils.String(kubernetesVersion),
			LinuxProfile:           linuxProfile,
			WindowsProfile:         windowsProfile,
			NetworkProfile:         networkProfile,
			NodeResourceGroup:      utils.String(nodeResourceGroup),
		},
		Tags: tags.Expand(t),
	}

	if v := d.Get("automatic_channel_upgrade").(string); v != "" {
		parameters.ManagedClusterProperties.AutoUpgradeProfile = &containerservice.ManagedClusterAutoUpgradeProfile{
			UpgradeChannel: containerservice.UpgradeChannel(v),
		}
	}

	managedClusterIdentityRaw := d.Get("identity").([]interface{})
	kubernetesClusterIdentityRaw := d.Get("kubelet_identity").([]interface{})
	servicePrincipalProfileRaw := d.Get("service_principal").([]interface{})

	if len(managedClusterIdentityRaw) == 0 && len(servicePrincipalProfileRaw) == 0 {
		return fmt.Errorf("either an `identity` or `service_principal` block must be specified for cluster authentication")
	}

	if len(managedClusterIdentityRaw) > 0 {
		parameters.Identity = expandKubernetesClusterManagedClusterIdentity(managedClusterIdentityRaw)
		parameters.ManagedClusterProperties.ServicePrincipalProfile = &containerservice.ManagedClusterServicePrincipalProfile{
			ClientID: utils.String("msi"),
		}
	}
	if len(kubernetesClusterIdentityRaw) > 0 {
		parameters.ManagedClusterProperties.IdentityProfile = expandKubernetesClusterIdentityProfile(kubernetesClusterIdentityRaw)
	}

	servicePrincipalSet := false
	if len(servicePrincipalProfileRaw) > 0 {
		servicePrincipalProfileVal := servicePrincipalProfileRaw[0].(map[string]interface{})
		parameters.ManagedClusterProperties.ServicePrincipalProfile = &containerservice.ManagedClusterServicePrincipalProfile{
			ClientID: utils.String(servicePrincipalProfileVal["client_id"].(string)),
			Secret:   utils.String(servicePrincipalProfileVal["client_secret"].(string)),
		}
		servicePrincipalSet = true
	}

	if v, ok := d.GetOk("private_dns_zone_id"); ok {
		if (parameters.Identity == nil && !servicePrincipalSet) || (v.(string) != "System" && v.(string) != "None" && (!servicePrincipalSet && parameters.Identity.Type != containerservice.ResourceIdentityTypeUserAssigned)) {
			return fmt.Errorf("a user assigned identity or a service principal must be used when using a custom private dns zone")
		}
		apiAccessProfile.PrivateDNSZone = utils.String(v.(string))
	}

	if v, ok := d.GetOk("dns_prefix_private_cluster"); ok {
		if !enablePrivateCluster || apiAccessProfile.PrivateDNSZone == nil || *apiAccessProfile.PrivateDNSZone == "System" || *apiAccessProfile.PrivateDNSZone == "None" {
			return fmt.Errorf("`dns_prefix_private_cluster` should only be set for private cluster with custom private dns zone")
		}
		parameters.FqdnSubdomain = utils.String(v.(string))
	}

	if v, ok := d.GetOk("disk_encryption_set_id"); ok && v.(string) != "" {
		parameters.ManagedClusterProperties.DiskEncryptionSetID = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("creating Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Managed Kubernetes Cluster %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read ID for Managed Kubernetes Cluster %q (Resource Group %q)", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceKubernetesClusterRead(d, meta)
}

func resourceKubernetesClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	containersClient := meta.(*clients.Client).Containers
	nodePoolsClient := containersClient.AgentPoolsClient
	clusterClient := containersClient.KubernetesClustersClient
	env := containersClient.Environment
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	tenantId := meta.(*clients.Client).Account.TenantId

	log.Printf("[INFO] preparing arguments for Managed Kubernetes Cluster update.")

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	d.Partial(true)

	// we need to conditionally update the cluster
	existing, err := clusterClient.Get(ctx, id.ResourceGroup, id.ManagedClusterName)
	if err != nil {
		return fmt.Errorf("retrieving existing Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
	}
	if existing.ManagedClusterProperties == nil {
		return fmt.Errorf("retrieving existing Kubernetes Cluster %q (Resource Group %q): `properties` was nil", id.ManagedClusterName, id.ResourceGroup)
	}

	if err := validateKubernetesCluster(d, &existing, id.ResourceGroup, id.ManagedClusterName); err != nil {
		return err
	}

	// when update, we should set the value of `Identity.UserAssignedIdentities` empty
	// otherwise the rest api will report error - this is tracked here: https://github.com/Azure/azure-rest-api-specs/issues/13631
	if existing.Identity != nil && existing.Identity.UserAssignedIdentities != nil {
		for k := range existing.Identity.UserAssignedIdentities {
			existing.Identity.UserAssignedIdentities[k] = &containerservice.ManagedClusterIdentityUserAssignedIdentitiesValue{}
		}
	}

	if d.HasChange("service_principal") && !d.HasChange("identity") {
		log.Printf("[DEBUG] Updating the Service Principal for Kubernetes Cluster %q (Resource Group %q)..", id.ManagedClusterName, id.ResourceGroup)
		servicePrincipals := d.Get("service_principal").([]interface{})
		// we'll be rotating the Service Principal - removing the SP block is handled by the validate function
		servicePrincipalRaw := servicePrincipals[0].(map[string]interface{})

		clientId := servicePrincipalRaw["client_id"].(string)
		clientSecret := servicePrincipalRaw["client_secret"].(string)
		params := containerservice.ManagedClusterServicePrincipalProfile{
			ClientID: utils.String(clientId),
			Secret:   utils.String(clientSecret),
		}

		future, err := clusterClient.ResetServicePrincipalProfile(ctx, id.ResourceGroup, id.ManagedClusterName, params)
		if err != nil {
			return fmt.Errorf("updating Service Principal for Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, clusterClient.Client); err != nil {
			return fmt.Errorf("waiting for update of Service Principal for Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
		}
		log.Printf("[DEBUG] Updated the Service Principal for Kubernetes Cluster %q (Resource Group %q).", id.ManagedClusterName, id.ResourceGroup)

		// since we're patching it, re-retrieve the latest version of the cluster
		existing, err = clusterClient.Get(ctx, id.ResourceGroup, id.ManagedClusterName)
		if err != nil {
			return fmt.Errorf("retrieving existing Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
		}
		if existing.ManagedClusterProperties == nil {
			return fmt.Errorf("retrieving existing Kubernetes Cluster %q (Resource Group %q): `properties` was nil", id.ManagedClusterName, id.ResourceGroup)
		}
	}

	// since there's multiple reasons why we could be called into Update, we use this to only update if something's changed that's not SP/Version
	updateCluster := false

	// RBAC profile updates need to be handled atomically before any call to createUpdate as a diff there will create a PropertyChangeNotAllowed error
	if d.HasChange("role_based_access_control") {
		props := existing.ManagedClusterProperties
		// check if we can determine current EnableRBAC state - don't do anything destructive if we can't be sure
		if props.EnableRBAC == nil {
			return fmt.Errorf("reading current state of RBAC Enabled, expected bool got %+v", props.EnableRBAC)
		}
		rbacRaw := d.Get("role_based_access_control").([]interface{})
		rbacEnabled, azureADProfile, err := expandKubernetesClusterRoleBasedAccessControl(rbacRaw, tenantId)
		if err != nil {
			return err
		}

		// changing rbacEnabled must still force cluster recreation
		if *props.EnableRBAC == rbacEnabled {
			props.AadProfile = azureADProfile
			props.EnableRBAC = utils.Bool(rbacEnabled)

			// Reset AAD profile is only possible if not managed
			if props.AadProfile.Managed == nil || !*props.AadProfile.Managed {
				log.Printf("[DEBUG] Updating the RBAC AAD profile")
				future, err := clusterClient.ResetAADProfile(ctx, id.ResourceGroup, id.ManagedClusterName, *props.AadProfile)
				if err != nil {
					return fmt.Errorf("updating Managed Kubernetes Cluster AAD Profile in cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
				}

				if err = future.WaitForCompletionRef(ctx, clusterClient.Client); err != nil {
					return fmt.Errorf("waiting for update of RBAC AAD profile of Managed Cluster %q (Resource Group %q):, %+v", id.ManagedClusterName, id.ResourceGroup, err)
				}
			}
		} else {
			updateCluster = true
		}

		if props.AadProfile != nil && props.AadProfile.Managed != nil && *props.AadProfile.Managed {
			existing.ManagedClusterProperties.AadProfile = azureADProfile
			updateCluster = true
		}
	}

	if d.HasChange("addon_profile") {
		updateCluster = true
		addOnProfilesRaw := d.Get("addon_profile").([]interface{})
		addonProfiles, err := expandKubernetesAddOnProfiles(addOnProfilesRaw, env)
		if err != nil {
			return err
		}

		existing.ManagedClusterProperties.AddonProfiles = *addonProfiles
	}

	if d.HasChange("api_server_authorized_ip_ranges") {
		updateCluster = true
		apiServerAuthorizedIPRangesRaw := d.Get("api_server_authorized_ip_ranges").(*pluginsdk.Set).List()

		enablePrivateCluster := false
		if v, ok := d.GetOk("private_link_enabled"); ok {
			enablePrivateCluster = v.(bool)
		}
		if v, ok := d.GetOk("private_cluster_enabled"); ok {
			enablePrivateCluster = v.(bool)
		}
		existing.ManagedClusterProperties.APIServerAccessProfile = &containerservice.ManagedClusterAPIServerAccessProfile{
			AuthorizedIPRanges:   utils.ExpandStringSlice(apiServerAuthorizedIPRangesRaw),
			EnablePrivateCluster: &enablePrivateCluster,
		}
		if v, ok := d.GetOk("private_dns_zone_id"); ok {
			existing.ManagedClusterProperties.APIServerAccessProfile.PrivateDNSZone = utils.String(v.(string))
		}
	}

	if d.HasChange("auto_scaler_profile") {
		updateCluster = true
		autoScalerProfileRaw := d.Get("auto_scaler_profile").([]interface{})

		autoScalerProfile := expandKubernetesClusterAutoScalerProfile(autoScalerProfileRaw)
		existing.ManagedClusterProperties.AutoScalerProfile = autoScalerProfile
	}

	if d.HasChange("enable_pod_security_policy") && d.Get("enable_pod_security_policy").(bool) {
		return fmt.Errorf("The AKS API has removed support for this field on 2020-10-15 and is no longer possible to configure this the Pod Security Policy - as such you'll need to set `enable_pod_security_policy` to `false`")
	}

	if d.HasChange("linux_profile") {
		updateCluster = true
		linuxProfileRaw := d.Get("linux_profile").([]interface{})
		linuxProfile := expandKubernetesClusterLinuxProfile(linuxProfileRaw)
		existing.ManagedClusterProperties.LinuxProfile = linuxProfile
	}

	if d.HasChange("network_profile") {
		updateCluster = true

		networkProfile := *existing.ManagedClusterProperties.NetworkProfile
		if networkProfile.LoadBalancerProfile == nil {
			// an existing LB Profile must be present, since it's Optional & Computed
			return fmt.Errorf("`loadBalancerProfile` was nil in Azure")
		}

		loadBalancerProfile := *networkProfile.LoadBalancerProfile

		if key := "network_profile.0.load_balancer_profile.0.effective_outbound_ips"; d.HasChange(key) {
			effectiveOutboundIPs := idsToResourceReferences(d.Get(key))
			loadBalancerProfile.EffectiveOutboundIPs = effectiveOutboundIPs
		}

		if key := "network_profile.0.load_balancer_profile.0.idle_timeout_in_minutes"; d.HasChange(key) {
			idleTimeoutInMinutes := d.Get(key).(int)
			loadBalancerProfile.IdleTimeoutInMinutes = utils.Int32(int32(idleTimeoutInMinutes))
		}

		if key := "network_profile.0.load_balancer_profile.0.managed_outbound_ip_count"; d.HasChange(key) {
			managedOutboundIPCount := d.Get(key).(int)
			loadBalancerProfile.ManagedOutboundIPs = &containerservice.ManagedClusterLoadBalancerProfileManagedOutboundIPs{
				Count: utils.Int32(int32(managedOutboundIPCount)),
			}

			// fixes: Load balancer profile must specify one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs.
			loadBalancerProfile.OutboundIPs = nil
			loadBalancerProfile.OutboundIPPrefixes = nil
		}

		if key := "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"; d.HasChange(key) {
			publicIPAddressIDs := idsToResourceReferences(d.Get(key))
			loadBalancerProfile.OutboundIPs = &containerservice.ManagedClusterLoadBalancerProfileOutboundIPs{
				PublicIPs: publicIPAddressIDs,
			}

			// fixes: Load balancer profile must specify one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs.
			loadBalancerProfile.ManagedOutboundIPs = nil
			loadBalancerProfile.OutboundIPPrefixes = nil
		}

		if key := "network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids"; d.HasChange(key) {
			outboundIPPrefixIDs := idsToResourceReferences(d.Get(key))
			loadBalancerProfile.OutboundIPPrefixes = &containerservice.ManagedClusterLoadBalancerProfileOutboundIPPrefixes{
				PublicIPPrefixes: outboundIPPrefixIDs,
			}

			// fixes: Load balancer profile must specify one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs.
			loadBalancerProfile.ManagedOutboundIPs = nil
			loadBalancerProfile.OutboundIPs = nil
		}

		if key := "network_profile.0.load_balancer_profile.0.outbound_ports_allocated"; d.HasChange(key) {
			allocatedOutboundPorts := d.Get(key).(int)
			loadBalancerProfile.AllocatedOutboundPorts = utils.Int32(int32(allocatedOutboundPorts))
		}

		existing.ManagedClusterProperties.NetworkProfile.LoadBalancerProfile = &loadBalancerProfile
	}

	if d.HasChange("tags") {
		updateCluster = true
		t := d.Get("tags").(map[string]interface{})
		existing.Tags = tags.Expand(t)
	}

	if d.HasChange("windows_profile") {
		updateCluster = true
		windowsProfileRaw := d.Get("windows_profile").([]interface{})
		windowsProfile := expandKubernetesClusterWindowsProfile(windowsProfileRaw)
		existing.ManagedClusterProperties.WindowsProfile = windowsProfile
	}

	if d.HasChange("identity") {
		updateCluster = true
		managedClusterIdentityRaw := d.Get("identity").([]interface{})
		existing.Identity = expandKubernetesClusterManagedClusterIdentity(managedClusterIdentityRaw)
	}

	if d.HasChange("sku_tier") {
		updateCluster = true
		existing.Sku.Tier = containerservice.ManagedClusterSKUTier(d.Get("sku_tier").(string))
	}

	if d.HasChange("automatic_channel_upgrade") {
		updateCluster = true
		if existing.ManagedClusterProperties.AutoUpgradeProfile == nil {
			existing.ManagedClusterProperties.AutoUpgradeProfile = &containerservice.ManagedClusterAutoUpgradeProfile{}
		}

		channel := containerservice.UpgradeChannelNone
		if v := d.Get("automatic_channel_upgrade").(string); v != "" {
			channel = containerservice.UpgradeChannel(v)
		}

		existing.ManagedClusterProperties.AutoUpgradeProfile.UpgradeChannel = channel
	}

	if updateCluster {
		log.Printf("[DEBUG] Updating the Kubernetes Cluster %q (Resource Group %q)..", id.ManagedClusterName, id.ResourceGroup)
		future, err := clusterClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedClusterName, existing)
		if err != nil {
			return fmt.Errorf("updating Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, clusterClient.Client); err != nil {
			return fmt.Errorf("waiting for update of Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
		}
		log.Printf("[DEBUG] Updated the Kubernetes Cluster %q (Resource Group %q)..", id.ManagedClusterName, id.ResourceGroup)
	}

	// then roll the version of Kubernetes if necessary
	if d.HasChange("kubernetes_version") {
		existing, err = clusterClient.Get(ctx, id.ResourceGroup, id.ManagedClusterName)
		if err != nil {
			return fmt.Errorf("retrieving existing Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
		}
		if existing.ManagedClusterProperties == nil {
			return fmt.Errorf("retrieving existing Kubernetes Cluster %q (Resource Group %q): `properties` was nil", id.ManagedClusterName, id.ResourceGroup)
		}

		kubernetesVersion := d.Get("kubernetes_version").(string)
		log.Printf("[DEBUG] Upgrading the version of Kubernetes to %q..", kubernetesVersion)
		existing.ManagedClusterProperties.KubernetesVersion = utils.String(kubernetesVersion)

		future, err := clusterClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedClusterName, existing)
		if err != nil {
			return fmt.Errorf("updating Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, clusterClient.Client); err != nil {
			return fmt.Errorf("waiting for update of Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
		}

		log.Printf("[DEBUG] Upgraded the version of Kubernetes to %q..", kubernetesVersion)
	}

	// update the node pool using the separate API
	if d.HasChange("default_node_pool") {
		log.Printf("[DEBUG] Updating of Default Node Pool..")

		agentProfiles, err := ExpandDefaultNodePool(d)
		if err != nil {
			return fmt.Errorf("expanding `default_node_pool`: %+v", err)
		}

		agentProfile := ConvertDefaultNodePoolToAgentPool(agentProfiles)
		nodePoolName := *agentProfile.Name

		// if a users specified a version - confirm that version is supported on the cluster
		if nodePoolVersion := agentProfile.ManagedClusterAgentPoolProfileProperties.OrchestratorVersion; nodePoolVersion != nil {
			if err := validateNodePoolSupportsVersion(ctx, containersClient, id.ResourceGroup, id.ManagedClusterName, nodePoolName, *nodePoolVersion); err != nil {
				return err
			}
		}

		agentPool, err := nodePoolsClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedClusterName, nodePoolName, agentProfile)
		if err != nil {
			return fmt.Errorf("updating Default Node Pool %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
		}

		if err := agentPool.WaitForCompletionRef(ctx, nodePoolsClient.Client); err != nil {
			return fmt.Errorf("waiting for update of Default Node Pool %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
		}
		log.Printf("[DEBUG] Updated Default Node Pool.")
	}

	d.Partial(false)

	return resourceKubernetesClusterRead(d, meta)
}

func resourceKubernetesClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ManagedClusterName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Managed Kubernetes Cluster %q was not found in Resource Group %q - removing from state!", id.ManagedClusterName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
	}

	profile, err := client.GetAccessProfile(ctx, id.ResourceGroup, id.ManagedClusterName, "clusterUser")
	if err != nil {
		return fmt.Errorf("retrieving Access Profile for Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	skuTier := string(containerservice.ManagedClusterSKUTierFree)
	if resp.Sku != nil && resp.Sku.Tier != "" {
		skuTier = string(resp.Sku.Tier)
	}
	d.Set("sku_tier", skuTier)

	if props := resp.ManagedClusterProperties; props != nil {
		d.Set("dns_prefix", props.DNSPrefix)
		d.Set("dns_prefix_private_cluster", props.FqdnSubdomain)
		d.Set("fqdn", props.Fqdn)
		d.Set("private_fqdn", props.PrivateFQDN)
		d.Set("disk_encryption_set_id", props.DiskEncryptionSetID)
		d.Set("kubernetes_version", props.KubernetesVersion)
		d.Set("node_resource_group", props.NodeResourceGroup)
		d.Set("enable_pod_security_policy", props.EnablePodSecurityPolicy)

		upgradeChannel := ""
		if profile := props.AutoUpgradeProfile; profile != nil && profile.UpgradeChannel != containerservice.UpgradeChannelNone {
			upgradeChannel = string(profile.UpgradeChannel)
		}
		d.Set("automatic_channel_upgrade", upgradeChannel)

		// TODO: 2.0 we should introduce a access_profile block to match the new API design,
		if accessProfile := props.APIServerAccessProfile; accessProfile != nil {
			apiServerAuthorizedIPRanges := utils.FlattenStringSlice(accessProfile.AuthorizedIPRanges)
			if err := d.Set("api_server_authorized_ip_ranges", apiServerAuthorizedIPRanges); err != nil {
				return fmt.Errorf("setting `api_server_authorized_ip_ranges`: %+v", err)
			}

			d.Set("private_link_enabled", accessProfile.EnablePrivateCluster)
			d.Set("private_cluster_enabled", accessProfile.EnablePrivateCluster)
			d.Set("private_dns_zone_id", accessProfile.PrivateDNSZone)
		}

		addonProfiles := flattenKubernetesAddOnProfiles(props.AddonProfiles)
		if err := d.Set("addon_profile", addonProfiles); err != nil {
			return fmt.Errorf("setting `addon_profile`: %+v", err)
		}

		autoScalerProfile, err := flattenKubernetesClusterAutoScalerProfile(props.AutoScalerProfile)
		if err != nil {
			return err
		}
		if err := d.Set("auto_scaler_profile", autoScalerProfile); err != nil {
			return fmt.Errorf("setting `auto_scaler_profile`: %+v", err)
		}

		flattenedDefaultNodePool, err := FlattenDefaultNodePool(props.AgentPoolProfiles, d)
		if err != nil {
			return fmt.Errorf("flattening `default_node_pool`: %+v", err)
		}
		if err := d.Set("default_node_pool", flattenedDefaultNodePool); err != nil {
			return fmt.Errorf("setting `default_node_pool`: %+v", err)
		}

		kubeletIdentity, err := flattenKubernetesClusterIdentityProfile(props.IdentityProfile)
		if err != nil {
			return err
		}
		if err := d.Set("kubelet_identity", kubeletIdentity); err != nil {
			return fmt.Errorf("setting `kubelet_identity`: %+v", err)
		}

		linuxProfile := flattenKubernetesClusterLinuxProfile(props.LinuxProfile)
		if err := d.Set("linux_profile", linuxProfile); err != nil {
			return fmt.Errorf("setting `linux_profile`: %+v", err)
		}

		networkProfile := flattenKubernetesClusterNetworkProfile(props.NetworkProfile)
		if err := d.Set("network_profile", networkProfile); err != nil {
			return fmt.Errorf("setting `network_profile`: %+v", err)
		}

		roleBasedAccessControl := flattenKubernetesClusterRoleBasedAccessControl(props, d)
		if err := d.Set("role_based_access_control", roleBasedAccessControl); err != nil {
			return fmt.Errorf("setting `role_based_access_control`: %+v", err)
		}

		servicePrincipal := flattenAzureRmKubernetesClusterServicePrincipalProfile(props.ServicePrincipalProfile, d)
		if err := d.Set("service_principal", servicePrincipal); err != nil {
			return fmt.Errorf("setting `service_principal`: %+v", err)
		}

		windowsProfile := flattenKubernetesClusterWindowsProfile(props.WindowsProfile, d)
		if err := d.Set("windows_profile", windowsProfile); err != nil {
			return fmt.Errorf("setting `windows_profile`: %+v", err)
		}

		// adminProfile is only available for RBAC enabled clusters with AAD
		if props.AadProfile != nil {
			adminProfile, err := client.GetAccessProfile(ctx, id.ResourceGroup, id.ManagedClusterName, "clusterAdmin")
			if err != nil {
				return fmt.Errorf("retrieving Admin Access Profile for Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
			}

			adminKubeConfigRaw, adminKubeConfig := flattenKubernetesClusterAccessProfile(adminProfile)
			d.Set("kube_admin_config_raw", adminKubeConfigRaw)
			if err := d.Set("kube_admin_config", adminKubeConfig); err != nil {
				return fmt.Errorf("setting `kube_admin_config`: %+v", err)
			}
		} else {
			d.Set("kube_admin_config_raw", "")
			d.Set("kube_admin_config", []interface{}{})
		}
	}

	identity, err := flattenKubernetesClusterManagedClusterIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	kubeConfigRaw, kubeConfig := flattenKubernetesClusterAccessProfile(profile)
	d.Set("kube_config_raw", kubeConfigRaw)
	if err := d.Set("kube_config", kubeConfig); err != nil {
		return fmt.Errorf("setting `kube_config`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceKubernetesClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ClusterID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ManagedClusterName)
	if err != nil {
		return fmt.Errorf("deleting Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of Managed Kubernetes Cluster %q (Resource Group %q): %+v", id.ManagedClusterName, id.ResourceGroup, err)
	}

	return nil
}

func flattenKubernetesClusterAccessProfile(profile containerservice.ManagedClusterAccessProfile) (*string, []interface{}) {
	if accessProfile := profile.AccessProfile; accessProfile != nil {
		if kubeConfigRaw := accessProfile.KubeConfig; kubeConfigRaw != nil {
			rawConfig := string(*kubeConfigRaw)
			var flattenedKubeConfig []interface{}

			if strings.Contains(rawConfig, "apiserver-id:") {
				kubeConfigAAD, err := kubernetes.ParseKubeConfigAAD(rawConfig)
				if err != nil {
					return utils.String(rawConfig), []interface{}{}
				}

				flattenedKubeConfig = flattenKubernetesClusterKubeConfigAAD(*kubeConfigAAD)
			} else {
				kubeConfig, err := kubernetes.ParseKubeConfig(rawConfig)
				if err != nil {
					return utils.String(rawConfig), []interface{}{}
				}

				flattenedKubeConfig = flattenKubernetesClusterKubeConfig(*kubeConfig)
			}

			return utils.String(rawConfig), flattenedKubeConfig
		}
	}
	return nil, []interface{}{}
}

func expandKubernetesClusterLinuxProfile(input []interface{}) *containerservice.LinuxProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})

	adminUsername := config["admin_username"].(string)
	linuxKeys := config["ssh_key"].([]interface{})

	keyData := ""
	if key, ok := linuxKeys[0].(map[string]interface{}); ok {
		keyData = key["key_data"].(string)
	}

	return &containerservice.LinuxProfile{
		AdminUsername: &adminUsername,
		SSH: &containerservice.SSHConfiguration{
			PublicKeys: &[]containerservice.SSHPublicKey{
				{
					KeyData: &keyData,
				},
			},
		},
	}
}

func expandKubernetesClusterIdentityProfile(input []interface{}) map[string]*containerservice.ManagedClusterPropertiesIdentityProfileValue {
	identityProfile := make(map[string]*containerservice.ManagedClusterPropertiesIdentityProfileValue)
	if len(input) == 0 || input[0] == nil {
		return identityProfile
	}

	values := input[0].(map[string]interface{})

	if containerservice.ResourceIdentityType(values["user_assigned_identity_id"].(string)) != "" {
		identityProfile["kubeletidentity"] = &containerservice.ManagedClusterPropertiesIdentityProfileValue{
			ResourceID: utils.String(values["user_assigned_identity_id"].(string)),
			ClientID:   utils.String(values["client_id"].(string)),
			ObjectID:   utils.String(values["object_id"].(string)),
		}
	}

	return identityProfile
}

func flattenKubernetesClusterIdentityProfile(profile map[string]*containerservice.ManagedClusterPropertiesIdentityProfileValue) ([]interface{}, error) {
	if profile == nil {
		return []interface{}{}, nil
	}

	kubeletIdentity := make([]interface{}, 0)
	if kubeletidentity := profile["kubeletidentity"]; kubeletidentity != nil {
		clientId := ""
		if clientid := kubeletidentity.ClientID; clientid != nil {
			clientId = *clientid
		}

		objectId := ""
		if objectid := kubeletidentity.ObjectID; objectid != nil {
			objectId = *objectid
		}

		userAssignedIdentityId := ""
		if resourceid := kubeletidentity.ResourceID; resourceid != nil {
			parsedId, err := msiparse.UserAssignedIdentityID(*resourceid)
			if err != nil {
				return nil, err
			}

			userAssignedIdentityId = parsedId.ID()
		}

		kubeletIdentity = append(kubeletIdentity, map[string]interface{}{
			"client_id":                 clientId,
			"object_id":                 objectId,
			"user_assigned_identity_id": userAssignedIdentityId,
		})
	}

	return kubeletIdentity, nil
}

func flattenKubernetesClusterLinuxProfile(profile *containerservice.LinuxProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	adminUsername := ""
	if username := profile.AdminUsername; username != nil {
		adminUsername = *username
	}

	sshKeys := make([]interface{}, 0)
	if ssh := profile.SSH; ssh != nil {
		if keys := ssh.PublicKeys; keys != nil {
			for _, sshKey := range *keys {
				keyData := ""
				if kd := sshKey.KeyData; kd != nil {
					keyData = *kd
				}
				sshKeys = append(sshKeys, map[string]interface{}{
					"key_data": keyData,
				})
			}
		}
	}

	return []interface{}{
		map[string]interface{}{
			"admin_username": adminUsername,
			"ssh_key":        sshKeys,
		},
	}
}

func expandKubernetesClusterWindowsProfile(input []interface{}) *containerservice.ManagedClusterWindowsProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})

	adminUsername := config["admin_username"].(string)
	adminPassword := config["admin_password"].(string)

	profile := containerservice.ManagedClusterWindowsProfile{
		AdminUsername: &adminUsername,
		AdminPassword: &adminPassword,
	}

	return &profile
}

func flattenKubernetesClusterWindowsProfile(profile *containerservice.ManagedClusterWindowsProfile, d *pluginsdk.ResourceData) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	adminUsername := ""
	if username := profile.AdminUsername; username != nil {
		adminUsername = *username
	}

	// admin password isn't returned, so let's look it up
	adminPassword := ""
	if v, ok := d.GetOk("windows_profile.0.admin_password"); ok {
		adminPassword = v.(string)
	}

	return []interface{}{
		map[string]interface{}{
			"admin_password": adminPassword,
			"admin_username": adminUsername,
		},
	}
}

func expandKubernetesClusterNetworkProfile(input []interface{}) (*containerservice.NetworkProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	config := input[0].(map[string]interface{})

	networkPlugin := config["network_plugin"].(string)
	networkMode := config["network_mode"].(string)
	if networkPlugin != "azure" && networkMode != "" {
		return nil, fmt.Errorf("`network_mode` cannot be set if `network_plugin` is not `azure`")
	}
	networkPolicy := config["network_policy"].(string)
	loadBalancerProfileRaw := config["load_balancer_profile"].([]interface{})
	loadBalancerSku := config["load_balancer_sku"].(string)
	outboundType := config["outbound_type"].(string)

	networkProfile := containerservice.NetworkProfile{
		NetworkPlugin:   containerservice.NetworkPlugin(networkPlugin),
		NetworkMode:     containerservice.NetworkMode(networkMode),
		NetworkPolicy:   containerservice.NetworkPolicy(networkPolicy),
		LoadBalancerSku: containerservice.LoadBalancerSku(loadBalancerSku),
		OutboundType:    containerservice.OutboundType(outboundType),
	}

	if len(loadBalancerProfileRaw) > 0 {
		if !strings.EqualFold(loadBalancerSku, "standard") {
			return nil, fmt.Errorf("only load balancer SKU 'Standard' supports load balancer profiles. Provided load balancer type: %s", loadBalancerSku)
		}

		networkProfile.LoadBalancerProfile = expandLoadBalancerProfile(loadBalancerProfileRaw)
	}

	if v, ok := config["dns_service_ip"]; ok && v.(string) != "" {
		dnsServiceIP := v.(string)
		networkProfile.DNSServiceIP = utils.String(dnsServiceIP)
	}

	if v, ok := config["pod_cidr"]; ok && v.(string) != "" {
		podCidr := v.(string)
		networkProfile.PodCidr = utils.String(podCidr)
	}

	if v, ok := config["docker_bridge_cidr"]; ok && v.(string) != "" {
		dockerBridgeCidr := v.(string)
		networkProfile.DockerBridgeCidr = utils.String(dockerBridgeCidr)
	}

	if v, ok := config["service_cidr"]; ok && v.(string) != "" {
		serviceCidr := v.(string)
		networkProfile.ServiceCidr = utils.String(serviceCidr)
	}

	return &networkProfile, nil
}

func expandLoadBalancerProfile(d []interface{}) *containerservice.ManagedClusterLoadBalancerProfile {
	if d[0] == nil {
		return nil
	}

	config := d[0].(map[string]interface{})

	profile := &containerservice.ManagedClusterLoadBalancerProfile{}

	if mins, ok := config["idle_timeout_in_minutes"]; ok && mins.(int) != 0 {
		profile.IdleTimeoutInMinutes = utils.Int32(int32(mins.(int)))
	}

	if port, ok := config["outbound_ports_allocated"].(int); ok {
		profile.AllocatedOutboundPorts = utils.Int32(int32(port))
	}

	if ipCount := config["managed_outbound_ip_count"]; ipCount != nil {
		if c := int32(ipCount.(int)); c > 0 {
			profile.ManagedOutboundIPs = &containerservice.ManagedClusterLoadBalancerProfileManagedOutboundIPs{Count: &c}
		}
	}

	if ipPrefixes := idsToResourceReferences(config["outbound_ip_prefix_ids"]); ipPrefixes != nil {
		profile.OutboundIPPrefixes = &containerservice.ManagedClusterLoadBalancerProfileOutboundIPPrefixes{PublicIPPrefixes: ipPrefixes}
	}

	if outIps := idsToResourceReferences(config["outbound_ip_address_ids"]); outIps != nil {
		profile.OutboundIPs = &containerservice.ManagedClusterLoadBalancerProfileOutboundIPs{PublicIPs: outIps}
	}

	return profile
}

func idsToResourceReferences(set interface{}) *[]containerservice.ResourceReference {
	if set == nil {
		return nil
	}

	s := set.(*pluginsdk.Set)
	results := make([]containerservice.ResourceReference, 0)

	for _, element := range s.List() {
		id := element.(string)
		results = append(results, containerservice.ResourceReference{ID: &id})
	}

	if len(results) > 0 {
		return &results
	}

	return nil
}

func resourceReferencesToIds(refs *[]containerservice.ResourceReference) []string {
	if refs == nil {
		return nil
	}

	ids := make([]string, 0)

	for _, ref := range *refs {
		if ref.ID != nil {
			ids = append(ids, *ref.ID)
		}
	}

	if len(ids) > 0 {
		return ids
	}

	return nil
}

func flattenKubernetesClusterNetworkProfile(profile *containerservice.NetworkProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	dnsServiceIP := ""
	if profile.DNSServiceIP != nil {
		dnsServiceIP = *profile.DNSServiceIP
	}

	dockerBridgeCidr := ""
	if profile.DockerBridgeCidr != nil {
		dockerBridgeCidr = *profile.DockerBridgeCidr
	}

	serviceCidr := ""
	if profile.ServiceCidr != nil {
		serviceCidr = *profile.ServiceCidr
	}

	podCidr := ""
	if profile.PodCidr != nil {
		podCidr = *profile.PodCidr
	}

	lbProfiles := make([]interface{}, 0)
	if lbp := profile.LoadBalancerProfile; lbp != nil {
		lb := make(map[string]interface{})

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

		lb["effective_outbound_ips"] = resourceReferencesToIds(profile.LoadBalancerProfile.EffectiveOutboundIPs)
		lbProfiles = append(lbProfiles, lb)
	}

	return []interface{}{
		map[string]interface{}{
			"dns_service_ip":        dnsServiceIP,
			"docker_bridge_cidr":    dockerBridgeCidr,
			"load_balancer_sku":     string(profile.LoadBalancerSku),
			"load_balancer_profile": lbProfiles,
			"network_plugin":        string(profile.NetworkPlugin),
			"network_mode":          string(profile.NetworkMode),
			"network_policy":        string(profile.NetworkPolicy),
			"pod_cidr":              podCidr,
			"service_cidr":          serviceCidr,
			"outbound_type":         string(profile.OutboundType),
		},
	}
}

func expandKubernetesClusterRoleBasedAccessControl(input []interface{}, providerTenantId string) (bool, *containerservice.ManagedClusterAADProfile, error) {
	if len(input) == 0 {
		return false, nil, nil
	}

	val := input[0].(map[string]interface{})

	rbacEnabled := val["enabled"].(bool)
	azureADsRaw := val["azure_active_directory"].([]interface{})

	var aad *containerservice.ManagedClusterAADProfile

	if len(azureADsRaw) > 0 {
		azureAdRaw := azureADsRaw[0].(map[string]interface{})

		clientAppId := azureAdRaw["client_app_id"].(string)
		serverAppId := azureAdRaw["server_app_id"].(string)
		serverAppSecret := azureAdRaw["server_app_secret"].(string)
		tenantId := azureAdRaw["tenant_id"].(string)
		managed := azureAdRaw["managed"].(bool)
		azureRbacEnabled := azureAdRaw["azure_rbac_enabled"].(bool)
		adminGroupObjectIdsRaw := azureAdRaw["admin_group_object_ids"].(*pluginsdk.Set).List()
		adminGroupObjectIds := utils.ExpandStringSlice(adminGroupObjectIdsRaw)

		if tenantId == "" {
			tenantId = providerTenantId
		}

		if managed {
			aad = &containerservice.ManagedClusterAADProfile{
				TenantID:            utils.String(tenantId),
				Managed:             utils.Bool(managed),
				AdminGroupObjectIDs: adminGroupObjectIds,
				EnableAzureRBAC:     utils.Bool(azureRbacEnabled),
			}

			if clientAppId != "" || serverAppId != "" || serverAppSecret != "" {
				return false, nil, fmt.Errorf("Can't specify client_app_id or server_app_id or server_app_secret when using managed aad rbac (managed = true)")
			}
		} else {
			aad = &containerservice.ManagedClusterAADProfile{
				ClientAppID:     utils.String(clientAppId),
				ServerAppID:     utils.String(serverAppId),
				ServerAppSecret: utils.String(serverAppSecret),
				TenantID:        utils.String(tenantId),
				Managed:         utils.Bool(managed),
			}

			if len(*adminGroupObjectIds) > 0 {
				return false, nil, fmt.Errorf("Can't specify admin_group_object_ids when using managed aad rbac (managed = false)")
			}

			if clientAppId == "" || serverAppId == "" || serverAppSecret == "" {
				return false, nil, fmt.Errorf("You must specify client_app_id and server_app_id and server_app_secret when using managed aad rbac (managed = false)")
			}

			if azureRbacEnabled {
				return false, nil, fmt.Errorf("You must enable Managed AAD before Azure RBAC can be enabled")
			}
		}
	}

	return rbacEnabled, aad, nil
}

func expandKubernetesClusterManagedClusterIdentity(input []interface{}) *containerservice.ManagedClusterIdentity {
	if len(input) == 0 || input[0] == nil {
		return &containerservice.ManagedClusterIdentity{
			Type: containerservice.ResourceIdentityTypeNone,
		}
	}

	values := input[0].(map[string]interface{})

	if containerservice.ResourceIdentityType(values["type"].(string)) == containerservice.ResourceIdentityTypeUserAssigned {
		userAssignedIdentities := map[string]*containerservice.ManagedClusterIdentityUserAssignedIdentitiesValue{
			values["user_assigned_identity_id"].(string): {},
		}

		return &containerservice.ManagedClusterIdentity{
			Type:                   containerservice.ResourceIdentityType(values["type"].(string)),
			UserAssignedIdentities: userAssignedIdentities,
		}
	}

	return &containerservice.ManagedClusterIdentity{
		Type: containerservice.ResourceIdentityType(values["type"].(string)),
	}
}

func flattenKubernetesClusterRoleBasedAccessControl(input *containerservice.ManagedClusterProperties, d *pluginsdk.ResourceData) []interface{} {
	rbacEnabled := false
	if input.EnableRBAC != nil {
		rbacEnabled = *input.EnableRBAC
	}

	results := make([]interface{}, 0)
	if profile := input.AadProfile; profile != nil {
		adminGroupObjectIds := utils.FlattenStringSlice(profile.AdminGroupObjectIDs)

		clientAppId := ""
		if profile.ClientAppID != nil {
			clientAppId = *profile.ClientAppID
		}

		managed := false
		if profile.Managed != nil {
			managed = *profile.Managed
		}

		azureRbacEnabled := false
		if profile.EnableAzureRBAC != nil {
			azureRbacEnabled = *profile.EnableAzureRBAC
		}

		serverAppId := ""
		if profile.ServerAppID != nil {
			serverAppId = *profile.ServerAppID
		}

		serverAppSecret := ""
		// since input.ServerAppSecret isn't returned we're pulling this out of the existing state (which won't work for Imports)
		// role_based_access_control.0.azure_active_directory.0.server_app_secret
		if existing, ok := d.GetOk("role_based_access_control"); ok {
			rbacRawVals := existing.([]interface{})
			if len(rbacRawVals) > 0 {
				rbacRawVal := rbacRawVals[0].(map[string]interface{})
				if azureADVals, ok := rbacRawVal["azure_active_directory"].([]interface{}); ok && len(azureADVals) > 0 {
					azureADVal := azureADVals[0].(map[string]interface{})
					v := azureADVal["server_app_secret"]
					if v != nil {
						serverAppSecret = v.(string)
					}
				}
			}
		}

		tenantId := ""
		if profile.TenantID != nil {
			tenantId = *profile.TenantID
		}

		results = append(results, map[string]interface{}{
			"admin_group_object_ids": pluginsdk.NewSet(pluginsdk.HashString, adminGroupObjectIds),
			"client_app_id":          clientAppId,
			"managed":                managed,
			"server_app_id":          serverAppId,
			"server_app_secret":      serverAppSecret,
			"tenant_id":              tenantId,
			"azure_rbac_enabled":     azureRbacEnabled,
		})
	}

	return []interface{}{
		map[string]interface{}{
			"enabled":                rbacEnabled,
			"azure_active_directory": results,
		},
	}
}

func flattenAzureRmKubernetesClusterServicePrincipalProfile(profile *containerservice.ManagedClusterServicePrincipalProfile, d *pluginsdk.ResourceData) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	clientId := ""
	if v := profile.ClientID; v != nil {
		clientId = *v
	}

	if strings.EqualFold(clientId, "msi") {
		return []interface{}{}
	}

	// client secret isn't returned by the API so pass the existing value along
	clientSecret := ""
	if sp, ok := d.GetOk("service_principal"); ok {
		var val []interface{}

		// prior to 1.34 this was a *pluginsdk.Set, now it's a List - try both
		if v, ok := sp.([]interface{}); ok {
			val = v
		} else if v, ok := sp.(*pluginsdk.Set); ok {
			val = v.List()
		}

		if len(val) > 0 && val[0] != nil {
			raw := val[0].(map[string]interface{})
			clientSecret = raw["client_secret"].(string)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"client_id":     clientId,
			"client_secret": clientSecret,
		},
	}
}

func flattenKubernetesClusterKubeConfig(config kubernetes.KubeConfig) []interface{} {
	// we don't size-check these since they're validated in the Parse method
	cluster := config.Clusters[0].Cluster
	user := config.Users[0].User
	name := config.Users[0].Name

	return []interface{}{
		map[string]interface{}{
			"client_certificate":     user.ClientCertificteData,
			"client_key":             user.ClientKeyData,
			"cluster_ca_certificate": cluster.ClusterAuthorityData,
			"host":                   cluster.Server,
			"password":               user.Token,
			"username":               name,
		},
	}
}

func flattenKubernetesClusterKubeConfigAAD(config kubernetes.KubeConfigAAD) []interface{} {
	// we don't size-check these since they're validated in the Parse method
	cluster := config.Clusters[0].Cluster
	name := config.Users[0].Name

	return []interface{}{
		map[string]interface{}{
			"client_certificate":     "",
			"client_key":             "",
			"cluster_ca_certificate": cluster.ClusterAuthorityData,
			"host":                   cluster.Server,
			"password":               "",
			"username":               name,
		},
	}
}

func flattenKubernetesClusterManagedClusterIdentity(input *containerservice.ManagedClusterIdentity) ([]interface{}, error) {
	// if it's none, omit the block
	if input == nil || input.Type == containerservice.ResourceIdentityTypeNone {
		return []interface{}{}, nil
	}

	identity := make(map[string]interface{})

	identity["principal_id"] = ""
	if input.PrincipalID != nil {
		identity["principal_id"] = *input.PrincipalID
	}

	identity["tenant_id"] = ""
	if input.TenantID != nil {
		identity["tenant_id"] = *input.TenantID
	}

	identity["user_assigned_identity_id"] = ""
	if input.UserAssignedIdentities != nil {
		keys := []string{}
		for key := range input.UserAssignedIdentities {
			keys = append(keys, key)
		}
		if len(keys) > 0 {
			parsedId, err := msiparse.UserAssignedIdentityID(keys[0])
			if err != nil {
				return nil, err
			}
			identity["user_assigned_identity_id"] = parsedId.ID()
		}
	}

	identity["type"] = string(input.Type)

	return []interface{}{identity}, nil
}

func flattenKubernetesClusterAutoScalerProfile(profile *containerservice.ManagedClusterPropertiesAutoScalerProfile) ([]interface{}, error) {
	if profile == nil {
		return []interface{}{}, nil
	}

	balanceSimilarNodeGroups := false
	if profile.BalanceSimilarNodeGroups != nil {
		// @tombuildsstuff: presumably this'll get converted to a Boolean at some point
		//					at any rate we should use the proper type users expect here
		balanceSimilarNodeGroups = strings.EqualFold(*profile.BalanceSimilarNodeGroups, "true")
	}

	maxGracefulTerminationSec := ""
	if profile.MaxGracefulTerminationSec != nil {
		maxGracefulTerminationSec = *profile.MaxGracefulTerminationSec
	}

	maxNodeProvisionTime := ""
	if profile.MaxNodeProvisionTime != nil {
		maxNodeProvisionTime = *profile.MaxNodeProvisionTime
	}

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

	scanInterval := ""
	if profile.ScanInterval != nil {
		scanInterval = *profile.ScanInterval
	}

	skipNodesWithLocalStorage := true
	if profile.SkipNodesWithLocalStorage != nil {
		skipNodesWithLocalStorage = strings.EqualFold(*profile.SkipNodesWithLocalStorage, "true")
	}

	skipNodesWithSystemPods := true
	if profile.SkipNodesWithSystemPods != nil {
		skipNodesWithSystemPods = strings.EqualFold(*profile.SkipNodesWithSystemPods, "true")
	}

	return []interface{}{
		map[string]interface{}{
			"balance_similar_node_groups":      balanceSimilarNodeGroups,
			"expander":                         string(profile.Expander),
			"max_graceful_termination_sec":     maxGracefulTerminationSec,
			"max_node_provisioning_time":       maxNodeProvisionTime,
			"max_unready_nodes":                maxUnreadyNodes,
			"max_unready_percentage":           maxUnreadyPercentage,
			"new_pod_scale_up_delay":           newPodScaleUpDelay,
			"scale_down_delay_after_add":       scaleDownDelayAfterAdd,
			"scale_down_delay_after_delete":    scaleDownDelayAfterDelete,
			"scale_down_delay_after_failure":   scaleDownDelayAfterFailure,
			"scale_down_unneeded":              scaleDownUnneededTime,
			"scale_down_unready":               scaleDownUnreadyTime,
			"scale_down_utilization_threshold": scaleDownUtilizationThreshold,
			"empty_bulk_delete_max":            emptyBulkDeleteMax,
			"scan_interval":                    scanInterval,
			"skip_nodes_with_local_storage":    skipNodesWithLocalStorage,
			"skip_nodes_with_system_pods":      skipNodesWithSystemPods,
		},
	}, nil
}

func expandKubernetesClusterAutoScalerProfile(input []interface{}) *containerservice.ManagedClusterPropertiesAutoScalerProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})

	balanceSimilarNodeGroups := config["balance_similar_node_groups"].(bool)
	expander := config["expander"].(string)
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

	return &containerservice.ManagedClusterPropertiesAutoScalerProfile{
		BalanceSimilarNodeGroups:      utils.String(strconv.FormatBool(balanceSimilarNodeGroups)),
		Expander:                      containerservice.Expander(expander),
		MaxGracefulTerminationSec:     utils.String(maxGracefulTerminationSec),
		MaxNodeProvisionTime:          utils.String(maxNodeProvisionTime),
		MaxTotalUnreadyPercentage:     utils.String(maxUnreadyPercentage),
		NewPodScaleUpDelay:            utils.String(newPodScaleUpDelay),
		OkTotalUnreadyCount:           utils.String(maxUnreadyNodes),
		ScaleDownDelayAfterAdd:        utils.String(scaleDownDelayAfterAdd),
		ScaleDownDelayAfterDelete:     utils.String(scaleDownDelayAfterDelete),
		ScaleDownDelayAfterFailure:    utils.String(scaleDownDelayAfterFailure),
		ScaleDownUnneededTime:         utils.String(scaleDownUnneededTime),
		ScaleDownUnreadyTime:          utils.String(scaleDownUnreadyTime),
		ScaleDownUtilizationThreshold: utils.String(scaleDownUtilizationThreshold),
		MaxEmptyBulkDelete:            utils.String(emptyBulkDeleteMax),
		ScanInterval:                  utils.String(scanInterval),
		SkipNodesWithLocalStorage:     utils.String(strconv.FormatBool(skipNodesWithLocalStorage)),
		SkipNodesWithSystemPods:       utils.String(strconv.FormatBool(skipNodesWithSystemPods)),
	}
}
