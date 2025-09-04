// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cosmos

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.CosmosAccountName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceCosmosDbMongoDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cosmos.MongoDbClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewMongodbDatabaseID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	resp, err := client.GetMongoDBDatabase(ctx, id.ResourceGroup, id.DatabaseAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("account_name", id.DatabaseAccountName)
	d.Set("resource_group_name", id.ResourceGroup)

	return tags.FlattenAndSet(d, resp.Tags)
}
