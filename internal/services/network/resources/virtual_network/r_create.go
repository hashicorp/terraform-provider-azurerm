// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewVirtualNetworkID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, virtualnetworks.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_network", id.ID())
		}
	}

	vnetProperties, routeTables, err := expandVirtualNetworkProperties(ctx, *client, id, d)
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	locks.MultipleByID(routeTables)
	defer locks.UnlockMultipleByID(routeTables)

	vnet := virtualnetworks.VirtualNetwork{
		Name:             pointer.To(id.VirtualNetworkName),
		ExtendedLocation: helpers.ExpandEdgeZoneModel(d.Get("edge_zone").(string)),
		Location:         pointer.To(location.Normalize(d.Get("location").(string))),
		Properties:       vnetProperties,
		Tags:             tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("flow_timeout_in_minutes"); ok {
		vnet.Properties.FlowTimeoutInMinutes = pointer.To(int64(v.(int)))
	}

	networkSecurityGroupIds := make([]string, 0)
	for _, subnet := range *vnet.Properties.Subnets {
		if subnet.Properties != nil && subnet.Properties.NetworkSecurityGroup != nil {
			parsedNsgID, err := networksecuritygroups.ParseNetworkSecurityGroupID(*subnet.Properties.NetworkSecurityGroup.Id)
			if err != nil {
				return err
			}
			networkSecurityGroupIds = append(networkSecurityGroupIds, parsedNsgID.ID())
		}
	}

	locks.MultipleByID(&networkSecurityGroupIds)
	defer locks.UnlockMultipleByID(&networkSecurityGroupIds)

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, vnet, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(virtualnetworks.ProvisioningStateUpdating)},
		Target:     []string{string(virtualnetworks.ProvisioningStateSucceeded)},
		Refresh:    helpers.VirtualNetworkProvisioningStateRefreshFunc(ctx, client, id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of %s: %+v", id, err)
	}

	return resourceVirtualNetworkRead(d, meta)
}
