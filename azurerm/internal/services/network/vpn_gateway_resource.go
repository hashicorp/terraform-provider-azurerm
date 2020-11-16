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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var VPNGatewayResourceName = "azurerm_vpn_gateway"

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

						"instance_0_bgp_peering_address": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
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

						"instance_1_bgp_peering_address": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
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

	locks.ByName(name, VPNGatewayResourceName)
	defer locks.UnlockByName(name, VPNGatewayResourceName)

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
	if len(bgpSettingsRaw) > 0 && resp.VpnGatewayProperties != nil && resp.VpnGatewayProperties.BgpSettings != nil && resp.VpnGatewayProperties.BgpSettings.BgpPeeringAddresses != nil {
		val := bgpSettingsRaw[0].(map[string]interface{})
		input0 := val["instance_0_bgp_peering_address"].([]interface{})
		input1 := val["instance_1_bgp_peering_address"].([]interface{})

		if len(input0) > 0 || len(input1) > 0 {
			if len(input0) > 0 {
				val := input0[0].(map[string]interface{})
				(*resp.VpnGatewayProperties.BgpSettings.BgpPeeringAddresses)[0].CustomBgpIPAddresses = utils.ExpandStringSlice(val["custom_ips"].(*schema.Set).List())
			}
			if len(input1) > 0 {
				val := input1[0].(map[string]interface{})
				(*resp.VpnGatewayProperties.BgpSettings.BgpPeeringAddresses)[1].CustomBgpIPAddresses = utils.ExpandStringSlice(val["custom_ips"].(*schema.Set).List())
			}
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

	locks.ByName(name, VPNGatewayResourceName)
	defer locks.UnlockByName(name, VPNGatewayResourceName)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving for presence of existing VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if d.HasChange("scale_unit") {
		existing.VpnGatewayScaleUnit = utils.Int32(int32(d.Get("scale_unit").(int)))
	}
	if d.HasChange("tags") {
		existing.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	bgpSettingsRaw := d.Get("bgp_settings").([]interface{})
	if len(bgpSettingsRaw) > 0 {
		val := bgpSettingsRaw[0].(map[string]interface{})

		if d.HasChange("bgp_settings.0.instance_0_bgp_peering_address") {
			if input := val["instance_0_bgp_peering_address"].([]interface{}); len(input) > 0 {
				val := input[0].(map[string]interface{})
				(*existing.VpnGatewayProperties.BgpSettings.BgpPeeringAddresses)[0].CustomBgpIPAddresses = utils.ExpandStringSlice(val["custom_ips"].(*schema.Set).List())
			}
		}
		if d.HasChange("bgp_settings.0.instance_1_bgp_peering_address") {
			if input := val["instance_1_bgp_peering_address"].([]interface{}); len(input) > 0 {
				val := input[0].(map[string]interface{})
				(*existing.VpnGatewayProperties.BgpSettings.BgpPeeringAddresses)[1].CustomBgpIPAddresses = utils.ExpandStringSlice(val["custom_ips"].(*schema.Set).List())
			}
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

	id, err := parse.VPNGatewayID(d.Id())
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

	id, err := parse.VPNGatewayID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, VPNGatewayResourceName)
	defer locks.UnlockByName(id.Name, VPNGatewayResourceName)

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

	var instance0BgpPeeringAddress, instance1BgpPeeringAddress []interface{}
	if input.BgpPeeringAddresses != nil && len(*input.BgpPeeringAddresses) > 0 {
		instance0BgpPeeringAddress = flattenVPNGatewayIPConfigurationBgpPeeringAddress((*input.BgpPeeringAddresses)[0])
	}
	if input.BgpPeeringAddresses != nil && len(*input.BgpPeeringAddresses) > 1 {
		instance1BgpPeeringAddress = flattenVPNGatewayIPConfigurationBgpPeeringAddress((*input.BgpPeeringAddresses)[1])
	}

	return []interface{}{
		map[string]interface{}{
			"asn":                            asn,
			"bgp_peering_address":            bgpPeeringAddress,
			"instance_0_bgp_peering_address": instance0BgpPeeringAddress,
			"instance_1_bgp_peering_address": instance1BgpPeeringAddress,
			"peer_weight":                    peerWeight,
		},
	}
}

func flattenVPNGatewayIPConfigurationBgpPeeringAddress(input network.IPConfigurationBgpPeeringAddress) []interface{} {
	ipConfigurationID := ""
	if input.IpconfigurationID != nil {
		ipConfigurationID = *input.IpconfigurationID
	}

	return []interface{}{
		map[string]interface{}{
			"ip_configuration_id": ipConfigurationID,
			"custom_ips":          utils.FlattenStringSlice(input.CustomBgpIPAddresses),
			"default_ips":         utils.FlattenStringSlice(input.DefaultBgpIPAddresses),
			"tunnel_ips":          utils.FlattenStringSlice(input.TunnelIPAddresses),
		},
	}
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
