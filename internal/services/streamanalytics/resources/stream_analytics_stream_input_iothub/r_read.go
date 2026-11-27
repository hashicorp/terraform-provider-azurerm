// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_stream_input_iothub

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsStreamInputIoTHubRead(d *pluginsdk.ResourceData, meta any) error {
	client := streamanalyticsClient.From(meta).InputsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := inputs.ParseInputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.InputName)
	d.Set("stream_analytics_job_name", id.StreamingJobName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			streamInput, ok := props.(inputs.StreamInputProperties)
			if !ok {
				return fmt.Errorf("converting %s to a Stream Input", *id)
			}

			streamIotHubInput, ok := streamInput.Datasource.(inputs.IoTHubStreamInputDataSource)
			if !ok {
				return fmt.Errorf("converting %s Stream Input Blob to an Stream Input", *id)
			}

			if streamIotHubInputProps := streamIotHubInput.Properties; streamIotHubInputProps != nil {
				d.Set("eventhub_consumer_group_name", pointer.From(streamIotHubInputProps.ConsumerGroupName))

				d.Set("endpoint", pointer.From(streamIotHubInputProps.Endpoint))

				d.Set("iothub_namespace", pointer.From(streamIotHubInputProps.IotHubNamespace))

				d.Set("shared_access_policy_name", pointer.From(streamIotHubInputProps.SharedAccessPolicyName))

				if err := d.Set("serialization", helpers.FlattenStreamAnalyticsStreamInputSerialization(streamInput.Serialization)); err != nil {
					return fmt.Errorf("setting `serialization`: %+v", err)
				}
			}
		}
	}

	return nil
}
