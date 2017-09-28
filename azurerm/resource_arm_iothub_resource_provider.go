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

			"location": locationSchema(),

			"subscription_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"etag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tags": tagsSchema(),

			"sku": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(iothub.F1),
								string(iothub.S1),
								string(iothub.S2),
								string(iothub.S3),
							}, true),
						},

						"tier": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(iothub.Free),
								string(iothub.Standard),
							}, true),
						},

						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
					},
				},
			},
		},
	}

}
func resourceArmIothubCreate(d *schema.ResourceData, meta interface{}) error {

	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient

	rg := d.Get("resource_group").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)
	subscriptionId := d.Get("subscription_id").(string)
	skuInfo := expandAzureRmIotHubSku(d)

	desc := iothub.Description{
		Name:           &name,
		Location:       &location,
		Subscriptionid: &subscriptionId,
		Sku:            &skuInfo,
	}

	cancel := make(chan struct{})

	RespChan, errChan := iothubClient.CreateOrUpdate(rg, name, desc, cancel)
	resp := <-RespChan
	err := <-errChan

	if err != nil {
		return err
	}

	d.SetId(*resp.ID)
	return resourceArmIothubRead(d, meta)

}

func expandAzureRmIotHubSku(d *schema.ResourceData) iothub.SkuInfo {
	skuList := d.Get("sku").([]interface{})
	skuMap := skuList[0].(map[string]interface{})
	cap := int64(skuMap["capacity"].(int))

	return iothub.SkuInfo{
		Name:     iothub.Sku(skuMap["name"].(string)),
		Tier:     iothub.SkuTier(skuMap["tier"].(string)),
		Capacity: &cap,
	}

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
