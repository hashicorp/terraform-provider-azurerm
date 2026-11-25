// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlServerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.ServersClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSqlServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, servers.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Error retrieving SQL Server %s - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving SQL Server %s: %v", id, err)
	}
	return resourceMssqlServerSetFlatten(ctx, d, id, resp.Model, meta.(*clients.Client))
}
