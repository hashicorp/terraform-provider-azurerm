package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmBatchCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmBatchCertificateRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validateAzureRMBatchCertificateName,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMBatchAccountName,
			},
			"resource_group_name": resourceGroupNameSchema(),
			"public_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"thumbprint_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmBatchCertificateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).batchCertificateClient

	resourceGroupName := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)
	name := d.Get("name").(string)

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Get(ctx, resourceGroupName, accountName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Batch certificate %q (Resource Group %q) was not found", name, resourceGroupName)
		}
		return fmt.Errorf("Error making Read request on AzureRM Batch certificate %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("account_name", accountName)
	d.Set("resource_group_name", resourceGroupName)

	if format := resp.Format; format != "" {
		d.Set("format", format)
	}
	if publicData := resp.PublicData; publicData != nil {
		d.Set("public_data", publicData)
	}
	if thumbprint := resp.Thumbprint; thumbprint != nil {
		d.Set("thumbprint", thumbprint)
	}
	if thumbprintAlgorithm := resp.ThumbprintAlgorithm; thumbprintAlgorithm != nil {
		d.Set("thumbprint_algorithm", thumbprintAlgorithm)
	}

	return nil
}
