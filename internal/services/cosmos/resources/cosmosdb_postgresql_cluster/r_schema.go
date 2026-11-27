// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_postgresql_cluster

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresqlhsc/2022-11-08/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 260),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"node_count": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ValidateFunc: validation.All(
				validation.IntBetween(0, 20),
				validation.IntNotInSlice([]int{1}),
			),
		},

		"administrator_login_password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringLenBetween(8, 256),
		},

		"citus_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				"8.3",
				"9.0",
				"9.1",
				"9.2",
				"9.3",
				"9.4",
				"9.5",
				"10.0",
				"10.1",
				"10.2",
				"11.0",
				"11.1",
				"11.2",
				"11.3",
				"12.1",
			}, false),
		},

		"coordinator_public_ip_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		// Once the issue https://github.com/Azure/azure-rest-api-specs/issues/23317 is fixed, we would submit PR to improve the validation
		"coordinator_server_edition": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "GeneralPurpose",
			ValidateFunc: validation.StringInSlice([]string{
				"BurstableGeneralPurpose",
				"BurstableMemoryOptimized",
				"GeneralPurpose",
				"MemoryOptimized",
			}, false),
		},

		"coordinator_storage_quota_in_mb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ValidateFunc: validation.All(
				validation.IntBetween(32768, 16777216),
				validation.IntDivisibleBy(1024),
			),
		},

		"coordinator_vcore_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ValidateFunc: validation.IntInSlice([]int{
				1,
				2,
				4,
				8,
				16,
				32,
				64,
				96,
			}),
		},

		"ha_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
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

		"node_public_ip_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		// Once the issue https://github.com/Azure/azure-rest-api-specs/issues/23317 is fixed, we would submit PR to improve the validation
		"node_server_edition": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "MemoryOptimized",
			ValidateFunc: validation.StringInSlice([]string{
				"BurstableGeneralPurpose",
				"BurstableMemoryOptimized",
				"GeneralPurpose",
				"MemoryOptimized",
			}, false),
		},

		"node_storage_quota_in_mb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.All(
				validation.IntBetween(32768, 16777216),
				validation.IntDivisibleBy(1024),
			),
		},

		"node_vcores": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntInSlice([]int{
				1,
				2,
				4,
				8,
				16,
				32,
				64,
				96,
				104,
			}),
		},

		"point_in_time_in_utc": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsRFC3339Time,
			RequiredWith: []string{"source_location", "source_resource_id"},
		},

		"preferred_primary_zone": commonschema.ZoneSingleOptional(),

		"shards_on_coordinator_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"source_location": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
			RequiredWith:     []string{"source_resource_id"},
		},

		"source_resource_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: clusters.ValidateServerGroupsv2ID,
			RequiredWith: []string{"source_location"},
		},

		"sql_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				"11",
				"12",
				"13",
				"14",
				"15",
				"16",
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"earliest_restore_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"servers": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"fqdn": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}
