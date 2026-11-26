// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group_policy_exemption

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceResourceGroupPolicyExemptionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.ExemptionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ResourceGroupPolicyExemptionID(d.Id())
	if err != nil {
		return err
	}

	resourceGroupId := commonids.NewResourceGroupID(id.SubscriptionId, id.ResourceGroup)

	if _, err := client.Delete(ctx, resourceGroupId.ID(), id.PolicyExemptionName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id.ID(), err)
	}

	return nil
}
