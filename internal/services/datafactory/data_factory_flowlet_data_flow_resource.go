// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

func resourceDataFactoryFlowletDataFlow() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryFlowletDataFlowCreateUpdate,
		Read:   resourceDataFactoryFlowletDataFlowRead,
		Update: resourceDataFactoryFlowletDataFlowCreateUpdate,
		Delete: resourceDataFactoryFlowletDataFlowDelete,

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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewDataFlowID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_data_factory_flowlet_data_flow", id.ID())
		}
	}

	flowLet := datafactory.Flowlet{
		FlowletTypeProperties: &datafactory.FlowletTypeProperties{
			Script:          utils.String(d.Get("script").(string)),
			Sinks:           expandDataFactoryDataFlowSink(d.Get("sink").([]interface{})),
			Sources:         expandDataFactoryDataFlowSource(d.Get("source").([]interface{})),
			Transformations: expandDataFactoryDataFlowTransformation(d.Get("transformation").([]interface{})),
		},
		Description: utils.String(d.Get("description").(string)),
		Type:        datafactory.TypeBasicDataFlowTypeFlowlet,
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotations := v.([]interface{})
		flowLet.Annotations = &annotations
	}

	if v, ok := d.GetOk("folder"); ok {
		flowLet.Folder = &datafactory.DataFlowFolder{
			Name: utils.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("script_lines"); ok {
		flowLet.ScriptLines = utils.ExpandStringSlice(v.([]interface{}))
	}

	dataFlow := datafactory.DataFlowResource{
		Properties: &flowLet,
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.Name, dataFlow, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryFlowletDataFlowRead(d, meta)
}

func resourceDataFactoryFlowletDataFlowRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DataFlowClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataFlowID(d.Id())
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

	flowLet, ok := resp.Properties.AsFlowlet()
	if !ok {
		return fmt.Errorf("classifying type of %s: Expected: %q", id, datafactory.TypeBasicDataFlowTypeFlowlet)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", factories.NewFactoryID(id.SubscriptionId, id.ResourceGroup, id.FactoryName).ID())
	d.Set("description", flowLet.Description)

	if err := d.Set("annotations", flattenDataFactoryAnnotations(flowLet.Annotations)); err != nil {
		return fmt.Errorf("setting `annotations`: %+v", err)
	}

	folder := ""
	if flowLet.Folder != nil && flowLet.Folder.Name != nil {
		folder = *flowLet.Folder.Name
	}
	d.Set("folder", folder)

	if prop := flowLet.FlowletTypeProperties; prop != nil {
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

	return nil
}

func resourceDataFactoryFlowletDataFlowDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.DataFlowClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataFlowID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
