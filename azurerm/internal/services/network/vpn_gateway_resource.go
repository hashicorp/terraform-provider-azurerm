package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	commonValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVPNGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVPNGatewayCreate,
		Read:   resourceArmVPNGatewayRead,
		Update: resourceArmVPNGatewayUpdate,
		Delete: resourceArmVPNGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"virtual_hub_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateVirtualHubID,
			},

			"bgp_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},

						"peer_weight": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},

						"bgp_peering_address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"instance_bgp_peering_address": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MinItems: 2,
							MaxItems: 2,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"custom_ips": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: commonValidate.IPv4Address,
										},
									},

									"ip_configuration_id": {
										Type:     schema.TypeString,
										Computed: true,
									},

									"default_ips": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},

									"tunnel_ips": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},

			"scale_unit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmVPNGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VpnGatewaysClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_vpn_gateway", *existing.ID)
	}

	bgpSettingsRaw := d.Get("bgp_settings").([]interface{})
	bgpSettings := expandVPNGatewayBGPSettings(bgpSettingsRaw)

	location := azure.NormalizeLocation(d.Get("location").(string))
	scaleUnit := d.Get("scale_unit").(int)
	virtualHubId := d.Get("virtual_hub_id").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := network.VpnGateway{
		Location: utils.String(location),
		VpnGatewayProperties: &network.VpnGatewayProperties{
			BgpSettings: bgpSettings,
			VirtualHub: &network.SubResource{
				ID: utils.String(virtualHubId),
			},
			VpnGatewayScaleUnit: utils.Int32(int32(scaleUnit)),
		},
		Tags: tags.Expand(t),
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err := waitForCompletion(d, ctx, client, resourceGroup, name); err != nil {
		return err
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// `vpnGatewayParameters.Properties.bgpSettings.bgpPeeringAddress` customer cannot provide this field during create. This will be set with default value once gateway is created.
	// it could only be updated
	if len(bgpSettingsRaw) > 0 {
		val := bgpSettingsRaw[0].(map[string]interface{})
		input := val["instance_bgp_peering_address"].([]interface{})
		if len(input) > 0 {
			expandVPNGatewayIPConfigurationBgpPeeringAddress(resp.VpnGatewayProperties.BgpSettings.BgpPeeringAddresses, input)
			if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, resp); err != nil {
				return fmt.Errorf("creating VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
			if err := waitForCompletion(d, ctx, client, resourceGroup, name); err != nil {
				return err
			}
		}
	}

	d.SetId(*resp.ID)

	return resourceArmVPNGatewayRead(d, meta)
}

func resourceArmVPNGatewayUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VpnGatewaysClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving for presence of existing VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if d.HasChange("tags") {
		existing.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}
	if d.HasChange("bgp_settings.0.instance_bgp_peering_address") {
		bgpSettingsRaw := d.Get("bgp_settings").([]interface{})
		if len(bgpSettingsRaw) > 0 {
			val := bgpSettingsRaw[0].(map[string]interface{})
			expandVPNGatewayIPConfigurationBgpPeeringAddress(existing.VpnGatewayProperties.BgpSettings.BgpPeeringAddresses, val["instance_bgp_peering_address"].([]interface{}))
		}
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, existing); err != nil {
		return fmt.Errorf("creating VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err := waitForCompletion(d, ctx, client, resourceGroup, name); err != nil {
		return err
	}

	return resourceArmVPNGatewayRead(d, meta)
}

func resourceArmVPNGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VpnGatewaysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseVPNGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] VPN Gateway %q was not found in Resource Group %q - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving VPN Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.VpnGatewayProperties; props != nil {
		if err := d.Set("bgp_settings", flattenVPNGatewayBGPSettings(props.BgpSettings)); err != nil {
			return fmt.Errorf("Error setting `bgp_settings`: %+v", err)
		}

		scaleUnit := 0
		if props.VpnGatewayScaleUnit != nil {
			scaleUnit = int(*props.VpnGatewayScaleUnit)
		}
		d.Set("scale_unit", scaleUnit)

		virtualHubId := ""
		if props.VirtualHub != nil && props.VirtualHub.ID != nil {
			virtualHubId = *props.VirtualHub.ID
		}
		d.Set("virtual_hub_id", virtualHubId)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmVPNGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VpnGatewaysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseVPNGatewayID(d.Id())
	if err != nil {
		return err
	}

	deleteFuture, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(deleteFuture.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting VPN Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	err = deleteFuture.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(deleteFuture.Response()) {
			return nil
		}

		return fmt.Errorf("Error waiting for deletion of VPN Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func waitForCompletion(d *schema.ResourceData, ctx context.Context, client *network.VpnGatewaysClient, resourceGroup, name string) error {
	log.Printf("[DEBUG] Waiting for Virtual Hub %q (Resource Group %q) to become available", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"pending"},
		Target:                    []string{"available"},
		Refresh:                   vpnGatewayWaitForCreatedRefreshFunc(ctx, client, resourceGroup, name),
		Delay:                     30 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 3,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for creation of Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func expandVPNGatewayBGPSettings(input []interface{}) *network.BgpSettings {
	if len(input) == 0 {
		return nil
	}

	val := input[0].(map[string]interface{})
	return &network.BgpSettings{
		Asn:        utils.Int64(int64(val["asn"].(int))),
		PeerWeight: utils.Int32(int32(val["peer_weight"].(int))),
	}
}

func expandVPNGatewayIPConfigurationBgpPeeringAddress(ipConfigurationBgpPeeringAddress *[]network.IPConfigurationBgpPeeringAddress, input []interface{}) error {
	if len(input) == 0 || ipConfigurationBgpPeeringAddress == nil {
		return nil
	}
	if len(input) != 2 || len(*ipConfigurationBgpPeeringAddress) != 2 {
		return fmt.Errorf("the size of block `instance_bgp_peering_address` must be 2")
	}

	for i, v := range input {
		val := v.(map[string]interface{})
		(*ipConfigurationBgpPeeringAddress)[i].CustomBgpIPAddresses = utils.ExpandStringSlice(val["custom_ips"].(*schema.Set).List())
	}
	return nil
}

func flattenVPNGatewayBGPSettings(input *network.BgpSettings) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	asn := 0
	if input.Asn != nil {
		asn = int(*input.Asn)
	}

	bgpPeeringAddress := ""
	if input.BgpPeeringAddress != nil {
		bgpPeeringAddress = *input.BgpPeeringAddress
	}

	peerWeight := 0
	if input.PeerWeight != nil {
		peerWeight = int(*input.PeerWeight)
	}

	return []interface{}{
		map[string]interface{}{
			"asn":                          asn,
			"bgp_peering_address":          bgpPeeringAddress,
			"instance_bgp_peering_address": flattenVPNGatewayIPConfigurationBgpPeeringAddress(input.BgpPeeringAddresses),
			"peer_weight":                  peerWeight,
		},
	}
}

func flattenVPNGatewayIPConfigurationBgpPeeringAddress(input *[]network.IPConfigurationBgpPeeringAddress) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	result := make([]interface{}, 0)
	for _, v := range *input {
		ipConfigurationID := ""
		if v.IpconfigurationID != nil {
			ipConfigurationID = *v.IpconfigurationID
		}

		result = append(result, map[string]interface{}{
			"ip_configuration_id": ipConfigurationID,
			"custom_ips":          utils.FlattenStringSlice(v.CustomBgpIPAddresses),
			"default_ips":         utils.FlattenStringSlice(v.DefaultBgpIPAddresses),
			"tunnel_ips":          utils.FlattenStringSlice(v.TunnelIPAddresses),
		})
	}
	return result
}

func vpnGatewayWaitForCreatedRefreshFunc(ctx context.Context, client *network.VpnGatewaysClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if VPN Gateway %q (Resource Group %q) has finished provisioning..", name, resourceGroup)

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			log.Printf("[DEBUG] Error retrieving VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
			return nil, "error", fmt.Errorf("Error retrieving VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if resp.VpnGatewayProperties == nil {
			log.Printf("[DEBUG] Error retrieving VPN Gateway %q (Resource Group %q): `properties` was nil", name, resourceGroup)
			return nil, "error", fmt.Errorf("Error retrieving VPN Gateway %q (Resource Group %q): `properties` was nil", name, resourceGroup)
		}

		log.Printf("[DEBUG] VPN Gateway %q (Resource Group %q) is %q..", name, resourceGroup, string(resp.VpnGatewayProperties.ProvisioningState))
		switch resp.VpnGatewayProperties.ProvisioningState {
		case network.Succeeded:
			return "available", "available", nil

		case network.Failed:
			return "error", "error", fmt.Errorf("VPN Gateway %q (Resource Group %q) is in provisioningState `Failed`", name, resourceGroup)

		default:
			return "pending", "pending", nil
		}
	}
}
