package iottimeseriesinsights

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/timeseriesinsights/mgmt/2018-08-15-preview/timeseriesinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	iothubValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iottimeseriesinsights/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iottimeseriesinsights/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmIoTTimeSeriesInsightsEventSourceIoTHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIoTTimeSeriesInsightsEventSourceIoTHubCreateUpdate,
		Read:   resourceArmIoTTimeSeriesInsightsEventSourceIoTHubRead,
		Update: resourceArmIoTTimeSeriesInsightsEventSourceIoTHubCreateUpdate,
		Delete: resourceArmIoTTimeSeriesInsightsEventSourceIoTHubDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.TimeSeriesInsightsEventSourceID(id)
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[-\w\._\(\)]+$`),
					"IoT Time Series Insights Event Source name must contain only word characters, periods, underscores, hyphens, and parentheses.",
				),
			},

			"time_series_insights_environment_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.TimeSeriesInsightsEnvironmentID,
			},

			"location": azure.SchemaLocation(),

			"iothub_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: iothubValidate.IoTHubName,
			},

			"event_source_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"consumer_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"key_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"shared_access_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Sensitive:    true,
			},

			"timestamp_property_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmIoTTimeSeriesInsightsEventSourceIoTHubCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EventSourcesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	environmentID := d.Get("time_series_insights_environment_id").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	id, err := parse.TimeSeriesInsightsEnvironmentID(environmentID)
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing IoT Time Series Insights Event Source IoTHub %q (Resource Group %q): %s", name, id.ResourceGroup, err)
			}
		}

		if existing.Value != nil {
			eventSource, ok := existing.Value.AsIoTHubEventSourceResource()
			if !ok {
				return fmt.Errorf("exisiting resource was not a standard IoT Time Series Insights IoTHub Event Source %q (Resource Group %q)", name, id.ResourceGroup)
			}

			if eventSource.ID != nil && *eventSource.ID != "" {
				return tf.ImportAsExistsError("azurerm_iot_time_series_insights_event_source_iothub", *eventSource.ID)
			}
		}
	}

	policy := timeseriesinsights.IoTHubEventSourceCreateOrUpdateParameters{
		Kind:     timeseriesinsights.KindMicrosoftIoTHub,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Location: &location,
		IoTHubEventSourceCreationProperties: &timeseriesinsights.IoTHubEventSourceCreationProperties{
			IotHubName:            utils.String(d.Get("iothub_name").(string)),
			EventSourceResourceID: utils.String(d.Get("event_source_resource_id").(string)),
			ConsumerGroupName:     utils.String(d.Get("consumer_group_name").(string)),
			KeyName:               utils.String(d.Get("key_name").(string)),
			SharedAccessKey:       utils.String(d.Get("shared_access_key").(string)),
			TimestampPropertyName: utils.String(d.Get("timestamp_property_name").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, name, policy); err != nil {
		return fmt.Errorf("creating/updating IoT Time Series Insights Event Source IoTHub %q (Resource Group %q): %+v", name, id.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving IoT Time Series Insights Event Source IoTHub %q (Resource Group %q): %+v", name, id.ResourceGroup, err)
	}

	eventSource, ok := resp.Value.AsEventSourceResource()
	return fmt.Errorf("%+v", eventSource)
	if !ok {
		return fmt.Errorf("%+v", eventSource)
		return fmt.Errorf("created resource was not a IoT Time Series Insights IoTHub Event Source %q (Resource Group %q)", name, id.ResourceGroup)
	}

	if eventSource.ID == nil || *eventSource.ID == "" {
		return fmt.Errorf("cannot read IoT Time Series Insights Event Source IoTHub %q (Resource Group %q) ID", name, id.ResourceGroup)
	}

	d.SetId(*eventSource.ID)

	return resourceArmIoTTimeSeriesInsightsEventSourceIoTHubRead(d, meta)
}

func resourceArmIoTTimeSeriesInsightsEventSourceIoTHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EventSourcesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TimeSeriesInsightsEventSourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving IoT Time Series Insights Event Source IoTHub %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	eventSource, ok := resp.Value.AsIoTHubEventSourceResource()
	if !ok {
		return fmt.Errorf("exisiting resource was not a standard IoT Time Series Insights IoTHub Event Source %q (Resource Group %q)", id.Name, id.ResourceGroup)
	}

	d.Set("name", eventSource.Name)
	d.Set("time_series_insights_environment_id", strings.Split(d.Id(), "/eventsources")[0])

	if props := eventSource.IoTHubEventSourceResourceProperties; props != nil {
		d.Set("iothub_name", props.IotHubName)
		d.Set("event_source_resource_id", props.EventSourceResourceID)
		d.Set("consumer_group_name", props.ConsumerGroupName)
		d.Set("key_name", props.KeyName)
		d.Set("timestamp_property_name", props.TimestampPropertyName)
	}

	return tags.FlattenAndSet(d, eventSource.Tags)
}

func resourceArmIoTTimeSeriesInsightsEventSourceIoTHubDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EventSourcesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TimeSeriesInsightsEventSourceID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting IoT Time Series Insights Event Source IoTHub %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}
