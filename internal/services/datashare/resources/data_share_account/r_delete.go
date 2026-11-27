// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_share_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account"
	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	datashareClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataShareAccountDelete(d *pluginsdk.ResourceData, meta any) error {
	client := datashareClient.From(meta).AccountClient
	ctx, cancel := timeouts.ForDelete(clientsAccount.StopContext(meta), d)
	defer cancel()

	id, err := account.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
