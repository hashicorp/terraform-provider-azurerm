// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_action_group

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-01-01/actiongroupsapis"
)

func expandMonitorActionGroupEmailReceiver(v []any) *[]actiongroupsapis.EmailReceiver {
	receivers := make([]actiongroupsapis.EmailReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]any)
		receiver := actiongroupsapis.EmailReceiver{
			Name:                 val["name"].(string),
			EmailAddress:         val["email_address"].(string),
			UseCommonAlertSchema: pointer.To(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupItsmReceiver(v []any) (*[]actiongroupsapis.ItsmReceiver, error) {
	receivers := make([]actiongroupsapis.ItsmReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]any)
		ticketConfiguration := val["ticket_configuration"].(string)
		receiver := actiongroupsapis.ItsmReceiver{
			Name:                val["name"].(string),
			WorkspaceId:         val["workspace_id"].(string),
			ConnectionId:        val["connection_id"].(string),
			TicketConfiguration: ticketConfiguration,
			Region:              location.Normalize(val["region"].(string)),
		}

		// https://github.com/Azure/azure-rest-api-specs/issues/20488 ticket_configuration should have `PayloadRevision` and `WorkItemType` keys

		j := make(map[string]any)
		if err := json.Unmarshal([]byte(ticketConfiguration), &j); err != nil {
			return nil, fmt.Errorf("`itsm_receiver.ticket_configuration` %s unmarshall json error: %+v", ticketConfiguration, err)
		}

		_, existKeyPayloadRevision := j["PayloadRevision"]
		_, existKeyWorkItemType := j["WorkItemType"]
		if !existKeyPayloadRevision || !existKeyWorkItemType {
			return nil, fmt.Errorf("`itsm_receiver.ticket_configuration` should be JSON blob with `PayloadRevision` and `WorkItemType` keys")
		}
		receivers = append(receivers, receiver)
	}
	return &receivers, nil
}

func expandMonitorActionGroupAzureAppPushReceiver(v []any) *[]actiongroupsapis.AzureAppPushReceiver {
	receivers := make([]actiongroupsapis.AzureAppPushReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]any)
		receiver := actiongroupsapis.AzureAppPushReceiver{
			Name:         val["name"].(string),
			EmailAddress: val["email_address"].(string),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupSmsReceiver(v []any) *[]actiongroupsapis.SmsReceiver {
	receivers := make([]actiongroupsapis.SmsReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]any)
		receiver := actiongroupsapis.SmsReceiver{
			Name:        val["name"].(string),
			CountryCode: val["country_code"].(string),
			PhoneNumber: val["phone_number"].(string),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupWebHookReceiver(tenantId string, v []any) *[]actiongroupsapis.WebhookReceiver {
	receivers := make([]actiongroupsapis.WebhookReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]any)
		receiver := actiongroupsapis.WebhookReceiver{
			Name:                 val["name"].(string),
			ServiceUri:           val["service_uri"].(string),
			UseCommonAlertSchema: pointer.To(val["use_common_alert_schema"].(bool)),
		}
		if v, ok := val["aad_auth"].([]any); ok && len(v) > 0 {
			secureWebhook := v[0].(map[string]any)
			receiver.UseAadAuth = pointer.To(true)
			receiver.ObjectId = pointer.To(secureWebhook["object_id"].(string))
			receiver.IdentifierUri = pointer.To(secureWebhook["identifier_uri"].(string))
			if v := secureWebhook["tenant_id"].(string); v != "" {
				receiver.TenantId = pointer.To(v)
			} else {
				receiver.TenantId = pointer.To(tenantId)
			}
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupAutomationRunbookReceiver(v []any) *[]actiongroupsapis.AutomationRunbookReceiver {
	receivers := make([]actiongroupsapis.AutomationRunbookReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]any)
		receiver := actiongroupsapis.AutomationRunbookReceiver{
			Name:                 pointer.To(val["name"].(string)),
			AutomationAccountId:  val["automation_account_id"].(string),
			RunbookName:          val["runbook_name"].(string),
			WebhookResourceId:    val["webhook_resource_id"].(string),
			IsGlobalRunbook:      val["is_global_runbook"].(bool),
			ServiceUri:           pointer.To(val["service_uri"].(string)),
			UseCommonAlertSchema: pointer.To(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupVoiceReceiver(v []any) *[]actiongroupsapis.VoiceReceiver {
	receivers := make([]actiongroupsapis.VoiceReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]any)
		receiver := actiongroupsapis.VoiceReceiver{
			Name:        val["name"].(string),
			CountryCode: val["country_code"].(string),
			PhoneNumber: val["phone_number"].(string),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupLogicAppReceiver(v []any) *[]actiongroupsapis.LogicAppReceiver {
	receivers := make([]actiongroupsapis.LogicAppReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]any)
		receiver := actiongroupsapis.LogicAppReceiver{
			Name:                 val["name"].(string),
			ResourceId:           val["resource_id"].(string),
			CallbackURL:          val["callback_url"].(string),
			UseCommonAlertSchema: pointer.To(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupAzureFunctionReceiver(v []any) *[]actiongroupsapis.AzureFunctionReceiver {
	receivers := make([]actiongroupsapis.AzureFunctionReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]any)
		receiver := actiongroupsapis.AzureFunctionReceiver{
			Name:                  val["name"].(string),
			FunctionAppResourceId: val["function_app_resource_id"].(string),
			FunctionName:          val["function_name"].(string),
			HTTPTriggerURL:        val["http_trigger_url"].(string),
			UseCommonAlertSchema:  pointer.To(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupRoleReceiver(v []any) *[]actiongroupsapis.ArmRoleReceiver {
	receivers := make([]actiongroupsapis.ArmRoleReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]any)
		receiver := actiongroupsapis.ArmRoleReceiver{
			Name:                 val["name"].(string),
			RoleId:               val["role_id"].(string),
			UseCommonAlertSchema: pointer.To(val["use_common_alert_schema"].(bool)),
		}
		receivers = append(receivers, receiver)
	}
	return &receivers
}

func expandMonitorActionGroupEventHubReceiver(tenantId string, subscriptionId string, v []any) *[]actiongroupsapis.EventHubReceiver {
	receivers := make([]actiongroupsapis.EventHubReceiver, 0)
	for _, receiverValue := range v {
		val := receiverValue.(map[string]any)

		eventHubNameSpace, eventHubName, subId := val["event_hub_namespace"].(string), val["event_hub_name"].(string), val["subscription_id"].(string)

		receiver := actiongroupsapis.EventHubReceiver{
			EventHubNameSpace:    eventHubNameSpace,
			EventHubName:         eventHubName,
			Name:                 val["name"].(string),
			UseCommonAlertSchema: pointer.To(val["use_common_alert_schema"].(bool)),
		}

		if v := val["tenant_id"].(string); v != "" {
			receiver.TenantId = pointer.To(v)
		} else {
			receiver.TenantId = pointer.To(tenantId)
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
