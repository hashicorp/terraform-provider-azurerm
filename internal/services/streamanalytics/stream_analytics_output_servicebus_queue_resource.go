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

func resourceStreamAnalyticsOutputServiceBusQueue() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStreamAnalyticsOutputServiceBusQueueCreateUpdate,
		Read:   resourceStreamAnalyticsOutputServiceBusQueueRead,
		Update: resourceStreamAnalyticsOutputServiceBusQueueCreateUpdate,
		Delete: resourceStreamAnalyticsOutputServiceBusQueueDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := outputs.ParseOutputID(id)
			return err
		}, importStreamAnalyticsOutput(outputs.ServiceBusQueueOutputDataSource{})),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsOutputServiceBusQueueV0ToV1{},
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

			"queue_name": {
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

			"serialization": schemaStreamAnalyticsOutputSerialization(),

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

			"authentication_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(outputs.AuthenticationModeConnectionString),
				ValidateFunc: validation.StringInSlice([]string{
					string(outputs.AuthenticationModeConnectionString),
					string(outputs.AuthenticationModeMsi),
				}, false),
			},
		},
	}
}

func resourceStreamAnalyticsOutputServiceBusQueueCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := outputs.NewOutputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
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

	queueName := d.Get("queue_name").(string)
	serviceBusNamespace := d.Get("servicebus_namespace").(string)
	sharedAccessPolicyKey := d.Get("shared_access_policy_key").(string)
	sharedAccessPolicyName := d.Get("shared_access_policy_name").(string)

	serializationRaw := d.Get("serialization").([]interface{})
	serialization, err := expandStreamAnalyticsOutputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	systemPropertyColumns := d.Get("system_property_columns")

	dataSourceProperties := &outputs.ServiceBusQueueOutputDataSourceProperties{
		QueueName:             pointer.To(queueName),
		ServiceBusNamespace:   pointer.To(serviceBusNamespace),
		PropertyColumns:       utils.ExpandStringSlice(d.Get("property_columns").([]interface{})),
		SystemPropertyColumns: &systemPropertyColumns,
		AuthenticationMode:    pointer.To(outputs.AuthenticationMode(d.Get("authentication_mode").(string))),
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

func resourceStreamAnalyticsOutputServiceBusQueueRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.OutputName)
	d.Set("stream_analytics_job_name", id.StreamingJobName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			output, ok := props.Datasource.(outputs.ServiceBusQueueOutputDataSource)
			if !ok {
				return fmt.Errorf("converting %s to a ServiceBus Queue Output", *id)
			}

			queue := ""
			if v := output.Properties.QueueName; v != nil {
				queue = *v
			}
			d.Set("queue_name", queue)

			namespace := ""
			if v := output.Properties.ServiceBusNamespace; v != nil {
				namespace = *v
			}
			d.Set("servicebus_namespace", namespace)

			policyName := ""
			if v := output.Properties.SharedAccessPolicyName; v != nil {
				policyName = *v
			}
			d.Set("shared_access_policy_name", policyName)

			var columns []string
			if v := output.Properties.PropertyColumns; v != nil {
				columns = *v
			}
			d.Set("property_columns", columns)

			var systemColumns interface{}
			if v := output.Properties.SystemPropertyColumns; v != nil {
				systemColumns = *v
			}
			d.Set("system_property_columns", systemColumns)

			authMode := ""
			if v := output.Properties.AuthenticationMode; v != nil {
				authMode = string(*v)
			}
			d.Set("authentication_mode", authMode)

			if err := d.Set("serialization", flattenStreamAnalyticsOutputSerialization(props.Serialization)); err != nil {
				return fmt.Errorf("setting `serialization`: %+v", err)
			}
		}
	}
	return nil
}

func resourceStreamAnalyticsOutputServiceBusQueueDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := outputs.ParseOutputID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}
