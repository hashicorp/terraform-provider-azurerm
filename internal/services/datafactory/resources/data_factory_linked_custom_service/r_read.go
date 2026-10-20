// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_custom_service

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

func resourceDataFactoryLinkedCustomServiceRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.LinkedServiceClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServiceID(d.Id())
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
	if err := d.Set("parameters", helpers.FlattenLinkedServiceParameters(parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	var integrationRuntime *datafactory.IntegrationRuntimeReference
	if v, ok := m["connectVia"]; ok && v != nil {
		integrationRuntime = &datafactory.IntegrationRuntimeReference{}
		if err := json.Unmarshal(*v, &integrationRuntime); err != nil {
			return err
		}
		delete(m, "connectVia")
	}
	if err := d.Set("integration_runtime", flattenDataFactoryLinkedServiceIntegrationRuntimeV2(integrationRuntime)); err != nil {
		return fmt.Errorf("setting `integration_runtime`: %+v", err)
	}

	delete(m, "typeProperties")

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

func flattenDataFactoryLinkedServiceIntegrationRuntimeV2(input *datafactory.IntegrationRuntimeReference) []any {
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
