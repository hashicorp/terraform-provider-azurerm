package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/monitor/mgmt/2018-03-01/insights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmActionGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmActionGroupCreateOrUpdate,
		Read:   resourceArmActionGroupRead,
		Update: resourceArmActionGroupCreateOrUpdate,
		Delete: resourceArmActionGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location":            locationSchema(),
			"resource_group_name": resourceGroupNameSchema(),

			"short_name": {
				Type:     schema.TypeString,
				Required: true,
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
							Type:     schema.TypeString,
							Required: true,
						},
						"email_address": {
							Type:     schema.TypeString,
							Required: true,
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
							Type:     schema.TypeString,
							Required: true,
						},
						"country_code": {
							Type:     schema.TypeString,
							Required: true,
						},
						"phone_number": {
							Type:     schema.TypeString,
							Required: true,
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
							Type:     schema.TypeString,
							Required: true,
						},
						"service_uri": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmActionGroupCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
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
		Location: utils.String(location),
		ActionGroup: &insights.ActionGroup{
			GroupShortName:   utils.String(shortName),
			Enabled:          utils.Bool(enabled),
			EmailReceivers:   expandActionGroupEmailReceiver(d),
			SmsReceivers:     expandActionGroupSmsReceiver(d),
			WebhookReceivers: expandActionGroupWebHookReceiver(d),
		},
		Tags: expandedTags,
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

	return resourceArmActionGroupRead(d, meta)
}

func resourceArmActionGroupRead(d *schema.ResourceData, meta interface{}) error {
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

	if err = d.Set("email_receiver", flattenActionGroupEmailReceiver(resp.EmailReceivers)); err != nil {
		return fmt.Errorf("Error setting `email_receiver` of action group %s (resource group %s): %+v", name, resGroup, err)
	}

	if err = d.Set("sms_receiver", flattenActionGroupSmsReceiver(resp.SmsReceivers)); err != nil {
		return fmt.Errorf("Error setting `sms_receiver` of action group %s (resource group %s): %+v", name, resGroup, err)
	}

	if err = d.Set("webhook_receiver", flattenActionGroupWebHookReceiver(resp.WebhookReceivers)); err != nil {
		return fmt.Errorf("Error setting `webhook_receiver` of action group %s (resource group %s): %+v", name, resGroup, err)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmActionGroupDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandActionGroupEmailReceiver(d *schema.ResourceData) *[]insights.EmailReceiver {
	v, ok := d.GetOk("email_receiver")
	if !ok {
		return nil
	}

	receivers := make([]insights.EmailReceiver, 0)
	for _, receiverValue := range v.([]interface{}) {
		val := receiverValue.(map[string]interface{})

		receiver := insights.EmailReceiver{
			Name:         utils.String(val["name"].(string)),
			EmailAddress: utils.String(val["email_address"].(string)),
		}

		receivers = append(receivers, receiver)
	}
	return &receivers
}

func flattenActionGroupEmailReceiver(receivers *[]insights.EmailReceiver) []interface{} {
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

func expandActionGroupSmsReceiver(d *schema.ResourceData) *[]insights.SmsReceiver {
	v, ok := d.GetOk("sms_receiver")
	if !ok {
		return nil
	}

	receivers := make([]insights.SmsReceiver, 0)
	for _, receiverValue := range v.([]interface{}) {
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

func flattenActionGroupSmsReceiver(receivers *[]insights.SmsReceiver) []interface{} {
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

func expandActionGroupWebHookReceiver(d *schema.ResourceData) *[]insights.WebhookReceiver {
	v, ok := d.GetOk("webhook_receiver")
	if !ok {
		return nil
	}

	receivers := make([]insights.WebhookReceiver, 0)
	for _, receiverValue := range v.([]interface{}) {
		val := receiverValue.(map[string]interface{})

		receiver := insights.WebhookReceiver{
			Name:       utils.String(val["name"].(string)),
			ServiceURI: utils.String(val["service_uri"].(string)),
		}

		receivers = append(receivers, receiver)
	}
	return &receivers
}

func flattenActionGroupWebHookReceiver(receivers *[]insights.WebhookReceiver) []interface{} {
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
