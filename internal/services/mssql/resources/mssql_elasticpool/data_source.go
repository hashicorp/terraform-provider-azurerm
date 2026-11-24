// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_elasticpool

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/elasticpools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing SQL elastic pool.",
		Read:        dataSourceMsSqlElasticpoolRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of the elastic pool.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"server_name": {
				Description: "The name of the SQL Server which contains the elastic pool.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"location": commonschema.LocationComputed(),

			"max_size_bytes": {
				Description: "The max data size of the elastic pool in bytes.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"max_size_gb": {
				Description: "The max data size of the elastic pool in gigabytes.",
				Type:        pluginsdk.TypeFloat,
				Computed:    true,
			},

			"per_db_min_capacity": {
				Description: "The minimum capacity all databases are guaranteed.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"per_db_max_capacity": {
				Description: "The maximum capacity any one database can consume.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"sku": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The SKU Name for this Elasticpool.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"capacity": {
							Description: "The scale up/out capacity, representing server's compute units.",
							Type:        pluginsdk.TypeInt,
							Computed:    true,
						},

						"tier": {
							Description: "The tier of the particular SKU.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"family": {
							Description: "The `family` of hardware.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"zone_redundant": {
				Description: "Whether or not this elastic pool is zone redundant.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"license_type": {
				Description: "The license type to apply for this elastic pool.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"enclave_type": {
				Description: "The type of enclave being used by the elastic pool.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"high_availability_replica_count": {
				Description: "The number of high availability replicas for the elastic pool.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceMsSqlElasticpoolRead(d *pluginsdk.ResourceData, meta any) error {
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
			d.Set("high_availability_replica_count", props.HighAvailabilityReplicaCount)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}
