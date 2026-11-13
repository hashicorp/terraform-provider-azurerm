// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server

import (
	"math"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func postgresqlFlexibleServerSchema() map[string]*pluginsdk.Schema {
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
			ValidateFunc: validation.All(validation.StringIsNotWhiteSpace, validate.AdminUsernames),
		},

		"administrator_password": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"administrator_password_wo"},
		},

		"administrator_password_wo": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			WriteOnly:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"administrator_password"},
			RequiredWith:  []string{"administrator_password_wo_version"},
		},

		"administrator_password_wo_version": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"administrator_password_wo"},
		},

		"authentication": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"active_directory_auth_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"password_auth_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
						RequiredWith: []string{
							"authentication.0.active_directory_auth_enabled",
						},
					},
				},
			},
		},

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.FlexibleServerSkuName,
		},

		"auto_grow_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"storage_mb": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(int(math.Exp2(15)), int(math.Exp2(26))),
		},

		"storage_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				string(servers.AzureManagedDiskPerformanceTierPFour),
				string(servers.AzureManagedDiskPerformanceTierPSix),
				string(servers.AzureManagedDiskPerformanceTierPOneZero),
				string(servers.AzureManagedDiskPerformanceTierPOneFive),
				string(servers.AzureManagedDiskPerformanceTierPTwoZero),
				string(servers.AzureManagedDiskPerformanceTierPThreeZero),
				string(servers.AzureManagedDiskPerformanceTierPFourZero),
				string(servers.AzureManagedDiskPerformanceTierPFiveZero),
				string(servers.AzureManagedDiskPerformanceTierPSixZero),
				string(servers.AzureManagedDiskPerformanceTierPSevenZero),
				string(servers.AzureManagedDiskPerformanceTierPEightZero),
			}, false),
		},

		"storage_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(servers.StorageTypePremiumLRS),
			ForceNew: true,
			// NOTE: not using `servers.PossibleValuesForStorageType()` because it includes `UltraSSD_LRS`, which is not GA yet
			ValidateFunc: validation.StringInSlice([]string{
				string(servers.StorageTypePremiumLRS),
				string(servers.StorageTypePremiumVTwoLRS),
			}, false),
		},

		"storage_iops": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C Azure computes IOPs for Premium_LRS or from source server
			Computed:     true,
			ValidateFunc: validation.IntBetween(3000, 80000),
		},

		"storage_throughput": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C Azure computes throughput for Premium_LRS or from source server
			Computed:     true,
			ValidateFunc: validation.IntBetween(125, 1200),
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C because Azure assigns a version when not explicitly set
			Computed:     true,
			ValidateFunc: validation.StringInSlice(servers.PossibleValuesForPostgresMajorVersion(), false),
		},

		"zone": commonschema.ZoneSingleOptional(),

		"create_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(servers.CreateModeDefault),
				string(servers.CreateModePointInTimeRestore),
				string(servers.CreateModeReplica),
				string(servers.CreateModeReviveDropped),
				string(servers.CreateModeGeoRestore),
				string(servers.CreateModeUpdate),
			}, false),
		},

		"delegated_subnet_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"private_dns_zone_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			// todo make this case sensitive when https://github.com/Azure/azure-rest-api-specs/issues/26346 is fixed
			DiffSuppressFunc: suppress.CaseDifference,
			// This is `computed`, because there is a breaking change to require this field when setting vnet.
			// For existing fs who don't want to be recreated, they could contact service team to manually migrate to the private dns zone
			// We need to ignore the diff when remote is set private dns zone
			ValidateFunc: privatezones.ValidatePrivateDnsZoneID,
		},

		"point_in_time_restore_time_in_utc": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"source_server_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: servers.ValidateFlexibleServerID,
		},

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

		"backup_retention_days": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(7, 35),
		},

		"geo_redundant_backup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},

		"high_availability": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(servers.PossibleValuesForHighAvailabilityMode(), false),
					},

					"standby_availability_zone": commonschema.ZoneSingleOptional(),
				},
			},
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"replication_role": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(servers.ReplicationRoleNone),
			}, false),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
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

		"cluster": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"size": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(1, 20),
					},
					"default_database_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						Default:      "postgres",
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
