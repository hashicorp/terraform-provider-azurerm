// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-04-02-preview/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-04-02-preview/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-04-02-preview/managedclusters"
	dnsValidate "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2020-06-01/privatezones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/migration"
	containerValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	keyVaultClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	resourcesClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKubernetesCluster() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceKubernetesClusterCreate,
		Read:   resourceKubernetesClusterRead,
		Update: resourceKubernetesClusterUpdate,
		Delete: resourceKubernetesClusterDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseKubernetesClusterID(id)
			return err
		}),

		CustomizeDiff: pluginsdk.CustomDiffInSequence(
			// Migration of `identity` to `service_principal` is not allowed, the other way around is
			pluginsdk.ForceNewIfChange("service_principal.0.client_id", func(ctx context.Context, old, new, meta interface{}) bool {
				return old == "msi" || old == ""
			}),
			pluginsdk.ForceNewIfChange("windows_profile.0.gmsa", func(ctx context.Context, old, new, meta interface{}) bool {
				return len(old.([]interface{})) != 0 && len(new.([]interface{})) == 0
			}),
			pluginsdk.ForceNewIfChange("windows_profile.0.gmsa.0.dns_server", func(ctx context.Context, old, new, meta interface{}) bool {
				return old != "" && new == ""
			}),
			pluginsdk.ForceNewIfChange("windows_profile.0.gmsa.0.root_domain", func(ctx context.Context, old, new, meta interface{}) bool {
				return old != "" && new == ""
			}),
			pluginsdk.ForceNewIfChange("api_server_access_profile.0.subnet_id", func(ctx context.Context, old, new, meta interface{}) bool {
				return old != "" && new == ""
			}),
			pluginsdk.ForceNewIf("default_node_pool.0.name", func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) bool {
				old, new := d.GetChange("default_node_pool.0.name")
				defaultName := d.Get("default_node_pool.0.name")
				tempName := d.Get("default_node_pool.0.temporary_name_for_rotation")

				// if the default node pool name has been set to temporary_name_for_rotation it means resizing failed
				// we should not try to recreate the cluster, another apply will attempt the resize again
				if old != "" && old == tempName {
					return new != defaultName
				}
				return true
			}),
			pluginsdk.ForceNewIfChange("network_profile.0.ebpf_data_plane", func(ctx context.Context, old, new, meta interface{}) bool {
				return old != ""
			}),
			func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
				if d.HasChange("oidc_issuer_enabled") {
					d.SetNewComputed("oidc_issuer_url")
				}
				return nil
			},
		),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KubernetesClusterV0ToV1{},
			1: migration.KubernetesClusterV1ToV2{},
		}),

		Schema: map[string]*pluginsdk.Schema{
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
				Computed: !features.FourPointOhBeta(),
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"vnet_integration_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: commonids.ValidateSubnetID,
						},

						"authorized_ip_ranges": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Computed: !features.FourPointOhBeta(),
							ConflictsWith: func() []string {
								if !features.FourPointOhBeta() {
									return []string{"api_server_authorized_ip_ranges"}
								}
								return []string{}
							}(),
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.CIDR,
							},
						},
					},
				},
			},

			"automatic_channel_upgrade": {
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
							Default:  string(managedclusters.ExpanderRandom),
							ValidateFunc: validation.StringInSlice([]string{
								string(managedclusters.ExpanderLeastNegativewaste),
								string(managedclusters.ExpanderMostNegativepods),
								string(managedclusters.ExpanderPriority),
								string(managedclusters.ExpanderRandom),
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
							Default:  !features.FourPointOhBeta(),
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
						"client_app_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
							AtLeastOneOf: []string{
								"azure_active_directory_role_based_access_control.0.client_app_id", "azure_active_directory_role_based_access_control.0.server_app_id",
								"azure_active_directory_role_based_access_control.0.server_app_secret", "azure_active_directory_role_based_access_control.0.tenant_id",
								"azure_active_directory_role_based_access_control.0.managed", "azure_active_directory_role_based_access_control.0.admin_group_object_ids",
							},
						},

						"server_app_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
							AtLeastOneOf: []string{
								"azure_active_directory_role_based_access_control.0.client_app_id", "azure_active_directory_role_based_access_control.0.server_app_id",
								"azure_active_directory_role_based_access_control.0.server_app_secret", "azure_active_directory_role_based_access_control.0.tenant_id",
								"azure_active_directory_role_based_access_control.0.managed", "azure_active_directory_role_based_access_control.0.admin_group_object_ids",
							},
						},

						"server_app_secret": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{
								"azure_active_directory_role_based_access_control.0.client_app_id", "azure_active_directory_role_based_access_control.0.server_app_id",
								"azure_active_directory_role_based_access_control.0.server_app_secret", "azure_active_directory_role_based_access_control.0.tenant_id",
								"azure_active_directory_role_based_access_control.0.managed", "azure_active_directory_role_based_access_control.0.admin_group_object_ids",
							},
						},

						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							// OrEmpty since this can be sourced from the client config if it's not specified
							ValidateFunc: validation.Any(validation.IsUUID, validation.StringIsEmpty),
							AtLeastOneOf: []string{
								"azure_active_directory_role_based_access_control.0.client_app_id", "azure_active_directory_role_based_access_control.0.server_app_id",
								"azure_active_directory_role_based_access_control.0.server_app_secret", "azure_active_directory_role_based_access_control.0.tenant_id",
								"azure_active_directory_role_based_access_control.0.managed", "azure_active_directory_role_based_access_control.0.admin_group_object_ids",
							},
						},

						"managed": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							AtLeastOneOf: []string{
								"azure_active_directory_role_based_access_control.0.client_app_id", "azure_active_directory_role_based_access_control.0.server_app_id",
								"azure_active_directory_role_based_access_control.0.server_app_secret", "azure_active_directory_role_based_access_control.0.tenant_id",
								"azure_active_directory_role_based_access_control.0.managed", "azure_active_directory_role_based_access_control.0.admin_group_object_ids",
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
								"azure_active_directory_role_based_access_control.0.client_app_id", "azure_active_directory_role_based_access_control.0.server_app_id",
								"azure_active_directory_role_based_access_control.0.server_app_secret", "azure_active_directory_role_based_access_control.0.tenant_id",
								"azure_active_directory_role_based_access_control.0.managed", "azure_active_directory_role_based_access_control.0.admin_group_object_ids",
							},
						},
					},
				},
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

			"default_node_pool": SchemaDefaultNodePool(),

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

			"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

			"enable_pod_security_policy": {
				Type:       pluginsdk.TypeBool,
				Deprecated: "The AKS API has removed support for this field on 2020-10-15 and is no longer possible to configure this the Pod Security Policy.",
				Optional:   true,
			},

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
							ForceNew: true,
						},
						"https_proxy": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
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
				Default:  false,
			},

			"image_cleaner_interval_hours": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      48,
				ValidateFunc: validation.IntBetween(24, 2160),
			},

			"web_app_routing": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dns_zone_id": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.Any(
								dnsValidate.ValidateDnsZoneID,
								validation.StringIsEmpty,
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

			"kubernetes_version": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
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
								false),
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
								false),
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
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice(
								maintenanceconfigurations.PossibleValuesForWeekDay(),
								false),
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
								false),
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

			"node_os_channel_upgrade": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(managedclusters.NodeOSUpgradeChannelNodeImage),
					string(managedclusters.NodeOSUpgradeChannelNone),
					string(managedclusters.NodeOSUpgradeChannelSecurityPatch),
					string(managedclusters.NodeOSUpgradeChannelUnmanaged),
				}, false),
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
							ValidateFunc: keyVaultValidate.NestedItemId,
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
								string(managedclusters.NetworkPluginAzure),
								string(managedclusters.NetworkPluginKubenet),
								string(managedclusters.NetworkPluginNone),
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
								string(managedclusters.NetworkModeBridge),
								string(managedclusters.NetworkModeTransparent),
							}, false),
						},

						"network_policy": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(managedclusters.NetworkPolicyCalico),
								string(managedclusters.NetworkPolicyAzure),
							}, false),
						},

						"dns_service_ip": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validate.IPv4Address,
						},

						"ebpf_data_plane": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(managedclusters.NetworkDataplaneCilium),
							}, false),
						},

						"network_plugin_mode": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
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
							ForceNew: true,
							Default:  string(managedclusters.OutboundTypeLoadBalancer),
							ValidateFunc: validation.StringInSlice([]string{
								string(managedclusters.OutboundTypeLoadBalancer),
								string(managedclusters.OutboundTypeUserDefinedRouting),
								string(managedclusters.OutboundTypeManagedNATGateway),
								string(managedclusters.OutboundTypeUserAssignedNATGateway),
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
					},
				},
			},

			"node_resource_group": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"node_resource_group_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"oidc_issuer_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
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
				Computed: true, // a Private Cluster is `System` by default even if unspecified
				ForceNew: true,
				ValidateFunc: validation.Any(
					privatezones.ValidatePrivateDnsZoneID,
					validation.StringInSlice([]string{
						"System",
						"None",
					}, false),
				),
			},

			"public_network_access_enabled": {
				Type:       pluginsdk.TypeBool,
				Optional:   true,
				Default:    true,
				Deprecated: "`public_network_access_enabled` is currently not functional and is not be passed to the API",
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
				Default:  string(managedclusters.ManagedClusterSKUTierFree),
				ValidateFunc: validation.StringInSlice([]string{
					string(managedclusters.ManagedClusterSKUTierFree),
					string(managedclusters.ManagedClusterSKUTierStandard),
				}, false),
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
						"disk_driver_version": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "v1",
							ValidateFunc: validation.StringInSlice([]string{
								"v1",
								"v2",
							}, false),
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

			"tags": commonschema.Tags(),

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
						// This needs to become Required in 4.0 - omitting it isn't accepted by the API
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
						"vertical_pod_autoscaler_update_mode": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"vertical_pod_autoscaler_controlled_values": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"workload_identity_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}

	for k, v := range schemaKubernetesAddOns() {
		resource.Schema[k] = v
	}

	if !features.FourPointOhBeta() {
		resource.Schema["api_server_authorized_ip_ranges"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.CIDR,
			},
			Deprecated:    "This property has been renamed to `authorized_ip_ranges` within the `api_server_access_profile` block and will be removed in v4.0 of the provider",
			ConflictsWith: []string{"api_server_access_profile.0.authorized_ip_ranges"},
		}
		resource.Schema["network_profile"].Elem.(*pluginsdk.Resource).Schema["docker_bridge_cidr"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			Deprecated:   "`docker_bridge_cidr` has been deprecated as the API no longer supports it and will be removed in version 4.0 of the provider.",
			ValidateFunc: validate.CIDR,
		}
		resource.Schema["network_profile"].Elem.(*pluginsdk.Resource).Schema["network_plugin_mode"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(managedclusters.NetworkPluginModeOverlay),
				"Overlay",
			}, false),
		}
		resource.Schema["windows_profile"].Elem.(*pluginsdk.Resource).Schema["admin_password"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringLenBetween(8, 123),
		}
	}

	return resource
}

func resourceKubernetesClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	tenantId := meta.(*clients.Client).Account.TenantId
	client := meta.(*clients.Client).Containers.KubernetesClustersClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
	env := meta.(*clients.Client).Containers.Environment
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Managed Kubernetes Cluster create.")

	id := commonids.NewKubernetesClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_kubernetes_cluster", id.ID())
	}

	if err := validateKubernetesCluster(d, nil, id.ResourceGroupName, id.ManagedClusterName); err != nil {
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

	// the AKS API will create the default node pool with the same version as the control plane regardless of what is
	// supplied by the user which will result in a diff in some cases, so if versions have been supplied check that they
	// are identical
	agentProfile := ConvertDefaultNodePoolToAgentPool(agentProfiles)
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

	networkProfileRaw := d.Get("network_profile").([]interface{})
	networkProfile, err := expandKubernetesClusterNetworkProfile(networkProfileRaw)
	if err != nil {
		return err
	}

	var azureADProfile *managedclusters.ManagedClusterAADProfile
	if v, ok := d.GetOk("azure_active_directory_role_based_access_control"); ok {
		azureADProfile, err = expandKubernetesClusterAzureActiveDirectoryRoleBasedAccessControl(v.([]interface{}), tenantId)
		if err != nil {
			return err
		}
	}

	t := d.Get("tags").(map[string]interface{})

	windowsProfileRaw := d.Get("windows_profile").([]interface{})
	windowsProfile := expandKubernetesClusterWindowsProfile(windowsProfileRaw)

	workloadAutoscalerProfileRaw := d.Get("workload_autoscaler_profile").([]interface{})
	workloadAutoscalerProfile := expandKubernetesClusterWorkloadAutoscalerProfile(workloadAutoscalerProfileRaw, d)

	apiAccessProfile := expandKubernetesClusterAPIAccessProfile(d)
	if !(*apiAccessProfile.EnablePrivateCluster) && dnsPrefix == "" {
		return fmt.Errorf("`dns_prefix` should be set if it is not a private cluster")
	}

	nodeResourceGroup := d.Get("node_resource_group").(string)

	if d.Get("enable_pod_security_policy").(bool) {
		return fmt.Errorf("the AKS API has removed support for this field on 2020-10-15 and is no longer possible to configure this the Pod Security Policy - as such you'll need to set `enable_pod_security_policy` to `false`")
	}

	autoScalerProfileRaw := d.Get("auto_scaler_profile").([]interface{})
	autoScalerProfile := expandKubernetesClusterAutoScalerProfile(autoScalerProfileRaw)

	azureMonitorKubernetesMetricsRaw := d.Get("monitor_metrics").([]interface{})
	azureMonitorProfile := expandKubernetesClusterAzureMonitorProfile(azureMonitorKubernetesMetricsRaw)

	httpProxyConfigRaw := d.Get("http_proxy_config").([]interface{})
	httpProxyConfig := expandKubernetesClusterHttpProxyConfig(httpProxyConfigRaw)

	enableOidcIssuer := false
	var oidcIssuerProfile *managedclusters.ManagedClusterOIDCIssuerProfile
	if v, ok := d.GetOk("oidc_issuer_enabled"); ok {
		enableOidcIssuer = v.(bool)
		oidcIssuerProfile = expandKubernetesClusterOidcIssuerProfile(enableOidcIssuer)
	}

	storageProfileRaw := d.Get("storage_profile").([]interface{})
	storageProfile := expandStorageProfile(storageProfileRaw)

	// assemble securityProfile (Defender, WorkloadIdentity, ImageCleaner, AzureKeyVaultKms)
	securityProfile := &managedclusters.ManagedClusterSecurityProfile{}

	microsoftDefenderRaw := d.Get("microsoft_defender").([]interface{})
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

	securityProfile.ImageCleaner = &managedclusters.ManagedClusterSecurityProfileImageCleaner{
		Enabled:       utils.Bool(d.Get("image_cleaner_enabled").(bool)),
		IntervalHours: utils.Int64(int64(d.Get("image_cleaner_interval_hours").(int))),
	}

	azureKeyVaultKmsRaw := d.Get("key_management_service").([]interface{})
	securityProfile.AzureKeyVaultKms, err = expandKubernetesClusterAzureKeyVaultKms(ctx, keyVaultsClient, resourcesClient, d, azureKeyVaultKmsRaw)
	if err != nil {
		return err
	}

	autoUpgradeProfile := &managedclusters.ManagedClusterAutoUpgradeProfile{}
	autoChannelUpgrade := d.Get("automatic_channel_upgrade").(string)
	nodeOsChannelUpgrade := d.Get("node_os_channel_upgrade").(string)

	// this check needs to be separate and gated since node_os_channel_upgrade is a preview feature
	if nodeOsChannelUpgrade != "" && autoChannelUpgrade != "" {
		if autoChannelUpgrade == string(managedclusters.UpgradeChannelNodeNegativeimage) && nodeOsChannelUpgrade != string(managedclusters.NodeOSUpgradeChannelNodeImage) {
			return fmt.Errorf("`node_os_channel_upgrade` cannot be set to a value other than `NodeImage` if `automatic_channel_upgrade` is set to `node-image`")
		}
	}

	if autoChannelUpgrade != "" {
		autoUpgradeProfile.UpgradeChannel = pointer.To(managedclusters.UpgradeChannel(autoChannelUpgrade))
	} else {
		autoUpgradeProfile.UpgradeChannel = pointer.To(managedclusters.UpgradeChannelNone)
	}

	if nodeOsChannelUpgrade != "" {
		autoUpgradeProfile.NodeOSUpgradeChannel = pointer.To(managedclusters.NodeOSUpgradeChannel(nodeOsChannelUpgrade))
	}

	if customCaTrustCertListRaw := d.Get("custom_ca_trust_certificates_base64").([]interface{}); len(customCaTrustCertListRaw) > 0 {
		securityProfile.CustomCATrustCertificates = convertCustomCaTrustCertsInput(customCaTrustCertListRaw)
	}

	parameters := managedclusters.ManagedCluster{
		ExtendedLocation: expandEdgeZone(d.Get("edge_zone").(string)),
		Location:         location,
		Sku: &managedclusters.ManagedClusterSKU{
			Name: utils.ToPtr(managedclusters.ManagedClusterSKUNameBase), // the only possible value at this point
			Tier: utils.ToPtr(managedclusters.ManagedClusterSKUTier(d.Get("sku_tier").(string))),
		},
		Properties: &managedclusters.ManagedClusterProperties{
			ApiServerAccessProfile:    apiAccessProfile,
			AadProfile:                azureADProfile,
			AddonProfiles:             addonProfiles,
			AgentPoolProfiles:         agentProfiles,
			AutoScalerProfile:         autoScalerProfile,
			AutoUpgradeProfile:        autoUpgradeProfile,
			AzureMonitorProfile:       azureMonitorProfile,
			DnsPrefix:                 utils.String(dnsPrefix),
			EnableRBAC:                utils.Bool(d.Get("role_based_access_control_enabled").(bool)),
			KubernetesVersion:         utils.String(kubernetesVersion),
			LinuxProfile:              linuxProfile,
			WindowsProfile:            windowsProfile,
			NetworkProfile:            networkProfile,
			NodeResourceGroup:         utils.String(nodeResourceGroup),
			DisableLocalAccounts:      utils.Bool(d.Get("local_account_disabled").(bool)),
			HTTPProxyConfig:           httpProxyConfig,
			OidcIssuerProfile:         oidcIssuerProfile,
			SecurityProfile:           securityProfile,
			StorageProfile:            storageProfile,
			WorkloadAutoScalerProfile: workloadAutoscalerProfile,
		},
		Tags: tags.Expand(t),
	}
	managedClusterIdentityRaw := d.Get("identity").([]interface{})
	kubernetesClusterIdentityRaw := d.Get("kubelet_identity").([]interface{})
	servicePrincipalProfileRaw := d.Get("service_principal").([]interface{})

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
		servicePrincipalProfileVal := servicePrincipalProfileRaw[0].(map[string]interface{})
		parameters.Properties.ServicePrincipalProfile = &managedclusters.ManagedClusterServicePrincipalProfile{
			ClientId: servicePrincipalProfileVal["client_id"].(string),
			Secret:   utils.String(servicePrincipalProfileVal["client_secret"].(string)),
		}
		servicePrincipalSet = true
	}

	if v, ok := d.GetOk("private_dns_zone_id"); ok {
		if (parameters.Identity == nil && !servicePrincipalSet) || (v.(string) != "System" && v.(string) != "None" && (!servicePrincipalSet && parameters.Identity.Type != identity.TypeUserAssigned)) {
			return fmt.Errorf("a user assigned identity or a service principal must be used when using a custom private dns zone")
		}
		apiAccessProfile.PrivateDNSZone = utils.String(v.(string))
	}

	if v, ok := d.GetOk("dns_prefix_private_cluster"); ok {
		if !(*apiAccessProfile.EnablePrivateCluster) || apiAccessProfile.PrivateDNSZone == nil || *apiAccessProfile.PrivateDNSZone == "System" || *apiAccessProfile.PrivateDNSZone == "None" {
			return fmt.Errorf("`dns_prefix_private_cluster` should only be set for private cluster with custom private dns zone")
		}
		parameters.Properties.FqdnSubdomain = utils.String(v.(string))
	}

	if v, ok := d.GetOk("disk_encryption_set_id"); ok && v.(string) != "" {
		parameters.Properties.DiskEncryptionSetID = utils.String(v.(string))
	}

	if ingressProfile := expandKubernetesClusterIngressProfile(d, d.Get("web_app_routing").([]interface{})); ingressProfile != nil {
		parameters.Properties.IngressProfile = ingressProfile
	}

	if serviceMeshProfile := expandKubernetesClusterServiceMeshProfile(d.Get("service_mesh_profile").([]interface{}), &managedclusters.ServiceMeshProfile{}); serviceMeshProfile != nil {
		parameters.Properties.ServiceMeshProfile = serviceMeshProfile
	}

	future, err := client.CreateOrUpdate(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	if maintenanceConfigRaw, ok := d.GetOk("maintenance_window"); ok {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		parameters := maintenanceconfigurations.MaintenanceConfiguration{
			Properties: expandKubernetesClusterMaintenanceConfigurationDefault(maintenanceConfigRaw.([]interface{})),
		}
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
		if _, err := client.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
			return fmt.Errorf("creating/updating default maintenance config for %s: %+v", id, err)
		}
	}

	if maintenanceConfigRaw, ok := d.GetOk("maintenance_window_auto_upgrade"); ok {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		parameters := maintenanceconfigurations.MaintenanceConfiguration{
			Properties: expandKubernetesClusterMaintenanceConfiguration(maintenanceConfigRaw.([]interface{})),
		}
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
		if _, err := client.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
			return fmt.Errorf("creating/updating auto upgrade schedule maintenance config for %s: %+v", id, err)
		}
	}

	if maintenanceConfigRaw, ok := d.GetOk("maintenance_window_node_os"); ok {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		parameters := maintenanceconfigurations.MaintenanceConfiguration{
			Properties: expandKubernetesClusterMaintenanceConfiguration(maintenanceConfigRaw.([]interface{})),
		}
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
		if _, err := client.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
			return fmt.Errorf("creating/updating node os upgrade schedule maintenance config for %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	return resourceKubernetesClusterRead(d, meta)
}

func resourceKubernetesClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	containersClient := meta.(*clients.Client).Containers
	nodePoolsClient := containersClient.AgentPoolsClient
	clusterClient := containersClient.KubernetesClustersClient
	keyVaultsClient := meta.(*clients.Client).KeyVault
	resourcesClient := meta.(*clients.Client).Resource
	env := containersClient.Environment
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKubernetesClusterID(d.Id())
	if err != nil {
		return err
	}

	d.Partial(true)

	// we need to conditionally update the cluster
	existing, err := clusterClient.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}
	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving existing %s: `properties` was nil", *id)
	}
	props := existing.Model.Properties

	if err := validateKubernetesCluster(d, existing.Model, id.ResourceGroupName, id.ManagedClusterName); err != nil {
		return err
	}

	// when update, we should set the value of `Identity.UserAssignedIdentities` empty
	// otherwise the rest api will report error - this is tracked here: https://github.com/Azure/azure-rest-api-specs/issues/13631
	if existing.Model.Identity != nil && existing.Model.Identity.IdentityIds != nil {
		for k := range existing.Model.Identity.IdentityIds {
			existing.Model.Identity.IdentityIds[k] = identity.UserAssignedIdentityDetails{}
		}
	}

	if d.HasChange("service_principal") && !d.HasChange("identity") {
		log.Printf("[DEBUG] Updating the Service Principal for %s..", *id)
		servicePrincipals := d.Get("service_principal").([]interface{})
		// we'll be rotating the Service Principal - removing the SP block is handled by the validate function
		servicePrincipalRaw := servicePrincipals[0].(map[string]interface{})

		clientId := servicePrincipalRaw["client_id"].(string)
		clientSecret := servicePrincipalRaw["client_secret"].(string)
		params := managedclusters.ManagedClusterServicePrincipalProfile{
			ClientId: clientId,
			Secret:   utils.String(clientSecret),
		}

		future, err := clusterClient.ResetServicePrincipalProfile(ctx, *id, params)
		if err != nil {
			return fmt.Errorf("updating Service Principal for %s: %+v", *id, err)
		}
		if err = future.Poller.PollUntilDone(); err != nil {
			return fmt.Errorf("waiting for update of Service Principal for %s: %+v", *id, err)
		}
		log.Printf("[DEBUG] Updated the Service Principal for %s.", *id)

		// since we're patching it, re-retrieve the latest version of the cluster
		existing, err = clusterClient.Get(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving updated %s: %+v", *id, err)
		}
		if existing.Model == nil || existing.Model.Properties == nil {
			return fmt.Errorf("retrieving updated %s: `properties` was nil", *id)
		}
	}

	// since there's multiple reasons why we could be called into Update, we use this to only update if something's changed that's not SP/Version
	updateCluster := false

	// RBAC profile updates need to be handled atomically before any call to createUpdate as a diff there will create a PropertyChangeNotAllowed error
	if d.HasChange("role_based_access_control_enabled") {
		// check if we can determine current EnableRBAC state - don't do anything destructive if we can't be sure
		if props.EnableRBAC == nil {
			return fmt.Errorf("updating %s: RBAC Enabled was nil", *id)
		}
		rbacEnabled := d.Get("role_based_access_control_enabled").(bool)

		// changing rbacEnabled must still force cluster recreation
		if *props.EnableRBAC == rbacEnabled {
			props.EnableRBAC = utils.Bool(rbacEnabled)
		} else {
			updateCluster = true
		}
	}

	if d.HasChange("azure_active_directory_role_based_access_control") {
		tenantId := meta.(*clients.Client).Account.TenantId
		azureADRaw := d.Get("azure_active_directory_role_based_access_control").([]interface{})
		azureADProfile, err := expandKubernetesClusterAzureActiveDirectoryRoleBasedAccessControl(azureADRaw, tenantId)
		if err != nil {
			return err
		}

		props.AadProfile = azureADProfile
		if props.AadProfile != nil && (props.AadProfile.Managed == nil || !*props.AadProfile.Managed) {
			log.Printf("[DEBUG] Updating the RBAC AAD profile")
			future, err := clusterClient.ResetAADProfile(ctx, *id, *props.AadProfile)
			if err != nil {
				return fmt.Errorf("updating Managed Kubernetes Cluster AAD Profile for %s: %+v", *id, err)
			}

			if err = future.Poller.PollUntilDone(); err != nil {
				return fmt.Errorf("waiting for update of RBAC AAD profile of %s: %+v", *id, err)
			}
		}

		if props.AadProfile != nil && props.AadProfile.Managed != nil && *props.AadProfile.Managed {
			existing.Model.Properties.AadProfile = azureADProfile
			updateCluster = true
		}
	}

	if d.HasChange("aci_connector_linux") || d.HasChange("azure_policy_enabled") || d.HasChange("confidential_computing") || d.HasChange("http_application_routing_enabled") || d.HasChange("oms_agent") || d.HasChange("ingress_application_gateway") || d.HasChange("open_service_mesh_enabled") || d.HasChange("key_vault_secrets_provider") {
		updateCluster = true
		addOns := collectKubernetesAddons(d)
		addonProfiles, err := expandKubernetesAddOns(d, addOns, env)
		if err != nil {
			return err
		}
		existing.Model.Properties.AddonProfiles = addonProfiles
	}

	if d.HasChange("api_server_authorized_ip_ranges") || d.HasChange("run_command_enabled") || d.HasChange("private_cluster_public_fqdn_enabled") || d.HasChange("api_server_access_profile") {
		updateCluster = true

		apiServerProfile := expandKubernetesClusterAPIAccessProfile(d)
		existing.Model.Properties.ApiServerAccessProfile = apiServerProfile
	}

	if d.HasChange("auto_scaler_profile") {
		updateCluster = true
		autoScalerProfileRaw := d.Get("auto_scaler_profile").([]interface{})

		autoScalerProfile := expandKubernetesClusterAutoScalerProfile(autoScalerProfileRaw)
		existing.Model.Properties.AutoScalerProfile = autoScalerProfile
	}

	if d.HasChange("monitor_metrics") {
		updateCluster = true
		azureMonitorKubernetesMetricsRaw := d.Get("monitor_metrics").([]interface{})

		azureMonitorProfile := expandKubernetesClusterAzureMonitorProfile(azureMonitorKubernetesMetricsRaw)
		existing.Model.Properties.AzureMonitorProfile = azureMonitorProfile
	}

	if d.HasChange("enable_pod_security_policy") && d.Get("enable_pod_security_policy").(bool) {
		return fmt.Errorf("The AKS API has removed support for this field on 2020-10-15 and is no longer possible to configure this the Pod Security Policy - as such you'll need to set `enable_pod_security_policy` to `false`")
	}

	if d.HasChange("linux_profile") {
		updateCluster = true
		linuxProfileRaw := d.Get("linux_profile").([]interface{})
		linuxProfile := expandKubernetesClusterLinuxProfile(linuxProfileRaw)
		existing.Model.Properties.LinuxProfile = linuxProfile
	}

	if d.HasChange("local_account_disabled") {
		updateCluster = true
		existing.Model.Properties.DisableLocalAccounts = utils.Bool(d.Get("local_account_disabled").(bool))
	}

	if d.HasChange("network_profile") {
		updateCluster = true

		if existing.Model.Properties.NetworkProfile == nil {
			return fmt.Errorf("updating %s: `network_profile` was nil", *id)
		}

		networkProfile := *existing.Model.Properties.NetworkProfile

		if networkProfile.LoadBalancerProfile != nil {
			loadBalancerProfile := *networkProfile.LoadBalancerProfile

			if key := "network_profile.0.load_balancer_profile.0.effective_outbound_ips"; d.HasChange(key) {
				effectiveOutboundIPs := idsToResourceReferences(d.Get(key))
				loadBalancerProfile.EffectiveOutboundIPs = effectiveOutboundIPs
			}

			if key := "network_profile.0.load_balancer_profile.0.idle_timeout_in_minutes"; d.HasChange(key) {
				idleTimeoutInMinutes := d.Get(key).(int)
				loadBalancerProfile.IdleTimeoutInMinutes = utils.Int64(int64(idleTimeoutInMinutes))
			}

			if key := "network_profile.0.load_balancer_profile.0.managed_outbound_ip_count"; d.HasChange(key) {
				managedOutboundIPCount := d.Get(key).(int)
				loadBalancerProfile.ManagedOutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileManagedOutboundIPs{
					Count: utils.Int64(int64(managedOutboundIPCount)),
				}

				// fixes: Load balancer profile must specify one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs.
				loadBalancerProfile.OutboundIPs = nil
				loadBalancerProfile.OutboundIPPrefixes = nil
			}

			if key := "network_profile.0.load_balancer_profile.0.managed_outbound_ipv6_count"; d.HasChange(key) {
				managedOutboundIPV6Count := d.Get(key).(int)
				if loadBalancerProfile.ManagedOutboundIPs == nil {
					loadBalancerProfile.ManagedOutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileManagedOutboundIPs{}
				}
				loadBalancerProfile.ManagedOutboundIPs.CountIPv6 = utils.Int64(int64(managedOutboundIPV6Count))

				// fixes: Load balancer profile must specify one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs.
				loadBalancerProfile.OutboundIPs = nil
				loadBalancerProfile.OutboundIPPrefixes = nil
			}

			if key := "network_profile.0.load_balancer_profile.0.outbound_ip_address_ids"; d.HasChange(key) {
				outboundIPAddress := d.Get(key)
				if v := outboundIPAddress.(*pluginsdk.Set).List(); len(v) == 0 {
					// sending [] to unset `outbound_ip_address_ids` results in 400 / Bad Request
					// instead we default back to AKS managed outbound which is the default of the AKS API when nothing is provided
					loadBalancerProfile.ManagedOutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileManagedOutboundIPs{
						Count: utils.Int64(1),
					}
					loadBalancerProfile.OutboundIPs = nil
					loadBalancerProfile.OutboundIPPrefixes = nil
				} else {
					publicIPAddressIDs := idsToResourceReferences(d.Get(key))
					loadBalancerProfile.OutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileOutboundIPs{
						PublicIPs: publicIPAddressIDs,
					}

					// fixes: Load balancer profile must specify one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs.
					loadBalancerProfile.ManagedOutboundIPs = nil
					loadBalancerProfile.OutboundIPPrefixes = nil
				}
			}

			if key := "network_profile.0.load_balancer_profile.0.outbound_ip_prefix_ids"; d.HasChange(key) {
				outboundIPPrefixes := d.Get(key)
				if v := outboundIPPrefixes.(*pluginsdk.Set).List(); len(v) == 0 {
					// sending [] to unset `outbound_ip_address_ids` results in 400 / Bad Request
					// instead we default back to AKS managed outbound which is the default of the AKS API when nothing is specified
					loadBalancerProfile.ManagedOutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileManagedOutboundIPs{
						Count: utils.Int64(1),
					}
					loadBalancerProfile.OutboundIPs = nil
					loadBalancerProfile.OutboundIPPrefixes = nil
				} else {
					outboundIPPrefixIDs := idsToResourceReferences(d.Get(key))
					loadBalancerProfile.OutboundIPPrefixes = &managedclusters.ManagedClusterLoadBalancerProfileOutboundIPPrefixes{
						PublicIPPrefixes: outboundIPPrefixIDs,
					}

					// fixes: Load balancer profile must specify one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs.
					loadBalancerProfile.ManagedOutboundIPs = nil
					loadBalancerProfile.OutboundIPs = nil
				}
			}

			if key := "network_profile.0.load_balancer_profile.0.outbound_ports_allocated"; d.HasChange(key) {
				allocatedOutboundPorts := d.Get(key).(int)
				loadBalancerProfile.AllocatedOutboundPorts = utils.Int64(int64(allocatedOutboundPorts))
			}

			existing.Model.Properties.NetworkProfile.LoadBalancerProfile = &loadBalancerProfile
		}

		if networkProfile.NatGatewayProfile != nil {
			natGatewayProfile := *networkProfile.NatGatewayProfile

			if key := "network_profile.0.nat_gateway_profile.0.idle_timeout_in_minutes"; d.HasChange(key) {
				idleTimeoutInMinutes := d.Get(key).(int)
				natGatewayProfile.IdleTimeoutInMinutes = utils.Int64(int64(idleTimeoutInMinutes))
			}

			if key := "network_profile.0.nat_gateway_profile.0.managed_outbound_ip_count"; d.HasChange(key) {
				managedOutboundIPCount := d.Get(key).(int)
				natGatewayProfile.ManagedOutboundIPProfile = &managedclusters.ManagedClusterManagedOutboundIPProfile{
					Count: utils.Int64(int64(managedOutboundIPCount)),
				}
				natGatewayProfile.EffectiveOutboundIPs = nil
			}

			existing.Model.Properties.NetworkProfile.NatGatewayProfile = &natGatewayProfile
		}

		if key := "network_profile.0.ebpf_data_plane"; d.HasChange(key) {
			ebpfDataPlane := d.Get(key).(string)
			existing.Model.Properties.NetworkProfile.NetworkDataplane = pointer.To(managedclusters.NetworkDataplane(ebpfDataPlane))
		}
	}
	if d.HasChange("service_mesh_profile") {
		updateCluster = true
		if serviceMeshProfile := expandKubernetesClusterServiceMeshProfile(d.Get("service_mesh_profile").([]interface{}), existing.Model.Properties.ServiceMeshProfile); serviceMeshProfile != nil {
			existing.Model.Properties.ServiceMeshProfile = serviceMeshProfile
		}
	}

	if d.HasChange("tags") {
		updateCluster = true
		t := d.Get("tags").(map[string]interface{})
		existing.Model.Tags = tags.Expand(t)
	}

	if d.HasChange("windows_profile") {
		updateCluster = true
		windowsProfileRaw := d.Get("windows_profile").([]interface{})
		windowsProfile := expandKubernetesClusterWindowsProfile(windowsProfileRaw)
		existing.Model.Properties.WindowsProfile = windowsProfile
	}

	if d.HasChange("identity") {
		updateCluster = true
		managedClusterIdentityRaw := d.Get("identity").([]interface{})

		expandedIdentity, err := expandKubernetesClusterManagedClusterIdentity(managedClusterIdentityRaw)
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		existing.Model.Identity = expandedIdentity
	}

	if d.HasChange("sku_tier") {
		updateCluster = true
		if existing.Model.Sku == nil {
			basic := managedclusters.ManagedClusterSKUNameBase
			existing.Model.Sku = &managedclusters.ManagedClusterSKU{
				Name: &basic,
			}
		}

		skuTier := managedclusters.ManagedClusterSKUTierFree
		if v := d.Get("sku_tier").(string); v != "" {
			skuTier = managedclusters.ManagedClusterSKUTier(v)
		}
		existing.Model.Sku.Tier = &skuTier
	}

	if d.HasChange("automatic_channel_upgrade") {
		updateCluster = true
		if existing.Model.Properties.AutoUpgradeProfile == nil {
			existing.Model.Properties.AutoUpgradeProfile = &managedclusters.ManagedClusterAutoUpgradeProfile{}
		}
		channel := d.Get("automatic_channel_upgrade").(string)
		if channel == "" {
			channel = string(managedclusters.UpgradeChannelNone)
		}
		existing.Model.Properties.AutoUpgradeProfile.UpgradeChannel = pointer.To(managedclusters.UpgradeChannel(channel))
	}

	if d.HasChange("node_os_channel_upgrade") {
		updateCluster = true
		if d.Get("automatic_channel_upgrade").(string) == string(managedclusters.UpgradeChannelNodeNegativeimage) && d.Get("node_os_channel_upgrade").(string) != string(managedclusters.NodeOSUpgradeChannelNodeImage) {
			return fmt.Errorf("`node_os_channel_upgrade` must be set to `NodeImage` if `automatic_channel_upgrade` is set to `node-image`")
		}
		if existing.Model.Properties.AutoUpgradeProfile == nil {
			existing.Model.Properties.AutoUpgradeProfile = &managedclusters.ManagedClusterAutoUpgradeProfile{}
		}
		existing.Model.Properties.AutoUpgradeProfile.NodeOSUpgradeChannel = pointer.To(managedclusters.NodeOSUpgradeChannel(d.Get("node_os_channel_upgrade").(string)))
	}

	if d.HasChange("http_proxy_config") {
		updateCluster = true
		httpProxyConfigRaw := d.Get("http_proxy_config").([]interface{})
		httpProxyConfig := expandKubernetesClusterHttpProxyConfig(httpProxyConfigRaw)
		existing.Model.Properties.HTTPProxyConfig = httpProxyConfig
	}

	if d.HasChange("oidc_issuer_enabled") {
		updateCluster = true
		oidcIssuerEnabled := d.Get("oidc_issuer_enabled").(bool)
		oidcIssuerProfile := expandKubernetesClusterOidcIssuerProfile(oidcIssuerEnabled)
		existing.Model.Properties.OidcIssuerProfile = oidcIssuerProfile
	}

	if d.HasChanges("key_management_service") {
		updateCluster = true
		azureKeyVaultKmsRaw := d.Get("key_management_service").([]interface{})
		azureKeyVaultKms, _ := expandKubernetesClusterAzureKeyVaultKms(ctx, keyVaultsClient, resourcesClient, d, azureKeyVaultKmsRaw)
		if existing.Model.Properties.SecurityProfile == nil {
			existing.Model.Properties.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
		}
		existing.Model.Properties.SecurityProfile.AzureKeyVaultKms = azureKeyVaultKms
	}

	if d.HasChanges("custom_ca_trust_certificates_base64") {
		updateCluster = true
		customCaTrustCertListRaw := d.Get("custom_ca_trust_certificates_base64").([]interface{})
		existing.Model.Properties.SecurityProfile.CustomCATrustCertificates = convertCustomCaTrustCertsInput(customCaTrustCertListRaw)
	}

	if d.HasChanges("microsoft_defender") {
		updateCluster = true
		microsoftDefenderRaw := d.Get("microsoft_defender").([]interface{})
		microsoftDefender := expandKubernetesClusterMicrosoftDefender(d, microsoftDefenderRaw)
		if existing.Model.Properties.SecurityProfile == nil {
			existing.Model.Properties.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
		}
		existing.Model.Properties.SecurityProfile.Defender = microsoftDefender
	}

	if d.HasChanges("storage_profile") {
		updateCluster = true
		storageProfileRaw := d.Get("storage_profile").([]interface{})
		clusterStorageProfile := expandStorageProfile(storageProfileRaw)
		existing.Model.Properties.StorageProfile = clusterStorageProfile
	}

	if d.HasChange("workload_autoscaler_profile") {
		updateCluster = true
		workloadAutoscalerProfileRaw := d.Get("workload_autoscaler_profile").([]interface{})
		workloadAutoscalerProfile := expandKubernetesClusterWorkloadAutoscalerProfile(workloadAutoscalerProfileRaw, d)
		if workloadAutoscalerProfile == nil {
			existing.Model.Properties.WorkloadAutoScalerProfile = &managedclusters.ManagedClusterWorkloadAutoScalerProfile{
				Keda: &managedclusters.ManagedClusterWorkloadAutoScalerProfileKeda{
					Enabled: false,
				},
			}
		} else {
			existing.Model.Properties.WorkloadAutoScalerProfile = workloadAutoscalerProfile
		}
	}

	if d.HasChanges("workload_identity_enabled") {
		updateCluster = true
		workloadIdentity := d.Get("workload_identity_enabled").(bool)
		if existing.Model.Properties.SecurityProfile == nil {
			existing.Model.Properties.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
		}
		existing.Model.Properties.SecurityProfile.WorkloadIdentity = &managedclusters.ManagedClusterSecurityProfileWorkloadIdentity{
			Enabled: &workloadIdentity,
		}
	}

	if d.HasChange("image_cleaner_enabled") || d.HasChange("image_cleaner_interval_hours") {
		updateCluster = true
		if existing.Model.Properties.SecurityProfile == nil {
			existing.Model.Properties.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
		}
		existing.Model.Properties.SecurityProfile.ImageCleaner = &managedclusters.ManagedClusterSecurityProfileImageCleaner{
			Enabled:       utils.Bool(d.Get("image_cleaner_enabled").(bool)),
			IntervalHours: utils.Int64(int64(d.Get("image_cleaner_interval_hours").(int))),
		}
	}

	if d.HasChange("web_app_routing") {
		updateCluster = true
		existing.Model.Properties.IngressProfile = expandKubernetesClusterIngressProfile(d, d.Get("web_app_routing").([]interface{}))
	}

	if updateCluster {
		// If Defender was explicitly disabled in a prior update then we should strip SecurityProfile.AzureDefender from the request
		// body to prevent errors in cases where Defender is disabled for the entire subscription
		if !d.HasChanges("microsoft_defender") && len(d.Get("microsoft_defender").([]interface{})) == 0 {
			if existing.Model.Properties.SecurityProfile == nil {
				existing.Model.Properties.SecurityProfile = &managedclusters.ManagedClusterSecurityProfile{}
			}
			existing.Model.Properties.SecurityProfile.Defender = nil
		}

		log.Printf("[DEBUG] Updating %s..", *id)
		future, err := clusterClient.CreateOrUpdate(ctx, *id, *existing.Model)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", *id, err)
		}

		if err = future.Poller.PollUntilDone(); err != nil {
			return fmt.Errorf("waiting for update of %s: %+v", *id, err)
		}
		log.Printf("[DEBUG] Updated %s..", *id)
	}

	// then roll the version of Kubernetes if necessary
	if d.HasChange("kubernetes_version") {
		existing, err = clusterClient.Get(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving existing %s: %+v", *id, err)
		}
		if existing.Model == nil || existing.Model.Properties == nil {
			return fmt.Errorf("retrieving existing %s: `properties` was nil", *id)
		}

		kubernetesVersion := d.Get("kubernetes_version").(string)
		log.Printf("[DEBUG] Upgrading the version of Kubernetes to %q..", kubernetesVersion)
		existing.Model.Properties.KubernetesVersion = utils.String(kubernetesVersion)

		future, err := clusterClient.CreateOrUpdate(ctx, *id, *existing.Model)
		if err != nil {
			return fmt.Errorf("updating Kubernetes Version for %s: %+v", *id, err)
		}

		if err = future.Poller.PollUntilDone(); err != nil {
			return fmt.Errorf("waiting for update of %s: %+v", *id, err)
		}

		log.Printf("[DEBUG] Upgraded the version of Kubernetes to %q..", kubernetesVersion)
	}

	// update the node pool using the separate API
	if d.HasChange("default_node_pool") {
		agentProfiles, err := ExpandDefaultNodePool(d)
		if err != nil {
			return fmt.Errorf("expanding `default_node_pool`: %+v", err)
		}
		agentProfile := ConvertDefaultNodePoolToAgentPool(agentProfiles)
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

			if err := validateNodePoolSupportsVersion(ctx, containersClient, currentNodePoolVersion, defaultNodePoolId, *nodePoolVersion); err != nil {
				return err
			}
		}

		cycleNodePoolProperties := []string{
			"default_node_pool.0.name",
			"default_node_pool.0.enable_host_encryption",
			"default_node_pool.0.enable_node_public_ip",
			"default_node_pool.0.kubelet_config",
			"default_node_pool.0.linux_os_config",
			"default_node_pool.0.max_pods",
			"default_node_pool.0.node_taints",
			"default_node_pool.0.only_critical_addons_enabled",
			"default_node_pool.0.os_disk_size_gb",
			"default_node_pool.0.os_disk_type",
			"default_node_pool.0.os_sku",
			"default_node_pool.0.pod_subnet_id",
			"default_node_pool.0.snapshot_id",
			"default_node_pool.0.ultra_ssd_enabled",
			"default_node_pool.0.vnet_subnet_id",
			"default_node_pool.0.vm_size",
			"default_node_pool.0.zones",
		}

		// if the default node pool name has changed, it means the initial attempt at resizing failed
		if d.HasChanges(cycleNodePoolProperties...) {
			log.Printf("[DEBUG] Cycling Default Node Pool..")
			// to provide a seamless updating experience for the vm size of the default node pool we need to cycle the default
			// node pool by provisioning a temporary system node pool, tearing down the former default node pool and then
			// bringing up the new one.

			if v := d.Get("default_node_pool.0.temporary_name_for_rotation").(string); v == "" {
				return fmt.Errorf("`temporary_name_for_rotation` must be specified when updating any of the following properties %q", cycleNodePoolProperties)
			}

			temporaryNodePoolName := d.Get("default_node_pool.0.temporary_name_for_rotation").(string)
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
			// if the temp node pool already exists due to a previous failure, don't bother spinning it up
			if tempExisting.Model == nil {
				if err := retrySystemNodePoolCreation(ctx, nodePoolsClient, tempNodePoolId, tempAgentProfile); err != nil {
					return fmt.Errorf("creating temporary %s: %+v", tempNodePoolId, err)
				}
			}

			ignorePodDisruptionBudget := true
			deleteOpts := agentpools.DeleteOperationOptions{
				IgnorePodDisruptionBudget: &ignorePodDisruptionBudget,
			}
			// delete the old default node pool if it exists
			if defaultExisting.Model != nil {
				if err := nodePoolsClient.DeleteThenPoll(ctx, defaultNodePoolId, deleteOpts); err != nil {
					return fmt.Errorf("deleting default %s: %+v", defaultNodePoolId, err)
				}
			}

			// create the default node pool with the new vm size
			if err := retrySystemNodePoolCreation(ctx, nodePoolsClient, defaultNodePoolId, agentProfile); err != nil {
				// if creation of the default node pool fails we automatically fall back to the temporary node pool
				// in func findDefaultNodePool
				log.Printf("[DEBUG] Creation of resized default node pool failed")
				return fmt.Errorf("creating default %s: %+v", defaultNodePoolId, err)
			}

			if err := nodePoolsClient.DeleteThenPoll(ctx, tempNodePoolId, deleteOpts); err != nil {
				return fmt.Errorf("deleting temporary %s: %+v", tempNodePoolId, err)
			}

			log.Printf("[DEBUG] Cycled Default Node Pool..")
		} else {
			log.Printf("[DEBUG] Updating of Default Node Pool..")

			if err := nodePoolsClient.CreateOrUpdateThenPoll(ctx, defaultNodePoolId, agentProfile); err != nil {
				return fmt.Errorf("updating Default Node Pool %s %+v", defaultNodePoolId, err)
			}

			log.Printf("[DEBUG] Updated Default Node Pool.")
		}
	}

	if d.HasChange("maintenance_window") {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		maintenanceWindowProperties := expandKubernetesClusterMaintenanceConfigurationDefault(d.Get("maintenance_window").([]interface{}))
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
		if maintenanceWindowProperties != nil {
			parameters := maintenanceconfigurations.MaintenanceConfiguration{
				Properties: maintenanceWindowProperties,
			}
			if _, err := client.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
				return fmt.Errorf("creating/updating Maintenance Configuration for Managed Kubernetes Cluster (%q): %+v", id, err)
			}
		} else {
			if _, err := client.Delete(ctx, maintenanceId); err != nil {
				return fmt.Errorf("deleting Maintenance Configuration for %s: %+v", id, err)
			}
		}
	}

	if d.HasChange("maintenance_window_auto_upgrade") {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		maintenanceWindowProperties := expandKubernetesClusterMaintenanceConfiguration(d.Get("maintenance_window_auto_upgrade").([]interface{}))
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
		if maintenanceWindowProperties != nil {
			parameters := maintenanceconfigurations.MaintenanceConfiguration{
				Properties: maintenanceWindowProperties,
			}
			if _, err := client.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
				return fmt.Errorf("creating/updating Auto Upgrade Schedule Maintenance Configuration for %s: %+v", id, err)
			}
		} else {
			if _, err := client.Delete(ctx, maintenanceId); err != nil {
				return fmt.Errorf("deleting Auto Upgrade Schedule Maintenance Configuration for %s: %+v", id, err)
			}
		}
	}

	if d.HasChange("maintenance_window_node_os") {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
		maintenanceWindowProperties := expandKubernetesClusterMaintenanceConfiguration(d.Get("maintenance_window_node_os").([]interface{}))
		if maintenanceWindowProperties != nil {
			parameters := maintenanceconfigurations.MaintenanceConfiguration{
				Properties: maintenanceWindowProperties,
			}
			if _, err := client.CreateOrUpdate(ctx, maintenanceId, parameters); err != nil {
				return fmt.Errorf("creating/updating Node OS Upgrade Schedule Maintenance Configuration for %s: %+v", id, err)
			}
		} else {
			if _, err := client.Delete(ctx, maintenanceId); err != nil {
				return fmt.Errorf("deleting Node OS Upgrade Schedule Maintenance Configuration for %s: %+v", id, err)
			}
		}
	}

	d.Partial(false)

	return resourceKubernetesClusterRead(d, meta)
}

func resourceKubernetesClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("retrieving User Credentials for %s: %+v", id, err)
	}
	if credentials.Model == nil {
		return fmt.Errorf("retrieving User Credentials for %s: payload is empty", id)
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
			d.Set("enable_pod_security_policy", props.EnablePodSecurityPolicy)
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
			d.Set("automatic_channel_upgrade", upgradeChannel)

			// the API returns `node_os_channel_upgrade` when `automatic_channel_upgrade` is set to `node-image`
			// since it's a preview feature we will only set this if it's explicitly been set in the config for the time being
			if v, ok := d.GetOk("node_os_channel_upgrade"); ok && v.(string) != "" {
				d.Set("node_os_channel_upgrade", nodeOSUpgradeChannel)
			}

			customCaTrustCertList := flattenCustomCaTrustCerts(props.SecurityProfile.CustomCATrustCertificates)
			d.Set("custom_ca_trust_certificates_base64", customCaTrustCertList)

			enablePrivateCluster := false
			enablePrivateClusterPublicFQDN := false
			runCommandEnabled := true

			apiServerAccessProfile := flattenKubernetesClusterAPIAccessProfile(props.ApiServerAccessProfile)
			if err := d.Set("api_server_access_profile", apiServerAccessProfile); err != nil {
				return fmt.Errorf("setting `api_server_access_profile`: %+v", err)
			}
			if accessProfile := props.ApiServerAccessProfile; accessProfile != nil {
				if !features.FourPointOhBeta() {
					apiServerAuthorizedIPRanges := utils.FlattenStringSlice(accessProfile.AuthorizedIPRanges)
					if err := d.Set("api_server_authorized_ip_ranges", apiServerAuthorizedIPRanges); err != nil {
						return fmt.Errorf("setting `api_server_authorized_ip_ranges`: %+v", err)
					}
				}
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
					d.Set("private_dns_zone_id", "System")
				case accessProfile.PrivateDNSZone != nil && strings.EqualFold("None", *accessProfile.PrivateDNSZone):
					d.Set("private_dns_zone_id", "None")
				default:
					d.Set("private_dns_zone_id", accessProfile.PrivateDNSZone)
				}
			}
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

			azureMonitorProfile := flattenKubernetesClusterAzureMonitorProfile(props.AzureMonitorProfile)
			if err := d.Set("monitor_metrics", azureMonitorProfile); err != nil {
				return fmt.Errorf("setting `monitor_metrics`: %+v", err)
			}

			serviceMeshProfile := flattenKubernetesClusterAzureServiceMeshProfile(props.ServiceMeshProfile)
			if err := d.Set("service_mesh_profile", serviceMeshProfile); err != nil {
				return fmt.Errorf("setting `service_mesh_profile`: %+v", err)
			}

			flattenedDefaultNodePool, err := FlattenDefaultNodePool(props.AgentPoolProfiles, d)
			if err != nil {
				return fmt.Errorf("flattening `default_node_pool`: %+v", err)
			}
			if err := d.Set("default_node_pool", flattenedDefaultNodePool); err != nil {
				return fmt.Errorf("setting `default_node_pool`: %+v", err)
			}

			kubeletIdentity := []interface{}{}
			if identityProfile := props.IdentityProfile; identityProfile != nil {
				kubeletIdentity, err = flattenKubernetesClusterIdentityProfile(*props.IdentityProfile)
				if err != nil {
					return err
				}
			}

			if err := d.Set("kubelet_identity", kubeletIdentity); err != nil {
				return fmt.Errorf("setting `kubelet_identity`: %+v", err)
			}

			linuxProfile := flattenKubernetesClusterLinuxProfile(props.LinuxProfile)
			if err := d.Set("linux_profile", linuxProfile); err != nil {
				return fmt.Errorf("setting `linux_profile`: %+v", err)
			}

			networkProfileRaw := d.Get("network_profile").([]interface{})
			networkProfile := flattenKubernetesClusterNetworkProfile(props.NetworkProfile, networkProfileRaw)
			if err := d.Set("network_profile", networkProfile); err != nil {
				return fmt.Errorf("setting `network_profile`: %+v", err)
			}

			rbacEnabled := true
			if props.EnableRBAC != nil {
				rbacEnabled = *props.EnableRBAC
			}
			d.Set("role_based_access_control_enabled", rbacEnabled)

			aadRbac := flattenKubernetesClusterAzureActiveDirectoryRoleBasedAccessControl(props, d)
			if err := d.Set("azure_active_directory_role_based_access_control", aadRbac); err != nil {
				return fmt.Errorf("setting `azure_active_directory_role_based_access_control`: %+v", err)
			}

			servicePrincipal := flattenAzureRmKubernetesClusterServicePrincipalProfile(props.ServicePrincipalProfile, d)
			if err := d.Set("service_principal", servicePrincipal); err != nil {
				return fmt.Errorf("setting `service_principal`: %+v", err)
			}

			windowsProfile := flattenKubernetesClusterWindowsProfile(props.WindowsProfile, d)
			if err := d.Set("windows_profile", windowsProfile); err != nil {
				return fmt.Errorf("setting `windows_profile`: %+v", err)
			}

			workloadAutoscalerProfile := flattenKubernetesClusterWorkloadAutoscalerProfile(props.WorkloadAutoScalerProfile)
			if err := d.Set("workload_autoscaler_profile", workloadAutoscalerProfile); err != nil {
				return fmt.Errorf("setting `workload_autoscaler_profile`: %+v", err)
			}

			if props.SecurityProfile != nil && props.SecurityProfile.ImageCleaner != nil {
				if props.SecurityProfile.ImageCleaner.Enabled != nil {
					d.Set("image_cleaner_enabled", props.SecurityProfile.ImageCleaner.Enabled)
				}
				if props.SecurityProfile.ImageCleaner.IntervalHours != nil {
					d.Set("image_cleaner_interval_hours", props.SecurityProfile.ImageCleaner.IntervalHours)
				}
			}

			httpProxyConfig := flattenKubernetesClusterHttpProxyConfig(props)
			if err := d.Set("http_proxy_config", httpProxyConfig); err != nil {
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

			microsoftDefender := flattenKubernetesClusterMicrosoftDefender(props.SecurityProfile)
			if err := d.Set("microsoft_defender", microsoftDefender); err != nil {
				return fmt.Errorf("setting `microsoft_defender`: %+v", err)
			}

			ingressProfile := flattenKubernetesClusterIngressProfile(props.IngressProfile)
			if err := d.Set("web_app_routing", ingressProfile); err != nil {
				return fmt.Errorf("setting `web_app_routing`: %+v", err)
			}

			workloadIdentity := false
			if props.SecurityProfile != nil && props.SecurityProfile.WorkloadIdentity != nil {
				workloadIdentity = *props.SecurityProfile.WorkloadIdentity.Enabled
			}
			d.Set("workload_identity_enabled", workloadIdentity)

			azureKeyVaultKms := flattenKubernetesClusterDataSourceKeyVaultKms(props.SecurityProfile)
			if err := d.Set("key_management_service", azureKeyVaultKms); err != nil {
				return fmt.Errorf("setting `key_management_service`: %+v", err)
			}

			// adminProfile is only available for RBAC enabled clusters with AAD and local account is not disabled
			var adminKubeConfigRaw *string
			adminKubeConfig := make([]interface{}, 0)
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

		maintenanceConfigurationsClient := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
		configResp, _ := maintenanceConfigurationsClient.Get(ctx, maintenanceId)
		if configurationBody := configResp.Model; configurationBody != nil && configurationBody.Properties != nil {
			d.Set("maintenance_window", flattenKubernetesClusterMaintenanceConfigurationDefault(configurationBody.Properties))
		}

		maintenanceId = maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedAutoUpgradeSchedule")
		configResp, _ = maintenanceConfigurationsClient.Get(ctx, maintenanceId)
		if configurationBody := configResp.Model; configurationBody != nil && configurationBody.Properties != nil && configurationBody.Properties.MaintenanceWindow != nil {
			d.Set("maintenance_window_auto_upgrade", flattenKubernetesClusterMaintenanceConfiguration(configurationBody.Properties.MaintenanceWindow))
		}

		maintenanceId = maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "aksManagedNodeOSUpgradeSchedule")
		configResp, _ = maintenanceConfigurationsClient.Get(ctx, maintenanceId)
		if configurationBody := configResp.Model; configurationBody != nil && configurationBody.Properties != nil && configurationBody.Properties.MaintenanceWindow != nil {
			d.Set("maintenance_window_node_os", flattenKubernetesClusterMaintenanceConfiguration(configurationBody.Properties.MaintenanceWindow))
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceKubernetesClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.KubernetesClustersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseKubernetesClusterID(d.Id())
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("maintenance_window"); ok {
		client := meta.(*clients.Client).Containers.MaintenanceConfigurationsClient
		maintenanceId := maintenanceconfigurations.NewMaintenanceConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName, "default")
		if _, err := client.Delete(ctx, maintenanceId); err != nil {
			return fmt.Errorf("deleting Maintenance Configuration for %s: %+v", *id, err)
		}
	}

	ignorePodDisruptionBudget := true
	future, err := client.Delete(ctx, *id, managedclusters.DeleteOperationOptions{
		IgnorePodDisruptionBudget: &ignorePodDisruptionBudget,
	})
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandKubernetesClusterLinuxProfile(input []interface{}) *managedclusters.ContainerServiceLinuxProfile {
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

	return &managedclusters.ContainerServiceLinuxProfile{
		AdminUsername: adminUsername,
		Ssh: managedclusters.ContainerServiceSshConfiguration{
			PublicKeys: []managedclusters.ContainerServiceSshPublicKey{
				{
					KeyData: keyData,
				},
			},
		},
	}
}

