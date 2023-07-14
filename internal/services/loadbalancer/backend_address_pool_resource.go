// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

var backendAddressPoolResourceName = "azurerm_lb_backend_address_pool"

func resourceArmLoadBalancerBackendAddressPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerBackendAddressPoolCreateUpdate,
		Update: resourceArmLoadBalancerBackendAddressPoolCreateUpdate,
		Read:   resourceArmLoadBalancerBackendAddressPoolRead,
		Delete: resourceArmLoadBalancerBackendAddressPoolDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*parse.LoadBalancerId, error) {
			id, err := parse.LoadBalancerBackendAddressPoolID(input)
			if err != nil {
				return nil, err
			}

			lbId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
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
				ValidateFunc: validate.LoadBalancerID,
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
								string(network.GatewayLoadBalancerTunnelInterfaceTypeNone),
								string(network.GatewayLoadBalancerTunnelInterfaceTypeInternal),
								string(network.GatewayLoadBalancerTunnelInterfaceTypeExternal),
							},
								false,
							),
						},

						"protocol": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.GatewayLoadBalancerTunnelProtocolNone),
								string(network.GatewayLoadBalancerTunnelProtocolNative),
								string(network.GatewayLoadBalancerTunnelProtocolVXLAN),
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
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerBackendAddressPoolsClient
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	loadBalancerId, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Load Balancer Name and Group: %+v", err)
	}

	name := d.Get("name").(string)
	id := parse.NewLoadBalancerBackendAddressPoolID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroup, loadBalancerId.Name, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Load Balancer Backend Address Pool %q: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_lb_backend_address_pool", id.ID())
		}
	}

	locks.ByName(name, backendAddressPoolResourceName)
	defer locks.UnlockByName(name, backendAddressPoolResourceName)

	locks.ByID(loadBalancerId.ID())
	defer locks.UnlockByID(loadBalancerId.ID())

	lb, err := lbClient.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(lb.Response) {
			return fmt.Errorf("Load Balancer %q for Backend Address Pool %q was not found", loadBalancerId, id)
		}
		return fmt.Errorf("failed to retrieve Load Balancer %q for Backend Address Pool %q: %+v", loadBalancerId, id, err)
	}

	param := network.BackendAddressPool{
		Name: &id.BackendAddressPoolName,
	}

	// Since API version 2020-05-01, there are two ways to CRUD backend address pool - either via the LB endpoint or via the
	// dedicated BAP endpoint. While based on different sku of the LB, users should insist on interacting one of the two endpoints:
	// - Basic sku: interact with LB endpoint for CUD
	// - Standard sku: interact with BAP endpoint for CUD
	// Particularly, the BAP endpoint can be used for R for bot cases.
	// See: https://github.com/Azure/azure-rest-api-specs/issues/11234 for details.
	sku := lb.Sku
	if sku == nil {
		return fmt.Errorf("nil or empty `sku` for Load Balancer %q for Backend Address Pool %q was not found", loadBalancerId, id)
	}

	if len(d.Get("tunnel_interface").([]interface{})) != 0 && sku.Name != network.LoadBalancerSkuNameGateway {
		return fmt.Errorf("only the Gateway (sku) Load Balancer allows IP based Backend Address Pool configuration,"+
			"whilst %q is of sku %s", id, sku.Name)
	}
	if len(d.Get("tunnel_interface").([]interface{})) == 0 && sku.Name == network.LoadBalancerSkuNameGateway {
		return fmt.Errorf("`tunnel_interface` is required for %q when sku is set to %s", id, sku.Name)
	}

	if v, ok := d.GetOk("virtual_network_id"); ok {
		param.BackendAddressPoolPropertiesFormat = &network.BackendAddressPoolPropertiesFormat{
			VirtualNetwork: &network.SubResource{
				ID: utils.String(v.(string)),
			}}
	}

	switch sku.Name {
	case network.LoadBalancerSkuNameBasic:
		if !d.IsNewResource() && d.HasChange("virtual_network_id") {
			return fmt.Errorf("updating the virtual_network_id of Backend Address Pool %q is not allowed for basic (sku) Load Balancer", id)
		}

		// Insert this BAP and update the LB since the dedicated BAP endpoint doesn't work for the Basic sku.
		backendAddressPools := append(*lb.LoadBalancerPropertiesFormat.BackendAddressPools, param)
		_, existingPoolIndex, exists := FindLoadBalancerBackEndAddressPoolByName(&lb, id.BackendAddressPoolName)
		if exists {
			// this pool is being updated/reapplied remove the old copy from the slice
			backendAddressPools = append(backendAddressPools[:existingPoolIndex], backendAddressPools[existingPoolIndex+1:]...)
		}

		lb.LoadBalancerPropertiesFormat.BackendAddressPools = &backendAddressPools

		future, err := lbClient.CreateOrUpdate(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, lb)
		if err != nil {
			return fmt.Errorf("updating Load Balancer %q for Backend Address Pool %q: %+v", loadBalancerId, id, err)
		}

		if err = future.WaitForCompletionRef(ctx, lbClient.Client); err != nil {
			return fmt.Errorf("waiting for update of Load Balancer %q for Backend Address Pool %q: %+v", loadBalancerId, id, err)
		}
	case network.LoadBalancerSkuNameStandard:
		if param.BackendAddressPoolPropertiesFormat == nil {
			param.BackendAddressPoolPropertiesFormat = &network.BackendAddressPoolPropertiesFormat{
				// NOTE: Backend Addresses are managed using `azurerm_lb_backend_pool_address`
			}
		}

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, param)
		if err != nil {
			return fmt.Errorf("creating/updating Load Balancer Backend Address Pool %q: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for Creating/Updating of Load Balancer Backend Address Pool %q: %+v", id, err)
		}
	case network.LoadBalancerSkuNameGateway:
		if param.BackendAddressPoolPropertiesFormat == nil {
			param.BackendAddressPoolPropertiesFormat = &network.BackendAddressPoolPropertiesFormat{}
		}
		param.BackendAddressPoolPropertiesFormat.TunnelInterfaces = expandGatewayLoadBalancerTunnelInterfaces(d.Get("tunnel_interface").([]interface{}))

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName, param)
		if err != nil {
			return fmt.Errorf("creating/updating %q: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for Creating/Updating of %q: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceArmLoadBalancerBackendAddressPoolRead(d, meta)
}

func resourceArmLoadBalancerBackendAddressPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerBackendAddressPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerBackendAddressPoolID(d.Id())
	if err != nil {
		return err
	}

	lbId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Backend Address Pool %q not found - removing from state", id)
			return nil
		}
		return fmt.Errorf("failed to retrieve Load Balancer Backend Address Pool %q: %+v", id, err)
	}

	d.Set("name", id.BackendAddressPoolName)
	d.Set("loadbalancer_id", lbId.ID())

	if props := resp.BackendAddressPoolPropertiesFormat; props != nil {
		if err := d.Set("tunnel_interface", flattenGatewayLoadBalancerTunnelInterfaces(props.TunnelInterfaces)); err != nil {
			return fmt.Errorf("setting `tunnel_interface`: %v", err)
		}

		var backendIPConfigurations []string
		if configs := props.BackendIPConfigurations; configs != nil {
			for _, backendConfig := range *configs {
				if backendConfig.ID == nil {
					continue
				}
				backendIPConfigurations = append(backendIPConfigurations, *backendConfig.ID)
			}
		}
		if err := d.Set("backend_ip_configurations", backendIPConfigurations); err != nil {
			return fmt.Errorf("setting `backend_ip_configurations`: %v", err)
		}

		network := ""
		if vnet := props.VirtualNetwork; vnet != nil && vnet.ID != nil {
			network = *vnet.ID
		}
		d.Set("virtual_network_id", network)

		var loadBalancingRules []string
		if rules := props.LoadBalancingRules; rules != nil {
			for _, rule := range *rules {
				if rule.ID == nil {
					continue
				}
				loadBalancingRules = append(loadBalancingRules, *rule.ID)
			}
		}
		if err := d.Set("load_balancing_rules", loadBalancingRules); err != nil {
			return fmt.Errorf("setting `load_balancing_rules`: %v", err)
		}

		var outboundRules []string
		if rules := props.OutboundRules; rules != nil {
			for _, rule := range *rules {
				if rule.ID == nil {
					continue
				}
				outboundRules = append(outboundRules, *rule.ID)
			}
		}
		if err := d.Set("outbound_rules", outboundRules); err != nil {
			return fmt.Errorf("setting `outbound_rules`: %v", err)
		}

		var inboundNATRules []string
		if rules := props.InboundNatRules; rules != nil {
			for _, rule := range *rules {
				if rule.ID == nil {
					continue
				}
				inboundNATRules = append(inboundNATRules, *rule.ID)
			}
		}
		if err := d.Set("inbound_nat_rules", inboundNATRules); err != nil {
			return fmt.Errorf("setting `inbound_nat_rules`: %v", err)
		}
	}

	return nil
}

func resourceArmLoadBalancerBackendAddressPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancerBackendAddressPoolsClient
	lbClient := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LoadBalancerBackendAddressPoolID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := parse.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName)
	loadBalancerID := loadBalancerId.ID()
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	locks.ByName(id.BackendAddressPoolName, backendAddressPoolResourceName)
	defer locks.UnlockByName(id.BackendAddressPoolName, backendAddressPoolResourceName)

	lb, err := lbClient.Get(ctx, loadBalancerId.ResourceGroup, loadBalancerId.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(lb.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to retrieve Load Balancer %q (resource group %q) for Backend Address Pool %q: %+v", loadBalancerId.Name, loadBalancerId.ResourceGroup, id.BackendAddressPoolName, err)
	}

	sku := lb.Sku
	if sku == nil {
		return fmt.Errorf("nil or empty `sku` for Load Balancer %q for Backend Address Pool %q was not found", loadBalancerId, id)
	}

	if sku.Name == network.LoadBalancerSkuNameBasic {
		_, index, exists := FindLoadBalancerBackEndAddressPoolByName(&lb, id.BackendAddressPoolName)
		if !exists {
			return nil
		}

		backEndPools := *lb.LoadBalancerPropertiesFormat.BackendAddressPools
		backEndPools = append(backEndPools[:index], backEndPools[index+1:]...)
		lb.LoadBalancerPropertiesFormat.BackendAddressPools = &backEndPools

		future, err := lbClient.CreateOrUpdate(ctx, id.ResourceGroup, id.LoadBalancerName, lb)
		if err != nil {
			return fmt.Errorf("updating Load Balancer %q (resource group %q) to remove Backend Address Pool %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.BackendAddressPoolName, err)
		}

		if err = future.WaitForCompletionRef(ctx, lbClient.Client); err != nil {
			return fmt.Errorf("waiting for update of Load Balancer %q (resource group %q) for Backend Address Pool %q: %+v", id.LoadBalancerName, id.ResourceGroup, id.BackendAddressPoolName, err)
		}
	} else {
		future, err := client.Delete(ctx, id.ResourceGroup, id.LoadBalancerName, id.BackendAddressPoolName)
		if err != nil {
			return fmt.Errorf("deleting Load Balancer Backend Address Pool %q: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for deleting of Load Balancer Backend Address Pool %q: %+v", id, err)
		}
	}

	return nil
}

func expandGatewayLoadBalancerTunnelInterfaces(input []interface{}) *[]network.GatewayLoadBalancerTunnelInterface {
	if len(input) == 0 {
		return nil
	}

	result := make([]network.GatewayLoadBalancerTunnelInterface, 0)

	for _, e := range input {
		e := e.(map[string]interface{})
		result = append(result, network.GatewayLoadBalancerTunnelInterface{
			Identifier: utils.Int32(int32(e["identifier"].(int))),
			Type:       network.GatewayLoadBalancerTunnelInterfaceType(e["type"].(string)),
			Protocol:   network.GatewayLoadBalancerTunnelProtocol(e["protocol"].(string)),
			Port:       utils.Int32(int32(e["port"].(int))),
		})
	}

	return &result
}

func flattenGatewayLoadBalancerTunnelInterfaces(input *[]network.GatewayLoadBalancerTunnelInterface) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		var identifier int
		if e.Identifier != nil {
			identifier = int(*e.Identifier)
		}

		t := string(e.Type)

		protocol := string(e.Protocol)

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
