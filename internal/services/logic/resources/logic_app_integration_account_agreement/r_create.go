// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_agreement

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountagreements"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogicAppIntegrationAccountAgreementCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Logic.IntegrationAccountAgreementClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := integrationaccountagreements.NewAgreementID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_logic_app_integration_account_agreement", id.ID())
			}
		}
	}

	agreementContent := integrationaccountagreements.AgreementContent{}
	content := d.Get("content").(string)
	if err := json.Unmarshal([]byte(content), &agreementContent); err != nil {
		return fmt.Errorf("parsing JSON: %+v", err)
	}

	parameters := integrationaccountagreements.IntegrationAccountAgreement{
		Properties: integrationaccountagreements.IntegrationAccountAgreementProperties{
			AgreementType: integrationaccountagreements.AgreementType(d.Get("agreement_type").(string)),
			GuestIdentity: expandIntegrationAccountAgreementBusinessIdentity(d.Get("guest_identity").([]any)),
			GuestPartner:  d.Get("guest_partner_name").(string),
			HostIdentity:  expandIntegrationAccountAgreementBusinessIdentity(d.Get("host_identity").([]any)),
			HostPartner:   d.Get("host_partner_name").(string),
			Content:       agreementContent,
		},
	}

	if v, ok := d.GetOk("metadata"); ok {
		parameters.Properties.Metadata = &v
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountAgreementRead(d, meta)
}

func expandIntegrationAccountAgreementBusinessIdentity(input []any) integrationaccountagreements.BusinessIdentity {
	if len(input) == 0 {
		return integrationaccountagreements.BusinessIdentity{}
	}
	v := input[0].(map[string]any)

	return integrationaccountagreements.BusinessIdentity{
		Qualifier: v["qualifier"].(string),
		Value:     v["value"].(string),
	}
}