func expandKubernetesClusterIdentityProfile(input []interface{}) *map[string]managedclusters.UserAssignedIdentity {
	identityProfile := make(map[string]managedclusters.UserAssignedIdentity)
	if len(input) == 0 || input[0] == nil {
		return &identityProfile
	}

	values := input[0].(map[string]interface{})

	if identity.Type(values["user_assigned_identity_id"].(string)) != "" {
		identityProfile["kubeletidentity"] = managedclusters.UserAssignedIdentity{
			ResourceId: utils.String(values["user_assigned_identity_id"].(string)),
			ClientId:   utils.String(values["client_id"].(string)),
			ObjectId:   utils.String(values["object_id"].(string)),
		}
	}

	return &identityProfile
}

func flattenKubernetesClusterIdentityProfile(profile map[string]managedclusters.UserAssignedIdentity) ([]interface{}, error) {
	if profile == nil {
		return []interface{}{}, nil
	}

	kubeletIdentity := make([]interface{}, 0)
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

		kubeletIdentity = append(kubeletIdentity, map[string]interface{}{
			"client_id":                 clientId,
			"object_id":                 objectId,
			"user_assigned_identity_id": userAssignedIdentityId,
		})
	}

	return kubeletIdentity, nil
}

func flattenKubernetesClusterLinuxProfile(profile *managedclusters.ContainerServiceLinuxProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	adminUsername := profile.AdminUsername

	sshKeys := make([]interface{}, 0)
	ssh := profile.Ssh
	if keys := ssh.PublicKeys; keys != nil {
		for _, sshKey := range keys {
			keyData := ""
			if kd := sshKey.KeyData; kd != "" {
				keyData = kd
			}
			sshKeys = append(sshKeys, map[string]interface{}{
				"key_data": keyData,
			})
		}
	}

	return []interface{}{
		map[string]interface{}{
			"admin_username": adminUsername,
			"ssh_key":        sshKeys,
		},
	}
}

