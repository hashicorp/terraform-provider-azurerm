package mariadb

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mariadb/mgmt/2018-06-01/mariadb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mariadb/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mariadb/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMariaDbDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMariaDbDatabaseCreateUpdate,
		Read:   resourceMariaDbDatabaseRead,
		Delete: resourceMariaDbDatabaseDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.MariaDBDatabaseID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[_a-zA-Z0-9]{1,64}$"),
					"name must be 1 - 64 characters long, and contain only letters, numbers and underscores",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},

			"charset": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[a-z0-9]{3,8}$"),
					"charset must be 3 - 8 characters long, and contain only lowercase letters and numbers",
				),
			},

			"collation": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[A-Za-z0-9_. ]+$"),
					"collation must contain only alphanumeric, underscore, space and dot characters",
				),
			},
		},
	}
}

func resourceMariaDbDatabaseCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.DatabasesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM MariaDB database creation")

	id := parse.NewMariaDBDatabaseID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.DatabaseName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_mariadb_database", id.ID())
		}
	}

	charset := d.Get("charset").(string)
	collation := d.Get("collation").(string)

	properties := mariadb.Database{
		DatabaseProperties: &mariadb.DatabaseProperties{
			Charset:   utils.String(charset),
			Collation: utils.String(collation),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.DatabaseName, properties)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMariaDbDatabaseRead(d, meta)
}

func resourceMariaDbDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MariaDBDatabaseID(d.Id())
	if err != nil {
		return fmt.Errorf("cannot parse MariaDB database %q ID:\n%+v", d.Id(), err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.DatabaseName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request on %s:\n%+v", *id, err)
	}

	d.Set("name", id.DatabaseName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("server_name", id.ServerName)

	if properties := resp.DatabaseProperties; properties != nil {
		d.Set("charset", properties.Charset)
		d.Set("collation", properties.Collation)
	}

	return nil
}

func resourceMariaDbDatabaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MariaDBDatabaseID(d.Id())
	if err != nil {
		return fmt.Errorf("cannot parse MariaDB database %q ID:\n%+v", d.Id(), err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.DatabaseName)
	if err != nil {
		return fmt.Errorf("making delete request on %s:\n%+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s:\n%+v", *id, err)
	}

	return nil
}
