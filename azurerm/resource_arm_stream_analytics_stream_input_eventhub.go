package azurerm

import (
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"

	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStreamAnalyticsStreamInputEventHub() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStreamAnalyticsStreamInputEventHubCreateUpdate,
		Read:   resourceArmStreamAnalyticsStreamInputEventHubRead,
		Update: resourceArmStreamAnalyticsStreamInputEventHubCreateUpdate,
		Delete: resourceArmStreamAnalyticsStreamInputEventHubDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"stream_analytics_job_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"eventhub_consumer_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"eventhub_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"servicebus_namespace": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"shared_access_policy_key": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"shared_access_policy_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"serialization": azure.SchemaStreamAnalyticsStreamInputSerialization(),
		},
	}
}

func resourceArmStreamAnalyticsStreamInputEventHubCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).streamanalytics.InputsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Stream Input EventHub creation.")
	name := d.Get("name").(string)
	jobName := d.Get("stream_analytics_job_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Stream Analytics Stream Input %q (Job %q / Resource Group %q): %s", name, jobName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_stream_analytics_stream_input_eventhub", *existing.ID)
		}
	}

	consumerGroupName := d.Get("eventhub_consumer_group_name").(string)
	eventHubName := d.Get("eventhub_name").(string)
	serviceBusNamespace := d.Get("servicebus_namespace").(string)
	sharedAccessPolicyKey := d.Get("shared_access_policy_key").(string)
	sharedAccessPolicyName := d.Get("shared_access_policy_name").(string)

	serializationRaw := d.Get("serialization").([]interface{})
	serialization, err := azure.ExpandStreamAnalyticsStreamInputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("Error expanding `serialization`: %+v", err)
	}

	props := streamanalytics.Input{
		Name: utils.String(name),
		Properties: &streamanalytics.StreamInputProperties{
			Type: streamanalytics.TypeStream,
			Datasource: &streamanalytics.EventHubStreamInputDataSource{
				Type: streamanalytics.TypeBasicStreamInputDataSourceTypeMicrosoftServiceBusEventHub,
				EventHubStreamInputDataSourceProperties: &streamanalytics.EventHubStreamInputDataSourceProperties{
					ConsumerGroupName:      utils.String(consumerGroupName),
					EventHubName:           utils.String(eventHubName),
					ServiceBusNamespace:    utils.String(serviceBusNamespace),
					SharedAccessPolicyKey:  utils.String(sharedAccessPolicyKey),
					SharedAccessPolicyName: utils.String(sharedAccessPolicyName),
				},
			},
			Serialization: serialization,
		},
	}

	if d.IsNewResource() {
		if _, err := client.CreateOrReplace(ctx, props, resourceGroup, jobName, name, "", ""); err != nil {
			return fmt.Errorf("Error Creating Stream Analytics Stream Input EventHub %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}

		read, err := client.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			return fmt.Errorf("Error retrieving Stream Analytics Stream Input EventHub %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}
		if read.ID == nil {
			return fmt.Errorf("Cannot read ID of Stream Analytics Stream Input EventHub %q (Job %q / Resource Group %q)", name, jobName, resourceGroup)
		}

		d.SetId(*read.ID)
	} else {
		if _, err := client.Update(ctx, props, resourceGroup, jobName, name, ""); err != nil {
			return fmt.Errorf("Error Updating Stream Analytics Stream Input EventHub %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}
	}

	return resourceArmStreamAnalyticsStreamInputEventHubRead(d, meta)
}

func resourceArmStreamAnalyticsStreamInputEventHubRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).streamanalytics.InputsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	jobName := id.Path["streamingjobs"]
	name := id.Path["inputs"]

	resp, err := client.Get(ctx, resourceGroup, jobName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Stream Input EventHub %q was not found in Stream Analytics Job %q / Resource Group %q - removing from state!", name, jobName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Stream Input EventHub %q (Stream Analytics Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("stream_analytics_job_name", jobName)

	if props := resp.Properties; props != nil {
		v, ok := props.AsStreamInputProperties()
		if !ok {
			return fmt.Errorf("Error converting Stream Input EventHub to an Stream Input: %+v", err)
		}

		eventHub, ok := v.Datasource.AsEventHubStreamInputDataSource()
		if !ok {
			return fmt.Errorf("Error converting Stream Input EventHub to an EventHub Stream Input: %+v", err)
		}

		d.Set("eventhub_consumer_group_name", eventHub.ConsumerGroupName)
		d.Set("eventhub_name", eventHub.EventHubName)
		d.Set("servicebus_namespace", eventHub.ServiceBusNamespace)
		d.Set("shared_access_policy_name", eventHub.SharedAccessPolicyName)

		if err := d.Set("serialization", azure.FlattenStreamAnalyticsStreamInputSerialization(v.Serialization)); err != nil {
			return fmt.Errorf("Error setting `serialization`: %+v", err)
		}
	}

	return nil
}

func resourceArmStreamAnalyticsStreamInputEventHubDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).streamanalytics.InputsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	jobName := id.Path["streamingjobs"]
	name := id.Path["inputs"]

	if resp, err := client.Delete(ctx, resourceGroup, jobName, name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Stream Input EventHub %q (Stream Analytics Job %q / Resource Group %q) %+v", name, jobName, resourceGroup, err)
		}
	}

	return nil
}
