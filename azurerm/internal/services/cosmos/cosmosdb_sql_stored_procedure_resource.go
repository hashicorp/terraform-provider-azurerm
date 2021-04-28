package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-01-15/documentdb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/hashicorp/go-azure-helpers/response"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cosmos/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCosmosDbSQLStoredProcedure() *schema.Resource {
	return &schema.Resource{
		Create: resourceCosmosDbSQLStoredProcedureCreate,
		Read:   resourceCosmosDbSQLStoredProcedureRead,
		Update: resourceCosmosDbSQLStoredProcedureUpdate,
		Delete: resourceCosmosDbSQLStoredProcedureDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"body": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"container_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},
		},
	}
}

func resourceCosmosDbSQLStoredProcedureCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	containerName := d.Get("container_name").(string)
	databaseName := d.Get("database_name").(string)
	accountName := d.Get("account_name").(string)
	storedProcBody := d.Get("body").(string)

	existing, err := client.GetSQLStoredProcedure(ctx, resourceGroupName, accountName, databaseName, containerName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of creating SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", name, containerName, databaseName, accountName, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("Error generating import ID for Cosmos  SQL Stored Proecdure '%q' (Container %q / Database %q / Account %q)", name, containerName, databaseName, accountName)
		}

		return tf.ImportAsExistsError("azurerm_cosmosdb_sql_stored_procedure", *existing.ID)
	}

	storedProcParams := documentdb.SQLStoredProcedureCreateUpdateParameters{
		SQLStoredProcedureCreateUpdateProperties: &documentdb.SQLStoredProcedureCreateUpdateProperties{
			Resource: &documentdb.SQLStoredProcedureResource{
				ID:   &name,
				Body: &storedProcBody,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	future, err := client.CreateUpdateSQLStoredProcedure(ctx, resourceGroupName, accountName, databaseName, containerName, name, storedProcParams)
	if err != nil {
		return fmt.Errorf("Error creating SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", name, containerName, databaseName, accountName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", name, containerName, databaseName, accountName, err)
	}

	resp, err := client.GetSQLStoredProcedure(ctx, resourceGroupName, accountName, databaseName, containerName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", name, containerName, databaseName, accountName, err)
	}

	d.SetId(*resp.ID)

	return resourceCosmosDbSQLStoredProcedureRead(d, meta)
}

func resourceCosmosDbSQLStoredProcedureUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlStoredProcedureID(d.Id())
	if err != nil {
		return err
	}

	containerName := id.ContainerName
	databaseName := id.SqlDatabaseName
	accountName := id.DatabaseAccountName
	name := id.StoredProcedureName

	storedProcParams := documentdb.SQLStoredProcedureCreateUpdateParameters{
		SQLStoredProcedureCreateUpdateProperties: &documentdb.SQLStoredProcedureCreateUpdateProperties{
			Resource: &documentdb.SQLStoredProcedureResource{
				ID:   utils.String(name),
				Body: utils.String(d.Get("body").(string)),
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	future, err := client.CreateUpdateSQLStoredProcedure(ctx, id.ResourceGroup, accountName, databaseName, containerName, name, storedProcParams)
	if err != nil {
		return fmt.Errorf("Error updating SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", name, containerName, databaseName, accountName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", name, containerName, databaseName, accountName, err)
	}

	return resourceCosmosDbSQLStoredProcedureRead(d, meta)
}

func resourceCosmosDbSQLStoredProcedureRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlStoredProcedureID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetSQLStoredProcedure(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.StoredProcedureName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] SQL Stored Procedure %q (Container %q / Database %q / Account %q) was not found - removing from state", id.StoredProcedureName, id.ContainerName, id.SqlDatabaseName, id.DatabaseAccountName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", id.StoredProcedureName, id.ContainerName, id.SqlDatabaseName, id.DatabaseAccountName, err)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("database_name", id.SqlDatabaseName)
	d.Set("container_name", id.ContainerName)
	d.Set("name", id.StoredProcedureName)

	if props := resp.SQLStoredProcedureGetProperties; props != nil {
		if props.Resource != nil {
			d.Set("body", props.Resource.Body)
		}
	}

	return nil
}

func resourceCosmosDbSQLStoredProcedureDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SqlStoredProcedureID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.DeleteSQLStoredProcedure(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.StoredProcedureName)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", id.StoredProcedureName, id.ContainerName, id.SqlDatabaseName, id.DatabaseAccountName, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for deletion of SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", id.StoredProcedureName, id.ContainerName, id.SqlDatabaseName, id.DatabaseAccountName, err)
	}

	return nil
}
