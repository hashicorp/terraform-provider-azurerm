// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package confidential_ledger

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceConfidentialLedgerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ConfidentialLedger.ConfidentialLedgerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := confidentialledger.ParseLedgerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.LedgerGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.LedgerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("azuread_based_service_principal", flattenAADBasedSecurityPrincipal(props.AadBasedSecurityPrincipals)); err != nil {
				return fmt.Errorf("setting `azuread_based_service_principal`: %+v", err)
			}
			if err := d.Set("certificate_based_security_principal", flattenCertBasedSecurityPrincipal(props.CertBasedSecurityPrincipals)); err != nil {
				return fmt.Errorf("setting `certificate_based_security_principal`: %+v", err)
			}

			ledgerType := ""
			if props.LedgerType != nil {
				ledgerType = string(*props.LedgerType)
			}
			d.Set("ledger_type", ledgerType)

			d.Set("ledger_endpoint", props.LedgerUri)
			d.Set("identity_service_endpoint", props.IdentityServiceUri)
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func flattenAADBasedSecurityPrincipal(input *[]confidentialledger.AADBasedSecurityPrincipal) []any {
	output := make([]any, 0)
	if input == nil {
		return output
	}

	for _, item := range *input {
		ledgerRoleName := ""
		if item.LedgerRoleName != nil {
			ledgerRoleName = string(*item.LedgerRoleName)
		}

		output = append(output, map[string]any{
			"ledger_role_name": ledgerRoleName,
			"principal_id":     pointer.From(item.PrincipalId),
			"tenant_id":        pointer.From(item.TenantId),
		})
	}

	return output
}

func flattenCertBasedSecurityPrincipal(input *[]confidentialledger.CertBasedSecurityPrincipal) []any {
	output := make([]any, 0)
	if input == nil {
		return []any{}
	}

	for _, item := range *input {
		ledgerRoleName := ""
		if item.LedgerRoleName != nil {
			ledgerRoleName = string(*item.LedgerRoleName)
		}

		output = append(output, map[string]any{
			"ledger_role_name": ledgerRoleName,
			"pem_public_key":   pointer.From(item.Cert),
		})
	}

	return output
}
