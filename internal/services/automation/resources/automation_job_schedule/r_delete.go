// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_job_schedule

import (
	"fmt"

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

func resourceAutomationJobScheduleDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.JobSchedule
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	tfID, err := commonids.ParseCompositeResourceID(d.Id(), &schedule.ScheduleId{}, &runbook.RunbookId{})
	if err != nil {
		return err
	}
	js, err := GetJobScheduleFromTFID(ctx, client, tfID)
	if err != nil {
		return err
	}

	if js == nil {
		return nil
	}

	id, err := jobschedule.ParseJobScheduleID(pointer.From(js.Id))
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
