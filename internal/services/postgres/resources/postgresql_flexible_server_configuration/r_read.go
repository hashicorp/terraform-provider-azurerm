// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_configuration

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2026-10-01/configurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2026-10-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	postgresClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFlexibleServerConfigurationRead(d *pluginsdk.ResourceData, meta any) error {
	client := postgresClient.From(meta).FlexibleServersConfigurationsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := configurations.ParseConfigurationID(d.Id())
	if err != nil {
		return err
	}

	// if the server has already been deleted, we need also mark it as gone
	serverID := servers.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName)
	if exists, _ := postgresClient.From(meta).FlexibleServersClient.Get(ctx, serverID); response.WasNotFound(exists.HttpResponse) {
		log.Printf("[WARN] server %s was not found, removing from state", serverID)
		d.SetId("")
		return nil
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found, removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for %s: %+v", id, err)
	}

	return resourcePostgresqlFlexibleServerConfigurationFlatten(d, id, resp.Model)
}
