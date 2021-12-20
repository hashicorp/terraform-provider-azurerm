package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var expressRouteCircuitResourceName = "azurerm_express_route_circuit"

func resourceExpressRouteCircuit() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceExpressRouteCircuitCreateUpdate,
		Read:   resourceExpressRouteCircuitRead,
		Update: resourceExpressRouteCircuitCreateUpdate,
		Delete: resourceExpressRouteCircuitDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ExpressRouteCircuitID(id)
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

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
								string(network.ExpressRouteCircuitSkuTierBasic),
								string(network.ExpressRouteCircuitSkuTierLocal),
								string(network.ExpressRouteCircuitSkuTierStandard),
								string(network.ExpressRouteCircuitSkuTierPremium),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"family": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.ExpressRouteCircuitSkuFamilyMeteredData),
								string(network.ExpressRouteCircuitSkuFamilyUnlimitedData),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
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
				ValidateFunc:  validate.ExpressRoutePortID,
			},

			"service_provider_provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"service_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceExpressRouteCircuitCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM ExpressRoute Circuit creation.")

	id := parse.NewExpressRouteCircuitID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.Name, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.Name, expressRouteCircuitResourceName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s : %s", id, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_express_route_circuit", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := expandExpressRouteCircuitSku(d)
	t := d.Get("tags").(map[string]interface{})
	allowRdfeOps := d.Get("allow_classic_operations").(bool)
	expandedTags := tags.Expand(t)

	// There is the potential for the express route circuit to become out of sync when the service provider updates
	// the express route circuit. We'll get and update the resource in place as per https://aka.ms/erRefresh
	// We also want to keep track of the resource obtained from the api and pass down any attributes not
	// managed by Terraform.
	erc := network.ExpressRouteCircuit{}
	if !d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(erc.Response) {
				return fmt.Errorf("checking for presence of existing %s : %s", id, err)
			}
		}

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, existing)
		if err != nil {
			return fmt.Errorf("Creating/Updating %s : %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Creating/Updating %s: %+v", id, err)
		}
		erc = existing
	}

	erc.Name = &id.Name
	erc.Location = &location
	erc.Sku = sku
	erc.Tags = expandedTags

	if !d.IsNewResource() {
		erc.ExpressRouteCircuitPropertiesFormat.AllowClassicOperations = &allowRdfeOps
	} else {
		erc.ExpressRouteCircuitPropertiesFormat = &network.ExpressRouteCircuitPropertiesFormat{}

		// ServiceProviderProperties and expressRoutePorts/bandwidthInGbps properties are mutually exclusive
		if _, ok := d.GetOk("express_route_port_id"); ok {
			erc.ExpressRouteCircuitPropertiesFormat.ExpressRoutePort = &network.SubResource{}
		} else {
			erc.ExpressRouteCircuitPropertiesFormat.ServiceProviderProperties = &network.ExpressRouteCircuitServiceProviderProperties{}
		}
	}

	if erc.ExpressRouteCircuitPropertiesFormat.ServiceProviderProperties != nil {
		erc.ExpressRouteCircuitPropertiesFormat.ServiceProviderProperties.ServiceProviderName = utils.String(d.Get("service_provider_name").(string))
		erc.ExpressRouteCircuitPropertiesFormat.ServiceProviderProperties.PeeringLocation = utils.String(d.Get("peering_location").(string))
		erc.ExpressRouteCircuitPropertiesFormat.ServiceProviderProperties.BandwidthInMbps = utils.Int32(int32(d.Get("bandwidth_in_mbps").(int)))
	} else {
		erc.ExpressRouteCircuitPropertiesFormat.ExpressRoutePort.ID = utils.String(d.Get("express_route_port_id").(string))
		erc.ExpressRouteCircuitPropertiesFormat.BandwidthInGbps = utils.Float(d.Get("bandwidth_in_gbps").(float64))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, erc)
	if err != nil {
		return fmt.Errorf("Creating/Updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Creating/Updating %s: %+v", id, err)
	}

	// API has bug, which appears to be eventually consistent on creation. Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/10148
	log.Printf("[DEBUG] Waiting for %s to be able to be queried", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"NotFound"},
		Target:                    []string{"Exists"},
		Refresh:                   expressRouteCircuitCreationRefreshFunc(ctx, client, id.ResourceGroup, id.Name),
		PollInterval:              3 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("for %s to be able to be queried: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceExpressRouteCircuitRead(d, meta)
}

func resourceExpressRouteCircuitRead(d *pluginsdk.ResourceData, meta interface{}) error {
	ercClient := meta.(*clients.Client).Network.ExpressRouteCircuitsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteCircuitID(d.Id())
	if err != nil {
		return fmt.Errorf("Parsing Azure Resource ID -: %+v", err)
	}

	resp, err := ercClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if resp.Sku != nil {
		sku := flattenExpressRouteCircuitSku(resp.Sku)
		if err := d.Set("sku", sku); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}
	}

	if resp.ExpressRoutePort != nil {
		d.Set("bandwidth_in_gbps", resp.BandwidthInGbps)

		if resp.ExpressRoutePort.ID != nil {
			portID, err := parse.ExpressRoutePortID(*resp.ExpressRoutePort.ID)
			if err != nil {
				return err
			}
			d.Set("express_route_port_id", portID.ID())
		}
	}

	if props := resp.ServiceProviderProperties; props != nil {
		d.Set("service_provider_name", props.ServiceProviderName)
		d.Set("peering_location", props.PeeringLocation)
		d.Set("bandwidth_in_mbps", props.BandwidthInMbps)
	}

	d.Set("service_provider_provisioning_state", string(resp.ServiceProviderProvisioningState))
	d.Set("service_key", resp.ServiceKey)
	d.Set("allow_classic_operations", resp.AllowClassicOperations)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceExpressRouteCircuitDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteCircuitID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Azure Resource ID -: %+v", err)
	}

	locks.ByName(id.Name, expressRouteCircuitResourceName)
	defer locks.UnlockByName(id.Name, expressRouteCircuitResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s : %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}

func expandExpressRouteCircuitSku(d *pluginsdk.ResourceData) *network.ExpressRouteCircuitSku {
	skuSettings := d.Get("sku").([]interface{})
	v := skuSettings[0].(map[string]interface{}) // [0] is guarded by MinItems in pluginsdk.
	tier := v["tier"].(string)
	family := v["family"].(string)
	name := fmt.Sprintf("%s_%s", tier, family)

	return &network.ExpressRouteCircuitSku{
		Name:   &name,
		Tier:   network.ExpressRouteCircuitSkuTier(tier),
		Family: network.ExpressRouteCircuitSkuFamily(family),
	}
}

func flattenExpressRouteCircuitSku(sku *network.ExpressRouteCircuitSku) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"tier":   string(sku.Tier),
			"family": string(sku.Family),
		},
	}
}

func expressRouteCircuitCreationRefreshFunc(ctx context.Context, client *network.ExpressRouteCircuitsClient, resGroup, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return nil, "NotFound", nil
			}

			return nil, "", fmt.Errorf("polling to check if the Express Route Circuit has been created: %+v", err)
		}

		return res, "Exists", nil
	}
}
