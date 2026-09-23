// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceStorageAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceStorageAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.StorageAccountName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"account_kind": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"account_tier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"account_replication_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"provisioned_billing_model_version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"access_tier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"blob_properties": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"change_feed_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"change_feed_retention_in_days": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"container_delete_retention_policy": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"days": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"cors_rule": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"allowed_headers": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"allowed_methods": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"allowed_origins": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"exposed_headers": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"max_age_in_seconds": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"default_service_version": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"delete_retention_policy": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"days": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},

									"permanent_delete_enabled": {
										Type:     pluginsdk.TypeBool,
										Computed: true,
									},
								},
							},
						},

						"last_access_time_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"restore_policy": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"days": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"versioning_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"cross_tenant_replication_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"custom_domain": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"customer_managed_key": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_key_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"user_assigned_identity_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"https_traffic_only_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"default_to_oauth_authentication": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"edge_zone": commonschema.EdgeZoneComputed(),

			"immutability_policy": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"allow_protected_append_writes": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"period_since_creation_in_days": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"state": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"min_tls_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"allow_nested_items_to_be_public": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"allowed_copy_scope": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"is_hns_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"large_file_share_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"local_user_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"network_rules": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"bypass": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							Set: pluginsdk.HashString,
						},

						"default_action": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"ip_rules": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							Set: pluginsdk.HashString,
						},

						"private_link_access": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"endpoint_resource_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"endpoint_tenant_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},

						"virtual_network_subnet_ids": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							Set: pluginsdk.HashString,
						},
					},
				},
			},

			"nfsv3_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"public_network_access": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"routing": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"choice": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"publish_internet_endpoints": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"publish_microsoft_endpoints": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"sas_policy": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"expiration_action": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"expiration_period": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"secondary_location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sftp_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"share_properties": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cors_rule": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"allowed_headers": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"allowed_methods": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"allowed_origins": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"exposed_headers": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},

									"max_age_in_seconds": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"retention_policy": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"days": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
								},
							},
						},

						"smb": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"authentication_types": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
										Set: pluginsdk.HashString,
									},

									"channel_encryption_type": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
										Set: pluginsdk.HashString,
									},

									"kerberos_ticket_encryption_type": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
										Set: pluginsdk.HashString,
									},

									"multichannel_enabled": {
										Type:     pluginsdk.TypeBool,
										Computed: true,
									},

									"versions": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
										Set: pluginsdk.HashString,
									},
								},
							},
						},
					},
				},
			},

			"shared_access_key_enabled": {
				Type:     pluginsdk.TypeBool,
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

			"queue_encryption_key_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"table_encryption_key_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"infrastructure_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"azure_files_authentication": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"directory_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"active_directory": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"domain_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"netbios_domain_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"forest_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"domain_guid": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"domain_sid": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"storage_sid": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},
						"default_share_level_permission": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"dns_endpoint_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceStorageAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	storageClient := meta.(*clients.Client).Storage.ResourceManager
	client := storageClient.StorageAccounts
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageDomainSuffix, ok := meta.(*clients.Client).Account.Environment.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Storage domain suffix for environment %q", meta.(*clients.Client).Account.Environment.Name)
	}

	id := commonids.NewStorageAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.GetProperties(ctx, id, storageaccounts.DefaultGetPropertiesOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	var accountKind storageaccounts.Kind
	var accountTier storageaccounts.SkuTier
	accountReplicationType := ""

	listKeysOpts := storageaccounts.DefaultListKeysOperationOptions()
	listKeysOpts.Expand = pointer.To(storageaccounts.ExpandKerb)
	keys, err := client.ListKeys(ctx, id, listKeysOpts)
	if err != nil {
		hasWriteLock := response.WasConflict(keys.HttpResponse)
		doesntHavePermissions := response.WasForbidden(keys.HttpResponse) || response.WasStatusCode(keys.HttpResponse, http.StatusUnauthorized)
		if !hasWriteLock && !doesntHavePermissions {
			return fmt.Errorf("listing Keys for %s: %+v", id, err)
		}
	}

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("edge_zone", flattenEdgeZone(model.ExtendedLocation))

		accountKind = pointer.From(model.Kind)
		d.Set("account_kind", string(accountKind))

		if sku := model.Sku; sku != nil {
			skuNameSegments := strings.Split(string(sku.Name), "_")
			accountTier = pointer.From(sku.Tier)
			accountReplicationType = skuNameSegments[1]

			d.Set("account_tier", accountTier)
			d.Set("account_replication_type", accountReplicationType)

			provisionedBillingModelVersion := ""
			if accountTier != "" {
				provisionedBillingModelVersion = strings.TrimPrefix(skuNameSegments[0], string(accountTier))
			}
			d.Set("provisioned_billing_model_version", provisionedBillingModelVersion)
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

		if props := model.Properties; props != nil {
			d.Set("access_tier", pointer.From(props.AccessTier))
			d.Set("allow_nested_items_to_be_public", pointer.From(props.AllowBlobPublicAccess))
			d.Set("allowed_copy_scope", pointer.FromEnum(props.AllowedCopyScope))
			if err := d.Set("custom_domain", flattenAccountCustomDomain(props.CustomDomain)); err != nil {
				return fmt.Errorf("setting `custom_domain`: %+v", err)
			}

			customerManagedKey, err := flattenAccountCustomerManagedKey(props.Encryption)
			if err != nil {
				return fmt.Errorf("flattening `customer_managed_key`: %+v", err)
			}
			if err := d.Set("customer_managed_key", customerManagedKey); err != nil {
				return fmt.Errorf("setting `customer_managed_key`: %+v", err)
			}

			d.Set("cross_tenant_replication_enabled", pointer.From(props.AllowCrossTenantReplication))
			d.Set("default_to_oauth_authentication", pointer.From(props.DefaultToOAuthAuthentication))
			d.Set("https_traffic_only_enabled", pointer.From(props.SupportsHTTPSTrafficOnly))

			if err := d.Set("immutability_policy", flattenAccountImmutabilityPolicy(props.ImmutableStorageWithVersioning)); err != nil {
				return fmt.Errorf("setting `immutability_policy`: %+v", err)
			}

			d.Set("is_hns_enabled", pointer.From(props.IsHnsEnabled))

			largeFileShareEnabled := pointer.From(props.LargeFileSharesState) == storageaccounts.LargeFileSharesStateEnabled
			d.Set("large_file_share_enabled", largeFileShareEnabled)

			isLocalEnabled := true
			if props.IsLocalUserEnabled != nil {
				isLocalEnabled = *props.IsLocalUserEnabled
			}
			d.Set("local_user_enabled", isLocalEnabled)

			if err := d.Set("network_rules", flattenAccountNetworkRules(props.NetworkAcls)); err != nil {
				return fmt.Errorf("setting `network_rules`: %+v", err)
			}

			d.Set("nfsv3_enabled", pointer.From(props.IsNfsV3Enabled))
			d.Set("primary_location", location.NormalizeNilable(props.PrimaryLocation))

			if err := d.Set("routing", flattenAccountRoutingPreference(props.RoutingPreference)); err != nil {
				return fmt.Errorf("setting `routing`: %+v", err)
			}

			if err := d.Set("sas_policy", flattenAccountSASPolicy(props.SasPolicy)); err != nil {
				return fmt.Errorf("setting `sas_policy`: %+v", err)
			}

			d.Set("secondary_location", location.NormalizeNilable(props.SecondaryLocation))
			d.Set("sftp_enabled", pointer.From(props.IsSftpEnabled))

			allowSharedKeyAccess := true
			if props.AllowSharedKeyAccess != nil {
				allowSharedKeyAccess = *props.AllowSharedKeyAccess
			}
			d.Set("shared_access_key_enabled", allowSharedKeyAccess)

			d.Set("public_network_access", pointer.FromEnum(props.PublicNetworkAccess))

			// Setting the encryption key type to "Service" in PUT. The following GET will not return the queue/table in the service list of its response.
			// So defaults to setting the encryption key type to "Service" if it is absent in the GET response. Also, define the default value as "Service" in the schema.
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

			// For storage account created using old API, the response of GET call will not return "min_tls_version"
			minTlsVersion := string(storageaccounts.MinimumTlsVersionTLSOneZero)
			if props.MinimumTlsVersion != nil {
				minTlsVersion = string(*props.MinimumTlsVersion)
			}
			d.Set("min_tls_version", minTlsVersion)

			// DNSEndpointType is null when unconfigured - therefore default this to Standard
			dnsEndpointType := storageaccounts.DnsEndpointTypeStandard
			if props.DnsEndpointType != nil {
				dnsEndpointType = *props.DnsEndpointType
			}
			d.Set("dns_endpoint_type", string(dnsEndpointType))

			if err := d.Set("azure_files_authentication", flattenAccountAzureFilesAuthentication(props.AzureFilesIdentityBasedAuthentication)); err != nil {
				return fmt.Errorf("setting `azure_files_authentication`: %+v", err)
			}
		}
	}

	var primaryEndpoints *storageaccounts.Endpoints
	var secondaryEndpoints *storageaccounts.Endpoints
	var routingPreference *storageaccounts.RoutingPreference
	if model := resp.Model; model != nil && model.Properties != nil {
		primaryEndpoints = model.Properties.PrimaryEndpoints
		routingPreference = model.Properties.RoutingPreference
		secondaryEndpoints = model.Properties.SecondaryEndpoints
	}
	endpoints := flattenAccountEndpoints(primaryEndpoints, secondaryEndpoints, routingPreference)
	endpoints.set(d)

	storageAccountKeys := make([]storageaccounts.StorageAccountKey, 0)
	if keys.Model != nil && keys.Model.Keys != nil {
		storageAccountKeys = *keys.Model.Keys
	}
	keysAndConnectionStrings := flattenAccountAccessKeysAndConnectionStrings(id.StorageAccountName, *storageDomainSuffix, storageAccountKeys, endpoints)
	keysAndConnectionStrings.set(d)

	supportLevel := availableFunctionalityForAccount(accountKind, accountTier, accountReplicationType)

	blobProperties := make([]interface{}, 0)
	if supportLevel.supportBlob {
		blobProps, err := storageClient.BlobServices.GetServiceProperties(ctx, id)
		if err != nil {
			return fmt.Errorf("reading blob properties for %s: %+v", id, err)
		}

		blobProperties = flattenAccountBlobServiceProperties(blobProps.Model)
	}
	if err := d.Set("blob_properties", blobProperties); err != nil {
		return fmt.Errorf("setting `blob_properties` for %s: %+v", id, err)
	}

	shareProperties := make([]interface{}, 0)
	if supportLevel.supportShare {
		shareProps, err := storageClient.FileServices.GetServiceProperties(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving share properties for %s: %+v", id, err)
		}

		shareProperties = flattenAccountShareProperties(shareProps.Model)
	}
	if err := d.Set("share_properties", shareProperties); err != nil {
		return fmt.Errorf("setting `share_properties` for %s: %+v", id, err)
	}

	return nil
}
