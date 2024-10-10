// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/automationaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-01-01/actiongroupsapis"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorActionGroup() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceMonitorActionGroupCreateUpdate,
		Read:   resourceMonitorActionGroupRead,
		Update: resourceMonitorActionGroupCreateUpdate,
		Delete: resourceMonitorActionGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := actiongroupsapis.ParseActionGroupID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ActionGroupUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "global",
				ValidateFunc: validation.Any(
					location.EnhancedValidate,
					validation.StringInSlice([]string{
						"global",
					}, false),
				),
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},

			"short_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 12),
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"email_receiver": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"email_address": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"use_common_alert_schema": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"itsm_receiver": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"workspace_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.WorkspaceID,
						},
						"connection_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"ticket_configuration": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
						},
						"region": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							ValidateFunc:     validation.StringIsNotEmpty,
							DiffSuppressFunc: location.DiffSuppressFunc,
						},
					},
				},
			},

			"azure_app_push_receiver": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"email_address": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"sms_receiver": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"country_code": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"phone_number": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"webhook_receiver": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"service_uri": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
						},
						"use_common_alert_schema": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"aad_auth": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"object_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.IsUUID,
									},

									"identifier_uri": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.IsURLWithScheme([]string{"api", "https"}),
									},

									"tenant_id": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Computed:     true,
										ValidateFunc: validation.IsUUID,
									},
								},
							},
						},
					},
				},
			},

			"automation_runbook_receiver": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"automation_account_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: automationaccount.ValidateAutomationAccountID,
						},
						"runbook_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"webhook_resource_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"is_global_runbook": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},
						"service_uri": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
						},
						"use_common_alert_schema": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"voice_receiver": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"country_code": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"phone_number": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"logic_app_receiver": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"resource_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"callback_url": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
						},
						"use_common_alert_schema": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"azure_function_receiver": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						// TODO: this should be `_id` and not `_resource_id`
						"function_app_resource_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: commonids.ValidateFunctionAppID,
						},
						"function_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"http_trigger_url": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsURLWithScheme([]string{"http", "https"}),
						},
						"use_common_alert_schema": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"arm_role_receiver": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"role_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"use_common_alert_schema": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"event_hub_receiver": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"event_hub_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"event_hub_namespace": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"tenant_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IsUUID,
						},
						"use_common_alert_schema": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
						"subscription_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}

	return resource
}

func resourceMonitorActionGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionGroupsClient
	tenantId := meta.(*clients.Client).Account.TenantId
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	location := location.Normalize(d.Get("location").(string))
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := actiongroupsapis.NewActionGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.ActionGroupsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_action_group", id.ID())
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
	eventHubReceiversRaw := d.Get("event_hub_receiver").([]interface{})

	expandedItsmReceiver, err := expandMonitorActionGroupItsmReceiver(itsmReceiversRaw)
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]interface{})

	parameters := actiongroupsapis.ActionGroupResource{
		Location: location,
		Properties: &actiongroupsapis.ActionGroup{
			GroupShortName:             shortName,
			Enabled:                    enabled,
			EmailReceivers:             expandMonitorActionGroupEmailReceiver(emailReceiversRaw),
			AzureAppPushReceivers:      expandMonitorActionGroupAzureAppPushReceiver(azureAppPushReceiversRaw),
			ItsmReceivers:              expandedItsmReceiver,
			SmsReceivers:               expandMonitorActionGroupSmsReceiver(smsReceiversRaw),
			WebhookReceivers:           expandMonitorActionGroupWebHookReceiver(tenantId, webhookReceiversRaw),
			AutomationRunbookReceivers: expandMonitorActionGroupAutomationRunbookReceiver(automationRunbookReceiversRaw),
			VoiceReceivers:             expandMonitorActionGroupVoiceReceiver(voiceReceiversRaw),
			LogicAppReceivers:          expandMonitorActionGroupLogicAppReceiver(logicAppReceiversRaw),
			AzureFunctionReceivers:     expandMonitorActionGroupAzureFunctionReceiver(azureFunctionReceiversRaw),
			ArmRoleReceivers:           expandMonitorActionGroupRoleReceiver(armRoleReceiversRaw),
			EventHubReceivers:          expandMonitorActionGroupEventHubReceiver(tenantId, subscriptionId, eventHubReceiversRaw),
		},
		Tags: utils.ExpandPtrMapStringString(t),
	}

	if _, err := client.ActionGroupsCreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorActionGroupRead(d, meta)
}

func resourceMonitorActionGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := actiongroupsapis.ParseActionGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ActionGroupsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ActionGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {

		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("short_name", props.GroupShortName)
			d.Set("enabled", props.Enabled)

			if err = d.Set("email_receiver", flattenMonitorActionGroupEmailReceiver(props.EmailReceivers)); err != nil {
				return fmt.Errorf("setting `email_receiver`: %+v", err)
			}

			if err = d.Set("itsm_receiver", flattenMonitorActionGroupItsmReceiver(props.ItsmReceivers)); err != nil {
				return fmt.Errorf("setting `itsm_receiver`: %+v", err)
			}

			if err = d.Set("azure_app_push_receiver", flattenMonitorActionGroupAzureAppPushReceiver(props.AzureAppPushReceivers)); err != nil {
				return fmt.Errorf("setting `azure_app_push_receiver`: %+v", err)
			}

			if err = d.Set("sms_receiver", flattenMonitorActionGroupSmsReceiver(props.SmsReceivers)); err != nil {
				return fmt.Errorf("setting `sms_receiver`: %+v", err)
			}

			if err = d.Set("webhook_receiver", flattenMonitorActionGroupWebHookReceiver(props.WebhookReceivers)); err != nil {
				return fmt.Errorf("setting `webhook_receiver`: %+v", err)
			}

			if err = d.Set("automation_runbook_receiver", flattenMonitorActionGroupAutomationRunbookReceiver(props.AutomationRunbookReceivers)); err != nil {
				return fmt.Errorf("setting `automation_runbook_receiver`: %+v", err)
			}

			if err = d.Set("voice_receiver", flattenMonitorActionGroupVoiceReceiver(props.VoiceReceivers)); err != nil {
				return fmt.Errorf("setting `voice_receiver`: %+v", err)
			}

			if err = d.Set("logic_app_receiver", flattenMonitorActionGroupLogicAppReceiver(props.LogicAppReceivers)); err != nil {
				return fmt.Errorf("setting `logic_app_receiver`: %+v", err)
			}

			if err = d.Set("azure_function_receiver", flattenMonitorActionGroupAzureFunctionReceiver(props.AzureFunctionReceivers)); err != nil {
				return fmt.Errorf("setting `azure_function_receiver`: %+v", err)
			}
			if err = d.Set("arm_role_receiver", flattenMonitorActionGroupRoleReceiver(props.ArmRoleReceivers)); err != nil {
				return fmt.Errorf("setting `arm_role_receiver`: %+v", err)
			}
			if err = d.Set("event_hub_receiver", flattenMonitorActionGroupEventHubReceiver(props.EventHubReceivers)); err != nil {
				return fmt.Errorf("setting `event_hub_receiver`: %+v", err)
			}
		}
		if err = d.Set("tags", utils.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}
	}
	return nil
}

func resourceMonitorActionGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := actiongroupsapis.ParseActionGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ActionGroupsDelete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func expandMonitorActionGroupEmailReceiver(v []interface{}) *[]actiongroupsapis.EmailReceiver {
	receivers := make([]actiongroupsapis.EmailReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := actiongroupsapis.EmailReceiver{
			Name:                 val["name"].(string),
			EmailAddress:         val["email_address"].(string),
			UseCommonAlertSchema: utils.Bool(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupItsmReceiver(v []interface{}) (*[]actiongroupsapis.ItsmReceiver, error) {
	receivers := make([]actiongroupsapis.ItsmReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		ticketConfiguration := val["ticket_configuration"].(string)
		receiver := actiongroupsapis.ItsmReceiver{
			Name:                val["name"].(string),
			WorkspaceId:         val["workspace_id"].(string),
			ConnectionId:        val["connection_id"].(string),
			TicketConfiguration: ticketConfiguration,
			Region:              azure.NormalizeLocation(val["region"].(string)),
		}

		// https://github.com/Azure/azure-rest-api-specs/issues/20488 ticket_configuration should have `PayloadRevision` and `WorkItemType` keys

		j := make(map[string]interface{})
		err := json.Unmarshal([]byte(ticketConfiguration), &j)
		if err != nil {
			return nil, fmt.Errorf("`itsm_receiver.ticket_configuration` %s unmarshall json error: %+v", ticketConfiguration, err)
		}

		_, existKeyPayloadRevision := j["PayloadRevision"]
		_, existKeyWorkItemType := j["WorkItemType"]
		if !(existKeyPayloadRevision && existKeyWorkItemType) {
			return nil, fmt.Errorf("`itsm_receiver.ticket_configuration` should be JSON blob with `PayloadRevision` and `WorkItemType` keys")
		}
		receivers = append(receivers, receiver)
	}
	return &receivers, nil
}

func expandMonitorActionGroupAzureAppPushReceiver(v []interface{}) *[]actiongroupsapis.AzureAppPushReceiver {
	receivers := make([]actiongroupsapis.AzureAppPushReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := actiongroupsapis.AzureAppPushReceiver{
			Name:         val["name"].(string),
			EmailAddress: val["email_address"].(string),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupSmsReceiver(v []interface{}) *[]actiongroupsapis.SmsReceiver {
	receivers := make([]actiongroupsapis.SmsReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := actiongroupsapis.SmsReceiver{
			Name:        val["name"].(string),
			CountryCode: val["country_code"].(string),
			PhoneNumber: val["phone_number"].(string),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupWebHookReceiver(tenantId string, v []interface{}) *[]actiongroupsapis.WebhookReceiver {
	receivers := make([]actiongroupsapis.WebhookReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := actiongroupsapis.WebhookReceiver{
			Name:                 val["name"].(string),
			ServiceUri:           val["service_uri"].(string),
			UseCommonAlertSchema: utils.Bool(val["use_common_alert_schema"].(bool)),
		}
		if v, ok := val["aad_auth"].([]interface{}); ok && len(v) > 0 {
			secureWebhook := v[0].(map[string]interface{})
			receiver.UseAadAuth = utils.Bool(true)
			receiver.ObjectId = utils.String(secureWebhook["object_id"].(string))
			receiver.IdentifierUri = utils.String(secureWebhook["identifier_uri"].(string))
			if v := secureWebhook["tenant_id"].(string); v != "" {
				receiver.TenantId = utils.String(v)
			} else {
				receiver.TenantId = utils.String(tenantId)
			}
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupAutomationRunbookReceiver(v []interface{}) *[]actiongroupsapis.AutomationRunbookReceiver {
	receivers := make([]actiongroupsapis.AutomationRunbookReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := actiongroupsapis.AutomationRunbookReceiver{
			Name:                 utils.String(val["name"].(string)),
			AutomationAccountId:  val["automation_account_id"].(string),
			RunbookName:          val["runbook_name"].(string),
			WebhookResourceId:    val["webhook_resource_id"].(string),
			IsGlobalRunbook:      val["is_global_runbook"].(bool),
			ServiceUri:           utils.String(val["service_uri"].(string)),
			UseCommonAlertSchema: utils.Bool(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupVoiceReceiver(v []interface{}) *[]actiongroupsapis.VoiceReceiver {
	receivers := make([]actiongroupsapis.VoiceReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := actiongroupsapis.VoiceReceiver{
			Name:        val["name"].(string),
			CountryCode: val["country_code"].(string),
			PhoneNumber: val["phone_number"].(string),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupLogicAppReceiver(v []interface{}) *[]actiongroupsapis.LogicAppReceiver {
	receivers := make([]actiongroupsapis.LogicAppReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := actiongroupsapis.LogicAppReceiver{
			Name:                 val["name"].(string),
			ResourceId:           val["resource_id"].(string),
			CallbackURL:          val["callback_url"].(string),
			UseCommonAlertSchema: utils.Bool(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupAzureFunctionReceiver(v []interface{}) *[]actiongroupsapis.AzureFunctionReceiver {
	receivers := make([]actiongroupsapis.AzureFunctionReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := actiongroupsapis.AzureFunctionReceiver{
			Name:                  val["name"].(string),
			FunctionAppResourceId: val["function_app_resource_id"].(string),
			FunctionName:          val["function_name"].(string),
			HTTPTriggerURL:        val["http_trigger_url"].(string),
			UseCommonAlertSchema:  utils.Bool(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupRoleReceiver(v []interface{}) *[]actiongroupsapis.ArmRoleReceiver {
	receivers := make([]actiongroupsapis.ArmRoleReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})
		receiver := actiongroupsapis.ArmRoleReceiver{
			Name:                 val["name"].(string),
			RoleId:               val["role_id"].(string),
			UseCommonAlertSchema: utils.Bool(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupEventHubReceiver(tenantId string, subscriptionId string, v []interface{}) *[]actiongroupsapis.EventHubReceiver {
	receivers := make([]actiongroupsapis.EventHubReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]interface{})

		eventHubNameSpace, eventHubName, subId := val["event_hub_namespace"].(string), val["event_hub_name"].(string), val["subscription_id"].(string)

		receiver := actiongroupsapis.EventHubReceiver{
			EventHubNameSpace:    eventHubNameSpace,
			EventHubName:         eventHubName,
			Name:                 val["name"].(string),
			UseCommonAlertSchema: utils.Bool(val["use_common_alert_schema"].(bool)),
		}

		if v := val["tenant_id"].(string); v != "" {
			receiver.TenantId = utils.String(v)
		} else {
			receiver.TenantId = utils.String(tenantId)
		}

		if subId != "" {
			receiver.SubscriptionId = subId
		} else {
			receiver.SubscriptionId = subscriptionId
		}

		receivers = append(receivers, receiver)
	}
	return &receivers
}

func flattenMonitorActionGroupEmailReceiver(receivers *[]actiongroupsapis.EmailReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})

			val["name"] = receiver.Name
			val["email_address"] = receiver.EmailAddress

			if receiver.UseCommonAlertSchema != nil {
				val["use_common_alert_schema"] = *receiver.UseCommonAlertSchema
			}
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupItsmReceiver(receivers *[]actiongroupsapis.ItsmReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})

			val["name"] = receiver.Name
			val["workspace_id"] = receiver.WorkspaceId
			val["connection_id"] = receiver.ConnectionId
			val["ticket_configuration"] = receiver.TicketConfiguration
			val["region"] = azure.NormalizeLocation(receiver.Region)

			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupAzureAppPushReceiver(receivers *[]actiongroupsapis.AzureAppPushReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})

			val["name"] = receiver.Name
			val["email_address"] = receiver.EmailAddress

			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupSmsReceiver(receivers *[]actiongroupsapis.SmsReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})

			val["name"] = receiver.Name
			val["country_code"] = receiver.CountryCode
			val["phone_number"] = receiver.PhoneNumber

			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupWebHookReceiver(receivers *[]actiongroupsapis.WebhookReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			var useCommonAlert bool
			if receiver.UseCommonAlertSchema != nil {
				useCommonAlert = *receiver.UseCommonAlertSchema
			}

			result = append(result, map[string]interface{}{
				"name":                    receiver.Name,
				"service_uri":             receiver.ServiceUri,
				"use_common_alert_schema": useCommonAlert,
				"aad_auth":                flattenMonitorActionGroupSecureWebHookReceiver(receiver),
			})
		}
	}
	return result
}

func flattenMonitorActionGroupSecureWebHookReceiver(receiver actiongroupsapis.WebhookReceiver) []interface{} {
	if receiver.UseAadAuth == nil || !*receiver.UseAadAuth {
		return []interface{}{}
	}

	var objectId, identifierUri, tenantId string

	if v := receiver.ObjectId; v != nil {
		objectId = *v
	}
	if v := receiver.IdentifierUri; v != nil {
		identifierUri = *v
	}
	if v := receiver.TenantId; v != nil {
		tenantId = *v
	}
	return []interface{}{
		map[string]interface{}{
			"object_id":      objectId,
			"identifier_uri": identifierUri,
			"tenant_id":      tenantId,
		},
	}
}

func flattenMonitorActionGroupAutomationRunbookReceiver(receivers *[]actiongroupsapis.AutomationRunbookReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})
			if receiver.Name != nil {
				val["name"] = *receiver.Name
			}

			val["automation_account_id"] = receiver.AutomationAccountId
			val["runbook_name"] = receiver.RunbookName
			val["webhook_resource_id"] = receiver.WebhookResourceId
			val["is_global_runbook"] = receiver.IsGlobalRunbook

			if receiver.ServiceUri != nil {
				val["service_uri"] = *receiver.ServiceUri
			}
			if receiver.UseCommonAlertSchema != nil {
				val["use_common_alert_schema"] = *receiver.UseCommonAlertSchema
			}
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupVoiceReceiver(receivers *[]actiongroupsapis.VoiceReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})

			val["name"] = receiver.Name
			val["country_code"] = receiver.CountryCode
			val["phone_number"] = receiver.PhoneNumber

			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupLogicAppReceiver(receivers *[]actiongroupsapis.LogicAppReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})

			val["name"] = receiver.Name
			val["resource_id"] = receiver.ResourceId
			val["callback_url"] = receiver.CallbackURL

			if receiver.UseCommonAlertSchema != nil {
				val["use_common_alert_schema"] = *receiver.UseCommonAlertSchema
			}
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupAzureFunctionReceiver(receivers *[]actiongroupsapis.AzureFunctionReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})

			val["name"] = receiver.Name
			val["function_app_resource_id"] = receiver.FunctionAppResourceId
			val["function_name"] = receiver.FunctionName
			val["http_trigger_url"] = receiver.HTTPTriggerURL

			if receiver.UseCommonAlertSchema != nil {
				val["use_common_alert_schema"] = *receiver.UseCommonAlertSchema
			}
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupRoleReceiver(receivers *[]actiongroupsapis.ArmRoleReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})

			val["name"] = receiver.Name
			val["role_id"] = receiver.RoleId

			if receiver.UseCommonAlertSchema != nil {
				val["use_common_alert_schema"] = *receiver.UseCommonAlertSchema
			}
			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupEventHubReceiver(receivers *[]actiongroupsapis.EventHubReceiver) []interface{} {
	result := make([]interface{}, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]interface{})

			val["name"] = receiver.Name

			eventHubNamespace := receiver.EventHubNameSpace
			eventHubName := receiver.EventHubName
			subscriptionId := receiver.SubscriptionId
			val["subscription_id"], val["event_hub_namespace"], val["event_hub_name"] = subscriptionId, eventHubNamespace, eventHubName

			if receiver.UseCommonAlertSchema != nil {
				val["use_common_alert_schema"] = *receiver.UseCommonAlertSchema
			}
			if receiver.TenantId != nil {
				val["tenant_id"] = *receiver.TenantId
			}
			result = append(result, val)
		}
	}
	return result
}
