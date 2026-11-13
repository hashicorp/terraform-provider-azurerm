// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_virtual_machine

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sqlvirtualmachine/2023-10-01/sqlvirtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlVirtualMachineRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.VirtualMachinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := sqlvirtualmachines.ParseSqlVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, sqlvirtualmachines.GetOperationOptions{Expand: pointer.To("*")})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Sql Virtual Machine %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}
	return resourceMssqlVirtualMachineSetFlatten(d, id, resp.Model)
}

func resourceMsSqlVirtualMachineCustomDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) error {
	// ForceNew when removing the auto_backup block.
	// See https://github.com/Azure/azure-rest-api-specs/issues/12818#issuecomment-773727756
	old, new := d.GetChange("auto_backup")
	if len(old.([]any)) == 1 && len(new.([]any)) == 0 {
		return d.ForceNew("auto_backup")
	}

	return nil
}

func mssqlVMCredentialNameDiffSuppressFunc(_, old, new string, _ *pluginsdk.ResourceData) bool {
	oldNamelist := strings.SplitSeq(old, ",")
	for n := range oldNamelist {
		cur := strings.Split(n, ":")
		if len(cur) > 1 && cur[1] == new {
			return true
		}
	}
	return false
}
