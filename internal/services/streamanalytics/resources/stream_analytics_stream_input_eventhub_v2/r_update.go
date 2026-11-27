// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_stream_input_eventhub_v2

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.InputsClient
			id, err := inputs.ParseInputID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state StreamInputEventHubV2ResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			d := metadata.ResourceData

			if d.HasChangesExcept("name", "stream_analytics_job_id") {
				props := &inputs.EventHubStreamInputDataSourceProperties{
					ServiceBusNamespace: pointer.To(state.ServiceBusNamespace),
					EventHubName:        pointer.To(state.EventHubName),
					ConsumerGroupName:   pointer.To(state.EventHubConsumerGroupName),
					AuthenticationMode:  pointer.ToEnum[inputs.AuthenticationMode](state.AuthenticationMode),
				}

				serialization, err := expandStreamAnalyticsStreamInputSerializationTyped(state.Serialization)
				if err != nil {
					return fmt.Errorf("expanding `serialization`: %+v", err)
				}

				payload := inputs.Input{
					Name: pointer.To(state.Name),
					Properties: &inputs.StreamInputProperties{
						Datasource: &inputs.EventHubV2StreamInputDataSource{
							Properties: props,
						},
						Serialization: serialization,
						PartitionKey:  pointer.To(state.PartitionKey),
					},
				}

				var opts inputs.UpdateOperationOptions
				if _, err := client.Update(ctx, *id, payload, opts); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
