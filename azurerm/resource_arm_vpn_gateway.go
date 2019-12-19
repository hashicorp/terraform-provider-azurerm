package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	networkSvc "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVPNGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVPNGatewayCreateUpdate,
		Read:   resourceArmVPNGatewayRead,
		Update: resourceArmVPNGatewayCreateUpdate,
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
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"virtual_hub_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkSvc.ValidateVirtualHubID,
			},

			"bgp_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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

func resourceArmVPNGatewayCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VpnGatewaysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_vpn_gateway", *existing.ID)
		}
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

	log.Printf("[DEBUG] Waiting for Virtual Hub %q (Resource Group %q) to become available", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"pending"},
		Target:                    []string{"available"},
		Refresh:                   vpnGatewayWaitForCreatedRefreshFunc(ctx, client, resourceGroup, name),
		Delay:                     30 * time.Second,
		PollInterval:              10 * time.Second,
		ContinuousTargetOccurence: 3,
	}

	if features.SupportsCustomTimeouts() {
		if d.IsNewResource() {
			stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
		} else {
			stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)
		}
	} else {
		stateConf.Timeout = 90 * time.Minute
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for creation of Virtual Hub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmVPNGatewayRead(d, meta)
}

func resourceArmVPNGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VpnGatewaysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := networkSvc.ParseVPNGatewayID(d.Id())
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

	id, err := networkSvc.ParseVPNGatewayID(d.Id())
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

func expandVPNGatewayBGPSettings(input []interface{}) *network.BgpSettings {
	if len(input) == 0 {
		return nil
	}

	val := input[0].(map[string]interface{})

	asn := val["asn"].(int)
	peerWeight := val["peer_weight"].(int)

	return &network.BgpSettings{
		Asn:        utils.Int64(int64(asn)),
		PeerWeight: utils.Int32(int32(peerWeight)),
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

	return []interface{}{
		map[string]interface{}{
			"asn":                 asn,
			"bgp_peering_address": bgpPeeringAddress,
			"peer_weight":         peerWeight,
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
