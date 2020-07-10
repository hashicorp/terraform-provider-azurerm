package compute

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPlatformImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPlatformImageRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"location": azure.SchemaLocation(),

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
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceArmPlatformImageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMImageClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := azure.NormalizeLocation(d.Get("location").(string))
	publisher := d.Get("publisher").(string)
	offer := d.Get("offer").(string)
	sku := d.Get("sku").(string)

	result, err := client.List(ctx, location, publisher, offer, sku, "", utils.Int32(int32(1000)), "name")
	if err != nil {
		return fmt.Errorf("Error reading Platform Images: %+v", err)
	}

	var image *compute.VirtualMachineImageResource
	if v, ok := d.GetOk("version"); ok {
		version := v.(string)
		for _, item := range *result.Value {
			if item.Name != nil && *item.Name == version {
				image = &item
				break
			}
		}
		if image == nil {
			return fmt.Errorf("could not find image (location %q / publisher %q / offer %q / sku %q / version % q): %+v", location, publisher, offer, sku, version, err)
		}
	} else {
		// get the latest image
		// the last value is the latest, apparently.
		image = &(*result.Value)[len(*result.Value)-1]
	}

	d.SetId(*image.ID)
	if location := image.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("publisher", publisher)
	d.Set("offer", offer)
	d.Set("sku", sku)
	d.Set("version", image.Name)

	return nil
}
