package azurerm

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPlatformImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPlatformImageRead,
		Schema: map[string]*schema.Schema{
			"location": locationSchema(),

			"publisher": {
				Type:     schema.TypeString,
				Required: true,
			},

			"offer": {
				Type:     schema.TypeString,
				Required: true,
			},

			"sku": {
				Type:     schema.TypeString,
				Required: true,
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmPlatformImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vmImageClient
	ctx := meta.(*ArmClient).StopContext

	location := azure.NormalizeLocation(d.Get("location").(string))
	publisher := d.Get("publisher").(string)
	offer := d.Get("offer").(string)
	sku := d.Get("sku").(string)

	result, err := client.List(ctx, location, publisher, offer, sku, "", utils.Int32(int32(1000)), "name")
	if err != nil {
		return fmt.Errorf("Error reading Platform Images: %+v", err)
	}

	// the last value is the latest, apparently.
	resp := (*result.Value)[len(*result.Value)-1]

	d.SetId(*resp.ID)
	if err := azure.FlattenAndSetLocation(d, resp.Location); err != nil {
		return err
	}
	d.Set("publisher", publisher)
	d.Set("offer", offer)
	d.Set("sku", sku)
	d.Set("version", resp.Name)

	return nil
}
