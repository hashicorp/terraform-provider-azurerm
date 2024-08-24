// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/localnetworkgateways"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLocalNetworkGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLocalNetworkGatewayCreate,
		Read:   resourceLocalNetworkGatewayRead,
		Update: resourceLocalNetworkGatewayUpdate,
		Delete: resourceLocalNetworkGatewayDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := localnetworkgateways.ParseLocalNetworkGatewayID(id)
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
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceLocalNetworkGatewayCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.LocalNetworkGateways
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := localnetworkgateways.NewLocalNetworkGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_local_network_gateway", id.ID())
	}

	gateway := localnetworkgateways.LocalNetworkGateway{
		Name:     pointer.To(id.LocalNetworkGatewayName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: localnetworkgateways.LocalNetworkGatewayPropertiesFormat{
			LocalNetworkAddressSpace: &localnetworkgateways.AddressSpace{},
			BgpSettings:              expandLocalNetworkGatewayBGPSettings(d),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	ipAddress := d.Get("gateway_address").(string)
	fqdn := d.Get("gateway_fqdn").(string)
	if ipAddress != "" {
		gateway.Properties.GatewayIPAddress = &ipAddress
	} else {
		gateway.Properties.Fqdn = &fqdn
	}

	// This custompoller can be removed once https://github.com/hashicorp/go-azure-sdk/issues/989 has been fixed
	pollerType := custompollers.NewLocalNetworkGatewayPoller(client, id)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)

	gateway.Properties.LocalNetworkAddressSpace = expandLocalNetworkGatewayAddressSpaces(d)

	if _, err := client.CreateOrUpdate(ctx, id, gateway); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceLocalNetworkGatewayRead(d, meta)
}

func resourceLocalNetworkGatewayUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.LocalNetworkGateways
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := localnetworkgateways.ParseLocalNetworkGatewayID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("gateway_address") {
		payload.Properties.GatewayIPAddress = pointer.To(d.Get("gateway_address").(string))
	}

	if d.HasChange("gateway_fqdn") {
		payload.Properties.Fqdn = pointer.To(d.Get("gateway_fqdn").(string))
	}

	if d.HasChange("bgp_settings") {
		payload.Properties.BgpSettings = expandLocalNetworkGatewayBGPSettings(d)
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	// This custompoller can be removed once https://github.com/hashicorp/go-azure-sdk/issues/989 has been fixed
	pollerType := custompollers.NewLocalNetworkGatewayPoller(client, *id)
	poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)

	// There is a bug in the provider where the address space ordering doesn't change as expected.
	// In the UI we have to remove the current list of addresses in the address space and re-add them in the new order and we'll copy that here.
	if d.HasChange("address_space") {
		// since the local network gateway cannot have both empty address prefix and empty BGP setting(confirmed with service team, it is by design),
		// replace the empty address prefix with the first address prefix in the "address_space" list to avoid error.
		if v := d.Get("address_space").([]interface{}); len(v) > 0 {
			payload.Properties.LocalNetworkAddressSpace = &localnetworkgateways.AddressSpace{
				AddressPrefixes: &[]string{v[0].(string)},
			}
		}

		// This can be switched back over to CreateOrUpdateThenPoll once https://github.com/hashicorp/go-azure-sdk/issues/989 has been fixed
		if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
			return fmt.Errorf("removing %s: %+v", id, err)
		}
		if err := poller.PollUntilDone(ctx); err != nil {
			return err
		}
	}

	payload.Properties.LocalNetworkAddressSpace = expandLocalNetworkGatewayAddressSpaces(d)

	if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err := poller.PollUntilDone(ctx); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceLocalNetworkGatewayRead(d, meta)
}

func resourceLocalNetworkGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.LocalNetworkGateways
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := localnetworkgateways.ParseLocalNetworkGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.LocalNetworkGatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		props := model.Properties
		d.Set("gateway_address", props.GatewayIPAddress)
		d.Set("gateway_fqdn", props.Fqdn)

		if lnas := props.LocalNetworkAddressSpace; lnas != nil {
			d.Set("address_space", lnas.AddressPrefixes)
		}
		flattenedSettings := flattenLocalNetworkGatewayBGPSettings(props.BgpSettings)
		if err := d.Set("bgp_settings", flattenedSettings); err != nil {
			return err
		}

		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceLocalNetworkGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.LocalNetworkGateways
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := localnetworkgateways.ParseLocalNetworkGatewayID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func localNetworkGatewayFromId(localNetworkGatewayId string) (string, error) {
	id, err := localnetworkgateways.ParseLocalNetworkGatewayID(localNetworkGatewayId)
	if err != nil {
		return "", err
	}

	return id.LocalNetworkGatewayName, nil
}

func expandLocalNetworkGatewayBGPSettings(d *pluginsdk.ResourceData) *localnetworkgateways.BgpSettings {
	v, exists := d.GetOk("bgp_settings")
	if !exists {
		return nil
	}

	settings := v.([]interface{})
	setting := settings[0].(map[string]interface{})

	bgpSettings := localnetworkgateways.BgpSettings{
		Asn:               pointer.To(int64(setting["asn"].(int))),
		BgpPeeringAddress: pointer.To(setting["bgp_peering_address"].(string)),
		PeerWeight:        pointer.To(int64(setting["peer_weight"].(int))),
	}

	return &bgpSettings
}

func expandLocalNetworkGatewayAddressSpaces(d *pluginsdk.ResourceData) *localnetworkgateways.AddressSpace {
	prefixes := make([]string, 0)

	for _, pref := range d.Get("address_space").([]interface{}) {
		prefixes = append(prefixes, pref.(string))
	}

	return &localnetworkgateways.AddressSpace{
		AddressPrefixes: &prefixes,
	}
}

func flattenLocalNetworkGatewayBGPSettings(input *localnetworkgateways.BgpSettings) []interface{} {
	output := make(map[string]interface{})

	if input == nil {
		return []interface{}{}
	}

	output["asn"] = int(*input.Asn)
	output["bgp_peering_address"] = *input.BgpPeeringAddress
	output["peer_weight"] = int(*input.PeerWeight)

	return []interface{}{output}
}
