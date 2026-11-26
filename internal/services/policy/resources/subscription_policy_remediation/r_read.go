// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription_policy_remediation

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/policyinsights/2021-10-01/remediations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubscriptionPolicyRemediationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.RemediationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := remediations.ParseRemediationID(d.Id())
	if err != nil {
		return fmt.Errorf("reading Policy Remediation: %+v", err)
	}

	subscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)

	resp, err := client.GetAtSubscription(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", id.ID())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id.ID(), err)
	}

	d.Set("name", id.RemediationName)
	d.Set("subscription_id", subscriptionId.ID())

	return helpers.SetRemediationProperties(d, resp.Model.Properties)
}

func subscriptionPolicyRemediationCancellationRefreshFunc(ctx context.Context, client *remediations.RemediationsClient, id remediations.RemediationId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		resp, err := client.GetAtSubscription(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("issuing read request for %s: %+v", id.ID(), err)
		}

		if resp.Model.Properties == nil {
			return nil, "", fmt.Errorf("`properties` was nil")
		}
		if resp.Model.Properties.ProvisioningState == nil {
			return nil, "", fmt.Errorf("`properties.ProvisioningState` was nil")
		}
		return resp, *resp.Model.Properties.ProvisioningState, nil
	}
}
