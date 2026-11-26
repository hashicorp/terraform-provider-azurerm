// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_agreement

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountagreements"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogicAppIntegrationAccountAgreementRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountAgreementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccountagreements.ParseAgreementID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AgreementName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if model := resp.Model; model != nil {
		props := model.Properties

		d.Set("agreement_type", props.AgreementType)
		d.Set("guest_partner_name", props.GuestPartner)
		d.Set("host_partner_name", props.HostPartner)

		content, err := json.Marshal(props.Content)
		if err != nil {
			return err
		}
		d.Set("content", string(content))

		if err := d.Set("guest_identity", flattenIntegrationAccountAgreementBusinessIdentity(props.GuestIdentity)); err != nil {
			return fmt.Errorf("setting `guest_identity`: %+v", err)
		}

		if err := d.Set("host_identity", flattenIntegrationAccountAgreementBusinessIdentity(props.HostIdentity)); err != nil {
			return fmt.Errorf("setting `host_identity`: %+v", err)
		}

		if props.Metadata != nil {
			d.Set("metadata", props.Metadata)
		}
	}

	return nil
}

func flattenIntegrationAccountAgreementBusinessIdentity(input integrationaccountagreements.BusinessIdentity) []any {
	return []any{
		map[string]any{
			"qualifier": input.Qualifier,
			"value":     input.Value,
		},
	}
}
