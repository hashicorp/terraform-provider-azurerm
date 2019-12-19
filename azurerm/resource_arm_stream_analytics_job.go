package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStreamAnalyticsJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStreamAnalyticsJobCreateUpdate,
		Read:   resourceArmStreamAnalyticsJobRead,
		Update: resourceArmStreamAnalyticsJobCreateUpdate,
		Delete: resourceArmStreamAnalyticsJobDelete,
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
				ValidateFunc: validate.NoEmptyStrings,
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
					//"1.2",
				}, false),
			},

			"data_locale": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.NoEmptyStrings,
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
				ValidateFunc: validate.NoEmptyStrings,
			},

			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmStreamAnalyticsJobCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	transformationsClient := meta.(*clients.Client).StreamAnalytics.TransformationsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Job creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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
		Name: utils.String("Transformation"),
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

		if _, err := transformationsClient.Update(ctx, transformation, resourceGroup, name, "Transformation", ""); err != nil {
			return fmt.Errorf("Error Updating Transformation for Stream Analytics Job %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return resourceArmStreamAnalyticsJobRead(d, meta)
}

func resourceArmStreamAnalyticsJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	transformationsClient := meta.(*clients.Client).StreamAnalytics.TransformationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["streamingjobs"]

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Stream Analytics Job %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Stream Analytics Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	transformation, err := transformationsClient.Get(ctx, resourceGroup, name, "Transformation")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Transformation for Stream Analytics Job %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Transformation for Stream Analytics Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

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
	}

	if props := transformation.TransformationProperties; props != nil {
		if units := props.StreamingUnits; units != nil {
			d.Set("streaming_units", int(*units))
		}
		d.Set("transformation_query", props.Query)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmStreamAnalyticsJobDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["streamingjobs"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Stream Analytics Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion for Stream Analytics Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}
