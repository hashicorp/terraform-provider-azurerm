// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourcePostgresqlFlexibleServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmPostgresqlFlexibleServerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"administrator_login": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"auto_grow_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_mb": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"delegated_subnet_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"backup_retention_days": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceArmPostgresqlFlexibleServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			d.Set("version", string(pointer.From(props.Version)))
			d.Set("fqdn", props.FullyQualifiedDomainName)

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
