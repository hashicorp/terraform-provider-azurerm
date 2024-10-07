// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2023-08-15/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	kustoValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceKustoDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceKustoDatabaseRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: kustoValidate.DatabaseName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"cluster_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: kustoValidate.ClusterName,
			},

			"location": commonschema.LocationComputed(),

			"soft_delete_period": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"hot_cache_period": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"size": {
				Type:     pluginsdk.TypeFloat,
				Computed: true,
			},
		},
	}
}

func dataSourceKustoDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.DatabasesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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
