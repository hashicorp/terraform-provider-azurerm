// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_custom_dataset

import (
	"encoding/json"
	"fmt"
	"maps"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func resourceDataFactoryCustomDatasetCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.DatasetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewDataSetID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_data_factory_custom_dataset", id.ID())
			}
		}
	}

	props := map[string]any{
		"type":              d.Get("type").(string),
		"linkedServiceName": expandDataFactoryLinkedService(d.Get("linked_service").([]any)),
	}

	typePropertiesJson := fmt.Sprintf(`{ "typeProperties": %s }`, d.Get("type_properties_json").(string))
	if err = json.Unmarshal([]byte(typePropertiesJson), &props); err != nil {
		return err
	}

	additionalProperties := d.Get("additional_properties").(map[string]any)
	maps.Copy(props, additionalProperties)

	if v, ok := d.GetOk("annotations"); ok {
		props["annotations"] = v.([]any)
	}

	if v, ok := d.GetOk("description"); ok {
		props["description"] = v.(string)
	}

	if v, ok := d.GetOk("folder"); ok {
		props["folder"] = &datafactory.DatasetFolder{
			Name: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("parameters"); ok {
		props["parameters"] = helpers.ExpandDataSetParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("schema_json"); ok {
		schemaJson := fmt.Sprintf(`{ "schema": %s }`, v.(string))
		if err = json.Unmarshal([]byte(schemaJson), &props); err != nil {
			return err
		}
	}

	jsonData, err := json.Marshal(map[string]any{
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

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceDataFactoryCustomDatasetRead(d, meta)
}

func expandDataFactoryLinkedService(input []any) *datafactory.LinkedServiceReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	return &datafactory.LinkedServiceReference{
		ReferenceName: pointer.To(v["name"].(string)),
		Type:          pointer.To("LinkedServiceReference"),
		Parameters:    v["parameters"].(map[string]any),
	}
}
