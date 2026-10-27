// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription_policy_remediation

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubscriptionPolicyRemediationCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionId, err := commonids.ParseSubscriptionID(d.Get("subscription_id").(string))
	if err != nil {
		return err
	}

	id := remediations.NewRemediationID(subscriptionId.SubscriptionId, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.GetAtSubscription(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id.ID(), err)
				}
			}
			if existing.Model != nil {
				return tf.ImportAsExistsError("azurerm_subscription_policy_remediation", id.ID())
			}
		}
	}

	parameters := remediations.Remediation{
		Properties: helpers.ReadRemediationProperties(d),
	}

	if _, err = client.CreateOrUpdateAtSubscription(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id.ID(), err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceSubscriptionPolicyRemediationRead(d, meta)
}
