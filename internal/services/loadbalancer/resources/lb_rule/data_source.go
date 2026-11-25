// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package lb_rule

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Load Balancer Rule.",
		Read:        dataSourceArmLoadBalancerRuleRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of this Load Balancer Rule.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.RuleName,
			},

			"loadbalancer_id": {
				Description:  "The ID of the Load Balancer Rule.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: loadbalancers.ValidateLoadBalancerID,
			},

			"frontend_ip_configuration_name": {
				Description: "The name of the frontend IP configuration to which the rule is associated.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"protocol": {
				Description: "The transport protocol for the external endpoint.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"frontend_port": {
				Description: "The port for the external endpoint.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"backend_port": {
				Description: "The port used for internal connections on the endpoint.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"backend_address_pool_id": {
				Description: "A reference to a Backend Address Pool over which this Load Balancing Rule operates.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"probe_id": {
				Description: "A reference to a Probe used by this Load Balancing Rule.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"floating_ip_enabled": {
				Description: "If Floating IPs are enabled for this Load Balancer Rule.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"tcp_reset_enabled": {
				Description: "If TCP Reset is enabled for this Load Balancer Rule.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"disable_outbound_snat": {
				Description: "If outbound SNAT is enabled for this Load Balancer Rule.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"idle_timeout_in_minutes": {
				Description: "The idle timeout in minutes for TCP connections.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"load_distribution": {
				Description: "The load balancing distribution type used by the Load Balancer.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceArmLoadBalancerRuleRead(d *pluginsdk.ResourceData, meta any) error {
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

	id := loadbalancers.NewLoadBalancerLoadBalancingRuleID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroupName, loadBalancerId.LoadBalancerName, name)
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
			if err := d.Set("disable_outbound_snat", pointer.From(props.DisableOutboundSnat)); err != nil {
				return fmt.Errorf("setting `disable_outbound_snat`: %+v", err)
			}

			if err := d.Set("idle_timeout_in_minutes", int(pointer.From(props.IdleTimeoutInMinutes))); err != nil {
				return fmt.Errorf("setting `idle_timeout_in_minutes`: %+v", err)
			}

			if err := d.Set("load_distribution", pointer.FromEnum(props.LoadDistribution)); err != nil {
				return fmt.Errorf("setting `load_distribution`: %+v", err)
			}
		}
	}
	return nil
}
