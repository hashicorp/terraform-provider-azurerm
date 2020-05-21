package datafactory

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
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
				ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/5788
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"data_factory_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryName(),
			},

			// This time can only be  represented in UTC.
			// An issue has been filed in the SDK for the timezone attribute that doesn't seem to work
			// https://github.com/Azure/azure-sdk-for-go/issues/6244
			"start_time": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time, // times in the past just start immediately
			},

			// This time can only be  represented in UTC.
			// An issue has been filed in the SDK for the timezone attribute that doesn't seem to work
			// https://github.com/Azure/azure-sdk-for-go/issues/6244
			"end_time": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time, // times in the past just start immediately
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
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"pipeline_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
			},

			"pipeline_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"annotations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceArmDataFactoryTriggerScheduleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Data Factory Trigger Schedule creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	triggerName := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)

	if d.IsNewResource() {
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
		},
	}

	if v, ok := d.GetOk("start_time"); ok {
		t, _ := time.Parse(time.RFC3339, v.(string)) // should be validated by the schema
		props.Recurrence.StartTime = &date.Time{Time: t}
	} else {
		props.Recurrence.StartTime = &date.Time{Time: time.Now()}
	}

	if v, ok := d.GetOk("end_time"); ok {
		t, _ := time.Parse(time.RFC3339, v.(string)) // should be validated by the schema
		props.Recurrence.EndTime = &date.Time{Time: t}
	}

	reference := &datafactory.PipelineReference{
		ReferenceName: utils.String(d.Get("pipeline_name").(string)),
		Type:          utils.String("PipelineReference"),
	}

	scheduleProps := &datafactory.ScheduleTrigger{
		ScheduleTriggerTypeProperties: props,
		Pipelines: &[]datafactory.TriggerPipelineReference{
			{
				PipelineReference: reference,
				Parameters:        d.Get("pipeline_parameters").(map[string]interface{}),
			},
		},
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
	client := meta.(*clients.Client).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
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
				d.Set("start_time", v.Format(time.RFC3339))
			}
			if v := recurrence.EndTime; v != nil {
				d.Set("end_time", v.Format(time.RFC3339))
			}
			d.Set("frequency", recurrence.Frequency)
			d.Set("interval", recurrence.Interval)
		}

		if pipelines := scheduleTriggerProps.Pipelines; pipelines != nil {
			if len(*pipelines) > 0 {
				pipeline := *pipelines
				if reference := pipeline[0].PipelineReference; reference != nil {
					d.Set("pipeline_name", reference.ReferenceName)
				}
				d.Set("pipeline_parameters", pipeline[0].Parameters)
			}
		}

		annotations := flattenDataFactoryAnnotations(scheduleTriggerProps.Annotations)
		if err := d.Set("annotations", annotations); err != nil {
			return fmt.Errorf("Error setting `annotations`: %+v", err)
		}
	}

	return nil
}

func resourceArmDataFactoryTriggerScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
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
