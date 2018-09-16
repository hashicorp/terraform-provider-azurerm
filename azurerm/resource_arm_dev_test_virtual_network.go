package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDevTestVirtualNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDevTestVirtualNetworkCreateUpdate,
		Read:   resourceArmDevTestVirtualNetworkRead,
		Update: resourceArmDevTestVirtualNetworkCreateUpdate,
		Delete: resourceArmDevTestVirtualNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDevTestVirtualNetworkName(),
			},

			"lab_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDevTestLabName(),
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tagsSchema(),

			"unique_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmDevTestVirtualNetworkCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devTestVirtualNetworksClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for DevTest Virtual Network creation")

	name := d.Get("name").(string)
	labName := d.Get("lab_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	description := d.Get("description").(string)
	tags := d.Get("tags").(map[string]interface{})

	parameters := dtl.VirtualNetwork{
		Location: utils.String(location),
		Tags:     expandTags(tags),
		VirtualNetworkProperties: &dtl.VirtualNetworkProperties{
			Description: utils.String(description),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, labName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating DevTest Virtual Network %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation/update of DevTest Virtual Network %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, labName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving DevTest Virtual Network %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read DevTest Virtual Network %q (Lab %q / Resource Group %q) ID", name, labName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmDevTestVirtualNetworkRead(d, meta)
}

func resourceArmDevTestVirtualNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devTestVirtualNetworksClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	labName := id.Path["labs"]
	name := id.Path["virtualnetworks"]

	read, err := client.Get(ctx, resourceGroup, labName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] DevTest Virtual Network %q was not found in Lab %q / Resource Group %q - removing from state!", name, labName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on DevTest Virtual Network %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	d.Set("name", read.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := read.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := read.VirtualNetworkProperties; props != nil {
		d.Set("description", props.Description)

		// Computed fields
		d.Set("unique_identifier", props.UniqueIdentifier)
	}

	flattenAndSetTags(d, read.Tags)

	return nil
}

func resourceArmDevTestVirtualNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devTestVirtualNetworksClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	labName := id.Path["labs"]
	name := id.Path["virtualnetworks"]

	read, err := client.Get(ctx, resourceGroup, labName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] DevTest Virtual Network %q was not found in Lab %q / Resource Group %q - assuming removed!", name, labName, resourceGroup)
			return nil
		}

		return fmt.Errorf("Error retrieving DevTest Virtual Network %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	future, err := client.Delete(ctx, resourceGroup, labName, name)
	if err != nil {
		return fmt.Errorf("Error deleting DevTest Virtual Network %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the deletion of DevTest Virtual Network %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	return err
}

func validateDevTestVirtualNetworkName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[A-Za-z0-9_-]+$"),
		"Virtual Network Name can only include alphanumeric characters, underscores, hyphens.")
}
