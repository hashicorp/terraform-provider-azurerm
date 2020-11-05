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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataFactoryTriggerTumblingWindow() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactoryTriggerTumblingWindowCreateUpdate,
		Read:   resourceArmDataFactoryTriggerTumblingWindowRead,
		Update: resourceArmDataFactoryTriggerTumblingWindowCreateUpdate,
		Delete: resourceArmDataFactoryTriggerTumblingWindowDelete,
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
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time,
			},

			// This time can only be  represented in UTC.
			// An issue has been filed in the SDK for the timezone attribute that doesn't seem to work
			// https://github.com/Azure/azure-sdk-for-go/issues/6244
			"end_time": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time,
			},

			"frequency": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"max_concurrency": {
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

			"delay": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.TriggerDelayTimespan(),
			},

			"pipeline_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"trigger_dependency": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.TriggerDelayTimespan(),
						},
						"offset": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.TriggerDelayTimespan(),
						},
						"trigger": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"retry": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
					},
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

func expandTriggerDependencies(d *schema.ResourceData) []datafactory.BasicDependencyReference {
	dependencies := d.Get("trigger_dependency").(*schema.Set).List()
	var expandedDependencies []datafactory.BasicDependencyReference

	for _, k := range dependencies {
		dep := k.(map[string]interface{})
		var trigger interface{}

		if target := dep["trigger"].(string); target != "" {
			trigger = &datafactory.TumblingWindowTriggerDependencyReference{
				ReferenceTrigger: &datafactory.TriggerReference{
					ReferenceName: utils.String(target),
					Type:          utils.String("TriggerReference"),
				},
			}
		} else {
			trigger = &datafactory.SelfDependencyTumblingWindowTriggerReference{}
		}

		var offset, size *string
		if v := dep["offset"].(string); v != "" {
			offset = utils.String(v)
		}
		if v := dep["size"].(string); v != "" {
			size = utils.String(v)
		}

		var basicDependency datafactory.BasicDependencyReference
		switch trigger := trigger.(type) {
		case *datafactory.TumblingWindowTriggerDependencyReference:
			trigger.Offset = offset
			trigger.Size = size
			basicDependency, _ = trigger.AsBasicDependencyReference()
		case *datafactory.SelfDependencyTumblingWindowTriggerReference:
			trigger.Offset = offset
			trigger.Size = size
			basicDependency, _ = trigger.AsBasicDependencyReference()
		}
		expandedDependencies = append(expandedDependencies, basicDependency)
	}

	return expandedDependencies
}

func flattenTriggerrDependencies(depRefs *[]datafactory.BasicDependencyReference) []interface{} {
	outputs := make([]interface{}, 0)
	for _, v := range *depRefs {
		var size, offset, trigger = "", "", ""
		var p_size, p_offset *string
		if t, ok := v.AsSelfDependencyTumblingWindowTriggerReference(); ok {
			p_size = t.Size
			p_offset = t.Offset
			trigger = ""
		} else if t, ok := v.AsTumblingWindowTriggerDependencyReference(); ok {
			p_size = t.Size
			p_offset = t.Offset
			trigger = *t.ReferenceTrigger.ReferenceName
		}

		if p_size != nil {
			size = *p_size
		}

		if p_offset != nil {
			offset = *p_offset
		}
		outputs = append(outputs, map[string]interface{}{
			"size":    size,
			"offset":  offset,
			"trigger": trigger,
		})
	}

	return outputs
}

func flattenRetryPolicy(r *datafactory.RetryPolicy) []interface{} {
	/*
		Sometimes the API returns RetryPolicy = nil sometimes it
		returns RetryPolicy = RetryPolicy{Count=nil, IntervalInSeconds=nil}
		in either case leaver retry as nil
	*/
	if r.Count == nil && r.IntervalInSeconds == nil {
		return nil
	}

	var count float64 = 0
	if r.Count != nil {
		count = r.Count.(float64)
	}

	var interval int32 = 0
	if r.IntervalInSeconds != nil {
		interval = *r.IntervalInSeconds
	}

	retry := map[string]interface{}{
		"count":    int(count),
		"interval": interval,
	}
	return []interface{}{retry}
}

func expandRetryPolicy(d *schema.ResourceData) *datafactory.RetryPolicy {
	policy := &datafactory.RetryPolicy{}
	if v, ok := d.GetOk("retry.0.count"); ok {
		policy.Count = utils.Int32(int32(v.(int)))
	}

	if v, ok := d.GetOk("retry.0.interval"); ok {
		policy.IntervalInSeconds = utils.Int32(int32(v.(int)))
	}

	return policy
}

