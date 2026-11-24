// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseLoadBalancerID(d.Id())
	if err != nil {
		return err
	}

	// Casting because of a bug in provider ID generation due to a data issue that needs investigation.
	resp, err := client.Get(ctx, loadbalancers.ProviderLoadBalancerId(*id), loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.LoadBalancerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("edge_zone", flattenEdgeZone(model.ExtendedLocation))
		if sku := model.Sku; sku != nil {
			d.Set("sku", pointer.FromEnum(sku.Name))
			d.Set("sku_tier", pointer.FromEnum(sku.Tier))
		}

		if props := model.Properties; props != nil {
			if feipConfigs := props.FrontendIPConfigurations; feipConfigs != nil {
				if err := d.Set("frontend_ip_configuration", flattenLoadBalancerFrontendIpConfiguration(feipConfigs)); err != nil {
					return fmt.Errorf("flattening `frontend_ip_configuration`: %+v", err)
				}

				privateIpAddress := ""
				privateIpAddresses := make([]string, 0)
				for _, config := range *feipConfigs {
					if feipProps := config.Properties; feipProps != nil {
						if ip := feipProps.PrivateIPAddress; ip != nil {
							if privateIpAddress == "" {
								privateIpAddress = pointer.From(feipProps.PrivateIPAddress)
							}

							privateIpAddresses = append(privateIpAddresses, *feipProps.PrivateIPAddress)
						}
					}
				}

				d.Set("private_ip_address", privateIpAddress)
				d.Set("private_ip_addresses", privateIpAddresses)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

func flattenLoadBalancerFrontendIpConfiguration(ipConfigs *[]loadbalancers.FrontendIPConfiguration) []any {
	result := make([]any, 0)
	if ipConfigs == nil {
		return result
	}

	for _, config := range *ipConfigs {
		name := pointer.From(config.Name)

		id := pointer.From(config.Id)

		var inboundNatRules []any
		var loadBalancingRules []any
		var outboundRules []any
		gatewayLoadBalancerId := ""
		publicIpPrefixId := ""
		privateIPAllocationMethod := ""
		publicIpAddressId := ""
		privateIpAddressVersion := ""
		subnetId := ""
		privateIpAddress := ""

		if props := config.Properties; props != nil {
			privateIPAllocationMethod = pointer.FromEnum(props.PrivateIPAllocationMethod)

			if props.GatewayLoadBalancer != nil {
				gatewayLoadBalancerId = pointer.From(props.GatewayLoadBalancer.Id)
			}

			if subnet := props.Subnet; subnet != nil {
				subnetId = pointer.From(subnet.Id)
			}
			privateIpAddress = pointer.From(props.PrivateIPAddress)
			privateIpAddressVersion = pointer.FromEnum(props.PrivateIPAddressVersion)

			if pip := props.PublicIPAddress; pip != nil {
				publicIpAddressId = pointer.From(pip.Id)
			}

			if pip := props.PublicIPPrefix; pip != nil {
				publicIpPrefixId = pointer.From(pip.Id)
			}

			if rules := props.LoadBalancingRules; rules != nil {
				for _, rule := range *rules {
					if rule.Id == nil {
						continue
					}

					loadBalancingRules = append(loadBalancingRules, pointer.From(rule.Id))
				}
			}

			if rules := props.InboundNatRules; rules != nil {
				for _, rule := range *rules {
					inboundNatRules = append(inboundNatRules, pointer.From(rule.Id))
				}
			}

			if rules := props.OutboundRules; rules != nil {
				for _, rule := range *rules {
					outboundRules = append(outboundRules, pointer.From(rule.Id))
				}
			}
		}

		out := map[string]any{
			"gateway_load_balancer_frontend_ip_configuration_id": gatewayLoadBalancerId,
			"id":                            id,
			"inbound_nat_rules":             pluginsdk.NewSet(pluginsdk.HashString, inboundNatRules),
			"load_balancer_rules":           pluginsdk.NewSet(pluginsdk.HashString, loadBalancingRules),
			"name":                          name,
			"outbound_rules":                pluginsdk.NewSet(pluginsdk.HashString, outboundRules),
			"public_ip_address_id":          publicIpAddressId,
			"private_ip_address":            privateIpAddress,
			"private_ip_address_version":    privateIpAddressVersion,
			"private_ip_address_allocation": privateIPAllocationMethod,
			"public_ip_prefix_id":           publicIpPrefixId,
			"subnet_id":                     subnetId,
		}

		flattenedZones := zones.FlattenUntyped(config.Zones)

		if len(flattenedZones) > 0 {
			out["zones"] = flattenedZones
		}

		result = append(result, out)
	}
	return result
}

func flattenEdgeZone(input *edgezones.Model) string {
	if input == nil || input.Name == "" {
		return ""
	}
	return edgezones.NormalizeNilable(&input.Name)
}
