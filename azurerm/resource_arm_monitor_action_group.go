package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2018-03-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorActionGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorActionGroupCreateOrUpdate,
		Read:   resourceArmMonitorActionGroupRead,
		Update: resourceArmMonitorActionGroupCreateOrUpdate,
		Delete: resourceArmMonitorActionGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"short_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"email_receiver": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"email_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
					},
				},
			},

			"sms_receiver": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"country_code": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"phone_number": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
					},
				},
			},

			"webhook_receiver": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"service_uri": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmMonitorActionGroupCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).actionGroupsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)

	shortName := d.Get("short_name").(string)
	enabled := d.Get("enabled").(bool)

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	parameters := insights.ActionGroupResource{
		Location: &location,
		ActionGroup: &insights.ActionGroup{
			GroupShortName: &shortName,
			Enabled:        &enabled,
		},
		Tags: expandedTags,
	}

	if v, ok := d.GetOk("email_receiver"); ok {
		parameters.ActionGroup.EmailReceivers = expandMonitorActionGroupEmailReceiver(v.([]interface{}))
	}

	if v, ok := d.GetOk("sms_receiver"); ok {
		parameters.ActionGroup.SmsReceivers = expandMonitorActionGroupSmsReceiver(v.([]interface{}))
	}

	if v, ok := d.GetOk("webhook_receiver"); ok {
		parameters.ActionGroup.WebhookReceivers = expandMonitorActionGroupWebHookReceiver(v.([]interface{}))
	}

	_, err := client.CreateOrUpdate(ctx, resGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating or updating action group %s (resource group %s): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error getting action group %s (resource group %s) after creation: %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Action group %s (resource group %s) ID is empty", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMonitorActionGroupRead(d, meta)
}

func resourceArmMonitorActionGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).actionGroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing action group resource ID \"%s\" during get: %+v", d.Id(), err)
	}
	resGroup := id.ResourceGroup
	name := id.Path["actionGroups"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting action group %s (resource group %s): %+v", name, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("short_name", *resp.GroupShortName)
	d.Set("enabled", *resp.Enabled)

	if err = d.Set("email_receiver", flattenMonitorActionGroupEmailReceiver(resp.EmailReceivers)); err != nil {
		return fmt.Errorf("Error setting `email_receiver` of action group %s (resource group %s): %+v", name, resGroup, err)
	}

	if err = d.Set("sms_receiver", flattenMonitorActionGroupSmsReceiver(resp.SmsReceivers)); err != nil {
		return fmt.Errorf("Error setting `sms_receiver` of action group %s (resource group %s): %+v", name, resGroup, err)
	}

	if err = d.Set("webhook_receiver", flattenMonitorActionGroupWebHookReceiver(resp.WebhookReceivers)); err != nil {
		return fmt.Errorf("Error setting `webhook_receiver` of action group %s (resource group %s): %+v", name, resGroup, err)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmMonitorActionGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).actionGroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing action group resource ID \"%s\" during delete: %+v", d.Id(), err)
	}
	resGroup := id.ResourceGroup
	name := id.Path["actionGroups"]

	resp, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("Error deleting action group %s (resource group %s): %+v", name, resGroup, err)
	}

	return nil
}

func expandMonitorActionGroupEmailReceiver(v []interface{}) *[]insights.EmailReceiver {
	receivers := make([]insights.EmailReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := insights.EmailReceiver{
			Name:         utils.String(val["name"].(string)),
			EmailAddress: utils.String(val["email_address"].(string)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupSmsReceiver(v []interface{}) *[]insights.SmsReceiver {
	receivers := make([]insights.SmsReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := insights.SmsReceiver{
			Name:        utils.String(val["name"].(string)),
			CountryCode: utils.String(val["country_code"].(string)),
			PhoneNumber: utils.String(val["phone_number"].(string)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupWebHookReceiver(v []interface{}) *[]insights.WebhookReceiver {
	receivers := make([]insights.WebhookReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := insights.WebhookReceiver{
			Name:       utils.String(val["name"].(string)),
			ServiceURI: utils.String(val["service_uri"].(string)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func flattenMonitorActionGroupEmailReceiver(receivers *[]insights.EmailReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{}, 0)
			val["name"] = *receiver.Name
			val["email_address"] = *receiver.EmailAddress
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupSmsReceiver(receivers *[]insights.SmsReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{}, 0)
			val["name"] = *receiver.Name
			val["country_code"] = *receiver.CountryCode
			val["phone_number"] = *receiver.PhoneNumber
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupWebHookReceiver(receivers *[]insights.WebhookReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{}, 0)
			val["name"] = *receiver.Name
			val["service_uri"] = *receiver.ServiceURI
			result = append(result, val)
		}
	}
	return result
}
