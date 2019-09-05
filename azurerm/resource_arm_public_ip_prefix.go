package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPublicIpPrefix() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPublicIpPrefixCreateUpdate,
		Read:   resourceArmPublicIpPrefixRead,
		Update: resourceArmPublicIpPrefixCreateUpdate,
		Delete: resourceArmPublicIpPrefixDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(network.Standard),
				ValidateFunc: validation.StringInSlice([]string{
					string(network.Standard),
				}, false),
			},

			"prefix_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      28,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(24, 31),
			},

			"ip_prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"zones": azure.SchemaSingleZone(),

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPublicIpPrefixCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PublicIPPrefixesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Public IP Prefix creation.")

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	sku := d.Get("sku").(string)
	prefix_length := d.Get("prefix_length").(int)
	t := d.Get("tags").(map[string]interface{})
	zones := azure.ExpandZones(d.Get("zones").([]interface{}))

	publicIpPrefix := network.PublicIPPrefix{
		Name:     &name,
		Location: &location,
		Sku: &network.PublicIPPrefixSku{
			Name: network.PublicIPPrefixSkuName(sku),
		},
		PublicIPPrefixPropertiesFormat: &network.PublicIPPrefixPropertiesFormat{
			PrefixLength: utils.Int32(int32(prefix_length)),
		},
		Tags:  tags.Expand(t),
		Zones: zones,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, publicIpPrefix)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Public IP Prefix %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Public IP Prefix %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Public IP Prefix %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPublicIpPrefixRead(d, meta)
}

func resourceArmPublicIpPrefixRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PublicIPPrefixesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["publicIPPrefixes"]

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Public IP Prefix %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("zones", resp.Zones)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Name))
	}

	if props := resp.PublicIPPrefixPropertiesFormat; props != nil {
		d.Set("prefix_length", props.PrefixLength)
		d.Set("ip_prefix", props.IPPrefix)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmPublicIpPrefixDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PublicIPPrefixesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["publicIPPrefixes"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Public IP Prefix %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Public IP Prefix %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}
