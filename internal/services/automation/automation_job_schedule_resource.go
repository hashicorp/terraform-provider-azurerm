// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/jobschedule"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceAutomationJobSchedule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationJobScheduleCreate,
		Read:   resourceAutomationJobScheduleRead,
		Delete: resourceAutomationJobScheduleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseCompositeResourceID(id, &schedule.ScheduleId{}, &runbook.RunbookId{})
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AutomationJobScheduleV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"automation_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"runbook_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RunbookName(),
			},

			"schedule_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ScheduleName(),
			},

			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
				ValidateFunc: validate.ParameterNames,
			},

			"run_on": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"job_schedule_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// NOTE: O+C this can remain as this can change if the runbook is updated but cannot be updated by the user
				Computed:     true,
				ForceNew:     features.FourPointOhBeta(),
				ValidateFunc: validation.IsUUID,
			},

			"resource_manager_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAutomationJobScheduleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.JobSchedule
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Job Schedule creation.")

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

	if d.IsNewResource() {
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
		for k, v := range v.(map[string]interface{}) {
			value := v.(string)
			jsParameters[k] = value
		}
		parameters.Properties.Parameters = &jsParameters
	}

	if v, ok := d.GetOk("run_on"); ok {
		value := v.(string)
		parameters.Properties.RunOn = &value
	}

	if _, err := client.Create(ctx, id, parameters); err != nil {
		return err
	}

	d.SetId(tfID.ID())
	d.Set("resource_manager_id", id.ID())

	return resourceAutomationJobScheduleRead(d, meta)
}

func resourceAutomationJobScheduleRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
				jsParameters := make(map[string]interface{})
				for key, value := range v {
					jsParameters[strings.ToLower(key)] = value
				}
				d.Set("parameters", jsParameters)
			}
		}
	}

	return nil
}

func resourceAutomationJobScheduleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