func expandKubernetesClusterWindowsProfile(input []interface{}) *managedclusters.ManagedClusterWindowsProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})

	license := managedclusters.LicenseTypeNone
	if v := config["license"].(string); v != "" {
		license = managedclusters.LicenseType(v)
	}

	gmsaProfile := expandGmsaProfile(config["gmsa"].([]interface{}))

	return &managedclusters.ManagedClusterWindowsProfile{
		AdminUsername: config["admin_username"].(string),
		AdminPassword: utils.String(config["admin_password"].(string)),
		LicenseType:   &license,
		GmsaProfile:   gmsaProfile,
	}
}

func expandKubernetesClusterAPIAccessProfile(d *pluginsdk.ResourceData) *managedclusters.ManagedClusterAPIServerAccessProfile {
	apiServerAuthorizedIPRangesRaw := []interface{}{}
	if !features.FourPointOhBeta() {
		apiServerAuthorizedIPRangesRaw = d.Get("api_server_authorized_ip_ranges").(*pluginsdk.Set).List()
	}
	apiServerAuthorizedIPRanges := utils.ExpandStringSlice(apiServerAuthorizedIPRangesRaw)

	enablePrivateCluster := false
	if v, ok := d.GetOk("private_cluster_enabled"); ok {
		enablePrivateCluster = v.(bool)
	}

	apiAccessProfile := &managedclusters.ManagedClusterAPIServerAccessProfile{
		EnablePrivateCluster:           &enablePrivateCluster,
		AuthorizedIPRanges:             apiServerAuthorizedIPRanges,
		EnablePrivateClusterPublicFQDN: utils.Bool(d.Get("private_cluster_public_fqdn_enabled").(bool)),
		DisableRunCommand:              utils.Bool(!d.Get("run_command_enabled").(bool)),
	}

	apiServerAccessProfileRaw := d.Get("api_server_access_profile").([]interface{})
	if len(apiServerAccessProfileRaw) == 0 {
		return apiAccessProfile
	}

	config := apiServerAccessProfileRaw[0].(map[string]interface{})

	if v := config["authorized_ip_ranges"]; v != nil {
		apiServerAuthorizedIPRangesRaw := v.(*pluginsdk.Set).List()
		if apiServerAuthorizedIPRanges := utils.ExpandStringSlice(apiServerAuthorizedIPRangesRaw); len(*apiServerAuthorizedIPRanges) > 0 {
			apiAccessProfile.AuthorizedIPRanges = apiServerAuthorizedIPRanges
		}
	}

	enableVnetIntegration := false
	if v := config["vnet_integration_enabled"]; v != nil {
		enableVnetIntegration = v.(bool)
	}
	apiAccessProfile.EnableVnetIntegration = utils.Bool(enableVnetIntegration)

	subnetId := ""
	if v := config["subnet_id"]; v != nil {
		subnetId = v.(string)
	}
	apiAccessProfile.SubnetId = utils.String(subnetId)

	return apiAccessProfile
}

