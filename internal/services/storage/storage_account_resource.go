// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/blobservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/fileservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	keyVaultsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	managedHsmParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	managedHsmValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/jackofallops/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/jackofallops/giovanni/storage/2023-11-03/queue/queues"
)

var (
	storageAccountResourceName  = "azurerm_storage_account"
	storageKindsSupportsSkuTier = map[storageaccounts.Kind]struct{}{
		storageaccounts.KindBlobStorage: {},
		storageaccounts.KindFileStorage: {},
		storageaccounts.KindStorageVTwo: {},
	}
	storageKindsSupportHns = map[storageaccounts.Kind]struct{}{
		storageaccounts.KindBlobStorage:      {},
		storageaccounts.KindBlockBlobStorage: {},
		storageaccounts.KindStorageVTwo:      {},
	}
	storageKindsSupportLargeFileShares = map[storageaccounts.Kind]struct{}{
		storageaccounts.KindFileStorage: {},
		storageaccounts.KindStorageVTwo: {},
	}
)

func resourceStorageAccount() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceStorageAccountCreate,
		Read:   resourceStorageAccountRead,
		Update: resourceStorageAccountUpdate,
		Delete: resourceStorageAccountDelete,

		SchemaVersion: 4,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AccountV0ToV1{},
			1: migration.AccountV1ToV2{},
			2: migration.AccountV2ToV3{},
			3: migration.AccountV3ToV4{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseStorageAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"account_kind": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForKind(), false),
				Default:      string(storageaccounts.KindStorageVTwo),
			},

			"account_tier": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForSkuTier(), false),
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
			},

			// Only valid for FileStorage, BlobStorage & StorageV2 accounts, defaults to "Hot" in create function
			"access_tier": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
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
							Computed: true,
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
							Optional:     true,
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
							ExactlyOneOf: []string{"customer_managed_key.0.managed_hsm_key_id", "customer_managed_key.0.key_vault_key_id"},
						},

						"managed_hsm_key_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.Any(managedHsmValidate.ManagedHSMDataPlaneVersionedKeyID, managedHsmValidate.ManagedHSMDataPlaneVersionlessKeyID),
							ExactlyOneOf: []string{"customer_managed_key.0.managed_hsm_key_id", "customer_managed_key.0.key_vault_key_id"},
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
				MaxItems: 1,
				ForceNew: true,
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
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(storageaccounts.MinimumTlsVersionTLSOneTwo),
				ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForMinimumTlsVersion(), false),
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

			"allow_nested_items_to_be_public": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"shared_access_key_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"dns_endpoint_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(storageaccounts.DnsEndpointTypeStandard),
				ValidateFunc: validation.StringInSlice([]string{
					string(storageaccounts.DnsEndpointTypeStandard),
					string(storageaccounts.DnsEndpointTypeAzureDnsZone),
				}, false),
			},

			"default_to_oauth_authentication": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
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
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForBypass(), false),
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
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							Set: pluginsdk.HashString,
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
										Computed:     true,
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
				Computed: true,
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

						"cors_rule": helpers.SchemaStorageAccountCorsRule(true),

						"default_service_version": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
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
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"days": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 365),
									},
								},
							},
							RequiredWith: []string{"blob_properties.0.delete_retention_policy"},
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
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForRoutingChoice(), false),
							Default:      string(storageaccounts.RoutingChoiceMicrosoftRouting),
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
				ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForKeyType(), false),
				Default:      string(storageaccounts.KeyTypeService),
			},

			"table_encryption_key_type": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(storageaccounts.PossibleValuesForKeyType(), false),
				Default:      string(storageaccounts.KeyTypeService),
			},

			"infrastructure_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"sas_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"expiration_action": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(storageaccounts.ExpirationActionLog),
							ValidateFunc: validation.StringInSlice([]string{
								string(storageaccounts.ExpirationActionLog),
							}, false),
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
				Computed: true,
			},

			"local_user_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
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
				// TODO: introduce/refactor this to use a `commonschema.TagsOptionalWith(a, b, c)` to enable us to handle this in one place
				Type:         pluginsdk.TypeMap,
				Optional:     true,
				ValidateFunc: validate.StorageAccountTags,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
				if d.HasChange("account_kind") {
					accountKind, changedKind := d.GetChange("account_kind")

					if accountKind != string(storageaccounts.KindStorage) && changedKind != string(storageaccounts.KindStorageVTwo) {
						log.Printf("[DEBUG] recreate storage account, could't be migrated from %q to %q", accountKind, changedKind)
						d.ForceNew("account_kind")
						return nil
					} else {
						log.Printf("[DEBUG] storage account can be upgraded from %q to %q", accountKind, changedKind)
					}
				}

				if d.HasChange("large_file_share_enabled") {
					lfsEnabled, changedEnabled := d.GetChange("large_file_share_enabled")
					if lfsEnabled.(bool) && !changedEnabled.(bool) {
						d.ForceNew("large_file_share_enabled")
					}
				}

				if d.Get("access_tier") != "" {
					accountKind := storageaccounts.Kind(d.Get("account_kind").(string))
					if _, ok := storageKindsSupportsSkuTier[accountKind]; !ok {
						keys := sortedKeysFromSlice(storageKindsSupportsSkuTier)
						return fmt.Errorf("`access_tier` is only available for accounts where `kind` is set to one of: %+v", strings.Join(keys, " / "))
					}
				}

				if !features.FivePointOhBeta() && !v.(*clients.Client).Features.Storage.DataPlaneAvailable {
					if _, ok := d.GetOk("queue_properties"); ok {
						return errors.New("cannot configure 'queue_properties' when the Provider Feature 'data_plane_available' is set to 'false'")
					}
					if _, ok := d.GetOk("static_website"); ok {
						return errors.New("cannot configure 'static_website' when the Provider Feature 'data_plane_available' is set to 'false'")
					}
				}

				return nil
			}),
			pluginsdk.ForceNewIfChange("account_replication_type", func(ctx context.Context, old, new, meta interface{}) bool {
				newAccRep := strings.ToUpper(new.(string))

				switch strings.ToUpper(old.(string)) {
				case "LRS", "GRS", "RAGRS":
					if newAccRep == "GZRS" || newAccRep == "RAGZRS" || newAccRep == "ZRS" {
						return true
					}
				case "ZRS", "GZRS", "RAGZRS":
					if newAccRep == "LRS" || newAccRep == "GRS" || newAccRep == "RAGRS" {
						return true
					}
				}
				return false
			}),
		),
	}

	if !features.FivePointOhBeta() {
		// lintignore:XS003
		resource.Schema["static_website"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"error_404_document": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"index_document": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			Deprecated: "this block has been deprecated and superseded by the `azurerm_storage_account_static_website` resource and will be removed in v5.0 of the AzureRM provider",
		}
	}

	resource.Schema["queue_properties"] = &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"cors_rule": helpers.SchemaStorageAccountCorsRule(false),
				"hour_metrics": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
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
				"logging": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
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
				"minute_metrics": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Computed: true,
					MaxItems: 1,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"version": {
								Type:         pluginsdk.TypeString,
								Required:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							// TODO 4.0: Remove this property and determine whether to enable based on existence of the out side block.
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
		Deprecated: "this block has been deprecated and superseded by the `azurerm_storage_account_queue_properties` resource and will be removed in v5.0 of the AzureRM provider",
	}

	return resource
}

func resourceStorageAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	tenantId := meta.(*clients.Client).Account.TenantId
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	storageUtils := meta.(*clients.Client).Storage
	storageClient := meta.(*clients.Client).Storage.ResourceManager
	client := storageClient.StorageAccounts
	keyVaultClient := meta.(*clients.Client).KeyVault
	dataPlaneAvailable := meta.(*clients.Client).Features.Storage.DataPlaneAvailable
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewStorageAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	existing, err := client.GetProperties(ctx, id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_storage_account", id.ID())
	}

	accountKind := storageaccounts.Kind(d.Get("account_kind").(string))
	accountTier := storageaccounts.SkuTier(d.Get("account_tier").(string))
	replicationType := d.Get("account_replication_type").(string)

	publicNetworkAccess := storageaccounts.PublicNetworkAccessDisabled
	if d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = storageaccounts.PublicNetworkAccessEnabled
	}
	expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	httpsTrafficOnlyEnabled := true
	// nolint staticcheck
	if v, ok := d.GetOkExists("https_traffic_only_enabled"); ok {
		httpsTrafficOnlyEnabled = v.(bool)
	}

	dnsEndpointType := d.Get("dns_endpoint_type").(string)
	isHnsEnabled := d.Get("is_hns_enabled").(bool)
	nfsV3Enabled := d.Get("nfsv3_enabled").(bool)
	payload := storageaccounts.StorageAccountCreateParameters{
		ExtendedLocation: expandEdgeZone(d.Get("edge_zone").(string)),
		Kind:             accountKind,
		Identity:         expandedIdentity,
		Location:         location.Normalize(d.Get("location").(string)),
		Properties: &storageaccounts.StorageAccountPropertiesCreateParameters{
			AllowBlobPublicAccess:        pointer.To(d.Get("allow_nested_items_to_be_public").(bool)),
			AllowCrossTenantReplication:  pointer.To(d.Get("cross_tenant_replication_enabled").(bool)),
			AllowSharedKeyAccess:         pointer.To(d.Get("shared_access_key_enabled").(bool)),
			DnsEndpointType:              pointer.To(storageaccounts.DnsEndpointType(dnsEndpointType)),
			DefaultToOAuthAuthentication: pointer.To(d.Get("default_to_oauth_authentication").(bool)),
			SupportsHTTPSTrafficOnly:     pointer.To(httpsTrafficOnlyEnabled),
			IsNfsV3Enabled:               pointer.To(nfsV3Enabled),
			IsHnsEnabled:                 pointer.To(isHnsEnabled),
			IsLocalUserEnabled:           pointer.To(d.Get("local_user_enabled").(bool)),
			IsSftpEnabled:                pointer.To(d.Get("sftp_enabled").(bool)),
			MinimumTlsVersion:            pointer.To(storageaccounts.MinimumTlsVersion(d.Get("min_tls_version").(string))),
			NetworkAcls:                  expandAccountNetworkRules(d.Get("network_rules").([]interface{}), tenantId),
			PublicNetworkAccess:          pointer.To(publicNetworkAccess),
			SasPolicy:                    expandAccountSASPolicy(d.Get("sas_policy").([]interface{})),
		},
		Sku: storageaccounts.Sku{
			Name: storageaccounts.SkuName(fmt.Sprintf("%s_%s", string(accountTier), replicationType)),
			Tier: pointer.To(accountTier),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v := d.Get("allowed_copy_scope").(string); v != "" {
		payload.Properties.AllowedCopyScope = pointer.To(storageaccounts.AllowedCopyScope(v))
	}
	if v, ok := d.GetOk("azure_files_authentication"); ok {
		expandAADFilesAuthentication, err := expandAccountAzureFilesAuthentication(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("parsing `azure_files_authentication`: %v", err)
		}
		payload.Properties.AzureFilesIdentityBasedAuthentication = expandAADFilesAuthentication
	}
	if _, ok := d.GetOk("custom_domain"); ok {
		payload.Properties.CustomDomain = expandAccountCustomDomain(d.Get("custom_domain").([]interface{}))
	}
	if v, ok := d.GetOk("immutability_policy"); ok {
		payload.Properties.ImmutableStorageWithVersioning = expandAccountImmutabilityPolicy(v.([]interface{}))
	}

	// BlobStorage does not support ZRS
	if accountKind == storageaccounts.KindBlobStorage && string(payload.Sku.Name) == string(storageaccounts.SkuNameStandardZRS) {
		return fmt.Errorf("`account_replication_type` of `ZRS` isn't supported for Blob Storage accounts")
	}

	accessTier, accessTierSetInConfig := d.GetOk("access_tier")
	_, skuTierSupported := storageKindsSupportsSkuTier[accountKind]
	if !skuTierSupported && accessTierSetInConfig {
		keys := sortedKeysFromSlice(storageKindsSupportsSkuTier)
		return fmt.Errorf("`access_tier` is only available for accounts of kind set to one of: %+v", strings.Join(keys, " / "))
	}
	if skuTierSupported {
		if !accessTierSetInConfig {
			// default to "Hot"
			accessTier = string(storageaccounts.AccessTierHot)
		}
		payload.Properties.AccessTier = pointer.To(storageaccounts.AccessTier(accessTier.(string)))
	}

	if _, supportsHns := storageKindsSupportHns[accountKind]; !supportsHns && isHnsEnabled {
		keys := sortedKeysFromSlice(storageKindsSupportHns)
		return fmt.Errorf("`is_hns_enabled` can only be used for accounts with `kind` set to one of: %+v", strings.Join(keys, " / "))
	}

	// NFSv3 is supported for standard general-purpose v2 storage accounts and for premium block blob storage accounts.
	// (https://docs.microsoft.com/en-us/azure/storage/blobs/network-file-system-protocol-support-how-to#step-5-create-and-configure-a-storage-account)
	if nfsV3Enabled {
		if !isHnsEnabled {
			return fmt.Errorf("`nfsv3_enabled` can only be used when `is_hns_enabled` is `true`")
		}

		isPremiumTierAndBlockBlobStorageKind := accountTier == storageaccounts.SkuTierPremium && accountKind == storageaccounts.KindBlockBlobStorage
		isStandardTierAndStorageV2Kind := accountTier == storageaccounts.SkuTierStandard && accountKind == storageaccounts.KindStorageVTwo
		if !isPremiumTierAndBlockBlobStorageKind && !isStandardTierAndStorageV2Kind {
			return fmt.Errorf("`nfsv3_enabled` can only be used with account tier `Standard` and account kind `StorageV2`, or account tier `Premium` and account kind `BlockBlobStorage`")
		}
	}

	// AccountTier must be Premium for FileStorage
	if accountKind == storageaccounts.KindFileStorage && accountTier != storageaccounts.SkuTierPremium {
		return fmt.Errorf("`account_tier` must be `Premium` for File Storage accounts")
	}

	// nolint staticcheck
	if v, ok := d.GetOkExists("large_file_share_enabled"); ok {
		// @tombuildsstuff: we can't set this to `false` because the API returns:
		//
		// performing Create: unexpected status 400 (400 Bad Request) with error: InvalidRequestPropertyValue: The
		// value 'Disabled' is not allowed for property largeFileSharesState. For more information, see -
		// https://aka.ms/storageaccountlargefilesharestate
		if v.(bool) {
			if _, ok := storageKindsSupportLargeFileShares[accountKind]; !ok {
				keys := sortedKeysFromSlice(storageKindsSupportLargeFileShares)
				return fmt.Errorf("`large_file_shares_enabled` can only be set to `true` with `account_kind` set to one of: %+v", strings.Join(keys, " / "))
			}
			payload.Properties.LargeFileSharesState = pointer.To(storageaccounts.LargeFileSharesStateEnabled)
		}
	}

	if v, ok := d.GetOk("routing"); ok {
		payload.Properties.RoutingPreference = expandAccountRoutingPreference(v.([]interface{}))
	}

	// TODO look into standardizing this across resources that support CMK and at the very least look at improving the UX
	// for encryption of blob, file, table and queue
	//
	// By default (by leaving empty), the table and queue encryption key type is set to "Service". While users can change it to "Account" so that
	// they can further use CMK to encrypt table/queue data. Only the StorageV2 account kind supports the Account key type.
	// Also noted that the blob and file are always using the "Account" key type.
	// See: https://docs.microsoft.com/en-gb/azure/storage/common/account-encryption-key-create?tabs=portal
	queueEncryptionKeyType := storageaccounts.KeyType(d.Get("queue_encryption_key_type").(string))
	tableEncryptionKeyType := storageaccounts.KeyType(d.Get("table_encryption_key_type").(string))
	encryptionRaw := d.Get("customer_managed_key").([]interface{})
	encryption, err := expandAccountCustomerManagedKey(ctx, keyVaultClient, id.SubscriptionId, encryptionRaw, accountTier, accountKind, *expandedIdentity, queueEncryptionKeyType, tableEncryptionKeyType)
	if err != nil {
		return fmt.Errorf("expanding `customer_managed_key`: %+v", err)
	}

	infrastructureEncryption := d.Get("infrastructure_encryption_enabled").(bool)

	if infrastructureEncryption {
		validPremiumConfiguration := accountTier == storageaccounts.SkuTierPremium && (accountKind == storageaccounts.KindBlockBlobStorage) || accountKind == storageaccounts.KindFileStorage
		validV2Configuration := accountKind == storageaccounts.KindStorageVTwo
		if !(validPremiumConfiguration || validV2Configuration) {
			return fmt.Errorf("`infrastructure_encryption_enabled` can only be used with account kind `StorageV2`, or account tier `Premium` and account kind is one of `BlockBlobStorage` or `FileStorage`")
		}
		encryption.RequireInfrastructureEncryption = &infrastructureEncryption
	}

	payload.Properties.Encryption = encryption

	if err := client.CreateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// populate the cache
	account, err := client.GetProperties(ctx, id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if account.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if err := storageUtils.AddToCache(id, *account.Model); err != nil {
		return fmt.Errorf("populating cache for %s: %+v", id, err)
	}

	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, replicationType)
	// Start of Data Plane access - this entire block can be removed for 5.0, as the data_plane_available flag becomes redundant at that time.
	if !features.FivePointOhBeta() && dataPlaneAvailable {
		dataPlaneClient := meta.(*clients.Client).Storage
		dataPlaneAccount, err := storageUtils.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if dataPlaneAccount == nil {
			return fmt.Errorf("unable to locate %q", id)
		}

		if err := waitForDataPlaneToBecomeAvailableForAccount(ctx, dataPlaneClient, dataPlaneAccount, supportLevel); err != nil {
			return fmt.Errorf("waiting for the Data Plane for %s to become available: %+v", id, err)
		}

		if val, ok := d.GetOk("queue_properties"); ok {
			if !supportLevel.supportQueue {
				return fmt.Errorf("`queue_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
			}

			queueClient, err := dataPlaneClient.QueuesDataPlaneClient(ctx, *dataPlaneAccount, dataPlaneClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Queues Client: %s", err)
			}

			queueProperties, err := expandAccountQueueProperties(val.([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `queue_properties`: %+v", err)
			}

			if err = queueClient.UpdateServiceProperties(ctx, *queueProperties); err != nil {
				return fmt.Errorf("updating Queue Properties: %+v", err)
			}

			if err = d.Set("queue_properties", val); err != nil {
				return fmt.Errorf("setting `queue_properties`: %+v", err)
			}
		}

		if val, ok := d.GetOk("static_website"); ok {
			if !supportLevel.supportStaticWebsite {
				return fmt.Errorf("`static_website` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
			}

			accountsClient, err := dataPlaneClient.AccountsDataPlaneClient(ctx, *dataPlaneAccount, dataPlaneClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Accounts Data Plane Client: %s", err)
			}

			staticWebsiteProps := expandAccountStaticWebsiteProperties(val.([]interface{}))

			if _, err = accountsClient.SetServiceProperties(ctx, id.StorageAccountName, staticWebsiteProps); err != nil {
				return fmt.Errorf("updating `static_website`: %+v", err)
			}

			if err = d.Set("static_website", val); err != nil {
				return fmt.Errorf("setting `static_website`: %+v", err)
			}
		}
	}

	if val, ok := d.GetOk("blob_properties"); ok {
		if !supportLevel.supportBlob {
			return fmt.Errorf("`blob_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		blobProperties, err := expandAccountBlobServiceProperties(accountKind, val.([]interface{}))
		if err != nil {
			return err
		}

		// See: https://learn.microsoft.com/en-us/azure/storage/blobs/versioning-overview#:~:text=Storage%20accounts%20with%20a%20hierarchical%20namespace%20enabled%20for%20use%20with%20Azure%20Data%20Lake%20Storage%20Gen2%20are%20not%20currently%20supported.
		isVersioningEnabled := pointer.From(blobProperties.Properties.IsVersioningEnabled)
		if isVersioningEnabled && isHnsEnabled {
			return fmt.Errorf("`versioning_enabled` can't be true when `is_hns_enabled` is true")
		}

		if !isVersioningEnabled {
			if blobProperties.Properties.RestorePolicy != nil && blobProperties.Properties.RestorePolicy.Enabled {
				// Otherwise, API returns: "Conflicting feature 'restorePolicy' is enabled. Please disable it and retry."
				return fmt.Errorf("`blob_properties.restore_policy` can't be set when `versioning_enabled` is false")
			}

			immutableStorageWithVersioningEnabled := false
			if props := account.Model.Properties; props != nil {
				if versioning := props.ImmutableStorageWithVersioning; versioning != nil {
					if versioning.ImmutabilityPolicy != nil && versioning.Enabled != nil {
						immutableStorageWithVersioningEnabled = *versioning.Enabled
					}
				}
			}
			if immutableStorageWithVersioningEnabled {
				// Otherwise, API returns: "Conflicting feature 'Account level WORM' is enabled. Please disable it and retry."
				// See: https://learn.microsoft.com/en-us/azure/storage/blobs/immutable-policy-configure-version-scope?tabs=azure-portal#prerequisites
				return fmt.Errorf("`immutability_policy` can't be set when `versioning_enabled` is false")
			}
		}

		// TODO: This is a temporary limitation on Storage service. Remove this check once the API supports this scenario.
		// See https://github.com/hashicorp/terraform-provider-azurerm/pull/25450#discussion_r1542471667 for the context.
		if dnsEndpointType == string(storageaccounts.DnsEndpointTypeAzureDnsZone) {
			if blobProperties.Properties.RestorePolicy != nil && blobProperties.Properties.RestorePolicy.Enabled {
				// Otherwise, API returns: "Required feature Global Dns is disabled"
				// This is confirmed with the SRP team, where they said:
				// > restorePolicy feature is incompatible with partitioned DNS
				return fmt.Errorf("`blob_properties.restore_policy` can't be set when `dns_endpoint_type` is set to `%s`", storageaccounts.DnsEndpointTypeAzureDnsZone)
			}
		}

		if _, err = storageClient.BlobService.SetServiceProperties(ctx, id, *blobProperties); err != nil {
			return fmt.Errorf("updating `blob_properties`: %+v", err)
		}
	}

	if val, ok := d.GetOk("share_properties"); ok {
		if !supportLevel.supportShare {
			return fmt.Errorf("`share_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		sharePayload := expandAccountShareProperties(val.([]interface{}))

		// The API complains if any multichannel info is sent on non premium fileshares. Even if multichannel is set to false
		if accountTier != storageaccounts.SkuTierPremium && sharePayload.Properties != nil && sharePayload.Properties.ProtocolSettings != nil {
			// Error if the user has tried to enable multichannel on a standard tier storage account
			smb := sharePayload.Properties.ProtocolSettings.Smb
			if smb != nil && smb.Multichannel != nil {
				if smb.Multichannel.Enabled != nil && *smb.Multichannel.Enabled {
					return fmt.Errorf("`multichannel_enabled` isn't supported for Standard tier Storage accounts")
				}

				sharePayload.Properties.ProtocolSettings.Smb.Multichannel = nil
			}
		}

		if _, err = storageClient.FileService.SetServiceProperties(ctx, id, sharePayload); err != nil {
			return fmt.Errorf("updating `share_properties`: %+v", err)
		}
	}

	return resourceStorageAccountRead(d, meta)
}

func resourceStorageAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	tenantId := meta.(*clients.Client).Account.TenantId
	storageClient := meta.(*clients.Client).Storage.ResourceManager
	client := storageClient.StorageAccounts
	keyVaultClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	accountTier := storageaccounts.SkuTier(d.Get("account_tier").(string))
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)
	accountKind := storageaccounts.Kind(d.Get("account_kind").(string))

	if accountKind == storageaccounts.KindBlobStorage || accountKind == storageaccounts.KindStorage {
		if storageType == string(storageaccounts.SkuNameStandardZRS) {
			return fmt.Errorf("an `account_replication_type` of `ZRS` isn't supported for Blob Storage accounts")
		}
	}

	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Kind == nil {
		return fmt.Errorf("retrieving %s: `model.Kind` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model.Properties` was nil", id)
	}
	if existing.Model.Sku == nil {
		return fmt.Errorf("retrieving %s: `model.Sku` was nil", id)
	}

	props := storageaccounts.StorageAccountPropertiesCreateParameters{
		AccessTier:                            existing.Model.Properties.AccessTier,
		AllowBlobPublicAccess:                 existing.Model.Properties.AllowBlobPublicAccess,
		AllowedCopyScope:                      existing.Model.Properties.AllowedCopyScope,
		AllowSharedKeyAccess:                  existing.Model.Properties.AllowSharedKeyAccess,
		AllowCrossTenantReplication:           existing.Model.Properties.AllowCrossTenantReplication,
		AzureFilesIdentityBasedAuthentication: existing.Model.Properties.AzureFilesIdentityBasedAuthentication,
		CustomDomain:                          existing.Model.Properties.CustomDomain,
		DefaultToOAuthAuthentication:          existing.Model.Properties.DefaultToOAuthAuthentication,
		DnsEndpointType:                       existing.Model.Properties.DnsEndpointType,
		Encryption:                            existing.Model.Properties.Encryption,
		KeyPolicy:                             existing.Model.Properties.KeyPolicy,
		ImmutableStorageWithVersioning:        existing.Model.Properties.ImmutableStorageWithVersioning,
		IsNfsV3Enabled:                        existing.Model.Properties.IsNfsV3Enabled,
		IsSftpEnabled:                         existing.Model.Properties.IsSftpEnabled,
		IsLocalUserEnabled:                    existing.Model.Properties.IsLocalUserEnabled,
		IsHnsEnabled:                          existing.Model.Properties.IsHnsEnabled,
		MinimumTlsVersion:                     existing.Model.Properties.MinimumTlsVersion,
		NetworkAcls:                           existing.Model.Properties.NetworkAcls,
		PublicNetworkAccess:                   existing.Model.Properties.PublicNetworkAccess,
		RoutingPreference:                     existing.Model.Properties.RoutingPreference,
		SasPolicy:                             existing.Model.Properties.SasPolicy,
		SupportsHTTPSTrafficOnly:              existing.Model.Properties.SupportsHTTPSTrafficOnly,
	}

	if existing.Model.Properties.LargeFileSharesState != nil && *existing.Model.Properties.LargeFileSharesState == storageaccounts.LargeFileSharesStateEnabled {
		// We can only set this if it's Enabled, else the API complains during Update that we're sending Disabled, even if it's always been off
		props.LargeFileSharesState = existing.Model.Properties.LargeFileSharesState
	}

	expandedIdentity := existing.Model.Identity
	if d.HasChange("identity") {
		expandedIdentity, err = identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
	}

	if d.HasChange("access_tier") {
		props.AccessTier = pointer.To(storageaccounts.AccessTier(d.Get("access_tier").(string)))
	}
	if d.HasChange("allowed_copy_scope") {
		props.AllowedCopyScope = pointer.To(storageaccounts.AllowedCopyScope(d.Get("allowed_copy_scope").(string)))
	}
	if d.HasChange("allow_nested_items_to_be_public") {
		props.AllowBlobPublicAccess = pointer.To(d.Get("allow_nested_items_to_be_public").(bool))
	}
	if d.HasChange("cross_tenant_replication_enabled") {
		props.AllowCrossTenantReplication = pointer.To(d.Get("cross_tenant_replication_enabled").(bool))
	}
	if d.HasChange("custom_domain") {
		props.CustomDomain = expandAccountCustomDomain(d.Get("custom_domain").([]interface{}))
	}
	if d.HasChange("customer_managed_key") {
		queueEncryptionKeyType := storageaccounts.KeyType(d.Get("queue_encryption_key_type").(string))
		tableEncryptionKeyType := storageaccounts.KeyType(d.Get("table_encryption_key_type").(string))
		encryptionRaw := d.Get("customer_managed_key").([]interface{})
		encryption, err := expandAccountCustomerManagedKey(ctx, keyVaultClient, id.SubscriptionId, encryptionRaw, accountTier, accountKind, *expandedIdentity, queueEncryptionKeyType, tableEncryptionKeyType)
		if err != nil {
			return fmt.Errorf("expanding `customer_managed_key`: %+v", err)
		}

		// When updating CMK the existing value for `RequireInfrastructureEncryption` gets overwritten which results in
		// an error from the API so we set this back into encryption after it's been overwritten by this update
		existingEnc := existing.Model.Properties.Encryption
		if existingEnc != nil && existingEnc.RequireInfrastructureEncryption != nil {
			encryption.RequireInfrastructureEncryption = existingEnc.RequireInfrastructureEncryption
		}

		props.Encryption = encryption
	}
	if d.HasChange("shared_access_key_enabled") {
		props.AllowSharedKeyAccess = pointer.To(d.Get("shared_access_key_enabled").(bool))
	} else {
		// If AllowSharedKeyAccess is nil that breaks the Portal UI as reported in https://github.com/hashicorp/terraform-provider-azurerm/issues/11689
		// currently the Portal UI reports nil as false, and per the ARM API documentation nil is true. This manifests itself in the Portal UI
		// when a storage account is created by terraform that the AllowSharedKeyAccess is Disabled when it is actually Enabled, thus confusing out customers
		// to fix this, I have added this code to explicitly to set the value to true if is nil to workaround the Portal UI bug for our customers.
		// this is designed as a passive change, meaning the change will only take effect when the existing storage account is modified in some way if the
		// account already exists. since I have also switched up the default behaviour for net new storage accounts to always set this value as true, this issue
		// should automatically correct itself over time with these changes.
		// TODO: Remove code when Portal UI team fixes their code
		if sharedKeyAccess := props.AllowSharedKeyAccess; sharedKeyAccess == nil {
			props.AllowSharedKeyAccess = pointer.To(true)
		}
	}
	if d.HasChange("default_to_oauth_authentication") {
		props.DefaultToOAuthAuthentication = pointer.To(d.Get("default_to_oauth_authentication").(bool))
	}

	if d.HasChange("https_traffic_only_enabled") {
		props.SupportsHTTPSTrafficOnly = pointer.To(d.Get("https_traffic_only_enabled").(bool))
	}

	if d.HasChange("large_file_share_enabled") {
		// largeFileSharesState can only be set to `Enabled` and not `Disabled`, even if it is currently `Disabled`
		if oldValue, newValue := d.GetChange("large_file_share_enabled"); oldValue.(bool) && !newValue.(bool) {
			return fmt.Errorf("`large_file_share_enabled` cannot be disabled once it's been enabled")
		}

		if _, ok := storageKindsSupportLargeFileShares[accountKind]; !ok {
			keys := sortedKeysFromSlice(storageKindsSupportLargeFileShares)
			return fmt.Errorf("`large_file_shares_enabled` can only be set to `true` with `account_kind` set to one of: %+v", strings.Join(keys, " / "))
		}
		props.LargeFileSharesState = pointer.To(storageaccounts.LargeFileSharesStateEnabled)
	}
	if d.HasChange("local_user_enabled") {
		props.IsLocalUserEnabled = pointer.To(d.Get("local_user_enabled").(bool))
	}
	if d.HasChange("min_tls_version") {
		props.MinimumTlsVersion = pointer.To(storageaccounts.MinimumTlsVersion(d.Get("min_tls_version").(string)))
	}
	if d.HasChange("network_rules") {
		props.NetworkAcls = expandAccountNetworkRules(d.Get("network_rules").([]interface{}), tenantId)
	}
	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := storageaccounts.PublicNetworkAccessDisabled
		if d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = storageaccounts.PublicNetworkAccessEnabled
		}
		props.PublicNetworkAccess = pointer.To(publicNetworkAccess)
	}
	if d.HasChange("routing") {
		props.RoutingPreference = expandAccountRoutingPreference(d.Get("routing").([]interface{}))
	}
	if d.HasChange("sas_policy") {
		// TODO: Currently, there is no way to represent a `null` value in the payload - instead it will be omitted, `sas_policy` can not be disabled once enabled.
		props.SasPolicy = expandAccountSASPolicy(d.Get("sas_policy").([]interface{}))
	}
	if d.HasChange("sftp_enabled") {
		props.IsSftpEnabled = pointer.To(d.Get("sftp_enabled").(bool))
	}

	payload := storageaccounts.StorageAccountCreateParameters{
		ExtendedLocation: existing.Model.ExtendedLocation,
		Kind:             *existing.Model.Kind,
		Location:         existing.Model.Location,
		Identity:         existing.Model.Identity,
		Properties:       &props,
		Sku:              *existing.Model.Sku,
		Tags:             existing.Model.Tags,
	}

	// ensure any top-level properties are updated
	if d.HasChange("account_kind") {
		payload.Kind = accountKind
	}
	if d.HasChange("account_replication_type") {
		// storageType is derived from "account_replication_type" and "account_tier" (force-new)
		payload.Sku = storageaccounts.Sku{
			Name: storageaccounts.SkuName(storageType),
		}
	}
	if d.HasChange("identity") {
		payload.Identity = expandedIdentity
	}
	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// azure_files_authentication must be the last to be updated, cause it'll occupy the storage account for several minutes after receiving the response 200 OK. Issue: https://github.com/Azure/azure-rest-api-specs/issues/11272
	if d.HasChange("azure_files_authentication") {
		// due to service issue: https://github.com/Azure/azure-rest-api-specs/issues/12473, we need to update to None before changing its DirectoryServiceOptions
		old, new := d.GetChange("azure_files_authentication.0.directory_type")
		if old != new && new != string(storageaccounts.DirectoryServiceOptionsNone) {
			log.Print("[DEBUG] Disabling AzureFilesIdentityBasedAuthentication prior to changing DirectoryServiceOptions")
			dsNone := storageaccounts.StorageAccountUpdateParameters{
				Properties: &storageaccounts.StorageAccountPropertiesUpdateParameters{
					AzureFilesIdentityBasedAuthentication: &storageaccounts.AzureFilesIdentityBasedAuthentication{
						DirectoryServiceOptions: storageaccounts.DirectoryServiceOptionsNone,
					},
				},
			}
			if _, err := client.Update(ctx, *id, dsNone); err != nil {
				return fmt.Errorf("updating `azure_files_authentication` for %s: %+v", *id, err)
			}
		}

		expandAADFilesAuthentication, err := expandAccountAzureFilesAuthentication(d.Get("azure_files_authentication").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `azure_files_authentication`: %+v", err)
		}
		opts := storageaccounts.StorageAccountUpdateParameters{
			Properties: &storageaccounts.StorageAccountPropertiesUpdateParameters{
				AzureFilesIdentityBasedAuthentication: expandAADFilesAuthentication,
			},
		}

		if _, err := client.Update(ctx, *id, opts); err != nil {
			return fmt.Errorf("updating `azure_files_authentication` for %s: %+v", *id, err)
		}
	}

	// Followings are updates to the sub-services
	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, replicationType)

	if d.HasChange("blob_properties") {
		if !supportLevel.supportBlob {
			return fmt.Errorf("`blob_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		blobProperties, err := expandAccountBlobServiceProperties(accountKind, d.Get("blob_properties").([]interface{}))
		if err != nil {
			return err
		}

		if blobProperties.Properties.IsVersioningEnabled != nil && *blobProperties.Properties.IsVersioningEnabled && d.Get("is_hns_enabled").(bool) {
			return fmt.Errorf("`versioning_enabled` can't be true when `is_hns_enabled` is true")
		}

		// Disable restore_policy first. Disabling restore_policy and while setting delete_retention_policy.allow_permanent_delete to true cause error.
		// Issue : https://github.com/Azure/azure-rest-api-specs/issues/11237
		if v := d.Get("blob_properties.0.restore_policy"); d.HasChange("blob_properties.0.restore_policy") && len(v.([]interface{})) == 0 {
			log.Print("[DEBUG] Disabling RestorePolicy prior to changing DeleteRetentionPolicy")
			blobPayload := blobservice.BlobServiceProperties{
				Properties: &blobservice.BlobServicePropertiesProperties{
					RestorePolicy: expandAccountBlobPropertiesRestorePolicy(v.([]interface{})),
				},
			}
			if _, err := storageClient.BlobService.SetServiceProperties(ctx, *id, blobPayload); err != nil {
				return fmt.Errorf("updating Azure Storage Account blob restore policy %q: %+v", id.StorageAccountName, err)
			}
		}

		if d.Get("dns_endpoint_type").(string) == string(storageaccounts.DnsEndpointTypeAzureDnsZone) {
			if blobProperties.Properties.RestorePolicy != nil && blobProperties.Properties.RestorePolicy.Enabled {
				// Otherwise, API returns: "Required feature Global Dns is disabled"
				// This is confirmed with the SRP team, where they said:
				// > restorePolicy feature is incompatible with partitioned DNS
				return fmt.Errorf("`blob_properties.restore_policy` can't be set when `dns_endpoint_type` is set to `%s`", storageaccounts.DnsEndpointTypeAzureDnsZone)
			}
		}

		if _, err = storageClient.BlobService.SetServiceProperties(ctx, *id, *blobProperties); err != nil {
			return fmt.Errorf("updating `blob_properties` for %s: %+v", *id, err)
		}
	}

	if !features.FivePointOhBeta() {
		dataPlaneClient := meta.(*clients.Client).Storage
		if d.HasChange("queue_properties") {
			if !supportLevel.supportQueue {
				return fmt.Errorf("`queue_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
			}

			account, err := dataPlaneClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if account == nil {
				return fmt.Errorf("unable to locate %s", *id)
			}

			queueClient, err := dataPlaneClient.QueuesDataPlaneClient(ctx, *account, dataPlaneClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Queues Client: %s", err)
			}

			queueProperties, err := expandAccountQueueProperties(d.Get("queue_properties").([]interface{}))
			if err != nil {
				return fmt.Errorf("expanding `queue_properties` for %s: %+v", *id, err)
			}

			if err = queueClient.UpdateServiceProperties(ctx, *queueProperties); err != nil {
				return fmt.Errorf("updating Queue Properties for %s: %+v", *id, err)
			}
		}

		if d.HasChange("static_website") {
			if !supportLevel.supportStaticWebsite {
				return fmt.Errorf("`static_website` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
			}

			account, err := dataPlaneClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if account == nil {
				return fmt.Errorf("unable to locate %s", *id)
			}

			accountsClient, err := dataPlaneClient.AccountsDataPlaneClient(ctx, *account, dataPlaneClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Data Plane client for %s: %+v", *id, err)
			}

			staticWebsiteProps := expandAccountStaticWebsiteProperties(d.Get("static_website").([]interface{}))

			if _, err = accountsClient.SetServiceProperties(ctx, id.StorageAccountName, staticWebsiteProps); err != nil {
				return fmt.Errorf("updating `static_website` for %s: %+v", *id, err)
			}
		}
	}

	if d.HasChange("share_properties") {
		if !supportLevel.supportShare {
			return fmt.Errorf("`share_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		sharePayload := expandAccountShareProperties(d.Get("share_properties").([]interface{}))
		// The API complains if any multichannel info is sent on non premium fileshares. Even if multichannel is set to false
		if accountTier != storageaccounts.SkuTierPremium {
			// Error if the user has tried to enable multichannel on a standard tier storage account
			if sharePayload.Properties.ProtocolSettings.Smb.Multichannel != nil && sharePayload.Properties.ProtocolSettings.Smb.Multichannel.Enabled != nil {
				if *sharePayload.Properties.ProtocolSettings.Smb.Multichannel.Enabled {
					return fmt.Errorf("`multichannel_enabled` isn't supported for Standard tier Storage accounts")
				}
			}

			sharePayload.Properties.ProtocolSettings.Smb.Multichannel = nil
		}

		if _, err = storageClient.FileService.SetServiceProperties(ctx, *id, sharePayload); err != nil {
			return fmt.Errorf("updating File Share Properties for %s: %+v", *id, err)
		}
	}

	return resourceStorageAccountRead(d, meta)
}

func resourceStorageAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageUtils := meta.(*clients.Client).Storage
	storageClient := meta.(*clients.Client).Storage.ResourceManager
	client := storageClient.StorageAccounts
	dataPlaneAvailable := meta.(*clients.Client).Features.Storage.DataPlaneAvailable
	env := meta.(*clients.Client).Account.Environment
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageDomainSuffix, ok := meta.(*clients.Client).Account.Environment.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Storage domain suffix for environment %q", meta.(*clients.Client).Account.Environment.Name)
	}

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	// we then need to find the storage account
	account, err := storageUtils.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	listKeysOpts := storageaccounts.DefaultListKeysOperationOptions()
	listKeysOpts.Expand = pointer.To(storageaccounts.ListKeyExpandKerb)
	keys, err := client.ListKeys(ctx, *id, listKeysOpts)
	if err != nil {
		hasWriteLock := response.WasConflict(keys.HttpResponse)
		doesntHavePermissions := response.WasForbidden(keys.HttpResponse) || response.WasStatusCode(keys.HttpResponse, http.StatusUnauthorized)
		if !hasWriteLock && !doesntHavePermissions {
			return fmt.Errorf("listing Keys for %s: %+v", id, err)
		}
	}

	d.Set("name", id.StorageAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	supportLevel := storageAccountServiceSupportLevel{
		supportBlob:          false,
		supportQueue:         false,
		supportShare:         false,
		supportStaticWebsite: false,
	}
	var accountKind storageaccounts.Kind
	var primaryEndpoints *storageaccounts.Endpoints
	var secondaryEndpoints *storageaccounts.Endpoints
	var routingPreference *storageaccounts.RoutingPreference
	if model := resp.Model; model != nil {
		if model.Kind != nil {
			accountKind = *model.Kind
		}
		d.Set("account_kind", string(accountKind))

		var accountTier storageaccounts.SkuTier
		accountReplicationType := ""
		if sku := model.Sku; sku != nil {
			accountReplicationType = strings.Split(string(sku.Name), "_")[1]
			if sku.Tier != nil {
				accountTier = *sku.Tier
			}
		}
		d.Set("account_tier", string(accountTier))
		d.Set("account_replication_type", accountReplicationType)

		d.Set("edge_zone", flattenEdgeZone(model.ExtendedLocation))
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			primaryEndpoints = props.PrimaryEndpoints
			routingPreference = props.RoutingPreference
			secondaryEndpoints = props.SecondaryEndpoints

			d.Set("access_tier", pointer.From(props.AccessTier))
			d.Set("allowed_copy_scope", pointer.From(props.AllowedCopyScope))
			if err := d.Set("azure_files_authentication", flattenAccountAzureFilesAuthentication(props.AzureFilesIdentityBasedAuthentication)); err != nil {
				return fmt.Errorf("setting `azure_files_authentication`: %+v", err)
			}
			d.Set("cross_tenant_replication_enabled", pointer.From(props.AllowCrossTenantReplication))
			d.Set("https_traffic_only_enabled", pointer.From(props.SupportsHTTPSTrafficOnly))
			d.Set("is_hns_enabled", pointer.From(props.IsHnsEnabled))
			d.Set("nfsv3_enabled", pointer.From(props.IsNfsV3Enabled))
			d.Set("primary_location", pointer.From(props.PrimaryLocation))
			if err := d.Set("routing", flattenAccountRoutingPreference(props.RoutingPreference)); err != nil {
				return fmt.Errorf("setting `routing`: %+v", err)
			}
			d.Set("secondary_location", pointer.From(props.SecondaryLocation))
			d.Set("sftp_enabled", pointer.From(props.IsSftpEnabled))

			// NOTE: The Storage API returns `null` rather than the default value in the API response for existing
			// resources when a new field gets added - meaning we need to default the values below.
			allowBlobPublicAccess := true
			if props.AllowBlobPublicAccess != nil {
				allowBlobPublicAccess = *props.AllowBlobPublicAccess
			}
			d.Set("allow_nested_items_to_be_public", allowBlobPublicAccess)

			defaultToOAuthAuthentication := false
			if props.DefaultToOAuthAuthentication != nil {
				defaultToOAuthAuthentication = *props.DefaultToOAuthAuthentication
			}
			d.Set("default_to_oauth_authentication", defaultToOAuthAuthentication)

			dnsEndpointType := storageaccounts.DnsEndpointTypeStandard
			if props.DnsEndpointType != nil {
				dnsEndpointType = *props.DnsEndpointType
			}
			d.Set("dns_endpoint_type", dnsEndpointType)

			isLocalEnabled := true
			if props.IsLocalUserEnabled != nil {
				isLocalEnabled = *props.IsLocalUserEnabled
			}
			d.Set("local_user_enabled", isLocalEnabled)

			largeFileShareEnabled := false
			if props.LargeFileSharesState != nil {
				largeFileShareEnabled = *props.LargeFileSharesState == storageaccounts.LargeFileSharesStateEnabled
			}
			d.Set("large_file_share_enabled", largeFileShareEnabled)

			minTlsVersion := string(storageaccounts.MinimumTlsVersionTLSOneZero)
			if props.MinimumTlsVersion != nil {
				minTlsVersion = string(*props.MinimumTlsVersion)
			}
			d.Set("min_tls_version", minTlsVersion)

			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == storageaccounts.PublicNetworkAccessDisabled {
				publicNetworkAccessEnabled = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

			allowSharedKeyAccess := true
			if props.AllowSharedKeyAccess != nil {
				allowSharedKeyAccess = *props.AllowSharedKeyAccess
			}
			d.Set("shared_access_key_enabled", allowSharedKeyAccess)

			if err := d.Set("custom_domain", flattenAccountCustomDomain(props.CustomDomain)); err != nil {
				return fmt.Errorf("setting `custom_domain`: %+v", err)
			}
			if err := d.Set("immutability_policy", flattenAccountImmutabilityPolicy(props.ImmutableStorageWithVersioning)); err != nil {
				return fmt.Errorf("setting `immutability_policy`: %+v", err)
			}
			if err := d.Set("network_rules", flattenAccountNetworkRules(props.NetworkAcls)); err != nil {
				return fmt.Errorf("setting `network_rules`: %+v", err)
			}

			// When the encryption key type is "Service", the queue/table is not returned in the service list, so we default
			// the encryption key type to "Service" if it is absent (must also be the default value for "Service" in the schema)
			infrastructureEncryption := false
			queueEncryptionKeyType := string(storageaccounts.KeyTypeService)
			tableEncryptionKeyType := string(storageaccounts.KeyTypeService)
			if encryption := props.Encryption; encryption != nil {
				infrastructureEncryption = pointer.From(encryption.RequireInfrastructureEncryption)
				if encryption.Services != nil {
					if encryption.Services.Queue != nil && encryption.Services.Queue.KeyType != nil {
						queueEncryptionKeyType = string(*encryption.Services.Queue.KeyType)
					}
					if encryption.Services.Table != nil && encryption.Services.Table.KeyType != nil {
						tableEncryptionKeyType = string(*encryption.Services.Table.KeyType)
					}
				}
			}
			d.Set("infrastructure_encryption_enabled", infrastructureEncryption)
			d.Set("queue_encryption_key_type", queueEncryptionKeyType)
			d.Set("table_encryption_key_type", tableEncryptionKeyType)

			customerManagedKey := flattenAccountCustomerManagedKey(props.Encryption, env)
			if err := d.Set("customer_managed_key", customerManagedKey); err != nil {
				return fmt.Errorf("setting `customer_managed_key`: %+v", err)
			}

			if err := d.Set("sas_policy", flattenAccountSASPolicy(props.SasPolicy)); err != nil {
				return fmt.Errorf("setting `sas_policy`: %+v", err)
			}

			supportLevel = availableFunctionalityForAccount(accountKind, accountTier, accountReplicationType)
		}

		flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	endpoints := flattenAccountEndpoints(primaryEndpoints, secondaryEndpoints, routingPreference)
	endpoints.set(d)

	storageAccountKeys := make([]storageaccounts.StorageAccountKey, 0)
	if keys.Model != nil && keys.Model.Keys != nil {
		storageAccountKeys = *keys.Model.Keys
	}
	keysAndConnectionStrings := flattenAccountAccessKeysAndConnectionStrings(id.StorageAccountName, *storageDomainSuffix, storageAccountKeys, endpoints)
	keysAndConnectionStrings.set(d)

	blobProperties := make([]interface{}, 0)
	if supportLevel.supportBlob {
		blobProps, err := storageClient.BlobService.GetServiceProperties(ctx, *id)
		if err != nil {
			return fmt.Errorf("reading blob properties for %s: %+v", *id, err)
		}

		blobProperties = flattenAccountBlobServiceProperties(blobProps.Model)
	}
	if err := d.Set("blob_properties", blobProperties); err != nil {
		return fmt.Errorf("setting `blob_properties` for %s: %+v", *id, err)
	}

	shareProperties := make([]interface{}, 0)
	if supportLevel.supportShare {
		shareProps, err := storageClient.FileService.GetServiceProperties(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving share properties for %s: %+v", *id, err)
		}

		shareProperties = flattenAccountShareProperties(shareProps.Model)
	}
	if err := d.Set("share_properties", shareProperties); err != nil {
		return fmt.Errorf("setting `share_properties` for %s: %+v", *id, err)
	}

	if !features.FivePointOhBeta() && dataPlaneAvailable {
		dataPlaneClient := meta.(*clients.Client).Storage
		queueProperties := make([]interface{}, 0)
		if supportLevel.supportQueue {
			queueClient, err := dataPlaneClient.QueuesDataPlaneClient(ctx, *account, dataPlaneClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Queues Client: %s", err)
			}

			queueProps, err := queueClient.GetServiceProperties(ctx)
			if err != nil {
				// Queue properties is a data plane only service, so we tolerate connection errors here in case of
				// firewalls and other connectivity issues that are not guaranteed.
				if !connectionError(err) {
					return fmt.Errorf("retrieving queue properties for %s: %+v", *id, err)
				}
			}

			queueProperties = flattenAccountQueueProperties(queueProps)
		}
		if err := d.Set("queue_properties", queueProperties); err != nil {
			return fmt.Errorf("setting `queue_properties`: %+v", err)
		}

		staticWebsiteProperties := make([]interface{}, 0)
		if supportLevel.supportStaticWebsite {
			accountsClient, err := dataPlaneClient.AccountsDataPlaneClient(ctx, *account, dataPlaneClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("building Accounts Data Plane Client: %s", err)
			}

			staticWebsiteProps, err := accountsClient.GetServiceProperties(ctx, id.StorageAccountName)
			if err != nil {
				if !connectionError(err) {
					return fmt.Errorf("retrieving static website properties for %s: %+v", *id, err)
				}
			}

			staticWebsiteProperties = flattenAccountStaticWebsiteProperties(staticWebsiteProps)
		}

		if err = d.Set("static_website", staticWebsiteProperties); err != nil {
			return fmt.Errorf("setting `static_website`: %+v", err)
		}
	}

	return nil
}

func resourceStorageAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageUtils := meta.(*clients.Client).Storage
	client := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	existing, err := client.GetProperties(ctx, *id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// the networking api's only allow a single change to be made to a network layout at once, so let's lock to handle that
	virtualNetworkNames := make([]string, 0)
	if model := existing.Model; model != nil && model.Properties != nil {
		if acls := model.Properties.NetworkAcls; acls != nil {
			if vnr := acls.VirtualNetworkRules; vnr != nil {
				for _, v := range *vnr {
					subnetId, err := commonids.ParseSubnetIDInsensitively(v.Id)
					if err != nil {
						return err
					}

					networkName := subnetId.VirtualNetworkName
					for _, virtualNetworkName := range virtualNetworkNames {
						if networkName == virtualNetworkName {
							continue
						}
					}
					virtualNetworkNames = append(virtualNetworkNames, networkName)
				}
			}
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// remove this from the cache
	storageUtils.RemoveAccountFromCache(*id)

	return nil
}

func expandAccountCustomDomain(input []interface{}) *storageaccounts.CustomDomain {
	if len(input) == 0 {
		return &storageaccounts.CustomDomain{
			Name: "",
		}
	}

	domain := input[0].(map[string]interface{})
	return &storageaccounts.CustomDomain{
		Name:             domain["name"].(string),
		UseSubDomainName: pointer.To(domain["use_subdomain"].(bool)),
	}
}

func flattenAccountCustomDomain(input *storageaccounts.CustomDomain) []interface{} {
	output := make([]interface{}, 0)
	if input != nil {
		output = append(output, map[string]interface{}{
			// use_subdomain isn't returned
			"name": input.Name,
		})
	}
	return output
}

func expandAccountCustomerManagedKey(ctx context.Context, keyVaultClient *keyVaultsClient.Client, subscriptionId string, input []interface{}, accountTier storageaccounts.SkuTier, accountKind storageaccounts.Kind, expandedIdentity identity.LegacySystemAndUserAssignedMap, queueEncryptionKeyType, tableEncryptionKeyType storageaccounts.KeyType) (*storageaccounts.Encryption, error) {
	if accountKind == storageaccounts.KindStorage {
		if queueEncryptionKeyType == storageaccounts.KeyTypeAccount {
			return nil, fmt.Errorf("`queue_encryption_key_type = %q` cannot be used with account kind `%q`", string(storageaccounts.KeyTypeAccount), string(storageaccounts.KindStorage))
		}
		if tableEncryptionKeyType == storageaccounts.KeyTypeAccount {
			return nil, fmt.Errorf("`table_encryption_key_type = %q` cannot be used with account kind `%q`", string(storageaccounts.KeyTypeAccount), string(storageaccounts.KindStorage))
		}
	}
	if len(input) == 0 {
		return &storageaccounts.Encryption{
			KeySource: pointer.To(storageaccounts.KeySourceMicrosoftPointStorage),
			Services: &storageaccounts.EncryptionServices{
				Queue: &storageaccounts.EncryptionService{
					KeyType: pointer.To(queueEncryptionKeyType),
				},
				Table: &storageaccounts.EncryptionService{
					KeyType: pointer.To(tableEncryptionKeyType),
				},
			},
		}, nil
	}

	if accountTier != storageaccounts.SkuTierPremium && accountKind != storageaccounts.KindStorageVTwo {
		return nil, fmt.Errorf("customer managed key can only be used with account kind `StorageV2` or account tier `Premium`")
	}

	if expandedIdentity.Type != identity.TypeUserAssigned && expandedIdentity.Type != identity.TypeSystemAssignedUserAssigned {
		return nil, fmt.Errorf("customer managed key can only be configured when the storage account uses a `UserAssigned` or `SystemAssigned, UserAssigned` managed identity but got %q", string(expandedIdentity.Type))
	}

	v := input[0].(map[string]interface{})

	var keyName, keyVersion, keyVaultURI *string
	if keyVaultKeyId, ok := v["key_vault_key_id"]; ok && keyVaultKeyId != "" {
		keyId, err := keyVaultParse.ParseOptionallyVersionedNestedItemID(keyVaultKeyId.(string))
		if err != nil {
			return nil, err
		}

		subscriptionResourceId := commonids.NewSubscriptionID(subscriptionId)
		keyVaultIdRaw, err := keyVaultClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, keyId.KeyVaultBaseUrl)
		if err != nil {
			return nil, err
		}
		if keyVaultIdRaw == nil {
			return nil, fmt.Errorf("unable to find the Resource Manager ID for the Key Vault URI %q in %s", keyId.KeyVaultBaseUrl, subscriptionResourceId)
		}
		keyVaultId, err := commonids.ParseKeyVaultID(*keyVaultIdRaw)
		if err != nil {
			return nil, err
		}

		vaultsClient := keyVaultClient.VaultsClient
		keyVault, err := vaultsClient.Get(ctx, *keyVaultId)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", *keyVaultId, err)
		}

		softDeleteEnabled := false
		purgeProtectionEnabled := false
		if model := keyVault.Model; model != nil {
			if esd := model.Properties.EnableSoftDelete; esd != nil {
				softDeleteEnabled = *esd
			}
			if epp := model.Properties.EnablePurgeProtection; epp != nil {
				purgeProtectionEnabled = *epp
			}
		}
		if !softDeleteEnabled || !purgeProtectionEnabled {
			return nil, fmt.Errorf("%s must be configured for both Purge Protection and Soft Delete", *keyVaultId)
		}

		keyName = pointer.To(keyId.Name)
		keyVersion = pointer.To(keyId.Version)
		keyVaultURI = pointer.To(keyId.KeyVaultBaseUrl)
	} else if managedHSMKeyId, ok := v["managed_hsm_key_id"]; ok && managedHSMKeyId != "" {
		if keyId, err := managedHsmParse.ManagedHSMDataPlaneVersionedKeyID(managedHSMKeyId.(string), nil); err == nil {
			keyName = pointer.To(keyId.KeyName)
			keyVersion = pointer.To(keyId.KeyVersion)
			keyVaultURI = pointer.To(keyId.BaseUri())
		} else if keyId, err := managedHsmParse.ManagedHSMDataPlaneVersionlessKeyID(managedHSMKeyId.(string), nil); err == nil {
			keyName = utils.String(keyId.KeyName)
			keyVersion = utils.String("")
			keyVaultURI = utils.String(keyId.BaseUri())
		} else {
			return nil, fmt.Errorf("parsing %q as HSM key ID", managedHSMKeyId.(string))
		}
	}

	encryption := &storageaccounts.Encryption{
		Services: &storageaccounts.EncryptionServices{
			Blob: &storageaccounts.EncryptionService{
				Enabled: pointer.To(true),
				KeyType: pointer.To(storageaccounts.KeyTypeAccount),
			},
			File: &storageaccounts.EncryptionService{
				Enabled: pointer.To(true),
				KeyType: pointer.To(storageaccounts.KeyTypeAccount),
			},
			Queue: &storageaccounts.EncryptionService{
				KeyType: pointer.To(queueEncryptionKeyType),
			},
			Table: &storageaccounts.EncryptionService{
				KeyType: pointer.To(tableEncryptionKeyType),
			},
		},
		Identity: &storageaccounts.EncryptionIdentity{
			UserAssignedIdentity: utils.String(v["user_assigned_identity_id"].(string)),
		},
		KeySource: pointer.To(storageaccounts.KeySourceMicrosoftPointKeyvault),
		Keyvaultproperties: &storageaccounts.KeyVaultProperties{
			Keyname:     keyName,
			Keyversion:  keyVersion,
			Keyvaulturi: keyVaultURI,
		},
	}

	return encryption, nil
}

func flattenAccountCustomerManagedKey(input *storageaccounts.Encryption, env environments.Environment) []interface{} {
	output := make([]interface{}, 0)

	if input != nil && input.KeySource != nil && *input.KeySource == storageaccounts.KeySourceMicrosoftPointKeyvault {
		userAssignedIdentityId := ""
		if props := input.Identity; props != nil {
			userAssignedIdentityId = pointer.From(props.UserAssignedIdentity)
		}

		customerManagedKey := flattenCustomerManagedKey(input.Keyvaultproperties, env.KeyVault, env.ManagedHSM)
		output = append(output, map[string]interface{}{
			"key_vault_key_id":          customerManagedKey.keyVaultKeyUri,
			"managed_hsm_key_id":        customerManagedKey.managedHsmKeyUri,
			"user_assigned_identity_id": userAssignedIdentityId,
		})
	}

	return output
}

func expandAccountImmutabilityPolicy(input []interface{}) *storageaccounts.ImmutableStorageAccount {
	if len(input) == 0 {
		return &storageaccounts.ImmutableStorageAccount{}
	}

	v := input[0].(map[string]interface{})
	return &storageaccounts.ImmutableStorageAccount{
		Enabled: utils.Bool(true),
		ImmutabilityPolicy: &storageaccounts.AccountImmutabilityPolicyProperties{
			AllowProtectedAppendWrites:            pointer.To(v["allow_protected_append_writes"].(bool)),
			ImmutabilityPeriodSinceCreationInDays: pointer.To(int64(v["period_since_creation_in_days"].(int))),
			State:                                 pointer.To(storageaccounts.AccountImmutabilityPolicyState(v["state"].(string))),
		},
	}
}

func flattenAccountImmutabilityPolicy(input *storageaccounts.ImmutableStorageAccount) []interface{} {
	if input == nil || input.ImmutabilityPolicy == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"allow_protected_append_writes": input.ImmutabilityPolicy.AllowProtectedAppendWrites,
			"period_since_creation_in_days": input.ImmutabilityPolicy.ImmutabilityPeriodSinceCreationInDays,
			"state":                         input.ImmutabilityPolicy.State,
		},
	}
}

func expandAccountActiveDirectoryProperties(input []interface{}) *storageaccounts.ActiveDirectoryProperties {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})

	output := &storageaccounts.ActiveDirectoryProperties{
		DomainGuid: m["domain_guid"].(string),
		DomainName: m["domain_name"].(string),
	}
	if v := m["storage_sid"]; v != "" {
		output.AzureStorageSid = utils.String(v.(string))
	}
	if v := m["domain_sid"]; v != "" {
		output.DomainSid = utils.String(v.(string))
	}
	if v := m["forest_name"]; v != "" {
		output.ForestName = utils.String(v.(string))
	}
	if v := m["netbios_domain_name"]; v != "" {
		output.NetBiosDomainName = utils.String(v.(string))
	}
	return output
}

func flattenAccountActiveDirectoryProperties(input *storageaccounts.ActiveDirectoryProperties) []interface{} {
	output := make([]interface{}, 0)
	if input != nil {
		output = append(output, map[string]interface{}{
			"domain_guid":         input.DomainGuid,
			"domain_name":         input.DomainName,
			"domain_sid":          pointer.From(input.DomainSid),
			"forest_name":         pointer.From(input.ForestName),
			"netbios_domain_name": pointer.From(input.NetBiosDomainName),
			"storage_sid":         pointer.From(input.AzureStorageSid),
		})
	}
	return output
}

func expandAccountAzureFilesAuthentication(input []interface{}) (*storageaccounts.AzureFilesIdentityBasedAuthentication, error) {
	if len(input) == 0 {
		return &storageaccounts.AzureFilesIdentityBasedAuthentication{
			DirectoryServiceOptions: storageaccounts.DirectoryServiceOptionsNone,
		}, nil
	}

	v := input[0].(map[string]interface{})
	output := storageaccounts.AzureFilesIdentityBasedAuthentication{
		DirectoryServiceOptions: storageaccounts.DirectoryServiceOptions(v["directory_type"].(string)),
	}
	if output.DirectoryServiceOptions == storageaccounts.DirectoryServiceOptionsAD ||
		output.DirectoryServiceOptions == storageaccounts.DirectoryServiceOptionsAADDS ||
		output.DirectoryServiceOptions == storageaccounts.DirectoryServiceOptionsAADKERB {
		ad := expandAccountActiveDirectoryProperties(v["active_directory"].([]interface{}))

		if output.DirectoryServiceOptions == storageaccounts.DirectoryServiceOptionsAD {
			if ad == nil {
				return nil, fmt.Errorf("`active_directory` is required when `directory_type` is `AD`")
			}
			if ad.AzureStorageSid == nil {
				return nil, fmt.Errorf("`active_directory.0.storage_sid` is required when `directory_type` is `AD`")
			}
			if ad.DomainSid == nil {
				return nil, fmt.Errorf("`active_directory.0.domain_sid` is required when `directory_type` is `AD`")
			}
			if ad.ForestName == nil {
				return nil, fmt.Errorf("`active_directory.0.forest_name` is required when `directory_type` is `AD`")
			}
			if ad.NetBiosDomainName == nil {
				return nil, fmt.Errorf("`active_directory.0.netbios_domain_name` is required when `directory_type` is `AD`")
			}
		}

		output.ActiveDirectoryProperties = ad
		output.DefaultSharePermission = pointer.To(storageaccounts.DefaultSharePermission(v["default_share_level_permission"].(string)))
	}

	return &output, nil
}

func flattenAccountAzureFilesAuthentication(input *storageaccounts.AzureFilesIdentityBasedAuthentication) []interface{} {
	if input == nil || input.DirectoryServiceOptions == storageaccounts.DirectoryServiceOptionsNone {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"active_directory":               flattenAccountActiveDirectoryProperties(input.ActiveDirectoryProperties),
			"directory_type":                 input.DirectoryServiceOptions,
			"default_share_level_permission": input.DefaultSharePermission,
		},
	}
}

func expandAccountRoutingPreference(input []interface{}) *storageaccounts.RoutingPreference {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &storageaccounts.RoutingPreference{
		PublishMicrosoftEndpoints: pointer.To(v["publish_microsoft_endpoints"].(bool)),
		PublishInternetEndpoints:  pointer.To(v["publish_internet_endpoints"].(bool)),
		RoutingChoice:             pointer.To(storageaccounts.RoutingChoice(v["choice"].(string))),
	}
}

func flattenAccountRoutingPreference(input *storageaccounts.RoutingPreference) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		routingChoice := ""
		if input.RoutingChoice != nil {
			routingChoice = string(*input.RoutingChoice)
		}

		output = append(output, map[string]interface{}{
			"choice":                      routingChoice,
			"publish_internet_endpoints":  pointer.From(input.PublishInternetEndpoints),
			"publish_microsoft_endpoints": pointer.From(input.PublishMicrosoftEndpoints),
		})
	}

	return output
}

func expandAccountBlobServiceProperties(kind storageaccounts.Kind, input []interface{}) (*blobservice.BlobServiceProperties, error) {
	props := blobservice.BlobServicePropertiesProperties{
		Cors: &blobservice.CorsRules{
			CorsRules: &[]blobservice.CorsRule{},
		},
		DeleteRetentionPolicy: &blobservice.DeleteRetentionPolicy{
			Enabled: utils.Bool(false),
		},
	}

	// `Storage` (v1) kind doesn't support:
	// - LastAccessTimeTrackingPolicy: Confirmed by SRP.
	// - ChangeFeed: See https://learn.microsoft.com/en-us/azure/storage/blobs/storage-blob-change-feed?tabs=azure-portal#enable-and-disable-the-change-feed.
	// - Versioning: See https://learn.microsoft.com/en-us/azure/storage/blobs/versioning-overview#how-blob-versioning-works
	// - Restore Policy: See https://learn.microsoft.com/en-us/azure/storage/blobs/point-in-time-restore-overview#prerequisites-for-point-in-time-restore
	if kind != storageaccounts.KindStorage {
		props.LastAccessTimeTrackingPolicy = &blobservice.LastAccessTimeTrackingPolicy{
			Enable: false,
		}
		props.ChangeFeed = &blobservice.ChangeFeed{
			Enabled: pointer.To(false),
		}
		props.IsVersioningEnabled = pointer.To(false)
	}

	if len(input) > 0 {
		v := input[0].(map[string]interface{})

		deletePolicyRaw := v["delete_retention_policy"].([]interface{})
		props.DeleteRetentionPolicy = expandAccountBlobDeleteRetentionPolicy(deletePolicyRaw)

		containerDeletePolicyRaw := v["container_delete_retention_policy"].([]interface{})
		props.ContainerDeleteRetentionPolicy = expandAccountBlobContainerDeleteRetentionPolicy(containerDeletePolicyRaw)

		corsRaw := v["cors_rule"].([]interface{})
		props.Cors = expandAccountBlobPropertiesCors(corsRaw)

		props.IsVersioningEnabled = pointer.To(v["versioning_enabled"].(bool))

		if version, ok := v["default_service_version"].(string); ok && version != "" {
			props.DefaultServiceVersion = pointer.To(version)
		}

		// `Storage` (v1) kind doesn't support:
		// - LastAccessTimeTrackingPolicy
		// - ChangeFeed
		// - Versioning
		// - RestorePolicy
		lastAccessTimeEnabled := v["last_access_time_enabled"].(bool)
		changeFeedEnabled := v["change_feed_enabled"].(bool)
		changeFeedRetentionInDays := v["change_feed_retention_in_days"].(int)
		restorePolicyRaw := v["restore_policy"].([]interface{})
		versioningEnabled := v["versioning_enabled"].(bool)
		if kind != storageaccounts.KindStorage {
			props.LastAccessTimeTrackingPolicy = &blobservice.LastAccessTimeTrackingPolicy{
				Enable: lastAccessTimeEnabled,
			}
			props.ChangeFeed = &blobservice.ChangeFeed{
				Enabled: pointer.To(changeFeedEnabled),
			}
			if changeFeedRetentionInDays != 0 {
				props.ChangeFeed.RetentionInDays = pointer.To(int64(changeFeedRetentionInDays))
			}
			props.RestorePolicy = expandAccountBlobPropertiesRestorePolicy(restorePolicyRaw)
			props.IsVersioningEnabled = &versioningEnabled
		} else {
			if lastAccessTimeEnabled {
				return nil, fmt.Errorf("`last_access_time_enabled` can not be configured when `kind` is set to `Storage` (v1)")
			}
			if changeFeedEnabled {
				return nil, fmt.Errorf("`change_feed_enabled` can not be configured when `kind` is set to `Storage` (v1)")
			}
			if changeFeedRetentionInDays != 0 {
				return nil, fmt.Errorf("`change_feed_retention_in_days` can not be configured when `kind` is set to `Storage` (v1)")
			}
			if len(restorePolicyRaw) != 0 {
				return nil, fmt.Errorf("`restore_policy` can not be configured when `kind` is set to `Storage` (v1)")
			}
			if versioningEnabled {
				return nil, fmt.Errorf("`versioning_enabled` can not be configured when `kind` is set to `Storage` (v1)")
			}
		}

		// Sanity check for the prerequisites of restore_policy
		// Ref: https://learn.microsoft.com/en-us/azure/storage/blobs/point-in-time-restore-overview#prerequisites-for-point-in-time-restore
		if p := props.RestorePolicy; p != nil && p.Enabled {
			if props.ChangeFeed == nil || props.ChangeFeed.Enabled == nil || !*props.ChangeFeed.Enabled {
				return nil, fmt.Errorf("`change_feed_enabled` must be `true` when `restore_policy` is set")
			}
			if props.IsVersioningEnabled == nil || !*props.IsVersioningEnabled {
				return nil, fmt.Errorf("`versioning_enabled` must be `true` when `restore_policy` is set")
			}
		}
	}

	return &blobservice.BlobServiceProperties{
		Properties: &props,
	}, nil
}

func flattenAccountBlobServiceProperties(input *blobservice.BlobServiceProperties) []interface{} {
	if input == nil || input.Properties == nil {
		return []interface{}{}
	}

	flattenedCorsRules := make([]interface{}, 0)
	if corsRules := input.Properties.Cors; corsRules != nil {
		flattenedCorsRules = flattenAccountBlobPropertiesCorsRule(corsRules)
	}

	flattenedDeletePolicy := make([]interface{}, 0)
	if deletePolicy := input.Properties.DeleteRetentionPolicy; deletePolicy != nil {
		flattenedDeletePolicy = flattenAccountBlobDeleteRetentionPolicy(deletePolicy)
	}

	flattenedRestorePolicy := make([]interface{}, 0)
	if restorePolicy := input.Properties.RestorePolicy; restorePolicy != nil {
		flattenedRestorePolicy = flattenAccountBlobPropertiesRestorePolicy(restorePolicy)
	}

	flattenedContainerDeletePolicy := make([]interface{}, 0)
	if containerDeletePolicy := input.Properties.ContainerDeleteRetentionPolicy; containerDeletePolicy != nil {
		flattenedContainerDeletePolicy = flattenAccountBlobContainerDeleteRetentionPolicy(containerDeletePolicy)
	}

	versioning, changeFeedEnabled, changeFeedRetentionInDays := false, false, 0
	if input.Properties.IsVersioningEnabled != nil {
		versioning = *input.Properties.IsVersioningEnabled
	}

	if v := input.Properties.ChangeFeed; v != nil {
		if v.Enabled != nil {
			changeFeedEnabled = *v.Enabled
		}
		if v.RetentionInDays != nil {
			changeFeedRetentionInDays = int(*v.RetentionInDays)
		}
	}

	var defaultServiceVersion string
	if input.Properties.DefaultServiceVersion != nil {
		defaultServiceVersion = *input.Properties.DefaultServiceVersion
	}

	var LastAccessTimeTrackingPolicy bool
	if v := input.Properties.LastAccessTimeTrackingPolicy; v != nil {
		LastAccessTimeTrackingPolicy = v.Enable
	}

	return []interface{}{
		map[string]interface{}{
			"change_feed_enabled":               changeFeedEnabled,
			"change_feed_retention_in_days":     changeFeedRetentionInDays,
			"container_delete_retention_policy": flattenedContainerDeletePolicy,
			"cors_rule":                         flattenedCorsRules,
			"default_service_version":           defaultServiceVersion,
			"delete_retention_policy":           flattenedDeletePolicy,
			"last_access_time_enabled":          LastAccessTimeTrackingPolicy,
			"restore_policy":                    flattenedRestorePolicy,
			"versioning_enabled":                versioning,
		},
	}
}

func expandAccountBlobDeleteRetentionPolicy(input []interface{}) *blobservice.DeleteRetentionPolicy {
	result := blobservice.DeleteRetentionPolicy{
		Enabled: pointer.To(false),
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]interface{})

	return &blobservice.DeleteRetentionPolicy{
		Enabled:              pointer.To(true),
		AllowPermanentDelete: pointer.To(policy["permanent_delete_enabled"].(bool)),
		Days:                 pointer.To(int64(policy["days"].(int))),
	}
}

func flattenAccountBlobDeleteRetentionPolicy(input *blobservice.DeleteRetentionPolicy) []interface{} {
	deleteRetentionPolicy := make([]interface{}, 0)

	if input == nil {
		return deleteRetentionPolicy
	}

	if enabled := input.Enabled; enabled != nil && *enabled {
		days := 0
		if input.Days != nil {
			days = int(*input.Days)
		}

		var permanentDeleteEnabled bool
		if input.AllowPermanentDelete != nil {
			permanentDeleteEnabled = *input.AllowPermanentDelete
		}

		deleteRetentionPolicy = append(deleteRetentionPolicy, map[string]interface{}{
			"days":                     days,
			"permanent_delete_enabled": permanentDeleteEnabled,
		})
	}

	return deleteRetentionPolicy
}

func expandAccountBlobContainerDeleteRetentionPolicy(input []interface{}) *blobservice.DeleteRetentionPolicy {
	result := blobservice.DeleteRetentionPolicy{
		Enabled: pointer.To(false),
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]interface{})

	return &blobservice.DeleteRetentionPolicy{
		Enabled: pointer.To(true),
		Days:    pointer.To(int64(policy["days"].(int))),
	}
}

func flattenAccountBlobContainerDeleteRetentionPolicy(input *blobservice.DeleteRetentionPolicy) []interface{} {
	deleteRetentionPolicy := make([]interface{}, 0)

	if input == nil {
		return deleteRetentionPolicy
	}

	if enabled := input.Enabled; enabled != nil && *enabled {
		days := 0
		if input.Days != nil {
			days = int(*input.Days)
		}

		deleteRetentionPolicy = append(deleteRetentionPolicy, map[string]interface{}{
			"days": days,
		})
	}

	return deleteRetentionPolicy
}

func expandAccountBlobPropertiesRestorePolicy(input []interface{}) *blobservice.RestorePolicyProperties {
	result := blobservice.RestorePolicyProperties{
		Enabled: false,
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]interface{})

	return &blobservice.RestorePolicyProperties{
		Enabled: true,
		Days:    pointer.To(int64(policy["days"].(int))),
	}
}

