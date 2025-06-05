// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/dataflows"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataFactoryFlowletDataFlow() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryFlowletDataFlowCreateUpdate,
		Read:   resourceDataFactoryFlowletDataFlowRead,
		Update: resourceDataFactoryFlowletDataFlowCreateUpdate,
		Delete: resourceDataFactoryFlowletDataFlowDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := dataflows.ParseDataflowID(id)
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

			"data_factory_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: factories.ValidateFactoryID,
			},

			"script": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				AtLeastOneOf: []string{"script", "script_lines"},
			},

			"script_lines": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"script", "script_lines"},
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"source": SchemaForDataFlowletSourceAndSink(),

			"sink": SchemaForDataFlowletSourceAndSink(),

			"transformation": SchemaForDataFlowSourceTransformation(),

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

func resourceDataFactoryFlowletDataFlowCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DataFlowClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := dataflows.NewDataflowID(dataFactoryId.SubscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, dataflows.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_factory_flowlet_data_flow", id.ID())
		}
	}

	flowLet := dataflows.Flowlet{
		TypeProperties: &dataflows.FlowletTypeProperties{
			Script:          pointer.To(d.Get("script").(string)),
			Sinks:           expandDataFactoryDataFlowSink(d.Get("sink").([]interface{})),
			Sources:         expandDataFactoryDataFlowSource(d.Get("source").([]interface{})),
			Transformations: expandDataFactoryDataFlowTransformation(d.Get("transformation").([]interface{})),
		},
		Description: pointer.To(d.Get("description").(string)),
		Type:        helper.DataFlowTypeFlowlet,
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		flowLet.Annotations = &annotations
	}

	if v, ok := d.GetOk("folder"); ok {
		flowLet.Folder = &dataflows.DataFlowFolder{
			Name: pointer.To(v.(string)),
		}
	}

	if v, ok := d.GetOk("script_lines"); ok {
		flowLet.TypeProperties.ScriptLines = utils.ExpandStringSlice(v.([]interface{}))
	}

	dataFlow := dataflows.DataFlowResource{
		Properties: &flowLet,
	}

	if _, err := client.CreateOrUpdate(ctx, id, dataFlow, dataflows.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryFlowletDataFlowRead(d, meta)
}

func resourceDataFactoryFlowletDataFlowRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		flowLet, ok := model.Properties.(dataflows.Flowlet)
		if !ok {
			return fmt.Errorf("classifying type of %s: Expected: %q", id, helper.DataFlowTypeFlowlet)
		}

		d.Set("name", id.DataflowName)
		d.Set("data_factory_id", factories.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName).ID())
		d.Set("description", flowLet.Description)

		if err := d.Set("annotations", flattenDataFactoryAnnotations(flowLet.Annotations)); err != nil {
			return fmt.Errorf("setting `annotations`: %+v", err)
		}

		folder := ""
		if flowLet.Folder != nil {
			folder = pointer.From(flowLet.Folder.Name)
		}
		d.Set("folder", folder)

		if prop := flowLet.TypeProperties; prop != nil {
			d.Set("script", prop.Script)
			d.Set("script_lines", prop.ScriptLines)

			if err := d.Set("source", flattenDataFactoryDataFlowSource(prop.Sources)); err != nil {
				return fmt.Errorf("setting `source`: %+v", err)
			}
			if err := d.Set("sink", flattenDataFactoryDataFlowSink(prop.Sinks)); err != nil {
				return fmt.Errorf("setting `sink`: %+v", err)
			}
			if err := d.Set("transformation", flattenDataFactoryDataFlowTransformation(prop.Transformations)); err != nil {
				return fmt.Errorf("setting `transformation`: %+v", err)
			}
		}
	}

	return nil
}

func resourceDataFactoryFlowletDataFlowDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DataFlowClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataflows.ParseDataflowID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
