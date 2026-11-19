// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app_cosmosdb_association

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSpringCloudAppCosmosDBAssociationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).AppPlatform.BindingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	id, err := appplatform_rm.ParseBindingIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.AppName, id.BindingName)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("[INFO] Spring Cloud App Association %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.Set("name", id.BindingName)
	d.Set("spring_cloud_app_id", appplatform_rm.NewAppID(id.SubscriptionId, id.ResourceGroupName, id.SpringName, id.AppName).ID())
	if props := resp.Properties; props != nil {
		d.Set("cosmosdb_account_id", props.ResourceID)

		apiType := ""
		if v, ok := props.BindingParameters[springCloudAppCosmosDbAssociationKeyAPIType]; ok && v != nil {
			apiType = *v
		}
		d.Set("api_type", apiType)

		cassandraKeyspaceName := ""
		if v, ok := props.BindingParameters[springCloudAppCosmosDbAssociationKeyKeySpace]; ok && v != nil {
			cassandraKeyspaceName = *v
		}
		d.Set("cosmosdb_cassandra_keyspace_name", cassandraKeyspaceName)

		mongoDatabaseName := ""
		sqlDatabaseName := ""
		gremlinDatabaseName := ""
		if v, ok := props.BindingParameters[springCloudAppCosmosDbAssociationKeyDatabaseName]; ok && v != nil {
			switch apiType {
			case springCloudAppCosmosDbAssociationAPITypeMongo:
				mongoDatabaseName = *v
			case springCloudAppCosmosDbAssociationAPITypeSql:
				sqlDatabaseName = *v
			case springCloudAppCosmosDbAssociationAPITypeGremlin:
				gremlinDatabaseName = *v
			}
		}
		d.Set("cosmosdb_gremlin_database_name", gremlinDatabaseName)
		d.Set("cosmosdb_mongo_database_name", mongoDatabaseName)
		d.Set("cosmosdb_sql_database_name", sqlDatabaseName)

		gremlinGraphName := ""
		if v, ok := props.BindingParameters[springCloudAppCosmosDbAssociationKeyCollectionName]; ok && v != nil {
			gremlinGraphName = *v
		}
		d.Set("cosmosdb_gremlin_graph_name", gremlinGraphName)
	}
	return nil
}

const (
	springCloudAppCosmosDbAssociationKeyAPIType        = "apiType"
	springCloudAppCosmosDbAssociationKeyCollectionName = "collectionName"
	springCloudAppCosmosDbAssociationKeyDatabaseName   = "databaseName"
	springCloudAppCosmosDbAssociationKeyKeySpace       = "keySpace"

	springCloudAppCosmosDbAssociationAPITypeCassandra = "cassandra"
	springCloudAppCosmosDbAssociationAPITypeGremlin   = "gremlin"
	springCloudAppCosmosDbAssociationAPITypeMongo     = "mongo"
	springCloudAppCosmosDbAssociationAPITypeSql       = "sql"
	springCloudAppCosmosDbAssociationAPITypeTable     = "table"
)
