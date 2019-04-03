package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/cosmos"
	"log"
	"regexp"
	"strings"

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
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"offer_throughput": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntBetweenAndDivisibleBy(400, 1000000, 100),
			},
		},
	}
}

func cosmosBuildResourceId(accountName, path string) string {
	return "cosmos/" + accountName + "/" + path
}

func cosmosSplitResourceId(id string) (account, path string) {
	stripped := strings.TrimPrefix(id, "cosmos/")
	parts := strings.SplitN(stripped, "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
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
			id := cosmosBuildResourceId(account, cosmos.BuildDatabasePath(name).String)
			return tf.ImportAsExistsError("azurerm_cosmosdb_database", id)
		}
	}

	db := cosmos.Database{}

	if v, ok := d.GetOkExists("offer_throughput"); ok {
		db.OfferThroughput = utils.Int(v.(int))
	}

	if _, err := client.Create(ctx, name, db); err != nil {
		return fmt.Errorf("Error creating Cosmos Database %s (Account %s): %+v", name, account, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("Error creating Cosmos Database %s (Account %s): %+v", name, account, err)
	}

	//path is unique, id is actually the name
	d.SetId(cosmosBuildResourceId(client.AccountName, resp.Path))

	return resourceArmCosmosDatabaseRead(d, meta)
}

func resourceArmCosmosDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext

	account, path := cosmosSplitResourceId(d.Id())
	p, err := cosmos.ParseDatabasePath(path)
	if err != nil {
		return fmt.Errorf("Unable to parsing ID (%s) into a Cosmos Database path into a path: %+v", d.Id(), err)
	}

	key := d.Get("account_key").(string)
	client := meta.(*ArmClient).getCosmosDatabasesClient(account, key)

	resp, err := client.Get(ctx, p.Database)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Cosmos Database %s (Account %s) - removing from state", p.Database, account)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Cosmos Database %s (Account %s): %+v", p.Database, account, err)
	}

	d.Set("name", resp.ID)
	d.Set("offer_throughput", resp.OfferThroughput)

	return nil
}

func resourceArmCosmosDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	ctx := meta.(*ArmClient).StopContext

	account, path := cosmosSplitResourceId(d.Id())
	p, err := cosmos.ParseDatabasePath(path)
	if err != nil {
		return fmt.Errorf("Unable to parsing ID (%s) into a Cosmos Database path into a path: %+v", d.Id(), err)
	}

	key := d.Get("account_key").(string)
	client := meta.(*ArmClient).getCosmosDatabasesClient(account, key)

	resp, err := client.Delete(ctx, p.Database)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Cosmos Database %s (Account %s): %+v", p.Database, account, err)
		}
	}

	return nil
}
