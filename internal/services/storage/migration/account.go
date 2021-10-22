package migration

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-01-01/storage"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	msiValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ pluginsdk.StateUpgrade = AccountV0ToV1{}

type AccountV0ToV1 struct{}

func (AccountV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return accountSchemaForV0AndV1()
}

func (AccountV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// this should have been applied from pre-0.12 migration system; backporting just in-case
		accountType := rawState["account_type"].(string)
		split := strings.Split(accountType, "_")
		rawState["account_tier"] = split[0]
		rawState["account_replication_type"] = split[1]
		return rawState, nil
	}
}

var _ pluginsdk.StateUpgrade = AccountV1ToV2{}

type AccountV1ToV2 struct{}

func (AccountV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return accountSchemaForV0AndV1()
}

func (AccountV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	// this should have been applied from pre-0.12 migration system; backporting just in-case
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		rawState["account_encryption_source"] = "Microsoft.Storage"
		return rawState, nil
	}
}

func accountSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"account_kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "Storage",
		},

		"account_type": {
			Type:       pluginsdk.TypeString,
			Optional:   true,
			Computed:   true,
			Deprecated: "This field has been split into `account_tier` and `account_replication_type`",
		},

		"account_tier": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"account_replication_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		// Only valid for BlobStorage accounts, defaults to "Hot" in create function
		"access_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"custom_domain": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"use_subdomain": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"enable_blob_encryption": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"enable_file_encryption": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"enable_https_traffic_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"primary_location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_blob_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_blob_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_queue_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_queue_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_table_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_table_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		// NOTE: The API does not appear to expose a secondary file endpoint
		"primary_file_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_access_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_access_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_blob_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_blob_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
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

