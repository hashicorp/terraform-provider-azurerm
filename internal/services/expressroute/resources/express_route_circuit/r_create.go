// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuits"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceExpressRouteCircuitCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ExpressRoute.ExpressRouteCircuits
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id := expressroutecircuits.NewExpressRouteCircuitID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s : %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_express_route_circuit", id.ID())
		}
	}

	erc := expressroutecircuits.ExpressRouteCircuit{
		Name:     &id.ExpressRouteCircuitName,
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Sku:      expandExpressRouteCircuitSku(d.Get("sku").([]any)),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	erc.Properties = &expressroutecircuits.ExpressRouteCircuitPropertiesFormat{
		AuthorizationKey: pointer.To(d.Get("authorization_key").(string)),
	}

	if v, ok := d.GetOk("allow_classic_operations"); ok {
		erc.Properties.AllowClassicOperations = pointer.To(v.(bool))
	}

	if v, ok := d.GetOk("rate_limiting_enabled"); ok {
		erc.Properties.EnableDirectPortRateLimit = pointer.To(v.(bool))
	}

	// ServiceProviderProperties and expressRoutePorts/bandwidthInGbps properties are mutually exclusive
	if _, ok := d.GetOk("express_route_port_id"); ok {
		erc.Properties.ExpressRoutePort = &expressroutecircuits.SubResource{}
	} else {
		erc.Properties.ServiceProviderProperties = &expressroutecircuits.ExpressRouteCircuitServiceProviderProperties{}
	}

	if erc.Properties.ServiceProviderProperties != nil {
		erc.Properties.ServiceProviderProperties.ServiceProviderName = pointer.To(d.Get("service_provider_name").(string))
		erc.Properties.ServiceProviderProperties.PeeringLocation = pointer.To(d.Get("peering_location").(string))
		erc.Properties.ServiceProviderProperties.BandwidthInMbps = pointer.To(int64(d.Get("bandwidth_in_mbps").(int)))
	} else {
		erc.Properties.ExpressRoutePort.Id = pointer.To(d.Get("express_route_port_id").(string))
		erc.Properties.BandwidthInGbps = pointer.To(d.Get("bandwidth_in_gbps").(float64))
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, erc, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	// API has bug, which appears to be eventually consistent on creation. Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/10148
	log.Printf("[DEBUG] Waiting for %s to be able to be queried", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"NotFound", "Waiting"},
		Target:                    []string{"Exists"},
		Refresh:                   expressRouteCircuitCreationRefreshFunc(ctx, client, id),
		PollInterval:              3 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("for %s to be able to be queried: %+v", id, err)
	}

	//  authorization_key can only be set after Circuit is created
	if erc.Properties.AuthorizationKey != nil && *erc.Properties.AuthorizationKey != "" {
		if err := client.CreateOrUpdateThenPoll(ctx, id, erc); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
	}

	return resourceExpressRouteCircuitRead(d, meta)
}

func expandExpressRouteCircuitSku(input []any) *expressroutecircuits.ExpressRouteCircuitSku {
	v := input[0].(map[string]any) // [0] is guarded by MinItems in pluginsdk.
	tier := v["tier"].(string)
	family := v["family"].(string)

	return &expressroutecircuits.ExpressRouteCircuitSku{
		Name:   pointer.To(fmt.Sprintf("%s_%s", tier, family)),
		Tier:   pointer.ToEnum[expressroutecircuits.ExpressRouteCircuitSkuTier](tier),
		Family: pointer.ToEnum[expressroutecircuits.ExpressRouteCircuitSkuFamily](family),
	}
}
