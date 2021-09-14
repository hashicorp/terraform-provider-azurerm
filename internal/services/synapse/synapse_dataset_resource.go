package synapse

import (
	"encoding/json"
	"fmt"
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

func resourceSynapseDataset() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseDatasetCreateUpdate,
		Read:   resourceSynapseDatasetRead,
		Update: resourceSynapseDatasetCreateUpdate,
		Delete: resourceSynapseDatasetDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DatasetID(id)
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

			"linked_service": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"parameters": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type_properties_json": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				StateFunc:        utils.NormalizeJson,
				DiffSuppressFunc: suppressJsonOrderingDifference,
			},

			"additional_properties": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
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

			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"schema_json": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				StateFunc:        utils.NormalizeJson,
				DiffSuppressFunc: suppressJsonOrderingDifference,
			},
		},
	}
}

func resourceSynapseDatasetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	client, err := synapseClient.DatasetClient(workspaceId.Name, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	id := parse.NewDatasetID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.GetDataset(ctx, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_dataset", id.ID())
		}
	}

	props := map[string]interface{}{
		"type":              d.Get("type").(string),
		"linkedServiceName": expandSynapseLinkedService(d.Get("linked_service").([]interface{})),
	}

	typePropertiesJson := fmt.Sprintf(`{ "typeProperties": %s }`, d.Get("type_properties_json").(string))
	if err = json.Unmarshal([]byte(typePropertiesJson), &props); err != nil {
		return err
	}

	additionalProperties := d.Get("additional_properties").(map[string]interface{})
	for k, v := range additionalProperties {
		props[k] = v
	}

	if v, ok := d.GetOk("annotations"); ok {
		props["annotations"] = v.([]interface{})
	}

	if v, ok := d.GetOk("description"); ok {
		props["description"] = v.(string)
	}

	if v, ok := d.GetOk("folder"); ok {
		props["folder"] = &artifacts.DatasetFolder{
			Name: utils.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		props["parameters"] = expandSynapseParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("schema_json"); ok {
		schemaJson := fmt.Sprintf(`{ "schema": %s }`, v.(string))
		if err = json.Unmarshal([]byte(schemaJson), &props); err != nil {
			return err
		}
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"properties": props,
	})
	if err != nil {
		return err
	}

	dataset := &artifacts.DatasetResource{}
	if err := dataset.UnmarshalJSON(jsonData); err != nil {
		return err
	}

	future, err := client.CreateOrUpdateDataset(ctx, id.Name, *dataset, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSynapseDatasetRead(d, meta)
}

func resourceSynapseDatasetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.DatasetID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.DatasetClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	resp, err := client.GetDataset(ctx, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("synapse_workspace_id", parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())

	byteArr, err := json.Marshal(resp.Properties)
	if err != nil {
		return err
	}

	var m map[string]*json.RawMessage
	if err = json.Unmarshal(byteArr, &m); err != nil {
		return err
	}

	description := ""
	if v, ok := m["description"]; ok && v != nil {
		if err := json.Unmarshal(*v, &description); err != nil {
			return err
		}
		delete(m, "description")
	}
	d.Set("description", description)

	t := ""
	if v, ok := m["type"]; ok && v != nil {
		if err := json.Unmarshal(*v, &t); err != nil {
			return err
		}
		delete(m, "type")
	}
	d.Set("type", t)

	folder := ""
	if v, ok := m["folder"]; ok && v != nil {
		datasetFolder := &artifacts.DatasetFolder{}
		if err := json.Unmarshal(*v, datasetFolder); err != nil {
			return err
		}
		if datasetFolder.Name != nil {
			folder = *datasetFolder.Name
		}
		delete(m, "folder")
	}
	d.Set("folder", folder)

	annotations := make([]interface{}, 0)
	if v, ok := m["annotations"]; ok && v != nil {
		if err := json.Unmarshal(*v, &annotations); err != nil {
			return err
		}
		delete(m, "annotations")
	}
	d.Set("annotations", annotations)

	parameters := make(map[string]*artifacts.ParameterSpecification)
	if v, ok := m["parameters"]; ok && v != nil {
		if err := json.Unmarshal(*v, &parameters); err != nil {
			return err
		}
		delete(m, "parameters")
	}
	if err := d.Set("parameters", flattenSynapseParameters(parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	var linkedService *artifacts.LinkedServiceReference
	if v, ok := m["linkedServiceName"]; ok && v != nil {
		linkedService = &artifacts.LinkedServiceReference{}
		if err := json.Unmarshal(*v, linkedService); err != nil {
			return err
		}
		delete(m, "linkedServiceName")
	}
	if err := d.Set("linked_service", flattenSynapseLinkedService(linkedService)); err != nil {
		return fmt.Errorf("setting `linked_service`: %+v", err)
	}

	// set "schema"
	schemaJson := ""
	if v, ok := m["schema"]; ok {
		schemaBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		schemaJson = string(schemaBytes)
		delete(m, "schema")
	}
	d.Set("schema_json", schemaJson)

	// set "type_properties_json"
	typePropertiesJson := ""
	if v, ok := m["typeProperties"]; ok {
		typePropertiesBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		typePropertiesJson = string(typePropertiesBytes)
		delete(m, "typeProperties")
	}
	d.Set("type_properties_json", typePropertiesJson)

	delete(m, "structure")

	// set "additional_properties"
	additionalProperties := make(map[string]interface{})
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &additionalProperties); err != nil {
		return err
	}
	d.Set("additional_properties", additionalProperties)

	return nil
}

func resourceSynapseDatasetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.DatasetID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.DatasetClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	future, err := client.DeleteDataset(ctx, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}
	return nil
}

func expandSynapseLinkedService(input []interface{}) *artifacts.LinkedServiceReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &artifacts.LinkedServiceReference{
		ReferenceName: utils.String(v["name"].(string)),
		Type:          utils.String("LinkedServiceReference"),
		Parameters:    v["parameters"].(map[string]interface{}),
	}
}

func flattenSynapseLinkedService(input *artifacts.LinkedServiceReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	name := ""
	if input.ReferenceName != nil {
		name = *input.ReferenceName
	}

	return []interface{}{
		map[string]interface{}{
			"name":       name,
			"parameters": input.Parameters,
		},
	}
}
