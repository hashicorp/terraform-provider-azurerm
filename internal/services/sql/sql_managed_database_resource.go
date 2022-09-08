package sql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2018-06-01-preview/sql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmSqlManagedDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlManagedDatabaseCreateUpdate,
		Read:   resourceArmSqlManagedDatabaseRead,
		Delete: resourceArmSqlManagedDatabaseDelete,

		DeprecationMessage: "The `azurerm_sql_managed_database` resource is deprecated and will be removed in version 4.0 of the AzureRM provider. Please use the `azurerm_mssql_managed_database` resource instead.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedDatabaseID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(24 * time.Hour),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(24 * time.Hour),
			Delete: schema.DefaultTimeout(24 * time.Hour),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": azure.SchemaLocation(),

			"sql_managed_instance_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceArmSqlManagedDatabaseCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ManagedDatabasesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	managedInstanceID, err := parse.ManagedInstanceID(d.Get("sql_managed_instance_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewManagedDatabaseID(managedInstanceID.SubscriptionId, managedInstanceID.ResourceGroup, managedInstanceID.Name, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Managed Database %q: %s", id.ID(), err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_sql_managed_database", id.ID())
		}
	}

	database := sql.ManagedDatabase{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName, database)
	if err != nil {
		return fmt.Errorf("creating/updating SQL Managed Database %q: %+v", id.ID(), err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creating/updating to complete for SQL Managed Database %q: %+v", id.ID(), err)
	}

	d.SetId(id.ID())

	return resourceArmSqlManagedDatabaseRead(d, meta)
}

func resourceArmSqlManagedDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ManagedDatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading SQL Managed Database %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading SQL Managed Database %q: %v", id.ID(), err)
	}

	d.Set("sql_managed_instance_id", parse.NewManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName).ID())
	d.Set("name", resp.Name)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	return nil
}

func resourceArmSqlManagedDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.ManagedDatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedDatabaseID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ManagedInstanceName, id.DatabaseName)
	if err != nil {
		return fmt.Errorf("deleting SQL Managed Database %q: %+v", id.ID(), err)
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}
