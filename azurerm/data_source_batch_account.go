package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
			"resource_group_name": resourceGroupNameDiffSuppressSchema(),
			"location":            locationForDataSourceSchema(),
			"storage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pool_allocation_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          string(batch.BatchService),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(batch.BatchService),
					string(batch.UserSubscription),
				}, true),
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

	// todo : fetch properties

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}
