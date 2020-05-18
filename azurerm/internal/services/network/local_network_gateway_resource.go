package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLocalNetworkGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLocalNetworkGatewayCreateUpdate,
		Read:   resourceArmLocalNetworkGatewayRead,
		Update: resourceArmLocalNetworkGatewayCreateUpdate,
		Delete: resourceArmLocalNetworkGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"gateway_address": {
				Type:     schema.TypeString,
				Required: true,
			},

			"address_space": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"bgp_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"bgp_peering_address": {
							Type:     schema.TypeString,
							Required: true,
						},

						"peer_weight": {
							Type:     schema.TypeInt,
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

func resourceArmLocalNetworkGatewayCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.LocalNetworkGatewaysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Local Network Gateway %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_local_network_gateway", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	ipAddress := d.Get("gateway_address").(string)

	addressSpaces := expandLocalNetworkGatewayAddressSpaces(d)

	t := d.Get("tags").(map[string]interface{})

	gateway := network.LocalNetworkGateway{
		Name:     &name,
		Location: &location,
		LocalNetworkGatewayPropertiesFormat: &network.LocalNetworkGatewayPropertiesFormat{
			LocalNetworkAddressSpace: &network.AddressSpace{
				AddressPrefixes: &addressSpaces,
			},
			GatewayIPAddress: &ipAddress,
			BgpSettings:      expandLocalNetworkGatewayBGPSettings(d),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, gateway)
	if err != nil {
		return fmt.Errorf("Error creating Local Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Local Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Local Network Gateway ID %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmLocalNetworkGatewayRead(d, meta)
}

func resourceArmLocalNetworkGatewayRead(d *schema.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("Error reading the state of Local Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.LocalNetworkGatewayPropertiesFormat; props != nil {
		d.Set("gateway_address", props.GatewayIPAddress)

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

func resourceArmLocalNetworkGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.LocalNetworkGatewaysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup, name, err := resourceGroupAndLocalNetworkGatewayFromId(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error issuing delete request for local network gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error waiting for completion of local network gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func resourceGroupAndLocalNetworkGatewayFromId(localNetworkGatewayId string) (string, string, error) {
	id, err := azure.ParseAzureResourceID(localNetworkGatewayId)
	if err != nil {
		return "", "", err
	}
	name := id.Path["localNetworkGateways"]
	resGroup := id.ResourceGroup

	return resGroup, name, nil
}

func expandLocalNetworkGatewayBGPSettings(d *schema.ResourceData) *network.BgpSettings {
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

func expandLocalNetworkGatewayAddressSpaces(d *schema.ResourceData) []string {
	prefixes := make([]string, 0)

	for _, pref := range d.Get("address_space").([]interface{}) {
		prefixes = append(prefixes, pref.(string))
	}

	return prefixes
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
