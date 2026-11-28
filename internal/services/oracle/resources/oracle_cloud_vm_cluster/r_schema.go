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
			Description:  "The name which should be used for this Cloud VM Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CloudVMClusterName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// Required
		"cloud_exadata_infrastructure_id": {
			Description:  "The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the Cloud Exadata infrastructure.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cloudexadatainfrastructures.ValidateCloudExadataInfrastructureID,
		},

		"cpu_core_count": {
			Description:  "The number of CPU cores enabled on the Cloud VM Cluster.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(2),
		},

		"data_storage_size_in_tbs": {
			Description:  "The data disk group size to be allocated in TBs.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validate.DataStorageSizeInTbs,
		},

		"db_node_storage_size_in_gbs": {
			Description: "The local node storage to be allocated in GBs.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
		},

		"db_servers": {
			Description: "The list of DB servers.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"display_name": {
			Description:  "The user-friendly name for the Cloud VM Cluster. The name does not need to be unique.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CloudVMClusterName,
		},

		"gi_version": {
			Description: "A valid Oracle Grid Infrastructure (GI) software version.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"hostname": {
			Description: "The hostname for the Cloud VM Cluster without suffix.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"license_model": {
			Description:  "The Oracle license model that applies to the Cloud VM Cluster, either `BringYourOwnLicense` or `LicenseIncluded`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LicenseModel,
		},

		"memory_size_in_gbs": {
			Description: "The memory to be allocated in GBs.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
		},

		"ssh_public_keys": {
			Description: "The public key portion of one or more key pairs used for SSH access to the Cloud VM Cluster.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"subnet_id": {
			Description:  "The ID of the subnet associated with the Cloud VM Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"virtual_network_id": {
			Description:  "The ID of the Virtual Network associated with the Cloud VM Cluster.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		// Optional
		"backup_subnet_cidr": {
			Description:  "The backup subnet CIDR of the Virtual Network associated with the Cloud VM Cluster.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsCIDR,
		},

		"cluster_name": {
			Description: "The cluster name for Cloud VM Cluster.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
		},

		"data_collection_options": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"diagnostics_events_enabled": {
						Description: "Whether diagnostic collection is enabled for the VM Cluster/Cloud VM Cluster/VMBM DBCS. Enabling diagnostic collection allows you to receive Events service notifications for guest VM issues. Diagnostic collection also allows Oracle to provide enhanced service and proactive support for your Exadata system. You can enable diagnostic collection during VM Cluster/Cloud VM Cluster provisioning. You can also disable or enable it at any time using the `UpdateVmCluster` or `updateCloudVmCluster` API.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ForceNew:    true,
					},

					"health_monitoring_enabled": {
						Description: "Whether health monitoring is enabled for the VM Cluster / Cloud VM Cluster / VMBM DBCS. Enabling health monitoring allows Oracle to collect diagnostic data and share it with its operations and support personnel. You may also receive notifications for some events. Collecting health diagnostics enables Oracle to provide proactive support and enhanced service for your system. Optionally enable health monitoring while provisioning a system. You can also disable or enable health monitoring anytime using the `UpdateVmCluster`, `UpdateCloudVmCluster` or `updateDbsystem` API.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ForceNew:    true,
					},

					"incident_logs_enabled": {
						Description: "Whether incident logs and trace collection are enabled for the VM Cluster / Cloud VM Cluster / VMBM DBCS. Enabling incident logs collection allows Oracle to receive Events service notifications for guest VM issues, collect incident logs and traces, and use them to diagnose issues and resolve them. Optionally enable incident logs collection while provisioning a system. You can also disable or enable incident logs collection anytime using the `UpdateVmCluster`, `updateCloudVmCluster` or `updateDbsystem` API.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						ForceNew:    true,
					},
				},
			},
		},

		"data_storage_percentage": {
			Description:  "The percentage assigned to DATA storage (user data and database files). The remaining percentage is assigned to RECO storage (database redo logs, archive logs, and recovery manager backups).",
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.IntInSlice([]int{35, 40, 60, 80}),
		},

		"domain": {
			Description: "The name of the OCI Private DNS Zone to be associated with the Cloud VM Cluster.",
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
		},

		"local_backup_enabled": {
			Description: "If true, database backup on local Exadata storage is configured for the Cloud VM Cluster. If `false`, database backup on local Exadata storage is not available in the Cloud VM Cluster.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
		},

		"sparse_diskgroup_enabled": {
			Description: "If true, the sparse disk group is configured for the Cloud VM Cluster. If `false`, the sparse disk group is not created.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
		},

		"scan_listener_port_tcp": {
			Description:  "The TCP Single Client Access Name (SCAN) port.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      1521,
			ValidateFunc: validation.IntBetween(1024, 8999),
		},

		"scan_listener_port_tcp_ssl": {
			Description:  "The TCPS Single Client Access Name (SCAN) port.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      2484,
			ValidateFunc: validation.IntBetween(1024, 8999),
		},

		"system_version": {
			Description:  "Operating system version of the Exadata image. System version must be <= Db server major version (the first two parts of the DB server version e.g. 23.1.X.X.XXXX). For Grid Infrastructure (GI) version 23.0.0.0 allowed system versions are 23.1.19.0.0.241015, 23.1.20.0.0.241112, 23.1.22.0.0.250119, 23.1.21.0.0.241204, 23.1.23.0.0.250207.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validate.SystemVersion,
		},

		"time_zone": {
			Description: "The time zone of the Cloud VM Cluster. For details, see [Exadata Infrastructure Time Zones](https://docs.cloud.oracle.com/iaas/Content/Database/References/timezones.htm).",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
		},

		"zone_id": {
			Description: "The OCID of the OCI Private DNS Zone to be associated with the Cloud VM Cluster.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ForceNew:    true,
		},

		"tags": commonschema.Tags(),

		"file_system_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mount_point": {
						Description: "The mount path of the file system.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
					"size_in_gb": {
						Description: "The size of the Virtual Machine's file system.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
					},
				},
			},
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"hostname_actual": {
			Description: "The hostname for the Cloud VM Cluster with suffix.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"ocid": {
			Description: "The [OCID](https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the Cloud VM Cluster.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
