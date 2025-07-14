// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/expressrouteports"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/expressroutecircuits"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var expressRouteCircuitResourceName = "azurerm_express_route_circuit"

func resourceExpressRouteCircuit() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRouteCircuitCreate,
		Read:   resourceExpressRouteCircuitRead,
		Update: resourceExpressRouteCircuitUpdate,
		Delete: resourceExpressRouteCircuitDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := expressroutecircuits.ParseExpressRouteCircuitID(id)
			return err
		}),

		CustomizeDiff: pluginsdk.CustomDiffInSequence(
			// If bandwidth is reduced force a new resource
			pluginsdk.ForceNewIfChange("bandwidth_in_mbps", func(ctx context.Context, old, new, meta interface{}) bool {
				return new.(int) < old.(int)
			}),
		),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"sku": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"tier": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(expressroutecircuits.ExpressRouteCircuitSkuTierBasic),
								string(expressroutecircuits.ExpressRouteCircuitSkuTierLocal),
								string(expressroutecircuits.ExpressRouteCircuitSkuTierStandard),
								string(expressroutecircuits.ExpressRouteCircuitSkuTierPremium),
							}, false),
						},

						"family": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(expressroutecircuits.ExpressRouteCircuitSkuFamilyMeteredData),
								string(expressroutecircuits.ExpressRouteCircuitSkuFamilyUnlimitedData),
							}, false),
						},
					},
				},
			},

			"allow_classic_operations": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"service_provider_name": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				RequiredWith:     []string{"bandwidth_in_mbps", "peering_location"},
				ConflictsWith:    []string{"bandwidth_in_gbps", "express_route_port_id"},
			},

			"peering_location": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				RequiredWith:     []string{"bandwidth_in_mbps", "service_provider_name"},
				ConflictsWith:    []string{"bandwidth_in_gbps", "express_route_port_id"},
			},

			"bandwidth_in_mbps": {
				Type:          pluginsdk.TypeInt,
				Optional:      true,
				RequiredWith:  []string{"peering_location", "service_provider_name"},
				ConflictsWith: []string{"bandwidth_in_gbps", "express_route_port_id"},
			},

			"bandwidth_in_gbps": {
				Type:          pluginsdk.TypeFloat,
				Optional:      true,
				RequiredWith:  []string{"express_route_port_id"},
				ConflictsWith: []string{"bandwidth_in_mbps", "peering_location", "service_provider_name"},
			},

			"express_route_port_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				RequiredWith:  []string{"bandwidth_in_gbps"},
				ConflictsWith: []string{"bandwidth_in_mbps", "peering_location", "service_provider_name"},
				ValidateFunc:  expressrouteports.ValidateExpressRoutePortID,
			},

			"rate_limiting_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"service_provider_provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"authorization_key": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"service_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceExpressRouteCircuitCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuits
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM ExpressRoute Circuit creation.")

	id := expressroutecircuits.NewExpressRouteCircuitID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s : %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_express_route_circuit", id.ID())
	}

	erc := expressroutecircuits.ExpressRouteCircuit{
		Name:     &id.ExpressRouteCircuitName,
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Sku:      expandExpressRouteCircuitSku(d.Get("sku").([]interface{})),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
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
		erc.Properties.BandwidthInGbps = utils.Float(d.Get("bandwidth_in_gbps").(float64))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, erc); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// API has bug, which appears to be eventually consistent on creation. Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/10148
	log.Printf("[DEBUG] Waiting for %s to be able to be queried", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"NotFound"},
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
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceExpressRouteCircuitRead(d, meta)
}

func resourceExpressRouteCircuitUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuits
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM ExpressRoute Circuit update.")

	id, err := expressroutecircuits.ParseExpressRouteCircuitID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)

	// There is the potential for the express route circuit to become out of sync when the service provider updates
	// the express route circuit. We'll get and update the resource in place as per https://aka.ms/erRefresh
	// We also want to keep track of the resource obtained from the api and pass down any attributes not
	// managed by Terraform.

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s : %s", id, err)
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s : %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	payload := *existing.Model

	if d.HasChange("sku") {
		payload.Sku = expandExpressRouteCircuitSku(d.Get("sku").([]interface{}))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("allow_classic_operations") {
		payload.Properties.AllowClassicOperations = pointer.To(d.Get("allow_classic_operations").(bool))
	}

	if d.HasChange("rate_limiting_enabled") {
		payload.Properties.EnableDirectPortRateLimit = pointer.To(d.Get("rate_limiting_enabled").(bool))
	}

	if d.HasChange("bandwidth_in_gbps") {
		payload.Properties.BandwidthInGbps = pointer.To(d.Get("bandwidth_in_gbps").(float64))
	}

	if d.HasChange("bandwidth_in_mbps") {
		payload.Properties.ServiceProviderProperties.BandwidthInMbps = pointer.To(int64(d.Get("bandwidth_in_mbps").(int)))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// API has bug, which appears to be eventually consistent on creation. Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/10148
	log.Printf("[DEBUG] Waiting for %s to be able to be queried", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"NotFound"},
		Target:                    []string{"Exists"},
		Refresh:                   expressRouteCircuitCreationRefreshFunc(ctx, client, *id),
		PollInterval:              3 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("for %s to be able to be queried: %+v", id, err)
	}

	//  authorization_key can only be set after Circuit is created
	if payload.Properties.AuthorizationKey != nil && *payload.Properties.AuthorizationKey != "" {
		if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceExpressRouteCircuitRead(d, meta)
}

func resourceExpressRouteCircuitRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuits
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressroutecircuits.ParseExpressRouteCircuitID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ExpressRouteCircuitName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		sku := flattenExpressRouteCircuitSku(model.Sku)
		if err := d.Set("sku", sku); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}
		if props := model.Properties; props != nil {
			d.Set("bandwidth_in_gbps", props.BandwidthInGbps)

			if props.ExpressRoutePort != nil && props.ExpressRoutePort.Id != nil {
				portID, err := expressrouteports.ParseExpressRoutePortIDInsensitively(*props.ExpressRoutePort.Id)
				if err != nil {
					return err
				}
				d.Set("express_route_port_id", portID.ID())
			}

			d.Set("service_provider_provisioning_state", string(pointer.From(props.ServiceProviderProvisioningState)))
			d.Set("service_key", props.ServiceKey)
			d.Set("allow_classic_operations", props.AllowClassicOperations)
			d.Set("rate_limiting_enabled", props.EnableDirectPortRateLimit)

			if serviceProviderProps := props.ServiceProviderProperties; serviceProviderProps != nil {
				d.Set("service_provider_name", serviceProviderProps.ServiceProviderName)
				d.Set("peering_location", serviceProviderProps.PeeringLocation)
				d.Set("bandwidth_in_mbps", serviceProviderProps.BandwidthInMbps)
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceExpressRouteCircuitDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuits
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := expressroutecircuits.ParseExpressRouteCircuitID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.ExpressRouteCircuitName, expressRouteCircuitResourceName)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s : %+v", *id, err)
	}

	return err
}

func expandExpressRouteCircuitSku(input []interface{}) *expressroutecircuits.ExpressRouteCircuitSku {
	v := input[0].(map[string]interface{}) // [0] is guarded by MinItems in pluginsdk.
	tier := v["tier"].(string)
	family := v["family"].(string)

	return &expressroutecircuits.ExpressRouteCircuitSku{
		Name:   pointer.To(fmt.Sprintf("%s_%s", tier, family)),
		Tier:   pointer.To(expressroutecircuits.ExpressRouteCircuitSkuTier(tier)),
		Family: pointer.To(expressroutecircuits.ExpressRouteCircuitSkuFamily(family)),
	}
}

func flattenExpressRouteCircuitSku(sku *expressroutecircuits.ExpressRouteCircuitSku) []interface{} {
	if sku == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"tier":   string(pointer.From(sku.Tier)),
			"family": string(pointer.From(sku.Family)),
		},
	}
}

func expressRouteCircuitCreationRefreshFunc(ctx context.Context, client *expressroutecircuits.ExpressRouteCircuitsClient, id expressroutecircuits.ExpressRouteCircuitId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return nil, "NotFound", nil
			}

			return nil, "", fmt.Errorf("polling to check if the Express Route Circuit has been created: %+v", err)
		}

		return res, "Exists", nil
	}
}
