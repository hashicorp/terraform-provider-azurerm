package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPostgreSQLDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPostgreSQLDatabaseCreate,
		Read:   resourceArmPostgreSQLDatabaseRead,
		Delete: resourceArmPostgreSQLDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"charset": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ForceNew:         true,
			},

			"collation": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatabaseCollation,
			},
		},
	}
}

func resourceArmPostgreSQLDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgres.DatabasesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Database creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	charset := d.Get("charset").(string)
	collation := d.Get("collation").(string)

	if features.ShouldResourcesBeImported() {
		existing, err := client.Get(ctx, resGroup, serverName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing PostgreSQL Database %s (resource group %s) ID", name, resGroup)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_postgresql_database", *existing.ID)
		}
	}

	properties := postgresql.Database{
		DatabaseProperties: &postgresql.DatabaseProperties{
			Charset:   utils.String(charset),
			Collation: utils.String(collation),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, serverName, name, properties)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read PostgreSQL Database %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPostgreSQLDatabaseRead(d, meta)
}

func resourceArmPostgreSQLDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgres.DatabasesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["databases"]

	resp, err := client.Get(ctx, resGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] PostgreSQL Database '%s' was not found (resource group '%s')", name, resGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure PostgreSQL Database %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("server_name", serverName)

	if props := resp.DatabaseProperties; props != nil {
		d.Set("charset", props.Charset)

		if collation := props.Collation; collation != nil {
			v := strings.Replace(*collation, "-", "_", -1)
			d.Set("collation", v)
		}
	}

	return nil
}

func resourceArmPostgreSQLDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgres.DatabasesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
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

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	return nil
}
