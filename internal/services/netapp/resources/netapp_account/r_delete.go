// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package netapp_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/netappaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	netappClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetAppAccountDelete(d *pluginsdk.ResourceData, meta any) error {
	client := netappClient.From(meta).AccountClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := netappaccounts.ParseNetAppAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if err := client.AccountsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
