// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_action_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-01-01/actiongroupsapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorActionGroupCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.ActionGroupsClient
	tenantId := meta.(*clients.Client).Account.TenantId
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	location := location.Normalize(d.Get("location").(string))
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := actiongroupsapis.NewActionGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
	}

	shortName := d.Get("short_name").(string)
	enabled := d.Get("enabled").(bool)

	emailReceiversRaw := d.Get("email_receiver").([]any)
	itsmReceiversRaw := d.Get("itsm_receiver").([]any)
	azureAppPushReceiversRaw := d.Get("azure_app_push_receiver").([]any)
	smsReceiversRaw := d.Get("sms_receiver").([]any)
	webhookReceiversRaw := d.Get("webhook_receiver").([]any)
	automationRunbookReceiversRaw := d.Get("automation_runbook_receiver").([]any)
	voiceReceiversRaw := d.Get("voice_receiver").([]any)
	logicAppReceiversRaw := d.Get("logic_app_receiver").([]any)
	azureFunctionReceiversRaw := d.Get("azure_function_receiver").([]any)
	armRoleReceiversRaw := d.Get("arm_role_receiver").([]any)
	eventHubReceiversRaw := d.Get("event_hub_receiver").([]any)

	expandedItsmReceiver, err := expandMonitorActionGroupItsmReceiver(itsmReceiversRaw)
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]any)

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
		Tags: pluginsdk.ExpandPtrMapStringString(t),
	}

	if _, err := client.ActionGroupsCreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceMonitorActionGroupRead(d, meta)
}