func flattenAccountBlobPropertiesRestorePolicy(input *blobservice.RestorePolicyProperties) []interface{} {
	restorePolicy := make([]interface{}, 0)

	if input == nil {
		return restorePolicy
	}

	if enabled := input.Enabled; enabled {
		days := 0
		if input.Days != nil {
			days = int(*input.Days)
		}

		restorePolicy = append(restorePolicy, map[string]interface{}{
			"days": days,
		})
	}

	return restorePolicy
}

func expandAccountBlobPropertiesCors(input []interface{}) *blobservice.CorsRules {
	blobCorsRules := blobservice.CorsRules{}

	if len(input) > 0 {
		corsRules := make([]blobservice.CorsRule, 0)
		for _, raw := range input {
			item := raw.(map[string]interface{})

			allowedMethods := make([]blobservice.AllowedMethods, 0)
			for _, val := range *utils.ExpandStringSlice(item["allowed_methods"].([]interface{})) {
				allowedMethods = append(allowedMethods, blobservice.AllowedMethods(val))
			}
			corsRules = append(corsRules, blobservice.CorsRule{
				AllowedHeaders:  *utils.ExpandStringSlice(item["allowed_headers"].([]interface{})),
				AllowedOrigins:  *utils.ExpandStringSlice(item["allowed_origins"].([]interface{})),
				AllowedMethods:  allowedMethods,
				ExposedHeaders:  *utils.ExpandStringSlice(item["exposed_headers"].([]interface{})),
				MaxAgeInSeconds: int64(item["max_age_in_seconds"].(int)),
			})
		}
		blobCorsRules.CorsRules = &corsRules
	}
	return &blobCorsRules
}

