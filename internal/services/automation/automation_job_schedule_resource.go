package automation

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAutomationJobSchedule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationJobScheduleCreate,
		Read:   resourceAutomationJobScheduleRead,
		Delete: resourceAutomationJobScheduleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.JobScheduleID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

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
	client := meta.(*clients.Client).Automation.JobScheduleClient
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

	id := parse.NewJobScheduleID(client.SubscriptionID, d.Get("resource_group_name").(string), d.Get("automation_account_name").(string), jobScheduleUUID.String())

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, jobScheduleUUID)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_automation_job_schedule", id.ID())
		}
	}

	// fix issue: https://github.com/hashicorp/terraform-provider-azurerm/issues/7130
	// When the runbook has some updates, it'll update all related job schedule id, so the elder job schedule will not exist
	// We need to delete the job schedule id if exists to recreate the job schedule
	for jsIterator, err := client.ListByAutomationAccountComplete(ctx, id.ResourceGroup, id.AutomationAccountName, ""); jsIterator.NotDone(); err = jsIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("loading Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
		}
		if props := jsIterator.Value().JobScheduleProperties; props != nil {
			if props.Schedule.Name != nil && *props.Schedule.Name == scheduleName && props.Runbook.Name != nil && *props.Runbook.Name == runbookName {
				if jsIterator.Value().JobScheduleID == nil || *jsIterator.Value().JobScheduleID == "" {
					return fmt.Errorf("job schedule Id is nil or empty listed by Automation Account %q Job Schedule List: %+v", id.AutomationAccountName, err)
				}
				jsId, err := uuid.FromString(*jsIterator.Value().JobScheduleID)
				if err != nil {
					return fmt.Errorf("parsing job schedule Id listed by Automation Account %q Job Schedule List:%v", id.AutomationAccountName, err)
				}
				if _, err := client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, jsId); err != nil {
					return fmt.Errorf("deleting job schedule Id listed by Automation Account %q Job Schedule List:%v", id.AutomationAccountName, err)
				}
			}
		}
	}

	parameters := automation.JobScheduleCreateParameters{
		JobScheduleCreateProperties: &automation.JobScheduleCreateProperties{
			Schedule: &automation.ScheduleAssociationProperty{
				Name: &scheduleName,
			},
			Runbook: &automation.RunbookAssociationProperty{
				Name: &runbookName,
			},
		},
	}
	properties := parameters.JobScheduleCreateProperties

	// parameters to be passed into the runbook
	if v, ok := d.GetOk("parameters"); ok {
		jsParameters := make(map[string]*string)
		for k, v := range v.(map[string]interface{}) {
			value := v.(string)
			jsParameters[k] = &value
		}
		properties.Parameters = jsParameters
	}

	if v, ok := d.GetOk("run_on"); ok {
		value := v.(string)
		properties.RunOn = &value
	}

	if _, err := client.Create(ctx, id.ResourceGroup, id.AutomationAccountName, jobScheduleUUID, parameters); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceAutomationJobScheduleRead(d, meta)
}

func resourceAutomationJobScheduleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.JobScheduleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.JobScheduleID(d.Id())
	if err != nil {
		return err
	}

	jobScheduleUUID := uuid.FromStringOrNil(id.Name)

	resp, err := client.Get(ctx, id.ResourceGroup, id.AutomationAccountName, jobScheduleUUID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on AzureRM Automation Job Schedule '%s': %+v", id.Name, err)
	}

	d.Set("job_schedule_id", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("automation_account_name", id.AutomationAccountName)
	d.Set("runbook_name", resp.JobScheduleProperties.Runbook.Name)
	d.Set("schedule_name", resp.JobScheduleProperties.Schedule.Name)

	if v := resp.JobScheduleProperties.RunOn; v != nil {
		d.Set("run_on", v)
	}

	if v := resp.JobScheduleProperties.Parameters; v != nil {
		jsParameters := make(map[string]interface{})
		for key, value := range v {
			jsParameters[strings.ToLower(key)] = value
		}
		d.Set("parameters", jsParameters)
	}

	return nil
}

func resourceAutomationJobScheduleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.JobScheduleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.JobScheduleID(d.Id())
	if err != nil {
		return err
	}

	jobScheduleUUID := uuid.FromStringOrNil(id.Name)

	resp, err := client.Delete(ctx, id.ResourceGroup, id.AutomationAccountName, jobScheduleUUID)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("issuing AzureRM delete request for Automation Job Schedule '%s': %+v", id.Name, err)
		}
	}

	return nil
}
