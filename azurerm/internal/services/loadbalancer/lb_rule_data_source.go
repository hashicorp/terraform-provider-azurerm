package loadbalancer

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmLoadBalancerRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmLoadBalancerRuleRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: ValidateArmLoadBalancerRuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"loadbalancer_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.LoadBalancerID,
			},

			"frontend_ip_configuration_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"frontend_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"backend_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"backend_address_pool_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"probe_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enable_floating_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enable_tcp_reset": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"disable_outbound_snat": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"idle_timeout_in_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"load_distribution": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmLoadBalancerRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	loadBalancerID, err := parse.LoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}

	loadBalancer, exists, err := retrieveLoadBalancerById(ctx, client, *loadBalancerID)
	if err != nil {
		return fmt.Errorf("retrieving Load Balancer by ID: %+v", err)
	}
	if !exists {
		return fmt.Errorf("Load Balancer %q (Resource Group %q) was not found", loadBalancerID.Name, loadBalancerID.ResourceGroup)
	}

	lbRuleClient := meta.(*clients.Client).LoadBalancers.LoadBalancingRulesClient
	ctx, cancel = timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := lbRuleClient.Get(ctx, resourceGroup, *loadBalancer.Name, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Load Balancer Rule %q was not found in Load Balancer %q (Resource Group: %q)", name, *loadBalancer.Name, resourceGroup)
		}

		return fmt.Errorf("retrieving Load Balancer %s: %s", name, err)
	}

	d.SetId(*resp.ID)

	if props := resp.LoadBalancingRulePropertiesFormat; props != nil {
		frontendIPConfigurationName, err := parse.LoadBalancerFrontendIpConfigurationID(*props.FrontendIPConfiguration.ID)
		if err != nil {
			return err
		}

		d.Set("frontend_ip_configuration_name", frontendIPConfigurationName.FrontendIPConfigurationName)
		d.Set("protocol", props.Protocol)
		d.Set("frontend_port", props.FrontendPort)
		d.Set("backend_port", props.BackendPort)

		if props.BackendAddressPool != nil {
			if err := d.Set("backend_address_pool_id", props.BackendAddressPool.ID); err != nil {
				return fmt.Errorf("setting `backend_address_pool_id`: %+v", err)
			}
		}

		if props.Probe != nil {
			if err := d.Set("probe_id", props.Probe.ID); err != nil {
				return fmt.Errorf("setting `probe_id`: %+v", err)
			}
		}

		if err := d.Set("enable_floating_ip", props.EnableFloatingIP); err != nil {
			return fmt.Errorf("setting `enable_floating_ip`: %+v", err)
		}

		if err := d.Set("enable_tcp_reset", props.EnableTCPReset); err != nil {
			return fmt.Errorf("setting `enable_tcp_reset`: %+v", err)
		}

		if err := d.Set("disable_outbound_snat", props.DisableOutboundSnat); err != nil {
			return fmt.Errorf("setting `disable_outbound_snat`: %+v", err)
		}

		if err := d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes); err != nil {
			return fmt.Errorf("setting `idle_timeout_in_minutes`: %+v", err)
		}

		if err := d.Set("load_distribution", props.LoadDistribution); err != nil {
			return fmt.Errorf("setting `load_distribution`: %+v", err)
		}
	}

	return nil
}
