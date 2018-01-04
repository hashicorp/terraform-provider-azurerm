package azurerm

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/arm/iothub"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceArmIotHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIotHubCreateAndUpdate,
		Read:   resourceArmIotHubRead,
		Update: resourceArmIotHubCreateAndUpdate,
		Delete: resourceArmIotHubDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags":     tagsSchema(),
			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),
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

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"etag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"shared_access_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"primary_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secondary_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permissions": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}

}

func resourceArmIotHubCreateAndUpdate(d *schema.ResourceData, meta interface{}) error {

	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient

	rg := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	res, err := iothubClient.CheckNameAvailability(iothub.OperationInputs{
		Name: &name,
	})

	if err != nil {
		return err
	}

	if !*res.NameAvailable {
		return errors.New(string(res.Reason))
	}

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

	if tagsI, ok := d.GetOk("tags"); ok {
		tags := tagsI.(map[string]interface{})
		desc.Tags = expandTags(tags)
	}

	if etagI, ok := d.GetOk("etag"); ok {
		etag := etagI.(string)
		desc.Etag = &etag
	}

	cancel := make(chan struct{})

	_, errChan := iothubClient.CreateOrUpdate(rg, name, desc, cancel)
	err = <-errChan

	if err != nil {
		return err
	}

	desc, err = iothubClient.Get(rg, name)
	if err != nil {
		return err
	}

	d.SetId(*desc.ID)
	return resourceArmIotHubRead(d, meta)

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

func resourceArmIotHubRead(d *schema.ResourceData, meta interface{}) error {

	id, err := parseAzureResourceID(d.Id())

	if err != nil {
		return err
	}

	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient
	iothubName := id.Path["IotHubs"]
	desc, err := iothubClient.Get(id.ResourceGroup, iothubName)

	if err != nil {
		return err
	}

	properties := desc.Properties

	keysResp, err := iothubClient.ListKeys(id.ResourceGroup, iothubName)

	var keys []map[string]interface{}
	for _, key := range *keysResp.Value {
		keyMap := make(map[string]interface{})
		keyMap["key_name"] = *key.KeyName
		keyMap["primary_key"] = *key.PrimaryKey
		keyMap["secondary_key"] = *key.SecondaryKey
		keyMap["permissions"] = string(key.Rights)
		keys = append(keys, keyMap)
	}

	if err != nil {
		return err
	}

	d.Set("shared_access_policy", keys)
	d.Set("hostname", *properties.HostName)
	d.Set("etag", *desc.Etag)
	d.Set("type", *desc.Type)
	flattenAndSetTags(d, desc.Tags)

	return nil
}

func resourceArmIotHubDelete(d *schema.ResourceData, meta interface{}) error {

	id, err := parseAzureResourceID(d.Id())

	if err != nil {
		return err
	}

	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient

	_, errChan := iothubClient.Delete(id.ResourceGroup, id.Path["IotHubs"], make(chan struct{}))
	err = <-errChan
	return err
}
