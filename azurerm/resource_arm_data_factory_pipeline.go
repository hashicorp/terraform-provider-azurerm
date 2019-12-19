package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataFactoryPipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataFactoryPipelineCreateUpdate,
		Read:   resourceArmDataFactoryPipelineRead,
		Update: resourceArmDataFactoryPipelineCreateUpdate,
		Delete: resourceArmDataFactoryPipelineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
			},

			"data_factory_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryName(),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/5788
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"variables": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"annotations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmDataFactoryPipelineCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.PipelinesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Data Factory Pipeline creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, dataFactoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Data Factory Pipeline %q (Resource Group %q / Data Factory %q): %s", name, resourceGroupName, dataFactoryName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_data_factory_pipeline", *existing.ID)
		}
	}

	description := d.Get("description").(string)
	pipeline := &datafactory.Pipeline{
		Parameters:  expandDataFactoryParameters(d.Get("parameters").(map[string]interface{})),
		Variables:   expandDataFactoryVariables(d.Get("variables").(map[string]interface{})),
		Description: &description,
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		pipeline.Annotations = &annotations
	} else {
		annotations := make([]interface{}, 0)
		pipeline.Annotations = &annotations
	}

	config := datafactory.PipelineResource{
		Pipeline: pipeline,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroupName, dataFactoryName, name, config, ""); err != nil {
		return fmt.Errorf("Error creating Data Factory Pipeline %q (Resource Group %q / Data Factory %q): %+v", name, resourceGroupName, dataFactoryName, err)
	}

	read, err := client.Get(ctx, resourceGroupName, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Data Factory Pipeline %q (Resource Group %q / Data Factory %q): %+v", name, resourceGroupName, dataFactoryName, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Data Factory Pipeline %q (Resource Group %q / Data Factory %q) ID", name, resourceGroupName, dataFactoryName)
	}

	d.SetId(*read.ID)

	return resourceArmDataFactoryPipelineRead(d, meta)
}

func resourceArmDataFactoryPipelineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.PipelinesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	dataFactoryName := id.Path["factories"]
	name := id.Path["pipelines"]

	resp, err := client.Get(ctx, id.ResourceGroup, dataFactoryName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			log.Printf("[DEBUG] Data Factory Pipeline %q was not found in Resource Group %q - removing from state!", name, id.ResourceGroup)
			return nil
		}
		return fmt.Errorf("Error reading the state of Data Factory Pipeline %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	if props := resp.Pipeline; props != nil {
		d.Set("description", props.Description)

		parameters := flattenDataFactoryParameters(props.Parameters)
		if err := d.Set("parameters", parameters); err != nil {
			return fmt.Errorf("Error setting `parameters`: %+v", err)
		}

		annotations := flattenDataFactoryAnnotations(props.Annotations)
		if err := d.Set("annotations", annotations); err != nil {
			return fmt.Errorf("Error setting `annotations`: %+v", err)
		}

		variables := flattenDataFactoryVariables(props.Variables)
		if err := d.Set("variables", variables); err != nil {
			return fmt.Errorf("Error setting `variables`: %+v", err)
		}
	}

	return nil
}

func resourceArmDataFactoryPipelineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.PipelinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	dataFactoryName := id.Path["factories"]
	name := id.Path["pipelines"]
	resourceGroupName := id.ResourceGroup

	if _, err = client.Delete(ctx, resourceGroupName, dataFactoryName, name); err != nil {
		return fmt.Errorf("Error deleting Data Factory Pipeline %q (Resource Group %q / Data Factory %q): %+v", name, resourceGroupName, dataFactoryName, err)
	}

	return nil
}
