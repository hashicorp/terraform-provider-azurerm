package azurerm

import (
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
			"streaming_unit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"job_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"inputs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(streamanalytics.TypeReference),
								string(streamanalytics.TypeStream),
							}, false),
						},
						"serialization": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(streamanalytics.TypeAvro),
											string(streamanalytics.TypeCsv),
											string(streamanalytics.TypeJSON),
										}, false),
									},
									"field_delimiter": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"encoding": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Default:  string(streamanalytics.UTF8),
										ValidateFunc: validation.StringInSlice([]string{
											string(streamanalytics.UTF8),
										}, false),
									},
								},
							},
						},
						"datasource": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"blob": &schema.Schema{
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										MinItems: 1,
										ConflictsWith: []string{
											"event_hub",
											"iot_hub",
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"storage_account_name": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"storage_account_key": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"container": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"path_pattern": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"date_format": &schema.Schema{
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
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

	jobChan, errChan := client.streamingJobClient.CreateOrReplace(job, rg, jobName, "", "", nil)
	err := <-errChan

	if err != nil {
		return err
	}
	jobResp := <-jobChan

	d.SetId(*jobResp.ID)

	return nil
}

func resourceArmStreamAnalyticsJobRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmStreamAnalyticsJobUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
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
