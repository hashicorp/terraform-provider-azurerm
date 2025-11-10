// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

var backendAddressPoolResourceName = "azurerm_lb_backend_address_pool"

func resourceArmLoadBalancerBackendAddressPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerBackendAddressPoolCreateUpdate,
		Update: resourceArmLoadBalancerBackendAddressPoolCreateUpdate,
		Read:   resourceArmLoadBalancerBackendAddressPoolRead,
		Delete: resourceArmLoadBalancerBackendAddressPoolDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*loadbalancers.LoadBalancerId, error) {
			id, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(input)
			if err != nil {
				return nil, err
			}

			lbId := loadbalancers.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName)
			return &lbId, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"loadbalancer_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: loadbalancers.ValidateLoadBalancerID,
			},

			"synchronous_mode": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(loadbalancers.PossibleValuesForSyncMode(), false),
				RequiredWith: []string{"virtual_network_id"},
			},

			"tunnel_interface": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"identifier": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},

						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(loadbalancers.GatewayLoadBalancerTunnelInterfaceTypeNone),
								string(loadbalancers.GatewayLoadBalancerTunnelInterfaceTypeInternal),
								string(loadbalancers.GatewayLoadBalancerTunnelInterfaceTypeExternal),
							},
								false,
							),
						},

						"protocol": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(loadbalancers.GatewayLoadBalancerTunnelProtocolNone),
								string(loadbalancers.GatewayLoadBalancerTunnelProtocolNative),
								string(loadbalancers.GatewayLoadBalancerTunnelProtocolVXLAN),
							},
								false,
							),
						},

						"port": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},
					},
				},
			},

			"virtual_network_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: commonids.ValidateVirtualNetworkID,
			},

			"backend_ip_configurations": {
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
		},
	}
}