func flattenKubernetesClusterAPIAccessProfile(profile *managedclusters.ManagedClusterAPIServerAccessProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	apiServerAuthorizedIPRanges := utils.FlattenStringSlice(profile.AuthorizedIPRanges)

	enableVnetIntegration := false
	if profile.EnableVnetIntegration != nil {
		enableVnetIntegration = *profile.EnableVnetIntegration
	}
	subnetId := ""
	if profile.SubnetId != nil && *profile.SubnetId != "" {
		subnetId = *profile.SubnetId
	}

	return []interface{}{
		map[string]interface{}{
			"authorized_ip_ranges":     apiServerAuthorizedIPRanges,
			"subnet_id":                subnetId,
			"vnet_integration_enabled": enableVnetIntegration,
		},
	}
}

func expandKubernetesClusterWorkloadAutoscalerProfile(input []interface{}, d *pluginsdk.ResourceData) *managedclusters.ManagedClusterWorkloadAutoScalerProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})

	var workloadAutoscalerProfile managedclusters.ManagedClusterWorkloadAutoScalerProfile

	if v := config["keda_enabled"].(bool); v || d.HasChange("workload_autoscaler_profile.0.keda_enabled") {
		workloadAutoscalerProfile.Keda = &managedclusters.ManagedClusterWorkloadAutoScalerProfileKeda{
			Enabled: v,
		}
	}

	if v := config["vertical_pod_autoscaler_enabled"].(bool); v || d.HasChange("workload_autoscaler_profile.0.vertical_pod_autoscaler_enabled") {
		workloadAutoscalerProfile.VerticalPodAutoscaler = &managedclusters.ManagedClusterWorkloadAutoScalerProfileVerticalPodAutoscaler{
			Enabled: v,
		}
	}

	return &workloadAutoscalerProfile
}

func expandGmsaProfile(input []interface{}) *managedclusters.WindowsGmsaProfile {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})
	return &managedclusters.WindowsGmsaProfile{
		Enabled:        utils.Bool(true),
		DnsServer:      utils.String(config["dns_server"].(string)),
		RootDomainName: utils.String(config["root_domain"].(string)),
	}
}

func flattenKubernetesClusterWindowsProfile(profile *managedclusters.ManagedClusterWindowsProfile, d *pluginsdk.ResourceData) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	adminUsername := profile.AdminUsername

	// admin password isn't returned, so let's look it up
	adminPassword := ""
	if v, ok := d.GetOk("windows_profile.0.admin_password"); ok {
		adminPassword = v.(string)
	}

	license := ""
	if profile.LicenseType != nil && *profile.LicenseType != managedclusters.LicenseTypeNone {
		license = string(*profile.LicenseType)
	}

	gmsaProfile := flattenGmsaProfile(profile.GmsaProfile)

	return []interface{}{
		map[string]interface{}{
			"admin_password": adminPassword,
			"admin_username": adminUsername,
			"license":        license,
			"gmsa":           gmsaProfile,
		},
	}
}

func flattenKubernetesClusterWorkloadAutoscalerProfile(profile *managedclusters.ManagedClusterWorkloadAutoScalerProfile) []interface{} {
	// The API always returns an empty WorkloadAutoScalerProfile object even if none of these values have ever been set
	if profile == nil || (profile.Keda == nil && profile.VerticalPodAutoscaler == nil) {
		return []interface{}{}
	}

	kedaEnabled := false
	if v := profile.Keda; v != nil && v.Enabled {
		kedaEnabled = v.Enabled
	}

	vpaEnabled := false
	controlledValues := ""
	updateMode := ""
	if v := profile.VerticalPodAutoscaler; v != nil && v.Enabled {
		vpaEnabled = v.Enabled
		controlledValues = string(v.ControlledValues)
		updateMode = string(v.UpdateMode)
	}

	return []interface{}{
		map[string]interface{}{
			"keda_enabled":                              kedaEnabled,
			"vertical_pod_autoscaler_enabled":           vpaEnabled,
			"vertical_pod_autoscaler_update_mode":       updateMode,
			"vertical_pod_autoscaler_controlled_values": controlledValues,
		},
	}
}

func flattenGmsaProfile(profile *managedclusters.WindowsGmsaProfile) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	dnsServer := ""
	if dns := profile.DnsServer; dns != nil {
		dnsServer = *dns
	}

	rootDomainName := ""
	if domain := profile.RootDomainName; domain != nil {
		rootDomainName = *domain
	}

	return []interface{}{
		map[string]interface{}{
			"dns_server":  dnsServer,
			"root_domain": rootDomainName,
		},
	}
}

func expandKubernetesClusterNetworkProfile(input []interface{}) (*managedclusters.ContainerServiceNetworkProfile, error) {
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
	natGatewayProfileRaw := config["nat_gateway_profile"].([]interface{})
	outboundType := config["outbound_type"].(string)
	ipVersions, err := expandIPVersions(config["ip_versions"].([]interface{}))
	if err != nil {
		return nil, err
	}

	networkProfile := managedclusters.ContainerServiceNetworkProfile{
		NetworkPlugin:   utils.ToPtr(managedclusters.NetworkPlugin(networkPlugin)),
		NetworkMode:     utils.ToPtr(managedclusters.NetworkMode(networkMode)),
		NetworkPolicy:   utils.ToPtr(managedclusters.NetworkPolicy(networkPolicy)),
		LoadBalancerSku: utils.ToPtr(managedclusters.LoadBalancerSku(loadBalancerSku)),
		OutboundType:    utils.ToPtr(managedclusters.OutboundType(outboundType)),
		IPFamilies:      ipVersions,
	}

	if ebpfDataPlane := config["ebpf_data_plane"].(string); ebpfDataPlane != "" {
		networkProfile.NetworkDataplane = utils.ToPtr(managedclusters.NetworkDataplane(ebpfDataPlane))
	}
	if networkPluginMode := config["network_plugin_mode"].(string); networkPluginMode != "" {
		networkProfile.NetworkPluginMode = utils.ToPtr(managedclusters.NetworkPluginMode(networkPluginMode))
	}

	if len(loadBalancerProfileRaw) > 0 {
		if !strings.EqualFold(loadBalancerSku, "standard") {
			return nil, fmt.Errorf("only load balancer SKU 'Standard' supports load balancer profiles. Provided load balancer type: %s", loadBalancerSku)
		}

		networkProfile.LoadBalancerProfile = expandLoadBalancerProfile(loadBalancerProfileRaw)
	}

	if len(natGatewayProfileRaw) > 0 {
		if !strings.EqualFold(loadBalancerSku, "standard") {
			return nil, fmt.Errorf("only load balancer SKU 'Standard' supports NAT Gateway profiles. Provided load balancer type: %s", loadBalancerSku)
		}

		networkProfile.NatGatewayProfile = expandNatGatewayProfile(natGatewayProfileRaw)
	}

	if v, ok := config["dns_service_ip"]; ok && v.(string) != "" {
		dnsServiceIP := v.(string)
		networkProfile.DnsServiceIP = utils.String(dnsServiceIP)
	}

	if v, ok := config["pod_cidr"]; ok && v.(string) != "" {
		podCidr := v.(string)
		networkProfile.PodCidr = utils.String(podCidr)
	}

	if v, ok := config["pod_cidrs"]; ok {
		networkProfile.PodCidrs = utils.ExpandStringSlice(v.([]interface{}))
	}

	if v, ok := config["service_cidr"]; ok && v.(string) != "" {
		serviceCidr := v.(string)
		networkProfile.ServiceCidr = utils.String(serviceCidr)
	}

	if v, ok := config["service_cidrs"]; ok {
		networkProfile.ServiceCidrs = utils.ExpandStringSlice(v.([]interface{}))
	}

	return &networkProfile, nil
}

func expandLoadBalancerProfile(d []interface{}) *managedclusters.ManagedClusterLoadBalancerProfile {
	if d[0] == nil {
		return nil
	}

	config := d[0].(map[string]interface{})

	profile := &managedclusters.ManagedClusterLoadBalancerProfile{}

	if mins, ok := config["idle_timeout_in_minutes"]; ok && mins.(int) != 0 {
		profile.IdleTimeoutInMinutes = utils.Int64(int64(mins.(int)))
	}

	if port, ok := config["outbound_ports_allocated"].(int); ok {
		profile.AllocatedOutboundPorts = utils.Int64(int64(port))
	}

	if ipCount := config["managed_outbound_ip_count"]; ipCount != nil {
		if c := int64(ipCount.(int)); c > 0 {
			profile.ManagedOutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileManagedOutboundIPs{Count: &c}
		}
	}

	if ipv6Count := config["managed_outbound_ipv6_count"]; ipv6Count != nil {
		if c := int64(ipv6Count.(int)); c > 0 {
			if profile.ManagedOutboundIPs == nil {
				profile.ManagedOutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileManagedOutboundIPs{}
			}
			profile.ManagedOutboundIPs.CountIPv6 = &c
		}
	}

	if ipPrefixes := idsToResourceReferences(config["outbound_ip_prefix_ids"]); ipPrefixes != nil {
		profile.OutboundIPPrefixes = &managedclusters.ManagedClusterLoadBalancerProfileOutboundIPPrefixes{PublicIPPrefixes: ipPrefixes}
	}

	if outIps := idsToResourceReferences(config["outbound_ip_address_ids"]); outIps != nil {
		profile.OutboundIPs = &managedclusters.ManagedClusterLoadBalancerProfileOutboundIPs{PublicIPs: outIps}
	}

	return profile
}

func expandIPVersions(input []interface{}) (*[]managedclusters.IPFamily, error) {
	if len(input) == 0 {
		return nil, nil
	}

	ipv := make([]managedclusters.IPFamily, 0)
	for _, data := range input {
		ipv = append(ipv, managedclusters.IPFamily(data.(string)))
	}

	if len(ipv) == 1 && ipv[0] == managedclusters.IPFamilyIPvSix {
		return nil, fmt.Errorf("`ip_versions` must be `IPv4` or `IPv4` and `IPv6`. `IPv6` alone is not supported")
	}

	return &ipv, nil
}

func expandNatGatewayProfile(d []interface{}) *managedclusters.ManagedClusterNATGatewayProfile {
	if d[0] == nil {
		return nil
	}

	config := d[0].(map[string]interface{})

	profile := &managedclusters.ManagedClusterNATGatewayProfile{}

	if mins, ok := config["idle_timeout_in_minutes"]; ok && mins.(int) != 0 {
		profile.IdleTimeoutInMinutes = utils.Int64(int64(mins.(int)))
	}

	if ipCount := config["managed_outbound_ip_count"]; ipCount != nil {
		if c := int64(ipCount.(int)); c > 0 {
			profile.ManagedOutboundIPProfile = &managedclusters.ManagedClusterManagedOutboundIPProfile{Count: &c}
		}
	}

	return profile
}

