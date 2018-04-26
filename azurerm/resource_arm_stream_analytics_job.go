package azurerm

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
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
				ForceNew: true,
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

func setFunctions(d *schema.ResourceData, client *ArmClient, rg, jobName string) error {
	ctx := client.StopContext
	if functions, ok := d.GetOk("function"); ok {
		functionList := functions.([]interface{})
		for _, functionSchema := range functionList {
			function := streamAnalyticsFunctionFromSchema(functionSchema)

			result, err := client.streamAnalyticsFunctionsClient.CreateOrReplace(ctx, function, rg, jobName, *function.Name, "", "")
			if err != nil {
				return err
			}
			log.Printf("[TRACE] Result from function creation is %#v \n", result)
		}
	}

	return nil
}

func setInputs(d *schema.ResourceData, client *ArmClient, rg, jobName string) error {
	ctx := client.StopContext
	if inputs, ok := d.GetOk("job_input"); ok {
		inputList := inputs.([]interface{})
		for _, inputSchema := range inputList {
			input, err := streamAnalyticsInputfromSchema(inputSchema)
			if err != nil {
				return err
			}
			result, err := client.streamAnalyticsInputsClient.CreateOrReplace(ctx, *input, rg, jobName, *input.Name, "", "")
			if err != nil {
				return err
			}
			log.Printf("[TRACE] Result from input creation is %#v \n", result)

		}

	}

	return nil

}
func setOutputs(d *schema.ResourceData, client *ArmClient, rg, jobName string) error {
	ctx := client.StopContext
	if outputs, ok := d.GetOk("job_output"); ok {
		outputList := outputs.([]interface{})
		for _, outputSchema := range outputList {
			output, err := streamAnalyticsOutputFromSchema(outputSchema)
			if err != nil {
				return err
			}
			result, err := client.streamAnalyticsOutputsClient.CreateOrReplace(ctx, *output, rg, jobName, *output.Name, "", "")
			if err != nil {
				return err
			}
			log.Printf("[TRACE] Result from output creation is %#v \n", result)
		}
	}

	return nil
}

func setTransformation(d *schema.ResourceData, client *ArmClient, rg, jobName string) error {
	ctx := client.StopContext
	if transformationI, ok := d.GetOk("transformation"); ok {
		transformationList := transformationI.([]interface{})
		transformationMap := transformationList[0].(map[string]interface{})
		transformation := streamAnalyticsTransformationFromSchema(transformationMap)
		result, err := client.streamAnalyticsTransformationsClient.CreateOrReplace(ctx, *transformation, rg, jobName, *transformation.Name, "", "")
		if err != nil {
			return err
		}
		log.Printf("Created transformation with fields %#v", result)

	}
	return nil
}

func setJobState(d *schema.ResourceData, client *ArmClient, rg, jobName string) error {
	if err := handleRunningJobState(d, client, rg, jobName); err != nil {
		return err
	}

	return handleStoppedJobState(d, client, rg, jobName)
}

func handleRunningJobState(d *schema.ResourceData, client *ArmClient, rg, jobName string) error {
	ctx := client.StopContext
	streamClient := client.streamAnalyticsJobsClient
	if jobState, ok := d.GetOk("job_state"); ok {
		jobStateStr := jobState.(string)

		job, err := client.streamAnalyticsJobsClient.Get(ctx, rg, jobName, "")
		if err != nil {
			return err
		}

		jobParams := &streamanalytics.StartStreamingJobParameters{
			OutputStartMode: job.OutputStartMode,
			OutputStartTime: job.OutputStartTime,
		}

		switch jobStateStr {
		case "Running":
			future, err := client.streamAnalyticsJobsClient.Start(ctx, rg, jobName, jobParams)
			if err != nil {
				if response.WasNotFound(future.Response()) {
					return nil
				}
				return err
			}
			err = future.WaitForCompletion(ctx, streamClient.Client)
			if err != nil {
				if response.WasNotFound(future.Response()) {
					return nil
				}
				return fmt.Errorf("Error starting job %q", jobName)
			}
		}
	}
	return nil
}

func handleStoppedJobState(d *schema.ResourceData, client *ArmClient, rg, jobName string) error {
	ctx := client.StopContext
	streamClient := client.streamAnalyticsJobsClient
	if jobState, ok := d.GetOk("job_state"); ok {
		jobStateStr := jobState.(string)

		cancelChan := make(chan struct{})
		defer close(cancelChan)

		switch jobStateStr {
		case "Stopped":
			future, err := client.streamAnalyticsJobsClient.Stop(ctx, rg, jobName)
			if err != nil {
				if response.WasNotFound(future.Response()) {
					return nil
				}
				return err
			}
			err = future.WaitForCompletion(ctx, streamClient.Client)
			if err != nil {
				if response.WasNotFound(future.Response()) {
					return nil
				}
				return fmt.Errorf("Error stopping job %q", jobName)
			}
		}
	}
	return nil
}

func resourceArmStreamAnalyticsJobCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ArmClient)
	streamClient := client.streamAnalyticsJobsClient
	ctx := client.StopContext

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

	if tagsInf, ok := d.GetOk("tags"); ok {
		job.Tags = *expandTags(tagsInf.(map[string]interface{}))
	}

	// TODO: try to make this whole creation as atomic as possible
	future, err := client.streamAnalyticsJobsClient.CreateOrReplace(ctx, job, rg, jobName, "", "")
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, streamClient.Client)
	if err != nil {
		return fmt.Errorf("Error creating or updating Stream Analytics %q (Resource Group %q): %+v", jobName, rg, err)
	}

	jobResp, err := client.streamAnalyticsJobsClient.Get(ctx, rg, jobName, "")
	if err != nil {
		return err
	}

	// The reason that we set the id of the job here i.e. before creation of the related resource
	// is because if any of the child resource creation fail then the delete lifecycle method will
	// clean them up as deleting of job will remove all the child resources as well.
	// In retrospect if the setId is called after all the related resource are created then in case
	// of failure the delete method will not remove anything hence leaking some resources.
	d.SetId(*jobResp.ID)

	err = setFunctions(d, client, rg, jobName)
	if err != nil {
		return err
	}

	err = setInputs(d, client, rg, jobName)
	if err != nil {
		return err
	}

	err = setOutputs(d, client, rg, jobName)
	if err != nil {
		return err
	}

	err = setTransformation(d, client, rg, jobName)
	if err != nil {
		return err
	}

	// This solves the chicken and egg situation going on
	err = setJobState(d, client, rg, jobName)
	if err != nil {
		return err
	}

	return resourceArmStreamAnalyticsJobRead(d, meta)
}

func resourceArmStreamAnalyticsJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)

	jobName := d.Get("name").(string)
	rg := d.Get("resource_group_name").(string)

	// if state is to stop the job then it should happen in the beginning
	if d.HasChange("job_state") {

		err := handleStoppedJobState(d, client, rg, jobName)
		if err != nil {
			return err
		}
	}

	if d.HasChange("job_input") {
		err := setInputs(d, client, rg, jobName)
		if err != nil {
			return err
		}
	}

	if d.HasChange("job_output") {
		err := setOutputs(d, client, rg, jobName)
		if err != nil {
			return err
		}
	}

	if d.HasChange("function") {
		err := setFunctions(d, client, rg, jobName)
		if err != nil {
			return err
		}
	}

	err := manageProps(d, client, rg, jobName)
	if err != nil {
		return err
	}

	if d.HasChange("transformation") {
		err := setTransformation(d, client, rg, jobName)
		if err != nil {
			return err
		}
	}

	// if state is to start the job then that should start at the end else no changes can be applied
	if d.HasChange("job_state") {
		err := handleRunningJobState(d, client, rg, jobName)
		if err != nil {
			return err
		}
	}

	return nil

}

func manageProps(d *schema.ResourceData, client *ArmClient, rg, jobName string) error {
	ctx := client.StopContext
	jobProps := &streamanalytics.StreamingJobProperties{}

	// Get our current Stream analytics instance and update a property if one is not set
	streamID := d.Id()
	resourceID, err := parseAzureResourceID(streamID)
	if err != nil {
		return err
	}

	currentJob, err := client.streamAnalyticsJobsClient.Get(ctx, resourceID.ResourceGroup, resourceID.Path["streamingjobs"], "")
	if err != nil {
		return err
	}

	if sec, ok := d.GetOk("events_out_of_order_max_delay_in_seconds"); ok {
		seci := int32(sec.(int))
		jobProps.EventsOutOfOrderMaxDelayInSeconds = &seci
	} else {
		jobProps.EventsOutOfOrderMaxDelayInSeconds = currentJob.StreamingJobProperties.EventsOutOfOrderMaxDelayInSeconds
	}

	if evpolicy, ok := d.GetOk("events_out_of_order_policy"); ok {
		jobProps.EventsOutOfOrderPolicy = streamanalytics.EventsOutOfOrderPolicy(evpolicy.(string))
	}

	job := streamanalytics.StreamingJob{
		StreamingJobProperties: jobProps,
	}

	if tagsInf, ok := d.GetOk("tags"); ok {
		job.Tags = *expandTags(tagsInf.(map[string]interface{}))
	}

	// This causes an error each time it runs with no changes. We will display the error and continue
	_, err = client.streamAnalyticsJobsClient.Update(ctx, job, rg, jobName, "")
	return err
}

func resourceArmStreamAnalyticsJobRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ArmClient)
	ctx := client.StopContext

	streamID := d.Id()
	resourceID, err := parseAzureResourceID(streamID)

	if err != nil {
		return err
	}
	job, err := client.streamAnalyticsJobsClient.Get(ctx, resourceID.ResourceGroup, resourceID.Path["streamingjobs"], "")

	if err != nil {
		return err
	}

	flattenAndSetTags(d, &job.Tags)

	d.Set("job_state", *job.JobState)
	return nil
}

func resourceArmStreamAnalyticsJobDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*ArmClient)
	ctx := client.StopContext

	// TODO check if job exists or not in the first place

	jobName := d.Get("name").(string)
	rg := d.Get("resource_group_name").(string)

	_, err := client.streamAnalyticsJobsClient.Delete(ctx, rg, jobName)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
