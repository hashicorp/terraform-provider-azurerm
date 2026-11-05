// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_customer_managed_key

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCognitiveAccountCustomerManagedKeyDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cognitiveservicesaccounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	resp, err := client.AccountsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.Encryption == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	// set key source to Microsoft.CognitiveServices to disable customer managed key
	props := cognitiveservicesaccounts.Account{
		Properties: &cognitiveservicesaccounts.AccountProperties{
			Encryption: &cognitiveservicesaccounts.Encryption{
				KeySource: pointer.To(cognitiveservicesaccounts.KeySourceMicrosoftPointCognitiveServices),
			},
		},
	}

	if err := client.AccountsUpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("removing Customer Managed Key for %s: %+v", *id, err)
	}

	return nil
}
