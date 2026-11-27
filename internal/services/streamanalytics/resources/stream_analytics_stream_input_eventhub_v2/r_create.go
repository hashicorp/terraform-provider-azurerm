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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model StreamInputEventHubV2ResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := streamanalyticsClient.From(metadata.Client).InputsClient
			subscriptionId := account.From(metadata.Client).SubscriptionId

			streamingJobStruct, err := streamingjobs.ParseStreamingJobID(model.StreamAnalyticsJobId)
			if err != nil {
				return err
			}
			id := inputs.NewInputID(subscriptionId, streamingJobStruct.ResourceGroupName, streamingJobStruct.StreamingJobName, model.Name)

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			props := &inputs.EventHubStreamInputDataSourceProperties{
				ServiceBusNamespace: pointer.To(model.ServiceBusNamespace),
				EventHubName:        pointer.To(model.EventHubName),
				ConsumerGroupName:   pointer.To(model.EventHubConsumerGroupName),
				AuthenticationMode:  pointer.ToEnum[inputs.AuthenticationMode](model.AuthenticationMode),
			}

			if v := model.SharedAccessPolicyKey; v != "" {
				props.SharedAccessPolicyKey = pointer.To(v)
			}

			if v := model.SharedAccessPolicyName; v != "" {
				props.SharedAccessPolicyName = pointer.To(v)
			}

			serialization, err := expandStreamAnalyticsStreamInputSerializationTyped(model.Serialization)
			if err != nil {
				return fmt.Errorf("expanding `serialization`: %+v", err)
			}

			payload := inputs.Input{
				Name: pointer.To(model.Name),
				Properties: &inputs.StreamInputProperties{
					Datasource: &inputs.EventHubV2StreamInputDataSource{
						Properties: props,
					},
					Serialization: serialization,
					PartitionKey:  pointer.To(model.PartitionKey),
				},
			}

			var opts inputs.CreateOrReplaceOperationOptions
			if _, err = client.CreateOrReplace(ctx, id, payload, opts); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
