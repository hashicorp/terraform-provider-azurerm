// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceSqlDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmSqlDatabaseRead,

		DeprecationMessage: "The `azurerm_sql_database` data source is deprecated and will be removed in version 4.0 of the AzureRM provider. Please use the `azurerm_mssql_database` data source instead.",

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"server_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"edition": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"collation": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"elastic_pool_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"default_secondary_location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"read_scale": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": tags.Schema(),

			"failover_group_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmSqlDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.DatabasesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDatabaseID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.DatabaseProperties; props != nil {
		d.Set("collation", props.Collation)
		d.Set("default_secondary_location", props.DefaultSecondaryLocation)
		d.Set("edition", string(props.Edition))
		d.Set("elastic_pool_name", props.ElasticPoolName)
		d.Set("failover_group_id", props.FailoverGroupID)
		d.Set("read_scale", props.ReadScale == sql.ReadScaleEnabled)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
