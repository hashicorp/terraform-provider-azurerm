package azurerm

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/streamanalytics"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceArmStreamAnalyticsJob() *schema.Resource {

	return &schema.Resource{

		Create: resourceArmStreamAnalyticsJobCreate,
		Read:   resourceArmStreamAnalyticsJobRead,
		Update: resourceArmStreamAnalyticsJobUpdate,
		Delete: resourceArmStreamAnalyticsJobDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sku": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(streamanalytics.Standard),
				}, false),
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location": locationSchema(),
			"tags":     tagsSchema(),
			"events_out_of_order_max_delay_in_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				// non-negative interval
				ValidateFunc: validation.IntAtLeast(0),
			},
			"events_out_of_order_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(streamanalytics.Adjust),
					string(streamanalytics.Drop),
				}, false),
			},
			"job_state": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"transformation": streamAnalyticsTransformationSchema(),
			"job_input":      streamAnalyticsInputSchema(),
			"job_output":     streamAnalyticsOutputSchema(),
			"function":       streamAnalyticsFunctionSchema(),
		},
	}

}

func resourceArmStreamAnalyticsJobCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ArmClient)

	jobName := d.Get("name").(string)
	sku := d.Get("sku").(string)
	rg := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)

	jobProps := &streamanalytics.StreamingJobProperties{
		Sku: &streamanalytics.Sku{
			Name: streamanalytics.SkuName(sku),
		},
	}

	if sec, ok := d.GetOk("events_out_of_order_max_delay_in_seconds"); ok {
		seci := int32(sec.(int))
		jobProps.EventsOutOfOrderMaxDelayInSeconds = &seci
	}

	if evpolicy, ok := d.GetOk("events_out_of_order_policy"); ok {
		jobProps.EventsOutOfOrderPolicy = streamanalytics.EventsOutOfOrderPolicy(evpolicy.(string))
	}

	job := streamanalytics.StreamingJob{
		Name:                   &jobName,
		Location:               &location,
		StreamingJobProperties: jobProps,
	}

	if jobState, ok := d.GetOk("job_state"); ok {
		jobStateStr := jobState.(string)
		jobProps.JobState = &jobStateStr
	}

	// TODO: try to make this whole creation as atomic as possible
	jobChan, errChan := client.streamingJobClient.CreateOrReplace(job, rg, jobName, "", "", nil)
	err := <-errChan
	jobResp := <-jobChan

	// The reason that we set the id of the job here i.e. before creation of the related resource
	// is because if any of the child resource creation fail then the delete lifecycle method will
	// clean them up as deleting of job will remove all the child resources as well.
	// In retrospect if the setId is called after all the related resource are created then in case
	// of failure the delete method will not remove anything hence leaking some resources.
	d.SetId(*jobResp.ID)

	if functions, ok := d.GetOk("function"); ok {
		functionList := functions.([]interface{})
		for _, functionSchema := range functionList {
			function := streamAnalyticsFunctionFromSchema(functionSchema)

			result, err := client.functionClient.CreateOrReplace(function, rg, jobName, *function.Name, "", "")
			if err != nil {
				return err
			}
			log.Printf("[TRACE] Result from function creation is %#v \n", result)
		}
	}

	if inputs, ok := d.GetOk("job_input"); ok {
		inputList := inputs.([]interface{})
		for _, inputSchema := range inputList {
			input, err := streamAnalyticsInputfromSchema(inputSchema)
			if err != nil {
				return err
			}
			result, err := client.inputsClient.CreateOrReplace(*input, rg, jobName, *input.Name, "", "")
			if err != nil {
				return err
			}
			log.Printf("[TRACE] Result from input creation is %#v \n", result)

		}

	}

	if outputs, ok := d.GetOk("job_output"); ok {
		outputList := outputs.([]interface{})
		for _, outputSchema := range outputList {
			output, err := streamAnalyticsOutputFromSchema(outputSchema)
			if err != nil {
				return err
			}
			result, err := client.outputsClient.CreateOrReplace(*output, rg, jobName, *output.Name, "", "")
			if err != nil {
				return err
			}
			log.Printf("[TRACE] Result from output creation is %#v \n", result)

		}
	}

	if transformationI, ok := d.GetOk("transformation"); ok {
		transformationList := transformationI.([]interface{})
		transformationMap := transformationList[0].(map[string]interface{})
		transformation := streamAnalyticsTransformationFromSchema(transformationMap)
		result, err := client.trasformationsClient.CreateOrReplace(*transformation, rg, jobName, *transformation.Name, "", "")
		if err != nil {
			return err
		}
		log.Printf("Created transformation with fields %#v", result)

	}

	// This solves the chicken and egg situation going on
	if jobState, ok := d.GetOk("job_state"); ok {
		jobStateStr := jobState.(string)
		jobProps.JobState = &jobStateStr
		jobChan, errChan := client.streamingJobClient.CreateOrReplace(job, rg, jobName, "", "", nil)
		err := <-errChan
	}

	if err != nil {
		return err
	}

	return resourceArmStreamAnalyticsJobRead(d, meta)
}

func resourceArmStreamAnalyticsJobRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ArmClient)

	streamID := d.Id()
	resourceId, err := parseAzureResourceID(streamID)

	if err != nil {
		return err
	}
	job, err := client.streamingJobClient.Get(resourceId.ResourceGroup, resourceId.Path["streamingjobs"], "")

	if err != nil {
		return err
	}

	d.Set("job_state", *job.JobState)
	return nil
}

func resourceArmStreamAnalyticsJobUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceArmStreamAnalyticsJobCreate(d, meta)
}
func resourceArmStreamAnalyticsJobDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ArmClient)

	// TODO check if job exists or not in the first place

	jobName := d.Get("name").(string)
	rg := d.Get("resource_group_name").(string)

	_, errChan := client.streamingJobClient.Delete(rg, jobName, nil)
	err := <-errChan

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
