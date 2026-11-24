// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_cloud_vm_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudexadatainfrastructures"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.CloudVMClusterName,
			ForceNew:     true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// Required
		"cloud_exadata_infrastructure_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cloudexadatainfrastructures.ValidateCloudExadataInfrastructureID,
		},

		"cpu_core_count": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(2),
		},

		"data_storage_size_in_tbs": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validate.DataStorageSizeInTbs,
		},

		"db_node_storage_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
		},

		"db_servers": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CloudVMClusterName,
		},

		"gi_version": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"hostname": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"license_model": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LicenseModel,
		},

		"memory_size_in_gbs": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
		},

		"ssh_public_keys": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		// Optional
		"backup_subnet_cidr": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsCIDR,
		},

		"cluster_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
		},

		"data_collection_options": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"diagnostics_events_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						ForceNew: true,
					},

					"health_monitoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						ForceNew: true,
					},

					"incident_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						ForceNew: true,
					},
				},
			},
		},

		"data_storage_percentage": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.IntInSlice([]int{35, 40, 60, 80}),
		},

		"domain": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
		},

		"local_backup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
		},

		"sparse_diskgroup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
		},

		"scan_listener_port_tcp": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1521,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1024, 8999),
		},

		"scan_listener_port_tcp_ssl": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      2484,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1024, 8999),
		},

		"system_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validate.SystemVersion,
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
		},

		"zone_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
		},

		"tags": commonschema.Tags(),

		"file_system_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mount_point": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"size_in_gb": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"hostname_actual": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ocid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
