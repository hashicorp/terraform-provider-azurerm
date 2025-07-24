// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerProbe() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmLoadBalancerProbeCreateUpdate,
		Read:   resourceArmLoadBalancerProbeRead,
		Update: resourceArmLoadBalancerProbeCreateUpdate,
		Delete: resourceArmLoadBalancerProbeDelete,

		Importer: loadBalancerSubResourceImporter(func(input string) (*loadbalancers.LoadBalancerId, error) {
			id, err := loadbalancers.ParseProbeID(input)
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

			"protocol": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(loadbalancers.ProbeProtocolTcp),
				ValidateFunc: validation.StringInSlice([]string{
					string(loadbalancers.ProbeProtocolHTTP),
					string(loadbalancers.ProbeProtocolHTTPS),
					string(loadbalancers.ProbeProtocolTcp),
				}, false),
			},

			"port": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumber,
			},

			"probe_threshold": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"request_path": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"interval_in_seconds": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      15,
				ValidateFunc: validation.IntAtLeast(5),
			},

			"number_of_probes": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Default:  2,
			},

			"load_balancer_rules": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Set: pluginsdk.HashString,
			},
		},
	}
}

func resourceArmLoadBalancerProbeCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	loadBalancerId, err := loadbalancers.ParseLoadBalancerID(d.Get("loadbalancer_id").(string))
	if err != nil {
		return err
	}
	loadBalancerIDRaw := loadBalancerId.ID()
	id := loadbalancers.NewProbeID(subscriptionId, loadBalancerId.ResourceGroupName, loadBalancerId.LoadBalancerName, d.Get("name").(string))
	locks.ByID(loadBalancerIDRaw)
	defer locks.UnlockByID(loadBalancerIDRaw)

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	loadBalancer, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(loadBalancer.HttpResponse) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer %q not found. Removing Proe %q from state", id.LoadBalancerName, id.ProbeName)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *loadBalancerId, err)
	}

	if model := loadBalancer.Model; model != nil {
		if props := model.Properties; props != nil {
			newProbe := expandAzureRmLoadBalancerProbe(d)
			probes := append(*props.Probes, *newProbe)
			existingProbe, existingProbeIndex, exists := FindLoadBalancerProbeByName(model, id.ProbeName)
			if exists {
				if id.ProbeName == *existingProbe.Name {
					if d.IsNewResource() {
						return tf.ImportAsExistsError("azurerm_lb_probe", *existingProbe.Id)
					}

					// this probe is being updated/reapplied remove old copy from the slice
					probes = append(probes[:existingProbeIndex], probes[existingProbeIndex+1:]...)
				}
			}

			props.Probes = &probes

			err := client.CreateOrUpdateThenPoll(ctx, plbId, *model)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
		}
	}

	d.SetId(id.ID())

	return resourceArmLoadBalancerProbeRead(d, meta)
}

func resourceArmLoadBalancerProbeRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseProbeID(d.Id())
	if err != nil {
		return err
	}

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	loadBalancer, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(loadBalancer.HttpResponse) {
			d.SetId("")
			log.Printf("[INFO] Load Balancer %q not found. Removing from state", id.LoadBalancerName)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", plbId, err)
	}

	if model := loadBalancer.Model; model != nil {
		config, _, exists := FindLoadBalancerProbeByName(model, id.ProbeName)
		if !exists {
			d.SetId("")
			log.Printf("[INFO] Load Balancer Probe %q not found. Removing from state", id.ProbeName)
			return nil
		}

		d.Set("name", config.Name)

		if props := config.Properties; props != nil {
			d.Set("interval_in_seconds", int(pointer.From(props.IntervalInSeconds)))
			d.Set("number_of_probes", int(pointer.From(props.NumberOfProbes)))
			d.Set("port", int(props.Port))
			d.Set("protocol", string(props.Protocol))
			d.Set("request_path", pointer.From(props.RequestPath))
			d.Set("probe_threshold", int(pointer.From(props.ProbeThreshold)))

			// TODO: parse/make these consistent
			var loadBalancerRules []string
			if rules := props.LoadBalancingRules; rules != nil {
				for _, ruleConfig := range *rules {
					if id := ruleConfig.Id; id != nil {
						loadBalancerRules = append(loadBalancerRules, *id)
					}
				}
			}
			if err := d.Set("load_balancer_rules", loadBalancerRules); err != nil {
				return fmt.Errorf("setting `load_balancer_rules` (Load Balancer %q Probe %q): %+v", id.LoadBalancerName, id.ProbeName, err)
			}
		}
	}
	return nil
}

func resourceArmLoadBalancerProbeDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LoadBalancers.LoadBalancersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := loadbalancers.ParseProbeID(d.Id())
	if err != nil {
		return err
	}

	loadBalancerId := loadbalancers.NewLoadBalancerID(id.SubscriptionId, id.ResourceGroupName, id.LoadBalancerName)
	loadBalancerID := loadBalancerId.ID()
	locks.ByID(loadBalancerID)
	defer locks.UnlockByID(loadBalancerID)

	plbId := loadbalancers.ProviderLoadBalancerId{SubscriptionId: id.SubscriptionId, ResourceGroupName: id.ResourceGroupName, LoadBalancerName: id.LoadBalancerName}
	loadBalancer, err := client.Get(ctx, plbId, loadbalancers.GetOperationOptions{})
	if err != nil {
		if response.WasNotFound(loadBalancer.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", loadBalancerId, err)
	}
	if model := loadBalancer.Model; model != nil {
		_, index, exists := FindLoadBalancerProbeByName(model, id.ProbeName)
		if !exists {
			return nil
		}

		if props := model.Properties; props != nil {
			probes := *props.Probes
			probes = append(probes[:index], probes[index+1:]...)
			props.Probes = &probes

			err := client.CreateOrUpdateThenPoll(ctx, plbId, *model)
			if err != nil {
				return fmt.Errorf("updating Load Balancer %q (Resource Group %q) for deletion of Probe %q: %+v", id.LoadBalancerName, id.ResourceGroupName, id.ProbeName, err)
			}
		}
	}
	return nil
}

func expandAzureRmLoadBalancerProbe(d *pluginsdk.ResourceData) *loadbalancers.Probe {
	properties := loadbalancers.ProbePropertiesFormat{
		NumberOfProbes:    pointer.To(int64(d.Get("number_of_probes").(int))),
		IntervalInSeconds: pointer.To(int64(d.Get("interval_in_seconds").(int))),
		Port:              int64(d.Get("port").(int)),
	}

	if v, ok := d.GetOk("protocol"); ok {
		properties.Protocol = loadbalancers.ProbeProtocol(v.(string))
	}

	if v, ok := d.GetOk("request_path"); ok {
		properties.RequestPath = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("probe_threshold"); ok {
		properties.ProbeThreshold = pointer.To(int64(v.(int)))
	}

	return &loadbalancers.Probe{
		Name:       pointer.To(d.Get("name").(string)),
		Properties: &properties,
	}
}
