// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_runbook

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/jobschedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runbook"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationRunbookRead(d *pluginsdk.ResourceData, meta any) error {
	autoCli := meta.(*clients.Client).Automation
	client := autoCli.Runbook
	jsClient := autoCli.JobSchedule
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := runbook.ParseRunbookID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on AzureRM Automation Runbook %q (Account %q / Resource Group %q): %+v", id.RunbookName, id.AutomationAccountName, id.ResourceGroupName, err)
	}

	d.Set("name", id.RunbookName)
	d.Set("resource_group_name", id.ResourceGroupName)
	model := resp.Model
	d.Set("location", location.Normalize(model.Location))

	d.Set("automation_account_name", id.AutomationAccountName)
	if props := model.Properties; props != nil {
		d.Set("log_verbose", props.LogVerbose)
		d.Set("log_progress", props.LogProgress)
		d.Set("runbook_type", pointer.FromEnum(props.RunbookType))
		d.Set("description", props.Description)
		d.Set("log_activity_trace_level", props.LogActivityTrace)
		d.Set("runtime_environment_name", pointer.From(props.RuntimeEnvironment))
	}

	// GetContent need to use preview version client RunbookClientHack
	// move to stable Runbook once this issue fixed: https://github.com/Azure/azure-sdk-for-go/issues/17591#issuecomment-1233676539
	contentResp, err := autoCli.Runbook.GetContent(ctx, *id)
	if err != nil {
		if response.WasNotFound(contentResp.HttpResponse) {
			d.Set("content", "")
		} else {
			return fmt.Errorf("retrieving content for Automation Runbook %s: %+v", id, err)
		}
	}
	d.Set("content", string(pointer.From(contentResp.Model)))

	jsMap := make(map[uuid.UUID]jobschedule.JobScheduleProperties)
	automationAccountId := jobschedule.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName)

	filter := fmt.Sprintf("properties/runbook/name eq '%s'", id.RunbookName)
	jsIterator, err := jsClient.ListByAutomationAccount(ctx, automationAccountId, jobschedule.ListByAutomationAccountOperationOptions{Filter: &filter})
	if err != nil {
		return fmt.Errorf("loading Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
	}
	for _, item := range pointer.From(jsIterator.Model) {
		if itemProps := item.Properties; itemProps != nil {
			if itemProps.JobScheduleId == nil || *itemProps.JobScheduleId == "" {
				return fmt.Errorf("job schedule Id is nil or empty listed by Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
			}
			jsId, err := uuid.FromString(*itemProps.JobScheduleId)
			if err != nil {
				return fmt.Errorf("parsing job schedule Id listed by Automation Account %q Job Schedule List: %v", id.AutomationAccountName, err)
			}
			// get job schedule from GET API, `ListByAutomationAccountComplete` lost parameters
			jobscheduleID, err := jobschedule.ParseJobScheduleID(pointer.From(item.Id))
			if err != nil {
				return fmt.Errorf("parsing job schedule Id listed by Automation Account %q Job Schedule List: %v", id.AutomationAccountName, err)
			}
			jsResult, err := jsClient.Get(ctx, *jobscheduleID)
			if err != nil {
				return fmt.Errorf("retrieving job schedule by %s: %v", *jobscheduleID, err)
			}
			if jsResult.Model != nil && jsResult.Model.Properties != nil {
				jsMap[jsId] = *jsResult.Model.Properties
			}
		}
	}

	if err := d.Set("job_schedule", helper.FlattenAutomationJobSchedule(jsMap)); err != nil {
		return fmt.Errorf("setting `job_schedule`: %+v", err)
	}

	if err := tags.FlattenAndSet(d, model.Tags); err != nil {
		return err
	}

	return nil
}

// if job in jsIterator but not in jsMap, then delete it
// if job in both jsIterator and jsMap, remove the entry in jsMap
// at last, create jobs still in jsMap
func updatedLinkedJobSchedules(ctx context.Context, subscriptionID string, client *jobschedule.JobScheduleClient, id *runbook.RunbookId, jsMap map[string]jobschedule.JobScheduleCreateParameters) error {
	automationAccountId := jobschedule.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName)
	filter := fmt.Sprintf("properties/runbook/name eq '%s'", id.RunbookName)
	jsIterator, err := client.ListByAutomationAccount(ctx, automationAccountId, jobschedule.ListByAutomationAccountOperationOptions{Filter: &filter})
	if err != nil {
		return fmt.Errorf("loading Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
	}

	for _, item := range pointer.From(jsIterator.Model) {
		prop := item.Properties
		jobDigest := helper.ResourceAutomationJobScheduleDigest(prop)

		if _, ok := jsMap[jobDigest]; ok {
			delete(jsMap, jobDigest)
		} else {
			if prop == nil || prop.JobScheduleId == nil || *prop.JobScheduleId == "" {
				return fmt.Errorf("job schedule Id is nil or empty listed by %s Job Schedule List: %+v", id, err)
			}
			parsedId := jobschedule.NewJobScheduleID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, pointer.From(item.Properties.JobScheduleId))
			if resp, err := client.Delete(ctx, parsedId); err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("deleting job schedule Id listed by %s Job Schedule List:%v", id, err)
				}
			}
		}
	}

	// create jobs still in jsMap
	for _, js := range jsMap {
		// skip if the schedule name is empty
		if pointer.From(js.Properties.Schedule.Name) == "" {
			continue
		}
		jsuuid, err := uuid.NewV4()
		if err != nil {
			return fmt.Errorf("creating job schedule Id(UUID) for %s: %+v", id, err)
		}

		jsId := jobschedule.NewJobScheduleID(subscriptionID, id.ResourceGroupName, id.AutomationAccountName, jsuuid.String())
		if _, err := client.Create(ctx, jsId, js); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}
	}

	return nil
}
