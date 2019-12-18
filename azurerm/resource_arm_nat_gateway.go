package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	aznet "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var natGatewayResourceName = "azurerm_nat_gateway"

func resourceArmNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNatGatewayCreateUpdate,
		Read:   resourceArmNatGatewayRead,
		Update: resourceArmNatGatewayCreateUpdate,
		Delete: resourceArmNatGatewayDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: aznet.ValidateNatGatewayName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"idle_timeout_in_minutes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      4,
				ValidateFunc: validation.IntBetween(4, 120),
			},

			"public_ip_address_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},

			"public_ip_prefix_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},

			"sku_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(network.Standard),
				ValidateFunc: validation.StringInSlice([]string{
					string(network.Standard),
				}, false),
			},

			"zones": azure.SchemaZones(),

			"resource_guid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmNatGatewayCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.NatGatewayClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing NAT Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if resp.ID != nil && *resp.ID != "" {
			return tf.ImportAsExistsError("azurerm_nat_gateway", *resp.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	idleTimeoutInMinutes := d.Get("idle_timeout_in_minutes").(int)
	publicIpAddressIds := d.Get("public_ip_address_ids").(*schema.Set).List()
	publicIpPrefixIds := d.Get("public_ip_prefix_ids").(*schema.Set).List()
	skuName := d.Get("sku_name").(string)
	zones := d.Get("zones").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	parameters := network.NatGateway{
		Location: utils.String(location),
		NatGatewayPropertiesFormat: &network.NatGatewayPropertiesFormat{
			IdleTimeoutInMinutes: utils.Int32(int32(idleTimeoutInMinutes)),
			PublicIPAddresses:    expandArmNatGatewaySubResourceID(publicIpAddressIds),
			PublicIPPrefixes:     expandArmNatGatewaySubResourceID(publicIpPrefixIds),
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

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving NAT Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Cannot read NAT Gateway %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmNatGatewayRead(d, meta)
}

func resourceArmNatGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.NatGatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["natGateways"]

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] NAT Gateway %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading NAT Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("sku_name", resp.Sku.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.NatGatewayPropertiesFormat; props != nil {
		d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)
		d.Set("resource_guid", props.ResourceGUID)

		if err := d.Set("public_ip_address_ids", flattenArmNatGatewaySubResourceID(props.PublicIPAddresses)); err != nil {
			return fmt.Errorf("Error setting `public_ip_address_ids`: %+v", err)
		}

		if err := d.Set("public_ip_prefix_ids", flattenArmNatGatewaySubResourceID(props.PublicIPPrefixes)); err != nil {
			return fmt.Errorf("Error setting `public_ip_prefix_ids`: %+v", err)
		}
	}

	if err := d.Set("zones", utils.FlattenStringSlice(resp.Zones)); err != nil {
		return fmt.Errorf("Error setting `zones`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.NatGatewayClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["natGateways"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting NAT Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting NAT Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmNatGatewaySubResourceID(input []interface{}) *[]network.SubResource {
	results := make([]network.SubResource, 0)
	for _, item := range input {
		id := item.(string)

		results = append(results, network.SubResource{
			ID: utils.String(id),
		})
	}
	return &results
}

func flattenArmNatGatewaySubResourceID(input *[]network.SubResource) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if item.ID != nil {
			results = append(results, *item.ID)
		}
	}

	return results
}
