// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2026-10-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing PostgreSQL Flexible Server.",
		Read:        dataSourceArmPostgresqlFlexibleServerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description: "The name of this PostgreSQL Flexible Server.",
				Type:        pluginsdk.TypeString,
				Required:    true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"administrator_login": {
				Description: "The Administrator login for the PostgreSQL Flexible Server.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"auto_grow_enabled": {
				Description: "Whether the storage auto grow for PostgreSQL Flexible Server is enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"sku_name": {
				Description: "The SKU Name for the PostgreSQL Flexible Server. The name of the SKU, follows the `tier` + `name` pattern (e.g. `B_Standard_B1ms`, `GP_Standard_D2s_v3`, `MO_Standard_E4s_v3`).",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"storage_mb": {
				Description: "The max storage allowed for the PostgreSQL Flexible Server.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"version": {
				Description: "The version of PostgreSQL Flexible Server to use.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"delegated_subnet_id": {
				Description: "The ID of the Virtual Network subnet to create the PostgreSQL Flexible Server.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"backup_retention_days": {
				Description: "The backup retention days for the PostgreSQL Flexible Server.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"fqdn": {
				Description: "The FQDN of the PostgreSQL Flexible Server.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"public_network_access_enabled": {
				Description: "Whether public network access is enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"tags": commonschema.TagsDataSource(),

			"high_availability": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"mode": {
							Description: "The high availability mode of the PostgreSQL Flexible Server.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"standby_availability_zone": {
							Description: "The Availability Zone of the standby Flexible Server.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"zone": commonschema.ZoneSingleComputed(),
		},
	}
}

func dataSourceArmPostgresqlFlexibleServerRead(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := servers.NewFlexibleServerID(subscriptionId, resourceGroup, name)

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s does not exist", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.FlexibleServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))

		if props := model.Properties; props != nil {
			d.Set("administrator_login", props.AdministratorLogin)
			d.Set("version", pointer.FromEnum(props.Version))
			d.Set("fqdn", props.FullyQualifiedDomainName)
			d.Set("zone", props.AvailabilityZone)

			if storage := props.Storage; storage != nil {
				if storage.AutoGrow != nil {
					d.Set("auto_grow_enabled", *storage.AutoGrow == servers.StorageAutoGrowEnabled)
				}

				if storage.StorageSizeGB != nil {
					d.Set("storage_mb", (*props.Storage.StorageSizeGB * 1024))
				}
			}

			if backup := props.Backup; backup != nil {
				d.Set("backup_retention_days", props.Backup.BackupRetentionDays)
			}

			if network := props.Network; network != nil {
				d.Set("delegated_subnet_id", network.DelegatedSubnetResourceId)
				publicNetworkAccess := false
				if network.PublicNetworkAccess != nil {
					publicNetworkAccess = *network.PublicNetworkAccess == servers.ServerPublicNetworkAccessStateEnabled
				}
				d.Set("public_network_access_enabled", publicNetworkAccess)
			}

			if err := d.Set("high_availability", flattenFlexibleServerHighAvailability(props.HighAvailability)); err != nil {
				return fmt.Errorf("flattening `high_availability`: %+v", err)
			}
		}

		sku, err := flattenFlexibleServerSku(model.Sku)
		if err != nil {
			return fmt.Errorf("flattening `sku_name`: %+v", err)
		}

		d.Set("sku_name", sku)

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}
