package batch

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2020-03-01/batch"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceBatchAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBatchAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName,
			},
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
			"location":            azure.SchemaLocationForDataSource(),
			"storage_account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"pool_allocation_mode": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"key_vault_reference": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"url": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"account_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceBatchAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.AccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Batch account %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on AzureRM Batch account %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("account_endpoint", resp.AccountEndpoint)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.AccountProperties; props != nil {
		if autoStorage := props.AutoStorage; autoStorage != nil {
			d.Set("storage_account_id", autoStorage.StorageAccountID)
		}
		d.Set("pool_allocation_mode", props.PoolAllocationMode)
		poolAllocationMode := d.Get("pool_allocation_mode").(string)

		if poolAllocationMode == string(batch.BatchService) {
			keys, err := client.GetKeys(ctx, resourceGroup, name)
			if err != nil {
				return fmt.Errorf("Cannot read keys for Batch account %q (resource group %q): %v", name, resourceGroup, err)
			}

			d.Set("primary_access_key", keys.Primary)
			d.Set("secondary_access_key", keys.Secondary)

			// set empty keyvault reference which is not needed in Batch Service allocation mode.
			d.Set("key_vault_reference", []interface{}{})
		} else if poolAllocationMode == string(batch.UserSubscription) {
			if err := d.Set("key_vault_reference", flattenBatchAccountKeyvaultReference(props.KeyVaultReference)); err != nil {
				return fmt.Errorf("Error flattening `key_vault_reference`: %+v", err)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
