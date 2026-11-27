// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package lb_outbound_rule

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	loadbalancerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Load Balancer Outbound Rule.",
		Read:        dataSourceArmLoadBalancerOutboundRuleRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of this Load Balancer Outbound Rule.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"loadbalancer_id": {
				Description:  "The ID of the Load Balancer in which the Outbound Rule exists.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: loadbalancers.ValidateLoadBalancerID,
			},

			"frontend_ip_configuration": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the Frontend IP Configuration.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"id": {
							Description: "The ID of the Frontend IP Configuration.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"backend_address_pool_id": {
				Description: "The ID of the Backend Address Pool. Outbound traffic is randomly load balanced across IPs in the backend IPs.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"protocol": {
				Description: "The transport protocol for the external endpoint.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"tcp_reset_enabled": {
				Description: "Whether the bidirectional TCP Reset on TCP flow idle timeout or unexpected connection termination is enabled. This value is useful when the protocol is set to TCP.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"allocated_outbound_ports": {
				Description: "The number of outbound ports used for NAT.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},

			"idle_timeout_in_minutes": {
				Description: "The timeout for the TCP idle connection.",
				Type:        pluginsdk.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourceArmLoadBalancerOutboundRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := loadbalancerClient.From(meta).LoadBalancersClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
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

	id := loadbalancers.NewOutboundRuleID(loadBalancerId.SubscriptionId, loadBalancerId.ResourceGroupName, loadBalancerId.LoadBalancerName, name)

	if model := loadBalancer.Model; model != nil {
		config, _, exists := FindLoadBalancerOutboundRuleByName(model, id.OutboundRuleName)
		if !exists {
			return fmt.Errorf("%s was not found", id)
		}

		d.SetId(id.ID())

		if props := config.Properties; props != nil {
			d.Set("allocated_outbound_ports", int(pointer.From(props.AllocatedOutboundPorts)))

			backendAddressPoolId := ""
			if props.BackendAddressPool.Id != nil {
				bapid, err := loadbalancers.ParseLoadBalancerBackendAddressPoolID(*props.BackendAddressPool.Id)
				if err != nil {
					return err
				}

				backendAddressPoolId = bapid.ID()
			}
			d.Set("backend_address_pool_id", backendAddressPoolId)
			d.Set("tcp_reset_enabled", pointer.From(props.EnableTcpReset))

			frontendIpConfigurations := make([]any, 0)
			if configs := props.FrontendIPConfigurations; configs != nil {
				for _, feConfig := range configs {
					if feConfig.Id == nil {
						continue
					}
					feid, err := loadbalancers.ParseFrontendIPConfigurationID(*feConfig.Id)
					if err != nil {
						return err
					}

					frontendIpConfigurations = append(frontendIpConfigurations, map[string]any{
						"id":   feid.ID(),
						"name": feid.FrontendIPConfigurationName,
					})
				}
			}
			d.Set("frontend_ip_configuration", frontendIpConfigurations)
			d.Set("idle_timeout_in_minutes", int(pointer.From(props.IdleTimeoutInMinutes)))
			d.Set("protocol", string(props.Protocol))
		}
	}
	return nil
}
