// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kubernetes_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2025-11-01/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2026-05-01/managedclusters"
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
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDRIPv4,
						},
					},
					"virtual_network_integration_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"automatic_upgrade_channel": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(managedclusters.ExpanderRandom),
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForExpander(), false),
					},

					"ignore_daemonsets_utilization_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"max_graceful_termination_sec": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
					},

					"max_node_provisioning_time": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "15m",
						ValidateFunc: validate.Duration,
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
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scan_interval": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scale_down_delay_after_add": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scale_down_delay_after_delete": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scale_down_delay_after_failure": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scale_down_unneeded": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scale_down_unready": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.Duration,
					},

					"scale_down_utilization_threshold": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
					},

					"empty_bulk_delete_max": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
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

		"azure_active_directory_role_based_access_control": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						// OrEmpty since this can be sourced from the client config if it's not specified
						ValidateFunc: validation.Any(validation.IsUUID, validation.StringIsEmpty),
						AtLeastOneOf: []string{
							"azure_active_directory_role_based_access_control.0.tenant_id",
							"azure_active_directory_role_based_access_control.0.admin_group_object_ids",
						},
					},

					"azure_rbac_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"admin_group_object_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
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
			Type:     pluginsdk.TypeBool,
			Optional: true,
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

		"default_node_pool": helpers.SchemaDefaultNodePool(),

		"disk_encryption_set_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseDiskEncryptionSetIDInsensitively),
		},

		"dns_prefix": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"dns_prefix", "dns_prefix_private_cluster"},
			ValidateFunc: validate.KubernetesDNSPrefix,
		},

		"dns_prefix_private_cluster": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ExactlyOneOf: []string{"dns_prefix", "dns_prefix_private_cluster"},
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
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

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"image_cleaner_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"image_cleaner_interval_hours": {
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
						Type:     pluginsdk.TypeList,
						Required: true,
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
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  managedclusters.NginxIngressControllerTypeAnnotationControlled,
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
			Computed: true, // azignore:AZS007 - pre-existing violation
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
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
						Computed: true, // azignore:AZS007 - pre-existing violation
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
						Computed: true, // azignore:AZS007 - pre-existing violation
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

		"kubernetes_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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

		"local_account_disabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(maintenanceconfigurations.PossibleValuesForWeekDay(), false),
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
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice(
							maintenanceconfigurations.PossibleValuesForWeekDay(),
							false,
						),
					},

					"duration": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(4, 24),
					},

					"week_index": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice(
							maintenanceconfigurations.PossibleValuesForType(),
							false,
						),
					},

					"day_of_month": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 31),
					},

					"start_date": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
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
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice(
							maintenanceconfigurations.PossibleValuesForWeekDay(),
							false,
						),
					},

					"duration": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(4, 24),
					},

					"week_index": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice(
							maintenanceconfigurations.PossibleValuesForType(),
							false,
						),
					},

					"day_of_month": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 31),
					},

					"start_date": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Computed:         true, // azignore:AZS007 - pre-existing violation
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

		"monitor_metrics": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
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

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"network_plugin": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForNetworkPlugin(), false),
					},

					"network_mode": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForNetworkMode(), false),
					},

					"network_policy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice([]string{
							string(managedclusters.NetworkPolicyCalico),
							string(managedclusters.NetworkPolicyAzure),
							string(managedclusters.NetworkPolicyCilium),
						}, false),
					},

					"dns_service_ip": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ForceNew:     true,
						ValidateFunc: validation.IsIPv4Address,
					},

					"network_data_plane": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(managedclusters.NetworkDataplaneAzure),
						ValidateFunc: validation.StringInSlice(
							managedclusters.PossibleValuesForNetworkDataplane(),
							false,
						),
					},

					"network_plugin_mode": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForNetworkPluginMode(), false),
					},

					"pod_cidr": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IsCIDRIPv4,
					},

					"pod_cidrs": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"service_cidr": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ForceNew:     true,
						ValidateFunc: validation.IsCIDRIPv4,
					},

					"service_cidrs": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"load_balancer_sku": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(managedclusters.LoadBalancerSkuStandard),
						ForceNew:     true,
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForLoadBalancerSku(), false),
					},

					"outbound_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(managedclusters.OutboundTypeLoadBalancer),
						ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForOutboundType(), false),
					},

					"load_balancer_profile": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						ForceNew: true,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
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
									Computed:      true, // azignore:AZS007 - pre-existing violation
									ValidateFunc:  validation.IntBetween(1, 100),
									ConflictsWith: []string{"network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids", "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"},
								},

								"managed_outbound_ipv6_count": {
									Type:          pluginsdk.TypeInt,
									Optional:      true,
									Computed:      true, // azignore:AZS007 - pre-existing violation
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
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      string(managedclusters.BackendPoolTypeNodeIPConfiguration),
									ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForBackendPoolType(), false),
								},
							},
						},
					},

					"nat_gateway_profile": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						ForceNew: true,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
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
									Computed:     true, // azignore:AZS007 - pre-existing violation
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
						Computed: true, // azignore:AZS007 - pre-existing violation
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

		"node_os_upgrade_channel": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(managedclusters.NodeOSUpgradeChannelNodeImage),
			ValidateFunc: validation.StringInSlice(managedclusters.PossibleValuesForNodeOSUpgradeChannel(), false),
		},

		"node_resource_group": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
		},

		"current_kubernetes_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"node_resource_group_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"oidc_issuer_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"oidc_issuer_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_fqdn": { // privateFqdn
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"portal_fqdn": { // azurePortalFqdn
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"private_cluster_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"private_cluster_public_fqdn_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"private_dns_zone_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},

		"run_command_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
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
						ValidateFunc: validate.ClientID,
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

		"support_plan": {
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

		"windows_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
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

		"workload_autoscaler_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"keda_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"vertical_pod_autoscaler_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"ai_toolchain_operator_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
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
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}
