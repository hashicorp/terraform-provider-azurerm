package synapse

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/sdk/2021-06-01-preview/artifacts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSynapseDataFlow() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseDataFlowCreateUpdate,
		Read:   resourceSynapseDataFlowRead,
		Update: resourceSynapseDataFlowCreateUpdate,
		Delete: resourceSynapseDataFlowDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DataFlowID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"script": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"source": SchemaForDataFlowSourceAndSink(),

			"sink": SchemaForDataFlowSourceAndSink(),

			"transformation": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"description": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"folder": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceSynapseDataFlowCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	client, err := synapseClient.DataFlowClient(workspaceId.Name, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	id := parse.NewDataFlowID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.GetDataFlow(ctx, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_data_flow", id.ID())
		}
	}

	mappingDataFlow := artifacts.MappingDataFlow{
		MappingDataFlowTypeProperties: &artifacts.MappingDataFlowTypeProperties{
			Script:          utils.String(d.Get("script").(string)),
			Sinks:           expandSynapseDataFlowSink(d.Get("sink").([]interface{})),
			Sources:         expandSynapseDataFlowSource(d.Get("source").([]interface{})),
			Transformations: expandSynapseDataFlowTransformation(d.Get("transformation").([]interface{})),
		},
		Description: utils.String(d.Get("description").(string)),
		Type:        artifacts.TypeMappingDataFlow,
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		mappingDataFlow.Annotations = &annotations
	}

	if v, ok := d.GetOk("folder"); ok {
		mappingDataFlow.Folder = &artifacts.DataFlowFolder{
			Name: utils.String(v.(string)),
		}
	}

	dataFlow := artifacts.DataFlowResource{
		Properties: &mappingDataFlow,
	}

	future, err := client.CreateOrUpdateDataFlow(ctx, id.Name, dataFlow, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSynapseDataFlowRead(d, meta)
}

func resourceSynapseDataFlowRead(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.DataFlowID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.DataFlowClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}
	resp, err := client.GetDataFlow(ctx, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	mappingDataFlow, ok := resp.Properties.AsMappingDataFlow()
	if !ok {
		return fmt.Errorf("classifying type of %s: Expected: %q", id, artifacts.TypeMappingDataFlow)
	}

	d.Set("name", id.Name)
	d.Set("synapse_workspace_id", parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())
	d.Set("description", mappingDataFlow.Description)

	if err := d.Set("annotations", flattenSynapseAnnotations(mappingDataFlow.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	folder := ""
	if mappingDataFlow.Folder != nil && mappingDataFlow.Folder.Name != nil {
		folder = *mappingDataFlow.Folder.Name
	}
	d.Set("folder", folder)

	if prop := mappingDataFlow.MappingDataFlowTypeProperties; prop != nil {
		d.Set("script", prop.Script)

		if err := d.Set("source", flattenSynapseDataFlowSource(prop.Sources)); err != nil {
			return fmt.Errorf("setting `source`: %+v", err)
		}
		if err := d.Set("sink", flattenSynapseDataFlowSink(prop.Sinks)); err != nil {
			return fmt.Errorf("setting `sink`: %+v", err)
		}
		if err := d.Set("transformation", flattenSynapseDataFlowTransformation(prop.Transformations)); err != nil {
			return fmt.Errorf("setting `transformation`: %+v", err)
		}
	}

	return nil
}

func resourceSynapseDataFlowDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.DataFlowID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.DataFlowClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	future, err := client.DeleteDataFlow(ctx, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func expandSynapseDataFlowTransformation(input []interface{}) *[]artifacts.Transformation {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	result := make([]artifacts.Transformation, 0)
	for _, v := range input {
		raw := v.(map[string]interface{})
		result = append(result, artifacts.Transformation{
			Description: utils.String(raw["description"].(string)),
			Name:        utils.String(raw["name"].(string)),
		})
	}
	return &result
}

func flattenSynapseDataFlowTransformation(input *[]artifacts.Transformation) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, v := range *input {
		name := ""
		description := ""
		if v.Name != nil {
			name = *v.Name
		}
		if v.Description != nil {
			description = *v.Description
		}
		result = append(result, map[string]interface{}{
			"name":        name,
			"description": description,
		})
	}
	return result
}

func flattenSynapseAnnotations(input *[]interface{}) []string {
	annotations := make([]string, 0)
	if input == nil {
		return annotations
	}

	for _, annotation := range *input {
		val, ok := annotation.(string)
		if !ok {
			log.Printf("[DEBUG] Skipping annotation %q since it's not a string", val)
		}
		annotations = append(annotations, val)
	}
	return annotations
}
