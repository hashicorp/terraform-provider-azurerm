package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var natGatewayResourceName = "azurerm_nat_gateway"

func resourceNatGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNatGatewayCreate,
		Read:   resourceNatGatewayRead,
		Update: resourceNatGatewayUpdate,
		Delete: resourceNatGatewayDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NatGatewayName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"idle_timeout_in_minutes": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      4,
				ValidateFunc: validation.IntBetween(4, 120),
			},

			"public_ip_address_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
				// TODO: remove in 3.0
				Deprecated: "Inline Public IP Address ID Deprecations have been deprecated in favour of the `azurerm_nat_gateway_public_ip_association` pluginsdk. This field will be removed in the next major version of the Azure Provider.",
			},

			"public_ip_prefix_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(network.NatGatewaySkuNameStandard),
				ValidateFunc: validation.StringInSlice([]string{
					string(network.NatGatewaySkuNameStandard),
				}, false),
			},

			"zones": azure.SchemaZones(),

			"resource_guid": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceNatGatewayCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(name, natGatewayResourceName)
	defer locks.UnlockByName(name, natGatewayResourceName)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error checking for present of existing NAT Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if resp.ID != nil && *resp.ID != "" {
		return tf.ImportAsExistsError("azurerm_nat_gateway", *resp.ID)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	idleTimeoutInMinutes := d.Get("idle_timeout_in_minutes").(int)
	publicIpAddressIds := d.Get("public_ip_address_ids").(*pluginsdk.Set).List()
	publicIpPrefixIds := d.Get("public_ip_prefix_ids").(*pluginsdk.Set).List()
	skuName := d.Get("sku_name").(string)
	zones := d.Get("zones").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	parameters := network.NatGateway{
		Location: utils.String(location),
		NatGatewayPropertiesFormat: &network.NatGatewayPropertiesFormat{
			IdleTimeoutInMinutes: utils.Int32(int32(idleTimeoutInMinutes)),
			PublicIPAddresses:    expandNetworkSubResourceID(publicIpAddressIds),
			PublicIPPrefixes:     expandNetworkSubResourceID(publicIpPrefixIds),
		},
		Sku: &network.NatGatewaySku{
			Name: network.NatGatewaySkuName(skuName),
		},
		Tags:  tags.Expand(t),
		Zones: utils.ExpandStringSlice(zones),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating NAT Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of NAT Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err = client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving NAT Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read NAT Gateway %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceNatGatewayRead(d, meta)
}

func resourceNatGatewayUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NatGatewayID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, natGatewayResourceName)
	defer locks.UnlockByName(id.Name, natGatewayResourceName)

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("NAT Gateway %q (Resource Group %q) was not found!", id.Name, id.ResourceGroup)
		}

		return fmt.Errorf("retrieving NAT Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if existing.NatGatewayPropertiesFormat == nil {
		return fmt.Errorf("retrieving NAT Gateway %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
	}
	props := *existing.NatGatewayPropertiesFormat

	// intentionally building a new object rather than reusing due to the additional read-only fields
	parameters := network.NatGateway{
		Location: existing.Location,
		NatGatewayPropertiesFormat: &network.NatGatewayPropertiesFormat{
			IdleTimeoutInMinutes: props.IdleTimeoutInMinutes,
			PublicIPAddresses:    props.PublicIPAddresses, // note: these can be managed via the separate resource
			PublicIPPrefixes:     props.PublicIPPrefixes,
		},
		Sku:   existing.Sku,
		Tags:  existing.Tags,
		Zones: existing.Zones,
	}

	if d.HasChange("idle_timeout_in_minutes") {
		timeout := d.Get("idle_timeout_in_minutes").(int)
		parameters.NatGatewayPropertiesFormat.IdleTimeoutInMinutes = utils.Int32(int32(timeout))
	}

	if d.HasChange("sku_name") {
		skuName := d.Get("sku_name").(string)
		parameters.Sku = &network.NatGatewaySku{
			Name: network.NatGatewaySkuName(skuName),
		}
	}

	if d.HasChange("public_ip_address_ids") {
		publicIpAddressIds := d.Get("public_ip_address_ids").(*pluginsdk.Set).List()
		parameters.NatGatewayPropertiesFormat.PublicIPAddresses = expandNetworkSubResourceID(publicIpAddressIds)
	}

	if d.HasChange("public_ip_prefix_ids") {
		publicIpPrefixIds := d.Get("public_ip_prefix_ids").(*pluginsdk.Set).List()
		parameters.NatGatewayPropertiesFormat.PublicIPPrefixes = expandNetworkSubResourceID(publicIpPrefixIds)
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]interface{})
		parameters.Tags = tags.Expand(t)
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("updating NAT Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of NAT Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceNatGatewayRead(d, meta)
}

func resourceNatGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NatGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] NAT Gateway %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading NAT Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.NatGatewayPropertiesFormat; props != nil {
		d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)
		d.Set("resource_guid", props.ResourceGUID)

		if err := d.Set("public_ip_address_ids", flattenNetworkSubResourceID(props.PublicIPAddresses)); err != nil {
			return fmt.Errorf("Error setting `public_ip_address_ids`: %+v", err)
		}

		if err := d.Set("public_ip_prefix_ids", flattenNetworkSubResourceID(props.PublicIPPrefixes)); err != nil {
			return fmt.Errorf("Error setting `public_ip_prefix_ids`: %+v", err)
		}
	}

	if err := d.Set("zones", utils.FlattenStringSlice(resp.Zones)); err != nil {
		return fmt.Errorf("Error setting `zones`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceNatGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NatGatewayID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, natGatewayResourceName)
	defer locks.UnlockByName(id.Name, natGatewayResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting NAT Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting NAT Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
