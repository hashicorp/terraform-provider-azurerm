// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage" // nolint: staticcheck
	azautorest "github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/fileservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	keyVaultClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	managedHsmHelpers "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/helpers"
	managedHsmParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	managedHsmValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/blob/accounts"
	"github.com/tombuildsstuff/giovanni/storage/2023-11-03/queue/queues"
)

var (
	storageAccountResourceName  = "azurerm_storage_account"
	storageKindsSupportsSkuTier = []storage.Kind{
		storage.KindBlobStorage,
		storage.KindStorageV2,
		storage.KindFileStorage,
	}
	storageKindsSupportHns = []storage.Kind{
		storage.KindBlobStorage,
		storage.KindStorageV2,
		storage.KindBlockBlobStorage,
	}
)

type storageAccountServiceSupportLevel struct {
	supportBlob          bool
	supportQueue         bool
	supportShare         bool
	supportStaticWebsite bool
}

func resourceStorageAccount() *pluginsdk.Resource {
	upgraders := map[int]pluginsdk.StateUpgrade{
		0: migration.AccountV0ToV1{},
		1: migration.AccountV1ToV2{},
		2: migration.AccountV2ToV3{},
		3: migration.AccountV3ToV4{},
	}
	schemaVersion := 4

	return &pluginsdk.Resource{
		Create: resourceStorageAccountCreate,
		Read:   resourceStorageAccountRead,
		Update: resourceStorageAccountUpdate,
		Delete: resourceStorageAccountDelete,

		SchemaVersion:  schemaVersion,
		StateUpgraders: pluginsdk.StateUpgrades(upgraders),

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
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.KindStorage),
					string(storage.KindBlobStorage),
					string(storage.KindBlockBlobStorage),
					string(storage.KindFileStorage),
					string(storage.KindStorageV2),
				}, false),
				Default: string(storage.KindStorageV2),
			},

			"account_tier": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.SkuTierStandard),
					string(storage.SkuTierPremium),
				}, false),
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
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.AccessTierCool),
					string(storage.AccessTierHot),
				}, false),
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
								string(storageaccounts.DirectoryServiceOptionsAADKERB),
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
					},
				},
			},
			"cross_tenant_replication_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
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

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_https_traffic_only": {
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
						"period_since_creation_in_days": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},
						"state": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(storage.AccountImmutabilityPolicyStateDisabled),
								string(storage.AccountImmutabilityPolicyStateLocked),
								string(storage.AccountImmutabilityPolicyStateUnlocked),
							}, false),
						},
						"allow_protected_append_writes": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
					},
				},
			},

			"min_tls_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(storage.MinimumTLSVersionTLS12),
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.MinimumTLSVersionTLS10),
					string(storage.MinimumTLSVersionTLS11),
					string(storage.MinimumTLSVersionTLS12),
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

			// TODO: document this new field in 3.0
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
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.DNSEndpointTypeStandard),
					string(storage.DNSEndpointTypeAzureDNSZone),
				}, false),
				Default:  string(storage.DNSEndpointTypeStandard),
				ForceNew: true,
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
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(storage.BypassAzureServices),
									string(storage.BypassLogging),
									string(storage.BypassMetrics),
									string(storage.BypassNone),
								}, false),
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

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

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
									"permanent_delete_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
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
								string(storage.RoutingChoiceMicrosoftRouting),
								string(storage.RoutingChoiceInternetRouting),
							}, false),
							Default: string(storage.RoutingChoiceMicrosoftRouting),
						},
					},
				},
			},

			"share_properties": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				// (@jackofallops) TODO - This should not be computed, however, this would be a breaking change with unknown implications for user data so needs to be addressed for 4.0
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
									"multichannel_enabled": {
										Type:     pluginsdk.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
					},
				},
			},

			// lintignore:XS003
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

			"queue_encryption_key_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.KeyTypeService),
					string(storage.KeyTypeAccount),
				}, false),
				Default: string(storage.KeyTypeService),
			},

			"table_encryption_key_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.KeyTypeService),
					string(storage.KeyTypeAccount),
				}, false),
				Default: string(storage.KeyTypeService),
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
						"expiration_period": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"expiration_action": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "Log",
							ValidateFunc: validation.StringInSlice([]string{
								// There is no definition of this enum in the Track1 SDK due to: https://github.com/Azure/azure-sdk-for-go/issues/14589
								"Log",
							}, false),
						},
					},
				},
			},

			"allowed_copy_scope": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(storage.AllowedCopyScopePrivateLink),
					string(storage.AllowedCopyScopeAAD),
				}, false),
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

					if accountKind != string(storage.KindStorage) && changedKind != string(storage.KindStorageV2) {
						log.Printf("[DEBUG] recreate storage account, could't be migrated from %s to %s", accountKind, changedKind)
						d.ForceNew("account_kind")
					} else {
						log.Printf("[DEBUG] storage account can be upgraded from %s to %s", accountKind, changedKind)
					}
				}

				if d.HasChange("large_file_share_enabled") {
					lfsEnabled, changedEnabled := d.GetChange("large_file_share_enabled")
					if lfsEnabled.(bool) && !changedEnabled.(bool) {
						return fmt.Errorf("`large_file_share_enabled` cannot be disabled once it's been enabled")
					}
				}

				if d.Get("access_tier") != "" {
					if accountKind := storage.Kind(d.Get("account_kind").(string)); !slices.Contains(storageKindsSupportsSkuTier, accountKind) {
						return fmt.Errorf("`access_tier` is only available for accounts of kind: %v", storageKindsSupportsSkuTier)
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
}

func resourceStorageAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	tenantId := meta.(*clients.Client).Account.TenantId
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	storageAccountsClient := meta.(*clients.Client).Storage.ResourceManager.StorageAccounts
	legacyClient := meta.(*clients.Client).Storage.AccountsClient
	storageClient := meta.(*clients.Client).Storage
	keyVaultClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewStorageAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	existing, err := storageAccountsClient.GetProperties(ctx, id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %s", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_storage_account", id.ID())
	}

	accountKind := storage.Kind(d.Get("account_kind").(string))
	isHnsEnabled := d.Get("is_hns_enabled").(bool)
	nfsV3Enabled := d.Get("nfsv3_enabled").(bool)
	defaultToOAuthAuthentication := d.Get("default_to_oauth_authentication").(bool)
	publicNetworkAccess := storage.PublicNetworkAccessDisabled
	if d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = storage.PublicNetworkAccessEnabled
	}
	dnsEndpointType := d.Get("dns_endpoint_type").(string)

	accountTier := storage.SkuTier(d.Get("account_tier").(string))
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)

	parameters := storage.AccountCreateParameters{
		ExtendedLocation: expandEdgeZone(d.Get("edge_zone").(string)),
		Kind:             accountKind,
		Location:         pointer.To(location.Normalize(d.Get("location").(string))),
		AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{
			AllowBlobPublicAccess:        pointer.To(d.Get("allow_nested_items_to_be_public").(bool)),
			AllowCrossTenantReplication:  pointer.To(d.Get("cross_tenant_replication_enabled").(bool)),
			AllowSharedKeyAccess:         pointer.To(d.Get("shared_access_key_enabled").(bool)),
			DNSEndpointType:              storage.DNSEndpointType(dnsEndpointType),
			DefaultToOAuthAuthentication: pointer.To(defaultToOAuthAuthentication),
			EnableHTTPSTrafficOnly:       pointer.To(d.Get("enable_https_traffic_only").(bool)),
			EnableNfsV3:                  pointer.To(nfsV3Enabled),
			IsHnsEnabled:                 pointer.To(isHnsEnabled),
			IsLocalUserEnabled:           pointer.To(d.Get("local_user_enabled").(bool)),
			IsSftpEnabled:                pointer.To(d.Get("sftp_enabled").(bool)),
			MinimumTLSVersion:            storage.MinimumTLSVersion(d.Get("min_tls_version").(string)),
			NetworkRuleSet:               expandStorageAccountNetworkRules(d, tenantId),
			PublicNetworkAccess:          publicNetworkAccess,
			SasPolicy:                    expandStorageAccountSASPolicy(d.Get("sas_policy").([]interface{})),
		},
		Sku: &storage.Sku{
			Name: storage.SkuName(storageType),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v := d.Get("allowed_copy_scope").(string); v != "" {
		parameters.AccountPropertiesCreateParameters.AllowedCopyScope = storage.AllowedCopyScope(v)
	}

	storageAccountIdentity, err := expandAzureRmStorageAccountIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return err
	}
	parameters.Identity = storageAccountIdentity

	if v, ok := d.GetOk("azure_files_authentication"); ok {
		expandAADFilesAuthentication, err := expandArmStorageAccountAzureFilesAuthentication(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("parsing `azure_files_authentication`: %v", err)
		}
		parameters.AzureFilesIdentityBasedAuthentication = expandAADFilesAuthentication
	}

	if _, ok := d.GetOk("custom_domain"); ok {
		parameters.CustomDomain = expandStorageAccountCustomDomain(d.Get("custom_domain").([]interface{}))
	}

	if v, ok := d.GetOk("immutability_policy"); ok {
		parameters.ImmutableStorageWithVersioning = expandStorageAccountImmutabilityPolicy(v.([]interface{}))
	}

	// BlobStorage does not support ZRS
	if accountKind == storage.KindBlobStorage {
		if string(parameters.Sku.Name) == string(storage.SkuNameStandardZRS) {
			return fmt.Errorf("a `account_replication_type` of `ZRS` isn't supported for Blob Storage accounts")
		}
	}

	accessTier, ok := d.GetOk("access_tier")
	if slices.Contains(storageKindsSupportsSkuTier, accountKind) {
		if !ok {
			// default to "Hot"
			accessTier = string(storage.AccessTierHot)
		}
		parameters.AccountPropertiesCreateParameters.AccessTier = storage.AccessTier(accessTier.(string))
	} else if ok {
		return fmt.Errorf("`access_tier` is only available for accounts of kind: %v", storageKindsSupportsSkuTier)
	}

	if isHnsEnabled && !slices.Contains(storageKindsSupportHns, accountKind) {
		return fmt.Errorf("`is_hns_enabled` can only be used with account of kinds: %v", storageKindsSupportHns)
	}

	// NFSv3 is supported for standard general-purpose v2 storage accounts and for premium block blob storage accounts.
	// (https://docs.microsoft.com/en-us/azure/storage/blobs/network-file-system-protocol-support-how-to#step-5-create-and-configure-a-storage-account)
	if nfsV3Enabled &&
		!((accountTier == storage.SkuTierPremium && accountKind == storage.KindBlockBlobStorage) ||
			(accountTier == storage.SkuTierStandard && accountKind == storage.KindStorageV2)) {
		return fmt.Errorf("`nfsv3_enabled` can only be used with account tier `Standard` and account kind `StorageV2`, or account tier `Premium` and account kind `BlockBlobStorage`")
	}
	if nfsV3Enabled && !isHnsEnabled {
		return fmt.Errorf("`nfsv3_enabled` can only be used when `is_hns_enabled` is `true`")
	}

	// AccountTier must be Premium for FileStorage
	if accountKind == storage.KindFileStorage {
		if string(parameters.Sku.Tier) == string(storage.SkuNameStandardLRS) {
			return fmt.Errorf("a `account_tier` of `Standard` is not supported for FileStorage accounts")
		}
	}

	// nolint staticcheck
	if v, ok := d.GetOkExists("large_file_share_enabled"); ok {
		parameters.LargeFileSharesState = storage.LargeFileSharesStateDisabled
		if v.(bool) {
			parameters.LargeFileSharesState = storage.LargeFileSharesStateEnabled
		}
	}

	if v, ok := d.GetOk("routing"); ok {
		parameters.RoutingPreference = expandArmStorageAccountRouting(v.([]interface{}))
	}

	// TODO 4.0
	// look into standardizing this across resources that support CMK and at the very least look at improving the UX
	// for encryption of blob, file, table and queue
	var encryption *storage.Encryption
	if v, ok := d.GetOk("customer_managed_key"); ok {
		if accountTier != storage.SkuTierPremium && accountKind != storage.KindStorageV2 {
			return fmt.Errorf("customer managed key can only be used with account kind `StorageV2` or account tier `Premium`")
		}
		if storageAccountIdentity.Type != storage.IdentityTypeUserAssigned && storageAccountIdentity.Type != storage.IdentityTypeSystemAssignedUserAssigned {
			return fmt.Errorf("customer managed key can only be used with identity type `UserAssigned` or `SystemAssigned, UserAssigned`")
		}
		encryption, err = expandStorageAccountCustomerManagedKey(ctx, keyVaultClient, id.SubscriptionId, v.([]interface{}))
		if err != nil {
			return err
		}
	}

	// By default (by leaving empty), the table and queue encryption key type is set to "Service". While users can change it to "Account" so that
	// they can further use CMK to encrypt table/queue data. Only the StorageV2 account kind supports the Account key type.
	// Also noted that the blob and file are always using the "Account" key type.
	// See: https://docs.microsoft.com/en-gb/azure/storage/common/account-encryption-key-create?tabs=portal
	queueEncryptionKeyType := d.Get("queue_encryption_key_type").(string)
	tableEncryptionKeyType := d.Get("table_encryption_key_type").(string)

	if accountKind != storage.KindStorageV2 {
		if queueEncryptionKeyType == string(storage.KeyTypeAccount) {
			return fmt.Errorf("`queue_encryption_key_type = %q` can only be used with account kind `%q`", string(storage.KeyTypeAccount), string(storage.KindStorageV2))
		}
		if tableEncryptionKeyType == string(storage.KeyTypeAccount) {
			return fmt.Errorf("`table_encryption_key_type = %q` can only be used with account kind `%q`", string(storage.KeyTypeAccount), string(storage.KindStorageV2))
		}
	}

	// if CMK is not supplied then only set storage encryption for queue and table, otherwise add it to the existing encryption block for CMK
	if encryption == nil {
		encryption = &storage.Encryption{
			KeySource: storage.KeySourceMicrosoftStorage,
			Services: &storage.EncryptionServices{
				Queue: &storage.EncryptionService{
					KeyType: storage.KeyType(queueEncryptionKeyType),
				},
				Table: &storage.EncryptionService{
					KeyType: storage.KeyType(tableEncryptionKeyType),
				},
			},
		}
	} else {
		encryption.Services.Queue = &storage.EncryptionService{
			KeyType: storage.KeyType(queueEncryptionKeyType),
		}
		encryption.Services.Table = &storage.EncryptionService{
			KeyType: storage.KeyType(tableEncryptionKeyType),
		}
	}

	infrastructureEncryption := d.Get("infrastructure_encryption_enabled").(bool)

	if infrastructureEncryption {
		if !((accountTier == storage.SkuTierPremium && (accountKind == storage.KindBlockBlobStorage) || accountKind == storage.KindFileStorage) ||
			(accountKind == storage.KindStorageV2)) {
			return fmt.Errorf("`infrastructure_encryption_enabled` can only be used with account kind `StorageV2`, or account tier `Premium` and account kind is one of `BlockBlobStorage` or `FileStorage`")
		}
		encryption.RequireInfrastructureEncryption = &infrastructureEncryption
	}

	parameters.Encryption = encryption

	// Create
	future, err := legacyClient.Create(ctx, id.ResourceGroupName, id.StorageAccountName, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, legacyClient.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// populate the cache
	account, err := storageAccountsClient.GetProperties(ctx, id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if account.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if err := storageClient.AddToCache(id, *account.Model); err != nil {
		return fmt.Errorf("populating cache for %s: %+v", id, err)
	}

	dataPlaneAccount, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if dataPlaneAccount == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	supportLevel := resolveStorageAccountServiceSupportLevel(accountKind, accountTier, replicationType)
	if err := resourceStorageAccountWaitForDataPlaneToBecomeAvailable(ctx, storageClient, dataPlaneAccount, supportLevel); err != nil {
		return fmt.Errorf("waiting for the Data Plane for %s to become available: %+v", id, err)
	}

	if val, ok := d.GetOk("blob_properties"); ok {
		if !supportLevel.supportBlob {
			return fmt.Errorf("`blob_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}
		blobClient := meta.(*clients.Client).Storage.BlobServicesClient

		blobProperties, err := expandBlobProperties(accountKind, val.([]interface{}))
		if err != nil {
			return err
		}

		// See: https://learn.microsoft.com/en-us/azure/storage/blobs/versioning-overview#:~:text=Storage%20accounts%20with%20a%20hierarchical%20namespace%20enabled%20for%20use%20with%20Azure%20Data%20Lake%20Storage%20Gen2%20are%20not%20currently%20supported.
		if blobProperties.IsVersioningEnabled != nil && *blobProperties.IsVersioningEnabled && isHnsEnabled {
			return fmt.Errorf("`versioning_enabled` can't be true when `is_hns_enabled` is true")
		}

		if blobProperties.IsVersioningEnabled != nil && !*blobProperties.IsVersioningEnabled {
			if blobProperties.RestorePolicy != nil && blobProperties.RestorePolicy.Enabled != nil && *blobProperties.RestorePolicy.Enabled {
				// Otherwise, API returns: "Conflicting feature 'restorePolicy' is enabled. Please disable it and retry."
				return fmt.Errorf("`blob_properties.restore_policy` can't be set when `versioning_enabled` is false")
			}
			if account.Model.Properties != nil &&
				account.Model.Properties.ImmutableStorageWithVersioning != nil &&
				account.Model.Properties.ImmutableStorageWithVersioning.ImmutabilityPolicy != nil &&
				account.Model.Properties.ImmutableStorageWithVersioning.Enabled != nil &&
				*account.Model.Properties.ImmutableStorageWithVersioning.Enabled {
				// Otherwise, API returns: "Conflicting feature 'Account level WORM' is enabled. Please disable it and retry."
				// See: https://learn.microsoft.com/en-us/azure/storage/blobs/immutable-policy-configure-version-scope?tabs=azure-portal#prerequisites
				return fmt.Errorf("`immutability_policy` can't be set when `versioning_enabled` is false")
			}
		}

		// TODO: This is a temporary limitation on Storage service. Remove this check once the API supports this scenario.
		// See https://github.com/hashicorp/terraform-provider-azurerm/pull/25450#discussion_r1542471667 for the context.
		if dnsEndpointType == string(storage.DNSEndpointTypeAzureDNSZone) {
			if blobProperties.RestorePolicy != nil && blobProperties.RestorePolicy.Enabled != nil && *blobProperties.RestorePolicy.Enabled {
				// Otherwise, API returns: "Required feature Global Dns is disabled"
				// This is confirmed with the SRP team, where they said:
				// > restorePolicy feature is incompatible with partitioned DNS
				return fmt.Errorf("`blob_properties.restore_policy` can't be set when `dns_endpoint_type` is set to `%s`", storage.DNSEndpointTypeAzureDNSZone)
			}
		}

		if _, err = blobClient.SetServiceProperties(ctx, id.ResourceGroupName, id.StorageAccountName, *blobProperties); err != nil {
			return fmt.Errorf("updating `blob_properties`: %+v", err)
		}
	}

	if val, ok := d.GetOk("queue_properties"); ok {
		if !supportLevel.supportQueue {
			return fmt.Errorf("`queue_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}
		accountDetails, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if accountDetails == nil {
			return fmt.Errorf("unable to locate %q", id)
		}

		queueClient, err := storageClient.QueuesDataPlaneClient(ctx, *accountDetails, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Queues Client: %s", err)
		}

		queueProperties, err := expandQueueProperties(val.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `queue_properties`: %+v", err)
		}

		if err = queueClient.UpdateServiceProperties(ctx, queueProperties); err != nil {
			return fmt.Errorf("updating Queue Properties: %+v", err)
		}
	}

	if val, ok := d.GetOk("share_properties"); ok {
		if !supportLevel.supportShare {
			return fmt.Errorf("`share_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}
		fileServiceClient := meta.(*clients.Client).Storage.FileServicesClient

		shareProperties := expandShareProperties(val.([]interface{}))

		// The API complains if any multichannel info is sent on non premium fileshares. Even if multichannel is set to false
		if accountTier != storage.SkuTierPremium && shareProperties.FileServicePropertiesProperties != nil && shareProperties.FileServicePropertiesProperties.ProtocolSettings != nil {
			// Error if the user has tried to enable multichannel on a standard tier storage account
			smb := shareProperties.FileServicePropertiesProperties.ProtocolSettings.Smb
			if smb != nil && smb.Multichannel != nil {
				if smb.Multichannel.Enabled != nil {
					if *shareProperties.FileServicePropertiesProperties.ProtocolSettings.Smb.Multichannel.Enabled {
						return fmt.Errorf("`multichannel_enabled` isn't supported for Standard tier Storage accounts")
					}
				}

				shareProperties.FileServicePropertiesProperties.ProtocolSettings.Smb.Multichannel = nil
			}
		}

		if _, err = fileServiceClient.SetServiceProperties(ctx, id.ResourceGroupName, id.StorageAccountName, shareProperties); err != nil {
			return fmt.Errorf("updating `share_properties`: %+v", err)
		}
	}

	if val, ok := d.GetOk("static_website"); ok {
		if !supportLevel.supportStaticWebsite {
			return fmt.Errorf("`static_website` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if account == nil {
			return fmt.Errorf("unable to locate %s", id)
		}

		accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Accounts Data Plane Client: %s", err)
		}

		staticWebsiteProps := expandStaticWebsiteProperties(val.([]interface{}))

		if _, err = accountsClient.SetServiceProperties(ctx, id.StorageAccountName, staticWebsiteProps); err != nil {
			return fmt.Errorf("updating `static_website`: %+v", err)
		}
	}

	return resourceStorageAccountRead(d, meta)
}

func resourceStorageAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	tenantId := meta.(*clients.Client).Account.TenantId
	client := meta.(*clients.Client).Storage.AccountsClient
	keyVaultClient := meta.(*clients.Client).KeyVault
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	accountTier := storage.SkuTier(d.Get("account_tier").(string))
	replicationType := d.Get("account_replication_type").(string)
	storageType := fmt.Sprintf("%s_%s", accountTier, replicationType)
	accountKind := storage.Kind(d.Get("account_kind").(string))

	if accountKind == storage.KindBlobStorage || accountKind == storage.KindStorage {
		if storageType == string(storage.SkuNameStandardZRS) {
			return fmt.Errorf("an `account_replication_type` of `ZRS` isn't supported for Blob Storage accounts")
		}
	}

	existing, err := client.GetProperties(ctx, id.ResourceGroupName, id.StorageAccountName, "")
	if err != nil {
		return fmt.Errorf("reading for %s: %+v", id, err)
	}

	if existing.AccountProperties == nil {
		return fmt.Errorf("unexpected nil AccountProperties of %s", id)
	}

	params := storage.AccountCreateParameters{
		Sku:              existing.Sku,
		Kind:             existing.Kind,
		Location:         existing.Location,
		ExtendedLocation: existing.ExtendedLocation,
		Tags:             existing.Tags,
		Identity:         existing.Identity,
		AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{
			AllowedCopyScope:                      existing.AccountProperties.AllowedCopyScope,
			PublicNetworkAccess:                   existing.AccountProperties.PublicNetworkAccess,
			SasPolicy:                             existing.AccountProperties.SasPolicy,
			KeyPolicy:                             existing.AccountProperties.KeyPolicy,
			CustomDomain:                          existing.AccountProperties.CustomDomain,
			Encryption:                            existing.AccountProperties.Encryption,
			NetworkRuleSet:                        existing.AccountProperties.NetworkRuleSet,
			AccessTier:                            existing.AccountProperties.AccessTier,
			AzureFilesIdentityBasedAuthentication: existing.AccountProperties.AzureFilesIdentityBasedAuthentication,
			EnableHTTPSTrafficOnly:                existing.AccountProperties.EnableHTTPSTrafficOnly,
			IsSftpEnabled:                         existing.AccountProperties.IsSftpEnabled,
			IsLocalUserEnabled:                    existing.AccountProperties.IsLocalUserEnabled,
			IsHnsEnabled:                          existing.AccountProperties.IsHnsEnabled,
			LargeFileSharesState:                  existing.AccountProperties.LargeFileSharesState,
			RoutingPreference:                     existing.AccountProperties.RoutingPreference,
			AllowBlobPublicAccess:                 existing.AccountProperties.AllowBlobPublicAccess,
			MinimumTLSVersion:                     existing.AccountProperties.MinimumTLSVersion,
			AllowSharedKeyAccess:                  existing.AccountProperties.AllowSharedKeyAccess,
			EnableNfsV3:                           existing.AccountProperties.EnableNfsV3,
			AllowCrossTenantReplication:           existing.AccountProperties.AllowCrossTenantReplication,
			DefaultToOAuthAuthentication:          existing.AccountProperties.DefaultToOAuthAuthentication,
			ImmutableStorageWithVersioning:        existing.AccountProperties.ImmutableStorageWithVersioning,
			DNSEndpointType:                       existing.AccountProperties.DNSEndpointType,
		},
	}

	props := params.AccountPropertiesCreateParameters

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

	if d.HasChange("cross_tenant_replication_enabled") {
		props.AllowCrossTenantReplication = pointer.To(d.Get("cross_tenant_replication_enabled").(bool))
	}

	if d.HasChange("account_replication_type") {
		// storageType is derived from "account_replication_type" and "account_tier" (force-new)
		params.Sku = &storage.Sku{
			Name: storage.SkuName(storageType),
		}
	}

	if d.HasChange("account_kind") {
		params.Kind = accountKind
	}

	if d.HasChange("access_tier") {
		props.AccessTier = storage.AccessTier(d.Get("access_tier").(string))
	}

	if d.HasChange("tags") {
		params.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("custom_domain") {
		props.CustomDomain = expandStorageAccountCustomDomain(d.Get("custom_domain").([]interface{}))
	}

	if d.HasChange("identity") {
		params.Identity, err = expandAzureRmStorageAccountIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return err
		}
	}

	if d.HasChange("customer_managed_key") {
		props.Encryption, err = expandStorageAccountCustomerManagedKey(ctx, keyVaultClient, id.SubscriptionId, d.Get("customer_managed_key").([]interface{}))
		if err != nil {
			return err
		}
	}

	if d.HasChange("local_user_enabled") {
		props.IsLocalUserEnabled = pointer.To(d.Get("local_user_enabled").(bool))
	}

	if d.HasChange("sftp_enabled") {
		props.IsSftpEnabled = pointer.To(d.Get("sftp_enabled").(bool))
	}

	if d.HasChange("enable_https_traffic_only") {
		props.EnableHTTPSTrafficOnly = pointer.To(d.Get("enable_https_traffic_only").(bool))
	}

	if d.HasChange("min_tls_version") {
		props.MinimumTLSVersion = storage.MinimumTLSVersion(d.Get("min_tls_version").(string))
	}

	if d.HasChange("allow_nested_items_to_be_public") {
		props.AllowBlobPublicAccess = pointer.To(d.Get("allow_nested_items_to_be_public").(bool))
	}

	if d.HasChange("public_network_access_enabled") {
		publicNetworkAccess := storage.PublicNetworkAccessDisabled
		if d.Get("public_network_access_enabled").(bool) {
			publicNetworkAccess = storage.PublicNetworkAccessEnabled
		}
		props.PublicNetworkAccess = publicNetworkAccess
	}

	if d.HasChange("network_rules") {
		props.NetworkRuleSet = expandStorageAccountNetworkRules(d, tenantId)
	}

	if d.HasChange("large_file_share_enabled") {
		isEnabled := storage.LargeFileSharesStateDisabled
		if v := d.Get("large_file_share_enabled").(bool); v {
			isEnabled = storage.LargeFileSharesStateEnabled
		}
		props.LargeFileSharesState = isEnabled
	}

	if d.HasChange("routing") {
		props.RoutingPreference = expandArmStorageAccountRouting(d.Get("routing").([]interface{}))
	}

	if d.HasChange("sas_policy") {
		// TODO: Currently, due to Track1 SDK has no way to represent a `null` value in the payload - instead it will be omitted, `sas_policy` can not be disabled once enabled.
		props.SasPolicy = expandStorageAccountSASPolicy(d.Get("sas_policy").([]interface{}))
	}

	if d.HasChange("allowed_copy_scope") {
		// TODO: Currently, due to Track1 SDK has no way to represent a `null` value in the payload - instead it will be omitted, `allowed_copy_scope` can not be disabled once enabled.
		props.AllowedCopyScope = storage.AllowedCopyScope(d.Get("allowed_copy_scope").(string))
	}

	// Update (via PUT) for the above changes
	future, err := client.Create(ctx, id.ResourceGroupName, id.StorageAccountName, params)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	// azure_files_authentication must be the last to be updated, cause it'll occupy the storage account for several minutes after receiving the response 200 OK. Issue: https://github.com/Azure/azure-rest-api-specs/issues/11272
	if d.HasChange("azure_files_authentication") {
		// due to service issue: https://github.com/Azure/azure-rest-api-specs/issues/12473, we need to update to None before changing its DirectoryServiceOptions
		old, new := d.GetChange("azure_files_authentication.0.directory_type")
		if old != new && new != string(storage.DirectoryServiceOptionsNone) {
			log.Print("[DEBUG] Disabling AzureFilesIdentityBasedAuthentication prior to changing DirectoryServiceOptions")
			dsNone := storage.AccountUpdateParameters{
				AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
					AzureFilesIdentityBasedAuthentication: &storage.AzureFilesIdentityBasedAuthentication{
						DirectoryServiceOptions: storage.DirectoryServiceOptionsNone,
					},
				},
			}
			if _, err := client.Update(ctx, id.ResourceGroupName, id.StorageAccountName, dsNone); err != nil {
				return fmt.Errorf("updating `azure_files_authentication` for %s: %+v", *id, err)
			}
		}

		expandAADFilesAuthentication, err := expandArmStorageAccountAzureFilesAuthentication(d.Get("azure_files_authentication").([]interface{}))
		if err != nil {
			return err
		}
		opts := storage.AccountUpdateParameters{
			AccountPropertiesUpdateParameters: &storage.AccountPropertiesUpdateParameters{
				AzureFilesIdentityBasedAuthentication: expandAADFilesAuthentication,
			},
		}

		if _, err := client.Update(ctx, id.ResourceGroupName, id.StorageAccountName, opts); err != nil {
			return fmt.Errorf("updating `azure_files_authentication` for %s: %+v", *id, err)
		}
	}

	// Followings are updates to the sub-services
	supportLevel := resolveStorageAccountServiceSupportLevel(accountKind, accountTier, replicationType)

	if d.HasChange("blob_properties") {
		if !supportLevel.supportBlob {
			return fmt.Errorf("`blob_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		blobClient := meta.(*clients.Client).Storage.BlobServicesClient
		blobProperties, err := expandBlobProperties(accountKind, d.Get("blob_properties").([]interface{}))
		if err != nil {
			return err
		}

		if blobProperties.IsVersioningEnabled != nil && *blobProperties.IsVersioningEnabled && d.Get("is_hns_enabled").(bool) {
			return fmt.Errorf("`versioning_enabled` can't be true when `is_hns_enabled` is true")
		}

		// Disable restore_policy first. Disabling restore_policy and while setting delete_retention_policy.allow_permanent_delete to true cause error.
		// Issue : https://github.com/Azure/azure-rest-api-specs/issues/11237
		if v := d.Get("blob_properties.0.restore_policy"); d.HasChange("blob_properties.0.restore_policy") && len(v.([]interface{})) == 0 {
			log.Print("[DEBUG] Disabling RestorePolicy prior to changing DeleteRetentionPolicy")
			props := storage.BlobServiceProperties{
				BlobServicePropertiesProperties: &storage.BlobServicePropertiesProperties{
					RestorePolicy: expandBlobPropertiesRestorePolicy(v.([]interface{})),
				},
			}
			if _, err := blobClient.SetServiceProperties(ctx, id.ResourceGroupName, id.StorageAccountName, props); err != nil {
				return fmt.Errorf("updating Azure Storage Account blob restore policy %q: %+v", id.StorageAccountName, err)
			}
		}

		if d.Get("dns_endpoint_type").(string) == string(storage.DNSEndpointTypeAzureDNSZone) {
			if blobProperties.RestorePolicy != nil && blobProperties.RestorePolicy.Enabled != nil && *blobProperties.RestorePolicy.Enabled {
				// Otherwise, API returns: "Required feature Global Dns is disabled"
				// This is confirmed with the SRP team, where they said:
				// > restorePolicy feature is incompatible with partitioned DNS
				return fmt.Errorf("`blob_properties.restore_policy` can't be set when `dns_endpoint_type` is set to `%s`", storage.DNSEndpointTypeAzureDNSZone)
			}
		}

		if _, err = blobClient.SetServiceProperties(ctx, id.ResourceGroupName, id.StorageAccountName, *blobProperties); err != nil {
			return fmt.Errorf("updating `blob_properties` for %s: %+v", *id, err)
		}
	}

	if d.HasChange("queue_properties") {
		if !supportLevel.supportQueue {
			return fmt.Errorf("`queue_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		storageClient := meta.(*clients.Client).Storage
		account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}
		if account == nil {
			return fmt.Errorf("unable to locate %s", *id)
		}

		queueClient, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Queues Client: %s", err)
		}

		queueProperties, err := expandQueueProperties(d.Get("queue_properties").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `queue_properties` for %s: %+v", *id, err)
		}

		if err = queueClient.UpdateServiceProperties(ctx, queueProperties); err != nil {
			return fmt.Errorf("updating Queue Properties for %s: %+v", *id, err)
		}
	}

	if d.HasChange("share_properties") {
		if !supportLevel.supportShare {
			return fmt.Errorf("`share_properties` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		fileServiceClient := meta.(*clients.Client).Storage.FileServicesClient

		shareProperties := expandShareProperties(d.Get("share_properties").([]interface{}))
		// The API complains if any multichannel info is sent on non premium fileshares. Even if multichannel is set to false
		if accountTier != storage.SkuTierPremium {
			// Error if the user has tried to enable multichannel on a standard tier storage account
			if shareProperties.FileServicePropertiesProperties.ProtocolSettings.Smb.Multichannel != nil && shareProperties.FileServicePropertiesProperties.ProtocolSettings.Smb.Multichannel.Enabled != nil {
				if *shareProperties.FileServicePropertiesProperties.ProtocolSettings.Smb.Multichannel.Enabled {
					return fmt.Errorf("`multichannel_enabled` isn't supported for Standard tier Storage accounts")
				}
			}

			shareProperties.FileServicePropertiesProperties.ProtocolSettings.Smb.Multichannel = nil
		}

		if _, err = fileServiceClient.SetServiceProperties(ctx, id.ResourceGroupName, id.StorageAccountName, shareProperties); err != nil {
			return fmt.Errorf("updating File Share Properties for %s: %+v", *id, err)
		}
	}

	if d.HasChange("static_website") {
		if !supportLevel.supportStaticWebsite {
			return fmt.Errorf("`static_website` aren't supported for account kind %q in sku tier %q", accountKind, accountTier)
		}

		storageClient := meta.(*clients.Client).Storage

		account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}
		if account == nil {
			return fmt.Errorf("unable to locate %s", *id)
		}

		accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Data Plane client for %s: %+v", *id, err)
		}

		staticWebsiteProps := expandStaticWebsiteProperties(d.Get("static_website").([]interface{}))

		if _, err = accountsClient.SetServiceProperties(ctx, id.StorageAccountName, staticWebsiteProps); err != nil {
			return fmt.Errorf("updating `static_website` for %s: %+v", *id, err)
		}
	}

	return resourceStorageAccountRead(d, meta)
}

func resourceStorageAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.AccountsClient
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

	resp, err := client.GetProperties(ctx, id.ResourceGroupName, id.StorageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
	}

	// handle the user not having permissions to list the keys
	d.Set("primary_connection_string", "")
	d.Set("secondary_connection_string", "")
	d.Set("primary_blob_connection_string", "")
	d.Set("secondary_blob_connection_string", "")
	d.Set("primary_access_key", "")
	d.Set("secondary_access_key", "")

	keys, err := client.ListKeys(ctx, id.ResourceGroupName, id.StorageAccountName, storage.ListKeyExpandKerb)
	if err != nil {
		// the API returns a 200 with an inner error of a 409..
		var hasWriteLock bool
		var doesntHavePermissions bool
		if e, ok := err.(azautorest.DetailedError); ok {
			if status, ok := e.StatusCode.(int); ok {
				hasWriteLock = status == http.StatusConflict
				doesntHavePermissions = (status == http.StatusUnauthorized || status == http.StatusForbidden)
			}
		}

		if !hasWriteLock && !doesntHavePermissions {
			return fmt.Errorf("listing Keys for %s: %s", *id, err)
		}
	}

	d.Set("name", id.StorageAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("edge_zone", flattenEdgeZone(resp.ExtendedLocation))
	d.Set("account_kind", resp.Kind)

	if sku := resp.Sku; sku != nil {
		d.Set("account_tier", sku.Tier)
		d.Set("account_replication_type", strings.Split(fmt.Sprintf("%v", sku.Name), "_")[1])
	}

	if props := resp.AccountProperties; props != nil {
		d.Set("access_tier", props.AccessTier)
		if err := d.Set("azure_files_authentication", flattenArmStorageAccountAzureFilesAuthentication(props.AzureFilesIdentityBasedAuthentication)); err != nil {
			return fmt.Errorf("setting `azure_files_authentication`: %+v", err)
		}
		if err := d.Set("routing", flattenArmStorageAccountRouting(props.RoutingPreference)); err != nil {
			return fmt.Errorf("setting `routing`: %+v", err)
		}
		d.Set("enable_https_traffic_only", props.EnableHTTPSTrafficOnly)
		d.Set("is_hns_enabled", props.IsHnsEnabled)
		d.Set("nfsv3_enabled", props.EnableNfsV3)

		publicNetworkAccessEnabled := true
		if props.PublicNetworkAccess == storage.PublicNetworkAccessDisabled {
			publicNetworkAccessEnabled = false
		}
		d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

		// DNSEndpointType is null when unconfigured - therefore default this to Standard
		dnsEndpointType := storage.DNSEndpointTypeStandard
		if props.DNSEndpointType != "" {
			// TODO: when this is ported over to `hashicorp/go-azure-sdk` this should be able to become != nil
			dnsEndpointType = props.DNSEndpointType
		}
		d.Set("dns_endpoint_type", dnsEndpointType)

		if crossTenantReplication := props.AllowCrossTenantReplication; crossTenantReplication != nil {
			d.Set("cross_tenant_replication_enabled", crossTenantReplication)
		}

		// There is a certain edge case that could result in the Azure API returning a null value for AllowBLobPublicAccess.
		// Since the field is a pointer, this gets marshalled to a nil value instead of a boolean.

		allowBlobPublicAccess := true
		if props.AllowBlobPublicAccess != nil {
			allowBlobPublicAccess = *props.AllowBlobPublicAccess
		}
		// lintignore:R001
		d.Set("allow_nested_items_to_be_public", allowBlobPublicAccess)

		// For storage account created using old API, the response of GET call will not return "min_tls_version"
		minTlsVersion := string(storage.MinimumTLSVersionTLS10)
		if props.MinimumTLSVersion != "" {
			minTlsVersion = string(props.MinimumTLSVersion)
		}
		d.Set("min_tls_version", minTlsVersion)

		if err := d.Set("custom_domain", flattenStorageAccountCustomDomain(props.CustomDomain)); err != nil {
			return fmt.Errorf("setting `custom_domain`: %+v", err)
		}

		if immutabilityPolicy := props.ImmutableStorageWithVersioning; immutabilityPolicy != nil && immutabilityPolicy.ImmutabilityPolicy != nil {
			if err := d.Set("immutability_policy", flattenStorageAccountImmutabilityPolicy(props.ImmutableStorageWithVersioning)); err != nil {
				return fmt.Errorf("setting `immutability_policy`: %+v", err)
			}
		}

		// Computed
		d.Set("primary_location", props.PrimaryLocation)
		d.Set("secondary_location", props.SecondaryLocation)

		if accessKeys := keys.Keys; accessKeys != nil {
			storageAccountKeys := *accessKeys
			if len(storageAccountKeys) > 0 {
				pcs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *storageAccountKeys[0].Value, *storageDomainSuffix)
				d.Set("primary_connection_string", pcs)
			}

			if len(storageAccountKeys) > 1 {
				scs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *storageAccountKeys[1].Value, *storageDomainSuffix)
				d.Set("secondary_connection_string", scs)
			}
		}

		if err := flattenAndSetAzureRmStorageAccountPrimaryEndpoints(d, props.PrimaryEndpoints, resp.RoutingPreference); err != nil {
			return fmt.Errorf("setting internet, microsoft primary endpoints and hosts for blob, queue, table and file: %+v", err)
		}

		if accessKeys := keys.Keys; accessKeys != nil {
			storageAccountKeys := *accessKeys
			var primaryBlobConnectStr string
			if v := props.PrimaryEndpoints; v != nil {
				primaryBlobConnectStr = getBlobConnectionString(v.Blob, resp.Name, storageAccountKeys[0].Value)
			}
			d.Set("primary_blob_connection_string", primaryBlobConnectStr)
		}

		if err := flattenAndSetAzureRmStorageAccountSecondaryEndpoints(d, props.SecondaryEndpoints, resp.RoutingPreference); err != nil {
			return fmt.Errorf("setting internet, microsoft secondary endpoints and hosts for blob, queue, table: %+v", err)
		}

		if accessKeys := keys.Keys; accessKeys != nil {
			storageAccountKeys := *accessKeys
			var secondaryBlobConnectStr string
			if v := props.SecondaryEndpoints; v != nil {
				secondaryBlobConnectStr = getBlobConnectionString(v.Blob, resp.Name, storageAccountKeys[1].Value)
			}
			d.Set("secondary_blob_connection_string", secondaryBlobConnectStr)
		}

		if err := d.Set("network_rules", flattenStorageAccountNetworkRules(props.NetworkRuleSet)); err != nil {
			return fmt.Errorf("setting `network_rules`: %+v", err)
		}

		if props.LargeFileSharesState != "" {
			d.Set("large_file_share_enabled", props.LargeFileSharesState == storage.LargeFileSharesStateEnabled)
		}

		// local_user_enabled defaults to true at service side when not specified in the API request.
		isLocalEnabled := true
		if props.IsLocalUserEnabled != nil {
			isLocalEnabled = *props.IsLocalUserEnabled
		}
		d.Set("local_user_enabled", isLocalEnabled)

		allowSharedKeyAccess := true
		if props.AllowSharedKeyAccess != nil {
			allowSharedKeyAccess = *props.AllowSharedKeyAccess
		}
		d.Set("shared_access_key_enabled", allowSharedKeyAccess)

		defaultToOAuthAuthentication := false
		if props.DefaultToOAuthAuthentication != nil {
			defaultToOAuthAuthentication = *props.DefaultToOAuthAuthentication
		}
		d.Set("default_to_oauth_authentication", defaultToOAuthAuthentication)

		// Setting the encryption key type to "Service" in PUT. The following GET will not return the queue/table in the service list of its response.
		// So defaults to setting the encryption key type to "Service" if it is absent in the GET response. Also, define the default value as "Service" in the schema.
		var (
			queueEncryptionKeyType = string(storage.KeyTypeService)
			tableEncryptionKeyType = string(storage.KeyTypeService)
		)
		if encryption := props.Encryption; encryption != nil && encryption.Services != nil {
			if encryption.Services.Queue != nil {
				queueEncryptionKeyType = string(encryption.Services.Queue.KeyType)
			}
			if encryption.Services.Table != nil {
				tableEncryptionKeyType = string(encryption.Services.Table.KeyType)
			}
		}
		d.Set("table_encryption_key_type", tableEncryptionKeyType)
		d.Set("queue_encryption_key_type", queueEncryptionKeyType)

		customerManagedKey, err := flattenStorageAccountCustomerManagedKey(id, props.Encryption, env)
		if err != nil {
			return err
		}

		if err := d.Set("customer_managed_key", customerManagedKey); err != nil {
			return fmt.Errorf("setting `customer_managed_key`: %+v", err)
		}

		infrastructureEncryption := false
		if encryption := props.Encryption; encryption != nil && encryption.RequireInfrastructureEncryption != nil {
			infrastructureEncryption = *encryption.RequireInfrastructureEncryption
		}
		d.Set("infrastructure_encryption_enabled", infrastructureEncryption)

		if err := d.Set("sas_policy", flattenStorageAccountSASPolicy(props.SasPolicy)); err != nil {
			return fmt.Errorf("setting `sas_policy`: %+v", err)
		}

		d.Set("allowed_copy_scope", props.AllowedCopyScope)
		d.Set("sftp_enabled", props.IsSftpEnabled)
	}

	if accessKeys := keys.Keys; accessKeys != nil {
		storageAccountKeys := *accessKeys
		d.Set("primary_access_key", storageAccountKeys[0].Value)
		d.Set("secondary_access_key", storageAccountKeys[1].Value)
	}

	identity, err := flattenAzureRmStorageAccountIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	storageClient := meta.(*clients.Client).Storage
	account, err := storageClient.FindAccount(ctx, id.SubscriptionId, id.StorageAccountName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if account == nil {
		return fmt.Errorf("unable to locate %q", id)
	}

	var tier storage.SkuTier
	if resp.Sku != nil {
		tier = resp.Sku.Tier
	}
	supportLevel := resolveStorageAccountServiceSupportLevel(resp.Kind, tier, d.Get("account_replication_type").(string))

	if supportLevel.supportBlob {
		blobClient := storageClient.BlobServicesClient
		blobProps, err := blobClient.GetServiceProperties(ctx, id.ResourceGroupName, id.StorageAccountName)
		if err != nil {
			return fmt.Errorf("reading blob properties for %s: %+v", *id, err)
		}
		if err := d.Set("blob_properties", flattenBlobProperties(blobProps)); err != nil {
			return fmt.Errorf("setting `blob_properties` for %s: %+v", *id, err)
		}
	}

	if supportLevel.supportQueue {
		queueClient, err := storageClient.QueuesDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Queues Client: %s", err)
		}

		queueProps, err := queueClient.GetServiceProperties(ctx)
		if err != nil {
			return fmt.Errorf("retrieving queue properties for %s: %+v", *id, err)
		}

		if err := d.Set("queue_properties", flattenQueueProperties(queueProps)); err != nil {
			return fmt.Errorf("setting `queue_properties`: %+v", err)
		}
	}

	if supportLevel.supportShare {
		fileServiceClient := storageClient.ResourceManager.FileService

		shareProps, err := fileServiceClient.GetServiceProperties(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving share properties for %s: %+v", *id, err)
		}

		if err := d.Set("share_properties", flattenShareProperties(shareProps)); err != nil {
			return fmt.Errorf("setting `share_properties` for %s: %+v", *id, err)
		}
	}

	if supportLevel.supportStaticWebsite {
		accountsClient, err := storageClient.AccountsDataPlaneClient(ctx, *account, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
		if err != nil {
			return fmt.Errorf("building Accounts Data Plane Client: %s", err)
		}

		staticWebsiteProps, err := accountsClient.GetServiceProperties(ctx, id.StorageAccountName)
		if err != nil {
			return fmt.Errorf("retrieving static website properties for %s: %+v", *id, err)
		}
		staticWebsite := flattenStaticWebsiteProperties(staticWebsiteProps)
		if err := d.Set("static_website", staticWebsite); err != nil {
			return fmt.Errorf("setting `static_website`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceStorageAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage
	client := storageClient.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseStorageAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.StorageAccountName, storageAccountResourceName)
	defer locks.UnlockByName(id.StorageAccountName, storageAccountResourceName)

	read, err := client.GetProperties(ctx, id.ResourceGroupName, id.StorageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	// the networking api's only allow a single change to be made to a network layout at once, so let's lock to handle that
	virtualNetworkNames := make([]string, 0)
	if props := read.AccountProperties; props != nil {
		if rules := props.NetworkRuleSet; rules != nil {
			if vnr := rules.VirtualNetworkRules; vnr != nil {
				for _, v := range *vnr {
					if v.VirtualNetworkResourceID == nil {
						continue
					}

					id, err2 := commonids.ParseSubnetID(*v.VirtualNetworkResourceID)
					if err2 != nil {
						return err2
					}

					networkName := id.VirtualNetworkName
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

	resp, err := client.Delete(ctx, id.ResourceGroupName, id.StorageAccountName)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("issuing delete request for %s: %+v", *id, err)
		}
	}

	// remove this from the cache
	storageClient.RemoveAccountFromCache(*id)

	return nil
}

func resourceStorageAccountWaitForDataPlaneToBecomeAvailable(ctx context.Context, client *client.Client, account *client.AccountDetails, supportLevel storageAccountServiceSupportLevel) error {
	initialDelayDuration := 10 * time.Second

	if supportLevel.supportBlob {
		log.Printf("[DEBUG] waiting for the Blob Service to become available")
		pollerType, err := custompollers.NewDataPlaneBlobContainersAvailabilityPoller(ctx, client, account)
		if err != nil {
			return fmt.Errorf("building Blob Service Poller: %+v", err)
		}
		poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for the Blob Service to become available: %+v", err)
		}
	}

	if supportLevel.supportQueue {
		log.Printf("[DEBUG] waiting for the Queues Service to become available")
		pollerType, err := custompollers.NewDataPlaneQueuesAvailabilityPoller(ctx, client, account)
		if err != nil {
			return fmt.Errorf("building Queues Poller: %+v", err)
		}
		poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for the Queues Service to become available: %+v", err)
		}
	}

	if supportLevel.supportShare {
		log.Printf("[DEBUG] waiting for the File Service to become available")
		pollerType, err := custompollers.NewDataPlaneFileShareAvailabilityPoller(client, account)
		if err != nil {
			return fmt.Errorf("building File Share Poller: %+v", err)
		}
		poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for the File Service to become available: %+v", err)
		}
	}

	if supportLevel.supportStaticWebsite {
		log.Printf("[DEBUG] waiting for the Static Website to become available")
		pollerType, err := custompollers.NewDataPlaneStaticWebsiteAvailabilityPoller(ctx, client, account)
		if err != nil {
			return fmt.Errorf("building Static Website Poller: %+v", err)
		}
		poller := pollers.NewPoller(pollerType, initialDelayDuration, pollers.DefaultNumberOfDroppedConnectionsToAllow)
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for the Static Website to become available: %+v", err)
		}
	}

	return nil
}

func expandStorageAccountCustomDomain(input []interface{}) *storage.CustomDomain {
	if len(input) == 0 {
		return &storage.CustomDomain{
			Name: utils.String(""),
		}
	}

	domain := input[0].(map[string]interface{})
	name := domain["name"].(string)
	useSubDomain := domain["use_subdomain"].(bool)
	return &storage.CustomDomain{
		Name:             utils.String(name),
		UseSubDomainName: utils.Bool(useSubDomain),
	}
}

func flattenStorageAccountCustomDomain(input *storage.CustomDomain) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	domain := make(map[string]interface{})

	if v := input.Name; v != nil {
		domain["name"] = *v
	}

	// use_subdomain isn't returned
	return []interface{}{domain}
}

func expandStorageAccountCustomerManagedKey(ctx context.Context, keyVaultClient *keyVaultClient.Client, subscriptionId string, input []interface{}) (*storage.Encryption, error) {
	if len(input) == 0 {
		return &storage.Encryption{}, nil
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
			return nil, fmt.Errorf("unexpected nil Key Vault ID retrieved at URL %s", keyId.KeyVaultBaseUrl)
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

		keyName = utils.String(keyId.Name)
		keyVersion = utils.String(keyId.Version)
		keyVaultURI = utils.String(keyId.KeyVaultBaseUrl)
	} else if managedHSMKeyId, ok := v["managed_hsm_key_id"]; ok && managedHSMKeyId != "" {
		if keyId, err := managedHsmParse.ManagedHSMDataPlaneVersionedKeyID(managedHSMKeyId.(string), nil); err == nil {
			keyName = utils.String(keyId.KeyName)
			keyVersion = utils.String(keyId.KeyVersion)
			keyVaultURI = utils.String(keyId.BaseUri())
		} else if keyId, err := managedHsmParse.ManagedHSMDataPlaneVersionlessKeyID(managedHSMKeyId.(string), nil); err == nil {
			keyName = utils.String(keyId.KeyName)
			keyVersion = utils.String("")
			keyVaultURI = utils.String(keyId.BaseUri())
		} else {
			return nil, fmt.Errorf("Failed to parse '%s' as HSM key ID", managedHSMKeyId.(string))
		}
	}

	encryption := &storage.Encryption{
		Services: &storage.EncryptionServices{
			Blob: &storage.EncryptionService{
				Enabled: utils.Bool(true),
				KeyType: storage.KeyTypeAccount,
			},
			File: &storage.EncryptionService{
				Enabled: utils.Bool(true),
				KeyType: storage.KeyTypeAccount,
			},
		},
		EncryptionIdentity: &storage.EncryptionIdentity{
			EncryptionUserAssignedIdentity: utils.String(v["user_assigned_identity_id"].(string)),
		},
		KeySource: storage.KeySourceMicrosoftKeyvault,
		KeyVaultProperties: &storage.KeyVaultProperties{
			KeyName:     keyName,
			KeyVersion:  keyVersion,
			KeyVaultURI: keyVaultURI,
		},
	}

	return encryption, nil
}

func expandStorageAccountImmutabilityPolicy(input []interface{}) *storage.ImmutableStorageAccount {
	if len(input) == 0 {
		return &storage.ImmutableStorageAccount{}
	}

	v := input[0].(map[string]interface{})

	immutableStorageAccount := storage.ImmutableStorageAccount{
		Enabled: utils.Bool(true),
		ImmutabilityPolicy: &storage.AccountImmutabilityPolicyProperties{
			AllowProtectedAppendWrites:            utils.Bool(v["allow_protected_append_writes"].(bool)),
			State:                                 storage.AccountImmutabilityPolicyState(v["state"].(string)),
			ImmutabilityPeriodSinceCreationInDays: utils.Int32(int32(v["period_since_creation_in_days"].(int))),
		},
	}

	return &immutableStorageAccount
}

func flattenStorageAccountImmutabilityPolicy(policy *storage.ImmutableStorageAccount) []interface{} {
	if policy == nil || policy.ImmutabilityPolicy == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"period_since_creation_in_days": policy.ImmutabilityPolicy.ImmutabilityPeriodSinceCreationInDays,
			"state":                         policy.ImmutabilityPolicy.State,
			"allow_protected_append_writes": policy.ImmutabilityPolicy.AllowProtectedAppendWrites,
		},
	}
}

func flattenStorageAccountCustomerManagedKey(storageAccountId *commonids.StorageAccountId, input *storage.Encryption, env environments.Environment) ([]interface{}, error) {
	if input == nil || input.KeySource == storage.KeySourceMicrosoftStorage {
		return make([]interface{}, 0), nil
	}

	userAssignedIdentityId := ""
	keyName := ""
	keyVaultURI := ""
	keyVersion := ""

	if props := input.EncryptionIdentity; props != nil {
		if props.EncryptionUserAssignedIdentity != nil {
			userAssignedIdentityId = *props.EncryptionUserAssignedIdentity
		}
	}

	if props := input.KeyVaultProperties; props != nil {
		if props.KeyName != nil {
			keyName = *props.KeyName
		}
		if props.KeyVaultURI != nil {
			keyVaultURI = *props.KeyVaultURI
		}
		if props.KeyVersion != nil {
			keyVersion = *props.KeyVersion
		}
	}

	if keyVaultURI == "" {
		return nil, fmt.Errorf("retrieving %s: `properties.encryption.keyVaultProperties.keyVaultURI` was nil", *storageAccountId)
	}

	ret := map[string]interface{}{
		"user_assigned_identity_id": userAssignedIdentityId,
	}

	isHSMURI, err, instanceName, domainSuffix := managedHsmHelpers.IsManagedHSMURI(env, keyVaultURI)
	if err != nil {
		return nil, err
	}

	switch {
	case isHSMURI && keyVersion == "":
		{
			keyId := managedHsmParse.NewManagedHSMDataPlaneVersionlessKeyID(instanceName, domainSuffix, keyName)
			ret["managed_hsm_key_id"] = keyId.ID()
		}
	case isHSMURI && keyVersion != "":
		{
			keyId := managedHsmParse.NewManagedHSMDataPlaneVersionedKeyID(instanceName, domainSuffix, keyName, keyVersion)
			ret["managed_hsm_key_id"] = keyId.ID()
		}
	case !isHSMURI:
		{
			keyId, err := keyVaultParse.NewNestedItemID(keyVaultURI, keyVaultParse.NestedItemTypeKey, keyName, keyVersion)
			if err != nil {
				return nil, err
			}
			ret["key_vault_key_id"] = keyId.ID()
		}
	}

	return []interface{}{
		ret,
	}, nil
}

func expandArmStorageAccountAzureFilesAuthentication(input []interface{}) (*storage.AzureFilesIdentityBasedAuthentication, error) {
	if len(input) == 0 {
		return &storage.AzureFilesIdentityBasedAuthentication{
			DirectoryServiceOptions: storage.DirectoryServiceOptionsNone,
		}, nil
	}

	v := input[0].(map[string]interface{})

	ad := expandArmStorageAccountActiveDirectoryProperties(v["active_directory"].([]interface{}))

	directoryOption := storage.DirectoryServiceOptions(v["directory_type"].(string))

	if directoryOption == storage.DirectoryServiceOptionsAD {
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

	return &storage.AzureFilesIdentityBasedAuthentication{
		DirectoryServiceOptions:   directoryOption,
		ActiveDirectoryProperties: ad,
	}, nil
}

func expandArmStorageAccountActiveDirectoryProperties(input []interface{}) *storage.ActiveDirectoryProperties {
	if len(input) == 0 {
		return nil
	}
	m := input[0].(map[string]interface{})

	output := &storage.ActiveDirectoryProperties{
		DomainGUID: utils.String(m["domain_guid"].(string)),
		DomainName: utils.String(m["domain_name"].(string)),
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

func expandArmStorageAccountRouting(input []interface{}) *storage.RoutingPreference {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &storage.RoutingPreference{
		RoutingChoice:             storage.RoutingChoice(v["choice"].(string)),
		PublishMicrosoftEndpoints: utils.Bool(v["publish_microsoft_endpoints"].(bool)),
		PublishInternetEndpoints:  utils.Bool(v["publish_internet_endpoints"].(bool)),
	}
}

func expandStorageAccountNetworkRules(d *pluginsdk.ResourceData, tenantId string) *storage.NetworkRuleSet {
	networkRules := d.Get("network_rules").([]interface{})
	if len(networkRules) == 0 {
		// Default access is enabled when no network rules are set.
		return &storage.NetworkRuleSet{DefaultAction: storage.DefaultActionAllow}
	}

	networkRule := networkRules[0].(map[string]interface{})
	networkRuleSet := &storage.NetworkRuleSet{
		IPRules:             expandStorageAccountIPRules(networkRule),
		VirtualNetworkRules: expandStorageAccountVirtualNetworks(networkRule),
		Bypass:              expandStorageAccountBypass(networkRule),
		ResourceAccessRules: expandStorageAccountPrivateLinkAccess(networkRule["private_link_access"].([]interface{}), tenantId),
	}

	if v := networkRule["default_action"]; v != nil {
		networkRuleSet.DefaultAction = storage.DefaultAction(v.(string))
	}

	return networkRuleSet
}

func expandStorageAccountIPRules(networkRule map[string]interface{}) *[]storage.IPRule {
	ipRulesInfo := networkRule["ip_rules"].(*pluginsdk.Set).List()
	ipRules := make([]storage.IPRule, len(ipRulesInfo))

	for i, ipRuleConfig := range ipRulesInfo {
		attrs := ipRuleConfig.(string)
		ipRule := storage.IPRule{
			IPAddressOrRange: utils.String(attrs),
			Action:           storage.ActionAllow,
		}
		ipRules[i] = ipRule
	}

	return &ipRules
}

func expandStorageAccountVirtualNetworks(networkRule map[string]interface{}) *[]storage.VirtualNetworkRule {
	virtualNetworkInfo := networkRule["virtual_network_subnet_ids"].(*pluginsdk.Set).List()
	virtualNetworks := make([]storage.VirtualNetworkRule, len(virtualNetworkInfo))

	for i, virtualNetworkConfig := range virtualNetworkInfo {
		attrs := virtualNetworkConfig.(string)
		virtualNetwork := storage.VirtualNetworkRule{
			VirtualNetworkResourceID: utils.String(attrs),
			Action:                   storage.ActionAllow,
		}
		virtualNetworks[i] = virtualNetwork
	}

	return &virtualNetworks
}

func expandStorageAccountBypass(networkRule map[string]interface{}) storage.Bypass {
	bypassInfo := networkRule["bypass"].(*pluginsdk.Set).List()

	var bypassValues []string
	for _, bypassConfig := range bypassInfo {
		bypassValues = append(bypassValues, bypassConfig.(string))
	}

	return storage.Bypass(strings.Join(bypassValues, ", "))
}

func expandStorageAccountPrivateLinkAccess(inputs []interface{}, tenantId string) *[]storage.ResourceAccessRule {
	privateLinkAccess := make([]storage.ResourceAccessRule, 0)
	if len(inputs) == 0 {
		return &privateLinkAccess
	}
	for _, input := range inputs {
		accessRule := input.(map[string]interface{})
		if v := accessRule["endpoint_tenant_id"].(string); v != "" {
			tenantId = v
		}
		privateLinkAccess = append(privateLinkAccess, storage.ResourceAccessRule{
			TenantID:   utils.String(tenantId),
			ResourceID: utils.String(accessRule["endpoint_resource_id"].(string)),
		})
	}

	return &privateLinkAccess
}

func expandBlobProperties(kind storage.Kind, input []interface{}) (*storage.BlobServiceProperties, error) {
	props := storage.BlobServiceProperties{
		BlobServicePropertiesProperties: &storage.BlobServicePropertiesProperties{
			Cors: &storage.CorsRules{
				CorsRules: &[]storage.CorsRule{},
			},
			DeleteRetentionPolicy: &storage.DeleteRetentionPolicy{
				Enabled: utils.Bool(false),
			},
		},
	}

	// `Storage` (v1) kind doesn't support:
	// - LastAccessTimeTrackingPolicy: Confirmed by SRP.
	// - ChangeFeed: See https://learn.microsoft.com/en-us/azure/storage/blobs/storage-blob-change-feed?tabs=azure-portal#enable-and-disable-the-change-feed.
	// - Versioning: See https://learn.microsoft.com/en-us/azure/storage/blobs/versioning-overview#how-blob-versioning-works
	// - Restore Policy: See https://learn.microsoft.com/en-us/azure/storage/blobs/point-in-time-restore-overview#prerequisites-for-point-in-time-restore
	if kind != storage.KindStorage {
		props.LastAccessTimeTrackingPolicy = &storage.LastAccessTimeTrackingPolicy{
			Enable: utils.Bool(false),
		}
		props.ChangeFeed = &storage.ChangeFeed{
			Enabled: utils.Bool(false),
		}
		props.IsVersioningEnabled = utils.Bool(false)
	}

	if len(input) == 0 || input[0] == nil {
		return &props, nil
	}

	v := input[0].(map[string]interface{})

	deletePolicyRaw := v["delete_retention_policy"].([]interface{})
	props.BlobServicePropertiesProperties.DeleteRetentionPolicy = expandBlobPropertiesDeleteRetentionPolicy(deletePolicyRaw)

	containerDeletePolicyRaw := v["container_delete_retention_policy"].([]interface{})
	props.BlobServicePropertiesProperties.ContainerDeleteRetentionPolicy = expandBlobPropertiesContainerDeleteRetentionPolicyWithoutPermDeleteOption(containerDeletePolicyRaw)

	corsRaw := v["cors_rule"].([]interface{})
	props.BlobServicePropertiesProperties.Cors = expandBlobPropertiesCors(corsRaw)

	props.IsVersioningEnabled = utils.Bool(v["versioning_enabled"].(bool))

	if version, ok := v["default_service_version"].(string); ok && version != "" {
		props.DefaultServiceVersion = utils.String(version)
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
	if kind != storage.KindStorage {
		props.BlobServicePropertiesProperties.LastAccessTimeTrackingPolicy = &storage.LastAccessTimeTrackingPolicy{
			Enable: utils.Bool(lastAccessTimeEnabled),
		}
		props.BlobServicePropertiesProperties.ChangeFeed = &storage.ChangeFeed{
			Enabled: utils.Bool(changeFeedEnabled),
		}
		if changeFeedRetentionInDays != 0 {
			props.BlobServicePropertiesProperties.ChangeFeed.RetentionInDays = utils.Int32(int32(changeFeedRetentionInDays))
		}
		props.BlobServicePropertiesProperties.RestorePolicy = expandBlobPropertiesRestorePolicy(restorePolicyRaw)
		props.BlobServicePropertiesProperties.IsVersioningEnabled = &versioningEnabled
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
	if p := props.BlobServicePropertiesProperties.RestorePolicy; p != nil && p.Enabled != nil && *p.Enabled {
		if props.ChangeFeed == nil || props.ChangeFeed.Enabled == nil || !*props.ChangeFeed.Enabled {
			return nil, fmt.Errorf("`change_feed_enabled` must be `true` when `restore_policy` is set")
		}
		if props.IsVersioningEnabled == nil || !*props.IsVersioningEnabled {
			return nil, fmt.Errorf("`versioning_enabled` must be `true` when `restore_policy` is set")
		}
	}

	return &props, nil
}

func expandBlobPropertiesDeleteRetentionPolicy(input []interface{}) *storage.DeleteRetentionPolicy {
	result := storage.DeleteRetentionPolicy{
		Enabled: utils.Bool(false),
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]interface{})

	return &storage.DeleteRetentionPolicy{
		Enabled:              utils.Bool(true),
		Days:                 utils.Int32(int32(policy["days"].(int))),
		AllowPermanentDelete: utils.Bool(policy["permanent_delete_enabled"].(bool)),
	}
}

func expandBlobPropertiesContainerDeleteRetentionPolicyWithoutPermDeleteOption(input []interface{}) *storage.DeleteRetentionPolicy {
	result := storage.DeleteRetentionPolicy{
		Enabled: utils.Bool(false),
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]interface{})

	return &storage.DeleteRetentionPolicy{
		Enabled: utils.Bool(true),
		Days:    utils.Int32(int32(policy["days"].(int))),
	}
}

func expandBlobPropertiesRestorePolicy(input []interface{}) *storage.RestorePolicyProperties {
	result := storage.RestorePolicyProperties{
		Enabled: utils.Bool(false),
	}
	if len(input) == 0 || input[0] == nil {
		return &result
	}

	policy := input[0].(map[string]interface{})

	return &storage.RestorePolicyProperties{
		Enabled: utils.Bool(true),
		Days:    utils.Int32(int32(policy["days"].(int))),
	}
}

func expandBlobPropertiesCors(input []interface{}) *storage.CorsRules {
	blobCorsRules := storage.CorsRules{}

	if len(input) == 0 {
		return &blobCorsRules
	}

	corsRules := make([]storage.CorsRule, 0)
	for _, attr := range input {
		corsRuleAttr := attr.(map[string]interface{})
		corsRule := storage.CorsRule{}

		allowedOrigins := *utils.ExpandStringSlice(corsRuleAttr["allowed_origins"].([]interface{}))
		allowedHeaders := *utils.ExpandStringSlice(corsRuleAttr["allowed_headers"].([]interface{}))
		allowedMethods := *utils.ExpandStringSlice(corsRuleAttr["allowed_methods"].([]interface{}))
		exposedHeaders := *utils.ExpandStringSlice(corsRuleAttr["exposed_headers"].([]interface{}))
		maxAgeInSeconds := int32(corsRuleAttr["max_age_in_seconds"].(int))

		corsRule.AllowedOrigins = &allowedOrigins
		corsRule.AllowedHeaders = &allowedHeaders
		corsRule.AllowedMethods = &allowedMethods
		corsRule.ExposedHeaders = &exposedHeaders
		corsRule.MaxAgeInSeconds = &maxAgeInSeconds

		corsRules = append(corsRules, corsRule)
	}

	blobCorsRules.CorsRules = &corsRules

	return &blobCorsRules
}

func expandShareProperties(input []interface{}) storage.FileServiceProperties {
	props := storage.FileServiceProperties{
		FileServicePropertiesProperties: &storage.FileServicePropertiesProperties{
			Cors: &storage.CorsRules{
				CorsRules: &[]storage.CorsRule{},
			},
			ShareDeleteRetentionPolicy: &storage.DeleteRetentionPolicy{
				Enabled: utils.Bool(false),
			},
		},
	}

	if len(input) == 0 || input[0] == nil {
		return props
	}

	v := input[0].(map[string]interface{})

	props.FileServicePropertiesProperties.ShareDeleteRetentionPolicy = expandBlobPropertiesContainerDeleteRetentionPolicyWithoutPermDeleteOption(v["retention_policy"].([]interface{}))

	props.FileServicePropertiesProperties.Cors = expandBlobPropertiesCors(v["cors_rule"].([]interface{}))

	props.ProtocolSettings = &storage.ProtocolSettings{
		Smb: expandSharePropertiesSMB(v["smb"].([]interface{})),
	}

	return props
}

func expandSharePropertiesSMB(input []interface{}) *storage.SmbSetting {
	if len(input) == 0 || input[0] == nil {
		return &storage.SmbSetting{
			Versions:                 utils.String(""),
			AuthenticationMethods:    utils.String(""),
			KerberosTicketEncryption: utils.String(""),
			ChannelEncryption:        utils.String(""),
		}
	}

	v := input[0].(map[string]interface{})

	return &storage.SmbSetting{
		Versions:                 utils.ExpandStringSliceWithDelimiter(v["versions"].(*pluginsdk.Set).List(), ";"),
		AuthenticationMethods:    utils.ExpandStringSliceWithDelimiter(v["authentication_types"].(*pluginsdk.Set).List(), ";"),
		KerberosTicketEncryption: utils.ExpandStringSliceWithDelimiter(v["kerberos_ticket_encryption_type"].(*pluginsdk.Set).List(), ";"),
		ChannelEncryption:        utils.ExpandStringSliceWithDelimiter(v["channel_encryption_type"].(*pluginsdk.Set).List(), ";"),
		Multichannel: &storage.Multichannel{
			Enabled: utils.Bool(v["multichannel_enabled"].(bool)),
		},
	}
}

func expandQueueProperties(input []interface{}) (queues.StorageServiceProperties, error) {
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
		return properties, nil
	}

	attrs := input[0].(map[string]interface{})

	properties.Cors = expandQueuePropertiesCors(attrs["cors_rule"].([]interface{}))
	properties.Logging = expandQueuePropertiesLogging(attrs["logging"].([]interface{}))
	properties.MinuteMetrics, err = expandQueuePropertiesMetrics(attrs["minute_metrics"].([]interface{}))
	if err != nil {
		return properties, fmt.Errorf("expanding `minute_metrics`: %+v", err)
	}
	properties.HourMetrics, err = expandQueuePropertiesMetrics(attrs["hour_metrics"].([]interface{}))
	if err != nil {
		return properties, fmt.Errorf("expanding `hour_metrics`: %+v", err)
	}

	return properties, nil
}

func expandQueuePropertiesMetrics(input []interface{}) (*queues.MetricsConfig, error) {
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
		Version: metricsAttr["version"].(string),
		Enabled: metricsAttr["enabled"].(bool),
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

func expandQueuePropertiesLogging(input []interface{}) *queues.LoggingConfig {
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
		Version: loggingAttr["version"].(string),
		Delete:  loggingAttr["delete"].(bool),
		Read:    loggingAttr["read"].(bool),
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

func expandQueuePropertiesCors(input []interface{}) *queues.Cors {
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

func expandStaticWebsiteProperties(input []interface{}) accounts.StorageServiceProperties {
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

func flattenArmStorageAccountAzureFilesAuthentication(input *storage.AzureFilesIdentityBasedAuthentication) []interface{} {
	if input == nil || input.DirectoryServiceOptions == storage.DirectoryServiceOptionsNone {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"directory_type":   input.DirectoryServiceOptions,
			"active_directory": flattenArmStorageAccountActiveDirectoryProperties(input.ActiveDirectoryProperties),
		},
	}
}

func flattenArmStorageAccountActiveDirectoryProperties(input *storage.ActiveDirectoryProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var azureStorageSid string
	if input.AzureStorageSid != nil {
		azureStorageSid = *input.AzureStorageSid
	}
	var domainGuid string
	if input.DomainGUID != nil {
		domainGuid = *input.DomainGUID
	}
	var domainName string
	if input.DomainName != nil {
		domainName = *input.DomainName
	}
	var domainSid string
	if input.DomainSid != nil {
		domainSid = *input.DomainSid
	}
	var forestName string
	if input.ForestName != nil {
		forestName = *input.ForestName
	}
	var netBiosDomainName string
	if input.NetBiosDomainName != nil {
		netBiosDomainName = *input.NetBiosDomainName
	}
	return []interface{}{
		map[string]interface{}{
			"storage_sid":         azureStorageSid,
			"domain_guid":         domainGuid,
			"domain_name":         domainName,
			"domain_sid":          domainSid,
			"forest_name":         forestName,
			"netbios_domain_name": netBiosDomainName,
		},
	}
}

func flattenArmStorageAccountRouting(input *storage.RoutingPreference) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var publishInternetEndpoints bool
	if input.PublishInternetEndpoints != nil {
		publishInternetEndpoints = *input.PublishInternetEndpoints
	}
	var publishMicrosoftEndpoints bool
	if input.PublishMicrosoftEndpoints != nil {
		publishMicrosoftEndpoints = *input.PublishMicrosoftEndpoints
	}

	return []interface{}{
		map[string]interface{}{
			"publish_internet_endpoints":  publishInternetEndpoints,
			"publish_microsoft_endpoints": publishMicrosoftEndpoints,
			"choice":                      input.RoutingChoice,
		},
	}
}

func flattenStorageAccountNetworkRules(input *storage.NetworkRuleSet) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	networkRules := make(map[string]interface{})

	networkRules["ip_rules"] = pluginsdk.NewSet(pluginsdk.HashString, flattenStorageAccountIPRules(input.IPRules))
	networkRules["virtual_network_subnet_ids"] = pluginsdk.NewSet(pluginsdk.HashString, flattenStorageAccountVirtualNetworks(input.VirtualNetworkRules))
	networkRules["bypass"] = pluginsdk.NewSet(pluginsdk.HashString, flattenStorageAccountBypass(input.Bypass))
	networkRules["default_action"] = string(input.DefaultAction)
	networkRules["private_link_access"] = flattenStorageAccountPrivateLinkAccess(input.ResourceAccessRules)

	return []interface{}{networkRules}
}

func flattenStorageAccountIPRules(input *[]storage.IPRule) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	ipRules := make([]interface{}, 0)
	for _, ipRule := range *input {
		if ipRule.IPAddressOrRange == nil {
			continue
		}

		ipRules = append(ipRules, *ipRule.IPAddressOrRange)
	}

	return ipRules
}

func flattenStorageAccountVirtualNetworks(input *[]storage.VirtualNetworkRule) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	virtualNetworks := make([]interface{}, 0)
	for _, virtualNetwork := range *input {
		if virtualNetwork.VirtualNetworkResourceID == nil {
			continue
		}

		virtualNetworks = append(virtualNetworks, *virtualNetwork.VirtualNetworkResourceID)
	}

	return virtualNetworks
}

func flattenStorageAccountPrivateLinkAccess(inputs *[]storage.ResourceAccessRule) []interface{} {
	if inputs == nil || len(*inputs) == 0 {
		return []interface{}{}
	}

	accessRules := make([]interface{}, 0)
	for _, input := range *inputs {
		var resourceId, tenantId string
		if input.ResourceID != nil {
			resourceId = *input.ResourceID
		}

		if input.TenantID != nil {
			tenantId = *input.TenantID
		}

		accessRules = append(accessRules, map[string]interface{}{
			"endpoint_resource_id": resourceId,
			"endpoint_tenant_id":   tenantId,
		})
	}

	return accessRules
}

func flattenBlobProperties(input storage.BlobServiceProperties) []interface{} {
	if input.BlobServicePropertiesProperties == nil {
		return []interface{}{}
	}

	flattenedCorsRules := make([]interface{}, 0)
	if corsRules := input.BlobServicePropertiesProperties.Cors; corsRules != nil {
		flattenedCorsRules = flattenBlobPropertiesCorsRule(corsRules)
	}

	flattenedDeletePolicy := make([]interface{}, 0)
	if deletePolicy := input.BlobServicePropertiesProperties.DeleteRetentionPolicy; deletePolicy != nil {
		flattenedDeletePolicy = flattenBlobPropertiesDeleteRetentionPolicy(deletePolicy)
	}

	flattenedRestorePolicy := make([]interface{}, 0)
	if restorePolicy := input.BlobServicePropertiesProperties.RestorePolicy; restorePolicy != nil {
		flattenedRestorePolicy = flattenBlobPropertiesRestorePolicy(restorePolicy)
	}

	flattenedContainerDeletePolicy := make([]interface{}, 0)
	if containerDeletePolicy := input.BlobServicePropertiesProperties.ContainerDeleteRetentionPolicy; containerDeletePolicy != nil {
		flattenedContainerDeletePolicy = flattenBlobPropertiesDeleteRetentionPolicyWithoutPermDeleteOption(containerDeletePolicy)
	}

	versioning, changeFeedEnabled, changeFeedRetentionInDays := false, false, 0
	if input.BlobServicePropertiesProperties.IsVersioningEnabled != nil {
		versioning = *input.BlobServicePropertiesProperties.IsVersioningEnabled
	}

	if v := input.BlobServicePropertiesProperties.ChangeFeed; v != nil {
		if v.Enabled != nil {
			changeFeedEnabled = *v.Enabled
		}
		if v.RetentionInDays != nil {
			changeFeedRetentionInDays = int(*v.RetentionInDays)
		}
	}

	var defaultServiceVersion string
	if input.BlobServicePropertiesProperties.DefaultServiceVersion != nil {
		defaultServiceVersion = *input.BlobServicePropertiesProperties.DefaultServiceVersion
	}

	var LastAccessTimeTrackingPolicy bool
	if v := input.BlobServicePropertiesProperties.LastAccessTimeTrackingPolicy; v != nil && v.Enable != nil {
		LastAccessTimeTrackingPolicy = *v.Enable
	}

	return []interface{}{
		map[string]interface{}{
			"cors_rule":                         flattenedCorsRules,
			"delete_retention_policy":           flattenedDeletePolicy,
			"restore_policy":                    flattenedRestorePolicy,
			"versioning_enabled":                versioning,
			"change_feed_enabled":               changeFeedEnabled,
			"change_feed_retention_in_days":     changeFeedRetentionInDays,
			"default_service_version":           defaultServiceVersion,
			"last_access_time_enabled":          LastAccessTimeTrackingPolicy,
			"container_delete_retention_policy": flattenedContainerDeletePolicy,
		},
	}
}

func flattenBlobPropertiesCorsRule(input *storage.CorsRules) []interface{} {
	corsRules := make([]interface{}, 0)

	if input == nil || input.CorsRules == nil {
		return corsRules
	}

	for _, corsRule := range *input.CorsRules {
		corsRules = append(corsRules, map[string]interface{}{
			"allowed_headers":    corsRule.AllowedHeaders,
			"allowed_origins":    corsRule.AllowedOrigins,
			"allowed_methods":    corsRule.AllowedMethods,
			"exposed_headers":    corsRule.ExposedHeaders,
			"max_age_in_seconds": int(*corsRule.MaxAgeInSeconds),
		})
	}

	return corsRules
}

func flattenBlobPropertiesDeleteRetentionPolicy(input *storage.DeleteRetentionPolicy) []interface{} {
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

func flattenBlobPropertiesDeleteRetentionPolicyWithoutPermDeleteOption(input *storage.DeleteRetentionPolicy) []interface{} {
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

func flattenBlobPropertiesRestorePolicy(input *storage.RestorePolicyProperties) []interface{} {
	restorePolicy := make([]interface{}, 0)

	if input == nil {
		return restorePolicy
	}

	if enabled := input.Enabled; enabled != nil && *enabled {
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

func flattenQueueProperties(input *queues.StorageServiceProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	queueProperties := make(map[string]interface{})

	if cors := input.Cors; cors != nil {
		if len(cors.CorsRule) > 0 {
			if cors.CorsRule[0].AllowedOrigins != "" {
				queueProperties["cors_rule"] = flattenQueuePropertiesCorsRule(input.Cors.CorsRule)
			}
		}
	}

	if logging := input.Logging; logging != nil {
		if logging.Version != "" {
			queueProperties["logging"] = flattenQueuePropertiesLogging(*logging)
		}
	}

	if hourMetrics := input.HourMetrics; hourMetrics != nil {
		if hourMetrics.Version != "" {
			queueProperties["hour_metrics"] = flattenQueuePropertiesMetrics(*hourMetrics)
		}
	}

	if minuteMetrics := input.MinuteMetrics; minuteMetrics != nil {
		if minuteMetrics.Version != "" {
			queueProperties["minute_metrics"] = flattenQueuePropertiesMetrics(*minuteMetrics)
		}
	}

	if len(queueProperties) == 0 {
		return []interface{}{}
	}
	return []interface{}{queueProperties}
}

func flattenQueuePropertiesMetrics(input queues.MetricsConfig) []interface{} {
	metrics := make(map[string]interface{})

	metrics["version"] = input.Version
	metrics["enabled"] = input.Enabled

	if input.IncludeAPIs != nil {
		metrics["include_apis"] = *input.IncludeAPIs
	}

	if input.RetentionPolicy.Enabled {
		metrics["retention_policy_days"] = input.RetentionPolicy.Days
	}

	return []interface{}{metrics}
}

func flattenQueuePropertiesCorsRule(input []queues.CorsRule) []interface{} {
	corsRules := make([]interface{}, 0)

	for _, corsRule := range input {
		attr := make(map[string]interface{})

		attr["allowed_origins"] = flattenCorsProperty(corsRule.AllowedOrigins)
		attr["allowed_methods"] = flattenCorsProperty(corsRule.AllowedMethods)
		attr["allowed_headers"] = flattenCorsProperty(corsRule.AllowedHeaders)
		attr["exposed_headers"] = flattenCorsProperty(corsRule.ExposedHeaders)
		attr["max_age_in_seconds"] = corsRule.MaxAgeInSeconds

		corsRules = append(corsRules, attr)
	}

	return corsRules
}

func flattenQueuePropertiesLogging(input queues.LoggingConfig) []interface{} {
	logging := make(map[string]interface{})

	logging["version"] = input.Version
	logging["delete"] = input.Delete
	logging["read"] = input.Read
	logging["write"] = input.Write

	if input.RetentionPolicy.Enabled {
		logging["retention_policy_days"] = input.RetentionPolicy.Days
	}

	return []interface{}{logging}
}

func flattenCorsProperty(input string) []interface{} {
	results := make([]interface{}, 0, len(input))

	origins := strings.Split(input, ",")
	for _, origin := range origins {
		results = append(results, origin)
	}

	return results
}

func flattenSharePropertiesCorsRule(input *fileservice.CorsRules) []interface{} {
	corsRules := make([]interface{}, 0)

	if input == nil || input.CorsRules == nil {
		return corsRules
	}

	for _, corsRule := range *input.CorsRules {
		corsRules = append(corsRules, map[string]interface{}{
			"allowed_headers":    corsRule.AllowedHeaders,
			"allowed_origins":    corsRule.AllowedOrigins,
			"allowed_methods":    corsRule.AllowedMethods,
			"exposed_headers":    corsRule.ExposedHeaders,
			"max_age_in_seconds": int(corsRule.MaxAgeInSeconds),
		})
	}

	return corsRules
}

func flattenShareProperties(input fileservice.GetServicePropertiesOperationResponse) []interface{} {
	output := make([]interface{}, 0)

	if model := input.Model; model != nil {
		if props := model.Properties; props != nil {
			output = append(output, map[string]interface{}{
				"cors_rule":        flattenSharePropertiesCorsRule(props.Cors),
				"retention_policy": flattenSharePropertiesDeleteRetentionPolicyWithoutPermDeleteOption(props.ShareDeleteRetentionPolicy),
				"smb":              flattenedSharePropertiesSMB(props.ProtocolSettings),
			})
		}
	}

	return output
}

func flattenSharePropertiesDeleteRetentionPolicyWithoutPermDeleteOption(input *fileservice.DeleteRetentionPolicy) []interface{} {
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

func flattenedSharePropertiesSMB(input *fileservice.ProtocolSettings) []interface{} {
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

	if len(versions) == 0 && len(authenticationMethods) == 0 && len(kerberosTicketEncryption) == 0 && len(channelEncryption) == 0 && input.Smb.Multichannel == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"versions":                        versions,
			"authentication_types":            authenticationMethods,
			"kerberos_ticket_encryption_type": kerberosTicketEncryption,
			"channel_encryption_type":         channelEncryption,
			"multichannel_enabled":            multichannelEnabled,
		},
	}
}

func flattenStaticWebsiteProperties(input accounts.GetServicePropertiesResult) []interface{} {
	if staticWebsite := input.StaticWebsite; staticWebsite != nil {
		if !staticWebsite.Enabled {
			return []interface{}{}
		}

		return []interface{}{
			map[string]interface{}{
				"index_document":     staticWebsite.IndexDocument,
				"error_404_document": staticWebsite.ErrorDocument404Path,
			},
		}
	}
	return []interface{}{}
}

func flattenStorageAccountBypass(input storage.Bypass) []interface{} {
	bypassValues := strings.Split(string(input), ", ")
	bypass := make([]interface{}, len(bypassValues))

	for i, value := range bypassValues {
		bypass[i] = value
	}

	return bypass
}

func expandStorageAccountSASPolicy(input []interface{}) *storage.SasPolicy {
	if len(input) == 0 {
		return nil
	}

	e := input[0].(map[string]interface{})

	return &storage.SasPolicy{
		ExpirationAction:    utils.String(e["expiration_action"].(string)),
		SasExpirationPeriod: utils.String(e["expiration_period"].(string)),
	}
}

func flattenStorageAccountSASPolicy(input *storage.SasPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var expirationAction string
	if input.ExpirationAction != nil {
		expirationAction = *input.ExpirationAction
	}

	var expirationPeriod string
	if input.SasExpirationPeriod != nil {
		expirationPeriod = *input.SasExpirationPeriod
	}

	return []interface{}{
		map[string]interface{}{
			"expiration_action": expirationAction,
			"expiration_period": expirationPeriod,
		},
	}
}
