// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package confidential_ledger

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/confidentialledger/2022-05-13/confidentialledger"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceConfidentialLedgerCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ConfidentialLedger.ConfidentialLedgerClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := confidentialledger.NewLedgerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.LedgerGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_confidential_ledger", id.ID())
		}
	}

	aadBasedUsers := expandAADBasedSecurityPrincipal(d.Get("azuread_based_service_principal").([]any))
	certBasedUsers := expandCertBasedSecurityPrincipal(d.Get("certificate_based_security_principal").([]any))
	ledgerType := confidentialledger.LedgerType(d.Get("ledger_type").(string))
	location := location.Normalize(d.Get("location").(string))
	parameters := confidentialledger.ConfidentialLedger{
		Location: pointer.To(location),
		Properties: &confidentialledger.LedgerProperties{
			AadBasedSecurityPrincipals:  aadBasedUsers,
			CertBasedSecurityPrincipals: certBasedUsers,
			LedgerType:                  &ledgerType,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.LedgerCreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id.ID(), err)
	}

	d.SetId(id.ID())
	return resourceConfidentialLedgerRead(d, meta)
}

func expandAADBasedSecurityPrincipal(input []any) *[]confidentialledger.AADBasedSecurityPrincipal {
	output := make([]confidentialledger.AADBasedSecurityPrincipal, 0)

	for _, item := range input {
		v := item.(map[string]any)
		principalId := v["principal_id"].(string)
		tenantId := v["tenant_id"].(string)

		result := confidentialledger.AADBasedSecurityPrincipal{
			LedgerRoleName: pointer.ToEnum[confidentialledger.LedgerRoleName](v["ledger_role_name"].(string)),
			PrincipalId:    pointer.To(principalId),
			TenantId:       pointer.To(tenantId),
		}

		output = append(output, result)
	}

	return &output
}

func expandCertBasedSecurityPrincipal(input []any) *[]confidentialledger.CertBasedSecurityPrincipal {
	output := make([]confidentialledger.CertBasedSecurityPrincipal, 0)

	for _, item := range input {
		v := item.(map[string]any)

		output = append(output, confidentialledger.CertBasedSecurityPrincipal{
			Cert:           pointer.To(v["pem_public_key"].(string)),
			LedgerRoleName: pointer.ToEnum[confidentialledger.LedgerRoleName](v["ledger_role_name"].(string)),
		})
	}

	return &output
}
