package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
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

		Schema: map[string]*schema.Schema{

			"resource_group_name": resourceGroupNameSchema(),

			"automation_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"runbook_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"schedule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"run_on": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceArmAutomationJobScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationJobScheduleClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Automation Job Schedule creation.")

	nameUUID := uuid.NewV4()
	resGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("automation_account_name").(string)

	runbookName := d.Get("runbook_name").(string)
	scheduleName := d.Get("schedule_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, accountName, nameUUID)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Automation Job Schedule %q (Account %q / Resource Group %q): %s", nameUUID, accountName, resGroup, err)
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

	if _, err := client.Create(ctx, resGroup, accountName, nameUUID, parameters); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, accountName, nameUUID)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Job Schedule '%s' (Account %q / Resource Group %s) ID", nameUUID, accountName, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAutomationJobScheduleRead(d, meta)
}

func resourceArmAutomationJobScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationJobScheduleClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["jobSchedules"]
	nameUUID := uuid.FromStringOrNil(name)
	resGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]

	resp, err := client.Get(ctx, resGroup, accountName, nameUUID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Automation Job Schedule '%s': %+v", nameUUID, err)
	}

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
			jsParameters[key] = value
		}
		d.Set("parameters", jsParameters)
	}

	return nil
}

func resourceArmAutomationJobScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationJobScheduleClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["jobSchedules"]
	nameUUID := uuid.FromStringOrNil(name)
	resGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]

	resp, err := client.Delete(ctx, resGroup, accountName, nameUUID)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing AzureRM delete request for Automation Job Schedule '%s': %+v", nameUUID, err)
		}
	}

	return nil
}
