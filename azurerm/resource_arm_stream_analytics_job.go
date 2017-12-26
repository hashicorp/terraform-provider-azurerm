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
			},
			"inputs": streamAnalyticsInputSchema(),
			// "output": streamAnalyticsOutputSchema(),
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

	if inputs, ok := d.GetOk("inputs"); ok {
		inputList := inputs.([]interface{})
		for _, inputSchema := range inputList {
			input, err := streamAnalyticsInputfromSchema(inputSchema)
			if err != nil {
				return err
			}
			result, err := client.inputsClient.CreateOrReplace(*input, rg, jobName, *input.Name, "", "")
			log.Printf("%#v \n", result)

		}

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
