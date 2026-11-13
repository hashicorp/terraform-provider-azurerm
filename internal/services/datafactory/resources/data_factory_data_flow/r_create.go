// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_data_flow

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/dataflows"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryDataFlowCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.DataFlowClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := dataflows.NewDataflowID(dataFactoryId.SubscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, dataflows.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_data_factory_data_flow", id.ID())
			}
		}
	}

	mappingDataFlow := dataflows.MappingDataFlow{
		TypeProperties: &dataflows.MappingDataFlowTypeProperties{
			Script:          pointer.To(d.Get("script").(string)),
			Sinks:           helpers.ExpandDataFactoryDataFlowSink(d.Get("sink").([]any)),
			Sources:         helpers.ExpandDataFactoryDataFlowSource(d.Get("source").([]any)),
			Transformations: helpers.ExpandDataFactoryDataFlowTransformation(d.Get("transformation").([]any)),
		},
		Description: pointer.To(d.Get("description").(string)),
		Type:        helpers.DataFlowTypeMappingDataFlow,
	}

	if v, ok := d.GetOk("annotations"); ok {
		mappingDataFlow.Annotations = pointer.To(v.([]any))
	}

	if v, ok := d.GetOk("folder"); ok {
		mappingDataFlow.Folder = &dataflows.DataFlowFolder{
			Name: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("script_lines"); ok {
		mappingDataFlow.TypeProperties.ScriptLines = pluginsdk.ExpandStringSlice(v.([]any))
	}

	dataFlow := dataflows.DataFlowResource{
		Properties: &mappingDataFlow,
	}

	if _, err := client.CreateOrUpdate(ctx, id, dataFlow, dataflows.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf(" creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceDataFactoryDataFlowRead(d, meta)
}
