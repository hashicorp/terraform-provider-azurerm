// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb_probe

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmLoadBalancerProbeCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
						if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
							return tf.ImportAsExistsError(azureLoadBalancerProbeResourceName, *existingProbe.Id)
						}
					}

					// this probe is being updated/reapplied remove old copy from the slice
					probes = append(probes[:existingProbeIndex], probes[existingProbeIndex+1:]...)
				}
			}

			props.Probes = &probes

			if d.IsNewResource() {
				if err := client.CreateOrUpdateCallbackThenPoll(ctx, plbId, *model, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
					return fmt.Errorf("creating %s: %+v", id, err)
				}
				d.SetId(id.ID())
			} else {
				if err := client.CreateOrUpdateThenPoll(ctx, plbId, *model); err != nil {
					return fmt.Errorf("updating %s: %+v", id, err)
				}
			}
		}
	}

	return resourceArmLoadBalancerProbeRead(d, meta)
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

	if v, ok := d.GetOk("no_healthy_backends_behavior"); ok {
		properties.NoHealthyBackendsBehavior = pointer.ToEnum[loadbalancers.ProbeNoHealthyBackendsBehavior](v.(string))
	}

	return &loadbalancers.Probe{
		Name:       pointer.To(d.Get("name").(string)),
		Properties: &properties,
	}
}
