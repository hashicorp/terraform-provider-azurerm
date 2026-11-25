// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package management_lock

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2020-05-01/managementlocks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagementLockRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Resource.LocksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managementlocks.ParseScopedLockID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetByScope(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.LockName)
	d.Set("scope", id.Scope)

	if model := resp.Model; model != nil {
		d.Set("lock_level", string(model.Properties.Level))
		d.Set("notes", model.Properties.Notes)
	}

	return nil
}
