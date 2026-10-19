// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_job_schedule

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/jobschedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationJobScheduleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.JobSchedule
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// the jobSchedule ID may be updated by Runbook, so need to get the real id by list API
	tfID, err := commonids.ParseCompositeResourceID(d.Id(), &schedule.ScheduleId{}, &runbook.RunbookId{})
	if err != nil {
		return err
	}

	js, err := GetJobScheduleFromTFID(ctx, client, tfID)
	if err != nil {
		return err
	}
	if js == nil {
		d.SetId("")
		return nil
	}

	id, err := jobschedule.ParseJobScheduleID(pointer.From(js.Id))
	if err != nil {
		return err
	}

	d.Set("resource_manager_id", id.ID())
	d.Set("job_schedule_id", id.JobScheduleId)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	// The response from the list API has no parameter field, so use Get API to get the JobSchedule
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return err
	}

	if resp.Model != nil && resp.Model.Properties != nil {
		props := resp.Model.Properties
		d.Set("runbook_name", props.Runbook.Name)
		d.Set("schedule_name", props.Schedule.Name)

		if v := props.RunOn; v != nil {
			d.Set("run_on", v)
		}

		if props.Parameters != nil {
			if v := *props.Parameters; v != nil {
				jsParameters := make(map[string]any)
				for key, value := range v {
					jsParameters[strings.ToLower(key)] = value
				}
				d.Set("parameters", jsParameters)
			}
		}
	}

	return nil
}

func GetJobScheduleFromTFID(ctx context.Context, client *jobschedule.JobScheduleClient, id *commonids.CompositeResourceID[*schedule.ScheduleId, *runbook.RunbookId]) (js *jobschedule.JobSchedule, err error) {
	accountID := jobschedule.NewAutomationAccountID(id.First.SubscriptionId, id.First.ResourceGroupName, id.First.AutomationAccountName)
	filter := fmt.Sprintf("properties/schedule/name eq '%s' and properties/runbook/name eq '%s'", id.First.ScheduleName, id.Second.RunbookName)
	jsList, err := client.ListByAutomationAccountComplete(ctx, accountID, jobschedule.ListByAutomationAccountOperationOptions{Filter: &filter})
	if err != nil {
		if response.WasNotFound(jsList.LatestHttpResponse) {
			return nil, nil
		}
		return nil, fmt.Errorf("loading Automation Account %q Job Schedule List: %+v", accountID.AutomationAccountName, err)
	}

	if len(jsList.Items) > 0 {
		js = &jsList.Items[0]
	}

	return js, nil
}
