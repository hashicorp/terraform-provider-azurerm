// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription_policy_exemption

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubscriptionPolicyExemptionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.ExemptionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionPolicyExemptionID(d.Id())
	if err != nil {
		return err
	}

	subscriptionId := commonids.NewSubscriptionID(id.SubscriptionId)

	if _, err := client.Delete(ctx, subscriptionId.ID(), id.PolicyExemptionName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id.ID(), err)
	}

	return nil
}
