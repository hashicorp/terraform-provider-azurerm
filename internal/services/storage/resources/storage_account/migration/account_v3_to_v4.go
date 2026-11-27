// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AccountV3ToV4 struct{}

func (AccountV3ToV4) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"access_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},
		"account_kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"account_replication_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"account_tier": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"allow_nested_items_to_be_public": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"allowed_copy_scope": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"azure_files_authentication": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"active_directory": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"domain_guid": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"domain_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"domain_sid": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"forest_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"netbios_domain_name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"storage_sid": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					}},
				},
				"directory_type": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
			}},
		},
		"blob_properties": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"change_feed_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
				"change_feed_retention_in_days": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
				},
				"container_delete_retention_policy": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{"days": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					}}},
				},
				"cors_rule": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"allowed_headers": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"allowed_methods": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"allowed_origins": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"exposed_headers": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"max_age_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},
					}},
				},
				"default_service_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
				},
				"delete_retention_policy": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{"days": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					}}},
				},
				"last_access_time_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
				"restore_policy": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{"days": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					}}},
				},
				"versioning_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			}},
		},
		"cross_tenant_replication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"custom_domain": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"use_subdomain": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			}},
		},
		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"key_vault_key_id": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"user_assigned_identity_id": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
			}},
		},
		"default_to_oauth_authentication": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"edge_zone": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"enable_https_traffic_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"identity_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},
				"principal_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"tenant_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
			}},
		},
		"immutability_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"allow_protected_append_writes": {
					Type:     pluginsdk.TypeBool,
					Required: true,
				},
				"period_since_creation_in_days": {
					Type:     pluginsdk.TypeInt,
					Required: true,
				},
				"state": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
			}},
		},
		"infrastructure_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"is_hns_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"large_file_share_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},
		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"min_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"network_rules": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"bypass": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Computed: true,
					Set:      pluginsdk.HashString,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},
				"default_action": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"ip_rules": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Computed: true,
					Set:      pluginsdk.HashString,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},
				"private_link_access": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"endpoint_resource_id": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"endpoint_tenant_id": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
						},
					}},
				},
				"virtual_network_subnet_ids": {
					Type:     pluginsdk.TypeSet,
					Optional: true,
					Computed: true,
					Set:      pluginsdk.HashString,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				},
			}},
		},
		"nfsv3_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"primary_access_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_blob_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_blob_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_blob_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_dfs_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_dfs_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_file_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_file_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_queue_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_queue_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_table_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_table_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_web_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"primary_web_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"queue_encryption_key_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"queue_properties": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"cors_rule": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"allowed_headers": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"allowed_methods": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"allowed_origins": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"exposed_headers": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"max_age_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},
					}},
				},
				"hour_metrics": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
						"include_apis": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
						"retention_policy_days": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},
						"version": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					}},
				},
				"logging": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"delete": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
						"read": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
						"retention_policy_days": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},
						"version": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"write": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
					}},
				},
				"minute_metrics": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
						"include_apis": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
						"retention_policy_days": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},
						"version": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					}},
				},
			}},
		},
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"routing": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"choice": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"publish_internet_endpoints": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
				"publish_microsoft_endpoints": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
				},
			}},
		},
		"sas_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"expiration_action": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"expiration_period": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
			}},
		},
		"secondary_access_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_blob_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_blob_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_blob_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_dfs_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_dfs_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_file_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_file_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_queue_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_queue_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_table_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_table_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_web_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"secondary_web_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"sftp_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"share_properties": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"cors_rule": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"allowed_headers": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"allowed_methods": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"allowed_origins": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"exposed_headers": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"max_age_in_seconds": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},
					}},
				},
				"retention_policy": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{"days": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					}}},
				},
				"smb": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"authentication_types": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"channel_encryption_type": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"kerberos_ticket_encryption_type": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"multichannel_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
						"versions": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
					}},
				},
			}},
		},
		"shared_access_key_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"static_website": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"error_404_document": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
				"index_document": {
					Type:     pluginsdk.TypeString,
					Optional: true,
				},
			}},
		},
		"table_encryption_key_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
	}
}

func (AccountV3ToV4) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
		if _, ok := rawState["cross_tenant_replication_enabled"]; !ok {
			rawState["cross_tenant_replication_enabled"] = true
		}

		return rawState, nil
	}
}
