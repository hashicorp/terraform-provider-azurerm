// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_application_gateway_backend_address_pool_association

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appgateway/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceIPConfigurationId{}, &parse.ApplicationGatewayBackendAddressPoolId{})
	if err != nil {
		return err
	}

	networkInterfaceId := commonids.NewNetworkInterfaceID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.NetworkInterfaceName)

	locks.ByID(strings.ToLower(networkInterfaceId.ID()))
	defer locks.UnlockByID(strings.ToLower(networkInterfaceId.ID()))

	resp, err := client.Get(ctx, networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", networkInterfaceId)
		}

		return fmt.Errorf("retrieving %s : %+v", networkInterfaceId, err)
	}

	if resp.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if resp.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}
	if resp.Model.Properties.IPConfigurations == nil {
		return fmt.Errorf("retrieving %s: `properties.IPConfigurations` was nil", id)
	}

	config := helpers.FindNetworkInterfaceIPConfiguration(resp.Model.Properties.IPConfigurations, id.First.IpConfigurationName)
	if config == nil {
		return fmt.Errorf("IP Configuration %q was not found on %s", id.First.IpConfigurationName, id)
	}
	if config.Properties == nil {
		return fmt.Errorf("`IPConfiguration.properties` was nil for %s", id)
	}
	props := resp.Model.Properties

	ipConfigProps := config.Properties

	backendAddressPools := make([]networkinterfaces.ApplicationGatewayBackendAddressPool, 0)
	if backendPools := ipConfigProps.ApplicationGatewayBackendAddressPools; backendPools != nil {
		for _, pool := range *backendPools {
			if pool.Id == nil {
				continue
			}

			if *pool.Id != id.Second.ID() {
				backendAddressPools = append(backendAddressPools, pool)
			}
		}
	}
	ipConfigProps.ApplicationGatewayBackendAddressPools = &backendAddressPools
	props.IPConfigurations = helpers.UpdateNetworkInterfaceIPConfiguration(*config, props.IPConfigurations)

	if err := client.CreateOrUpdateThenPoll(ctx, networkInterfaceId, *resp.Model); err != nil {
		return fmt.Errorf("removing %s Association for %s: %+v", id.Second, id.First, err)
	}

	return nil
}
