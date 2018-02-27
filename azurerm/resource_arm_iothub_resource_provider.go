package azurerm

import (
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2017-07-01/devices"
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
								string(devices.F1),
								string(devices.S1),
								string(devices.S2),
								string(devices.S3),
							}, true),
						},

						"tier": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.Free),
								string(devices.Standard),
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
	ctx := armClient.StopContext

	rg := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	res, err := iothubClient.CheckNameAvailability(ctx, devices.OperationInputs{
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

	desc := devices.IotHubDescription{
		Resourcegroup:  &rg,
		Name:           &name,
		Location:       &location,
		Subscriptionid: &subscriptionID,
		Sku:            &skuInfo,
	}

	if tagsI, ok := d.GetOk("tags"); ok {
		tags := tagsI.(map[string]interface{})
		desc.Tags = *expandTags(tags)
	}

	match := ""

	if etagI, ok := d.GetOk("etag"); ok {
		etag := etagI.(string)
		desc.Etag = &etag
		match = etag
	}

	future, err := iothubClient.CreateOrUpdate(ctx, rg, name, desc, match)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, iothubClient.Client)
	if err != nil {
		return fmt.Errorf("Error creating or updating IotHub %q (Resource Group %q): %+v", name, rg, err)
	}

	desc, err = iothubClient.Get(ctx, rg, name)
	if err != nil {
		return err
	}

	d.SetId(*desc.ID)
	return resourceArmIotHubRead(d, meta)

}

func expandAzureRmIotHubSku(d *schema.ResourceData) devices.IotHubSkuInfo {
	skuList := d.Get("sku").([]interface{})
	skuMap := skuList[0].(map[string]interface{})
	cap := int64(skuMap["capacity"].(int))

	return devices.IotHubSkuInfo{
		Name:     devices.IotHubSku(skuMap["name"].(string)),
		Tier:     devices.IotHubSkuTier(skuMap["tier"].(string)),
		Capacity: &cap,
	}

}

func resourceArmIotHubRead(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient
	ctx := armClient.StopContext

	id, err := parseAzureResourceID(d.Id())

	if err != nil {
		return err
	}

	iothubName := id.Path["IotHubs"]
	desc, err := iothubClient.Get(ctx, id.ResourceGroup, iothubName)

	if err != nil {
		return err
	}

	properties := desc.Properties

	keysResp, err := iothubClient.ListKeys(ctx, id.ResourceGroup, iothubName)
	keyList := keysResp.Response()

	var keys []map[string]interface{}
	for _, key := range *keyList.Value {
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
	flattenAndSetTags(d, &desc.Tags)

	return nil
}

func resourceArmIotHubDelete(d *schema.ResourceData, meta interface{}) error {

	id, err := parseAzureResourceID(d.Id())

	if err != nil {
		return err
	}

	armClient := meta.(*ArmClient)
	iothubClient := armClient.iothubResourceClient
	ctx := armClient.StopContext

	future, err := iothubClient.Delete(ctx, id.ResourceGroup, id.Path["IotHubs"])
	if err != nil {
		return err
	}

	status := future.Status()
	if status == "unkown" {
		return fmt.Errorf("Error Waiting for the deletion of IoTHub %q (Resource Group %q): %+v", id.Path["IotHubs"], id.ResourceGroup, err)
	}

	return nil
}
