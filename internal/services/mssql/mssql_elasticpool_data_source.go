// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/elasticpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"zone_redundant": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"license_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"enclave_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceMsSqlElasticpoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ElasticPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewSqlElasticPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("server_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("name", id.ElasticPoolName)
		d.Set("resource_group_name", id.ResourceGroupName)
		d.Set("server_name", id.ServerName)
		d.Set("location", location.Normalize(model.Location))

		if err := d.Set("sku", flattenMsSqlElasticPoolSku(model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("max_size_gb", float64(*props.MaxSizeBytes/int64(1073741824)))
			d.Set("max_size_bytes", props.MaxSizeBytes)
			d.Set("zone_redundant", props.ZoneRedundant)

			licenseType := ""
			if props.LicenseType != nil {
				licenseType = string(*props.LicenseType)
			}
			d.Set("license_type", licenseType)

			if perDbSettings := props.PerDatabaseSettings; perDbSettings != nil {
				d.Set("per_db_min_capacity", perDbSettings.MinCapacity)
				d.Set("per_db_max_capacity", perDbSettings.MaxCapacity)
			}

			enclaveType := ""
			if props.PreferredEnclaveType != nil && *props.PreferredEnclaveType != elasticpools.AlwaysEncryptedEnclaveTypeDefault {
				enclaveType = string(elasticpools.AlwaysEncryptedEnclaveTypeVBS)
			}
			d.Set("enclave_type", enclaveType)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}
