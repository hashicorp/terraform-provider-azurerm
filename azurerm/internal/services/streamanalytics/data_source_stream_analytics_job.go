package streamanalytics

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmStreamAnalyticsJob() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmStreamAnalyticsJobRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"compatibility_level": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"data_locale": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"events_late_arrival_max_delay_in_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"events_out_of_order_max_delay_in_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"events_out_of_order_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"output_error_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"streaming_units": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"transformation_query": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmStreamAnalyticsJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	transformationsClient := meta.(*clients.Client).StreamAnalytics.TransformationsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Stream Analytics Job %q was not found in Resource Group %q!", name, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Stream Analytics Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	transformation, err := transformationsClient.Get(ctx, resourceGroup, name, "Transformation")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Transformation for Stream Analytics Job %q was not found in Resource Group %q!", name, resourceGroup)
		}

		return fmt.Errorf("Error retrieving Transformation for Stream Analytics Job %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

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
		d.Set("job_id", props.JobID)
		d.Set("output_error_policy", string(props.OutputErrorPolicy))
	}

	if props := transformation.TransformationProperties; props != nil {
		if units := props.StreamingUnits; units != nil {
			d.Set("streaming_units", int(*units))
		}
		d.Set("transformation_query", props.Query)
	}

	return nil
}
