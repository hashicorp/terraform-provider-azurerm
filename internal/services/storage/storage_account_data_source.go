// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-09-01/storage" // nolint: staticcheck
	azautorest "github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"access_tier": {
				Type:     pluginsdk.TypeString,
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

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_https_traffic_only": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"min_tls_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"allow_nested_items_to_be_public": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"is_hns_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"nfsv3_enabled": {
				Type:     pluginsdk.TypeBool,
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
					},
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceStorageAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.AccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storageDomainSuffix, ok := meta.(*clients.Client).Account.Environment.Storage.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Storage domain suffix for environment %q", meta.(*clients.Client).Account.Environment.Name)
	}

	id := commonids.NewStorageAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.GetProperties(ctx, id.ResourceGroupName, id.StorageAccountName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

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
				doesntHavePermissions = status == http.StatusUnauthorized || status == http.StatusForbidden
			}
		}

		if !hasWriteLock && !doesntHavePermissions {
			return fmt.Errorf("listing Keys for %s: %+v", id, err)
		}
	}

	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("account_kind", resp.Kind)

	if sku := resp.Sku; sku != nil {
		d.Set("account_tier", sku.Tier)
		d.Set("account_replication_type", strings.Split(string(sku.Name), "_")[1])
	}

	if props := resp.AccountProperties; props != nil {
		d.Set("access_tier", props.AccessTier)
		d.Set("enable_https_traffic_only", props.EnableHTTPSTrafficOnly)
		d.Set("min_tls_version", string(props.MinimumTLSVersion))
		d.Set("is_hns_enabled", props.IsHnsEnabled)
		d.Set("nfsv3_enabled", props.EnableNfsV3)
		d.Set("allow_nested_items_to_be_public", props.AllowBlobPublicAccess)

		if err := d.Set("custom_domain", flattenStorageAccountCustomDomain(props.CustomDomain)); err != nil {
			return fmt.Errorf("setting `custom_domain`: %+v", err)
		}

		// Computed
		d.Set("primary_location", props.PrimaryLocation)
		d.Set("secondary_location", props.SecondaryLocation)

		if accessKeys := keys.Keys; accessKeys != nil {
			storageAccessKeys := *accessKeys
			if len(storageAccessKeys) > 0 {
				pcs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *storageAccessKeys[0].Value, *storageDomainSuffix)
				d.Set("primary_connection_string", pcs)
			}

			if len(storageAccessKeys) > 1 {
				scs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *storageAccessKeys[1].Value, *storageDomainSuffix)
				d.Set("secondary_connection_string", scs)
			}
		}

		if err := flattenAndSetAzureRmStorageAccountPrimaryEndpoints(d, props.PrimaryEndpoints); err != nil {
			return fmt.Errorf("setting primary endpoints and hosts for blob, queue, table and file: %+v", err)
		}

		if accessKeys := keys.Keys; accessKeys != nil {
			var primaryBlobConnectStr string
			if v := props.PrimaryEndpoints; v != nil {
				primaryBlobConnectStr = getBlobConnectionString(v.Blob, resp.Name, (*accessKeys)[0].Value)
			}
			d.Set("primary_blob_connection_string", primaryBlobConnectStr)
		}

		if err := flattenAndSetAzureRmStorageAccountSecondaryEndpoints(d, props.SecondaryEndpoints); err != nil {
			return fmt.Errorf("setting secondary endpoints and hosts for blob, queue, table: %+v", err)
		}

		if accessKeys := keys.Keys; accessKeys != nil {
			var secondaryBlobConnectStr string
			if v := props.SecondaryEndpoints; v != nil {
				secondaryBlobConnectStr = getBlobConnectionString(v.Blob, resp.Name, (*accessKeys)[1].Value)
			}
			d.Set("secondary_blob_connection_string", secondaryBlobConnectStr)
		}

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

		infrastructureEncryption := false
		if encryption := props.Encryption; encryption != nil && encryption.RequireInfrastructureEncryption != nil {
			infrastructureEncryption = *encryption.RequireInfrastructureEncryption
		}
		d.Set("infrastructure_encryption_enabled", infrastructureEncryption)

		if err := d.Set("azure_files_authentication", flattenArmStorageAccountAzureFilesAuthentication(props.AzureFilesIdentityBasedAuthentication)); err != nil {
			return fmt.Errorf("setting `azure_files_authentication`: %+v", err)
		}
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

	return tags.FlattenAndSet(d, resp.Tags)
}
