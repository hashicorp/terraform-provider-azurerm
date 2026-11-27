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
			Description: "The name of the Service Fabric Cluster.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"reliability_level": {
			Description:  "The Reliability Level of the Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(cluster.PossibleValuesForReliabilityLevel(), false),
		},

		"upgrade_mode": {
			Description:  "The Upgrade Mode of the cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(cluster.PossibleValuesForUpgradeMode(), false),
		},

		"service_fabric_zonal_upgrade_mode": {
			Description:  "The logical grouping of VMs in upgrade domains.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(cluster.PossibleValuesForSfZonalUpgradeMode(), false),
		},

		"vmss_zonal_upgrade_mode": {
			Description:  "The upgrade mode for the Virtual Machine Scale Set updates that happen in all Availability Zones at once.",
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
			Description: "The Management Endpoint of the cluster such as `http://example.com`.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"vm_image": {
			Description: "The Image expected for the Service Fabric Cluster, such as `Windows`.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"add_on_features": {
			Description: "A list of one or more features which should be enabled, such as `DnsService`.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Set:         pluginsdk.HashString,
			Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"azure_active_directory": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tenant_id": {
						Description:  "The Azure Active Directory Tenant ID.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"cluster_application_id": {
						Description:  "The Azure Active Directory Cluster Application ID.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"client_application_id": {
						Description:  "The Azure Active Directory Client ID which should be used for the Client Application.",
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
						Description: "The Thumbprint of the Certificate.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"thumbprint_secondary": {
						Description: "The Secondary Thumbprint of the Certificate.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"x509_store_name": {
						Description: "The X509 Store where the Certificate Exists, such as `My`.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
									Description:  "The common or subject name of the certificate.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"certificate_issuer_thumbprint": {
									Description:  "The Issuer Thumbprint of the Certificate.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"x509_store_name": {
						Description:  "The X509 Store where the Certificate Exists, such as `My`.",
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
						Description: "The Thumbprint of the Certificate.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"thumbprint_secondary": {
						Description: "The Secondary Thumbprint of the Certificate.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"x509_store_name": {
						Description: "The X509 Store where the Certificate Exists, such as `My`.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
									Description:  "The common or subject name of the certificate.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"certificate_issuer_thumbprint": {
									Description:  "The Issuer Thumbprint of the Certificate.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"x509_store_name": {
						Description:  "The X509 Store where the Certificate Exists, such as `My`.",
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
						Description: "The Thumbprint associated with the Client Certificate.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"is_admin": {
						Description: "Whether the Client Certificate has Admin Access to the cluster. Non-admin clients can only perform read-only operations on the cluster.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
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
						Description:  "The common or subject name of the certificate.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"issuer_thumbprint": {
						Description:  "The Issuer Thumbprint of the Certificate.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						// todo remove this when https://github.com/Azure/azure-sdk-for-go/issues/17744 is fixed
						DiffSuppressFunc: suppress.CaseDifference,
					},
					"is_admin": {
						Description: "Whether the Client Certificate has Admin Access to the cluster. Non-admin clients can only perform read-only operations on the cluster.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
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
						Description: "The name of the Storage Account where the Diagnostics should be sent to.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"protected_account_key_name": {
						Description: "The protected diagnostics storage key name, such as `StorageAccountKey1`.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"blob_endpoint": {
						Description: "The Blob Endpoint of the Storage Account.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"queue_endpoint": {
						Description: "The Queue Endpoint of the Storage Account.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"table_endpoint": {
						Description: "The Table Endpoint of the Storage Account.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
						Description: "Whether to restart the Service Fabric node even if only dynamic configurations have changed.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"health_check_retry_timeout": {
						Description:  "The duration, in `hh:mm:ss` string format, after which Service Fabric retries the health check if the previous health check fails.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "00:45:00",
						ValidateFunc: validate.UpgradeTimeout,
					},
					"health_check_stable_duration": {
						Description:  "The duration, in `hh:mm:ss` string format, that Service Fabric waits to verify that the cluster is stable before it continues to the next upgrade domain or completes the upgrade. This wait duration prevents undetected changes of health right after the health check is performed.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "00:01:00",
						ValidateFunc: validate.UpgradeTimeout,
					},
					"health_check_wait_duration": {
						Description:  "The duration, in `hh:mm:ss` string format, that Service Fabric waits before it performs the initial health check after it finishes the upgrade on the upgrade domain.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "00:00:30",
						ValidateFunc: validate.UpgradeTimeout,
					},
					"upgrade_domain_timeout": {
						Description:  "The duration, in `hh:mm:ss` string format, that Service Fabric takes to upgrade a single upgrade domain. After this period, the upgrade fails.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "02:00:00",
						ValidateFunc: validate.UpgradeTimeout,
					},
					"upgrade_replica_set_check_timeout": {
						Description:  "The duration, in `hh:mm:ss` string format, that Service Fabric waits for a replica set to reconfigure into a safe state, if it is not already in a safe state, before Service Fabric proceeds with the upgrade.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "10675199.02:48:05.4775807",
						ValidateFunc: validate.UpgradeTimeout,
					},
					"upgrade_timeout": {
						Description:  "The duration, in `hh:mm:ss` string format, that Service Fabric takes for the entire upgrade. After this period, the upgrade fails.",
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
									Description:  "The maximum tolerated percentage of applications that can have aggregated health state of error. If the upgrade exceeds this percentage, the cluster is unhealthy.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntBetween(0, 100),
								},
								"max_unhealthy_nodes_percent": {
									Description:  "The maximum tolerated percentage of nodes that can have aggregated health states of error. If an upgrade exceeds this percentage, the cluster is unhealthy.",
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
									Description:  "The maximum tolerated percentage of delta unhealthy applications that can have aggregated health states of error. If the current unhealthy applications do not respect the percentage relative to the state at the beginning of the upgrade, the cluster is unhealthy.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntBetween(0, 100),
								},
								"max_delta_unhealthy_nodes_percent": {
									Description:  "The maximum tolerated percentage of delta unhealthy nodes that can have aggregated health states of error. If the current unhealthy nodes do not respect the percentage relative to the state at the beginning of the upgrade, the cluster is unhealthy.",
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      0,
									ValidateFunc: validation.IntBetween(0, 100),
								},
								"max_upgrade_domain_delta_unhealthy_nodes_percent": {
									Description:  "The maximum tolerated percentage of upgrade domain delta unhealthy nodes that can have aggregated health state of error. If there is any upgrade domain where the current unhealthy nodes do not respect the percentage relative to the state at the beginning of the upgrade, the cluster is unhealthy.",
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
						Description: "The name of the Fabric Setting, such as `Security` or `Federation`.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"parameters": {
						Description: "A map containing settings for the specified Fabric Setting.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
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
						Description: "The name of the Node Type.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
					"placement_properties": {
						Description: "The placement tags applied to nodes in the node type, which can be used to indicate where certain services (workload) should run.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"capacities": {
						Description: "The capacity tags applied to the nodes in the node type, the cluster Resource Manager uses these tags to understand how much resource a node has.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"instance_count": {
						Description: "The number of nodes for this Node Type.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},
					"is_primary": {
						Description: "Whether this is the primary node type.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
					"is_stateless": {
						Description: "Whether this node type should only run stateless services.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"multiple_availability_zones": {
						Description: "Whether this node type spans across Availability Zones.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
					"client_endpoint_port": {
						Description: "The Port used for the Client Endpoint for this Node Type.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},
					"http_endpoint_port": {
						Description: "The Port used for the HTTP Endpoint for this Node Type.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},
					"reverse_proxy_endpoint_port": {
						Description:  "The Port used for the Reverse Proxy Endpoint for this Node Type.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IsPortNumber,
					},
					"durability_level": {
						Description:  "The Durability Level for this Node Type.",
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
									Description: "The start of the Application Port Range on this Node Type.",
									Type:        pluginsdk.TypeInt,
									Required:    true,
								},
								"end_port": {
									Description: "The end of the Application Port Range on this Node Type.",
									Type:        pluginsdk.TypeInt,
									Required:    true,
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
									Description: "The start of the Ephemeral Port Range on this Node Type.",
									Type:        pluginsdk.TypeInt,
									Required:    true,
								},
								"end_port": {
									Description: "The end of the Ephemeral Port Range on this Node Type.",
									Type:        pluginsdk.TypeInt,
									Required:    true,
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"cluster_endpoint": {
			Description: "The Cluster Endpoint for this Service Fabric Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
