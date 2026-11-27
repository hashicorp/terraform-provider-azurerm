// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_eventhub

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsOutputEventHubCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := outputs.NewOutputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_stream_analytics_output_eventhub", id.ID())
			}
		}
	}

	eventHubName := d.Get("eventhub_name").(string)
	serviceBusNamespace := d.Get("servicebus_namespace").(string)
	sharedAccessPolicyKey := d.Get("shared_access_policy_key").(string)
	sharedAccessPolicyName := d.Get("shared_access_policy_name").(string)
	propertyColumns := d.Get("property_columns").([]any)
	partitionKey := d.Get("partition_key").(string)

	serializationRaw := d.Get("serialization").([]any)
	serialization, err := helpers.ExpandStreamAnalyticsOutputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	eventHubOutputDataSourceProps := &outputs.EventHubOutputDataSourceProperties{
		PartitionKey:        pointer.To(partitionKey),
		PropertyColumns:     pluginsdk.ExpandStringSlice(propertyColumns),
		EventHubName:        pointer.To(eventHubName),
		ServiceBusNamespace: pointer.To(serviceBusNamespace),
		AuthenticationMode:  pointer.ToEnum[outputs.AuthenticationMode](d.Get("authentication_mode").(string)),
	}

	if sharedAccessPolicyKey != "" {
		eventHubOutputDataSourceProps.SharedAccessPolicyKey = &sharedAccessPolicyKey
	}

	if sharedAccessPolicyName != "" {
		eventHubOutputDataSourceProps.SharedAccessPolicyName = &sharedAccessPolicyName
	}

	props := outputs.Output{
		Name: pointer.To(id.OutputName),
		Properties: &outputs.OutputProperties{
			Datasource: &outputs.EventHubOutputDataSource{
				Properties: eventHubOutputDataSourceProps,
			},
			Serialization: serialization,
		},
	}

	var createOpts outputs.CreateOrReplaceOperationOptions
	var updateOpts outputs.UpdateOperationOptions
	if d.IsNewResource() {
		if _, err := client.CreateOrReplace(ctx, id, props, createOpts); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
	} else if _, err := client.Update(ctx, id, props, updateOpts); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStreamAnalyticsOutputEventHubRead(d, meta)
}
