// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_data_flow

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/dataflows"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryDataFlowRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.DataFlowClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataflows.ParseDataflowID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, dataflows.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		mappingDataFlow, ok := model.Properties.(dataflows.MappingDataFlow)
		if !ok {
			return fmt.Errorf("classifying type of %s: Expected: %q", id, helpers.DataFlowTypeMappingDataFlow)
		}

		d.Set("name", id.DataflowName)
		d.Set("data_factory_id", factories.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName).ID())
		d.Set("description", mappingDataFlow.Description)

		if err := d.Set("annotations", helpers.FlattenDataFactoryAnnotations(mappingDataFlow.Annotations)); err != nil {
			return fmt.Errorf("setting `annotations`: %+v", err)
		}

		folder := ""
		if mappingDataFlow.Folder != nil {
			folder = pointer.From(mappingDataFlow.Folder.Name)
		}
		d.Set("folder", folder)

		if props := mappingDataFlow.TypeProperties; props != nil {
			d.Set("script", props.Script)
			d.Set("script_lines", props.ScriptLines)

			if err := d.Set("source", helpers.FlattenDataFactoryDataFlowSource(props.Sources)); err != nil {
				return fmt.Errorf("setting `source`: %+v", err)
			}
			if err := d.Set("sink", helpers.FlattenDataFactoryDataFlowSink(props.Sinks)); err != nil {
				return fmt.Errorf("setting `sink`: %+v", err)
			}
			if err := d.Set("transformation", helpers.FlattenDataFactoryDataFlowTransformation(props.Transformations)); err != nil {
				return fmt.Errorf("setting `transformation`: %+v", err)
			}
		}
	}

	return nil
}
