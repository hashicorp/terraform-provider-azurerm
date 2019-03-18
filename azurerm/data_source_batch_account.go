package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmBatchAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmBatchAccountRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMBatchAccountName,
			},
			"resource_group_name": resourceGroupNameForDataSourceSchema(),
			"location":            locationForDataSourceSchema(),
			"storage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pool_allocation_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_shared_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"secondary_shared_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"account_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmBatchAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batchAccountClient

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	ctx := meta.(*ArmClient).StopContext
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

	if d.Get("pool_allocation_mode") == string(batch.BatchService) {
		keys, err := client.GetKeys(ctx, resourceGroup, name)

		if err != nil {
			return err
		}

		d.Set("primary_shared_key", keys.Primary)
		d.Set("secondary_shared_key", keys.Secondary)
	}

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.AccountProperties; props != nil {
		if autoStorage := props.AutoStorage; autoStorage != nil {
			d.Set("storage_account_id", autoStorage.StorageAccountID)
		}
		d.Set("pool_allocation_mode", props.PoolAllocationMode)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}
