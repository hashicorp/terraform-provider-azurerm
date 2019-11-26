package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	networkSvc "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPointToSiteVPNGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPointToSiteVPNGatewayCreateUpdate,
		Read:   resourceArmPointToSiteVPNGatewayRead,
		Update: resourceArmPointToSiteVPNGatewayCreateUpdate,
		Delete: resourceArmPointToSiteVPNGatewayDelete,
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

			"vpn_server_configuration_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkSvc.ValidateVpnServerConfigurationID,
			},

			"scale_unit": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPointToSiteVPNGatewayCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PointToSiteVpnGatewaysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Point-to-Site VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_point_to_site_vpn_gateway", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	scaleUnit := d.Get("scale_unit").(int)
	virtualHubId := d.Get("virtual_hub_id").(string)
	vpnServerConfigurationId := d.Get("vpn_server_configuration_id").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := network.P2SVpnGateway{
		Location: utils.String(location),
		P2SVpnGatewayProperties: &network.P2SVpnGatewayProperties{
			VpnServerConfiguration: &network.SubResource{
				ID: utils.String(vpnServerConfigurationId),
			},
			VirtualHub: &network.SubResource{
				ID: utils.String(virtualHubId),
			},
			VpnGatewayScaleUnit: utils.Int32(int32(scaleUnit)),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Point-to-Site VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of Point-to-Site VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Point-to-Site VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmPointToSiteVPNGatewayRead(d, meta)
}

func resourceArmPointToSiteVPNGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PointToSiteVpnGatewaysClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := networkSvc.ParsePointToSiteVPNGatewayID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.Base.ResourceGroup

	resp, err := client.Get(ctx, resourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Point-to-Site VPN Gateway %q was not found in Resource Group %q - removing from state!", id.Name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Point-to-Site VPN Gateway %q (Resource Group %q): %+v", id.Name, resourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.P2SVpnGatewayProperties; props != nil {
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

		vpnServerConfigurationId := ""
		if props.VpnServerConfiguration != nil && props.VpnServerConfiguration.ID != nil {
			vpnServerConfigurationId = *props.VpnServerConfiguration.ID
		}
		d.Set("vpn_server_configuration_id", vpnServerConfigurationId)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmPointToSiteVPNGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PointToSiteVpnGatewaysClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := networkSvc.ParsePointToSiteVPNGatewayID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.Base.ResourceGroup

	future, err := client.Delete(ctx, resourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Point-to-Site VPN Gateway %q (Resource Group %q): %+v", id.Name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Point-to-Site VPN Gateway %q (Resource Group %q): %+v", id.Name, resourceGroup, err)
	}

	return nil
}
