// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_stream_input_eventhub

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsStreamInputEventHubCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := streamanalyticsClient.From(meta).InputsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id := inputs.NewInputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_stream_analytics_stream_input_eventhub", id.ID())
			}
		}
	}

	serializationRaw := d.Get("serialization").([]any)
	serialization, err := helpers.ExpandStreamAnalyticsStreamInputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	eventHubDataSourceProps := &inputs.EventHubStreamInputDataSourceProperties{
		EventHubName:        pointer.To(d.Get("eventhub_name").(string)),
		ServiceBusNamespace: pointer.To(d.Get("servicebus_namespace").(string)),
		ConsumerGroupName:   pointer.To(d.Get("eventhub_consumer_group_name").(string)),
		AuthenticationMode:  pointer.ToEnum[inputs.AuthenticationMode](d.Get("authentication_mode").(string)),
	}

	if v, ok := d.GetOk("shared_access_policy_key"); ok {
		eventHubDataSourceProps.SharedAccessPolicyKey = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("shared_access_policy_name"); ok {
		eventHubDataSourceProps.SharedAccessPolicyName = pointer.To(v.(string))
	}

	props := inputs.Input{
		Name: pointer.To(id.InputName),
		Properties: &inputs.StreamInputProperties{
			Datasource: &inputs.EventHubStreamInputDataSource{
				Properties: eventHubDataSourceProps,
			},
			Serialization: serialization,
			PartitionKey:  pointer.To(d.Get("partition_key").(string)),
		},
	}

	var createOpts inputs.CreateOrReplaceOperationOptions
	var updateOpts inputs.UpdateOperationOptions
	if d.IsNewResource() {
		if _, err := client.CreateOrReplace(ctx, id, props, createOpts); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
	} else if _, err := client.Update(ctx, id, props, updateOpts); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStreamAnalyticsStreamInputEventHubRead(d, meta)
}
