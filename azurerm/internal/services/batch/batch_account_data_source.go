package batch

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2019-08-01/batch"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmBatchAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmBatchAccountRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateAzureRMBatchAccountName,
			},
			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
			"location":            azure.SchemaLocationForDataSource(),
			"storage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pool_allocation_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_vault_reference": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"primary_access_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"secondary_access_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"account_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmBatchAccountRead(d *schema.ResourceData, meta interface{}) error {
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
