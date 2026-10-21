// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_backend_address_pool_association

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceBackendAddressPoolAssociationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceIPConfigurationId{}, &loadbalancers.BackendAddressPoolId{})
	if err != nil {
		return err
	}

	networkInterfaceId := commonids.NewNetworkInterfaceID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.NetworkInterfaceName)

	locks.ByID(strings.ToLower(networkInterfaceId.ID()))
	defer locks.UnlockByID(strings.ToLower(networkInterfaceId.ID()))

	read, err := client.Get(ctx, networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			return fmt.Errorf("%s was not found", networkInterfaceId)
		}
		return fmt.Errorf("retrieving %s: %+v", networkInterfaceId, err)
	}

	if read.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", networkInterfaceId)
	}
	if read.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", networkInterfaceId)
	}
	if read.Model.Properties.IPConfigurations == nil {
		return fmt.Errorf("retrieving %s: `properties.ipConfigurations` was nil", networkInterfaceId)
	}

	props := read.Model.Properties
	config := helpers.FindNetworkInterfaceIPConfiguration(read.Model.Properties.IPConfigurations, id.First.IpConfigurationName)
	if config == nil {
		return fmt.Errorf("IP Configuration %q was not found on %s", id.First.IpConfigurationName, networkInterfaceId)
	}

	ipConfigProps := config.Properties
	if ipConfigProps == nil {
		return fmt.Errorf("retrieving %s: `properties` for IPConfiguration %q was nil", id.First.IpConfigurationName, networkInterfaceId)
	}

	backendAddressPools := make([]networkinterfaces.BackendAddressPool, 0)
	if backendPools := ipConfigProps.LoadBalancerBackendAddressPools; backendPools != nil {
		for _, pool := range *backendPools {
			if pool.Id == nil {
				continue
			}

			if *pool.Id != id.Second.ID() {
				backendAddressPools = append(backendAddressPools, pool)
			}
		}
	}
	ipConfigProps.LoadBalancerBackendAddressPools = &backendAddressPools
	props.IPConfigurations = helpers.UpdateNetworkInterfaceIPConfiguration(*config, props.IPConfigurations)

	if err := client.CreateOrUpdateThenPoll(ctx, networkInterfaceId, *read.Model); err != nil {
		return fmt.Errorf("removing Backend Address Pool Association for %s: %+v", networkInterfaceId, err)
	}

	return nil
}
