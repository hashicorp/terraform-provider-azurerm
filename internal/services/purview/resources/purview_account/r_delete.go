// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package purview_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-12-01/account"
	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	purviewClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/purview/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePurviewAccountDelete(d *pluginsdk.ResourceData, meta any) error {
	client := purviewClient.From(meta).AccountsClient
	ctx, cancel := timeouts.ForDelete(clientsAccount.StopContext(meta), d)
	defer cancel()

	id, err := account.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
