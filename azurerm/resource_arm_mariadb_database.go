package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/preview/mariadb/mgmt/2018-06-01-preview/mariadb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMariaDbDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMariaDbDatabaseCreateOrUpdate,
		Read:   resourceArmMariaDbDatabaseRead,
		Delete: resourceArmMariaDbDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[_a-zA-Z0-9]{1,64}$"),
					"name must be 1 - 64 characters long, and contain only letters, numbers and underscores.",
				),
			},

			"resource_group_name": resourceGroupNameSchema(),

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{3,50}$"),
					"server_name must be 3 - 50 characters long, and contain only letters, numbers and hyphens.",
				),
			},

			"charset": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-z0-9]{3,8}$"),
					"charset must be 3 - 8 characters long, and contain only lowercase letters and numbers.",
				),
			},

			"collation": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCollation(),
			},
		},
	}
}

func resourceArmMariaDbDatabaseCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mariadbDatabasesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM MariaDB Database creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	charset := d.Get("charset").(string)
	collation := d.Get("collation").(string)

	properties := mariadb.Database{
		DatabaseProperties: &mariadb.DatabaseProperties{
			Charset:   utils.String(charset),
			Collation: utils.String(collation),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, serverName, name, properties)
	if err != nil {
		return fmt.Errorf("Error creating Maria DB %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Maria DB %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read MariaDB Database %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMariaDbDatabaseRead(d, meta)
}

func resourceArmMariaDbDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mariadbDatabasesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["databases"]

	resp, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] MariaDB Database '%s' was not found (resource group '%s')", name, resGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure MariaDB Database %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("server_name", serverName)

	if props := resp.DatabaseProperties; props != nil {
		d.Set("charset", props.Charset)

		if collation := props.Collation; collation != nil {
			d.Set("collation", collation)
		}
	}

	return nil
}

func resourceArmMariaDbDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mariadbDatabasesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["databases"]

	future, err := client.Delete(ctx, resGroup, serverName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	return nil
}