func flattenAccountBlobPropertiesCorsRule(input *blobservice.CorsRules) []interface{} {
	corsRules := make([]interface{}, 0)

	if input == nil || input.CorsRules == nil {
		return corsRules
	}

	for _, corsRule := range *input.CorsRules {
		corsRules = append(corsRules, map[string]interface{}{
			"allowed_headers":    corsRule.AllowedHeaders,
			"allowed_methods":    corsRule.AllowedMethods,
			"allowed_origins":    corsRule.AllowedOrigins,
			"exposed_headers":    corsRule.ExposedHeaders,
			"max_age_in_seconds": int(corsRule.MaxAgeInSeconds),
		})
	}

	return corsRules
}

func expandAccountShareProperties(input []interface{}) fileservice.FileServiceProperties {
	props := fileservice.FileServiceProperties{
		Properties: &fileservice.FileServicePropertiesProperties{
			Cors: &fileservice.CorsRules{
				CorsRules: &[]fileservice.CorsRule{},
			},
			ShareDeleteRetentionPolicy: &fileservice.DeleteRetentionPolicy{
				Enabled: pointer.To(false),
			},
		},
	}

	if len(input) > 0 && input[0] != nil {
		v := input[0].(map[string]interface{})

		props.Properties.ShareDeleteRetentionPolicy = expandAccountShareDeleteRetentionPolicy(v["retention_policy"].([]interface{}))

		props.Properties.Cors = expandAccountSharePropertiesCorsRule(v["cors_rule"].([]interface{}))

		props.Properties.ProtocolSettings = &fileservice.ProtocolSettings{
			Smb: expandAccountSharePropertiesSMB(v["smb"].([]interface{})),
		}
	}

	return props
}

