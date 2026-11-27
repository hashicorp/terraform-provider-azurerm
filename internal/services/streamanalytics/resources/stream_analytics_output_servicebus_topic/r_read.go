// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_servicebus_topic

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsOutputServiceBusTopicRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := outputs.ParseOutputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.OutputName)
	d.Set("stream_analytics_job_name", id.StreamingJobName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			output, ok := props.Datasource.(outputs.ServiceBusTopicOutputDataSource)
			if !ok {
				return fmt.Errorf("converting %s to a ServiceBus Topic Output", *id)
			}

			d.Set("topic_name", pointer.From(output.Properties.TopicName))

			d.Set("servicebus_namespace", pointer.From(output.Properties.ServiceBusNamespace))

			d.Set("shared_access_policy_name", pointer.From(output.Properties.SharedAccessPolicyName))

			d.Set("property_columns", pointer.From(output.Properties.PropertyColumns))

			authMode := ""
			if v := output.Properties.AuthenticationMode; v != nil {
				authMode = string(*v)
			}
			d.Set("authentication_mode", authMode)

			if err = d.Set("system_property_columns", output.Properties.SystemPropertyColumns); err != nil {
				return err
			}

			if err := d.Set("serialization", helpers.FlattenStreamAnalyticsOutputSerialization(props.Serialization)); err != nil {
				return fmt.Errorf("setting `serialization`: %+v", err)
			}
		}
	}
	return nil
}
