package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2021-06-01/postgresqlflexibleservers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePostgresqlFlexibleServerDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePostgresqlFlexibleServerDatabaseCreate,
		Read:   resourcePostgresqlFlexibleServerDatabaseRead,
		Delete: resourcePostgresqlFlexibleServerDatabaseDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FlexibleServerDatabaseID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FlexibleServerDatabaseName,
			},

			"server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FlexibleServerID,
			},

			"charset": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.DatabaseCharset,
				Default:          "UTF8",
			},

			"collation": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatabaseCollation,
				Default:      "en_US.utf8",
			},
		},
	}
}

func resourcePostgresqlFlexibleServerDatabaseCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Postgres.FlexibleServerDatabaseClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	serverId, err := parse.FlexibleServerID(d.Get("server_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFlexibleServerDatabaseID(subscriptionId, serverId.ResourceGroup, serverId.Name, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %q: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_postgresql_flexible_server_database", id.ID())
		}
	}

	properties := postgresqlflexibleservers.Database{
		DatabaseProperties: &postgresqlflexibleservers.DatabaseProperties{
			Charset:   utils.String(d.Get("charset").(string)),
			Collation: utils.String(d.Get("collation").(string)),
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.FlexibleServerName, id.DatabaseName, properties)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourcePostgresqlFlexibleServerDatabaseRead(d, meta)
}

func resourcePostgresqlFlexibleServerDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Postgres.FlexibleServerDatabaseClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FlexibleServerDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FlexibleServerName, id.DatabaseName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.DatabaseName)
	d.Set("server_id", parse.NewFlexibleServerID(subscriptionId, id.ResourceGroup, id.FlexibleServerName).ID())

	if props := resp.DatabaseProperties; props != nil {
		d.Set("charset", props.Charset)
		d.Set("collation", props.Collation)
	}

	return nil
}

func resourcePostgresqlFlexibleServerDatabaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.FlexibleServerDatabaseClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FlexibleServerDatabaseID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.FlexibleServerName, id.DatabaseName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deleting of %s: %+v", *id, err)
	}

	return nil
}
