// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_backend_address_pool_address

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			lbClient := metadata.Client.LoadBalancers.LoadBalancersClient

			id, err := parse.BackendAddressPoolAddressID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			poolId := loadbalancers.NewLoadBalancerBackendAddressPoolID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			pool, err := lbClient.LoadBalancerBackendAddressPoolsGet(ctx, poolId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", poolId, err)
			}
			if pool.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", poolId)
			}
			if pool.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", poolId)
			}

			plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroup, LoadBalancerName: id.LoadBalancerName}
			lb, err := lbClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", plbId, err)
			}

			var backendAddress *loadbalancers.LoadBalancerBackendAddress
			if model := pool.Model; model != nil {
				if pool := model.Properties; pool != nil {
					for _, address := range *pool.LoadBalancerBackendAddresses {
						if address.Name == nil {
							continue
						}

						if *address.Name == id.AddressName {
							backendAddress = &address
							break
						}
					}
				}
			}

			if backendAddress == nil {
				return metadata.MarkAsGone(id)
			}

			model := BackendAddressPoolAddressModel{
				Name:                 id.AddressName,
				BackendAddressPoolId: poolId.ID(),
			}

			if props := backendAddress.Properties; props != nil {
				if lb.Model != nil && pointer.From(lb.Model.Sku.Tier) == loadbalancers.LoadBalancerSkuTierGlobal {
					if props.LoadBalancerFrontendIPConfiguration != nil && props.LoadBalancerFrontendIPConfiguration.Id != nil {
						model.FrontendIPConfiguration = *props.LoadBalancerFrontendIPConfiguration.Id
					}
				} else {
					if props.IPAddress != nil {
						model.IPAddress = *props.IPAddress
					}

					if props.VirtualNetwork != nil && props.VirtualNetwork.Id != nil {
						model.VirtualNetworkId = *props.VirtualNetwork.Id
					}
				}
				var inboundNATRulePortMappingList []inboundNATRulePortMapping
				if rules := props.InboundNatRulesPortMapping; rules != nil {
					for _, rule := range *rules {
						rulePortMapping := inboundNATRulePortMapping{}

						if rule.InboundNatRuleName != nil {
							rulePortMapping.Name = *rule.InboundNatRuleName
						}
						if rule.FrontendPort != nil {
							rulePortMapping.FrontendPort = *rule.FrontendPort
						}

						if rule.BackendPort != nil {
							rulePortMapping.BackendPort = *rule.BackendPort
						}
						inboundNATRulePortMappingList = append(inboundNATRulePortMappingList, rulePortMapping)
					}
					model.PortMapping = inboundNATRulePortMappingList
				}
			}

			return metadata.Encode(&model)
		},
		Timeout: 5 * time.Minute,
	}
}

type inboundNATRulePortMapping struct {
	Name         string `tfschema:"inbound_nat_rule_name"`
	FrontendPort int64  `tfschema:"frontend_port"`
	BackendPort  int64  `tfschema:"backend_port"`
}

func portMapping() *pluginsdk.Schema {
	return &pluginsdk.Schema{
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
	}
}

func loadbalacnerProvisioningStatusRefreshFunc(ctx context.Context, client *loadbalancers.LoadBalancersClient, id parse.BackendAddressPoolAddressId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroup, LoadBalancerName: id.LoadBalancerName}
		lbClient, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
		if err != nil {
			return nil, "", fmt.Errorf("retrieving load balancer error： %+v", err)
		}
		if model := lbClient.Model; model != nil {
			if props := model.Properties; props != nil {
				return lbClient, pointer.FromEnum(props.ProvisioningState), nil
			}
		}
		return lbClient, "", nil
	}
}
