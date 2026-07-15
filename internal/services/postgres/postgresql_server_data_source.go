// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package postgres

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourcePostgreSqlServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourcePostgreSqlServerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(10 * time.Minute),
		},

		DeprecationMessage: "The `azurerm_postgresql_server` data source is deprecated and will be removed in v5.0 of the AzureRM Provider. Azure Database for PostgreSQL Single Server and its sub resources have been retired as of 2025-03-28, please use the `azurerm_postgresql_flexible_server` data source instead. For more information, see https://techcommunity.microsoft.com/blog/adforpostgresql/retiring-azure-database-for-postgresql-single-server-in-2025/3783783.",

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"administrator_login": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedIdentityComputed(),

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourcePostgreSqlServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := servers.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if err := d.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("fqdn", props.FullyQualifiedDomainName)
			d.Set("version", string(pointer.From(props.Version)))
			d.Set("administrator_login", props.AdministratorLogin)
		}

		if sku := model.Sku; sku != nil {
			d.Set("sku_name", sku.Name)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}
