// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/restorabledroppeddatabases" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/servers"                    // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceMsSqlServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceMsSqlServerRead,

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

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"administrator_login": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"fully_qualified_domain_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

			"restorable_dropped_database_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceMsSqlServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	restorableDroppedDatabasesClient := meta.(*clients.Client).MSSQL.RestorableDroppedDatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	serverId := servers.ServerId{
		SubscriptionId:    id.SubscriptionId,
		ResourceGroupName: id.ResourceGroup,
		ServerName:        id.Name,
	}

	resp, err := client.Get(ctx, serverId, servers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %s", id, err)
	}

	if model := resp.Model; model != nil {
		t := make(map[string]interface{}, 0)
		if model.Tags != nil {
			t = tags.ExpandFrom(model.Tags)
		}

		d.SetId(id.ID())
		d.Set("location", location.NormalizeNilable(pointer.To(model.Location)))

		if props := model.Properties; props != nil {
			d.Set("version", props.Version)
			d.Set("administrator_login", props.AdministratorLogin)
			d.Set("fully_qualified_domain_name", props.FullyQualifiedDomainName)
		}

		identity, err := flattenSqlServerIdentity(model.Identity)
		if err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if err := d.Set("identity", identity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		restorableListPage, err := restorableDroppedDatabasesClient.ListByServerComplete(ctx, restorabledroppeddatabases.ServerId(serverId))
		if err != nil {
			return fmt.Errorf("listing %s Restorable Dropped Databases: %v", id, err)
		}

		if err := d.Set("restorable_dropped_database_ids", flattenSqlServerRestorableDatabases(restorableListPage)); err != nil {
			return fmt.Errorf("setting `restorable_dropped_database_ids`: %+v", err)
		}

		return tags.FlattenAndSet(d, tags.Expand(t))
	}

	return fmt.Errorf("model was `nil`")
}
