// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package purview_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-12-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePurviewAccountUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := account.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	parameters := account.AccountUpdateParameters{
		Properties: &account.AccountProperties{},
	}

	if d.HasChange("public_network_enabled") {
		if d.Get("public_network_enabled").(bool) {
			parameters.Properties.PublicNetworkAccess = pointer.To(account.PublicNetworkAccessEnabled)
		} else {
			parameters.Properties.PublicNetworkAccess = pointer.To(account.PublicNetworkAccessDisabled)
		}
	}

	if d.HasChange("managed_event_hub_enabled") {
		if d.Get("managed_event_hub_enabled").(bool) {
			parameters.Properties.ManagedEventHubState = pointer.To(account.ManagedEventHubStateEnabled)
		} else {
			parameters.Properties.ManagedEventHubState = pointer.To(account.ManagedEventHubStateDisabled)
		}
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("identity") {
		expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		parameters.Identity = expandedIdentity
	}

	if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourcePurviewAccountRead(d, meta)
}
