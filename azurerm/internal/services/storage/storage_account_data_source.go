package storage

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2021-01-01/storage"
	azautorest "github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

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

			"enable_https_traffic_only": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"min_tls_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"allow_blob_public_access": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"is_hns_enabled": {
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

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceStorageAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Storage.AccountsClient
	endpointSuffix := meta.(*clients.Client).Account.Environment.StorageEndpointSuffix
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.GetProperties(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Storage Account %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading the state of AzureRM Storage Account %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	// handle the user not having permissions to list the keys
	d.Set("primary_connection_string", "")
	d.Set("secondary_connection_string", "")
	d.Set("primary_blob_connection_string", "")
	d.Set("secondary_blob_connection_string", "")
	d.Set("primary_access_key", "")
	d.Set("secondary_access_key", "")

	keys, err := client.ListKeys(ctx, resourceGroup, name, storage.Kerb)
	if err != nil {
		// the API returns a 200 with an inner error of a 409..
		var hasWriteLock bool
		var doesntHavePermissions bool
		if e, ok := err.(azautorest.DetailedError); ok {
			if status, ok := e.StatusCode.(int); ok {
				hasWriteLock = status == http.StatusConflict
				doesntHavePermissions = status == http.StatusUnauthorized
			}
		}

		if !hasWriteLock && !doesntHavePermissions {
			return fmt.Errorf("Error listing Keys for Storage Account %q (Resource Group %q): %s", name, resourceGroup, err)
		}
	}

	accountKeys := keys.Keys
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
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
		d.Set("allow_blob_public_access", props.AllowBlobPublicAccess)

		if customDomain := props.CustomDomain; customDomain != nil {
			if err := d.Set("custom_domain", flattenStorageAccountCustomDomain(customDomain)); err != nil {
				return fmt.Errorf("Error setting `custom_domain`: %+v", err)
			}
		}

		// Computed
		d.Set("primary_location", props.PrimaryLocation)
		d.Set("secondary_location", props.SecondaryLocation)

		if accessKeys := accountKeys; accessKeys != nil {
			storageAccessKeys := *accessKeys
			if len(storageAccessKeys) > 0 {
				pcs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *storageAccessKeys[0].Value, endpointSuffix)
				d.Set("primary_connection_string", pcs)
			}

			if len(storageAccessKeys) > 1 {
				scs := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", *resp.Name, *storageAccessKeys[1].Value, endpointSuffix)
				d.Set("secondary_connection_string", scs)
			}
		}

		if err := flattenAndSetAzureRmStorageAccountPrimaryEndpoints(d, props.PrimaryEndpoints); err != nil {
			return fmt.Errorf("error setting primary endpoints and hosts for blob, queue, table and file: %+v", err)
		}

		if accessKeys := accountKeys; accessKeys != nil {
			var primaryBlobConnectStr string
			if v := props.PrimaryEndpoints; v != nil {
				primaryBlobConnectStr = getBlobConnectionString(v.Blob, resp.Name, (*accessKeys)[0].Value)
			}
			d.Set("primary_blob_connection_string", primaryBlobConnectStr)
		}

		if err := flattenAndSetAzureRmStorageAccountSecondaryEndpoints(d, props.SecondaryEndpoints); err != nil {
			return fmt.Errorf("error setting secondary endpoints and hosts for blob, queue, table: %+v", err)
		}

		if accessKeys := accountKeys; accessKeys != nil {
			var secondaryBlobConnectStr string
			if v := props.SecondaryEndpoints; v != nil {
				secondaryBlobConnectStr = getBlobConnectionString(v.Blob, resp.Name, (*accessKeys)[1].Value)
			}
			d.Set("secondary_blob_connection_string", secondaryBlobConnectStr)
		}
	}

	if accessKeys := accountKeys; accessKeys != nil {
		storageAccountKeys := *accessKeys
		d.Set("primary_access_key", storageAccountKeys[0].Value)
		d.Set("secondary_access_key", storageAccountKeys[1].Value)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
