// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package confidential_ledger

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	confidentialledgerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceConfidentialLedgerDelete(d *pluginsdk.ResourceData, meta any) error {
	client := confidentialledgerClient.From(meta).ConfidentialLedgerClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := confidentialledger.ParseLedgerID(d.Id())
	if err != nil {
		return err
	}

	if err := client.LedgerDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
