// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStreamAnalyticsOutputServiceBusTopic() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsOutputServiceBusTopicCreateUpdate,
		Read:   resourceStreamAnalyticsOutputServiceBusTopicRead,
		Update: resourceStreamAnalyticsOutputServiceBusTopicCreateUpdate,
		Delete: resourceStreamAnalyticsOutputServiceBusTopicDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := outputs.ParseOutputID(id)
			return err
		}, importStreamAnalyticsOutput(outputs.ServiceBusTopicOutputDataSource{})),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsOutputServiceBusTopicV0ToV1{},
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

			"topic_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"servicebus_namespace": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"shared_access_policy_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"shared_access_policy_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"property_columns": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"system_property_columns": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"serialization": schemaStreamAnalyticsOutputSerialization(),

			"authentication_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(outputs.AuthenticationModeConnectionString),
				ValidateFunc: validation.StringInSlice([]string{
					string(outputs.AuthenticationModeMsi),
					string(outputs.AuthenticationModeConnectionString),
				}, false),
			},
		},
	}
}

func resourceStreamAnalyticsOutputServiceBusTopicCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Output ServiceBus Topic creation.")
	id := outputs.NewOutputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_stream_analytics_output_servicebus_topic", id.ID())
		}
	}

	serializationRaw := d.Get("serialization").([]interface{})
	serialization, err := expandStreamAnalyticsOutputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	systemPropertyColumns := d.Get("system_property_columns").(map[string]interface{})
	dataSourceProperties := &outputs.ServiceBusTopicOutputDataSourceProperties{
		TopicName:             pointer.To(d.Get("topic_name").(string)),
		ServiceBusNamespace:   pointer.To(d.Get("servicebus_namespace").(string)),
		PropertyColumns:       utils.ExpandStringSlice(d.Get("property_columns").([]interface{})),
		SystemPropertyColumns: expandSystemPropertyColumns(systemPropertyColumns),
		AuthenticationMode:    pointer.To(outputs.AuthenticationMode(d.Get("authentication_mode").(string))),
	}

	// Add shared access policy key/name only if required by authentication mode
	if *dataSourceProperties.AuthenticationMode == outputs.AuthenticationModeConnectionString {
		dataSourceProperties.SharedAccessPolicyKey = pointer.To(d.Get("shared_access_policy_key").(string))
		dataSourceProperties.SharedAccessPolicyName = pointer.To(d.Get("shared_access_policy_name").(string))
	}

	props := outputs.Output{
		Name: pointer.To(id.OutputName),
		Properties: &outputs.OutputProperties{
			Datasource: &outputs.ServiceBusTopicOutputDataSource{
				Properties: dataSourceProperties,
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
	} else if _, err := client.Update(ctx, id, props, updateOpts); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceStreamAnalyticsOutputServiceBusTopicRead(d, meta)
}

func resourceStreamAnalyticsOutputServiceBusTopicRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

			topicName := ""
			if v := output.Properties.TopicName; v != nil {
				topicName = *v
			}
			d.Set("topic_name", topicName)

			namespace := ""
			if v := output.Properties.ServiceBusNamespace; v != nil {
				namespace = *v
			}
			d.Set("servicebus_namespace", namespace)

			accessPolicy := ""
			if v := output.Properties.SharedAccessPolicyName; v != nil {
				accessPolicy = *v
			}
			d.Set("shared_access_policy_name", accessPolicy)

			var propertyColumns []string
			if v := output.Properties.PropertyColumns; v != nil {
				propertyColumns = *v
			}
			d.Set("property_columns", propertyColumns)

			authMode := ""
			if v := output.Properties.AuthenticationMode; v != nil {
				authMode = string(*v)
			}
			d.Set("authentication_mode", authMode)

			if err = d.Set("system_property_columns", output.Properties.SystemPropertyColumns); err != nil {
				return err
			}

			if err := d.Set("serialization", flattenStreamAnalyticsOutputSerialization(props.Serialization)); err != nil {
				return fmt.Errorf("setting `serialization`: %+v", err)
			}
		}
	}
	return nil
}

func resourceStreamAnalyticsOutputServiceBusTopicDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := outputs.ParseOutputID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}

func expandSystemPropertyColumns(input map[string]interface{}) *map[string]string {
	output := make(map[string]string)
	for k, v := range input {
		output[k] = v.(string)
	}
	return &output
}
