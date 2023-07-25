// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStreamAnalyticsStreamInputIoTHub() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsStreamInputIoTHubCreateUpdate,
		Read:   resourceStreamAnalyticsStreamInputIoTHubRead,
		Update: resourceStreamAnalyticsStreamInputIoTHubCreateUpdate,
		Delete: resourceStreamAnalyticsStreamInputIoTHubDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := inputs.ParseInputID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsStreamInputIotHubV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"stream_analytics_job_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"endpoint": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"iothub_namespace": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"eventhub_consumer_group_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"shared_access_policy_key": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"shared_access_policy_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"serialization": schemaStreamAnalyticsStreamInputSerialization(),
		},
	}
}

func resourceStreamAnalyticsStreamInputIoTHubCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Stream Input IoTHub creation.")

	id := inputs.NewInputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
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

	consumerGroupName := d.Get("eventhub_consumer_group_name").(string)
	endpoint := d.Get("endpoint").(string)
	iotHubNamespace := d.Get("iothub_namespace").(string)
	sharedAccessPolicyKey := d.Get("shared_access_policy_key").(string)
	sharedAccessPolicyName := d.Get("shared_access_policy_name").(string)

	serializationRaw := d.Get("serialization").([]interface{})
	serialization, err := expandStreamAnalyticsStreamInputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	props := inputs.Input{
		Name: utils.String(id.InputName),
		Properties: &inputs.StreamInputProperties{
			Datasource: &inputs.IoTHubStreamInputDataSource{
				Properties: &inputs.IoTHubStreamInputDataSourceProperties{
					ConsumerGroupName:      utils.String(consumerGroupName),
					SharedAccessPolicyKey:  utils.String(sharedAccessPolicyKey),
					SharedAccessPolicyName: utils.String(sharedAccessPolicyName),
					Endpoint:               utils.String(endpoint),
					IotHubNamespace:        utils.String(iotHubNamespace),
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

func resourceStreamAnalyticsStreamInputIoTHubRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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
			input, ok := props.(inputs.InputProperties) // nolint: gosimple
			if !ok {
				return fmt.Errorf("converting %s to an Input", *id)
			}

			streamInput, ok := input.(inputs.StreamInputProperties)
			if !ok {
				return fmt.Errorf("converting %s to a Stream Input", *id)
			}

			streamIotHubInput, ok := streamInput.Datasource.(inputs.IoTHubStreamInputDataSource)
			if !ok {
				return fmt.Errorf("converting %s Stream Input Blob to an Stream Input", *id)
			}

			if streamIotHubInputProps := streamIotHubInput.Properties; streamIotHubInputProps != nil {
				eventHubConsumerGroupName := ""
				if v := streamIotHubInputProps.ConsumerGroupName; v != nil {
					eventHubConsumerGroupName = *v
				}
				d.Set("eventhub_consumer_group_name", eventHubConsumerGroupName)

				endpoint := ""
				if v := streamIotHubInputProps.Endpoint; v != nil {
					endpoint = *v
				}
				d.Set("endpoint", endpoint)

				iothubNamespace := ""
				if v := streamIotHubInputProps.IotHubNamespace; v != nil {
					iothubNamespace = *v
				}
				d.Set("iothub_namespace", iothubNamespace)

				sharedAccessPolicyName := ""
				if v := streamIotHubInputProps.SharedAccessPolicyName; v != nil {
					sharedAccessPolicyName = *v
				}
				d.Set("shared_access_policy_name", sharedAccessPolicyName)

				if err := d.Set("serialization", flattenStreamAnalyticsStreamInputSerialization(streamInput.Serialization)); err != nil {
					return fmt.Errorf("setting `serialization`: %+v", err)
				}
			}
		}
	}

	return nil
}

func resourceStreamAnalyticsStreamInputIoTHubDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := inputs.ParseInputID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
