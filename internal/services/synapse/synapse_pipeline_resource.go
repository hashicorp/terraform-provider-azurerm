package synapse

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/sdk/2021-06-01-preview/artifacts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSynapsePipeline() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapsePipelineCreateUpdate,
		Read:   resourceSynapsePipelineRead,
		Update: resourceSynapsePipelineCreateUpdate,
		Delete: resourceSynapsePipelineDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PipelineID(id)
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SynapsePipelineAndTriggerName(),
			},

			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

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
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"activities_json": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressJsonOrderingDifference, // bespoke since the ordering of keys matters
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceSynapsePipelineCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	client, err := synapseClient.PipelinesClient(workspaceId.Name, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	log.Printf("[INFO] preparing arguments for Synapse Pipeline create or update.")

	id := parse.NewPipelineID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.GetPipeline(ctx, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_pipeline", id.ID())
		}
	}

	pipeline := &artifacts.Pipeline{
		Parameters:  expandSynapseParameters(d.Get("parameters").(map[string]interface{})),
		Variables:   expandSynapseVariables(d.Get("variables").(map[string]interface{})),
		Description: utils.String(d.Get("description").(string)),
	}

	if v, ok := d.GetOk("activities_json"); ok {
		activities, err := deserializeSynapsePipelineActivities(v.(string))
		if err != nil {
			return fmt.Errorf("parsing 'activities_json' for %s", err)
		}
		pipeline.Activities = activities
	}

	annotations := d.Get("annotations").([]interface{})
	pipeline.Annotations = &annotations

	config := artifacts.PipelineResource{
		Pipeline: pipeline,
	}

	future, err := client.CreateOrUpdatePipeline(ctx, id.Name, config, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation/updation for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSynapsePipelineRead(d, meta)
}

func resourceSynapsePipelineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.PipelineID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.PipelinesClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	resp, err := client.GetPipeline(ctx, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.Set("synapse_workspace_id", parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())
	d.Set("name", id.Name)

	if props := resp.Pipeline; props != nil {
		d.Set("description", props.Description)

		parameters := flattenSynapseParameters(props.Parameters)
		if err := d.Set("parameters", parameters); err != nil {
			return fmt.Errorf("setting `parameters`: %+v", err)
		}

		annotations := flattenSynapseAnnotations(props.Annotations)
		if err := d.Set("annotations", annotations); err != nil {
			return fmt.Errorf("setting `annotations`: %+v", err)
		}

		variables := flattenSynapseVariables(props.Variables)
		if err := d.Set("variables", variables); err != nil {
			return fmt.Errorf("setting `variables`: %+v", err)
		}

		activitiesJson, err := serializeSynapsePipelineActivities(props.Activities)
		if err != nil {
			return fmt.Errorf("serializing `activities_json`: %+v", err)
		}
		if err := d.Set("activities_json", activitiesJson); err != nil {
			return fmt.Errorf("setting `activities_json`: %+v", err)
		}
	}

	return nil
}

func resourceSynapsePipelineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.PipelineID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.PipelinesClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	if _, err = client.DeletePipeline(ctx, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandSynapseParameters(input map[string]interface{}) map[string]*artifacts.ParameterSpecification {
	output := make(map[string]*artifacts.ParameterSpecification)

	for k, v := range input {
		output[k] = &artifacts.ParameterSpecification{
			Type:         artifacts.ParameterTypeString,
			DefaultValue: v.(string),
		}
	}

	return output
}

func flattenSynapseParameters(input map[string]*artifacts.ParameterSpecification) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		if v != nil {
			// we only support string parameters at this moment
			val, ok := v.DefaultValue.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping parameter %q since it's not a string", k)
			}

			output[k] = val
		}
	}

	return output
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

func expandSynapseVariables(input map[string]interface{}) map[string]*artifacts.VariableSpecification {
	output := make(map[string]*artifacts.VariableSpecification)

	for k, v := range input {
		output[k] = &artifacts.VariableSpecification{
			Type:         artifacts.VariableTypeString,
			DefaultValue: v.(string),
		}
	}

	return output
}

func flattenSynapseVariables(input map[string]*artifacts.VariableSpecification) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		if v != nil {
			// we only support string parameters at this moment
			val, ok := v.DefaultValue.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping variable %q since it's not a string", k)
			}

			output[k] = val
		}
	}

	return output
}

func deserializeSynapsePipelineActivities(jsonData string) (*[]artifacts.BasicActivity, error) {
	jsonData = fmt.Sprintf(`{ "activities": %s }`, jsonData)
	pipeline := &artifacts.Pipeline{}
	err := pipeline.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		return nil, err
	}
	return pipeline.Activities, nil
}

func serializeSynapsePipelineActivities(activities *[]artifacts.BasicActivity) (string, error) {
	pipeline := &artifacts.Pipeline{Activities: activities}
	result, err := pipeline.MarshalJSON()
	if err != nil {
		return "nil", err
	}

	var m map[string]*json.RawMessage
	err = json.Unmarshal(result, &m)
	if err != nil {
		return "", err
	}

	activitiesJson, err := json.Marshal(m["activities"])
	if err != nil {
		return "", err
	}

	return string(activitiesJson), nil
}

func suppressJsonOrderingDifference(_, old, new string, _ *pluginsdk.ResourceData) bool {
	return utils.NormalizeJson(old) == utils.NormalizeJson(new)
}
