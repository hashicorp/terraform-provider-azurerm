// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server

import (
	"math"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2026-10-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	postgresqlVersion11 = "11"
	postgresqlVersion12 = "12"
	postgresqlVersion13 = "13"
	postgresqlVersion14 = "14"
	postgresqlVersion15 = "15"
	postgresqlVersion16 = "16"
	postgresqlVersion17 = "17"
	postgresqlVersion18 = "18"
)

func postgresqlFlexibleServerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this PostgreSQL Flexible Server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FlexibleServerName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"administrator_login": {
			Description:  "The Administrator login for the PostgreSQL Flexible Server.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.All(validation.StringIsNotWhiteSpace, validate.AdminUsernames),
		},

		"administrator_password": {
			Description:   "The Password associated with the `administrator_login` for the PostgreSQL Flexible Server.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"administrator_password_wo"},
		},

		"administrator_password_wo": {
			Description:   "The Password associated with the `administrator_login` for the PostgreSQL Flexible Server.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			WriteOnly:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			RequiredWith:  []string{"administrator_password_wo_version"},
			ConflictsWith: []string{"administrator_password"},
		},

		"administrator_password_wo_version": {
			Description:  "An integer value used to trigger an update for `administrator_password_wo`. This property should be incremented when updating `administrator_password_wo`.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"administrator_password_wo"},
		},

		"authentication": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"active_directory_auth_enabled": {
						Description: "Whether Active Directory authentication is allowed to access the PostgreSQL Flexible Server.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"password_auth_enabled": {
						Description: "Whether password authentication is allowed to access the PostgreSQL Flexible Server.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"tenant_id": {
						Description:  "The Tenant ID of the Azure Active Directory which is used by the Active Directory authentication. `active_directory_auth_enabled` must be set to `true`.",
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
			Description:  "The SKU Name for the PostgreSQL Flexible Server. The name of the SKU, follows the `tier` + `name` pattern (e.g. `B_Standard_B1ms`, `B_Standard_B2ats_v2`, `GP_Standard_D2s_v3`, `MO_Standard_E4s_v3`).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validate.FlexibleServerSkuName,
		},

		"auto_grow_enabled": {
			Description: "Whether the storage auto grow for PostgreSQL Flexible Server is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"storage_mb": {
			Description:  "The max storage allowed for the PostgreSQL Flexible Server.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(int(math.Exp2(15)), int(math.Exp2(26))),
		},

		"storage_tier": {
			Description: "The name of storage performance tier for IOPS of the PostgreSQL Flexible Server. See the `storage_tier` defaults based on `storage_mb` table below.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
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
			Description: "The type of storage used for the PostgreSQL Flexible Server.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true, Default: string(servers.StorageTypePremiumLRS),
			// NOTE: not using `servers.PossibleValuesForStorageType()` because it includes `UltraSSD_LRS`, which is not GA yet
			ValidateFunc: validation.StringInSlice([]string{
				string(servers.StorageTypePremiumLRS),
				string(servers.StorageTypePremiumVTwoLRS),
			}, false),
		},

		"storage_iops": {
			Description: "The maximum IOPS supported for storage.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// NOTE: O+C Azure computes IOPs for Premium_LRS or from source server
			Computed:     true,
			ValidateFunc: validation.IntBetween(3000, 80000),
		},

		"storage_throughput": {
			Description: "The maximum throughput supported for storage in MB/s.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			// NOTE: O+C Azure computes throughput for Premium_LRS or from source server
			Computed:     true,
			ValidateFunc: validation.IntBetween(125, 1200),
		},

		"version": {
			Description: "The version of PostgreSQL Flexible Server to use.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because Azure assigns a version when not explicitly set
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				postgresqlVersion11,
				postgresqlVersion12,
				postgresqlVersion13,
				postgresqlVersion14,
				postgresqlVersion15,
				postgresqlVersion16,
				postgresqlVersion17,
				postgresqlVersion18,
			}, false),
		},

		"zone": commonschema.ZoneSingleOptional(),

		"create_mode": {
			Description: "The creation mode which can be used to restore or replicate existing servers.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
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
			Description:  "The ID of the Virtual Network subnet to create the PostgreSQL Flexible Server. The provided subnet should not have any other resource deployed in it and this subnet will be delegated to the PostgreSQL Flexible Server, if not already delegated.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"private_dns_zone_id": {
			Description: "The ID of the private DNS zone to create the PostgreSQL Flexible Server.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			// This is `computed`, because there is a breaking change to require this field when setting vnet.
			// For existing fs who don't want to be recreated, they could contact service team to manually migrate to the private dns zone
			// We need to ignore the diff when remote is set private dns zone
			ValidateFunc: privatezones.ValidatePrivateDnsZoneID,
			// todo make this case sensitive when https://github.com/Azure/azure-rest-api-specs/issues/26346 is fixed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"point_in_time_restore_time_in_utc": {
			Description:  "The point in time to restore from `source_server_id` when `create_mode` is `GeoRestore`, `PointInTimeRestore`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"source_server_id": {
			Description:  "The resource ID of the source PostgreSQL Flexible Server to be restored.",
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
						Description:  "The day of week for maintenance window, where the week starts on a Sunday, i.e. Sunday = `0`, Monday = `1`.",
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

		"backup_retention_days": {
			Description:  "The backup retention days for the PostgreSQL Flexible Server.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(7, 35),
		},

		"geo_redundant_backup_enabled": {
			Description: "Whether Geo-Redundant backup is enabled on the PostgreSQL Flexible Server.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true, Default: false,
		},

		"high_availability": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"mode": {
						Description:  "The high availability mode for the PostgreSQL Flexible Server.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(servers.PossibleValuesForHighAvailabilityMode(), false),
					},

					"standby_availability_zone": commonschema.ZoneSingleOptional(),
				},
			},
		},

		"fqdn": {
			Description: "The FQDN of the PostgreSQL Flexible Server.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"public_network_access_enabled": {
			Description: "Whether this PostgreSQL Flexible Server is publicly accessible.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"replication_role": {
			Description: "The replication role for the PostgreSQL Flexible Server.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(servers.ReplicationRoleNone),
			}, false),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Description:  "The versioned/versionless ID of the Key Vault Key.",
						Type:         pluginsdk.TypeString,
						Required:     true,
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
						Description:  "The versioned/versionless ID of the geo backup Key Vault Key.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
						RequiredWith: []string{
							"identity",
							"customer_managed_key.0.geo_backup_user_assigned_identity_id",
						},
					},
					"geo_backup_user_assigned_identity_id": {
						Description:  "The geo backup user Managed Identity ID for a Customer Managed Key.",
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
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"size": {
						Description:  "The number of nodes in the cluster.",
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
