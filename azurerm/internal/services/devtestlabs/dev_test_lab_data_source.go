package devtestlabs

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmDevTestLab() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDevTestLabRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DevTestLabName(),
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"storage_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),

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
	client := meta.(*clients.Client).DevTestLabs.LabsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	read, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			return fmt.Errorf("DevTest Lab %q was not found in Resource Group %q", name, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on DevTest Lab %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*read.ID)

	d.Set("name", read.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := read.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := read.LabProperties; props != nil {
		d.Set("storage_type", string(props.LabStorageType))

		// Computed fields
		d.Set("artifacts_storage_account_id", props.ArtifactsStorageAccount)
		d.Set("default_storage_account_id", props.DefaultStorageAccount)
		d.Set("default_premium_storage_account_id", props.DefaultPremiumStorageAccount)
		d.Set("key_vault_id", props.VaultName)
		d.Set("premium_data_disk_storage_account_id", props.PremiumDataDiskStorageAccount)
		d.Set("unique_identifier", props.UniqueIdentifier)
	}

	return tags.FlattenAndSet(d, read.Tags)
}
