// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceLocalNetworkGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLocalNetworkGatewayCreateUpdate,
		Read:   resourceLocalNetworkGatewayRead,
		Update: resourceLocalNetworkGatewayCreateUpdate,
		Delete: resourceLocalNetworkGatewayDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LocalNetworkGatewayID(id)
			return err
		}),

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

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"gateway_address": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"gateway_address", "gateway_fqdn"},
			},

			"gateway_fqdn": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"gateway_address", "gateway_fqdn"},
			},

			"address_space": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"bgp_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"asn": {
							Type:     pluginsdk.TypeInt,
							Required: true,
						},

						"bgp_peering_address": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},

						"peer_weight": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceLocalNetworkGatewayCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.LocalNetworkGatewaysClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewLocalNetworkGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_local_network_gateway", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	gateway := network.LocalNetworkGateway{
		Name:     &id.Name,
		Location: &location,
		LocalNetworkGatewayPropertiesFormat: &network.LocalNetworkGatewayPropertiesFormat{
			LocalNetworkAddressSpace: &network.AddressSpace{},
			BgpSettings:              expandLocalNetworkGatewayBGPSettings(d),
		},
		Tags: tags.Expand(t),
	}

	ipAddress := d.Get("gateway_address").(string)
	fqdn := d.Get("gateway_fqdn").(string)
	if ipAddress != "" {
		gateway.LocalNetworkGatewayPropertiesFormat.GatewayIPAddress = &ipAddress
	} else {
		gateway.LocalNetworkGatewayPropertiesFormat.Fqdn = &fqdn
	}

	// There is a bug in the provider where the address space ordering doesn't change as expected.
	// In the UI we have to remove the current list of addresses in the address space and re-add them in the new order and we'll copy that here.
	if !d.IsNewResource() && d.HasChange("address_space") {
		// since the local network gateway cannot have both empty address prefix and empty BGP setting(confirmed with service team, it is by design),
		// replace the empty address prefix with the first address prefix in the "address_space" list to avoid error.
		if v := d.Get("address_space").([]interface{}); len(v) > 0 {
			gateway.LocalNetworkGatewayPropertiesFormat.LocalNetworkAddressSpace = &network.AddressSpace{
				AddressPrefixes: &[]string{v[0].(string)},
			}
		}

		future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, gateway)
		if err != nil {
			return fmt.Errorf("removing %s: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for completion of %s: %+v", id, err)
		}
	}
	gateway.LocalNetworkGatewayPropertiesFormat.LocalNetworkAddressSpace = expandLocalNetworkGatewayAddressSpaces(d)

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, gateway)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceLocalNetworkGatewayRead(d, meta)
}

func resourceLocalNetworkGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.LocalNetworkGatewaysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup, name, err := resourceGroupAndLocalNetworkGatewayFromId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading the state of Local Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.LocalNetworkGatewayPropertiesFormat; props != nil {
		d.Set("gateway_address", props.GatewayIPAddress)
		d.Set("gateway_fqdn", props.Fqdn)

		if lnas := props.LocalNetworkAddressSpace; lnas != nil {
			d.Set("address_space", lnas.AddressPrefixes)
		}
		flattenedSettings := flattenLocalNetworkGatewayBGPSettings(props.BgpSettings)
		if err := d.Set("bgp_settings", flattenedSettings); err != nil {
			return err
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceLocalNetworkGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.LocalNetworkGatewaysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup, name, err := resourceGroupAndLocalNetworkGatewayFromId(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("issuing delete request for local network gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of local network gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func resourceGroupAndLocalNetworkGatewayFromId(localNetworkGatewayId string) (string, string, error) {
	id, err := parse.LocalNetworkGatewayID(localNetworkGatewayId)
	if err != nil {
		return "", "", err
	}

	return id.ResourceGroup, id.Name, nil
}

func expandLocalNetworkGatewayBGPSettings(d *pluginsdk.ResourceData) *network.BgpSettings {
	v, exists := d.GetOk("bgp_settings")
	if !exists {
		return nil
	}

	settings := v.([]interface{})
	setting := settings[0].(map[string]interface{})

	bgpSettings := network.BgpSettings{
		Asn:               utils.Int64(int64(setting["asn"].(int))),
		BgpPeeringAddress: utils.String(setting["bgp_peering_address"].(string)),
		PeerWeight:        utils.Int32(int32(setting["peer_weight"].(int))),
	}

	return &bgpSettings
}

func expandLocalNetworkGatewayAddressSpaces(d *pluginsdk.ResourceData) *network.AddressSpace {
	prefixes := make([]string, 0)

	for _, pref := range d.Get("address_space").([]interface{}) {
		prefixes = append(prefixes, pref.(string))
	}

	return &network.AddressSpace{
		AddressPrefixes: &prefixes,
	}
}

func flattenLocalNetworkGatewayBGPSettings(input *network.BgpSettings) []interface{} {
	output := make(map[string]interface{})

	if input == nil {
		return []interface{}{}
	}

	output["asn"] = int(*input.Asn)
	output["bgp_peering_address"] = *input.BgpPeeringAddress
	output["peer_weight"] = int(*input.PeerWeight)

	return []interface{}{output}
}
