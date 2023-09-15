// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2022-08-08/jobschedule"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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
			_, err := jobschedule.ParseJobScheduleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

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
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
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

	runbookName := d.Get("runbook_name").(string)
	scheduleName := d.Get("schedule_name").(string)

	jobScheduleUUID, err := uuid.NewV4()
	if err != nil {
		return err
	}
	if jobScheduleID, ok := d.GetOk("job_schedule_id"); ok {
		jobScheduleUUID = uuid.FromStringOrNil(jobScheduleID.(string))
	}

	id := jobschedule.NewJobScheduleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), jobScheduleUUID.String())

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_automation_job_schedule", id.ID())
		}
	}

	automationAccountId := jobschedule.NewAutomationAccountID(subscriptionId, id.ResourceGroupName, id.AutomationAccountName)

	// fix issue: https://github.com/hashicorp/terraform-provider-azurerm/issues/7130
	// When the runbook has some updates, it'll update all related job schedule id, so the elder job schedule will not exist
	// We need to delete the job schedule id if exists to recreate the job schedule
	jsIterator, err := client.ListByAutomationAccountComplete(ctx, automationAccountId, jobschedule.ListByAutomationAccountOperationOptions{})
	if err != nil {
		return fmt.Errorf("loading Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
	}

	for _, item := range jsIterator.Items {
		if itemProps := item.Properties; itemProps != nil {
			if itemProps.Schedule != nil && itemProps.Schedule.Name != nil && *itemProps.Schedule.Name == scheduleName && itemProps.Runbook != nil && itemProps.Runbook.Name != nil && *itemProps.Runbook.Name == runbookName {
				if itemProps.JobScheduleId == nil || *itemProps.JobScheduleId == "" {
					return fmt.Errorf("job schedule Id is nil or empty listed by Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
				}

				jsId := jobschedule.NewJobScheduleID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName, *itemProps.JobScheduleId)
				if _, err := client.Delete(ctx, jsId); err != nil {
					return fmt.Errorf("deleting job schedule Id listed by Automation Account %q Job Schedule List:%v", id.AutomationAccountName, err)
				}
			}
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

	d.SetId(id.ID())

	return resourceAutomationJobScheduleRead(d, meta)
}

func resourceAutomationJobScheduleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.JobSchedule
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := jobschedule.ParseJobScheduleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("job_schedule_id", id.JobScheduleId)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("automation_account_name", id.AutomationAccountName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
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
	}

	return nil
}

func resourceAutomationJobScheduleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.JobSchedule
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := jobschedule.ParseJobScheduleID(d.Id())
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