func flattenAccountShareProperties(input *fileservice.FileServiceProperties) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if props := input.Properties; props != nil {
			output = append(output, map[string]interface{}{
				"cors_rule":        flattenAccountSharePropertiesCorsRule(props.Cors),
				"retention_policy": flattenAccountShareDeleteRetentionPolicy(props.ShareDeleteRetentionPolicy),
				"smb":              flattenAccountSharePropertiesSMB(props.ProtocolSettings),
			})
		}
	}

	return output
}

func expandAccountSharePropertiesCorsRule(input []interface{}) *fileservice.CorsRules {
	blobCorsRules := fileservice.CorsRules{}

	if len(input) > 0 {
		corsRules := make([]fileservice.CorsRule, 0)
		for _, raw := range input {
			item := raw.(map[string]interface{})

			allowedMethods := make([]fileservice.AllowedMethods, 0)
			for _, val := range *utils.ExpandStringSlice(item["allowed_methods"].([]interface{})) {
				allowedMethods = append(allowedMethods, fileservice.AllowedMethods(val))
			}
			corsRules = append(corsRules, fileservice.CorsRule{
				AllowedHeaders:  *utils.ExpandStringSlice(item["allowed_headers"].([]interface{})),
				AllowedMethods:  allowedMethods,
				AllowedOrigins:  *utils.ExpandStringSlice(item["allowed_origins"].([]interface{})),
				ExposedHeaders:  *utils.ExpandStringSlice(item["exposed_headers"].([]interface{})),
				MaxAgeInSeconds: int64(item["max_age_in_seconds"].(int)),
			})
		}
		blobCorsRules.CorsRules = &corsRules
	}
	return &blobCorsRules
}

