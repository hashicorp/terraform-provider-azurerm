package datafactory

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datafactory/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataFactoryCustomDataset() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryCustomDatasetCreateUpdate,
		Read:   resourceDataFactoryCustomDatasetRead,
		Update: resourceDataFactoryCustomDatasetCreateUpdate,
		Delete: resourceDataFactoryCustomDatasetDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DataSetID(id)
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
				ValidateFunc: validate.LinkedServiceDatasetName,
			},

			"data_factory_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryID,
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

func resourceDataFactoryCustomDatasetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := parse.DataFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewDataSetID(subscriptionId, dataFactoryId.ResourceGroup, dataFactoryId.FactoryName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_custom_dataset", id.ID())
		}
	}

	props := map[string]interface{}{
		"type":              d.Get("type").(string),
		"linkedServiceName": expandDataFactoryLinkedService(d.Get("linked_service").([]interface{})),
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
		props["folder"] = &datafactory.DatasetFolder{
			Name: utils.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		props["parameters"] = expandDataFactoryParameters(v.(map[string]interface{}))
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

	dataset := &datafactory.DatasetResource{}
	if err := dataset.UnmarshalJSON(jsonData); err != nil {
		return err
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, *dataset, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryCustomDatasetRead(d, meta)
}

func resourceDataFactoryCustomDatasetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", parse.NewDataFactoryID(subscriptionId, id.ResourceGroup, id.FactoryName).ID())

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
		datasetFolder := &datafactory.DatasetFolder{}
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

	parameters := make(map[string]*datafactory.ParameterSpecification)
	if v, ok := m["parameters"]; ok && v != nil {
		if err := json.Unmarshal(*v, &parameters); err != nil {
			return err
		}
		delete(m, "parameters")
	}
	if err := d.Set("parameters", flattenDataFactoryParameters(parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	var linkedService *datafactory.LinkedServiceReference
	if v, ok := m["linkedServiceName"]; ok && v != nil {
		linkedService = &datafactory.LinkedServiceReference{}
		if err := json.Unmarshal(*v, linkedService); err != nil {
			return err
		}
		delete(m, "linkedServiceName")
	}
	if err := d.Set("linked_service", flattenDataFactoryLinkedService(linkedService)); err != nil {
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

func resourceDataFactoryCustomDatasetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataSetID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandDataFactoryLinkedService(input []interface{}) *datafactory.LinkedServiceReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &datafactory.LinkedServiceReference{
		ReferenceName: utils.String(v["name"].(string)),
		Type:          utils.String("LinkedServiceReference"),
		Parameters:    v["parameters"].(map[string]interface{}),
	}
}

func flattenDataFactoryLinkedService(input *datafactory.LinkedServiceReference) []interface{} {
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
