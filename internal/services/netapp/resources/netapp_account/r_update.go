// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/netappaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppAccountUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.AccountClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := netappaccounts.ParseNetAppAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	existing, err := client.AccountsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	if ad := existing.Model.Properties.ActiveDirectories; ad != nil && len(*ad) > 0 {
		// The API doesn't return these, so we'll set these values based on config.
		// If there are no changes to the `active_directory` block this ensures we don't unintentionally wipe these values
		existingAD := (*ad)[0]
		existingAD.ServerRootCACertificate = pointer.To(d.Get("active_directory.0.server_root_ca_certificate").(string))
		existingAD.Password = pointer.To(d.Get("active_directory.0.password").(string))
		existing.Model.Properties.ActiveDirectories = pointer.To([]netappaccounts.ActiveDirectory{existingAD})
	}

	if d.HasChange("active_directory") {
		existing.Model.Properties.ActiveDirectories = expandNetAppActiveDirectories(d.Get("active_directory").([]any))
	}

	if d.HasChange("tags") {
		existing.Model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("identity") {
		anfAccountIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		existing.Model.Identity = anfAccountIdentity
	}

	if err := client.AccountsCreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id.ID(), err)
	}

	return resourceNetAppAccountRead(d, meta)
}
