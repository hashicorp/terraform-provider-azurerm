package azurerm

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmDevTestLab() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDevTestLabRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DevTestLabName(),
			},

			"location": locationForDataSourceSchema(),

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"storage_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsForDataSourceSchema(),

			"artifacts_storage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_storage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_premium_storage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"key_vault_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"premium_data_disk_storage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"unique_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmDevTestLabRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devTestLabsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("DevTest Lab %q was not found in Resource Group %q", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on DevTest Lab %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if err := azure.FlattenAndSetLocation(d, resp.Location); err != nil {
		return err
	}

	if props := resp.LabProperties; props != nil {
		d.Set("storage_type", string(props.LabStorageType))

		// Computed fields
		d.Set("artifacts_storage_account_id", props.ArtifactsStorageAccount)
		d.Set("default_storage_account_id", props.DefaultStorageAccount)
		d.Set("default_premium_storage_account_id", props.DefaultPremiumStorageAccount)
		d.Set("key_vault_id", props.VaultName)
		d.Set("premium_data_disk_storage_account_id", props.PremiumDataDiskStorageAccount)
		d.Set("unique_identifier", props.UniqueIdentifier)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}
