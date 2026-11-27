// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_database

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2026-10-01/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	postgresClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePostgresqlFlexibleServerDatabaseRead(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := account.From(meta).SubscriptionId
	client := postgresClient.From(meta).FlexibleServerDatabaseClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := databases.ParseDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.DatabaseName)
	d.Set("server_id", databases.NewFlexibleServerID(subscriptionId, id.ResourceGroupName, id.FlexibleServerName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("charset", props.Charset)
			d.Set("collation", props.Collation)
		}
	}

	return nil
}
