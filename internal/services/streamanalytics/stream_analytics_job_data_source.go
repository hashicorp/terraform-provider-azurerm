// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceStreamAnalyticsJob() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceStreamAnalyticsJobRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"location": commonschema.LocationComputed(),

			"compatibility_level": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"data_locale": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"events_late_arrival_max_delay_in_seconds": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"events_out_of_order_max_delay_in_seconds": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"events_out_of_order_policy": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"job_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

			"last_output_time": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"output_error_policy": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"start_mode": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"start_time": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"streaming_units": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"transformation_query": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStreamAnalyticsJobRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.JobsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := streamingjobs.NewStreamingJobID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	opts := streamingjobs.GetOperationOptions{
		Expand: pointer.To("transformation"),
	}
	resp, err := client.Get(ctx, id, opts)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.StreamingJobName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		flattenedIdentity, err := identity.FlattenSystemOrUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %v", err)
		}

		if props := model.Properties; props != nil {
			compatibilityLevel := ""
			if v := props.CompatibilityLevel; v != nil {
				compatibilityLevel = string(*v)
			}
			d.Set("compatibility_level", compatibilityLevel)

			dataLocale := ""
			if v := props.DataLocale; v != nil {
				dataLocale = *v
			}
			d.Set("data_locale", dataLocale)

			var lateArrival int64
			if v := props.EventsLateArrivalMaxDelayInSeconds; v != nil {
				lateArrival = *v
			}
			d.Set("events_late_arrival_max_delay_in_seconds", lateArrival)

			var maxDelay int64
			if v := props.EventsLateArrivalMaxDelayInSeconds; v != nil {
				maxDelay = *v
			}
			d.Set("events_out_of_order_max_delay_in_seconds", maxDelay)

			orderPolicy := ""
			if v := props.EventsOutOfOrderPolicy; v != nil {
				orderPolicy = string(*v)
			}
			d.Set("events_out_of_order_policy", orderPolicy)

			outputPolicy := ""
			if v := props.OutputErrorPolicy; v != nil {
				outputPolicy = string(*v)
			}
			d.Set("output_error_policy", outputPolicy)

			lastOutputTime := ""
			if v := props.LastOutputEventTime; v != nil {
				lastOutputTime = *v
			}
			d.Set("last_output_time", lastOutputTime)

			startTime := ""
			if v := props.OutputStartTime; v != nil {
				startTime = *v
			}
			d.Set("start_time", startTime)

			startMode := ""
			if v := props.OutputStartMode; v != nil {
				startMode = string(*v)
			}
			d.Set("start_mode", startMode)

			jobId := ""
			if v := props.JobId; v != nil {
				jobId = *v
			}
			d.Set("job_id", jobId)

			sku := ""
			if props.Sku != nil && props.Sku.Name != nil {
				sku = string(*props.Sku.Name)
			}
			d.Set("sku_name", sku)

			if props.Transformation != nil && props.Transformation.Properties != nil {
				var streamingUnits int64
				if v := props.Transformation.Properties.StreamingUnits; v != nil {
					streamingUnits = *v
				}
				d.Set("streaming_units", streamingUnits)

				query := ""
				if v := props.Transformation.Properties.Query; v != nil {
					query = *v
				}
				d.Set("transformation_query", query)
			}
		}
	}
	return nil
}