func accountSchemaForV2() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageAccountName,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": azure.SchemaLocation(),

		"account_kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(storage.Storage),
				string(storage.BlobStorage),
				string(storage.BlockBlobStorage),
				string(storage.FileStorage),
				string(storage.StorageV2),
			}, true),
			Default: string(storage.StorageV2),
		},

		"account_tier": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Standard",
				"Premium",
			}, true),
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"account_replication_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"LRS",
				"ZRS",
				"GRS",
				"RAGRS",
				"GZRS",
				"RAGZRS",
			}, true),
			DiffSuppressFunc: suppress.CaseDifference,
		},

		// Only valid for BlobStorage & StorageV2 accounts, defaults to "Hot" in create function
		"access_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(storage.Cool),
				string(storage.Hot),
			}, true),
		},

		"azure_files_authentication": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"directory_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(storage.DirectoryServiceOptionsAADDS),
							string(storage.DirectoryServiceOptionsAD),
						}, false),
					},

					"active_directory": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"storage_sid": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"domain_guid": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsUUID,
								},

								"domain_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"domain_sid": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"forest_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"netbios_domain_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},

		"custom_domain": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"use_subdomain": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"enable_https_traffic_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"min_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(storage.TLS10),
			ValidateFunc: validation.StringInSlice([]string{
				string(storage.TLS10),
				string(storage.TLS11),
				string(storage.TLS12),
			}, false),
		},

		"is_hns_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},

		"nfsv3_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},

		"allow_blob_public_access": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"shared_access_key_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"network_rules": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"bypass": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(storage.AzureServices),
								string(storage.Logging),
								string(storage.Metrics),
								string(storage.None),
							}, true),
						},
						Set: pluginsdk.HashString,
					},

					"ip_rules": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.StorageAccountIpRule,
						},
						Set: pluginsdk.HashString,
					},

					"virtual_network_subnet_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Computed: true,
						Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						Set:      pluginsdk.HashString,
					},

					"default_action": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(storage.DefaultActionAllow),
							string(storage.DefaultActionDeny),
						}, false),
					},

					"private_link_access": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"endpoint_resource_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: azure.ValidateResourceID,
								},

								"endpoint_tenant_id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Computed:     true,
									ValidateFunc: validation.IsUUID,
								},
							},
						},
					},
				},
			},
		},

		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						DiffSuppressFunc: suppress.CaseDifference,
						ValidateFunc: validation.StringInSlice([]string{
							string(storage.IdentityTypeSystemAssigned),
							string(storage.IdentityTypeSystemAssignedUserAssigned),
							string(storage.IdentityTypeUserAssigned),
						}, true),
					},
					"principal_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"identity_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: msiValidate.UserAssignedIdentityID,
						},
					},
				},
			},
		},

		"blob_properties": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cors_rule": helpers.SchemaStorageAccountCorsRule(true),
					"delete_retention_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"days": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      7,
									ValidateFunc: validation.IntBetween(1, 365),
								},
							},
						},
					},

					"versioning_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"change_feed_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"default_service_version": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validate.BlobPropertiesDefaultServiceVersion,
					},

					"last_access_time_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"container_delete_retention_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"days": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      7,
									ValidateFunc: validation.IntBetween(1, 365),
								},
							},
						},
					},
				},
			},
		},

		"queue_properties": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cors_rule": helpers.SchemaStorageAccountCorsRule(false),
					"logging": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"version": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"delete": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"read": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"write": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"retention_policy_days": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntBetween(1, 365),
								},
							},
						},
					},
					"hour_metrics": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"version": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"include_apis": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},
								"retention_policy_days": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntBetween(1, 365),
								},
							},
						},
					},
					"minute_metrics": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"version": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"enabled": {
									Type:     pluginsdk.TypeBool,
									Required: true,
								},
								"include_apis": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
								},
								"retention_policy_days": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntBetween(1, 365),
								},
							},
						},
					},
				},
			},
		},

		"routing": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"publish_internet_endpoints": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"publish_microsoft_endpoints": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"choice": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(storage.MicrosoftRouting),
							string(storage.InternetRouting),
						}, false),
						Default: string(storage.MicrosoftRouting),
					},
				},
			},
		},

		"share_properties": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cors_rule": helpers.SchemaStorageAccountCorsRule(true),

					"retention_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"days": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									Default:      7,
									ValidateFunc: validation.IntBetween(1, 365),
								},
							},
						},
					},

					"smb": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"versions": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"SMB2.1",
											"SMB3.0",
											"SMB3.1.1",
										}, false),
									},
								},

								"authentication_types": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"NTLMv2",
											"Kerberos",
										}, false),
									},
								},

								"kerberos_ticket_encryption_type": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"RC4-HMAC",
											"AES-256",
										}, false),
									},
								},

								"channel_encryption_type": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"AES-128-CCM",
											"AES-128-GCM",
											"AES-256-GCM",
										}, false),
									},
								},
							},
						},
					},
				},
			},
		},

		//lintignore:XS003
		"static_website": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"index_document": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"error_404_document": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"large_file_share_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"primary_location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_location": {
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

		"secondary_blob_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_blob_host": {
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

		"secondary_queue_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_queue_host": {
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

		"secondary_table_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_table_host": {
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

		"secondary_web_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_web_host": {
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

		"secondary_dfs_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_dfs_host": {
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

		"secondary_file_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_file_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_access_key": {
			Type:      pluginsdk.TypeString,
			Sensitive: true,
			Computed:  true,
		},

		"secondary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_blob_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_blob_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"tags": {
			Type:         pluginsdk.TypeMap,
			Optional:     true,
			ValidateFunc: validate.StorageAccountTags,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

type AccountV2ToV3 struct{}

func (AccountV2ToV3) Schema() map[string]*pluginsdk.Schema {
	return accountSchemaForV2()
}

func (AccountV2ToV3) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		x, ok := rawState["allow_blob_public_access"]
		if ok {
			rawState["allow_nested_items_to_be_public"] = x
			delete(rawState, "allow_blob_public_access")
		}
		return rawState, nil
	}
}
