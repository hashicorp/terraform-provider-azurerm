package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"log"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCosmosDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCosmosDatabaseCreate,
		Read:   resourceArmCosmosDatabaseRead,
		Delete: resourceArmCosmosDatabaseDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},

			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,50}$"),
					"Cosmos DB Account name must be 3 - 50 characters long, contain only letters, numbers and hyphens.",
				),
			},

			"account_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				//Sensitive: true,
			},

			//x offer throughput
		},
	}
}

func resourceArmCosmosDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	account := d.Get("account_name").(string)
	key := d.Get("account_key").(string)

	client := meta.(*ArmClient).getCosmosDatabasesClient(account, key)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing resource group: %+v", err)
			}
		}

		if existing.Path != "" {
			return tf.ImportAsExistsError("azurerm_cosmosdb_database", existing.Path)
		}
	}

	if _, err := client.Create(ctx, name); err != nil {
		return fmt.Errorf("Error creating Cosmos Database %s (Account %s): %+v", name, account, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("Error creating Cosmos Database %s (Account %s): %+v", name, account, err)
	}

	//path is unique, id is actually the name
	d.SetId(resp.Path)

	return resourceArmCosmosDatabaseRead(d, meta)
}

func resourceArmCosmosDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	account := d.Get("account_name").(string)
	key := d.Get("account_key").(string)

	client := meta.(*ArmClient).getCosmosDatabasesClient(account, key)

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Database %s (Account %s) - removing from state", name, account)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Database %s (Account %s): %+v", name, account, err)
	}

	d.Set("name", resp.ID)
	d.Set("account_name", account)
	d.Set("account_key", key)

	//d.Set("name", resp.Path)

	return nil
}

func resourceArmCosmosDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext

	/*id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Azure Resource ID %q: %+v", d.Id(), err)
	}*/

	name := d.Get("name").(string)
	account := d.Get("account_name").(string)
	key := d.Get("account_key").(string)

	client := meta.(*ArmClient).getCosmosDatabasesClient(account, key)

	resp, err := client.Delete(ctx, name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Cosmos Database %s (Account %s): %+v", name, account, err)
		}
	}

	return nil
}
