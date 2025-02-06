// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
			ValidateFunc: loadbalancers.ValidateLoadBalancerBackendAddressPoolID,
		},

		"virtual_network_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"backend_address_ip_configuration_id"},
			ValidateFunc:  commonids.ValidateVirtualNetworkID,
		},

		"ip_address": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsIPAddress,
		},

		"backend_address_ip_configuration_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"virtual_network_id"},
			ValidateFunc:  loadbalancers.ValidateFrontendIPConfigurationID,
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
			lbClient := metadata.Client.LoadBalancers.LoadBalancersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model BackendAddressPoolAddressModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			poolId, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(model.BackendAddressPoolId)
			if err != nil {
				return err
			}

			locks.ByName(poolId.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(poolId.BackendAddressPoolName, backendAddressPoolResourceName)

			// Backend Addresses can not be created for Basic sku, so we have to check
			plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: subscriptionId, ResourceGroupName: poolId.ResourceGroupName, LoadBalancerName: poolId.LoadBalancerName}
			lb, err := metadata.Client.LoadBalancers.LoadBalancersClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", plbId, err)
			}

			isBasicSku := true
			if lb := lb.Model; lb != nil {
				if lb.Sku != nil && pointer.From(lb.Sku.Name) != loadbalancers.LoadBalancerSkuNameBasic {
					isBasicSku = false
				}
				if isBasicSku {
					return fmt.Errorf("Backend Addresses are not supported on Basic SKU Load Balancers")
				}
				if lb.Sku != nil && pointer.From(lb.Sku.Tier) == loadbalancers.LoadBalancerSkuTierGlobal {
					if model.FrontendIPConfiguration == "" {
						return fmt.Errorf("Please set a Regional Backend Address Pool Addresses for the Global load balancer")
					}
				}

				id := parse.NewBackendAddressPoolAddressID(subscriptionId, poolId.ResourceGroupName, poolId.LoadBalancerName, poolId.BackendAddressPoolName, model.Name)
				pool, err := lbClient.LoadBalancerBackendAddressPoolsGet(ctx, *poolId)
				if err != nil {
					return fmt.Errorf("retrieving %s: %+v", *poolId, err)
				}

				if pool.Model == nil {
					return fmt.Errorf("retrieving %s: `model` was nil", *poolId)
				}

				if pool.Model.Properties == nil {
					return fmt.Errorf("retrieving %s: `properties` was nil", *poolId)
				}

				if lb.Sku != nil && pointer.From(lb.Sku.Tier) == loadbalancers.LoadBalancerSkuTierRegional {
					if pointer.From(pool.Model.Properties.SyncMode) != "Manual" && (model.IPAddress != "" && model.VirtualNetworkId == "" || model.IPAddress == "" && model.VirtualNetworkId != "") {
						return fmt.Errorf("For regional load balancer, `ip_address` and `virtual_network_id` should be specified when sync mode is not `Manual`")
					}
				}

				addresses := make([]loadbalancers.LoadBalancerBackendAddress, 0)
				if pool.Model.Properties.LoadBalancerBackendAddresses != nil {
					addresses = *pool.Model.Properties.LoadBalancerBackendAddresses
				}

				metadata.Logger.Infof("checking for existing %s..", id)
				for _, address := range addresses {
					if address.Name == nil {
						continue
					}

					if *address.Name == model.Name {
						return metadata.ResourceRequiresImport(r.ResourceType(), id)
					}
				}

				if pointer.From(lb.Sku.Tier) == loadbalancers.LoadBalancerSkuTierGlobal {
					addresses = append(addresses, loadbalancers.LoadBalancerBackendAddress{
						Name: pointer.To(model.Name),
						Properties: &loadbalancers.LoadBalancerBackendAddressPropertiesFormat{
							LoadBalancerFrontendIPConfiguration: &loadbalancers.SubResource{
								Id: pointer.To(model.FrontendIPConfiguration),
							},
						},
					})
				} else {
					address := loadbalancers.LoadBalancerBackendAddress{
						Properties: &loadbalancers.LoadBalancerBackendAddressPropertiesFormat{},
						Name:       pointer.To(model.Name),
					}
					if model.IPAddress != "" {
						address.Properties.IPAddress = pointer.To(model.IPAddress)
					}
					if model.VirtualNetworkId != "" {
						address.Properties.VirtualNetwork = &loadbalancers.SubResource{
							Id: pointer.To(model.VirtualNetworkId),
						}
					}
					addresses = append(addresses, address)
				}

				pool.Model.Properties.LoadBalancerBackendAddresses = &addresses

				metadata.Logger.Infof("adding %s..", id)
				err = lbClient.LoadBalancerBackendAddressPoolsCreateOrUpdateThenPoll(ctx, *poolId, *pool.Model)
				if err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
				metadata.Logger.Infof("waiting for update %s..", id)

				metadata.SetID(id)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r BackendAddressPoolAddressResource) Read() sdk.ResourceFunc {
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

func (r BackendAddressPoolAddressResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			lbClient := metadata.Client.LoadBalancers.LoadBalancersClient
			id, err := parse.BackendAddressPoolAddressID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
			defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

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

			timeout, _ := ctx.Deadline()
			lbStatus := &pluginsdk.StateChangeConf{
				Pending:                   []string{string(loadbalancers.ProvisioningStateUpdating)},
				Target:                    []string{string(loadbalancers.ProvisioningStateSucceeded)},
				MinTimeout:                5 * time.Second,
				Refresh:                   loadbalacnerProvisioningStatusRefreshFunc(ctx, lbClient, *id),
				ContinuousTargetOccurence: 10,
				Timeout:                   time.Until(timeout),
			}

			if _, err := lbStatus.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for parent resource loadbalancer status to be ready error: %+v", err)
			}

			addresses := make([]loadbalancers.LoadBalancerBackendAddress, 0)
			if pool.Model.Properties.LoadBalancerBackendAddresses != nil {
				addresses = *pool.Model.Properties.LoadBalancerBackendAddresses
			}

			newAddresses := make([]loadbalancers.LoadBalancerBackendAddress, 0)
			for _, address := range addresses {
				if address.Name == nil {
					continue
				}

				if *address.Name != id.AddressName {
					newAddresses = append(newAddresses, address)
				}
			}

			metadata.Logger.Infof("removing %s..", *id)
			pool.Model.Properties.LoadBalancerBackendAddresses = &newAddresses

			err = lbClient.LoadBalancerBackendAddressPoolsCreateOrUpdateThenPoll(ctx, poolId, *pool.Model)
			if err != nil {
				return fmt.Errorf("removing %s: %+v", *id, err)
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

			poolId := loadbalancers.NewLoadBalancerBackendAddressPoolID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
			pool, err := lbClient.LoadBalancerBackendAddressPoolsGet(ctx, poolId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if pool.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if pool.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroup, LoadBalancerName: id.LoadBalancerName}
			lb, err := lbClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", plbId, err)
			}

			addresses := make([]loadbalancers.LoadBalancerBackendAddress, 0)
			if backendAddress := pool.Model.Properties.LoadBalancerBackendAddresses; backendAddress != nil {
				addresses = *backendAddress
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

			if lb.Model != nil && pointer.From(lb.Model.Sku.Tier) == loadbalancers.LoadBalancerSkuTierGlobal {
				addresses[index] = loadbalancers.LoadBalancerBackendAddress{
					Name: pointer.To(model.Name),
					Properties: &loadbalancers.LoadBalancerBackendAddressPropertiesFormat{
						LoadBalancerFrontendIPConfiguration: &loadbalancers.SubResource{
							Id: pointer.To(model.FrontendIPConfiguration),
						},
					},
				}
			} else {
				addresses[index] = loadbalancers.LoadBalancerBackendAddress{
					Properties: &loadbalancers.LoadBalancerBackendAddressPropertiesFormat{
						IPAddress: pointer.To(model.IPAddress),
						VirtualNetwork: &loadbalancers.SubResource{
							Id: pointer.To(model.VirtualNetworkId),
						},
					},
					Name: pointer.To(id.AddressName),
				}
			}

			pool.Model.Properties.LoadBalancerBackendAddresses = &addresses

			timeout, _ := ctx.Deadline()
			lbStatus := &pluginsdk.StateChangeConf{
				Pending:                   []string{string(loadbalancers.ProvisioningStateUpdating)},
				Target:                    []string{string(loadbalancers.ProvisioningStateSucceeded)},
				MinTimeout:                5 * time.Minute,
				PollInterval:              10 * time.Second,
				Refresh:                   loadbalacnerProvisioningStatusRefreshFunc(ctx, lbClient, *id),
				ContinuousTargetOccurence: 10,
				Timeout:                   time.Until(timeout),
			}

			if _, err := lbStatus.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for parent resource loadbalancer status to be ready error: %+v", err)
			}

			err = lbClient.LoadBalancerBackendAddressPoolsCreateOrUpdateThenPoll(ctx, poolId, *pool.Model)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func loadbalacnerProvisioningStatusRefreshFunc(ctx context.Context, client *loadbalancers.LoadBalancersClient, id parse.BackendAddressPoolAddressId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroup, LoadBalancerName: id.LoadBalancerName}
		lbClient, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
		if err != nil {
			return nil, "", fmt.Errorf("retrieving load balancer error： %+v", err)
		}
		if model := lbClient.Model; model != nil {
			if props := model.Properties; props != nil {
				return lbClient, string(pointer.From(props.ProvisioningState)), nil
			}
		}
		return lbClient, "", nil
	}
}