func flattenAccountSharePropertiesCorsRule(input *fileservice.CorsRules) []interface{} {
	corsRules := make([]interface{}, 0)

	if input == nil || input.CorsRules == nil {
		return corsRules
	}

	for _, corsRule := range *input.CorsRules {
		corsRules = append(corsRules, map[string]interface{}{
			"allowed_headers":    corsRule.AllowedHeaders,
			"allowed_methods":    corsRule.AllowedMethods,
			"allowed_origins":    corsRule.AllowedOrigins,
			"exposed_headers":    corsRule.ExposedHeaders,
			"max_age_in_seconds": int(corsRule.MaxAgeInSeconds),
		})
	}

	return corsRules
}

func expandAccountShareDeleteRetentionPolicy(input []interface{}) *fileservice.DeleteRetentionPolicy {
	result := fileservice.DeleteRetentionPolicy{
		Enabled: pointer.To(false),
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]interface{})

	return &fileservice.DeleteRetentionPolicy{
		Enabled: pointer.To(true),
		Days:    pointer.To(int64(policy["days"].(int))),
	}
}

func flattenAccountShareDeleteRetentionPolicy(input *fileservice.DeleteRetentionPolicy) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		if enabled := input.Enabled; enabled != nil && *enabled {
			days := 0
			if input.Days != nil {
				days = int(*input.Days)
			}

			output = append(output, map[string]interface{}{
				"days": days,
			})
		}
	}

	return output
}

