package batch

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/batch/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceBatchCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBatchCertificateRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.CertificateName,
			},

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"public_data": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"format": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"thumbprint_algorithm": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBatchCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.CertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
