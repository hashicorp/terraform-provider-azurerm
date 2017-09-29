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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// TODO: is this needed
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"location": locationSchema(),

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

	rg := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)
	subscriptionID := armClient.subscriptionId
	skuInfo := expandAzureRmIotHubSku(d)

	desc := iothub.Description{
		Resourcegroup:  &rg,
		Name:           &name,
		Location:       &location,
		Subscriptionid: &subscriptionID,
		Sku:            &skuInfo,
	}

	if etagI, ok := d.GetOk("etag"); ok {
		etag := etagI.(string)
		desc.Etag = &etag
	}

	cancel := make(chan struct{})

	_, errChan := iothubClient.CreateOrUpdate(rg, name, desc, cancel)
	err := <-errChan

	if err != nil {
		return err
	}

	desc, err = iothubClient.Get(rg, name)
	if err != nil {
		return err
	}

	d.SetId(*desc.ID)
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

	id, err := parseAzureResourceID(d.Id())

	if err != nil {
		return err
	}

	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient
	desc, err := iothubClient.Get(id.ResourceGroup, id.Path["IotHubs"])
	if err != nil {
		return err
	}

	d.Set("etag", *desc.Etag)
	d.Set("type", *desc.Type)

	return nil
}

func resourceArmIothubUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceArmIothubDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
