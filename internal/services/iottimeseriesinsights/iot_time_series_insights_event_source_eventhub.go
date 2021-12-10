package iottimeseriesinsights

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/timeseriesinsights/mgmt/2020-05-15/timeseriesinsights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iottimeseriesinsights/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIoTTimeSeriesInsightsEventSourceEventhub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIoTTimeSeriesInsightsEventSourceEventhubCreateUpdate,
		Read:   resourceIoTTimeSeriesInsightsEventSourceEventhubRead,
		Update: resourceIoTTimeSeriesInsightsEventSourceEventhubCreateUpdate,
		Delete: resourceIoTTimeSeriesInsightsEventSourceEventhubDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EventSourceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[-\w\._\(\)]+$`),
					"IoT Time Series Insights Event Source name must contain only word characters, periods, underscores, and parentheses.",
				),
			},

			"location": azure.SchemaLocation(),

			"environment_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"eventhub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: eventhubValidate.ValidateEventHubName(),
			},

			"namespace_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"shared_access_key": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"consumer_group_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"shared_access_key_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"event_source_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"timestamp_property_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceIoTTimeSeriesInsightsEventSourceEventhubCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EventSourcesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := azure.NormalizeLocation(d.Get("location").(string))
	environmentID, err := parse.EnvironmentID(d.Get("environment_id").(string))
	if err != nil {
		return fmt.Errorf("unable to parse `environment_id`: %+v", err)
	}

	id := parse.NewEventSourceID(environmentID.SubscriptionId, environmentID.ResourceGroup, environmentID.Name, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing IoT Time Series Insights Event Source %q: %s", id, err)
			}
		}

		if existing.Value != nil {
			eventSource, ok := existing.Value.AsEventHubEventSourceResource()
			if !ok {
				return fmt.Errorf("exisiting resource was not an IoT Time Series Insights IoTHub Event Source %q", id)
			}

			if eventSource.ID != nil && *eventSource.ID != "" {
				return tf.ImportAsExistsError("azurerm_iot_time_series_insights_event_source_eventhub", *eventSource.ID)
			}
		}
	}

	eventSource := timeseriesinsights.EventHubEventSourceCreateOrUpdateParameters{
		Location: &location,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		EventHubEventSourceCreationProperties: &timeseriesinsights.EventHubEventSourceCreationProperties{
			EventHubName:          utils.String(d.Get("eventhub_name").(string)),
			ServiceBusNamespace:   utils.String(d.Get("namespace_name").(string)),
			SharedAccessKey:       utils.String(d.Get("shared_access_key").(string)),
			ConsumerGroupName:     utils.String(d.Get("consumer_group_name").(string)),
			KeyName:               utils.String(d.Get("shared_access_key_name").(string)),
			EventSourceResourceID: utils.String(d.Get("event_source_resource_id").(string)),
			TimestampPropertyName: utils.String(d.Get("timestamp_property_name").(string)),
		},
	}

	if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.EnvironmentName, id.Name, eventSource); err != nil {
		return fmt.Errorf("creating/updating IoT Time Series Insights EventHub Event Source %q: %+v", id, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving IoT Time Series Insights EventHub EventSource %q: %+v", id, err)
	}

	if _, ok := resp.Value.AsEventHubEventSourceResource(); !ok {
		return fmt.Errorf("created resource was not an IoT Time Series Insights EventHub Event Source %q", id)
	}

	d.SetId(id.ID())

	return resourceIoTTimeSeriesInsightsEventSourceEventhubRead(d, meta)
}

func resourceIoTTimeSeriesInsightsEventSourceEventhubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EventSourcesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EventSourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
	if err != nil || resp.Value == nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving IoT Time Series Insights EventHub EventSource %q: %+v", id, err)
	}

	eventSource, ok := resp.Value.AsEventHubEventSourceResource()
	if !ok {
		return fmt.Errorf("exisiting resource was not a IoT Time Series Insights EventHub EventSource %q", id)
	}

	d.Set("name", eventSource.Name)
	d.Set("environment_id", parse.NewEnvironmentID(id.SubscriptionId, id.ResourceGroup, id.EnvironmentName).ID())
	if location := eventSource.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := eventSource.EventHubEventSourceResourceProperties; props != nil {
		d.Set("eventhub_name", props.EventHubName)
		d.Set("namespace_name", props.ServiceBusNamespace)
		d.Set("consumer_group_name", props.ConsumerGroupName)
		d.Set("shared_access_key_name", props.KeyName)
		d.Set("event_source_resource_id", props.EventSourceResourceID)
		d.Set("timestamp_property_name", props.TimestampPropertyName)
	}

	return tags.FlattenAndSet(d, eventSource.Tags)
}

func resourceIoTTimeSeriesInsightsEventSourceEventhubDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EventSourcesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EventSourceID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting IoT Time Series Insights EventHub Event Source %q: %+v", id, err)
		}
	}

	return nil
}