func idsToResourceReferences(set interface{}) *[]managedclusters.ResourceReference {
	if set == nil {
		return nil
	}

	s := set.(*pluginsdk.Set)
	results := make([]managedclusters.ResourceReference, 0)

	for _, element := range s.List() {
		id := element.(string)
		results = append(results, managedclusters.ResourceReference{Id: &id})
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

func flattenKubernetesClusterNetworkProfile(profile *managedclusters.ContainerServiceNetworkProfile, raw []interface{}) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	dnsServiceIP := ""
	if profile.DnsServiceIP != nil {
		dnsServiceIP = *profile.DnsServiceIP
	}

	dockerBridgeCidr := ""
	if len(raw) != 0 && raw[0] != nil {
		config := raw[0].(map[string]interface{})
		if v, ok := config["docker_bridge_cidr"]; ok && v.(string) != "" {
			dockerBridgeCidr = v.(string)
		}
	}

	serviceCidr := ""
	if profile.ServiceCidr != nil {
		serviceCidr = *profile.ServiceCidr
	}

	podCidr := ""
	if profile.PodCidr != nil {
		podCidr = *profile.PodCidr
	}

	networkPlugin := ""
	if profile.NetworkPlugin != nil {
		networkPlugin = string(*profile.NetworkPlugin)
	}

	networkMode := ""
	if profile.NetworkMode != nil {
		networkMode = string(*profile.NetworkMode)
	}

	networkPolicy := ""
	if profile.NetworkPolicy != nil {
		networkPolicy = string(*profile.NetworkPolicy)
	}

	outboundType := ""
	if profile.OutboundType != nil {
		outboundType = string(*profile.OutboundType)
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

			if countIPv6 := ips.CountIPv6; countIPv6 != nil {
				lb["managed_outbound_ipv6_count"] = countIPv6
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

	ngwProfiles := make([]interface{}, 0)
	if ngwp := profile.NatGatewayProfile; ngwp != nil {
		ng := make(map[string]interface{})

		if v := ngwp.IdleTimeoutInMinutes; v != nil {
			ng["idle_timeout_in_minutes"] = v
		}

		if ips := ngwp.ManagedOutboundIPProfile; ips != nil {
			if count := ips.Count; count != nil {
				ng["managed_outbound_ip_count"] = count
			}
		}

		ng["effective_outbound_ips"] = resourceReferencesToIds(profile.NatGatewayProfile.EffectiveOutboundIPs)
		ngwProfiles = append(ngwProfiles, ng)
	}

	ipVersions := make([]interface{}, 0)
	if ipfs := profile.IPFamilies; ipfs != nil {
		for _, item := range *ipfs {
			ipVersions = append(ipVersions, item)
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
	ebpfDataPlane := ""
	// 2023-02-02-preview returns the default value `azure` for this new property
	// @stephybun: we should look into replacing `ebpf_data_plane` with a new property called `network_data_plane`
	if v := profile.NetworkDataplane; v != nil && *v != managedclusters.NetworkDataplaneAzure {
		ebpfDataPlane = string(*v)
	}

	return []interface{}{
		map[string]interface{}{
			"dns_service_ip":        dnsServiceIP,
			"docker_bridge_cidr":    dockerBridgeCidr,
			"ebpf_data_plane":       ebpfDataPlane,
			"load_balancer_sku":     string(*sku),
			"load_balancer_profile": lbProfiles,
			"nat_gateway_profile":   ngwProfiles,
			"ip_versions":           ipVersions,
			"network_plugin":        networkPlugin,
			"network_plugin_mode":   networkPluginMode,
			"network_mode":          networkMode,
			"network_policy":        networkPolicy,
			"pod_cidr":              podCidr,
			"pod_cidrs":             utils.FlattenStringSlice(profile.PodCidrs),
			"service_cidr":          serviceCidr,
			"service_cidrs":         utils.FlattenStringSlice(profile.ServiceCidrs),
			"outbound_type":         outboundType,
		},
	}
}

func expandKubernetesClusterAzureActiveDirectoryRoleBasedAccessControl(input []interface{}, providerTenantId string) (*managedclusters.ManagedClusterAADProfile, error) {
	if len(input) == 0 {
		return nil, nil
	}

	var aad *managedclusters.ManagedClusterAADProfile

	azureAdRaw := input[0].(map[string]interface{})

	clientAppId := azureAdRaw["client_app_id"].(string)
	serverAppId := azureAdRaw["server_app_id"].(string)
	serverAppSecret := azureAdRaw["server_app_secret"].(string)
	tenantId := azureAdRaw["tenant_id"].(string)
	managed := azureAdRaw["managed"].(bool)
	azureRbacEnabled := azureAdRaw["azure_rbac_enabled"].(bool)
	adminGroupObjectIdsRaw := azureAdRaw["admin_group_object_ids"].([]interface{})
	adminGroupObjectIds := utils.ExpandStringSlice(adminGroupObjectIdsRaw)

	if tenantId == "" {
		tenantId = providerTenantId
	}

	if managed {
		aad = &managedclusters.ManagedClusterAADProfile{
			TenantID:            utils.String(tenantId),
			Managed:             utils.Bool(managed),
			AdminGroupObjectIDs: adminGroupObjectIds,
			EnableAzureRBAC:     utils.Bool(azureRbacEnabled),
		}

		if clientAppId != "" || serverAppId != "" || serverAppSecret != "" {
			return nil, fmt.Errorf("can't specify client_app_id or server_app_id or server_app_secret when using managed aad rbac (managed = true)")
		}
	} else {
		aad = &managedclusters.ManagedClusterAADProfile{
			ClientAppID:     utils.String(clientAppId),
			ServerAppID:     utils.String(serverAppId),
			ServerAppSecret: utils.String(serverAppSecret),
			TenantID:        utils.String(tenantId),
			Managed:         utils.Bool(managed),
		}

		if len(*adminGroupObjectIds) > 0 {
			return nil, fmt.Errorf("can't specify admin_group_object_ids when using managed aad rbac (managed = false)")
		}

		if clientAppId == "" || serverAppId == "" || serverAppSecret == "" {
			return nil, fmt.Errorf("you must specify client_app_id and server_app_id and server_app_secret when using managed aad rbac (managed = false)")
		}

		if azureRbacEnabled {
			return nil, fmt.Errorf("you must enable Managed AAD before Azure RBAC can be enabled")
		}
	}

	return aad, nil
}

func expandKubernetesClusterManagedClusterIdentity(input []interface{}) (*identity.SystemOrUserAssignedMap, error) {
	expanded, err := identity.ExpandSystemOrUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := identity.SystemOrUserAssignedMap{
		Type: identity.Type(string(expanded.Type)),
	}
	if expanded.Type == identity.TypeUserAssigned {
		out.IdentityIds = make(map[string]identity.UserAssignedIdentityDetails)
		for k := range expanded.IdentityIds {
			out.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				// intentionally empty
			}
		}
	}
	return &out, nil
}

func flattenKubernetesClusterAzureActiveDirectoryRoleBasedAccessControl(input *managedclusters.ManagedClusterProperties, d *pluginsdk.ResourceData) []interface{} {
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
		// azure_active_directory_role_based_access_control.0.server_app_secret
		if existing, ok := d.GetOk("azure_active_directory_role_based_access_control"); ok {
			aadRbacRaw := existing.([]interface{})
			if len(aadRbacRaw) > 0 {
				aadRbac := aadRbacRaw[0].(map[string]interface{})
				if v := aadRbac["server_app_secret"]; v != nil {
					serverAppSecret = v.(string)
				}
			}
		}

		tenantId := ""
		if profile.TenantID != nil {
			tenantId = *profile.TenantID
		}

		results = append(results, map[string]interface{}{
			"admin_group_object_ids": adminGroupObjectIds,
			"client_app_id":          clientAppId,
			"managed":                managed,
			"server_app_id":          serverAppId,
			"server_app_secret":      serverAppSecret,
			"tenant_id":              tenantId,
			"azure_rbac_enabled":     azureRbacEnabled,
		})
	}

	return results
}

func flattenAzureRmKubernetesClusterServicePrincipalProfile(profile *managedclusters.ManagedClusterServicePrincipalProfile, d *pluginsdk.ResourceData) []interface{} {
	if profile == nil {
		return []interface{}{}
	}

	clientId := profile.ClientId

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

func flattenKubernetesClusterAutoScalerProfile(profile *managedclusters.ManagedClusterPropertiesAutoScalerProfile) ([]interface{}, error) {
	if profile == nil {
		return []interface{}{}, nil
	}

	balanceSimilarNodeGroups := false
	if profile.BalanceSimilarNodeGroups != nil {
		// @tombuildsstuff: presumably this'll get converted to a Boolean at some point
		//					at any rate we should use the proper type users expect here
		balanceSimilarNodeGroups = strings.EqualFold(*profile.BalanceSimilarNodeGroups, "true")
	}

	expander := ""
	if profile.Expander != nil {
		expander = string(*profile.Expander)
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
			"expander":                         expander,
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

func expandKubernetesClusterAutoScalerProfile(input []interface{}) *managedclusters.ManagedClusterPropertiesAutoScalerProfile {
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

	return &managedclusters.ManagedClusterPropertiesAutoScalerProfile{
		BalanceSimilarNodeGroups:      utils.String(strconv.FormatBool(balanceSimilarNodeGroups)),
		Expander:                      utils.ToPtr(managedclusters.Expander(expander)),
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

func expandKubernetesClusterAzureKeyVaultKms(ctx context.Context, keyVaultsClient *keyVaultClient.Client, resourcesClient *resourcesClient.Client, d *pluginsdk.ResourceData, input []interface{}) (*managedclusters.AzureKeyVaultKms, error) {
	if ((input == nil) || len(input) == 0) && d.HasChanges("key_management_service") {
		return &managedclusters.AzureKeyVaultKms{
			Enabled: utils.Bool(false),
		}, nil
	} else if (input == nil) || len(input) == 0 {
		return nil, nil
	}

	raw := input[0].(map[string]interface{})
	kvAccess := managedclusters.KeyVaultNetworkAccessTypes(raw["key_vault_network_access"].(string))

	azureKeyVaultKms := &managedclusters.AzureKeyVaultKms{
		Enabled:               utils.Bool(true),
		KeyId:                 utils.String(raw["key_vault_key_id"].(string)),
		KeyVaultNetworkAccess: &kvAccess,
	}

	// Set Key vault Resource ID in case public access is disabled
	if kvAccess == managedclusters.KeyVaultNetworkAccessTypesPrivate {
		keyVaultKeyId, err := keyVaultParse.ParseNestedItemID(*azureKeyVaultKms.KeyId)
		if err != nil {
			return nil, err
		}
		keyVaultID, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, keyVaultKeyId.KeyVaultBaseUrl)
		if err != nil {
			return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", keyVaultKeyId.KeyVaultBaseUrl, err)
		}

		azureKeyVaultKms.KeyVaultResourceId = keyVaultID
	}

	return azureKeyVaultKms, nil
}

func expandKubernetesClusterMaintenanceConfigurationDefault(input []interface{}) *maintenanceconfigurations.MaintenanceConfigurationProperties {
	if len(input) == 0 {
		return nil
	}
	value := input[0].(map[string]interface{})
	return &maintenanceconfigurations.MaintenanceConfigurationProperties{
		NotAllowedTime: expandKubernetesClusterMaintenanceConfigurationTimeSpans(value["not_allowed"].(*pluginsdk.Set).List()),
		TimeInWeek:     expandKubernetesClusterMaintenanceConfigurationTimeInWeeks(value["allowed"].(*pluginsdk.Set).List()),
	}
}

func expandKubernetesClusterMaintenanceConfiguration(input []interface{}) *maintenanceconfigurations.MaintenanceConfigurationProperties {
	if len(input) == 0 {
		return nil
	}
	value := input[0].(map[string]interface{})

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
			UtcOffset:       utils.String(value["utc_offset"].(string)),
			NotAllowedDates: expandKubernetesClusterMaintenanceConfigurationDateSpans(value["not_allowed"].(*pluginsdk.Set).List()),
			Schedule:        schedule,
		},
	}

	if startDateRaw := value["start_date"]; startDateRaw != nil && startDateRaw.(string) != "" {
		startDate, _ := time.Parse(time.RFC3339, startDateRaw.(string))
		output.MaintenanceWindow.StartDate = utils.String(startDate.Format("2006-01-02"))
	}

	if duration := value["duration"]; duration != nil && duration.(int) != 0 {
		output.MaintenanceWindow.DurationHours = int64(duration.(int))
	}

	return output
}

func expandKubernetesClusterMaintenanceConfigurationTimeSpans(input []interface{}) *[]maintenanceconfigurations.TimeSpan {
	results := make([]maintenanceconfigurations.TimeSpan, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		start, _ := time.Parse(time.RFC3339, v["start"].(string))
		end, _ := time.Parse(time.RFC3339, v["end"].(string))
		results = append(results, maintenanceconfigurations.TimeSpan{
			Start: utils.ToPtr(start.Format("2006-01-02T15:04:05Z07:00")),
			End:   utils.ToPtr(end.Format("2006-01-02T15:04:05Z07:00")),
		})
	}
	return &results
}

func expandKubernetesClusterMaintenanceConfigurationDateSpans(input []interface{}) *[]maintenanceconfigurations.DateSpan {
	results := make([]maintenanceconfigurations.DateSpan, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		start, _ := time.Parse(time.RFC3339, v["start"].(string))
		end, _ := time.Parse(time.RFC3339, v["end"].(string))
		results = append(results, maintenanceconfigurations.DateSpan{
			Start: start.Format("2006-01-02"),
			End:   end.Format("2006-01-02"),
		})
	}
	return &results
}

func expandKubernetesClusterMaintenanceConfigurationTimeInWeeks(input []interface{}) *[]maintenanceconfigurations.TimeInWeek {
	results := make([]maintenanceconfigurations.TimeInWeek, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, maintenanceconfigurations.TimeInWeek{
			Day:       utils.ToPtr(maintenanceconfigurations.WeekDay(v["day"].(string))),
			HourSlots: utils.ExpandInt64Slice(v["hours"].(*pluginsdk.Set).List()),
		})
	}
	return &results
}

func flattenKubernetesClusterMaintenanceConfiguration(input *maintenanceconfigurations.MaintenanceWindow) interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	startDate := ""
	if input.StartDate != nil {
		startDate = *input.StartDate + "T00:00:00Z"
	}
	utcOfset := ""
	if input.UtcOffset != nil {
		utcOfset = *input.UtcOffset
	}

	windowProperties := map[string]interface{}{
		"not_allowed": flattenKubernetesClusterMaintenanceConfigurationDateSpans(input.NotAllowedDates),
		"duration":    int(input.DurationHours),
		"start_date":  startDate,
		"start_time":  input.StartTime,
		"utc_offset":  utcOfset,
	}
	// Add flattened schedule properties
	for k, v := range flattenKubernetesClusterMaintenanceConfigurationSchedule(input.Schedule) {
		windowProperties[k] = v
	}

	return append(results, windowProperties)
}

func flattenKubernetesClusterMaintenanceConfigurationSchedule(input maintenanceconfigurations.Schedule) map[string]interface{} {
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

	dayOfMonth := 0
	if input.AbsoluteMonthly != nil {
		frequency = "AbsoluteMonthly"
		interval = input.AbsoluteMonthly.IntervalMonths
		dayOfMonth = int(input.AbsoluteMonthly.DayOfMonth)
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

func flattenKubernetesClusterMaintenanceConfigurationDefault(input *maintenanceconfigurations.MaintenanceConfigurationProperties) interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}
	results = append(results, map[string]interface{}{
		"not_allowed": flattenKubernetesClusterMaintenanceConfigurationTimeSpans(input.NotAllowedTime),
		"allowed":     flattenKubernetesClusterMaintenanceConfigurationTimeInWeeks(input.TimeInWeek),
	})
	return results
}

func flattenKubernetesClusterMaintenanceConfigurationTimeSpans(input *[]maintenanceconfigurations.TimeSpan) []interface{} {
	results := make([]interface{}, 0)
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
		results = append(results, map[string]interface{}{
			"end":   end,
			"start": start,
		})
	}
	return results
}

func flattenKubernetesClusterMaintenanceConfigurationDateSpans(input *[]maintenanceconfigurations.DateSpan) []interface{} {
	results := make([]interface{}, 0)
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
		results = append(results, map[string]interface{}{
			"end":   end + "T00:00:00Z",
			"start": start + "T00:00:00Z",
		})
	}
	return results
}

func flattenKubernetesClusterMaintenanceConfigurationTimeInWeeks(input *[]maintenanceconfigurations.TimeInWeek) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		day := ""
		if item.Day != nil {
			day = string(*item.Day)
		}
		results = append(results, map[string]interface{}{
			"day":   day,
			"hours": utils.FlattenInt64Slice(item.HourSlots),
		})
	}
	return results
}

