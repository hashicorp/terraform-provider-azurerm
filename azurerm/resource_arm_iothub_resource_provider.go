package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/arm/iothub"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceArmIothub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIothubCreate,
		Read:   resourceArmIothubRead,
		Update: resourceArmIothubUpdate,
		Delete: resourceArmIothubDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"location": {
				Type:     schema.TypeString,
				Required: true,
			},

			"subscription_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group": {
				Type:     schema.TypeString,
				Required: true,
			},

			"etag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"skuinfo": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type: schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(iothub.F1),
								string(iothub.S1),
								string(iothub.S2),
								string(iothub.S3),
							}, true),
						},

						"tier": {
							Type: schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(iothub.Free),
								string(iothub.Standard),
							}, true),
						},

						"capacity": {
							Type:         schema.TypeInt,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
		},
	}

}
func resourceArmIothubCreate(d *schema.ResourceData, meta interface{}) error {

	return nil

}

func resourceArmIothubRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceArmIothubUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceArmIothubDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
