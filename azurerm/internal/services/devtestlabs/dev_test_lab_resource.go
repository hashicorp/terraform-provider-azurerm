package devtestlabs

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDevTest() *schema.Resource {
	return &schema.Resource{
		Create: resourceDevTestCreateUpdate,
		Read:   resourceDevTestRead,
		Update: resourceDevTestCreateUpdate,
		Delete: resourceDevTestDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DevTestLabName(),
			},

			"location": azure.SchemaLocation(),

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"storage_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(dtl.Premium),
				ValidateFunc: validation.StringInSlice([]string{
					string(dtl.Standard),
					string(dtl.Premium),
				}, false),
			},

			"tags": tags.Schema(),

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

func resourceDevTestCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for DevTest Lab creation")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Dev Test Lab %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dev_test_lab", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	storageType := d.Get("storage_type").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := dtl.Lab{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
		LabProperties: &dtl.LabProperties{
			LabStorageType: dtl.StorageType(storageType),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating DevTest Lab %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of DevTest Lab %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving DevTest Lab %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read DevTest Lab %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceDevTestRead(d, meta)
}

func resourceDevTestRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["labs"]

	read, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] DevTest Lab %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on DevTest Lab %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

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

func resourceDevTestDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.LabsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["labs"]

	read, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] DevTest Lab %q was not found in Resource Group %q - assuming removed!", name, resourceGroup)
			return nil
		}

		return fmt.Errorf("Error retrieving DevTest Lab %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting DevTest Lab %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of DevTest Lab %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return err
}
