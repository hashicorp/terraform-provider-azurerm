package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutomationJobSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationJobScheduleCreate,
		Read:   resourceArmAutomationJobScheduleRead,
		Delete: resourceArmAutomationJobScheduleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"resource_group_name": azure.SchemaResourceGroupName(),

			"automation_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateAutomationAccountName(),
			},

			"runbook_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateAutomationRunbookName(),
			},

			"schedule_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateAutomationScheduleName(),
			},

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ValidateFunc: func(v interface{}, _ string) (warnings []string, errors []error) {
					m := v.(map[string]interface{})

					for k := range m {
						if k != strings.ToLower(k) {
							errors = append(errors, fmt.Errorf("Due to a bug in the implementation of Runbooks in Azure, the parameter names need to be specified in lowercase only. See: \"https://github.com/Azure/azure-sdk-for-go/issues/4780\" for more information."))
						}
					}

					return warnings, errors
				},
			},

			"run_on": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"job_schedule_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArmAutomationJobScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.JobScheduleClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Automation Job Schedule creation.")

	jobScheduleUUID := uuid.NewV4()
	resGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("automation_account_name").(string)

	runbookName := d.Get("runbook_name").(string)
	scheduleName := d.Get("schedule_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, accountName, jobScheduleUUID)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Automation Job Schedule %q (Account %q / Resource Group %q): %s", jobScheduleUUID, accountName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_automation_job_schedule", *existing.ID)
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

	if _, err := client.Create(ctx, resGroup, accountName, jobScheduleUUID, parameters); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, accountName, jobScheduleUUID)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Job Schedule '%s' (Account %q / Resource Group %s) ID", jobScheduleUUID, accountName, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAutomationJobScheduleRead(d, meta)
}

func resourceArmAutomationJobScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.JobScheduleClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	jobScheduleID := id.Path["jobSchedules"]
	jobScheduleUUID := uuid.FromStringOrNil(jobScheduleID)
	resGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]

	resp, err := client.Get(ctx, resGroup, accountName, jobScheduleUUID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Automation Job Schedule '%s': %+v", jobScheduleUUID, err)
	}

	d.Set("job_schedule_id", resp.JobScheduleID)
	d.Set("resource_group_name", resGroup)
	d.Set("automation_account_name", accountName)
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

func resourceArmAutomationJobScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.JobScheduleClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	jobScheduleID := id.Path["jobSchedules"]
	jobScheduleUUID := uuid.FromStringOrNil(jobScheduleID)
	resGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]

	resp, err := client.Delete(ctx, resGroup, accountName, jobScheduleUUID)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing AzureRM delete request for Automation Job Schedule '%s': %+v", jobScheduleUUID, err)
		}
	}

	return nil
}
