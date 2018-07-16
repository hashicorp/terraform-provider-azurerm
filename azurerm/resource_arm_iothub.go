package azurerm

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2017-07-01/devices"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								string(devices.F1),
								string(devices.S1),
								string(devices.S2),
								string(devices.S3),
							}, true),
						},

						"tier": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"secondary_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						"permissions": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}

}

func resourceArmIotHubCreateAndUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothubResourceClient
	ctx := meta.(*ArmClient).StopContext
	subscriptionID := meta.(*ArmClient).subscriptionId

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	res, err := client.CheckNameAvailability(ctx, devices.OperationInputs{
		Name: &name,
	})

	if err != nil {
		return fmt.Errorf("An error occurred checking if the IoTHub name was unique: %+v", err)
	}

	if !*res.NameAvailable {
		_, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("An IoTHub already exists with the name %q - please choose an alternate name: %s", name, string(res.Reason))
		}
	}

	location := azureRMNormalizeLocation(d.Get("location").(string))
	skuInfo := expandIoTHubSku(d)
	tags := d.Get("tags").(map[string]interface{})

	properties := devices.IotHubDescription{
		Name:           utils.String(name),
		Location:       utils.String(location),
		Resourcegroup:  utils.String(resourceGroup),
		Subscriptionid: utils.String(subscriptionID),
		Sku:            &skuInfo,
		Tags:           expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, properties, "")
	if err != nil {
		return fmt.Errorf("Error creating/updating IotHub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the completion of the creating/updating of IotHub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)
	return resourceArmIotHubRead(d, meta)
}

func resourceArmIotHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothubResourceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["IotHubs"]
	resourceGroup := id.ResourceGroup
	hub, err := client.Get(ctx, id.ResourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(hub.Response) {
			log.Printf("[DEBUG] IoTHub %q (Resource Group %q) was not found!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving IotHub Client %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	keysResp, err := client.ListKeys(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error listing keys for IoTHub %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	keyList := keysResp.Response()
	keys := flattenIoTHubSharedAccessPolicy(keyList.Value)

	if err := d.Set("shared_access_policy", keys); err != nil {
		return fmt.Errorf("Error flattening `shared_access_policy` in IoTHub %q: %+v", name, err)
	}

	if properties := hub.Properties; properties != nil {
		d.Set("hostname", properties.HostName)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := hub.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	sku := flattenIoTHubSku(hub.Sku)
	if err := d.Set("sku", sku); err != nil {
		return fmt.Errorf("Error flattening `sku`: %+v", err)
	}
	d.Set("type", hub.Type)
	flattenAndSetTags(d, hub.Tags)

	return nil
}

func resourceArmIotHubDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*ArmClient).iothubResourceClient
	ctx := meta.(*ArmClient).StopContext

	name := id.Path["IotHubs"]
	resourceGroup := id.ResourceGroup

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return err
	}

	return waitForIotHubToBeDeleted(ctx, client, resourceGroup, name)
}

func waitForIotHubToBeDeleted(ctx context.Context, client devices.IotHubResourceClient, resourceGroup, name string) error {
	// we can't use the Waiter here since the API returns a 404 once it's deleted which is considered a polling status code..
	log.Printf("[DEBUG] Waiting for IotHub (%q in Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: iothubStateStatusCodeRefreshFunc(ctx, client, resourceGroup, name),
		Timeout: 40 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for IotHub (%q in Resource Group %q) to be deleted: %+v", name, resourceGroup, err)
	}

	return nil
}

func iothubStateStatusCodeRefreshFunc(ctx context.Context, client devices.IotHubResourceClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroup, name)

		log.Printf("Retrieving IoTHub %q (Resource Group %q) returned Status %d", resourceGroup, name, res.StatusCode)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, strconv.Itoa(res.StatusCode), nil
			}
			return nil, "", fmt.Errorf("Error polling for the status of the IotHub %q (RG: %q): %+v", name, resourceGroup, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func expandIoTHubSku(d *schema.ResourceData) devices.IotHubSkuInfo {
	skuList := d.Get("sku").([]interface{})
	skuMap := skuList[0].(map[string]interface{})
	cap := int64(skuMap["capacity"].(int))

	name := skuMap["name"].(string)
	tier := skuMap["tier"].(string)

	return devices.IotHubSkuInfo{
		Name:     devices.IotHubSku(name),
		Tier:     devices.IotHubSkuTier(tier),
		Capacity: &cap,
	}
}

func flattenIoTHubSku(input *devices.IotHubSkuInfo) []interface{} {
	output := make(map[string]interface{}, 0)

	output["name"] = string(input.Name)
	output["tier"] = string(input.Tier)
	if capacity := input.Capacity; capacity != nil {
		output["capacity"] = int(*capacity)
	}

	return []interface{}{output}
}

func flattenIoTHubSharedAccessPolicy(input *[]devices.SharedAccessSignatureAuthorizationRule) []interface{} {
	results := make([]interface{}, 0)

	if keys := input; keys != nil {
		for _, key := range *keys {
			keyMap := make(map[string]interface{})

			if keyName := key.KeyName; keyName != nil {
				keyMap["key_name"] = *keyName
			}

			if primaryKey := key.PrimaryKey; primaryKey != nil {
				keyMap["primary_key"] = *primaryKey
			}

			if secondaryKey := key.SecondaryKey; secondaryKey != nil {
				keyMap["secondary_key"] = *secondaryKey
			}

			keyMap["permissions"] = string(key.Rights)
			results = append(results, keyMap)
		}
	}

	return results
}