func resourceArmLoadBalancerBackendAddressPoolCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	loadBalancerId, err := loadbalancers.ParseLoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Load Balancer Name and Group: %+v", err)
	}

	name := d.Get("name").(string)
	id := loadbalancers.NewLoadBalancerBackendAddressPoolID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroupName, loadBalancerId.LoadBalancerName, name)

	if d.IsNewResource() {
		existing, err := lbClient.LoadBalancerBackendAddressPoolsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_lb_backend_address_pool", id.ID())
		}
	}

	locks.ByName(name, backendAddressPoolResourceName)
	defer locks.UnlockByName(name, backendAddressPoolResourceName)

	locks.ByID(loadBalancerId.ID())
	defer locks.UnlockByID(loadBalancerId.ID())

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: loadBalancerId.SubscriptionId, ResourceGroupName: loadBalancerId.ResourceGroupName, LoadBalancerName: loadBalancerId.LoadBalancerName}
	lb, err := lbClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(lb.HttpResponse) {
			return fmt.Errorf("%s was not found", *loadBalancerId)
		}
		return fmt.Errorf("retrieving %s: %+v", *loadBalancerId, err)
	}

	param := loadbalancers.BackendAddressPool{
		Name: &id.BackendAddressPoolName,
	}

	// Since API version 2020-05-01, there are two ways to CRUD backend address pool - either via the LB endpoint or via the
	// dedicated BAP endpoint. While based on different sku of the LB, users should insist on interacting one of the two endpoints:
	// - Basic sku: interact with LB endpoint for CUD
	// - Standard sku: interact with BAP endpoint for CUD
	// Particularly, the BAP endpoint can be used for R for bot cases.
	// See: https://github.com/Azure/azure-rest-api-specs/issues/11234 for details.
	if lb.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *loadBalancerId)
	}
	sku := lb.Model.Sku
	if sku == nil {
		return fmt.Errorf("nil or empty `sku` for Load Balancer %q for Backend Address Pool %q was not found", *loadBalancerId, id)
	}

	if len(d.Get("tunnel_interface").([]interface{})) != 0 && *sku.Name != loadbalancers.LoadBalancerSkuNameGateway {
		return fmt.Errorf("only the Gateway (sku) Load Balancer allows IP based Backend Address Pool configuration,"+
			"whilst %q is of sku %s", id, *sku.Name)
	}
	if len(d.Get("tunnel_interface").([]interface{})) == 0 && *sku.Name == loadbalancers.LoadBalancerSkuNameGateway {
		return fmt.Errorf("`tunnel_interface` is required for %q when sku is set to %s", id, *sku.Name)
	}

	if _, ok := d.GetOk("synchronous_mode"); ok && *sku.Name != loadbalancers.LoadBalancerSkuNameStandard {
		return fmt.Errorf("`synchronous_mode` can set only for Load Balancer with `Standard` SKU")
	}

	if v, ok := d.GetOk("virtual_network_id"); ok {
		param.Properties = &loadbalancers.BackendAddressPoolPropertiesFormat{
			VirtualNetwork: &loadbalancers.SubResource{
				Id: pointer.To(v.(string)),
			},
		}
	}

	if v, ok := d.GetOk("synchronous_mode"); ok {
		if param.Properties == nil {
			param.Properties = &loadbalancers.BackendAddressPoolPropertiesFormat{
				SyncMode: pointer.To(loadbalancers.SyncMode(v.(string))),
			}
		} else {
			param.Properties.SyncMode = pointer.To(loadbalancers.SyncMode(v.(string)))
		}
	}

	if properties := lb.Model.Properties; properties != nil {
		switch *sku.Name {
		case loadbalancers.LoadBalancerSkuNameBasic:
			if !d.IsNewResource() && d.HasChange("virtual_network_id") {
				return fmt.Errorf("updating the virtual_network_id of Backend Address Pool %q is not allowed for basic (sku) Load Balancer", id)
			}

			// Insert this BAP and update the LB since the dedicated BAP endpoint doesn't work for the Basic sku.
			backendAddressPools := append(*properties.BackendAddressPools, param)
			_, existingPoolIndex, exists := FindLoadBalancerBackEndAddressPoolByName(lb.Model, id.BackendAddressPoolName)
			if exists {
				// this pool is being updated/reapplied remove the old copy from the slice
				backendAddressPools = append(backendAddressPools[:existingPoolIndex], backendAddressPools[existingPoolIndex+1:]...)
			}

			properties.BackendAddressPools = &backendAddressPools

			err := lbClient.CreateOrUpdateThenPoll(ctx, plbId, *lb.Model)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *loadBalancerId, err)
			}
		case loadbalancers.LoadBalancerSkuNameStandard:
			if param.Properties == nil {
				param.Properties = &loadbalancers.BackendAddressPoolPropertiesFormat{
					// NOTE: Backend Addresses are managed using `azurerm_lb_backend_pool_address`
				}
			}

			err := lbClient.LoadBalancerBackendAddressPoolsCreateOrUpdateThenPoll(ctx, id, param)
			if err != nil {
				return fmt.Errorf("creating/updating Load Balancer Backend Address Pool %q: %+v", id, err)
			}
		case loadbalancers.LoadBalancerSkuNameGateway:
			if param.Properties == nil {
				param.Properties = &loadbalancers.BackendAddressPoolPropertiesFormat{}
			}
			param.Properties.TunnelInterfaces = expandGatewayLoadBalancerTunnelInterfaces(d.Get("tunnel_interface").([]interface{}))

			err := lbClient.LoadBalancerBackendAddressPoolsCreateOrUpdateThenPoll(ctx, id, param)
			if err != nil {
				return fmt.Errorf("creating/updating %q: %+v", id, err)
			}
		}

		d.SetId(id.ID())
	}

	return resourceArmLoadBalancerBackendAddressPoolRead(d, meta)
}

func resourceArmLoadBalancerBackendAddressPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(d.Id())
	if err != nil {
		return err
	}

	lbId := loadbalancers.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName)

	resp, err := lbClient.LoadBalancerBackendAddressPoolsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			log.Printf("[INFO] %s was not found - removing from state", *id)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.BackendAddressPoolName)
	d.Set("loadbalancer_id", lbId.ID())

	if model := resp.Model; model != nil {
		if properties := model.Properties; properties != nil {
			if err := d.Set("tunnel_interface", flattenGatewayLoadBalancerTunnelInterfaces(properties.TunnelInterfaces)); err != nil {
				return fmt.Errorf("setting `tunnel_interface`: %v", err)
			}

			var backendIPConfigurations []string
			if configs := properties.BackendIPConfigurations; configs != nil {
				for _, backendConfig := range *configs {
					if backendConfig.Id == nil {
						continue
					}
					backendIPConfigurations = append(backendIPConfigurations, *backendConfig.Id)
				}
			}
			if err := d.Set("backend_ip_configurations", backendIPConfigurations); err != nil {
				return fmt.Errorf("setting `backend_ip_configurations`: %v", err)
			}

			d.Set("synchronous_mode", pointer.From(properties.SyncMode))

			network := ""
			if vnet := properties.VirtualNetwork; vnet != nil && vnet.Id != nil {
				network = *vnet.Id
			}
			d.Set("virtual_network_id", network)

			var loadBalancingRules []string
			if rules := properties.LoadBalancingRules; rules != nil {
				for _, rule := range *rules {
					if rule.Id == nil {
						continue
					}
					loadBalancingRules = append(loadBalancingRules, *rule.Id)
				}
			}
			if err := d.Set("load_balancing_rules", loadBalancingRules); err != nil {
				return fmt.Errorf("setting `load_balancing_rules`: %v", err)
			}

			var outboundRules []string
			if rules := properties.OutboundRules; rules != nil {
				for _, rule := range *rules {
					if rule.Id == nil {
						continue
					}
					outboundRules = append(outboundRules, *rule.Id)
				}
			}
			if err := d.Set("outbound_rules", outboundRules); err != nil {
				return fmt.Errorf("setting `outbound_rules`: %v", err)
			}

			var inboundNATRules []string
			if rules := properties.InboundNatRules; rules != nil {
				for _, rule := range *rules {
					if rule.Id == nil {
						continue
					}
					inboundNATRules = append(inboundNATRules, *rule.Id)
				}
			}
			if err := d.Set("inbound_nat_rules", inboundNATRules); err != nil {
				return fmt.Errorf("setting `inbound_nat_rules`: %v", err)
			}
		}
	}

	return nil
}

func resourceArmLoadBalancerBackendAddressPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := loadbalancers.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName)
	loadBalancerID := loadBalancerId.ID()
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
	defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	lb, err := lbClient.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(lb.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", loadBalancerId, err)
	}

	if lb.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", plbId)
	}
	sku := lb.Model.Sku
	if sku == nil {
		return fmt.Errorf("nil or empty `sku` for Load Balancer %q for Backend Address Pool %q was not found", loadBalancerId, id)
	}

	if *sku.Name == loadbalancers.LoadBalancerSkuNameBasic {
		_, index, exists := FindLoadBalancerBackEndAddressPoolByName(lb.Model, id.BackendAddressPoolName)
		if !exists {
			return nil
		}

		if lb.Model.Properties == nil {
			return fmt.Errorf("retrieving %s: properties was nil", *id)
		}

		backEndPools := *lb.Model.Properties.BackendAddressPools
		backEndPools = append(backEndPools[:index], backEndPools[index+1:]...)
		lb.Model.Properties.BackendAddressPools = &backEndPools

		err := lbClient.CreateOrUpdateThenPoll(ctx, plbId, *lb.Model)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", loadBalancerId, err)
		}
	} else {
		err := lbClient.LoadBalancerBackendAddressPoolsDeleteThenPoll(ctx, *id)
		if err != nil {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func expandGatewayLoadBalancerTunnelInterfaces(input []interface{}) *[]loadbalancers.GatewayLoadBalancerTunnelInterface {
	if len(input) == 0 {
		return nil
	}

	result := make([]loadbalancers.GatewayLoadBalancerTunnelInterface, 0)

	for _, e := range input {
		e := e.(map[string]interface{})
		result = append(result, loadbalancers.GatewayLoadBalancerTunnelInterface{
			Identifier: pointer.To(int64(e["identifier"].(int))),
			Type:       pointer.To(loadbalancers.GatewayLoadBalancerTunnelInterfaceType(e["type"].(string))),
			Protocol:   pointer.To(loadbalancers.GatewayLoadBalancerTunnelProtocol(e["protocol"].(string))),
			Port:       pointer.To(int64(e["port"].(int))),
		})
	}

	return &result
}

func flattenGatewayLoadBalancerTunnelInterfaces(input *[]loadbalancers.GatewayLoadBalancerTunnelInterface) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		var identifier int
		if e.Identifier != nil {
			identifier = int(*e.Identifier)
		}

		t := string(pointer.From(e.Type))

		protocol := string(pointer.From(e.Protocol))

		var port int
		if e.Port != nil {
			port = int(*e.Port)
		}

		output = append(output, map[string]interface{}{
			"identifier": identifier,
			"type":       t,
			"protocol":   protocol,
			"port":       port,
		})
	}

	return output
}