func expandKubernetesClusterHttpProxyConfig(input []interface{}) *managedclusters.ManagedClusterHTTPProxyConfig {
	httpProxyConfig := managedclusters.ManagedClusterHTTPProxyConfig{}
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]interface{})

	httpProxyConfig.HTTPProxy = utils.String(config["http_proxy"].(string))
	httpProxyConfig.HTTPSProxy = utils.String(config["https_proxy"].(string))
	if value := config["trusted_ca"].(string); len(value) != 0 {
		httpProxyConfig.TrustedCa = utils.String(value)
	}

	noProxyRaw := config["no_proxy"].(*pluginsdk.Set).List()
	httpProxyConfig.NoProxy = utils.ExpandStringSlice(noProxyRaw)

	return &httpProxyConfig
}

func expandKubernetesClusterOidcIssuerProfile(input bool) *managedclusters.ManagedClusterOIDCIssuerProfile {
	oidcIssuerProfile := managedclusters.ManagedClusterOIDCIssuerProfile{}
	oidcIssuerProfile.Enabled = &input

	return &oidcIssuerProfile
}

func flattenKubernetesClusterHttpProxyConfig(props *managedclusters.ManagedClusterProperties) []interface{} {
	if props == nil || props.HTTPProxyConfig == nil {
		return []interface{}{}
	}

	httpProxyConfig := props.HTTPProxyConfig

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

	results := []interface{}{}
	return append(results, map[string]interface{}{
		"http_proxy":  httpProxy,
		"https_proxy": httpsProxy,
		"no_proxy":    noProxyList,
		"trusted_ca":  trustedCa,
	})
}

func expandKubernetesClusterMicrosoftDefender(d *pluginsdk.ResourceData, input []interface{}) *managedclusters.ManagedClusterSecurityProfileDefender {
	if (len(input) == 0 || input[0] == nil) && d.HasChange("microsoft_defender") {
		return &managedclusters.ManagedClusterSecurityProfileDefender{
			SecurityMonitoring: &managedclusters.ManagedClusterSecurityProfileDefenderSecurityMonitoring{
				Enabled: utils.Bool(false),
			},
		}
	} else if len(input) == 0 || input[0] == nil {
		return nil
	}

	config := input[0].(map[string]interface{})
	return &managedclusters.ManagedClusterSecurityProfileDefender{
		SecurityMonitoring: &managedclusters.ManagedClusterSecurityProfileDefenderSecurityMonitoring{
			Enabled: utils.Bool(true),
		},
		LogAnalyticsWorkspaceResourceId: utils.String(config["log_analytics_workspace_id"].(string)),
	}
}

func flattenKubernetesClusterMicrosoftDefender(input *managedclusters.ManagedClusterSecurityProfile) []interface{} {
	if input == nil || input.Defender == nil || input.Defender.SecurityMonitoring == nil || input.Defender.SecurityMonitoring.Enabled == nil || !*input.Defender.SecurityMonitoring.Enabled {
		return []interface{}{}
	}

	logAnalyticsWorkspace := ""
	if v := input.Defender.LogAnalyticsWorkspaceResourceId; v != nil {
		logAnalyticsWorkspace = *v
	}

	return []interface{}{
		map[string]interface{}{
			"log_analytics_workspace_id": logAnalyticsWorkspace,
		},
	}
}

func expandStorageProfile(input []interface{}) *managedclusters.ManagedClusterStorageProfile {
	if (input == nil) || len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})

	profile := managedclusters.ManagedClusterStorageProfile{
		BlobCSIDriver: &managedclusters.ManagedClusterStorageProfileBlobCSIDriver{
			Enabled: utils.Bool(raw["blob_driver_enabled"].(bool)),
		},
		DiskCSIDriver: &managedclusters.ManagedClusterStorageProfileDiskCSIDriver{
			Enabled: utils.Bool(raw["disk_driver_enabled"].(bool)),
			Version: utils.String(raw["disk_driver_version"].(string)),
		},
		FileCSIDriver: &managedclusters.ManagedClusterStorageProfileFileCSIDriver{
			Enabled: utils.Bool(raw["file_driver_enabled"].(bool)),
		},
		SnapshotController: &managedclusters.ManagedClusterStorageProfileSnapshotController{
			Enabled: utils.Bool(raw["snapshot_controller_enabled"].(bool)),
		},
	}

	return &profile
}

func expandEdgeZone(input string) *edgezones.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &edgezones.Model{
		Name: normalized,
	}
}

func flattenEdgeZone(input *edgezones.Model) string {
	// As the `extendedLocation.type` returned by API is always lower case, so it has to use `Normalize` function while comparing them
	if input == nil || input.Name == "" {
		return ""
	}
	return edgezones.NormalizeNilable(&input.Name)
}

func base64Decode(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(data)
}

func base64IsEncoded(data string) bool {
	_, err := base64.StdEncoding.DecodeString(data)
	return err == nil
}

func expandKubernetesClusterServiceMeshProfile(input []interface{}, existing *managedclusters.ServiceMeshProfile) *managedclusters.ServiceMeshProfile {
	if (input == nil) || len(input) == 0 {
		// explicitly disable istio if it was enabled before
		if existing != nil && existing.Mode == managedclusters.ServiceMeshModeIstio {
			return &managedclusters.ServiceMeshProfile{
				Mode: managedclusters.ServiceMeshModeDisabled,
			}
		}

		return nil
	}

	raw := input[0].(map[string]interface{})

	mode := raw["mode"].(string)

	profile := managedclusters.ServiceMeshProfile{}

	if managedclusters.ServiceMeshMode(mode) == managedclusters.ServiceMeshModeIstio {
		profile.Mode = managedclusters.ServiceMeshMode(mode)
		profile.Istio = &managedclusters.IstioServiceMesh{}

		profile.Istio.Components = &managedclusters.IstioComponents{}

		istioIngressGatewaysList := make([]managedclusters.IstioIngressGateway, 0)

		if raw["internal_ingress_gateway_enabled"] != nil {

			ingressGatewayElementInternal := managedclusters.IstioIngressGateway{
				Enabled: raw["internal_ingress_gateway_enabled"].(bool),
				Mode:    managedclusters.IstioIngressGatewayModeInternal,
			}

			istioIngressGatewaysList = append(istioIngressGatewaysList, ingressGatewayElementInternal)
		}

		if raw["external_ingress_gateway_enabled"] != nil {

			ingressGatewayElementExternal := managedclusters.IstioIngressGateway{
				Enabled: raw["external_ingress_gateway_enabled"].(bool),
				Mode:    managedclusters.IstioIngressGatewayModeExternal,
			}

			istioIngressGatewaysList = append(istioIngressGatewaysList, ingressGatewayElementExternal)
		}

		profile.Istio.Components.IngressGateways = &istioIngressGatewaysList
	}

	return &profile
}

func expandKubernetesClusterIngressProfile(d *pluginsdk.ResourceData, input []interface{}) *managedclusters.ManagedClusterIngressProfile {
	if len(input) == 0 && d.HasChange("web_app_routing") {
		return &managedclusters.ManagedClusterIngressProfile{
			WebAppRouting: &managedclusters.ManagedClusterIngressProfileWebAppRouting{
				Enabled: utils.Bool(false),
			},
		}
	} else if len(input) == 0 {
		return nil
	}

	out := managedclusters.ManagedClusterIngressProfile{
		WebAppRouting: &managedclusters.ManagedClusterIngressProfileWebAppRouting{
			Enabled: utils.Bool(true),
		},
	}
	if input[0] != nil {
		config := input[0].(map[string]interface{})
		dnsZoneResourceId := config["dns_zone_id"].(string)
		if dnsZoneResourceId != "" {
			out.WebAppRouting.DnsZoneResourceId = utils.String(dnsZoneResourceId)
		}
	}
	return &out
}

func flattenKubernetesClusterIngressProfile(input *managedclusters.ManagedClusterIngressProfile) []interface{} {
	if input == nil || input.WebAppRouting == nil || (input.WebAppRouting.Enabled != nil && !*input.WebAppRouting.Enabled) {
		return []interface{}{}
	}

	dnsZoneId := ""
	if v := input.WebAppRouting.DnsZoneResourceId; v != nil {
		dnsZoneId = *v
	}

	webAppRoutingIdentity := []interface{}{}

	if v := input.WebAppRouting.Identity; v != nil {
		webAppRoutingIdentity = flattenKubernetesClusterAddOnIdentityProfile(v)
	}

	return []interface{}{
		map[string]interface{}{
			"dns_zone_id":              dnsZoneId,
			"web_app_routing_identity": webAppRoutingIdentity,
		},
	}
}

func expandKubernetesClusterAzureMonitorProfile(input []interface{}) *managedclusters.ManagedClusterAzureMonitorProfile {
	if len(input) == 0 {
		return &managedclusters.ManagedClusterAzureMonitorProfile{
			Metrics: &managedclusters.ManagedClusterAzureMonitorProfileMetrics{
				Enabled: false,
			},
		}
	}
	if input[0] == nil {
		return &managedclusters.ManagedClusterAzureMonitorProfile{
			Metrics: &managedclusters.ManagedClusterAzureMonitorProfileMetrics{
				Enabled: true,
			},
		}
	}
	config := input[0].(map[string]interface{})
	return &managedclusters.ManagedClusterAzureMonitorProfile{
		Metrics: &managedclusters.ManagedClusterAzureMonitorProfileMetrics{
			Enabled: true,
			KubeStateMetrics: &managedclusters.ManagedClusterAzureMonitorProfileKubeStateMetrics{
				MetricAnnotationsAllowList: utils.String(config["annotations_allowed"].(string)),
				MetricLabelsAllowlist:      utils.String(config["labels_allowed"].(string)),
			},
		},
	}
}

func flattenKubernetesClusterAzureServiceMeshProfile(input *managedclusters.ServiceMeshProfile) []interface{} {
	if input == nil || input.Mode != managedclusters.ServiceMeshModeIstio {
		return nil
	}

	returnMap := map[string]interface{}{
		"mode": string(managedclusters.ServiceMeshModeIstio),
	}

	if (input.Istio.Components.IngressGateways != nil) && len(*input.Istio.Components.IngressGateways) > 0 {

		for _, value := range *input.Istio.Components.IngressGateways {

			mode := value.Mode
			enabled := value.Enabled

			var currentIngressKey string

			if mode == managedclusters.IstioIngressGatewayModeExternal {
				currentIngressKey = "external_ingress_gateway_enabled"
			} else {
				currentIngressKey = "internal_ingress_gateway_enabled"
			}

			returnMap[currentIngressKey] = enabled
		}
	}

	return []interface{}{returnMap}
}

func flattenKubernetesClusterAzureMonitorProfile(input *managedclusters.ManagedClusterAzureMonitorProfile) []interface{} {
	if input == nil || input.Metrics == nil || !input.Metrics.Enabled {
		return nil
	}
	if input.Metrics.KubeStateMetrics == nil {
		return []interface{}{
			map[string]interface{}{},
		}
	}
	annotationAllowList := ""
	if input.Metrics.KubeStateMetrics.MetricAnnotationsAllowList != nil {
		annotationAllowList = *input.Metrics.KubeStateMetrics.MetricAnnotationsAllowList
	}
	labelAllowList := ""
	if input.Metrics.KubeStateMetrics.MetricLabelsAllowlist != nil {
		labelAllowList = *input.Metrics.KubeStateMetrics.MetricLabelsAllowlist
	}
	return []interface{}{
		map[string]interface{}{
			"annotations_allowed": annotationAllowList,
			"labels_allowed":      labelAllowList,
		},
	}
}

func retrySystemNodePoolCreation(ctx context.Context, client *agentpools.AgentPoolsClient, id agentpools.AgentPoolId, profile agentpools.AgentPool) error {
	// retries the creation of a system node pool 3 times
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		if err = client.CreateOrUpdateThenPoll(ctx, id, profile); err == nil {
			return nil
		}
	}

	return err
}

func convertCustomCaTrustCertsInput(input []interface{}) *[]string {
	if len(input) == 0 {
		return nil
	}

	customCaTrustCertList := make([]string, 0)

	for _, value := range input {
		customCaTrustCertList = append(customCaTrustCertList, value.(string))
	}

	return &customCaTrustCertList

}
