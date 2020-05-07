package monitor

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorActionGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorActionGroupCreateUpdate,
		Read:   resourceArmMonitorActionGroupRead,
		Update: resourceArmMonitorActionGroupCreateUpdate,
		Delete: resourceArmMonitorActionGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"short_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 12),
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
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"email_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"use_common_alert_schema": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"itsm_receiver": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"workspace_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"connection_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"ticket_configuration": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: structure.SuppressJsonDiff,
						},
						"region": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validation.StringIsNotEmpty,
							DiffSuppressFunc: location.DiffSuppressFunc,
						},
					},
				},
			},

			"azure_app_push_receiver": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"email_address": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"country_code": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"phone_number": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"service_uri": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
						},
						"use_common_alert_schema": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"automation_runbook_receiver": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"automation_account_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"runbook_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"webhook_resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"is_global_runbook": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"service_uri": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
						},
						"use_common_alert_schema": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"voice_receiver": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"country_code": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"phone_number": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"logic_app_receiver": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"callback_url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
						},
						"use_common_alert_schema": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"azure_function_receiver": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"function_app_resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"function_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"http_trigger_url": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
						},
						"use_common_alert_schema": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"arm_role_receiver": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"role_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"use_common_alert_schema": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			"tags": tags.Schema(),
		},
	}
}

func resourceArmMonitorActionGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionGroupsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Monitor Action Group Service Plan %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_monitor_action_group", *existing.ID)
		}
	}

	shortName := d.Get("short_name").(string)
	enabled := d.Get("enabled").(bool)

	emailReceiversRaw := d.Get("email_receiver").([]interface{})
	itsmReceiversRaw := d.Get("itsm_receiver").([]interface{})
	azureAppPushReceiversRaw := d.Get("azure_app_push_receiver").([]interface{})
	smsReceiversRaw := d.Get("sms_receiver").([]interface{})
	webhookReceiversRaw := d.Get("webhook_receiver").([]interface{})
	automationRunbookReceiversRaw := d.Get("automation_runbook_receiver").([]interface{})
	voiceReceiversRaw := d.Get("voice_receiver").([]interface{})
	logicAppReceiversRaw := d.Get("logic_app_receiver").([]interface{})
	azureFunctionReceiversRaw := d.Get("azure_function_receiver").([]interface{})
	armRoleReceiversRaw := d.Get("arm_role_receiver").([]interface{})

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	parameters := insights.ActionGroupResource{
		Location: utils.String(azure.NormalizeLocation("Global")),
		ActionGroup: &insights.ActionGroup{
			GroupShortName:             utils.String(shortName),
			Enabled:                    utils.Bool(enabled),
			EmailReceivers:             expandMonitorActionGroupEmailReceiver(emailReceiversRaw),
			AzureAppPushReceivers:      expandMonitorActionGroupAzureAppPushReceiver(azureAppPushReceiversRaw),
			ItsmReceivers:              expandMonitorActionGroupItsmReceiver(itsmReceiversRaw),
			SmsReceivers:               expandMonitorActionGroupSmsReceiver(smsReceiversRaw),
			WebhookReceivers:           expandMonitorActionGroupWebHookReceiver(webhookReceiversRaw),
			AutomationRunbookReceivers: expandMonitorActionGroupAutomationRunbookReceiver(automationRunbookReceiversRaw),
			VoiceReceivers:             expandMonitorActionGroupVoiceReceiver(voiceReceiversRaw),
			LogicAppReceivers:          expandMonitorActionGroupLogicAppReceiver(logicAppReceiversRaw),
			AzureFunctionReceivers:     expandMonitorActionGroupAzureFunctionReceiver(azureFunctionReceiversRaw),
			ArmRoleReceivers:           expandMonitorActionGroupArmRoleReceiver(armRoleReceiversRaw),
		},
		Tags: expandedTags,
	}

	if _, err := client.CreateOrUpdate(ctx, resGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating or updating action group %q (resource group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error getting action group %q (resource group %q) after creation: %+v", name, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Action group %q (resource group %q) ID is empty", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMonitorActionGroupRead(d, meta)
}

func resourceArmMonitorActionGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["actionGroups"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting action group %q (resource group %q): %+v", name, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)

	if group := resp.ActionGroup; group != nil {
		d.Set("short_name", group.GroupShortName)
		d.Set("enabled", group.Enabled)

		if err = d.Set("email_receiver", flattenMonitorActionGroupEmailReceiver(group.EmailReceivers)); err != nil {
			return fmt.Errorf("Error setting `email_receiver`: %+v", err)
		}

		if err = d.Set("itsm_receiver", flattenMonitorActionGroupItsmReceiver(group.ItsmReceivers)); err != nil {
			return fmt.Errorf("Error setting `itsm_receiver`: %+v", err)
		}

		if err = d.Set("azure_app_push_receiver", flattenMonitorActionGroupAzureAppPushReceiver(group.AzureAppPushReceivers)); err != nil {
			return fmt.Errorf("Error setting `azure_app_push_receiver`: %+v", err)
		}

		if err = d.Set("sms_receiver", flattenMonitorActionGroupSmsReceiver(group.SmsReceivers)); err != nil {
			return fmt.Errorf("Error setting `sms_receiver`: %+v", err)
		}

		if err = d.Set("webhook_receiver", flattenMonitorActionGroupWebHookReceiver(group.WebhookReceivers)); err != nil {
			return fmt.Errorf("Error setting `webhook_receiver`: %+v", err)
		}

		if err = d.Set("automation_runbook_receiver", flattenMonitorActionGroupAutomationRunbookReceiver(group.AutomationRunbookReceivers)); err != nil {
			return fmt.Errorf("Error setting `automation_runbook_receiver`: %+v", err)
		}

		if err = d.Set("voice_receiver", flattenMonitorActionGroupVoiceReceiver(group.VoiceReceivers)); err != nil {
			return fmt.Errorf("Error setting `voice_receiver`: %+v", err)
		}

		if err = d.Set("logic_app_receiver", flattenMonitorActionGroupLogicAppReceiver(group.LogicAppReceivers)); err != nil {
			return fmt.Errorf("Error setting `logic_app_receiver`: %+v", err)
		}

		if err = d.Set("azure_function_receiver", flattenMonitorActionGroupAzureFunctionReceiver(group.AzureFunctionReceivers)); err != nil {
			return fmt.Errorf("Error setting `azure_function_receiver`: %+v", err)
		}
		if err = d.Set("arm_role_receiver", flattenMonitorActionGroupArmRoleReceiver(group.ArmRoleReceivers)); err != nil {
			return fmt.Errorf("Error setting `arm_role_receiver`: %+v", err)
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMonitorActionGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["actionGroups"]

	resp, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting action group %q (resource group %q): %+v", name, resGroup, err)
		}
	}

	return nil
}

func expandMonitorActionGroupEmailReceiver(v []interface{}) *[]insights.EmailReceiver {
	receivers := make([]insights.EmailReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := insights.EmailReceiver{
			Name:                 utils.String(val["name"].(string)),
			EmailAddress:         utils.String(val["email_address"].(string)),
			UseCommonAlertSchema: utils.Bool(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupItsmReceiver(v []interface{}) *[]insights.ItsmReceiver {
	receivers := make([]insights.ItsmReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := insights.ItsmReceiver{
			Name:                utils.String(val["name"].(string)),
			WorkspaceID:         utils.String(val["workspace_id"].(string)),
			ConnectionID:        utils.String(val["connection_id"].(string)),
			TicketConfiguration: utils.String(val["ticket_configuration"].(string)),
			Region:              utils.String(azure.NormalizeLocation(val["region"].(string))),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupAzureAppPushReceiver(v []interface{}) *[]insights.AzureAppPushReceiver {
	receivers := make([]insights.AzureAppPushReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := insights.AzureAppPushReceiver{
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
			Name:                 utils.String(val["name"].(string)),
			ServiceURI:           utils.String(val["service_uri"].(string)),
			UseCommonAlertSchema: utils.Bool(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupAutomationRunbookReceiver(v []interface{}) *[]insights.AutomationRunbookReceiver {
	receivers := make([]insights.AutomationRunbookReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := insights.AutomationRunbookReceiver{
			Name:                 utils.String(val["name"].(string)),
			AutomationAccountID:  utils.String(val["automation_account_id"].(string)),
			RunbookName:          utils.String(val["runbook_name"].(string)),
			WebhookResourceID:    utils.String(val["webhook_resource_id"].(string)),
			IsGlobalRunbook:      utils.Bool(val["is_global_runbook"].(bool)),
			ServiceURI:           utils.String(val["service_uri"].(string)),
			UseCommonAlertSchema: utils.Bool(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupVoiceReceiver(v []interface{}) *[]insights.VoiceReceiver {
	receivers := make([]insights.VoiceReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := insights.VoiceReceiver{
			Name:        utils.String(val["name"].(string)),
			CountryCode: utils.String(val["country_code"].(string)),
			PhoneNumber: utils.String(val["phone_number"].(string)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupLogicAppReceiver(v []interface{}) *[]insights.LogicAppReceiver {
	receivers := make([]insights.LogicAppReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := insights.LogicAppReceiver{
			Name:                 utils.String(val["name"].(string)),
			ResourceID:           utils.String(val["resource_id"].(string)),
			CallbackURL:          utils.String(val["callback_url"].(string)),
			UseCommonAlertSchema: utils.Bool(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupAzureFunctionReceiver(v []interface{}) *[]insights.AzureFunctionReceiver {
	receivers := make([]insights.AzureFunctionReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := insights.AzureFunctionReceiver{
			Name:                  utils.String(val["name"].(string)),
			FunctionAppResourceID: utils.String(val["function_app_resource_id"].(string)),
			FunctionName:          utils.String(val["function_name"].(string)),
			HTTPTriggerURL:        utils.String(val["http_trigger_url"].(string)),
			UseCommonAlertSchema:  utils.Bool(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupArmRoleReceiver(v []interface{}) *[]insights.ArmRoleReceiver {
	receivers := make([]insights.ArmRoleReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := insights.ArmRoleReceiver{
			Name:                 utils.String(val["name"].(string)),
			RoleID:               utils.String(val["role_id"].(string)),
			UseCommonAlertSchema: utils.Bool(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func flattenMonitorActionGroupEmailReceiver(receivers *[]insights.EmailReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})
			if receiver.Name != nil {
				val["name"] = *receiver.Name
			}
			if receiver.EmailAddress != nil {
				val["email_address"] = *receiver.EmailAddress
			}
			if receiver.UseCommonAlertSchema != nil {
				val["use_common_alert_schema"] = *receiver.UseCommonAlertSchema
			}
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupItsmReceiver(receivers *[]insights.ItsmReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})
			if receiver.Name != nil {
				val["name"] = *receiver.Name
			}
			if receiver.WorkspaceID != nil {
				val["workspace_id"] = *receiver.WorkspaceID
			}
			if receiver.ConnectionID != nil {
				val["connection_id"] = *receiver.ConnectionID
			}
			if receiver.TicketConfiguration != nil {
				val["ticket_configuration"] = *receiver.TicketConfiguration
			}
			if receiver.Region != nil {
				val["region"] = azure.NormalizeLocation(*receiver.Region)
			}
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupAzureAppPushReceiver(receivers *[]insights.AzureAppPushReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})
			if receiver.Name != nil {
				val["name"] = *receiver.Name
			}
			if receiver.EmailAddress != nil {
				val["email_address"] = *receiver.EmailAddress
			}
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupSmsReceiver(receivers *[]insights.SmsReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})
			if receiver.Name != nil {
				val["name"] = *receiver.Name
			}
			if receiver.CountryCode != nil {
				val["country_code"] = *receiver.CountryCode
			}
			if receiver.PhoneNumber != nil {
				val["phone_number"] = *receiver.PhoneNumber
			}

			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupWebHookReceiver(receivers *[]insights.WebhookReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})
			if receiver.Name != nil {
				val["name"] = *receiver.Name
			}
			if receiver.ServiceURI != nil {
				val["service_uri"] = *receiver.ServiceURI
			}
			if receiver.UseCommonAlertSchema != nil {
				val["use_common_alert_schema"] = *receiver.UseCommonAlertSchema
			}

			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupAutomationRunbookReceiver(receivers *[]insights.AutomationRunbookReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})
			if receiver.Name != nil {
				val["name"] = *receiver.Name
			}
			if receiver.AutomationAccountID != nil {
				val["automation_account_id"] = *receiver.AutomationAccountID
			}
			if receiver.RunbookName != nil {
				val["runbook_name"] = *receiver.RunbookName
			}
			if receiver.WebhookResourceID != nil {
				val["webhook_resource_id"] = *receiver.WebhookResourceID
			}
			if receiver.IsGlobalRunbook != nil {
				val["is_global_runbook"] = *receiver.IsGlobalRunbook
			}
			if receiver.ServiceURI != nil {
				val["service_uri"] = *receiver.ServiceURI
			}
			if receiver.UseCommonAlertSchema != nil {
				val["use_common_alert_schema"] = *receiver.UseCommonAlertSchema
			}
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupVoiceReceiver(receivers *[]insights.VoiceReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})
			if receiver.Name != nil {
				val["name"] = *receiver.Name
			}
			if receiver.CountryCode != nil {
				val["country_code"] = *receiver.CountryCode
			}
			if receiver.PhoneNumber != nil {
				val["phone_number"] = *receiver.PhoneNumber
			}
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupLogicAppReceiver(receivers *[]insights.LogicAppReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})
			if receiver.Name != nil {
				val["name"] = *receiver.Name
			}
			if receiver.ResourceID != nil {
				val["resource_id"] = *receiver.ResourceID
			}
			if receiver.CallbackURL != nil {
				val["callback_url"] = *receiver.CallbackURL
			}
			if receiver.UseCommonAlertSchema != nil {
				val["use_common_alert_schema"] = *receiver.UseCommonAlertSchema
			}
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupAzureFunctionReceiver(receivers *[]insights.AzureFunctionReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})
			if receiver.Name != nil {
				val["name"] = *receiver.Name
			}
			if receiver.FunctionAppResourceID != nil {
				val["function_app_resource_id"] = *receiver.FunctionAppResourceID
			}
			if receiver.FunctionName != nil {
				val["function_name"] = *receiver.FunctionName
			}
			if receiver.HTTPTriggerURL != nil {
				val["http_trigger_url"] = *receiver.HTTPTriggerURL
			}
			if receiver.UseCommonAlertSchema != nil {
				val["use_common_alert_schema"] = *receiver.UseCommonAlertSchema
			}
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupArmRoleReceiver(receivers *[]insights.ArmRoleReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})
			if receiver.Name != nil {
				val["name"] = *receiver.Name
			}
			if receiver.RoleID != nil {
				val["role_id"] = *receiver.RoleID
			}
			if receiver.UseCommonAlertSchema != nil {
				val["use_common_alert_schema"] = *receiver.UseCommonAlertSchema
			}
			result = append(result, val)
		}
	}
	return result
}
