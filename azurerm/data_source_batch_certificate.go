package azurerm

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmBatchCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmBatchCertificateRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAzureRMBatchCertificateName,
			},

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAzureRMBatchAccountName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

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
	ctx := meta.(*ArmClient).StopContext
	client := meta.(*ArmClient).batch.CertificateClient

	resourceGroupName := d.Get("resource_group_name").(string)
	accountName := d.Get("account_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroupName, accountName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Batch certificate %q (Account %q / Resource Group %q) was not found", name, accountName, resourceGroupName)
		}
		return fmt.Errorf("Error making Read request on AzureRM Batch certificate %q: %+v", name, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("account_name", accountName)
	d.Set("resource_group_name", resourceGroupName)

	if props := resp.CertificateProperties; props != nil {
		d.Set("format", props.Format)
		d.Set("public_data", props.PublicData)
		d.Set("thumbprint", props.Thumbprint)
		d.Set("thumbprint_algorithm", props.ThumbprintAlgorithm)
	}

	return nil
}
func validateAzureRMBatchCertificateName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if !regexp.MustCompile(`^[\w]+-[\w]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"must be made up of algorithm and thumbprint separated by a dash in %q: %q", k, value))
	}

	return warnings, errors
}
