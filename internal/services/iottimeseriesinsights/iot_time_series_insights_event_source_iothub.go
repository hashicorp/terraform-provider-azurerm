// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iottimeseriesinsights

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/environments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/eventsources"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIoTTimeSeriesInsightsEventSourceIoTHub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIoTTimeSeriesInsightsEventSourceIoTHubCreateUpdate,
		Read:   resourceIoTTimeSeriesInsightsEventSourceIoTHubRead,
		Update: resourceIoTTimeSeriesInsightsEventSourceIoTHubCreateUpdate,
		Delete: resourceIoTTimeSeriesInsightsEventSourceIoTHubDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := eventsources.ParseEventSourceID(id)
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

			"location": commonschema.Location(),

			"environment_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: environments.ValidateEnvironmentID,
			},

			"iothub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: iothubValidate.IoTHubName,
			},

			"shared_access_key": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
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

			"tags": commonschema.Tags(),
		},
	}
}

func resourceIoTTimeSeriesInsightsEventSourceIoTHubCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EventSources
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	environmentID, err := environments.ParseEnvironmentID(d.Get("environment_id").(string))
	if err != nil {
		return fmt.Errorf("unable to parse `environment_id`: %+v", err)
	}

	id := eventsources.NewEventSourceID(environmentID.SubscriptionId, environmentID.ResourceGroupName, environmentID.EnvironmentName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing IoT Time Series Insights Event Source %q: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_iot_time_series_insights_event_source_iothub", id.ID())
		}
	}

	eventSource := eventsources.IoTHubEventSourceCreateOrUpdateParameters{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: eventsources.IoTHubEventSourceCreationProperties{
			IotHubName:            d.Get("iothub_name").(string),
			SharedAccessKey:       d.Get("shared_access_key").(string),
			ConsumerGroupName:     d.Get("consumer_group_name").(string),
			KeyName:               d.Get("shared_access_key_name").(string),
			EventSourceResourceId: utils.String(d.Get("event_source_resource_id").(string)),
			TimestampPropertyName: utils.String(d.Get("timestamp_property_name").(string)),
		},
	}

	if _, err = client.CreateOrUpdate(ctx, id, eventSource); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIoTTimeSeriesInsightsEventSourceIoTHubRead(d, meta)
}

func resourceIoTTimeSeriesInsightsEventSourceIoTHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EventSources
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := eventsources.ParseEventSourceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.EventSourceName)
	d.Set("environment_id", environments.NewEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.EnvironmentName).ID())
	if model := resp.Model; model != nil {
		eventSource, ok := (*model).(eventsources.IoTHubEventSourceResource)
		if !ok {
			return fmt.Errorf("retrieving %s: expected an IoTHubEventSourceResource but got: %+v", *id, err)
		}

		d.Set("location", location.Normalize(eventSource.Location))

		d.Set("iothub_name", eventSource.Properties.IotHubName)
		d.Set("consumer_group_name", eventSource.Properties.ConsumerGroupName)
		d.Set("shared_access_key_name", eventSource.Properties.KeyName)
		d.Set("event_source_resource_id", eventSource.Properties.EventSourceResourceId)
		d.Set("timestamp_property_name", eventSource.Properties.TimestampPropertyName)

		if err := tags.FlattenAndSet(d, eventSource.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceIoTTimeSeriesInsightsEventSourceIoTHubDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EventSources
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := eventsources.ParseEventSourceID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
