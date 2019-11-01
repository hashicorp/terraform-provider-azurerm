package azurerm

import (
	"fmt"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataFactoryTriggerSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactoryTriggerScheduleCreateUpdate,
		Read:   resourceArmDataFactoryTriggerScheduleRead,
		Update: resourceArmDataFactoryTriggerScheduleCreateUpdate,
		Delete: resourceArmDataFactoryTriggerScheduleDelete,
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMDataFactoryPipelineName,
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/5788
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"data_factory_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
					`Invalid data_factory_name, see https://docs.microsoft.com/en-us/azure/data-factory/naming-rules`,
				),
			},

			"start_time": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validate.RFC3339Time, //times in the past just start immediately
			},

			"end_time": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validate.RFC3339Time, //times in the past just start immediately
			},

			"frequency": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(datafactory.Minute),
				ValidateFunc: validation.StringInSlice([]string{
					string(datafactory.Minute),
					string(datafactory.Hour),
					string(datafactory.Day),
					string(datafactory.Week),
					string(datafactory.Month),
				}, false),
			},

			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UTC",
				ValidateFunc: validate.NoEmptyStrings,
			},

			"annotations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmDataFactoryTriggerScheduleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Data Factory Trigger Schedule creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	triggerName := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, dataFactoryName, triggerName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory Trigger Schedule %q (Resource Group %q / Data Factory %q): %s", triggerName, resourceGroupName, dataFactoryName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_trigger_schedule", *existing.ID)
		}
	}

	props := &datafactory.ScheduleTriggerTypeProperties{
		Recurrence: &datafactory.ScheduleTriggerRecurrence{
			Frequency: datafactory.RecurrenceFrequency(d.Get("frequency").(string)),
			Interval:  utils.Int32(int32(d.Get("interval").(int))),
			TimeZone:  utils.String(d.Get("timezone").(string)),
		},
	}

	if v, ok := d.GetOk("start_time"); ok {
		t, _ := time.Parse(time.RFC3339, v.(string)) //should be validated by the schema
		props.Recurrence.StartTime = &date.Time{Time: t}
	} else {
		props.Recurrence.StartTime = &date.Time{Time: time.Now()}
	}

	if v, ok := d.GetOk("end_time"); ok {
		t, _ := time.Parse(time.RFC3339, v.(string)) //should be validated by the schema
		props.Recurrence.EndTime = &date.Time{Time: t}
	}

	scheduleProps := &datafactory.ScheduleTrigger{
		ScheduleTriggerTypeProperties: props,
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		scheduleProps.Annotations = &annotations
	}

	trigger := datafactory.TriggerResource{
		Properties: scheduleProps,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroupName, dataFactoryName, triggerName, trigger, ""); err != nil {
		return fmt.Errorf("Error creating Data Factory Trigger Schedule %q (Resource Group %q / Data Factory %q): %+v", triggerName, resourceGroupName, dataFactoryName, err)
	}

	read, err := client.Get(ctx, resourceGroupName, dataFactoryName, triggerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Trigger Schedule %q (Resource Group %q / Data Factory %q): %+v", triggerName, resourceGroupName, dataFactoryName, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Trigger Schedule %q (Resource Group %q / Data Factory %q) ID", triggerName, resourceGroupName, dataFactoryName)
	}

	d.SetId(*read.ID)

	return resourceArmDataFactoryTriggerScheduleRead(d, meta)
}

func resourceArmDataFactoryTriggerScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	dataFactoryName := id.Path["factories"]
	triggerName := id.Path["triggers"]

	resp, err := client.Get(ctx, id.ResourceGroup, dataFactoryName, triggerName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			log.Printf("[DEBUG] Data Factory Trigger Schedule %q was not found in Resource Group %q - removing from state!", triggerName, id.ResourceGroup)
			return nil
		}
		return fmt.Errorf("Error reading the state of Data Factory Trigger Schedule %q: %+v", triggerName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	scheduleTriggerProps, ok := resp.Properties.AsScheduleTrigger()
	if !ok {
		return fmt.Errorf("Error classifiying Data Factory Trigger Schedule %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", triggerName, dataFactoryName, id.ResourceGroup, datafactory.TypeScheduleTrigger, *resp.Type)
	}

	if scheduleTriggerProps != nil {
		if recurrence := scheduleTriggerProps.Recurrence; recurrence != nil {
			if v := recurrence.StartTime; v != nil {
				d.Set("start_time", (*v).Format(time.RFC3339))
			}
			if v := recurrence.EndTime; v != nil {
				d.Set("end_time", (*v).Format(time.RFC3339))
			}
			d.Set("frequency", recurrence.Frequency)
			d.Set("interval", recurrence.Interval)
			d.Set("timezone", recurrence.TimeZone)
		}

		annotations := flattenDataFactoryAnnotations(scheduleTriggerProps.Annotations)
		if err := d.Set("annotations", annotations); err != nil {
			return fmt.Errorf("Error setting `annotations`: %+v", err)
		}
	}

	return nil
}

func resourceArmDataFactoryTriggerScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	dataFactoryName := id.Path["factories"]
	triggerName := id.Path["triggers"]

	if _, err = client.Delete(ctx, id.ResourceGroup, dataFactoryName, triggerName); err != nil {
		return fmt.Errorf("Error deleting Data Factory Trigger Schedule %q (Resource Group %q / Data Factory %q): %+v", triggerName, id.ResourceGroup, dataFactoryName, err)
	}

	return nil
}
