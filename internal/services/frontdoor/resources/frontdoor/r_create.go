// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFrontDoorCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := frontdoors.NewFrontDoorID(subscriptionId, resourceGroup, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_frontdoor", id.ID())
		}
	}

	var backendCertNameCheck bool
	var backendPoolsSendReceiveTimeoutSeconds int64
	if bps, ok := d.Get("backend_pool_settings").([]any); ok && len(bps) > 0 {
		bpsMap := bps[0].(map[string]any)
		if v, ok := bpsMap["enforce_backend_pools_certificate_name_check"].(bool); ok {
			backendCertNameCheck = v
		}
		if v, ok := bpsMap["backend_pools_send_receive_timeout_seconds"].(int); ok {
			backendPoolsSendReceiveTimeoutSeconds = int64(v)
		}
	}

	friendlyName := d.Get("friendly_name").(string)
	routingRules := d.Get("routing_rule").([]any)
	loadBalancingSettings := d.Get("backend_pool_load_balancing").([]any)
	healthProbeSettings := d.Get("backend_pool_health_probe").([]any)
	backendPools := d.Get("backend_pool").([]any)
	frontendEndpoints := d.Get("frontend_endpoint").([]any)

	enabledState := expandFrontDoorEnabledState(d.Get("load_balancer_enabled").(bool))
	t := d.Get("tags").(map[string]any)

	frontDoorParameters := frontdoors.FrontDoor{
		Location: pointer.To("Global"),
		Properties: &frontdoors.FrontDoorProperties{
			FriendlyName:          pointer.To(friendlyName),
			RoutingRules:          expandFrontDoorRoutingRule(routingRules, id, nil),
			BackendPools:          expandFrontDoorBackendPools(backendPools, id),
			BackendPoolsSettings:  expandFrontDoorBackendPoolsSettings(backendCertNameCheck, backendPoolsSendReceiveTimeoutSeconds),
			FrontendEndpoints:     expandFrontDoorFrontendEndpoint(frontendEndpoints, id),
			HealthProbeSettings:   expandFrontDoorHealthProbeSettingsModel(healthProbeSettings, id),
			LoadBalancingSettings: expandFrontDoorLoadBalancingSettingsModel(loadBalancingSettings, id),
			EnabledState:          &enabledState,
		},
		Tags: tags.Expand(t),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, frontDoorParameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	d.Set("explicit_resource_order", flattenExplicitResourceOrder(backendPools, frontendEndpoints, routingRules, loadBalancingSettings, healthProbeSettings, id))

	return resourceFrontDoorRead(d, meta)
}
