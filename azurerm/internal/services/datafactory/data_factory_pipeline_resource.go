package datafactory

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryPipeline() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryPipelineCreateUpdate,
		Read:   resourceDataFactoryPipelineRead,
		Update: resourceDataFactoryPipelineCreateUpdate,
		Delete: resourceDataFactoryPipelineDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
				ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
			},

			"data_factory_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryName(),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/5788
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"variables": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"activities_json": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				StateFunc:        utils.NormalizeJson,
				DiffSuppressFunc: suppressJsonOrderingDifference,
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"folder": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceDataFactoryPipelineCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.PipelinesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Data Factory Pipeline creation.")

	resourceGroupName := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	dataFactoryName := d.Get("data_factory_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, dataFactoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Data Factory Pipeline %q (Resource Group %q / Data Factory %q): %s", name, resourceGroupName, dataFactoryName, err)
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

	if v, ok := d.GetOk("activities_json"); ok {
		activities, err := deserializeDataFactoryPipelineActivities(v.(string))
		if err != nil {
			return fmt.Errorf("parsing 'activities_json' for Data Factory Pipeline %q (Resource Group %q / Data Factory %q) ID: %+v", name, resourceGroupName, dataFactoryName, err)
		}
		pipeline.Activities = activities
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		pipeline.Annotations = &annotations
	} else {
		annotations := make([]interface{}, 0)
		pipeline.Annotations = &annotations
	}

	if v, ok := d.GetOk("folder"); ok {
		name := v.(string)
		pipeline.Folder = &datafactory.PipelineFolder{
			Name: &name,
		}
	}

	config := datafactory.PipelineResource{
		Pipeline: pipeline,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroupName, dataFactoryName, name, config, ""); err != nil {
		return fmt.Errorf("creating Data Factory Pipeline %q (Resource Group %q / Data Factory %q): %+v", name, resourceGroupName, dataFactoryName, err)
	}

	read, err := client.Get(ctx, resourceGroupName, dataFactoryName, name, "")
	if err != nil {
		return fmt.Errorf("retrieving Data Factory Pipeline %q (Resource Group %q / Data Factory %q): %+v", name, resourceGroupName, dataFactoryName, err)
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read Data Factory Pipeline %q (Resource Group %q / Data Factory %q) ID", name, resourceGroupName, dataFactoryName)
	}

	d.SetId(*read.ID)

	return resourceDataFactoryPipelineRead(d, meta)
}

func resourceDataFactoryPipelineRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("reading the state of Data Factory Pipeline %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("data_factory_name", dataFactoryName)

	if props := resp.Pipeline; props != nil {
		d.Set("description", props.Description)

		parameters := flattenDataFactoryParameters(props.Parameters)
		if err := d.Set("parameters", parameters); err != nil {
			return fmt.Errorf("setting `parameters`: %+v", err)
		}

		annotations := flattenDataFactoryAnnotations(props.Annotations)
		if err := d.Set("annotations", annotations); err != nil {
			return fmt.Errorf("setting `annotations`: %+v", err)
		}

		if folder := props.Folder; folder != nil {
			if folder.Name != nil {
				d.Set("folder", folder.Name)
			}
		}

		variables := flattenDataFactoryVariables(props.Variables)
		if err := d.Set("variables", variables); err != nil {
			return fmt.Errorf("setting `variables`: %+v", err)
		}

		if activities := props.Activities; activities != nil {
			activitiesJson, err := serializeDataFactoryPipelineActivities(activities)
			if err != nil {
				return fmt.Errorf("serializing `activities_json`: %+v", err)
			}
			if err := d.Set("activities_json", activitiesJson); err != nil {
				return fmt.Errorf("setting `activities_json`: %+v", err)
			}
		}
	}

	return nil
}

func resourceDataFactoryPipelineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("deleting Data Factory Pipeline %q (Resource Group %q / Data Factory %q): %+v", name, resourceGroupName, dataFactoryName, err)
	}

	return nil
}
