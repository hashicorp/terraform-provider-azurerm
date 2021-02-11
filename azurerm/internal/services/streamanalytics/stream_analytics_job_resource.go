package streamanalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/streamanalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/streamanalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceStreamAnalyticsJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceStreamAnalyticsJobCreateUpdate,
		Read:   resourceStreamAnalyticsJobRead,
		Update: resourceStreamAnalyticsJobCreateUpdate,
		Delete: resourceStreamAnalyticsJobDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.StreamingJobID(id)
			return err
		}),

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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"compatibility_level": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					// values found in the other API the portal uses
					string(streamanalytics.OneFullStopZero),
					"1.1",
					// TODO: support for 1.2 when this is fixed:
					// https://github.com/Azure/azure-rest-api-specs/issues/5604
					// "1.2",
				}, false),
			},

			"data_locale": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"events_late_arrival_max_delay_in_seconds": {
				// portal allows for up to 20d 23h 59m 59s
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(-1, 1814399),
				Default:      5,
			},

			"events_out_of_order_max_delay_in_seconds": {
				// portal allows for up to 9m 59s
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 599),
				Default:      0,
			},

			"events_out_of_order_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(streamanalytics.Adjust),
					string(streamanalytics.Drop),
				}, false),
				Default: string(streamanalytics.Adjust),
			},

			"output_error_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(streamanalytics.OutputErrorPolicyDrop),
					string(streamanalytics.OutputErrorPolicyStop),
				}, false),
				Default: string(streamanalytics.OutputErrorPolicyDrop),
			},

			"streaming_units": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.StreamAnalyticsJobStreamingUnits,
			},

			"transformation_query": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceStreamAnalyticsJobCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	transformationsClient := meta.(*clients.Client).StreamAnalytics.TransformationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Job creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Stream Analytics Job %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_stream_analytics_job", *existing.ID)
		}
	}

	compatibilityLevel := d.Get("compatibility_level").(string)
	eventsLateArrivalMaxDelayInSeconds := d.Get("events_late_arrival_max_delay_in_seconds").(int)
	eventsOutOfOrderMaxDelayInSeconds := d.Get("events_out_of_order_max_delay_in_seconds").(int)
	eventsOutOfOrderPolicy := d.Get("events_out_of_order_policy").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	outputErrorPolicy := d.Get("output_error_policy").(string)
	streamingUnits := d.Get("streaming_units").(int)
	transformationQuery := d.Get("transformation_query").(string)
	t := d.Get("tags").(map[string]interface{})

	// needs to be defined inline for a Create but via a separate API for Update
	transformation := streamanalytics.Transformation{
		Name: utils.String("main"),
		TransformationProperties: &streamanalytics.TransformationProperties{
			StreamingUnits: utils.Int32(int32(streamingUnits)),
			Query:          utils.String(transformationQuery),
		},
	}

	props := streamanalytics.StreamingJob{
		Name:     utils.String(name),
		Location: utils.String(location),
		StreamingJobProperties: &streamanalytics.StreamingJobProperties{
			Sku: &streamanalytics.Sku{
				Name: streamanalytics.Standard,
			},
			CompatibilityLevel:                 streamanalytics.CompatibilityLevel(compatibilityLevel),
			EventsLateArrivalMaxDelayInSeconds: utils.Int32(int32(eventsLateArrivalMaxDelayInSeconds)),
			EventsOutOfOrderMaxDelayInSeconds:  utils.Int32(int32(eventsOutOfOrderMaxDelayInSeconds)),
			EventsOutOfOrderPolicy:             streamanalytics.EventsOutOfOrderPolicy(eventsOutOfOrderPolicy),
			OutputErrorPolicy:                  streamanalytics.OutputErrorPolicy(outputErrorPolicy),
		},
		Tags: tags.Expand(t),
	}

	if dataLocale, ok := d.GetOk("data_locale"); ok {
		props.StreamingJobProperties.DataLocale = utils.String(dataLocale.(string))
	}

	if d.IsNewResource() {
		props.StreamingJobProperties.Transformation = &transformation

		future, err := client.CreateOrReplace(ctx, props, resourceGroup, name, "", "")
		if err != nil {
			return fmt.Errorf("Error Creating Stream Analytics Job %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for creation of Stream Analytics Job %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		read, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			return err
		}
		if read.ID == nil {
			return fmt.Errorf("Cannot read ID of Stream Analytics Job %q (Resource Group %q)", name, resourceGroup)
		}

		d.SetId(*read.ID)
	} else {
		if _, err := client.Update(ctx, props, resourceGroup, name, ""); err != nil {
			return fmt.Errorf("Error Updating Stream Analytics Job %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		job, err := client.Get(ctx, resourceGroup, name, "transformation")
		if err != nil {
			return err
		}

		if readTransformation := job.Transformation; readTransformation != nil {
			if _, err := transformationsClient.Update(ctx, transformation, resourceGroup, name, *readTransformation.Name, ""); err != nil {
				return fmt.Errorf("Error Updating Transformation for Stream Analytics Job %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
	}

	return resourceStreamAnalyticsJobRead(d, meta)
}

func resourceStreamAnalyticsJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamingJobID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "transformation")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if resp.Location != nil {
		d.Set("location", azure.NormalizeLocation(*resp.Location))
	}

	if props := resp.StreamingJobProperties; props != nil {
		d.Set("compatibility_level", string(props.CompatibilityLevel))
		d.Set("data_locale", props.DataLocale)
		if props.EventsLateArrivalMaxDelayInSeconds != nil {
			d.Set("events_late_arrival_max_delay_in_seconds", int(*props.EventsLateArrivalMaxDelayInSeconds))
		}
		if props.EventsOutOfOrderMaxDelayInSeconds != nil {
			d.Set("events_out_of_order_max_delay_in_seconds", int(*props.EventsOutOfOrderMaxDelayInSeconds))
		}
		d.Set("events_out_of_order_policy", string(props.EventsOutOfOrderPolicy))
		d.Set("output_error_policy", string(props.OutputErrorPolicy))

		// Computed
		d.Set("job_id", props.JobID)

		if transformation := props.Transformation; transformation != nil {
			if units := transformation.StreamingUnits; units != nil {
				d.Set("streaming_units", int(*units))
			}
			d.Set("transformation_query", transformation.Query)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceStreamAnalyticsJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StreamingJobID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}
