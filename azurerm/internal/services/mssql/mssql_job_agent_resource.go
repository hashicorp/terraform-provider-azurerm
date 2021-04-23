package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMsSqlJobAgent() *schema.Resource {
	return &schema.Resource{
		Create: resourceMsSqlJobAgentCreateUpdate,
		Read:   resourceMsSqlJobAgentRead,
		Update: resourceMsSqlJobAgentCreateUpdate,
		Delete: resourceMsSqlJobAgentDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.JobAgentID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateMsSqlJobAgentName,
			},

			"database_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatabaseID,
			},

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),
		},
	}
}

func resourceMsSqlJobAgentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.JobAgentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Job Agent creation.")

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	databaseId := d.Get("database_id").(string)
	dbId, _ := parse.DatabaseID(databaseId)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, dbId.ResourceGroup, dbId.ServerName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Failed to check for presence of existing Job Agent %q (MsSql Server %q / Resource Group %q): %s", name, dbId.ServerName, dbId.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_mssql_job_agent", *existing.ID)
		}
	}

	params := sql.JobAgent{
		Name:     &name,
		Location: utils.String(location),
		JobAgentProperties: &sql.JobAgentProperties{
			DatabaseID: &databaseId,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, dbId.ResourceGroup, dbId.ServerName, name, params)
	if err != nil {
		return fmt.Errorf("creating MsSql Job Agent %q (Sql Server %q / Resource Group %q): %+v", name, dbId.ServerName, dbId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Job Agent %q (MsSql Server Name %q / Resource Group %q): %+v", name, dbId.ServerName, dbId.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, dbId.ResourceGroup, dbId.ServerName, name)
	if err != nil {
		return fmt.Errorf("reading request for Job Agent %q (MsSql Server Name %q / Resource Group %q): %+v", name, dbId.ServerName, dbId.ResourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceMsSqlJobAgentRead(d, meta)
}

func resourceMsSqlJobAgentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.JobAgentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.JobAgentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading MsSql Job Agent %s (MsSql Server Name %q / Resource Group %q): %s", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("database_id", resp.DatabaseID)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMsSqlJobAgentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.JobAgentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.JobAgentID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Job Agent %s: %+v", id.Name, err)
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}
