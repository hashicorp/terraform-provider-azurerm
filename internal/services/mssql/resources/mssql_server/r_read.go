// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server

import (
	"context"
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

func msSqlMinimumTLSVersionDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) (err error) {
	old, new := d.GetChange("minimum_tls_version")
	// todo remove `old != "None"` when https://github.com/Azure/azure-rest-api-specs/issues/24348 is addressed
	if old != "" && old != "None" && old != "Disabled" && new == "Disabled" {
		err = fmt.Errorf("`minimum_tls_version` cannot be removed once set, please set a valid value for this property")
	}
	return
}

func msSqlPasswordChangeWhenAADAuthOnly(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) (err error) {
	old, _ := d.GetChange("azuread_administrator.0.azuread_authentication_only")
	if old.(bool) && d.HasChange("administrator_login_password") {
		err = fmt.Errorf("`administrator_login_password` cannot be changed once `azuread_administrator.0.azuread_authentication_only = true`")
	}
	return
}
