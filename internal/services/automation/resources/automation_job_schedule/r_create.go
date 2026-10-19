// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation_job_schedule

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/jobschedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationJobScheduleCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Automation.JobSchedule
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("automation_account_name").(string)
	runbookName := d.Get("runbook_name").(string)
	scheduleName := d.Get("schedule_name").(string)

	jobScheduleUUID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	if jobScheduleID, ok := d.GetOk("job_schedule_id"); ok {
		jobScheduleUUID = uuid.FromStringOrNil(jobScheduleID.(string))
	}

	scheduleID := schedule.NewScheduleID(subscriptionId, resourceGroup, accountName, scheduleName)
	runbookID := runbook.NewRunbookID(subscriptionId, resourceGroup, accountName, runbookName)
	id := jobschedule.NewJobScheduleID(subscriptionId, resourceGroup, accountName, jobScheduleUUID.String())

	tfID := &commonids.CompositeResourceID[*schedule.ScheduleId, *runbook.RunbookId]{
		First:  &scheduleID,
		Second: &runbookID,
	}

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := GetJobScheduleFromTFID(ctx, client, tfID)
		if err != nil {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}

		if existing != nil {
			return tf.ImportAsExistsError("azurerm_automation_job_schedule", tfID.ID())
		}
	}

	parameters := jobschedule.JobScheduleCreateParameters{
		Properties: jobschedule.JobScheduleCreateProperties{
			Schedule: jobschedule.ScheduleAssociationProperty{
				Name: &scheduleName,
			},
			Runbook: jobschedule.RunbookAssociationProperty{
				Name: &runbookName,
			},
		},
	}

	// parameters to be passed into the runbook
	if v, ok := d.GetOk("parameters"); ok {
		jsParameters := make(map[string]string)
		for k, v := range v.(map[string]any) {
			jsParameters[k] = v.(string)
		}
		parameters.Properties.Parameters = &jsParameters
	}

	if v, ok := d.GetOk("run_on"); ok {
		parameters.Properties.RunOn = pointer.To(v.(string))
	}

	if _, err := client.Create(ctx, id, parameters); err != nil {
		return err
	}

	d.SetId(tfID.ID())
	d.Set("resource_manager_id", id.ID())

	return resourceAutomationJobScheduleRead(d, meta)
}
