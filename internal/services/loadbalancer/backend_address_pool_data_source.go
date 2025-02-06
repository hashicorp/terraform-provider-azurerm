// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceArmLoadBalancerBackendAddressPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmLoadBalancerBackendAddressPoolRead,
		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"loadbalancer_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: loadbalancers.ValidateLoadBalancerID,
			},

			"backend_address": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"virtual_network_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"ip_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"inbound_nat_rule_port_mapping": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"inbound_nat_rule_name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"frontend_port": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
									"backend_port": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"backend_ip_configurations": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"load_balancing_rules": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"outbound_rules": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"inbound_nat_rules": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceArmLoadBalancerBackendAddressPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	loadBalancerId, err := loadbalancers.ParseLoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}
	id := loadbalancers.NewLoadBalancerBackendAddressPoolID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroupName, loadBalancerId.LoadBalancerName, name)

	resp, err := lbClient.LoadBalancerBackendAddressPoolsGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("backend_address", flattenArmLoadBalancerBackendAddresses(props.LoadBalancerBackendAddresses)); err != nil {
				return fmt.Errorf("setting `backend_address`: %v", err)
			}

			var backendIPConfigurations []interface{}
			if beipConfigs := props.BackendIPConfigurations; beipConfigs != nil {
				for _, config := range *beipConfigs {
					ipConfig := make(map[string]interface{})
					if id := config.Id; id != nil {
						ipConfig["id"] = pointer.From(id)
						backendIPConfigurations = append(backendIPConfigurations, ipConfig)
					}
				}
			}
			if err := d.Set("backend_ip_configurations", backendIPConfigurations); err != nil {
				return fmt.Errorf("setting `backend_ip_configurations`: %v", err)
			}

			var loadBalancingRules []string
			if rules := props.LoadBalancingRules; rules != nil {
				for _, rule := range *rules {
					if rule.Id == nil {
						continue
					}
					loadBalancingRules = append(loadBalancingRules, pointer.From(rule.Id))
				}
			}
			if err := d.Set("load_balancing_rules", loadBalancingRules); err != nil {
				return fmt.Errorf("setting `load_balancing_rules`: %v", err)
			}

			var outboundRules []string
			if rules := props.OutboundRules; rules != nil {
				for _, rule := range *rules {
					if rule.Id == nil {
						continue
					}
					outboundRules = append(outboundRules, pointer.From(rule.Id))
				}
			}
			if err := d.Set("outbound_rules", outboundRules); err != nil {
				return fmt.Errorf("setting `outbound_rules`: %v", err)
			}

			var inboundNATRules []string
			if rules := props.InboundNatRules; rules != nil {
				for _, rule := range *rules {
					if rule.Id == nil {
						continue
					}
					inboundNATRules = append(inboundNATRules, pointer.From(rule.Id))
				}
			}
			if err := d.Set("inbound_nat_rules", inboundNATRules); err != nil {
				return fmt.Errorf("setting `inbound_nat_rules`: %v", err)
			}
		}
	}

	return nil
}

func flattenArmLoadBalancerBackendAddresses(input *[]loadbalancers.LoadBalancerBackendAddress) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		var name string
		if e.Name != nil {
			name = *e.Name
		}

		var (
			ipAddress string
			vnetId    string
		)
		var inboundNATRulePortMappingList []interface{}
		if prop := e.Properties; prop != nil {
			ipAddress = pointer.From(prop.IPAddress)

			if prop.VirtualNetwork != nil {
				vnetId = pointer.From(prop.VirtualNetwork.Id)
			}
			if prop.InboundNatRulesPortMapping != nil {
				rules := prop.InboundNatRulesPortMapping
				for _, rule := range *rules {
					rulePortMapping := make(map[string]interface{})

					if rule.InboundNatRuleName != nil {
						rulePortMapping["inbound_nat_rule_name"] = pointer.From(rule.InboundNatRuleName)
					}
					if rule.FrontendPort != nil {
						rulePortMapping["frontendPort"] = *rule.FrontendPort
					}

					if rule.BackendPort != nil {
						rulePortMapping["backendPort"] = *rule.BackendPort
					}
					inboundNATRulePortMappingList = append(inboundNATRulePortMappingList, rulePortMapping)
				}
			}
		}

		v := map[string]interface{}{
			"name":                          name,
			"virtual_network_id":            vnetId,
			"ip_address":                    ipAddress,
			"inbound_nat_rule_port_mapping": inboundNATRulePortMappingList,
		}
		output = append(output, v)
	}

	return output
}
