package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmStreamAnalyticsOutputServiceBusQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmStreamAnalyticsOutputServiceBusQueueCreateUpdate,
		Read:   resourceArmStreamAnalyticsOutputServiceBusQueueRead,
		Update: resourceArmStreamAnalyticsOutputServiceBusQueueCreateUpdate,
		Delete: resourceArmStreamAnalyticsOutputServiceBusQueueDelete,
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

			"queue_name": {
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

			"serialization": azure.SchemaStreamAnalyticsOutputSerialization(),
		},
	}
}

func resourceArmStreamAnalyticsOutputServiceBusQueueCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).streamanalytics.OutputsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure Stream Analytics Output ServiceBus Queue creation.")
	name := d.Get("name").(string)
	jobName := d.Get("stream_analytics_job_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Stream Analytics Output ServiceBus Queue %q (Job %q / Resource Group %q): %s", name, jobName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_stream_analytics_output_servicebus_queue", *existing.ID)
		}
	}

	queueName := d.Get("queue_name").(string)
	serviceBusNamespace := d.Get("servicebus_namespace").(string)
	sharedAccessPolicyKey := d.Get("shared_access_policy_key").(string)
	sharedAccessPolicyName := d.Get("shared_access_policy_name").(string)

	serializationRaw := d.Get("serialization").([]interface{})
	serialization, err := azure.ExpandStreamAnalyticsOutputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("Error expanding `serialization`: %+v", err)
	}

	props := streamanalytics.Output{
		Name: utils.String(name),
		OutputProperties: &streamanalytics.OutputProperties{
			Datasource: &streamanalytics.ServiceBusQueueOutputDataSource{
				Type: streamanalytics.TypeMicrosoftServiceBusQueue,
				ServiceBusQueueOutputDataSourceProperties: &streamanalytics.ServiceBusQueueOutputDataSourceProperties{
					QueueName:              utils.String(queueName),
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
			return fmt.Errorf("Error Creating Stream Analytics Output ServiceBus Queue %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}

		read, err := client.Get(ctx, resourceGroup, jobName, name)
		if err != nil {
			return fmt.Errorf("Error retrieving Stream Analytics Output ServiceBus Queue %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}
		if read.ID == nil {
			return fmt.Errorf("Cannot read ID of Stream Analytics Output ServiceBus Queue %q (Job %q / Resource Group %q)", name, jobName, resourceGroup)
		}

		d.SetId(*read.ID)
	} else {
		if _, err := client.Update(ctx, props, resourceGroup, jobName, name, ""); err != nil {
			return fmt.Errorf("Error Updating Stream Analytics Output ServiceBus Queue %q (Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
		}
	}

	return resourceArmStreamAnalyticsOutputServiceBusQueueRead(d, meta)
}

func resourceArmStreamAnalyticsOutputServiceBusQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).streamanalytics.OutputsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	jobName := id.Path["streamingjobs"]
	name := id.Path["outputs"]

	resp, err := client.Get(ctx, resourceGroup, jobName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Output ServiceBus Queue %q was not found in Stream Analytics Job %q / Resource Group %q - removing from state!", name, jobName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Stream Output EventHub %q (Stream Analytics Job %q / Resource Group %q): %+v", name, jobName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("stream_analytics_job_name", jobName)

	if props := resp.OutputProperties; props != nil {
		v, ok := props.Datasource.AsServiceBusQueueOutputDataSource()
		if !ok {
			return fmt.Errorf("Error converting Output Data Source to a ServiceBus Queue Output: %+v", err)
		}

		d.Set("queue_name", v.QueueName)
		d.Set("servicebus_namespace", v.ServiceBusNamespace)
		d.Set("shared_access_policy_name", v.SharedAccessPolicyName)

		if err := d.Set("serialization", azure.FlattenStreamAnalyticsOutputSerialization(props.Serialization)); err != nil {
			return fmt.Errorf("Error setting `serialization`: %+v", err)
		}
	}

	return nil
}

func resourceArmStreamAnalyticsOutputServiceBusQueueDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).streamanalytics.OutputsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	jobName := id.Path["streamingjobs"]
	name := id.Path["outputs"]

	if resp, err := client.Delete(ctx, resourceGroup, jobName, name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting Output ServiceBus Queue %q (Stream Analytics Job %q / Resource Group %q) %+v", name, jobName, resourceGroup, err)
		}
	}

	return nil
}
