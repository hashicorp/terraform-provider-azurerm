// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_stream_input_iothub

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsStreamInputIoTHubCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := inputs.NewInputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_stream_analytics_stream_input_iothub", id.ID())
			}
		}
	}

	consumerGroupName := d.Get("eventhub_consumer_group_name").(string)
	endpoint := d.Get("endpoint").(string)
	iotHubNamespace := d.Get("iothub_namespace").(string)
	sharedAccessPolicyKey := d.Get("shared_access_policy_key").(string)
	sharedAccessPolicyName := d.Get("shared_access_policy_name").(string)

	serializationRaw := d.Get("serialization").([]any)
	serialization, err := helpers.ExpandStreamAnalyticsStreamInputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	props := inputs.Input{
		Name: pointer.To(id.InputName),
		Properties: &inputs.StreamInputProperties{
			Datasource: &inputs.IoTHubStreamInputDataSource{
				Properties: &inputs.IoTHubStreamInputDataSourceProperties{
					ConsumerGroupName:      pointer.To(consumerGroupName),
					SharedAccessPolicyKey:  pointer.To(sharedAccessPolicyKey),
					SharedAccessPolicyName: pointer.To(sharedAccessPolicyName),
					Endpoint:               pointer.To(endpoint),
					IotHubNamespace:        pointer.To(iotHubNamespace),
				},
			},
			Serialization: serialization,
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

	return resourceStreamAnalyticsStreamInputIoTHubRead(d, meta)
}
