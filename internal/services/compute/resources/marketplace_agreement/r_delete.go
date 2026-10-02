// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package marketplace_agreement

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/marketplaceordering/2015-06-01/agreements"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMarketplaceAgreementDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.MarketplaceAgreementsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := agreements.ParsePlanID(d.Id())
	if err != nil {
		return err
	}

	if _, err = client.MarketplaceAgreementsCancel(ctx, *id); err != nil {
		return fmt.Errorf("cancelling agreement for %s: %s", *id, err)
	}

	return nil
}
