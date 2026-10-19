// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package confidential_ledger

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceConfidentialLedgerUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ConfidentialLedger.ConfidentialLedgerClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := confidentialledger.ParseLedgerID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.LedgerGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving existing %s: model was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving existing %s: `properties` was nil", *id)
	}

	ledger := confidentialledger.ConfidentialLedger{
		Location: existing.Model.Location,
		Properties: &confidentialledger.LedgerProperties{
			AadBasedSecurityPrincipals:  existing.Model.Properties.AadBasedSecurityPrincipals,
			CertBasedSecurityPrincipals: existing.Model.Properties.CertBasedSecurityPrincipals,
			LedgerType:                  existing.Model.Properties.LedgerType,
		},
		Tags: existing.Model.Tags,
	}

	if d.HasChange("azuread_based_service_principal") {
		ledger.Properties.AadBasedSecurityPrincipals = expandAADBasedSecurityPrincipal(d.Get("azuread_based_service_principal").([]any))
	}

	if d.HasChange("certificate_based_security_principal") {
		ledger.Properties.CertBasedSecurityPrincipals = expandCertBasedSecurityPrincipal(d.Get("certificate_based_security_principal").([]any))
	}

	if d.HasChange("tags") {
		ledger.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err := client.LedgerUpdateThenPoll(ctx, *id, ledger); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceConfidentialLedgerRead(d, meta)
}
