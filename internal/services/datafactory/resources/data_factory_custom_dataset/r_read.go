// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_custom_dataset

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryCustomDatasetRead(d *pluginsdk.ResourceData, meta any) error {
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
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", factories.NewFactoryID(subscriptionId, id.ResourceGroup, id.FactoryName).ID())

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

	annotations := make([]any, 0)
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
	if err := d.Set("parameters", helpers.FlattenDataSetParameters(parameters)); err != nil {
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
	additionalProperties := make(map[string]any)
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

func flattenDataFactoryLinkedService(input *datafactory.LinkedServiceReference) []any {
	if input == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"name":       pointer.From(input.ReferenceName),
			"parameters": input.Parameters,
		},
	}
}
