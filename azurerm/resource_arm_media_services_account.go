package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2018-07-01/media"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
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
					regexp.MustCompile("^[-a-z0-9]{3,24}$"),
					"Media Services Account name must be 3 - 24 characters long, contain only lowercase letters and numbers.",
				),
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"storage_account": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"is_primary": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"tags": tagsSchema(),

			"media_service_account_id": {
				Type:     schema.TypeString,
				Computed: true,
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

	storageAccountsRaw := d.Get("storage_account").(*schema.Set).List()
	storageAccounts, err := expandMediaServicesAccountStorageAccounts(storageAccountsRaw)
	if err != nil {
		return err
	}

	parameters := media.Service{
		ServiceProperties: &media.ServiceProperties{
			StorageAccounts: storageAccounts,
		},
		Location: utils.String(location),
		Tags:     expandTags(tags),
	}

	service, err := client.CreateOrUpdate(ctx, resourceGroup, accountName, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Media Service Account %q (Resource Group %q): %+v", accountName, resourceGroup, err)
	}

	d.SetId(*service.ID)

	return nil
}

func resourceArmMediaServicesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mediaServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["mediaservices"]
	resourceGroup := id.ResourceGroup

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Media Services Account %q was not found in Resource Group %q - removing from state", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Media Services Account %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.ServiceProperties; props != nil {
		d.Set("media_service_account_id", props.MediaServiceID.String())
		accounts := flattenMediaServicesAccountStorageAccounts(props.StorageAccounts)
		if e := d.Set("storage_account", accounts); e != nil {
			return fmt.Errorf("Error flattening `storage_account`: %s", e)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmMediaServicesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mediaServicesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["mediaservices"]
	resourceGroup := id.ResourceGroup

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for Media Services Account '%s': %+v", name, err)
	}

	return nil
}

func expandMediaServicesAccountStorageAccounts(input []interface{}) (*[]media.StorageAccount, error) {
	results := make([]media.StorageAccount, 0)

	foundPrimary := false
	for _, accountMapRaw := range input {
		accountMap := accountMapRaw.(map[string]interface{})

		id := accountMap["id"].(string)

		storageType := media.Secondary
		if accountMap["is_primary"].(bool) {
			if foundPrimary {
				return nil, fmt.Errorf("Only one Storage Account can be set as Primary")
			}

			storageType = media.Primary
			foundPrimary = true
		}

		storageAccount := media.StorageAccount{
			ID:   utils.String(id),
			Type: storageType,
		}

		results = append(results, storageAccount)
	}

	return &results, nil
}

func flattenMediaServicesAccountStorageAccounts(input *[]media.StorageAccount) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, storageAccount := range *input {
		output := make(map[string]interface{}, 0)

		if storageAccount.ID != nil {
			output["id"] = *storageAccount.ID
		}

		output["is_primary"] = storageAccount.Type == media.Primary

		results = append(results, output)
	}

	return results
}