func expandAccountSharePropertiesSMB(input []interface{}) *fileservice.SmbSetting {
	if len(input) == 0 || input[0] == nil {
		return &fileservice.SmbSetting{
			AuthenticationMethods:    pointer.To(""),
			ChannelEncryption:        pointer.To(""),
			KerberosTicketEncryption: pointer.To(""),
			Versions:                 pointer.To(""),
			Multichannel:             nil,
		}
	}

	v := input[0].(map[string]interface{})

	return &fileservice.SmbSetting{
		AuthenticationMethods:    utils.ExpandStringSliceWithDelimiter(v["authentication_types"].(*pluginsdk.Set).List(), ";"),
		ChannelEncryption:        utils.ExpandStringSliceWithDelimiter(v["channel_encryption_type"].(*pluginsdk.Set).List(), ";"),
		KerberosTicketEncryption: utils.ExpandStringSliceWithDelimiter(v["kerberos_ticket_encryption_type"].(*pluginsdk.Set).List(), ";"),
		Versions:                 utils.ExpandStringSliceWithDelimiter(v["versions"].(*pluginsdk.Set).List(), ";"),
		Multichannel: &fileservice.Multichannel{
			Enabled: pointer.To(v["multichannel_enabled"].(bool)),
		},
	}
}

func flattenAccountSharePropertiesSMB(input *fileservice.ProtocolSettings) []interface{} {
	if input == nil || input.Smb == nil {
		return []interface{}{}
	}

	versions := make([]interface{}, 0)
	if input.Smb.Versions != nil {
		versions = utils.FlattenStringSliceWithDelimiter(input.Smb.Versions, ";")
	}

	authenticationMethods := make([]interface{}, 0)
	if input.Smb.AuthenticationMethods != nil {
		authenticationMethods = utils.FlattenStringSliceWithDelimiter(input.Smb.AuthenticationMethods, ";")
	}

	kerberosTicketEncryption := make([]interface{}, 0)
	if input.Smb.KerberosTicketEncryption != nil {
		kerberosTicketEncryption = utils.FlattenStringSliceWithDelimiter(input.Smb.KerberosTicketEncryption, ";")
	}

	channelEncryption := make([]interface{}, 0)
	if input.Smb.ChannelEncryption != nil {
		channelEncryption = utils.FlattenStringSliceWithDelimiter(input.Smb.ChannelEncryption, ";")
	}

	multichannelEnabled := false
	if input.Smb.Multichannel != nil && input.Smb.Multichannel.Enabled != nil {
		multichannelEnabled = *input.Smb.Multichannel.Enabled
	}

	if len(versions) == 0 && len(authenticationMethods) == 0 && len(kerberosTicketEncryption) == 0 && len(channelEncryption) == 0 && (input.Smb.Multichannel == nil || input.Smb.Multichannel.Enabled == nil) {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"authentication_types":            authenticationMethods,
			"channel_encryption_type":         channelEncryption,
			"kerberos_ticket_encryption_type": kerberosTicketEncryption,
			"multichannel_enabled":            multichannelEnabled,
			"versions":                        versions,
		},
	}
}