func resourceArmDataFactoryTriggerTumblingWindowCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.TriggersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Data Factory Trigger Tumbling Window creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	triggerName := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, dataFactoryName, triggerName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory Trigger Tumbling Window %q (Resource Group %q / Data Factory %q): %s", triggerName, resourceGroupName, dataFactoryName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_trigger_tumbling_window", *existing.ID)
		}
	}

	props := &datafactory.TumblingWindowTriggerTypeProperties{
		Frequency: datafactory.TumblingWindowFrequency(d.Get("frequency").(string)),
		Interval:  utils.Int32(int32(d.Get("interval").(int))),
	}

	if v, ok := d.GetOk("delay"); ok {
		props.Delay = datafactory.TumblingWindowFrequency(v.(string))
	}

	if v, ok := d.GetOk("max_concurrency"); ok {
		props.MaxConcurrency = utils.Int32(int32(v.(int)))
	}

	if _, ok := d.GetOk("trigger_dependency"); ok {
		deps := expandTriggerDependencies(d)
		props.DependsOn = &deps
	}

	if _, ok := d.GetOk("retry"); ok {
		props.RetryPolicy = expandRetryPolicy(d)
	}

	if v, ok := d.GetOk("start_time"); ok {
		t, _ := time.Parse(time.RFC3339, v.(string)) // should be validated by the schema
		props.StartTime = &date.Time{Time: t}
	} else {
		props.StartTime = &date.Time{Time: time.Now()}
	}

	if v, ok := d.GetOk("end_time"); ok {
		t, _ := time.Parse(time.RFC3339, v.(string)) // should be validated by the schema
		props.EndTime = &date.Time{Time: t}
	}

	triggerProps := &datafactory.TumblingWindowTrigger{
		TumblingWindowTriggerTypeProperties: props,
		Pipeline: &datafactory.TriggerPipelineReference{
			PipelineReference: &datafactory.PipelineReference{
				ReferenceName: utils.String(d.Get("pipeline_name").(string)),
				Type:          utils.String("PipelineReference"),
			},
			Parameters: d.Get("pipeline_parameters").(map[string]interface{}),
		},
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		triggerProps.Annotations = &annotations
	}

	trigger := datafactory.TriggerResource{
		Properties: triggerProps,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroupName, dataFactoryName, triggerName, trigger, ""); err != nil {
		return fmt.Errorf("Error creating Data Factory Trigger Tumbling Window %q (Resource Group %q / Data Factory %q): %+v", triggerName, resourceGroupName, dataFactoryName, err)
	}

	read, err := client.Get(ctx, resourceGroupName, dataFactoryName, triggerName, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Trigger Tumbling Window  %q (Resource Group %q / Data Factory %q): %+v", triggerName, resourceGroupName, dataFactoryName, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Trigger Tumbling Window  %q (Resource Group %q / Data Factory %q) ID", triggerName, resourceGroupName, dataFactoryName)
	}

	d.SetId(*read.ID)

	return resourceArmDataFactoryTriggerTumblingWindowRead(d, meta)
}

func resourceArmDataFactoryTriggerTumblingWindowRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Data Factory Trigger Tumbling Window %q was not found in Resource Group %q - removing from state!", triggerName, id.ResourceGroup)
			return nil
		}
		return fmt.Errorf("Error reading the state of Data Factory Trigger Tumbling Window %q: %+v", triggerName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	tumblingTrigger, ok := resp.Properties.AsTumblingWindowTrigger()
	if !ok {
		return fmt.Errorf("Error classifiying Data Factory Trigger Tumbling Window %q (Data Factory %q / Resource Group %q): Expected: %q Received: %q", triggerName, dataFactoryName, id.ResourceGroup, datafactory.TypeTumblingWindowTrigger, *resp.Type)
	}

	if tumblingTrigger != nil {
		if tumblingTriggerProps := tumblingTrigger.TumblingWindowTriggerTypeProperties; tumblingTriggerProps != nil {
			if v := tumblingTriggerProps.StartTime; v != nil {
				d.Set("start_time", v.Format(time.RFC3339))
			}
			if v := tumblingTriggerProps.EndTime; v != nil {
				d.Set("end_time", v.Format(time.RFC3339))
			}
			d.Set("frequency", tumblingTriggerProps.Frequency)
			d.Set("interval", tumblingTriggerProps.Interval)
			d.Set("max_concurrency", tumblingTriggerProps.MaxConcurrency)
			d.Set("delay", tumblingTriggerProps.Delay)

			if v := tumblingTriggerProps.RetryPolicy; v != nil {
				d.Set("retry", flattenRetryPolicy(v))
			}

			if v := tumblingTriggerProps.DependsOn; v != nil {
				d.Set("trigger_dependency", flattenTriggerrDependencies(v))
			}
		}

		if pipeline := tumblingTrigger.Pipeline; pipeline != nil {
			if reference := pipeline.PipelineReference; reference != nil {
				d.Set("pipeline_name", reference.ReferenceName)
			}
			d.Set("pipeline_parameters", pipeline.Parameters)
		}

		annotations := flattenDataFactoryAnnotations(tumblingTrigger.Annotations)
		if err := d.Set("annotations", annotations); err != nil {
			return fmt.Errorf("Error setting `annotations`: %+v", err)
		}
	}

	return nil
}

func resourceArmDataFactoryTriggerTumblingWindowDelete(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error deleting Data Factory Trigger Tumbling Window %q (Resource Group %q / Data Factory %q): %+v", triggerName, id.ResourceGroup, dataFactoryName, err)
	}

	return nil
}
