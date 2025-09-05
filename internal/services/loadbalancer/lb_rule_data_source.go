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
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceArmLoadBalancerRule() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Read: dataSourceArmLoadBalancerRuleRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.RuleName,
			},

			"loadbalancer_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: loadbalancers.ValidateLoadBalancerID,
			},

			"frontend_ip_configuration_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"protocol": {
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

			"backend_address_pool_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"probe_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"floating_ip_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tcp_reset_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"disable_outbound_snat": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"idle_timeout_in_minutes": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"load_distribution": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}

	if !features.FivePointOh() {
		resource.Schema["enable_floating_ip"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeBool,
			Computed:   true,
			Deprecated: "The property `enable_floating_ip` has been deprecated in favour of `floating_ip_enabled` and will be removed in version 5.0 of the provider",
		}

		resource.Schema["enable_tcp_reset"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeBool,
			Computed:   true,
			Deprecated: "The property `enable_tcp_reset` has been deprecated in favour of `tcp_reset_enabled` and will be removed in version 5.0 of the provider",
		}
	}

	return resource
}

func dataSourceArmLoadBalancerRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	loadBalancerId, err := loadbalancers.ParseLoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: loadBalancerId.SubscriptionId, ResourceGroupName: loadBalancerId.ResourceGroupName, LoadBalancerName: loadBalancerId.LoadBalancerName}
	loadBalancer, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(loadBalancer.HttpResponse) {
			return fmt.Errorf("%s was not found", *loadBalancerId)
		}
		return fmt.Errorf("retrieving %s: %+v", *loadBalancerId, err)
	}

	id := loadbalancers.NewLoadBalancingRuleID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroupName, loadBalancerId.LoadBalancerName, name)
	resp, err := client.LoadBalancerLoadBalancingRulesGet(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			frontendIPConfigurationName, err := loadbalancers.ParseFrontendIPConfigurationID(*props.FrontendIPConfiguration.Id)
			if err != nil {
				return err
			}

			d.Set("frontend_ip_configuration_name", frontendIPConfigurationName.FrontendIPConfigurationName)
			d.Set("protocol", string(props.Protocol))
			d.Set("frontend_port", props.FrontendPort)
			d.Set("backend_port", pointer.From(props.BackendPort))

			if props.BackendAddressPool != nil {
				if err := d.Set("backend_address_pool_id", pointer.From(props.BackendAddressPool.Id)); err != nil {
					return fmt.Errorf("setting `backend_address_pool_id`: %+v", err)
				}
			}

			if props.Probe != nil {
				if err := d.Set("probe_id", pointer.From(props.Probe.Id)); err != nil {
					return fmt.Errorf("setting `probe_id`: %+v", err)
				}
			}
			if err := d.Set("floating_ip_enabled", pointer.From(props.EnableFloatingIP)); err != nil {
				return fmt.Errorf("setting `floating_ip_enabled`: %+v", err)
			}
			if err := d.Set("tcp_reset_enabled", pointer.From(props.EnableTcpReset)); err != nil {
				return fmt.Errorf("setting `tcp_reset_enabled`: %+v", err)
			}
			if !features.FivePointOh() {
				if err := d.Set("enable_floating_ip", pointer.From(props.EnableFloatingIP)); err != nil {
					return fmt.Errorf("setting `enable_floating_ip`: %+v", err)
				}
				if err := d.Set("enable_tcp_reset", pointer.From(props.EnableTcpReset)); err != nil {
					return fmt.Errorf("setting `enable_tcp_reset`: %+v", err)
				}
			}

			if err := d.Set("disable_outbound_snat", pointer.From(props.DisableOutboundSnat)); err != nil {
				return fmt.Errorf("setting `disable_outbound_snat`: %+v", err)
			}

			if err := d.Set("idle_timeout_in_minutes", int(pointer.From(props.IdleTimeoutInMinutes))); err != nil {
				return fmt.Errorf("setting `idle_timeout_in_minutes`: %+v", err)
			}

			if err := d.Set("load_distribution", string(pointer.From(props.LoadDistribution))); err != nil {
				return fmt.Errorf("setting `load_distribution`: %+v", err)
			}
		}
	}
	return nil
}
