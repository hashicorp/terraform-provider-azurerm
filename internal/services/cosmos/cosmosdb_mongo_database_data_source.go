package cosmos

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceCosmosDbMongoDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceCosmosDbMongoDatabaseRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceCosmosDbMongoDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.MongoDbClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	account := d.Get("account_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.GetMongoDBDatabase(ctx, resourceGroup, account, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Cosmos DB Mongo Database %q (Account Name %q) was not found", name, account)
		}
		return fmt.Errorf("making Read request on AzureRM Cosmos DB Mongo Database %s (Account Name %q): %+v", name, account, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", resp.Name)
	d.Set("account_name", account)
	d.Set("resource_group_name", resourceGroup)

	return tags.FlattenAndSet(d, resp.Tags)
}
