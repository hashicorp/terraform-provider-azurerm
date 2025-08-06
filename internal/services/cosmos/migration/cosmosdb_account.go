// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = CosmosDBAccountV0toV1{}

type CosmosDBAccountV0toV1 struct{}

func (c CosmosDBAccountV0toV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"offer_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"analytical_storage": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"schema_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"capacity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"total_throughput_limit": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		"minimal_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"create_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"default_identity_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"ip_range_filter": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"free_tier_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"analytical_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"automatic_failover_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"key_vault_key_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"consistency_policy": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"consistency_level": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"max_interval_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"max_staleness_prefix": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"geo_location": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"location": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"failover_priority": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"zone_redundant": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"capabilities": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"is_virtual_network_filter_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"virtual_network_rule": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"ignore_missing_vnet_service_endpoint": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"access_key_metadata_writes_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"local_authentication_disabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"mongo_server_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"multiple_write_locations_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"network_acl_bypass_for_azure_services": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"network_acl_bypass_ids": {
			Type:     pluginsdk.TypeList,
			Optional: true, Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"partition_merge_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"burst_capacity_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"backup": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"tier": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"interval_in_minutes": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},

					"retention_in_hours": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},

					"storage_redundancy": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"identity_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"cors_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed_origins": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"exposed_headers": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"allowed_headers": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"allowed_methods": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"max_age_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"connection_strings": {
			Type:      pluginsdk.TypeList,
			Computed:  true,
			Sensitive: true,
			Elem: &pluginsdk.Schema{
				Type:      pluginsdk.TypeString,
				Sensitive: true,
			},
		},

		"enable_multiple_write_locations": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"enable_free_tier": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"enable_automatic_failover": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"restore": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"source_cosmosdb_account_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},

					"restore_timestamp_in_utc": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"database": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"collection_names": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"gremlin_database": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},

								"graph_names": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},

					"tables_to_restore": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"read_endpoints": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"write_endpoints": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"primary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_readonly_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_readonly_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_sql_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_sql_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_readonly_sql_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_readonly_sql_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_mongodb_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_mongodb_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_readonly_mongodb_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_readonly_mongodb_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (c CosmosDBAccountV0toV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		if ipfString, ok := rawState["ip_range_filter"].(string); ok {
			ipfSet := make([]string, 0)
			if ipfString != "" {
				ipfSet = strings.Split(ipfString, ",")
			}
			rawState["ip_range_filter"] = ipfSet
		}

		return rawState, nil
	}
}
