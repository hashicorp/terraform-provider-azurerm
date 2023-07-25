// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2021-10-15/documentdb" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCosmosDbSQLStoredProcedure() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbSQLStoredProcedureCreate,
		Read:   resourceCosmosDbSQLStoredProcedureRead,
		Update: resourceCosmosDbSQLStoredProcedureUpdate,
		Delete: resourceCosmosDbSQLStoredProcedureDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SqlStoredProcedureID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"body": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"container_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},

			"database_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CosmosEntityName,
			},
		},
	}
}

func resourceCosmosDbSQLStoredProcedureCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.SqlClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	storedProcBody := d.Get("body").(string)
	id := parse.NewSqlStoredProcedureID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("database_name").(string), d.Get("container_name").(string), d.Get("name").(string))

	existing, err := client.GetSQLStoredProcedure(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.StoredProcedureName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of %s: %+v", id, err)
		}
	} else {
		if existing.ID == nil && *existing.ID == "" {
			return fmt.Errorf("generating import ID for %s", id)
		}

		return tf.ImportAsExistsError("azurerm_cosmosdb_sql_stored_procedure", *existing.ID)
	}

	storedProcParams := documentdb.SQLStoredProcedureCreateUpdateParameters{
		SQLStoredProcedureCreateUpdateProperties: &documentdb.SQLStoredProcedureCreateUpdateProperties{
			Resource: &documentdb.SQLStoredProcedureResource{
				ID:   &id.StoredProcedureName,
				Body: &storedProcBody,
			},
			Options: &documentdb.CreateUpdateOptions{},
		},
	}

	future, err := client.CreateUpdateSQLStoredProcedure(ctx, id.ResourceGroup, id.DatabaseAccountName, id.SqlDatabaseName, id.ContainerName, id.StoredProcedureName, storedProcParams)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbSQLStoredProcedureRead(d, meta)
}

func resourceCosmosDbSQLStoredProcedureUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("updating SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", name, containerName, databaseName, accountName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", name, containerName, databaseName, accountName, err)
	}

	return resourceCosmosDbSQLStoredProcedureRead(d, meta)
}

func resourceCosmosDbSQLStoredProcedureRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("retrieving SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", id.StoredProcedureName, id.ContainerName, id.SqlDatabaseName, id.DatabaseAccountName, err)
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

func resourceCosmosDbSQLStoredProcedureDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("deleting SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", id.StoredProcedureName, id.ContainerName, id.SqlDatabaseName, id.DatabaseAccountName, err)
		}
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("waiting for deletion of SQL Stored Procedure %q (Container %q / Database %q / Account %q): %+v", id.StoredProcedureName, id.ContainerName, id.SqlDatabaseName, id.DatabaseAccountName, err)
	}

	return nil
}
