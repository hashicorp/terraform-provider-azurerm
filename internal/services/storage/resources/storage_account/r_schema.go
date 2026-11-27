// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func storageAccountSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StorageAccountName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"account_kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(storageaccounts.KindStorageVTwo),
			ValidateFunc: validation.StringInSlice([]string{
				string(storageaccounts.KindBlockBlobStorage),
				string(storageaccounts.KindFileStorage),
				string(storageaccounts.KindStorageVTwo),
			}, false),
		},

		"account_tier": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForSkuTier(), false),
		},

		"provisioned_billing_model_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			// There are two billing model in Azure Storage (at least for File): Provisioned (V1 or V2) and PAYG.
			// (https://learn.microsoft.com/azure/storage/files/understanding-billing)
			// The version only applies to Provisioned model. When reflecting to the API payload, Provisioned V1 and PAYG
			// look the same, i.e. the `sku.name` will be no suffix for it's first part (e.g. "Standard" in "Standard_LRS").
			// Only for Provisioned V2, there will be a "V2" suffix added, e.g. "StandardV2" in "StandardV2_LRS".
			// Due to this, we don't differentiate the empty string and "V1" here, as they ends up to be an empty suffix.
			//
			// Aside, at least for File, the diff between Provisioned and PYAG is determined by the "account_kind":
			// - account_kind = Storage/StorageV2: PAYG
			// - account_kind = FileStorage: Provisioned
			ValidateFunc: validation.StringInSlice([]string{"V2"}, false),
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
			}, false),
			DiffSuppressFunc: func(_, _, n string, d *schema.ResourceData) bool {
				// Migration of a storage account to/from zonal can take days
				// so the provider doesn't poll and wait for it to complete.
				// While in queue for migration/actively migrating, the Storage Account will return the old `account_replication_type`
				// which we'll suppress here if the configured type matches the target migration type.
				return d.Get("account_replication_type_migration_in_progress").(bool) && n == d.Get("account_replication_type_migrating_to").(string)
			},
		},

		"access_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C Only valid for FileStorage, BlobStorage & StorageV2 accounts, defaults to "Hot" in create function
			Computed:     true,
			ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForAccessTier(), false), // TODO: docs for `Premium`
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
							string(storageaccounts.DirectoryServiceOptionsAADDS),
							string(storageaccounts.DirectoryServiceOptionsAADKERB),
							string(storageaccounts.DirectoryServiceOptionsAD),
						}, false),
					},

					"active_directory": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
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

								"storage_sid": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"domain_sid": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"forest_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"netbios_domain_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"default_share_level_permission": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(storageaccounts.DefaultSharePermissionNone),
						ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForDefaultSharePermission(), false),
					},
				},
			},
		},

		"smb_oauth_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"cross_tenant_replication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
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

		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"key_vault_key_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeKey),
					},

					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"https_traffic_only_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"immutability_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allow_protected_append_writes": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"period_since_creation_in_days": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"state": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForAccountImmutabilityPolicyState(), false),
					},
				},
			},
		},

		"min_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(storageaccounts.MinimumTlsVersionTLSOneTwo),
			ValidateFunc: validation.StringInSlice([]string{
				string(storageaccounts.MinimumTlsVersionTLSOneTwo),
			}, false),
		},

		"hns_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"nfsv3_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"allow_nested_items_to_be_public": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"shared_access_key_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"public_network_access": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(storageaccounts.PublicNetworkAccessEnabled),
			ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForPublicNetworkAccess(), false),
		},

		"dns_endpoint_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(storageaccounts.DnsEndpointTypeStandard),
			ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForDnsEndpointType(), false),
		},

		"default_to_oauth_authentication": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"network_rules": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"bypass": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						Set:      pluginsdk.HashString,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForBypass(), false),
						},
					},

					"ip_rules": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						Set:      pluginsdk.HashString,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.StorageAccountIpRule,
						},
					},

					"virtual_network_subnet_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						Set:      pluginsdk.HashString,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"default_action": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForDefaultAction(), false),
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
									Computed:     true, // azignore:AZS007 - pre-existing violation
									ValidateFunc: validation.IsUUID,
								},
							},
						},
					},
				},
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"blob_properties": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"change_feed_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"change_feed_retention_in_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 146000),
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

					"cors_rule": SchemaStorageAccountCorsRule(true),

					"default_service_version": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validate.BlobPropertiesDefaultServiceVersion,
					},

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
								"permanent_delete_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
							},
						},
					},
					"last_access_time_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"restore_policy": {
						Type:         pluginsdk.TypeList,
						Optional:     true,
						MaxItems:     1,
						RequiredWith: []string{"blob_properties.0.delete_retention_policy"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"days": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
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
				},
			},
		},

		"routing": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(storageaccounts.RoutingChoiceMicrosoftRouting),
						ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForRoutingChoice(), false),
					},
				},
			},
		},

		"share_properties": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cors_rule": SchemaStorageAccountCorsRule(true),

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
								"authentication_types": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"Kerberos",
											"NTLMv2",
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

								"kerberos_ticket_encryption_type": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"AES-256",
											"RC4-HMAC",
										}, false),
									},
								},

								"multichannel_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

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
							},
						},
					},
				},
			},
		},

		"queue_encryption_key_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(storageaccounts.KeyTypeService),
			ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForKeyType(), false),
		},

		"table_encryption_key_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(storageaccounts.KeyTypeService),
			ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForKeyType(), false),
		},

		"infrastructure_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"sas_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expiration_action": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(storageaccounts.ExpirationActionLog),
						ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForExpirationAction(), false),
					},
					"expiration_period": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"allowed_copy_scope": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForAllowedCopyScope(), false),
		},

		"sftp_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"large_file_share_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"local_user_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"account_replication_type_migration_in_progress": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"account_replication_type_migrating_to": {
			Type:     pluginsdk.TypeString,
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

		"primary_blob_internet_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_blob_internet_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_blob_microsoft_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_blob_microsoft_host": {
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

		"secondary_blob_internet_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_blob_internet_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_blob_microsoft_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_blob_microsoft_host": {
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

		"primary_queue_microsoft_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_queue_microsoft_host": {
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

		"secondary_queue_microsoft_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_queue_microsoft_host": {
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

		"primary_table_microsoft_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_table_microsoft_host": {
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

		"secondary_table_microsoft_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_table_microsoft_host": {
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

		"primary_web_internet_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_web_internet_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_web_microsoft_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_web_microsoft_host": {
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

		"secondary_web_internet_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_web_internet_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_web_microsoft_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_web_microsoft_host": {
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

		"primary_dfs_internet_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_dfs_internet_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_dfs_microsoft_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_dfs_microsoft_host": {
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

		"secondary_dfs_internet_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_dfs_internet_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_dfs_microsoft_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_dfs_microsoft_host": {
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

		"primary_file_internet_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_file_internet_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_file_microsoft_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_file_microsoft_host": {
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

		"secondary_file_internet_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_file_internet_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_file_microsoft_host": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_file_microsoft_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
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
			// TODO: introduce/refactor this to use a `commonschema.TagsOptionalWith(a, b, c)` to enable us to handle this in one place
			Type:         pluginsdk.TypeMap,
			Optional:     true,
			ValidateFunc: validate.StorageAccountTags,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}
