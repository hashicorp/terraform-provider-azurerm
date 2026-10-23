// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface_nat_rule_association

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceNatRuleAssociationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseCompositeResourceID(d.Id(), &commonids.NetworkInterfaceIPConfigurationId{}, &loadbalancers.InboundNatRuleId{})
	if err != nil {
		return err
	}

	networkInterfaceId := commonids.NewNetworkInterfaceID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.NetworkInterfaceName)

	read, err := client.Get(ctx, networkInterfaceId, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("%s was not found - removing from state!", networkInterfaceId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", networkInterfaceId, err)
	}

	if model := read.Model; model != nil {
		if props := model.Properties; props != nil {
			ipConfigs := props.IPConfigurations
			if ipConfigs == nil {
				return fmt.Errorf("`properties.ipConfigurations` was nil for %s", networkInterfaceId)
			}

			config := helpers.FindNetworkInterfaceIPConfiguration(props.IPConfigurations, id.First.IpConfigurationName)
			if config == nil {
				log.Printf("IP Configuration %q was not found in %s - removing from state!", id.First.IpConfigurationName, networkInterfaceId)
				d.SetId("")
				return nil
			}

			found := false
			if ipConfigProps := config.Properties; ipConfigProps != nil {
				if rules := ipConfigProps.LoadBalancerInboundNatRules; rules != nil {
					for _, rule := range *rules {
						if rule.Id == nil {
							continue
						}

						if *rule.Id == id.Second.ID() {
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

	d.Set("ip_configuration_name", id.First.IpConfigurationName)
	d.Set("nat_rule_id", id.Second.ID())
	d.Set("network_interface_id", networkInterfaceId.ID())

	return nil
}
