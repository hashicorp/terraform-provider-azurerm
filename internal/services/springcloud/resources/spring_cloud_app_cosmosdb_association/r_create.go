// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app_cosmosdb_association

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudAppCosmosDBAssociationCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := springcloudClient.From(meta).BindingsClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	appId, err := appplatform_rm.ParseAppIDInsensitively(d.Get("spring_cloud_app_id").(string))
	if err != nil {
		return err
	}

	id := appplatform_rm.NewBindingID(appId.SubscriptionId, appId.ResourceGroupName, appId.SpringName, appId.AppName, d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.AppName, id.BindingName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_spring_cloud_app_cosmosdb_association", id.ID())
			}
		}
	}

	apiType := d.Get("api_type").(string)
	cassandraKeyspaceName := d.Get("cosmosdb_cassandra_keyspace_name")
	gremlinDatabaseName := d.Get("cosmosdb_gremlin_database_name")
	gremlinGraphName := d.Get("cosmosdb_gremlin_graph_name")
	mongoDatabaseName := d.Get("cosmosdb_mongo_database_name")
	sqlDatabaseName := d.Get("cosmosdb_sql_database_name")

	bindingParameters := map[string]*string{
		springCloudAppCosmosDbAssociationKeyAPIType: pointer.To(apiType),
	}

	switch apiType {
	case springCloudAppCosmosDbAssociationAPITypeCassandra:
		if cassandraKeyspaceName == "" {
			return fmt.Errorf("`cosmosdb_cassandra_keyspace_name` should be set if `api_type` is `%s`", apiType)
		}
		bindingParameters[springCloudAppCosmosDbAssociationKeyKeySpace] = pointer.To(cassandraKeyspaceName.(string))
	case springCloudAppCosmosDbAssociationAPITypeGremlin:
		if gremlinDatabaseName == "" || gremlinGraphName == "" {
			return fmt.Errorf("`cosmosdb_gremlin_database_name` and `cosmosdb_gremlin_graph_name` should be set if `api_type` is `%s`", apiType)
		}
		bindingParameters[springCloudAppCosmosDbAssociationKeyDatabaseName] = pointer.To(gremlinDatabaseName.(string))
		bindingParameters[springCloudAppCosmosDbAssociationKeyCollectionName] = pointer.To(gremlinGraphName.(string))
	case springCloudAppCosmosDbAssociationAPITypeMongo:
		if mongoDatabaseName == "" {
			return fmt.Errorf("`cosmosdb_mongo_database_name` should be set if `api_type` is `%s`", apiType)
		}
		bindingParameters[springCloudAppCosmosDbAssociationKeyDatabaseName] = pointer.To(mongoDatabaseName.(string))
	case springCloudAppCosmosDbAssociationAPITypeSql:
		if sqlDatabaseName == "" {
			return fmt.Errorf("`cosmosdb_sql_database_name` should be set if `api_type` is `%s`", apiType)
		}
		bindingParameters[springCloudAppCosmosDbAssociationKeyDatabaseName] = pointer.To(sqlDatabaseName.(string))
	case springCloudAppCosmosDbAssociationAPITypeTable:
		if cassandraKeyspaceName != "" || gremlinDatabaseName != "" || gremlinGraphName != "" || mongoDatabaseName != "" || sqlDatabaseName != "" {
			return fmt.Errorf("`cosmosdb_cassandra_keyspace_name`, `cosmosdb_gremlin_database_name`, `cosmosdb_gremlin_graph_name`, `cosmosdb_mongo_database_name`, `cosmosdb_sql_database_name` should not be set if `api_type` is `%s`", apiType)
		}
	}

	bindingResource := appplatform.BindingResource{
		Properties: &appplatform.BindingResourceProperties{
			BindingParameters: bindingParameters,
			Key:               pointer.To(d.Get("cosmosdb_access_key").(string)),
			ResourceID:        pointer.To(d.Get("cosmosdb_account_id").(string)),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, id.AppName, id.BindingName, bindingResource)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %q: %+v", id, err)
	}

	return resourceSpringCloudAppCosmosDBAssociationRead(d, meta)
}
