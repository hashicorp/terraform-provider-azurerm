package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var expressRouteCircuitResourceName = "azurerm_express_route_circuit"

func resourceExpressRouteCircuit() *schema.Resource {
	return &schema.Resource{
		Create: resourceExpressRouteCircuitCreateUpdate,
		Read:   resourceExpressRouteCircuitRead,
		Update: resourceExpressRouteCircuitCreateUpdate,
		Delete: resourceExpressRouteCircuitDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: customdiff.Sequence(
			// If bandwidth is reduced force a new resource
			customdiff.ForceNewIfChange("bandwidth_in_mbps", func(old, new, meta interface{}) bool {
				return new.(int) < old.(int)
			}),
		),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"service_provider_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"peering_location": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"bandwidth_in_mbps": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tier": {
							Type:     schema.TypeString,
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
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.MeteredData),
								string(network.UnlimitedData),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},

			"allow_classic_operations": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"service_provider_provisioning_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceExpressRouteCircuitCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM ExpressRoute Circuit creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	locks.ByName(name, expressRouteCircuitResourceName)
	defer locks.UnlockByName(name, expressRouteCircuitResourceName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing ExpressRoute Circuit %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_express_route_circuit", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	serviceProviderName := d.Get("service_provider_name").(string)
	peeringLocation := d.Get("peering_location").(string)
	bandwidthInMbps := int32(d.Get("bandwidth_in_mbps").(int))
	sku := expandExpressRouteCircuitSku(d)
	allowRdfeOps := d.Get("allow_classic_operations").(bool)
	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	// There is the potential for the express route circuit to become out of sync when the service provider updates
	// the express route circuit. We'll get and update the resource in place as per https://aka.ms/erRefresh
	// We also want to keep track of the resource obtained from the api and pass down any attributes not
	// managed by Terraform.
	erc := network.ExpressRouteCircuit{}
	if !d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(erc.Response) {
				return fmt.Errorf("Error checking for presence of existing ExpressRoute Circuit %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		future, err := client.CreateOrUpdate(ctx, resGroup, name, existing)
		if err != nil {
			return fmt.Errorf("Error Creating/Updating ExpressRouteCircuit %q (Resource Group %q): %+v", name, resGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error Creating/Updating ExpressRouteCircuit %q (Resource Group %q): %+v", name, resGroup, err)
		}
		erc = existing
	}

	erc.Name = &name
	erc.Location = &location
	erc.Sku = sku
	erc.Tags = expandedTags

	if erc.ExpressRouteCircuitPropertiesFormat != nil {
		erc.ExpressRouteCircuitPropertiesFormat.AllowClassicOperations = &allowRdfeOps
		if erc.ExpressRouteCircuitPropertiesFormat.ServiceProviderProperties != nil {
			erc.ExpressRouteCircuitPropertiesFormat.ServiceProviderProperties.ServiceProviderName = &serviceProviderName
			erc.ExpressRouteCircuitPropertiesFormat.ServiceProviderProperties.PeeringLocation = &peeringLocation
			erc.ExpressRouteCircuitPropertiesFormat.ServiceProviderProperties.BandwidthInMbps = &bandwidthInMbps
		}
	} else {
		erc.ExpressRouteCircuitPropertiesFormat = &network.ExpressRouteCircuitPropertiesFormat{
			AllowClassicOperations: &allowRdfeOps,
			ServiceProviderProperties: &network.ExpressRouteCircuitServiceProviderProperties{
				ServiceProviderName: &serviceProviderName,
				PeeringLocation:     &peeringLocation,
				BandwidthInMbps:     &bandwidthInMbps,
			},
		}
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, erc)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating ExpressRouteCircuit %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error Creating/Updating ExpressRouteCircuit %q (Resource Group %q): %+v", name, resGroup, err)
	}

	// API has bug, which appears to be eventually consistent on creation. Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/10148
	log.Printf("[DEBUG] Waiting for Express Route Circuit %q (Resource Group %q) to be able to be queried", name, resGroup)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"NotFound"},
		Target:                    []string{"Exists"},
		Refresh:                   expressRouteCircuitCreationRefreshFunc(ctx, client, resGroup, name),
		PollInterval:              3 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   d.Timeout(schema.TimeoutCreate),
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error for Express Route Circuit %q (Resource Group %q) to be able to be queried: %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error Retrieving ExpressRouteCircuit %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ExpressRouteCircuit %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceExpressRouteCircuitRead(d, meta)
}

func resourceExpressRouteCircuitRead(d *schema.ResourceData, meta interface{}) error {
	ercClient := meta.(*clients.Client).Network.ExpressRouteCircuitsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error Parsing Azure Resource ID -: %+v", err)
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["expressRouteCircuits"]

	resp, err := ercClient.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Express Route Circuit %q (Resource Group %q) was not found - removing from state", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Express Route Circuit %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if resp.Sku != nil {
		sku := flattenExpressRouteCircuitSku(resp.Sku)
		if err := d.Set("sku", sku); err != nil {
			return fmt.Errorf("Error setting `sku`: %+v", err)
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

func resourceExpressRouteCircuitDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteCircuitsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error Parsing Azure Resource ID -: %+v", err)
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["expressRouteCircuits"]

	locks.ByName(name, expressRouteCircuitResourceName)
	defer locks.UnlockByName(name, expressRouteCircuitResourceName)

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}

func expandExpressRouteCircuitSku(d *schema.ResourceData) *network.ExpressRouteCircuitSku {
	skuSettings := d.Get("sku").([]interface{})
	v := skuSettings[0].(map[string]interface{}) // [0] is guarded by MinItems in schema.
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

func expressRouteCircuitCreationRefreshFunc(ctx context.Context, client *network.ExpressRouteCircuitsClient, resGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return nil, "NotFound", nil
			}

			return nil, "", fmt.Errorf("Error polling to check if the Express Route Circuit has been created: %+v", err)
		}

		return res, "Exists", nil
	}
}
