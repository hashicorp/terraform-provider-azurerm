package azurerm

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRMDataFactoryPipelineName,
			},

			"data_factory_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"parameters": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			/*
				"activity": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{},
					},
				},*/
		},
	}
}

func resourceArmDataFactoryPipelineCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactoryPipelineClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Data Factory Pipeline creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
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

	parameters := expandDataFactoryPipelineParameters(d.Get("parameters").(map[string]interface{}))

	pipeline := &datafactory.Pipeline{
		Parameters: parameters,
	}

	config := datafactory.PipelineResource{
		Pipeline: pipeline,
	}

	_, err := client.CreateOrUpdate(ctx, resourceGroupName, dataFactoryName, name, config, "")
	if err != nil {
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
	client := meta.(*ArmClient).dataFactoryPipelineClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	dataFactoryName := id.Path["factories"]
	name := id.Path["pipelines"]
	resourceGroupName := id.ResourceGroup

	resp, err := client.Get(ctx, resourceGroupName, dataFactoryName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			log.Printf("[DEBUG] Data Factory Pipeline %q was not found in Resource Group %q - removing from state!", name, resourceGroupName)
			return nil
		}
		return fmt.Errorf("Error reading the state of Data Factory Pipeline %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroupName)
	d.Set("data_factory_name", dataFactoryName)

	if props := resp.Pipeline; props != nil {
		parameters := flattenDataFactoryPipelineParameters(props.Parameters)
		if err := d.Set("parameters", parameters); err != nil {
			return fmt.Errorf("Error setting `parameters`: %+v", err)
		}
	}

	return nil
}

func resourceArmDataFactoryPipelineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).dataFactoryPipelineClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	dataFactoryName := id.Path["factories"]
	name := id.Path["pipelines"]
	resourceGroupName := id.ResourceGroup

	_, err = client.Delete(ctx, resourceGroupName, dataFactoryName, name)
	if err != nil {
		return fmt.Errorf("Error deleting Data Factory Pipeline %q (Resource Group %q / Data Factory %q): %+v", name, resourceGroupName, dataFactoryName, err)
	}

	return nil
}

func validateAzureRMDataFactoryPipelineName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	if regexp.MustCompile(`^[.+?/<>*%&:\\]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("any of '.', '+', '?', '/', '<', '>', '*', '%%', '&', ':', '\\', are not allowed in %q: %q", k, value))
	}

	return warnings, errors
}

func expandDataFactoryPipelineParameters(input map[string]interface{}) map[string]*datafactory.ParameterSpecification {
	output := make(map[string]*datafactory.ParameterSpecification)

	for k, v := range input {
		output[k] = &datafactory.ParameterSpecification{
			Type:         datafactory.ParameterTypeString,
			DefaultValue: v.(string),
		}
	}

	return output
}

func flattenDataFactoryPipelineParameters(input map[string]*datafactory.ParameterSpecification) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		if v != nil {
			// we only support string parameters at this time
			val, ok := v.DefaultValue.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping parameter %q since it's not a string", k)
			}

			output[k] = val
		}
	}

	return output
}
