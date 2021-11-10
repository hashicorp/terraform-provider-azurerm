package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const virtualHubResourceName = "azurerm_virtual_hub"

func resourceVirtualHub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubCreateUpdate,
		Read:   resourceVirtualHubRead,
		Update: resourceVirtualHubCreateUpdate,
		Delete: resourceVirtualHubDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualHubID(id)
			return err
		}),

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

			"default_route_table_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVirtualHubCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if _, ok := ctx.Deadline(); !ok {
		return fmt.Errorf("deadline is not properly set for Virtual Hub")
	}

	id := parse.NewVirtualHubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.Name, virtualHubResourceName)
	defer locks.UnlockByName(id.Name, virtualHubResourceName)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
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

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	// Hub returns provisioned while the routing state is still "provisining". This might cause issues with following hubvnet connection operations.
	// https://github.com/Azure/azure-rest-api-specs/issues/10391
	// As a workaround, we will poll the routing state and ensure it is "Provisioned".

	// deadline is checked at the entry point of this function
	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Provisioning"},
		Target:                    []string{"Provisioned", "Failed", "None"},
		Refresh:                   virtualHubCreateRefreshFunc(ctx, client, id.ResourceGroup, id.Name),
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 3,
		Timeout:                   time.Until(timeout),
	}
	respRaw, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for %s provisioning route: %+v", id, err)
	}

	resp := respRaw.(network.VirtualHub)
	if resp.ID == nil {
		return fmt.Errorf("cannot read %s ID", id)
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
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.VirtualHubProperties; props != nil {
		d.Set("address_prefix", props.AddressPrefix)
		d.Set("sku", props.Sku)

		if err := d.Set("route", flattenVirtualHubRoute(props.RouteTable)); err != nil {
			return fmt.Errorf("setting `route`: %+v", err)
		}

		var virtualWanId *string
		if props.VirtualWan != nil {
			virtualWanId = props.VirtualWan.ID
		}
		d.Set("virtual_wan_id", virtualWanId)
	}

	defaultRouteTable := parse.NewHubRouteTableID(id.SubscriptionId, id.ResourceGroup, id.Name, "defaultRouteTable")
	d.Set("default_route_table_id", defaultRouteTable.ID())

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
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting %s: %+v", *id, err)
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
