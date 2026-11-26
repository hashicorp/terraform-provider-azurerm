// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMysqlFlexibleServerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Servers
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseFlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			if d.Get("rename_enabled").(bool) {
				return fmt.Errorf("retrieving %s: server not found with `rename_enabled` set; state has been retained because a rename may have changed its ID. Verify any rename operation and reconcile the confirmed server ID with Terraform state before retrying", *id)
			}
			log.Printf("[INFO] Mysql Flexible Server %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceMysqlFlexibleServerFlatten(d, id, resp.Model)
}

const (
	ServerMaintenanceWindowEnabled  = "Enabled"
	ServerMaintenanceWindowDisabled = "Disabled"
)
