// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

var (
	_ sdk.Resource           = BackendAddressPoolAddressResource{}
	_ sdk.ResourceWithUpdate = BackendAddressPoolAddressResource{}
)

type BackendAddressPoolAddressResource struct{}

type BackendAddressPoolAddressModel struct {
	Name                    string                      `tfschema:"name"`
	BackendAddressPoolId    string                      `tfschema:"backend_address_pool_id"`
	VirtualNetworkId        string                      `tfschema:"virtual_network_id"`
	IPAddress               string                      `tfschema:"ip_address"`
	FrontendIPConfiguration string                      `tfschema:"backend_address_ip_configuration_id"`
	PortMapping             []inboundNATRulePortMapping `tfschema:"inbound_nat_rule_port_mapping"`
}

type inboundNATRulePortMapping struct {
	Name         string `tfschema:"inbound_nat_rule_name"`
	FrontendPort int32  `tfschema:"frontend_port"`
	BackendPort  int32  `tfschema:"backend_port"`
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

func (r BackendAddressPoolAddressResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"backend_address_pool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LoadBalancerBackendAddressPoolID,
		},

		"virtual_network_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			RequiredWith:  []string{"ip_address"},
			ConflictsWith: []string{"backend_address_ip_configuration_id"},
			ValidateFunc:  commonids.ValidateVirtualNetworkID,
			Description:   "For regional load balancer, user needs to specify `virtual_network_id` and `ip_address`",
		},

		"ip_address": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			RequiredWith: []string{"virtual_network_id"},
			ValidateFunc: validation.IsIPAddress,
		},

		"backend_address_ip_configuration_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"virtual_network_id"},
			ValidateFunc:  validate.LoadBalancerFrontendIpConfigurationID,
			Description:   "For global load balancer, user needs to specify the `backend_address_ip_configuration_id` of the added regional load balancers",
		},
	}
}

func (r BackendAddressPoolAddressResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"inbound_nat_rule_port_mapping": portMapping(),
	}
}

func (r BackendAddressPoolAddressResource) ModelObject() interface{} {
	return &BackendAddressPoolAddressModel{}
}

func (r BackendAddressPoolAddressResource) ResourceType() string {
	return "azurerm_lb_backend_address_pool_address"
}

func (r BackendAddressPoolAddressResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadBalancers.LoadBalancerBackendAddressPoolsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model BackendAddressPoolAddressModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			poolId, err := parse.LoadBalancerBackendAddressPoolID(model.BackendAddressPoolId)
			if err != nil {
				return err
			}

			locks.ByName(poolId.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(poolId.BackendAddressPoolName, backendAddressPoolResourceName)

			// Backend Addresses can not be created for Basic sku, so we have to check
			lb, err := metadata.Client.LoadBalancers.LoadBalancersClient.Get(ctx, poolId.ResourceGroup, poolId.LoadBalancerName, "")
			if err != nil {
				return fmt.Errorf("retrieving Load Balancer %q (Resource Group %q): %+v", poolId.LoadBalancerName, poolId.ResourceGroup, err)
			}
			isBasicSku := true
			if lb.Sku != nil && lb.Sku.Name != network.LoadBalancerSkuNameBasic {
				isBasicSku = false
			}
			if isBasicSku {
				return fmt.Errorf("Backend Addresses are not supported on Basic SKU Load Balancers")
			}
			if lb.Sku != nil && lb.Sku.Tier == network.LoadBalancerSkuTierGlobal {
				if model.FrontendIPConfiguration == "" {
					return fmt.Errorf("Please set a Regional Backend Address Pool Addresses for the Global load balancer")
				}
			}

			id := parse.NewBackendAddressPoolAddressID(subscriptionId, poolId.ResourceGroup, poolId.LoadBalancerName, poolId.BackendAddressPoolName, model.Name)
			pool, err := client.Get(ctx, poolId.ResourceGroup, poolId.LoadBalancerName, poolId.BackendAddressPoolName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *poolId, err)
			}
			if pool.BackendAddressPoolPropertiesFormat == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *poolId)
			}

			addresses := make([]network.LoadBalancerBackendAddress, 0)
			if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
				addresses = *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses
			}

			metadata.Logger.Infof("checking for existing %s..", id)
			for _, address := range addresses {
				if address.Name == nil {
					continue
				}

				if *address.Name == id.AddressName {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			if lb.Sku.Tier == network.LoadBalancerSkuTierGlobal {
				addresses = append(addresses, network.LoadBalancerBackendAddress{
					Name: utils.String(id.AddressName),
					LoadBalancerBackendAddressPropertiesFormat: &network.LoadBalancerBackendAddressPropertiesFormat{
						LoadBalancerFrontendIPConfiguration: &network.SubResource{
							ID: utils.String(model.FrontendIPConfiguration),
						},
					},
				})
			} else {
				addresses = append(addresses, network.LoadBalancerBackendAddress{
					LoadBalancerBackendAddressPropertiesFormat: &network.LoadBalancerBackendAddressPropertiesFormat{
						IPAddress: utils.String(model.IPAddress),
						VirtualNetwork: &network.SubResource{
							ID: utils.String(model.VirtualNetworkId),
						},
					},
					Name: utils.String(id.AddressName),
				})
			}
			pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = &addresses

			metadata.Logger.Infof("adding %s..", id)
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			metadata.Logger.Infof("waiting for update %s..", id)
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
			}
			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r BackendAddressPoolAddressResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadBalancers.LoadBalancerBackendAddressPoolsClient
			lbClient := metadata.Client.LoadBalancers.LoadBalancersClient
			id, err := parse.BackendAddressPoolAddressID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if pool.BackendAddressPoolPropertiesFormat == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			lb, err := lbClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
			if err != nil {
				return fmt.Errorf("retrieving Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
			}

			var backendAddress *network.LoadBalancerBackendAddress
			if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
				for _, address := range *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses {
					if address.Name == nil {
						continue
					}

					if *address.Name == id.AddressName {
						backendAddress = &address
						break
					}
				}
			}
			if backendAddress == nil {
				return metadata.MarkAsGone(id)
			}

			backendAddressPoolId := parse.NewLoadBalancerBackendAddressPoolID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			model := BackendAddressPoolAddressModel{
				Name:                 id.AddressName,
				BackendAddressPoolId: backendAddressPoolId.ID(),
			}

			if props := backendAddress.LoadBalancerBackendAddressPropertiesFormat; props != nil {
				if lb.Sku.Tier == network.LoadBalancerSkuTierGlobal {
					if props.LoadBalancerFrontendIPConfiguration != nil && props.LoadBalancerFrontendIPConfiguration.ID != nil {
						model.FrontendIPConfiguration = *props.LoadBalancerFrontendIPConfiguration.ID
					}
				} else {
					if props.IPAddress != nil {
						model.IPAddress = *props.IPAddress
					}

					if props.VirtualNetwork != nil && props.VirtualNetwork.ID != nil {
						model.VirtualNetworkId = *props.VirtualNetwork.ID
					}
				}
			}

			var inboundNATRulePortMappingList []inboundNATRulePortMapping
			if rules := backendAddress.LoadBalancerBackendAddressPropertiesFormat.InboundNatRulesPortMapping; rules != nil {
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
			return metadata.Encode(&model)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r BackendAddressPoolAddressResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadBalancers.LoadBalancerBackendAddressPoolsClient
			lbClient := metadata.Client.LoadBalancers.LoadBalancersClient
			id, err := parse.BackendAddressPoolAddressID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

			pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if pool.BackendAddressPoolPropertiesFormat == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			timeout, _ := ctx.Deadline()
			lbStatus := &pluginsdk.StateChangeConf{
				Pending:                   []string{string(network.ProvisioningStateUpdating)},
				Target:                    []string{string(network.ProvisioningStateSucceeded)},
				MinTimeout:                5 * time.Second,
				Refresh:                   loadbalacnerProvisioningStatusRefreshFunc(ctx, lbClient, *id),
				ContinuousTargetOccurence: 10,
				Timeout:                   time.Until(timeout),
			}

			if _, err := lbStatus.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for parent resource loadbalancer status to be ready error: %+v", err)
			}

			addresses := make([]network.LoadBalancerBackendAddress, 0)
			if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
				addresses = *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses
			}

			newAddresses := make([]network.LoadBalancerBackendAddress, 0)
			for _, address := range addresses {
				if address.Name == nil {
					continue
				}

				if *address.Name != id.AddressName {
					newAddresses = append(newAddresses, address)
				}
			}
			pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = &newAddresses

			metadata.Logger.Infof("removing %s..", *id)
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
			if err != nil {
				return fmt.Errorf("removing %s: %+v", *id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for removal of %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r BackendAddressPoolAddressResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.BackendAddressPoolAddressID
}

func (r BackendAddressPoolAddressResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadBalancers.LoadBalancerBackendAddressPoolsClient
			lbClient := metadata.Client.LoadBalancers.LoadBalancersClient
			id, err := parse.BackendAddressPoolAddressID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

			var model BackendAddressPoolAddressModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			pool, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if pool.BackendAddressPoolPropertiesFormat == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			lb, err := lbClient.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
			if err != nil {
				return fmt.Errorf("retrieving Load Balancer %q (Resource Group %q): %+v", id.LoadBalancerName, id.ResourceGroup, err)
			}

			addresses := make([]network.LoadBalancerBackendAddress, 0)
			if pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses != nil {
				addresses = *pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses
			}
			index := -1
			for i, address := range addresses {
				if address.Name == nil {
					continue
				}

				if *address.Name == id.AddressName {
					index = i
					break
				}
			}
			if index == -1 {
				return fmt.Errorf("%s was not found", *id)
			}

			if lb.Sku.Tier == network.LoadBalancerSkuTierGlobal {
				addresses[index] = network.LoadBalancerBackendAddress{
					Name: utils.String(model.Name),
					LoadBalancerBackendAddressPropertiesFormat: &network.LoadBalancerBackendAddressPropertiesFormat{
						LoadBalancerFrontendIPConfiguration: &network.SubResource{
							ID: utils.String(model.FrontendIPConfiguration),
						},
					},
				}
			} else {
				addresses[index] = network.LoadBalancerBackendAddress{
					LoadBalancerBackendAddressPropertiesFormat: &network.LoadBalancerBackendAddressPropertiesFormat{
						IPAddress: utils.String(model.IPAddress),
						VirtualNetwork: &network.SubResource{
							ID: utils.String(model.VirtualNetworkId),
						},
					},
					Name: utils.String(id.AddressName),
				}
			}
			pool.BackendAddressPoolPropertiesFormat.LoadBalancerBackendAddresses = &addresses

			timeout, _ := ctx.Deadline()
			lbStatus := &pluginsdk.StateChangeConf{
				Pending:                   []string{string(network.ProvisioningStateUpdating)},
				Target:                    []string{string(network.ProvisioningStateSucceeded)},
				MinTimeout:                5 * time.Minute,
				PollInterval:              10 * time.Second,
				Refresh:                   loadbalacnerProvisioningStatusRefreshFunc(ctx, lbClient, *id),
				ContinuousTargetOccurence: 10,
				Timeout:                   time.Until(timeout),
			}

			if _, err := lbStatus.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for parent resource loadbalancer status to be ready error: %+v", err)
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, pool)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func loadbalacnerProvisioningStatusRefreshFunc(ctx context.Context, client *network.LoadBalancersClient, id parse.BackendAddressPoolAddressId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lbClient, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, "")
		if err != nil {
			return nil, "", fmt.Errorf("retrieving load balancer error： %+v", err)
		}
		return lbClient, string(lbClient.ProvisioningState), nil
	}
}
