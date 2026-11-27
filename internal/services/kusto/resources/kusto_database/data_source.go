// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package kusto_database

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	kustoClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Kusto Database.",
		Read:        dataSourceKustoDatabaseRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Kusto Database.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DatabaseName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"cluster_name": {
				Description:  "The name of the Kusto Cluster this database is added to.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ClusterName,
			},

			"location": commonschema.LocationComputed(),

			"soft_delete_period": {
				Description: "The time the data should be kept before it stops being accessible to queries as ISO 8601 timespan.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"hot_cache_period": {
				Description: "The time the data that should be kept in cache for fast queries as ISO 8601 timespan.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"size": {
				Description: "The size of the database in bytes.",
				Type:        pluginsdk.TypeFloat,
				Computed:    true,
			},
		},
	}
}

func dataSourceKustoDatabaseRead(d *pluginsdk.ResourceData, meta any) error {
	client := kustoClient.From(meta).DatabasesClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id := commonids.NewKustoDatabaseID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s does not exist", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: response was nil", id)
	}

	database, ok := resp.Model.(databases.ReadWriteDatabase)
	if !ok {
		return fmt.Errorf("%s was not a Read/Write Database", id)
	}

	d.SetId(id.ID())

	d.Set("name", id.KustoDatabaseName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("cluster_name", id.KustoClusterName)
	d.Set("location", location.NormalizeNilable(database.Location))

	if props := database.Properties; props != nil {
		d.Set("hot_cache_period", props.HotCachePeriod)
		d.Set("soft_delete_period", props.SoftDeletePeriod)

		if statistics := props.Statistics; statistics != nil {
			d.Set("size", statistics.Size)
		}
	}

	return nil
}
