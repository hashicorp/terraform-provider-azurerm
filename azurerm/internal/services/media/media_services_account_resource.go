package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2018-07-01/media"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMediaServicesAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMediaServicesAccountCreateUpdate,
		Read:   resourceArmMediaServicesAccountRead,
		Update: resourceArmMediaServicesAccountCreateUpdate,
		Delete: resourceArmMediaServicesAccountDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.MediaServicesAccountID(id)
			return err
		}),

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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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

			// TODO: support Tags when this bug is fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/5249
			// "tags": tags.Schema(),
		},
	}
}

func resourceArmMediaServicesAccountCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.ServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	accountName := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
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
	}

	if _, e := client.CreateOrUpdate(ctx, resourceGroup, accountName, parameters); e != nil {
		return fmt.Errorf("Error creating Media Service Account %q (Resource Group %q): %+v", accountName, resourceGroup, e)
	}

	service, err := client.Get(ctx, resourceGroup, accountName)
	if err != nil {
		return fmt.Errorf("Error retrieving Media Service Account %q (Resource Group %q): %+v", accountName, resourceGroup, err)
	}
	d.SetId(*service.ID)

	return resourceArmMediaServicesAccountRead(d, meta)
}

func resourceArmMediaServicesAccountRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.ServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MediaServicesAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Media Services Account %q was not found in Resource Group %q - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Media Services Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ServiceProperties; props != nil {
		accounts := flattenMediaServicesAccountStorageAccounts(props.StorageAccounts)
		if e := d.Set("storage_account", accounts); e != nil {
			return fmt.Errorf("Error flattening `storage_account`: %s", e)
		}
	}

	// TODO: support Tags when this bug is fixed:
	// https://github.com/Azure/azure-rest-api-specs/issues/5249
	// return tags.FlattenAndSet(d, resp.Tags)
	return nil
}

func resourceArmMediaServicesAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.ServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MediaServicesAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for Media Services Account '%s': %+v", id.Name, err)
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
		output := make(map[string]interface{})

		if storageAccount.ID != nil {
			output["id"] = *storageAccount.ID
		}

		output["is_primary"] = storageAccount.Type == media.Primary

		results = append(results, output)
	}

	return results
}
