// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package service_fabric_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicefabric/2021-06-01/cluster"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicefabric/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func serviceFabricClusterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"reliability_level": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(cluster.PossibleValuesForReliabilityLevel(), false),
		},

		"upgrade_mode": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(cluster.PossibleValuesForUpgradeMode(), false),
		},

		"service_fabric_zonal_upgrade_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(cluster.PossibleValuesForSfZonalUpgradeMode(), false),
		},

		"vmss_zonal_upgrade_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(cluster.PossibleValuesForVMSSZonalUpgradeMode(), false),
		},

		"cluster_code_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"management_endpoint": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"vm_image": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"add_on_features": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Set:      pluginsdk.HashString,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"azure_active_directory": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"cluster_application_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"client_application_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"certificate": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"certificate_common_names"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"thumbprint": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"thumbprint_secondary": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"x509_store_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"certificate_common_names": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"certificate"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"common_names": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"certificate_common_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"certificate_issuer_thumbprint": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"x509_store_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"reverse_proxy_certificate": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"reverse_proxy_certificate_common_names"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"thumbprint": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"thumbprint_secondary": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"x509_store_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"reverse_proxy_certificate_common_names": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"reverse_proxy_certificate"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"common_names": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"certificate_common_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"certificate_issuer_thumbprint": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"x509_store_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"client_certificate_thumbprint": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"thumbprint": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"is_admin": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},

		"client_certificate_common_name": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"common_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"issuer_thumbprint": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						// todo remove this when https://github.com/Azure/azure-sdk-for-go/issues/17744 is fixed
						DiffSuppressFunc: suppress.CaseDifference,
					},
					"is_admin": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},

		"diagnostics_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"storage_account_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"protected_account_key_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"blob_endpoint": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"queue_endpoint": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"table_endpoint": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"upgrade_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"force_restart_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"health_check_retry_timeout": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "00:45:00",
						ValidateFunc: validate.UpgradeTimeout,
					},
					"health_check_stable_duration": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "00:01:00",
						ValidateFunc: validate.UpgradeTimeout,
					},
					"health_check_wait_duration": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "00:00:30",
						ValidateFunc: validate.UpgradeTimeout,
					},
					"upgrade_domain_timeout": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "02:00:00",
						ValidateFunc: validate.UpgradeTimeout,
					},
					"upgrade_replica_set_check_timeout": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "10675199.02:48:05.4775807",
						ValidateFunc: validate.UpgradeTimeout,
					},
					"upgrade_timeout": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "12:00:00",
						ValidateFunc: validate.UpgradeTimeout,
					},
					"health_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max_unhealthy_applications_percent": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntBetween(0, 100),
								},
								"max_unhealthy_nodes_percent": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntBetween(0, 100),
								},
							},
						},
					},
					"delta_health_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"max_delta_unhealthy_applications_percent": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntBetween(0, 100),
								},
								"max_delta_unhealthy_nodes_percent": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntBetween(0, 100),
								},
								"max_upgrade_domain_delta_unhealthy_nodes_percent": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntBetween(0, 100),
								},
							},
						},
					},
				},
			},
		},

		"fabric_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"parameters": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"node_type": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"placement_properties": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"capacities": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"instance_count": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"is_primary": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"is_stateless": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"multiple_availability_zones": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"client_endpoint_port": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"http_endpoint_port": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"reverse_proxy_endpoint_port": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IsPortNumber,
					},
					"durability_level": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(cluster.DurabilityLevelBronze),
						ValidateFunc: validation.StringInSlice(cluster.PossibleValuesForDurabilityLevel(), false),
					},

					"application_ports": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"start_port": {
									Type:     pluginsdk.TypeInt,
									Required: true,
								},
								"end_port": {
									Type:     pluginsdk.TypeInt,
									Required: true,
								},
							},
						},
					},

					"ephemeral_ports": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"start_port": {
									Type:     pluginsdk.TypeInt,
									Required: true,
								},
								"end_port": {
									Type:     pluginsdk.TypeInt,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"cluster_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
