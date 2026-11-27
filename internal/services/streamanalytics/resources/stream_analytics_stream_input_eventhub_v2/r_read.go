// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_stream_input_eventhub_v2

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.InputsClient
			id, err := inputs.ParseInputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			streamingJobId := streamingjobs.NewStreamingJobID(id.SubscriptionId, id.ResourceGroupName, id.StreamingJobName)

			state := StreamInputEventHubV2ResourceModel{
				Name:                 id.InputName,
				StreamAnalyticsJobId: streamingJobId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					streamInput, ok := props.(inputs.StreamInputProperties)
					if !ok {
						return fmt.Errorf("converting %s to a Stream Input", *id)
					}

					eventHubV2Input, ok := streamInput.Datasource.(inputs.EventHubV2StreamInputDataSource)
					if !ok {
						return fmt.Errorf("converting %s to an EventHub V2 Stream Input", *id)
					}

					if eventHubV2InputProps := eventHubV2Input.Properties; eventHubV2InputProps != nil {
						authenticationMode := ""
						if v := eventHubV2InputProps.AuthenticationMode; v != nil {
							authenticationMode = string(*v)
						}

						serialization := flattenStreamAnalyticsStreamInputSerializationTyped(streamInput.Serialization)

						state.ServiceBusNamespace = pointer.From(eventHubV2InputProps.ServiceBusNamespace)
						state.EventHubName = pointer.From(eventHubV2InputProps.EventHubName)
						state.EventHubConsumerGroupName = pointer.From(eventHubV2InputProps.ConsumerGroupName)
						state.AuthenticationMode = authenticationMode
						state.SharedAccessPolicyName = pointer.From(eventHubV2InputProps.SharedAccessPolicyName)
						state.SharedAccessPolicyKey = metadata.ResourceData.Get("shared_access_policy_key").(string)
						state.Serialization = []Serialization{serialization}
						state.PartitionKey = pointer.From(streamInput.PartitionKey)

						return metadata.Encode(&state)
					}
				}
			}
			return nil
		},
	}
}

type Serialization struct {
	Type           string `tfschema:"type"`
	FieldDelimiter string `tfschema:"field_delimiter"`
	Encoding       string `tfschema:"encoding"`
}
