package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const virtualHubResourceName = "azurerm_virtual_hub"

func resourceVirtualHub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubCreateUpdate,
		Read:   resourceVirtualHubRead,
		Update: resourceVirtualHubCreateUpdate,
		Delete: resourceVirtualHubDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.VirtualHubName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"address_prefix": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.CIDR,
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Basic",
					"Standard",
				}, false),
			},

			"virtual_wan_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"route": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"address_prefixes": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.CIDR,
							},
						},
						"next_hop_ip_address": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.IPv4Address,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceVirtualHubCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if _, ok := ctx.Deadline(); !ok {
		return fmt.Errorf("deadline is not properly set for Virtual Hub")
	}

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(name, virtualHubResourceName)
	defer locks.UnlockByName(name, virtualHubResourceName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_hub", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	route := d.Get("route").(*pluginsdk.Set).List()
	t := d.Get("tags").(map[string]interface{})

	parameters := network.VirtualHub{
		Location: utils.String(location),
		VirtualHubProperties: &network.VirtualHubProperties{
			RouteTable: expandVirtualHubRoute(route),
		},
		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("address_prefix"); ok {
		parameters.VirtualHubProperties.AddressPrefix = utils.String(v.(string))
	}

	if v, ok := d.GetOk("sku"); ok {
		parameters.VirtualHubProperties.Sku = utils.String(v.(string))
	}

	if v, ok := d.GetOk("virtual_wan_id"); ok {
		parameters.VirtualHubProperties.VirtualWan = &network.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// Hub returns provisioned while the routing state is still "provisining". This might cause issues with following hubvnet connection operations.
	// https://github.com/Azure/azure-rest-api-specs/issues/10391
	// As a workaround, we will poll the routing state and ensure it is "Provisioned".

	// deadline is checked at the entry point of this function
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Provisioning"},
		Target:                    []string{"Provisioned", "Failed", "None"},
		Refresh:                   virtualHubCreateRefreshFunc(ctx, client, resourceGroup, name),
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   time.Until(timeout),
	}
	respRaw, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for Virtual Hub %q (Host Group Name %q) provisioning route: %+v", name, resourceGroup, err)
	}

	resp := respRaw.(network.VirtualHub)
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Virtual Hub %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceVirtualHubRead(d, meta)
}

func resourceVirtualHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualHubID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Virtual Hub %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Virtual Hub %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.VirtualHubProperties; props != nil {
		d.Set("address_prefix", props.AddressPrefix)
		d.Set("sku", props.Sku)

		if err := d.Set("route", flattenVirtualHubRoute(props.RouteTable)); err != nil {
			return fmt.Errorf("Error setting `route`: %+v", err)
		}

		var virtualWanId *string
		if props.VirtualWan != nil {
			virtualWanId = props.VirtualWan.ID
		}
		d.Set("virtual_wan_id", virtualWanId)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVirtualHubDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualHubID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, virtualHubResourceName)
	defer locks.UnlockByName(id.Name, virtualHubResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Virtual Hub %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Virtual Hub %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandVirtualHubRoute(input []interface{}) *network.VirtualHubRouteTable {
	if len(input) == 0 {
		return nil
	}

	results := make([]network.VirtualHubRoute, 0)
	for _, item := range input {
		if item == nil {
			continue
		}

		v := item.(map[string]interface{})
		addressPrefixes := v["address_prefixes"].([]interface{})
		nextHopIpAddress := v["next_hop_ip_address"].(string)

		results = append(results, network.VirtualHubRoute{
			AddressPrefixes:  utils.ExpandStringSlice(addressPrefixes),
			NextHopIPAddress: utils.String(nextHopIpAddress),
		})
	}

	result := network.VirtualHubRouteTable{
		Routes: &results,
	}

	return &result
}

func flattenVirtualHubRoute(input *network.VirtualHubRouteTable) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || input.Routes == nil {
		return results
	}

	for _, item := range *input.Routes {
		addressPrefixes := utils.FlattenStringSlice(item.AddressPrefixes)
		nextHopIpAddress := ""

		if item.NextHopIPAddress != nil {
			nextHopIpAddress = *item.NextHopIPAddress
		}

		results = append(results, map[string]interface{}{
			"address_prefixes":    addressPrefixes,
			"next_hop_ip_address": nextHopIpAddress,
		})
	}

	return results
}

func virtualHubCreateRefreshFunc(ctx context.Context, client *network.VirtualHubsClient, resourceGroup, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("Virtual Hub %q (Resource Group %q) doesn't exist", resourceGroup, name)
			}

			return nil, "", fmt.Errorf("retrieving Virtual Hub %q (Resource Group %q): %+v", resourceGroup, name, err)
		}
		if res.VirtualHubProperties == nil {
			return nil, "", fmt.Errorf("unexpected nil properties of Virtual Hub %q (Resource Group %q)", resourceGroup, name)
		}

		state := res.VirtualHubProperties.RoutingState
		if state == "Failed" {
			return nil, "", fmt.Errorf("failed to provision routing on Virtual Hub %q (Resource Group %q)", resourceGroup, name)
		}
		return res, string(res.VirtualHubProperties.RoutingState), nil
	}
}
