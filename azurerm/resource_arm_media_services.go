package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2018-07-01/media"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMediaServices() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMediaServicesCreateUpdate,
		Read:   resourceArmMediaServicesRead,
		Update: resourceArmMediaServicesCreateUpdate,
		Delete: resourceArmMediaServicesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,50}$"),
					"Cosmos DB Account name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tags": tagsSchema(),

			"storage_account": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceArmMediaServicesCreateUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ArmClient).mediaServicesClient
	ctx := meta.(*ArmClient).StopContext

	accountName := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})
	resourceGroup := d.Get("resource_group_name").(string)
	storageAccount := d.Get("storage_account").(string)

	parameters := media.Service{
		ServiceProperties: &media.ServiceProperties{
			StorageAccounts: &[]media.StorageAccount{
				{
					ID:   utils.String(storageAccount),
					Type: media.Primary,
				},
			},
		},
		Location: utils.String(location),
		Tags:     expandTags(tags),
	}

	service, err := client.CreateOrUpdate(ctx, resourceGroup, accountName, parameters)
	if err != nil {
		return fmt.Errorf("Error creating media service account: %+v", err)
	}

	d.SetId(*service.ID)

	return nil
}

func resourceArmMediaServicesRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ArmClient).mediaServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q: %+v", d.Id(), err)
	}

	name := id.Path["mediaservices"]
	resourceGroup := id.ResourceGroup
	
	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Media Services Account %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Media Services Account: %+v", err)
	}

	d.Set("name", resp.Name)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmMediaServicesDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ArmClient).mediaServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q: %+v", d.Id(), err)
	}

	name := id.Path["mediaservices"]
	resourceGroup := id.ResourceGroup

	httpResponse, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(httpResponse.Response) {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for Media Services Account '%s': %+v", name, err)
	}

	return nil
}
