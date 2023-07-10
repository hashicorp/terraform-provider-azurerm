// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceMsSqlElasticpool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceMsSqlElasticpoolRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"server_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"location": commonschema.LocationComputed(),

			"max_size_bytes": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"max_size_gb": {
				Type:     pluginsdk.TypeFloat,
				Computed: true,
			},

			"per_db_min_capacity": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"per_db_max_capacity": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"sku": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"capacity": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"tier": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"family": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),

			"zone_redundant": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"license_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMsSqlElasticpoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ElasticPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewElasticPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("server_name", id.ServerName)

	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.ElasticPoolProperties; props != nil {
		d.Set("max_size_gb", float64(*props.MaxSizeBytes/int64(1073741824)))
		d.Set("max_size_bytes", props.MaxSizeBytes)

		d.Set("zone_redundant", props.ZoneRedundant)
		d.Set("license_type", props.LicenseType)

		if perDbSettings := props.PerDatabaseSettings; perDbSettings != nil {
			d.Set("per_db_min_capacity", perDbSettings.MinCapacity)
			d.Set("per_db_max_capacity", perDbSettings.MaxCapacity)
		}
	}

	if err := d.Set("sku", flattenMsSqlElasticPoolSku(resp.Sku)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
