// Copyright (c) HashiCorp, Inc.
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
			Computed: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"account_kind": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"account_replication_type": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"account_tier": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"allow_nested_items_to_be_public": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"allowed_copy_scope": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"azure_files_authentication": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"active_directory": {
					Computed: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"domain_guid": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"domain_name": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"domain_sid": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"forest_name": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"netbios_domain_name": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"storage_sid": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"directory_type": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"blob_properties": {
			Computed: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"change_feed_enabled": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
				"change_feed_retention_in_days": {
					Optional: true,
					Type:     pluginsdk.TypeInt,
				},
				"container_delete_retention_policy": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{"days": {
						Optional: true,
						Type:     pluginsdk.TypeInt,
					}}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"cors_rule": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"allowed_headers": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"allowed_methods": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"allowed_origins": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"exposed_headers": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"max_age_in_seconds": {
							Required: true,
							Type:     pluginsdk.TypeInt,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"default_service_version": {
					Computed: true,
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"delete_retention_policy": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{"days": {
						Optional: true,
						Type:     pluginsdk.TypeInt,
					}}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"last_access_time_enabled": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
				"restore_policy": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{"days": {
						Required: true,
						Type:     pluginsdk.TypeInt,
					}}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"versioning_enabled": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"cross_tenant_replication_enabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"custom_domain": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"name": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"use_subdomain": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"customer_managed_key": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"key_vault_key_id": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"user_assigned_identity_id": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"default_to_oauth_authentication": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"edge_zone": {
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"enable_https_traffic_only": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"identity": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"identity_ids": {
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Optional: true,
					Type:     pluginsdk.TypeSet,
				},
				"principal_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"tenant_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"type": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"immutability_policy": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"allow_protected_append_writes": {
					Required: true,
					Type:     pluginsdk.TypeBool,
				},
				"period_since_creation_in_days": {
					Required: true,
					Type:     pluginsdk.TypeInt,
				},
				"state": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"infrastructure_encryption_enabled": {
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"is_hns_enabled": {
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"large_file_share_enabled": {
			Computed: true,
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"location": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"min_tls_version": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"network_rules": {
			Computed: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"bypass": {
					Computed: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Optional: true,
					Set:      pluginsdk.HashString,
					Type:     pluginsdk.TypeSet,
				},
				"default_action": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"ip_rules": {
					Computed: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Optional: true,
					Set:      pluginsdk.HashString,
					Type:     pluginsdk.TypeSet,
				},
				"private_link_access": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"endpoint_resource_id": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"endpoint_tenant_id": {
							Computed: true,
							Optional: true,
							Type:     pluginsdk.TypeString,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"virtual_network_subnet_ids": {
					Computed: true,
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Optional: true,
					Set:      pluginsdk.HashString,
					Type:     pluginsdk.TypeSet,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"nfsv3_enabled": {
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"primary_access_key": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_blob_connection_string": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_blob_endpoint": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_blob_host": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_connection_string": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_dfs_endpoint": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_dfs_host": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_file_endpoint": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_file_host": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_location": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_queue_endpoint": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_queue_host": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_table_endpoint": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_table_host": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_web_endpoint": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"primary_web_host": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"public_network_access_enabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"queue_encryption_key_type": {
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"queue_properties": {
			Computed: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"cors_rule": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"allowed_headers": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"allowed_methods": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"allowed_origins": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"exposed_headers": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"max_age_in_seconds": {
							Required: true,
							Type:     pluginsdk.TypeInt,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"hour_metrics": {
					Computed: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Required: true,
							Type:     pluginsdk.TypeBool,
						},
						"include_apis": {
							Optional: true,
							Type:     pluginsdk.TypeBool,
						},
						"retention_policy_days": {
							Optional: true,
							Type:     pluginsdk.TypeInt,
						},
						"version": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"logging": {
					Computed: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"delete": {
							Required: true,
							Type:     pluginsdk.TypeBool,
						},
						"read": {
							Required: true,
							Type:     pluginsdk.TypeBool,
						},
						"retention_policy_days": {
							Optional: true,
							Type:     pluginsdk.TypeInt,
						},
						"version": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
						"write": {
							Required: true,
							Type:     pluginsdk.TypeBool,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"minute_metrics": {
					Computed: true,
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Required: true,
							Type:     pluginsdk.TypeBool,
						},
						"include_apis": {
							Optional: true,
							Type:     pluginsdk.TypeBool,
						},
						"retention_policy_days": {
							Optional: true,
							Type:     pluginsdk.TypeInt,
						},
						"version": {
							Required: true,
							Type:     pluginsdk.TypeString,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"resource_group_name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"routing": {
			Computed: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"choice": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"publish_internet_endpoints": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
				"publish_microsoft_endpoints": {
					Optional: true,
					Type:     pluginsdk.TypeBool,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"sas_policy": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"expiration_action": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"expiration_period": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"secondary_access_key": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_blob_connection_string": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_blob_endpoint": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_blob_host": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_connection_string": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_dfs_endpoint": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_dfs_host": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_file_endpoint": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_file_host": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_location": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_queue_endpoint": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_queue_host": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_table_endpoint": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_table_host": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_web_endpoint": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"secondary_web_host": {
			Computed: true,
			Type:     pluginsdk.TypeString,
		},
		"sftp_enabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"share_properties": {
			Computed: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"cors_rule": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"allowed_headers": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"allowed_methods": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"allowed_origins": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"exposed_headers": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Required: true,
							Type:     pluginsdk.TypeList,
						},
						"max_age_in_seconds": {
							Required: true,
							Type:     pluginsdk.TypeInt,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"retention_policy": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{"days": {
						Optional: true,
						Type:     pluginsdk.TypeInt,
					}}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"smb": {
					Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
						"authentication_types": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Optional: true,
							Type:     pluginsdk.TypeSet,
						},
						"channel_encryption_type": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Optional: true,
							Type:     pluginsdk.TypeSet,
						},
						"kerberos_ticket_encryption_type": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Optional: true,
							Type:     pluginsdk.TypeSet,
						},
						"multichannel_enabled": {
							Optional: true,
							Type:     pluginsdk.TypeBool,
						},
						"versions": {
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Optional: true,
							Type:     pluginsdk.TypeSet,
						},
					}},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"shared_access_key_enabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"static_website": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"error_404_document": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"index_document": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"table_encryption_key_type": {
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"tags": {
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Optional: true,
			Type:     pluginsdk.TypeMap,
		},
	}
}

func (AccountV3ToV4) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {

		if _, ok := rawState["cross_tenant_replication_enabled"]; !ok {
			rawState["cross_tenant_replication_enabled"] = true
		}

		return rawState, nil
	}
}
