// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package attestation_provider

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/attestation/2020-10-01/attestationproviders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	attestationClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAttestationProviderDelete(d *pluginsdk.ResourceData, meta any) error {
	client := attestationClient.From(meta).ProviderClient
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := attestationproviders.ParseAttestationProvidersID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return nil
}
