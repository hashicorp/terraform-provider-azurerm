// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_ddos_protection_plan

import (
	"fmt"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ddosprotectionplans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkDDoSProtectionPlanCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.DdosProtectionPlans
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vnetsToLock, err := expandNetworkDDoSProtectionPlanVnetIDs(d.Get("virtual_network_ids").([]any))
	if err != nil {
		return fmt.Errorf("extracting IDs of Virtual Network: %+v", err)
	}

	id := ddosprotectionplans.NewDdosProtectionPlanID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())
	locks.MultipleByID(vnetsToLock)
	defer locks.UnlockMultipleByID(vnetsToLock)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_network_ddos_protection_plan", id.ID())
		}
	}

	payload := ddosprotectionplans.DdosProtectionPlan{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceNetworkDDoSProtectionPlanRead(d, meta)
}

func expandNetworkDDoSProtectionPlanVnetIDs(input []any) (*[]string, error) {
	vnetIDs := make([]string, 0)

	for _, vnetID := range input {
		vnetResourceID, err := commonids.ParseVirtualNetworkID(vnetID.(string))
		if err != nil {
			return nil, err
		}

		if !slices.Contains(vnetIDs, vnetResourceID.ID()) {
			vnetIDs = append(vnetIDs, vnetResourceID.ID())
		}
	}

	return &vnetIDs, nil
}