func expandAccountQueueProperties(input []interface{}) (*queues.StorageServiceProperties, error) {
	var err error
	properties := queues.StorageServiceProperties{
		Cors: &queues.Cors{
			CorsRule: []queues.CorsRule{},
		},
		HourMetrics: &queues.MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		},
		MinuteMetrics: &queues.MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		},
		Logging: &queues.LoggingConfig{
			Version: "1.0",
			Delete:  false,
			Read:    false,
			Write:   false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		},
	}
	if len(input) == 0 {
		return &properties, nil
	}

	attrs := input[0].(map[string]interface{})

	properties.Cors = expandAccountQueuePropertiesCors(attrs["cors_rule"].([]interface{}))
	properties.Logging = expandAccountQueuePropertiesLogging(attrs["logging"].([]interface{}))
	properties.MinuteMetrics, err = expandAccountQueuePropertiesMetrics(attrs["minute_metrics"].([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("expanding `minute_metrics`: %+v", err)
	}
	properties.HourMetrics, err = expandAccountQueuePropertiesMetrics(attrs["hour_metrics"].([]interface{}))
	if err != nil {
		return nil, fmt.Errorf("expanding `hour_metrics`: %+v", err)
	}

	return &properties, nil
}

func flattenAccountQueueProperties(input *queues.StorageServiceProperties) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		corsRules := flattenAccountQueuePropertiesCors(input.Cors)
		logging := flattenAccountQueuePropertiesLogging(input.Logging)
		hourMetrics := flattenAccountQueuePropertiesMetrics(input.HourMetrics)
		minuteMetrics := flattenAccountQueuePropertiesMetrics(input.MinuteMetrics)

		if len(corsRules) > 0 || len(logging) > 0 || len(hourMetrics) > 0 || len(minuteMetrics) > 0 {
			output = append(output, map[string]interface{}{
				"cors_rule":      corsRules,
				"hour_metrics":   hourMetrics,
				"logging":        logging,
				"minute_metrics": minuteMetrics,
			})
		}
	}

	return output
}

func expandAccountQueuePropertiesLogging(input []interface{}) *queues.LoggingConfig {
	if len(input) == 0 {
		return &queues.LoggingConfig{
			Version: "1.0",
			Delete:  false,
			Read:    false,
			Write:   false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		}
	}

	loggingAttr := input[0].(map[string]interface{})
	logging := &queues.LoggingConfig{
		Delete:  loggingAttr["delete"].(bool),
		Read:    loggingAttr["read"].(bool),
		Version: loggingAttr["version"].(string),
		Write:   loggingAttr["write"].(bool),
	}

	if v, ok := loggingAttr["retention_policy_days"]; ok {
		if days := v.(int); days > 0 {
			logging.RetentionPolicy = queues.RetentionPolicy{
				Days:    days,
				Enabled: true,
			}
		}
	}

	return logging
}

func flattenAccountQueuePropertiesLogging(input *queues.LoggingConfig) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	retentionPolicyDays := 0
	if input.RetentionPolicy.Enabled {
		retentionPolicyDays = input.RetentionPolicy.Days
	}

	return []interface{}{
		map[string]interface{}{
			"delete":                input.Delete,
			"read":                  input.Read,
			"retention_policy_days": retentionPolicyDays,
			"version":               input.Version,
			"write":                 input.Write,
		},
	}
}

func expandAccountQueuePropertiesMetrics(input []interface{}) (*queues.MetricsConfig, error) {
	if len(input) == 0 {
		return &queues.MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: queues.RetentionPolicy{
				Enabled: false,
			},
		}, nil
	}

	metricsAttr := input[0].(map[string]interface{})

	metrics := &queues.MetricsConfig{
		Enabled: metricsAttr["enabled"].(bool),
		Version: metricsAttr["version"].(string),
	}

	if v, ok := metricsAttr["retention_policy_days"]; ok {
		if days := v.(int); days > 0 {
			metrics.RetentionPolicy = queues.RetentionPolicy{
				Days:    days,
				Enabled: true,
			}
		}
	}

	if v, ok := metricsAttr["include_apis"]; ok {
		includeAPIs := v.(bool)
		if metrics.Enabled {
			metrics.IncludeAPIs = &includeAPIs
		} else if includeAPIs {
			return nil, fmt.Errorf("`include_apis` may only be set when `enabled` is true")
		}
	}

	return metrics, nil
}

func flattenAccountQueuePropertiesMetrics(input *queues.MetricsConfig) []interface{} {
	output := make([]interface{}, 0)

	if input != nil && input.Version != "" {
		retentionPolicyDays := 0
		if input.RetentionPolicy.Enabled {
			retentionPolicyDays = input.RetentionPolicy.Days
		}

		output = append(output, map[string]interface{}{
			"enabled":               input.Enabled,
			"include_apis":          pointer.From(input.IncludeAPIs),
			"retention_policy_days": retentionPolicyDays,
			"version":               input.Version,
		})
	}

	return output
}

func expandAccountQueuePropertiesCors(input []interface{}) *queues.Cors {
	if len(input) == 0 {
		return &queues.Cors{}
	}

	corsRules := make([]queues.CorsRule, 0)
	for _, attr := range input {
		corsRuleAttr := attr.(map[string]interface{})
		corsRule := queues.CorsRule{}

		corsRule.AllowedOrigins = strings.Join(*utils.ExpandStringSlice(corsRuleAttr["allowed_origins"].([]interface{})), ",")
		corsRule.ExposedHeaders = strings.Join(*utils.ExpandStringSlice(corsRuleAttr["exposed_headers"].([]interface{})), ",")
		corsRule.AllowedHeaders = strings.Join(*utils.ExpandStringSlice(corsRuleAttr["allowed_headers"].([]interface{})), ",")
		corsRule.AllowedMethods = strings.Join(*utils.ExpandStringSlice(corsRuleAttr["allowed_methods"].([]interface{})), ",")
		corsRule.MaxAgeInSeconds = corsRuleAttr["max_age_in_seconds"].(int)

		corsRules = append(corsRules, corsRule)
	}

	cors := &queues.Cors{
		CorsRule: corsRules,
	}
	return cors
}

func flattenAccountQueuePropertiesCors(input *queues.Cors) []interface{} {
	output := make([]interface{}, 0)

	if input == nil || len(input.CorsRule) == 0 || input.CorsRule[0].AllowedOrigins == "" {
		return output
	}

	for _, item := range input.CorsRule {
		output = append(output, map[string]interface{}{
			"allowed_headers":    flattenAccountQueuePropertiesCorsRule(item.AllowedHeaders),
			"allowed_methods":    flattenAccountQueuePropertiesCorsRule(item.AllowedMethods),
			"allowed_origins":    flattenAccountQueuePropertiesCorsRule(item.AllowedOrigins),
			"exposed_headers":    flattenAccountQueuePropertiesCorsRule(item.ExposedHeaders),
			"max_age_in_seconds": item.MaxAgeInSeconds,
		})
	}

	return output
}

func flattenAccountQueuePropertiesCorsRule(input string) []interface{} {
	results := make([]interface{}, 0)

	components := strings.Split(input, ",")
	for _, item := range components {
		results = append(results, item)
	}

	return results
}

func expandAccountStaticWebsiteProperties(input []interface{}) accounts.StorageServiceProperties {
	properties := accounts.StorageServiceProperties{
		StaticWebsite: &accounts.StaticWebsite{
			Enabled: false,
		},
	}
	if len(input) == 0 {
		return properties
	}

	properties.StaticWebsite.Enabled = true

	// @tombuildsstuff: this looks weird, doesn't it?
	// Since the presence of this block signifies the website's enabled however all fields within it are optional
	// TF Core returns a nil object when there's no keys defined within the block, rather than an empty map. As
	// such this hack allows us to have a Static Website block with only Enabled configured, without the optional
	// inner properties.
	if val := input[0]; val != nil {
		attr := val.(map[string]interface{})
		if v, ok := attr["index_document"]; ok {
			properties.StaticWebsite.IndexDocument = v.(string)
		}

		if v, ok := attr["error_404_document"]; ok {
			properties.StaticWebsite.ErrorDocument404Path = v.(string)
		}
	}

	return properties
}

func flattenAccountStaticWebsiteProperties(input accounts.GetServicePropertiesResult) []interface{} {
	if staticWebsite := input.StaticWebsite; staticWebsite != nil {
		if !staticWebsite.Enabled {
			return []interface{}{}
		}

		return []interface{}{
			map[string]interface{}{
				"error_404_document": staticWebsite.ErrorDocument404Path,
				"index_document":     staticWebsite.IndexDocument,
			},
		}
	}
	return []interface{}{}
}

func expandAccountSASPolicy(input []interface{}) *storageaccounts.SasPolicy {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})
	return &storageaccounts.SasPolicy{
		ExpirationAction:    storageaccounts.ExpirationAction(raw["expiration_action"].(string)),
		SasExpirationPeriod: raw["expiration_period"].(string),
	}
}

func flattenAccountSASPolicy(input *storageaccounts.SasPolicy) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		output = append(output, map[string]interface{}{
			"expiration_action": string(input.ExpirationAction),
			"expiration_period": input.SasExpirationPeriod,
		})
	}

	return output
}

func expandEdgeZone(input string) *edgezones.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &edgezones.Model{
		Name: normalized,
	}
}

func flattenEdgeZone(input *edgezones.Model) string {
	output := ""
	if input != nil {
		output = edgezones.Normalize(input.Name)
	}
	return output
}
