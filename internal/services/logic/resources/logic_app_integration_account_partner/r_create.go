// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_partner

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountpartners"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogicAppIntegrationAccountPartnerCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Logic.IntegrationAccountPartnerClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := integrationaccountpartners.NewPartnerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_logic_app_integration_account_partner", id.ID())
			}
		}
	}

	parameters := integrationaccountpartners.IntegrationAccountPartner{
		Properties: integrationaccountpartners.IntegrationAccountPartnerProperties{
			Content: integrationaccountpartners.PartnerContent{
				B2b: &integrationaccountpartners.B2BPartnerContent{
					BusinessIdentities: expandIntegrationAccountPartnerBusinessIdentity(d.Get("business_identity").(*pluginsdk.Set).List()),
				},
			},
			PartnerType: integrationaccountpartners.PartnerTypeBTwoB,
		},
	}

	if v, ok := d.GetOk("metadata"); ok {
		parameters.Properties.Metadata = &v
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountPartnerRead(d, meta)
}

func expandIntegrationAccountPartnerBusinessIdentity(input []any) *[]integrationaccountpartners.BusinessIdentity {
	results := make([]integrationaccountpartners.BusinessIdentity, 0)

	for _, item := range input {
		v := item.(map[string]any)

		results = append(results, integrationaccountpartners.BusinessIdentity{
			Qualifier: v["qualifier"].(string),
			Value:     v["value"].(string),
		})
	}

	return &results
}
