// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCosmosDbSQLStoredProcedure() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCosmosDbSQLStoredProcedureCreate,
		Read:   resourceCosmosDbSQLStoredProcedureRead,
		Update: resourceCosmosDbSQLStoredProcedureUpdate,
		Delete: resourceCosmosDbSQLStoredProcedureDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := cosmosdb.ParseStoredProcedureID(id)
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
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cosmosdb.NewStoredProcedureID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("database_name").(string), d.Get("container_name").(string), d.Get("name").(string))

	existing, err := client.SqlResourcesGetSqlStoredProcedure(ctx, id)
	if !response.WasNotFound(existing.HttpResponse) {
		if err != nil {
			return fmt.Errorf("checking for presence of %s: %+v", id, err)
		}
		return tf.ImportAsExistsError("azurerm_cosmosdb_sql_stored_procedure", id.ID())
	}

	storedProcParams := cosmosdb.SqlStoredProcedureCreateUpdateParameters{
		Properties: cosmosdb.SqlStoredProcedureCreateUpdateProperties{
			Resource: cosmosdb.SqlStoredProcedureResource{
				Id:   id.StoredProcedureName,
				Body: pointer.To(d.Get("body").(string)),
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if err := client.SqlResourcesCreateUpdateSqlStoredProcedureThenPoll(ctx, id, storedProcParams); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCosmosDbSQLStoredProcedureRead(d, meta)
}

func resourceCosmosDbSQLStoredProcedureUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseStoredProcedureID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.SqlResourcesGetSqlStoredProcedure(ctx, *id); err != nil {
		return fmt.Errorf("retrieving %s: %w", id, err)
	}

	storedProcParams := cosmosdb.SqlStoredProcedureCreateUpdateParameters{
		Properties: cosmosdb.SqlStoredProcedureCreateUpdateProperties{
			Resource: cosmosdb.SqlStoredProcedureResource{
				Id:   id.StoredProcedureName,
				Body: pointer.To(d.Get("body").(string)),
			},
			Options: &cosmosdb.CreateUpdateOptions{},
		},
	}

	if err := client.SqlResourcesCreateUpdateSqlStoredProcedureThenPoll(ctx, *id, storedProcParams); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceCosmosDbSQLStoredProcedureRead(d, meta)
}

func resourceCosmosDbSQLStoredProcedureRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseStoredProcedureID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SqlResourcesGetSqlStoredProcedure(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("database_name", id.SqlDatabaseName)
	d.Set("container_name", id.ContainerName)
	d.Set("name", id.StoredProcedureName)

	if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Resource != nil {
		d.Set("body", resp.Model.Properties.Resource.Body)
	}

	return nil
}

func resourceCosmosDbSQLStoredProcedureDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.CosmosDBClient

	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cosmosdb.ParseStoredProcedureID(d.Id())
	if err != nil {
		return err
	}

	if err := client.SqlResourcesDeleteSqlStoredProcedureThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
