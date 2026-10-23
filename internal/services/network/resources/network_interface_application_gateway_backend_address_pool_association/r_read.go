// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_application_gateway_backend_address_pool_association

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appgateway/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceIPConfigurationId{}, &parse.ApplicationGatewayBackendAddressPoolId{})
	if err != nil {
		return err
	}

	networkInterfaceId := commonids.NewNetworkInterfaceID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.NetworkInterfaceName)

	resp, err := client.Get(ctx, networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("%s was not found - removing from state!", networkInterfaceId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", networkInterfaceId, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			config := helpers.FindNetworkInterfaceIPConfiguration(props.IPConfigurations, id.First.IpConfigurationName)
			if config == nil {
				log.Printf("%s was not found - removing from state!", id.First.ID())
				d.SetId("")
				return nil
			}

			found := false
			if ipConfigProps := config.Properties; ipConfigProps != nil {
				if backendPools := ipConfigProps.ApplicationGatewayBackendAddressPools; backendPools != nil {
					for _, pool := range *backendPools {
						if pool.Id == nil {
							continue
						}

						if *pool.Id == id.Second.ID() {
							found = true
							break
						}
					}
				}
			}

			if !found {
				log.Printf("[DEBUG] Association between %s and %s was not found - removing from state!", id.First, id.Second)
				d.SetId("")
				return nil
			}
		}
	}

	d.Set("backend_address_pool_id", id.Second.ID())
	d.Set("ip_configuration_name", id.First.IpConfigurationName)
	d.Set("network_interface_id", networkInterfaceId.ID())

	return nil
}
