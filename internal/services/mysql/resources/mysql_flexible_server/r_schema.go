// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func mysqlFlexibleServerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this MySQL Flexible Server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FlexibleServerName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"administrator_login": {
			Description:  "The Administrator login for the MySQL Flexible Server.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validate.FlexibleServerAdministratorLogin,
		},

		"administrator_password": {
			Description:   "The Password associated with the `administrator_login` for the MySQL Flexible Server.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validate.FlexibleServerAdministratorPassword,
			ConflictsWith: []string{"administrator_password_wo"},
		},

		"administrator_password_wo": {
			Description:   "The Password associated with the `administrator_login` for the MySQL Flexible Server.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			WriteOnly:     true,
			ConflictsWith: []string{"administrator_password"},
			RequiredWith:  []string{"administrator_password_wo_version"},
			ValidateFunc:  validate.FlexibleServerAdministratorPassword,
		},

		"administrator_password_wo_version": {
			Description:  "An integer value used to trigger an update for `administrator_password_wo`. This property should be incremented when updating `administrator_password_wo`.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"administrator_password_wo"},
		},

		"backup_retention_days": {
			Description:  "The backup retention days for the MySQL Flexible Server.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      7,
			ValidateFunc: validation.IntBetween(1, 35),
		},

		"create_mode": {
			Description:  "The creation mode which can be used to restore or replicate existing servers.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(servers.PossibleValuesForCreateMode(), false),
		},

		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Description:  "The ID of the Key Vault Key.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
						RequiredWith: []string{
							"identity",
							"customer_managed_key.0.primary_user_assigned_identity_id",
						},
					},
					"primary_user_assigned_identity_id": {
						Description:  "The primary user Managed Identity ID for a Customer Managed Key.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
					"geo_backup_key_vault_key_id": {
						Description:  "The ID of the geo backup Key Vault Key. It can't cross region and need Customer Managed Key in same region as geo backup.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
						RequiredWith: []string{
							"identity",
							"customer_managed_key.0.geo_backup_user_assigned_identity_id",
						},
					},
					"geo_backup_user_assigned_identity_id": {
						Description:  "The geo backup user Managed Identity ID for a Customer Managed Key. It can't cross region and need identity in same region as geo backup.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},

		"delegated_subnet_id": {
			Description:  "The ID of the Virtual Network subnet to create the MySQL Flexible Server.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"geo_redundant_backup_enabled": {
			Description: "Should geo redundant backup be enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},

		"high_availability": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Description: "The high availability mode for the MySQL Flexible Server.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(servers.HighAvailabilityModeZoneRedundant),
							string(servers.HighAvailabilityModeSameZone),
						}, false),
					},

					"standby_availability_zone": commonschema.ZoneSingleOptionalComputed(),
				},
			},
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"maintenance_window": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"day_of_week": {
						Description:  "The day of week for maintenance window.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 6),
					},

					"start_hour": {
						Description:  "The start hour for maintenance window.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 23),
					},

					"start_minute": {
						Description:  "The start minute for maintenance window.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 59),
					},
				},
			},
		},

		"point_in_time_restore_time_in_utc": {
			Description:  "The point in time to restore from `creation_source_server_id` when `create_mode` is `PointInTimeRestore`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"private_dns_zone_id": {
			Description:  "The ID of the private DNS zone to create the MySQL Flexible Server.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: privatezones.ValidatePrivateDnsZoneID,
		},

		"public_network_access": {
			Description: "Whether approved public traffic is allowed through the firewall to this server.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C because Azure normally defaults this to `Enabled` unless values are provided for `delegated_subnet_id` and `private_dns_zone_id`
			Computed:     true,
			ValidateFunc: validation.StringInSlice(servers.PossibleValuesForEnableStatusEnum(), false),
		},

		"replication_role": {
			Description: "The replication role.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				string(servers.ReplicationRoleNone),
			}, false),
		},

		"sku_name": {
			Description:  "The SKU Name for the MySQL Flexible Server.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.FlexibleServerSkuName,
		},

		"source_server_id": {
			Description:  "The resource ID of the source MySQL Flexible Server to be restored.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: servers.ValidateFlexibleServerID,
		},

		"storage": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"auto_grow_enabled": {
						Description: "Should Storage Auto Grow be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"iops": {
						Description:  "The storage IOPS for the MySQL Flexible Server.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(360, 48000),
					},

					"log_on_disk_enabled": {
						Description: "Should Storage Log On Disk be enabled?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"size_gb": {
						Description:  "The max storage allowed for the MySQL Flexible Server.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(20, 16384),
					},
					"io_scaling_enabled": {
						Description: "Should IOPS be scaled automatically? If `true`, `iops` can not be set.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"version": {
			Description: "The version of the MySQL Flexible Server to use.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because Azure assigns a version when not explicitly set
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(servers.ServerVersionFivePointSeven),
				string(servers.ServerVersionEightPointZeroPointTwoOne),
				string(validate.ServerVersionEightPointFour),
			}, false),
		},

		"zone": commonschema.ZoneSingleOptionalComputed(),

		"fqdn": {
			Description: "The fully qualified domain name of the MySQL Flexible Server.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"replica_capacity": {
			Description: "The maximum number of replicas that a primary MySQL Flexible Server can have.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
