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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FlexibleServerName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"administrator_login": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validate.FlexibleServerAdministratorLogin,
		},

		"administrator_password": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validate.FlexibleServerAdministratorPassword,
			ConflictsWith: []string{"administrator_password_wo"},
		},

		"administrator_password_wo": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			WriteOnly:     true,
			ConflictsWith: []string{"administrator_password"},
			RequiredWith:  []string{"administrator_password_wo_version"},
			ValidateFunc:  validate.FlexibleServerAdministratorPassword,
		},

		"administrator_password_wo_version": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"administrator_password_wo"},
		},

		"backup_retention_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      7,
			ValidateFunc: validation.IntBetween(1, 35),
		},

		"create_mode": {
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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
						RequiredWith: []string{
							"identity",
							"customer_managed_key.0.primary_user_assigned_identity_id",
						},
					},
					"primary_user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
					"geo_backup_key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
						RequiredWith: []string{
							"identity",
							"customer_managed_key.0.geo_backup_user_assigned_identity_id",
						},
					},
					"geo_backup_user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},

		"delegated_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"geo_redundant_backup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"high_availability": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Type:     pluginsdk.TypeString,
						Required: true,
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
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 6),
					},

					"start_hour": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 23),
					},

					"start_minute": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntBetween(0, 59),
					},
				},
			},
		},

		"point_in_time_restore_time_in_utc": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"private_dns_zone_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: privatezones.ValidatePrivateDnsZoneID,
		},

		"public_network_access": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C because Azure normally defaults this to `Enabled` unless values are provided for `delegated_subnet_id` and `private_dns_zone_id`
			Computed:     true,
			ValidateFunc: validation.StringInSlice(servers.PossibleValuesForEnableStatusEnum(), false),
		},

		"replication_role": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				string(servers.ReplicationRoleNone),
			}, false),
		},

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.FlexibleServerSkuName,
		},

		"source_server_id": {
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
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"iops": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(360, 48000),
					},

					"log_on_disk_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"size_gb": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.IntBetween(20, 16384),
					},
					"io_scaling_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
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
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"replica_capacity": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
