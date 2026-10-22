// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_action_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-01-01/actiongroupsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceMonitorActionGroupFlatten(d *pluginsdk.ResourceData, id *actiongroupsapis.ActionGroupId, model *actiongroupsapis.ActionGroupResource) error {
	var err error
	d.Set("name", id.ActionGroupName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
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
		if err = d.Set("tags", pluginsdk.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenMonitorActionGroupEmailReceiver(receivers *[]actiongroupsapis.EmailReceiver) []any {
	result := make([]any, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]any)

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

func flattenMonitorActionGroupItsmReceiver(receivers *[]actiongroupsapis.ItsmReceiver) []any {
	result := make([]any, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]any)

			val["name"] = receiver.Name
			val["workspace_id"] = receiver.WorkspaceId
			val["connection_id"] = receiver.ConnectionId
			val["ticket_configuration"] = receiver.TicketConfiguration
			val["region"] = location.Normalize(receiver.Region)

			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupAzureAppPushReceiver(receivers *[]actiongroupsapis.AzureAppPushReceiver) []any {
	result := make([]any, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]any)

			val["name"] = receiver.Name
			val["email_address"] = receiver.EmailAddress

			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupSmsReceiver(receivers *[]actiongroupsapis.SmsReceiver) []any {
	result := make([]any, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]any)

			val["name"] = receiver.Name
			val["country_code"] = receiver.CountryCode
			val["phone_number"] = receiver.PhoneNumber

			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupWebHookReceiver(receivers *[]actiongroupsapis.WebhookReceiver) []any {
	result := make([]any, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			result = append(result, map[string]any{
				"name":                    receiver.Name,
				"service_uri":             receiver.ServiceUri,
				"use_common_alert_schema": pointer.From(receiver.UseCommonAlertSchema),
				"aad_auth":                flattenMonitorActionGroupSecureWebHookReceiver(receiver),
			})
		}
	}
	return result
}

func flattenMonitorActionGroupSecureWebHookReceiver(receiver actiongroupsapis.WebhookReceiver) []any {
	if receiver.UseAadAuth == nil || !*receiver.UseAadAuth {
		return []any{}
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
	return []any{
		map[string]any{
			"object_id":      objectId,
			"identifier_uri": identifierUri,
			"tenant_id":      tenantId,
		},
	}
}

func flattenMonitorActionGroupAutomationRunbookReceiver(receivers *[]actiongroupsapis.AutomationRunbookReceiver) []any {
	result := make([]any, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]any)
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

func flattenMonitorActionGroupVoiceReceiver(receivers *[]actiongroupsapis.VoiceReceiver) []any {
	result := make([]any, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]any)

			val["name"] = receiver.Name
			val["country_code"] = receiver.CountryCode
			val["phone_number"] = receiver.PhoneNumber

			result = append(result, val)
		}
	}
	return result
}

func flattenMonitorActionGroupLogicAppReceiver(receivers *[]actiongroupsapis.LogicAppReceiver) []any {
	result := make([]any, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]any)

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

func flattenMonitorActionGroupAzureFunctionReceiver(receivers *[]actiongroupsapis.AzureFunctionReceiver) []any {
	result := make([]any, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]any)

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

func flattenMonitorActionGroupRoleReceiver(receivers *[]actiongroupsapis.ArmRoleReceiver) []any {
	result := make([]any, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]any)

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

func flattenMonitorActionGroupEventHubReceiver(receivers *[]actiongroupsapis.EventHubReceiver) []any {
	result := make([]any, 0)
	if receivers != nil {
		for _, receiver := range *receivers {
			val := make(map[string]any)

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
