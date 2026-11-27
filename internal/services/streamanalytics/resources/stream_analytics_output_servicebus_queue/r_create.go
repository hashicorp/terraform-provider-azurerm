// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_servicebus_queue

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsOutputServiceBusQueueCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := streamanalyticsClient.From(meta).OutputsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id := outputs.NewOutputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_stream_analytics_output_servicebus_queue", id.ID())
			}
		}
	}

	queueName := d.Get("queue_name").(string)
	serviceBusNamespace := d.Get("servicebus_namespace").(string)
	sharedAccessPolicyKey := d.Get("shared_access_policy_key").(string)
	sharedAccessPolicyName := d.Get("shared_access_policy_name").(string)

	serializationRaw := d.Get("serialization").([]any)
	serialization, err := helpers.ExpandStreamAnalyticsOutputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	systemPropertyColumns := d.Get("system_property_columns")

	dataSourceProperties := &outputs.ServiceBusQueueOutputDataSourceProperties{
		QueueName:             pointer.To(queueName),
		ServiceBusNamespace:   pointer.To(serviceBusNamespace),
		PropertyColumns:       pluginsdk.ExpandStringSlice(d.Get("property_columns").([]any)),
		SystemPropertyColumns: &systemPropertyColumns,
		AuthenticationMode:    pointer.ToEnum[outputs.AuthenticationMode](d.Get("authentication_mode").(string)),
	}

	// Add shared access policy key/name only if required by authentication mode
	if *dataSourceProperties.AuthenticationMode == outputs.AuthenticationModeConnectionString {
		dataSourceProperties.SharedAccessPolicyName = pointer.To(sharedAccessPolicyName)
		dataSourceProperties.SharedAccessPolicyKey = pointer.To(sharedAccessPolicyKey)
	}

	props := outputs.Output{
		Name: pointer.To(id.OutputName),
		Properties: &outputs.OutputProperties{
			Datasource: &outputs.ServiceBusQueueOutputDataSource{
				Properties: dataSourceProperties,
			},
			Serialization: serialization,
		},
	}

	// TODO: split the create/update functions to allow for ignore changes etc
	var createOpts outputs.CreateOrReplaceOperationOptions
	var updateOpts outputs.UpdateOperationOptions
	if d.IsNewResource() {
		if _, err := client.CreateOrReplace(ctx, id, props, createOpts); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
	} else if _, err := client.Update(ctx, id, props, updateOpts); err != nil {
		return fmt.Errorf("uUpdating %s: %+v", id, err)
	}

	return resourceStreamAnalyticsOutputServiceBusQueueRead(d, meta)
}
