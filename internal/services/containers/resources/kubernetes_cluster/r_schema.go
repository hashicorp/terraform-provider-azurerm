// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-07-01/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func kubernetesClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Managed Kubernetes Cluster to create.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"api_server_access_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authorized_ip_ranges": {
						Description: "Set of authorized IP ranges to allow access to API server, e.g. [`198.51.100.0/24`].",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDRIPv4,
						},
					},
					"virtual_network_integration_enabled": {
						Description: "Whether to enable Virtual Network integration for the API Server.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"subnet_id": {
						Description:  "The ID of the Subnet where the API server endpoint is delegated to.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"automatic_upgrade_channel": {
			Description: "The upgrade channel for this Kubernetes Cluster. Omitting this field sets this value to `none`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(managedclusters.UpgradeChannelPatch),
				string(managedclusters.UpgradeChannelRapid),
				string(managedclusters.UpgradeChannelStable),
				string(managedclusters.UpgradeChannelNodeNegativeimage),
			}, false),
		},

		"auto_scaler_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"balance_similar_node_groups": {
						Description: "Detect similar node groups and balance the number of nodes between them.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"daemonset_eviction_for_empty_nodes_enabled": {
						Description: "Whether DaemonSet pods will be gracefully terminated from empty nodes.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"daemonset_eviction_for_occupied_nodes_enabled": {
						Description: "Whether DaemonSet pods will be gracefully terminated from non-empty nodes.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"expander": {
						Description:  "Expander to use.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(managedclusters.ExpanderRandom),
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForExpander(), false),
					},

					"ignore_daemonsets_utilization_enabled": {
						Description: "Whether DaemonSet pods will be ignored when calculating resource utilization for scale down.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"max_graceful_termination_sec": {
						Description: "Maximum number of seconds the cluster autoscaler waits for pod termination when trying to scale down a node.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},

					"max_node_provisioning_time": {
						Description:  "Maximum time the autoscaler waits for a node to be provisioned.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "15m",
						ValidateFunc: validate.Duration,
					},

					"max_unready_nodes": {
						Description:  "Maximum Number of allowed unready nodes.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      3,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"max_unready_percentage": {
						Description:  "Maximum percentage of unready nodes the cluster autoscaler will stop if the percentage is exceeded.",
						Type:         pluginsdk.TypeFloat,
						Optional:     true,
						Default:      45,
						ValidateFunc: validation.FloatBetween(0, 100),
					},

					"new_pod_scale_up_delay": {
						Description:  "For scenarios like burst/batch scale where you don't want CA to act before the kubernetes scheduler could schedule all the pods, you can tell CA to ignore unscheduled pods before they're a certain age.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scan_interval": {
						Description:  "How often the AKS Cluster should be re-evaluated for scale up/down.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scale_down_delay_after_add": {
						Description:  "How long after the scale up of AKS nodes the scale down evaluation resumes.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scale_down_delay_after_delete": {
						Description:  "How long after node deletion that scale down evaluation resumes.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scale_down_delay_after_failure": {
						Description:  "How long after scale down failure that scale down evaluation resumes.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scale_down_unneeded": {
						Description:  "How long a node should be unneeded before it is eligible for scale down.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scale_down_unready": {
						Description:  "How long an unready node should be unneeded before it is eligible for scale down.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scale_down_utilization_threshold": {
						Description: "Node utilization level, defined as sum of requested resources divided by capacity, below which a node can be considered for scale down.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},

					"empty_bulk_delete_max": {
						Description: "Maximum number of empty nodes that can be deleted at the same time.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},

					"skip_nodes_with_local_storage": {
						Description: "If `true` cluster autoscaler will never delete nodes with pods with local storage, for example, EmptyDir or HostPath.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"skip_nodes_with_system_pods": {
						Description: "If `true` cluster autoscaler will never delete nodes with pods from kube-system (except for DaemonSet or mirror pods).",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
				},
			},
		},

		"azure_active_directory_role_based_access_control": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tenant_id": {
						Description: "The Tenant ID used for Azure Active Directory Application. If this isn't specified the Tenant ID of the current Subscription is used.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						// OrEmpty since this can be sourced from the client config if it's not specified
						ValidateFunc: validation.Any(validation.IsUUID, validation.StringIsEmpty),
						AtLeastOneOf: []string{
							"azure_active_directory_role_based_access_control.0.tenant_id",
							"azure_active_directory_role_based_access_control.0.admin_group_object_ids",
						},
					},

					"azure_rbac_enabled": {
						Description: "Whether Role Based Access Control based on Azure AD is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"admin_group_object_ids": {
						Description: "A list of Object IDs of Azure Active Directory Groups which should have Admin Role on the Cluster.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
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

		"cost_analysis_enabled": {
			Description: "Whether cost analysis should be enabled for this Kubernetes Cluster. The `sku_tier` must be set to `Standard` or `Premium` to enable this feature. Enabling this will add Kubernetes Namespace and Deployment details to the Cost Analysis views in the Azure portal.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"custom_ca_trust_certificates_base64": {
			Description: "A list of up to 10 base64 encoded CA certificates that will be added to the trust store on nodes.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    10,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsBase64,
			},
		},

		"default_node_pool": helpers.SchemaDefaultNodePool(),

		"disk_encryption_set_id": {
			Description:  "The ID of the Disk Encryption Set which should be used for the Nodes and Volumes. More information [can be found in the documentation](https://learn.microsoft.com/azure/aks/azure-disk-customer-managed-keys).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseDiskEncryptionSetIDInsensitively),
		},

		"dns_prefix": {
			Description:  "DNS prefix specified when creating the managed cluster.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"dns_prefix", "dns_prefix_private_cluster"},
			ValidateFunc: validate.KubernetesDNSPrefix,
		},

		"dns_prefix_private_cluster": {
			Description:  "The DNS prefix to use with private clusters.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"dns_prefix", "dns_prefix_private_cluster"},
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"fqdn": {
			Description: "The FQDN of the Azure Kubernetes Managed Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"http_proxy_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"http_proxy": {
						Description: "The proxy address to be used when communicating over HTTP.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"https_proxy": {
						Description: "The proxy address to be used when communicating over HTTPS.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"no_proxy": {
						Description: "The list of domains that will not use the proxy for communication.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"trusted_ca": {
						Description: "The base64 encoded alternative CA certificate content in PEM format.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"image_cleaner_enabled": {
			Description: "Whether Image Cleaner is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"image_cleaner_interval_hours": {
			Description:  "The interval in hours when images should be cleaned up.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(24, 2160),
		},

		"web_app_routing": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dns_zone_ids": {
						Description: "The list of the DNS Zone IDs in which DNS entries are created for applications deployed to the cluster when Web App Routing is enabled. If not using Bring-Your-Own DNS zones this property should be set to an empty list.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								zones.ValidateDnsZoneID,
								privatezones.ValidatePrivateDnsZoneID,
								validation.StringIsEmpty,
							),
						},
					},

					"default_nginx_controller": {
						Description: "The ingress type for the default `NginxIngressController` custom resource.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     managedclusters.NginxIngressControllerTypeAnnotationControlled,
						ValidateFunc: validation.StringInSlice(
							managedclusters.PossibleValuesForNginxIngressControllerType(), false,
						),
					},

					"web_app_routing_identity": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"client_id": {
									Description: "The Client ID of the user-defined Managed Identity used for Web App Routing.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"object_id": {
									Description: "The Object ID of the user-defined Managed Identity used for Web App Routing.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
								"user_assigned_identity_id": {
									Description: "The ID of the User Assigned Identity used for Web App Routing.",
									Type:        pluginsdk.TypeString,
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},

		"kube_admin_config": {
			Description: "This is only available when Role Based Access Control with Azure Active Directory is enabled and local accounts enabled.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Sensitive:   true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Description: "The Kubernetes cluster server host.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"username": {
						Description: "A username used to authenticate to the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"password": {
						Description: "A password or token used to authenticate to the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"client_certificate": {
						Description: "Base64 encoded public certificate used by clients to authenticate to the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"client_key": {
						Description: "Base64 encoded private key used by clients to authenticate to the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"cluster_ca_certificate": {
						Description: "Base64 encoded public CA certificate used as the root of trust for the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"kube_admin_config_raw": {
			Description: "Raw Kubernetes config for the admin account to be used by [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) and other compatible tools. This is only available when Role Based Access Control with Azure Active Directory is enabled and local accounts enabled.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"kube_config": {
			Type:      pluginsdk.TypeList,
			Computed:  true,
			Sensitive: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"host": {
						Description: "The Kubernetes cluster server host.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"username": {
						Description: "A username used to authenticate to the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"password": {
						Description: "A password or token used to authenticate to the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"client_certificate": {
						Description: "Base64 encoded public certificate used by clients to authenticate to the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"client_key": {
						Description: "Base64 encoded private key used by clients to authenticate to the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
					"cluster_ca_certificate": {
						Description: "Base64 encoded public CA certificate used as the root of trust for the Kubernetes cluster.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
		},

		"kube_config_raw": {
			Description: "Raw Kubernetes config to be used by [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) and other compatible tools.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"kubelet_identity": {
			Type:     pluginsdk.TypeList,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Description: "The Client ID of the user-defined Managed Identity to be assigned to the Kubelets. If not specified a Managed Identity is created automatically.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ForceNew:    true,
						RequiredWith: []string{
							"kubelet_identity.0.object_id",
							"kubelet_identity.0.user_assigned_identity_id",
							"identity.0.identity_ids",
						},
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"object_id": {
						Description: "The Object ID of the user-defined Managed Identity assigned to the Kubelets. If not specified a Managed Identity is created automatically.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ForceNew:    true,
						RequiredWith: []string{
							"kubelet_identity.0.client_id",
							"kubelet_identity.0.user_assigned_identity_id",
							"identity.0.identity_ids",
						},
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"user_assigned_identity_id": {
						Description: "The ID of the User Assigned Identity assigned to the Kubelets. If not specified a Managed Identity is created automatically.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ForceNew:    true,
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

		"kubernetes_version": {
			Description: "Version of Kubernetes specified when creating the AKS managed cluster. If not specified, the latest recommended version will be used at provisioning time (but won't auto-upgrade). AKS does not require an exact patch version to be specified, minor version aliases such as `1.22` are also supported. - The minor version's latest GA patch is automatically chosen in that case. More details can be found in [the documentation](https://learn.microsoft.com/azure/aks/supported-kubernetes-versions?tabs=azure-cli#alias-minor-version).",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because AKS assigns a Kubernetes version when not specified
			Computed:         true,
			ValidateFunc:     validation.StringIsNotEmpty,
			DiffSuppressFunc: helpers.SuppressKubernetesVersionDiff,
		},

		"linux_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Description:  "The Admin Username for the Cluster.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.KubernetesAdminUserName,
					},
					"ssh_key": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,

						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_data": {
									Description:  "The Public SSH Key used to access the cluster.",
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

		"bootstrap_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// Note: O+C because the API returns a default value for `bootstrapProfile` when it is omitted in the API request
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"artifact_source": {
						Description:  "The artifact source. The source where the artifacts are downloaded from.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForArtifactSource(), false),
						Default:      managedclusters.ArtifactSourceDirect,
					},

					"container_registry_id": {
						Description:  "The resource ID of Azure Container Registry.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: registries.ValidateRegistryID,
					},
				},
			},
		},

		"local_account_disabled": {
			Description: "If `true` local accounts will be disabled. See [the documentation](https://learn.microsoft.com/azure/aks/managed-aad#disable-local-accounts) for more information.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
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
									Description:  "A day in a week.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(maintenanceconfigurations.PossibleValuesForWeekDay(), false),
								},

								"hours": {
									Description: "An array of hour slots in a day. For example, specifying `1` will allow maintenance from 1:00am to 2:00am. Specifying `1`, `2` will allow maintenance from 1:00am to 3:00am.",
									Type:        pluginsdk.TypeSet,
									Required:    true,
									MinItems:    1,
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
									Description:      "The end of a time span, formatted as an RFC3339 string.",
									Type:             pluginsdk.TypeString,
									Required:         true,
									DiffSuppressFunc: suppress.RFC3339Time,
									ValidateFunc:     validation.IsRFC3339Time,
								},

								"start": {
									Description:      "The start of a time span, formatted as an RFC3339 string.",
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
						Description: "Frequency of maintenance.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"Daily",
							"Weekly",
							"RelativeMonthly",
							"AbsoluteMonthly",
						}, false),
					},

					"interval": {
						Description: "The interval for maintenance runs. Depending on the frequency this interval is week or month based.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},

					"day_of_week": {
						Description: "The day of the week for the maintenance run. Required in combination with weekly frequency.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice(
							maintenanceconfigurations.PossibleValuesForWeekDay(),
							false,
						),
					},

					"duration": {
						Description:  "The duration of the window for maintenance to run in hours.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(4, 24),
					},

					"week_index": {
						Description: "The instance of the allowed days specified in `day_of_week` on which the maintenance occurs. Required in combination with relative monthly frequency.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice(
							maintenanceconfigurations.PossibleValuesForType(),
							false,
						),
					},

					"day_of_month": {
						Description:  "The day of the month for the maintenance run. Required in combination with AbsoluteMonthly frequency. Value between 0 and 31 (inclusive).",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 31),
					},

					"start_date": {
						Description: "The date on which the maintenance window begins to take effect.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
					},

					"start_time": {
						Description: "The time for maintenance to begin, based on the timezone determined by `utc_offset`. Format is `HH:mm`.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},

					"utc_offset": {
						Description: "Used to determine the timezone for cluster maintenance.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},

					"not_allowed": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"end": {
									Description:      "The end of a time span, formatted as an RFC3339 string.",
									Type:             pluginsdk.TypeString,
									Required:         true,
									DiffSuppressFunc: suppress.RFC3339Time,
									ValidateFunc:     validation.IsRFC3339Time,
								},

								"start": {
									Description:      "The start of a time span, formatted as an RFC3339 string.",
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
						Description: "Frequency of maintenance.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"Weekly",
							"RelativeMonthly",
							"AbsoluteMonthly",
							"Daily",
						}, false),
					},

					"interval": {
						Description: "The interval for maintenance runs. Depending on the frequency this interval is week or month based.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},

					"day_of_week": {
						Description: "The day of the week for the maintenance run. Required in combination with weekly frequency.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice(
							maintenanceconfigurations.PossibleValuesForWeekDay(),
							false,
						),
					},

					"duration": {
						Description:  "The duration of the window for maintenance to run in hours.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(4, 24),
					},

					"week_index": {
						Description: "The week in the month used for the maintenance run.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						ValidateFunc: validation.StringInSlice(
							maintenanceconfigurations.PossibleValuesForType(),
							false,
						),
					},

					"day_of_month": {
						Description:  "The day of the month for the maintenance run. Required in combination with AbsoluteMonthly frequency. Value between 0 and 31 (inclusive).",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 31),
					},

					"start_date": {
						Description:      "The date on which the maintenance window begins to take effect.",
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Computed:         true, // azignore:AZS007 - pre-existing violation
						DiffSuppressFunc: suppress.RFC3339Time,
						ValidateFunc:     validation.IsRFC3339Time,
					},

					"start_time": {
						Description: "The time for maintenance to begin, based on the timezone determined by `utc_offset`. Format is `HH:mm`.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},

					"utc_offset": {
						Description: "Used to determine the timezone for cluster maintenance.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},

					"not_allowed": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"end": {
									Description:      "The end of a time span, formatted as an RFC3339 string.",
									Type:             pluginsdk.TypeString,
									Required:         true,
									DiffSuppressFunc: suppress.RFC3339Time,
									ValidateFunc:     validation.IsRFC3339Time,
								},

								"start": {
									Description:      "The start of a time span, formatted as an RFC3339 string.",
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

		"monitor_metrics": {
			Description: "A Prometheus add-on profile for the Kubernetes Cluster.",
			Type:        pluginsdk.TypeList,
			MaxItems:    1,
			Optional:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"annotations_allowed": {
						Description:  "A comma-separated list of Kubernetes annotation keys that will be used in the resource's labels metric.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"labels_allowed": {
						Description:  "A Comma-separated list of additional Kubernetes label keys that will be used in the resource's labels metric.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"key_management_service": {
			Description: "For more details, visit [Key Management Service (KMS) etcd encryption to an AKS cluster](https://learn.microsoft.com/azure/aks/use-kms-etcd-encryption).",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Description: "Identifier of Azure Key Vault Key. See [key identifier format](https://learn.microsoft.com/azure/key-vault/general/about-keys-secrets-certificates#vault-name-and-object-name) for more details. A versioned key ID uses the legacy KMS experience, while a versionless key ID uses the new KMS experience with automatic key rotation. For more details, visit [KMS data encryption concepts](https://learn.microsoft.com/azure/aks/kms-data-encryption-concepts).",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.Any(
							keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
						),
					},
					"key_vault_id": {
						Description:  "The ID of the Key Vault in which the key is located.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateKeyVaultID,
					},
					"key_vault_network_access": {
						Description:  "Network access of the Key Vault. `Public` means the Key Vault allows public access from all networks. `Private` means the Key Vault disables public access and enables Private Link.",
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
						Description:  "The ID of the Log Analytics Workspace where the audit logs collected by Microsoft Defender should be sent to.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: workspaces.ValidateWorkspaceID,
					},
				},
			},
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_plugin": {
						Description:  "Network plugin to use for networking.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForNetworkPlugin(), false),
					},

					"network_mode": {
						Description:  "Network mode to be used with Azure CNI.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForNetworkMode(), false),
					},

					"network_policy": {
						Description: "Sets up network policy to be used with Azure CNI. [Network policy allows us to control the traffic flow between pods](https://learn.microsoft.com/azure/aks/use-network-policies).",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.NetworkPolicyCalico),
							string(managedclusters.NetworkPolicyAzure),
							string(managedclusters.NetworkPolicyCilium),
						}, false),
					},

					"dns_service_ip": {
						Description:  "IP address within the Kubernetes service address range that will be used by cluster service discovery (kube-dns).",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ForceNew:     true,
						ValidateFunc: validation.IsIPv4Address,
					},

					"network_data_plane": {
						Description: "The data plane used for building the Kubernetes network. Disabling this forces a new resource to be created.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     string(managedclusters.NetworkDataplaneAzure),
						ValidateFunc: validation.StringInSlice(
							managedclusters.PossibleValuesForNetworkDataplane(),
							false,
						),
					},

					"network_plugin_mode": {
						Description:  "The network plugin mode used for building the Kubernetes network.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForNetworkPluginMode(), false),
					},

					"pod_cidr": {
						Description:  "The CIDR to use for pod IP addresses. This field can only be set when `network_plugin` is set to `kubenet` or `network_plugin_mode` is set to `overlay`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IsCIDRIPv4,
					},

					"pod_cidrs": {
						Description: "A list of CIDRs to use for pod IP addresses. For single-stack networking a single IPv4 CIDR is expected. For dual-stack networking an IPv4 and IPv6 CIDR are expected.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"service_cidr": {
						Description:  "The Network Range used by the Kubernetes service.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ForceNew:     true,
						ValidateFunc: validation.IsCIDRIPv4,
					},

					"service_cidrs": {
						Description: "A list of CIDRs to use for Kubernetes services. For single-stack networking a single IPv4 CIDR is expected. For dual-stack networking an IPv4 and IPv6 CIDR are expected.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ForceNew:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"load_balancer_sku": {
						Description:  "The SKU of the Load Balancer used for this Kubernetes Cluster.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(managedclusters.LoadBalancerSkuStandard),
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForLoadBalancerSku(), false),
					},

					"outbound_type": {
						Description:  "The outbound (egress) routing method which should be used for this Kubernetes Cluster.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(managedclusters.OutboundTypeLoadBalancer),
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForOutboundType(), false),
					},

					"load_balancer_profile": {
						Description: "This can only be specified when `load_balancer_sku` is set to `standard`.",
						Type:        pluginsdk.TypeList,
						MaxItems:    1,
						ForceNew:    true,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"outbound_ports_allocated": {
									Description:  "Number of desired SNAT port for each VM in the clusters load balancer.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntBetween(0, 64000),
								},

								"idle_timeout_in_minutes": {
									Description:  "Desired outbound flow idle timeout in minutes for the cluster load balancer.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      30,
									ValidateFunc: validation.IntBetween(4, 100),
								},

								"managed_outbound_ip_count": {
									Description:   "Count of desired managed outbound IPs for the cluster load balancer.",
									Type:          pluginsdk.TypeInt,
									Optional:      true,
									Computed:      true, // azignore:AZS007 - pre-existing violation
									ValidateFunc:  validation.IntBetween(1, 100),
									ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids", "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"},
								},

								"managed_outbound_ipv6_count": {
									Description:   "The desired number of IPv6 outbound IPs created and managed by Azure for the cluster load balancer.",
									Type:          pluginsdk.TypeInt,
									Optional:      true,
									Computed:      true, // azignore:AZS007 - pre-existing violation
									ValidateFunc:  validation.IntBetween(1, 100),
									ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids", "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"},
								},

								"outbound_ip_prefix_ids": {
									Description:   "The ID of the outbound Public IP Address Prefixes which should be used for the cluster load balancer.",
									Type:          pluginsdk.TypeSet,
									Optional:      true,
									ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.managed_outbound_ip_count", "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"},
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: azure.ValidateResourceID,
									},
								},

								"outbound_ip_address_ids": {
									Description:   "The ID of the Public IP Addresses which should be used for outbound communication for the cluster load balancer.",
									Type:          pluginsdk.TypeSet,
									Optional:      true,
									ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.managed_outbound_ip_count", "network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids"},
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: azure.ValidateResourceID,
									},
								},

								"effective_outbound_ips": {
									Description: "The outcome (resource IDs) of the specified arguments.",
									Type:        pluginsdk.TypeSet,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},

								"backend_pool_type": {
									Description:  "The type of the managed inbound Load Balancer Backend Pool. See [the documentation](https://learn.microsoft.com/azure/aks/load-balancer-standard#change-the-inbound-pool-type) for more information.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(managedclusters.BackendPoolTypeNodeIPConfiguration),
									ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForBackendPoolType(), false),
								},
							},
						},
					},

					"nat_gateway_profile": {
						Description: "This can only be specified when `load_balancer_sku` is set to `standard` and `outbound_type` is set to `managedNATGateway` or `userAssignedNATGateway`.",
						Type:        pluginsdk.TypeList,
						MaxItems:    1,
						ForceNew:    true,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"idle_timeout_in_minutes": {
									Description:  "Desired outbound flow idle timeout in minutes for the managed NAT gateway.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      4,
									ValidateFunc: validation.IntBetween(4, 120),
								},
								"managed_outbound_ip_count": {
									Description:  "Count of desired managed outbound IPs for the managed NAT gateway.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Computed:     true, // azignore:AZS007 - pre-existing violation
									ValidateFunc: validation.IntBetween(1, 100),
								},
								"effective_outbound_ips": {
									Description: "The outcome (resource IDs) of the specified arguments.",
									Type:        pluginsdk.TypeSet,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"ip_versions": {
						Description: "A list of IP versions the Kubernetes Cluster will use to assign IP addresses to its nodes and pods. `IPv4` must always be specified.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						ForceNew:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForIPFamily(), false),
						},
					},

					"advanced_networking": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"observability_enabled": {
									Description:  "Whether observability is enabled.",
									Type:         pluginsdk.TypeBool,
									Optional:     true,
									Default:      false,
									AtLeastOneOf: []string{"network_profile.0.advanced_networking.0.observability_enabled", "network_profile.0.advanced_networking.0.security_enabled"},
								},
								"security_enabled": {
									Description:  "Whether security is enabled. This can only be enabled (set to `true`) when `network_plugin` is set to `azure` and `network_data_plane` is set to `cilium`.",
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

		"node_os_upgrade_channel": {
			Description:  "The upgrade channel for this Kubernetes Cluster Nodes' OS Image.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(managedclusters.NodeOSUpgradeChannelNodeImage),
			ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForNodeOSUpgradeChannel(), false),
		},

		"node_resource_group": {
			Description: "The name of the Resource Group where the Kubernetes Nodes should exist.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
		},

		"current_kubernetes_version": {
			Description: "The current version running on the Azure Kubernetes Managed Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"node_resource_group_id": {
			Description: "The ID of the Resource Group containing the resources for this Managed Kubernetes Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"oidc_issuer_enabled": {
			Description: "Whether to enable the [OIDC issuer feature](https://learn.microsoft.com/azure/aks/use-oidc-issuer).",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"oidc_issuer_url": {
			Description: "The OIDC issuer URL that is associated with the cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"private_fqdn": {
			Description: "The FQDN for the Kubernetes Cluster when Private Link has been enabled, which is only resolvable inside the Virtual Network used by the Kubernetes Cluster.",
			// privateFqdn
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"portal_fqdn": {
			Description: "The FQDN for the Azure Portal resources when Private Link has been enabled, which is only resolvable inside the Virtual Network used by the Kubernetes Cluster.",
			// azurePortalFqdn
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_cluster_enabled": {
			Description: "Whether this Kubernetes Cluster should have its API server only exposed on internal IP addresses. This provides a Private IP Address for the Kubernetes API on the Virtual Network where the Kubernetes Cluster is located.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"private_cluster_public_fqdn_enabled": {
			Description: "Whether a Public FQDN for this Private Cluster should be added.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"private_dns_zone_id": {
			Description: "Either the ID of Private DNS Zone which should be delegated to this Cluster, `System` to have AKS manage this or `None`. In case of `None` you will need to bring your own DNS server and set up resolving, otherwise, the cluster will have issues after provisioning.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because a Private Cluster is `System` by default even if unspecified
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

		"role_based_access_control_enabled": {
			Description: "Whether Role Based Access Control for the Kubernetes Cluster should be enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			ForceNew:    true,
		},

		"run_command_enabled": {
			Description: "Whether to enable run command for the cluster.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"service_mesh_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Description: "The mode of the service mesh.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.ServiceMeshModeIstio),
						}, false),
					},
					"internal_ingress_gateway_enabled": {
						Description: "Whether Istio Internal Ingress Gateway is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"external_ingress_gateway_enabled": {
						Description: "Whether Istio External Ingress Gateway is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"certificate_authority": {
						Description: "When this property is specified, `key_vault_secrets_provider` is also required to be set. This configuration allows you to bring your own root certificate and keys for Istio CA in the Istio-based service mesh add-on for Azure Kubernetes Service.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						MaxItems:    1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"key_vault_id": commonschema.ResourceIDReferenceRequired(&commonids.KeyVaultId{}, "Key Vault"),
								"root_cert_object_name": {
									Description:  "The root certificate object name in Azure Key Vault.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"cert_chain_object_name": {
									Description:  "The certificate chain object name in Azure Key Vault.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"cert_object_name": {
									Description:  "The intermediate certificate object name in Azure Key Vault.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"key_object_name": {
									Description:  "The intermediate certificate private key object name in Azure Key Vault.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"revisions": {
						Description: "One to two Istio control plane revisions for managing minor upgrades using the canary upgrade process. For example, create the resource with `revisions` set to `[\"asm-1-27\"]`. To start the canary upgrade, change `revisions` to `[\"asm-1-27\", \"asm-1-28\"]`. To roll back the canary upgrade, revert to `[\"asm-1-27\"]`. To confirm the upgrade, change to `[\"asm-1-28\"]`.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						MinItems:    1,
						MaxItems:    2,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
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
						Description:  "The Client ID for the Service Principal.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ClientID,
					},

					"client_secret": {
						Description:  "The Client Secret for the Service Principal.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"sku_tier": {
			Description:  "The SKU Tier that should be used for this Kubernetes Cluster.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(managedclusters.ManagedClusterSKUTierFree),
			ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForManagedClusterSKUTier(), false),
		},

		"storage_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"blob_driver_enabled": {
						Description: "Whether the Blob CSI driver is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"disk_driver_enabled": {
						Description: "Whether the Disk CSI driver is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"file_driver_enabled": {
						Description: "Whether the File CSI driver is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"snapshot_controller_enabled": {
						Description: "Whether the Snapshot Controller is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
				},
			},
		},

		"support_plan": {
			Description:  "The support plan which should be used for this Kubernetes Cluster.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(managedclusters.KubernetesSupportPlanKubernetesOfficial),
			ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForKubernetesSupportPlan(), false),
		},

		"tags": commonschema.Tags(),

		"upgrade_override": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"force_upgrade_enabled": {
						Description: "Whether to force upgrade the cluster.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},

					"effective_until": {
						Description:  "The duration, in RFC 3339 format (e.g., `2025-10-01T13:00:00Z`), the `upgrade_override` values are effective. This field must be set for the `upgrade_override` values to take effect. The date-time must be within the next 30 days.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
				},
			},
		},

		"windows_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_username": {
						Description: "The Admin Username for Windows VMs.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ForceNew:    true,
					},
					"admin_password": {
						Description:  "The Admin Password for Windows VMs.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringLenBetween(8, 123),
					},
					"license": {
						Description: "The type of on-premise license which should be used for Node Pool Windows Virtual Machine.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
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
								"dns_server": {
									Description: "The DNS server for Windows gMSA. Set this to an empty string if you have configured the DNS server in the VNet which was used to create the managed cluster.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"root_domain": {
									Description: "The root domain name for Windows gMSA. Set this to an empty string if you have configured the DNS server in the VNet which was used to create the managed cluster.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
							},
						},
					},
				},
			},
		},

		"workload_autoscaler_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"keda_enabled": {
						Description: "Whether KEDA Autoscaler can be used for workloads.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"vertical_pod_autoscaler_enabled": {
						Description: "Whether Vertical Pod Autoscaler should be enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"ai_toolchain_operator_enabled": {
			Description: "Whether the AI Toolchain Operator should be enabled for the Cluster.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"node_provisioning_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      managedclusters.NodeProvisioningModeManual,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForNodeProvisioningMode(), false),
						AtLeastOneOf: []string{"node_provisioning_profile.0.mode", "node_provisioning_profile.0.default_node_pools"},
					},
					"default_node_pools": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      managedclusters.NodeProvisioningDefaultNodePoolsAuto,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForNodeProvisioningDefaultNodePools(), false),
						AtLeastOneOf: []string{"node_provisioning_profile.0.mode", "node_provisioning_profile.0.default_node_pools"},
					},
				},
			},
		},

		"workload_identity_enabled": {
			Description: "Whether Azure AD Workload Identity should be enabled for the Cluster.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},
	}
}
